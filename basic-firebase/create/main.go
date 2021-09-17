package main

import (
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("../path/to/serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// データ追加
	_, err = client.Collection("users").Doc("user1").Set(ctx, map[string]interface{}{
		"first":  "koki",
		"last":   "nakayama",
		"born":   2001,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}

	// データ読み取り
	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
	defer client.Close()
}