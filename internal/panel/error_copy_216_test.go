package panel

// Ticket 216 AC#5's instruments. The branch point errorSentence216 / the
// dispatch refusal sentence now ride on observe.ClassOf (the D37 class pulled
// by errors.As over the chain); these tests pin the three things the ticket
// names:
//
//  1. THE CODE→SENTENCE TABLE (对照测试): one error CODE selects one leading
//     sentence - 17 classes, 17 distinct sentences, picked by the class
//     constant, never by anything found inside err.Error().
//  1b. THE WIRE: the dispatch refusal sentence (RefusedEnvelopeForUser, the
//     only page-visible error sentence this package renders for a refused
//     composer request) carries the class's own leading sentence when the
//     chain carries a D37 class, through the real config.set write route
//     whose ApplySetting errors are observe.ClassConfig wrapped (%w) upstream
//     (internal/config/sectionwrite.go), so the branch is reached through
//     production wiring, not a decoration.
//  2. THE SUBSTRING FORBIDDEN SHAPE (正控): the same error text fed under 17
//     different classes must yield 17 DIFFERENT leading sentences - a
//     substring-of-err.Error() implementation collapses those 17 to one
//     identical branch (the openchamber ChatInput.tsx:2205 "413" shape) and
//     goes red here. This is the discriminating half: an implementation that
//     branches on the error string cannot satisfy both this and test 1.
//
// The negative half AC#5 forbids is already scanned by the a2/v2 re-runs
// (zero substring-branch sites in panel production code); this file adds the
// positive branch and the instrument that keeps it that way.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// wantHead216 is the test's OWN copy of the 17 expected leading sentences -
// one Go string literal per class, written out by hand, never read back from
// the production table errorCopy216. The production table and this table are
// two independent artifacts: change the production sentence without changing
// this one and TestErrorCodeSelectsTheSentence216 goes red. The old shape
// (want := errorCopy216[c] + err.Error()) was a tautology - both sides of the
// comparison read the same table, so mutating the table mutated the
// expectation and the ruler had no teeth (ticket 216 AC#5 re-run finding).
var wantHead216 = map[observe.ErrorClass]string{
	observe.ClassConfig:           "配置这一步出了问题，先看下面这句具体原因：",
	observe.ClassAuth:             "身份或密钥没有通过，先检查密钥配置：",
	observe.ClassNetwork:          "网络这一步不通，先确认机器能连上外网：",
	observe.ClassRateLimit:        "上游限流了，稍等再试：",
	observe.ClassProvider:         "上游服务这一步出错，具体原因见下面这句：",
	observe.ClassModel:            "模型这一步出了问题，具体原因见下面这句：",
	observe.ClassAudioDevice:      "音频设备这一步出了问题，先检查麦克风与扬声器：",
	observe.ClassASR:              "语音识别这一步出了问题，具体原因见下面这句：",
	observe.ClassTool:             "工具这一步执行出错，具体原因见下面这句：",
	observe.ClassPermissionDenied: "这一步被权限拦下了，先看需要哪条许可：",
	observe.ClassUserRejected:     "这一步被你拒绝了，所以没有继续：",
	observe.ClassCancelled:        "这一步被取消了，所以没有继续：",
	observe.ClassBudget:           "预算这一步用完了，先确认额度再继续：",
	observe.ClassLoop:             "任务这一步卡住了，先看下面这句具体卡在哪：",
	observe.ClassInjection:        "这一步被判定为疑似注入，已按安全规则拦下：",
	observe.ClassInternal:         "程序内部这一步出了错，具体原因见下面这句：",
	observe.ClassResource:         "资源这一步超限了，先看下面这句具体是哪一项：",
}

// TestErrorCodeSelectsTheSentence216 is the code→sentence table test. Each of
// the 17 D37 classes feeds the SAME error text; the HEAD (the class's leading
// sentence, errorHead216) must equal the test's own literal for that class -
// an expectation that does not live in the production table.
func TestErrorCodeSelectsTheSentence216(t *testing.T) {
	// The same underlying text for every class - the branch key must be the
	// class, never this text.
	const sameDetail = "upstream 413 payload too large"
	for _, c := range observe.AllClasses() {
		err := observe.New(c, sameDetail)
		got := errorHead216(err)
		want := wantHead216[c]
		if got != want {
			t.Errorf("class %s: head = %q, want literal %q", c, got, want)
		}
		// The whole sentence keeps the same shape it always had: the head
		// followed by the cause verbatim.
		if whole := errorSentence216(err); whole != want+err.Error() {
			t.Errorf("class %s: sentence = %q, want head+cause %q", c, whole, want+err.Error())
		}
	}
	// The table must be exhaustive over the taxonomy: an 18th class added
	// later without a copy row would fall back silently - this ruler makes
	// that red, so a class cannot drift out of the table.
	if complete := errorCopy216Complete(); !complete {
		t.Errorf("errorCopy216 is not exhaustive over observe.AllClasses()")
	}
}

