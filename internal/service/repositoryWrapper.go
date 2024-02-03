package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Falokut/admin_cinema_service/internal/repository"
	admin_cinema_service "github.com/Falokut/admin_cinema_service/pkg/admin_cinema_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

type GetMoviesDurationer interface {
	GetMoviesDuration(ctx context.Context, ids []int32) (map[int32]uint32, error)
}

type cinemaRepositoryWrapper struct {
	errorHandler
	logger        *logrus.Logger
	moviesService GetMoviesDurationer

	repo repository.CinemaRepository
}

func NewCinemaRepositoryWrapper(logger *logrus.Logger,
	cinemaRepository repository.CinemaRepository, moviesService GetMoviesDurationer) *cinemaRepositoryWrapper {
	errorHandler := newErrorHandler(logger)
	return &cinemaRepositoryWrapper{
		logger:        logger,
		errorHandler:  errorHandler,
		repo:          cinemaRepository,
		moviesService: moviesService,
	}
}
func (w *cinemaRepositoryWrapper) UpdateHallType(ctx context.Context, id int32, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.UpdateHallType")
	defer span.Finish()

	exist, err := w.repo.IsHallTypeAlreadyExist(ctx, name)
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if exist {
		return w.createErrorResponseWithSpan(span, ErrAlreadyExist, fmt.Sprintf("hall type with %s name already exist", name))
	}

	err = w.repo.UpdateHallType(ctx, id, name)
	if errors.Is(err, repository.ErrNotFound) {
		return w.createErrorResponseWithSpan(span, ErrNotFound, fmt.Sprintf("hall type with %d id not exist", id))
	}
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return nil
}
func (w *cinemaRepositoryWrapper) CreateHallType(ctx context.Context, name string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.CreateHallType")
	defer span.Finish()

	exist, err := w.repo.IsHallTypeAlreadyExist(ctx, name)
	if err != nil {
		return 0, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if exist {
		return 0, w.createErrorResponseWithSpan(span, ErrAlreadyExist, fmt.Sprintf("hall type with %s name already exist", name))
	}

	id, err := w.repo.CreateHallType(ctx, name)
	if err != nil {
		return 0, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return id, nil
}
func (w *cinemaRepositoryWrapper) DeleteHallType(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.DeleteHallType")
	defer span.Finish()

	err := w.repo.DeleteHallType(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return w.createErrorResponseWithSpan(span, ErrNotFound, fmt.Sprintf("hall type with %d id not exist", id))
	}
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	return nil
}

func (w *cinemaRepositoryWrapper) UpdateScreeningType(ctx context.Context, id int32, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"cinemaRepositoryWrapper.UpdateScreeningType")
	defer span.Finish()

	exist, err := w.repo.IsScreeningTypeAlreadyExist(ctx, name)
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if exist {
		return w.createErrorResponseWithSpan(span, ErrAlreadyExist,
			fmt.Sprintf("screening type with %s name already exist", name))
	}

	err = w.repo.UpdateScreeningType(ctx, id, name)
	if errors.Is(err, repository.ErrNotFound) {
		return w.createErrorResponseWithSpan(span, ErrNotFound, fmt.Sprintf("screening type with %d id not exist", id))
	}
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return nil
}

func (w *cinemaRepositoryWrapper) CreateScreeningType(ctx context.Context, name string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.CreateHallType")
	defer span.Finish()

	exist, err := w.repo.IsScreeningTypeAlreadyExist(ctx, name)
	if err != nil {
		return 0, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if exist {
		return 0, w.createErrorResponseWithSpan(span, ErrAlreadyExist, fmt.Sprintf("screening type with %s name already exist",
			name))
	}

	id, err := w.repo.CreateScreeningType(ctx, name)
	if err != nil {
		return 0, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return id, nil
}

func (w *cinemaRepositoryWrapper) DeleteScreeningType(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"cinemaRepositoryWrapper.DeleteScreeningType")
	defer span.Finish()

	err := w.repo.DeleteScreeningType(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return w.createErrorResponseWithSpan(span, ErrNotFound, fmt.Sprintf("screening type with %d id not exist", id))
	}
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	return nil
}

func (w *cinemaRepositoryWrapper) GetCinemasInCity(ctx context.Context, id int32) (*admin_cinema_service.Cinemas, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.GetCinemasInCity")
	defer span.Finish()

	cinemas, err := w.repo.GetCinemasInCity(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(cinemas) == 0 {
		return nil, w.createErrorResponseWithSpan(span, ErrNotFound, "")
	}

	return convertToCinemas(cinemas), nil
}

func convertToCinemas(cinemas []repository.Cinema) *admin_cinema_service.Cinemas {
	res := &admin_cinema_service.Cinemas{}
	res.Cinemas = make([]*admin_cinema_service.Cinema, len(cinemas))
	for i, cinema := range cinemas {
		res.Cinemas[i] = &admin_cinema_service.Cinema{
			CinemaId: cinema.ID,
			Name:     cinema.Name,
			Address:  cinema.Address,
			Coordinates: &admin_cinema_service.Coordinates{
				Latityde:  cinema.Coordinates.Latityde,
				Longitude: cinema.Coordinates.Longitude,
			},
		}
	}
	return res
}

func (w *cinemaRepositoryWrapper) GetMoviesScreenings(ctx context.Context, cinemaID int32,
	startPeriod, endPeriod time.Time) (*admin_cinema_service.PreviewScreenings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"cinemaRepositoryWrapper.GetMoviesScreenings")
	defer span.Finish()

	previews, err := w.repo.GetMoviesScreenings(ctx, cinemaID, startPeriod, endPeriod)
	if err != nil {
		return nil, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if len(previews) == 0 {
		return nil, w.createErrorResponseWithSpan(span, ErrNotFound, "")
	}

	res := &admin_cinema_service.PreviewScreenings{}
	res.Screenings = make([]*admin_cinema_service.PreviewScreening, len(previews))
	for i, preview := range previews {
		res.Screenings[i] = &admin_cinema_service.PreviewScreening{
			MovieID:         preview.MovieID,
			ScreeningsTypes: preview.ScreeningsTypes,
			HallsTypes:      preview.HallsTypes,
		}
	}

	return res, nil
}

func (w *cinemaRepositoryWrapper) GetScreenings(ctx context.Context, cinemaID, movieID int32,
	startPeriod, endPeriod time.Time) (*admin_cinema_service.Screenings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"cinemaRepositoryWrapper.GetScreenings")
	defer span.Finish()

	screenings, err := w.repo.GetScreenings(ctx, cinemaID, movieID,
		startPeriod, endPeriod)
	if err != nil {
		return nil, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if len(screenings) == 0 {
		return nil, w.createErrorResponseWithSpan(span, ErrNotFound, "")
	}

	res := &admin_cinema_service.Screenings{}
	res.Screenings = make([]*admin_cinema_service.Screening, len(screenings))
	for i, screening := range screenings {
		res.Screenings[i] = &admin_cinema_service.Screening{
			ScreeningID:   screening.ScreeningID,
			ScreeningType: screening.ScreeningType,
			MovieID:       screening.MovieID,
			HallId:        screening.HallID,
			StartTime:     &admin_cinema_service.Timestamp{FormattedTimestamp: screening.StartTime.Format(time.RFC3339)},
			TicketPrice:   priceFromFloat(screening.TicketPrice),
		}
	}

	return res, nil
}

func (w *cinemaRepositoryWrapper) GetCinemasCities(ctx context.Context) (*admin_cinema_service.Cities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"cinemaRepositoryWrapper.GetCinemasCities")
	defer span.Finish()

	cities, err := w.repo.GetCinemasCities(ctx)
	if err != nil {
		return nil, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if len(cities) == 0 {
		return nil, w.createErrorResponseWithSpan(span, ErrNotFound, "")
	}

	return convertToCities(cities), nil
}

func convertToCities(cities []repository.City) *admin_cinema_service.Cities {
	res := &admin_cinema_service.Cities{}
	res.Cities = make([]*admin_cinema_service.City, len(cities))
	for i, city := range cities {
		res.Cities[i] = &admin_cinema_service.City{CityID: city.ID, Name: city.Name}
	}
	return res
}

func (w *cinemaRepositoryWrapper) GetHallConfiguraion(ctx context.Context,
	id int32) (*admin_cinema_service.HallConfiguration, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"cinemaRepositoryWrapper.GetHallConfiguraion")
	defer span.Finish()

	places, err := w.repo.GetHallConfiguraion(ctx, id)
	if err != nil {
		return nil, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if len(places) == 0 {
		return nil, w.createErrorResponseWithSpan(span, ErrNotFound, "")
	}

	return convertToHallConfiguration(places), nil
}

func convertToHallConfiguration(places []repository.Place) *admin_cinema_service.HallConfiguration {
	res := &admin_cinema_service.HallConfiguration{}
	res.Place = make([]*admin_cinema_service.Place, len(places))
	for i, place := range places {
		res.Place[i] = &admin_cinema_service.Place{
			Row:      place.Row,
			Seat:     place.Seat,
			GridPosX: place.GridPosX,
			GridPosY: place.GridPosY,
		}
	}
	return res
}

func (w *cinemaRepositoryWrapper) CreateScreenings(ctx context.Context, in *admin_cinema_service.CreateScreeningsRequest) ([]int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"cinemaRepositoryWrapper.CreateScreenings")
	defer span.Finish()
	exist, err := w.IsHallExist(ctx, in.HallId)
	if err != nil {
		return nil, w.errorHandler.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if !exist {
		return nil, w.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument,
			fmt.Sprintf("hall with id %d not exist", in.HallId))
	}

	ids := make(map[int32]struct{}, len(in.Screenings))
	var screeningDTOs = make([]repository.ScreeningDTO, len(in.Screenings))

	var startPeriod, endPeriod time.Time
	for i, screening := range in.Screenings {
		if screening.StartTime == nil {
			return nil, w.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, "screening must have start time")
		}
		if screening.TicketPrice == nil {
			return nil, w.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, "screening must have ticket price")
		}

		t, err := time.Parse(time.RFC3339, screening.StartTime.FormattedTimestamp)
		if err != nil {
			return nil, w.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument,
				fmt.Sprintf("invalid start period value, it must be RFC3339 layout current value: %s",
					screening.StartTime.FormattedTimestamp))
		}
		screeningDTOs[i].StartTime = t
		screeningDTOs[i].TicketPrice = floatFromPrice(screening.TicketPrice)
		screeningDTOs[i].ScreeningTypeId = screening.ScreeningTypeId
		screeningDTOs[i].MovieId = screening.MovieId
		ids[screening.MovieId] = struct{}{}
		if startPeriod.After(t) || startPeriod.IsZero() {
			startPeriod = t
		}
		if endPeriod.Before(t) {
			endPeriod = t
		}
	}
	if startPeriod.Before(time.Now()) {
		return nil, w.createErrorResponseWithSpan(span, ErrInvalidArgument,
			fmt.Sprintf("invalid screening time, time now: %s, the earliest screening time: %s,"+
				"screening time must be later than the current time",
				time.Now().Format(time.RFC3339), startPeriod.Format(time.RFC3339)))
	}

	screenings, err := w.repo.GetScreeningsInHallInfo(ctx, in.HallId, startPeriod, endPeriod)
	if err != nil {
		return nil, w.errorHandler.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	var alreadyInRepo = make([]repository.ScreeningDTO, len(screenings))
	for i, screening := range screenings {
		ids[screening.MovieID] = struct{}{}
		alreadyInRepo[i].StartTime = screening.StartTime
		alreadyInRepo[i].MovieId = screening.MovieID
	}

	moviesDuration, err := w.moviesService.GetMoviesDuration(ctx, maps.Keys(ids))
	if err != nil {
		return nil, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	var notFoundIds = make([]int32, 0, len(ids))
	for id := range ids {
		_, ok := moviesDuration[id]
		if !ok {
			notFoundIds = append(notFoundIds, id)
		}
	}
	if len(notFoundIds) != 0 {
		return nil, w.createErrorResponseWithSpan(span, ErrInvalidArgument, fmt.Sprintf("movies with ids %s doesn't exist",
			convertIntSliceTOString(notFoundIds)))
	}

	if err := isScreeningsDtoTimesValid(screeningDTOs, alreadyInRepo, moviesDuration); err != nil {
		return nil, w.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, err.Error())
	}

	return w.repo.CreateScreenings(ctx, repository.CreateScreeningDTO{
		HallId:     in.HallId,
		Screenings: screeningDTOs,
	})
}

