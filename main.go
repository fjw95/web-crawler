package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/fjw95/web-crawler/util"
)

var (
	countChildUrl      int
	countRootUrl       int
	maxGoroutinesSpawn int
	mu                 sync.Mutex
	rootUrl            string
	content_file       []string
)

const (
	path_file = "./work/src/github.com/fjw95/web-crawler/file/save/file-save.txt"
)

func fetchSite(url string, wg *sync.WaitGroup, slots chan bool) {
	defer wg.Done()

	respBody, _ := util.GetRespBody(url)

	pattern := regexp.MustCompile(url + "[a-z]+")
	bodyStr := string(respBody[:])
	found := pattern.FindAllString(bodyStr, -1)

	urls := util.RemoveDuplicates(found)
	urls = append(urls, url)

	mu.Lock()
	countChildUrlStr := strconv.Itoa(len(urls))
	content_file = append(content_file, "Found "+countChildUrlStr+" URL, From "+url)
	for _, urlList := range urls {

		content_file = append(content_file, urlList)
		countChildUrl++
	}
	content_file = append(content_file, "")
	mu.Unlock()

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

	respBody, _ := util.GetRespBody(mainURL)
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

	util.WriteFile(content_file, path_file)
	fmt.Println("\nFound", countChildUrl, "unique urls\n")
	fmt.Println("From", countRootUrl, "Root url\n")
}

func main() {

	flag.StringVar(&rootUrl, "url", "", "root url")
	flag.IntVar(&maxGoroutinesSpawn, "max", 1, "max goroutines")
	flag.Parse()

	if rootUrl == "" {
		fmt.Println("Cannot null URL Parameter")
		os.Exit(-1)
	} else {
		getSiteURL(rootUrl, maxGoroutinesSpawn)
	}
}
