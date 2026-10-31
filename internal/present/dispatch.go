package present

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
)

// The channel ports. Each one is the narrowest face the router needs, so the
// host adapters stay thin and the real implementations keep their own
// semantics:
//
//   - Speaker hands the line to the voice leg's announce queue (D38b owner
//     goroutine). The ChTTS exfil gate lives on that announce path in the
//     host (ticket 26 wiring) and is NOT duplicated here.
//   - Notifier is the system toast. Its failure must never silence the rest
//     of tier 3 (D42#8).
//   - Badger is the ball badge and TaskLister the panel task list entry.
//     Tier 3 drives both unconditionally - Focus Assist may swallow the
//     toast, so the dual-channel update is the degradation path, never a
//     fallback that can be skipped.
//   - PanelSink shows the finished result; artifact "" means the text is the
//     whole result. The C17 bridge itself is untouched by this package.
//   - ClipboardWriter copies the tier-4 result. The adapter MUST run
//     risk.Provenance.CheckText(scope, risk.ChClipboard, text) before the
//     real write (the same gate the clipboard tool carries); this package
//     cannot run it because the provenance scope is a host fact.
//   - Artifacts is satisfied by *memory.Store as-is: PutArtifact carries the
//     name guard, the private DACL and the on-write 500MB quota check, so
//     the router can never write a second, unaccounted file tree.
type Speaker interface {
	// Speak queues one line for announcement; false = not taken (voice
	// chain not assembled or queue full). It does not block for playback.
	Speak(text string) bool
}

type Notifier interface {
	Notify(title, body string) error
}

type Badger interface {
	Badge(text string)
}

type TaskLister interface {
	TaskEntry(text string)
}

type PanelSink interface {
	ShowResult(text, artifact string)
}

type ClipboardWriter interface {
	WriteClipboard(text string) error
}

type Artifacts interface {
	PutArtifact(ctx context.Context, name string, data []byte) error
	ArtifactsDir() string
}

// Channels is the router's full port set. A nil port is a channel this host
// did not wire: Dispatch records it as skipped and moves on.
type Channels struct {
	Speak     Speaker
	Notify    Notifier
	Badge     Badger
	TaskList  TaskLister
	Panel     PanelSink
	Clipboard ClipboardWriter
	Artifacts Artifacts
	// Logger receives the router's diagnostic lines (nil = stay quiet). The
	// artifacts line carries the file path, which is the "path written to
	// the log" half of ticket 30's AC.
	Logger *slog.Logger
}

// Dispatcher routes finished results through Channels.
type Dispatcher struct {
	ch Channels
}

// New builds a router over the given ports.
func New(ch Channels) *Dispatcher { return &Dispatcher{ch: ch} }

// Input is one finished task's routable slice (the fields of agent.Result
// this layer may read; the host maps them, so this package does not import
// the loop).
type Input struct {
	// TaskID names the task in receipts and artifact names. "" is allowed:
	// the artifact name then falls back to its content hash.
	TaskID string
	// Text is the final assistant text (agent.Result.Text). The router
	// classifies and presents THIS text; failures and control words never
	// reach it because the host only routes a completed result.
	Text string
}

// Note is one channel's outcome on one dispatch. Status is one of
// "sent", "skipped", "failed"; Detail carries the reason or the error word.
type Note struct {
	Channel string
	Status  string
	Detail  string
}

// Receipt is what one dispatch did: the tier and the rule that matched, the
// lines the user gets on each surface, the artifact path, and the
// per-channel trace. OK is true when no channel failed; skipped channels do
// not spoil it, because "not wired" is a host fact, not a routing failure.
type Receipt struct {
	Tier     Tier
	Rule     string
	Spoken   string
	Display  string
	Artifact string
	Notes    []Note
	OK       bool
}

// ok reports whether any note failed.
func (r *Receipt) ok() bool {
	for _, n := range r.Notes {
		if n.Status == "failed" {
			return false
		}
	}
	return true
}

// Dispatch routes one finished result. It never starts a goroutine and never
// blocks on playback: Speak only queues. The channel call order per tier is
// fixed so the receipt trace is comparable across runs.
func (d *Dispatcher) Dispatch(ctx context.Context, in Input) Receipt {
	display, summary := ExtractSummary(in.Text)
	tier, rule := Classify(display)
	rc := Receipt{Tier: tier, Rule: rule, Display: display}
	d.route(ctx, in, display, summary, &rc)
	rc.OK = rc.ok()
	if d.ch.Logger != nil {
		d.ch.Logger.Info("present: routed",
			"task", in.TaskID, "tier", string(tier), "rule", rule,
			"ok", rc.OK, "artifact", rc.Artifact)
	}
	return rc
}

