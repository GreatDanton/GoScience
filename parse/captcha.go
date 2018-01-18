package parse

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// CaptchaStruct holds all variables related to captcha template rendering
type CaptchaStruct struct {
	Captcha     string
	CaptchaID   string
	ArticleDoi  string
	ArticleLink string
}

// TODO: rewrite this part and use OOP
// getCaptchaURL returns url of captcha image used for downloading captcha image,
// from scihub website
func getCaptchaURL(html string, articleDirectLink string) (string, error) {

	baseURLStart := strings.Index(html, `id="pdf"`)
	if baseURLStart < 0 {
		// iframe with id="pdf"	does not exist,
		// create captcha url from direct pdf link and parse captcha relative url from input html
		captchaURL, err := createCaptchaURL(html, articleDirectLink)
		return captchaURL, err
	}
	html = html[baseURLStart:]
	baseURLStart = strings.Index(html, "src")
	html = html[baseURLStart+len(`src="`):]
	baseURLEnd := strings.Index(html, `"`)

	pdfDirectLink := html[:baseURLEnd] // directPDFLink = http://dacemirror.sci-hub.xx/journal-article/xxxxxx
	captchaURL, err := createCaptchaURL(html, pdfDirectLink)
	return captchaURL, err
}

func createCaptchaURL(html string, pdfDirectLink string) (string, error) {
	arr := strings.Split(pdfDirectLink, "/")
	baseURL := fmt.Sprintf("%s%s//%s", arr[0], arr[1], arr[2])

	imgTagStart := strings.Index(html, "captcha")
	if imgTagStart < 0 {
		return "", fmt.Errorf("Could not parse captcha image from scihub server")
	}
	html = html[imgTagStart:]

	start := strings.Index(html, "src")
	html = html[start+len(`src="`):] // (/img/number.jpg"...some more html)
	end := strings.Index(html, `"`)
	captchaRelativeURL := html[:end]           // (/img/number.jpg)
	captchaURL := baseURL + captchaRelativeURL // http://dacemirror.scihub.org/img/captcha_number.jpg
	return captchaURL, nil
}

// parse captcha id for hidden field in captcha template from provided captcha html
// returned from scihub website
func getCaptchaID(captchaURL string) (string, error) {
	arr := strings.Split(captchaURL, "/")

	captchaID := arr[len(arr)-1] // number.jpg

	idArr := strings.Split(captchaID, ".")
	if len(idArr) < 2 {
		return "", fmt.Errorf("getCaptchaID: Captcha id is wrongly formatted, Scihub website design probably changed")
	}
	id := idArr[0]
	return id, nil
}

// DownloadCaptcha downloads captcha image and turns it into base64 string, that is later
// embedded into captcha template (the end user can see captcha image and solve it via input field)
func DownloadCaptcha(captchaHTML string, articleDirectLink string) (CaptchaStruct, error) {
	c := CaptchaStruct{}
	captchaURL, err := getCaptchaURL(captchaHTML, articleDirectLink)
	if err != nil {
		return c, err
	}

	id, err := getCaptchaID(captchaURL)
	if err != nil {
		return c, err
	}

	resp, err := http.Get(captchaURL)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return c, fmt.Errorf("Scihub server status code: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c, err
	}
	encodedStr := base64.StdEncoding.EncodeToString(body)

	c.Captcha = encodedStr
	c.CaptchaID = id
	c.ArticleLink = articleDirectLink

	return c, nil
}
