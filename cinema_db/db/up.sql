CREATE EXTENSION postgis;

CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE cinemas (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    city_id INT REFERENCES cities(id) ON UPDATE CASCADE ON DELETE SET NULL,
    --address without city
    address TEXT NOT NULL,
    coordinates geography(POINT,4326) NOT NULL
);

CREATE TABLE halls_types (
    type_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE halls (
    id SERIAL PRIMARY KEY,
    cinema_id INT REFERENCES cinemas(id) ON UPDATE CASCADE ON DELETE SET NULL,
    hall_type_id INT REFERENCES halls_types(type_id) ON UPDATE CASCADE ON DELETE SET NULL,
    name TEXT NOT NULL,
    hall_size INT NOT NULL DEFAULT 0
);

CREATE TABLE halls_configurations (
    hall_id INT REFERENCES halls(id) ON UPDATE CASCADE ON DELETE CASCADE,
    row INT CHECK(row > 0),
    seat INT CHECK(seat > 0),
    grid_pos_x FLOAT NOT NULL,
    grid_pos_y FLOAT NOT NULL,
    PRIMARY KEY(hall_id, row, seat)
);

CREATE OR REPLACE FUNCTION update_hall_size()
RETURNS TRIGGER
AS $$
BEGIN
    UPDATE halls SET hall_size=( SELECT COUNT(seat) FROM halls_configurations WHERE hall_id=id);
    RETURN NEW;
END; $$
LANGUAGE PLPGSQL;

CREATE TRIGGER hall_place_insert_trigger
            AFTER INSERT ON halls_configurations
            EXECUTE FUNCTION update_hall_size();

CREATE TRIGGER hall_place_delete_trigger
            AFTER DELETE ON halls_configurations
            FOR EACH ROW
            EXECUTE FUNCTION update_hall_size();


CREATE TABLE screenings_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE screenings (
    id BIGSERIAL PRIMARY KEY,
    screening_type_id INT REFERENCES screenings_types(id) ON UPDATE CASCADE ON DELETE SET NULL,
    movie_id INT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL CHECK(start_time > clock_timestamp()),
    hall_id INT REFERENCES halls(id) ON UPDATE CASCADE ON DELETE SET NULL,
    ticket_price DECIMAL(8,2) CHECK(ticket_price>0.0)
);

GRANT USAGE, SELECT ON SEQUENCE  cities_id_seq TO admin_cinema_service;
GRANT  SELECT,UPDATE,DELETE,INSERT ON cities TO admin_cinema_service;

GRANT USAGE, SELECT ON SEQUENCE  cinemas_id_seq TO admin_cinema_service;
GRANT  SELECT,UPDATE,DELETE,INSERT ON cinemas TO admin_cinema_service;
GRANT  SELECT,UPDATE,DELETE,INSERT ON halls_configurations TO admin_cinema_service;

GRANT USAGE, SELECT ON SEQUENCE  halls_types_type_id_seq TO admin_cinema_service;
GRANT  SELECT,UPDATE,DELETE,INSERT ON halls_types TO admin_cinema_service;

GRANT USAGE, SELECT ON SEQUENCE  halls_id_seq TO admin_cinema_service;
GRANT  SELECT,UPDATE,DELETE,INSERT ON halls TO admin_cinema_service;

GRANT USAGE, SELECT ON SEQUENCE  screenings_id_seq TO admin_cinema_service;
GRANT  SELECT,UPDATE,DELETE,INSERT ON screenings TO admin_cinema_service;
GRANT USAGE, SELECT ON SEQUENCE  screenings_types_id_seq TO admin_cinema_service;
GRANT  SELECT,UPDATE,DELETE,INSERT ON screenings_types TO admin_cinema_service;





