package markdown

import (
	"regexp"
	"testing"

	"strings"

	"github.com/eternalsad/markdownify/html"
	"github.com/eternalsad/markdownify/parser"
)

func TestEmphasis(t *testing.T) {
	tests := readTestFile2(t, "emphasis.test")
	doTestsInlineParam(t, tests, TestParams{})
}

func TestBug309(t *testing.T) {
	var tests = []string{
		`*f*—`,
		"<p><em>f</em>—</p>\n",
	}
	p := TestParams{}
	p.extensions = parser.NoIntraEmphasis
	doTestsInlineParam(t, tests, p)
}

func TestReferenceOverride(t *testing.T) {
	var tests = []string{
		"test [ref1][]\n",
		"<p>test <a href=\"http://www.ref1.com/\" title=\"Reference 1\">ref1</a></p>\n",

		"test [my ref][ref1]\n",
		"<p>test <a href=\"http://www.ref1.com/\" title=\"Reference 1\">my ref</a></p>\n",

		"test [ref2][]\n\n[ref2]: http://www.leftalone.com/ (Ref left alone)\n",
		"<p>test <a href=\"http://www.overridden.com/\" title=\"Reference Overridden\">ref2</a></p>\n",

		"test [ref3][]\n\n[ref3]: http://www.leftalone.com/ (Ref left alone)\n",
		"<p>test <a href=\"http://www.leftalone.com/\" title=\"Ref left alone\">ref3</a></p>\n",

		"test [ref4][]\n\n[ref4]: http://zombo.com/ (You can do anything)\n",
		"<p>test [ref4][]</p>\n",

		"test [!(*http.ServeMux).ServeHTTP][] complicated ref\n",
		"<p>test <a href=\"http://localhost:6060/pkg/net/http/#ServeMux.ServeHTTP\" title=\"ServeHTTP docs\">!(*http.ServeMux).ServeHTTP</a> complicated ref</p>\n",

		"test [ref5][]\n",
		"<p>test <a href=\"http://www.ref5.com/\" title=\"Reference 5\">Moo</a></p>\n",

		"test [ref6]\n",
		"<p>test <a href=\"http://www.ref6.com/\" title=\"Reference 6\">Moo</a></p>\n",
	}
	doTestsInlineParam(t, tests, TestParams{
		referenceOverride: func(reference string) (rv *parser.Reference, overridden bool) {
			switch reference {
			case "ref1":
				// just an overridden reference exists without definition
				return &parser.Reference{
					Link:  "http://www.ref1.com/",
					Title: "Reference 1"}, true
			case "ref2":
				// overridden exists and reference defined
				return &parser.Reference{
					Link:  "http://www.overridden.com/",
					Title: "Reference Overridden"}, true
			case "ref3":
				// not overridden and reference defined
				return nil, false
			case "ref4":
				// overridden missing and defined
				return nil, true
			case "!(*http.ServeMux).ServeHTTP":
				return &parser.Reference{
					Link:  "http://localhost:6060/pkg/net/http/#ServeMux.ServeHTTP",
					Title: "ServeHTTP docs"}, true
			case "ref5":
				return &parser.Reference{
					Link:  "http://www.ref5.com/",
					Title: "Reference 5",
					Text:  "Moo",
				}, true
			case "ref6":
				return &parser.Reference{
					Link:  "http://www.ref6.com/",
					Title: "Reference 6",
					Text:  "Moo",
				}, true
			}
			return nil, false
		},
	})
}

func TestStrong(t *testing.T) {
	var tests = []string{
		"nothing inline\n",
		"<p>nothing inline</p>\n",

		"simple **inline** test\n",
		"<p>simple <strong>inline</strong> test</p>\n",

		"**at the** beginning\n",
		"<p><strong>at the</strong> beginning</p>\n",

		"at the **end**\n",
		"<p>at the <strong>end</strong></p>\n",

		"**try two** in **one line**\n",
		"<p><strong>try two</strong> in <strong>one line</strong></p>\n",

		"over **two\nlines** test\n",
		"<p>over <strong>two\nlines</strong> test</p>\n",

		"odd **number of** markers** here\n",
		"<p>odd <strong>number of</strong> markers** here</p>\n",

		"odd **number\nof** markers** here\n",
		"<p>odd <strong>number\nof</strong> markers** here</p>\n",

		"simple __inline__ test\n",
		"<p>simple <strong>inline</strong> test</p>\n",

		"__at the__ beginning\n",
		"<p><strong>at the</strong> beginning</p>\n",

		"at the __end__\n",
		"<p>at the <strong>end</strong></p>\n",

		"__try two__ in __one line__\n",
		"<p><strong>try two</strong> in <strong>one line</strong></p>\n",

		"over __two\nlines__ test\n",
		"<p>over <strong>two\nlines</strong> test</p>\n",

		"odd __number of__ markers__ here\n",
		"<p>odd <strong>number of</strong> markers__ here</p>\n",

		"odd __number\nof__ markers__ here\n",
		"<p>odd <strong>number\nof</strong> markers__ here</p>\n",

		"mix of **markers__\n",
		"<p>mix of **markers__</p>\n",

		"**`/usr`** : this folder is named `usr`\n",
		"<p><strong><code>/usr</code></strong> : this folder is named <code>usr</code></p>\n",

		"**`/usr`** :\n\n this folder is named `usr`\n",
		"<p><strong><code>/usr</code></strong> :</p>\n\n<p>this folder is named <code>usr</code></p>\n",
	}
	doTestsInline(t, tests)
}

