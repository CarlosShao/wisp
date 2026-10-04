# 票 256 `256-v1` 非实现者验收表 —— 攻 `256-r1` 那四枚提交

**腿**＝`256-v1`（验收腿，未参与 256 任何一程的写作）
**被验对象**＝`256-r1` 四枚：`5ad8ec27`（骨架 10:10:10）→ `2a10134e`（落地 10:20:13）→ `2708289f`（证据件写满 10:26:16）→ `826d7a54`（收尾 10:27:18）
**票面**＝`.scratch/wisp/issues/256-resident-gate-built-before-session-grants.md`
**它的证据件**＝`.scratch/wisp/probes/256/r1/evidence.md`（我 `wc -l -c` 复量＝**465 行／44,337 字节**，14:58 现跑，与它 §9 自量逐字相等）⛔ **其中每一句对本腿都是待验断言**
**本腿起手 HEAD**＝`afbeb2c7`；交件时 HEAD＝`ed892e6e`（并发腿在推进，与本腿无关，见 §4）
**本腿全部读数取于 2026-10-04 14:0x–14:5x +0800，件同目录 `probes/256/v1/**`**

一句话结论：**代码这一半是真的、判据真的在问门（我用两形突变各自证了）；但这发落地顺带改掉两项常驻腿的行为（L1 窗口 3s→2s、L2 超时不再恒 300s 且无上界），而这两项没有任何一枚上游件裁过——本腿判 AC#1／AC#2 停勾、AC#3 可翻，并把那两枚读数按甲／乙／不做摆给编排者转机主（§8）。**

---

## §1 名册与现量复认（全部本腿自己重拉）

### 1.1 逐 sha 名册（尺＝`git show --name-only --format= <sha>`，14:1x 现跑）

| sha | 触及面 |
|---|---|
| `5ad8ec27` | `.scratch/wisp/probes/256/r1/evidence.md`（1 枚） |
| `2a10134e` | `cmd/wisp/resident_approval_risk_256_windows_test.go`（新，552/0）＋ `cmd/wisp/resident_approval_windows.go`（140/11）＋ `cmd/wisp/resident_windows.go`（7/1） |
| `2708289f` | `.scratch/wisp/probes/256/r1/evidence.md`（1 枚） |
| `826d7a54` | `.scratch/wisp/probes/256/r1/evidence.md`（1 枚） |

⇒ **probes 之外名册＝三枚，与编排者那份逐名一致**（新测试件／`resident_approval_windows.go`／`resident_windows.go`），**没有第四枚、没有漏计**。

### 1.2 票面现量复认

- 框数：`grep -o '\- \[[ x]\]'` 现跑＝**3 枚未勾／1 枚已勾**（已勾那枚＝AC#0），与编排者那份现量一致。
- 它 §2.2 那句"产码唯一调用点"我现读复认＝`cmd/wisp/resident_windows.go:132` 逐字 `ra := newResidentApprovalWithConfig(rt.Layout.DataDir)`（行号与它写的 `:132` 相等，这枚没漂）。
- ⚠ **它 §2.1 那五行 `file:line` 今天全漂了**，但不是它的错：`resident_approval_windows.go` 在它之后被 **`260-r3`／`260-r4`**（`170e0459`／`79c579e2`）在同一枚文件里加了约 129 行。本腿现读对照＝它写 `:117-119` 的委托函数现在在 **`:170`**、它写 `:160` 的新函数在 **`:213`**、它写 `:176-181` 的 `slog.Info` 在 **`:229`**、它写 `:204-218` 的三枚 provenance 常量在 **`:280` 起**、它写 `:226-241` 的 `residentRiskGateValues` 在 **`:297`**。⇒ **引它那些行号前先重拉**（票面对"现量"的这条通用要求对它自己也成立）。
- 它 §9 那把占位尺我复跑同一形状（三个圆括号词面＋一枚英文大写词的 alternation，字面抄自票面 §7-1；本行刻意不重复其字面，否则这一行会把自己的扫描命中）＝**0 命中**；另记一句：它全文里另有 **4 枚"方括号三字待验标记"**（＝它自己声明"这一格我没跑"的写法），那种标记不在前两枚词面的射程里 ⇒ **它的 0 命中成立，但"这一格没实测"在它文里是 4 处、不是 0 处**（本表 §6 用同一口径，但改成具名散文，不用那种标记）。
- 它 §2.3／§7-4／§7-5 三处 `grep` 我全部复跑：`grep -rn "只作用于跑任务的进程\|常驻腿今天用常量" cmd internal tools`（14:58）＝**0 命中**（ⓑ 那句至今没落到任何产码）；`internal/config/unwired.go:121-122` 现读仍逐字只写 `cmd/wisp/run.go` ⇒ 复认成立；`grep -cn '"risk' cmd/wisp/config_readers_255.go`＝**0** ⇒ `hotRowClaims` 今天没有 `risk` 这一格，复认成立。
- 它 §7-4 那句"今天不会因此判红"我从读码升成实跑：`go test -count=1 -run 'TestEveryLockedSectionKeyIsAccountedFor' ./internal/config/`（14:59:21）＝**`ok ... 0.028s`** ⇒ 那两行少说一个消费者确实不判红。

---

## §2 ★那一格：schema 的 2s 与门自己的 3s 不是同一枚默认值

### ① 改期望是对的还是把缺陷洗成常态？——**判：是对的（但不是"只改期望"那么小）**

逐字读那枚用例（`cmd/wisp/resident_approval_risk_256_windows_test.go:198-208` 的 case 定义 ＋ `:253-255` 的断言本体）：

