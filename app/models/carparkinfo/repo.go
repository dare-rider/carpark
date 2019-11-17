package carparkinfo

import (
	"fmt"
	"github.com/dare-rider/carpark/app/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Model struct {
	CarparkID     int    `db:"carpark_id"`
	LotType       string `db:"lot_type"`
	LotsAvailable int    `db:"lots_available"`
	TotalLots     int    `db:"total_lots"`
	// Non Db Fields
	CarParkNo string `db:"car_park_no"`
}

const (
	insertOrUpdateByCarParkNoQuery = `
		INSERT into carpark_infos
			(carpark_id, lot_type, lots_available, total_lots)
		SELECT
			cp.id as carpark_id, :lot_type, :lots_available, :total_lots
		FROM
			carparks cp WHERE cp.car_park_no = :car_park_no
		ON DUPLICATE KEY UPDATE
			lots_available = :lots_available,
			total_lots = :total_lots
	`

	selectAllByCarparkIdsQuery = `
		SELECT
			carpark_id, lot_type, lots_available, total_lots
		FROM
			carpark_infos where carpark_id in (%s)
	`
)

type RepoI interface {
	InsertOrUpdateByCarParkNo(mod *Model, tx ...*sqlx.Tx) error
	FindAllByCarparkIDs(cpIds []int) ([]Model, error)
}

type repo struct {
	models.BaseRepo
}

func NewRepo(db *sqlx.DB) RepoI {
	rp := &repo{}
	rp.Db = db
	return rp
}

func (rp *repo) InsertOrUpdateByCarParkNo(mod *Model, tx ...*sqlx.Tx) error {
	db := rp.DbOrTx(tx...)
	rows, err := db.NamedQuery(insertOrUpdateByCarParkNoQuery, mod)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (rp *repo) FindAllByCarparkIDs(cpIds []int) ([]Model, error) {
	preparedBlanks := strings.TrimRight(strings.Repeat("?,", len(cpIds)), ",")
	query := fmt.Sprintf(selectAllByCarparkIdsQuery, preparedBlanks)
	var res []Model
	// converting []int to []interface for `sqlx.Select`
	var data []interface{}
	for _, cpId := range cpIds {
		data = append(data, cpId)
	}
	err := rp.Db.Select(&res, query, data...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
