package tests

import (
	"errors"
	"github.com/bxcodec/faker"
	"github.com/dare-rider/carpark/app/models/carparkinfo"
	"github.com/dare-rider/carpark/app/models/carparkinfo/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UsecaseTestSuite struct {
	suite.Suite
	rp              *mocks.RepoI
	fakeCpInfoModel *carparkinfo.Model
}

func (suite *UsecaseTestSuite) SetupTest() {
	suite.rp = &mocks.RepoI{}
	suite.fakeCpInfoModel = new(carparkinfo.Model)
	_ = faker.FakeData(suite.fakeCpInfoModel)
}

func TestUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}

func (suite UsecaseTestSuite) TestInsertOrUpdateByCarParkNo() {
	t := suite.T()
	t.Run("should return an error if rp.InsertOrUpdateByCarParkNo return error", func(t *testing.T) {
		suite.SetupTest()
		expectedErr := errors.New("test-error")
		suite.rp.On("InsertOrUpdateByCarParkNo", suite.fakeCpInfoModel).Return(expectedErr)
		uc := carparkinfo.NewUsecase(suite.rp)
		err := uc.InsertOrUpdateByCarParkNo(suite.fakeCpInfoModel)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should not return error if rp.InsertOrUpdateByCarParkNo is success", func(t *testing.T) {
		suite.SetupTest()
		suite.rp.On("InsertOrUpdateByCarParkNo", suite.fakeCpInfoModel).Return(nil)
		uc := carparkinfo.NewUsecase(suite.rp)
		err := uc.InsertOrUpdateByCarParkNo(suite.fakeCpInfoModel)
		require.NoError(t, err)
	})
}

func (suite UsecaseTestSuite) TestFindAllByCarparkIDs() {
	t := suite.T()
	t.Run("should return an error if rp.FindAllByCarparkIDs return error", func(t *testing.T) {
		suite.SetupTest()
		expectedErr := errors.New("test-error")
		cpIds := []int{1, 2, 3, 4}
		suite.rp.On("FindAllByCarparkIDs", cpIds).Return(nil, expectedErr)
		uc := carparkinfo.NewUsecase(suite.rp)
		_, err := uc.FindAllByCarparkIDs(cpIds)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should return the same result if rp.FindAllByCarparkIDs is success", func(t *testing.T) {
		suite.SetupTest()
		suite.rp.On("FindAllByCarparkIDs", []int{suite.fakeCpInfoModel.CarparkID}).Return([]carparkinfo.Model{*suite.fakeCpInfoModel}, nil)
		uc := carparkinfo.NewUsecase(suite.rp)
		res, err := uc.FindAllByCarparkIDs([]int{suite.fakeCpInfoModel.CarparkID})
		require.NoError(t, err)
		assert.Equal(t, res[0].CarparkID, suite.fakeCpInfoModel.CarparkID)
	})
}