- 断言本体＝`if got := ra.gate.Window(); got != c.wantWin { t.Errorf("Window() = %v, want %v for seed %q; what this reading proves: %s", ...) }`，其中 `wantWin: 2 * time.Second`、case 名逐字 `"timeout only - the window then comes from the schema tag, not from the gate's compiled constant"`、`winIsWhat` 逐字点名来源（`internal/config/schema.go` 的 `default:"2"`）并写着 "Before ticket 256 this process always held 3s"。
- ⇒ **它断的是"等于 schema 默认"，不是"随配置变"**——而对一枚**没被种值**的键，"随配置变"这句根本无法构造（没有可比的两发）。所以这不是一枚被换弱的断言，是一枚**换了问法**的断言：从"门自己的默认是多少"改成"这份可读文件进门之后门到底拿着多少"。
- 本腿不复述它的读数，我自己跑了一发独立尺（件 `probe-gateface.txt` 14:08:42／`probe-gateface-2.txt` 14:47:32）：
  - `P1 file readable, no [risk] section at all` ⇒ **`gate_window=2s`**、`provenance="config"`；
  - `P2 dataDir 存在但没有 config.toml` ⇒ **`gate_window=3s`**、`provenance="defaults (config.toml unreadable)"`；
  - `P3 [risk] 只写了 confirm_timeout_sec = 45` ⇒ **`gate_window=2s`** ＋ `gate_timeout=45s`；
  - `P4 用户自己写 l1_window_sec = 3` ⇒ `gate_window=3s`。
  ⇒ **它 §2.3 那枚"没被任何上游件量过"的读数成立，且它第一版 3s 的期望是被产码判错的、不是它想绿的**。
- 我还用一发突变证明这枚 case 不是恒绿（`mut-V2-window-hardcoded.txt` 14:10:31，把 `residentRiskGateValues` 的窗口那一行改成字面 `2 * time.Second`，即"文件里的窗口值根本不流进门"）：**五枚 case 里恰有一枚红**（`l1_window_sec = 99` 那枚，`Window() = 2s, want 3s`）⇒ 期望值本身没有把断言变成词面复读。
- **但两件事它没做，本腿点名**：(a) 期望写成**字面量 `2 * time.Second`**，而不是 `config.NewDefaults().Risk.L1WindowSec`——schema tag 哪天真被改动，这枚钉会红在 `cmd/wisp`、却把原因写在失败文案的第 3 行里，追因成本留给别人（这是口径建议，不是缺陷判决）；(b) 这枚 case 现在**同时**是"未来任何裁定想让常驻腿守 3s"的障碍件，而它身上**没有任何机器可读的"本读数待裁定"标记**，只有 `winIsWhat` 里那句散文（"Whatever the 255/265 wording for the ⓑ sentence ends up saying, it has to say 2s here and not 3s"）。⇒ **没洗白，但把一件待裁的事钉成了已实现的常态**，这一层要靠 §8 那枚拍板补，不该由本腿或它自己补。

### ② 这次落地有没有"顺带改变一项默认值"？——**判：有两枚，都具名，且都不在任何上游裁定里**

**第一枚（它自己量到并上报的那枚）**：常驻腿 L1 执行前阻止窗口，在"config.toml 读得到"这一形里 **3s → 2s**。
钳位带是 `[MinL1Window=2s, MaxL1Window=3s]`（`queue.go:120`／`:122`），2s 落在带内 ⇒ 钳位不动它（本腿 `gate.go:143-151` 现读复认）。

**第二枚（它没量、本腿现跑出来的）**：`confirm_timeout_sec` 现在**没有上界地**进了常驻腿。凭据（`probe-gateface*.txt`）：

| 种子 | 门实际拿着 | 读数 |
|---|---|---|
| `confirm_timeout_sec = 99999` | `gate_timeout=27h46m39s` | 越界值**读到的是原值**，不是钳位 |
| `confirm_timeout_sec = 10` | `gate_timeout=10s`、`warning_lead=30s` ⇒ **`c18_warning_scheduled=false`** | C18 那句"超时前 30s 醒目提示"从此不再排期 |
| `confirm_timeout_sec = 45` | `gate_timeout=45s`、`c18_warning_scheduled=true` | 30s 是那道坎 |

机制逐处现读：`internal/config/validate.go` 里 `timeout`/`window` **零命中**（配置文件对这两枚键**没有任何范围校验**）；`queue.go:84-87` 只管 `timeout <= 0 → 300s`；`queue.go:88-90` 把 `warnBefore <= 0 || warnBefore >= timeout` 兜成 `DefaultApprovalWarning`，而 `gate.go:528` 是 `if lead := g.q.Timeout() - g.q.WarningLead(); lead > 0` ⇒ **超时被配到 ≤30s 时那句提示静默消失**。256-r1 之前常驻腿恒 300s ⇒ 对常驻腿而言这份免疫力是**本腿交掉的**（会话腿 `run.go:615-616` 一直如此，属既有形状，不是它新造）。

