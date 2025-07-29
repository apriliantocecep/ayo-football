package config

import (
	"errors"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/entity"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
)

type JwtWrapper struct {
	SecretKey                    string
	Issuer                       string
	AccessTokenExpirationMinutes int64
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (w *JwtWrapper) GenerateAccessToken(user *entity.User) (signedToken string, claims *Claims, err error) {
	claims = &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    w.Issuer,
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(w.AccessTokenExpirationMinutes) * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(w.SecretKey))

	if err != nil {
		return "", nil, err
	}

	return signedToken, claims, nil
}

func (w *JwtWrapper) ValidateToken(signedToken string) (c *Claims, err error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(w.SecretKey), nil
	})

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	return claims, nil
}

func NewJwt(vaultClient *shared.VaultClient) *JwtWrapper {
	secret := utils.GetVaultSecretAuthSvc(vaultClient)

	jwtSecretKey := secret["JWT_SECRET_KEY"]
	if jwtSecretKey == nil || jwtSecretKey == "" {
		log.Fatalf("JWT_SECRET_KEY is not set")
	}

	jwtIssuer := secret["JWT_ISSUER"]
	if jwtIssuer == nil || jwtIssuer == "" {
		log.Fatalf("JWT_ISSUER is not set")
	}

	jwtAccessTokenExpirationMinutes := secret["JWT_ACCESS_TOKEN_EXPIRATION_MINUTES"]
	if jwtAccessTokenExpirationMinutes == nil || jwtAccessTokenExpirationMinutes == "" {
		log.Fatalf("JWT_ACCESS_TOKEN_EXPIRATION_MINUTES is not set")
	}

	expiration, err := strconv.ParseInt(jwtAccessTokenExpirationMinutes.(string), 10, 64)
	if err != nil {
		log.Fatalf("Error converting JWT_ACCESS_TOKEN_EXPIRATION_MINUTES")
	}

	return &JwtWrapper{
		SecretKey:                    jwtSecretKey.(string),
		Issuer:                       jwtIssuer.(string),
		AccessTokenExpirationMinutes: expiration,
	}
}
