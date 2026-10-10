```
$ date
Sun Oct 11 07:52:27 CST 2026
$ git log -1 --oneline
ec87076f 306-v1 起手锚（验收腿 commit-first 第 1 笔）：台面＋重锚 HEAD=e4740e35（派单写的 6284a489 已腐烂，具名报回）＋scoped porcelain＋进程闸门自陈
$ git rev-parse HEAD
ec87076f1c6a90c45facf3558635582fb0ead989
$ git rev-parse --abbrev-ref HEAD
dev
$ git status --porcelain -- .scratch/wisp/probes/risk-attrib docs/reports
（0 行）rc=0
$ git status --porcelain | wc -l
841
```

# risk-attrib/1 —— 起手锚（A817 那 12 枚 CI 红的只读归因普查腿）

## 1. 台面

- 腿名＝`risk-red-attrib-1`；射程＝CI `--scope=windows` 在托管 `windows-latest` 上的 **12 枚恒红**（台账 `A817` §3/§5）。
- 派发时刻 HEAD **现量** ＝ `ec87076f1c6a90c45facf3558635582fb0ead989`（短号 `ec87076f`），分支 `dev`。
  ⚠ 本单派发文字里没给 HEAD 号，只有 `A817` 正文里的 `bcd0a543`／`cf46c24a`／`05db4bc6`／`e4740e35` 那几枚 —— 那些是**别的发的 headSha**，⛔ 当本腿台面引。
- 工作树别家脏面现量 ＝ **841 行**（尺＝`git status --porcelain | wc -l`）。本腿自家 scoped 区间（`probes/risk-attrib`＋`docs/reports`）起手 **0 行**。
  ⇒ 别家脏面（`.gitignore`、`design/**`、别家 `probes/**`、票面）**⛔ 碰／⛔ 暂存／⛔ 提交**，本腿 scoped pathspec 只落 `.scratch/wisp/probes/risk-attrib/1/**`。
- 本机既有事实（编排者派单给的，本腿不重测、只声明不据此下判）：机主仓根**确实存在 8.3 别名**、TEMP **没有**、生产码**不调** `GetShortPathNameW`。
  ★本腿会自己复核"仓根 8.3 别名"那一枚（它直接决定表② 的甲形能否在本机复现），复核读数落 `logs/`。

## 2. 我允许跑的 Go 命令（自陈，越出这一行的都算越界）

- ✅ `go test ./internal/risk/...`（定向 `-run '<具名>'`／`-count=1`／`-v`）
- ✅ 表① 量出的**这 12 枚实际所在包**的 `go test`（表① 现量＝只有 `internal/risk` 一枚包，见 `10-ownership-roster.md`）
- ✅ `go build` / `go vet` **只限 `./internal/risk/...`**
- ✅ `go list` / `go env`
- ✅ 只读 `git` / `grep` / `sed`
- ⛔ `./cmd/wisp/`（含整包 `./...`、`-tags winlive`）＝编排者独占车道
- ⛔ `scripts/d22scan.sh`／⛔ `gofumpt`／⛔ `gofmt -l` 名册类
- ⛔ 任何写产码／写 `_test.go`／新建 `.go`；⛔ 为变绿放宽断言；⛔ 造 `--- SKIP`
- ⛔ 真窗／真机／GUI（机主在用这台电脑）；⛔ 读 `%APPDATA%\wisp-dev\secrets\`；⛔ 读机主 config 的值
- ⛔ `git add -A`／`.`、⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`
- ✅ 只 commit、⛔ push
- 复现 runner 形状（假 `%USERPROFILE%`、8.3 短名、under-profile）⇒ **只在 `C:/Users/swq/tmp/<具名>/` 里造**，进程内 env 覆盖；⛔ 写进仓、⛔ 改全局/用户级配置、⛔ 动 git config；临时件**只建不删**、落本目录 `logs/`。

## 3. 重锚：这枚欠账在盘上是什么

- `A817` 标题行现量＝`docs/reports/pending-and-issues.md:15313`（行号一律当快照）。
- `A817` §3 的三档读数（既有读数，本腿**复跑、⛔ 照抄**）：5 枚红句逐字含 `C:\Users\RUNNER~1`／2 枚"同文件同族但红句不含那串"／5 枚 `syncdirs_test.go` under-profile fallback 家族。
  ⚠ 本腿起手就把 §3 那句"**2 枚同文件同族**"当**待核**而不是当事实（见 `10-ownership-roster.md` §4 的冲突具名）。
- 尺的坑（`A817` §3 自陈踩过）：Go 的 `-v` 把 `t.Errorf` 明细写在 `--- FAIL:` **之前** ⇒ 正确射程＝`=== RUN` ↔ `--- FAIL` 之间。本腿所有"红句里有没有某串"的存在性尺**一律用这一把**，并在件里写死射程与方向。
- CI 名册出处＝已入库的 `.scratch/wisp/probes/303/orch/r9-ci-after-red-roster.txt`（83 行）。⛔ 读 MB 级原始日志。

## 4. 本腿要交的三张表

1. `10-ownership-roster.md`＝表① 归属与包名册（分母）
2. `20-mechanism-per-red.md`＝表② 逐枚机制归因（甲/乙/丙/丁 ＋〔已证／读码推／缺世界〕标签）
3. `30-routing.md`＝表③ 归口建议（该立／该挂／该登记 三选一，⛔ 动手、⛔ 新造票号）
4. `90-final.md`＝现态＋⛔ 做到的名册＋逐笔越界尺＋与本单不符处
