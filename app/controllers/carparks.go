package controllers

import (
	"github.com/dare-rider/carpark/app/models/carpark"
	"github.com/dare-rider/carpark/app/presentors"
	"github.com/dare-rider/carpark/app/requests"
	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	"net/http"
)

type carparkController struct {
	base      BaseController
	carparkUc carpark.Usecase
}

func NewCarparkController(base BaseController, carparkUc carpark.Usecase) *carparkController {
	return &carparkController{
		base:      base,
		carparkUc: carparkUc,
	}
}

// Router loads package specific routes
func (ctrl *carparkController) Router(r chi.Router) {
	r.Get("/carparks/nearest", ctrl.nearestCarparks)
}

func (ctrl *carparkController) nearestCarparks(w http.ResponseWriter, r *http.Request) {
	var inReq requests.NearestCarparksRequest
	err := schema.NewDecoder().Decode(&inReq, r.URL.Query())
	if err != nil {
		ctrl.base.WriteErrorJSONWithStatus(w, http.StatusBadRequest, err)
		return
	}
	if err = ctrl.base.RequestValidator().Struct(inReq); err != nil {
		ctrl.base.WriteErrorJSONWithStatus(w, http.StatusBadRequest, err)
		return
	}
	modResps, err := ctrl.carparkUc.FetchNearestWithInfo(&inReq)
	if err != nil {
		ctrl.base.WriteErrorJSONWithStatus(w, http.StatusBadRequest, err)
		return
	}
	// preparing required results
	var results []presentors.NearestCarparkResponse
	for _, modResp := range modResps {
		res := new(presentors.NearestCarparkResponse)
		results = append(results, *res.SerializeFromModel(&modResp))
	}
	ctrl.base.WriteJSONWithStatus(w, http.StatusOK, results)
}