// TestUnclassifiedKeepsThePlainSentence216: ClassOf ok=false keeps the
// pre-AC#5 plain sentence byte-for-byte, so the unclassified path (the shape
// every existing panel fixture feeds) does not move.
func TestUnclassifiedKeepsThePlainSentence216(t *testing.T) {
	plain := errors.New("配置读不了")
	if got := errorSentence216(plain); got != plain.Error() {
		t.Errorf("unclassified error = %q, want the plain text verbatim %q", got, plain.Error())
	}
	// A wrapped unclassified chain stays unclassified: errors.As finds no
	// observe.Error anywhere in it.
	wrapped := fmt.Errorf("panel: 配置读取失败（requestId=%q）：%w", "r-1", plain)
	if _, ok := observe.ClassOf(wrapped); ok {
		t.Fatalf("the wrapped chain suddenly classifies; the test premise broke")
	}
	if got := errorSentence216(wrapped); got != wrapped.Error() {
		t.Errorf("wrapped unclassified = %q, want verbatim %q", got, wrapped.Error())
	}
}

// TestErrorCopySentencesAreAllDistinct216 is the discriminating half: the same
// error text under 17 different classes must yield 17 DIFFERENT HEADS. The
// old shape compared whole sentences (head + err.Error()); err.Error() embeds
// the class name, so every class differed even when every head was identical,
// and a table that collapsed two classes onto one sentence slipped through.
// Comparing errorHead216 alone removes that blindfold. The prefix self-check
// at the bottom additionally forbids any head being a prefix of (or equal to)
// another - a "distinct" pair that only differs inside the appended cause can
// no longer masquerade as separation.
func TestErrorCopySentencesAreAllDistinct216(t *testing.T) {
	const sameDetail = "upstream 413 payload too large"
	seen := map[observe.ErrorClass]string{}
	for _, c := range observe.AllClasses() {
		seen[c] = errorHead216(observe.New(c, sameDetail))
	}
	classes := observe.AllClasses()
	for i, a := range classes {
		for _, b := range classes[i+1:] {
			if seen[a] == seen[b] {
				t.Errorf("classes %s and %s produce the same head %q - the branch no longer separates codes", a, b, seen[a])
			}
			if strings.HasPrefix(seen[a], seen[b]) {
				t.Errorf("head of %s is a prefix of head of %s (%q) - the heads do not separate on their own", b, a, seen[b])
			}
			if strings.HasPrefix(seen[b], seen[a]) {
				t.Errorf("head of %s is a prefix of head of %s (%q) - the heads do not separate on their own", a, b, seen[a])
			}
		}
	}
}

// TestRefusedEnvelopeCarriesTheClassSentence216 is the wire half: a real
// config.set write through the production dispatch, with the leg returning a
// classified error wrapped exactly as internal/config's ApplySetting wraps it
// (observe.ClassConfig + %w through fmt.Errorf), must come back to the page
// with the class's own leading sentence in the refusal.
func TestRefusedEnvelopeCarriesTheClassSentence216(t *testing.T) {
	legErr := fmt.Errorf("%w", observe.New(observe.ClassConfig, "role_chat_thinking 传入了名册外的档位"))
	leg := &cfg248Leg{applyEr: legErr}
	dispatch := cfg248Dispatch(leg)
	raw := cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-216", "source": ComposerRequestSource,
		"configField": FieldRoleChatThinking, "configValue": "high",
	})
	reply, err := dispatch.Handle(context.Background(), raw)
	if err == nil {
		t.Fatalf("a refused write came back nil-error; the door did not refuse")
	}
	if leg.applyCalls != 1 {
		t.Errorf("leg apply calls = %d, want 1 - the refusal was made before the leg", leg.applyCalls)
	}
	want := errorCopy216[observe.ClassConfig]
	if !strings.Contains(reply, want) {
		t.Errorf("refusal reply %q does not carry the config-class sentence %q", reply, want)
	}
	// And the raw cause still follows - the branch adds a head, it never
	// replaces the cause (ticket 147's "读不到 without a name" rule).
	if !strings.Contains(reply, "名册外的档位") {
		t.Errorf("refusal reply %q lost the underlying cause", reply)
	}
}
