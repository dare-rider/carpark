package config

import (
	"github.com/dare-rider/carpark/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type DbConfig struct {
	Dsn           string `mapstructure:"dsn"`
	MigrationPath string `mapstructure:"migration_path"`
	SeedPath      string `mapstructure:"seed_path"`
}

type MiscConfig struct {
	Environment string `mapstructure:"environment"`
}

// Production returns true if the environment is production
func (misc *MiscConfig) Production() bool {
	return misc.Environment == "production"
}

// Staging returns true if the environment is stage
func (misc *MiscConfig) Staging() bool {
	return misc.Environment == "staging"
}

// Development returns true if the environment is development
func (misc *MiscConfig) Development() bool {
	return misc.Environment == "development"
}

type GovSgService struct {
	BaseUrl string `mapstructure:"base_url"`
}

type GeneralConfig struct {
	DbConfig     *DbConfig
	MiscConfig   *MiscConfig
	GovSgService *GovSgService
}

func LoadConfig(filePath string) *GeneralConfig {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	utils.HandleError(err)
	config := new(GeneralConfig)
	config.DbConfig = new(DbConfig)
	config.MiscConfig = new(MiscConfig)
	config.GovSgService = new(GovSgService)
	_ = mapstructure.Decode(viper.GetStringMap("db"), config.DbConfig)
	_ = mapstructure.Decode(viper.GetStringMap("misc"), config.MiscConfig)
	_ = mapstructure.Decode(viper.GetStringMap("gov_sg_service"), config.GovSgService)
	return config
}
