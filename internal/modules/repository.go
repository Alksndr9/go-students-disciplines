package modules

import (
	user "gitgub.com/Alksndr9/go-students-disciplines/internal/modules/user/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storages struct {
	UserStorage *user.Repo
}

func NewStorages(db *pgxpool.Pool) *Storages {
	return &Storages{
		UserStorage: user.NewRepo(db),
	}
}
