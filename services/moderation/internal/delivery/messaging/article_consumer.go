package messaging

import (
	"context"
	"encoding/json"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/usecase"
	sharedmodel "github.com/apriliantocecep/posfin-blog/shared/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type ArticleConsumer struct {
	MetadataUseCase *usecase.MetadataUseCase
}

func (c *ArticleConsumer) Consume(delivery amqp.Delivery) error {
	articleEvent := new(sharedmodel.ArticleEvent)
	err := json.Unmarshal(delivery.Body, articleEvent)
	if err != nil {
		return err
	}
	log.Printf("start consume article id: %s", articleEvent.ID)

	// TODO
	req := model.MetadataRequest{
		ArticleId: articleEvent.ID,
		Title:     articleEvent.Title,
		Author:    articleEvent.Author,
	}
	_, err = c.MetadataUseCase.Save(context.Background(), &req)
	if err != nil {
		return err
	}

	//res.MetadataId

	return nil
}

func NewArticleConsumer(metadataUseCase *usecase.MetadataUseCase) *ArticleConsumer {
	return &ArticleConsumer{MetadataUseCase: metadataUseCase}
}
