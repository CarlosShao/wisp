# 418-AC#4 — `GOOS=linux GOARCH=amd64 go build ./...` 逐包扫描读数（2026-11-10 19:3x +08，腿 `418-ac4`）

## 0. 环境与命令（可逐字复跑）

- 仓锚：HEAD `477a2ed00`；418 落地笔 `f54c972c2` **已核在祖先**（`git merge-base --is-ancestor f54c972c2 HEAD` ⇒ 真）。
- 前缀（派单铁律）：`cd "/d/work/workspace/projects plans/Wisp"; export LC_ALL=C; unset GOTMPDIR TMPDIR TEMP TMP; export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`。
- ⛔ 未设 `CGO_ENABLED=0`（派单明禁）；Windows 宿主交叉 `GOOS=linux` 时 Go 工具链**自动关 cgo** ⇒ 下面读到的形状是「cgo 依赖整包被 build constraints 排除」这一既有交叉档机制（台账 `A1032` 第③节定式在册）。

## 1. 第一把尺：整仓 `GOOS=linux GOARCH=amd64 go build ./... 2>&1`（完整输出，逐字）

```
package github.com/CarlosShao/wisp/cmd/wisp
	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8
RC=1
```

⚠ 整仓发只暴露**一条失败链**（加载序在 `cmd/wisp` 处停）⇒ 不能据此说"只有一枚包编不过"，故有第 2 把尺。

## 2. 第二把尺：逐包单扫（分母＝`GOOS=linux GOARCH=amd64 go list ./...`，**40 枚包**）

分母清单（40 枚，存档 `/d/tmp/wisp/20261110/418ac4-pkglist.txt`）：

```
github.com/CarlosShao/wisp/cmd/llmrecord
github.com/CarlosShao/wisp/cmd/wisp
github.com/CarlosShao/wisp/frontend
github.com/CarlosShao/wisp/internal/agent
github.com/CarlosShao/wisp/internal/agent/approval
github.com/CarlosShao/wisp/internal/agent/scheduler
github.com/CarlosShao/wisp/internal/audio
github.com/CarlosShao/wisp/internal/audio/aec
github.com/CarlosShao/wisp/internal/audio/aecspike
github.com/CarlosShao/wisp/internal/ball
github.com/CarlosShao/wisp/internal/buildinfo
github.com/CarlosShao/wisp/internal/command
github.com/CarlosShao/wisp/internal/config
github.com/CarlosShao/wisp/internal/llm
github.com/CarlosShao/wisp/internal/llm/adaptertest
github.com/CarlosShao/wisp/internal/llm/anthropic
github.com/CarlosShao/wisp/internal/llm/golden
github.com/CarlosShao/wisp/internal/llm/openaichat
github.com/CarlosShao/wisp/internal/llm/openairesponses
github.com/CarlosShao/wisp/internal/memory
github.com/CarlosShao/wisp/internal/models
github.com/CarlosShao/wisp/internal/observe
github.com/CarlosShao/wisp/internal/panel
github.com/CarlosShao/wisp/internal/perm
github.com/CarlosShao/wisp/internal/plugin
github.com/CarlosShao/wisp/internal/plugin/manifestslot
github.com/CarlosShao/wisp/internal/present
github.com/CarlosShao/wisp/internal/proc
github.com/CarlosShao/wisp/internal/projctx
github.com/CarlosShao/wisp/internal/risk
github.com/CarlosShao/wisp/internal/secret
github.com/CarlosShao/wisp/internal/session
github.com/CarlosShao/wisp/internal/speech
github.com/CarlosShao/wisp/internal/speech/cer
github.com/CarlosShao/wisp/internal/statemachine
github.com/CarlosShao/wisp/internal/streamkey
github.com/CarlosShao/wisp/internal/tools
github.com/CarlosShao/wisp/internal/watchdog
github.com/CarlosShao/wisp/internal/winsec
github.com/CarlosShao/wisp/tools/signmodels
```

循环＝对每枚单独 `GOOS=linux GOARCH=amd64 go build "$p"` 记 rc 与输出。**TOTAL=40 FAIL=2**。失败包完整输出（逐字，全部输出仅此两块，其余 38 枚 rc=0 零输出）：

