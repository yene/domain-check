package checker

import (
	"fmt"
	"net/http"
	"strings"
)

type requestResult struct {
	finalURL   string
	statusCode int
}

// CheckDomain takes a "example.org" and returns some analytics.
func CheckDomain(domain string) {
	// make sure the user has not forgotten to provide no prefix
	domain = domainPreparation(domain)

	wwwHTTP := "http://www." + domain
	wwwHTTPS := "https://www." + domain
	nowwwHTTP := "http://" + domain
	nowwwHTTPS := "https://" + domain
	urls := []string{wwwHTTP, wwwHTTPS, nowwwHTTP, nowwwHTTPS}
	expected := "https://www." + domain

	for _, url := range urls {
		rs := request(url)
		if rs.finalURL != expected {
			fmt.Println("Error: expected final URL to be", expected, "but got", rs.finalURL)
		}
		if rs.statusCode != 200 {
			fmt.Println("Error: statuscode is not 200, it is", rs.statusCode)
		}
	}
}

func domainPreparation(domain string) string {
	domain = strings.Replace(domain, "/", "", -1)
	domain = strings.Replace(domain, "www.", "", -1)
	domain = strings.Replace(domain, "https:", "", -1)
	domain = strings.Replace(domain, "http:", "", -1)
	return domain
}

func request(url string) requestResult {
	rs, err := http.Get(url)
	if err != nil {
		panic(err) // TODO: More idiomatic way would be to print the error and die unless it's a serious error
	}
	fu := strings.TrimRight(rs.Request.URL.String(), "/")
	return requestResult{finalURL: fu, statusCode: rs.StatusCode}
}
