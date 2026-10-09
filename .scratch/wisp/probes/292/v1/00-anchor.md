# 292-v1 00 — 锚点与尺名册

## 锚（起手第一条命令链的 stdout 逐字）

- `date` => `Fri Oct  9 18:19:45 CST 2026`
- `git log --oneline -1` => `bccafc30 296-a1 只读普查交回：AC#0 射程名册（6 真落点＋4 注释档，排除 258 变异拷贝 8 枚；多量到 resident_windows.go:213 第二形 live=0）＋AC#2 甲乙代价表与 none 整族穷举（判契约空白）；票面仅追加 Progress log 一行，零翻框零产码`
- 被验改动 commit = `4405e0f4 票 292 腿 292-r1 AC#2：:598-600 三行注释改写成说实话的三节（零语义）`
- 另两笔（证据与票面 Progress log）= `28292c05` / `e390fdd9`（本件起手时尚未在 `git log --oneline -8 -- cmd/wisp/subagent_carrier_197_test.go` 中出现，待 10 复核）
- 分支 = `dev`；写面仅 `.scratch/wisp/probes/292/v1/*.md` ＋票面末尾 Progress log 一行

## 起手环境读数

- `tasklist //FI "IMAGENAME eq wisp.exe"` 与 `tasklist //FI "IMAGENAME eq balldebug.exe"`：见 10-anchor-env 之后各发读数；本文件记录起手那次调用（详见 `20-runs.md`）。

## 我要用的尺名册（逐把写明口径）

| 尺号 | 用途 | 逐字命令 | 口径声明 |
|---|---|---|---|
| M1 | 被验改动正文 | `git show 4405e0f4` | 全量 |
| M2 | 逐枚同枚数（AC#2 主证） | `git show --numstat 4405e0f4` | 整族穷举（该 commit 名册全集） |
| M3 | 纯注释行枚数＝全量枚数 | `git show -U0 4405e0f4` 取非 `+++`/`---` 行，数 `+`/`-` 前缀行 | 整族穷举 |
| M4 | 三行新注释逐字 | `sed -n '598,605p' cmd/wisp/subagent_carrier_197_test.go` | 抽样（行号窗） |
| M5 | join 来源枚数（ⓐ） | `sed -n '205,240p' internal/panel/subagent_roster_197.go` | 整族（建表函数体全窗） |
| M6 | card corr 形状（ⓑ） | `sed -n '595,610p;670,685p' internal/agent/loop.go` | 整族（callCorr 两支） |
| M7 | 生产装配 L1Windows（ⓒ） | `grep -c L1Windows cmd/wisp/run.go` ＋ `grep -rn "L1Windows" --include=*.go .` | 整族穷举（全仓 .go） |
| M8 | tool_choice 两环（ⓓ） | `grep -rn "ToolChoice" internal/agent/*.go` ＋ `grep -rn "tool_calls\|ToolCalls" internal/llm/adaptertest/mockllm.go` | 整族穷举（那两枚件全文） |
| M9 | 整包两发三数尺 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v`；`grep -cE '^[[:space:]]*--- PASS'` / `--- FAIL` / `--- SKIP` | 每态 2 发，红名册逐名作差 |
| M10 | 门禁四把 | `GOFLAGS= go build ./...`；`sh scripts/d22scan.sh`；`gofmt -l`；`$(go env GOPATH)/bin/gofumpt.exe -l` | 整族（HEAD blob 与工作树各一把） |
| M11 | 归口越界 | `git show --name-status <三笔>` | 整族穷举（三笔名册） |
| M12 | CRLF 假枚复核 | `git config core.autocrlf`；`git ls-files --eol cmd/wisp/models.go`；`tr -cd '\r' \| wc -c`；`git status --short cmd/wisp/models.go` | 单枚 |

## 预算自缚

上限 45 次工具调用；两发整包各计 1 次。逼近上限先 commit 再具名写未攻完的问。
