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
	"net/url"		// to fix our URLs
	"crypto/tls"	// to get access to some low-level transport customizations
	//"github.com/jackdanger/collectlinks" // return to you a slice of all the href links found.
	"net/http" 		// using to retrieve a page
	"io/ioutil" 	// to maek reading and printing the html page
	"fmt"
	"flag" 			// 'flag' helps you parse command line arguments
	"os" 			// 'os' gives you access to system calls
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
	for uri := range queue {
		// pass each URL we find off to be read & enqueued
		enqueue(uri, queue)
	}
}

func enqueue(uri string, queue chan string) {
	fmt.Println("fetching", uri)
	visited[uri] = true // Record that we're going to visit this page

	// ~ new(thing(a: b))
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	// also provides a way to override defaults (like 'http.Get')
	client := http.Client{Transport: transport}

	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	//links := collectlinks.All(resp.Body)

	// SHOW REQUIMENT by NORDIC-CODER
	body, _ := ioutil.ReadAll(resp.Body)
	getLitteContentPage(string(body))

	// for _, link := range(links) {
	// 	// Don't enqueue the raw thing we find
	// 	// set invalid URLs to blank strings
	// 	absolute := fixUrl(link, uri)

	// 	// so let's never send those to the channel
	// 	if uri != "" {
	// 		// Don't enqueue a page twice!
	// 		if !visited[absolute] {
	// 			// We asynchronously enqueue what we've found
	// 			go func() {
	// 				queue <- link
	// 			}()
	// 		}
	// 	}
	// }
}

func getLitteContentPage(bodyStr string) {
	fmt.Println(bodyStr)
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

func retrieve(uri string) {
	// ex: https://6brand.com/
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	// need to close the resource we opened
	// `defer` delays an operation until the function ends.
	defer resp.Body.Close()

	// resp.Body isn't a string, it's more like a reference
	// to a stream of data. So we use the 'ioutil'
	// package to read it into memory.
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("read error is: ", err)

	// cast the html body to a string because
	// Go hands it to us as a byte array
	fmt.Println(string(body))

}