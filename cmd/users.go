package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

const (
	cursor         = "AQHR_QntGMVOBk3ZPbLFCk8F0ZDCIgQTJV6iKIrbdWCsa-17t9Bwm5PHaaFoe6eNPYf5kS8OhjJwSd3ZiyvMbkyeYw"
	id             = "YXBwX2NvbGxlY3Rpb246MTAwMDg2NTU2MDAzMTMwOjIzNTYzMTgzNDk6Mg"
	usrsPath       = "data.node.pageItems.edges"
	namePath       = "node.title.text"
	idPath         = "node.id"
	urlPath        = "node.url"
	nextCursorPath = "data.node.pageItems.page_info.end_cursor"
	lastpagePath   = "data.node.pageItems.page_info.has_next_page"
	cursorPath     = "cursor"
	globalIdPath   = "require.0.3.0.__bbox.require.11.3.1].__bbox.result.data.node.id"
	susrsPath      = "require.0.3.0.__bbox.require.11.3.1.__bbox.result.data.node.all_collections.nodes.0.style_renderer.collection.pageItems.edges"

	filename     = "debug_output.json"
	startingURL  = "https://www.facebook.com/profile.php?id=100078255380484&sk=friends"
	fbNextCursor = "require.0.3..0.__bbox.require.11.3.1.__bbox.result.data.user.timeline_nav_app_sections.page_info.end_cursor"
)

type ScrapedData struct {
	url  string
	body io.Reader
}

type ResultReciever[T any] struct {
	dataCh chan T
	errCh  chan error
}

type FetchResultReciever ResultReciever[io.Reader]

func NewReciever[T any](data T) ResultReciever[T] {
	return ResultReciever[T]{
		dataCh: make(chan T),
		errCh:  make(chan error),
	}
}

type User struct {
	id         string
	urls       []string
	name       string
	occupation Ocupation
	interests  []string
	cursor     string // temp value to retrieve user data
}

type Ocupation struct {
	id     string
	name   string
	skills []string
}

type BodyOptions struct {
	cursor string
	id     string
}

type FetchReq struct {
	url    string
	body   io.Reader
	method string
}

func dqGraph(r *http.Request) *http.Request {
	r.Header.Add("authority", "www.facebook.com")
	r.Header.Add("accept", "*/*")
	r.Header.Add("accept-language", "en-US,en;q=0.9")
	r.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.Header.Add("cookie", "sb=a_DyZfH3k1NQ2qLhPkE_LtCx; datr=bPDyZc1CbSLKMP4L3z_JCrog; c_user=61550755120281; ps_n=0; ps_l=0; dpr=2; locale=en_US; vpd=v1%3B667x375x2; wl_cbv=v2%3Bclient_version%3A2447%3Btimestamp%3A1711222828; xs=23%3A6Or-C-QVa7FXXg%3A2%3A1710420144%3A-1%3A-1%3A%3AAcUJZypyOjlY_QLL8xQxhfgkrTv1dNuoyN7fFow6t50; fr=1dGRo50GEfGvF9a6t.AWWpv_I7P7u6WlynslapJRFq_Yw.Bl_5Bn..AAA.0.0.Bl_5HL.AWXHrBjiRLM; usida=eyJ2ZXIiOjEsImlkIjoiQXNhdHo4b2lzNXpjciIsInRpbWUiOjE3MTEyNDc4MjF9; wd=814x927; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1711342375817%2C%22v%22%3A1%7D")
	r.Header.Add("dpr", "1")
	r.Header.Add("origin", "https://www.facebook.com")
	r.Header.Add("referer", "https://www.facebook.com/profile.php?id=100078255380484&sk=friends_all")
	r.Header.Add("sec-ch-prefers-color-scheme", "dark")
	r.Header.Add("sec-ch-ua", "'Not(A:Brand';v='24', 'Chromium';v='122'")
	r.Header.Add("sec-ch-ua-full-version-list", "'Not(A:Brand';v='24.0.0.0', 'Chromium';v='122.0.6261.128'")
	r.Header.Add("sec-ch-ua-mobile", "?0")
	r.Header.Add("sec-ch-ua-model", "''")
	r.Header.Add("sec-ch-ua-platform", "'Linux'")
	r.Header.Add("sec-ch-ua-platform-version", "'6.7.9'")
	r.Header.Add("sec-fetch-dest", "empty")
	r.Header.Add("sec-fetch-mode", "cors")
	r.Header.Add("sec-fetch-site", "same-origin")
	r.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	r.Header.Add("viewport-width", "814")
	r.Header.Add("x-asbd-id", "129477")
	r.Header.Add("x-fb-friendly-name", "ProfileCometAppCollectionListRendererPaginationQuery")
	r.Header.Add("x-fb-lsd", "JuNwYI5mcy28LM_Mjd2Kvx")
	return r
}

