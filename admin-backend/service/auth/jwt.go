package auth

import (
	"errors"
	"i18n-flow/config"
	"i18n-flow/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaim 定义JWT的claim
type JWTClaim struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(user *model.User) (string, error) {
	cfg := config.GetConfig()

	// 设置token有效期
	expirationTime := time.Now().Add(time.Hour * time.Duration(cfg.JWT.ExpirationHours))

	// 创建claims
	claims := &JWTClaim{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "i18n-flow-admin",
			Subject:   user.Username,
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken 生成刷新token
func GenerateRefreshToken(user *model.User) (string, error) {
	cfg := config.GetConfig()

	// 设置refresh token有效期(更长)
	expirationTime := time.Now().Add(time.Hour * time.Duration(cfg.JWT.RefreshExpirationHours))

	// 创建claims
	claims := &JWTClaim{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "i18n-flow-admin-refresh",
			Subject:   user.Username,
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(cfg.JWT.RefreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT token
func ValidateToken(tokenString string) (*JWTClaim, error) {
	cfg := config.GetConfig()

	// 解析token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// 验证token是否有效
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 验证并返回claims
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	// 检查token是否过期
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// ValidateRefreshToken 验证刷新token
func ValidateRefreshToken(tokenString string) (*JWTClaim, error) {
	cfg := config.GetConfig()

	// 解析token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.RefreshSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// 验证token是否有效
	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// 验证并返回claims
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.New("couldn't parse refresh token claims")
	}

	// 检查token是否过期
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	return claims, nil
}
