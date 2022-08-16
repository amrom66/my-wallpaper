package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

const IMAGES = "images"

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
	var enddate string
	var url2 string
	var title string

	for _, image := range images.([]interface{}) {
		url := image.(map[string]interface{})["url"]
		url2 := strings.Split(url.(string), "&")[0]
		title = image.(map[string]interface{})["title"].(string)
		enddate = image.(map[string]interface{})["enddate"].(string)
		text = "![" + title + "](" + BING_URL + url2 + ")" + "\n" +
			"<center>" + title + "</center>" + "\n\n"
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("file open failed", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(text)
	writer.Flush()
	err = os.MkdirAll(IMAGES+"/"+enddate, os.ModePerm)
	if err != nil {
		fmt.Println("mkdir failed", err)
		return
	}
	filename, err := DownloadImage(BING_URL+url2, strings.Join([]string{IMAGES, enddate}, "/"), title)
	if err != nil {
		fmt.Println("download failed", err)
		return
	}
	fmt.Println("download success", filename)
}

func DownloadImage(imgUrl string, path string, name string) (filename string, err error) {
	filename = path + "/" + name + ".jpg"
	res, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	out, _ := os.Create(filename)
	io.Copy(out, bytes.NewReader(body))

	out.Close()
	return
}