func TestStrongShort(t *testing.T) {
	var tests = []string{
		"**`/usr`** :\n\n this folder is named `usr`\n",
		"<p><strong><code>/usr</code></strong> :</p>\n\n<p>this folder is named <code>usr</code></p>\n",
	}
	doTestsInline(t, tests)

}
func TestEmphasisMix(t *testing.T) {
	var tests = []string{
		"***triple emphasis***\n",
		"<p><strong><em>triple emphasis</em></strong></p>\n",

		"***triple\nemphasis***\n",
		"<p><strong><em>triple\nemphasis</em></strong></p>\n",

		"___triple emphasis___\n",
		"<p><strong><em>triple emphasis</em></strong></p>\n",

		"***triple emphasis___\n",
		"<p>***triple emphasis___</p>\n",

		"*italics **and bold** end*\n",
		"<p><em>italics <strong>and bold</strong> end</em></p>\n",

		"*italics **and bold***\n",
		"<p><em>italics <strong>and bold</strong></em></p>\n",

		"***bold** and italics*\n",
		"<p><em><strong>bold</strong> and italics</em></p>\n",

		"*start **bold** and italics*\n",
		"<p><em>start <strong>bold</strong> and italics</em></p>\n",

		"*__triple emphasis__*\n",
		"<p><em><strong>triple emphasis</strong></em></p>\n",

		"__*triple emphasis*__\n",
		"<p><strong><em>triple emphasis</em></strong></p>\n",

		"**improper *nesting** is* bad\n",
		"<p><strong>improper *nesting</strong> is* bad</p>\n",

		"*improper **nesting* is** bad\n",
		"<p><em>improper **nesting</em> is** bad</p>\n",
	}
	doTestsInline(t, tests)
}

func TestEmphasisLink(t *testing.T) {
	var tests = []string{
		"[first](before) *text[second] (inside)text* [third](after)\n",
		"<p><a href=\"before\">first</a> <em>text<a href=\"inside\">second</a>text</em> <a href=\"after\">third</a></p>\n",

		"*incomplete [link] definition*\n",
		"<p><em>incomplete [link] definition</em></p>\n",

		"*it's [emphasis*] (not link)\n",
		"<p><em>it's [emphasis</em>] (not link)</p>\n",

		"*it's [emphasis*] and *[asterisk]\n",
		"<p><em>it's [emphasis</em>] and *[asterisk]</p>\n",
	}
	doTestsInline(t, tests)
}

func TestStrikeThrough(t *testing.T) {
	var tests = []string{
		"nothing inline\n",
		"<p>nothing inline</p>\n",

		"simple ~~inline~~ test\n",
		"<p>simple <del>inline</del> test</p>\n",

		"~~at the~~ beginning\n",
		"<p><del>at the</del> beginning</p>\n",

		"at the ~~end~~\n",
		"<p>at the <del>end</del></p>\n",

		"~~try two~~ in ~~one line~~\n",
		"<p><del>try two</del> in <del>one line</del></p>\n",

		"over ~~two\nlines~~ test\n",
		"<p>over <del>two\nlines</del> test</p>\n",

		"odd ~~number of~~ markers~~ here\n",
		"<p>odd <del>number of</del> markers~~ here</p>\n",

		"odd ~~number\nof~~ markers~~ here\n",
		"<p>odd <del>number\nof</del> markers~~ here</p>\n",
	}
	doTestsInline(t, tests)
}

func TestCodeSpan(t *testing.T) {
	var tests = []string{
		"`source code`\n",
		"<p><code>source code</code></p>\n",

		"` source code with spaces `\n",
		"<p><code>source code with spaces</code></p>\n",

		"` source code with spaces `not here\n",
		"<p><code>source code with spaces</code>not here</p>\n",

		"a `single marker\n",
		"<p>a `single marker</p>\n",

		"a single multi-tick marker with ``` no text\n",
		"<p>a single multi-tick marker with ``` no text</p>\n",

		"markers with ` ` a space\n",
		"<p>markers with  a space</p>\n",

		"`source code` and a `stray\n",
		"<p><code>source code</code> and a `stray</p>\n",

		"`source *with* _awkward characters_ in it`\n",
		"<p><code>source *with* _awkward characters_ in it</code></p>\n",

		"`split over\ntwo lines`\n",
		"<p><code>split over\ntwo lines</code></p>\n",

		"```multiple ticks``` for the marker\n",
		"<p><code>multiple ticks</code> for the marker</p>\n",

		"```multiple ticks `with` ticks inside```\n",
		"<p><code>multiple ticks `with` ticks inside</code></p>\n",

		"`@param {(string|number)} n`",
		"<p><code>@param {(string|number)} n</code></p>\n",
	}
	doTestsInlineParam(t, tests, TestParams{})

	tp := TestParams{
		extensions: parser.Mmark,
	}
	doTestsInlineParam(t, tests, tp)
}

func TestCodeBlock(t *testing.T) {
	var tests = []string{
		"1. This is an item\n" +
			"   ```java\n" +
			"   int a = 1;\n" +
			"   ```\n" +
			"1. This is another item\n",
		"<ol>\n<li>This is an item\n\n<pre><code class=\"language-java\">\nint a = 1;\n</code></pre>\n</li>\n<li>This is another item</li>\n</ol>\n",
	}
	doTestsInlineParam(t, tests, TestParams{
		extensions: parser.FencedCode,
	})
}

func TestLineBreak(t *testing.T) {
	var tests = []string{
		"this line  \nhas a break\n",
		"<p>this line<br />\nhas a break</p>\n",

		"this line \ndoes not\n",
		"<p>this line\ndoes not</p>\n",

		"this line\\\ndoes not\n",
		"<p>this line\\\ndoes not</p>\n",

		"this line\\ \ndoes not\n",
		"<p>this line\\\ndoes not</p>\n",

		"this has an   \nextra space\n",
		"<p>this has an<br />\nextra space</p>\n",
	}
	doTestsInline(t, tests)

	tests = []string{
		"this line  \nhas a break\n",
		"<p>this line<br />\nhas a break</p>\n",

		"this line \ndoes not\n",
		"<p>this line\ndoes not</p>\n",

		"this line\\\nhas a break\n",
		"<p>this line<br />\nhas a break</p>\n",

		"this line\\ \ndoes not\n",
		"<p>this line\\\ndoes not</p>\n",

		"this has an   \nextra space\n",
		"<p>this has an<br />\nextra space</p>\n",
	}
	doTestsInlineParam(t, tests, TestParams{
		extensions: parser.BackslashLineBreak})
}

