package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type cinemaRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

func NewCinemaRepository(logger *logrus.Logger, db *sqlx.DB) *cinemaRepository {
	return &cinemaRepository{
		logger: logger,
		db:     db,
	}
}

const (
	cinemasTableName             = "cinemas"
	citiesTableName              = "cities"
	screeningTypeTableName       = "screenings_types"
	hallsTypesTableName          = "halls_types"
	hallsTableName               = "halls"
	screeningsTableName          = "screenings"
	hallsConfigurationsTableName = "halls_configurations"
)

func (r *cinemaRepository) GetCinemasInCity(ctx context.Context, id int32) ([]Cinema, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.GetCinemasInCity")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`SELECT id,name,address, ST_AsText(coordinates) AS coordinates
								FROM %s
								WHERE city_id=$1
								ORDER BY id`,
		cinemasTableName)

	var cinemas []Cinema
	err = r.db.SelectContext(ctx, &cinemas, query, id)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return []Cinema{}, err
	}

	return cinemas, nil
}

func (r *cinemaRepository) GetCinemasCities(ctx context.Context) ([]City, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.GetCinemasCities")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=ANY(SELECT DISTINCT city_id FROM %s) ORDER BY id",
		citiesTableName, cinemasTableName)

	var cities []City
	err = r.db.SelectContext(ctx, &cities, query)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return []City{}, err
	}

	return cities, nil
}

type previewScreening struct {
	MovieID         int32  `json:"movie_id" db:"movie_id"`
	ScreeningsTypes string `json:"screenings_types" db:"screenings_types"`
	HallsTypes      string `json:"halls_types" db:"halls_types"`
}

func (r *cinemaRepository) GetMoviesScreenings(ctx context.Context,
	cinemaID int32, startPeriod, endPeriod time.Time) ([]MovieScreening, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.GetMoviesScreenings")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`SELECT movie_id, ARRAY_AGG(DISTINCT %[1]s.name) AS screenings_types,
		ARRAY_AGG(DISTINCT %[2]s.name) AS halls_types 
		FROM %[3]s JOIN %[1]s ON screening_type_id=%[1]s.id 
		JOIN %[4]s ON hall_id=%[4]s.id JOIN %[2]s ON hall_type_id=%[2]s.type_id 
		WHERE cinema_id=$1 AND start_time>=$2 AND start_time<=$3 
		GROUP BY movie_id`,
		screeningTypeTableName, hallsTypesTableName, screeningsTableName, hallsTableName)

	var previews []previewScreening
	err = r.db.SelectContext(ctx, &previews, query, cinemaID, startPeriod, endPeriod)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return []MovieScreening{}, err
	}

	res := make([]MovieScreening, len(previews))
	for i, screening := range previews {
		res[i] = MovieScreening{
			MovieID:         screening.MovieID,
			HallsTypes:      convertSQLArray(screening.HallsTypes),
			ScreeningsTypes: convertSQLArray(screening.ScreeningsTypes),
		}
	}

	return res, nil
}

func (r *cinemaRepository) GetScreenings(ctx context.Context,
	cinemaID, movieID int32, startPeriod, endPeriod time.Time) ([]Screening, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.GetScreenings")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`
		SELECT %[1]s.id, movie_id, %[2]s.name AS screening_type, hall_id, ticket_price,start_time
		FROM %[1]s JOIN %[2]s ON screening_type_id=%[2]s.id 
		WHERE hall_id=ANY(SELECT id FROM %[3]s WHERE cinema_id=$1) AND movie_id=$2 AND start_time>=$3 AND start_time<=$4
		ORDER BY start_time`,
		screeningsTableName, screeningTypeTableName, hallsTableName)

	var screenings []Screening
	err = r.db.SelectContext(ctx, &screenings, query, cinemaID, movieID, startPeriod, endPeriod)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return []Screening{}, err
	}

	return screenings, nil
}

func (r *cinemaRepository) GetHallConfiguraion(ctx context.Context, id int32) ([]Place, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.GetHallConfiguraion")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`SELECT row, seat, grid_pos_x, grid_pos_y
								FROM %s
								WHERE hall_id=$1
								ORDER BY row,seat`,
		hallsConfigurationsTableName)
	var places []Place
	err = r.db.SelectContext(ctx, &places, query, id)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return []Place{}, err
	}

	return places, nil
}

func convertSQLArray(str string) []string {
	if strings.EqualFold(str, "{NULL}") {
		return []string{}
	}

	str = strings.Trim(str, "{}")
	return strings.Split(str, ",")
}

func (r *cinemaRepository) IsHallExist(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsHallExist")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf(`SELECT id FROM %s WHERE id=$1`, hallsTableName)
	err = r.db.GetContext(ctx, &id, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *cinemaRepository) CreateScreenings(ctx context.Context, dto CreateScreeningDTO) ([]int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsHallExist")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	var placeholders, args = make([]string, len(dto.Screenings)), make([]any, 0, len(dto.Screenings)*4+1)
	args = append(args, dto.HallId)

	for i, screening := range dto.Screenings {
		placeholders[i] = fmt.Sprintf("($1,$%d,$%d,$%d,$%d)", len(args)+1, len(args)+2, len(args)+3, len(args)+4)
		args = append(args, screening.MovieId, screening.ScreeningTypeId, screening.StartTime, screening.TicketPrice)
	}

	query := fmt.Sprintf(`INSERT INTO %s (hall_id,movie_id,screening_type_id,start_time,ticket_price)
	VALUES %s RETURNING id`, screeningsTableName, strings.Join(placeholders, ","))
	var ids = make([]int32, len(dto.Screenings))
	err = r.db.SelectContext(ctx, &ids, query, args...)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return []int32{}, err
	}

	return ids, nil
}

func (r *cinemaRepository) IsHallAreadyExist(ctx context.Context, cinemaId, hallTypeId int32, name string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsHallAreadyExist")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf(`SELECT id FROM %s WHERE cinema_id=$1 AND hall_type_id=$2 AND name=$3`, hallsTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, cinemaId, hallTypeId, name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *cinemaRepository) DeleteHall(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.DeleteHall")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("DELETE FROM %s WHERE hall_id=$1 RETURNING hall_id", hallsTableName)
	err = r.db.GetContext(ctx, &id, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return err
	}
	return nil
}

func (r *cinemaRepository) DeleteCity(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.DeleteCity")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING id", citiesTableName)
	err = r.db.GetContext(ctx, &id, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return err
	}
	return nil
}

func (r *cinemaRepository) DeleteCinema(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.DeleteCinema")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING id", cinemasTableName)
	err = r.db.GetContext(ctx, &id, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return err
	}

	return nil
}

func (r *cinemaRepository) UpdateHall(ctx context.Context, id int32, dto UpdateHallDTO) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.UpdateHall")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	placeholders, args := getUpdateStatement(&dto)
	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s %s WHERE id=$%d", hallsTableName, placeholders, len(args))
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return err
	}
	if num, _ := res.RowsAffected(); num != 1 {
		return ErrNotFound
	}

	if len(dto.Configuration) != 0 {
		query = fmt.Sprintf("DELETE FROM %s WHERE hall_id=$1", hallsConfigurationsTableName)
		_, err = tx.ExecContext(ctx, query, id)
		if err != nil {
			r.logger.Errorf("err: %v query: %s", err.Error(), query)
			return err
		}

		var placeholders, args = placesToInsertStatement(dto.Configuration, id)
		query = fmt.Sprintf("INSERT INTO %s (hall_id,row,seat,grid_pos_x,grid_pos_y) VALUES %s ON CONFLICT DO NOTHING",
			hallsConfigurationsTableName, placeholders)

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			r.logger.Errorf("err: %v query: %s", err.Error(), query)
			return err
		}
	}

	tx.Commit()
	return err
}

func (r *cinemaRepository) CreateCinema(ctx context.Context, dto CreateCinemaDTO) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.CreateCinema")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("INSERT INTO %s (name,city_id,address,coordinates) VALUES($1,$2,$3,ST_GeographyFromText($4)) RETURNING id",
		cinemasTableName)

	var id int32
	err = r.db.GetContext(ctx, &id, query, dto.Name, dto.CityId, dto.Address, dto.Coordinates)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return 0, err
	}

	return id, nil
}

func (r *cinemaRepository) IsCinemaAlreadyExist(ctx context.Context, coord GeoPoint) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsCinemaAlreadyExist")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf("SELECT id FROM %s WHERE coordinates=ST_GeographyFromText($1);",
		cinemasTableName)

	r.logger.Info(coord.Value())
	var id int32
	err = r.db.GetContext(ctx, &id, query, coord)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *cinemaRepository) GetHalls(ctx context.Context, ids []int32) ([]Hall, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.GetHalls")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`SELECT id, COALESCE(%[1]s.name,'') AS hall_type, %[2]s.name AS name, hall_size AS size
	FROM %[2]s LEFT JOIN %[1]s ON hall_type_id=type_id
	WHERE id=ANY($1)`, hallsTypesTableName, hallsTableName)
	var halls []Hall
	err = r.db.SelectContext(ctx, &halls, query, ids)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return []Hall{}, err
	}
	return halls, nil
}

func (r *cinemaRepository) CreateHall(ctx context.Context, dto HallDTO) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.CreateHall")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	query := fmt.Sprintf("INSERT INTO %s (cinema_id,hall_type_id,name) VALUES($1,$2,$3) RETURNING id;", hallsTableName)

	var id int32
	err = tx.GetContext(ctx, &id, query, dto.CinemaId, dto.TypeId, dto.Name)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return 0, err
	}

	var placeholders, args = placesToInsertStatement(dto.Configuration, id)
	query = fmt.Sprintf("INSERT INTO %s (hall_id,row,seat,grid_pos_x,grid_pos_y) VALUES %s ON CONFLICT DO NOTHING",
		hallsConfigurationsTableName, placeholders)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return 0, err
	}
	tx.Commit()

	return id, nil
}

func (r *cinemaRepository) GetScreeningsInHallInfo(ctx context.Context, hallId int32,
	startPeriod, endPeriod time.Time) ([]ScreeningInHallInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.GetScreeningsInHallInfo")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`SELECT movie_id, start_time
	FROM %[1]s WHERE hall_id=$1 AND start_time>=$2 AND start_time<=$3`, screeningsTableName)

	var screenings []ScreeningInHallInfo
	err = r.db.SelectContext(ctx, &screenings, query, hallId, startPeriod, endPeriod)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return []ScreeningInHallInfo{}, err
	}

	return screenings, nil
}

func getUpdateStatement[T comparable](val T) (string, []any) {
	rv := reflect.ValueOf(val)
	rt := rv.Type()

	statements := make([]string, 0, rt.NumField())
	args := make([]any, 0, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		if _, ok := rt.Field(i).Tag.Lookup("db"); !ok {
			continue
		}

		v := rv.Field(i).Interface()
		if isDefaultValue(v) {
			continue
		}
		args = append(args, v)
		statements = append(statements, fmt.Sprintf("%s=$%d", rt.Field(i).Tag.Get("db"), len(args)))
	}

	return " SET " + strings.Join(statements, ","), args
}

func isDefaultValue(field interface{}) bool {
	fieldVal := reflect.ValueOf(field)

	return !fieldVal.IsValid() || fieldVal.Interface() == reflect.Zero(fieldVal.Type()).Interface()
}

func placesToInsertStatement(places []Place, hallId int32) (string, []any) {
	var placeholders, args = make([]string, len(places)), make([]any, 0, len(places)*4+1)
	args = append(args, hallId)

	for i, place := range places {
		placeholders[i] = fmt.Sprintf("($1,$%d,$%d,$%d,$%d)", len(args)+1, len(args)+2, len(args)+3, len(args)+4)
		args = append(args, place.Row, place.Seat, place.GridPosX, place.GridPosY)
	}

	return strings.Join(placeholders, ","), args
}

func (r *cinemaRepository) CreateCity(ctx context.Context, name string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.CreateCity")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`INSERT INTO %s (name) VALUES($1) RETURNING id`, citiesTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return 0, err
	}

	return id, nil
}

func (r *cinemaRepository) IsCinemaExist(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsCinemaExist")
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf(`SELECT id FROM %s WHERE id=$1`, cinemasTableName)
	err = r.db.GetContext(ctx, &id, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *cinemaRepository) IsCityExist(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsCityExist")
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf(`SELECT id FROM %s WHERE id=$1`, citiesTableName)
	err = r.db.GetContext(ctx, &id, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *cinemaRepository) IsCityAlreadyExist(ctx context.Context, name string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsCityAlreadyExist")
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf(`SELECT id FROM %s WHERE name=$1`, citiesTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *cinemaRepository) UpdateCity(ctx context.Context, id int32, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.UpdateCity")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`UPDATE %s SET name=$1 WHERE id=$2`, citiesTableName)
	res, err := r.db.ExecContext(ctx, query, name, id)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return err
	}
	if num, _ := res.RowsAffected(); num != 1 {
		return ErrNotFound
	}

	return nil
}

func (r *cinemaRepository) UpdateHallType(ctx context.Context, id int32, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.UpdateHallType")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`UPDATE %s SET name=$1 WHERE type_id=$2`, hallsTypesTableName)
	res, err := r.db.ExecContext(ctx, query, name, id)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return err
	}
	if num, _ := res.RowsAffected(); num != 1 {
		return ErrNotFound
	}

	return nil
}
func (r *cinemaRepository) CreateHallType(ctx context.Context, name string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.CreateHallType")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("INSERT INTO %s (name) VALUES($1) RETURNING type_id;", hallsTypesTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return 0, err
	}
	return id, nil
}

func (r *cinemaRepository) DeleteHallType(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.DeleteHallType")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`DELETE FROM %s WHERE type_id=$1`, hallsTypesTableName)
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return err
	}
	if num, _ := res.RowsAffected(); num != 1 {
		return ErrNotFound
	}

	return nil
}
func (r *cinemaRepository) IsHallTypeExist(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsHallTypeExist")
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf(`SELECT type_id FROM %s WHERE type_id=$1`, hallsTypesTableName)
	err = r.db.GetContext(ctx, &id, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *cinemaRepository) IsHallTypeAlreadyExist(ctx context.Context, name string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsHallTypeAlreadyExist")
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf(`SELECT type_id FROM %s WHERE name=$1`, hallsTypesTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *cinemaRepository) UpdateScreeningType(ctx context.Context, id int32, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.UpdateScreeningType")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`UPDATE %s SET name=$1 WHERE id=$2`, screeningTypeTableName)
	res, err := r.db.ExecContext(ctx, query, name, id)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return err
	}
	if num, _ := res.RowsAffected(); num != 1 {
		return ErrNotFound
	}

	return nil
}

func (r *cinemaRepository) CreateScreeningType(ctx context.Context, name string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.CreateScreeningType")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("INSERT INTO %s (name) VALUES($1) RETURNING id;", screeningTypeTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return 0, err
	}
	return id, nil
}

func (r *cinemaRepository) DeleteScreeningType(ctx context.Context, id int32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.DeleteScreeningType")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, screeningTypeTableName)
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return err
	}
	if num, _ := res.RowsAffected(); num != 1 {
		return ErrNotFound
	}

	return nil
}

func (r *cinemaRepository) IsScreeningTypeExist(ctx context.Context, id int32) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cinemaRepository.IsScreeningTypeExist")
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf(`SELECT id FROM %s WHERE id=$1`, screeningTypeTableName)
	err = r.db.GetContext(ctx, &id, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}

func (r *cinemaRepository) IsScreeningTypeAlreadyExist(ctx context.Context, name string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"cinemaRepository.IsScreeningTypeAlreadyExist")
	var err error
	defer span.SetTag("error", err != nil && !errors.Is(err, sql.ErrNoRows))

	query := fmt.Sprintf(`SELECT id FROM %s WHERE name=$1`, screeningTypeTableName)
	var id int32
	err = r.db.GetContext(ctx, &id, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		r.logger.Errorf("err: %v query: %s", err.Error(), query)
		return false, err
	}

	return true, nil
}
