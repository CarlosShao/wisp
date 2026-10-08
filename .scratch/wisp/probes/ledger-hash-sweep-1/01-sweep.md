# ledger-hash-sweep-1 · 01-sweep（A720–A736 机械对账）

- 区段：`docs/reports/pending-and-issues.md` 行 `14093`–`14393`（A720 节头→文件尾，共 301 行；A736 为末节）。
- 起手锚：HEAD `5f52f310`（10-08 19:01），腿起时刻 19:01:45；复核窗口内 HEAD 有前移（如 `97cb748a` 10-08 19:03「票 259 收口」——影响见表三）。
- 尺：hex＝`grep -oE '\b[0-9a-f]{7,40}\b'` 66 枚（联合大小写 `[0-9A-Fa-f]` 复扫同 66，差集 ∅＝无遗漏）；路径＝反引号包住且带 `\.go|\.md|\.sh|\.py|\.json|\.yml` 的去重 52 枚。
- 全程 rc 必检；无 `2>/dev/null`；零 Go 命令；零改盘（写面仅本目录）。

## 表一：commit／hex 号（66 枚：存在 57／不存在 9）

### 存在 57 枚
| 号 | 类别 | 现量 |
|---|---|---|
| `05992d05` | commit | 10-08 12:05 242-v2 对抗验收 |
| `082bc7ba` | commit | 10-08 16:21 comment-truth-2 起手锚 |
| `0af6594d` | commit | 10-08 18:44 259-v1 起手锚 |
| `0c9726f9` | commit | 10-08 11:50 111-ciif1 起手锚 |
| `1309757b` | commit | 10-06 14:03 probe(111/r5) 骨架 |
| `14115e8d` | commit | 10-08 18:27 comment-fix-prep-1 起手锚 |
| `16901acb` | commit | 10-05 11:36 票 181 AC#7 生产侧落地(181-r3) |
| `1ae79210` | commit | 10-08 18:33 parking-2 停车点刷新 |
| `2405be70` | commit | 10-08 18:52 card-proof-prep-1 起手锚 |
| `351ba362` | commit | 10-08 18:20 253-v1 起手锚 |
| `351e5a5e` | commit | 10-08 09:01 票 111 r6 AC#11 |
| `3a343bc7` | commit | 10-08 09:37 票 35 35-r5 AC#8 |
| `4240c0b0` | commit | 10-08 17:17 ticket-181-status-1 就地打旧 |
| `4263812a` | commit | 10-08 18:32 parking-2 起手锚 |
| `440dd88765`→`440dd887` | commit | 09-21 18:57 docs(95-done,107→rejected,109) |
| `4db510ce` | **blob** | 种前/还原 hash（git log 静默 rc=0，cat-file -t=blob） |
| `50b34971` | commit | 10-08 18:57 gate-snapshot-1 起手锚 |
| `58e2b155` | **blob** | 种前 hash（同上判法） |
| `5e8748b3` | commit | 10-03 21:07 212-r1＋258-r1 代笔收尾 |
| `6410a062` | commit | 10-08 18:40 comment-fix-prep-1 交料 |
| `6538b3c1` | commit | 10-08 18:30 ledger-audit-1 起手锚 |
| `6547fd30` | commit | 10-08 14:41 A718 落账 |
| `65f4c968` | commit | 10-08 15:48 253-r1 第 1 笔 |
| `73b2438c` | commit | 10-08 18:32 253-v1 对抗验收件 |
| `762b694e` | commit | 10-08 17:12 comment-truth-2 交件 |
| `7a452d0a` | commit | 10-08 16:10 A721 落账 |
| `89443d0a` | commit | 10-08 16:40 ticket-181-status-1 起手锚 |
| `8ace2272` | commit | 10-08 18:59 gate-snapshot-1 门禁快照 |
| `8d865932` | commit | 10-08 17:28 253-r1 第 2 笔 |
| `8e085ad4` | commit | 10-08 15:56 ci-if-eval-1 收档 |
| `99e14860` | commit | 10-08 15:52 stale-claim-1 交件代提 |
| `9e8477ed` | commit | 10-08 18:26 167-a2b 两格读数 |
| `a04a095f` | commit | 10-04 19:53 265-r1c 收编 |
| `a1f0d2ff` | commit | 10-08 18:42 comment-fix-check-1 起手锚 |
| `a818df46` | commit | 09-30 12:05 test(197-r1) 正控 |
| `a966b564` | **blob** | 票 253 工作树件 hash（同 HEAD:票面 blob） |
| `b4602f13` | commit | 10-08 18:42 181-v3 证据件更正 |
| `bb1d5ef8` | commit | 10-08 18:42 246-raisercensus-2 起手锚 |
| `bf185665` | commit | 10-08 18:36 181-v3 起手锚 |
| `c9820fb0` | commit | 10-08 16:55 ticket-181-status-1 线索复跑件 |
| `cc31526165`→`cc315261` | commit | 10-06 10:58 ## A634 |
| `cc31526165734e612de297848bb2080bd459ccba` | **commit（40 位）** | = `cc31526165` 全号 |
| `ccf16baa` | commit | 10-08 18:53 259-v1 终裁 |
| `d33d492d` | commit | 10-08 17:58 167-a2 死腿遗产代提 |
| `d568ad8d` | commit | 10-08 18:37 ledger-audit-1 复核件 |
| `d892ecba` | commit | 10-08 18:38 ledger-audit-1 时点窗核对 |
| `e060cef6` | commit | 10-08 18:58 card-proof-prep-1 01-recipe |
| `e6dc79ed` | commit | 10-08 12:06 111-ciif1 ci.yml 补挂 23 枚 |
| `e6dc79ed5bd14380cce1507f94886d65a51ee63f` | **commit（40 位）** | = `e6dc79ed` 全号 |
| `e7cd0c6e` | commit | 10-08 18:52 comment-fix-check-1 交料 |
| `eb2a217a` | commit | 10-08 18:50 246-raisercensus-2 交件 |
| `ec78addf` | commit | 10-08 18:34 parking-2 读数补记 |
| `f6b79ab0` | commit | 10-07 09:23 ci(111 r5) lint 名册 |
| `f8810238` | commit | 10-08 09:29 ci.yml 注释更正 |
| `f9e7584a` | commit | 10-08 17:38 246-raisercensus 起手锚 |
| `fa86e2d5` | commit | 10-08 18:24 167-a2b 起手锚 |
| `fe1a9ce6` | commit | 10-08 18:42 181-v3 证据件 |

