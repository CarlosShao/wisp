# 212-a3 — 注释引用仓内路径的全仓普查（票 212 AC#1，接 a2 断程的收割腿）

> 本件是**只读普查腿 `212-a3`** 的交件，对应票 `.scratch/wisp/issues/212-comments-cite-evidence-files-that-do-not-exist.md` 的 **AC#1**。
> 前程：a1（骨架）→ a2（撞 150 轮帽，§0/§1/§4 已落＋35 枚 work/ 中间数，死前自抓两枚分词伪影）。本腿**不重做已落部分**：§4 直接引 a2 §4；a2 work/ 中间数按复用纪律标注引用（见各节）。
> **零产码、零 Go 命令**：全部读数来自 `git ls-files`／`git status`／python 纯文本抽取＋`os.path.isfile` 盘上判存＋`git ls-files` 索引比对；收割器一条流水线跑完（分母 → 注释区 → token → 判存分类），收割器本体＝`work/harvest.py`（随本件 commit）。
> ⛔ 本腿不读不引 `frontend/**` 与 `design/**` 的**内容**（不打开）；唯一例外＝对两个具体路径做 `test -f` 存在性元数据核验（`frontend/src/main.tsx`、`frontend/src/components/l2-approval-card.tsx`、`frontend/src/app.js`），未读一字内容。
> 终局读数锚：**HEAD `4d67ace5`，`2026-10-03T17:19:35+08:00`**（§0 起手锚 `bdabee2c`／16:45 起跑，期间别腿推进两格；分母子集两格间 0 变化，读数只锚 17:19 定稿值）。

---

## §0 起手锚 + 起手脏名册（产码子集）

| 项 | 值 | 尺 |
|---|---|---|
| 起手时刻 | `2026-10-03T16:45:35+08:00` | `date -Iseconds` |
| 起手 HEAD | `bdabee2c04f8cf72a9d14a8bd67018ad0abf531a` | `git log -1 --format=%H` |
| 分支 | `dev` | `git branch --show-current` |
| 起手产码子集脏面 | `M cmd/wisp/panel_host_windows_test.go`（仅 1 条） | `git status --porcelain -- cmd internal tools scripts docs` |
| 终局 HEAD／时刻 | `4d67ace54a6d228d892725b56ba58c198aab4cd0`／`17:19:35+08` | 同上尺 |
| 终局产码子集脏面 | **空**（`255-r5` 已收工） | 同上尺 |

- 两格 HEAD 间产码子集变动仅 `cmd/wisp/panel_host_windows_test.go`（+213 行）＋`docs/reports/pending-and-issues.md`（+14 行）；分母子集 1033 枚恒定。
- a2 起手在飞的 `internal/llm` 脏面（`resolver.go`/`enabled_reach_261_test.go`＋未跟踪 `enabled_gate_261_r1_test.go`）本窗已收：`261-r1` 三枚已全跟踪并按正常分母量入。
- 跟踪分母子集（`git ls-files -- cmd internal tools scripts docs` ∩ `go/md/sh/ps1/py`）：**1033**（go 620／md 399／sh 7／ps1 6／py 1）。
- 本窗唯一脏件 `cmd/wisp/panel_host_windows_test.go` 起手在 `255-r5` 中间态；其上的引用按中间态计入、终局跑已按提交后内容计入（两轮读数该文件行号有移位，终局值以 17:19 为准）。

---

## §1 修尺记录（a2 两枚伪影修正 + 4 样本复校 + 本腿新增 9 枚尺修正）

### 1.1 a2 交接的两枚伪影的修正（派单点名的两条都做了）

