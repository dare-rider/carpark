package models

import (
	"github.com/dare-rider/carpark/types/sqlxwrap"
	"github.com/jmoiron/sqlx"
)

type BaseRepo struct {
	Db *sqlx.DB
}

type BaseRepoI interface {
	DbOrTx(tx ...*sqlx.Tx) sqlxwrap.DBOrTx
}

func (rp *BaseRepo) DbOrTx(tx ...*sqlx.Tx) sqlxwrap.DBOrTx {
	if len(tx) > 0 {
		return tx[0]
	}
	return rp.Db
}