**"改契约"这一问的判定**：**不是改契约文本，但也不是"修 bug 不需批准"那一档。**
- 不是改契约：`[risk]` 的 schema 默认 **2** 就是契约里写着的数（`docs/PLAN.md:2738` 与 `docs/specs/SPEC-03-config-secrets-envs.md:34` 逐字 `l1_window_sec(int)=2`），L1 窗口的允许带在 `PLAN.md:128`／`:284`／`:1174` 与 `SPEC-06 §2` 都写作 **2–3 秒** ⇒ 2s 合法，`tools/d22scan` 与 `internal/observe/thresholds.go` 里没有 L1 窗口阈值（尺＝`grep -n "l1\|L1 window" internal/observe/thresholds.go`＝0 命中）。
- 但它**推翻了四份上游件共同写下的前提**：`docs/evidence/s1/246-resident-task-source-r2.md:85`、`246-resident-task-source-v2.md:212`、`248-settings-write-path-r1.md:305`、台账 `docs/reports/pending-and-issues.md:10107` 都逐字断言"**窗口那一项即便接了也钳在 3s ⇒ 只有超时是真差异**"。这句话现在**是错的**（接了之后是 2s，不是钳在 3s）。而 `A534` 裁的 ⓑ 那句回执文案（"常驻腿今天用常量 300s／3s"）**正是长在这句错前提上**的。
- ⇒ **判：需要一枚具名批准／`Q##`**。批准的对象**不是** schema（产码里 `internal/config/**` 一字未动），而是两件事：**常驻腿要不要吃 `[risk]`／schema 默认（含 2s）**，以及**ⓑ 那句文案照实改写**。它自己没有替谁裁（§7-7 明写"交给编排者的读数，不是本腿自行修的文案"），这一点它守住了；**它没守住的是"未定义即停"**——撞上"顺带改掉两项既有行为、且推翻四份裁过的前提"这一形时它继续交付了，按 `AGENTS.md §0.1` 与票面 §排程那枚"未定义即停"，这一格该停下来问。（它有权继续做ⓐ那半，因为 `[risk]` 两枚可派是 `§8-1` 明裁的；被推翻的是那条裁定**赖以成立的理由**。）

### 用户看得见什么差别（不写成"只是内部默认值"）

1. **球上那行 L1 倒计时少 1 秒**。窗口是"执行前阻止"（`gate.go:247-251`：窗口走完无否决＝执行），界面上是 `PLAN.md:3488` 那枚 20px 环形倒计时＋"单击球 / Esc 取消"提示 ⇒ **用户可反应的时间从 3s 变 2s**，且是**默认装机态就变**：任何跑过一次 `wisp run` 或被设置页写过的宿主，`config.toml` 都存在（`run.go:245` 调 `ensureFirstRunConfig` → `SaveFile(cfgPath, config.NewDefaults())`，`[risk]` 会被逐枚写进文件）；常驻腿自己不建那份文件（`firstrun_198_test.go:281` 正是这么钉的："the resident leg reaches the assembly root directly, and first-run creation shared there writes a config.toml its pins forbid"）。⇒ **"文件可读但没写 `[risk]`"不是边缘态；真实宿主上常驻腿基本恒 2s，3s 只剩"从未跑过 run 且没有文件／文件读不到"那一支。**
2. **同一份配置能让球上那张 L2 卡挂几小时**（第一枚"超时一律判拒绝"的 C18 语义在常驻腿被配置覆盖），**或让它不再给"快超时"那句提示**（第二枚，≤30s 时静默）。两者 256-r1 之前对常驻腿都不可能。
3. **没有任何一句界面文案说这些变了**：ⓑ 那句话今天 0 命中（§1.2 我复跑），`residentStatusLine()` 它一个字没改（本腿复读：状态句没出现 `config`/`defaults` 字样，provenance 只进 `slog`）。⇒ 用户看到的还是同一句话，脚下的数已经换了一轮。
4. ⚠ **方向读数记我一笔，且具名两枚都不是本腿射程**：`internal/config/schema.go:451-453` 的注释与 `manager.go:442-445` 把"调高 `l1_window_sec`"记作 **loosening**，而 `gate.go:247-251` 的真语义是"窗口走完才执行"⇒ 按 `gate.go`，调高＝更晚执行＝**更紧**。这两处对方向的读法**相反**（既有形状，归 `internal/config`／D45 那一族）。但机主要拍的那件事在两种读法下都成立：**数字动了，且没人说。**

---

## §3 AC#1／AC#2／AC#3 逐格判语（每格只用"成立／不成立／量不到"）

### AC#1（`Options` 补 `Grants` ＋ `[risk]` 两枚）——**整格：不成立；`[risk]` 那一半：成立；`Grants` 那一半：不成立**

- `Grants` 那一半**今天确实没做**，且不是我读它注释读出来的：本腿现读 `approval.Options` 字面量（`resident_approval_windows.go:219-225`）实传字段＝`UI`/`Channels`/`Window`/`ApprovalTimeout`/`Logf`，**没有 `Grants`**；两形突变里也没有任何一枚能造出它（V1/V2 都只动 `[risk]` 那两枚的值）。
- 那枚"在位的钉断的是 Grants 不在场"＝`resident_approval_risk_256_windows_test.go:452-456`（`if got["Grants"] { t.Errorf(...) }`）。⚠ **具名一句：那枚钉不是既有资产，是 `2a10134e` 这一发自己新增的**（尺＝`grep -rn "Grants" cmd/wisp --include=*_test.go` 只有这一枚文件里的 AST 断言是负向钉，其余全是 `ListGrantsBySession` 那类账本读面）。⇒ 它的方向是对的（自我封锁越界），但**不能拿它当"别人早已钉住这件事"的凭据**；真正在位的外部钉是票 265 §现量.2 那条读码（`run.go:618` 带 `Grants`、`resident_approval_windows.go` 不带）。
- **票面有没有把"只做到哪一半"说到明处？——半有半没有。**明处在 §8-1／§8-5（编排者自己的裁定）与它 §7-1（"本票 AC#1 整格不许由本腿翻勾"）；**AC 那一行本身（票面 `:15`）仍逐字写着"补上 `Grants` ＋ 来自 `[risk]` 的两项"，勾旁的凭据没写半字**。⇒ 翻勾时**不要**把这格当"完成了一半可勾一半"，要么另立 AC#1a/AC#1b（编排者的写面），要么停勾（本腿建议见 §7）。
- `[risk]` 那一半：成立。凭据不是它的自述，是本腿 §2① 那六发独立读数（P1–P4）＋ §4 两形突变的红名集合。

