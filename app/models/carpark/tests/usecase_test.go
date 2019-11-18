package tests

import (
	"errors"
	"github.com/bxcodec/faker"
	"github.com/dare-rider/carpark/app/models/carpark"
	"github.com/dare-rider/carpark/app/models/carpark/mocks"
	"github.com/dare-rider/carpark/app/models/carparkinfo"
	cpInfoMocks "github.com/dare-rider/carpark/app/models/carparkinfo/mocks"
	"github.com/dare-rider/carpark/app/requests"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UsecaseTestSuite struct {
	suite.Suite
	rp                         *mocks.RepoI
	carparkInfoUc              *cpInfoMocks.Usecase
	fakeCarparkMod             *carpark.Model
	fakeNearestCarparksRequest *requests.NearestCarparksRequest
	fakeCarparkInfoMod         *carparkinfo.Model
}

func (suite *UsecaseTestSuite) SetupTest() {
	suite.rp = &mocks.RepoI{}
	suite.carparkInfoUc = &cpInfoMocks.Usecase{}
	suite.fakeCarparkMod = new(carpark.Model)
	_ = faker.FakeData(suite.fakeCarparkMod)
	suite.fakeNearestCarparksRequest = new(requests.NearestCarparksRequest)
	_ = faker.FakeData(suite.fakeNearestCarparksRequest)
	suite.fakeCarparkInfoMod = new(carparkinfo.Model)
	_ = faker.FakeData(suite.fakeCarparkInfoMod)
}

func TestUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}

func (suite *UsecaseTestSuite) TestInsertOrUpdate() {
	t := suite.T()
	t.Run("should return error if rp.InsertOrUpdate return error", func(t *testing.T) {
		suite.SetupTest()
		expectedErr := errors.New("test-error")
		suite.rp.On("InsertOrUpdate", suite.fakeCarparkMod).Return(expectedErr)
		uc := carpark.NewUsecase(suite.rp, suite.carparkInfoUc)
		err := uc.InsertOrUpdate(suite.fakeCarparkMod)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should not return error if rp.InsertOrUpdate does not error", func(t *testing.T) {
		suite.SetupTest()
		suite.rp.On("InsertOrUpdate", suite.fakeCarparkMod).Return(nil)
		uc := carpark.NewUsecase(suite.rp, suite.carparkInfoUc)
		err := uc.InsertOrUpdate(suite.fakeCarparkMod)
		require.NoError(t, err)
	})
}

func (suite *UsecaseTestSuite) TestFetchNearestWithInfo() {
	t := suite.T()
	t.Run("should return error if rp.FetchNearestWithInfo return error", func(t *testing.T) {
		suite.SetupTest()
		uc := carpark.NewUsecase(suite.rp, suite.carparkInfoUc)
		expectedErr := errors.New("test-error")

		limit, offset := uc.LimitOffset(suite.fakeNearestCarparksRequest.Page, suite.fakeNearestCarparksRequest.PerPage)
		suite.rp.On("FetchNearest", suite.fakeNearestCarparksRequest.Latitude, suite.fakeNearestCarparksRequest.Longitude, limit, offset).Return(nil, expectedErr)

		cpIds := []int{2, 3, 4, 5}
		suite.carparkInfoUc.On("FindAllByCarparkIDs", cpIds).Return(suite.fakeCarparkInfoMod, nil)

		_, err := uc.FetchNearestWithInfo(suite.fakeNearestCarparksRequest)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should return error if carparkInfoUc.FindAllByCarparkIDs return error", func(t *testing.T) {
		suite.SetupTest()
		uc := carpark.NewUsecase(suite.rp, suite.carparkInfoUc)
		expectedErr := errors.New("test-error")

		limit, offset := uc.LimitOffset(suite.fakeNearestCarparksRequest.Page, suite.fakeNearestCarparksRequest.PerPage)
		suite.rp.On("FetchNearest", suite.fakeNearestCarparksRequest.Latitude, suite.fakeNearestCarparksRequest.Longitude, limit, offset).Return([]carpark.Model{*suite.fakeCarparkMod}, nil)

		cpIds := []int{suite.fakeCarparkMod.ID}
		suite.carparkInfoUc.On("FindAllByCarparkIDs", cpIds).Return(nil, expectedErr)

		_, err := uc.FetchNearestWithInfo(suite.fakeNearestCarparksRequest)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should not return error if rp.FetchNearestWithInfo does not error", func(t *testing.T) {
		suite.SetupTest()
		uc := carpark.NewUsecase(suite.rp, suite.carparkInfoUc)

		limit, offset := uc.LimitOffset(suite.fakeNearestCarparksRequest.Page, suite.fakeNearestCarparksRequest.PerPage)
		suite.rp.On("FetchNearest", suite.fakeNearestCarparksRequest.Latitude, suite.fakeNearestCarparksRequest.Longitude, limit, offset).Return([]carpark.Model{*suite.fakeCarparkMod}, nil)

		cpIds := []int{suite.fakeCarparkMod.ID}
		suite.carparkInfoUc.On("FindAllByCarparkIDs", cpIds).Return([]carparkinfo.Model{*suite.fakeCarparkInfoMod}, nil)

		res, err := uc.FetchNearestWithInfo(suite.fakeNearestCarparksRequest)
		require.NoError(t, err)
		assert.Equal(t, len(res), len(cpIds))
	})
}