// moviesDuration - map which key - movie id, value - movie duration in munutes
// toCheck - screenings, that need check before sending to the repository
// already in repo - screenings int repository
func isScreeningsDtoTimesValid(toCheck []repository.ScreeningDTO,
	alreadyInRepo []repository.ScreeningDTO, moviesDuration map[int32]uint32) error {
	sort.Slice(toCheck, func(i, j int) bool { return toCheck[i].StartTime.After(toCheck[j].StartTime) })

	for i := 0; i < len(toCheck); i++ {
		// Checking start screening time conflicts with screenings in the input screenings
		for j := i + 1; j < len(toCheck); j++ {
			endTime := toCheck[i].StartTime.Add(time.Minute * time.Duration(moviesDuration[toCheck[i].MovieId]))
			err := checkTimeConflict(toCheck[i].StartTime, endTime, toCheck[j].StartTime)
			if err != nil {
				return err
			}
		}

		// Checking start screening time conflicts with screenings that are already in repository.
		for j := 0; j < len(alreadyInRepo); j++ {
			endTime := alreadyInRepo[j].StartTime.Add(
				time.Minute * time.Duration(moviesDuration[alreadyInRepo[j].MovieId]))
			err := checkTimeConflict(alreadyInRepo[j].StartTime, endTime, toCheck[i].StartTime)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Return error if toCheck is in the time period from start to end(start <= toCheck <= end)
// end must be later than start, function doesn't check this.
func checkTimeConflict(start, end, toCheck time.Time) error {
	if (toCheck.After(start) && toCheck.Before(end)) || toCheck.Equal(start) || toCheck.Equal(end) {
		return fmt.Errorf(`screening time conflict: %s, other screening: start time %s end time: %s`,
			toCheck.Format(time.RFC3339),
			start.Format(time.RFC3339),
			end.Format(time.RFC3339))
	}
	return nil
}

func (w *cinemaRepositoryWrapper) IsHallExist(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"cinemaRepositoryWrapper.IsHallExist")
	defer span.Finish()

	return w.repo.IsHallExist(ctx, id)
}

func (w *cinemaRepositoryWrapper) CreateHall(ctx context.Context,
	in *admin_cinema_service.CreateHallRequest) (*admin_cinema_service.CreateHallResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.CreateHall")
	defer span.Finish()

	hallInfo := repository.HallDTO{
		CinemaId: in.CinemaId,
		TypeId:   in.TypeId,
		Name:     in.Name,
	}

	exist, err := w.repo.IsHallTypeExist(ctx, in.TypeId)
	if err != nil {
		return nil, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if !exist {
		return nil, w.createErrorResponseWithSpan(span, ErrInvalidArgument, fmt.Sprintf("hall type with id %d not exist", in.TypeId))
	}

	exist, err = w.repo.IsCinemaExist(ctx, in.CinemaId)
	if err != nil {
		return nil, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if !exist {
		return nil, w.createErrorResponseWithSpan(span, ErrInvalidArgument, fmt.Sprintf("cinema with id %d not exist", in.CinemaId))
	}
	hallInfo.Configuration = make([]repository.Place, 0, len(in.Configuration.Place))
	for _, place := range in.Configuration.Place {
		hallInfo.Configuration = append(hallInfo.Configuration, repository.Place{
			Row:      place.Row,
			Seat:     place.Seat,
			GridPosX: place.GridPosX,
			GridPosY: place.GridPosY,
		})
	}

	exist, err = w.repo.IsHallAreadyExist(ctx, in.CinemaId, in.TypeId, in.Name)
	if err != nil {
		return nil, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if exist {
		return nil, w.createErrorResponseWithSpan(span, ErrAlreadyExist, "hall already exist")
	}

	id, err := w.repo.CreateHall(ctx, hallInfo)
	if err != nil {
		return nil, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return &admin_cinema_service.CreateHallResponse{HallId: id}, nil
}

func (w *cinemaRepositoryWrapper) CreateCinema(ctx context.Context,
	in *admin_cinema_service.CreateCinemaRequest) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.CreateCinema")
	defer span.Finish()

	exist, err := w.repo.IsCinemaAlreadyExist(ctx, convertCoordinatesToGeoPoint(in.Coordinates))
	if err != nil {
		return 0, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if exist {
		return 0, w.createErrorResponseWithSpan(span, ErrAlreadyExist, "cinema with this coordinates already exist")
	}

	exist, err = w.repo.IsCityExist(ctx, in.CityId)
	if err != nil {
		return 0, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if !exist {
		return 0, w.createErrorResponseWithSpan(span, ErrInvalidArgument, fmt.Sprintf("city with id %d not found", in.CityId))
	}

	id, err := w.repo.CreateCinema(ctx, repository.CreateCinemaDTO{
		Name:        in.Name,
		Address:     in.Address,
		Coordinates: convertCoordinatesToGeoPoint(in.Coordinates),
		CityId:      in.CityId,
	})
	if err != nil {
		return 0, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return id, nil
}

func (w *cinemaRepositoryWrapper) DeleteCity(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.DeleteCity")
	defer span.Finish()

	err := w.repo.DeleteCity(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return w.createErrorResponseWithSpan(span, ErrNotFound, fmt.Sprintf("cinema with id %d not exist", id))
	}
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return nil
}

func (w *cinemaRepositoryWrapper) UpdateHall(ctx context.Context,
	in *admin_cinema_service.UpdateHallRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.UpdateHall")
	defer span.Finish()

	if in.CinemaId != nil {
		exist, err := w.repo.IsCinemaExist(ctx, in.GetCinemaId())
		if err != nil {
			return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
		}
		if !exist {
			return w.createErrorResponseWithSpan(span, ErrInvalidArgument,
				fmt.Sprintf("cinema with id %d not exist", in.CinemaId))
		}

	}

	if in.TypeId != nil {
		exist, err := w.repo.IsHallTypeExist(ctx, in.GetTypeId())
		if err != nil {
			return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
		}
		if !exist {
			return w.createErrorResponseWithSpan(span, ErrInvalidArgument,
				fmt.Sprintf("hall type with id %d not exist", in.TypeId))
		}
	}

	places, err := convertToPlaces(in.Configuration)
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInvalidArgument, err.Error())
	}

	dto := repository.UpdateHallDTO{
		CinemaId:      in.GetCinemaId(),
		TypeId:        in.GetTypeId(),
		Name:          in.GetName(),
		Configuration: places,
	}

	err = w.repo.UpdateHall(ctx, in.HallId, dto)
	if errors.Is(err, repository.ErrNotFound) {
		return w.createErrorResponseWithSpan(span, ErrNotFound, fmt.Sprintf("hall with id %d not exist", in.HallId))
	}
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return nil
}

func convertToPlaces(in *admin_cinema_service.HallConfiguration) ([]repository.Place, error) {
	if in == nil {
		return []repository.Place{}, nil
	}

	var res = make([]repository.Place, len(in.Place))
	for i, place := range in.Place {
		if place == nil {
			return []repository.Place{}, fmt.Errorf("invalid place index: %d, place can't be empty", i)
		}
		res[i] = repository.Place{
			Row:      place.Row,
			Seat:     place.Seat,
			GridPosX: place.GridPosX,
			GridPosY: place.GridPosY,
		}
	}
	return res, nil
}

func (w *cinemaRepositoryWrapper) DeleteCinema(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.DeleteCinema")
	defer span.Finish()

	err := w.repo.DeleteCinema(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return w.createErrorResponseWithSpan(span, ErrNotFound, fmt.Sprintf("cinema with id %d not exist", id))
	}
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return nil
}

func (w *cinemaRepositoryWrapper) DeleteHall(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.DeleteHall")
	defer span.Finish()

	err := w.repo.DeleteHall(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return w.createErrorResponseWithSpan(span, ErrNotFound, fmt.Sprintf("hall with id %d not exist", id))
	}
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return nil
}

func (w *cinemaRepositoryWrapper) GetHalls(ctx context.Context,
	ids []int32) (*admin_cinema_service.Halls, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.GetHalls")
	defer span.Finish()

	halls, err := w.repo.GetHalls(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(halls) == 0 {
		return nil, w.createErrorResponseWithSpan(span, ErrNotFound, "")
	}

	res := &admin_cinema_service.Halls{}
	res.Halls = make([]*admin_cinema_service.Hall, len(halls))
	for i, hall := range halls {
		res.Halls[i] = &admin_cinema_service.Hall{
			HallId:   hall.Id,
			HallSize: hall.Size,
			Name:     hall.Name,
			Type:     hall.Type,
		}
	}

	return res, nil
}

func priceFromFloat(n string) *admin_cinema_service.Price {
	nums := strings.Split(n, ".")
	if len(nums) != 2 {
		return nil
	}
	units, _ := strconv.ParseInt(nums[0], 10, 64)
	nanos, _ := strconv.ParseInt(nums[1], 10, 32)
	return &admin_cinema_service.Price{Value: int32(units)*100 + int32(nanos)}
}

func (w *cinemaRepositoryWrapper) IsCityExist(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.IsCityExist")
	defer span.Finish()

	exist, err := w.repo.IsCityExist(ctx, id)
	if err != nil {
		return false, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if !exist {
		return false, nil
	}

	return true, nil
}

func (w *cinemaRepositoryWrapper) IsCityAlreadyExist(ctx context.Context, name string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.IsCityExist")
	defer span.Finish()

	exist, err := w.repo.IsCityAlreadyExist(ctx, name)
	if err != nil {
		return false, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if !exist {
		return false, nil
	}

	return true, nil
}

func (w *cinemaRepositoryWrapper) CreateCity(ctx context.Context, name string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.CreateCity")
	defer span.Finish()

	exist, err := w.IsCityAlreadyExist(ctx, name)
	if err != nil {
		return 0, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if exist {
		return 0, w.createErrorResponseWithSpan(span, ErrAlreadyExist, fmt.Sprintf("city with name %s already exist", name))
	}

	id, err := w.repo.CreateCity(ctx, name)
	if err != nil {
		return 0, w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}

	return id, nil
}

func (w *cinemaRepositoryWrapper) UpdateCity(ctx context.Context, id int32, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepositoryWrapper.UpdateCity")
	defer span.Finish()

	exist, err := w.IsCityAlreadyExist(ctx, name)
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	if exist {
		return w.createErrorResponseWithSpan(span, ErrAlreadyExist, fmt.Sprintf("city with name %s already exist", name))
	}

	err = w.repo.UpdateCity(ctx, id, name)
	if errors.Is(err, repository.ErrNotFound) {
		return w.createErrorResponseWithSpan(span, ErrNotFound, fmt.Sprintf("city with id %d not found", id))
	}
	if err != nil {
		return w.createErrorResponseWithSpan(span, ErrInternal, err.Error())
	}
	return nil
}

func floatFromPrice(in *admin_cinema_service.Price) float64 {
	if in == nil {
		return float64(0)
	}
	return float64(in.Value) / 100.0
}

func convertCoordinatesToGeoPoint(coord *admin_cinema_service.Coordinates) repository.GeoPoint {
	return repository.GeoPoint{
		Latityde:  coord.Latityde,
		Longitude: coord.Longitude,
	}
}
