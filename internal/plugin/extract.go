package plugin

// Declarative output extraction (D46, C14 `extract`).
//
// Why this exists rather than "return stdout": the extracted value is the ONLY
// thing that may enter the model context from a command call. Raw stdout is
// arbitrary bytes from an arbitrary program - it can be a 40MB JSON dump, a
// colour-coded progress bar, or a sentence written by whoever sent the last
// message. The extractor says, in the manifest, which part of the output is the
// answer, and the host caps that answer at CommandExtractCapBytes.
//
// The grammar is deliberately the small, strict subset a manifest author can
// get right without a parser library, and every construct is fail-closed:
//
//	extract = { kind = "json",  path = ".data.items[].summary" }
//	extract = { kind = "regex", pattern = `^sha:(?P<hex>[0-9a-f]{7,40})(?:\s|$)`, group = "hex", multiline = true }
//
// JSON path:   $ / @ root, .name / ["name"] member, [n] index, [] iterate,
//              ..name deep scan. A path that matches nothing is an EMPTY
//              result, not an error (a CLI that printed an empty list is not a
//              failure); a path that matches more than once yields the list.
// Regex:       Go regexp, optional named group (default: whole match, or the
//              group whose name is "1"), all matches collected, each trimmed.
//              A pattern that matches nothing is an empty result.
//
// The extractor NEVER returns raw output on failure: a parse error is a
// tool-level error, and its text names the extractor, not the payload.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Extract kinds.
const (
	ExtractKindJSON  = "json"
	ExtractKindRegex = "regex"
)

// Extractor turns a command's stdout into the text that may enter the context.
type Extractor interface {
	// Extract returns the extracted text (already capped by the caller, never
	// here) or an error describing WHY the output could not be read. It never
	// returns the payload alongside an error.
	Extract(stdout []byte) (string, error)
	// Kind is "json" or "regex" (diagnostics, audit rows).
	Kind() string
}

// NewExtractor builds the extractor a manifest declared; see newExtractor for
// the closed grammar. It is exported so the runtime leg (which lives in another
// package because of the import direction) reads the same declaration with the
// same rules instead of growing a second grammar.
func NewExtractor(decl map[string]any) (Extractor, error) { return newExtractor(decl) }

// newExtractor builds the extractor a manifest declared. The grammar is closed:
// an unknown key, an unknown kind or a missing required key is an error, so a
// typo can never leave a call with "no extractor" and raw output on the path.
func newExtractor(decl map[string]any) (Extractor, error) {
	if len(decl) == 0 {
		return nil, fmt.Errorf("extract is empty: declare kind = %q or %q", ExtractKindJSON, ExtractKindRegex)
	}
	kind, _ := decl["kind"].(string)
	switch kind {
	case ExtractKindJSON:
		return newJSONExtractor(decl)
	case ExtractKindRegex:
		return newRegexExtractor(decl)
	case "":
		return nil, fmt.Errorf("extract.kind is required (%q or %q)", ExtractKindJSON, ExtractKindRegex)
	default:
		return nil, fmt.Errorf("extract.kind %q is not one of %q/%q", kind, ExtractKindJSON, ExtractKindRegex)
	}
}