func TestInlineLink(t *testing.T) {
	var tests = []string{
		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",

		"[foo with a title](/bar/ \"title\")\n",
		"<p><a href=\"/bar/\" title=\"title\">foo with a title</a></p>\n",

		"[foo with a title](/bar/\t\"title\")\n",
		"<p><a href=\"/bar/\" title=\"title\">foo with a title</a></p>\n",

		"[foo with a title](/bar/ \"title\"  )\n",
		"<p><a href=\"/bar/\" title=\"title\">foo with a title</a></p>\n",

		"[foo with a title](/bar/ title with no quotes)\n",
		"<p><a href=\"/bar/ title with no quotes\">foo with a title</a></p>\n",

		"[foo]()\n",
		"<p><a href=\"\">foo</a></p>\n",

		"![foo](/bar/)\n",
		"<p><img src=\"/bar/\" alt=\"foo\" /></p>\n",

		"![foo with a title](/bar/ \"title\")\n",
		"<p><img src=\"/bar/\" alt=\"foo with a title\" title=\"title\" /></p>\n",

		"![foo with a title](/bar/\t\"title\")\n",
		"<p><img src=\"/bar/\" alt=\"foo with a title\" title=\"title\" /></p>\n",

		"![foo with a title](/bar/ \"title\"  )\n",
		"<p><img src=\"/bar/\" alt=\"foo with a title\" title=\"title\" /></p>\n",

		"![foo with a title](/bar/ title with no quotes)\n",
		"<p><img src=\"/bar/ title with no quotes\" alt=\"foo with a title\" /></p>\n",

		"![](img.jpg)\n",
		"<p><img src=\"img.jpg\" alt=\"\" /></p>\n",

		"[link](url)\n",
		"<p><a href=\"url\">link</a></p>\n",

		"![foo]()\n",
		"<p><img src=\"\" alt=\"foo\" /></p>\n",

		"[a link]\t(/with_a_tab/)\n",
		"<p><a href=\"/with_a_tab/\">a link</a></p>\n",

		"[a link]  (/with_spaces/)\n",
		"<p><a href=\"/with_spaces/\">a link</a></p>\n",

		"[text (with) [[nested] (brackets)]](/url/)\n",
		"<p><a href=\"/url/\">text (with) [[nested] (brackets)]</a></p>\n",

		"[text (with) [broken nested] (brackets)]](/url/)\n",
		"<p>[text (with) <a href=\"brackets\">broken nested</a>]](/url/)</p>\n",

		"[text\nwith a newline](/link/)\n",
		"<p><a href=\"/link/\">text\nwith a newline</a></p>\n",

		"[text in brackets] [followed](/by a link/)\n",
		"<p>[text in brackets] <a href=\"/by a link/\">followed</a></p>\n",

		"[link with\\] a closing bracket](/url/)\n",
		"<p><a href=\"/url/\">link with] a closing bracket</a></p>\n",

		"[link with\\[ an opening bracket](/url/)\n",
		"<p><a href=\"/url/\">link with[ an opening bracket</a></p>\n",

		"[link with\\) a closing paren](/url/)\n",
		"<p><a href=\"/url/\">link with) a closing paren</a></p>\n",

		"[link with\\( an opening paren](/url/)\n",
		"<p><a href=\"/url/\">link with( an opening paren</a></p>\n",

		"[link](  with whitespace)\n",
		"<p><a href=\"with whitespace\">link</a></p>\n",

		"[link](  with whitespace   )\n",
		"<p><a href=\"with whitespace\">link</a></p>\n",

		"[![image](someimage)](with image)\n",
		"<p><a href=\"with image\"><img src=\"someimage\" alt=\"image\" /></a></p>\n",

		"[link](url \"one quote - broken markdown)\n",
		"<p>[link](url &quot;one quote - broken markdown)</p>\n",

		"[link](url 'one quote - broken markdown)\n",
		"<p>[link](url 'one quote - broken markdown)</p>\n",

		"[link](<url>)\n",
		"<p><a href=\"url\">link</a></p>\n",

		"[link & ampersand](/url/)\n",
		"<p><a href=\"/url/\">link &amp; ampersand</a></p>\n",

		"[link &amp; ampersand](/url/)\n",
		"<p><a href=\"/url/\">link &amp; ampersand</a></p>\n",

		"[link](/url/&query)\n",
		"<p><a href=\"/url/&amp;query\">link</a></p>\n",

		"[[t]](/t)\n",
		"<p><a href=\"/t\">[t]</a></p>\n",

		"[link](</>)\n",
		"<p><a href=\"/\">link</a></p>\n",

		"[link](<./>)\n",
		"<p><a href=\"./\">link</a></p>\n",

		"[link](<../>)\n",
		"<p><a href=\"../\">link</a></p>\n",

		"[disambiguation](http://en.wikipedia.org/wiki/Disambiguation_(disambiguation))",
		"<p><a href=\"http://en.wikipedia.org/wiki/Disambiguation_(disambiguation)\">disambiguation</a></p>\n",
	}
	doLinkTestsInline(t, tests)

}

