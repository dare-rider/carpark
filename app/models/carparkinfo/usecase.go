package carparkinfo

import (
	"github.com/dare-rider/carpark/app/models"
	"github.com/jmoiron/sqlx"
)

type Usecase interface {
	InsertOrUpdateByCarParkNo(mod *Model, tx ...*sqlx.Tx) error
	FindAllByCarparkIDs(cpIds []int) ([]Model, error)
}

type usecase struct {
	models.BaseUsecase
	rp repoI
}

func NewUsecase(db *sqlx.DB) Usecase {
	return &usecase{
		rp: newRepo(db),
	}
}

func (uc *usecase) InsertOrUpdateByCarParkNo(mod *Model, tx ...*sqlx.Tx) error {
	return uc.rp.insertOrUpdateByCarParkNo(mod, tx...)
}

func (uc *usecase) FindAllByCarparkIDs(cpIds []int) ([]Model, error) {
	return uc.rp.findAllByCarparkIDs(cpIds)
}
