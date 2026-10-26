# lint-r1 — `scripts/ticket-lint.sh` v1（只报不拦）落地读数

**腿**：`lint-r1`（写腿／实现者）。**日期**：2026-10-26。**件名**：`docs/evidence/s1/lint-r1-impl.md`。
**规格来源**：`docs/audit/2026-10-25-harness-reform-plan.md` §3.2 第 94–103 行七条 ＋ `--self-test` 正控。
**写面**：只新建 `scripts/ticket-lint.sh`（675 行）与本件。⛔ 未改任何票面、⛔ 未删任何存量、⛔ 未进 CI、⛔ 未动 `scripts/check-path-length-budget.sh` 一字节。

## 一、七条读数（真票池 365 张，2026-10-26 现跑）

| 条 | 读法 | 我的数 | 审计的数 | 差 |
|---|---|---|---|---|
| 1 红 | `-done` ＋ 行首未勾框 | **3**（4 格） | 3（附录 A 旗标①） | 票号与格数**逐枚相同** |
| 2 红 | `-done` 而全部 Status 行无 done/结案字样 | **66** | 69（旗标②，69/189） | −3（−4.3%） |
| 3a 红 | 活票无任何 Status 行 | **20** | 20（旗标④） | 票号**20/20 逐枚相同** |
| 3b 红 | 活票 Status 取值在词表外 | **12** | 无对应旗标（审计只做定性「词表 20+ 种写法」） | — |
| 3c 记 | 旧写法经 alias 表归一（不是红） | **128** | 同上 | — |
| 4 黄 | 双文件同票号 / 票号断档与 `voided/` 对账 | **0 重号 ＋ 14 缺号** | 无重号；缺 14 号（§4 末） | 14/14 **同名同号** |
| 5 红 | 票名 > 帽（100，读 262 门） | **58**（262 名册覆盖 58、新犯 **0**） | 58（§4 末，最长 164） | 相等 |
| 6 黄 | `in-progress`>72h／`review`>14d／`blocked-on-owner`>21d | **15** | 无对应旗标（审计只有定性「老化票 7 张」） | 审计 7 张中 6 张被点名 |
| 7 统计 | done／活／void／等 owner／立即可派 | done 189／活 175／void 1／等 owner 1／立即可派 78 | 189／175／—／~42（深读口径）／~65（深读口径） | 见 §六 |

**票号清单（逐枚）**

- **1**：`07`(1 格) `115`(2 格) `92`(1 格) — 与审计「07/92/115 共 4 格」逐字对上。
- **2（66）**：20 82 111 115 116 117 118 119 120 124 126 127 129 130 131 132 133 137 141 142 143 144 146 147 151 161 178 183 198 212 221 231 235 241 242 243 250 251 252 253 254 255 257 259 260 263 265 267 272 287 288 295 299 317 318 319 322 325 326 331 332 338 341 343 349 353
- **3a（20）**：223 228 256 262 264 266 273 274 276 279 284 301 302 303 304 305 310 312 313 314 — **与审计旗标④的 20 枚完全同集**。
- **3b（12）**：21 `segment-1-landed-decision-layer`／169 `owner／那一支写着`／177 `改于`／201 `部分结案`／219 `接单＝做`／222 `产码已交并静态核过`／226 `终裁已下`／247 `未结案`／290 `AC#0/AC#1/…`／291 `AC#1`／293 `本票的最新态`／296 `AC#3／AC#4`。
- **4 缺号（14）**：199 202 203 204 205 206 207 208 209 210 217 218 240 374；`voided/` 只存 `154-first-draft-void-premise-falsified.md`（154 在池内），故 **14 枚全部「voided/ 解释不了」** —— 机器对账走了，答案是没有一份归档能解释。
- **5（58）**：126 127 128 129 134 135 136 137 139 140 141 142 143 144 145 147 148 150 151 152 153 155 156 158 169 171 172 173 174 175 178 179 180 181 182 184 185 186 191 192 193 194 195 196 197 235 237 243 245 246 250 251 252 255 258 261 272 275 —— 全部带「grandfathered by the 262 roster」标记，**新犯 0**。
- **6（15）**：in-progress 103(35d) 107(35d) 108(35d) 109(5d) 128(32d) 311(7d) 70(35d) 71(35d)；review 12(36d) 62(26d) 64(35d) 68(36d) 72(15d) 112(35d)；blocked-on-owner 65(36d)。审计 §5 的「老化 62/64/65/68/70/72/77」里 **6 张被点名**，`77` 不在内——它的值是 `blocked`，第 6 条只认 `blocked-on-owner`（口径差，非漏判）。

