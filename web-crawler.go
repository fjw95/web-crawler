package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func getSiteURL(mainURL string) <-chan string {

	client := &http.Client{}
	res, err := client.Get(mainURL)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	bodyStr := string(body[:])
	var regex = regexp.MustCompile("http://" + "[a-z]+" + ".ub.ac.id/en")
	var urlStr = regex.FindAllString(bodyStr, -1)

	channel := make(chan string)
	go func() {
		defer close(channel)
		for _, linkURL := range urlStr {
			channel <- linkURL
		}
	}()
	return channel

}

func main() {
	//channel := make(chan bool)

	//defer close(channel)

	seedUrls := os.Args[1:]

	for _, url := range seedUrls {
		getSiteURL(url)
	}

	//<-channel
}
