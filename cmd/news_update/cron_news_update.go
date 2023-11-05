package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopher_news/internal/articles"
)

// refactor this to a config file, change name to keywords
var languages = []string{"Go", "Rust", "Elixir", "TypeScript", "Ocaml", "Zig", "Java", "Haskell", "C++"}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	defer func() {
		if err = db.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	coll := db.Database("test").Collection("articles")

	indexArticle := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "url", Value: 1}, {Key: "keyword", Value: -1}},
		Options: options.Index().SetUnique(true),
	}
	index, err := coll.Indexes().CreateOne(ctx, indexArticle)
	if err != nil {
		fmt.Printf("Error creating index: %s\n", err)
	} else {
		fmt.Println(index)
	}

	//test
	articles := articles.GetArticles(500, languages)
	for _, article := range articles {
		fmt.Printf("%s: %s -> %s\n", article.Keyword, article.Name, article.URL)
		result, err := coll.InsertOne(ctx, article)
		if err != nil {
			fmt.Printf("Error inserting document: %s\n", err)
			fmt.Println("Article may already exist")
		}
		fmt.Println(result)
	}
}
