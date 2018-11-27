package main

import (
	"fmt"
	"log"

	"gopkg.in/headzoo/surf.v1"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Create a new browser and open reddit.
	bow := surf.NewBrowser()
	err := bow.Open("http://www.bcpa.net/RecAddr.asp")
	if err != nil {
		panic(err)
	}

	// Log in to the site.
	fm, _ := bow.Form("[name='homeind']")
	fm.Input("Situs_Street_Number", "515")
	fm.SelectByOptionValue("Situs_Street_Direction", "SW")
	fm.Input("Situs_Street_Name", "18")
	fm.SelectByOptionValue("Situs_Street_Type", "AVE")
	fm.Input("Situs_Street_Post_Dir", "")
	fm.Input("Situs_Unit_Number", "15")
	fm.SelectByOptionValue("Situs_City", "FL")

	if fm.Submit() != nil {
		panic(err)
	}

	//fmt.Println(bow.Body())

	fmt.Println(bow.Url())

	// Load the HTML document

	doc, err := goquery.NewDocument(bow.Url().String())
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf(doc.Html())

	var pageTitle string

	// use CSS selector found with the browser inspector
	// for each, use index and item
	pageTitle = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(1) > td:nth-child(2) > span > a > b").Contents().Text()

	fmt.Printf("Page Title: '%s'\n", pageTitle)

}
