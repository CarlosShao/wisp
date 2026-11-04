# 361 — `attrib.sh --tracked-only` 本机 35 枚红：**本机尺与 CI 尺不同形**（第 3 例）

- Status: **done**（`AC#0` 只读现量可先派）
- 立案：2026-10-22 16:3x +08 编排者立
- ★**来路（非实现者签名）**：`227-v1` 第 6 项（非实现者终裁腿，件 `docs/evidence/s1/227-v1-verdict.md`、提交 `6e838c16`）
- 关联：票 161（`attrib.sh` 载具）· `ci.yml:240`（那一步）· `docs/evidence/s1/pipe-format-gate-reads-tracked-blob-not-worktree.md`（格式门禁裁定）

## 这是什么（人话）

`ci.yml:240` 那一步跑 `sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`。**CI 上它是绿的**，但**本机跑它是红的**：

| 红枚数 | 是什么 | 归谁 |
|---|---|---|
| **31 枚** | `.scratch/**` 台件（**故意坏的样本**） | ★`attrib.sh` **自己的归属规则允许**它们红 |
| **4 枚** | **CRLF 假红**（本机 `core.autocrlf=true`） | ★已由格式门禁裁定件判过：**读入库 blob** 时它们是干净的 |

⇒ **本机 35 枚红 ⛔ 代表 CI 会红**。**这是"本机尺与 CI 尺不同形"的第 3 例**：
- 第 1 例：`cmd/wisp/models.go` 的 CRLF 假红（已裁：读 blob）
- 第 2 例：`cmd/wisp` 整包在**含空格路径**下被打死（已另立票）
- ★**第 3 例＝本票**：`attrib.sh --tracked-only` 的本机读数与 CI 读数**不同形**

### ★★更正块（2026-10-22 17:0x，编排者；依据 `361-a1` 普查腿 `5182cb74`）——**上面那段的前提是错的**

> **原写**：「31 枚是 `.scratch/**` 台件（`attrib.sh` 自己的归属规则允许）、4 枚是 CRLF 假红」⇒「本机 35 枚红 ⛔ 代表 CI 会红」。
> **现读**（普查腿逐枚现量、我复认关键两处）：**这个 31/4 是"路径分裂"，⛔ "原因分裂"**。真原因是**两类**：

| 类 | 枚数 | 判据（`blob_CR` 对 `wt_CR`） | 是什么 |
|---|---|---|---|
| **B** | **20** | `wt_CR=0` 且 `blob_CR=0`（**纯 LF 仍脏**） | ★**入库 blob 本身就脏**——**全是台件** |
| **C** | **15** | `blob_CR=0`、`wt_CR>0` | **CRLF 假红**——台件 11 ＋ 非台件 4 |

**★两条硬读数（决定性）**
1. ★**`--tracked-only` 下归属规则根本不执行**：`classify_line` 只在 `:248`（`--classify-line`）／`:328`（`--self-test`）／`:593`（`ruler_b`）被调用；
   `ruler_a`（`:539-583`）**从不调用它**，而 `--tracked-only` 只跑 `ruler_a`。`:52` 逐字 `rc=1 (A) printed anything at all -> tracked bytes are dirty; CI red`
   ⇒ **该模式"打印任何东西即红"、台件与 `internal/` 产码同权**。实测反证：`--classify-line .scratch/wisp/probes/163/a1/main.go` ⇒ `TRACKED<CI-RED> rc=1`。
2. ★**仓外 `git archive` 导出树（纯 LF、无台件状态）上仍红 20 枚**：`rc=1、lines=20 files=20`（35→20，掉的 15 枚正是全部 C 类）；
   这 20 枚与 HEAD blob **逐枚 sha256 相同（20/20）**、与已推送 `626fcc7b` 的 blob **也逐枚相同（20/20）**；
   **其中 15 枚连 `gofmt -l` 都判脏**；★**把 `ci.yml:172` 那步形状在同一树跑、得到同样 20 行**。

**⇒ 因此原票面的"第 3 例"定性也不成立**（这是**两件事**，⛔ 一件）：
- **15 枚 CRLF 假红** ⇒ ★**这才是"本机尺 ≠ CI 尺"的第 3 例**（`attrib.sh` 侧可修）；
- **20 枚入库 blob 真脏** ⇒ ★**这是"仓库里真有 20 枚没格式化"**（与 CI 无关、与 `autocrlf` 无关）⇒ **另立票、⛔ 归 `attrib.sh` 射程**。

