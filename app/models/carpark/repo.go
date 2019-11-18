package carpark

import (
	"fmt"
	"github.com/dare-rider/carpark/app/models"
	"github.com/dare-rider/carpark/app/models/carparkinfo"
	"github.com/jmoiron/sqlx"
)

type Model struct {
	ID                  int     `db:"id"`
	CarParkNo           string  `db:"car_park_no"`
	Address             string  `db:"address"`
	XCoord              float64 `db:"x_coord"`
	YCoord              float64 `db:"y_coord"`
	Latitude            float64 `db:"latitude"`
	Longitude           float64 `db:"longitude"`
	CarParkType         string  `db:"car_park_type"`
	TypeOfParkingSystem string  `db:"type_of_parking_system"`
	ShortTermParking    string  `db:"short_term_parking"`
	FreeParking         string  `db:"free_parking"`
	NightParking        bool    `db:"night_parking"`
	CarParkDecks        int     `db:"car_park_decks"`
	GantryHeight        float64 `db:"gantry_height"`
	CarParkBasement     bool    `db:"car_park_basement"`
	// Non DB field
	DistanceFromCurrentLocation float64 `db:"distance_from_current_location"`
	CarparkInfos                []carparkinfo.Model
}

const (
	insertOrUpdateQuery = `
		INSERT into carparks
			(car_park_no, address, x_coord, y_coord, car_park_type, type_of_parking_system, short_term_parking,
				free_parking, night_parking, car_park_decks, gantry_height, car_park_basement, latitude, longitude)
		VALUES
			(:car_park_no, :address, :x_coord, :y_coord, :car_park_type, :type_of_parking_system, :short_term_parking,
				:free_parking, :night_parking, :car_park_decks, :gantry_height, :car_park_basement, :latitude, :longitude)
		ON DUPLICATE KEY UPDATE
			address = :address,
			car_park_type = :car_park_type,
			type_of_parking_system = :type_of_parking_system,
			short_term_parking = :short_term_parking,
			free_parking = :free_parking,
			night_parking = :night_parking,
			car_park_decks = :car_park_decks,
			gantry_height = :gantry_height,
			car_park_basement = :car_park_basement
	`
)

var (
	// param 0 - current latitude
	// param 1 - current longitude
	// param 2 - current latitude
	twoPtDistanceQuery = `
		111.111 * DEGREES(
            ACOS(LEAST(
                  1.0,
                  COS(RADIANS(?)) * COS(RADIANS(latitude)) * COS(RADIANS(? - longitude))
                  + SIN(RADIANS(?)) * SIN(RADIANS(latitude)))
                )
          )
	`
	// param 0 - limit
	// param 1 - offset
	fetchNearest = fmt.Sprintf(`
		SELECT
			%s as distance_from_current_location, id, car_park_no, address, x_coord, y_coord, car_park_type,
			type_of_parking_system, short_term_parking, free_parking, night_parking, car_park_decks, gantry_height, car_park_basement,
			latitude, longitude
		FROM
			carparks
		ORDER BY
			distance_from_current_location
		LIMIT ? OFFSET ?
	`, twoPtDistanceQuery)
)

type RepoI interface {
	InsertOrUpdate(mod *Model, tx ...*sqlx.Tx) error
	FetchNearest(currentLat float64, currentLong float64, limit int, offset int) ([]Model, error)
}

type repo struct {
	models.BaseRepo
}

func NewRepo(db *sqlx.DB) RepoI {
	rp := &repo{}
	rp.Db = db
	return rp
}

func (rp *repo) InsertOrUpdate(mod *Model, tx ...*sqlx.Tx) error {
	db := rp.DbOrTx(tx...)
	rows, err := db.NamedQuery(insertOrUpdateQuery, mod)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (rp *repo) FetchNearest(currentLat float64, currentLong float64, limit int, offset int) ([]Model, error) {
	var results []Model
	err := rp.Db.Select(&results, fetchNearest, currentLat, currentLong, currentLat, limit, offset)
	if err != nil {
		return nil, err
	}
	return results, nil
}
