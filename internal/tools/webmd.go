package tools

import (
	"bytes"
	"encoding/json"
	"html"
	"strings"
	"unicode/utf8"
)

// webmd.go is the D34 "web.fetch → Markdown" leg and the F2 no-HTML rule
// (D33/F2: external content is never rendered as HTML anywhere). The rule is
// held BY CONSTRUCTION here: the converter reads HTML but can only ever emit
// plain markdown text — there is no code path that returns input markup, no
// raw-tag passthrough (svg/math/script bodies are dropped outright), and the
// result is what the bridge books as tool output text. A hostile page cannot
// hand a downstream renderer a live tag through this door, which is exactly
// the injection shape SPEC-06 §9 layer 2 defends.
//
// The scope is honest extraction, not a CommonMark round-trip: headings,
// paragraphs, lists, links, code blocks, line breaks; drop script/style/
// template/noscript/iframe/object/embed/svg/math entirely (their TEXT goes
// too — script bodies are payload, not content).

// htmlToMarkdown converts one HTML document (or fragment) to markdown text.
func htmlToMarkdown(in []byte) string {
	p := &mdParser{src: in}
	p.parse()
	out := normalizeMarkdownLines(p.b.String())
	return out
}

type mdParser struct {
	src []byte
	pos int
	b   bytes.Buffer
	// block state
	inPre   bool
	preLang string
	liDepth int
	pendNL  int // pending blank-line count, flushed lazily so output stays tidy
	links   []linkMark
}

// voidBlockTags force a paragraph break.
var voidBlockTags = map[string]bool{
	"p": true, "div": true, "section": true, "article": true, "header": true,
	"footer": true, "main": true, "aside": true, "nav": true, "blockquote": true,
	"ul": true, "ol": true, "dl": true, "table": true, "tr": true, "br": true,
	"hr": true, "figure": true, "figcaption": true, "form": true, "fieldset": true,
	"address": true, "details": true, "summary": true,
}

var dropEntirely = map[string]bool{
	"script": true, "style": true, "noscript": true, "template": true,
	"iframe": true, "object": true, "embed": true, "svg": true, "math": true,
	"head": true, "title": true,
}

// parse walks the byte stream. Any "<!" or "<?" run is skipped to ">", any
// tag handled by openTag; text between is emitted with whitespace folding.
func (p *mdParser) parse() {
	for {
		i := bytes.IndexByte(p.src[p.pos:], '<')
		if i < 0 {
			p.text(p.src[p.pos:])
			return
		}
		p.text(p.src[p.pos : p.pos+i])
		p.pos += i
		rest := p.src[p.pos:]
		switch {
		case bytes.HasPrefix(rest, []byte("<!--")):
			end := bytes.Index(rest[4:], []byte("-->"))
			if end < 0 {
				return
			}
			p.pos += 4 + end + 3
		case bytes.HasPrefix(rest, []byte("<!")), bytes.HasPrefix(rest, []byte("<?")):
			end := bytes.IndexByte(rest, '>')
			if end < 0 {
				return
			}
			p.pos += end + 1
		default:
			name, attrs, selfClose, n := readTag(rest)
			if n <= 0 { // malformed "<": emit as literal and continue
				p.textByte('<')
				p.pos++
				continue
			}
			p.pos += n
			if strings.HasPrefix(name, "/") {
				p.closeTag(strings.ToLower(name[1:]))
			} else {
				p.openTag(strings.ToLower(name), attrs, selfClose)
			}
		}
	}
}

// readTag parses "<name attrs>" returning the raw name (possibly with a
// leading slash), the concatenated attribute string, whether it is
// self-closing, and the consumed byte count (0 = not a tag).
func readTag(b []byte) (name, attrs string, selfClose bool, n int) {
	j := bytes.IndexByte(b, '>')
	if j < 0 {
		return "", "", false, 0
	}
	inner := string(b[1:j])
	if inner == "" {
		return "", "", false, 0
	}
	// a real tag starts with a letter or "/"
	c := inner[0]
	if !(c == '/' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
		return "", "", false, 0
	}
	sp := strings.IndexAny(inner, " \t\r\n\f")
	if sp < 0 {
		name, attrs = inner, ""
	} else {
		name, attrs = inner[:sp], inner[sp:]
	}
	if strings.HasSuffix(name, "/") {
		selfClose = true
		name = strings.TrimSuffix(name, "/")
	} else if strings.HasSuffix(strings.TrimSpace(attrs), "/") {
		selfClose = true
	}
	return name, attrs, selfClose, j + 1
}

