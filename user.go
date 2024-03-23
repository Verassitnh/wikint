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

func perr(err error) {
	fmt.Printf("error: %s \n", err)
}

// Retries given func t times, retuns result of fn in the form (T, error)
func retry[T interface{}, R interface{}](t int, fn func(...R) (T, error), args ...R) (T, error) {
	var v T
	var err error

	for i := 0; i < t; i++ {
		if v, err = fn(args...); err != nil {
			fmt.Printf("retrying (%v/%v): failed: %v \n", i, t, v)
		} else {
			break
		}
	}
	return v, nil
}

func usrScrape(startingUser User) {
	usrsCh, scrapedch, errCh := make(chan []User), make(chan ScrapedData), make(chan error)
	go fetchAll(startingUser.urls, scrapedch, errCh)
	go parse(usrsCh, scrapedch, errCh)

	// db, err := Database("my.db", errCh)

	// if err != nil {
	// 	perr(err)
	// 	db, err = retry[database](3, Database, "my.db", errCh)
	// }
	for {
		select {
		case users := <-usrsCh:
			for _, user := range users {
				fmt.Println(user.name)
			}
		case err := <-errCh:
			perr(err)
		}
	}

}

func generateId() string {
	return uuid.New().String()
}
