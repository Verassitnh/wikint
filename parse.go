package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

const fbJsonPath = "require[0][3][0].__bbox.require[11][3][1].__bbox.result.data.node.all_collections.nodes[0].style_renderer.collection.pageItems.edges"
const fbJsonNamePath = "node.title.text"
const fbJsonUrlPath = "node.url"

func parse(usrsCh chan []User, scrapedCh chan ScrapedData, errCh chan error) {
	for sd := range scrapedCh {
		if platformFacebook(sd.url) {
			handleFbdata(usrsCh, sd.body, errCh)
		}
	}
}

func platformFacebook(url string) bool {
	return strings.Contains(url, "facebook.com")
}

func handleFbdata(usrsCh chan []User, body io.Reader, errCh chan error) {
	// 1. Read the data into a document
	// 2. find the scripts with "application/json"
	// 3. parse the json until you find stuff about friends
	// 4. loop through array making sure every item matches form
	// 5. collect the users and send the slice back through the users channel
	var users []User

	// Read the data into a doc
	var doc *goquery.Document
	var err error
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		errCh <- err
		return
	}

	// Find the scripts with "application/json"
	doc.Find("script").Map(func(i int, s *goquery.Selection) string {
		fmt.Printf("parsing %v script block / %v\n", i, s.Length())

		jd := s.Text()

		if !json.Valid([]byte(jd)) {
			errCh <- fmt.Errorf("invalid json, moving on")
			return "" // if json isn't valid, move on
		}

		res := gjson.Get(jd, fbJsonPath)

		if !res.IsArray() {
			errCh <- fmt.Errorf("didn't find an array of friends at '%s', moving on \n json dump: %v", fbJsonPath, jd)
			return "" // if there are no friends, or the objet is wrong, move on
		}

		// Iterate through the array of json results
		for _, v := range res.Array() {
			fmt.Println("found valid json array: parsing friends")
			users = append(users, User{
				name: v.Get(fbJsonNamePath).String(),
				id:   generateId(),
				urls: append([]string{}, v.Get(fbJsonUrlPath).String()),
			})
		}
		return ""
	})

	if len(users) == 0 {
		errCh <- fmt.Errorf("no friends found for user, moving on")
	}

	usrsCh <- users

}
