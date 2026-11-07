package panel

// Ticket 190, line level: the top-level `files` key on the composer section.
//
// The reader-level pins live in the tree_190 family (root states, read-only
// surface, rewrite accounting). What nobody had pinned is the WIRE behavior of
// the key itself on ComposerState - the layer TestComposerContractTypesMatch-
// Frontend reconciles against frontend/src/lib/panel.ts. Per ledger A1487 the
// page-side declaration for this key belongs to the owner's external UI leg
// (the 181/197/248 order: the Go side lands the key first, the page declares
// it after), so until that declaration lands the contract test stays red on
// exactly this one key. These tests hold the Go half of that contract still:
// the key's name, its layer, its on/off-wire semantics, and the fact that a
// refusal rides the key instead of being dropped - so the page side lands
// against a wire that has not drifted underneath it.

import (
	"encoding/json"
	"testing"
)

// filesKeyOf190 marshals a fully assembled snapshot and returns the composer
// section exactly as it crossed the wire, so every assertion below reads the
// same JSON a renderer would receive.
func filesKeyOf190(t *testing.T, src PumpSources) map[string]any {
	t.Helper()
	wire, err := json.Marshal(NewSnapshotPump(src).Snapshot())
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	var packet map[string]any
	if err := json.Unmarshal(wire, &packet); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
	composer, ok := packet["composer"].(map[string]any)
	if !ok {
		t.Fatalf("snapshot wire carries no composer object: %s", wire)
	}
	return composer
}

// Test190ComposerFilesKeyAbsentWithoutReader pins the omitempty half of the
// key's contract (composer.go, on the Files field: "omitempty keeps a section
// with no reader off the wire entirely"). A pump assembled without a Files
// reader must not put a `files` key on the wire at all - not as null, not as
// an empty object - because key-present-but-empty is not a state this
// dimension has a sentence for.
func Test190ComposerFilesKeyAbsentWithoutReader(t *testing.T) {
	composer := filesKeyOf190(t, PumpSources{})
	if _, on := composer["files"]; on {
		t.Fatalf("no Files reader attached, yet the wire carries a files key: the omitempty promise is broken")
	}
}

// Test190ComposerFilesKeyWireShapeWithReader pins the on-wire half: with a
// Files reader attached, the key is named exactly `files`, sits at the TOP
// level of the composer section (a sibling of git and currentModel, not nested
// inside workspace or git), and its value is the reader's TreeView under the
// seven json names tree_types_190.go declares - no more, no fewer.
func Test190ComposerFilesKeyWireShapeWithReader(t *testing.T) {
	composer := filesKeyOf190(t, PumpSources{
		Files: func() TreeView {
			return TreeView{
				Kind:       TreeKindCounted,
				Reason:     "one file and one folder",
				Entries:    []TreeEntry{{Path: "notes/a.txt", IsDir: false, Size: 3}, {Path: "notes", IsDir: true}},
				FilesTotal: 1,
			}
		},
	})
	raw, on := composer["files"]
	if !on {
		t.Fatalf("Files reader attached, yet the wire carries no top-level files key")
	}
	files, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("files crossed the wire as %T, want an object", raw)
	}
	want := []string{"kind", "reason", "entries", "filesTotal", "filesTruncated", "bytesTotal", "recoverPointer"}
	if len(files) != len(want) {
		t.Fatalf("files object carries %d keys (%v), want exactly %v", len(files), files, want)
	}
	for _, k := range want {
		if _, ok := files[k]; !ok {
			t.Fatalf("files object is missing %q (want exactly %v)", k, want)
		}
	}
	if files["kind"] != TreeKindCounted {
		t.Fatalf("files.kind came back %v, want %q", files["kind"], TreeKindCounted)
	}
	if files["filesTotal"] != float64(1) {
		t.Fatalf("files.filesTotal came back %v, want 1", files["filesTotal"])
	}
	if _, ok := composer["git"]; !ok {
		t.Fatalf("files must sit beside git at the composer top level, but the git key is missing from the wire")
	}
	entries, ok := files["entries"].([]any)
	if !ok || len(entries) != 2 {
		t.Fatalf("files.entries crossed the wire as %T, want a 2-row array", files["entries"])
	}
	first, ok := entries[0].(map[string]any)
	if !ok {
		t.Fatalf("files.entries[0] is %T, want an object", entries[0])
	}
	if first["path"] != "notes/a.txt" || first["isDir"] != false || first["size"] != float64(3) {
		t.Fatalf("files.entries[0] shape drifted: %v (want path/isDir/size)", first)
	}
}

// Test190ComposerFilesKeyCarriesRefusalNotDropped pins the third-state half at
// the wire level, one layer up from the reader pins in the tree_190 family:
// when the reader answers "no workspace was narrowed", the key is STILL on the
// wire carrying kind, reason and an empty entries array. A refusal is a
// sentence on the key, never the key's absence - a page that distinguishes
// "this dimension was not read" from "no tree" needs both halves on the wire.
func Test190ComposerFilesKeyCarriesRefusalNotDropped(t *testing.T) {
	composer := filesKeyOf190(t, PumpSources{
		Files: func() TreeView { return ReadFileTree190("") },
	})
	raw, on := composer["files"]
	if !on {
		t.Fatalf("a not-chosen reading must ride the files key; the wire dropped the key instead")
	}
	files, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("files crossed the wire as %T, want an object", raw)
	}
	if files["kind"] != TreeKindNotChosen {
		t.Fatalf("refusal kind came back %v, want %q", files["kind"], TreeKindNotChosen)
	}
	if reason, _ := files["reason"].(string); reason == "" {
		t.Fatalf("refusal crossed the wire with an empty reason: %v", files)
	}
	entries, ok := files["entries"].([]any)
	if !ok || len(entries) != 0 {
		t.Fatalf("a not-chosen tree must carry an empty entries array on the wire, got %T %v", files["entries"], files["entries"])
	}
}
