package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const BING_API = "https://bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&nc=1612409408851&pid=hp&FORM=BEHPTB&uhd=1&uhdwidth=3840&uhdheight=2160"
const BING_URL = "https://bing.com"
const URL_ATTACH = "&pid=hp&w=384&h=216&rs=1&c=4"

var result map[string]interface{}

const filePath = "README.md"

func main() {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	resp, err := c.Get(BING_API)
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}
	images := result["images"]
	var text string
	for _, image := range images.([]interface{}) {
		url := image.(map[string]interface{})["url"]
		url2 := strings.Split(url.(string), "&")
		title := image.(map[string]interface{})["title"]
		// enddate := image.(map[string]interface{})["enddate"]
		text = "![" + title.(string) + "](" + BING_URL + url2[0] + ")" + "\n" +
			"<center>" + title.(string) + "</center>" + "\n\n"

	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("file open failed", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(text)
	writer.Flush()
}