### AC#2（反向判据＋正控）——**正控那一半：成立；`GRANT-DROPPED` 禁现那一半：不成立（今天不可钉）；"是不是真在问门"：是，但两枚尺例外**

本腿按编排者定形的唯一有效正控亲手种了两发（件 `mut-*`／`probe-gateface*.txt`），并用**两形独立突变**问同一件事："用例读的是门，还是配置解析结果？"

| 突变 | 形状 | 红名（本腿现跑） |
|---|---|---|
| **V1** | 配置照读、回执照写、日志照打，只把 `Options` 里那两枚值换成字面 `0`（＝"接了名字但值不流进门"） | ①`config.toml missing`（它的**配对正控**先红）、② 五枚 case 红 **4** 枚、③ 红（setup fail-fast）；**④⑤ 全绿** |
| **V2** | 窗口那一行换成字面 `2 * time.Second`（＝"文件里的窗口值永远不进门"，超时仍是真读） | **只有 ② 的 `l1_window_sec = 99` 那枚红**；①③④⑤ 全绿 |

⇒ 判语：
1. **①②③ 真在问门**（V1 之下它们红，且 ① 红的是它自己的配对正控那半，兜底那半照旧绿——形状与它 §4-M2 的自述一致）。**这条我复现了它的核心主张。**
2. **④⑤ 不背值的账**：`Window: 0` 这种"名字在、值是零"的形④全绿（它数的是字段名集合，`Window` 这个 key 还在）。它 §3④ 的字面主张（"缺 Window/ApprovalTimeout 就红"）**没夸大**，但**"落点自证"这个词会被读成值也自证了**——本腿建议把它的口径写成"形状尺／名册尺"。⚠ 另记：它自己的 M2 突变是**删掉那两行**，比本腿 V1 的"塞零"更强，所以**它没有暴露 ④ 对值不敏感这一面**；这是它证据件的覆盖缺口，不是假话。
3. **五枚 case 没有任何一枚同时看住 V1 与 V2**：`99 → MaxL1Window` 那枚在 V1 下绿（3s 恰是兜底值），`1 → MinL1Window` 与两枚 2s case 在 V2 下绿。⇒ 两枚合起来才看住窗口那一侧。这一格它写对了方向（"Window() 只能当钳位反控"），但**它没写出"哪枚反控在什么形状下会恒绿"**——本腿补在这里。
4. **`GRANT-DROPPED` 禁现那一半：不成立**，且理由与它 §5 一致（`g.grants` 仍为 nil ⇒ 票面 §7-6 的第一形今天照样出）。本腿没有把它读成通过，也没有去动 `ticket224_reply_grant_test.go:283`（既有仪器，⛔ 不改）。它归口票 265 的写法成立（票 265 已立，13:2x，账 A599/A598 那族）。
5. 另补一发本腿自种的越界读数，回答"种一发越界值 ⇒ 必须读到钳位而非原值"：**窗口这一枚读到了钳位**（`l1_window_sec = 99` ⇒ 3s；`= 1` ⇒ 2s；`= 0`/`= -1` ⇒ 3s，⚠ 这一枚语义上"配 0"读回**带内最大值**，是 `gate.go:143-145` 的既有形状，不是本腿能改的）；**超时这一枚读不到钳位、读到原值**（`= 99999` ⇒ 27h46m39s）。⇒ 两枚字段的"越界"命运不同，这条区别**票面与它 §8-3 都没写过**。

### AC#3（越界检查）——**成立**（凭据全在本腿自己手里）

- 禁区逐名对：四枚 commit 的并集＝**4 个路径**（一枚 md ＋ 三枚 Go），`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／任何 golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` **零命中**；三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）不在名册里。
- **D43 转移表我自己验，没接受它的自述**：`git show --name-only --format= <sha> -- docs/PLAN.md internal/statemachine` 四枚**全为空**；再加一把 blob 级尺＝`internal/statemachine/machine.go` 在 `2a10134e^` 与 HEAD 的 md5 都是 `10531f6330c7f92b95bc1e9a4c0d2bfa` ⇒ **一字未动**。
- **不新增状态**：名册里没有 `internal/statemachine/**`、没有 `cmd/wisp/models.go`（D43 那枚产码 `statemachine.New` 的唯一 `cmd/wisp` 调用者）⇒ **无从新增状态**，本腿复读到这一层为止；"有没有语义上偷偷多一枚态"属读码判断，本腿读到的是它只加了三枚字段＋两枚函数，没碰状态机。
- 工作树里今天躺着的 `design/**` 删除、`tools/d22scan/**`、`internal/agent/approval/pending_read.go` 等改动：本腿起手（14:0x `git status --porcelain`）就看到它们在飞，**逐条与它 §1.1③ 那份对照过，全是别家腿的**，且**不在它那四枚 commit 里**（尺＝§1.1 的 per-sha 名册）。

---

## §4 突变与还原证明（本腿自己的台件）

