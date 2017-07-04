package parse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ScihubURL is the entry point for downloading articles on scihub (front page url),
// it might change in the future that's why it's provided as constant at the top of
// the page
const ScihubURL = "http://sci-hub.cc/"

/* Example usage:
func main() {
	doi := "10.1080/09500340.2010.500105"
	GetPdf(doi)
}*/

// GetPdf fetches pdf from doi string, and returns byte stream(pdf), name of pdf
// article and error (in case of error)
func GetPdf(doi string) ([]byte, string, error) {
	d, err := parseDoiNumber(doi)
	// cannot parse doi numbers from doi string
	if err != nil {
		return nil, "", err
	}

	url := fmt.Sprintf("%v%s", ScihubURL, d)
	html, err := getHTML(url)
	// cannot get html page -> returning error
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	pdfLink, err := parseLink(html, "content")
	pdfName := parsePdfName(pdfLink)
	// link does not exist is doi correct? //TODO: check for captcha ?
	if err != nil {
		return nil, "", fmt.Errorf("GetPdf: Could not download pdf, provided doi '%v' does not exist", doi)
	}

	pdfResp, err := http.Get(pdfLink)
	if err != nil {
		return nil, "", err
	}
	defer pdfResp.Body.Close()

	if pdfResp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("%v", pdfResp.StatusCode)
	}

	// we got the stream, return it
	pdf, err := ioutil.ReadAll(pdfResp.Body)
	if err != nil {
		return nil, "", err
	}

	fmt.Println("File downloaded")
	// return bytes stream == pdf
	return pdf, pdfName, nil
}

// function for parsing doi number out of doi string
// returning just the number part from
// input: http://dx.doi.org/10.1080/09500340.2010.500105
// output: 10.1080/09500340.2010.500105
func parseDoiNumber(d string) (string, error) {
	if len(d) == 0 {
		return "", fmt.Errorf("Doi number does not exist")
	}

	// remove trailing white space
	d = strings.Trim(d, " ")

	// if http does not exist, it means we have just the
	//string of integers return string
	if strings.Index(d, "http") == -1 {
		return d, nil
	}

	// parse string of doi integers out of provided string
	domain := strings.Index(d, ".org")
	if domain == -1 {
		return "", fmt.Errorf("Could not parse doi string out of provided url: %v", d)
	}
	doi := d[domain+len(".org/"):]
	return doi, nil
}

//
// getHTML gets html page (string) from provided url
func getHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	text := string(body)
	return text, nil
}

// parseLink parses pdf link from provided html page and id tag, if link is not found
// it means the article was not found in html string and error should
// be returned
func parseLink(htmlString, tagID string) (string, error) {
	// when html is parsed from url all html tags are returned like:
	// <htmlTag id = "id">  <-- note the space between the = and "id"
	id := "id = " + `"` + tagID + `"`
	htmlTagStart := strings.Index(htmlString, id)
	// if htmlTag with id does not exist reutrn error.
	// Currently this is true, but this part should be rewritten
	// in case they decide to change their captcha implementation
	if htmlTagStart == -1 {
		return "", fmt.Errorf("%v does not exist in provided html string", tagID)
	}
	html := htmlString[htmlTagStart:]

	// get index of link starting
	startLink := strings.Index(html, "http")
	if startLink == -1 {
		return "", fmt.Errorf("parseLink `startLink` could not be found in provided html")
	}

	// get index of link ending
	endLink := strings.Index(html[startLink:], `"`)
	if endLink == -1 {
		return "", fmt.Errorf("parseLink `endLink` could not be found in provided html string")
	}

	// htmlString stays the same all the time that's why we are parsing it via [start:start+end]
	return html[startLink : startLink+endLink], nil
}

// parse pdf name from provided url string. Name is everything after last slash in url
func parsePdfName(url string) string {
	names := strings.Split(url, "/")
	name := names[len(names)-1]
	return name
}
