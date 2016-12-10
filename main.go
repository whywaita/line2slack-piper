package main

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"os"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	slack_api_url = "https://slack.com/api/chat.postMessage"
)

func main() {
	port := os.Getenv("PORT")
	line_channel_secret := os.Getenv("LINE_CHANNEL_SECRET")
	line_channel_access_token := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	slack_token := os.Getenv("SLACK_TOKEN")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	if line_channel_secret == "" {
		log.Fatal("$LINE_CHANNEL_SECRET must be set")
	}
	if line_channel_access_token == "" {
		log.Fatal("$LINE_CHANNEL_ACCESS_TOKEN must be set")
	}
	if slack_token == "" {
		log.Fatal("$SLACK_TOKEN must be set")
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

		for _, event := range received {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:

					// make slack data
					slackPostData := url.Values{}
					slackPostData.Set("token", slack_token)
					slackPostData.Add("channel", os.Getenv("SLACK_CHANNEL"))
					slackPostData.Add("username", "line2slack piper")
					slackPostData.Add("text", message.Text)

					r, _ := http.NewRequest("POST", fmt.Sprintf("%s", slack_api_url), bytes.NewBufferString(slackPostData.Encode()))
					r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

					client.Do(r)
					//resMessage := getResMessage(message.Text)
					//if resMessage != "" {
					//      postMessage := linebot.NewTextMessage(resMessage)
					//      _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do()
					//      if err != nil {
					//              log.Print(err)
					//      }
					//}
				}
			}
		}
	})

	router.Run(":" + port)
}