func (d *Dispatcher) route(ctx context.Context, in Input, display, summary string, rc *Receipt) {
	switch rc.Tier {
	case TierNone:
		rc.Notes = append(rc.Notes, Note{
			Channel: "all", Status: "skipped",
			Detail: rc.Rule,
		})
	case TierStructured, TierLong:
		path, saved := d.writeArtifact(ctx, in.TaskID, display, rc)
		if d.ch.Panel != nil {
			d.ch.Panel.ShowResult(display, path)
			rc.Notes = append(rc.Notes, Note{Channel: "panel", Status: "sent"})
		} else {
			rc.Notes = append(rc.Notes, Note{
				Channel: "panel", Status: "skipped",
				Detail: "not wired",
			})
		}
		d.speak(rc, spokenWithFile(summary, display, path, saved))
	case TierMedium:
		// D42#8: notify, badge and task list are one unconditional
		// dual-channel update. A toast failure never gates the badge, and a
		// nil port is a skip, not an excuse to drop the siblings.
		if d.ch.Notify != nil {
			if err := d.ch.Notify.Notify("Wisp", display); err != nil {
				rc.Notes = append(rc.Notes, Note{
					Channel: "notify", Status: "failed",
					Detail: err.Error(),
				})
			} else {
				rc.Notes = append(rc.Notes, Note{Channel: "notify", Status: "sent"})
			}
		} else {
			rc.Notes = append(rc.Notes, Note{
				Channel: "notify", Status: "skipped",
				Detail: "not wired",
			})
		}
		if d.ch.Badge != nil {
			d.ch.Badge.Badge(FirstSentence(display, MaxSpokenRunes))
			rc.Notes = append(rc.Notes, Note{Channel: "badge", Status: "sent"})
		} else {
			rc.Notes = append(rc.Notes, Note{
				Channel: "badge", Status: "skipped",
				Detail: "not wired",
			})
		}
		if d.ch.TaskList != nil {
			d.ch.TaskList.TaskEntry(display)
			rc.Notes = append(rc.Notes, Note{Channel: "task", Status: "sent"})
		} else {
			rc.Notes = append(rc.Notes, Note{
				Channel: "task", Status: "skipped",
				Detail: "not wired",
			})
		}
		if d.ch.Panel != nil {
			d.ch.Panel.ShowResult(display, "")
			rc.Notes = append(rc.Notes, Note{Channel: "panel", Status: "sent"})
		} else {
			rc.Notes = append(rc.Notes, Note{
				Channel: "panel", Status: "skipped",
				Detail: "not wired",
			})
		}
		d.speak(rc, FirstSentence(display, MaxSpokenRunes))
	case TierShort:
		// The clipboard write lives in this branch and nowhere else: tiers
		// 1-3 must never touch the user's clipboard (D10).
		if d.ch.Clipboard != nil {
			if err := d.ch.Clipboard.WriteClipboard(display); err != nil {
				rc.Notes = append(rc.Notes, Note{
					Channel: "clipboard", Status: "failed",
					Detail: err.Error(),
				})
			} else {
				rc.Notes = append(rc.Notes, Note{Channel: "clipboard", Status: "sent"})
			}
		} else {
			rc.Notes = append(rc.Notes, Note{
				Channel: "clipboard", Status: "skipped",
				Detail: "not wired",
			})
		}
		d.speak(rc, display)
	}
}

// speak queues the spoken line when the voice chain is wired.
func (d *Dispatcher) speak(rc *Receipt, line string) {
	rc.Spoken = line
	if line == "" {
		rc.Notes = append(rc.Notes, Note{
			Channel: "speak", Status: "skipped",
			Detail: "nothing to say",
		})
		return
	}
	if d.ch.Speak != nil {
		if !d.ch.Speak.Speak(line) {
			rc.Notes = append(rc.Notes, Note{
				Channel: "speak", Status: "failed",
				Detail: "announce not taken (voice chain not assembled or queue full)",
			})
		} else {
			rc.Notes = append(rc.Notes, Note{Channel: "speak", Status: "sent"})
		}
	} else {
		rc.Notes = append(rc.Notes, Note{
			Channel: "speak", Status: "skipped",
			Detail: "not wired",
		})
	}
}

// writeArtifact lands the display text in the artifacts directory through
// the store's own face - the same guarded, quota-accounted route every other
// artifact takes. The receipt and the log line carry the full path.
func (d *Dispatcher) writeArtifact(ctx context.Context, taskID, display string, rc *Receipt) (string, bool) {
	if d.ch.Artifacts == nil {
		rc.Notes = append(rc.Notes, Note{
			Channel: "artifact", Status: "skipped",
			Detail: "not wired",
		})
		return "", false
	}
	name := ArtifactName(taskID, display)
	if err := d.ch.Artifacts.PutArtifact(ctx, name, []byte(display)); err != nil {
		rc.Notes = append(rc.Notes, Note{
			Channel: "artifact", Status: "failed",
			Detail: err.Error(),
		})
		return "", false
	}
	path := filepath.Join(d.ch.Artifacts.ArtifactsDir(), name)
	rc.Artifact = path
	rc.Notes = append(rc.Notes, Note{Channel: "artifact", Status: "sent"})
	if d.ch.Logger != nil {
		d.ch.Logger.Info("present: artifact written",
			"task", taskID, "path", path, "bytes", len(display))
	}
	return path, true
}

// ArtifactName builds one bare artifacts file name for a result: the task id
// (sanitized) plus a content hash, so two different results never share a
// name and a re-read of the same result is idempotent. The name passes
// memory.ValidArtifactName: no separators, no trailing dot or space.
func ArtifactName(taskID, content string) string {
	id := sanitizeTaskID(taskID)
	sum := sha256.Sum256([]byte(content))
	return fmt.Sprintf("result-%s-%s.md", id, hex.EncodeToString(sum[:4]))
}

// sanitizeTaskID keeps the name-guard's charset: bare name, no path
// separators or colons, no trailing dot or space.
func sanitizeTaskID(taskID string) string {
	var b strings.Builder
	for _, r := range taskID {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
		if b.Len() >= 40 {
			break
		}
	}
	if b.Len() == 0 {
		return "task"
	}
	return b.String()
}
