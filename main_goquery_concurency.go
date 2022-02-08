package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var userAgents = []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/37.0.2062.94 Chrome/37.0.2062.94 Safari/537.36", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/37.0.2062.94 Chrome/37.0.2062.94 Safari/537.36"}

func randomUserAgent() {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents(randNum)
}

func discoverLinks(response *http.Request, baseURL string) []string {
	if response != nil {
		doc, _ := goquery.NewDocumentFromResponse(response)
		founUrls := []string{}

		if doc != nil {
			doc.Find("div.s-result-item ").Each(func(index int, item *goquery.Selection) {

				//interasi ke dalam link tiap item
				mainUrl, _ := item.Find("a.a-link-normal").Attr("href")

				detailUrl := "https://www.amazon.com" + mainUrl

				detailResponse, err := http.Get(detailUrl)
				check(err)

				if detailResponse.StatusCode > 400 {
					fmt.Println("Status Code: ", response.StatusCode)
				}

				detailDoc, err := goquery.NewDocumentFromReader(detailResponse.Body)
				check(err)

				//ambil detail data tiap item
				title := detailDoc.Find("span.a-size-large").Text()

				asin, _ := detailDoc.Find("div#averageCustomerReviews").Attr("data-asin")

				Brand = detailDoc.Find("tr:contains('Brand')").Find("td").First().Next().Text()

				ratings := detailDoc.Find("span#acrCustomerReviewText.a-size-base").Text()

				Price := detailDoc.Find("tr:contains('Price:')").Find("td span.a-price span.a-offscreen").Text()

				Sales := "belum ada"

				SalesGraph := "belum ada"

				fmt.Println("Product Details: "+title, "asin: "+asin, "brand: "+Brand, "Price: "+Price, "Ratings: "+ratings, "Sales: "+Sales,
					"SalesGraph:"+SalesGraph)

			})

		}

		return foundUrls
	} else {
		return []string{}
	}
}

func getRequest(targetURL string) (*http.Response, error) {
	client := &http.Client()
	request, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}

	request.Header.Set("User-Agent", randomUserAgent())

	ress, err := client.Do(request)

	if err != nil {
		return nil, err
	} else {
		return res, nil
	}

}

func resolveRelativeLinks() {}

func Crawl(targetURL string, baseURL string) []string {
	fmt.Println(targetURL)
	response, _ := getRequest(targetURL)

	links := discoverLinks(response, baseURL)
	foundUrls := []string{}

	for _, link := range links {
		ok, correctLinks := resolveRelativeLinks(link, baseURL)
		if ok {
			if correctLinks != "" {
				foundUrls = append(foundUrls, correctLinks)
			}
		}
	}

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
	go func() {
		worklist <- []string{"https://www.amazon.com/s?k=keyboard&crid=ICSAKQ8GAW3D&sprefix=keybo%2Caps%2C592&ref=nb_sb_noss_2"}
	}()

	seen := make(map[string]bool)

	for ; n > 0; n-- {
		list := worklist

		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string, baseUrl string) {
					foundLinks := Crawl(link, url)
					if foundLinks != nil {
						worklist <- foundLinks
					}
				}()
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
