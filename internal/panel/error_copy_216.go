package panel

// Ticket 216 AC#5: error sentences branch on the error CODE, never on
// substrings of the English error text. The forbidden shape is openchamber's
// ChatInput.tsx:2205 (three-substring OR including "413" deciding which copy
// to show) - docs/reports/survey-2026-09-28-composer-plus-menu.md:378-381
// transcribed it verbatim; this package must never grow one.
//
// The code source is the D37 taxonomy that already exists:
// observe.Error.Class (17 enum constants) extracted from an error chain by
// observe.ClassOf (errors.As). One class selects one leading sentence; the
// raw error text still follows it, because "读不到" without the cause is the
// sentence ticket 147 refused. Unclassified errors (ClassOf ok=false) keep
// the pre-existing plain sentence byte-for-byte - this layer adds a branch,
// it does not rewrite the old one.
//
// This is the consumption layer ticket 216-v2 judged missing: the hook
// (observe.ErrorClass.MessageKey, internal/observe/errors.go) existed with
// zero production consumers. Each branch key here is the MessageKey's class
// constant, not a string found inside err.Error().

import (
	"github.com/CarlosShao/wisp/internal/observe"
)

// errorCopy216 maps one D37 error class to the one leading sentence the panel
// shows before the raw error text. Keys are the 17 enum constants; the map is
// exhaustive over observe.AllClasses() and errorCopy216Complete pins that, so
// a class added to the taxonomy cannot silently fall out of the copy table.
var errorCopy216 = map[observe.ErrorClass]string{
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

// errorCopy216Complete reports whether every class in the taxonomy has a copy
// row. It is the ruler that keeps a future 18th class from silently falling
// out of the table (the test pins it).
func errorCopy216Complete() bool {
	for _, c := range observe.AllClasses() {
		if _, ok := errorCopy216[c]; !ok {
			return false
		}
	}
	return true
}

// errorSentence216 is the branch point AC#5 requires: class in, sentence out,
// via observe.ClassOf (errors.As over the chain). err.Error() is appended
// verbatim AFTER the branch, as the cause, never used to pick the sentence.
func errorSentence216(err error) string {
	if err == nil {
		return ""
	}
	if class, ok := observe.ClassOf(err); ok {
		if sentence, ok := errorCopy216[class]; ok {
			return sentence + err.Error()
		}
	}
	// Unclassified errors keep the plain sentence - the pre-AC#5 shape, byte
	// for byte, so existing fixtures do not move.
	return err.Error()
}

// userFacingErrorSentence216 is the one page-visible error sentence this
// package renders for a refused composer request (RefusedEnvelopeForUser's
// cause half). A classified error gains its class's leading sentence; an
// unclassified one is verbatim, byte for byte.
func userFacingErrorSentence216(err error) string {
	return errorSentence216(err)
}
