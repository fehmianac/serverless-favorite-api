package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fehmianac/serverless-favorite-api/domain/model"
	"github.com/fehmianac/serverless-favorite-api/domain/repository"
)

type favoriteRepository struct {
	db *dynamodb.DynamoDB
}

// Add implements repository.IFavoriteRepository
func (favoriteRepository *favoriteRepository) Add(ctx context.Context, entity model.Favorite) error {
	putItem, err := dynamodbattribute.MarshalMap(entity)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("favorites"),
		Item:      putItem,
	}

	_, err = favoriteRepository.db.PutItemWithContext(ctx, input)
	return err

}

// Delete implements repository.IFavoriteRepository
func (favoriteRepository *favoriteRepository) Delete(ctx context.Context, userId string, itemId string) error {
	favoriteRepository.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("favorites"),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(userId),
			},
			"sk": {
				S: aws.String(itemId),
			},
		},
	})
	return nil
}

// Get implements repository.IFavoriteRepository
func (favoriteRepository *favoriteRepository) Get(ctx context.Context, userId string, itemId string) (model.Favorite, error) {
	item, err := favoriteRepository.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("favorites"),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(userId),
			},
			"sk": {
				S: aws.String(itemId),
			},
		},
	})
	if err != nil {
		return model.Favorite{}, err
	}

	var favorite model.Favorite
	err = dynamodbattribute.UnmarshalMap(item.Item, &favorite)
	if err != nil {
		return model.Favorite{}, err
	}
	return favorite, nil
}

// GetItems implements repository.IFavoriteRepository
func (favoriteRepository *favoriteRepository) GetItems(ctx context.Context, userId string, itemIds []string) ([]model.Favorite, error) {

	requestedItem := []map[string]*dynamodb.AttributeValue{}

	for _, itemId := range itemIds {
		requestedItem = append(requestedItem, map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(userId),
			},
			"sk": {
				S: aws.String(itemId),
			},
		})
	}

	items, err := favoriteRepository.db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			"favorites": {
				Keys: requestedItem,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var favorites []model.Favorite
	err = dynamodbattribute.UnmarshalListOfMaps(items.Responses["favorites"], &favorites)
	if err != nil {
		return nil, err
	}

	return favorites, nil

}

// GetPaginated implements repository.IFavoriteRepository
func (favoriteRepository *favoriteRepository) GetPaginated(ctx context.Context, userId string, nextToken string, limit int) ([]model.Favorite, string, error) {
	response, err := favoriteRepository.db.Query(&dynamodb.QueryInput{
		TableName: aws.String("favorites"),
		KeyConditions: map[string]*dynamodb.Condition{
			"pk": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(userId),
					},
				},
			},
		},
		Limit:            aws.Int64(int64(limit)),
		ScanIndexForward: aws.Bool(false),
		ExclusiveStartKey: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(userId),
			},
			"sk": {
				S: aws.String(nextToken),
			},
		},
	})

	if err != nil {
		return nil, "", err
	}

	var favorites []model.Favorite
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &favorites)
	if err != nil {
		return nil, "", err
	}

	var nextTokenResponse string
	if response.LastEvaluatedKey != nil {
		nextTokenResponse = *response.LastEvaluatedKey["sk"].S
	}

	return favorites, nextTokenResponse, nil

}

// NewMysqlRepository ...
func NewFavoriteRepository(db *dynamodb.DynamoDB) repository.IFavoriteRepository {
	return &favoriteRepository{db}
}
