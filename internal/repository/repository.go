package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("entity not found")
var ErrInvalidArgument = errors.New("invalid input data")

func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", conStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type Cinema struct {
	ID          int32    `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Address     string   `json:"address" db:"address"`
	Coordinates GeoPoint `json:"coordinates" db:"coordinates"`
}

type City struct {
	ID   int32  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type MovieScreening struct {
	MovieID         int32    `json:"movie_id" db:"movie_id"`
	ScreeningsTypes []string `json:"screenings_types" db:"screenings_types"`
	HallsTypes      []string `json:"halls_types" db:"halls_types"`
}

type Screening struct {
	ScreeningID   int64     `db:"id"`
	MovieID       int32     `db:"movie_id"`
	ScreeningType string    `db:"screening_type"`
	HallID        int32     `db:"hall_id"`
	TicketPrice   string    `db:"ticket_price"`
	StartTime     time.Time `db:"start_time"`
}

type Place struct {
	Row      int32   `db:"row"`
	Seat     int32   `db:"seat"`
	GridPosX float32 `db:"grid_pos_x"`
	GridPosY float32 `db:"grid_pos_y"`
}

type ScreeningDTO struct {
	MovieId         int32
	StartTime       time.Time
	ScreeningTypeId int32
	TicketPrice     float64
}
type CreateScreeningDTO struct {
	HallId     int32
	Screenings []ScreeningDTO
}

type HallDTO struct {
	CinemaId      int32
	TypeId        int32
	Name          string
	Configuration []Place
}

type CreateCinemaDTO struct {
	Name        string
	Address     string
	CityId      int32
	Coordinates GeoPoint
}

type ScreeningInHallInfo struct {
	MovieID   int32     `db:"movie_id"`
	StartTime time.Time `db:"start_time"`
}

type UpdateHallDTO struct {
	CinemaId      int32  `db:"cinema_id"`
	TypeId        int32  `db:"hall_type_id"`
	Name          string `db:"name"`
	Configuration []Place
}

type Hall struct {
	Id   int32  `db:"id"`
	Type string `db:"hall_type"`
	Name string `db:"name"`
	Size uint32 `db:"size"`
}

type CinemaRepository interface {
	GetCinemasInCity(ctx context.Context, id int32) ([]Cinema, error)
	IsCinemaAlreadyExist(ctx context.Context, coordinates GeoPoint) (bool, error)
	CreateCinema(ctx context.Context, dto CreateCinemaDTO) (int32, error)
	DeleteCinema(ctx context.Context, id int32) error
	IsCinemaExist(ctx context.Context, id int32) (bool, error)

	CreateCity(ctx context.Context, name string) (int32, error)
	DeleteCity(ctx context.Context, id int32) error
	IsCityExist(ctx context.Context, id int32) (bool, error)
	GetCinemasCities(ctx context.Context) ([]City, error)
	IsCityAlreadyExist(ctx context.Context, name string) (bool, error)
	UpdateCity(ctx context.Context, id int32, name string) error

	GetMoviesScreenings(ctx context.Context, cinemaID int32, startPeriod, endPeriod time.Time) ([]MovieScreening, error)
	GetScreenings(ctx context.Context, cinemaID, movieID int32, startPeriod, endPeriod time.Time) ([]Screening, error)
	GetHallConfiguraion(ctx context.Context, id int32) ([]Place, error)
	IsHallExist(ctx context.Context, id int32) (bool, error)
	DeleteHall(ctx context.Context, id int32) error

	CreateScreenings(ctx context.Context, dto CreateScreeningDTO) ([]int32, error)
	GetScreeningsInHallInfo(ctx context.Context, hallId int32, startPeriod, endPeriod time.Time) ([]ScreeningInHallInfo, error)
	IsHallAreadyExist(ctx context.Context, cinemaId, hallTypeId int32, name string) (bool, error)
	CreateHall(ctx context.Context, dto HallDTO) (int32, error)

	UpdateHall(ctx context.Context, id int32, dto UpdateHallDTO) error
	GetHalls(ctx context.Context, ids []int32) ([]Hall, error)

	UpdateHallType(ctx context.Context, id int32, name string) error
	CreateHallType(ctx context.Context, name string) (int32, error)
	DeleteHallType(ctx context.Context, id int32) error
	IsHallTypeExist(ctx context.Context, id int32) (bool, error)
	IsHallTypeAlreadyExist(ctx context.Context, name string) (bool, error)

	UpdateScreeningType(ctx context.Context, id int32, name string) error
	CreateScreeningType(ctx context.Context, name string) (int32, error)
	DeleteScreeningType(ctx context.Context, id int32) error
	IsScreeningTypeExist(ctx context.Context, id int32) (bool, error)
	IsScreeningTypeAlreadyExist(ctx context.Context, name string) (bool, error)
}
