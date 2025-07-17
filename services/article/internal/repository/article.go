package repository

import (
	"context"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type ArticleRepository struct {
	Collection *mongo.Collection
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
