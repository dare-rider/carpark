package govsgcarpark

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dare-rider/carpark/constant"
	"github.com/dare-rider/carpark/types"
	"io/ioutil"
	"net/http"
	"time"
)

type errResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResp struct {
	ApiInfo *apiInfo `json:"api_info"`
	Items   []item   `json:"items"`
}

type apiInfo struct {
	Status string `json:"status"`
}

type item struct {
	Timestamp   time.Time     `json:"timestamp"`
	CarparkData []carparkData `json:"carpark_data"`
}

type carparkData struct {
	CarparkInfos   []carparkInfo           `json:"carpark_info"`
	CarparkNumber  string                  `json:"carpark_number"`
	UpdateDatetime types.GovSgResponseTime `json:"update_datetime"`
}

type carparkInfo struct {
	TotalLots     string `json:"total_lots"`
	LotsAvailable string `json:"lots_available"`
	LotType       string `json:"lot_type"`
}

type repo struct {
	baseURL string
	client  *http.Client
}

type repoI interface {
	carparkAvaialability() (*successResp, error)
}

func newRepo(baseUrl string, client *http.Client) repoI {
	return &repo{
		baseURL: baseUrl,
		client:  client,
	}
}

func (rp *repo) carparkAvaialability() (*successResp, error) {
	req, err := http.NewRequest(http.MethodGet, rp.carparkAvailabilitylURL(), nil)
	if err != nil {
		return nil, errors.New(constant.InvalidRequest)
	}
	req.Header.Add(constant.ContentTypeHeader, constant.ContentType)

	resp, err := rp.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp errResp
		err = json.NewDecoder(bytes.NewReader(resBody)).Decode(&errResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errResp.Message)
	}
	var response successResp
	err = json.NewDecoder(bytes.NewReader(resBody)).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
