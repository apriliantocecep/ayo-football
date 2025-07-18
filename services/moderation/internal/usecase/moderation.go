package usecase

import (
	"context"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/gateway/messaging"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/moderation/internal/repository"
	sharedmodel "github.com/apriliantocecep/posfin-blog/shared/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
	"strings"
)

type ModerationUseCase struct {
	DB                 *gorm.DB
	MetadataRepository *repository.MetadataRepository
	MetadataPublisher  *messaging.MetadataPublisher
}

func (u *ModerationUseCase) Process(ctx context.Context, request *model.ProcessRequest) (*model.ProcessResponse, error) {
	metadata, err := u.MetadataRepository.FindByArticleId(u.DB, request.ArticleId)
	if err != nil {
		return nil, err
	}
	//log.Printf("IsPass = %t", request.IsPass)

	moderationStatus := "review"
	if request.IsPass == true {
		moderationStatus = "accepted"
	} else {
		moderationStatus = "rejected"
	}
	metadata.ModerationStatus = moderationStatus
	err = u.MetadataRepository.Update(u.DB, metadata)
	if err != nil {
		return nil, err
	}

	// publish to broker
	event := sharedmodel.MetadataEvent{
		ArticleId:        request.ArticleId,
		ModerationStatus: moderationStatus,
	}
	err = u.MetadataPublisher.Publish(&event)
	if err != nil {
		log.Printf("failed publish moderation event : %+v", err)
		return nil, status.Errorf(codes.Aborted, "failed to publish moderation")
	}

	response := model.ProcessResponse{ModerationStatus: moderationStatus}
	return &response, nil
}

func (u *ModerationUseCase) CheckContent(ctx context.Context, request *model.CheckContentRequest) (*model.CheckContentResponse, error) {
	var badWords = []string{
		"anjing", "bangsat", "kontol", "tai", "babi", "goblok", "tolol", "ngentot",
	}
	lowered := strings.ToLower(request.Content)
	for _, word := range badWords {
		if strings.Contains(lowered, word) {
			return &model.CheckContentResponse{IsPass: false}, nil
		}
	}
	return &model.CheckContentResponse{IsPass: true}, nil
}

func NewModerationUseCase(DB *gorm.DB, metadataRepository *repository.MetadataRepository, metadataPublisher *messaging.MetadataPublisher) *ModerationUseCase {
	return &ModerationUseCase{
		DB:                 DB,
		MetadataRepository: metadataRepository,
		MetadataPublisher:  metadataPublisher,
	}
}
