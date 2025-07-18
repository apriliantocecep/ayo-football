package repository

import (
	"context"
	"errors"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ArticleRepository struct {
	Collection *mongo.Collection
}

func (a *ArticleRepository) GetByOwnedId(id string, userId string) (*entity.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	articleId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid objectId")
	}

	article := new(entity.Article)
	//filter := bson.D{{"_id", articleId}}
	filter := bson.M{
		"_id":     articleId,
		"user_id": userId,
	}
	err = a.Collection.FindOne(ctx, filter).Decode(article)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "article not found")
		} else {
			return nil, err
		}
	}

	return article, nil
}

func (a *ArticleRepository) Update(id string, updateData bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	articleId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid objectId")
	}

	filter := bson.M{"_id": articleId}
	count, err := a.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return status.Errorf(codes.Aborted, "failed to check existence: %v", err)
	}
	if count == 0 {
		return status.Errorf(codes.NotFound, "article not found")
	}

	update := bson.M{"$set": updateData}
	_, err = a.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return status.Errorf(codes.Aborted, "failed to update article: %v", err)
	}

	return nil
}

func (a *ArticleRepository) Create(article *entity.Article) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := a.Collection.InsertOne(ctx, bson.D{
		{"content", article.Content},
		{"user_id", article.UserId},
		{"status", article.Status},
	})
	id := res.InsertedID.(bson.ObjectID)
	idStr := id.Hex()
	return idStr, err
}

func NewArticleRepository(collection *mongo.Collection) *ArticleRepository {
	return &ArticleRepository{Collection: collection}
}
