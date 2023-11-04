package main

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var languages = []string{"Go", "Rust", "Elixir", "TypeScript", "Ocaml", "Zig"}

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
	lang := "Go"

	var result bson.M
	err = coll.FindOne(ctx, bson.M{"lang": lang}).Decode(&result)
	if err != nil {
		panic(err)
	}
	jsonData, err := bson.MarshalExtJSON(result, false, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))

	//test
	articles := getArticles(500)
	for _, article := range articles {
		for i, lang := range languages {
			if customContains(article.name, languages[i]) {
				fmt.Printf("%s: %s -> %s\n", lang, article.name, article.url)
			}
		}
	}
}

type article struct {
	url  string
	name string
}

// Returns array of articles, article might be empty
func getArticles(number int) []article {
	articles := make([]article, number)
	// Call my program end extract the data
	cmd := exec.Command("hackns", "--no-input", "--number", fmt.Sprintf("%d", number))
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// Parse the data
	// println(string(out))
	for n, line := range strings.Split(string(out), "\n") {
		article_url := strings.Split(line, "~")
		if len(article_url) != 2 {
			continue
		}
		articles[n] = article{article_url[1], article_url[0]}
	}
	return articles
}

func customContains(s, substr string) bool {
	// Create a regex pattern to match the word with whitespace or string boundaries.
	pattern := "\\b" + regexp.QuoteMeta(substr) + "\\b"
	re := regexp.MustCompile(pattern)
	return re.MatchString(s)
}