## 二、与审计 3／69／20 的对照（差在哪儿，为什么）

- **① 3 ＝ 3**：连「07 一格／92 一格／115 两格」的格数都对上。审计 §附录 A 注里那条「331 被机械计数为 1，实为正文引用」我的行首锚定同样不误报。
- **② 66 vs 69（−4.3%）**：同一分母（189 张 `-done`）。差的 3 张是「Status 行里 done 只出现在括号正文、不在取值里」的边界票，口径可调，量级不动。
- **④ 20 ＝ 20**：**逐票同集**（223 228 256 262 264 266 273 274 276 279 284 301 302 303 304 305 310 312 313 314）。
- **③＝25 是审计自己的「s1 无证据件」，不是规则 3**（审计 §4 六旗标 ①残余未勾／②Status 从未写 done／③s1 无证据件／④无状态行／⑤台账零提及／⑥无出处标记）。规则 3 的可比读数是旗标④＝20，本尺 20/20 命中。
- **规则 3 的红口径**：3a（无状态行，20）＋3b（词表外，12）＝**32**，比审计的 20 多出 3b 那 12 张「有 Status 行但取值是散文/AC 号/文件路径」。这不是尺松紧问题，是**规则 3 原文是两问**（"无 Status 行 / 词表外取值"），而审计只查了前一半。**该报警的确实该报**：`291 [AC#1]`、`293 [本票的最新态]` 这种值，任何调度器都读不出状态。
- **3c（128）不是红**：`ready-for-agent`(60) `已立`(19) `待派`(13) `open`(7) `立而不派`(5) `未开工`(4) `在派`(2) `OPEN`(2) `在派队列里`(1) `进行中`(1) `rejected-needs-fix`(3) `withdrawn`(1) `ready-for-review`(1) `in-progress`(1) 全部映射到 v2 词表值并**照条打印**。理由写在脚本头：让 3b 报「真缺陷」，把「已知迁移债」留在 3c 里可数可烧；两半都不许消失。

## 三、`--self-test` 输出（22 条控，rc=0）

```
ticket-lint.sh: positive control bench=/tmp/tmp.OeUaxfMLxx (kept on disk; this project never deletes temp artifacts)
ticket-lint.sh: control 1 ok - the bench baseline exits 0
ticket-lint.sh: control 2 ok - the healthy pool has no RED line at all
ticket-lint.sh: control 3 ok - the healthy pool reports all-zero findings
ticket-lint.sh: control 4 ok - the done ticket whose only done marker sits on an appended line is NOT reported by item 2
ticket-lint.sh: control 5 ok - the scan still exits 0 with planted violations on the table (v1 reports, never blocks)
ticket-lint.sh: control 6 ok - item 1 names the -done ticket that still has a line-start unchecked box
ticket-lint.sh: control 7 ok - item 2 names the -done ticket with no done word on any status line
ticket-lint.sh: control 8 ok - item 2 does not name 02-appended-done-done.md
ticket-lint.sh: control 9 ok - item 3a names the live ticket with no status line
ticket-lint.sh: control 10 ok - item 3b names the out-of-vocabulary value verbatim
ticket-lint.sh: control 11 ok - item 3c reports the legacy spelling as a note, not a red
ticket-lint.sh: control 12 ok - item 4 names the two files that share one ticket number
ticket-lint.sh: control 13 ok - item 4 reconciles gap 12 against voided/
ticket-lint.sh: control 14 ok - item 4 reports gap 14 as unexplained
ticket-lint.sh: control 15 ok - item 5 names the over-hat name and says it is not grandfathered
ticket-lint.sh: control 16 ok - item 6 ages the review ticket
ticket-lint.sh: control 17 ok - item 6 ages the blocked-on-owner ticket on its own threshold
ticket-lint.sh: control 18 ok - item 7 counts every bucket from the same records
ticket-lint.sh: control 19 ok - item 2 has teeth: removing the appended done word turns that ticket red
ticket-lint.sh: control 20 ok - the mutated scan still exits 0
ticket-lint.sh: control 21 ok - item 3b has teeth: a second bad value is caught and the note for the old one disappears
ticket-lint.sh: control 22 ok - item 3c notice is not stuck on: the normalised spelling is gone once the value is bad
ticket-lint.sh: positive control PASSED - 22 control(s) ran, 0 failed
```

