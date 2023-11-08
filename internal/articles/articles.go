package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Article struct {
	URL          string    `bson:"url"`
	Name         string    `bson:"name"`
	Keyword      string    `bson:"keyword"`
	CreationDate time.Time `bson:"creation_date"`
	Read         bool      `bson:"read"`
}

func NewArticle(url string, name string, keyword string) Article {
	return Article{URL: url, Name: name, Keyword: keyword, CreationDate: time.Now(), Read: false}
}

func GetArticles(number int, languages []string) []Article {
	articles := make([]Article, 0)
	// Call my program and extract the data
	cmd := exec.Command("/usr/bin/hackns", "--json", "--number", fmt.Sprintf("%d", number))
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Parse the data from an array of JSON objects
	var articlesFromJSON []map[string]interface{}
	if err := json.Unmarshal(out, &articlesFromJSON); err != nil {
		panic(err)
	}

	// Iterate through the articles and extract the relevant data
	for _, articleData := range articlesFromJSON {
		articleURL := articleData["url"].(string)
		articleTitle := articleData["title"].(string)

		// Check if the article's title contains one of the specified languages
		for _, lang := range languages {
			if StrictContains(articleTitle, lang) {
				articles = append(articles, NewArticle(articleURL, articleTitle, lang))
			}
		}
	}

	fmt.Printf("Total number of matching articles: %d\n", len(articles))
	return articles
}


func GetArticlesFromDatabase() []Article {
	articleList := make([]Article, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	defer func() {
		if err = db.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	coll := db.Database("test").Collection("articles")
	opts := options.Find().SetSort(bson.D{{"creation_date", -1}})
	cur, err := coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var article Article
		err := cur.Decode(&article)
		if err != nil {
			panic(err)
		}
		articleList = append(articleList, article)
	}
	if err := cur.Err(); err != nil {
		panic(err)
	}
	return articleList
}

func GetArticlesFromDatabaseKeyword(keyword string) []Article {
	articleList := make([]Article, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	defer func() {
		if err = db.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	coll := db.Database("test").Collection("articles")
	opts := options.Find().SetSort(bson.D{{"creation_date", -1}})
	cur, err := coll.Find(ctx, bson.D{{"keyword", keyword}}, opts)
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var article Article
		err := cur.Decode(&article)
		if err != nil {
			panic(err)
		}
		articleList = append(articleList, article)
	}
	if err := cur.Err(); err != nil {
		panic(err)
	}
	return articleList
}

func StrictContains(s, substr string) bool {
	// Regex is stupid, but one day i will need to implement ...
	// For example i dont want to match "Go" with "Google"
	return strings.Contains(s, substr)
}
