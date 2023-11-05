package articles

import (
    "time"
    "os/exec"
    "strings"
    "fmt"
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

func StrictContains(s, substr string) bool {
	// Regex is stupid, but one day i will need to implement ...
	// For example i dont want to match "Go" with "Google"
	return strings.Contains(s, substr)
}

