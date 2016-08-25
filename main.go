package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
)

const (
	main_url = "http://ub.ac.id/akademik/fakultas"
)

var count int
var mu sync.Mutex

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] != true {
			/// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func fetchSite(url string, wg *sync.WaitGroup) {

	defer wg.Done()
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
		mu.Lock()
		fmt.Println(urlList)
		count++
		mu.Unlock()
	}
}

func getSiteURL(mainURL string, wg *sync.WaitGroup) {

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
	var pattern = regexp.MustCompile("http://" + "[a-z]+" + ".ub.ac.id/en")
	var urlStr = pattern.FindAllString(bodyStr, -1)
	var regexRep = regexp.MustCompile("en")

	for _, linkURL := range urlStr {
		wg.Add(1)

		strRep := regexRep.ReplaceAllString(linkURL, "")
		linkURL := strRep

		go fetchSite(linkURL, wg)

	}

}

func main() {

	// Deklarasi variabel waitGroup
	var wg sync.WaitGroup
	getSiteURL(main_url, &wg)

	wg.Wait()
	fmt.Println("\nFound", count, "unique urls:\n")
}
