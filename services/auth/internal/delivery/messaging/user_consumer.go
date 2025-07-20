package messaging

import (
	"context"
	"encoding/json"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/usecase"
	sharedmodel "github.com/apriliantocecep/posfin-blog/shared/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type UserConsumer struct {
	RegisterUseCase *usecase.RegisterUseCase
}

func (c *UserConsumer) ConsumeUserCreated(delivery *amqp.Delivery) error {
	userEvent := new(sharedmodel.UserEvent)
	err := json.Unmarshal(delivery.Body, userEvent)
	if err != nil {
		return err
	}

	req := model.CreateUserRequest{
		Name:     userEvent.Name,
		Email:    userEvent.Email,
		Username: userEvent.Username,
		Password: userEvent.Password,
	}
	res, err := c.RegisterUseCase.CreateUser(context.Background(), &req)
	if err != nil {
		return err
	}

	log.Printf("New user created: %s", res.UserId)

	return nil
}

func NewUserConsumer(registerUseCase *usecase.RegisterUseCase) *UserConsumer {
	return &UserConsumer{RegisterUseCase: registerUseCase}
}
