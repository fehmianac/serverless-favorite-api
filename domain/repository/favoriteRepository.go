package repository

import (
	"context"

	"github.com/fehmianac/serverless-favorite-api/domain/model"
)

type IFavoriteRepository interface {
	Add(ctx context.Context, entity model.Favorite) error
	Delete(ctx context.Context, userId string, itemId string) error
	Get(ctx context.Context, userId string, itemId string) (model.Favorite, error)
	GetPaginated(ctx context.Context, userId string, nextToken string, limit int) ([]model.Favorite, string, error)
	GetItems(ctx context.Context, userId string, itemIds []string) ([]model.Favorite, error)
}