func (p *mdParser) openTag(name, attrs string, selfClose bool) {
	switch name {
	case "script", "style", "noscript", "template", "iframe", "object", "embed", "svg", "math", "head", "title":
		if !selfClose {
			p.skipUntilClose(name)
		}
		return
	case "pre":
		p.inPre = true
		p.nl()
		p.b.WriteString("```")
		p.nlRaw()
		return
	case "code":
		if !p.inPre {
			p.b.WriteString("`")
		}
		return
	case "blockquote":
		p.nl()
		p.b.WriteString("> ")
		return
	case "h1", "h2", "h3", "h4", "h5", "h6":
		p.nl()
		p.nl()
		p.b.WriteString(strings.Repeat("#", int(name[1]-'0')))
		p.b.WriteByte(' ')
		return
	case "li":
		p.nl()
		ind := strings.Repeat("  ", p.liDepth)
		p.b.WriteString(ind + "- ")
		if p.liDepth < 4 {
			p.liDepth++
		}
		return
	case "ul", "ol":
		p.nl()
		return
	case "td", "th":
		if p.b.Len() > 0 {
			p.b.WriteByte(' ')
		}
		return
	case "a":
		href := extractAttr(attrs, "href")
		if href != "" {
			p.links = append(p.links, linkMark{href: href, mark: p.b.Len()})
		}
		return
	case "img":
		alt := extractAttr(attrs, "alt")
		src := extractAttr(attrs, "src")
		if src != "" {
			p.nl()
			if alt == "" {
				alt = "图片"
			}
			p.b.WriteString("![" + alt + "](" + src + ")")
			p.nl()
		}
		return
	case "br":
		if p.inPre {
			p.nlRaw()
		} else {
			p.nl()
		}
		return
	case "hr":
		p.nl()
		p.b.WriteString("----")
		p.nl()
		return
	default:
		if voidBlockTags[name] && !selfClose {
			p.nl()
		}
	}
}

func (p *mdParser) closeTag(name string) {
	switch name {
	case "pre":
		if p.inPre {
			p.inPre = false
			p.nl()
			p.b.WriteString("```")
			p.nl()
		}
		return
	case "code":
		if !p.inPre {
			p.b.WriteString("`")
		}
		return
	case "li", "ul", "ol":
		if name == "li" && p.liDepth > 0 {
			p.liDepth--
		}
		p.nl()
		return
	case "a":
		if len(p.links) > 0 {
			lm := p.links[len(p.links)-1]
			p.links = p.links[:len(p.links)-1]
			// rewrite from the mark: "[text](href)"
			if lm.mark <= p.b.Len() {
				text := p.b.String()[lm.mark:]
				t := strings.TrimSpace(text)
				p.b.Truncate(lm.mark)
				if t == "" {
					t = lm.href
				}
				p.b.WriteString("[" + mdEscape(t) + "](" + lm.href + ")")
			}
		}
		return
	default:
		if voidBlockTags[name] {
			p.nl()
		}
	}
}

// skipUntilClose consumes the whole subtree of a dropped tag by matching the
// close tag of the same name (no nesting of the same tag is assumed inside
// script/style bodies, which is true: their content is raw text mode).
func (p *mdParser) skipUntilClose(name string) {
	close := []byte("</" + name)
	i := bytes.Index(p.src[p.pos:], close)
	if i < 0 {
		p.pos = len(p.src)
		return
	}
	j := bytes.IndexByte(p.src[p.pos+i:], '>')
	if j < 0 {
		p.pos = len(p.src)
		return
	}
	p.pos += i + j + 1
}

type linkMark struct {
	href string
	mark int
}

func (p *mdParser) text(seg []byte) {
	if len(seg) == 0 {
		return
	}
	s := html.UnescapeString(string(seg))
	if p.inPre {
		p.b.WriteString(s)
		return
	}
	// F2 belt: decoded text that SPELLS a tag ("<script" after entity
	// decoding) must not leave this door looking like markup either. The
	// output is markdown text and is never rendered as HTML, but escaping
	// the character here means even a careless downstream renderer cannot
	// revive a tag from page DATA.
	s = strings.ReplaceAll(s, "<", "&lt;")
	if strings.ContainsFunc(s, func(r rune) bool { return !isSpaceRune(r) }) {
		p.b.WriteString(foldInline(s))
	} else if strings.ContainsAny(s, " \t\n") && p.b.Len() > 0 && !p.pendingSpace() {
		p.b.WriteByte(' ')
	}
}

