package plugin

// D46 command plugins — the part of the runtime that is a pure function of the
// tool declaration and one call's arguments: parameter screening, argv building
// and the child environment.
//
// It lives in internal/plugin (not in the manifestslot runtime that spawns the
// process) because these three steps must be identical in every leg that ever
// executes a command tool, and because they are the part a reviewer has to be
// able to read without following a process-creation call. The process itself,
// the timeout, the output caps and the extraction hand-off are in
// manifestslot/command.go; C14 validation of the same fields is in
// validator.go and the install-time pin is in command.go here.
//
// The rule this file exists to enforce, in the plan's words (PLAN.md:2665-2667):
// argv is a VECTOR handed to process creation, never a command string; a
// parameter that fails its schema never reaches argv; and a parameter whose
// value carries a shell metacharacter is refused even though the vector would
// be safe, because the wrapped CLI may interpret its own arguments.

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Parameter refusals. Every one of them happens BEFORE a process exists.
var (
	// ErrParamSchema: the argument object violates the tool's declared schema.
	ErrParamSchema = errors.New("plugin: command parameters violate the declared schema")
	// ErrParamMetachar: a value carries a shell metacharacter (D46 ban list).
	ErrParamMetachar = errors.New("plugin: command parameter contains a shell metacharacter")
	// ErrParamUnsupported: the manifest asked for a schema keyword this build
	// does not implement. Refused rather than ignored: an unenforced keyword is
	// a promise the host is not keeping.
	ErrParamUnsupported = errors.New("plugin: command parameter schema uses an unsupported keyword")
)

// CommandMetacharacters is the D46 ban list (PLAN.md:2667, SPEC-07 §4):
//
//	|  &  ;  $  `  <  >  CR  LF
//
// Every value that reaches argv is screened for these. The list is deliberately
// short: a value like "Bob's notes" or "a|b" quoted sanely by the CLI itself is
// not the threat model - the threat model is a CLI that re-parses its arguments,
// and these are the characters that let it.
const CommandMetacharacters = "|&;$`<>\r\n"

// ContainsCommandMetachar reports whether v carries a banned character.
func ContainsCommandMetachar(v string) bool {
	return strings.ContainsAny(v, CommandMetacharacters)
}

// ParamValue is one screened parameter value.
type ParamValue struct {
	// Kind is "string", "number", "boolean", "array" or "object".
	Kind string
	// Scalar is the rendered text for string/number/boolean (what argv gets).
	Scalar string
	// Array holds the rendered elements when Kind == "array".
	Array []string
	// Object holds the rendered members when Kind == "object".
	Object map[string]string
}

// CommandCall is one screened call: the parameter values in a shape argv
// building can use, plus the metadata the risk gate needs.
type CommandCall struct {
	// Values are the placeholders the template substitutes, by name.
	Values map[string]ParamValue
	// Names are the placeholder names in template order.
	Names []string
}

// ScreenCommandParams is the pre-spawn gate for one call.
//
// Order matters and is contractual: parse -> schema -> metacharacter screen.
// A value that violates the schema is refused before anyone looks at its bytes,
// so the two refusals never compete (a 5000-character string with a ';' in it
// is a schema violation, and saying so is more useful than saying "metachar").
func ScreenCommandParams(t ToolDef, params json.RawMessage) (CommandCall, error) {
	raw := strings.TrimSpace(string(params))
	args := map[string]any{}
	if raw != "" {
		if err := json.Unmarshal([]byte(raw), &args); err != nil {
			return CommandCall{}, fmt.Errorf("%w: arguments are not a JSON object: %v", ErrParamSchema, err)
		}
		if args == nil {
			return CommandCall{}, fmt.Errorf("%w: arguments are JSON null", ErrParamSchema)
		}
	}
	props, err := schemaProperties(t.Parameters)
	if err != nil {
		return CommandCall{}, fmt.Errorf("%w: %v", ErrParamSchema, err)
	}
	for _, key := range sortedKeys(args) {
		if _, ok := props[key]; !ok {
			return CommandCall{}, fmt.Errorf("%w: parameter %q is not declared (the command schema is closed)", ErrParamSchema, key)
		}
	}
	for _, name := range requiredParameters(t.Parameters) {
		if _, ok := args[name]; !ok {
			return CommandCall{}, fmt.Errorf("%w: parameter %q is required", ErrParamSchema, name)
		}
	}

	names := substitutedPlaceholders(t.ArgvTemplate)
	call := CommandCall{Values: make(map[string]ParamValue, len(names)), Names: names}
	for _, name := range names {
		value, ok := args[name]
		if !ok {
			// Not required by the schema and not supplied: nothing is
			// substituted for it, and BuildCommandArgv refuses a template that
			// still needs it.
			continue
		}
		schema := props[name]
		pv, err := screenValue(name, value, schema)
		if err != nil {
			return CommandCall{}, err
		}
		call.Values[name] = pv
	}
	return call, nil
}

