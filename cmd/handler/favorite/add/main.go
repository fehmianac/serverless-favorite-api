package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/fehmianac/serverless-favorite-api/domain/model"
	"github.com/fehmianac/serverless-favorite-api/infra/repository"
)

type RequestModel struct {
	ItemId string `json:"itemId"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	requestModel := RequestModel{}
	err := json.Unmarshal([]byte(request.Body), &requestModel)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request body",
			StatusCode: 400,
		}, nil
	}

	userId := request.PathParameters["userId"]
	itemId := requestModel.ItemId

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoDb := dynamodb.New(sess)
	favoriteRepository := repository.NewFavoriteRepository(dynamoDb)

	err = favoriteRepository.Add(ctx, model.Favorite{
		Pk:        userId,
		Sk:        itemId,
		UserId:    userId,
		ItemId:    itemId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: 500,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("favorite %s added", itemId),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(handler)
}
