package risk

import "fmt"

// ruleDeclared implements R1: the tool manifest's declared RiskLevel is the
// lower bound of the decision, never the conclusion. A declaration above L0
// raises the floor; an out-of-range declaration is invalid input and
// fail-closes to R9 L2.
func ruleDeclared(ctx *assessCtx) *contribution {
	d := ctx.facts.Declared
	if d > Deny || d < L0 {
		return &contribution{
			rules:  []RuleID{R9},
			level:  L2,
			reason: fmt.Sprintf("R9: 工具声明风险级非法（%d），fail-closed 升 L2", int(d)),
		}
	}
	if d <= L0 {
		return nil
	}
	return &contribution{
		rules:  []RuleID{R1},
		level:  d,
		reason: fmt.Sprintf("R1: 工具声明为下界（%s）", d),
	}
}

// rulePathAllowlist implements R2: every affected path must land inside the
// authorized directories after C26 canonicalization; out of bounds -> L2.
// Dormant until ticket 18 wires the PathCanonicalizer. A WIRED resolver that
// errors on a path fail-closes to L2 (unresolvable = unjudgeable = out).
func rulePathAllowlist(ctx *assessCtx) *contribution {
	if len(ctx.facts.Paths) == 0 || ctx.canon == nil {
		return nil
	}
	for _, raw := range ctx.facts.Paths {
		canonical, err := ctx.canon.Canonicalize(raw)
		if err != nil {
			return &contribution{
				rules:  []RuleID{R2},
				level:  L2,
				reason: fmt.Sprintf("R2: 路径无法规范化，按越界处理（fail-closed）: %q", raw),
			}
		}
		if !ctx.canon.InAllowlist(canonical) {
			return &contribution{
				rules:  []RuleID{R2},
				level:  L2,
				reason: fmt.Sprintf("R2: 目标路径在授权目录之外: %s", canonical),
			}
		}
	}
	return nil
}

// ruleSensitivePath implements R3 on canonical paths: tier A -> Deny
// (absolute, not overridable by any authorization), tier B -> L2 (single-file
// exemption flow lands with the approval queue). Dormant until ticket 18
// wires resolver + classifier.
func ruleSensitivePath(ctx *assessCtx) *contribution {
	if len(ctx.facts.Paths) == 0 || ctx.canon == nil || ctx.classifier == nil {
		return nil
	}
	var worst *contribution
	for _, raw := range ctx.facts.Paths {
		canonical, err := ctx.canon.Canonicalize(raw)
		if err != nil {
			// Cannot classify a path we cannot canonicalize: fail closed to L2
			// (R2 already flags the same path; R3 records its own conservative verdict).
			worst = &contribution{
				rules:  []RuleID{R3},
				level:  L2,
				reason: fmt.Sprintf("R3: 路径无法规范化，敏感检查 fail-closed: %q", raw),
			}
			continue
		}
		switch ctx.classifier.Classify(canonical) {
		case TierA:
			return &contribution{
				rules:  []RuleID{R3},
				level:  Deny,
				reason: fmt.Sprintf("R3: 命中 A 档敏感路径，拒绝（不可豁免）: %s", canonical),
			}
		case TierB:
			if worst == nil || worst.level < L2 {
				worst = &contribution{
					rules:  []RuleID{R3},
					level:  L2,
					reason: fmt.Sprintf("R3: 命中 B 档敏感路径（可单文件豁免）: %s", canonical),
				}
			}
		}
	}
	return worst
}

// ruleTaint implements R4: an exfil channel parameter carrying a
// sensitive-source fragment -> L2, and the verdict is never overridable by a
// D45 session authorization (SPEC-06 §8.3). Dormant until ticket 19 wires the
// TaintDetector.
func ruleTaint(ctx *assessCtx) *contribution {
	if ctx.taint == nil {
		return nil
	}
	source, hit := ctx.taint.TaintHit(ctx.params)
	if !hit {
		return nil
	}
	return &contribution{
		rules:            []RuleID{R4},
		level:            L2,
		reason:           fmt.Sprintf("R4: 包含来自 %s 的内容", source),
		blockSessionAuth: true,
	}
}
