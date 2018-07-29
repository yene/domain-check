# domain-check

## features
* checking if domain is reachable, uses HTTPS, and uses `www.`


## planned features
- [ ] flag with file path, newline separated domains
- [ ] take list of domains from pipe as input
- [ ] check if a path is redirected or rewritten
- [ ] check if HTTPS certificate expires soon
- [ ] test if AAAA (ipv6) records are configured
- [ ] checking the content: handling 500 and 404 and other offline responses
- [ ] machine readable output
- [ ] menu bar UI for macOS


## Why www?
Without `www.` browsers send the cookies to all subdomains. For example if you login on 'example.org' the domain 'cdn.example.org' also receives the session cookie.

## Material
* https://jonathanmh.com/tracing-preventing-http-redirects-golang/