func screenValue(name string, value any, schema map[string]any) (ParamValue, error) {
	if err := checkKeywords(name, schema); err != nil {
		return ParamValue{}, err
	}
	want, _ := schema["type"].(string)
	switch want {
	case "string":
		s, ok := value.(string)
		if !ok {
			return ParamValue{}, fmt.Errorf("%w: parameter %q must be a string", ErrParamSchema, name)
		}
		if err := checkStringConstraints(name, s, schema); err != nil {
			return ParamValue{}, err
		}
		return ParamValue{Kind: "string", Scalar: s}, nil
	case "integer", "number":
		f, err := asNumber(name, value)
		if err != nil {
			return ParamValue{}, err
		}
		if want == "integer" && f != math.Trunc(f) {
			return ParamValue{}, fmt.Errorf("%w: parameter %q must be an integer", ErrParamSchema, name)
		}
		if err := checkNumberConstraints(name, f, schema); err != nil {
			return ParamValue{}, err
		}
		return ParamValue{Kind: "number", Scalar: strconv.FormatFloat(f, 'f', -1, 64)}, nil
	case "boolean":
		b, ok := value.(bool)
		if !ok {
			return ParamValue{}, fmt.Errorf("%w: parameter %q must be a boolean", ErrParamSchema, name)
		}
		return ParamValue{Kind: "boolean", Scalar: strconv.FormatBool(b)}, nil
	case "array":
		arr, ok := value.([]any)
		if !ok {
			return ParamValue{}, fmt.Errorf("%w: parameter %q must be an array", ErrParamSchema, name)
		}
		if maxLen, ok := schemaInt(schema, "maxItems"); ok && len(arr) > maxLen {
			return ParamValue{}, fmt.Errorf("%w: parameter %q has %d items (maxItems %d)", ErrParamSchema, name, len(arr), maxLen)
		}
		itemSchema, _ := schema["items"].(map[string]any)
		out := make([]string, 0, len(arr))
		for i, item := range arr {
			pv, err := screenValue(fmt.Sprintf("%s[%d]", name, i), item, itemSchema)
			if err != nil {
				return ParamValue{}, err
			}
			if pv.Kind == "array" || pv.Kind == "object" {
				return ParamValue{}, fmt.Errorf("%w: parameter %q[%d]: nested containers are refused (flatten it or declare it as a string)", ErrParamSchema, name, i)
			}
			if err := screenMetachars(fmt.Sprintf("%s[%d]", name, i), pv.Scalar); err != nil {
				return ParamValue{}, err
			}
			out = append(out, pv.Scalar)
		}
		return ParamValue{Kind: "array", Array: out}, nil
	case "object":
		obj, ok := value.(map[string]any)
		if !ok {
			return ParamValue{}, fmt.Errorf("%w: parameter %q must be an object", ErrParamSchema, name)
		}
		if maxLen, ok := schemaInt(schema, "maxProperties"); ok && len(obj) > maxLen {
			return ParamValue{}, fmt.Errorf("%w: parameter %q has %d members (maxProperties %d)", ErrParamSchema, name, len(obj), maxLen)
		}
		known := map[string]bool{}
		if propsRaw, ok := schema["properties"].(map[string]any); ok {
			for k := range propsRaw {
				known[k] = true
			}
		}
		out := make(map[string]string, len(obj))
		for _, k := range sortedKeys(obj) {
			if len(known) > 0 && !known[k] {
				return ParamValue{}, fmt.Errorf("%w: parameter %q member %q is not declared", ErrParamSchema, name, k)
			}
			if err := screenMetachars(fmt.Sprintf("%s.%s", name, k), k); err != nil {
				return ParamValue{}, err
			}
			pv, err := screenValue(fmt.Sprintf("%s.%s", name, k), obj[k], nil)
			if err != nil {
				return ParamValue{}, err
			}
			if pv.Kind == "array" || pv.Kind == "object" {
				return ParamValue{}, fmt.Errorf("%w: parameter %q.%s: nested containers are refused", ErrParamSchema, name, k)
			}
			out[k] = pv.Scalar
		}
		return ParamValue{Kind: "object", Object: out}, nil
	case "":
		return ParamValue{}, fmt.Errorf("%w: parameter %q declares no type (D46 requires type/enum/range/maxLength per placeholder)", ErrParamSchema, name)
	default:
		return ParamValue{}, fmt.Errorf("%w: parameter %q declares unsupported type %q", ErrParamUnsupported, name, want)
	}
}

