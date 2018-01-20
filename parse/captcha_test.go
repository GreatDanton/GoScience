package parse

import "testing"

func Test_getCaptchaURL(t *testing.T) {
	tests := []struct {
		input      string
		captchaURL string
		articleURL string
	}{
		{
			`
		<html>
		<body>
		<iframe id="pdf" src="http://dacemirror.sci-hub.hk/journal-article/741d97542057cb863df59bc6dcc699c6/sviridov2006.pdf">
			<form action="" method="POST">
				<img src="/img/wrongImage.jpg" />
				<p><img id="captcha" src="/img/5a566b72e229c.jpg"></p>
				<img src="/img/wrongImage2.jpg" />
				<input name="id" value="5a566b72e229c" type="hidden">
				<input maxlength="6" name="answer" style="width:256px;font-size:18px;height:36px;margin-top:18px;text-align:center" autofocus="" type="text"><br>
				<p style="margin-top:22px"><input value="send" type="submit"></p>
			</form>
		</iframe>
		</body>
		</html>
		`,
			`http://dacemirror.sci-hub.hk/img/5a566b72e229c.jpg`,
			`http://dacemirror.sci-hub.hk/journal-article/741d97542057cb863df59bc6dcc699c6/sviridov2006.pdf`,
		},
		{
			`
			<form action="" method="POST">
				<img src="/img/wrongImage.jpg" />
				<p><img id="captcha" src="/img/5a566b72e229c.jpg"></p>
				<img src="/img/wrongImage2.jpg" />
				<input name="id" value="5a566b72e229c" type="hidden">
				<input maxlength="6" name="answer" style="width:256px;font-size:18px;height:36px;margin-top:18px;text-align:center" autofocus="" type="text"><br>
				<p style="margin-top:22px"><input value="send" type="submit"></p>
			</form>
			`,
			`http://dacemirror.sci-hub.hk/img/5a566b72e229c.jpg`,
			`http://dacemirror.sci-hub.hk/journal-article/741d97542057cb863df59bc6dcc699c6/sviridov2006.pdf`,
		},
	}

	for _, test := range tests {
		c := Captcha{ArticleURL: test.articleURL}
		err := c.getCaptchaURL(test.input)
		if err != nil {
			t.Errorf("getCaptchaID() reported error that should not occur: %v", err)
		}
		if c.URL != test.captchaURL {
			t.Errorf("getCaptchaURL() = %v", c.URL)
			t.Errorf("Output should be: %v", test.captchaURL)
		}
	}
}

func Test_getCaptchaID(t *testing.T) {
	tests := []struct {
		captchaURL string
		output     string
	}{
		{
			`http://dacemirror.sci-hub.hk/img/5a566b72e229c.jpg`,
			"5a566b72e229c",
		},
		{
			`http://dacemirror.sci-hub.hk/img/1234.png`,
			"1234",
		},
	}

	for _, test := range tests {
		c := Captcha{URL: test.captchaURL}
		err := c.getCaptchaID()
		if err != nil {
			t.Errorf("getCaptchaID(%v) returned error: %v", test.captchaURL, err)
		}
		if c.ID != test.output {
			t.Errorf("getCaptchaID(%v) returned %v", test.captchaURL, c.ID)
			t.Errorf("It should return %v", test.output)
		}
	}
}
