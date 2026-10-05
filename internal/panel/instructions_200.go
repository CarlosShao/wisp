package panel

import "github.com/CarlosShao/wisp/internal/projctx"

// Ticket 200 AC#7: the panel carrier for "which project instruction files this
// run actually followed". The panel shows it; how it shows it is the UI leg's
// decision, not this file's. The view type is write-free on purpose, exactly
// like ModeView further up in composer.go: a view model with no setter cannot
// be written through from the panel side, and nothing here gives the panel a
// route back to a permission mode or an allowed-dir list.
//
// WHAT 200-r2 CHANGED, AND WHY IT WAS NOT COSMETIC. At r1 the carrier was a
// bare list with omitempty, and an empty load returned nil - so "you turned the
// feature off", "it is on but the tree holds no instruction file", and "there
// is no loader at all in this host" all marshalled as the SAME absent key. A
// renderer cannot say anything true about an absent key except "nothing", which
// is the 文案在、控件不在 disease ledger A408 names on the other side of the
// wire. The carrier is therefore now an object that always states its status
// once the pump has a reader, and the only way the key stays absent is the only
// way it honestly can: a pump assembled WITHOUT a reader.

// The wire vocabulary for "why is there nothing to show" (ticket 200 AC#7,
// 200-r2 §2(b)). These are the states; they are mutually exclusive, and a
// consumer must never have to guess between two of them.
const (
	// InstructionsStatusLoaded: at least one file was injected this turn.
	InstructionsStatusLoaded = "loaded"
	// InstructionsStatusNoneFound: enabled, looked, the tree and the data dir
	// hold no instruction file. This is NOT the same fact as "off".
	InstructionsStatusNoneFound = "none_found"
	// InstructionsStatusOff: agent.project_instructions_enabled=false.
	InstructionsStatusOff = "off"
	// InstructionsStatusRefused: the host refused to read any directory at all
	// because the workspace coordinate's ticket 102 rewrite account says its
	// spelling was expanded onto another tree (see
	// ProjectInstructionLoadRequest).
	InstructionsStatusRefused = "refused"
	// InstructionsStatusNotRun: a loader is wired but no turn has loaded yet -
	// an absence with a reason, not a silent empty.
	InstructionsStatusNotRun = "not_run"
)

// ProjectInstructionFile is one loaded instruction file, as the panel sees it.
type ProjectInstructionFile struct {
	// Path is the file the loader actually read.
	Path string `json:"path"`
	// Tier is "project" (found by walking up from the workspace) or "global"
	// (Wisp's own data dir, next to config.toml).
	Tier string `json:"tier"`
	// Depth is 0 for the workspace directory itself, larger further out, -1 for
	// the global tier.
	Depth int `json:"depth"`
	// Bytes is what this file contributed before any truncation.
	Bytes int `json:"bytes"`
	// TruncatedBytes is how much the context budget cut (omitted when zero).
	TruncatedBytes int `json:"truncatedBytes,omitempty"`
	// Dropped is true when the budget removed the whole file - the panel must
	// be able to say "this one did not take effect", not pretend it did.
	Dropped bool `json:"dropped,omitempty"`
	// DuplicateOf names the file this path collapsed into: one physical file
	// seen through two directory levels is injected only once.
	DuplicateOf string `json:"duplicateOf,omitempty"`
	// Source is the C25 provenance name stamped on the content ("fs.read").
	Source string `json:"source,omitempty"`
}

// InstructionsSection is the carrier's whole answer: which files, and if there
// are none, why. Status is never empty when the section exists, and Files is
// always an array (empty, never null), so a renderer may index it without a
// nil check and still be told the reason out loud.
type InstructionsSection struct {
	// Status is one of the InstructionsStatus* names above.
	Status string `json:"status"`
	// Reason is the loader's own sentence for a non-loaded state, verbatim from
	// the place that made the decision. Omitted when the run loaded files,
	// because then the list itself is the reason.
	Reason string `json:"reason,omitempty"`
	// Files is this turn's manifest. Empty array when nothing was loaded.
	Files []ProjectInstructionFile `json:"files"`
}

// ProjectInstructionLoadRequest is what the composition root hands the project
// instruction loader: the two directories it may read, plus the refusal sentence
// when it may read neither.
type ProjectInstructionLoadRequest struct {
	// WorkspaceDir is the tree the loader walks UP from ("" = none).
	WorkspaceDir string
	// GlobalDir is Wisp's own data dir, next to config.toml ("" = none).
	GlobalDir string
	// Refused is the non-empty reason nothing at all may be read; the loader
	// prints it instead of loading.
	Refused string
}

