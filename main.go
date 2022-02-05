package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

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
	url := "https://www.amazon.com/s?k=keyboard&crid=ICSAKQ8GAW3D&sprefix=keybo%2Caps%2C592&ref=nb_sb_noss_2"
	// url := "https://www.tokopedia.com"
	var Brand string

	// response, err := http.Get(url)
	request, err := http.NewRequest("GET", url, nil)
	check(err)

	request.Header.Set("User-Agent", "Not Firefox")
	// request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	response, err := http.DefaultClient.Do(request)
	check(err)

	defer response.Body.Close()
	// slurp, _ := ioutil.ReadAll(response.Body)
	// fmt.Printf("%v\n%v\n%s\n", response, response.Request, slurp)

	if response.StatusCode > 400 {
		fmt.Println("Status Code: ", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	check(err)

	//untuk cek kalo di block
	// html, err := doc.Find("body").Html()
	// check(err)
	//
	// fmt.Println(html)

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
