package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const fbJsonPath = "require[0][3][0].__bbox.require[11][3][1].__bbox.result.data.node.all_collections.nodes[0].style_renderer.collection.pageItems.edges"
const fbJsonNamePath = "node.title.text"
const fbJsonUrlPath = "node.url"

func parse(usrsCh chan []User, scrapedCh chan ScrapedData, errCh chan error) {
	for sd := range scrapedCh {
		if platformFacebook(sd.url) {
			handleFbdata(usrsCh, sd, errCh)
		}
	}
}

func platformFacebook(url string) bool {
	return strings.Contains(url, "facebook.com")
}

func containsAttribute(attr []html.Attribute, key string, value string) bool {
	for _, v := range attr {
		if v.Key == key && v.Val == value {
			return true
		}
	}
	return false
}

func handleFbdata(usrsCh chan []User, sd ScrapedData, errCh chan error) {
	// 1. Read the data into a document
	// 2. find the scripts with "application/json"
	// 3. parse the json until you find stuff about friends
	// 4. loop through array making sure every item matches form
	// 5. collect the users and send the slice back through the users channel
	var users []User

	// Read the data into a doc
	node, err := html.Parse(sd.body)
	if err != nil {
		errCh <- fmt.Errorf("cannot parse html for %v, moving on: %v", sd.url, err)
	}

	htmlRecurser(node, func(n *html.Node) bool {
		// Select script tag with the correct attributes
		if n.DataAtom != atom.Script && !containsAttribute(n.Attr, "type", "application/json") {
			// errCh <- fmt.Errorf("node not script, moving on")
			return false

		}
		if n.FirstChild != nil {
			if !json.Valid([]byte(n.FirstChild.Data)) {
				errCh <- fmt.Errorf("invalid json for user, moving on")
				return false // if json isn't valid, move on
			}
		} else {
			if !json.Valid([]byte(n.Data)) {
				errCh <- fmt.Errorf("invalid json for user, moving on")
				return false // if json isn't valid, move on
			}
		}

		res := gjson.Get(n.Data, fbJsonPath)

		if !res.IsArray() {
			errCh <- fmt.Errorf("didn't find an array of friends, moving on")
			return false // if there are no friends, or the objet is wrong, move on
		}

		// Iterate through the array of json results
		for _, v := range res.Array() {
			fmt.Println("found json: populating user")
			users = append(users, User{
				name: v.Get(fbJsonNamePath).String(),
				id:   generateId(),
				urls: append([]string{}, v.Get(fbJsonUrlPath).String()),
			})
		}
		return true
	})

	if len(users) == 0 {
		errCh <- fmt.Errorf("no friends found for user, moving on")
	}

	usrsCh <- users

}

// recurses through the siblings and children of n, calling f() on each element
// When f returns true, it will stop looking
func htmlRecurser(n *html.Node, f func(*html.Node) bool) *html.Node {

	if f(n) {
		return n
	}

	element := n.FirstChild
	if element == nil {
		return nil
	}

	for element != nil {
		n := htmlRecurser(element, f)
		if n != nil {
			return n
		}

		element = element.NextSibling
	}
	return nil
}
