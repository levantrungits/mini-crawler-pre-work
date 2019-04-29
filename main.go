/*
	Build a Web Crawler in Go
		Step 1. Starting from a specific page
		Step 2. Retrieving a page from the internet
		Step 3. Parsing hyperlinks from HTML
		Step 4. Concurrency
		Step 5. Data sanitization
		Step 6. Avoiding Loops
*/

package main

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	"net/url"		// to fix our URLs
	"github.com/jackdanger/collectlinks" // return a slice of all the href links found.
	"net/http" 		// using to retrieve a page
	"fmt"
	"flag" 			// 'flag' helps parse command line arguments
	"os" 			// 'os' gives access to system calls
	//"go-module/crawler"
)

// Let’s not fetch any page more than once.
// a global variable
// a map of string -> bool
var visited = make(map[string]bool)

func main()  {
	// converts the command line
	flag.Parse() 	
	// aruments to a new variable named 'args' --> args will be either an empty list			
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please specify start page...")
		os.Exit(1)
	}

	// new channel -> url string data fed into it.
	queue := make(chan string)
	// asynchronously -> put args[0] url into the channel
	go func() {
		queue <- args[0]
	}()
	// an effective iterator keyword
	for url := range queue {
		// pass each URL find off to be read & enqueued
		enqueue(url, queue)
		return //temp
	}
}

func enqueue(url string, queue chan string) {
	fmt.Println("fetching", url)
	visited[url] = true // Record that we're going to visit this page
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Get All Links Exist In Body URL
	links := collectlinks.All(resp.Body)

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var docTitle, docAuthor, docCreateDate = "1", "2", "3"
	// Find the review items
	document.Find("#ctl00_cphContent_lblTitleHtml").Each(func(index int, element *goquery.Selection) {
		// For each item found, get the band and title
		// See if the href attribute exists on the element
		docTitle = element.First().Text()
	})
	document.Find("#ctl00_cphContent_Lbl_Author").Each(func(index int, element *goquery.Selection) {
		docAuthor = element.First().Text()
	})
	document.Find("#ctl00_cphContent_lblCreateDate").Each(func(index int, element *goquery.Selection) {
		docCreateDate = element.First().Text()
	})
	fmt.Printf("Review %s: %s, %s, %s\n", url, docTitle, docAuthor, docCreateDate)

	for _, link := range(links) { // for (idx, valueString)
		if strings.Index(link, ".html") >= 0 {
			if strings.Index(link, "http") < 0 {
				link = "https://www.thesaigontimes.vn" + link
			}

			// fmt.Printf("idx %d  === %s\n", idx, link)

			// Don't enqueue the raw thing we find
			// set invalid URLs to blank strings
			absolute := fixUrl(link, url)

			// so let's never send those to the channel
			if url != "" {
				// Don't enqueue a page twice!
				if !visited[absolute] {
					// We asynchronously enqueue what we've found
					go func() {
						queue <- link
					}()
				}
			}
		}
	}
}

func fixUrl(href, base string) (string) {
	// given a relative link
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}

	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}

	// parsed url objects in this
	uri = baseUrl.ResolveReference(uri)
	// return a plain string.
	return uri.String()
}