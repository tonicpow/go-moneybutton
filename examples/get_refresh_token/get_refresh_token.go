package main

import (
	"context"
	"log"
	"os"

	"github.com/tonicpow/go-moneybutton"
)

func main() {
	client := moneybutton.NewClient(nil, nil)

	response, err := client.GetRefreshToken(
		context.Background(),
		os.Getenv("CLIENT_ID"),
		os.Getenv("AUTH_CODE"),
		os.Getenv("REDIRECT_URL"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("refresh response: ", response)
}