func (p *mdParser) textByte(c byte) {
	p.b.WriteByte(c)
}

// extractAttr pulls one attribute value out of a tag's raw attribute string
// (the text between the name and ">"). Handles "k=v", k="v", k='v' spellings;
// returns "" when absent. Values are entity-unescaped like text runs.
func extractAttr(attrs, key string) string {
	lk := strings.ToLower(key)
	// scan tokens: skip whitespace, read name up to '=' or space, read value.
	for i := 0; i < len(attrs); {
		for i < len(attrs) && isSpaceByte(attrs[i]) {
			i++
		}
		start := i
		for i < len(attrs) && attrs[i] != '=' && !isSpaceByte(attrs[i]) {
			i++
		}
		name := strings.ToLower(strings.TrimRight(attrs[start:i], "/"))
		if i < len(attrs) && attrs[i] == '=' {
			i++ // past '='
			vstart := i
			var val string
			if i < len(attrs) && (attrs[i] == '"' || attrs[i] == '\'') {
				q := attrs[i]
				i++
				vstart = i
				for i < len(attrs) && attrs[i] != q {
					i++
				}
				val = attrs[vstart:i]
				if i < len(attrs) {
					i++ // past closing quote
				}
			} else {
				for i < len(attrs) && !isSpaceByte(attrs[i]) && attrs[i] != '>' && attrs[i] != '/' {
					i++
				}
				val = attrs[vstart:i]
			}
			if name == lk {
				return html.UnescapeString(val)
			}
			continue
		}
		// valueless attribute (boolean), keep scanning
		if name == lk {
			return ""
		}
	}
	return ""
}

func (p *mdParser) pendingSpace() bool {
	b := p.b.Bytes()
	return len(b) > 0 && isSpaceByte(b[len(b)-1])
}

// nl ends a line (or forces a blank line in block context).
func (p *mdParser) nl() {
	if p.inPre {
		p.nlRaw()
		return
	}
	b := p.b.Bytes()
	if len(b) == 0 {
		return
	}
	if b[len(b)-1] == '\n' {
		if len(b) >= 2 && b[len(b)-2] == '\n' {
			return
		}
		p.b.WriteByte('\n')
		return
	}
	p.b.WriteByte('\n')
}

func (p *mdParser) nlRaw() {
	b := p.b.Bytes()
	if len(b) > 0 && b[len(b)-1] != '\n' {
		p.b.WriteByte('\n')
	}
}

func isSpaceRune(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f' || r == '\v'
}

func isSpaceByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// foldInline collapses runs of whitespace (incl. newlines) to single spaces,
// keeping a leading/trailing single space so adjacent tags keep separation.
func foldInline(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	prevSpace := false
	for _, r := range s {
		if isSpaceRune(r) {
			prevSpace = true
			continue
		}
		if prevSpace && b.Len() > 0 {
			b.WriteByte(' ')
		}
		prevSpace = false
		b.WriteRune(r)
	}
	if prevSpace {
		b.WriteByte(' ')
	}
	return b.String()
}

// normalizeMarkdownLines trims trailing spaces per line and collapses 3+
// blank lines to one empty line.
func normalizeMarkdownLines(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	var out []string
	blanks := 0
	for _, l := range lines {
		if l == "" {
			blanks++
			if blanks > 1 {
				continue
			}
		} else {
			blanks = 0
		}
		out = append(out, l)
	}
	return strings.Trim(strings.Join(out, "\n"), "\n")
}

// mdEscape neutralizes the few sequences that could turn extracted text back
// into markdown syntax from inside a link label.
func mdEscape(s string) string {
	r := strings.NewReplacer("[", "\\[", "]", "\\]")
	return r.Replace(s)
}

// prettyJSON reformats a JSON body inside a fenced block. Invalid UTF-8 or
// parse failures fall back to the raw text path (the caller decides by
// content-type, this is formatting only, never validation).
func prettyJSON(b []byte) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, b, "", "  "); err != nil {
		return ""
	}
	return "```json\n" + buf.String() + "\n```"
}

// ensureUTF8 replaces invalid bytes so the tool output never carries
// broken runes into the model context (the transcript, not the body log).
func ensureUTF8(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	return strings.ToValidUTF8(string(b), "\uFFFD")
}
