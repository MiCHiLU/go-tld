package tld

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

// IANA is the URL to the IANA TLD list.
var IANA = "https://data.iana.org/TLD/tlds-alpha-by-domain.txt"

// Check if tld is valid.
func Valid(tld []byte) bool {
	for _, t := range TLDs {
		if bytes.Equal(tld, t) {
			return true
		}
	}
	return false
}

// Update will update the default list of TLDs. You can use the included IANA
// URL or host your own TLD list. The format of the returned data is one TLD
// per line, lines that start with # are ignored, and unicode domains must be
// punycode encoded.
func Update(url string) (err error) {
	var TLDs = make([][]byte, 0)
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	r := bufio.NewReader(resp.Body)
	var completeLine []byte
	for {
		line, prefix, err := r.ReadLine()
		if prefix {
			completeLine = append(completeLine, line...)
			continue
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		completeLine = line
		if completeLine[0] != '#' {
			tld := bytes.ToLower(completeLine)
			TLDs = append(TLDs, tld)
		}
	}
	return
}