**起手副本（任何突变之前抽的，尺＝`git cat-file blob HEAD:<path> > <副本>`，14:05，当时 HEAD＝`afbeb2c7`）**
件＝`.scratch/wisp/probes/256/v1/pre-mut-resident_approval_windows.go.pristine`／`pre-mut-resident_windows.go.pristine`。

| 文件 | 起手副本 md5 | 起手时工作树 md5 | 终值（14:48:23 现跑） |  matches_start |
|---|---|---|---|---|
| `cmd/wisp/resident_approval_windows.go` | `7a26c7a990dbd2351bdf9898b5bdc192` | 同左 | **`7a26c7a990dbd2351bdf9898b5bdc192`** | **是** |
| `cmd/wisp/resident_windows.go` | `a5ad69a4f8300a14ee7fd596b5c3a2fb` | 同左 | **`a5ad69a4f8300a14ee7fd596b5c3a2fb`** | **是**（本腿**没有对它做任何突变**——V1/V2 都落在另一枚文件；起手值＝终值这条只是复认，不是还原动作） |

每发现场读数：
- V1：突变后 md5＝`c49759e2dc107889347996868618924e`（跑前当场打印，件 `mut-V1-gate-not-fed.txt` 第 2 行，14:09:22）；跑完**先比当前 md5＝c49759e2**（确认没有并发腿在这两分钟里写同一枚文件）⇒ 才用起手副本还原 ⇒ 还原后 `7a26c7a9`。
- V2：突变后 md5＝`687fae56a1613397f7c89c2cf6686b80`（件 `mut-V2-window-hardcoded.txt`，14:10:31）；还原后 `7a26c7a9`。
- 还原动作**只有一条**：`cp <起手副本> <工作树文件>`。⛔ 没有"finally 里重读当前文件再写回"那一形。
- **`git status --porcelain -- cmd internal` 终值＝空**（14:5x 复跑，本腿的台件 `zz256_v1_probe_windows_test.go` 已从 `cmd/wisp/` **移入** `probes/256/v1/…go.bench`＝只建不删、不在产码面）。⚠ 14:48:23 那一发曾打印出 ` M internal/ball/hotkey_borrow_refused_260r5_test.go`——**那是 `260-r5` 在飞的工作**（其后果然由它自己提交成 `629f146b`／`ed892e6e`），与本腿无关；本腿没有碰 `internal/ball/**`（见 §6-7）。
- 它的还原凭据我也独立验了一枚更硬的：`git cat-file blob 2a10134e:cmd/wisp/resident_approval_windows.go | md5sum` ＝ `6b33f0df963e0eed3a2e35a29aedcb15`，**与它 §4 报的那枚逐字相等，也与它留在 `probes/256/r1/` 里的那份 pristine 相等**（本腿现跑，14:1x）⇒ 它那两发突变在**已提交树上**确实还原干净了。（本腿自己那枚 `7a26c7a9` 与它不同，是因为 `260-r3`/`260-r4` 后来在同一枚文件里加了行，与本腿的还原无关——这条差异已在 §1.2 具名解释。）

---

## §5 门禁读数（本腿现跑，逐条带时刻；⛔ 整包 `./...` 没跑、`slo-check.ps1` 没跑、真窗没开）