// supportedKeywords is the JSON Schema subset D46 parameters may use. A
// declaration outside it is refused: a keyword the host ignores is an
// expectation the plugin author is entitled to hold and would silently lose.
var supportedKeywords = map[string]bool{
	"type": true, "enum": true, "description": true, "title": true,
	"minLength": true, "maxLength": true, "pattern": true, "format": true,
	"minimum": true, "maximum": true, "exclusiveMinimum": true, "exclusiveMaximum": true,
	"minimumLength": true, "maxItems": true, "minItems": true, "maxProperties": true,
	"items": true, "properties": true, "required": true, "additionalProperties": true,
	"default": true, "examples": true, "const": true,
}

func checkKeywords(name string, schema map[string]any) error {
	if len(schema) == 0 {
		return nil
	}
	for _, k := range sortedKeys(schema) {
		if !supportedKeywords[k] {
			return fmt.Errorf("%w: parameter %q uses %q (supported: %s)",
				ErrParamUnsupported, name, k, supportedKeywordList())
		}
	}
	return nil
}

func supportedKeywordList() string {
	keys := sortedKeys(supportedKeywords)
	return strings.Join(keys, ", ")
}

func checkStringConstraints(name, s string, schema map[string]any) error {
	if enum, ok := schemaStringList(schema, "enum"); ok && !containsString(enum, s) {
		return fmt.Errorf("%w: parameter %q must be one of %s; got %q", ErrParamSchema, name, strings.Join(enum, "|"), s)
	}
	if v, ok := schemaString(schema, "const"); ok && s != v {
		return fmt.Errorf("%w: parameter %q must equal %q", ErrParamSchema, name, v)
	}
	if n, ok := schemaInt(schema, "minLength"); ok && utf8.RuneCountInString(s) < n {
		return fmt.Errorf("%w: parameter %q is shorter than minLength %d", ErrParamSchema, name, n)
	}
	if n, ok := schemaInt(schema, "maxLength"); ok && utf8.RuneCountInString(s) > n {
		return fmt.Errorf("%w: parameter %q is longer than maxLength %d", ErrParamSchema, name, n)
	}
	if pat, ok := schemaString(schema, "pattern"); ok {
		re, err := regexp.Compile(pat)
		if err != nil {
			return fmt.Errorf("%w: parameter %q pattern is not valid: %v", ErrParamUnsupported, name, err)
		}
		if !re.MatchString(s) {
			return fmt.Errorf("%w: parameter %q does not match pattern %q", ErrParamSchema, name, pat)
		}
	}
	return screenMetachars(name, s)
}

func checkNumberConstraints(name string, f float64, schema map[string]any) error {
	if enum, ok := schemaAnyList(schema, "enum"); ok {
		found := false
		for _, e := range enum {
			if n, err := asNumber(name, e); err == nil && n == f {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%w: parameter %q is not one of the declared enum values", ErrParamSchema, name)
		}
	}
	if v, ok := schemaFloat(schema, "minimum"); ok && f < v {
		return fmt.Errorf("%w: parameter %q is below minimum %v", ErrParamSchema, name, v)
	}
	if v, ok := schemaFloat(schema, "maximum"); ok && f > v {
		return fmt.Errorf("%w: parameter %q is above maximum %v", ErrParamSchema, name, v)
	}
	return nil
}

// screenMetachars is D46's ban list applied to one value ("参数值若含
// | & ; $ ` < > \n \r -> 拒绝").
func screenMetachars(name, v string) error {
	if ContainsCommandMetachar(v) {
		shown, banned := firstMetachar(v)
		return fmt.Errorf("%w: parameter %q contains %q at byte %d (banned set %q)",
			ErrParamMetachar, name, string(banned), shown, CommandMetacharacters)
	}
	return nil
}

func firstMetachar(v string) (int, rune) {
	for i, r := range v {
		if strings.ContainsRune(CommandMetacharacters, r) {
			return i, r
		}
	}
	return 0, 0
}

func asNumber(name string, value any) (float64, error) {
	switch n := value.(type) {
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0, fmt.Errorf("%w: parameter %q is not a number", ErrParamSchema, name)
		}
		return f, nil
	case float64:
		return n, nil
	}
	return 0, fmt.Errorf("%w: parameter %q must be a number", ErrParamSchema, name)
}

