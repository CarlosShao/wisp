package main

// 票 223 非实现者验收腿 223-v1 的 AC#4 台件（**只经 go test -overlay 挂载，物理文件在
// .scratch/wisp/probes/223/v1/overlay/，仓内 cmd/wisp 目录里不存在这一枚文件**）。
//
// 目的：把同一批对抗输入喂给**生产分类器** describeReloadFailure（cmd/wisp/config_reload.go:315），
// 读的是 `wisp run` 那句话本身，不是 loader 的内部返回值。种法走真实路径：
// config.NewManager 装配一枚活 Manager → 手改文件 → CheckAndReload → 把错误交生产句子函数。

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
)

type v1SentenceCase struct {
	name string
	body string
}

func v1SentenceCases() []v1SentenceCase {
	return []v1SentenceCase{
		{"A1_票面基线_无版本且语法坏", "this is not toml [[[\n"},
		{"A2_票面基线_声明1且坏表头", "schema_version = 1\n[llm\nbroken ===\n"},
		{"B1_版本写在section里", "[app]\nschema_version = 2\n"},
		{"B2_版本写在section里且语法坏", "[app]\nschema_version = 2\nbroken [[[\n"},
		{"C1_注释行以方括号开头", "# [fs] 这不是表头\nschema_version = 2\nbroken [[[\n"},
		{"D1_多行字符串里的假表头_无版本", "note = \"\"\"\n[fake]\n\"\"\"\nbroken [[[\n"},
		{"D2_多行字符串里的假表头_版本在其后", "note = \"\"\"\n[fake]\n\"\"\"\nschema_version = 1\nbroken [[[\n"},
		{"E1_CRLF且声明版本", "schema_version = 2\r\nbroken [[[\r\n"},
		{"E2_CRLF无版本", "broken [[[\r\n"},
		{"F1_BOM且声明版本", "\ufeffschema_version = 2\nbroken [[[\n"},
		{"G1_无空格写法", "schema_version=2\nbroken [[[\n"},
		{"H1_同名键出现两次", "schema_version = 2\nschema_version = 1\n"},
		{"I1_首行就是fs表", "[fs]\nallowed_dirs = [\"D:/x\"]\nschema_version = 2\nbroken [[[\n"},
		{"J1_声明未来版本99且语法坏", "schema_version = 99\nbroken [[[\n"},
		{"K1_版本被写成字符串", "schema_version = \"2\"\nbroken [[[\n"},
		{"L1_行首缩进的版本行", "   schema_version = 2\nbroken [[[\n"},
		{"M1_合法v1需要真迁移", "schema_version = 1\n[app]\ntheme = \"dark\"\n"},
	}
}

func TestV1ProbeProductionSentence(t *testing.T) {
	for _, tc := range v1SentenceCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "config.toml")
			if err := config.SaveFile(path, config.NewDefaults()); err != nil {
				t.Fatalf("seed canonical config: %v", err)
			}
			mgr, err := config.NewManager(path, nil)
			if err != nil {
				t.Fatalf("NewManager on a canonical config must work: %v", err)
			}
			time.Sleep(30 * time.Millisecond)
			if err := os.WriteFile(path, []byte(tc.body), 0o600); err != nil {
				t.Fatalf("plant: %v", err)
			}
			time.Sleep(30 * time.Millisecond)
			rep, err := mgr.CheckAndReload()
			switch {
			case err != nil:
				t.Logf("SENTENCE %q", tc.name)
				t.Logf("  raw err = %v", err)
				t.Logf("  PRODUCTION LINE = config: HOT-RELOAD state=not-applied %s", describeReloadFailure(err))
			case rep == nil:
				t.Logf("SENTENCE %q -> (nil,nil)：轮询认为什么都没变（静默）", tc.name)
			default:
				t.Logf("SENTENCE %q -> APPLIED hot=%v reload=%v restart=%v locked=%d",
					tc.name, rep.Hot, rep.Reload, rep.Restart, len(rep.Locked))
			}
		})
	}
}