// ProjectInstructionLoadRequestFor is the account-aware answer to "which tree
// may the project-instruction loader read?" (ticket 200-r2, ticket 102's rule).
//
// This loader is the one place in the harness that puts a repository's own words
// into the model's system prompt. Handing it a coordinate whose spelling C26's
// expansion step MOVED onto another tree would let a %VAR% or ~ rewrite choose
// which repository talks to the model - the same fail-open ticket 102 closed for
// writes and seals, arriving through the context side instead of the path side.
// So the rewrite account is read HERE, before any path is used, and a rewritten
// view yields no directory at all - not the workspace, not the global tier -
// plus the sentence the host prints and the loader prints once per turn.
//
// MEASURED TODAY (re-measured by ticket 181 AC#7's producer leg, 181-r3; the
// honest note this paragraph carried at 200-r2 is now out of date): a view with
// Rewritten=true IS reachable through the snapshot path, because
// WorkspaceViewFromRoot now renders the C26 account internal/tools records for
// the narrowing in force instead of dropping it. What has NOT changed is the
// panel's own door: RequestWorkspaceSwitch still refuses a rewritten spelling
// at workspace.go through res.Actable(), so no composer request can put this
// branch in force - only a host that names the expanded tree can. That is
// ticket 102's rule, not a gap, and this branch stays non-decorative for the
// same reason it was written: TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll
// runs a real loader over real files and requires it to read neither tree, and
// TestWorkspaceAccountReachesEveryConsumer drives the account all the way from
// a path scope's answer to this refusal.
func ProjectInstructionLoadRequestFor(ws WorkspaceView, globalDir string) ProjectInstructionLoadRequest {
	if ws.Rewritten {
		return ProjectInstructionLoadRequest{
			Refused: "项目说明加载器本轮一份都没有读：这条工作区路径被 C26 的展开步改写过" +
				"（改写账户 rewritten=true），改写后的树 " + ws.Canonical +
				" 不是使用者命名的那棵 " + ws.Spelling + "；说明文件会把所读目录里的字" +
				"写进系统提示，所以既不读这棵工作区、也不读数据目录（票 200-r2／票 102）。",
		}
	}
	req := ProjectInstructionLoadRequest{GlobalDir: globalDir}
	if ws.Set {
		req.WorkspaceDir = ws.Canonical
	}
	return req
}

// ProjectInstructionsFromBundle turns one turn's load result into the file list.
// A bundle that read nothing yields nil here; the caller that builds the wire
// section turns that into an empty array WITH a status, so "nothing" never
// travels without a reason attached to it.
func ProjectInstructionsFromBundle(b *projctx.Bundle) []ProjectInstructionFile {
	if b == nil || len(b.Files) == 0 {
		return nil
	}
	out := make([]ProjectInstructionFile, 0, len(b.Files))
	for _, f := range b.Files {
		out = append(out, ProjectInstructionFile{
			Path:           f.Path,
			Tier:           f.Tier,
			Depth:          f.Depth,
			Bytes:          f.Bytes,
			TruncatedBytes: f.TruncatedBytes,
			Dropped:        f.Dropped,
			DuplicateOf:    f.DuplicateOf,
			Source:         f.Source,
		})
	}
	return out
}

// InstructionsSectionFromBundle is the one mapping from "what the loader did" to
// "what the panel is told". nil (no turn yet) is a state, not an omission.
func InstructionsSectionFromBundle(b *projctx.Bundle) *InstructionsSection {
	if b == nil {
		return &InstructionsSection{
			Status: InstructionsStatusNotRun,
			Reason: "加载器已经接线，但这一轮还没有跑过加载，所以这里是没有读数，不是没有说明文件。",
			Files:  []ProjectInstructionFile{},
		}
	}
	files := ProjectInstructionsFromBundle(b)
	if files == nil {
		files = []ProjectInstructionFile{}
	}
	sec := &InstructionsSection{Files: files, Reason: b.Skipped}
	switch b.SkipReason {
	case projctx.SkipConfigOff:
		sec.Status = InstructionsStatusOff
	case projctx.SkipRewriteRefused:
		sec.Status = InstructionsStatusRefused
	case "":
		if len(files) == 0 {
			sec.Status = InstructionsStatusNoneFound
			sec.Reason = "已启用并查过：工作区逐级向上的每一层目录和数据目录里都没有说明文件（不是被配置关掉的）。"
		} else {
			sec.Status = InstructionsStatusLoaded
		}
	default:
		// An unknown skip name is its own finding: state it, never fold it into
		// one of the five the renderer knows how to draw.
		sec.Status = InstructionsStatusNoneFound
		sec.Reason = "加载器报了一个本包不认识的跳过原因：" + b.SkipReason +
			"（按「没有说明文件」显示，原因逐字带在上面）。"
		if b.Skipped != "" {
			sec.Reason = b.Skipped + " 跳过原因名=" + b.SkipReason
		}
	}
	return sec
}

// WithInstructions returns s carrying the loaded-instruction list. Snapshots
// are value types handed to the pump, so this is a copy, not a mutation; the
// list itself is copied too, because projctx hands out its live slice and a
// snapshot that aliased it would let a later turn rewrite an earlier packet.
// An empty list is NOT an absent key: it is the none_found section.
func WithInstructions(s Snapshot, files []ProjectInstructionFile) Snapshot {
	if len(files) == 0 {
		s.Instructions = InstructionsSectionFromBundle(&projctx.Bundle{})
		return s
	}
	s.Instructions = &InstructionsSection{
		Status: InstructionsStatusLoaded,
		Files:  append([]ProjectInstructionFile(nil), files...),
	}
	return s
}
