package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	admin_cinema_service "github.com/Falokut/admin_cinema_service/pkg/admin_cinema_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CinemaRepository interface {
	GetCinemasInCity(ctx context.Context, id int32) (*admin_cinema_service.Cinemas, error)
	CreateCinema(ctx context.Context, in *admin_cinema_service.CreateCinemaRequest) (int32, error)
	DeleteCinema(ctx context.Context, id int32) error

	CreateCity(ctx context.Context, name string) (int32, error)
	UpdateCity(ctx context.Context, id int32, name string) error
	DeleteCity(ctx context.Context, id int32) error
	GetCinemasCities(ctx context.Context) (*admin_cinema_service.Cities, error)

	GetMoviesScreenings(ctx context.Context, cinemaID int32,
		startPeriod, endPeriod time.Time) (*admin_cinema_service.PreviewScreenings, error)
	GetScreenings(ctx context.Context, cinemaID, movieID int32,
		startPeriod, endPeriod time.Time) (*admin_cinema_service.Screenings, error)
	CreateScreenings(ctx context.Context, in *admin_cinema_service.CreateScreeningsRequest) ([]int32, error)

	CreateHall(ctx context.Context,
		in *admin_cinema_service.CreateHallRequest) (*admin_cinema_service.CreateHallResponse, error)
	DeleteHall(ctx context.Context, id int32) error
	UpdateHall(ctx context.Context, in *admin_cinema_service.UpdateHallRequest) error
	GetHalls(ctx context.Context, ids []int32) (*admin_cinema_service.Halls, error)
	IsHallExist(ctx context.Context, id int32) (bool, error)
	GetHallConfiguraion(ctx context.Context, id int32) (*admin_cinema_service.HallConfiguration, error)

	UpdateHallType(ctx context.Context, id int32, name string) error
	CreateHallType(ctx context.Context, name string) (int32, error)
	DeleteHallType(ctx context.Context, id int32) error

	UpdateScreeningType(ctx context.Context, id int32, name string) error
	CreateScreeningType(ctx context.Context, name string) (int32, error)
	DeleteScreeningType(ctx context.Context, id int32) error
}

type cinemaService struct {
	admin_cinema_service.UnimplementedCinemaServiceV1Server
	logger       *logrus.Logger
	cinemaRepo   CinemaRepository
	errorHandler errorHandler
}

func NewCinemaService(logger *logrus.Logger, cinemaRepo CinemaRepository) *cinemaService {
	errorHandler := newErrorHandler(logger)
	return &cinemaService{logger: logger,
		errorHandler: errorHandler,
		cinemaRepo:   cinemaRepo,
	}
}

func (s *cinemaService) GetCinemasInCity(ctx context.Context,
	in *admin_cinema_service.GetCinemasInCityRequest) (*admin_cinema_service.Cinemas, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.GetCinemasInCity")
	defer span.Finish()

	res, err := s.cinemaRepo.GetCinemasInCity(ctx, in.CityId)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}

	span.SetTag("grpc.status", codes.OK)
	return res, nil
}

func (s *cinemaService) GetMoviesScreenings(ctx context.Context,
	in *admin_cinema_service.GetMoviesScreeningsRequest) (*admin_cinema_service.PreviewScreenings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.GetMoviesScreenings")
	defer span.Finish()

	start, end, err := parsePeriods(in.StartPeriod, in.EndPeriod)
	if err != nil {
		return nil, s.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, err.Error())
	}
	res, err := s.cinemaRepo.GetMoviesScreenings(ctx, in.CinemaId, start, end)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}

	span.SetTag("grpc.status", codes.OK)
	return res, nil
}

func (s *cinemaService) GetScreenings(ctx context.Context,
	in *admin_cinema_service.GetScreeningsRequest) (*admin_cinema_service.Screenings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.GetScreenings")
	defer span.Finish()
	start, end, err := parsePeriods(in.StartPeriod, in.EndPeriod)
	if err != nil {
		return nil, s.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, err.Error())
	}

	res, err := s.cinemaRepo.GetScreenings(ctx, in.CinemaId, in.MovieId, start, end)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}

	span.SetTag("grpc.status", codes.OK)
	return res, nil
}

