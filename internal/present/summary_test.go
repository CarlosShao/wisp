package present

import (
	"strings"
	"testing"
)

// The same-response summary protocol: the marker line is lifted out of the
// display, the fallback is the first sentence, and the 60-rune ceiling is
// enforced on the summary - never by truncating the artifact path.

func TestExtractSummaryTakesTheTrailingMarkerLine(t *testing.T) {
	text := "正文第一段。\n正文第二段。\n口播摘要：结果已整理成表"
	display, summary := ExtractSummary(text)
	if summary != "结果已整理成表" {
		t.Errorf("summary = %q, want the marker payload", summary)
	}
	if display != "正文第一段。\n正文第二段。" {
		t.Errorf("display = %q, want the marker line stripped", display)
	}
	if strings.Contains(display, spokenMarker) {
		t.Errorf("display still carries the marker: %q", display)
	}
}

func TestExtractSummaryAcceptsHalfWidthColonAndTrailingBlanks(t *testing.T) {
	display, summary := ExtractSummary("正文。\n口播摘要: 简短\n\n")
	if summary != "简短" || display != "正文。" {
		t.Errorf("got display=%q summary=%q, want 正文。/简短", display, summary)
	}
}

func TestExtractSummaryIgnoresAMidTextMarker(t *testing.T) {
	// Only a TRAILING marker is protocol: a sentence that merely contains
	// the phrase stays content, or the router would eat user prose.
	text := "开头。\n口播摘要：这行是正文的一部分\n结尾。"
	display, summary := ExtractSummary(text)
	if summary != "" {
		t.Errorf("summary = %q, want empty: the marker is not the last line", summary)
	}
	if display != text {
		t.Errorf("display = %q, want the text untouched", display)
	}
}

func TestExtractSummaryWithoutMarkerLeavesEverything(t *testing.T) {
	display, summary := ExtractSummary(" plain answer ")
	if summary != "" || display != " plain answer " {
		t.Errorf("got display=%q summary=%q, want the text back with no summary", display, summary)
	}
}

func TestExtractSummaryOnMarkerOnlyTextLeavesNothingToPresent(t *testing.T) {
	display, _ := ExtractSummary("口播摘要：只有摘要")
	if tier, _ := Classify(display); tier != TierNone {
		t.Errorf("Classify(display) = %s, want none: the response carried nothing else", tier)
	}
}

func TestFirstSentenceCutsAndClamps(t *testing.T) {
	if got := FirstSentence("第一句。第二句。第三句", MaxSpokenRunes); got != "第一句" {
		t.Errorf("FirstSentence = %q, want the first sentence", got)
	}
	if got := FirstSentence("没有句读的一句话", MaxSpokenRunes); got != "没有句读的一句话" {
		t.Errorf("FirstSentence = %q, want the whole text", got)
	}
	long := strings.Repeat("字", 80)
	if got := FirstSentence(long, MaxSpokenRunes); len([]rune(got)) != MaxSpokenRunes {
		t.Errorf("FirstSentence length = %d, want the %d-rune ceiling", len([]rune(got)), MaxSpokenRunes)
	}
	if got := FirstSentence("换行\n算断句", MaxSpokenRunes); got != "换行" {
		t.Errorf("FirstSentence = %q, want the pre-newline head", got)
	}
}

func TestSpokenWithFileKeepsTheWholePathAndClampsOnlyTheSummary(t *testing.T) {
	path := `C:\Users\swq\AppData\Roaming\wisp\artifacts\result-task-abcd1234.md`
	got := spokenWithFile("一句话摘要", "正文", path, true)
	if !strings.HasSuffix(got, path) {
		t.Errorf("spoken = %q, want it to end with the full path: a truncated path is a path the user cannot open", got)
	}
	if !strings.Contains(got, spokenHaveFile) {
		t.Errorf("spoken = %q, want the saved tail", got)
	}
	if runes := len([]rune(got)) - len([]rune(path)) - len([]rune(spokenHaveFile)); runes > MaxSpokenRunes {
		t.Errorf("summary part = %d runes, over the %d ceiling", runes, MaxSpokenRunes)
	}

	// A runaway marker summary is clamped, not trusted.
	runaway := strings.Repeat("长", 200)
	clamped := spokenWithFile(runaway, "正文", path, true)
	summaryPart := len([]rune(clamped)) - len([]rune(path)) - len([]rune(spokenHaveFile))
	if summaryPart != MaxSpokenRunes {
		t.Errorf("clamped summary part = %d runes, want exactly %d", summaryPart, MaxSpokenRunes)
	}

	if got := spokenWithFile("摘要", "正文", "", false); !strings.HasSuffix(got, spokenNoFile) {
		t.Errorf("failed save spoken = %q, want the honest could-not-save tail", got)
	}
}
