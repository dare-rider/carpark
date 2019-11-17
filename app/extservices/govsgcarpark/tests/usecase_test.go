package tests

import (
	"errors"
	"github.com/bxcodec/faker"
	"github.com/dare-rider/carpark/app/extservices/govsgcarpark"
	"github.com/dare-rider/carpark/app/extservices/govsgcarpark/mocks"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UsecaseTestSuite struct {
	suite.Suite
	rp                           *mocks.RepoI
	fakeRpCarparkAvailabilityRes *govsgcarpark.SuccessResp
}

func (suite *UsecaseTestSuite) SetupTest() {
	suite.rp = &mocks.RepoI{}
	suite.fakeRpCarparkAvailabilityRes = new(govsgcarpark.SuccessResp)
	_ = faker.FakeData(suite.fakeRpCarparkAvailabilityRes)
}

func TestUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}

func (suite *UsecaseTestSuite) TestCarparkInfos() {
	t := suite.T()

	t.Run("should return error if rp.CarparkAvailability return error", func(t *testing.T) {
		suite.SetupTest()
		expectedErr := errors.New("test-error")
		suite.rp.On("CarparkAvailability").Return(nil, expectedErr)
		uc := govsgcarpark.NewUsecase(suite.rp)
		_, err := uc.CarparkInfos()
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should return same length carparkinfo.Model if rp.CarparkAvailability return result", func(t *testing.T) {
		suite.SetupTest()
		// length of carpark_data
		ctr := 0
		for _, item := range suite.fakeRpCarparkAvailabilityRes.Items {
			for _, cp := range item.CarparkData {
				ctr += len(cp.CarparkInfos)
			}
		}
		suite.rp.On("CarparkAvailability").Return(suite.fakeRpCarparkAvailabilityRes, nil)
		uc := govsgcarpark.NewUsecase(suite.rp)
		res, err := uc.CarparkInfos()
		require.NoError(t, err)
		assert.Equal(t, len(res), ctr)
	})
}