func (s *cinemaService) GetCinemasCities(ctx context.Context,
	in *emptypb.Empty) (*admin_cinema_service.Cities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.GetCinemasCities")
	defer span.Finish()

	res, err := s.cinemaRepo.GetCinemasCities(ctx)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return res, nil
}

func (s *cinemaService) GetHallConfiguration(ctx context.Context,
	in *admin_cinema_service.GetHallConfigurationRequest) (*admin_cinema_service.HallConfiguration, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.GetHallConfiguration")
	defer span.Finish()
	res, err := s.cinemaRepo.GetHallConfiguraion(ctx, in.HallId)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return res, nil
}

func (s *cinemaService) CreateScreenings(ctx context.Context,
	in *admin_cinema_service.CreateScreeningsRequest) (*admin_cinema_service.CreateScreeningsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.CreateScreenings")
	defer span.Finish()
	if len(in.Screenings) == 0 {
		return nil, s.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, "screenings mustn't be empty")
	}

	ids, err := s.cinemaRepo.CreateScreenings(ctx, in)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}

	span.SetTag("grpc.status", codes.OK)
	return &admin_cinema_service.CreateScreeningsResponse{ScreeningsIds: ids}, nil
}

func (s *cinemaService) CreateHall(ctx context.Context,
	in *admin_cinema_service.CreateHallRequest) (*admin_cinema_service.CreateHallResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.CreateHall")
	defer span.Finish()

	if in.Configuration == nil || len(in.Configuration.Place) == 0 {
		return nil, s.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, "hall configuration must't be empty")
	}
	res, err := s.cinemaRepo.CreateHall(ctx, in)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return res, nil
}

func (s *cinemaService) CreateCinema(ctx context.Context,
	in *admin_cinema_service.CreateCinemaRequest) (*admin_cinema_service.CreateCinemaResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.CreateCinema")
	defer span.Finish()

	if in.Coordinates == nil {
		return nil, s.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, "cinema must have coordinates")
	}
	id, err := s.cinemaRepo.CreateCinema(ctx, in)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &admin_cinema_service.CreateCinemaResponse{CinemaId: id}, nil
}

func (s *cinemaService) DeleteCinema(ctx context.Context,
	in *admin_cinema_service.DeleteCinemaRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.DeleteCinema")
	defer span.Finish()

	err := s.cinemaRepo.DeleteCinema(ctx, in.CinemaId)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *cinemaService) DeleteCity(ctx context.Context,
	in *admin_cinema_service.DeleteCityRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.DeleteCity")
	defer span.Finish()

	err := s.cinemaRepo.DeleteCity(ctx, in.CityId)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *cinemaService) DeleteHall(ctx context.Context,
	in *admin_cinema_service.DeleteHallRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.DeleteHall")
	defer span.Finish()

	err := s.cinemaRepo.DeleteHall(ctx, in.HallId)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *cinemaService) GetHalls(ctx context.Context,
	in *admin_cinema_service.GetHallsRequest) (*admin_cinema_service.Halls, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.GetHalls")
	defer span.Finish()

	in.HallsIds = strings.ReplaceAll(in.HallsIds, `"`, "")
	if err := checkIds(in.HallsIds); err != nil {
		return nil, s.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, err.Error())
	}

	ids := convertStringsSlice(strings.Split(in.HallsIds, ","))
	if len(ids) == 0 {
		return nil, s.errorHandler.createErrorResponseWithSpan(span, ErrInvalidArgument, "halls_ids musn't be empty")
	}

	res, err := s.cinemaRepo.GetHalls(ctx, ids)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return res, nil
}

