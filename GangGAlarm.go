package main

import (
	"bytes"
	"fmt"
	"github.com/solapi/solapi-go"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getTwitchAcessToken(clientId string, clientSecret string) string {
	url := "https://id.twitch.tv/oauth2/token?client_id=" + clientId + "&client_secret=" + clientSecret + "&grant_type=client_credentials"
	reqBody := bytes.NewBufferString("Post")
	resp, err := http.Post(url, "", reqBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(data)[17:47]
}

func GetStreamerLiveB(clientId string, twitchAccessTocken string, streamerId string) bool {
	url := "https://api.twitch.tv/helix/search/channels?query=" + streamerId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("client-id", clientId)
	req.Header.Add("Authorization", "Bearer "+twitchAccessTocken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes)

	if strings.Contains(str, "is_live") {
		pos := strings.Index(str, "is_live") + 9
		str = str[pos : pos+5]

		if strings.Contains(str, "true") {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("[ GangG Alarm System ]")

	client := solapi.NewClient()
	fmt.Println("> Initialize solapi done.")
	twitchToken := getTwitchAcessToken("***Twitch console client id***", "***Twitch console client secret key***")
	fmt.Println("> Generate twitch access token done.")

	fmt.Println("> Start watching...")

	var bLive bool = false
	var bSwitch bool = false
	for {
		bLive = GetStreamerLiveB("***Twitch console client id***", twitchToken, "rkdwl12")

		if bLive {
			if !bSwitch {
				bSwitch = true

				fmt.Println(" > GangG stream on!")

				message := make(map[string]interface{})
				message["to"] = "***발신 전화번호***"
				message["from"] = "***수신 전화번호***"
				message["text"] = "강지 방송이 시작대떠!"
				message["type"] = "SMS"

				params := make(map[string]interface{})
				params["message"] = message

				result, err := client.Messages.SendSimpleMessage(params)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf(" > solapi result: %+v\n", result)
			} else {
				fmt.Println(" > GangG stream ongoing...")
			}
		} else {
			fmt.Println(" > GangG stream off")
			if bSwitch {
				bSwitch = false
			}
		}

		time.Sleep(30 * time.Second)
	}
}
