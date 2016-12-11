package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	slackAPIURL = "https://slack.com/api/chat.postMessage"
	err         error
)

func validateENVValue() error {
	if os.Getenv("PORT") == "" {
		log.Println("$PORT must be set")
		return errors.New("error!")
	}
	if os.Getenv("LINE_CHANNEL_SECRET") == "" {
		log.Println("$LINE_CHANNEL_SECRET must be set")
		return errors.New("error!")
	}
	if os.Getenv("LINE_CHANNEL_ACCESS_TOKEN") == "" {
		log.Println("$LINE_CHANNEL_ACCESS_TOKEN must be set")
		return errors.New("error!")
	}
	if os.Getenv("SLACK_TOKEN") == "" {
		log.Println("$SLACK_TOKEN must be set")
		return errors.New("error!")
	}
	if os.Getenv("SLACK_CHANNEL") == "" {
		log.Println("$SLACK_CHANNEL must be set")
		return errors.New("error!")
	}

	return nil

}

func makeSlackData(text string) url.Values {
	// make slack post data

	slackData := url.Values{}
	slackData.Set("token", os.Getenv("SLACK_TOKEN"))
	slackData.Add("channel", os.Getenv("SLACK_CHANNEL"))
	slackData.Add("username", "line2slack piper")
	slackData.Add("text", text)

	return slackData

}

func main() {
	port := os.Getenv("PORT")
	line_channel_secret := os.Getenv("LINE_CHANNEL_SECRET")
	line_channel_access_token := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	err = validateENVValue()
	if err != nil {
		log.Fatal("all ENV must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.POST("/hook", func(c *gin.Context) {
		client := &http.Client{}
		bot, err := linebot.New(line_channel_secret, line_channel_access_token, linebot.WithHTTPClient(client))
		if err != nil {
			fmt.Println(err)
			return
		}
		received, err := bot.ParseRequest(c.Request)

		for _, event := range received { // wait new post
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage: // if new post type is text

					slackPostData := makeSlackData(message.Text)

					// post to slack
					r, _ := http.NewRequest("POST", fmt.Sprintf("%s", slackAPIURL), bytes.NewBufferString(slackPostData.Encode()))
					r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

					client.Do(r)
				}
			}
		}
	})

	router.Run(":" + port)
}