func (s *cinemaService) CreateCity(ctx context.Context,
	in *admin_cinema_service.CreateCityRequest) (*admin_cinema_service.CreateCityResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.CreateCity")
	defer span.Finish()

	id, err := s.cinemaRepo.CreateCity(ctx, in.Name)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &admin_cinema_service.CreateCityResponse{CityId: id}, nil
}

func (s *cinemaService) UpdateCity(ctx context.Context,
	in *admin_cinema_service.UpdateCityRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.UpdateCity")
	defer span.Finish()

	err := s.cinemaRepo.UpdateCity(ctx, in.CityId, in.Name)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func convertStringsSlice(str []string) []int32 {
	var nums = make([]int32, 0, len(str))
	for _, s := range str {
		num, err := strconv.Atoi(s)
		if err == nil {
			nums = append(nums, int32(num))
		}
	}
	return nums
}

func (s *cinemaService) UpdateHall(ctx context.Context,
	in *admin_cinema_service.UpdateHallRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.UpdateHall")
	defer span.Finish()

	err := s.cinemaRepo.UpdateHall(ctx, in)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *cinemaService) UpdateHallType(ctx context.Context,
	in *admin_cinema_service.UpdateHallTypeRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.UpdateHallType")
	defer span.Finish()

	err := s.cinemaRepo.UpdateHallType(ctx, in.TypeId, in.Name)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *cinemaService) CreateHallType(ctx context.Context,
	in *admin_cinema_service.CreateHallTypeRequest) (*admin_cinema_service.CreateHallTypeResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.CreateHallType")
	defer span.Finish()

	id, err := s.cinemaRepo.CreateHallType(ctx, in.Name)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &admin_cinema_service.CreateHallTypeResponse{TypeId: id}, nil
}

func (s *cinemaService) DeleteHallType(ctx context.Context,
	in *admin_cinema_service.DeleteHallTypeRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.DeleteHallType")
	defer span.Finish()

	err := s.cinemaRepo.DeleteHallType(ctx, in.TypeId)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *cinemaService) UpdateScreeningType(ctx context.Context,
	in *admin_cinema_service.UpdateScreeningTypeRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.UpdateScreeningType")
	defer span.Finish()

	err := s.cinemaRepo.UpdateScreeningType(ctx, in.TypeId, in.Name)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *cinemaService) CreateScreeningType(ctx context.Context,
	in *admin_cinema_service.CreateScreeningTypeRequest) (*admin_cinema_service.CreateScreeningTypeResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.CreateScreeningType")
	defer span.Finish()

	id, err := s.cinemaRepo.CreateScreeningType(ctx, in.Name)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &admin_cinema_service.CreateScreeningTypeResponse{TypeId: id}, nil
}

func (s *cinemaService) DeleteScreeningType(ctx context.Context,
	in *admin_cinema_service.DeleteScreeningTypeRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaService.DeleteScreeningType")
	defer span.Finish()

	err := s.cinemaRepo.DeleteHallType(ctx, in.TypeId)
	if err != nil {
		ext.LogError(span, err)
		span.SetTag("grpc.status", status.Code(err))
		return nil, err
	}
	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func parsePeriods(startPeriod, endPeriod *admin_cinema_service.Timestamp) (start time.Time, end time.Time, err error) {
	if startPeriod == nil || endPeriod == nil {
		err = fmt.Errorf("invalid period value, it mustn't be empty")
		return
	}
	start, err = time.Parse(time.RFC3339, startPeriod.FormattedTimestamp)
	if err != nil {
		err = fmt.Errorf("invalid start period value, it must be RFC3339 layout value: %s", startPeriod)
		return
	}
	end, err = time.Parse(time.RFC3339, endPeriod.FormattedTimestamp)
	if err != nil {
		err = fmt.Errorf("invalid start period value, it must be RFC3339 layout value: %s", endPeriod)
		return
	}

	return
}

func checkIds(val string) error {
	exp := regexp.MustCompile("^[!-&!+,0-9]+$")
	if !exp.Match([]byte(val)) {
		return errors.New("invalid ids value, ids must contains only digits and comma separator")
	}

	return nil
}
