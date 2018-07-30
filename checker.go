package checker

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const warnCertificateExpiringInDays = 80

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

	result := expirationCheck(domain + ":443")
	for _, err := range result.certErrors {
		fmt.Println(err.err)
	}

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

type certErrors struct {
	commonName string
	err        error
}

type hostResult struct {
	host       string
	err        error
	certErrors []certErrors
}

func expirationCheck(host string) (result hostResult) {
	result = hostResult{
		host:       host,
		certErrors: []certErrors{},
	}

	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		result.err = err
		fmt.Println("error", err)
		return
	}
	timeNow := time.Now()
	checkedCerts := make(map[string]struct{})
	for _, chain := range conn.ConnectionState().VerifiedChains {
		for _, cert := range chain {
			if _, checked := checkedCerts[string(cert.Signature)]; checked {
				continue
			}
			checkedCerts[string(cert.Signature)] = struct{}{}
			var cErrs error

			// Check the expiration.
			if timeNow.AddDate(0, 0, warnCertificateExpiringInDays).After(cert.NotAfter) {
				expiresIn := int64(cert.NotAfter.Sub(timeNow).Hours())
				const errExpiringSoon = "%s: '%s' (S/N %X) expires in roughly %d days."
				cErrs = fmt.Errorf(errExpiringSoon, host, cert.Subject.CommonName, cert.SerialNumber, expiresIn/24)
				result.certErrors = append(result.certErrors, certErrors{
					commonName: cert.Subject.CommonName,
					err:        cErrs,
				})
			}
		}
	}

	return result
}