所有 `go test ./cmd/wisp` 都带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`。**这枚坑本腿自己复现了**：裸跑 `go test -count=1 -run 'TestTicket256Resident' ./cmd/wisp`（14:58:38）＝`exit status 0xc0000135` ＋ `FAIL github.com/CarlosShao/wisp/cmd/wisp 0.033s`、**`--- FAIL` 行数＝0** ⇒ 用例根本没跑，那不是红。包路径形状也用 `./cmd/wisp`（裸 `cmd/wisp` 会被当 import path）。

| 尺 | 时刻 | 读数 | 件 |
|---|---|---|---|
| 定向 `-run TestTicket256` | 14:02:44–14:02:48 | `ok cmd/wisp 0.126s`，五枚判据全 PASS | `v1-baseline-256.txt` |
| 定向（它的收尾同形过滤 `TestAC246\|TestTicket224\|TestTicket255Roster\|TestTicket256Resident`） | 14:06:52–14:07:27 | `ok ... 31.382s`，**top-level PASS 25／FAIL 0／SKIP 0**，rc=0 | `t-cmdwisp-broadfilter.txt` |
| **红名集合逐名比它** | 14:4x | 它那发的 25 枚 PASS 名 vs 本腿这发的 25 枚 ⇒ `diff` **空输出**＝**名集逐名相同**；它的 `final-targeted.txt` 里 `--- FAIL\|--- SKIP` 计数＝**0** | `their-pass-names.txt`／`my-pass-names.txt` |
| **整包 `./cmd/wisp/`（本腿多做的一发，超出它的取数面）** | 14:16:42–14:20:45 | `ok github.com/CarlosShao/wisp/cmd/wisp 240.587s`，**rc=0，无一行 `--- FAIL`/`--- SKIP`/`panic`** ⇒ 它这发落地在 `cmd/wisp` 全包层面无连带红 | `t-cmdwisp-full.txt` |
| `./internal/agent/approval/` | 14:06:33–14:06:35 ＋ 14:06:49 | `ok ... 0.390s`；`-v` 复跑：**top-level `=== RUN` 78 枚、`--- FAIL` 0、`--- SKIP` 1**＝`TestDefaultDeadlineWallClockMeasurement`（受 `WISP_84_MEASURE` 门 ⇒ **SKIP 不当通过读**） | `t-approval-pkg.txt` |
| `sh scripts/d22scan.sh` | 14:05:46（突变台件在场**之前**）与 14:54:48（还原**之后**） | 两发都 `rc=0`、末句 `d22scan: clean - no D22 ban violations`；ban #8 覆盖 `cmd/` **100** 枚 Go 文件（含注释与 `_test.go`）、`internal/` 503→**504** 枚（多的那枚＝`260-r5`，与本腿无关）；自带 `runtests.sh: OK - packages=[./...] PASS=35 FAIL=0 SKIP=0, === RUN=77` ⇒ 与它 §6 那行**逐字同形**（那是 `tools/d22scan` 自己的读数，不是 `cmd/wisp` 的） | `gate-d22scan.txt`／`gate-d22scan-final.txt` |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 14:06:15 与 14:55:08 | 两发 `rc=0`、`VERDICT GREEN`；`tracked paths=5586 over-budget=57 covered=57 not in roster=0`；unit cross-check＝字节视图与字符视图选出**同样 57 枚**（⚠ 它 §6 引的是 `tracked paths=5310`，本腿这发是 5586——差的 276 枚是 10-04 上午到下午之间别人落的工单/证据件，不是它的读数错） | `gate-pathlen.txt`／`gate-pathlen-final.txt` |
| 还原后的收口一发 | 14:55:50–14:56:22 | 同形过滤 `ok ... 26.760s` rc=0；单跑 `-run TestTicket256Resident` ⇒ `--- PASS` 计数 **5** | `t-final-restored.txt` |
| `go vet -tags winlive ./cmd/wisp/` | 14:58:40–14:58:41 | `rc=0`、空输出 ⇒ 复认它 §1.1 那发（**只证明编得动，不证明跑得绿**） | 本行即读数 |
| `gofumpt -l` | 14:32:00 | 点名三枚＝**`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/approval/queue.go`** ⇒ 按任务书这三枚是**票 212／258 的账户**，本腿**一枚没修**，只具名上报；它那三枚文件（含新测试件）**gofumpt 干净** | 本行即读数 |
| `go test -run TestEveryLockedSectionKeyIsAccountedFor ./internal/config/` | 14:59:21 | `ok ... 0.028s`（为把 §1.2 那一格从读码升成实跑，只点这一枚用例，未跑整包） | 本行即读数 |

**卫生门红名集合比对结论**：两把门本腿各跑两发（突变前／还原后），**四发全绿、红名集合恒为空集**，与它 §6 那两发同形；唯一逐名不同处＝pathlen 的分母（5310 → 5586），归因＝时间流逝与他人落件，⛔ 不属它这发。

---

## §6 判不动／量不到（逐条具名，⛔ 没有一条用"应该没问题"填）

1. **真起常驻进程（双击／Explorer 那一支）里门到底吃到没吃到：量不到（本腿未实测）**。判据 ⑤ 是 AST 尺，它证明的是"boot 调了新函数并交了 `*.DataDir`"，不是"那一发跑起来后 `config.toml` 真在那个目录、真读到了"。真窗／真热键那批（`winlive`）本腿一律**未实测**，机主那个词还没给。
2. **ⓑ 那句文案落在哪、说成什么数：判不动**。归口是 `255-r2`／票 265（`§8-6` 已定），本腿只复跑 0 命中（§1.2）。
3. **`confirm_timeout_sec` 无上界这一格该不该由本票修：判不动**。修它的正确写面要么在 `internal/config/validate.go`（⛔ 不在 256 写面，且 `[risk]` 是 locked 段），要么在 `approval`（⛔ C18 数字，动它＝改契约）。本腿只把它记成"本腿交掉了常驻腿那枚免疫力"的读数并交给 §8。
4. **C18 与 `[risk].confirm_timeout_sec` 谁优先：判不动**。`PLAN.md:1368`/`:3210` 把 C18 写成"超时 = 300s，一律判拒绝"的**冻结契约**，而 `SPEC-03:34`/`PLAN.md:2738` 又把同一枚数做成用户可改的 locked 键。这一对矛盾不是 256 引入的（会话腿一直如此），但**本腿的 P8/P12 让它第一次在常驻腿上也成立**。本腿不自裁，摆给编排者（§8 末条）。
5. **"配 `l1_window_sec = 0` 读回 3s（带内最大）"这一格有没有人钉：量不到**。既有尺里本腿没找到钉"0/负值 ⇒ 读到最大窗口"的用例（尺＝§3 第 5 条那批读数＋`grep -rn "MinL1Window\|MaxL1Window" cmd internal --include=*_test.go`，命中只有 256 那枚新件与 approval 包内文件）；这一格归谁，是编排者的账，不是本腿能定的。
6. **窗口那一侧的"不随改而动"（构造期定值）它没单独钉**，本腿也**没有替它钉**：`Window()` 的可达值只有 2s/3s 两枚，钉"改了不跟"在窗口这一侧读出恒绿的概率高于读出事实（它 §7-9 的理由，本腿 V1/V2 两发读数与之相容）。
7. **`internal/ball` 那一包一律未跑、未碰、未归因**。260-r5 正在那包里新增测试（`629f146b`/`ed892e6e`），那包出的任何动静本腿按**时序读数**处理，不记到 256 头上。
8. **`winlive` 那三枚的运行时读数：未实测**（本腿只复认了"编得动"，§5）。
9. **schema tag 若被改动，那枚 2s 期望会红在哪、算谁的账：判不动**（本腿不改 `internal/config`，也没做假想突变去验那条链路）。
10. **`git status` 与 HEAD 并发漂移**：交件时 HEAD 已从起手 `afbeb2c7` 推进到 `ed892e6e`，`resident_approval_windows.go` 中途被 `260-r3`/`260-r4` 改过（§1.2）。⇒ 本腿所有 file:line 都是**现读号**，引用者请重拉；本腿没和任何别的会话联系，也没查任何别的会话。

---

## §7 给编排者的翻勾建议（⛔ 本腿一枚框都没碰）

| 框 | 建议 | 条件／理由 |
|---|---|---|
| **AC#0** | 已勾，本腿无异议 | 它是只读普查那一程，本腿这轮验的是落地四发 |
| **AC#1** | **必须停勾** | 票面那一行逐字要求 `Grants` ＋ `[risk]` 两枚；`Grants` 那半不成立（§3），且今天已有新家＝票 265。若要把"半格交付"记进史，**正确动作是另立 AC#1a/AC#1b 或改写那一行**（编排者的写面），⛔ 不是把 AC#1 翻成完成 |
| **AC#2** | **必须停勾**；若要留痕，翻成一枚新框"AC#2a＝正控（`Queue().Timeout()`／`Window()` 读门）"可成立 | 正控半成立（本腿两形突变亲验），禁现半不可钉。**带条件**：① 把 ④⑤ 在票面/证据里的口径从"落点自证"降为"形状尺"（V1 之下它们全绿）；② 补一枚"值真流进门"的钉（本腿 V1 那形就是缺口的现成凭据）；③ 记下"99 那枚反控在 V1 下恒绿、1 与两枚 2s 那三枚在 V2 下恒绿"这条敏感度事实 |
| **AC#3** | **可翻** | 凭据全在本腿自己手里（per-sha 名册、`docs/PLAN.md`＋`internal/statemachine` 零触及、`machine.go` blob md5 与 `2a10134e^` 相等）；D43 一字未动、未新增状态 |
| 附带 | 无论翻不翻，§2②那两枚（2s 与"超时不再恒 300s／≤30s 时 C18 提示消失"）各该落一枚台账 `A##`，并把 `A534` 那句 ⓑ 标成"前提已被 256-r1 证伪，待重写"；`A591` 里"窗口接了也钳在 3s"那句同族措辞同样待改 | 出处＝§2② 那四份被推翻的上游件 |

