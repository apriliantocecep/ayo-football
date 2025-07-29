package messaging

import (
	"context"
	"encoding/json"
	"github.com/apriliantocecep/ayo-football/services/article/internal/model"
	"github.com/apriliantocecep/ayo-football/services/article/internal/usecase"
	sharedmodel "github.com/apriliantocecep/ayo-football/shared/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type MetadataConsumer struct {
	ModerationUseCase *usecase.ModerationUseCase
}

func (c *MetadataConsumer) Consume(delivery *amqp.Delivery) error {
	metadataEvent := new(sharedmodel.MetadataEvent)
	err := json.Unmarshal(delivery.Body, metadataEvent)
	if err != nil {
		return err
	}

	req := model.PublishArticleRequest{
		ArticleId:        metadataEvent.ArticleId,
		ModerationStatus: metadataEvent.ModerationStatus,
	}
	res, err := c.ModerationUseCase.PublishArticle(context.Background(), &req)
	if err != nil {
		return err
	}

	log.Printf("article id '%s' status is '%s'", metadataEvent.ArticleId, res.Status)

	return nil
}

func NewMetadataConsumer(moderationUseCase *usecase.ModerationUseCase) *MetadataConsumer {
	return &MetadataConsumer{ModerationUseCase: moderationUseCase}
}
