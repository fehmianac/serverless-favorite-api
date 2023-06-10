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

	println("putItem", putItem)
	println("input", input)

	_, err = favoriteRepository.db.PutItemWithContext(ctx, input)
	if err != nil {
		println("err", err)
	}

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

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("favorites"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(userId),
			},
		},
		Limit: aws.Int64(int64(limit)),
	}
	if nextToken != "" {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"sk": {
				S: aws.String(nextToken),
			},
		}
	}
	response, err := favoriteRepository.db.QueryWithContext(ctx, queryInput)

	if err != nil {
		println("err", err)
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
