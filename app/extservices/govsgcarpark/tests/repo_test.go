package tests

import (
	"github.com/dare-rider/carpark/app/extservices/govsgcarpark"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sampleSuccessCarparkAvailablityHttpHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`
			{
			  "items": [
			    {
			      "timestamp": "2019-11-16T15:17:27+08:00",
			      "carpark_data": [
			        {
			          "carpark_info": [
			            {
			              "total_lots": "91",
			              "lot_type": "C",
			              "lots_available": "12"
			            }
			          ],
			          "carpark_number": "HE12",
			          "update_datetime": "2019-11-16T15:16:48"
			        },
							{
			          "carpark_info": [
			            {
			              "total_lots": "91",
			              "lot_type": "C",
			              "lots_available": "12"
			            }
			          ],
			          "carpark_number": "HE12",
			          "update_datetime": "2019-11-16T15:16:48"
			        }
						]
					}
				]
			}
		`))
}

func sampleFailCarparkAvailablityHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(`
			{
			  "code": 0,
			  "message": "Unable to fetch carpark info currently"
			}
		`))
}

func TestGovSgCarPark_carparkAvailablity(t *testing.T) {
	t.Run("should not return error for success response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(sampleSuccessCarparkAvailablityHttpHandler))
		defer server.Close()
		baseURL := server.URL
		rp := govsgcarpark.NewRepo(baseURL, server.Client())
		_, err := rp.CarparkAvailability()
		assert.NoError(t, err)
	})

	t.Run("should return 1 item for mock success response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(sampleSuccessCarparkAvailablityHttpHandler))
		defer server.Close()
		baseURL := server.URL
		rp := govsgcarpark.NewRepo(baseURL, server.Client())
		res, _ := rp.CarparkAvailability()
		assert.Equal(t, 1, len(res.Items))
	})

	t.Run("should return 2 carpark data inside 1 item for mock success response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(sampleSuccessCarparkAvailablityHttpHandler))
		defer server.Close()
		baseURL := server.URL
		rp := govsgcarpark.NewRepo(baseURL, server.Client())
		res, _ := rp.CarparkAvailability()
		assert.Equal(t, 2, len(res.Items[0].CarparkData))
	})

	t.Run("should return error in case of error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(sampleFailCarparkAvailablityHttpHandler))
		defer server.Close()
		baseURL := server.URL
		rp := govsgcarpark.NewRepo(baseURL, server.Client())
		_, err := rp.CarparkAvailability()
		assert.Error(t, err)
	})
}
