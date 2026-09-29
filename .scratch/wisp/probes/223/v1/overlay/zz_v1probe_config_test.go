package config

// 票 223 非实现者验收腿 223-v1 的 AC#4 台件（**只经 go test -overlay 挂载，物理文件在
// .scratch/wisp/probes/223/v1/overlay/，仓内 internal/config 目录里不存在这一枚文件**）。
//
// 目的：把 loader.go 新造的 declaredSchemaVersion 按行扫描器逼到边界上，实测
// "语法错" 与 "迁移失败" 这两句话会不会说反，以及 peekSchemaVersion 的解析错误
// 是不是真的不再被丢掉。本文件**不判对错、只读数**——判语写在裁决表里。

import (
	"os"
	"path/filepath"
	"testing"
)

type v1ProbeCase struct {
	name string
	body string
	why  string
}

func v1ProbeCases() []v1ProbeCase {
	return []v1ProbeCase{
		{"A1_票面基线_无版本且语法坏", "this is not toml [[[\n", "实现腿 AC#4 种的那一发"},
		{"A2_票面基线_声明1且坏表头", "schema_version = 1\n[llm\nbroken ===\n", "migrate_test.go:123 钉住的形状"},
		{"B1_版本写在section里", "[app]\nschema_version = 2\n", "任务①：非文件头"},
		{"B2_版本写在section里且语法坏", "[app]\nschema_version = 2\nbroken [[[\n", "任务①的对抗发"},
		{"C1_注释行以方括号开头", "# [fs] 这不是表头\nschema_version = 2\nbroken [[[\n", "任务②"},
		{"D1_多行字符串里有一行以方括号开头_无版本", "note = \"\"\"\n[fake]\n\"\"\"\nbroken [[[\n", "任务③"},
		{"D2_多行字符串里有一行以方括号开头_版本在其后", "note = \"\"\"\n[fake]\n\"\"\"\nschema_version = 1\nbroken [[[\n", "任务③的对抗发：扫描器会不会先撞假表头就停"},
		{"E1_CRLF且声明版本", "schema_version = 2\r\nbroken [[[\r\n", "任务④"},
		{"E2_CRLF无版本", "broken [[[\r\n", "任务④的对照"},
		{"F1_BOM且声明版本", "\ufeffschema_version = 2\nbroken [[[\n", "任务⑤"},
		{"G1_无空格写法", "schema_version=2\nbroken [[[\n", "任务⑥"},
		{"H1_同名键出现两次", "schema_version = 2\nschema_version = 1\n", "任务⑦：TOML 里重复键＝解析错"},
		{"I1_首行就是fs表", "[fs]\nallowed_dirs = [\"D:/x\"]\nschema_version = 2\nbroken [[[\n", "任务⑧"},
		{"J1_声明未来版本99且语法坏", "schema_version = 99\nbroken [[[\n", "版本探针读到 99 会走哪一句"},
		{"K1_版本被写成字符串", "schema_version = \"2\"\nbroken [[[\n", "Atoi 解不开时归谁"},
		{"L1_行首缩进的版本行", "   schema_version = 2\nbroken [[[\n", "TrimSpace 支不支持"},
		{"M1_合法v1需要真迁移", "schema_version = 1\n[app]\ntheme = \"dark\"\n", "对照组：这句话本就该归迁移管线"},
	}
}

func TestV1ProbeLoaderClassification(t *testing.T) {
	for _, tc := range v1ProbeCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.toml")
			if err := os.WriteFile(path, []byte(tc.body), 0o600); err != nil {
				t.Fatalf("write: %v", err)
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read back: %v", err)
			}
			dv, dok := declaredSchemaVersion(raw)
			pv, perr := peekSchemaVersion(raw)
			cfg, lerr := readConfigFile(path)
			after, _ := os.ReadFile(path)

			t.Logf("PROBE %q  (%s)", tc.name, tc.why)
			t.Logf("  declaredSchemaVersion -> value=%d ok=%v", dv, dok)
			t.Logf("  peekSchemaVersion     -> ver=%d peekErrNil=%v", pv, perr == nil)
			if lerr != nil {
				t.Logf("  readConfigFile        -> ERR class-text=%q", lerr.Error())
			} else {
				t.Logf("  readConfigFile        -> OK schema_version=%d ball.size=%d", cfg.SchemaVersion, cfg.Ball.Size)
			}
			if string(after) != string(raw) {
				t.Logf("  *** FILE WAS REWRITTEN by the load pipeline (%d -> %d bytes)", len(raw), len(after))
			} else {
				t.Logf("  file untouched")
			}
		})
	}
}
