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
	"github.com/PuerkitoBio/goquery"
	"net/url"		// to fix our URLs
	"github.com/jackdanger/collectlinks" // return a slice of all the href links found.
	"net/http" 		// using to retrieve a page
	"fmt"
	"flag" 			// 'flag' helps parse command line arguments
	"os" 			// 'os' gives access to system calls
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
	fmt.Println(args)
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

	var docTitle, docAuthor, docCreateDate = "", "", ""
	// Find the review items
	document.Find("span").Each(func (index  int, element *goquery.Selection) {
		// For each item found, get the band and title
		// See if the href attribute exists on the element
		value1 := element.Find(".m_Title").First().Text()
		if value1 != "" {
			docTitle = value1
			fmt.Println("value1: " + value1);
		}
	})
	document.Find("span").Each(func (index  int, element *goquery.Selection) {
		value2 := element.Find(".m_ReferenceSourceTG").First().Text()
		if value2 != "" {
			docAuthor = value2
			fmt.Println("value2: " + value2);
		}
	})
	document.Find("span").Each(func (index  int, element *goquery.Selection) {
		value3 := element.Find(".m_DateCreated").First().Text()
		if value3 != "" {
			docCreateDate = value3
			fmt.Println("value3: " + value3);
		}
	})
	fmt.Printf("Review %s: %s, %s, %s\n", url, docTitle, docAuthor, docCreateDate)

	for _, link := range(links) {
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