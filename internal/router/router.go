package router

import (
	"gitgub.com/Alksndr9/go-students-disciplines/internal/modules"
	"github.com/go-chi/chi"
)

func NewRouter(controllers *modules.Controllers) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Post("/", controllers.UserController.CreateUser)
		r.Get("/{id}", controllers.UserController.GetUserByID)
		r.Post("/update/{id}", controllers.UserController.UpdateUser)
		r.Get("/delete/{id}", controllers.UserController.DeleteUser)
	})
	return r
}
