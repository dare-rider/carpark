package constant

import (
	"github.com/dare-rider/carpark/config"
)

// Environment independent Application constants
const (
	DBConnTimeout      = 60 // In seconds
	DBResultLimit      = 50
	DefaultHttpTimeout = 10 // In seconds
)

// Environment dependent variables

// Don't export these variables rather create getters
// else they won't behave as constants
var (
	dbMaxIdleConns = 1
	dbMaxOpenConns = 5
)

// InitConstants initialises Environment specific constants
// Can also update above variables based on environment
func InitConstants(env *config.MiscConfig) {
	// update the constant values based on environments
	if env.Production() {
		dbMaxIdleConns = 5
		dbMaxOpenConns = 400
	} else if env.Staging() {

	}
}

// Getters for Environment dependent variables`
// DBMaxIdleConns -> dbMaxIdleConns
func DBMaxIdleConns() int {
	return dbMaxIdleConns
}

// DBMaxOpenConns -> dbMaxOpenConns
func DBMaxOpenConns() int {
	return dbMaxOpenConns
}
