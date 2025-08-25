package lark

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type Report struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Content struct {
	Post PostContent `json:"post"`
}

type PostContent struct {
	EnUs MessageWrapper `json:"en_us"`
}

type MessageWrapper struct {
	Title   string             `json:"title"`
	Content [][]MessageContent `json:"content"`
}
type MessageContent struct {
	Tag    string `json:"tag"`
	Text   string `json:"text"`
	Href   string `json:"href"`
	UserId string `json:"user_id"`
}

func InitReport(messageType string, title string) Report {
	return Report{
		MsgType: messageType,
		Content: Content{
			Post: PostContent{
				EnUs: MessageWrapper{
					Title: title,
				},
			},
		},
	}
}

func (r *Report) AddMessage(msgs []MessageContent) {
	r.Content.Post.EnUs.Content = append(r.Content.Post.EnUs.Content, msgs)
}

func (r *Report) Send(channel string) {
	url := "https://open.larksuite.com/open-apis/bot/v2/hook/" + channel
	payload := strings.NewReader(r.ToJson())
	log.Println("========= lark payload =========\n" + r.ToJson())
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, _ := client.Do(req)

	b, _ := io.ReadAll(resp.Body)
	output := &bytes.Buffer{}
	_ = json.Compact(output, b)
	log.Println("========= lark response =========\n" + output.String())
	_ = resp.Body.Close()
}

func (r *Report) ToJson() string {
	// Marshal the struct into JSON
	jsonData, _ := json.Marshal(r)

	// Print the JSON data
	return string(jsonData)
}

func (r *Report) Dump() {
	prettyJSON, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(prettyJSON))
}
