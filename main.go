package main

import (
	"fmt"
	"github.com/fjw95/web-crawler/util"
	"regexp"
	"sync"
	"time"
)

const (
	main_url = "http://ub.ac.id/akademik/fakultas"
)

var countChildUrl, countRootUrl int
var mu sync.Mutex

func fetchSite(url string, wg *sync.WaitGroup, done chan bool) {
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
	done <- true
}

func getSiteURL(mainURL string, wg *sync.WaitGroup, max int) {

	concurrentGoroutines := make(chan []string, max)
	done := make(chan bool)
	waitForAllJobs := make(chan bool)

	defer close(concurrentGoroutines)

	for i := 0; i < max; i++ {
		concurrentGoroutines <- []string{}
	}

	respBody, _ := util.GetRespBody(main_url)

	bodyStr := string(respBody[:])
	var pattern = regexp.MustCompile("http://" + "[a-z]+" + ".ub.ac.id/en")
	var regexRep = regexp.MustCompile("en")
	var urlStr = pattern.FindAllString(bodyStr, -1)

	countRootUrl = len(urlStr)

	wg.Add(countRootUrl)

	go func() {
		for i := 0; i < countRootUrl; i++ {
			<-done

			concurrentGoroutines <- []string{}
		}
		waitForAllJobs <- true
	}()

	for _, linkURL := range urlStr {
		strRep := regexRep.ReplaceAllString(linkURL, "")
		linkURL := strRep

		<-concurrentGoroutines
		fmt.Println("Launch URL : " + linkURL + "\n")
		go fetchSite(linkURL, wg, done)
	}
	<-waitForAllJobs

}

func main() {

	maxGoroutinesSpawn := 3
	var wg sync.WaitGroup
	getSiteURL(main_url, &wg, maxGoroutinesSpawn)

	wg.Wait()
	fmt.Println("Found", countChildUrl, "unique urls:\n")
	fmt.Println("From", countRootUrl, "Root url:\n")
}