---

## §8 待人拍板（三栏：甲／乙／不做＋每栏"没防住的形状"；⛔ 不许把"不批甲"读成"顺延乙"；⛔ 本腿没有直接找机主）

**要拍的那件事**：常驻（球＋托盘）那条腿的审批门，从今天起吃 `[risk]`／schema 的值——**L1 执行前阻止窗口在默认装机态从 3s 变 2s**，**L2 卡片超时不再恒 300s（可被配到几小时，也可被配到 ≤30s 从而丢掉"超时前醒目提示"）**。这两项都不在任何已裁文本里（`§8-1`/`§8-3` 反而写着"窗口接了也钳在 3s ⇒ 只有超时是真差异"，本腿已证伪）。

| 栏 | 内容 | 成本／要落哪 | **没防住的形状** |
|---|---|---|---|
| **甲：追认现状**（常驻腿吃 `[risk]`，含 schema 默认 2s） | 三件配套：① ⓑ 那句改写成"文件读得到 ⇒ 窗口取 `[risk].l1_window_sec`（未写时 2s）；读不到 ⇒ 常量 3s；超时同理，未写时 300s"；② 给 `confirm_timeout_sec` 补上界校验（或让 `WarningLead` 随超时缩短，别让 `timeout ≤ 30s` 静默丢掉 C18 提示）；③ 把 §3-AC#2 那枚"值真流进门"的钉补上（V1 那一形） | ①＝`255-r2`／票 265 写面；②＝`internal/config`（locked 段校验，⚠ **不是**改默认值）；③＝`cmd/wisp` 一枚新判据 | **没人改过配置的用户也被静默从 3s 改到 2s**：这一格甲只补"事后说清楚"，不补"事前有人同意"；且 `unwired.go:121-122` 那两行在甲之下仍少说一个消费者，直到有人去改 |
| **乙：只在文件真写了 `l1_window_sec` 时才吃，否则守门自己的 3s** | 需要"某个键有没有被写"的读取面。现读：`config.LoadFile` 只返回值与 `*Resolved`（`loader.go:41`，`Resolved` 里只有密钥），**没有 presence 出口** ⇒ 乙**是一枚新接缝，不是补字段**（与票 265 给 `Grants` 定的性是同一形，代价要照那枚票的口径量） | `internal/config` 加读取面 ＋ `cmd/wisp` 改取值；两枚都不在 256 写面 | **两条腿从此不一致**：会话腿 2s、常驻腿 3s（`run.go:615` 从不 presence 判断），而 ⓑ 那句要写两套话；并且**乙完全不解决超时那一枚**（P8 27h、P12 丢提示照样成立）⇒ 选乙必须另外处理 §2② 第二枚 |
| **不做：把 `[risk]` 那一半也退回，只留 AC#0 的普查** | 已交付的五枚判据＋`newResidentApprovalWithConfig` 全部拆回；票 265 的接口点（`§现量.2` 现在按 `:219` 读那枚不带 `Grants` 的门）要跟着改写 | `cmd/wisp` 反向一发 ＋ 台账更正两处 | **票面 §现量.3② 与票 248 AC#10 那条"改这两项配置对常驻腿不生效"继续是真的**，用户端等于"设置页里有两格球不看"；本腿读数说明这件事今天**只剩超时那一枚为真**（窗口那枚已经不等于是） |

