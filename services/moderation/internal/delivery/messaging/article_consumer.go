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
	MetadataUseCase   *usecase.MetadataUseCase
	ModerationUseCase *usecase.ModerationUseCase
}

func (c *ArticleConsumer) ConsumeArticleModeration(delivery *amqp.Delivery) error {
	articleEvent := new(sharedmodel.ArticleEvent)
	err := json.Unmarshal(delivery.Body, articleEvent)
	if err != nil {
		return err
	}
	//log.Printf("Content: %s", articleEvent.Content)

	req := model.CheckContentRequest{Content: articleEvent.Content}
	res, err := c.ModerationUseCase.CheckContent(context.Background(), &req)
	if err != nil {
		return err
	}

	processReq := model.ProcessRequest{
		ArticleId: articleEvent.GetId(),
		IsPass:    res.IsPass,
	}
	processRes, err := c.ModerationUseCase.Process(context.Background(), &processReq)
	if err != nil {
		return err
	}
	log.Printf("article id '%s' schedule to set '%s'", articleEvent.ID, processRes.ModerationStatus)

	return nil
}

func (c *ArticleConsumer) ConsumeArticleCreated(delivery *amqp.Delivery) error {
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

func NewArticleConsumer(metadataUseCase *usecase.MetadataUseCase, moderationUseCase *usecase.ModerationUseCase) *ArticleConsumer {
	return &ArticleConsumer{
		MetadataUseCase:   metadataUseCase,
		ModerationUseCase: moderationUseCase,
	}
}