func TestRelAttrLink(t *testing.T) {
	var nofollowTests = []string{
		"[foo](http://bar.com/foo/)\n",
		"<p><a href=\"http://bar.com/foo/\" rel=\"nofollow\">foo</a></p>\n",

		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",

		"[foo](/)\n",
		"<p><a href=\"/\">foo</a></p>\n",

		"[foo](./)\n",
		"<p><a href=\"./\">foo</a></p>\n",

		"[foo](../)\n",
		"<p><a href=\"../\">foo</a></p>\n",

		"[foo](../bar)\n",
		"<p><a href=\"../bar\">foo</a></p>\n",
	}
	doTestsInlineParam(t, nofollowTests, TestParams{
		Flags: html.Safelink | html.NofollowLinks,
	})

	var noreferrerTests = []string{
		"[foo](http://bar.com/foo/)\n",
		"<p><a href=\"http://bar.com/foo/\" rel=\"noreferrer\">foo</a></p>\n",

		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",
	}
	doTestsInlineParam(t, noreferrerTests, TestParams{
		Flags: html.Safelink | html.NoreferrerLinks,
	})

	var nofollownoreferrerTests = []string{
		"[foo](http://bar.com/foo/)\n",
		"<p><a href=\"http://bar.com/foo/\" rel=\"nofollow noreferrer\">foo</a></p>\n",

		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",
	}
	doTestsInlineParam(t, nofollownoreferrerTests, TestParams{
		Flags: html.Safelink | html.NofollowLinks | html.NoreferrerLinks,
	})

	var noopenernoreferrerTests = []string{
		"[foo](http://bar.com/foo/)\n",
		"<p><a href=\"http://bar.com/foo/\" rel=\"noreferrer noopener\">foo</a></p>\n",

		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",
	}
	doTestsInlineParam(t, noopenernoreferrerTests, TestParams{
		Flags: html.Safelink | html.NoopenerLinks | html.NoreferrerLinks,
	})
}

func TestHrefTargetBlank(t *testing.T) {
	var tests = []string{
		// internal link
		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",

		"[foo](/)\n",
		"<p><a href=\"/\">foo</a></p>\n",

		"[foo](./)\n",
		"<p><a href=\"./\">foo</a></p>\n",

		"[foo](./bar)\n",
		"<p><a href=\"./bar\">foo</a></p>\n",

		"[foo](../)\n",
		"<p><a href=\"../\">foo</a></p>\n",

		"[foo](../bar)\n",
		"<p><a href=\"../bar\">foo</a></p>\n",

		"[foo](http://example.com)\n",
		"<p><a href=\"http://example.com\" target=\"_blank\">foo</a></p>\n",
	}
	doTestsInlineParam(t, tests, TestParams{
		Flags: html.Safelink | html.HrefTargetBlank,
	})
}

func TestSafeInlineLink(t *testing.T) {
	var tests = []string{
		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",

		"[foo](/)\n",
		"<p><a href=\"/\">foo</a></p>\n",

		"[foo](./)\n",
		"<p><a href=\"./\">foo</a></p>\n",

		"[foo](../)\n",
		"<p><a href=\"../\">foo</a></p>\n",

		"[foo](http://bar/)\n",
		"<p><a href=\"http://bar/\">foo</a></p>\n",

		"[foo](https://bar/)\n",
		"<p><a href=\"https://bar/\">foo</a></p>\n",

		"[foo](ftp://bar/)\n",
		"<p><a href=\"ftp://bar/\">foo</a></p>\n",

		"[foo](mailto:bar/)\n",
		"<p><a href=\"mailto:bar/\">foo</a></p>\n",

		"[foo](monero:4AfUP827TeRZ1cck3tZThgZbRCEwBrpcJTkA1LCiyFVuMH4b5y59bKMZHGb9y58K3gSjWDCBsB4RkGsGDhsmMG5R2qmbLeW)\n",
		"<p><a href=\"monero:4AfUP827TeRZ1cck3tZThgZbRCEwBrpcJTkA1LCiyFVuMH4b5y59bKMZHGb9y58K3gSjWDCBsB4RkGsGDhsmMG5R2qmbLeW\">foo</a></p>\n",

		"[foo](bitcoin:bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh)\n",
		"<p><a href=\"bitcoin:bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh\">foo</a></p>\n",

		// Not considered safe
		"[foo](baz://bar/)\n",
		"<p><tt>foo</tt></p>\n",
	}
	doSafeTestsInline(t, tests)
}

func TestReferenceLink(t *testing.T) {
	var tests = []string{
		"[link][ref]\n",
		"<p>[link][ref]</p>\n",

		"[link][ref]\n   [ref]: /url/ \"title\"\n",
		"<p><a href=\"/url/\" title=\"title\">link</a></p>\n",

		"[link][ref]\n   [ref]: /url/\n",
		"<p><a href=\"/url/\">link</a></p>\n",

		"   [ref]: /url/\n",
		"",

		"   [ref]: /url/\n[ref2]: /url/\n [ref3]: /url/\n",
		"",

		"   [ref]: /url/\n[ref2]: /url/\n [ref3]: /url/\n    [4spaces]: /url/\n",
		"<pre><code>[4spaces]: /url/\n</code></pre>\n",

		"[hmm](ref2)\n   [ref]: /url/\n[ref2]: /url/\n [ref3]: /url/\n",
		"<p><a href=\"ref2\">hmm</a></p>\n",

		"[ref]\n",
		"<p>[ref]</p>\n",

		"[ref]\n   [ref]: /url/ \"title\"\n",
		"<p><a href=\"/url/\" title=\"title\">ref</a></p>\n",

		"[ref]\n   [ref]: ../url/ \"title\"\n",
		"<p><a href=\"../url/\" title=\"title\">ref</a></p>\n",

		"[link][ref]\n   [ref]: /url/",
		"<p><a href=\"/url/\">link</a></p>\n",
	}
	doLinkTestsInline(t, tests)
}

func parseTestCases(s string) []string {
	parts := strings.Split(s, "+++")
	if len(parts)%2 != 0 {
		panic("odd test tuples")
	}
	for i, s := range parts {
		parts[i] = strings.TrimLeft(s, "\n")
	}
	return parts
}

var testTagsTestCases = `a <span>tag</span>
+++
<p>a <span>tag</span></p>
+++
<span>tag</span>
+++
<p><span>tag</span></p>
+++
<span>mismatch</spandex>
+++
<p><span>mismatch</spandex></p>
+++
a <singleton /> tag
+++
<p>a <singleton /> tag</p>
`

func TestParseTestCases(t *testing.T) {
	var exp = []string{
		"a <span>tag</span>\n",
		"<p>a <span>tag</span></p>\n",

		"<span>tag</span>\n",
		"<p><span>tag</span></p>\n",

		"<span>mismatch</spandex>\n",
		"<p><span>mismatch</spandex></p>\n",

		"a <singleton /> tag\n",
		"<p>a <singleton /> tag</p>\n",
	}
	got := parseTestCases(testTagsTestCases)
	for i, sGot := range got {
		sExp := exp[i]
		if sExp != sGot {
			t.Errorf("\nExpected[%#v]\nGot     [%#v]\n", sExp, sGot)
		}
	}
}

