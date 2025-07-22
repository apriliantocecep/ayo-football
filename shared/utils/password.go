package utils

import (
	"context"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(ctx context.Context, password string) string {
	_, span := otel.Tracer("utils.HashPassword").Start(ctx, "HashPassword")
	defer span.End()

	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func ComparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
