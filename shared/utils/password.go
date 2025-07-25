package utils

import (
	"context"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResult struct {
	Password string
	Err      error
}

func HashPassword(ctx context.Context, password string) (string, error) {
	_, span := otel.Tracer("utils.HashPassword").Start(ctx, "HashPassword")
	defer span.End()

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func HashPasswordAsync(password string, ch chan PasswordResult) {
	var passwordResult PasswordResult
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		passwordResult.Err = err
		passwordResult.Password = ""
	}
	passwordResult.Err = nil
	passwordResult.Password = string(bytes)
	ch <- passwordResult
}

func ComparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
