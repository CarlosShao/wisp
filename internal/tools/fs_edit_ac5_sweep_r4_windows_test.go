//go:build windows

package tools

// Ticket 162 AC#6 的收尾件（派单 §3 那一枚〔现读码〕）。
//
// 162-r3 交回 AC#4b 时记下一笔没钉的账：`fs.edit` 被真硬杀留下的 `.wisp-tmp-*`
// 残件，会不会被下一次写盘的清扫器带走。当时的凭据只有〔现读码
// `fs_write.go:288`、`fs.edit` 经 `fs_edit.go:271` 走同一枚〕、零枚用例钉，
// 编排者 16:5x 把它记成"别让它蒸发成大家都读过码了"。
//
// 这一格钉的是**从 fs.edit 这一侧**的那条链：一枚可归因、创建者进程已死的孤儿
// 躺在目标目录里，`fs.edit` 的一次成功写盘必须把它带走，并且台账要把这件事
// 说出来——"清扫挂在每次写盘上"是一句能被读到的记录，不是注释。
//
// 另一侧（清扫只扫自己的、活进程的文件不动）由票 73 那族的
// `TestSweepReclaimsOnlyItsOwnStagingFiles` 钉着，本文件不另造第二把尺：
// 它复用同一批 helper（`deadPID` / `plantOrphan` / `stagingFiles`），
// 只把触发清扫的那一次写盘换成 fs.edit。
//
// ⚠ 持续性断言的范围：这里断言的是**用例窗内**那次写盘之后孤儿消失。它本来也
// 不会活到窗尾——Go 的 `t.TempDir()` 在清理阶段会带走整个目录（162-v2 就是
// 用这一点把 r3 的"残件留在盘上"收窄成窗内为真的）。

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFSEditWriteReclaimsADeadWritersOrphan(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "alpha\nbravo\n"
	target := filepath.Join(root, "has.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}

	// A staging file the sweeper CAN attribute and CANNOT find a live owner for:
	// exactly the shape a taskkilled write leaves behind (AC#4b named one).
	dead := deadPID(t)
	orphan := plantOrphan(t, root, dead, "4242")

	b, _ := fsEditBridge(t, mustCanonical(t, root))
	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "bravo", "new": "BRAVO"})))
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out.IsError {
		t.Fatalf("这一次写盘本身必须成功，否则清扫断言没有意义: %+v", out)
	}
	if got := readString(t, target); got != "alpha\nBRAVO\n" {
		t.Fatalf("编辑没落地: %q", got)
	}

	// (a) the orphan is gone — the reclaim really ran inside this write.
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatalf("fs.edit 的写盘没扫掉死进程留下的孤儿 %s（err=%v）⇒ 清扫器根本没跑在这一侧", orphan, err)
	}
	// (b) nothing of THIS write is left behind either: the rename consumed it.
	if left := stagingFiles(t, root); len(left) != 0 {
		t.Errorf("写盘后目录里仍留 %d 枚 .wisp-tmp-*：%v", len(left), left)
	}
	// (c) and the ledger says so, because "residue was reclaimed" must be
	// readable by someone who is not standing in this directory.
	var swept string
	for _, s := range out.AppliedSteps {
		if strings.Contains(s, "清扫") {
			swept = s
		}
	}
	if swept == "" {
		t.Fatalf("台账里没有清扫那一步，ledger=%q", out.AppliedSteps)
	}
	t.Logf("读数：fs.edit 一次成功写盘带走死进程孤儿 %s；台账=%q", filepath.Base(orphan), swept)
}
