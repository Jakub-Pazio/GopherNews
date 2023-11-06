package articles

import (
	"context"
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

// Returns array of articles, article might be empty
func GetArticles(number int, languages []string) []Article {
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
		for i, s := range article_url {
			article_url[i] = strings.TrimSpace(s)
		}
		if len(article_url) != 2 {
			continue
		}
		for _, lang := range languages {
			if StrictContains(article_url[0], lang) {
				articles = append(articles, NewArticle(article_url[1], article_url[0], lang))
			}
		}
	}
	println(len(articles))
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
	cur, err := coll.Find(ctx, bson.D{})
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
