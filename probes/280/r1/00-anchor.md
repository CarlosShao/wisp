# 280-r1 起手锚（本腿自己的取数，不引任何人的快照）

腿名：`280-r1`。票面：`.scratch/wisp/issues/280-seven-files-fail-plain-gofmt-decide-gofmt-vs-gofumpt-ownership.md`（20 行，三格 AC + 禁区，本腿整行读过）。
取数时刻：`2026-10-09 09:49 +0800`。

## 锚

- `git rev-parse HEAD` = `bf9974c229974b4dbd288d0c2aa49a41673d8243`
- `git rev-parse --abbrev-ref HEAD` = `dev`
- `git log --oneline -3`：
  - `bf9974c2` A751/A752 落账：收 281-r1（翻框审计 10 对 10 逐对同文＋Status 两枚原句逐字留档，07/115 一字未动）＋收 283-r1（新尺 227 行 +0/-0 提交形状即硬）⇒ 票 283 三格翻勾并 -done
  - `c87ad40e` 283-r1 收件：查重三处＋两处补尺落点＋三发正控四件套读数（种 A/B/C 皆红并逐字还原）
  - `58a4b2fe` 281-r1 收口件：八张归位逐笔凭据＋AC#4 复跑对拉（未勾 14→4、-done 98 不变、零撤名）

## 本腿自己的名册（射程逐条写清）

尺 A：`gofmt -l cmd internal`（票面射程，Windows 分隔符原样打印）＝ 7 命中：

```
cmd\wisp\models.go
cmd\wisp\panel_inbound_guards_35r3_test.go
cmd\wisp\panel_transport_35r2_test.go
internal\agent\approval\pending_read.go
internal\agent\tools.go
internal\risk\provenance.go
internal\tools\bridge.go
```

尺 B：`gofmt -l internal cmd tools`（编排者转述的射程）＝ **同一 7 枚**，只是打印顺序不同。
⇒ 具名读数：`tools/` 目录对裸 `gofmt` 的贡献是 **0 枚**，两条射程的名册**没有漂**，票面 7 枚与盘上 7 枚**同名同数**。

尺 C：`gofumpt -l <这 7 枚>`（本机 `gofumpt.exe`，见下）＝ **7 枚全中**，rc=0。
⇒ 本票的问题不是"gofumpt 是否也点这 7 枚"（点了），而是"点红的**理由**是否相同"。

## 尺本身（重要：编排者的转述在这里是错的，取盘上为准）

- `command -v gofumpt` → 不在 PATH。⚠ 但**尺存在**：`ls /d/work/base/gopath/bin` → `gofumpt.exe`、`staticcheck.exe`；
  `gofumpt --version` = **`v0.12.0 (go1.27.1)`**。
- `GOMODCACHE` = `D:\work\base\gopath\pkg\mod`，里面有 `mvdan.cc/gofumpt@v0.12.0` 与 `@v0.7.0` 两份源码。
- 本腿**没有** `go install`、没有下载、没有 `-w`；跑的全部是 `gofumpt -l` / `gofumpt -d` / `gofumpt < file`（只读）。
- 版本差登记：CI 命令是 `go install mvdan.cc/gofumpt@latest`（`.github/workflows/ci.yml:175`），本腿的尺钉在 `v0.12.0`。
  ⇒ 本件里凡"CI 那一步会不会点红"的判断，凭的是**这把尺对 CI 射程内同一份字节的发声**，不是对某次 runner 的观测。

## 禁区自查

写面只有 `probes/280/r1/*.md`；零 `-w`、零重排、零产码字节改动、不改票面与台账、不翻任何 AC 框、零 push、
不跑 `go build`/`go test`/`go vet`/`go list -deps`；`frontend/**`、`design/**` 零读零写；工作树里别人的在飞改动未触碰。