func TestTags(t *testing.T) {
	tests := parseTestCases(testTagsTestCases)
	doTestsInline(t, tests)
}

func TestAutoLink(t *testing.T) {
	var tests = []string{
		"http://foo.com/\n",
		"<p><a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"1 http://foo.com/\n",
		"<p>1 <a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"1http://foo.com/\n",
		"<p>1<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"1.http://foo.com/\n",
		"<p>1.<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"1. http://foo.com/\n",
		"<ol>\n<li><a href=\"http://foo.com/\">http://foo.com/</a></li>\n</ol>\n",

		"-http://foo.com/\n",
		"<p>-<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"- http://foo.com/\n",
		"<ul>\n<li><a href=\"http://foo.com/\">http://foo.com/</a></li>\n</ul>\n",

		"_http://foo.com/\n",
		"<p>_<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"令狐http://foo.com/\n",
		"<p>令狐<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"令狐 http://foo.com/\n",
		"<p>令狐 <a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"ahttp://foo.com/\n",
		"<p>ahttp://foo.com/</p>\n",

		">http://foo.com/\n",
		"<blockquote>\n<p><a href=\"http://foo.com/\">http://foo.com/</a></p>\n</blockquote>\n",

		"> http://foo.com/\n",
		"<blockquote>\n<p><a href=\"http://foo.com/\">http://foo.com/</a></p>\n</blockquote>\n",

		"go to <http://foo.com/>\n",
		"<p>go to <a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"a secure <https://link.org>\n",
		"<p>a secure <a href=\"https://link.org\">https://link.org</a></p>\n",

		"an email <mailto:some@one.com>\n",
		"<p>an email <a href=\"mailto:some@one.com\">some@one.com</a></p>\n",

		"an email <mailto:some@one.com>\n",
		"<p>an email <a href=\"mailto:some@one.com\">some@one.com</a></p>\n",

		"an email <some@one.com>\n",
		"<p>an email <a href=\"mailto:some@one.com\">some@one.com</a></p>\n",

		"an ftp <ftp://old.com>\n",
		"<p>an ftp <a href=\"ftp://old.com\">ftp://old.com</a></p>\n",

		"an ftp <ftp:old.com>\n",
		"<p>an ftp <a href=\"ftp:old.com\">ftp:old.com</a></p>\n",

		"a link with <http://new.com?query=foo&bar>\n",
		"<p>a link with <a href=\"http://new.com?query=foo&amp;bar\">" +
			"http://new.com?query=foo&amp;bar</a></p>\n",

		"quotes mean a tag <http://new.com?query=\"foo\"&bar>\n",
		"<p>quotes mean a tag <http://new.com?query=\"foo\"&bar></p>\n",

		"quotes mean a tag <http://new.com?query='foo'&bar>\n",
		"<p>quotes mean a tag <http://new.com?query='foo'&bar></p>\n",

		"unless escaped <http://new.com?query=\\\"foo\\\"&bar>\n",
		"<p>unless escaped <a href=\"http://new.com?query=&quot;foo&quot;&amp;bar\">" +
			"http://new.com?query=&quot;foo&quot;&amp;bar</a></p>\n",

		"even a > can be escaped <http://new.com?q=\\>&etc>\n",
		"<p>even a &gt; can be escaped <a href=\"http://new.com?q=&gt;&amp;etc\">" +
			"http://new.com?q=&gt;&amp;etc</a></p>\n",

		"<a href=\"http://fancy.com\">http://fancy.com</a>\n",
		"<p><a href=\"http://fancy.com\">http://fancy.com</a></p>\n",

		"<a href=\"http://fancy.com\">This is a link</a>\n",
		"<p><a href=\"http://fancy.com\">This is a link</a></p>\n",

		"<a href=\"http://www.fancy.com/A_B.pdf\">http://www.fancy.com/A_B.pdf</a>\n",
		"<p><a href=\"http://www.fancy.com/A_B.pdf\">http://www.fancy.com/A_B.pdf</a></p>\n",

		"(<a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a> (\n",
		"<p>(<a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a> (</p>\n",

		"(<a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a> (part two: <a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a>)).\n",
		"<p>(<a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a> (part two: <a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a>)).</p>\n",

		"http://www.foo.com<br />\n",
		"<p><a href=\"http://www.foo.com\">http://www.foo.com</a><br /></p>\n",

		"http://foo.com/viewtopic.php?f=18&amp;t=297",
		"<p><a href=\"http://foo.com/viewtopic.php?f=18&amp;t=297\">http://foo.com/viewtopic.php?f=18&amp;t=297</a></p>\n",

		"http://foo.com/viewtopic.php?param=&quot;18&quot;zz",
		"<p><a href=\"http://foo.com/viewtopic.php?param=&quot;18&quot;zz\">http://foo.com/viewtopic.php?param=&quot;18&quot;zz</a></p>\n",

		"http://foo.com/viewtopic.php?param=&quot;18&quot;",
		"<p><a href=\"http://foo.com/viewtopic.php?param=&quot;18&quot;\">http://foo.com/viewtopic.php?param=&quot;18&quot;</a></p>\n",

		"<a href=\"https://fancy.com\">https://fancy.com</a>\n",
		"<p><a href=\"https://fancy.com\">https://fancy.com</a></p>\n",
	}
	doLinkTestsInline(t, tests)
}

