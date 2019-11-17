package controllers

import (
	"github.com/dare-rider/carpark/app/tasks"
	"github.com/go-chi/chi"
	"net/http"
)

type taskController struct {
	base                BaseController
	carparkUploader     tasks.CarparkUploader
	carparkInfoUploader tasks.CarparkInfoUploader
}

func NewTaskController(base BaseController, cu tasks.CarparkUploader, ciu tasks.CarparkInfoUploader) *taskController {
	return &taskController{
		base:                base,
		carparkUploader:     cu,
		carparkInfoUploader: ciu,
	}
}

// Router loads package specific routes
func (ctrl *taskController) Router(r chi.Router) {
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/carpark_upload", ctrl.carparkUpload)
		r.Get("/carparkinfo_upload", ctrl.carparkInfoUpload)
	})
}

func (ctrl *taskController) carparkUpload(w http.ResponseWriter, r *http.Request) {
	err := ctrl.carparkUploader.Upload()
	if err != nil {
		ctrl.base.WriteErrorJSONWithStatus(w, http.StatusUnprocessableEntity, err)
		return
	}
	ctrl.base.DefaultSuccessResponse(w)
}

func (ctrl *taskController) carparkInfoUpload(w http.ResponseWriter, r *http.Request) {
	err := ctrl.carparkInfoUploader.Upload()
	if err != nil {
		ctrl.base.WriteErrorJSONWithStatus(w, http.StatusUnprocessableEntity, err)
		return
	}
	ctrl.base.DefaultSuccessResponse(w)
}
