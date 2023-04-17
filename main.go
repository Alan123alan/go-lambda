package main

import (
	"context"
	"fmt"

	"github.com/gocolly/colly"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Event struct {
	Name string `json:"Name" `
}

type Pokemon struct {
	Name        string
	Currency    string
	Price       string
	Description string
	Stock       string
}

func HandleRequest(ctx context.Context, event Event) (Event, error) {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	dynamoDB := dynamodb.New(session)
	// Create table Movies
	tableName := "Pokemons"

	pokemon := getScrapedPokemonData(event.Name)

	payload, err := dynamodbattribute.MarshalMap(pokemon)
	if err != nil {
		fmt.Println("Failed to marshal request")
		return Event{}, err
	}

	input := &dynamodb.PutItemInput{
		Item:      payload,
		TableName: aws.String(tableName),
	}
	_, err = dynamoDB.PutItem(input)
	if err != nil {
		fmt.Println("Failed to write to db")
		return Event{}, err
	}

	return event, nil
}

func main() {
	lambda.Start(HandleRequest)
}

func getScrapedPokemonData(pokemonName string) Pokemon {
	collector := colly.NewCollector()
	pokemon := Pokemon{}
	collector.OnHTML("h1[class='product_title entry-title']", func(element *colly.HTMLElement) {
		pokemon.Name = element.Text
	})
	collector.OnHTML("p[class='price']", func(element *colly.HTMLElement) {
		pokemon.Price = element.Text
	})
	collector.OnHTML("div[class='woocommerce-product-details__short-description']", func(element *colly.HTMLElement) {
		pokemon.Description = element.Text
	})
	collector.OnHTML("p[class='stock in-stock']", func(element *colly.HTMLElement) {
		pokemon.Stock = element.Text
	})
	collector.Visit(fmt.Sprintf("https://scrapeme.live/shop/%v/", pokemonName))
	return pokemon
}
