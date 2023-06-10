package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	contract "github.com/fehmianac/serverless-favorite-api/cmd/contract/favorite"
	"github.com/fehmianac/serverless-favorite-api/domain/model"
	"github.com/fehmianac/serverless-favorite-api/infra/repository"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	userId := request.PathParameters["userId"]
	itemIds := strings.Split(request.QueryStringParameters["itemIds"], ",")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoDb := dynamodb.New(sess)
	favoriteRepository := repository.NewFavoriteRepository(dynamoDb)

	for _, itemId := range itemIds {
		println("itemId", itemId)
	}

	if request.QueryStringParameters["itemIds"] == "" {
		nextToken := request.QueryStringParameters["nextToken"]
		limit := request.QueryStringParameters["limit"]
		if limit == "" {
			limit = "100"
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			println(err.Error())
			return events.APIGatewayProxyResponse{
				Body:       "invalid limit",
				StatusCode: 400,
			}, nil

		}

		items, nextToken, err := favoriteRepository.GetPaginated(ctx, userId, nextToken, limitInt)
		if err != nil {
			println(err.Error())
			return events.APIGatewayProxyResponse{
				Body:       "internal server error",
				StatusCode: 500,
			}, nil

		}

		result := []contract.ListFavoriteResponseItem{}
		for _, item := range items {
			result = append(result, contract.ListFavoriteResponseItem{
				UserId:    userId,
				ItemId:    item.ItemId,
				CreatedAt: item.CreatedAt,
			})
		}
		jsonString, _ := json.Marshal(contract.ListFavoriteResponse{
			Items:     result,
			NextToken: nextToken,
		})

		return events.APIGatewayProxyResponse{
			Body:       string(jsonString),
			StatusCode: 200,
		}, nil

	}
	items, err := favoriteRepository.GetItems(ctx, userId, itemIds)
	if err != nil {
		println(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: 500,
		}, nil

	}
	result := []contract.CheckFavoriteResponse{}
	for _, item := range itemIds {
		result = append(result, contract.CheckFavoriteResponse{
			UserId:     userId,
			ItemId:     item,
			IsFavorite: checkInList(item, items),
		})
	}
	jsonString, _ := json.Marshal(result)

	return events.APIGatewayProxyResponse{
		Body:       string(jsonString),
		StatusCode: 200,
	}, nil

}

func checkInList(item string, items []model.Favorite) bool {
	for _, i := range items {
		if i.ItemId == item {
			return true
		}
	}
	return false

}

func main() {
	lambda.Start(handler)
}