### 不存在 9 枚（rc=128，`fatal: Not a valid object name`／`ambiguous argument`）
| 号 | 类别 | 说明 |
|---|---|---|
| `76b694e` | **真异常**（7 位 commit 短号，少一位） | 节内 A722 行 82 已自纠：真身＝`762b694e`（存在）「它回报里写的是`76b694e`…`git show` 现量＝fatal” |
| `35591482293` | 非 commit 引用（11 位纯十进制 id） | ci 文件 :523 引用 |
| `37021179942`／`37166458550`／`37396530365`／`37405698188`／`37406757402`／`37703959747`／`37736935814` | 同上（采样 id 族，含 n296/n300/n301/n303/n305 标注） | 节内「采样 6 发完整 id」族 |

⚠ 转述差异具名：转述说“40 位如 `58e2b155…` 是磁盘 blob”——盘上 `58e2b155` 是 8 位短号且 =blob（属实），但**两枚真 40 位（`cc31…ccba`／`e6dc…e63f`）cat-file -t 均＝commit**，非 blob。

## 表二：文件路径（反引号取件 52 枚 = 带目录 28 ＋ 裸名 20 ＋ 片段 4）

### 带目录 28 枚：28/28 存在
- 根字面即中 20 枚：`cmd/wisp/{models,run,panel_host_windows,panel_dispatch_binding_roster_253r1_windows_test,panel_inbound_guards_35r3_test,panel_transport_35r2_test}.go`、`cmd/wisp/testdata/esclistener/main.go`、`cmd/wisp/subagent_selfapproval_197_test.go`（porcelain 形，含 `M ` 前缀）、`.scratch/wisp/probes/comment-fix-prep-1/01-ready-to-apply.md`、`internal/{agent/tools,panel/bridge,panel/composer_dispatch,risk/provenance,tools/bridge}.go`、`internal/tools/bridge.go`（另见“1280 internal/tools/bridge.go”输出行形）、`.github/workflows/ci.yml`（两个命令形 `git diff --stat…`／`git show e6dc79ed --…` 内嵌）、`scripts/slo-freshness.sh`、`scripts/d22scan.sh`（`sh …` 形）、`tools/signmodels/main.go`。
- 仅 `.scratch/wisp/` 前缀下存在 6 枚（根字面 rc=128）：`issues/README.md`、`probes/181/v3/evidence.md`、`probes/246/raisercensus2/01-census.md`、`probes/253/r1/01-anchor-correction.md`、`probes/259/v1/01-evidence.md`、`probes/card-proof-prep-1/01-recipe.md`（节内简写）。
- HEAD: 形 2 枚 `HEAD:cmd/wisp/{resident_task_source_windows,resident_windows}.go` 直接 rc=0。

