package tests

import (
	"errors"
	"github.com/bxcodec/faker"
	govSgCarparkMocks "github.com/dare-rider/carpark/app/extservices/govsgcarpark/mocks"
	"github.com/dare-rider/carpark/app/models/carparkinfo"
	carparkInfoMocks "github.com/dare-rider/carpark/app/models/carparkinfo/mocks"
	"github.com/dare-rider/carpark/app/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CarparkInfoUploaderTestSuite struct {
	suite.Suite
	govSgCarparkUc  *govSgCarparkMocks.Usecase
	carparkInfoUc   *carparkInfoMocks.Usecase
	fakeCpInfoModel *carparkinfo.Model
}

func (suite *CarparkInfoUploaderTestSuite) SetupTest() {
	suite.govSgCarparkUc = &govSgCarparkMocks.Usecase{}
	suite.carparkInfoUc = &carparkInfoMocks.Usecase{}
	suite.fakeCpInfoModel = new(carparkinfo.Model)
	_ = faker.FakeData(suite.fakeCpInfoModel)
}

func TestCarparkInfoUploaderTestSuite(t *testing.T) {
	suite.Run(t, new(CarparkInfoUploaderTestSuite))
}

func (suite CarparkInfoUploaderTestSuite) TestUpload() {
	t := suite.T()
	t.Run("should return error if govSgCarparkUc.CarparkInfos returns error", func(t *testing.T) {
		suite.SetupTest()
		expectedErr := errors.New("test-error")
		suite.govSgCarparkUc.On("CarparkInfos").Return(nil, expectedErr)
		suite.carparkInfoUc.On("InsertOrUpdateByCarParkNo", suite.fakeCpInfoModel).Return(nil)
		upldr := tasks.NewCarparkInfoUploader(suite.govSgCarparkUc, suite.carparkInfoUc)
		err := upldr.Upload()
		assert.Equal(t, err, expectedErr)
	})

	t.Run("ignore error from carparkInfoUc.InsertOrUpdateByCarParkNo by logging the err and return success", func(t *testing.T) {
		suite.SetupTest()
		expectedErr := errors.New("test-error")
		suite.govSgCarparkUc.On("CarparkInfos").Return([]carparkinfo.Model{*suite.fakeCpInfoModel}, nil)
		suite.carparkInfoUc.On("InsertOrUpdateByCarParkNo", suite.fakeCpInfoModel).Return(expectedErr)
		upldr := tasks.NewCarparkInfoUploader(suite.govSgCarparkUc, suite.carparkInfoUc)
		err := upldr.Upload()
		require.NoError(t, err)
	})

	t.Run("should not return error if govSgCarparkUc.CarparkInfos is success", func(t *testing.T) {
		suite.SetupTest()
		suite.govSgCarparkUc.On("CarparkInfos").Return([]carparkinfo.Model{*suite.fakeCpInfoModel}, nil)
		suite.carparkInfoUc.On("InsertOrUpdateByCarParkNo", suite.fakeCpInfoModel).Return(nil)
		upldr := tasks.NewCarparkInfoUploader(suite.govSgCarparkUc, suite.carparkInfoUc)
		err := upldr.Upload()
		require.NoError(t, err)
	})
}
