package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const BING_API = "https://bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&nc=1612409408851&pid=hp&FORM=BEHPTB&uhd=1&uhdwidth=3840&uhdheight=2160"
const BING_URL = "https://bing.com"
const URL_ATTACH = "&pid=hp&w=384&h=216&rs=1&c=4"
const IMAGES = "images"
const README = "README.md"

var result map[string]interface{}

var (
	filepath string
	random   bool
)

func main() {
	flag.Parse()
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	resp, _ := c.Get(BING_API)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}
	images := result["images"]
	var text string
	var enddate string
	var url string
	var url2 string
	var title string
	var name string

	for _, image := range images.([]interface{}) {
		url = image.(map[string]interface{})["url"].(string)
		url2 = strings.Split(url, "&")[0]
		title = image.(map[string]interface{})["title"].(string)
		enddate = image.(map[string]interface{})["enddate"].(string)
		text = "\n" + "![" + title + "](" + BING_URL + url2 + URL_ATTACH + ")" + "\n"
		name = strings.Split(url2, "=")[1]

		os.MkdirAll(IMAGES+"/"+enddate, os.ModePerm)
		fmt.Println("url: " + url)
		fmt.Println("url2: " + url2)
		URL := BING_URL + url2
		fileName := name
		err = download(URL, IMAGES+"/"+enddate+"/"+fileName)
		if err != nil {
			log.Fatal(err)
		}

		readme(text)
	}

}

func readme(text string) {
	file, err := os.OpenFile(README, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("file open failed", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(text)
	writer.Flush()
}

func download(URL, name string) error {
	fmt.Println(URL)
	response, _ := http.Get(URL)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil
	}
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	flag.StringVar(&filepath, "filepath", "README.md", "write to file")
	flag.BoolVar(&random, "random", false, "if true, random image")
}
