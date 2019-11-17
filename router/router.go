package router

import (
	"github.com/dare-rider/carpark/app/controllers"
	"github.com/dare-rider/carpark/app/extservices/govsgcarpark"
	"github.com/dare-rider/carpark/app/models/carpark"
	"github.com/dare-rider/carpark/app/models/carparkinfo"
	"github.com/dare-rider/carpark/app/tasks"
	"github.com/dare-rider/carpark/appmiddleware"
	"github.com/dare-rider/carpark/config"
	"github.com/dare-rider/carpark/constant"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

// InitRoutes initializes the router
func InitRoutes(config *config.GeneralConfig, db *sqlx.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(appmiddleware.SetJSON)
	router.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)
	// prepare controller routes
	vldtor := validator.New()
	govSgCarparkRp := govsgcarpark.NewRepo(config.GovSgService.BaseUrl, &http.Client{Timeout: constant.DefaultHttpTimeout * time.Second})
	carparkInfoRp := carparkinfo.NewRepo(db)
	carparkRp := carpark.NewRepo(db)

	govSgCarparkUc := govsgcarpark.NewUsecase(govSgCarparkRp)
	carparkInfoUc := carparkinfo.NewUsecase(carparkInfoRp)
	carparkUc := carpark.NewUsecase(carparkRp, carparkInfoUc)

	cpInfoUploader := tasks.NewCarparkInfoUploader(govSgCarparkUc, carparkInfoUc)
	cpUploader := tasks.NewCarparkUploader(carparkUc, config.DbConfig.SeedPath)

	base := controllers.NewBaseController(config, vldtor, db)
	taskController := controllers.NewTaskController(base, cpUploader, cpInfoUploader)
	carparkController := controllers.NewCarparkController(base, carparkUc)

	// Mounting controller routes
	router.Route("/", func(r chi.Router) {
		r.Group(base.Router)
		r.Group(taskController.Router)
		r.Group(carparkController.Router)
	})
	return router
}