var footnoteTests = []string{
	"testing footnotes.[^a]\n\n[^a]: This is the note\n",
	`<p>testing footnotes.<sup class="footnote-ref" id="fnref:a"><a href="#fn:a">1</a></sup></p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:a">This is the note</li>
</ol>

</div>
`,

	`testing long[^b] notes.

[^b]: Paragraph 1

	Paragraph 2

	` + "```\n\tsome code\n\t```" + `

	Paragraph 3

No longer in the footnote
`,
	`<p>testing long<sup class="footnote-ref" id="fnref:b"><a href="#fn:b">1</a></sup> notes.</p>

<p>No longer in the footnote</p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:b"><p>Paragraph 1</p>

<p>Paragraph 2</p>

<p>
<pre><code>
some code
</code></pre>
</p>

<p>Paragraph 3</p></li>
</ol>

</div>
`,

	`testing[^c] multiple[^d] notes.

[^c]: this is [note] c


omg

[^d]: this is note d

what happens here

[note]: /link/c

`,
	`<p>testing<sup class="footnote-ref" id="fnref:c"><a href="#fn:c">1</a></sup> multiple<sup class="footnote-ref" id="fnref:d"><a href="#fn:d">2</a></sup> notes.</p>

<p>omg</p>

<p>what happens here</p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:c">this is <a href="/link/c">note</a> c</li>

<li id="fn:d">this is note d</li>
</ol>

</div>
`,

	"testing inline^[this is the note] notes.\n",
	`<p>testing inline<sup class="footnote-ref" id="fnref:this-is-the-note"><a href="#fn:this-is-the-note">1</a></sup> notes.</p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:this-is-the-note">this is the note</li>
</ol>

</div>
`,

	"testing multiple[^1] types^[inline note] of notes[^2]\n\n[^2]: the second deferred note\n[^1]: the first deferred note\n\n\twhich happens to be a block\n",
	`<p>testing multiple<sup class="footnote-ref" id="fnref:1"><a href="#fn:1">1</a></sup> types<sup class="footnote-ref" id="fnref:inline-note"><a href="#fn:inline-note">2</a></sup> of notes<sup class="footnote-ref" id="fnref:2"><a href="#fn:2">3</a></sup></p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:1"><p>the first deferred note</p>

<p>which happens to be a block</p></li>

<li id="fn:inline-note">inline note</li>

<li id="fn:2">the second deferred note</li>
</ol>

</div>
`,

	`This is a footnote[^1]^[and this is an inline footnote]

[^1]: the footnote text.

    may be multiple paragraphs.
`,
	`<p>This is a footnote<sup class="footnote-ref" id="fnref:1"><a href="#fn:1">1</a></sup><sup class="footnote-ref" id="fnref:and-this-is-an-i"><a href="#fn:and-this-is-an-i">2</a></sup></p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:1"><p>the footnote text.</p>

<p>may be multiple paragraphs.</p></li>

<li id="fn:and-this-is-an-i">and this is an inline footnote</li>
</ol>

</div>
`,

	"empty footnote[^]\n\n[^]: fn text",
	"<p>empty footnote<sup class=\"footnote-ref\" id=\"fnref:\"><a href=\"#fn:\">1</a></sup></p>\n\n<div class=\"footnotes\">\n\n<hr />\n\n<ol>\n<li id=\"fn:\">fn text</li>\n</ol>\n\n</div>\n",

	"Some text.[^note1]\n\n[^note1]: fn1",
	"<p>Some text.<sup class=\"footnote-ref\" id=\"fnref:note1\"><a href=\"#fn:note1\">1</a></sup></p>\n\n<div class=\"footnotes\">\n\n<hr />\n\n<ol>\n<li id=\"fn:note1\">fn1</li>\n</ol>\n\n</div>\n",

	"Some text.[^note1][^note2]\n\n[^note1]: fn1\n[^note2]: fn2\n",
	"<p>Some text.<sup class=\"footnote-ref\" id=\"fnref:note1\"><a href=\"#fn:note1\">1</a></sup><sup class=\"footnote-ref\" id=\"fnref:note2\"><a href=\"#fn:note2\">2</a></sup></p>\n\n<div class=\"footnotes\">\n\n<hr />\n\n<ol>\n<li id=\"fn:note1\">fn1</li>\n\n<li id=\"fn:note2\">fn2</li>\n</ol>\n\n</div>\n",

	`Bla bla [^1] [WWW][w3]

[^1]: This is a footnote

[w3]: http://www.w3.org/
`,
	`<p>Bla bla <sup class="footnote-ref" id="fnref:1"><a href="#fn:1">1</a></sup> <a href="http://www.w3.org/">WWW</a></p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:1">This is a footnote</li>
</ol>

</div>
`,

	`This is exciting![^fn1]

[^fn1]: Fine print
`,
	`<p>This is exciting!<sup class="footnote-ref" id="fnref:fn1"><a href="#fn:fn1">1</a></sup></p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:fn1">Fine print</li>
</ol>

</div>
`,

	`This text does not reference a footnote.

[^footnote]: But it has a footnote! And it gets omitted.
`,
	"<p>This text does not reference a footnote.</p>\n",

	`testing footnotes.[^a]
 test footnotes the second.[^b]
 [^a]: This is the first note[^a].
[^b]: this is the second note.[^a]
`,
	`<p>testing footnotes.<sup class="footnote-ref" id="fnref:a"><a href="#fn:a">1</a></sup>
 test footnotes the second.<sup class="footnote-ref" id="fnref:b"><a href="#fn:b">2</a></sup></p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:a">This is the first note<sup class="footnote-ref" id="fnref:a"><a href="#fn:a">1</a></sup>.</li>

<li id="fn:b">this is the second note.<sup class="footnote-ref" id="fnref:a"><a href="#fn:a">1</a></sup></li>
</ol>

</div>
`,
}

func TestFootnotes(t *testing.T) {
	doTestsInlineParam(t, footnoteTests, TestParams{
		extensions: parser.Footnotes,
	})
}

