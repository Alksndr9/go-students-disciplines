package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"gitgub.com/Alksndr9/go-students-disciplines/internal/domain"
	"gitgub.com/Alksndr9/go-students-disciplines/internal/models"
	"gitgub.com/Alksndr9/go-students-disciplines/internal/modules/user/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
}

type Storage interface {
	SaveUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id uint64) (*domain.User, error)
	UpdateUser(ctx context.Context, id uint64, user *domain.User) error
	DeleteUser(ctx context.Context, id uint64) error
}

type User struct {
	log     *zap.Logger
	storage Storage
	Responder
}

func NewUserController(responder Responder, log *zap.Logger, storage Storage) *User {
	return &User{
		log:       log,
		storage:   storage,
		Responder: responder,
	}
}

func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.save.New"
	log := u.log.With(
		zap.String("op", op),
		zap.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req *domain.User

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", zap.Error(err))

		u.ErrorBadRequest(w, err)
		return
	}

	log.Info("request body decoded", zap.Any("request", req))

	err = u.storage.SaveUser(r.Context(), req)
	if errors.Is(err, repository.ErrUserExists) {
		log.Info("user already exists", zap.String("user", req.Username))

		u.ErrorBadRequest(w, err)
		return
	}
	if err != nil {
		log.Info("failed to add user", zap.Error(err))

		u.ErrorInternal(w, err)
		return
	}

	log.Info("user added", zap.String("user", req.Username))

	u.OutputJSON(w, models.Response{
		Success: true,
		Data: models.Data{
			Message: "user created successfully",
		},
	})
}

func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.get.ByID"
	log := u.log.With(
		zap.String("op", op),
		zap.String("request_id", middleware.GetReqID(r.Context())),
	)

	URLParamID := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(URLParamID, 10, 64)
	if err != nil {
		log.Error("failed to get url param", zap.Error(err))

		u.ErrorBadRequest(w, err)
		return
	}

	User, err := u.storage.GetUserByID(r.Context(), ID)
	if errors.Is(err, repository.ErrUserNotFound) {
		log.Info("user not found", zap.Uint64("id", ID))

		u.ErrorBadRequest(w, err)
		return
	}
	if err != nil {
		log.Info("failed to get user", zap.Error(err))

		u.ErrorInternal(w, err)
		return
	}

	log.Info("got user", zap.Uint64("id", ID))

	u.OutputJSON(w, models.Response{
		Success: true,
		Data:    User,
	})
}

func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.update.ByID"
	log := u.log.With(
		zap.String("op", op),
		zap.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req *domain.User

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", zap.Error(err))

		u.ErrorBadRequest(w, err)
		return
	}

	log.Info("request body decoded", zap.Any("request", req))

	URLParamID := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(URLParamID, 10, 64)
	if err != nil {
		log.Error("failed to get url param", zap.Error(err))

		u.ErrorBadRequest(w, err)
		return
	}

	err = u.storage.UpdateUser(r.Context(), ID, req)
	if errors.Is(err, repository.ErrUserNotFound) {
		log.Info("user not found", zap.Uint64("id", ID))

		u.ErrorBadRequest(w, err)
		return
	}
	if err != nil {
		log.Info("failed to update user", zap.Error(err))

		u.ErrorInternal(w, err)
		return
	}

	u.OutputJSON(w, models.Response{
		Success: true,
		Data: models.Data{
			Message: "user updated successfully",
		},
	})
}

func (u *User) DeleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.delete.ByID"
	log := u.log.With(
		zap.String("op", op),
		zap.String("request_id", middleware.GetReqID(r.Context())),
	)

	URLParamID := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(URLParamID, 10, 64)
	if err != nil {
		log.Error("failed to get url param", zap.Error(err))

		u.ErrorBadRequest(w, err)
		return
	}

	err = u.storage.DeleteUser(r.Context(), ID)
	if errors.Is(err, repository.ErrUserNotFound) {
		log.Info("user not found", zap.Uint64("id", ID))

		u.ErrorBadRequest(w, err)
		return
	}
	if err != nil {
		log.Info("failed to delete user", zap.Error(err))

		u.ErrorInternal(w, err)
		return
	}

	u.OutputJSON(w, models.Response{
		Success: true,
		Data: models.Data{
			Message: "user deleted successfully",
		},
	})
}
