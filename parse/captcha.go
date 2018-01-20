package parse

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Captcha struct holds captcha fields and provides convinient api for gathering
// captcha data
type Captcha struct {
	ID    string // captcha id number
	Image string // base64 encoded string for including in template
	URL   string // captcha full url
	// Both items below are needed for submitting captcha answer to Scihub servers
	ArticleURL string // direct link to article
	ArticleDoi string // doi of the article
}

// Download parses relevant captcha details and returns error if any
// download part fails
func (c *Captcha) Download(captchaHTML string) error {
	err := c.getCaptchaURL(captchaHTML)
	if err != nil {
		return err
	}

	err = c.getCaptchaID()
	if err != nil {
		return err
	}

	// fetch captcha image and turn it into base64 string
	err = c.getImage()
	if err != nil {
		return err
	}
	return nil
}

// getCaptchaURL creates captcha url from provided captcha html and
// article url. Parsed captcha url is used to download captcha image
func (c *Captcha) getCaptchaURL(captchaHTML string) error {
	arr := strings.Split(c.ArticleURL, "/")
	baseURL := fmt.Sprintf("%s%s//%s", arr[0], arr[1], arr[2])

	html := captchaHTML
	imgTagStart := strings.Index(html, "captcha")
	if imgTagStart < 0 {
		return fmt.Errorf("Could not parse captcha image from scihub server")
	}
	html = html[imgTagStart:]

	// parse captcha relative url "/img/captcha_number.jpg"
	start := strings.Index(html, "src")
	html = html[start+len(`src="`):] // "/img/number.jpg"+...some more html"
	end := strings.Index(html, `"`)
	captchaRelativeURL := html[:end] // "/img/number.jpg"
	// create full captcha url
	captchaURL := baseURL + captchaRelativeURL // "http://dacemirror.scihub.org/img/captcha_number.jpg"
	c.URL = captchaURL
	return nil
}

// getCaptchaID gets id from captcha url. Captcha ID is needed, since the ID has to be sent
// to scihub server as part of the captcha post request (Scihub server checks if user sent text
// is correct for the given Captcha ID)
func (c *Captcha) getCaptchaID() error {
	arr := strings.Split(c.URL, "/")
	captchaID := arr[len(arr)-1] // number.jpg

	idArr := strings.Split(captchaID, ".")
	if len(idArr) < 2 {
		return fmt.Errorf("getCaptchaID: Captcha id is wrongly formatted, Scihub website design probably changed")
	}
	id := idArr[0]
	c.ID = id
	return nil
}

// getImage fetches captcha image from Captcha.URL and turns it into
// base64 string for embedding into html template
func (c *Captcha) getImage() error {
	resp, err := http.Get(c.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return fmt.Errorf("Scihub server status code: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	encodedStr := base64.StdEncoding.EncodeToString(body)
	c.Image = encodedStr
	return nil
}