func TestFootnotesWithParameters(t *testing.T) {
	tests := make([]string, len(footnoteTests))

	prefix := "testPrefix"
	returnText := "ret"
	re := regexp.MustCompile(`(?ms)<li id="fn:(\S+?)">(.*?)</li>`)

	// Transform the test expectations to match the parameters we're using.
	for i, test := range footnoteTests {
		if i%2 == 1 {
			test = strings.Replace(test, "fn:", "fn:"+prefix, -1)
			test = strings.Replace(test, "fnref:", "fnref:"+prefix, -1)
			test = re.ReplaceAllString(test, `<li id="fn:$1">$2 <a class="footnote-return" href="#fnref:$1">ret</a></li>`)
		}
		tests[i] = test
	}

	params := html.RendererOptions{
		FootnoteAnchorPrefix:       prefix,
		FootnoteReturnLinkContents: returnText,
	}

	doTestsInlineParam(t, tests, TestParams{
		extensions:      parser.Footnotes,
		Flags:           html.FootnoteReturnLinks,
		RendererOptions: params,
	})
}

func TestNestedFootnotes(t *testing.T) {
	var tests = []string{
		`Paragraph.[^fn1]

[^fn1]:
  Asterisk[^fn2]

[^fn2]:
  Obelisk`,
		`<p>Paragraph.<sup class="footnote-ref" id="fnref:fn1"><a href="#fn:fn1">1</a></sup></p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:fn1">Asterisk<sup class="footnote-ref" id="fnref:fn2"><a href="#fn:fn2">2</a></sup></li>

<li id="fn:fn2">Obelisk</li>
</ol>

</div>
`,
		`This uses footnote A.[^A]

This uses footnote C.[^C]

[^A]:
  A note. use itself.[^A]
[^B]:
  B note, uses A to test duplicate.[^A]
[^C]:
  C note, uses B.[^B]
`,

		`<p>This uses footnote A.<sup class="footnote-ref" id="fnref:A"><a href="#fn:A">1</a></sup></p>

<p>This uses footnote C.<sup class="footnote-ref" id="fnref:C"><a href="#fn:C">2</a></sup></p>

<div class="footnotes">

<hr />

<ol>
<li id="fn:A">A note. use itself.<sup class="footnote-ref" id="fnref:A"><a href="#fn:A">1</a></sup></li>

<li id="fn:C">C note, uses B.<sup class="footnote-ref" id="fnref:B"><a href="#fn:B">3</a></sup></li>

<li id="fn:B">B note, uses A to test duplicate.<sup class="footnote-ref" id="fnref:A"><a href="#fn:A">1</a></sup></li>
</ol>

</div>
`,
	}
	doTestsInlineParam(t, tests, TestParams{extensions: parser.Footnotes})
}

func TestInlineComments(t *testing.T) {
	var tests = []string{
		"Hello <!-- there ->\n",
		"<p>Hello &lt;!&mdash; there &ndash;&gt;</p>\n",

		"Hello <!-- there -->\n",
		"<p>Hello <!-- there --></p>\n",

		"Hello <!-- there -->",
		"<p>Hello <!-- there --></p>\n",

		"Hello <!---->\n",
		"<p>Hello <!----></p>\n",

		"Hello <!-- there -->\na",
		"<p>Hello <!-- there -->\na</p>\n",

		"* list <!-- item -->\n",
		"<ul>\n<li>list <!-- item --></li>\n</ul>\n",

		"<!-- Front --> comment\n",
		"<p><!-- Front --> comment</p>\n",

		"blahblah\n<!--- foo -->\nrhubarb\n",
		"<p>blahblah\n<!--- foo -->\nrhubarb</p>\n",
	}
	doTestsInlineParam(t, tests, TestParams{Flags: html.Smartypants | html.SmartypantsDashes})
}

func TestSmartDoubleQuotes(t *testing.T) {
	var tests = []string{
		"this should be normal \"quoted\" text.\n",
		"<p>this should be normal &ldquo;quoted&rdquo; text.</p>\n",
		"this \" single double\n",
		"<p>this &ldquo; single double</p>\n",
		"two pair of \"some\" quoted \"text\".\n",
		"<p>two pair of &ldquo;some&rdquo; quoted &ldquo;text&rdquo;.</p>\n"}

	doTestsInlineParam(t, tests, TestParams{Flags: html.Smartypants})
}

func TestSmartDoubleQuotesNBSP(t *testing.T) {
	var tests = []string{
		"this should be normal \"quoted\" text.\n",
		"<p>this should be normal &ldquo;&nbsp;quoted&nbsp;&rdquo; text.</p>\n",
		"this \" single double\n",
		"<p>this &ldquo;&nbsp; single double</p>\n",
		"two pair of \"some\" quoted \"text\".\n",
		"<p>two pair of &ldquo;&nbsp;some&nbsp;&rdquo; quoted &ldquo;&nbsp;text&nbsp;&rdquo;.</p>\n"}

	doTestsInlineParam(t, tests, TestParams{Flags: html.Smartypants | html.SmartypantsQuotesNBSP})
}

func TestSmartAngledDoubleQuotes(t *testing.T) {
	var tests = []string{
		"this should be angled \"quoted\" text.\n",
		"<p>this should be angled &laquo;quoted&raquo; text.</p>\n",
		"this \" single double\n",
		"<p>this &laquo; single double</p>\n",
		"two pair of \"some\" quoted \"text\".\n",
		"<p>two pair of &laquo;some&raquo; quoted &laquo;text&raquo;.</p>\n"}

	doTestsInlineParam(t, tests, TestParams{Flags: html.Smartypants | html.SmartypantsAngledQuotes})
}

func TestSmartAngledDoubleQuotesNBSP(t *testing.T) {
	var tests = []string{
		"this should be angled \"quoted\" text.\n",
		"<p>this should be angled &laquo;&nbsp;quoted&nbsp;&raquo; text.</p>\n",
		"this \" single double\n",
		"<p>this &laquo;&nbsp; single double</p>\n",
		"two pair of \"some\" quoted \"text\".\n",
		"<p>two pair of &laquo;&nbsp;some&nbsp;&raquo; quoted &laquo;&nbsp;text&nbsp;&raquo;.</p>\n"}

	doTestsInlineParam(t, tests, TestParams{Flags: html.Smartypants | html.SmartypantsAngledQuotes | html.SmartypantsQuotesNBSP})
}

