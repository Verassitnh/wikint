package main

import (
	"fmt"
	"io"

	"github.com/google/uuid"
)

// possibly implement a name in future so we can track them on other social media sites
// possibly implement an occupation, and skills, to build a data profile on the person.
//
// currently those both remain nil
type User struct {
	id         string
	urls       []string
	name       string
	occupation Ocupation
	interests  []string
}

type Ocupation struct {
	id     string
	name   string
	skills []string
}

type ScrapedData struct {
	url  string
	body io.Reader
}

func usrScrape(startingUser User) {
	usrsCh, scrapedch, errCh := make(chan []User), make(chan ScrapedData), make(chan error)
	go fetchAll(startingUser.urls, scrapedch, errCh)
	go parse(usrsCh, scrapedch, errCh)

	select {
	case users := <-usrsCh:
		for _, user := range users {
			fmt.Println(user.name)
		}
	case err := <-errCh:
		fmt.Println(err)
	}

}

func generateId() string {
	return uuid.New().String()
}