// checkExtractKeys refuses keys the chosen kind does not read: a field the host
// ignores is a promise it is not keeping.
func checkExtractKeys(decl map[string]any, allowed ...string) error {
	for _, k := range sortedKeys(decl) {
		if k == "kind" {
			continue
		}
		ok := false
		for _, a := range allowed {
			if k == a {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("extract.%s is not a field of this kind (allowed: kind, %s)", k, strings.Join(allowed, ", "))
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// JSON path
// ---------------------------------------------------------------------------

type jsonExtractor struct {
	path []jsonStep
	raw  string
}

func (e *jsonExtractor) Kind() string { return ExtractKindJSON }

func newJSONExtractor(decl map[string]any) (Extractor, error) {
	if err := checkExtractKeys(decl, "path"); err != nil {
		return nil, err
	}
	raw, ok := decl["path"].(string)
	if !ok || strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("extract.path is required for kind = %q (e.g. \".data.items[].summary\")", ExtractKindJSON)
	}
	steps, err := parseJSONPath(raw)
	if err != nil {
		return nil, err
	}
	return &jsonExtractor{path: steps, raw: raw}, nil
}

type jsonStepKind int

const (
	jsonStepMember jsonStepKind = iota
	jsonStepIndex
	jsonStepIterate
	jsonStepDeep
)

type jsonStep struct {
	kind  jsonStepKind
	name  string
	index int
}

// parseJSONPath implements the grammar documented above.
func parseJSONPath(p string) ([]jsonStep, error) {
	i := 0
	if strings.HasPrefix(p, "$") || strings.HasPrefix(p, "@") {
		i = 1
	}
	var steps []jsonStep
	for i < len(p) {
		switch {
		case strings.HasPrefix(p[i:], ".."):
			i += 2
			j := i
			for j < len(p) && isPathIdentByte(p[j]) {
				j++
			}
			if j == i {
				return nil, fmt.Errorf("extract.path %q: '..' must be followed by a member name", p)
			}
			steps = append(steps, jsonStep{kind: jsonStepDeep, name: p[i:j]})
			i = j
		case p[i] == '.':
			i++
			j := i
			for j < len(p) && isPathIdentByte(p[j]) {
				j++
			}
			if j == i {
				return nil, fmt.Errorf("extract.path %q: '.' must be followed by a member name", p)
			}
			steps = append(steps, jsonStep{kind: jsonStepMember, name: p[i:j]})
			i = j
		case p[i] == '[':
			end := strings.IndexByte(p[i:], ']')
			if end < 0 {
				return nil, fmt.Errorf("extract.path %q: unclosed '['", p)
			}
			inner := p[i+1 : i+end]
			i += end + 1
			if inner == "" {
				steps = append(steps, jsonStep{kind: jsonStepIterate})
				continue
			}
			if n, err := strconv.Atoi(inner); err == nil {
				if n < 0 {
					return nil, fmt.Errorf("extract.path %q: negative index %d", p, n)
				}
				steps = append(steps, jsonStep{kind: jsonStepIndex, index: n})
				continue
			}
			quoted := strings.TrimSpace(inner)
			if len(quoted) >= 2 && (quoted[0] == '"' && quoted[len(quoted)-1] == '"' || quoted[0] == '\'' && quoted[len(quoted)-1] == '\'') {
				var name string
				if quoted[0] == '"' {
					if err := json.Unmarshal([]byte(quoted), &name); err != nil {
						return nil, fmt.Errorf("extract.path %q: %v", p, err)
					}
				} else {
					name = quoted[1 : len(quoted)-1]
				}
				steps = append(steps, jsonStep{kind: jsonStepMember, name: name})
				continue
			}
			return nil, fmt.Errorf("extract.path %q: [%s] must be an index, empty (iterate), or a quoted member name", p, inner)
		default:
			return nil, fmt.Errorf("extract.path %q: unexpected %q at offset %d", p, p[i], i)
		}
	}
	return steps, nil
}

func isPathIdentByte(b byte) bool {
	switch {
	case b >= 'a' && b <= 'z', b >= 'A' && b <= 'Z', b >= '0' && b <= '9', b == '_', b == '-':
		return true
	}
	return false
}

// Extract decodes the document and walks the path. A document that is not
// JSON is an error; a path that selects nothing is an empty result.
func (e *jsonExtractor) Extract(stdout []byte) (string, error) {
	trimmed := bytes.TrimSpace(stdout)
	if len(trimmed) == 0 {
		return "", nil
	}
	var doc any
	dec := json.NewDecoder(bytes.NewReader(trimmed))
	dec.UseNumber()
	if err := dec.Decode(&doc); err != nil {
		return "", fmt.Errorf("extract: stdout is not JSON (path %q): %v", e.raw, err)
	}
	values := walkJSON(doc, e.path)
	switch len(values) {
	case 0:
		return "", nil
	case 1:
		return renderJSONValue(values[0])
	}
	// Several hits stay a JSON ARRAY, rendered per element by the JSON
	// encoder: a plain string element is quoted, so the shape a manifest
	// declared (".data.items[].summary") is still JSON the model can read.
	// Joining raw values would turn ["a","b"] into "[a,b]", which is neither
	// the payload nor a valid document.
	b, err := json.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("extract: cannot render %d selected values: %v", len(values), err)
	}
	return string(b), nil
}

// walkJSON applies the steps to one value, flattening iteration and deep scan.
func walkJSON(doc any, steps []jsonStep) []any {
	cur := []any{doc}
	for _, st := range steps {
		var next []any
		for _, v := range cur {
			next = append(next, applyJSONStep(v, st)...)
		}
		cur = next
		if len(cur) == 0 {
			return nil
		}
	}
	return cur
}

func applyJSONStep(v any, st jsonStep) []any {
	switch st.kind {
	case jsonStepMember:
		if m, ok := v.(map[string]any); ok {
			if got, ok := m[st.name]; ok {
				return []any{got}
			}
		}
		return nil
	case jsonStepIndex:
		if arr, ok := v.([]any); ok && st.index < len(arr) {
			return []any{arr[st.index]}
		}
		return nil
	case jsonStepIterate:
		if arr, ok := v.([]any); ok {
			return arr
		}
		return nil
	case jsonStepDeep:
		var out []any
		var visit func(any)
		visit = func(node any) {
			switch n := node.(type) {
			case map[string]any:
				for _, k := range sortedKeys(n) {
					if k == st.name {
						out = append(out, n[k])
					}
					visit(n[k])
				}
			case []any:
				for _, item := range n {
					visit(item)
				}
			}
		}
		visit(v)
		return out
	}
	return nil
}

func renderJSONValue(v any) (string, error) {
	switch s := v.(type) {
	case string:
		return s, nil
	case json.Number:
		return s.String(), nil
	case nil:
		return "null", nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("extract: cannot render a selected value: %v", err)
	}
	return string(b), nil
}

// ---------------------------------------------------------------------------
// Regex
// ---------------------------------------------------------------------------

type regexExtractor struct {
	re        *regexp.Regexp
	group     string
	pattern   string
	multiline bool
}

func (e *regexExtractor) Kind() string { return ExtractKindRegex }

func newRegexExtractor(decl map[string]any) (Extractor, error) {
	if err := checkExtractKeys(decl, "pattern", "group", "multiline"); err != nil {
		return nil, err
	}
	pattern, ok := decl["pattern"].(string)
	if !ok || pattern == "" {
		return nil, fmt.Errorf("extract.pattern is required for kind = %q", ExtractKindRegex)
	}
	group := ""
	switch g := decl["group"].(type) {
	case nil:
	case string:
		group = g
	default:
		return nil, fmt.Errorf("extract.group must be a string (a named group, or \"0\")")
	}
	multiline := false
	switch m := decl["multiline"].(type) {
	case nil:
	case bool:
		multiline = m
	default:
		return nil, fmt.Errorf("extract.multiline must be a boolean")
	}
	src := pattern
	if multiline && !strings.HasPrefix(src, "(?m)") && !strings.HasPrefix(src, "(?") {
		src = "(?m)" + src
	}
	re, err := regexp.Compile(src)
	if err != nil {
		return nil, fmt.Errorf("extract.pattern is not a valid Go regular expression: %v", err)
	}
	if group != "" {
		// A group is addressable two ways: by name, or by the index a caller
		// can read off the pattern. Index 0 means the WHOLE match (and is not a
		// subexpression), so the bound is `0 <= n <= NumSubexp` - a check of
		// `n < NumSubexp` rejects "1" for a one-group pattern, which is the
		// declaration SPEC-07 §4's own example shape uses.
		found := false
		for _, name := range re.SubexpNames() {
			if name == group {
				found = true
				break
			}
		}
		if n, err := strconv.Atoi(group); err == nil && n >= 0 && n <= re.NumSubexp() {
			found = true
		}
		if !found {
			return nil, fmt.Errorf("extract.group %q names no group in the pattern", group)
		}
	}
	return &regexExtractor{re: re, group: group, pattern: pattern, multiline: multiline}, nil
}

// Extract runs the pattern over the whole stdout and collects every match
// (trimmed). No match is an empty result, not an error.
func (e *regexExtractor) Extract(stdout []byte) (string, error) {
	if len(stdout) == 0 {
		return "", nil
	}
	all := e.re.FindAllSubmatch(stdout, -1)
	if len(all) == 0 {
		return "", nil
	}
	idx := 0
	if e.group != "" {
		if n, err := strconv.Atoi(e.group); err == nil {
			idx = n
		} else {
			idx = e.re.SubexpIndex(e.group)
		}
	}
	var parts []string
	for _, m := range all {
		if idx >= len(m) {
			continue
		}
		parts = append(parts, strings.TrimSpace(string(m[idx])))
	}
	if e.group == "" && len(parts) == 1 {
		return parts[0], nil
	}
	return strings.Join(parts, "\n"), nil
}

// sortedKeys returns the keys of a map in stable order (diagnostics, error
// text and the deep-scan walk all need it).
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