func TestSmartFractions(t *testing.T) {
	var tests = []string{
		"1/2, 1/4 and 3/4; 1/4th and 3/4ths\n",
		"<p>&frac12;, &frac14; and &frac34;; &frac14;th and &frac34;ths</p>\n",
		"1/2/2015, 1/4/2015, 3/4/2015; 2015/1/2, 2015/1/4, 2015/3/4.\n",
		"<p>1/2/2015, 1/4/2015, 3/4/2015; 2015/1/2, 2015/1/4, 2015/3/4.</p>\n"}

	doTestsInlineParam(t, tests, TestParams{Flags: html.Smartypants})

	tests = []string{
		"1/2, 2/3, 81/100 and 1000000/1048576.\n",
		"<p><sup>1</sup>&frasl;<sub>2</sub>, <sup>2</sup>&frasl;<sub>3</sub>, <sup>81</sup>&frasl;<sub>100</sub> and <sup>1000000</sup>&frasl;<sub>1048576</sub>.</p>\n",
		"1/2/2015, 1/4/2015, 3/4/2015; 2015/1/2, 2015/1/4, 2015/3/4.\n",
		"<p>1/2/2015, 1/4/2015, 3/4/2015; 2015/1/2, 2015/1/4, 2015/3/4.</p>\n"}

	doTestsInlineParam(t, tests, TestParams{Flags: html.Smartypants | html.SmartypantsFractions})
}

func TestDisableSmartDashes(t *testing.T) {
	doTestsInlineParam(t, []string{
		"foo - bar\n",
		"<p>foo - bar</p>\n",
		"foo -- bar\n",
		"<p>foo -- bar</p>\n",
		"foo --- bar\n",
		"<p>foo --- bar</p>\n",
	}, TestParams{})
	doTestsInlineParam(t, []string{
		"foo - bar\n",
		"<p>foo &ndash; bar</p>\n",
		"foo -- bar\n",
		"<p>foo &mdash; bar</p>\n",
		"foo --- bar\n",
		"<p>foo &mdash;&ndash; bar</p>\n",
	}, TestParams{Flags: html.Smartypants | html.SmartypantsDashes})
	doTestsInlineParam(t, []string{
		"foo - bar\n",
		"<p>foo - bar</p>\n",
		"foo -- bar\n",
		"<p>foo &ndash; bar</p>\n",
		"foo --- bar\n",
		"<p>foo &mdash; bar</p>\n",
	}, TestParams{Flags: html.Smartypants | html.SmartypantsLatexDashes | html.SmartypantsDashes})
	doTestsInlineParam(t, []string{
		"foo - bar\n",
		"<p>foo - bar</p>\n",
		"foo -- bar\n",
		"<p>foo -- bar</p>\n",
		"foo --- bar\n",
		"<p>foo --- bar</p>\n",
	}, TestParams{Flags: html.Smartypants | html.SmartypantsLatexDashes})
}

func TestSkipLinks(t *testing.T) {
	doTestsInlineParam(t, []string{
		"[foo](gopher://foo.bar)",
		"<p><tt>foo</tt></p>\n",

		"[foo](mailto:bar/)\n",
		"<p><tt>foo</tt></p>\n",
	}, TestParams{
		Flags: html.SkipLinks,
	})
}

func TestSkipImages(t *testing.T) {
	doTestsInlineParam(t, []string{
		"![foo](/bar/)\n",
		"<p></p>\n",
	}, TestParams{
		Flags: html.SkipImages,
	})
}

func TestLazyLoadImages(t *testing.T) {
	doTestsInlineParam(t, []string{
		"![foo](/bar/)\n",
		"<p><img loading=\"lazy\" src=\"/bar/\" alt=\"foo\" /></p>\n",
	}, TestParams{
		Flags: html.LazyLoadImages,
	})
}

func TestUseXHTML(t *testing.T) {
	doTestsParam(t, []string{
		"---",
		"<hr>\n",
	}, TestParams{})
	doTestsParam(t, []string{
		"---",
		"<hr />\n",
	}, TestParams{Flags: html.UseXHTML})
}

func TestSkipHTML(t *testing.T) {
	doTestsParam(t, []string{
		"<div class=\"foo\"></div>\n\ntext\n\n<form>the form</form>",
		"<p>text</p>\n\n<p>the form</p>\n",

		"text <em>inline html</em> more text",
		"<p>text inline html more text</p>\n",
	}, TestParams{Flags: html.SkipHTML})
}

func TestInlineMath(t *testing.T) {
	doTestsParam(t, []string{
		"$a_b$",
		`<p><span class="math inline">\(a_b\)</span></p>
`,
	}, TestParams{Flags: html.SkipHTML, extensions: parser.CommonExtensions})
}

// TODO: not fixed yet. Need to change the logic and update the tests.
// https://github.com/gomarkdown/markdown/issues/327
func TestBug327(t *testing.T) {
	doTestsParam(t, []string{
		`[site](https://somesite.com/?"s"(b)h)`,
		`<p><a href="https://somesite.com/?&quot;s&quot;(b">site</a>h)</p>
`,
	}, TestParams{extensions: parser.CommonExtensions})
}

func TestSubSuper(t *testing.T) {
	var tests = []string{
		"H~2~O is a liquid, 2^10^ is 1024\n",
		"<p>H<sub>2</sub>O is a liquid, 2<sup>10</sup> is 1024</p>\n",
		"2^10^ is 1024, H~2~O is a liquid\n",
		"<p>2<sup>10</sup> is 1024, H<sub>2</sub>O is a liquid</p>\n",
		"2\\^10 is 2^10^ is 1024\n",
		"<p>2^10 is 2<sup>10</sup> is 1024</p>\n",
	}
	doTestsInlineParam(t, tests, TestParams{extensions: parser.SuperSubscript})
}

func BenchmarkSmartDoubleQuotes(b *testing.B) {
	params := TestParams{Flags: html.Smartypants}
	params.extensions |= parser.Autolink | parser.Strikethrough
	params.Flags |= html.UseXHTML

	for i := 0; i < b.N; i++ {
		runMarkdown("this should be normal \"quoted\" text.\n", params)
	}
}