func schemaString(schema map[string]any, key string) (string, bool) {
	if schema == nil {
		return "", false
	}
	s, ok := schema[key].(string)
	return s, ok
}

func schemaInt(schema map[string]any, key string) (int, bool) {
	if schema == nil {
		return 0, false
	}
	if f, ok := schemaFloat(schema, key); ok {
		return int(f), true
	}
	return 0, false
}

func schemaFloat(schema map[string]any, key string) (float64, bool) {
	if schema == nil {
		return 0, false
	}
	switch n := schema[key].(type) {
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case float64:
		return n, true
	case int64:
		return float64(n), true
	case int:
		return float64(n), true
	}
	return 0, false
}

func schemaStringList(schema map[string]any, key string) ([]string, bool) {
	list, ok := schemaAnyList(schema, key)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(list))
	for _, v := range list {
		s, ok := v.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

func schemaAnyList(schema map[string]any, key string) ([]any, bool) {
	if schema == nil {
		return nil, false
	}
	list, ok := schema[key].([]any)
	return list, ok
}

// BuildCommandArgv renders the argv VECTOR for one screened call.
//
// Every element of the result is ONE argument. Nothing is joined, nothing is
// quoted for a shell (there is no shell), and a placeholder that could not be
// filled is an error rather than an empty string - an empty argv slot silently
// becomes a different command.
//
// A placeholder that is the WHOLE element ("{{to}}") of an array-typed
// parameter expands to one argument per element; embedded placeholders
// ("--to={{to}}") accept only scalars.
func BuildCommandArgv(t ToolDef, call CommandCall) ([]string, error) {
	argv := make([]string, 0, len(t.ArgvTemplate))
	for i, element := range t.ArgvTemplate {
		if i == 0 {
			// argv[0] is the pinned exe, from the pin - never from a
			// parameter (C14 refuses a placeholder there; this is the second
			// half of the same rule).
			argv = append(argv, t.Exe)
			continue
		}
		whole := fullPlaceholder(element)
		if whole != "" {
			pv, ok := call.Values[whole]
			if !ok {
				return nil, fmt.Errorf("plugin: tool %q: argv[%d] needs parameter %q, which the call did not supply",
					t.Name, i, whole)
			}
			switch pv.Kind {
			case "array":
				if len(pv.Array) == 0 {
					return nil, fmt.Errorf("plugin: tool %q: argv[%d]: parameter %q is an empty array and its slot cannot be dropped silently",
						t.Name, i, whole)
				}
				argv = append(argv, pv.Array...)
			case "object":
				return nil, fmt.Errorf("plugin: tool %q: argv[%d]: parameter %q is an object and cannot be one argument",
					t.Name, i, whole)
			default:
				argv = append(argv, pv.Scalar)
			}
			continue
		}
		rendered, err := renderElement(t.Name, i, element, call)
		if err != nil {
			return nil, err
		}
		argv = append(argv, rendered)
	}
	if len(argv) == 0 || argv[0] == "" {
		return nil, fmt.Errorf("plugin: tool %q: empty argv", t.Name)
	}
	return argv, nil
}

// fullPlaceholder returns name when element is exactly "{{name}}".
func fullPlaceholder(element string) string {
	if !strings.HasPrefix(element, "{{") || !strings.HasSuffix(element, "}}") {
		return ""
	}
	inner := element[2 : len(element)-2]
	if inner == "" || strings.Contains(inner, "{{") || strings.Contains(inner, "}}") {
		return ""
	}
	if !placeholderPattern.MatchString(element) || placeholderPattern.FindString(element) != element {
		return ""
	}
	return inner
}

func renderElement(tool string, index int, element string, call CommandCall) (string, error) {
	var b strings.Builder
	last := 0
	for _, loc := range placeholderPattern.FindAllStringSubmatchIndex(element, -1) {
		name := element[loc[2]:loc[3]]
		b.WriteString(element[last:loc[0]])
		pv, ok := call.Values[name]
		if !ok {
			return "", fmt.Errorf("plugin: tool %q: argv[%d] embeds {{%s}}, which the call did not supply", tool, index, name)
		}
		if pv.Kind == "array" || pv.Kind == "object" {
			return "", fmt.Errorf("plugin: tool %q: argv[%d]: %q is a %s and cannot be embedded in one argument",
				tool, index, name, pv.Kind)
		}
		b.WriteString(pv.Scalar)
		last = loc[1]
	}
	b.WriteString(element[last:])
	return b.String(), nil
}

// ---------------------------------------------------------------------------
// Environment allowlist (D46: "子进程环境白名单制")
// ---------------------------------------------------------------------------

// baseEnvNames is the contract's always-present set (PLAN.md:2671, SPEC-07 §4).
// SystemRoot/SystemDrive/WINDIR/TEMP/TMP join it because a Windows process
// without SystemRoot cannot load its own DLLs, and a CLI without TEMP cannot
// write a temp file - refusing those is not a security win, it is a broken
// install. None of them carries a credential by design.
var baseEnvNames = []string{"PATH", "HOME", "USERPROFILE", "LANG", "SystemRoot", "SystemDrive", "WINDIR", "TEMP", "TMP"}

// pathValueCap bounds the PATH that reaches a child: a long PATH is not a
// credential, but it is a map of everything installed on the machine, and the
// wrapped CLI has no business knowing about most of it.
const pathValueCap = 8 * 1024

// envCaseInsensitive reports whether environment names collide case-folded
// (Windows), which is why the allowlist comparison does too.
const envCaseInsensitive = true

// CommandEnv is the child environment for one call: BASE + the manifest's
// allowlist, resolved from the host environment, and NOTHING else.
//
// It never calls os.Environ() wholesale. That is the whole point: the host
// process holds API keys, tokens and proxy credentials, and a wrapped CLI that
// inherited them would carry them off the machine (the ticket's "key-leak
// defense", SPEC-07 §7's echo-fixture proof). A name that is in the allowlist
// but absent from the host is simply absent from the child.
//
// lookup is injected so the test leg can run the same code over a sealed
// environment instead of whichever one the developer's shell happened to have.
func CommandEnv(t ToolDef, lookup func(string) (string, bool)) []string {
	if lookup == nil {
		lookup = os.LookupEnv
	}
	var names []string
	names = append(names, baseEnvNames...)
	names = append(names, t.EnvAllowlist...)

	seen := make(map[string]bool, len(names))
	out := make([]string, 0, len(names))
	for _, name := range names {
		key := name
		if envCaseInsensitive {
			key = strings.ToLower(name)
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		value, ok := lookup(name)
		if !ok {
			continue
		}
		if strings.EqualFold(name, "PATH") {
			value = trimPathValue(value)
		}
		out = append(out, name+"="+value)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// trimPathValue drops empty components and caps the total length. It does not
// reorder or rewrite the directories: following PATH must still find the same
// executables the user expects.
func trimPathValue(v string) string {
	parts := strings.Split(v, string(os.PathListSeparator))
	kept := make([]string, 0, len(parts))
	total := 0
	for _, p := range parts {
		if p == "" {
			continue
		}
		if total+len(p) > pathValueCap {
			break
		}
		total += len(p) + 1
		kept = append(kept, p)
	}
	return strings.Join(kept, string(os.PathListSeparator))
}

// ChildEnvNames returns the variable names a child of this tool may see. It is
// the same allowlist CommandEnv builds, exposed for the install review sheet and
// the audit row ("what could this plugin have read?").
//
// It is the ALLOWLIST, not host∩allowlist: the review sheet is shown BEFORE the
// install and must name what the plugin asked for, and an audit row that
// silently shrank with the operator's environment would answer a different
// question than the one asked. CommandEnv stays the only thing that decides what
// is actually passed (a name the host does not have is absent from the child).
func ChildEnvNames(t ToolDef) []string {
	names := make([]string, 0, len(baseEnvNames)+len(t.EnvAllowlist))
	seen := make(map[string]bool, cap(names))
	for _, name := range append(append([]string{}, baseEnvNames...), t.EnvAllowlist...) {
		key := name
		if envCaseInsensitive {
			key = strings.ToLower(name)
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// CommandTimeoutMs resolves the per-call budget: the manifest's value, clamped
// into the range the validator guarantees, or the D46 default.
func CommandTimeoutMs(t ToolDef) int {
	switch {
	case t.TimeoutMs == 0:
		return DefaultCommandTimeoutMS
	case t.TimeoutMs < MinCommandTimeoutMS:
		return MinCommandTimeoutMS
	case t.TimeoutMs > MaxCommandTimeoutMS:
		return MaxCommandTimeoutMS
	}
	return t.TimeoutMs
}
