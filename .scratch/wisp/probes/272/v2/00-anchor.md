# 票 272 · 腿 `272-v2`（非实现者验收腿）— 起手锚（第 1 笔）

派单：`.scratch/wisp/probes/orch/20261013/txt/272v2-dispatch.txt`（2026-10-14 14:1x 发出）。
本腿面＝**只读＋重跑**：⛔ 不翻框、⛔ 不动 `Status:`/票名、⛔ 不改任何产码/测试字节；突变一律 `-overlay`（仓外件）。

## §0 起手读数（本腿自跑，`date` 与各读数同批）

```
$ date
Wed Oct 14 13:46:33 CST 2026

$ git rev-parse HEAD          # 13:46:33 起手；⚠ recon 期间 HEAD 漂过（13:40 eaf846a0 → 4cd874b7 → 11cfbc93=174-r6 交件笔，均未触 cmd/internal）
11cfbc93fd490bbbddc0a79f8f244dddfb132882

$ git status --porcelain -- cmd internal scripts
 M scripts/portable-tests.sh      # 317-r1 在飞（scripts 面，非本腿）
 M scripts/winsec-tests.sh
 M scripts/wisp-cli-tests.sh

$ git status --porcelain -- cmd internal | wc -l
0                                 # cmd／internal 写面干净（本腿要的产码+被审尺两面零脏）

$ md5sum cmd/wisp/config_reload.go cmd/wisp/config_reload_223_test.go
348f5a67a73442d45f751079e3b80a6c *cmd/wisp/config_reload.go
2ea71578ba9ca015c535673bc9792d09 *cmd/wisp/config_reload_223_test.go
# 两枚与 `git show HEAD:` 逐字同值（316-r1 10-14 13:0x 件里的同值 2ea71578…/348f5a67…，非抄）

$ grep -c '^- \[ \]' <票272>  ⇒ 2      # AC#5、AC#7 未勾
$ grep -c '^- \[x\]' <票272>  ⇒ 5      # AC#1–#4、AC#6 已勾

$ go version        ⇒ go1.27.1 windows/amd64
$ gofumpt -version  ⇒ v0.12.0 (go1.27.1)   # D:\work\base\gopath\bin\gofumpt.exe
$ go env GOPATH     ⇒ D:\work\base\gopath
```

## §1 本腿要跑的四门（逐条现跑，rc 逐字进 `logs/`）＋AC#7（两发 overlay）

1. `PATH="$PWD/third_party/sherpa-onnx:$PATH" sh scripts/d22scan.sh`（正控先绿）。
2. `gofmt -l`／`gofumpt v0.12.0 -l` **只喂 `.go`**（被审 `cmd/wisp/config_reload_223_test.go` 与产码；⛔ 不喂 .md）。
3. `go vet ./cmd/wisp/` rc。
4. 名册：`go test -count=1 ./cmd/wisp/`（整包，314 死区如实记）＋`./internal/...`（逐包）。
   红名归因位＝在册族：panel 契约漂移／C21／35r8×2／`TestSessionListByWorkspace` 间歇／环境形
   （`TestAC2SealNoticeLandsInTheRunLegLogFile`=GOTMPDIR 形）／计时带载形（`TestResolvePerCallBudget`、`TestTenOpsInOneToolCallGetOneConfirm`）。
5. AC#7：以**现行 HEAD 产码**重建「单份 post-mark 欢迎语」突变（删 `reportRestartPending` 的真句＋在 `reportReload` 的 `state=applied` 后注入一份六 needle 的欢迎语打印）⇒ 落地尺上 **PASS 期望**；
   配**零份对照**（同位置不注入、只删真句）⇒ 期望 `:623` await 红（非空转证明）。两发全 `-overlay`，仓外 `D:/tmp/wisp272v2/`。

## §2 环境纪律自报

- `GOTMPDIR` 只用**反斜杠形** `C:\tmp\wisp272v2gotmp`（317 号票的环境坑；⛔ posix 形）。
- 只 commit、⛔ 不 push；显式 pathspec。他腿在飞：`317-r1`（scripts，工作树已见 3 枚 M）、`236-r7`（internal/tools）。
- ⚠ 起手发现一枚**仓外**脏目录 `D:\work\workspace\projects plans Wisp\`（无 `/`，生于 12:20:08，非本腿）——
  内含空 `.scratch/wisp/probes/272/v2/logs`（本腿 13:46 的一次路径笔误 `mkdir` 落进去的空目录）＋空 `probes/orch/20261013/txt/`。
  按「临时件只建不删」**不删、如实登记**；本腿真实件在仓内 `.scratch/wisp/probes/272/v2/`。
