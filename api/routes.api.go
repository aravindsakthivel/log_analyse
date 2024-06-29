package api

import (
	"log_analyse/controller"

	"github.com/go-chi/chi/v5"
)

type SRouter struct{}

func (r *SRouter) RoutesInit(router *chi.Mux) {
	routes := SRoutes{}
	DBCtrl := controller.SDBCtrl{}
	routes.setCtrl(&DBCtrl)
	router.Get("/", routes.main)
	router.Get("/health", routes.health)
	router.Get("/users/dummy", routes.createUserDummy)
	router.Post("/upload", routes.uploadFile)
	router.Post("/createUser", routes.createUser)
	router.Get("/unCompressedFiles", routes.findUnCompressedFiles)
}
