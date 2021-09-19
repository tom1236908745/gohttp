package main

import (
	"github.com/labstack/echo"
	"github.com/tom1236908745/gohttp/twitter"
	"html/template"
	"io"
	"net/http"
)

func main() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.GET("/tweets", tweets)

	e.Logger.Fatal(e.Start(":8000"))
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func tweets(c echo.Context) error {
	value := c.QueryParam("value")
	api := twitter.ConnectTwitterApi()
	//検索
	searchResult, _ := api.GetSearch(`"` + value + `"`, nil)
	tweets := make([]*TweetTempete, 0)
	for _, data := range searchResult.Statuses {
		tweet := new(TweetTempete)
		tweet.Text = data.FullText
		tweet.User = data.User.Name
		tweet.Id = data.User.IdStr
		tweet.ScreenName = data.User.ScreenName
		tweet.Date = data.CreatedAt
		tweet.TweetId = data.IdStr
		tweets = append(tweets, tweet)
	}

	return c.Render(http.StatusOK, "tweets.html", tweets)
}

type Template struct {
	templates *template.Template
}

//　ツイートの情報
type TweetTempete struct {
	User string `json:"user"`
	Text string `json:"text"`
	ScreenName string `json:"screenName"`
	Id string `json:"id"`
	Date string `json:"date"`
	TweetId string `json:"tweetId"`
}