package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/christer79/ccTray2Slack/cctray"
)

type Result struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

type Client struct {
	URL      string
	Username string
	Password string
}

func NewClient(url, username, password string) Client {
	return Client{url, username, password}
}

func (c *Client) Create(p cctray.Project, jiraProject string) {
	log.Println("post 2")
	var jsonStr = []byte(`{
    "fields": {
       "project":
       {
          "key": "PROJECT"
       },
       "summary": "SUMMARY",
       "description": "DESCTIPTION",
       "issuetype": {
          "name": "Bug"
       }
    }
  }`)
	log.Println("post 3")
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	log.Println("post 4")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)
	log.Println("post 5")
	client := &http.Client{}
	log.Println("post 6")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	log.Println("post 7")
	defer resp.Body.Close()
	var result Result
	log.Printf("Status: %v  \n", resp.Status)

	htmlData, err := ioutil.ReadAll(resp.Body) //<--- here!

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// print out
	fmt.Println(os.Stdout, string(htmlData)) //<-- here !
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}
	log.Println(result)
}
