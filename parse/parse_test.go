package parse

import (
	"testing"
)

// testing parseDoiNumber function for parsing doi numbers out
// of url string
func Test_parseDoiNumber(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			"http://dx.doi.org/10.1080/09500340.2010.500105",
			"10.1080/09500340.2010.500105",
		},
		{
			"10.1080/09500340.2010.500105",
			"10.1080/09500340.2010.500105",
		},
	}

	for _, test := range tests {
		out, _ := parseDoiNumber(test.input)

		if out != test.output {
			t.Errorf("parseDoiNumber(%v) = %v", test.input, out)
			t.Errorf("Output should be: %v", test.output)
		}
	}
}

func Test_parseLink(t *testing.T) {
	tests := []struct {
		html   string
		tagID  string
		output string
	}{
		{
			// there is space between class and =, since go
			// downloads website from scihub in that way
			// ex. id = "some id"
			// TODO: Find a cleaner way to solve this problem?
			`<div class = "some-class">
				<a href="http://www.website1.com"></a>
			</div>

			<div id = "some-id">
				<a href="http://www.website2.com"></a>
			</div>`,

			// tag id
			"some-id",

			// output
			"http://www.website2.com",
		},
	}

	for _, test := range tests {
		out, _ := parseLink(test.html, test.tagID)
		if out != test.output {
			t.Errorf("parseLink(input) = %v", out)
			t.Errorf("Output should be: %v", test.output)
		}
	}
}

func Test_parsePdfName(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"http://moscow.sci-hub.cc/ab00ac9007edb544d7251d3f6e6c6c0e/jiang2010.pdf", "jiang2010.pdf"},
		{"http://www.reddit.com/r/golang", "golang"},
		{"http://www.somewebsite.com/random/completely+user=someuser+randomstuff.pdf", "completely+user=someuser+randomstuff.pdf"},
	}

	for _, test := range tests {
		out := parsePdfName(test.input)
		if out != test.output {
			t.Errorf("parsePdfName(%v) = %v", test.input, out)
			t.Errorf("Output should be: %v", test.output)
		}
	}
}
