package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/fehmianac/serverless-favorite-api/infra/repository"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	userId := request.PathParameters["userId"]
	itemId := request.PathParameters["itemId"]

	if userId == "" || itemId == "" {
		return events.APIGatewayProxyResponse{
			Body:       "UserId and ItemId are required",
			StatusCode: 400,
		}, nil
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoDb := dynamodb.New(sess)
	favoriteRepository := repository.NewFavoriteRepository(dynamoDb)

	err := favoriteRepository.Delete(ctx, userId, itemId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: 500,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("favorite %s removed", itemId),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(handler)
}
