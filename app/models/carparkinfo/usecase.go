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
	rp RepoI
}

func NewUsecase(rp RepoI) Usecase {
	return &usecase{
		rp: rp,
	}
}

func (uc *usecase) InsertOrUpdateByCarParkNo(mod *Model, tx ...*sqlx.Tx) error {
	return uc.rp.InsertOrUpdateByCarParkNo(mod, tx...)
}

func (uc *usecase) FindAllByCarparkIDs(cpIds []int) ([]Model, error) {
	if len(cpIds) == 0 {
		return nil, nil
	}
	return uc.rp.FindAllByCarparkIDs(cpIds)
}
