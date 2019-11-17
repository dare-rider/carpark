package controllers

import (
	"encoding/json"
	"github.com/dare-rider/carpark/app/presentors"
	"github.com/dare-rider/carpark/config"
	"github.com/dare-rider/carpark/constant"
	"github.com/dare-rider/carpark/utils"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

// baseController controller contains the setup required by every other controller
type baseController struct {
	db           *sqlx.DB
	env          string
	reqValidator *validator.Validate
}

// BaseController interface used by external components
type BaseController interface {
	Router(r chi.Router)
	DB() *sqlx.DB
	Env() string
	RequestValidator() *validator.Validate
	WriteJSON(w http.ResponseWriter, result interface{}) error
	WriteJSONWithStatus(w http.ResponseWriter, statusCode int, v interface{}) error
	WriteErrorJSONWithStatus(w http.ResponseWriter, statusCode int, err error) error
	DefaultSuccessResponse(w http.ResponseWriter)
}

// NewBaseController return the instance prepopulated of `BaseController` struct
func NewBaseController(
	generalConfig *config.GeneralConfig,
	reqValidator *validator.Validate,
	db *sqlx.DB) BaseController {
	it := new(baseController)
	it.env = generalConfig.MiscConfig.Environment
	it.db = db
	it.reqValidator = reqValidator
	return it
}

// Router loads package specific routes
func (base *baseController) Router(r chi.Router) {
	r.Get("/ping", base.ping)
}

// DB returns the active db connection
func (base *baseController) DB() *sqlx.DB {
	return base.db
}

// Env returns the active db connection
func (base *baseController) Env() string {
	return base.env
}

// Env returns the active db connection
func (base *baseController) RequestValidator() *validator.Validate {
	return base.reqValidator
}

// WriteJSON is used to render json response without HTTP STATUS.
func (base *baseController) WriteJSON(w http.ResponseWriter, result interface{}) error {
	data, err := json.Marshal(result)
	utils.HandleError(err)
	_, err = w.Write(data)
	return err
}

// WriteJSONWithStatus is used to render json response with HTTP STATUS.
func (base *baseController) WriteJSONWithStatus(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.WriteHeader(statusCode)
	return base.WriteJSON(w, v)
}

func (base *baseController) WriteErrorJSONWithStatus(w http.ResponseWriter, statusCode int, err error) error {
	w.WriteHeader(statusCode)
	resp := &presentors.DefaultResponse{
		Status:  false,
		Message: err.Error(),
	}
	return base.WriteJSON(w, resp)
}

func (base *baseController) DefaultSuccessResponse(w http.ResponseWriter) {
	resp := &presentors.DefaultResponse{
		Status:  true,
		Message: constant.RespSuccessMessage,
	}
	base.WriteJSONWithStatus(w, http.StatusOK, resp)
}

// Application test request
func (base *baseController) ping(w http.ResponseWriter, r *http.Request) {
	status := &presentors.Ping{
		Status:      "ok",
		Environment: base.env}
	base.WriteJSONWithStatus(w, http.StatusOK, status)
}
