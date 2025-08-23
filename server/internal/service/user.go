package service

import (
	"time"

	"i18n-flow/internal/model"
	"i18n-flow/internal/pkg/jwt"
	"i18n-flow/internal/pkg/redis"
	"i18n-flow/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return repository.CreateUser(model.User{
		Username: username,
		Password: string(hash),
	})
}

func LoginUser(username, password string) (user *model.User, accessToken, refreshToken string, err error) {
	user, err = repository.GetUserByUsername(username)
	if err != nil || user == nil {
		return nil, "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err = jwt.GenerateTokens(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	// 保存 refresh token
	redis.Set("refresh_"+user.ID, refreshToken, time.Second*604800)
	return user, accessToken, refreshToken, nil
}
