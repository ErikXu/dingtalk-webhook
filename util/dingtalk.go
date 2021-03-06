package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SendMarkdownMsg(webhook string, secret string, title string, msg string, isAtAll bool, atMobiles []string, atUserIds []string) error {
	markdownRequest := MarkdownRequest{
		Markdown: MarkdownRequestContent{
			Title: title,
			Text:  msg,
		},
		Msgtype: "markdown",
		At: At{
			IsAtAll: isAtAll,
		},
	}

	if len(atMobiles) > 0 {
		markdownRequest.At.AtMobiles = atMobiles
	}

	if len(atUserIds) > 0 {
		markdownRequest.At.AtUserIds = atUserIds
	}

	jsonContent, err := json.Marshal(markdownRequest)
	if err != nil {
		return err
	}

	body := string(jsonContent)

	timestamp := time.Now().UnixNano() / 1000000
	strToSign := fmt.Sprintf("%d\n%s", timestamp, secret)

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(strToSign))
	sum := hash.Sum(nil)
	encode := base64.StdEncoding.EncodeToString(sum)
	urlEncode := url.QueryEscape(encode)

	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", webhook, timestamp, urlEncode)

	request, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}

	client := &http.Client{}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}

type MarkdownRequest struct {
	Markdown MarkdownRequestContent `json:"markdown"`
	Msgtype  string                 `json:"msgtype"`
	At       At                     `json:"at"`
}

type MarkdownRequestContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}
