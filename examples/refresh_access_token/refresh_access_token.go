package main

import (
	"context"
	"log"
	"os"

	"github.com/tonicpow/go-moneybutton"
)

func main() {
	client := moneybutton.NewClient(nil, nil)

	response, err := client.RefreshAccessToken(
		context.Background(),
		os.Getenv("CLIENT_ID"),
		os.Getenv("ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("refresh response: ", response)
}
