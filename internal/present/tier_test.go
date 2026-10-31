package present

import (
	"strings"
	"testing"
)

// Ticket 30 AC#1: the tier matrix. Crafted outputs for every tier, the
// boundary values D10 names (60/61, 400/401, 800) and the structured-beats-
// length precedence cases, including SPEC-05 §9's "40-character table".

func repeatHanzi(n int) string { return strings.Repeat("测", n) }

func TestClassifySendsEveryTierWhereD10PutsIt(t *testing.T) {
	cases := []struct {
		name string
		text string
		want Tier
		rule string
	}{
		{"empty", "", TierNone, "empty result"},
		{"blank", "   \n\t ", TierNone, "empty result"},
		{
			"60 hanzi stays short",
			repeatHanzi(59) + "字",
			TierShort, "at or under 60",
		},
		{
			"61 hanzi leaves short",
			repeatHanzi(61),
			TierMedium, "61-400",
		},
		{
			"400 hanzi stays medium",
			repeatHanzi(400),
			TierMedium, "61-400",
		},
		{
			"401 hanzi goes to file",
			repeatHanzi(401),
			TierLong, "over the 400",
		},
		{
			"800 latin chars stay short",
			strings.Repeat("a", 800),
			TierShort, "at or under 60",
		},
		{
			"801 latin chars go to file",
			strings.Repeat("a", 801),
			TierLong, "over the 800",
		},
		{
			"whitespace padding cannot push to file",
			strings.Repeat("a", 700) + strings.Repeat(" ", 200),
			TierShort, "at or under 60",
		},
		{
			"40-char table beats short length",
			"| 名称 | 数量 |\n|---|---|\n| 甲 | 1 |",
			TierStructured, "markdown table",
		},
		{
			"code fence beats long length",
			"```\n" + strings.Repeat("x", 900) + "\n```",
			TierStructured, "code fence",
		},
		{
			"three-item list beats short length",
			"- 甲\n- 乙\n- 丙",
			TierStructured, "list of 3 items",
		},
		{
			"two links beat short length",
			"看 [文档](https://a.example) 和 [指南](https://b.example) 再说",
			TierStructured, "2 links",
		},
		{
			"one markdown link is not a link list",
			"看 [文档](https://a.example) 就够",
			TierShort, "at or under 60",
		},
		{
			"one link plus one bare url counts two",
			"[文档](https://a.example) 或 https://b.example",
			TierStructured, "2 links",
		},
		{
			"two path lines beat short length",
			"C:\\data\\报告.docx\nD:\\out\\汇总.pdf",
			TierStructured, "2 path lines",
		},
		{
			"a single path inside prose stays short",
			"文件在 C:\\x\\a.txt 里",
			TierShort, "at or under 60",
		},
		{
			"tilde fence counts as a fence",
			"~~~\nraw\n~~~",
			TierStructured, "code fence",
		},
		{
			"numbered list of three counts",
			"1. 甲\n2. 乙\n3. 丙",
			TierStructured, "list of 3 items",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, rule := Classify(tc.text)
			if got != tc.want {
				t.Fatalf("Classify(%q) = %s (%s), want %s", truncFor(t, tc.text), got, rule, tc.want)
			}
			if !strings.Contains(rule, tc.rule) {
				t.Errorf("rule = %q, want it to name %q", rule, tc.rule)
			}
		})
	}
}

// First match wins: the ladder checks structure before length, so a long
// text that also carries a marker is rule 1 with the rule-1 channels.
func TestClassifyStructuredBeatsLengthBothWays(t *testing.T) {
	longTable := repeatHanzi(500) + "\n| 名称 | 数量 |\n|---|---|"
	got, rule := Classify(longTable)
	if got != TierStructured || !strings.Contains(rule, "table") {
		t.Fatalf("long text with a table = %s (%s), want structured: the 40-char-table rule beats length", got, rule)
	}
	longList := strings.Repeat("- 项\n", 5) + repeatHanzi(500)
	got, rule = Classify(longList)
	if got != TierStructured || !strings.Contains(rule, "list") {
		t.Fatalf("long text with a list = %s (%s), want structured", got, rule)
	}
}

func truncFor(t *testing.T, s string) string {
	t.Helper()
	r := []rune(s)
	if len(r) > 24 {
		return string(r[:24]) + "..."
	}
	return s
}
