# 票 254 — 第五档 `winsec` 今天**在 CI 里没有调用者**：`winsec-tests.sh` 仍走显式路径，而 `ci.yml` 的形状属契约级（251-r1 按派单停手没动）

**立票时刻**：2026-10-02 10:2x +08，锚点 HEAD `0cdd3dd3`（`dev`）
**来路**：`251-r1` 交件时的两条"甲（谁补得上）"＋一条"乙（停手未动）"；台账 `A523`。**不夹进票 251**（那票四格已做完，本票是"把新门接上电源"这一步）。
**性质**：仪器／CI 口径。owner 用软件撞不到；⚠ 但**不接上，票 251 那道新 GUARD C 在 CI 上永远不会被触发**。

## 现量（编排者本机自跑，2026-10-02 10:2x）

1. `bash scripts/portable-tests-selftest.sh` ⇒ **`18 case(s) ran, 0 assertion(s) failed`、`rc=0`**（载具已含第五档的两形正控）。
2. `bash scripts/portable-tests.sh --scope=winsecfoo` ⇒ 逐字 `unknown --scope=winsecfoo (known: core, windows, cli, winsec, census)`、**`rc=2`** ⇒ 第五档进名册了、未知档仍响亮拒答。
3. 腿报（我未复跑，属〔仅自述〕级，落地腿起手要复认）：`scripts/winsec-tests.sh:94` 那一行**把包路径显式传给** `portable-tests.sh`，CI 侧由 `.github/workflows/ci.yml:439` 调 `winsec-tests.sh` ⇒ **今天没有任何一步以 `--scope=winsec` 的名义去跑那档**。

## 要建什么

- [ ] **AC#1 把第五档接上电源（二选一，先交代价再动手）**：ⓐ 让 `scripts/winsec-tests.sh:94` 改点 `--scope=winsec`（写面只在 `scripts/`，最小）；ⓑ 保持显式路径不变，但**证明显式路径模式今天也真被 GUARD C 对账**（要给出"配钉后拿 glob 解析集"那一条的**同种子双形读数**：拆包必响、正常必绿，⚠ 这正是 `251-r1` 推翻我 AC#1 原判据那一格——**目录形会照绿，所以"我在显式路径上"不等于"我被对账了"**）。判据＝两种形都要有"改前必红"的发数摆着，选哪一支都要具名凭据。
- [ ] **AC#2 `core` 档里那行 `./internal/winsec/`（`scripts/portable-tests.sh:211`〔行号待验〕）与第五档的关系要定死**：⛔ 不许两档各自认领同一枚包而对账结果不同；要改 glob 就**同笔**动 `core_pin`，并给"漏计／双计"两形各一枚正控。
- [ ] **AC#3 unknown 消息里那份 known 列表（从 4 名扩到 5 名）的口径裁清**：要么"消息文本属射程、不许扩"（则第五档与它不能并存，需改文案生成方式），要么"列表随档自动生成、扩名不算改动"。**这一条由编排者裁，不由腿自选**；腿只交两形的代价与现量。

## 禁区

- ⛔ **`ci.yml` 的形状不属本票射程**：动 workflow 触发档＝契约级（票 134 的 C+B 形状、见 [[wisp-ci-selfhosted-topology]]）⇒ **停手上报，不许"顺手加一步"**。
- ⛔ 不许把 GUARD C 的**双向对账**改成单向或跳过；⛔ 不许把 unknown scope 的 `rc=2` 改软；⛔ 不许为了让第五档"被调用"而往 `core` 档塞包名（那是把两档的钉搅在一起）。
- ⛔ 不许动 `internal/**` 的任何测试或产码（本票射程只有 `scripts/` ＋证据件）。
- `frontend/**`／`design/**` **既不读也不引**；grep/find 必写显式根（`scripts docs tools cmd internal .scratch`），⛔ 不用 `.` 当根。
- Git 纪律：只 commit 不 push；显式 pathspec；⛔ `git add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 交件要求

起手 `date -Iseconds`＋`git log -1`＋`git status --porcelain scripts/`（⛔ `scripts/` 此刻若有别的写腿在飞就停手上报——票 250／251 都在这面上）；**第一轮内先落骨架并 commit**；跑到第 100 轮前先把「门禁读数」与「判不动的地方」两节写满并 commit（⛔ 空着不算交件）。占位符自查用字符类尺在命令行跑、⛔ 尺的字面文本不落进被扫文件；计数尺一律 `| wc -l` 收尾。⛔ **一枚 AC 勾选框都不许碰**（翻勾归编排者）。写明**它推翻编排者上面哪一句**（上面三行现量与所有行号都是待验断言）。

## 排程

- 写面＝`scripts/`（此刻**已空**：`251-r1` 交件退出）。
- ⚠ 与 `197-r3`（写 `cmd/wisp`）无关，可并发；⛔ 但**不许与另一枚会跑整包 `go test` 的腿并发**（`248-v1` 正在跑 `cmd/wisp`，本腿不跑 Go 测试即可）。
