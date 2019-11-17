package db

import (
	"database/sql"
	"github.com/dare-rider/carpark/constant"
	"github.com/dare-rider/carpark/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"time"

	"github.com/dare-rider/carpark/config"
	"github.com/jmoiron/sqlx"
)

var dbConn *sqlx.DB

// InitMysqlDb Initializes the DB connection
func InitMysqlDb(config *config.DbConfig) {
	var err error
	dbConn, err = sqlx.Connect("mysql", config.Dsn)
	utils.HandleError(err)
	dbConn.SetMaxIdleConns(constant.DBMaxIdleConns())
	dbConn.SetMaxOpenConns(constant.DBMaxOpenConns())
	dbConn.SetConnMaxLifetime(constant.DBConnTimeout * time.Second)
	var errPing = dbConn.Ping()
	utils.HandleError(errPing)
}

func InitMigrations(config *config.DbConfig) {
	db, err := sql.Open("mysql", config.Dsn)
	utils.HandleError(err)
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	utils.HandleError(err)
	m, err := migrate.NewWithDatabaseInstance(
		config.MigrationPath,
		"mysql",
		driver,
	)
	utils.HandleError(err)
	m.Steps(2)
}

// MysqlConn returns existing mysql *DB Conn
func MysqlConn() *sqlx.DB {
	return dbConn
}