正控的形状：坏样例**逐条验红**（控 6–18），好样例**不出红**（控 2/3/4/8/11），并且**两发牙齿**（控 19：把好样例的追加 done 拿掉 ⇒ 立刻红；控 21/22：换一个坏值 ⇒ 立刻红且旧 note 消失）。底板先跑一遍绿（控 1–3）再种坏样例，防「卡死在红」。自测全程 `mktemp -d` ＋ `git init` 的临时池，**⛔ 没碰真票池**。

## 四、退出码语义（v1 只报不拦）

| 场景 | rc | 依据 |
|---|---|---|
| 真票池跑完（有 32 红 101 黄也一样） | **0** | 规格「v1 只报、v2 进 CI」；`VERDICT reported, not enforced` |
| `--self-test` 22 控全过 | **0** | — |
| `--self-test` 任一控不符 | **1** | 尺没牙就不许印读数 |
| 拒绝作答（无工作树／无池／无 `HAT_NAME_LIMIT` 行／分母 0） | **2** | 不算「跑完了」，故非 0 |

实测：真池 `rc=0`；`--self-test` `rc=0`；`--with-self-test` `rc=0`（自测＋全量，4.69 s）。

## 五、性能（365 张全池，含 189 张 `-done`）

`1515 / 1626 / 1638 / 1721 / 1893 ms`（成品现跑五次，含 `--with-self-test` 那次 4691 ms 是自测＋全量）。一次 `awk` 走池 ＋ 一次 `awk` 出报告 ＋ 一次 `awk` 分组，无 `git log`／无逐票 fork，日常跑得起。

## 六、统计输出样例（第 7 条，实测逐字）

```
ticket-lint.sh: R7 STAT done=189 live=175 void=1 blocked-on-owner=1 blocked=4 ready-right-now=78 ready-but-blocked=43
ticket-lint.sh: R7 NOTE ready-right-now = live tickets whose status is ready and whose Blocked-by field is absent or empty; 52 of 175 live tickets carry a non-empty Blocked-by field, so that count is only as complete as the field is
ticket-lint.sh: SUMMARY R1=3 R2=66 R3a=20 R3b=12 R3c=128 R4dup=0 R4gap=14 R5=58 R6=15 tickets=365 hat=100 unit=bytes
```

**这替代了多少手数**：审计手数的「189 done／175 活／14 缺号／58 超长名／20 无状态行」全部一行出；`done=189 / live=175 / void=1` 与审计 §4 逐数相等。审计的「~65 立即可派／~42 等 owner」是**深读判定**，机读字段读不出来：`blocked-on-owner` 全池只有 1 张，`blocked` 4 张，`Blocked-by` 字段只有 78 张票有（52 张活票有真前置）。⇒ 78 是**上界**，42 那个数是人判的，两者不可互相替换——这一句已写进脚本的 R7 NOTE，防后人拿 78 当「等 owner 42」的反证。

## 七、踩到的坑（本腿最容易错的地方）

