package main

// 223-v2 任务③ J1 三形台件（overlay 件，物理件在 .scratch/wisp/probes/223/v2/overlay/）。
//
// 这一枚只读数、不判等：它把三形文件种进真 Manager、走真 CheckAndReload、
// 把生产分类器 describeReloadFailure 的那一句原文打到日志里，然后一律 PASS。
// 它不放宽任何断言，因为它压根没有断言 —— 判决靠读原文，不靠颜色。
//
// 三形（对应派单 ③ 的 a/b/c）：
//   (a) 声明未来版 + 正文语法坏          -> 期望 cause=syntax
//   (b) 声明未来版 + 正文解析得开        -> 读它到底报哪一句
//   (c) 声明未来版 + 未知键              -> 读它归哪一句
//
// 另附两形对照，只为把 (b)(c) 的归因钉死：
//   (d) 当前版   + 未知键                -> 现读它该是 cause=unknown-key
//   (e) 当前版   + 解析得开             -> 现读它该是干净重载（无错）
//
// PATH warning (ticket 98)：cmd/wisp 测试二进制链 sherpa-onnx，PATH 不带它会 0xc0000135。

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
)

func v2ProbeSentence(t *testing.T, label, body string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := config.SaveFile(path, config.NewDefaults()); err != nil {
		t.Fatalf("seed canonical config: %v", err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager on canonical config: %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("plant: %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	rep, err := mgr.CheckAndReload()
	if err != nil {
		t.Logf("SHAPE %s\n  declared body = %q\n  reload err    = %v\n  PRODUCTION LINE = %s",
			label, body, err, describeReloadFailure(err))
		return
	}
	t.Logf("SHAPE %s\n  declared body = %q\n  reload err    = <nil>\n  report        = hot=%v reload=%v restart=%v locked=%d",
		label, body, rep.Hot, rep.Reload, rep.Restart, len(rep.Locked))
}

func TestV2ProbeJ1Shapes(t *testing.T) {
	v2ProbeSentence(t, "a_未来版_正文语法坏", "schema_version = 99\nbroken [[[\n")
	v2ProbeSentence(t, "b_未来版_正文解析得开", "schema_version = 99\n\n[app]\nbogus_future_only_key = 1\n")
	v2ProbeSentence(t, "c_未来版_顶层未知键", "schema_version = 99\nthis_key_does_not_exist = 1\n")
	v2ProbeSentence(t, "d_当前版_顶层未知键_对照", "schema_version = 2\nthis_key_does_not_exist = 1\n")
	v2ProbeSentence(t, "e_当前版_解析得开_对照", "schema_version = 2\n")
	v2ProbeSentence(t, "f_未来版_合法键全对_对照", "schema_version = 99\n\n[ball]\nsize = 64\n")
}