func dqProfile(r *http.Request) *http.Request {

	r.Header.Add("authority", "www.facebook.com")
	r.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	r.Header.Add("accept-language", "en-US,en;q=0.9")
	r.Header.Add("cache-control", "max-age=0")
	r.Header.Add("cookie", "sb=a_DyZfH3k1NQ2qLhPkE_LtCx; datr=bPDyZc1CbSLKMP4L3z_JCrog; c_user=61550755120281; ps_n=0; ps_l=0; dpr=2; locale=en_US; vpd=v1%3B667x375x2; wl_cbv=v2%3Bclient_version%3A2447%3Btimestamp%3A1711222828; usida=eyJ2ZXIiOjEsImlkIjoiQXNhdHo4b2lzNXpjciIsInRpbWUiOjE3MTEyNDc4MjF9; wd=814x927; xs=23%3A6Or-C-QVa7FXXg%3A2%3A1710420144%3A-1%3A-1%3A%3AAcUkDIk7TOc9zosmNBf5yc6uST3as4z_PsraPS9ycW8; fr=1wElodwBRfIkotvX2.AWWGJwWKxJjq71duzSI0wkp5uRE.BmARAL..AAA.0.0.BmARAL.AWXyas0nifQ; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1711345798762%2C%22v%22%3A1%7D")
	r.Header.Add("dpr", "1")
	r.Header.Add("sec-ch-prefers-color-scheme", "dark")
	r.Header.Add("sec-ch-ua", "'Not(A:Brand';v='24', 'Chromium';v='122'")
	r.Header.Add("sec-ch-ua-full-version-list", "'Not(A:Brand';v='24.0.0.0', 'Chromium';v='122.0.6261.128'")
	r.Header.Add("sec-ch-ua-mobile", "?0")
	r.Header.Add("sec-ch-ua-model", "''")
	r.Header.Add("sec-ch-ua-platform", "'Linux'")
	r.Header.Add("sec-ch-ua-platform-version", "'6.7.9'")
	r.Header.Add("sec-fetch-dest", "document")
	r.Header.Add("sec-fetch-mode", "navigate")
	r.Header.Add("sec-fetch-site", "same-origin")
	r.Header.Add("sec-fetch-user", "?1")
	r.Header.Add("upgrade-insecure-requests", "1")
	r.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	r.Header.Add("viewport-width", "814")

	return r
}

func setBody(b BodyOptions) io.Reader {
	bs := `
	av=61550755120281&__aaid=0&__user=61550755120281&__a=1&__req=10&__hs=19807.HYP%3Acomet_pkg.2.1..2.1&dpr=1&__ccg=GOOD&__rev=1012280021&__s=pbw2be%3Abm12en%3Ayflj8o&__hsi=7350158820251186660&__dyn=7AzHK4HwkEng5K8G6EjBAg2owIxu13wFwnUW3q2ibwNwnof8boG0x8bo6u3y4o2Gwfi0LVEtwMw65xO2OU7m221Fwgo9oO0-E7m4oaEnxO0Bo7O2l2Utwwwi831wiE567Udo5qfK0zEkxe2GewyDwkUtxGm2SUbElxm3y11xfxmu3W3y261eBx_wHwdG7FoarCwLyES1Iwh888cA0z8c84q58jyUaUcojxK2B08-269wqQ1FwgU4q3G1eKufxa3m7E&__csr=gFs8layNROQOcQLO4_vqnd-zncORO99OihasGaLTisyj9rRTcG-myFT9sylpqGCF9BD-RKkOCJdeO5RQdyoKucjgzmEjhoCHzV8W9Q9y9AczAazFUCi68CeyQ4oV5yeq2ecwzVUC8Lwyz8cqypUx4xh1h0LwIG4EW2a2q6um321mwiU4byE562-2W3mfwnE8uu2S1wwIgiwgo4m0hSm1tx-10xu0V83Lw7gw3wo6F0gUcU2Vw15S00kVd02SUlglUV02JQ0WE0ggw1oAE4O0h63C0P804t-0bxw0o181Xo1o8mwSa0jy&__comet_req=15&fb_dtsg=NAcNuL21y2OinB-H8_00IujwF0S97BkmLb1etInAjGE2vQz04KY3IEQ%3A23%3A1710420144&jazoest=25081&lsd=wUNsBRQom6XDOR4vbNT_kh&__spin_r=1012280021&__spin_b=trunk&__spin_t=1711342209&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=ProfileCometAppCollectionListRendererPaginationQuery&variables=%7B%22count%22%3A8%2C%22cursor%22%3A%22` + b.cursor + `%22%2C%22scale%22%3A1%2C%22search%22%3Anull%2C%22id%22%3A%22` + b.id + `%22%2C%22__relay_internal__pv__VideoPlayerRelayReplaceDashManifestWithPlaylistrelayprovider%22%3Afalse%7D&server_timestamps=true&doc_id=7019108158194686
	`
	return strings.NewReader(bs)

}