1. **Status 不是两种写法，是四种＋一种续行**（2026-10-26 全池量）：
   `A **Status:** value`（174 票）／`B - Status: value`（137 票）／`C - **Status**：value`（17 票）／`D Status: value`（4 票）／`E → **accepted-done**（…）`（续行 11 行）。
   **只认 `- Status:` 会漏 174 票的 A 形**；**只认首行会撞 E 形**：`95 97 99 101 102` 五张 `-done` 票的 done 标记**只在追加行上**（首行还写着 `open`），尺若只读首行 ⇒ 5 张假红（71 vs 66）。
2. **E 形不能按「行首有箭头」认**：池里有 `  → within cap or alert`、`→ exit 1`、`→ \`sed -n …\``（22 行行首箭头、35 行缩进箭头，多是正文）。认法＝**行首 ＋ 接在真 head 之后 ＋ 载荷是单个加粗 token**（`→ **word**（`）。改前实测 25 张 3b，改后 12 张——那 13 张是我的尺在把散文当状态。
3. **多字节不能进中括号**：`[：]` 在字节模式 awk 里会吃到 `（`（EF BC 88）的首字节 `EF`，把 `- Status（…` 误判成 B 形并切出乱码取值。全脚本非 ASCII 字面量一律用**完整序列做交替** `(:|：)`，绝不进 `[…]`。`--self-test` 无覆盖此坑的样例，靠注释＋本件留档。
4. **零填充撞整数键**：缺号走整数循环，`"01"` 与 `1` 不是同一键 ⇒ 曾把 1–9 全报成缺号（合成池 4 张票报 4 个假缺号，控 3 当场抓住）。修法＝`numseen[num + 0]`。
5. **`in progress` 两词**：空格切值会切成 `in` ⇒ 加一条 `in progress → in-progress` 归一（票 71）。

## 八、边界与已知局限

- **v1 不进 CI**：`.github/**` 零改动；脚本本身 `exit 0`（跑完就算成功），红只体现在输出。`--with-self-test` 是给人用的本地形状，不是 CI 步。
- **第 5 条共用阈值、零复制**：`HAT_NAME_LIMIT` 由 `sed` 从 `scripts/check-path-length-budget.sh` 的 `HAT_NAME_LIMIT=100` 行**现读**，读不到就 `rc=2` 拒绝作答（不猜、不带第二份数）。262 的 roster 也现读，只用来给超长名打「已在册／新犯」标记。
- **不认的形状**：`> **Status …**：`（引用块前缀）未识别——实测无票只靠它、且审计那 20 张无一张带它，故今天不产生错判；v2 应补。
- **老化精度＝日历日**（票面只有日期、没有时刻），故 `72h` 读作「>3 天」。第 6 条不覆盖 `blocked`（规格只写 `blocked-on-owner`），`77` 因此不点名。
- **只读**：除 `mktemp -d` 的仓外工作目录，脚本对仓库零写；跑完 `git status --porcelain .scratch/wisp/issues/` 只有**起手就脏的 115**（非本腿所改），票池零变化。

## 九、门禁 rc

| 门 | 命令 | rc |
|---|---|---|
| 语法（bash） | `bash -n scripts/ticket-lint.sh` | **0** |
| 语法（POSIX sh） | `sh -n scripts/ticket-lint.sh` | **0** |
| 自测 | `sh scripts/ticket-lint.sh --self-test` | **0**（22/22） |
| 自测＋全量 | `sh scripts/ticket-lint.sh --with-self-test` | **0** |
| 全量（只报） | `sh scripts/ticket-lint.sh` | **0**（365 票，32 红不掉 rc） |
| 票池零改动 | `git status --porcelain .scratch/wisp/issues/` | 只有起手已脏的 `115` |
| 未碰 262 门 | `git diff --stat scripts/check-path-length-budget.sh` | 空 |

**没做（具名）**：⛔ 未把本脚本接进 `.github/workflows/ci.yml`（规格写明 v1 不进，那是 v2 的事）；⛔ 未回填任何票面 Status（回填是阶段 1 扫除腿的活，本腿只读）；⛔ 未给 `voided/` 补档；⛔ 未跑 `go` 任何门（本腿零 Go 改动）。
