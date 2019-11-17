package carpark

import (
	"github.com/dare-rider/carpark/app/models"
	"github.com/dare-rider/carpark/app/models/carparkinfo"
	"github.com/dare-rider/carpark/app/requests"
	"github.com/dare-rider/carpark/constant"
	"github.com/dare-rider/carpark/utils/geodist"
	"github.com/jmoiron/sqlx"
)

type Usecase interface {
	InsertOrUpdate(mod *Model, tx ...*sqlx.Tx) error
	FetchNearestWithInfo(req *requests.NearestCarparksRequest) ([]Model, error)
	LimitOffset(pgNo int, perPg int) (limit int, offset int)
}

type usecase struct {
	models.BaseUsecase
	rp            RepoI
	carparkInfoUc carparkinfo.Usecase
}

func NewUsecase(rp RepoI, carparkInfoUc carparkinfo.Usecase) Usecase {
	return &usecase{
		rp:            rp,
		carparkInfoUc: carparkInfoUc,
	}
}

func (uc *usecase) InsertOrUpdate(mod *Model, tx ...*sqlx.Tx) error {
	return uc.rp.InsertOrUpdate(mod, tx...)
}

func (uc *usecase) FetchNearestWithInfo(req *requests.NearestCarparksRequest) ([]Model, error) {
	currentDistFromCenter := geodist.Distance(req.Latitude, req.Longitude, constant.GeoDistSgLat, constant.GeoDistSgLong, constant.GeoDistUnit)
	limit, offset := uc.LimitOffset(req.Page, req.PerPage)
	cps, err := uc.rp.FetchNearest(currentDistFromCenter, limit, offset)
	if err != nil {
		return nil, err
	}
	// Fetching Infos without n+1 query
	var cpIds []int
	for _, cp := range cps {
		cpIds = append(cpIds, cp.ID)
	}
	cpInfos, err := uc.carparkInfoUc.FindAllByCarparkIDs(cpIds)
	if err != nil {
		return nil, err
	}
	cpIdcpInfoMap := make(map[int][]carparkinfo.Model)
	for _, cpInfo := range cpInfos {
		cpIdcpInfoMap[cpInfo.CarparkID] = append(cpIdcpInfoMap[cpInfo.CarparkID], cpInfo)
	}
	var res []Model
	// mapping infos to cp
	for _, cp := range cps {
		cp.CarparkInfos = cpIdcpInfoMap[cp.ID]
		res = append(res, cp)
	}
	return res, nil
}
