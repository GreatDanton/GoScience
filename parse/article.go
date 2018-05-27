package parse

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/greatdanton/goScience/global"
)

// ErrCaptchaPresent should be returned when the scihub servers return captcha
// instead of desired pdf
var ErrCaptchaPresent = errors.New("Scihub servers returned captcha, try again later")

// ErrGeneric is used when reporting error is necessary, but you don't want
// to expose app internals to the end user
var ErrGeneric = errors.New("GoScience: Internal application error, try again later")

// Article struct represents pdf article that will be fetched
// from scihub servers.
type Article struct {
	URL       string
	Doi       string
	Name      string
	PdfStream []byte
	Captcha   Captcha
}

// GetPdf will fetch the article from the Scihub servers and report an error
// if something goes wrong
func (a *Article) GetPdf(doi string) error {
	err := a.parseDoiNumber(doi)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Please check if doi string is correct")
	}

	url := fmt.Sprintf("%v%s", global.ScihubURL, a.Doi)
	htmlString, err := getHTMLStr(url)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Scihub servers are not available")
	}

	err = a.parseArticleURL(htmlString)
	if err != nil {
		return fmt.Errorf("Article with this doi does not exist")
	}

	a.parseName()

	// fetch pdf
	err = a.fetchPdf()
	if err != nil {
		return err
	}
	return nil
}

// parseDoiNumber helps with parsing doi number from user provided doi string
// and reports an error if string is not in correct format
func (a *Article) parseDoiNumber(doiStr string) error {
	// remove trailing white space
	doiStr = strings.Trim(doiStr, " ")
	if len(doiStr) == 0 {
		return fmt.Errorf("Doi number does not exist")
	}
	// if http does not exist, it means user provided just the string of doi integers
	if strings.Index(doiStr, "http") == -1 {
		a.Doi = doiStr
		return nil
	}
	// parse string of doi integers out of provided string
	domain := strings.Index(doiStr, ".org")
	if domain == -1 {
		return fmt.Errorf("Could not parse doi string out of provided url: %v", doiStr)
	}
	doi := doiStr[domain+len(".org/"):]
	a.Doi = doi
	return nil
}

// parse article url from provided html string or return an error
// if that is not possible
func (a *Article) parseArticleURL(htmlString string) error {
	// when html is parsed from url all html tags are returned like:
	// <htmlTag id = "id">  <-- note the space between the = and "id"
	tagID := "main_content"
	id := fmt.Sprintf(`id = "%s"`, tagID)
	htmlTagStart := strings.Index(htmlString, id)
	// if htmlTag with id does not exist return error.
	// Currently this is true, but this part should be rewritten
	// in case they decide to change their captcha implementation
	if htmlTagStart == -1 {
		return fmt.Errorf("'%v' does not exist in provided html string", tagID)
	}
	html := htmlString[htmlTagStart:]

	// get index of link starting
	startLink := strings.Index(html, "http")
	if startLink == -1 {
		return fmt.Errorf("`startLink` could not be found in provided html")
	}

	// get index of link ending (the link always ends with .pdf)
	endLink := strings.Index(html[startLink:], `.pdf`) + len(".pdf")
	if endLink == -1 {
		return fmt.Errorf("`endLink` could not be found in provided html string")
	}

	// htmlString stays the same all the time that's why we are parsing it via [start:start+end]
	articleURL := html[startLink : startLink+endLink]
	a.URL = articleURL
	return nil
}

// parses article name from article url. Make sure to execute
// parseArticleURL before this method.
func (a *Article) parseName() {
	names := strings.Split(a.URL, "/")
	name := names[len(names)-1]
	a.Name = name
}

// fetchPdf creates a get request on scihub servers and fetches the pdf bytes
// or returns an error if anything goes wrong (such as scihub displaying captcha)
func (a *Article) fetchPdf() error {
	pdfResp, err := http.Get(a.URL)
	if err != nil {
		fmt.Println(err)
		return ErrGeneric
	}
	defer pdfResp.Body.Close()

	// return http status code as error stream
	if pdfResp.StatusCode != http.StatusOK {
		if pdfResp.StatusCode == http.StatusBadGateway {
			return fmt.Errorf("Scihub servers are over capacity, try again later")
		}
		return fmt.Errorf("Scihub server status code: %v", pdfResp.Status)
	}

	// Captcha check: if captcha is present on scihub (Content-Type in headers
	// is text/html instead of application/pdf)
	content := pdfResp.Header.Get("Content-type")
	if strings.Contains(content, "text/html") {
		html, err := ioutil.ReadAll(pdfResp.Body)
		if err != nil {
			fmt.Println(err)
			return ErrGeneric
		}
		captcha := Captcha{ArticleDoi: a.Doi, ArticleURL: a.URL}
		// download captcha details
		err = captcha.Download(string(html))
		if err != nil {
			return err
		}
		a.Captcha = captcha // embed captcha inside article struct
		// return error about captcha being present so the outer layer can detect
		// captcha error and display new captcha template with relevant data
		return ErrCaptchaPresent
	}

	// everything is allright, we got the pdf byte stream, return it
	pdf, err := ioutil.ReadAll(pdfResp.Body)
	if err != nil {
		fmt.Println(err)
		return ErrGeneric
	}
	a.PdfStream = pdf
	return nil
}

// getHTMLStr fetches url and returns html string of website
func getHTMLStr(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Scihub server status code: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	text := string(body)
	return text, nil
}
