package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	DistanceMatrix string `json:"distance_matrix"`
	NumberOfPoints int32  `json:"number_of_points"`
}

type MyResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (MyResponse, error) {
	// Log the entire request object
	requestJSON, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
	} else {
		log.Printf("Received request: %s", requestJSON)
	}

	// Get the HTTP method
	method := request.HTTPMethod

	if method != "GET" {
		return MyResponse{
			StatusCode: 404,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"message": "Resource not found"}`,
		}, nil
	}
	log.Printf("HTTP method: %s", method)

	// Unmarshal the body into MyEvent
	var event MyEvent
	err = json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		log.Printf("Error unmarshaling request body: %v", err)
		return MyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"message": "Invalid request body"}`,
		}, nil
	}

	optimalSequence := "calculate optimal sequence"

	body, err := json.Marshal(map[string]string{
		"result": optimalSequence,
	})
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return MyResponse{}, err
	}

	response := MyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}

	// Log the response
	log.Printf("Response: %+v", response)

	return response, nil
}

func main() {
	log.Println("Starting Lambda function")
	lambda.Start(HandleRequest)
}
