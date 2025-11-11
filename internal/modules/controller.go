package modules

import (
	user "gitgub.com/Alksndr9/go-students-disciplines/internal/modules/user/controller"
	"gitgub.com/Alksndr9/go-students-disciplines/internal/responder"
	"go.uber.org/zap"
)

type Controllers struct {
	UserController *user.User
}

func NewControllers(responder responder.Responder, log *zap.Logger, storages *Storages) *Controllers {
	return &Controllers{
		UserController: user.NewUserController(responder, log, storages.UserStorage),
	}
}
