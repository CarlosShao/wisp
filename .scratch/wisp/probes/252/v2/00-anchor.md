# 252-v2（验收腿，非实现者）锚、方法与受扰记录

- 腿＝`252-v2`（非实现者；被验收＝`252-r3`：产码 `d3cccfda`／交付 `c8a360e8`／票尾 `ade60808`）。⛔ 未翻框、⛔ 未 push、⛔ 零 Go 改动。
- 起手钟点（`date -Iseconds`）＝`2026-10-14T12:24:54+08:00`；起手 HEAD＝`1d19b60b`（`编排者：A919…派 252-v2…`）；工作目录＝`D:\work\workspace\projects plans\Wisp`（win32 / Git Bash）；`go version go1.27.1 windows/amd64`。

## 对象锚（本腿自取）

- 改前锚＝`2b1a3071^`＝`62ca6c34`；本腿自己 `git archive 2b1a3071^ | tar -x` 到仓外 `C:\tmp\wisp252v2\tree-before`（不是 r3 的树）。
- 交付锚＝`d3cccfda`；本腿自己 `git archive d3cccfda` 到仓外 `C:\tmp\wisp252v2\tree-delivery`（封存树：不受在飞腿笔墨影响，AC#4 的绿以它为准）。
- 漂移核：`git diff d3cccfda 1d19b60b --stat -- cmd internal tools scripts`＝**0 行**（12:39 现读）。
- 交付件 sha256：`cmd/wisp/ticket224_assembly_test.go` 工作树与 `d3cccfda` 两处同＝`366da8f5f64488be4975a2b67f7a7e2f2e40bcbb86166c8859958936487180ff`；改前树同件＝`8ac227a80c03ab479a2910a74720f541e278eff7c14d529c2e217cf485f7df6f`。

## 短名皮肤（本腿自己复刻）

- 本腿自建 `C:\tmp\wisp252v2\short-root-252v2`；`cmd //c "dir /x C:\tmp\wisp252v2"` 现读 8.3 别名＝`SHORT-~1`（首读 12:25）。
- 每发一律：`unset GOTMPDIR`（r3 具名的透传坑）；`TMP=TEMP=C:\tmp\wisp252v2\SHORT-~1`；`PATH` 前置 `third_party/sherpa-onnx`（MSYS 形）；`go test -count=1 -v -run '^<name>$' ./cmd/wisp/`。
- 短名落地现证：各发日志 `grep -c 'SHORT-~1'`＝4~9 枚（含 `t.TempDir()` 返回串，如 W1 红句里的 `pattern="C:\\tmp\\wisp252v2\\SHORT-~1\\Test…"`）。

## 受扰窗口（他腿在飞，逐报）

- 12:24:54／12:27:48 两次现读 `git status --porcelain cmd/wisp internal/tools internal/risk`＝0 行。
- `cmd/wisp/run.go` mtime `12:31:35`、`cmd/wisp/config_readers_255.go` mtime `12:33:27`（279-r1 的笔墨）；本腿 12:38:46 现读 porcelain＝2 行（`M cmd/wisp/config_readers_255.go`／`M cmd/wisp/run.go`）。
- ⇒ 12:35 那发"工作树 W2 绿"（PASS 5.70s）带脏面，**弃用、不入判语**；AC#4 的绿一律改在封存树 `tree-delivery` 复跑（12:39:48 起，见 01 件）。
- 12:44:34 现读 porcelain 回 0 行、HEAD＝`a1cf63b2`（279-r1 已落件）。全程三枚 witness 件与交付 blob 逐枚同哈希（`run_test.go` `485d2203…`／`run_mode101_test.go` `c074a666…`／`ticket224_assembly_test.go` `366da8f5…`）。
- CI run 38508311376 在远端自跑，本机不可观测其扰动；本机各发 rc 见 01/02 件。

## 留盘（只建不删）

- 原始日志 9 枚在 `C:\tmp\wisp252v2\logs\`：`w1-before-short-v2.log` 4,615B／`w2-before-short-v2.log` 2,844B／`w3-before-short-v2.log` 3,448B／`w1-after-sealed-v2.log` 3,331B／`w2-after-sealed-v2.log` 3,219B／`w3-after-sealed-v2.log` 3,615B／`ac5-overlay-red-v2.log` 5,411B／`w1-after-short-v2.log` 3,332B（弃用）／`w2-after-short-v2.log` 3,227B（弃用）；门禁三枚 `gate-build-v2.log` 0B／`gate-d22scan-v2.log` 23,758B／`gate-length-v2.log` 1,457B。overlay JSON＝`C:\tmp\wisp252v2\overlay-ac5-v2.json`（内容逐字在 logs/20 件）；两棵封存树 `tree-before\`／`tree-delivery\` 亦留盘。
