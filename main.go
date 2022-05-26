package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/askcloudarchitech/mediumautopost/pkg/mediumautopost"
    // "github.com/SouliyaPPS/realtime-chat-go/pkg/websocket"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"github.com/SouliyaPPS/realtime-chat-go/pkg/websocket"
)

type RequestBody struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Context string `json:"context"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	requestBody := RequestBody{}
	json.Unmarshal([]byte(request.Body), &requestBody)

	if requestBody.Payload.Context == "production" {
		mediumautopost.Do("")
	} else {
		fmt.Println("context " + requestBody.Payload.Context + " detected, skipping")
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            "Success",
		IsBase64Encoded: false,
	}, nil
}


func Handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
    name := request.QueryStringParameters["name"]
    response := fmt.Sprintf("Hello %s!", name)

    return &events.APIGatewayProxyResponse{
        StatusCode: 200,
        Headers:    map[string]string{"Content-Type": "text/html; charset=UTF-8"},
        Body:       response,
    }, nil
}

func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("websocket endpoint reached")

	conn, err := websocket.Upgrade(w, r)

	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(pool, w, r)
	})
}

func main() {
	fmt.Println("Chat")
	setupRoutes()
	http.ListenAndServe(":9000", nil)

	lambda.Start(handler)
	lambda.Start(Handler)
}