# 批 A：[ball] size 配置上限 72 → 160（config 层与契约文字）

- 日期：2026-11-08 15:3x–15:4x +08
- 执行：Wisp 编队写腿（批 A，config＋docs 射程；`internal/ball/**` 零触碰，批 B 另派）
- 台账：`A1636`

## 1. 授权出处（改契约须留出处）

机主 2026-11-08 15:4x 逐字授权：

> 问：「① 球要不要放大？（推荐：放大到 120~160）」
> **机主答：「可放大，按你推荐的试试」**

⇒ 改契约已获机主批准（`SPEC-12 §4.1` 的「改契约＝人工批准」这一步完成）。
台账落账 `A1636`（`docs/reports/pending-and-issues.md:32171` 起）。

## 2. 裁定与边界

- **新范围 = 44..160**；**默认值仍 56，未动**（默认变大是行为变更；先放能力，机主试完满意再议默认）。
- 上限 160 依据：`WindowEdgePx(160,96)=204` ⇒ 绘制区 160px 够画精细粒子；再大在 1080p 占 1/5 宽。
- **`BallSizeLargePx` 档位语义不动**（仍 72）：上限 ≠ 档位。`internal/ball/tokens.go` 的 `BallSizeMaxPx=160` 改动属批 B，本腿未碰。

## 3. 起手三核（逐枚现跑核实，非转录）

| # | 位置 | 现跑原文 | 处置 |
|---|---|---|---|
| ① | `internal/config/validate.go:57,62,64` | 注释 `size 44-72`；`< 44 \|\| > 72`；`out of range [44, 72]` | 全部改 160 |
| ② | `docs/specs/SPEC-03-config-secrets-envs.md:27` | `` `size(int)=56`（44–72） `` | 只改数字 ⇒（44–160），默认 56 一字未动 |
| ③ | `docs/PLAN.md:3557` | 「格子内球体 **56px**（配置范围 44–72px，原型另出三档尺寸对照）」 | 只改数字 ⇒ 44–160px |
| ④ | `internal/config/validate_test.go:156-170` | bad `{43,73,0,-1}`／good `{44,56,72}`／`allowed 44-72` | 迁移＋加严（见 §4） |

**PLAN.md:3557 上下文裁定（如实报）**：该句在 §17.7「悬浮球原型（`ball.html`）」节里，主体讲的是原型画布格子（每态一格 160×160、格子内球体 56px）——**但括号里那句自述是「配置范围 44–72px」，指的是配置契约本身**（与 SPEC-03 同源），不是画布格子尺寸。⇒ **该改**，且只改括号内数字；「三档尺寸对照」（sm/md/lg 档位）不属配置范围，未动。

**本腿射程内追加一处**：`internal/config/schema.go:179` 结构体注释 `range 44-72 (SPEC-03 sec 3)` ⇒ `44-160`（同注释，注明 A1636 授权、默认留 56）。

## 4. 测试改动（迁移＋加严，⛔ 未删任何既有断言）

`TestValidateBallSizeRange` 重切为逐种子 subtest：

- 拒绝：`43 / 161 / 200 / 0 / -1`（原 `73` 从拒绝表迁出——新范围下 73 合法；新增 `161`、`200`）。
- 接受：`44 / 56 / 72 / 73 / 160`（新增 `73`、`160`）。
- 加严：错误必须含键名 `ball.size` 且含字面 band `out of range [44, 160]`（同 `TestValidateRiskConfirmTimeoutRange` 的字面 band 规矩）。
- 旧 `allowed 44-72` 文案 ⇒ `allowed 44-160`。

## 5. 牙（突变）读数

把上限临时改回 72（`MUTATION-160tooth` 标记）跑目标用例：

```
--- FAIL: TestValidateBallSizeRange
    --- FAIL: .../accept_73   ball.size=73 must validate: ... out of range [44, 72]
    --- FAIL: .../accept_160  ball.size=160 must validate: ... out of range [44, 72]
MUTATION_RC=1
```

`size=160 必须接受` 这枚新牙**必须红、真红**。还原后 `grep MUTATION-160tooth` 零命中、全包复跑绿。

## 6. 门禁读数

| 门禁 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l internal/config` | 输出空，rc0 |
| go vet | `go vet ./internal/config` | rc0 |
| go build | `go build ./internal/... ./cmd/...` | rc0（注：`go build ./...` rc1 落 `build/mut_rvw.go`——`build/` 是 gitignore 的本地临时突变目录、非我射程，`git status --porcelain build/` 空、无 git 历史，预存在） |
| go test | `go test ./internal/config -count=1` | `ok ... 2.4s`；目标用例 10/10 subtest PASS |
| d22scan | `cd tools/d22scan && go run . -root ../../` | `clean - no D22 ban violations`，rc0 |

环境：跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；未设 `CGO_ENABLED=0`。

## 7. 44..72 全仓普查（射程外逐枚列出，⛔ 本腿未改）

| 位置 | 性质 | 归属 |
|---|---|---|
| `docs/specs/SPEC-08-ui-ball-panel.md:29` | 「配置范围 44–72px（默认 56）」 | 别的 spec，须另有授权/别腿 |
| `internal/ball/tokens.go:523`（注释）与 `:528`（`BallSizeMaxPx=72`） | 产码＋钳位 | **批 B** |
| `internal/ball/testdata/tokens_c21.css:222` | tokens golden 面 | 批 B／不动 golden |
| `cmd/wisp/config_reload_223_test.go:367` | 注释 | cmd/wisp 射程 |
| `frontend/src/components/config-screen.tsx:105` | 设置页文案「悬浮球直径（44-72）」 | frontend 归机主外部 agent |
| `design/old/screens/ball.html:430`、`design/old/assets/tokens.css:199`、`design/doubao/demo/screens/config.js:295` | 原型/demo | design 归外部 agent |
| `tools/jarvis-export/README.md:39` | 工具文档 | tools 射程 |
| `docs/evidence/s1/`（62/65/69/c21 等多枚） | **历史证据件**，如实记录当时契约 | ⛔ 不改写历史 |
| `docs/reports/pending-and-issues.md:32171-32197` | 台账 A1636 本身 | 台账（已落账） |

## 8. 交付物

- `internal/config/validate.go`（范围＋文案＋注释）
- `internal/config/schema.go`（结构体注释）
- `internal/config/validate_test.go`（边界迁移＋新用例＋band 字面）
- `docs/specs/SPEC-03-config-secrets-envs.md:27`（机主授权）
- `docs/PLAN.md:3557`（机主授权）
- 本件