```
=== FAIL(1) github.com/CarlosShao/wisp/cmd/wisp ===
package github.com/CarlosShao/wisp/cmd/wisp
	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8
=== FAIL(1) github.com/CarlosShao/wisp/internal/speech ===
package github.com/CarlosShao/wisp/internal/speech
	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8
```

## 3. 第三把尺（佐证）：`internal/ball` 的 Linux 编译面已愈

```
GOOS=linux GOARCH=amd64 go vet ./internal/ball/                 ⇒ RC=0（零输出）
GOOS=linux GOARCH=amd64 go vet ./internal/proc/ ./internal/ball/ ⇒ RC=0（零输出）＝ TestCrossVetForLinux 同射程
```

⇒ 票 418 的根因（`orbdebug.go` 无 tag 引 `Ball`）在 build/vet 层面**确认消失**；`internal/ball` **不在**今天的 Linux 编不过名单里。

## 4. 逐枚归属判定

| # | 失败包 | 根因 | 归属判定 |
|---|---|---|---|
| 1 | `cmd/wisp` | `sherpa-onnx-go-linux@v1.13.8` cgo 依赖在交叉档被 build constraints 整包排除；错误文本**零行指向本仓源码** | **存量交叉档形状，非缺陷、非新红、归已有账**：台账 `A1032` 第③节定式（Windows 宿主交叉＝cgo 关 ⇒ 必然红，⛔ 与本机档成对判）；同形前科 `docs/evidence/s1/102:147`、`105:166`、`121:276`、`127:162`、`128-ac4-gates-and-ci-divergence.md:235`（修前锚同逐字红句＝差分证明既有）、`131-followup-1:280`、`135-ac1-ac6-impl.md:100`、`136-ac8-ac9-impl.md:322`、`137-ac2-ac5-impl.md:259`、`138-cancel-kills-children-r1.md:805`、`342-ac4-readings.md:64`（"⛔ 当新红"）；`pending-and-issues.md:4436` 也记 cmd/wisp＋cmd/balldebug 同形。CI 的 ubuntu 本机档带 cgo 不咬 ⇒ 是**测量口径**不是产品缺陷（`A1032` 原裁"不立票"）。 |
| 2 | `internal/speech` | 同上（`internal/speech/asr.go` 等 7 枚文件 import `sherpa_onnx`，加载同一外部 cgo 模块） | **同根因、同归属**：`15b-r1-impl.md:100` 逐字在册「`GOOS=linux CGO_ENABLED=0 go vet`：`internal/speech` rc 1——与未改动的 `cmd/wisp` 逐字同因 ⇒ 没引入新的破法」。零行指向本仓源码。 |

**不在名单里的（对照票 230 AC#5 在册与票 418 原始现量）**：

- `internal/ball`：418 起手现量那枚（`orbdebug.go:39/43 undefined: Ball`，CI run `129250565302` 真红）⇒ **本程已消失**（第 3 节两把 rc=0）＝**本票治的那枚，已治**。
- `internal/panel`：GOOS=linux **build** 干净（rc=0，在 40 枚分母内）。票 230 AC#5 在册的 panel 4 枚（`230-cells` 读数：`TestC21DesignTokensFourWayAgree` 等）是 **Windows 测试面红，不是 Linux 编译面红** —— 两本账不同尺，本程只对 Linux build 面作证，不动 AC#5 那册。
- 票 230 AC#5 在册的另一枚 `internal/ball 1`（`TestCrossVetForLinux` 真红，`230-cells` §②）⇒ 根因即 418 这枚，修后**应随 418 转绿**（本程 vet 同射程 rc=0 佐证；`go test` 面绿不归本腿跑，AC#5 复量为准）。

**新增账判定**：本程逐包扫**没有发现任何一枚新的、无主 Linux 编不过包** —— 2 枚失败全归存量交叉档口径，418 自己那枚已消失。按 AC#4 字面「需把新读数追加进 230 的账」⇒ 本读数表已追加进票 230 Progress log（只追加、⛔ 改 AC 原文）。

## 5. 本腿边界（诚实说明）

- 写面＝票 230 ＋ 票 418 ＋ 本证据件，**零产码**；`git status --porcelain -- '*.go'` 中本腿触碰＝0。
- 写面全 `.md` ⇒ `gofumpt` 射程＝空集（非"跳过"）。
- ⛔ 翻框（AC 判据勾归编排者收件）；⛔ push（只 commit）。
