# 297-a1 · 00 起手件：锚与尺具名册

腿号 `297-a1`（只读普查）。派单人＝编排者（人类）。本件是该腿**第一笔**交付，落笔时刻在任何长跑命令之前。

---

## 1. 锚

### 1.1 锚漂了（具名申报，这是本件最重要的一行）

- 派单给的锚：HEAD `cfacfd3c`。
- 我 `git rev-parse --short HEAD` 现量：**`e96da9f4`**。
- ⇒ **本腿的读数锚一律改为 `e96da9f4`**，下面所有"对象层/工作树"的对尺都以这枚号为准。

漂移关系（现量，非推算）：

```
$ git cat-file -t cfacfd3c
commit
$ git merge-base --is-ancestor cfacfd3c HEAD
YES_ANCESTOR
$ git rev-list --count cfacfd3c..HEAD
1
```

⇒ 漂的这 **1 枚**（尺＝`git rev-list --count cfacfd3c..HEAD`，射程＝全史）不是别人，正是派单里点名"此刻在改 `cmd/wisp/`"的那枚写腿：

```
$ git log --oneline -14
e96da9f4 票 293 · 写腿 293-r1 第 5 笔：AC#4 门禁件（改前两发基线自取＝7 枚红并集、改后两发逐名作差新增红 0、gofumpt 两口径冲突具名、d22scan rc=0、两发中间态作废留痕）＋ 15 枚原始 logs ＋ 票面 Progress log 追加一条（只追加，diff＝1 insertion／0 deletion，⛔ 未翻任何框）
cfacfd3c 收我自己跑的 299 `AC#6`（账 A798）⇒ 翻 `AC#0`／`AC#6` 两枚、`AC#1`..`AC#5` 五枚一枚未碰：★那枚冻结 CSP `<meta>` 在 `SetHtml` 产出的文档里真执行（同一台件两发：不带 meta＝经典内联／内联模块／blob: 三形全跑；带上 meta＝三形全不跑）⇒ 普查腿判"乙形＝内联整包结构性不成立"从今天起有实测支撑；★判活器＝宿主 `Eval` 通道两发都落到 Go ⇒ B 发三形沉默只能读成"页面脚本被拦"，⛔ 不是"桥断了"（没有这一枚整发白跑）；★`document.baseURI` 实测逐字 `about:blank` ⇒ `◈` 具名欠的那一枚销账（⚠ 射程＝取自不带 meta 那一发）。
```

⚠ 上面第二行的原文极长，我在本节里**只逐字抄到该行开头一段**，其余部分【⚠ 更正：本括号里的三串归属写错了，连同后面那句"截断点在……"一起作废，正确说明见 §5 更正①】（`AC#0`/`AC#6` 翻框那串、`36 枚命中`那把尺那句、以及"⛔ 零 push"尾部）**未抄全**。所以本块标的是「**逐字（截断）**」而不是「逐字」——截断点在 `⛔ 不是"桥断了"` 之后的原文仍在 `git log --oneline -14` 那一行里，要看全量自己现跑。⇒ 派单早上批评过"省行又写逐字"，这一处我按同一纪律自己先认。

### 1.2 分支与时刻（现量）

```
$ git rev-parse --abbrev-ref HEAD
dev
$ date "+%Y-%m-%d %H:%M:%S %z"
2026-10-10 10:39:39 +0800
```

---

## 2. 读数来源纪律（本腿自己给自己立的，后面四件都受这条管）

- **写腿 `293-r1` 的射程＝`cmd/wisp/`**（派单原文 + `e96da9f4` 的提交标题自证）。⇒ 我在 `cmd/wisp/` 里的**任何**读数，都必须满足下面两条之一，并**在引用处写明是哪一条**：
  - (A) **工作树读**，且同发证明该文件 `git status --porcelain -- <path>` 输出为空 ⇒ "工作树内容 == `e96da9f4` 该文件"（这是等价的，不是假设）；
  - (B) **对象层读**，逐字写 `git show e96da9f4:<path>`。
- `internal/audio/`、`.scratch/wisp/issues/`、`scripts/`、`.github/` 不在写腿射程内 ⇒ 默认工作树读，但对每张我要贴"逐字"块的件，**落笔前补一发 (A) 的 status 尺**证明我贴的与工作树/HEAD 一致。
- **票面行号一律当过期快照**（派单纪律）⇒ 我只用**内容锚**（逐字那一行 + 函数名/测试名），行号只作旁证；每处注明读的哪一层。
- 派单里我**不接受**默认成立的两条转述，都要自己复跑尺核（核不上就具名报回）：
  1. "机主已答'听见了'、已入账 `A793`" ⇒ 尺＝`grep -n '^## A793' docs/reports/pending-and-issues.md`。
  2. "`convertPacket` 唯一的生产调用点在真设备上 ⇒ 本仓内不可能被测到" ⇒ 尺＝`git grep` 调用点名册（见 10 件）。

---

## 3. 我已经跑过的尺（逐字命令，截至本件落笔）

编号供后面四件回指。**每一发的射程目录都写在右边**。