**★另倒一条断言**：原写"CI 上它是绿的"——★**同样未经我验证**。普查腿查到本机 CI 日志里 `lint` 作业**只留了 D22 正控与 staticcheck 两步、没有 `attrib.sh` 那一步的行**
⇒ **"CI 上它是绿"目前没有任何 CI 侧原件支撑**，按 `AC#3` 口径登记为〔待推送取数〕。

**★射程收窄**：本票**只做 15 枚 CRLF 假红那一半**（`attrib.sh` 侧）；20 枚真脏**另立票**。

---

## 结案判据

- [x] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：把 35 枚量成可复跑读数——
      ①**逐名**列出 35 枚（⛔ 只给计数），并**逐枚归类**（31 台件／4 CRLF）；
      ②**`attrib.sh` 自己的归属规则**逐字引出（它在哪一行、怎么判"这枚属台件"）；
      ③**正控**：造一枚**真脏**的 tracked `.go`（非台件、非 CRLF）⇒ `attrib.sh --tracked-only` **必须红并点名它** ⇒ 证明尺是活的、⛔ 恒红；
      ④**两侧对照**：在**仓外 `git archive` 导出树**（无 `.scratch` 台件、纯 LF）上跑同一发 ⇒ 给出读数（★这才是与 CI 同形的那个）。
- [x] **AC#1 裁形（编排者拍）**：三形各写代价＝〔甲〕**只写清"本机不可判"**（文档一句，⛔ 动脚本）；
      〔乙〕让 `attrib.sh` **本机也按入库 blob 读**（与格式门禁裁定件同口径）；〔丙〕两者都做。
      ★**未裁形任何实现都⛔ 开工**。
- [x] **AC#2 落地＋牙（写腿）**：按裁定那一形做；每处改动配**一枚今天响、不修不响**的读数；
      ★反向判据＝把新加的那一跳摘掉 ⇒ 必须转红并**逐名报**；三读数逐字贴 ＋ **还原核 `sha256`**。
- [x] **AC#3 门禁**：`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh`／`sh .scratch/wisp/probes/161/r5/attrib.sh --self-test`（★票 161 的自检载具必须仍绿、且它的正控仍红得起来）各 rc。
      ★★**CI 侧读数⛔ 本机冒充**（机主说过暂不推远程）⇒ 若判据要 CI 数，按〔待推送取数〕登记。
- [ ] **AC#4 纪律格**：⛔ 翻任何框（本票各格由非实现者裁）、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 161 任何原句。

## 硬约束

- ⛔ 改 `.github/**`（★那是**另一次具名解冻**，属"改 CI"、须机主点头）；本票只许动 `attrib.sh` 或其读数件。
- ⛔ 改 `thresholds.go`／golden／`allowlist.txt`／三枚冻结件；⛔ `frontend/**`／`design/**`（零读零写零转述）。
- ★**⛔ 把"本机红"当"CI 会红"**——本票存在的全部意义就是**这两者不同形**。

## 排程

`361-a1`（只读＝`AC#0`）→ 编排者裁 `AC#1` → `361-r1`（写腿＝`AC#2`）→ `361-v1`（非实现者裁）。

## Progress log

- [2026-10-22 16:3x +08] agent=`编排者` did=立案（来路＝`227-v1` 第 6 项；★票面写明"本机尺与 CI 尺不同形第 3 例"、并列前两例）next=派 `361-a1` 只读腿。

## 159轮现量更新（编排者，HEAD 12d3ae7f2）——票面前提再次收窄
- 现跑 `.scratch/wisp/probes/161/r5/attrib.sh --tracked-only` rc=1，仅 **2 枚**脏：internal/agent/approval/pending_read.go、internal/agent/tools.go（票面旧 35/20/15 已随 HEAD 演进）。
- 判真脏 vs 假红：两文件对 **HEAD blob 的 gofumpt -l=0（干净）、对 worktree=1（脏）**，core.autocrlf=true ⇒ **纯 CRLF 假红**（worktree 被 autocrlf 转 CRLF，gofumpt 读 worktree 误判；CI 读 LF blob 故绿）。⇒ 当前 **0 枚真脏**（票面「20 枚真脏」那一半在现 HEAD 已不成立/随历史修掉）。
- AC#1 裁形（编排者拍，按铁律最有价值/主流）＝**〔乙〕**：让 attrib.sh ruler_a 本机也按**入库 blob** 读（与格式门禁裁定同口径），⛔ 不选只写「不可判」的〔甲〕。AC#2 实现下次（须保 --self-test 正控仍红、真脏 blob 仍被点名）。AC#0 ③正控/④archive 对照待实现同程补。

