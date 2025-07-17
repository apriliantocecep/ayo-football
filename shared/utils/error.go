package utils

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcError(err error) error {
	st, ok := status.FromError(err)
	if ok {
		return fiber.NewError(GrpcCodeToHTTPStatus(st.Code()), st.Message())
	} else {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}

func GrpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return 200
	case codes.InvalidArgument:
		return 400
	case codes.NotFound:
		return 404
	case codes.AlreadyExists:
		return 409
	case codes.PermissionDenied:
		return 403
	case codes.Unauthenticated:
		return 401
	case codes.Unavailable:
		return 503
	default:
		return 500
	}
}
