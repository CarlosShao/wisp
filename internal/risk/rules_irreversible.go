package risk

import "strings"

// irreversibleClassReason maps the frozen R8 irreversibility classes
// (SPEC-06 §3 R8) to their card-renderable descriptions. The class set is
// contract: adding a class is a contract change.
var irreversibleClassReason = map[string]string{
	"delete":       "永久删除",
	"overwrite":    "覆盖已有内容",
	"power":        "关机/重启",
	"close-window": "关闭窗口",
	"send":         "发送类外发操作",
}

// ruleIrreversible implements R8: permanent delete / overwrite-existing /
// power off or restart / close-window / send-class operations -> L2.
// OverwriteExisting is the structured spelling of the "overwrite" class.
// An UNKNOWN class string fails closed to L2 — an unrecognized
// irreversibility claim must never be treated as reversible.
func ruleIrreversible(ctx *assessCtx) *contribution {
	classes := make([]string, 0, len(ctx.facts.Irreversible)+1)
	classes = append(classes, ctx.facts.Irreversible...)
	if ctx.facts.OverwriteExisting {
		classes = append(classes, "overwrite")
	}
	if len(classes) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(classes))
	var parts []string
	for _, class := range classes {
		if seen[class] {
			continue
		}
		seen[class] = true
		if desc, ok := irreversibleClassReason[class]; ok {
			parts = append(parts, desc)
			continue
		}
		parts = append(parts, "未知不可逆类别（fail-closed）: "+class)
	}
	return &contribution{
		rules:  []RuleID{R8},
		level:  L2,
		reason: "R8: 不可逆操作（" + strings.Join(parts, "、") + "）",
	}
}
