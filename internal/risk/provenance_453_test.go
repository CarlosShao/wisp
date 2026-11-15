package risk

// Ticket 453: the C25 source name a panel attachment is marked under. Two
// judgements live here, both about the DATA TABLE and the ENGINE rather than
// about a source string:
//
//  1. IsSensitiveSource(panel.attachment) is true. The roster is what tells the
//     rest of the tree this name is a SPEC-06 §5-class source; if the entry is
//     removed, this case reddens with the roster's own answer (false), not with a
//     missing literal somewhere.
//  2. A mark under that name is reported BY that name through the egress check.
//     This is the behaviour a mis-wired source would break: Mark under a
//     different name and Inspect answers a different SrcTool.
//
// SPEC-06 §5's enumeration does not name panel.attachment yet; adding that
// sentence is Q-456 (owner). This file deliberately does not touch docs/specs.

import "testing"

func Test453PanelAttachmentIsOnTheSensitiveSourceRoster(t *testing.T) {
	if !IsSensitiveSource(SrcPanelAttachment) {
		t.Fatalf("IsSensitiveSource(%q) = false: a panel attachment is outside content (ticket 214 AC#5) and the roster must say so", SrcPanelAttachment)
	}
	if SrcPanelAttachment != "panel.attachment" {
		t.Fatalf("SrcPanelAttachment = %q, want the contract spelling %q", SrcPanelAttachment, "panel.attachment")
	}
}

func Test453PanelAttachmentMarksAreReportedAsThatSource(t *testing.T) {
	prov := NewProvenance(ProvOptions{NoProbe: true})
	if scope := prov.OpenScope("task-453-roster"); scope == nil {
		t.Fatal("OpenScope refused the task scope")
	}
	const secret = "附件里的这一段文字不应该被无标记地带出网"
	if !prov.Mark("task-453-roster", SrcPanelAttachment, "artifact:attachment-abc.png", secret) {
		t.Fatal("Mark refused a source on the roster: the engine and the table disagree")
	}
	hit, ok := prov.Inspect("task-453-roster", "notify", map[string]any{"text": secret})
	if !ok {
		t.Fatalf("Inspect on the notify channel found no taint after a %q mark", SrcPanelAttachment)
	}
	if hit.SrcTool != SrcPanelAttachment {
		t.Fatalf("hit source = %q, want %q: the egress verdict must name the source the attachment was marked under", hit.SrcTool, SrcPanelAttachment)
	}
}
