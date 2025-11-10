package pgxutils

import (
	"gitgub.com/Alksndr9/go-students-disciplines/internal/domain"
	"github.com/jackc/pgx/v5"
)

func UserArgs(user *domain.User) pgx.NamedArgs {
	return pgx.NamedArgs{
		"username":      user.Username,
		"password_hash": user.PasswordHash,
		"email":         user.Email,
		"phone_number":  user.PhoneNumber,
		"role":          user.Role,
	}
}
