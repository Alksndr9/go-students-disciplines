package router

import (
	"gitgub.com/Alksndr9/go-students-disciplines/internal/modules/user/controller"
	"github.com/go-chi/chi"
)

func NewRouter(user *controller.User) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Post("/", user.CreateUser)
		r.Get("/{id}", user.GetUserByID)
		r.Post("/update/{id}", user.UpdateUser)
		r.Get("/delete/{id}", user.DeleteUser)
	})
	return r
}