| 号 | 逐字命令 | 射程目录 | 现量/算的 |
|---|---|---|---|
| R1 | `cd "D:/work/workspace/projects plans/Wisp" && git rev-parse --short HEAD && git rev-parse --abbrev-ref HEAD && date "+%Y-%m-%d %H:%M:%S %z"` | 仓根（只读 git 元数据） | 现量 |
| R2 | `cd "D:/work/workspace/projects plans/Wisp" && ls .scratch/wisp/ && echo "--- probes? ---" && ls -d .scratch/wisp/probes 2>/dev/null \|\| echo "NO probes dir"` | `.scratch/wisp/` 一层 | 现量 |
| R3 | `cd "D:/work/workspace/projects plans/Wisp" && echo "=== cat-file -t cfacfd3c ===" && git cat-file -t cfacfd3c 2>&1 && echo "=== is-ancestor cfacfd3c HEAD ===" && (git merge-base --is-ancestor cfacfd3c HEAD && echo YES_ANCESTOR \|\| echo NO_NOT_ANCESTOR) && echo "=== log oneline -14 ===" && git log --oneline -14 && echo "=== commits between cfacfd3c..HEAD ===" && git rev-list --count cfacfd3c..HEAD` | 全史（截 14 行标题） | 现量 |
| R4 | `cd "D:/work/workspace/projects plans/Wisp" && ls -R .scratch/wisp/probes 2>/dev/null \| head -40 && echo "=== gitignore line 8 ===" && sed -n '1,12p' .gitignore` | `.scratch/wisp/probes/**`（截 40 行）+ `.gitignore` 前 12 行 | 现量 |

R4 的两条用途：
- `.scratch/wisp/probes/` 下**已有 111/114/132/139/145/…的票号目录**（我只看了前 40 行输出，⛔ 未数全）⇒ 本腿落点 `probes/297/a1/` 是新目录，不与别人撞。
- `.gitignore` 第 8 行**确认为 `*.out`**（逐字：`*.out`）⇒ 派单那句"根 `.gitignore` 第 8 行是全仓 `*.out`，会被静默跳过"**核上了**，本腿⛔ 不在仓内建任何 `.out`；原始大输出一律落 `TMPDIR`（仓外）再读。

---

## 4. ⛔ 我没跑的尺（具名归口，落地腿不许引用我的旧读数）

| 没跑的东西 | 为什么没跑 | 直接后果 |
|---|---|---|
| `go build ./...` | 派单禁跑（写腿 `293-r1` 在 `cmd/wisp/` 编译面上工作） | **本腿零枚"改前基线颜色"读数**。谁要写"改前绿/改前 N 枚红"必须**自己现跑**，⛔ 不许引我这腿 |
| `go vet ./...` | 同上（禁跑编译面） | 同上 |
| `go test ./...`（含 `-run`、含 `-tags windows` / `-tags winlive`） | 同上 | 本腿对"某用例今天到底跑不跑、跑成什么颜色"**一律只能读码＋读 CI 配置**，读到就算"码层/配置层答案"，读不到就写"未证成" |
| `tools/d22scan` | 未派给我，且属编译面外的可执行程序；我没自行跑 | "emoji 禁带会不会被这批改动顶红"我在 40 件里**只能列名册、不能给 rc** |
| 真机音频（任何 WASAPI 采集/回放、`wisp run` 真链路、电平尺实跑） | 本腿是**只读普查**，派单⛔ 不修不跑；且真机读数属编排者那一格 | 本腿对 `0.3677` / `-8.6 dBFS` **只做溯源与形状判定，一律不复现、不背书** |
| `docs/` 下任何写动作、票面任何 `- [ ]` 框 | 派单⛔ | 本腿五件全落 `.scratch/wisp/probes/297/a1/`，票面只**末尾 append 一行** |

**允许跑但我需要自报的**：`go env` / `go list`。截至本件落笔，**我一枚都还没跑**（后面若跑了，会在本件 §5 追加并单独 commit）。

---

## 5. 追加区（只在真的发生"新增尺/新增 go env·go list 自报"时写，只追加不删）

### 5.1 更正①（本腿自查，落笔第 1 笔之前）

§1.1 那句"其余部分（`AC#0`/`AC#6` 翻框那串、`36 枚命中`那把尺那句、以及"⛔ 零 push"尾部）未抄全"**是我写错的**：
- `AC#0`／`AC#6` 翻框那串**就在第 1 句、我已经抄了**；
- `36 枚命中`那把尺、`⛔ 零 push` 尾部**根本不在 `cfacfd3c` 这一行里**（它们属更早的 `eb65db8e`，我把它串的标题混进来了）。

**正确的截断说明**：`cfacfd3c` 那一行我只抄到 ★`document.baseURI` 实测逐字 `about:blank` ⇒ `◈` 具名欠的那一枚销账（⚠ 射程＝取自不带 meta 那一发）。**该行为止**，即这一枚 commit 标题**我其实抄全了**；不完整的只是我没有再抄 `%b`（正文体）。要看全量自己现跑 `git log -1 --format=%s cfacfd3c`（尺＝现跑，射程＝单枚 commit 标题）。

⇒ 教训同源于派单点名的那两条："标了逐字的块必须真逐字"＋"引文不许混枚"。这一处**在合入前自查到**，未污染后面四件。


---

## 6. 本件自己的字节数与提交

- 本件由 `wc -c` 现量后在回报里给数（不在本文件里自报，免得自我循环）。
- 第 1 笔 commit 只带这一枚 pathspec：`git add -- .scratch/wisp/probes/297/a1/00-anchor-and-rulers.md`。
- ⛔ 未 `git add -A`、⛔ 未碰 `design/`、`frontend/`、`cmd/wisp/` 等工作树里不属于我的脏文件。