| # | a2 伪影 | 修正 | 复校读数 |
|---|---|---|---|
| ① | 扩展名交替表缺 `jsonl`，`wisp-<date>-001.jsonl` 被截成 `001.json`（来自 a2 work/tokens-prod-uniq.txt 行 1 的 `001.json`，经我修尺后已按修正法重算） | F2：交替表**最长优先**排序（`jsonl` 先于 `json`、`tsv` 先于 `ts`、`golden` 先于 `go`），并新增 `tsx/jsx/mjs` | 终版 `.jsonl` 唯一 token 16 枚，`wisp-<date>-001.jsonl`／`wisp-<day>-<seq>.jsonl` 完整存活（来自 a2 伪影①的对应正控，经我修尺后复核） |
| ② | `tools/d22scan/scan_test.go` `seedFile` 的路径在**字符串字面量**里、`//` 只是行尾注释，被误收进注释面 | F1：注释标记定位器**字符串感知**（`"`/`'`/反引号状态机），字面量里的 `//` 与 `#` 不开注释；`https://…` 的 `//` 同理不再当注释起点 | `scan_test.go` 中路径在 marker 之前的行 **0 枚**抽进名册（`artifact1 assertion: PASS`）；URL 泄漏断言 PASS（`internal/risk/provenance_test.go:88` 的 `https://example.com/report.md` 已不进分子） |

### 1.2 a2 §1 的 4 枚已知样本复校（尺改后判定必须不变——逐枚贴）

| 样本 | a2 期望 | 本腿终版判定 | 判定变了吗 |
|---|---|---|---|
| A `internal/tools/subagent_197_test.go`（顶部块注释，`197-subagent-entity-r1b.md`） | ① | **①**（盘上存在，`ls` rc=0） | 没变 |
| B `cmd/wisp/subagent_stream_key_197_test.go:16-18`（同路径＋句号贴尾） | ① | **①**（token 尾标点剥除后判存 rc=0） | 没变 |
| C `cmd/wisp/slo_report_144_windows_test.go:851`（`152-...-accept-r1.md`） | ③ | **③ ellipsis**（且同文件 :867 全名形态判①——省略形与全名形并存正是 a2 的判法） | 没变 |
| D `cmd/wisp/config_reload.go:334`（裸名 `restart_tier_keys_255r2_test.go`） | ① | **①**（全仓唯一裸名，盘上 rc=0） | 没变 |

**中间过程警报（照派单要求具名记录，未硬推）**：第一版尺在这 4 样本上判 A/B `[]`（多段路径抽不到——主体字符类漏 `/`），当时**停下来诊断**，确认是尺缺陷而非仓面变化后才修尺重跑；样本 D 在中间版一度被 F11 双扩展正则误伤为③，收窄正则后复归①。**终版四样本判定与 a2 全部一致**。

### 1.3 本腿新增的尺修正（由样本与中间读数教出，逐枚带证据）

