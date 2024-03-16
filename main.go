package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/segmentio/encoding/json"
)

func main() {
	startingUser := Friend{url: "https://www.facebook.com/henry.jolly.94"}
	startingUser.Find(startingUser.url)
}

func cherr(err error, msg any) {
	if err != nil {
		fmt.Printf("ERRROR: %v - %v", msg, err)
		os.Exit(1)
	}
}

func fetch(url string) io.ReadCloser {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	cherr(err, nil)

	// Disguise request
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Set("cookie", "sb=a_DyZfH3k1NQ2qLhPkE_LtCx; datr=bPDyZc1CbSLKMP4L3z_JCrog; c_user=61550755120281; xs=23%3A6Or-C-QVa7FXXg%3A2%3A1710420144%3A-1%3A-1; ps_n=0; ps_l=0; fr=0UIoaagQ53gRjzj1Y.AWVqh78_D7YbQwt2XpkF-Gkivbc.Bl8vBr..AAA.0.0.Bl8vDK.AWV9waHv65U; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1710420341047%2C%22v%22%3A1%7D; wd=1093x927")
	req.Header.Set("authority", "www.facebook.com")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("dpr", "1")
	req.Header.Set("referer", "https://www.facebook.com/login/device-based/regular/login/?login_attempt=1&lwv=120&lwc=1348028")
	req.Header.Set("sec-ch-prefers-color-scheme", "dark")
	req.Header.Set("sec-ch-ua", "\"Not(A:Brand\";v=\"24\", \"Chromium\";v=\"122\"")
	req.Header.Set("sec-ch-ua-full-version-list", "\"Not(A:Brand\";v=\"24.0.0.0\", \"Chromium\";v=\"122.0.6261.128\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-model", "\"\"")
	req.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Set("sec-ch-ua-platform-version", "\"6.6.1\"")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("viewport-width", "1093")

	res, err := http.DefaultClient.Do(req)
	cherr(err, "Failed to fetch")

	// Temporarily write to file for testing purposes
	documentstring, err := io.ReadAll(res.Body)
	cherr(err, "Failed to read network response into a file")
	os.WriteFile("h.html", documentstring, 0644)

	return res.Body
}

// get all the friends of Friend
func (f Friend) Find(url string) []Friend {
	body := fetch(url)
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	cherr(err, "Document didn't create")

	var flist []Friend

	doc.Find("script [type]=\"application/json\"").Each(func(i int, s *goquery.Selection) {
		contents, err := s.Contents().Html()
		cherr(err, "Cannot for the life of me find the contents of a script")

		bytes := []byte(contents)

		if json.Valid(bytes) {
			flist, err = f.Marshal(bytes)
			cherr(err, nil)
		} else {
			log.Fatal("Invalid JSON")
		}
	})

	// return all friends
	return flist

}

type Friend struct {
	url string
}

func (f *Friend) Marshal(b []byte) ([]Friend, error) {
	for t := json.NewTokenizer(b); t.Next(); {
		switch k := t.Kind(); k.Class() {
		case json.Array:
			continue
		case json.Object:
			continue
		case json.String:
			fmt.Print(t.Value)
		}
	}

	var fbUsers []Friend
	if err := json.Unmarshal(b, &fbUsers); err != nil {
		return nil, err
	}

	return fbUsers, nil
}

// Project Steps
// 1. Fetch facebook user ids
// 2. fetch facebook user posts
// 3. define interests
// 4. read posts to find users interests
// 5. build database to compare them
// 6. frontend presentation

type Person struct {
	url string
	ids []string
}
