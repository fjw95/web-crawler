package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	//"sync"
)

const (
	main_url = "http://ub.ac.id/akademik/fakultas"
)

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func fetchSite(url string, chUrl chan string, chSignal chan bool) {

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	pattern := regexp.MustCompile(url + "[a-z]+")
	bodyStr := string(respBody[:])
	found := pattern.FindAllString(bodyStr, -1)

	urls := removeDuplicates(found)
	urls = append(urls, url)

	for _, urlList := range urls {
		// fmt.Println(urlList)
		chUrl <- urlList
	}

	chSignal <- true

}

func main() {

	// Deklarasi variabel waitGroup
	// var wg sync.WaitGroup
	// getSiteURL(main_url, &wg)
	foundUrls := make(map[string]bool)
	channelUrl := make(chan string)
	channelSignal := make(chan bool)

	client := &http.Client{}
	res, err := client.Get(main_url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	bodyStr := string(body[:])
	var pattern = regexp.MustCompile("http://" + "[a-z]+" + ".ub.ac.id/en")
	var urlStr = pattern.FindAllString(bodyStr, -1)

	for _, linkURL := range urlStr {

		var regexRep = regexp.MustCompile("en")
		var strRep = regexRep.ReplaceAllString(linkURL, "")
		linkURL := strRep

		go fetchSite(linkURL, channelUrl, channelSignal)

	}

	for i := 0; i < len(urlStr); {
		select {
		case url := <-channelUrl:
			foundUrls[url] = true
		case <-channelSignal:
			i++
		}
	}

	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")
	close(channelUrl)
	// wg.Wait()
}