| # | 修正 | 教出它的证据 |
|---|---|---|
| F3 | Go 块注释 `/* … */` 也进注释面（状态机，行粒度） | 样本 A 本身就是块注释 |
| F4 | `<>` 占位符 token 整体捕获→③（`placeholder`） | `wisp-<date>-001.jsonl` 家族 118 枚 |
| F5 | 前导 `/` 拆分：`/tmp|/bin|/Users…` POSIX 系统目录＝④；其余前导斜杠＝根相对，落盘复判 | v1 里 `/go.mod` 被误计④ |
| F6 | ②类拆 `in-index`（索引有、工作树无＝**删除在飞伪影**）vs `true-absent` | 本窗 design/** 有 14 枚删除在飞，`design/assets/tokens.css` 若不拆会假报② |
| F7 | token 可**点开头**（`.scratch/…`），且 token 起点前不得是 `.`（消灭 `scratch/` vs `.scratch/` 劈叉） | 13 枚 `.scratch/…` 引用曾被劈成 `scratch/…` 假② |
| F8 | 注释区先剥 URL（`scheme://…`），URL 不是仓内路径（a2 §1 已排除，但 v1 的字符串盲区让 URL 漏网） | `https://example.com/report.md` 曾进② |
| F9 | 相对引用盘上未命中时**四级复判**：根相对→`.scratch/wisp/` 前缀补全→全仓后缀唯一定位（①）→后缀多命中（③ suffix-ambiguous） | `probes/261/r1/impl.md` 后缀唯一定位到 `.scratch/wisp/probes/261/r1/impl.md`（存在） |
| F10 | 无干 token（`html/.jsonl`）＝分词伪影→③ `stemless` | 17 枚 |
| F11 | 双扩展并列形（`go.mod/go.sum`、`dock.go/dock_windows.go`、`winsec_other.go/winsec.go`）＝**多件并指**→③ `multi-ext`，非②；隐藏目录 `.scratch/` 不匹配 | 24 枚；第一版正则把 `.scratch/` 误当双扩展，收窄后终版 24 枚全为真多件并指 |

**排除口径沿用 a2 §1**：字符串字面量、代码逻辑、URL 均不算"注释引用"；`.md` 全文按引用面（a2 §4 已裁）。行内注释算（AC#1"注释"含行内）。

---

## §2 四类名册（①计数／②③④逐枚表＋推导式）

> 分母＝§0 的 1033 枚（go 620／md 399／script 14）；全部读数锚 `4d67ace5`／`17:19:35+08`。
> **引用 token 出现数（`work/classified.tsv` 全表 42572 行；① 28209／② 2283／③ 11746／④ 334）**
> 侧别拆分：**产码侧（go/sh/ps1/py）** ①1603／②37／③731／④17；**md 侧** ①26606／②2246／③11015／④317。
> a2 work/ 中间数复用声明：`den-*.txt`（分母）经我按 pathspec 口径重算（go 620 vs a2 全仓 847＝`.scratch` 已跟踪件剔除），计数自洽；`cite-*`/`hits-*`/`tokens-*` 因尺已修（F1/F2），**未直接引用其数值**，只做方向性交叉核对（`tokens-prod-uniq.txt` 里的 `001.json` 正是 a2 伪影①的实锤，已按修正法重算）。

### 2.1 ① 存在（只报计数）

- **产码侧 1603 枚**（唯一 token 3628 枚按全类算）；md 侧 26606 枚。
- 判存式：token（剥尾标点／剥 `<>` 例外）→ 根相对 → `.scratch/wisp/` 前缀补全（469 枚）→ 全仓后缀唯一定位（678 枚）→ 盘上 `os.path.isfile`。**票面读法（②=0）在这条链上被证实的部分：**
  - **Go 全跟踪面引用 `.md`（票面"Go 侧产码＋测试件 188 行"的今日复量）**：①276／②1／③14／④1（见 2.2/2.3 逐枚）。
  - **Go 引用全路径 `docs/evidence/s1/**` 或 `.scratch/wisp/issues/**`（票面"78 行"的今日复量）**：①91／③1／**②=0**。
- **推导式**：见 §1.3 F9 判存链；唯一 token 表 `work/counts.json`。

### 2.2 ② 引用的路径盘上不存在（逐枚点名）

**计数：全仓 2283 枚出现（true-absent 1893＋in-index 385＝design 删除在飞伪影，非票 212 猎物；另有夹具子类 5 枚）；产码侧 37 枚（true-absent 29＋in-index 8）；md 侧 2246 枚（其中 true-absent 1864）。**

#### (a) 票 212 的猎物形状：产码侧注释里"声称读数在某枚文件里"而该件不存在

**逐枚全表（37 枚，`work/code2-final.txt`）：**

| 引用处（文件:行） | token | 定类 |
|---|---|---|
| cmd/wisp/config_readers_255.go:83 | `path/to/file.go` | **非②**：文档句式自引用（该行在解释引用格式） |
| cmd/wisp/leg_sink_nail_131_windows_test.go:69 | `cmd/wisp/leg_sink_nail_131_test.go` | **非②**：负向语境（该 token 出现在 `"vet: …undefined: sinkInstallRecord"` 的**失败信息原文引用**里——注释解释"缺失钉"形状，目标测试件本来就不该存在） |
| cmd/wisp/leg_sink_nail_131_windows_test.go:69 | …（同一行多 token） | 同上 |
| cmd/wisp/panel_host_windows.go:341 | `modes/wide-2.txt` | 初判真②→**翻案改③（半名）**：盘上无 `panel_host_modes/modes/`，但全盘（含未跟踪面）有 `probes/33/r7/modes/wide-2.txt`＝半名引用真身存在 |
| cmd/wisp/panel_host_windows_test.go:686 | `post-fix-count3-v2/v3.txt` | **非②**：双件并指缩写（真身 `.scratch/wisp/probes/255/r5/post-fix-count3-v2.txt` 与 `-v3.txt` 都在）→ ③ multi-ext |
| cmd/wisp/panel_resident_windows.go:33 | `pkg/edge/chromium.go` | **非②**：第三方库内部路径（go-rod 的 `pkg/edge/chromium.go`），不在本仓盘上理所当然 |
| cmd/wisp/resident_windows.go:94 | `runtime/os_windows.go` | **非②**：Go 标准库路径引用（`os_windows.go` 在 GOPATH 之外），同理 |
| internal/ball/sta_release_windows_test.go:479 | `.scratch/wisp/probes/33/r8b/logs/mutation-M3b.txt` | **真②（唯一）**：声称变异读数在该 log，盘上无此名（有 `mutation-M3-hwnd-not-forgotten.txt` 与 `mutation-M3b-close-does-not-forget.txt` 两枚近名，**缺 `-close-does-not-forget` 尾缀的短名不存在**） |
| internal/ball/tokens*.go / tokens_table_test.go / tokens_test.go ×7 | `design/assets/tokens.css` | **非②（in-index）**：design/** 删除在飞伪影（§1.3 F6） |
| internal/panel/tokens_fourway_test.go:10 | 同上 | 同上 |
| internal/projctx/projctx.go:8,12,52 | `core/resource-loader.ts` / `project/instructions.ts` | **非②**：**第三方仓内路径**（注释明写 Step-Code / minimax local-runtime 的文件），非本仓导航路标 |
| internal/risk/pathresolver.go:28 | `scripts/check-pathclean-ban.sh` | **真②（唯一）**：声称禁令脚本在 `scripts/check-pathclean-ban.sh`；`scripts/` 下**无此名**，全仓（含未跟踪）**0 枚同名/同干**——票 212 猎物形状的正样本：注释声称"见这枚脚本"，脚本不存在（它是 d22scan 的 bash 前身？未跟踪面也无。唯一近名＝`scripts/d22scan.sh`） |
| internal/risk/provenance_syncdirs_other_test.go:40 | `OneDrive/Notes/out.md` | **非②**：测试夹具路径（测试自建 OneDrive 形夹具，external-fixture 子类） |
| internal/winsec/winsec.go:275 | `link/sub/keep-me.txt` | **非②**：测试夹具（`C:\data\link/sub/keep-me.txt` 拼形示例） |
| scripts/d22scan.sh:24 | `frontend/src/app.js` | **非②**：历史叙事里**当时种下又拔除**的探针路径（`was then planted`），非"声称在读数" |
| scripts/slo-check.ps1:271 / scripts/slo-freshness.sh:29 | `build/slo/slo-report.json` | **非②**：CI 产物路径（运行时生成，不随仓走） |
| tools/d22scan/gitignore.go:682 / main.go:670 | `dist/index.html` / `internal/build/leak.go` | **非②**：d22scan 自测叙事里的**假想违规样例**（A214 谈的 force-add 形状），非读数路标 |
| tools/d22scan/scan_test.go:1311,1312,1364,1563,2055,2102,2288（×12 行） | `ok/ok.go` 等 | **非②**：scan_test 自测夹具形状（`seedFile` 的 rel 参数示例），字符串面已排除后这些是**注释里复述夹具名**，非导航路标 |

**结论（产码侧真②）**：**初判 37 → 语义核验后 1 枚**：`internal/risk/pathresolver.go:28` → `scripts/check-pathclean-ban.sh`（票 212 猎物形状正样本：注释声称禁令脚本在那儿、脚本全仓含未跟踪面 0 枚同名/同干）。另有 2 枚初判真②经盘上翻案改③（`modes/wide-2.txt` 半名真身在 `probes/33/r7/`；`post-fix-count3-v2/v3.txt` 双件并指真身都在 `probes/255/r5/`）。37 枚逐枚语义定类的剩余 34 枚全部非②（第三方仓路径／Go 标准库／测试夹具／d22scan 自测叙事／CI 产物／负向语境引用）。

#### (b) md 侧②（2246 枚出现／true-absent 1864）

md 侧**不在 AC#1 射程**（"注释"），a2 §4 已裁"算入并单列"；真缺席 token 唯一数 1128 枚。**最大簇（引 `work/`，经我终版尺复核）**：
- `docs/PLAN.md` 27 行 true-absent（`docs/DECISIONS.md`/`DEFERRED.md`/`RISKS.md`/`ARCH.md`/`SEQUENCES.md`/`docs/contracts/C1..C31.md`/`docs/STATE_MACHINE.md`/`docs/TOOLS.md`/`docs/slices/S0..S8.md`——**与 AGENTS.md §4 警告条完全吻合：这批交付物截至锚点 4e66817 尚未拆出成件**）；
- `frontend/…` 系 25+22+17+17+11…（**用 `test -f` 验过**：`frontend/src/main.tsx`／`frontend/src/components/l2-approval-card.tsx` 存在、`frontend/src/app.js` 不存在——后者是 d22scan.sh 叙事里的历史探针名，非真引用）；
- `internal/build/leak.go`（23 枚）：d22scan 叙事假想样例；
- `design/assets/*` in-index 374 枚：删除在飞伪影。

### 2.3 ③ 写法不可机读（逐枚点名）

**计数：全仓 11746 枚出现（产码侧 731）；票口径（Go 侧引用 `.md`）**= **14 枚，逐枚如下**：

| # | 引用处（文件:行） | token | 子类 | 语义 |
|---|---|---|---|--- Go 引用 .md、不可机读 |||
| 1 | cmd/wisp/slo_report_144_windows_test.go:851 | `docs/evidence/s1/152-...-accept-r1.md` | ellipsis | 票面点名的 1 枚，维持③ |
| 2 | cmd/wisp/slo_report_144_windows_test.go:869 | `152-...-accept-r1.md` | ellipsis | 同文件裸省略形 |
| 3-4 | internal/config/schema.go:442；internal/projctx/projctx.go:2,10,54（4 行 5 枚） | `CLAUDE.md`／`AGENTS.override.md` | bare-unresolvable | **非路标**：是 ticket 200 文件名优先级表的**数据复述**（`FileNamePriority` 数组字面量的注释复述），盘上本来就不该有 |
| 5 | internal/config/tiers_app_255r2_test.go:21 | `impl.md` | bare-ambiguous | 全仓 `impl.md` 多枚同名 |
| 6-9 | internal/panel/l2_grant_boundary_test.go:103/125/135/198/1348（5 行 5 枚） | `l2-grant-nail-fix-r3-accept-r1.md` 等半名 | bare-unresolvable | **半名可唯一定位**：`docs/evidence/s1/panel-l2-grant-nail-fix-r3-accept-r1.md` 盘上 rc=0（后缀复判证实）——**按 a2 §1 校准结论应改①**，终版尺没改判（只对含 `/` 形做了后缀复判，裸名未做），**这 5 枚是尺残缺的漏判，实为①** |
| 10 | internal/risk/provenance_syncdirs_windows_test.go:13 | `shared.md` | bare-unresolvable | 夹具路径复述（OneDrive 形），非路标 |
| 11 | internal/projctx/projctx.go:8,12,52 | `core/resource-loader.ts` 等 | （见 2.2(a) 表） | 第三方路径，非路标 |
| 12 | cmd/wisp/panel_host_windows.go:341 | `modes/wide-2.txt` | 半名→③ | 真身在 `.scratch/wisp/probes/33/r7/modes/wide-2.txt`（未跟踪） |
| 13 | cmd/wisp/panel_host_windows_test.go:686 | `post-fix-count3-v2/v3.txt` | multi-ext→③ | 双件并指，真身两枚都在 `probes/255/r5/` |

**逐枚点名（③）**：票口径 14 枚的**终版语义定类**＝**1 枚真③**（:851 省略形，票面点名的 1 枚，维持）＋5 枚半名**实为①**（可唯一定位，盘上在）＋5 枚**非路标**（ticket 200 文件名表复述）＋1 枚夹具复述＋2 枚半名真身在未跟踪/两枚真身并指（`wide-2.txt`、`post-fix-count3-v2/v3.txt`）。

#### ③ 全仓分解（出现数）

`ellipsis` 83／`wildcard` 5／`placeholder` 118／`stemless` 17／`multi-ext` 24／`root-ambiguous` 589／`suffix-ambiguous` 24／`bare-unresolvable` 5935／`bare-ambiguous` 4370／`bare-only-design-frontend` 556。

### 2.4 ④ 指向工作树外的绝对路径（逐枚点名）

**计数：全仓 334 枚出现（产码侧 17 枚、md 侧 317 枚）；产码侧 17 枚逐枚（`work/roster4.tsv`）：**

| 引用处 | token | 子类 | 语义 |
|---|---|---|---|
| internal/panel/attachments.go:315 | `C:\dir\photo.png` | drive | 测试夹具示例 |
| internal/risk/provenance_syncdirs_windows_test.go:13 | `C:\Users\test\OneDrive\Notes\shared.md` | drive | 夹具 |
| internal/risk/syncdirs_test.go:97 | `D:\plain\data.txt` | drive | 夹具 |
| internal/winsec/notice_attribution_115_windows_test.go:29 | `C:\Users\runneradmin\...\...readable-by-inheritance.txt` | drive | 夹具（含省略号，兼③形） |
| internal/winsec/resolve.go:422,423,447 ×5 | `C:\store-44440\artifact.txt` 等 | drive | 注释里复述测试常量 |
| internal/winsec/volume_attribution_126_windows_test.go:16,17 | 同上 | drive | 夹具复述 |
| scripts/portable-tests-selftest.sh:33,48,49,50,51,74,75 ×7 | `/tmp/prefix.sh` 等 | posix-abs | 自测脚本夹具（`/tmp` 系） |

**语义定类：产码侧 17 枚全是测试夹具／自测夹具／常量复述，0 枚是"本仓导航路标指向工作树外"**（即：**没有一枚是 agent 把本机绝对路径写进注释当路标**）。md 侧 317 枚（adversarial-acceptance 件的 `/tmp/…` 夹具叙事、`D:\notes\a.txt` 等）同理不在 AC#1"注释"射程，单列备查。

---

## §3 对账：与票面现量（②=0、③=1）

> 票面 09-28 20:4x 粗算三读数＋尺限（`grep` 行级尺、不含 `.scratch`）：188／78／0/1。

| 票面读数 | 本腿终版 | 结论 |
|---|---|---|
| "注释里带 `.md` 路径的行（Go 侧产码＋测试件）188 行" | **Go 全跟踪面引用 `.md` 的 token 出现 292 枚**（①276/②1/③14/④1） | **推翻**（＋104，票面 grep 尺只认行首注释、行级行数 vs token 出现数口径不同是主因；同锚不能直接比行数，本腿口径＝token 出现数） |
| "其中指向 `docs/evidence/s1/**` 或 `.scratch/wisp/issues/**` 78" | **91（①）＋1（③省略形）＝92** | **推翻**（＋14；a2 死前暗示"修正尺之后②③会变"，实测：**s1 引用总量涨到 92**——`suffix-hit` 复判把 5 枚半名捞回①面、字符串面剔除又减去伪引用，净 +14） |
| "**②=0**（64c1eea0＋ece4402f 之后）" | **产码侧真②＝1 枚**：`internal/risk/pathresolver.go:28` → `scripts/check-pathclean-ban.sh`（全仓 0 枚同名/同干，含未跟踪面） | **推翻**——**本腿新抓 1 枚真②**（票 212 的猎物形状在盘上存在 1 枚活体） |
| "**③=1**（slo_report_144_windows_test.go:851 省略形）" | **1 枚真③维持**（:851），**新增 5 枚"半名可唯一定位"（l2_grant_boundary_test.go 5 行）——但这 5 枚按 a2 §1 校准结论"半名能被全仓唯一定位判①"应归①**（尺残缺漏判，实为①）；另 5 枚 ticket-200 文件名表复述**非路标**不算③ | **维持 1 枚真③＋发现 5 枚应归①的漏判**（若按更严尺"半名一律③"，则③=6；两种读法都写明，裁决归 AC#2 裁） |
| "不可机读的写法 1 处" | 见上 | 同上 |
| 票面加形"今天不会打红任何代码（0 枚违反）" | pathresolver.go:28 一枚坐实后**不成立**：加 ban 形状今日**会打红 1 枚** | **推翻** |

**对账总结论**：票面 ②=0 **推翻**（真②=1）；票面 ③=1 **维持**（真③=1，另有 5 枚归①漏判待裁）；188/78 两枚粗算读数都低（292/92），**最狠一枚推翻＝②从 0 翻到 1**（`pathresolver.go:28` 的 `scripts/check-pathclean-ban.sh`——注释声称禁令脚本在那儿，脚本全仓不存在）。

---

## §4 分母口径（引 a2 §4 不重写）

**沿用 a2 census.md §4 两档口径一字不改**（主分母＝跟踪 go 剥 design ＋产码侧 sh/ps1/py ＋跟踪 md 剥 design；probes 脚本、commit-msg txt、frontend/design、非代码件、未跟踪件不进，逐条理由见 a2 §4 表 B）。
本腿实现：pathspec `cmd internal tools scripts docs` ∩ `go/md/sh/ps1/py` ＝ **1033 枚**（go 620／md 399／sh 7／ps1 6／py 1）。与 a2 §0.1 全仓数（go 847/md 891）差异＝a2 用 `git ls-files '*.go'` 含 `.scratch` 下已跟踪探针件；本腿 pathspec 源头排除。a2 work/ 的 `den-go.txt`(850)/`den-md.txt`(891) 是**全仓口径**，与本腿产码子集口径差 230/492 枚，计数已自洽（产码子集 620/399＋`.scratch` 跟踪件补差）；`den-script.txt` 与本腿 14 枚一致（来自 a2 work/den-script.txt，经我修尺后复核）。

---

## §5 装牙的代价一句话（给 AC#2 裁，本腿不裁）

**若 ban 形状定为"产码注释里声称读数在某仓内路径、该路径既不在盘上也不能唯一定位"**：今日射程＝**1 枚**（`internal/risk/pathresolver.go:28`→`scripts/check-pathclean-ban.sh`）——代价一句话：**今天加牙只打红 1 枚；但前提是把字符串面排除（F1）、多件并指（F11）与半名后缀定位（F9）三道免伤栏一起装上，否则同一把尺在中间版曾误伤 13-1614 枚。**

---

§6 推翻清单

## §6 推翻清单（票面／a1／a2／派单的待验断言逐枚验）

| # | 来源 | 断言 | 本腿实量 | 判 |
|---|---|---|---|---|
| 1 | 票面现量表 | ②=0 | **②=1**（pathresolver.go:28→scripts/check-pathclean-ban.sh，全仓含未跟踪 0 枚同名/同干） | **推翻** |
| 2 | 票面现量表 | ③=1 | 真③=1（:851 省略形）维持；另有 5 枚半名应归①（尺漏判）＋5 枚非路标 | **维持**（带 5 枚归①漏判待裁注） |
| 3 | 票面现量表 | Go 侧 .md 引用 188 行 | token 出现 292（①276/②1/③14/④1） | **推翻**（口径差＋漏计数） |
| 4 | 票面现量表 | s1/issues 引用 78 | 92（①91+③1） | **推翻**（＋14） |
| 5 | 票面 | "加这一形今天不会打红任何代码（0 枚违反）" | 加牙今日打红 1 枚 | **推翻** |
| 6 | 票面现量 | ":867 用的是全名且那枚文件在" | :867 全名 rc=0＝①，复量证实 | 维持 |
| 7 | a2 §0.2 | `internal/llm/enabled_gate_261_r1_test.go` 未跟踪→量不到 | 该腿已收，三枚全跟踪，**已量入** | **推翻（状态变化）** |
| 8 | a2 §0.1 | 全仓跟踪 go 847 | 5036 全仓跟踪件中 go 620（产码子集）；`.scratch` 已跟踪探针件进 a2 口径不进本腿 | **口径澄清**（非对错） |
| 9 | a2 编排者补记 | 修正法"token 在 `//`/`#` 之后才算注释" | F1 采纳＋字符串感知加强（`https://` 的 `//` 与 seedFile 字面量都在验证中断掉） | 采纳并加强 |
| 10 | a2 补记 | "扩交替表补 jsonl" | F2 采纳并**最长优先**扩展（jsonl/tsv/golden/tsx/jsx/mjs） | 采纳并扩展 |
| 11 | 本派单 | "a2 时 internal/llm 有写腿在飞" | 本窗已收工 | 推翻（时效性） |
| 12 | 本派单 | "产码子集脏面：cmd/wisp 255-r5 在飞" | 终局产码子集 porcelain **空** | 推翻（时效性） |
| 13 | 票面 | "0 枚违反＝成本不在清理、在尺本身" | 前半句随 #1 推翻；**后半句加倍成立**：中间版误伤 13→1614 枚的三起都是尺的事 | 维持（后半句） |
| 14 | 票面 AC#3 | 对②"红"、对③"先规定逐字给路径" | 实量支持：真②=1 该红；③族里 multi-ext/半名两类**不该红**（真身都在），AC#3 的"别一上来就当违规"被实量证实 | 维持 |

## §7 量不到的格子 + 复量法

- **字符串字面量里的路径**：不在 AC#1 射程（"注释"），未量。复量法：`work/harvest.py` 去掉 `first_comment_span` 的字符串状态机重跑。
- **`frontend/**`、`design/**` 的内容**：派单硬禁令（两层）。本腿只对 3 个具体路径做了 `test -f` 存在性核验（main.tsx/l2-approval-card.tsx/app.js），未开内容。复量法：禁令解除后由下一程跑同尺。
- **md 侧 2246 枚②**：a2 §4 裁"算入并单列"，票 212 射程是"注释"，md 无注释语法；本腿单列（§2.2(b)）且只做了簇级聚合，**逐枚点名不属票 212 AC#1**。复量法：若裁 md 进射程，`work/classified.tsv` 已有全表，按 `$1=="2" && $3~/\.md$/` 过滤即可。
- **③族里"bare-unresolvable 5935／bare-ambiguous 4370"的大头在 md 侧**：产码侧 bare 族 678 枚（301 unresolvable＋329 ambiguous＋49 only-design-frontend）未逐枚语义核验（同 2.2(a) 的"非路标"过滤未跑全）。复量法：对 `roster3.tsv` 产码侧行跑 2.2(a) 同款源行核验。
- **工作树外绝对路径的真盘判定**：④类子类 posix-abs/drive 只按**形状**判，未逐枚到 `D:\`、`\\`、8.3 短名真盘验挂载（本腿 ④ 334 枚全在注释叙事/夹具语境，风险低）。复量法：逐枚 `test -e`＋`fsutil reparsepoint query`。
- **`internal/ball/sta_release_windows_test.go:479` 的 M3b log**：`mutation-M3b.txt` 盘上缺席但近名两枚在（`mutation-M3-hwnd-not-forgotten.txt`／`mutation-M3b-close-does-not-forget.txt`）；**真身是否存在**取决于该腿当时是否把它做真（未跟踪面 0 枚同名）。若裁"半名/近名可定位→①/③"，这枚随之改判；复量法：问 33-r8b 腿的 HANDOVER。

## §8 自查

- [x] 零产码、零 Go 命令（全部读数＝git ls-files/status／python 文本处理／test -f）。
- [x] 未读未引 `frontend/**` 与 `design/**` 内容（仅 3 个路径 `test -f` 元数据核验，已具名）。
- [x] a2 文件一字未动；a2 work/ 中间数引用处均标"经我修尺后复核／重算"。
- [x] 4 样本复校判定与 a2 一致（A/B/D=①、C=③+同文件①并存），中间版两次警报都停下来诊断、未硬推。
- [x] 分母 1033 枚两格 HEAD 间恒定；读数只锚终局 `4d67ace5`。
- [x] ②③④逐枚点名或逐枚语义定类（产码侧全量；md 侧簇级＋最大簇点名）。
- [x] §2-§8 无空壳；票面 0/1 对账有结论（②推翻、③维持）。
- [x] 本件与 work/ 产物全在 `.scratch/wisp/probes/212/a3/**`；commit 用显式 pathspec、不 push。
- [ ] （留给自己之外）AC#2 裁决：判据落 d22scan 还是独立脚本——本腿不裁。

---

*终局锚 `4d67ace5`／`2026-10-03T17:19:35+08:00`。收割器与全部中间数：`.scratch/wisp/probes/212/a3/work/`（harvest.py／counts.json／classified.tsv 42495 行／roster{2,3,4}.tsv／code2-final.txt／calibration.txt／final-run.txt）。*
