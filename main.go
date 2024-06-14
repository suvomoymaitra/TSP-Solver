package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const MAX = 1e9

func tsp(dist [][]float64) (float64, []int) {
	n := len(dist)
	memo := make([][]float64, n)
	for i := range memo {
		memo[i] = make([]float64, 1<<n)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	// dp function to find the minimum cost
	var dp func(pos, mask int) float64
	dp = func(pos, mask int) float64 {
		if mask == (1<<n)-1 {
			return dist[pos][0]
		}
		if memo[pos][mask] != -1 {
			return memo[pos][mask]
		}
		minCost := MAX
		for city := 0; city < n; city++ {
			if city != pos && (mask&(1<<city)) == 0 {
				newCost := dist[pos][city] + dp(city, mask|(1<<city))
				if newCost < minCost {
					minCost = newCost
				}
			}
		}
		memo[pos][mask] = minCost
		return minCost
	}

	minDistance := dp(0, 1)

	// Function to reconstruct the path
	findPath := func() []int {
		mask := 1
		pos := 0
		path := []int{0}

		for i := 1; i < n; i++ {
			bestCity := -1
			bestCost := MAX
			for city := 0; city < n; city++ {
				if (mask & (1 << city)) == 0 {
					currentCost := dist[pos][city] + memo[city][mask|(1<<city)]
					if currentCost < bestCost {
						bestCost = currentCost
						bestCity = city
					}
				}
			}
			path = append(path, bestCity)
			pos = bestCity
			mask |= 1 << bestCity
		}
		path = append(path, 0)
		return path
	}

	minPath := findPath()
	return minDistance, minPath
}

func getCostMatrixFromString(costMatrixString string, numberOfPoints int) [][]float64 {

	strValues := strings.Split(costMatrixString, ",")

	tempArr := make([]float64, len(strValues))

	// Convert each substring to a float64
	for i, str := range strValues {
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			fmt.Println("Error converting string to float64:", err)
		}
		tempArr[i] = value
	}

	arr := make([][]float64, numberOfPoints)
	for i := range arr {
		arr[i] = make([]float64, numberOfPoints)
	}
	k := 0

	for i := 0; i < numberOfPoints; i++ {
		for j := 0; j < numberOfPoints; j++ {
			arr[i][j] = tempArr[k]
			k++
		}
	}
	return arr
}

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
