package model

import (
	"i18n-flow/internal/pkg/role"
)

type User struct {
	ID       string      `json:"id" bson:"_id,omitempty"`
	Username string      `json:"username" bson:"username"`
	Password string      `json:"password" bson:"password"`
	Roles    []role.Role `json:"roles" bson:"roles"`
}
