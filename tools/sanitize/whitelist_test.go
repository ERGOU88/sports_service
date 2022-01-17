package sanitize

import (
	"testing"
	"strings"
)

func TestSanitizeRendersDoctypeCorrectly(t *testing.T) {
	htmlDoc := `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`
	expectedOutput := `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`
	config := `
	{
		"elements": {
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestSanitizeRemoveRemovesNonWhitelistedNodes(t *testing.T) {
	htmlDoc := `<!DOCTYPE html>
				<html>
					<head>
						<title>My SkuName</title>
					</head>
					<body>
						<div>
							<b>Bold</b>
							<i>Italic</i>
							<em>Emphatic</em>
						</div>
						<div>
							<ul>
								<li>after</li>
							</ul>
						</div>
					</body>
				</html>`
	expectedOutput := `<!DOCTYPE html><html><head><title>My SkuName</title></head><body><div><i>Italic</i></div><div></div></body></html>`
	config := `{
		"stripWhitespace": true,
		"elements": {
			"html": [],
			"head": [],
			"title": [],
			"body": [],
			"div": [],
			"i": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestSanitizeUnwrapUnwrapsNonWhitelistedNodes(t *testing.T) {
	htmlDoc := `<!DOCTYPE html>
				<html>
					<head>
						<title>My SkuName</title>
					</head>
					<body>
						<div class="not-allowed">
							<b>Bold</b>
							<i>Italic</i>
							<em>Emphatic</em>
						</div>
						<div>
							<i>After</i>
						</div>
					</body>
				</html>`
	expectedOutput := `<!DOCTYPE html><html><head><title>My SkuName</title></head><body><b>Bold</b><i>Italic</i><em>Emphatic</em><i>After</i></body></html>`
	config := `{
		"stripWhitespace": true,
		"elements": {
			"html": [],
			"head": [],
			"title": [],
			"body": [],
			"b": [],
			"i": [],
			"em": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeUnwrap(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestStripWhitespace(t *testing.T) {
	htmlDoc := `<!DOCTYPE html>
				<html>
					<head>
					</head>
					<body>
					</body>
				</html>`
	expectedOutput := `<!DOCTYPE html><html><head></head><body></body></html>`
	config := `{
		"stripWhitespace": true,
		"elements": {
			"html": [],
			"head": [],
			"body": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestStripComments(t *testing.T) {
	htmlDoc := `<!DOCTYPE html><!-- hello world --><html><head></head><body></body></html>`
	expectedOutput := `<!DOCTYPE html><html><head></head><body></body></html>`
	config := `{
		"stripComments": true,
		"elements": {
			"html": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}


func TestCanSanitizeRemoveDocumentFragment(t *testing.T) {
	htmlDoc := `<div class="my-class" style="width:100%;"><ul><li>element</li></ul><b>bold</b></div>`
	expectedOutput := `<div class="my-class"><ul></ul><b>bold</b></div>`
	config := `{
		"stripComments": true,
		"elements": {
			"div": ["class"],
			"ul": [],
			"b": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemoveFragment(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestCanSanitizeRemoveDocumentFragmentHeadType(t *testing.T) {
	htmlDoc := `<script src="/path/to/script.js" type="text/javascript" disallowed=true>alert("something");</script><div></div>`
	expectedOutput := `<script src="/path/to/script.js" type="text/javascript">alert("something");</script>`
	config := `{
		"stripComments": true,
		"elements": {
			"script": ["src", "type"]
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeUnwrapFragment(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestCanSanitizeRemoveDocumentFragmentHeadTypeInline(t *testing.T) {
	htmlDoc := `<div></div><script src="/path/to/script.js" type="text/javascript" disallowed=true>alert("something");</script><div></div>`
	expectedOutput := `<div></div><script src="/path/to/script.js" type="text/javascript">alert("something");</script><div></div>`
	config := `{
		"stripComments": true,
		"elements": {
			"div": [],
			"script": ["src", "type"]
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeUnwrapFragment(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}


func TestCanSanitizeUnwrapDocumentFragment(t *testing.T) {
	htmlDoc := `<div class="my-class" style="width:100%;"><ul><li>element</li></ul><b>bold</b></div>`
	expectedOutput := `<div class="my-class"><ul>element</ul><b>bold</b></div>`
	config := `{
		"stripComments": true,
		"elements": {
			"div": ["class"],
			"ul": [],
			"b": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeUnwrapFragment(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestCanSanitizeUnwrapDocumentFragmentHeadType(t *testing.T) {
	htmlDoc := `<script src="/path/to/script.js" type="text/javascript" disallowed=true>alert("something");</script><div></div>`
	expectedOutput := `<script src="/path/to/script.js" type="text/javascript">alert("something");</script>`
	config := `{
		"stripComments": true,
		"elements": {
			"script": ["src", "type"]
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeUnwrapFragment(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestCanStripHead(t *testing.T) {
	htmlDoc := `<!DOCTYPE html><html><head></head></html>`
	expectedOutput := `<!DOCTYPE html><html></html>`
	config := `{
		"elements": {
			"html": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}

func TestCanStripBody(t *testing.T) {
	htmlDoc := `<!DOCTYPE html><html><body></body></html>`
	expectedOutput := `<!DOCTYPE html><html></html>`
	config := `{
		"elements": {
			"html": []
		}
	}`

	whitelist, _ := NewWhitelist([]byte(config))
	output, _ := whitelist.SanitizeRemove(strings.NewReader(htmlDoc))

	if output != expectedOutput {
		t.Errorf("failed: %s != %s", output, expectedOutput)
	}
}
