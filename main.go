package main

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/fjw95/web-crawler/util"
)

const (
	main_url = "http://ub.ac.id/akademik/fakultas"
)

var countChildUrl, countRootUrl int
var mu sync.Mutex

func fetchSite(url string, wg *sync.WaitGroup, slots chan bool) {
	defer wg.Done()

	respBody, _ := util.GetRespBody(url)

	pattern := regexp.MustCompile(url + "[a-z]+")
	bodyStr := string(respBody[:])
	found := pattern.FindAllString(bodyStr, -1)

	urls := util.RemoveDuplicates(found)
	urls = append(urls, url)

	for _, urlList := range urls {
		mu.Lock()
		fmt.Println("--> " + urlList)
		countChildUrl++
		mu.Unlock()
	}

	fmt.Println("\n---> Found :", len(urls), "URL \nFrom : "+url+"\n")

	<-time.After(3 * time.Second)

	// put back slot that occupied
	slots <- true
}

func getSiteURL(mainURL string, max int) {

	var wg sync.WaitGroup

	// define "max" concurrent slots
	concurrentGoroutines := make(chan bool, max)
	defer close(concurrentGoroutines)

	// fill initial slots
	for i := 0; i < max; i++ {
		concurrentGoroutines <- true
	}

	respBody, _ := util.GetRespBody(main_url)
	bodyStr := string(respBody[:])
	var pattern = regexp.MustCompile("http://" + "[a-z]+" + ".ub.ac.id/en")
	var regexRep = regexp.MustCompile("en")
	var urlStr = pattern.FindAllString(bodyStr, -1)

	// sync with wait group
	countRootUrl = len(urlStr)
	wg.Add(countRootUrl)

	for i, linkURL := range urlStr {
		strRep := regexRep.ReplaceAllString(linkURL, "")
		linkURL := strRep

		// wait available slots
		<-concurrentGoroutines

		fmt.Printf("%d. Launch URL : %s \n", i+1, linkURL)
		go fetchSite(linkURL, &wg, concurrentGoroutines)
	}

	wg.Wait()
}

func main() {

	maxGoroutinesSpawn := 3
	getSiteURL(main_url, maxGoroutinesSpawn)

	fmt.Println("Found", countChildUrl, "unique urls:\n")
	fmt.Println("From", countRootUrl, "Root url:\n")
}
