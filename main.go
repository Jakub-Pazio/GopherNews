package main

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var languages = []string{"Go", "Rust", "Elixir", "TypeScript", "Ocaml", "Zig", "Java"}

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
		Keys: bson.D{primitive.E{Key: "url", Value: 1},{Key: "keyword", Value: -1}},
		Options: options.Index().SetUnique(true),
	}
	index, err := coll.Indexes().CreateOne(ctx, indexArticle)
	if err != nil {
		fmt.Printf("Error creating index: %s\n", err)
	} else {
		fmt.Println(index)
	}

	//test
	articles := getArticles(500)
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

type Article struct {
	URL     string `bson:"url"`
	Name    string `bson:"name"`
	Keyword string `bson:"keyword"`
}

// Returns array of articles, article might be empty
func getArticles(number int) []Article {
	articles := make([]Article, 0)
	// Call my program end extract the data
	cmd := exec.Command("/usr/bin/hackns", "--no-input", "--number", fmt.Sprintf("%d", number))
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// Parse the data
	// println(string(out))
	for _, line := range strings.Split(string(out), "\n") {
		article_url := strings.Split(line, "~")
		if len(article_url) != 2 {
			continue
		}
		for _, lang := range languages {
			if customContains(article_url[0], lang) {
				articles = append(articles, Article{URL: article_url[1], Name: article_url[0], Keyword: lang})
			}
		}
	}
	return articles
}

func customContains(s, substr string) bool {
	// Create a regex pattern to match the word with whitespace or string boundaries.
	pattern := "\\b" + regexp.QuoteMeta(substr) + "\\b"
	re := regexp.MustCompile(pattern)
	return re.MatchString(s)
}
