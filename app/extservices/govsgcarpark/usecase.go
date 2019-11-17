package govsgcarpark

import (
	"github.com/dare-rider/carpark/app/models/carparkinfo"
	"strconv"
)

type usecase struct {
	rp RepoI
}

type Usecase interface {
	CarparkInfos() ([]carparkinfo.Model, error)
}

func NewUsecase(rp RepoI) Usecase {
	return &usecase{
		rp: rp,
	}
}

func (uc *usecase) CarparkInfos() ([]carparkinfo.Model, error) {
	rawResp, err := uc.rp.CarparkAvailability()
	if err != nil {
		return nil, err
	}
	var cpiMods []carparkinfo.Model
	for _, item := range rawResp.Items {
		for _, cpData := range item.CarparkData {
			cpiMods = append(cpiMods, uc.serializeToCarparkInfoMod(&cpData)...)
		}
	}
	return cpiMods, nil
}

func (uc *usecase) serializeToCarparkInfoMod(rawData *carparkData) []carparkinfo.Model {
	var cpiMods []carparkinfo.Model
	for _, cpInfo := range rawData.CarparkInfos {
		cpiMod := carparkinfo.Model{
			CarParkNo: rawData.CarparkNumber,
			LotType:   cpInfo.LotType,
		}
		cpiMod.LotsAvailable, _ = strconv.Atoi(cpInfo.LotsAvailable)
		cpiMod.TotalLots, _ = strconv.Atoi(cpInfo.TotalLots)
		cpiMods = append(cpiMods, cpiMod)
	}
	return cpiMods
}