func graphPing(c string, id string, fr FetchResultReciever) {
	f := FetchReq{
		url:    "https://www.facebook.com/api/graphql/",
		method: http.MethodPost,
		body: setBody(BodyOptions{
			cursor: c,
			id:     id,
		}),
	}

	go fetch(f, fr, dqGraph)
}

func profilePing(url string, fr FetchResultReciever) {
	f := FetchReq{
		url:    url,
		method: http.MethodGet,
		body:   nil,
	}

	go fetch(f, fr, dqProfile)
}

func fetch(f FetchReq, fr FetchResultReciever, dq func(*http.Request) *http.Request) {
	r, err := http.NewRequest(f.method, f.url, f.body)
	if err != nil {
		fr.errCh <- err
	}

	// disguise the request
	r = dq(r)

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		fr.errCh <- fmt.Errorf("killed req: %v", f.url)

	} else {
		fr.dataCh <- res.Body
	}
}

func main() {
	ScrapeUser(startingURL)
}

// High level function: put helper functions together
func ScrapeUser(url string) {
	fr := FetchResultReciever{
		dataCh: make(chan io.Reader),
		errCh:  make(chan error),
	}
	r := ResultReciever[User]{
		dataCh: make(chan User),
		errCh:  fr.errCh,
	}

	db, err := Database("../my.db", r.errCh)
	if err != nil {
		log.Fatal("failed to start database")
	}

	defer db.Destroy()
	go profilePing(url, fr)

	// Infintely listen for results
	for {
		select {
		case pd := <-fr.dataCh:
			go handleFbData(pd, r, fr)
		case usr := <-r.dataCh:
			fmt.Printf("Recieved user: %v\n", usr.name)
			go db.AppendUser(usr)
			go profilePing(usr.urls[0], fr)
		case err := <-r.errCh:
			fmt.Printf("error: %v\n", err)
		}
	}

}

func handleFbData(sd io.Reader, r ResultReciever[User], fr FetchResultReciever) {

	// Read the data into a doc
	var doc *goquery.Document
	var err error
	var cursor string
	var id string

	if doc, err = goquery.NewDocumentFromReader(sd); err != nil {
		r.errCh <- err
		return
	}

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		fmt.Print(".")
		if v, ok := s.Attr("type"); ok && v != "application/json" {
			// r.errCh <- fmt.Errorf("wasn't of type application/json")
			return
		}
		jd := s.Text()

		if !json.Valid([]byte(jd)) {
			// r.errCh <- fmt.Errorf("didn't find valid json")
			return
		}

		res := gjson.Get(jd, susrsPath)

		if !res.IsArray() {
			// r.errCh <- fmt.Errorf("didn't find json array")
			return
		}

		// Iterate through the array of json results
		for _, v := range res.Array() {

			// handle two different profile types
			url := v.Get(urlPath).Str
			if strings.Contains(url, "profile.php") {
				url += "&sk=friends"
			} else {
				url += "/friends"
			}

			r.dataCh <- User{
				name: v.Get(namePath).Str,
				id:   v.Get(idPath).Str,
				urls: []string{url},
			}
		}

		id = gjson.Get(jd, globalIdPath).Str
		cursor = res.Get(nextCursorPath).Str
	})

	fmt.Printf("calling graphPing: cursor: %v\nid: %v\n", cursor, id)
	go graphPing(cursor, id, fr)

}
