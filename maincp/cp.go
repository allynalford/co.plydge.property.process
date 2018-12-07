package maincp

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var res, html string
	err = c.Run(ctxt, submit(`https://officialrecords.broward.org/AcclaimWeb/search/SearchTypeParcel`, &res, &html))
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	p := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(p)

	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		tr.Find("td").Each(func(int int, s *goquery.Selection) {
			fmt.Println(s.Contents().Text())
		})

	})

	fmt.Println(doc.Text()) // Links:FooBarBazTEXT I WANT

	log.Printf("HTML: `%s`", html)
	log.Printf("First Value: `%s`", res)
}

func submit(urlstr string, res *string, html *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(`/html/body/div[2]/div/div/div/div[2]/form`),
		chromedp.Submit(`#btnButton`),
		chromedp.WaitVisible(`//input[@name="ParcelId"]`),
		chromedp.SendKeys(`//input[@name="ParcelId"]`, `504203310010`),
		chromedp.Click(`//*[@id="btnSearch"]`),
		chromedp.WaitReady(`//*[@id="0"]`),
		chromedp.OuterHTML(`//*[@id="RsltsGrid"]/div[4]/table`, html),
		chromedp.Text(`//*[@id="RsltsGrid"]/div[4]/table/tbody/tr[1]/td[3]`, res),
	}
}
