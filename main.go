package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type FoundUrls struct {
	ProductDetails string
	Asin           string
	Brand          string
	Price          string
	Ratings        string
	Sales          string
	SalesGraph     string
}

var userAgents = []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/37.0.2062.94 Chrome/37.0.2062.94 Safari/537.36", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/37.0.2062.94 Chrome/37.0.2062.94 Safari/537.36"}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func discoverLinks(response *http.Response, baseURL string) []string {
	if response != nil {
		doc, _ := goquery.NewDocumentFromResponse(response)
		foundUrls := []string{}

		if doc != nil {
			doc.Find("div.s-result-item ").Each(func(index int, item *goquery.Selection) {

				//interasi ke dalam link tiap item
				mainUrl, _ := item.Find("a.a-link-normal").Attr("href")

				foundUrls = append(foundUrls, mainUrl)
				// detailUrl := "https://www.amazon.com" + mainUrl
				// //
				// detailResponse, err := http.Get(detailUrl)
				// check(err)
                //
				// if detailResponse.StatusCode > 400 {
				// 	fmt.Println("Status Code: ", response.StatusCode)
				// }
                //
				// detailDoc, err := goquery.NewDocumentFromReader(detailResponse.Body)
				// check(err)
                //
				// // ambil detail data tiap item
				// title := detailDoc.Find("span.a-size-large").Text()
                //
				// foundUrls = append(foundUrls, title)
				// asin, _ := detailDoc.Find("div#averageCustomerReviews").Attr("data-asin")
				//
				// Brand := detailDoc.Find("tr:contains('Brand')").Find("td").First().Next().Text()
				//
				// ratings := detailDoc.Find("span#acrCustomerReviewText.a-size-base").Text()
				//
				// Price := detailDoc.Find("tr:contains('Price:')").Find("td span.a-price span.a-offscreen").Text()
				//
				// Sales := "belum ada"

				// SalesGraph := "belum ada"
				//
				// foundUrls = &FoundUrls{
				// 	ProductDetails: title,
				// 	Asin:           asin,
				// 	Brand:          Brand,
				// 	Price:          Price,
				// 	Ratings:        ratings,
				// 	Sales:          Sales,
				// 	SalesGraph:     SalesGraph,
				// }
				// fmt.Println("Product Details: "+title, "asin: "+asin, "brand: "+Brand, "Price: "+Price, "Ratings: "+ratings, "Sales: "+Sales,
				// "SalesGraph:"+SalesGraph)

			})

		}

		return foundUrls

	} else {
		return []string{}
	}
}

func getRequest(targetURL string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", randomUserAgent())

	res, err := client.Do(request)

	if err != nil {
		return nil, err
	} else {
		return res, nil
	}

}

func checkRelative(href string, baseUrl string) string {
	if strings.HasPrefix(href, "/") {
		return fmt.Sprintf("%s%s", baseUrl, href)
	} else {
		return href
	}
}

func resolveRelativeLinks(href string, baseUrl string) (bool, string) {
	resultHref := checkRelative(href, baseUrl)
	baseParse, _ := url.Parse(baseUrl)
	resultParse, _ := url.Parse(resultHref)

	if baseParse != nil && resultParse != nil {
		if baseParse.Host == resultParse.Host {
			return true, resultHref
		} else {
			return false, ""
		}
	}

	return false, ""
}

var tokens = make(chan struct{}, 5)

func Crawl(targetURL string, baseURL string) []string {
	fmt.Println(targetURL)

	tokens <- struct{}{}
	response, _ := getRequest(targetURL)
	<-tokens

	links := discoverLinks(response, baseURL)
	foundUrls := []string{}

	fmt.Println(links)
	// fmt.Println(baseURL)

	// for _, link := range links {
	// 	ok, correctLinks := resolveRelativeLinks(link, baseURL)
	// 	if ok {
	// 		if correctLinks != "" {
	// 			foundUrls = append(foundUrls, correctLinks)
	// 		}
	// 	}
    //
	// }

	return foundUrls
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}

}

func writeFile(data, filename string) {
	file, err := os.Create(filename)
	defer file.Close()
	check(err)

	file.WriteString(data)
}

func main() {

	worklist := make(chan []string)
	var n int
	n++

	url := "https://www.amazon.com/s?k=keyboard&crid=ICSAKQ8GAW3D&sprefix=keybo%2Caps%2C592&ref=nb_sb_noss_2"
	// url := "https://www.guardian.com"
	go func() {
		worklist <- []string{"https://www.amazon.com/s?k=keyboard&crid=ICSAKQ8GAW3D&sprefix=keybo%2Caps%2C592&ref=nb_sb_noss_2"}
		// worklist <- []string{"https://www.guardian.com"}
	}()

	seen := make(map[string]bool)

	for ; n > 0; n-- {
		list := <-worklist

		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string, baseUrl string) {
					foundLinks := Crawl(link, url)
					if foundLinks != nil {
						worklist <- foundLinks
					}
				}(link, url)
			}

		}

	}

	// var Brand string
	//
	// response, err := http.Get(url)
	//
	// check(err)
	//
	// defer response.Body.Close()
	//
	// if response.StatusCode > 400 {
	// 	fmt.Println("Status Code: ", response.StatusCode)
	// }
	//
	// doc, err := goquery.NewDocumentFromReader(response.Body)
	// check(err)
	//
	// doc.Find("div.s-result-item ").Each(func(index int, item *goquery.Selection) {
	//
	// 	//interasi ke dalam link tiap item
	// 	mainUrl, _ := item.Find("a.a-link-normal").Attr("href")
	//
	// 	detailUrl := "https://www.amazon.com" + mainUrl
	//
	// 	detailResponse, err := http.Get(detailUrl)
	// 	check(err)
	//
	// 	if detailResponse.StatusCode > 400 {
	// 		fmt.Println("Status Code: ", response.StatusCode)
	// 	}
	//
	// 	detailDoc, err := goquery.NewDocumentFromReader(detailResponse.Body)
	// 	check(err)
	//
	// 	//ambil detail data tiap item
	// 	title := detailDoc.Find("span.a-size-large").Text()
	//
	// 	asin, _ := detailDoc.Find("div#averageCustomerReviews").Attr("data-asin")
	//
	// 	Brand = detailDoc.Find("tr:contains('Brand')").Find("td").First().Next().Text()
	//
	// 	ratings := detailDoc.Find("span#acrCustomerReviewText.a-size-base").Text()
	//
	// 	Price := detailDoc.Find("tr:contains('Price:')").Find("td span.a-price span.a-offscreen").Text()
	//
	// 	Sales := "belum ada"
	//
	// 	SalesGraph := "belum ada"
	//
	// 	fmt.Println("Product Details: "+title, "asin: "+asin, "brand: "+Brand, "Price: "+Price, "Ratings: "+ratings, "Sales: "+Sales,
	// 		"SalesGraph:"+SalesGraph)
	//
	// })

}
