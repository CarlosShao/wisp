//go:build windows

package main

// Ticket 314 ruler audit, PROBE SET B. Same export-tree-only rule as probe A.
// This one feeds synthetic sources to the ruler's OWN matcher
// (panelProfileHitsInText -> panelProfileHitsInSource), so every shape gets an
// isolated reading in one run without nine recompiles. Samples are parsed, not
// compiled: names like panelControl here are spellings, not symbols.

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestAuditRulerProbeTable(t *testing.T) {
	cases := []struct {
		id     string
		note   string
		source string
		want   []string
	}{
		{
			id:   "A1a helper return, profile leaf, value to sink",
			note: "the audit's attack 1 verbatim",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mk() string { return filepath.Join(os.TempDir(), "wisp-audit1a-panel-profile") }

func use() { _ = NewPanelManager(nil, nil, mk()) }
`,
			want: []string{"wisp-audit1a-panel-profile"},
		},
		{
			id:   "A1b helper return, plain leaf, value to sink",
			note: "same flow, leaf avoids the word profile: only SINK could see it",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mk() string { return filepath.Join(os.TempDir(), "wisp-audit1b-quicklook") }

func use() { _ = NewPanelManager(nil, nil, mk()) }
`,
		},
		{
			id:   "A1c param helper (header escape 1)",
			note: "root inside the helper, literal at the call site",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mk(leaf string) string { return filepath.Join(os.TempDir(), leaf) }

func use() { _ = NewPanelManager(nil, nil, mk("wisp-audit1c-panel-profile")) }
`,
		},
		{
			id:   "A2a two half literals inside one Join",
			note: "at the sink, the most favourable position",
			source: `package main

import (
	"os"
	"path/filepath"
)

func use() { _ = NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-33r1"+"-panel-profile")) }
`,
		},
		{
			id:   "A2b whole path by + concatenation",
			note: "os.TempDir() + \"/wisp-33r1\" + \"-panel-profile\"",
			source: `package main

import "os"

func use() { _ = NewPanelManager(nil, nil, os.TempDir()+"/wisp-33r1"+"-panel-profile") }
`,
		},
		{
			id:   "A2c leaf by strings.Join",
			note: "Join(os.TempDir(), strings.Join([]string{...}, \"\"))",
			source: `package main

import (
	"os"
	"path/filepath"
	"strings"
)

func use() {
	_ = NewPanelManager(nil, nil, filepath.Join(os.TempDir(), strings.Join([]string{"wisp-33r1", "-panel-profile"}, "")))
}
`,
		},
		{
			id:   "A2d join then suffix, outside the Join",
			note: "the + sits between the Join call and the tail literal",
			source: `package main

import (
	"os"
	"path/filepath"
)

func use() { _ = NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-33r1")+"-panel-profile") }
`,
		},
		{
			id:   "A2e concatenated value bound to a local, then to the sink",
			note: "the binding the SINK rule does promote must be a bare Join",
			source: `package main

import (
	"os"
	"path/filepath"
)

func use() {
	p := filepath.Join(os.TempDir(), "wisp-33r1") + "-panel-profile"
	_ = NewPanelManager(nil, nil, p)
}
`,
		},
		{
			id:   "A3a root os.Getenv(\"TEMP\")",
			note: "not os.TempDir()",
			source: `package main

import (
	"os"
	"path/filepath"
)

func use() { _ = NewPanelManager(nil, nil, filepath.Join(os.Getenv("TEMP"), "wisp-audit3a-panel-profile")) }
`,
		},
		{
			id:   "A3b root home directory",
			note: "not os.TempDir()",
			source: `package main

import (
	"os"
	"path/filepath"
)

func use() {
	home, _ := os.UserHomeDir()
	_ = NewPanelManager(nil, nil, filepath.Join(home, "wisp-audit3b-panel-profile"))
}
`,
		},
		{
			id:   "A3c absolute literal root",
			note: "no root expression at all",
			source: "package main\n\nfunc use() { _ = NewPanelManager(nil, nil, `C:\\Temp\\wisp-audit3c-panel-profile`) }\n",
		},
		{
			id:   "X1 package-level DataPath key, plain leaf",
			note: "bonus: not in the header's escape list",
			source: `package main

import (
	"os"
	"path/filepath"
)

type ctl struct{ DataPath string }

var sink = ctl{DataPath: filepath.Join(os.TempDir(), "wisp-auditpkg-sink")}
`,
		},
		{
			id:   "C1 control: inline join, plain leaf, to sink",
			note: "proves SINK works at all",
			source: `package main

import (
	"os"
	"path/filepath"
)

func use() { _ = NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-plain")) }
`,
			want: []string{"wisp-plain"},
		},
		{
			id:   "C2 control: NAME alone, no sink",
			note: "profile leaf in a returned join",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mk() string { return filepath.Join(os.TempDir(), "wisp-name-only-profile") }
`,
			want: []string{"wisp-name-only-profile"},
		},
		{
			id:   "C3 control: DataPath field assignment",
			note: "one of the four SINK forms, has no case in the 14",
			source: `package main

import (
	"os"
	"path/filepath"
)

func use(c *ctl) { c.DataPath = filepath.Join(os.TempDir(), "wisp-field") }
`,
			want: []string{"wisp-field"},
		},
		{
			id:   "C4 control: DataPath key inside a function body",
			note: "another of the four SINK forms, also has no case in the 14",
			source: `package main

import (
	"os"
	"path/filepath"
)

func use() {
	_ = ctl{DataPath: filepath.Join(os.TempDir(), "wisp-key")}
}
`,
			want: []string{"wisp-key"},
		},
	}

	for _, tc := range cases {
		hits := panelProfileHitsInText(t, strings.ReplaceAll(tc.id, " ", "_")+".go", []byte(tc.source))
		got := make([]string, len(hits))
		lines := make([]string, len(hits))
		for i, hit := range hits {
			got[i] = hit.leaf
			lines[i] = fmt.Sprintf("%d", hit.line)
		}
		sort.Strings(got)
		want := append([]string(nil), tc.want...)
		sort.Strings(want)
		verdict := "SILENT"
		if len(got) > 0 {
			verdict = "HIT"
		}
		ok := strings.Join(got, "|") == strings.Join(want, "|")
		if !ok {
			verdict += "/UNEXPECTED"
		}
		t.Logf("[%s] %s | leaves=%q lines=%q want=%q | %s", verdict, tc.id, got, lines, want, tc.note)
		if !ok {
			t.Errorf("%s: got %q want %q", tc.id, got, want)
		}
	}
}