### 裸名 20 枚（HEAD 后缀匹配枚数）
`00-anchor.md`=39、`01-check.md`=1、`01-draft.md`=1、`02-crash-self-rescue.md`=1、`20-mutations.md`=1、`30-blind-spots.md`=1、`PLAN.md`=1、`bridge.go`=**9**、`ci.yml`=1、`config_receipt_255_test.go`=1、`doc.go`=20、`loader.go`=1、`pending_read.go`=1、`resident_windows.go`=3、`slo-fresh.yml`=1、`ticket259_denial_rulers_test.go`=1、`ticket259_panel_capability_rulers_test.go`=1。
- **HEAD=0 但工作树在 3 枚（未跟踪 `??`、非 ignore rc=1）**：`02-commit2-readings-correction.md`（`.scratch/wisp/probes/253/r1/`，该目录 HEAD 已入库 7 枚、独它未入库）、`90-msg-00.md`／`91-msg-01.md`（`.scratch/wisp/probes/comment-truth-2/`，HEAD 已入库 3 枚、此两枚未入库）。
- ⚠ 转述差异具名：转述“本仓 bridge.go 有两枚”——HEAD 现量 **9 枚**（产码 3：`internal/panel`、`internal/tools`、`internal/models`；余 6 为探针副本）。以盘上为准。

### 片段 4 枚（非路径形状，0 匹配不构成缺失）
`*_test.go`、`.md`、`_test.go`、`:(exclude)*_test.go`（git pathspec 魔字）。

## 表三：锚抽验

### ① 票面尺三处（现跑：`wc -l`＋`grep -cE '^[[:space:]]*- \[ \]'`＋`\[x\]`；工作树=HEAD 均 clean）
| 票 | 台账句（节） | 台账读数 | 我现量 | 判 |
|---|---|---|---|---|
| 181 | A731「94→101 行／未勾 2→1／已勾 5→6」 | 101/1/6 | 101/1/6 | **对上** |
| 246 | A723「97 行／勾 8／未勾 1」 | 97/8/1 | **105**/8/1 | 勾格对上；**行数对不上**（票面 10-08 18:51 由 `f13e7c37`「246 票面追加编排者裁定」＋8 行＝A732 流程，读数时点后漂移） |
| 259 | A734「59→69 行／未勾 5→1／已勾 1→5」 | 69/1/5 | **77/0/6** | **对不上**：被其后 `97cb748a`（10-08 19:03「票 259 收口：AC#4 翻勾」）取代；另票面 §12 自报「69→**79** 行」vs 复量 77（`wc -l`＝`awk NR`＝`grep -c ''` 三者一致，clean）⇒ 自报 79 疑误计（差 2） |

### ② HEAD 行号锚两处（现取）
- `cmd/wisp/resident_windows.go`：`:249` 装「…This call is the caller: it reads」✓、`:248` 装「…AskOnTaskRoot / askConfirmation still had zero product」✓ —— **对上**（A722 那枚更正成立）。
- `.github/workflows/ci.yml`：`:521`＝`- uses: actions/checkout@v4`，属 `test-windows:` 作业（:516 起）且该步无 `if:` ✓；`:469`＝`if: ${{ !cancelled() }}`（守卫本身），属 `test-core:` 作业（:402 起）的「Portable package tests」步（name :447、run :470）✓ —— **对上**。
- （bonus）票 181 `:79`＝「## 6. AC#7 现状核对…」✓。

## 附：三表汇总
- 表一：66 枚；存在 57（commit 54：52 短号＋2 枚 40 位；blob 3）／不存在 9（`76b694e`＋8 枚 11 位十进制 id）。
- 表二：52 枚；带目录 28/28 存在（6 枚仅 `.scratch/wisp/` 前缀下）；裸名 17/20 有匹配，3 枚 HEAD=0 但工作树在（未跟踪）；片段 4 枚非路径；不存在的字面路径 **0 枚**。
- 表三：票 181 对上；246 行数对不上（原因已具名）；259 对不上（被 19:03 收口取代＋§12 自报 79 vs 77 差 2）；两处 HEAD 行号锚均对上。
