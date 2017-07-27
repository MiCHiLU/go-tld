// Package tld provides tools for validating top level domain names.
// It comes with a predefined list however it can be updated at runtime
// by running tld.Update(url) where url points to a text file containing a list
// of acceptable TLDs.
package tld

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

// IANA is the URL to the IANA TLD list and should be used for Update() unless
// you need to specify your own custom valid TLDs.
var IANA = "https://data.iana.org/TLD/tlds-alpha-by-domain.txt"

// Valid checks if the TLD is found in the list. There are other requirements of a TLD
// but we assume that if it's in our list, it's valid.
func Valid(tld []byte) bool {
	tld = bytes.ToLower(tld) // convert to lowercase to match our list
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
