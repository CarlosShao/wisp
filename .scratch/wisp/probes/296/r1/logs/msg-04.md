296-r1 AC#4 门禁（build／d22scan／格式两把并排／整包改前改后各两发取交集＝新增红 0 枚）＋越界自查

GOFLAGS= go build ./... rc=0；sh scripts/d22scan.sh rc=0（clean；本腿新增那枚文件在 ban #8 的 cmd/ 118 枚射程内）。

格式按 10-09 新口径并排交、射程写清：
- 工作树尺（gofmt -l cmd/wisp 与 $(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp，两把逐名相同）＝3 枚：
  models.go（CRLF 假枚：工作树 CR 字节 334、HEAD blob 0）＋票 298 名下那两枚已知脏件（⛔ 没顺手洗）。
- 判据尺（同一把 git archive→gofmt 到仓外，逐名作差）：改前 28a2ff4e~1＝2 枚、改后 HEAD＝2 枚、逐名相同 ⇒ 新增 0 枚；
  本腿碰过的两枚文件单独跑两把尺都是空（rc=0）。
- 具名报：裸 gofumpt 不在 PATH＝rc=127（不是"跳过"）。

整包 PATH=sherpa-onnx:build go test ./cmd/wisp/ -count=1（每发起手 tasklist wisp.exe／balldebug.exe 均 0 枚＋date 记锚）：
改前两发 rc=1、^--- FAIL 顶层＝6（含本票那两枚夹具）；改后两发 rc=1、^--- FAIL 顶层＝4，两发逐名相同。
名册三数带尺名：^--- PASS 顶层 273→275、--- PASS 含子 383→385、^--- SKIP 顶层 2→2；改前那发无 -v 的一把天生零 PASS 行，只作红名册交集用。
逐名作差＝**新增红 0 枚**，本票两枚夹具由红转绿。留下的 4 枚红是 WebView2／页面投递那一族
（panel_resident_windows_test.go 的 AC#13/AC#14 三枚 "no report from the page within 15s" ＋ panel_host_windows_test.go:662 冷拉起一枚），
改前改后都在、与本腿那枚热键闭包无关，本腿未动也未放宽任何断言。

越界自查：三笔 commit 名册只含 cmd/wisp/**（resident_windows.go ＋ resident_hotkey_296_windows_test.go）与 .scratch/wisp/probes/296/**；
internal/config/schema.go、internal/ball/**、三枚冻结件、golden、thresholds.go、allowlist.txt、frontend/**、design/**、.gitignore、
D43 表与机主的 config.toml 全部零字节。四发 -v 原始日志（约 1.1 MB）按大日志不入库也不删留在 logs/，件里逐名带字节数。
AC#3 真机那一发归编排者跑。
