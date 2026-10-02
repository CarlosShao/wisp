package config

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// preset is a built-in provider default (SPEC-03 sec 3: presets supply
// protocol/base_url only; keys always come via api_key_ref and model entries
// are user-configured or discovered via /v1/models).
type preset struct {
	Protocol string
	BaseURL  string
}

// providerPresets is the built-in provider set (D8 + 2026-09-19 approved
// extension: minimax, mimo, stepfun). Endpoints are best-effort defaults -
// every field stays user-overridable in [llm.providers.<name>].
var providerPresets = map[string]preset{
	"openai":      {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.openai.com/v1"},
	"anthropic":   {Protocol: ProtocolAnthropic, BaseURL: "https://api.anthropic.com"},
	"deepseek":    {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.deepseek.com/v1"},
	"qwen":        {Protocol: ProtocolOpenAIChat, BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1"},
	"zhipu":       {Protocol: ProtocolOpenAIChat, BaseURL: "https://open.bigmodel.cn/api/paas/v4"},
	"moonshot":    {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.moonshot.cn/v1"},
	"siliconflow": {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.siliconflow.cn/v1"},
	"openrouter":  {Protocol: ProtocolOpenAIChat, BaseURL: "https://openrouter.ai/api/v1"},
	"ollama":      {Protocol: ProtocolOpenAIChat, BaseURL: "http://127.0.0.1:11434/v1"},
	"minimax":     {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.minimax.chat/v1"},
	"mimo":        {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.mimo.xiaomi.com/v1"},
	"stepfun":     {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.stepfun.com/v1"},
}

// PresetNames returns the built-in provider preset names (GUI provider
// picker, ticket 39). Order is alphabetical.
func PresetNames() []string {
	names := make([]string, 0, len(providerPresets))
	for n := range providerPresets {
		names = append(names, n)
	}
	slices.Sort(names)
	return names
}

// LookupPreset returns the built-in preset for name, if any.
func LookupPreset(name string) (preset, bool) {
	p, ok := providerPresets[name]
	return p, ok
}

// NewDefaults returns a Config carrying exactly the `default:"..."` tag
// values (D36 rule 3: the tags are the single source of defaults; the TOML
// file only overrides). Map-typed fields stay nil: go-toml merges into
// pre-existing maps, so pre-populating them would leak phantom entries.
func NewDefaults() *Config {
	c := &Config{}
	applyDefaults(reflect.ValueOf(c).Elem())
	return c
}

// applyDefaults walks a struct value (addressable) setting every tagged leaf.
func applyDefaults(v reflect.Value) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}
		def, hasDef := f.Tag.Lookup("default")
		switch fv.Kind() {
		case reflect.Struct:
			applyDefaults(fv)
		case reflect.Map:
			// leave nil (see NewDefaults)
		default:
			if hasDef {
				if err := setDefault(fv, def); err != nil {
					panic(fmt.Sprintf("config: bad default tag on %s.%s: %v",
						t.String(), f.Name, err))
				}
			}
		}
	}
}

// setDefault parses a tag into a leaf value. Slice defaults are
// comma-separated (elements trimmed; an empty tag yields nil).
func setDefault(fv reflect.Value, def string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(def)
	case reflect.Bool:
		b, err := strconv.ParseBool(def)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	case reflect.Int, reflect.Int64:
		n, err := strconv.ParseInt(def, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(n)
	case reflect.Float64:
		x, err := strconv.ParseFloat(def, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(x)
	case reflect.Slice:
		if def == "" {
			return nil
		}
		parts := strings.Split(def, ",")
		out := reflect.MakeSlice(fv.Type(), 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			ev := reflect.New(fv.Type().Elem()).Elem()
			if err := setDefault(ev, p); err != nil {
				return err
			}
			out = reflect.Append(out, ev)
		}
		fv.Set(out)
	default:
		return fmt.Errorf("unsupported default kind %s", fv.Kind())
	}
	return nil
}

// normalizeZero rewrites empty maps and slices to nil so that a config
// decoded from `key = []` / `[map]` with no entries compares equal (via
// reflect.DeepEqual) to one where the key was absent. This is what makes the
// load -> marshal -> load round-trip stable (GUI double-source guard).
func normalizeZero(v reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			if v.Field(i).CanSet() {
				normalizeZero(v.Field(i))
			}
		}
	case reflect.Map:
		if v.Len() == 0 && !v.IsNil() {
			v.SetZero() // empty map -> nil, matching the absent-key decode state
			return
		}
		for _, k := range v.MapKeys() {
			ev := v.MapIndex(k)
			if ev.Kind() == reflect.Struct {
				cp := reflect.New(ev.Type()).Elem()
				cp.Set(ev)
				normalizeZero(cp)
				v.SetMapIndex(k, cp)
			}
		}
	case reflect.Slice:
		if v.Len() == 0 && !v.IsNil() {
			v.SetZero()
		}
	}
}

// normalizeConfig applies normalizeZero to a whole Config.
func normalizeConfig(c *Config) {
	if c == nil {
		return
	}
	normalizeZero(reflect.ValueOf(c).Elem())
}

// deepCopyConfig returns a fully independent copy of c (maps and slices
// cloned). Config is treated as an immutable snapshot between copies; this
// is the only sanctioned way to hand the state out.
func deepCopyConfig(c *Config) *Config {
	if c == nil {
		return nil
	}
	out := &Config{}
	copyValue(reflect.ValueOf(out).Elem(), reflect.ValueOf(c).Elem())
	return out
}

// unwrapInterface peels interface wrapping off a reflect.Value (MapRange /
// Field values of interface kinds) so callers always see the dynamic value.
func unwrapInterface(v reflect.Value) reflect.Value {
	for v.IsValid() && v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

// copyValue copies src into dst for matching types (struct/map/slice/leaf).
func copyValue(dst, src reflect.Value) {
	src = unwrapInterface(src)
	if !src.IsValid() {
		return
	}
	switch src.Kind() {
	case reflect.Struct:
		t := src.Type()
		for i := 0; i < t.NumField(); i++ {
			if dst.Field(i).CanSet() {
				copyValue(dst.Field(i), src.Field(i))
			}
		}
	case reflect.Map:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
		iter := src.MapRange()
		for iter.Next() {
			ev := reflect.New(src.Type().Elem()).Elem()
			copyValue(ev, unwrapInterface(iter.Value()))
			dst.SetMapIndex(iter.Key(), ev)
		}
	case reflect.Slice:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Len()))
		for i := 0; i < src.Len(); i++ {
			copyValue(dst.Index(i), src.Index(i))
		}
	case reflect.Ptr:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.New(src.Type().Elem()))
		copyValue(dst.Elem(), src.Elem())
	default:
		dst.Set(src)
	}
}
