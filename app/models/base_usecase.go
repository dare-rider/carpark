package models

import "github.com/dare-rider/carpark/constant"

type BaseUsecase struct{}

func (bu BaseUsecase) LimitOffset(pgNo int, perPg int) (limit int, offset int) {
	if pgNo == 0 {
		pgNo = 1
	}
	if perPg == 0 {
		perPg = constant.DBResultLimit
	}
	limit = perPg
	offset = (pgNo - 1) * perPg
	return
}