## AC#2 实现规格＋基线（编排者 160轮；实现腿/fresh 照此一次做完）
- 改点：ruler_a `:585` 的 `A_OUT=$(git ls-files -z "*.go" "$A_N_EXCLUDE" | xargs -0 "$GOFUMPT" -l 2>&1)` 现读 **worktree(CRLF)** ⇒ 改读**入库 blob(LF)**。
- 法A（推荐）：`git archive HEAD` 导出到 tmp（.gitattributes `*.go text eol=lf` ⇒ 导出树 LF），tmp 内剔 `.scratch/**`，`cd tmp && "$GOFUMPT" -l <五生产目录或 .>`，输出即 repo-relative；A_OUT 语义不变、下游循环不改。
- 法B：逐文件 `git cat-file blob HEAD:$f` 到 tmp（1089 子进程，慢）。
- 必保 3 边：①`--self-test` rosterctl 正控（:568、`/tmp/wisp-275-r2-rosterctl.*`，要 gofumpt rc=2＋A-ROSTER 点名）不受影响；②parse-error rc>=2 分支（:586-609）语义不变；③A_FILES 的 `:N:N:` 剥离（:623）不变。
- 基线（160轮现跑）：`--self-test` rc=0（cases=10 failures=0 GREEN）；`--tracked-only` rc=1 tracked-dirty=2。目标：改后 `--tracked-only` rc=0（blob 干净）、`--self-test` 仍 rc=0、真脏 blob 仍被点名红。
- 牙：造一枚真脏 blob（`git hash-object -w` 脏内容＋临时 index 替换该路径）⇒ 改后 attrib 仍红点名；摘掉 blob 读取跳 ⇒ 回 CRLF 假红。
- ⚠ 本轮未实现（编排者上下文近枯＋改坏 gate 脚本风险高于收益）；spec＋基线已备，fresh/腿 一次做完再非实现者裁。

## AC#2 实现（编排者 161轮，〔乙〕已落：ruler_a 读 index(LF) 取代 worktree(CRLF)）
- 改动＝`attrib.sh` ruler_a A_OUT 源：`git -c core.autocrlf=false checkout-index -a --prefix=<tmp>` 写 LF 树（index 字节＝CI 同形）→ 剔 `.scratch` → gofumpt -l；含**防空 0 假绿**守卫（临时树 0 枚 .go 则回退 worktree 读，绝不报空绿）。self-test 抛树用 index(git add 无 commit)，走此路径不变。
- 四读（现跑，全在 HEAD 工作树）：
  ① `--self-test` rc=0 cases=10 failures=0 GREEN（证 index 读仍抓 parse-error CTL_BROKEN 与脏 CTL_UNFORMATTED 并点名）；
  ② 新版 `--tracked-only` rc=0 GREEN 0 脏（CRLF 假红消失，1089 枚）；
  ③ 真脏牙：临时 stage 畸形 `internal/agent/zz_tmp_tooth_dirty.go` ⇒ rc=1 逐名点名该文件（非假绿），已 `git rm --cached`+rm 清理、index 干净；
  ④ 反向牙：`git show HEAD:attrib.sh`（原版读 worktree）⇒ rc=1 红 2 枚（pending_read.go/tools.go）——差异归本次改动。
- 待非实现者 v1 独立复跑四读＋复核 diff 后裁收口（裁决者≠实现者）。

## 收口（161轮，编排者；v1 非实现者 PASS）
- AC#2 实现经 v1 独立复跑四读 PASS（self-test 10/0 绿／新版 tracked-only 0 脏绿／真脏点名红／原版红 2 枚）＋两路字节源互证（checkout-index 与 git archive 均 1089 枚 0 脏）。
- AC#3 门禁：d22scan rc=0 绿；attrib --self-test rc=0 绿；path-length-budget rc=1 **红与本票无关**——红的是 5 枚本会话 `-done` 改名后超 165 预算的票名（128/175/273/364/403），属 owner 待办 ⑤ path-budget roster 漂移，非 attrib.sh 改动引入。
- v1 遗留 2 条已记：①`:582` 自述行仍印旧管道（文案滞后）；②守卫只防"提取为 0"不防"部分失败致分母缩水"——留后续小票。
