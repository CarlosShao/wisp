package secret

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"

	"github.com/pelletier/go-toml/v2"
)

// ConfigRefs returns every api_key_ref in the config file at configPath as
// dotted TOML path -> ref value (e.g. "llm.providers.deepseek.api_key_ref" ->
// "dpapi:stepfun"). It reads the file with the same tree walk
// MigratePlaintext uses, so nested tables and arrays of tables are covered.
//
// This is the reference-index view `wisp secret unset` needs to refuse
// deleting a blob a live config still points at (ticket 63). The values are
// references, never secret material, so they are safe in diagnostics. A
// missing config file yields an empty map and no error: no config means no
// references.
func ConfigRefs(configPath string) (map[string]string, error) {
	raw, err := os.ReadFile(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("secret: config refs: read %s: %w", configPath, err)
	}
	var tree map[string]any
	if err := toml.Unmarshal(raw, &tree); err != nil {
		return nil, fmt.Errorf("secret: config refs: parse %s: %w", configPath, err)
	}
	out := map[string]string{}
	var walk func(node map[string]any, prefix []string)
	walk = func(node map[string]any, prefix []string) {
		for k, v := range node {
			switch t := v.(type) {
			case map[string]any:
				walk(t, append(append([]string{}, prefix...), k))
			case []any:
				for i, e := range t {
					if m, ok := e.(map[string]any); ok {
						walk(m, append(append([]string{}, prefix...), fmt.Sprintf("%s[%d]", k, i)))
					}
				}
			case string:
				if k == RefKey {
					out[joinPath(prefix, k)] = t
				}
			}
		}
	}
	walk(tree, nil)
	return out, nil
}

// RefFieldNames filters ConfigRefs down to the dotted paths holding ref,
// sorted for stable output. An empty result means nothing in the config
// points at that ref.
func RefFieldNames(configPath, ref string) ([]string, error) {
	refs, err := ConfigRefs(configPath)
	if err != nil {
		return nil, err
	}
	var out []string
	for path, v := range refs {
		if v == ref {
			out = append(out, path)
		}
	}
	slices.Sort(out)
	return out, nil
}
