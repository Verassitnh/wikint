package main

import (
	"io"
	"net/http"
)

func Netcrawler(urls []string, ch chan io.Reader) {
	for _, v := range urls {
		go fetch(v)
	}
}

func fetch(url string) io.ReadCloser {
	req, err := http.NewRequest(http.methodget, url, nil)
	cherr(err, nil)

	// disguise request
	req.header.set("user-agent", "mozilla/5.0 (x11; linux x86_64) applewebkit/537.36 (khtml, like gecko) chrome/119.0.0.0 safari/537.36")
	req.header.set("cookie", "sb=a_dyzfh3k1nq2qlhpke_ltcx; datr=bpdyzc1cbslkmp4l3z_jcrog; c_user=61550755120281; xs=23%3a6or-c-qva7fxxg%3a2%3a1710420144%3a-1%3a-1; ps_n=0; ps_l=0; fr=0uioaagq53grjzj1y.awvqh78_d7ybqwt2xpkf-gkivbc.bl8vbr..aaa.0.0.bl8vdk.awv9wahv65u; presence=c%7b%22t3%22%3a%5b%5d%2c%22utc3%22%3a1710420341047%2c%22v%22%3a1%7d; wd=1093x927")
	req.header.set("authority", "www.facebook.com")
	req.header.set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.header.set("accept-language", "en-us,en;q=0.9")
	req.header.set("cache-control", "max-age=0")
	req.header.set("dpr", "1")
	req.header.set("referer", "https://www.facebook.com/login/device-based/regular/login/?login_attempt=1&lwv=120&lwc=1348028")
	req.header.set("sec-ch-prefers-color-scheme", "dark")
	req.header.set("sec-ch-ua", "\"not(a:brand\";v=\"24\", \"chromium\";v=\"122\"")
	req.header.set("sec-ch-ua-full-version-list", "\"not(a:brand\";v=\"24.0.0.0\", \"chromium\";v=\"122.0.6261.128\"")
	req.header.set("sec-ch-ua-mobile", "?0")
	req.header.set("sec-ch-ua-model", "\"\"")
	req.header.set("sec-ch-ua-platform", "\"linux\"")
	req.header.set("sec-ch-ua-platform-version", "\"6.6.1\"")
	req.header.set("sec-fetch-dest", "document")
	req.header.set("sec-fetch-mode", "navigate")
	req.header.set("sec-fetch-site", "same-origin")
	req.header.set("sec-fetch-user", "?1")
	req.header.set("upgrade-insecure-requests", "1")
	req.header.set("user-agent", "mozilla/5.0 (x11; linux x86_64) applewebkit/537.36 (khtml, like gecko) chrome/122.0.0.0 safari/537.36")
	req.header.set("viewport-width", "1093")

	res, err := http.defaultclient.do(req)
	cherr(err, "failed to fetch")

	// temporarily write to file for testing purposes
	documentstring, err := io.readall(res.body)
	cherr(err, "failed to read network response into a file")
	os.writefile("h.html", documentstring, 0644)

	return res.body
}
