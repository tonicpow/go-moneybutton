package main

import (
	"context"
	"log"
	"os"

	"github.com/tonicpow/go-moneybutton"
)

func main() {
	client := moneybutton.NewClient(nil, nil)

	response, err := client.GetUserIdentity(
		context.Background(),
		os.Getenv("ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("identity: ", response.Data)
}
