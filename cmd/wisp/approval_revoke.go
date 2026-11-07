package main

// Ticket 226 AC#5 - the reply grammar's revoke verb: the answer side of the
// 撤销 half of a long-lived rule.
//
// WHY THIS FILE EXISTS. The 长期 branch (approval_always.go) stores one line in
// [fs] allowed_dirs behind a second, L2-level confirmation - and from the day it
// landed until this file, nothing in the reply grammar could remove that line
// again. The write path for removal has existed since ticket 201's census named
// the gap and ticket 226's landing gave it a guarded body:
// Manager.SetAllowedDirs (internal/config/allowdirs.go:109, whose own header
// calls it "the 撤销 half of AddAllowedDir") rewrites exactly that one key and
// leaves every other line of the file alone. What was missing was the entrance:
// SetAllowedDirs had ZERO production callers and the audit trail booked
// WIDEN-APPLIED dir=... lines nobody could act on. Ticket 226 AC#5's sentence is
// 「既能落库、又能被答复语法触发」; the 落库 half shipped with the write guard,
// this file is the 触发 half.
//
// WHY NO SECOND CARD. PLAN.md:1645-1646 freezes the re-confirmation requirement
// for a LOOSENING of [fs]; a removal is a tightening, and applyLocked's own
// bookkeeping (the sentence SetAllowedDirs' header quotes) hot-applies a
// tightening without asking. So revoke raises no card, exactly like the mode
// switcher never asks before narrowing. The confirmation asymmetry is the
// contract, not an oversight, and this file does not soften it in either
// direction: revoke can only shrink the list the process holds.
//
// WHY IT ADDRESSES THE RULE, NOT A CORRELATION ID. Every other verb here is
// typed against a card number, but the ledger forgets a card the moment its
// answer is spent (Replies.Allow/Reject/Veto all call Forget) - by the time an
// operator wants a rule gone, the card that offered it is long gone from the
// ledger. What survives is the rule's own text: the 长期 flow prints it on the
// screen (「要存的规则：... [fs] allowed_dirs += "D:/..."」) and books it in the
// audit (WIDEN-APPLIED dir=...). revoke takes that same text, or the bare
// directory, and names the line to remove. Fail-closed both ways: a rule the
// process does not hold is refused out loud with the current list spelled out,
// and a write that claims success while the rule is still in memory is booked
// REVOKE-FAILED instead of being answered with a happy sentence.
//
// LIKE 「一直」, THE EFFECT IS NEXT-START. The running process keeps the C26
// canonicalizer it built at assembly (run.go:624) and the 长期 card already says
// the stored rule applies on the next start; a revocation is that promise read
// backwards and gets the same honest sentence, not an over-claim of liveness.

import (
	"errors"
	"fmt"
	"strings"
)

// revokeSyntax is the spelling this file's own help shows, and the exact shape
// the 长期 flow prints and the audit books - so an operator can copy either one
// after the verb and name the same line.
const (
	revokeVerb  = "revoke"
	revokeRoute = "native/revoke"
	revokeTool  = "config.allow_dir"
)

const revokeSyntax = `[fs] allowed_dirs += "D:/某目录"`

// revokeRule removes one long-lived rule from [fs] allowed_dirs through
// Manager.SetAllowedDirs, the guarded single-key write path. It takes the rule
// text exactly as the 长期 flow printed it, or the bare directory; nothing else
// is accepted, because a removal that guessed at which line the operator meant
// would be a second, looser grammar growing next to the one that widens.
func (s *replySurface) revokeRule(arg string) (string, error) {
	dir, ok := revokeDirOf(arg)
	if !ok {
		s.record("REFUSED", "", revokeTool, revokeRoute, "没有指名要撤销的规则")
		return "", errors.New("要撤销哪一条长期规则？规则得指名目录：" + revokeVerb +
			" " + revokeSyntax + "（光敲 revoke 不撤销任何东西）")
	}
	if s.rt == nil || s.rt.mgr == nil {
		s.record("REFUSED", "", revokeTool, revokeRoute, "配置写面未装配")
		return "", errors.New("配置写面没有装配：撤销无处可写，配置一个字节都没动")
	}
	before := s.rt.mgr.Config()
	if !dirInList(before.FS.AllowedDirs, dir) {
		s.record("REFUSED", "", revokeTool, revokeRoute,
			fmt.Sprintf("dir=%q not in the in-memory allowlist %v", dir, before.FS.AllowedDirs))
		return "", fmt.Errorf("本进程账上没有 %q 这条规则（当前 [fs] allowed_dirs：%v）；"+
			"文件一个字节都没动", dir, before.FS.AllowedDirs)
	}
	if err := s.rt.mgr.SetAllowedDirs(listWithoutDir(before.FS.AllowedDirs, dir)); err != nil {
		s.record("REVOKE-WRITE-FAILED", "", revokeTool, revokeRoute,
			fmt.Sprintf("dir=%q err=%v", dir, err))
		return "", fmt.Errorf("撤销写盘失败，配置一个字节都没改：%v", err)
	}
	// The write path is guarded, but the reading this verb owes the operator is
	// "the rule is really gone", so it is asserted on the manager's own view
	// instead of inferred from a nil error.
	after := s.rt.mgr.Config()
	if dirInList(after.FS.AllowedDirs, dir) {
		s.record("REVOKE-FAILED", "", revokeTool, revokeRoute,
			fmt.Sprintf("dir=%q still in memory after the write: %v", dir, after.FS.AllowedDirs))
		return "", fmt.Errorf("撤销没有被遵守：%q 仍在内存名单里（这是程序故障，不是你的答复有错）", dir)
	}
	s.record("REVOKED", "", revokeTool, revokeRoute,
		fmt.Sprintf("dir=%q entries=%d->%d key=fs.allowed_dirs effect=next-start",
			dir, len(before.FS.AllowedDirs), len(after.FS.AllowedDirs)))
	return fmt.Sprintf("已撤销长期规则 %q：[fs] allowed_dirs 从 %d 条减为 %d 条，已写回 config.toml；"+
		"本次运行仍按旧名单判定（撤销与「一直」同轨，生效都在下一次启动）",
		dir, len(before.FS.AllowedDirs), len(after.FS.AllowedDirs)), nil
}

// revokeDirOf reads the rule argument. Two shapes are accepted, and they are the
// two the operator can actually see: the exact line the 长期 flow printed
// ([fs] allowed_dirs += "D:/x", also the audit's rule text) and the bare
// directory. Anything else - including an empty argument - is a refusal, not a
// guess.
func revokeDirOf(arg string) (string, bool) {
	s := strings.TrimSpace(arg)
	const head = `[fs] allowed_dirs += "`
	if strings.HasPrefix(s, head) && strings.HasSuffix(s, `"`) {
		s = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(s, head), `"`))
	}
	if s == "" {
		return "", false
	}
	return s, true
}

// dirInList and listWithoutDir are this file's one comparison, kept next to each
// other so a list spelling change cannot loosen one and not the other. They
// mirror the manager's own entry comparison (allowdirs.go trims before it
// compares): trimming here is reading, not deciding paths.
func dirInList(list []string, dir string) bool {
	for _, have := range list {
		if strings.TrimSpace(have) == dir {
			return true
		}
	}
	return false
}

func listWithoutDir(list []string, dir string) []string {
	next := make([]string, 0, len(list))
	for _, have := range list {
		if strings.TrimSpace(have) != dir {
			next = append(next, have)
		}
	}
	return next
}
