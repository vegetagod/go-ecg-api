package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"go-ecg/models"
	"go-ecg/utils"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func main() {

	//Pipe LINE User ID:U57802f0ee059e4b5fefe07e525192376
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	bot, err := messaging_api.NewMessagingApiAPI(
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		log.Println("/callback called...")

		cb, err := webhook.ParseRequest(channelSecret, req)
		if err != nil {
			log.Printf("Cannot parse request: %+v\n", err)
			if errors.Is(err, webhook.ErrInvalidSignature) {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		log.Println("Handling events...")

		for _, event := range cb.Events {
			log.Printf("/callback called%+v...\n", event)

			switch e := event.(type) {
			case webhook.FollowEvent:
				feb, err := e.MarshalJSON()

				if err != nil {
					log.Println("error marshall json in follow event: ", err)
					return
				}
				log.Println(string(feb))
				var evm models.FollowEvent
				if err := json.Unmarshal(feb, &evm); err != nil {
					log.Println("error unmarshall json to struct in follow event: ", err)
					return
				}

				log.Println("user id", evm.Source.UserID)

				requestUrl := fmt.Sprintf("https://api.line.me/v2/bot/profile/%s", evm.Source.UserID)

				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("LINE_CHANNEL_TOKEN")))

				if err != nil {
					log.Println("client: could not create request: %s\n", err)
					return
				}
				client := utils.GetHttpClient()
				resp, err := client.Do(req)
				if err != nil {
					log.Println("error cannot get user profile", err)
					return
				}
				defer resp.Body.Close()

				//if resp.StatusCode == http.StatusOK {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				bodyString := string(bodyBytes)
				log.Println(bodyString)
				//}

				if _, err = bot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text: "ขอบคุณที่รับแอดเราเป็นเพื่อน",
							},
						},
					},
				); err != nil {
					log.Print(err)
				} else {
					log.Println("Sent text reply.")
				}
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					if _, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: message.Text,
								},
							},
						},
					); err != nil {
						log.Print(err)
					} else {
						log.Println("Sent text reply.")
					}
				case webhook.StickerMessageContent:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerId, message.StickerResourceType)
					if _, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: replyMessage,
								},
							},
						}); err != nil {
						log.Print(err)
					} else {
						log.Println("Sent sticker reply.")
					}

				default:
					log.Printf("Unsupported message content: %T\n", e.Message)
				}
			default:
				log.Printf("Unsupported message: %T\n", event)
			}
		}
	})

	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Println("http://localhost:" + port + "/")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
