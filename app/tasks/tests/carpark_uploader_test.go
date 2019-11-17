package tests

import (
	"github.com/bxcodec/faker"
	"github.com/dare-rider/carpark/app/models/carpark"
	carparkMocks "github.com/dare-rider/carpark/app/models/carpark/mocks"
	"github.com/dare-rider/carpark/app/tasks"
	"github.com/dare-rider/carpark/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type CarparkUploaderTestSuite struct {
	suite.Suite
	carparkUc       *carparkMocks.Usecase
	appRootPath     string
	csvFileBasePath string
	fakeCpMod       *carpark.Model
}

func (suite *CarparkUploaderTestSuite) SetupTest() {
	suite.carparkUc = &carparkMocks.Usecase{}
	dir, _ := os.Getwd()
	suite.appRootPath = utils.JoinURL(dir, "/../../../")
	suite.csvFileBasePath = utils.JoinURL(suite.appRootPath, "db/migrations/seed")
	suite.fakeCpMod = &carpark.Model{}
	_ = faker.FakeData(suite.fakeCpMod)
}

func TestCarparkUploaderTestSuite(t *testing.T) {
	suite.Run(t, new(CarparkUploaderTestSuite))
}

func (suite *CarparkUploaderTestSuite) TestUpload() {
	t := suite.T()
	t.Run("should return an error if csv file not exists", func(t *testing.T) {
		suite.SetupTest()
		suite.carparkUc.On("InsertOrUpdate", suite.fakeCpMod).Return(nil)
		suite.csvFileBasePath = "incorrect/csv/path"
		upldr := tasks.NewCarparkUploader(suite.carparkUc, suite.csvFileBasePath)
		err := upldr.Upload()
		require.Error(t, err)
	})
}