**另附第二枚待拍（不属 256 射程，但本腿的读数让它第一次在常驻腿上也成立）**：C18 冻结文本"超时 = 300s，一律判拒绝"与 `[risk].confirm_timeout_sec`（可配、零校验、无上界）**谁优先**。要么在 C 表旁具名声明"300s 是默认值不是契约常量"，要么给那枚键加带界并让 `WarningLead` 跟着。**本腿两样都没动。**

---

## §9 记本腿自己写错的尺（全数留下，不抹）

1. **区间尺用错**：我第一次拉被验面用了 `git diff --name-only 5ad8ec27^..826d7a54` ⇒ 数出 **12 枚文件**，并把 `docs/reports/pending-and-issues.md` 与 `design/**` 误读成"禁区命中"。真相＝那四枚 commit 在历史上**不与并发腿相邻**，区间把 `260-r2`／`264-a1`／台账 commit 一起吞了。正解＝**逐 sha `git show --name-only --format=`**（§1.1 那张表）。⇒ 记：**共享工作树里验"谁碰了什么"永远逐 sha，不用区间。**
2. **把两发放反**：`0.390s` 我第一眼看成"整包 `cmd/wisp` 的读数"，实际它是 `internal/agent/approval` 那发；整包 `cmd/wisp` 是 `240.587s`（尺＝`t-cmdwisp-full.txt` 第 2 行 `cat -A` 复读）。⇒ 记：引数前回到件本身，不回自己的记忆。
3. **台件编译自查失败两处**：`zz256_v1_probe_windows_test.go` 第一版引用了不存在的 `approval_Min256()` helper、并 `import "time"` 未用 ⇒ 两处在跑之前修掉（件已移入 `probes/256/v1/…go.bench`，形状＝坏过的那版没留下，只记这件事本身）。
4. **还原证明差点少一发比对**：14:48:23 那发的 `git status --porcelain -- cmd internal` 打印出 ` M internal/ball/hotkey_borrow_refused_260r5_test.go`，我差点把它记成本腿未还原的痕迹（它属 `260-r5`，随后由它自己提交）。⇒ 记：**状态计数要现跑现抄并把每一行归到名上**；本腿最终那发才是空输出（§4）。
5. **一次编辑竞态自伤**：V2 还原后我用 `Edit` 改产码，工具回了"file changed since your last read"（因为我用 `cp` 还原过、文件 mtime 变了）。我没有靠那次提示继续瞎改，而是先 `git diff` 确认工作树相对 HEAD 只差我要的那一行，才跑 V2（凭据＝`mut-V2-window-hardcoded.txt` 上方那 12 行 diff 输出）。⇒ 记：**突变台件的形状校验要在跑之前做，不是事后。**
6. **一把尺的射程我说过头过一次**：§1.2 里 `grep -rn ... | head` 之后我差点把 `GREP_RC=0` 当成 grep 的退出码（那其实是管道末端的）。⇒ 本表里凡"0 命中"我只写自己直接看到空输出的那几发（ⓑ 那句、`hotRowClaims`、`thresholds.go`）。
7. **交件发上盘时 git 给了 18 条 `LF will be replaced by CRLF` 提示**（尺＝本腿 `git add` 那一条命令的 stderr）——**包括我那两份 pristine 台件**。⇒ 后果具名：这台机器带 autocrlf，**谁按 §4 去复跑 md5 时，比的必须是本腿交件时工作树里的那份，不是未来一次 `git checkout` 之后的那份**；未来复验 pristine 的 md5，请改用 `git cat-file blob <sha>:<path> | md5sum`（本腿 §4 验它那一发用的正是这把尺，不受行尾转换影响）。本腿交件后复跑：两枚产码文件与两份 pristine 的 md5 **仍逐字等于 §4 表里那两行**，`git status --porcelain -- cmd internal` 与 `-- .scratch/wisp/probes/256/v1` 都为空 ⇒ 这一次没被改写，但这条风险得记在案。

---

## §10 本件自量（成稿判据取盘上，不取工具回执；⛔ 占位词面 0 命中）

| 尺 | 读数 |
|---|---|
| 占位词面扫描（形状同 §1.2 那把，字面不在本行重复） | **0 命中**（`grep -cE` 打印 0；第一次跑出 1 命中，命中的就是我自己抄尺那行，已就地改写成不重复字面的写法——记进 §9 之前它先被本腿自己抓到） |
| 方括号待验标记枚数 | **0**（本表一律用具名散文写"未实测／量不到"，见 §6） |
| `wc -l` | **223** |
| `wc -c` | **38,187 字节** |
| 节次 | §1–§10 全填（§1 名册复认／§2 那一格三问＋用户看得见什么／§3 三格判语／§4 突变与还原／§5 门禁读数带时刻／§6 判不动十条／§7 翻勾建议／§8 甲乙不做三栏／§9 记本腿写错的尺／§10 自量） |
| 台件名册 | `probes/256/v1/` 下 19 枚：本件 ＋ 两枚起手 pristine ＋ 两发突变 ＋ 两发探针 ＋ 一发红名比对（两枚 txt 名集）＋ 六发门禁/测试读数 ＋ 一枚基线 ＋ 本腿台件 `zz256_v1_probe_windows_test.go.bench`（从 `cmd/wisp/` 移入，只建不删） |
| Git | 三发 commit 全带显式 pathspec（骨架发＝1 枚 md；本发＝同 1 枚 md），⛔ 无 `add -A`／`.`，⛔ 无 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`，⛔ 未 push，⛔ 未碰 `.scratch/wisp/issues/**` 任何一枚框，⛔ 未碰台账 `docs/reports/**` |

