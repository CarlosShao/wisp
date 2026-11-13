# 262 — 跟踪路径长度预算门禁：派单执行证据件（2026-11-13）

**腿**：262 派单写腿（写面＝本件＋票 262 Progress log；脚本零字节改动）。
**派单**：`.scratch/wisp/probes/orch/20261113/txt/262-dispatch.txt`（52 行，已从头读到尾）。
**现跑时刻**：2026-11-13 20:5x–21:0x +08；HEAD＝`02939b8a2`（dev）；跑尺时工作树对脚本/票面零 diff。
**结论先行**：派单要建的尺**已在仓、已接线、已验收、今天绿**。本腿没有重造它，而是把派单点名的每一项证据**自己现跑了一遍**，并把中间带（今日读数 **65 枚**，非派单引用的 57）逐枚对到了名册上。

## 1. 派单前提与票面现量的冲突（本腿的复认，⛔ 不照过时前提动手）

派单说「建 scripts/check-path-length-budget.sh（新）」。票 262 面的 Progress log 记载：

- 该脚本**2026-10-04 已由落地腿 `262-r1` 交付**（commit `519e1ade`＋`6a735939`，608 行现量）；
- **2026-10-07 已由验收腿 `262-v1` 对抗验收**（AC 框 8 勾／2 未勾，未勾两格＝AC#5ⓑ 结构不可能＋AC#5 父格，属编排者裁定遗留）；
- **已接进 CI**：`.github/workflows/ci.yml:176`（步名 `Tracked path-length budget (ticket 262)`）＋`:207`（`run: sh scripts/check-path-length-budget.sh --with-self-test`；本腿 `git grep` 现跑命中）；
- `docs/evidence/s1/262-close-r1.md`（2026-11-09）已裁「不可收口」并指出当时 HEAD 门红（四枚 `-done` 改名越帽）——该红随后已由名册跟进修绿（脚本最后一次内容 commit＝`4aee763af`，2026-11-13 11:58）。

⇒ 处置＝**不重造仪器**（重写一把已验收的尺＝无收益的复杂，且写面会撞 `.github/**` 邻域）；本腿改为**全量现跑复验＋把派单要的证据落成这一件**。脚本文件本身与 HEAD 逐字节一致（`git status --porcelain -- scripts/check-path-length-budget.sh`＝空）， hence **不进本次 commit**。

## 2. 口径（派单回话第 1 条）

**口径＝字符数（characters），不是字节数**；判据是 README 规则 9 的名帽按字符计、Windows 检出失败按 UTF-16 码元计，字符是物理上对的尺。脚本 `:75-96` 头注释写明：awk 的 `length()` 默认计**字节**，除非 multibyte-aware＋UTF-8 locale（CI 的 ubuntu-latest 默认 mawk 不是）——所以脚本**每次运行现场探测自己的 awk**（`:249-257`，5 字节探针 `a U+25FF b`），字符口径则继续、字节口径则**拒绝作答（exit 2）**而不是带假读数判树。

**中文名实测（本腿一手，2026-11-13 21:0x）**：

| 尺 | 读数 |
|---|---|
| 字符串 `中文路径测试票.md` | **10 字符** |
| `wc -c`（字节） | **24 字节** |
| `LC_ALL=C awk length()` | **24**（字节口径把中文名读长 2.4 倍） |
| `LC_ALL=en_US.UTF-8 awk length()`（GNU Awk 5.4.0） | **10**（与字符数一致） |
| 脚本探针（本机跑出） | `unit=bytes (probe: length of the 5-byte sample "a U+25FF b" = 5; LANG=unset LC_ALL=unset)` |

⇒ 本机默认 locale 下脚本以字节口径运行；它对此**不装傻**：本仓现量 **35 枚**跟踪路径含非 ASCII（`git ls-files | awk '$0 ~ /[^ -~]/' | wc -l`），若字节口径下出现越帽 offender 它会打 NOTE 并 REFUSE（exit 2），不会拿 2.4 倍虚长的读数去点别人的名。今天 69 枚 offender 全 ASCII，双口径 offender 集**同集合同数**（69＝69，两组清单已存 `/d/tmp/wisp/20261113/262/{char,byte}-view-offenders.txt`，`cmp` 全等）。

## 3. 分母行（派单回话第 2 条，逐字）

```
check-path-length-budget.sh: denominator: tracked paths=19530  over-budget=69  covered by roster=69  not in roster=0
```

同一次运行的其余读数（逐字）：

```
check-path-length-budget.sh: denominator read: 19530 tracked paths
check-path-length-budget.sh: longest=185 chars relative (.scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-short-name-stays-short-long-name-folds-long-so-in-root-writes-read-as-out-of-bounds-r2-l2-done.md)
check-path-length-budget.sh: worst full path on the self-hosted runner=229 chars (hat budget 165, wall open interval (206,217])
check-path-length-budget.sh: bands: over the hat=69  of which in the 122..180 middle=65  past the old debt line(180)=4  in the wall interval=0  roster entries=69
check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree
```

**分母不是常量**（票面 AC#2 既有裁定）：7433（10-07 编排者）→ 19218（11-09 close-r1）→ 19255（11-09 CI lint）→ **19530（本腿 21:0x）**。凡引用必带时刻。

## 4. 正控与负控（派单回话第 3 条，两次 exit code ＋逐字红句）

**负控（当前树，直读）**：`sh scripts/check-path-length-budget.sh` ⇒ **exit 0**（§3 的 GREEN 全文）。

**正控（种一枚超帽跟踪路径，本腿一手）**：方法＝**仓外一次性仓**（`D:/tmp/wisp/20261113/262/bench-direct`，先例＝票面 10-07 节「AC#3 那一发必须在仓外 clone 里种」——在本仓工作树留越帽跟踪路径会让每一条并发腿的尺读红）。在名册摘除后的门副本下三跳：

| 跳 | 读数 |
|---|---|
| 种前（1 枚短名跟踪路径） | `NEG_RC=0`（绿） |
| 种入 113 字符名（相对 134＝越帽 121 不越墙 206/217） | **`POS_RC=1`**，红句逐字点名（下引） |
| `git rm --cached` 拆掉 | `POST_UNPLANT_RC=0`（绿） |

红句逐字：

```
gate.sh: RED - over budget and NOT in the roster: .scratch/wisp/issues/seeded-xxxx…-over-the-hat.md  (relative 134 chars, name 113 chars, full path on the self-hosted runner 178 chars, budget 165)
```

（`name_len=113`、`rel_len=134` 由 `wc -c` 独立复核，与红句读数一致。）

**内建正控（`--with-self-test`，CI 每次推送真跑的那条命令）本腿一发全过**：control 1/3（bench 基线绿）→ control 2/3（种相对 125 的跟踪路径 ⇒ rc=1 且逐字点名 `relative 125 chars, name 104 chars`）→ control 3/3（拆掉转绿）→ `positive control PASSED` ⇒ **exit 0**（正控＋真扫全绿）。

## 5. 门禁四件套（派单门禁节）

| 门 | 命令 | 读数 |
|---|---|---|
| 语法 | `bash -n scripts/check-path-length-budget.sh` | **rc=0** |
| 负控 | `sh scripts/check-path-length-budget.sh` | **exit 0**（§3） |
| 正控 | 上节 | **exit 1 且点名** |
| 卫生 | `sh scripts/d22scan.sh` | **rc=0 clean**（`PASS=45 FAIL=0 SKIP=0`；scopes：bans #1-5 internal/=403、cmd/=77、ban #6 frontend/=123、ban #7 internal/tools/=51、ban #10 internal/risk/=21、ban #8 design/=39 frontend/=123 internal/=1005 cmd/=276） |

## 6. 中间带处置（派单回话第 4 条）——选了「逐枚进豁免名册＋理由」，且它已在场

**今天的数**：`git ls-files | awk 'length($0)>121 && length($0)<=180' | wc -l` ＝ **65**（不是派单引的 57——票面 §现量 6 自己就钉了「别照抄 57，先复跑」；57 是 2026-10-04 的读数，其后新收的长名票已逐枚入册）。

**对账（本腿一手，join 已修正引号剥法后复跑）**：
- 中间带 65 枚 **0 枚缺册**（`comm -23 middle roster` ＝ 0 行；清单 `/d/tmp/wisp/20261113/262/middle-band-65.txt`，65 行、区间 122–175）；
- 名册 69 ＝ 中间带 65 ＋ 债表线（>180）4 枚（`158`/`252`/`255`/`261`，正是相对 181–185 的那四枚）。

**选了哪个**：**名册方案**（派单「优先选」）。理由：
1. **零改名风险**：65 枚里 51 枚带 `-done`，改名要在保住防重领键的同时追改台账与证据件引文，风险全在语义层；名册方案对票池零扰动。
2. **它不豁免墙**：帽（>121）只是**预警线**不是安全线；真正打死 checkout 的是墙（相对 (206,217]）。最坏一枚 185 ⇒ 全路径 229，离墙全路径下界 250 还有 **21 字符**余量（`in the wall interval=0`）。改短的边际收益买不到任何一道新防线。
3. **帽的牙在别处**：新票名违反规则 9 会被守卫 A 当场点名（§4 正控就是 113 字符名的一发实弹），名册挡不住拦新票——名册只收**已历史存在**的枚。
4. 派单只授权「真超帽且必须收敛的类别」才改名；65 枚全部 ≤180，**没有一枚**落进那个类别。

名册理由现状（对票面 AC#4「类目句可核」口径的复认）：69 枚理由分三类——`added 2026-10-04` 历史类 **53 枚**（含 3 枚 dispatch note，规则 9 只管票名）、`-done` 追加致越帽的守卫 B 跟随类 **16 枚**、票 180 特记 1 枚（2026-11-13 收票加 `-done` 推过帽，理由具名写明这一笔）。**无空理由**（守卫 D 在场，票面 AC#4 有「删理由必红」的行为突变读数）。

## 7. 本腿没有做的事（红线对齐）

- ⛔ `.github/**`、`thresholds.go`、golden、`allowlist.txt`、契约、`models/**`、`cmd/wisp/**`、`internal/**`：**零写入**（本腿全部写入＝本件＋票 262 Progress log）。
- ⛔ 阈值未动一字：`HAT_NAME_LIMIT=100`（帽＝相对 121）原样；`WALL_LOW=206 / WALL_HIGH=217` 仅作读数打印、永不参与判红绿。
- ⛔ 未删任何票文件；未做任何 `git mv`（理由见 §6）；未动票 262 的 AC 框（终态裁定 2026-10-07：2 未勾＋8 已勾＋不加 `-done`，翻框权在编排者）。
- 未在本仓工作树种任何越帽路径（正控全部在仓外一次性仓；bench 保留不删，`/d/tmp/wisp/20261113/262/`）。

## 8. 交付物与 commit

- `scripts/check-path-length-budget.sh`：**已落盘且已在 HEAD**（`4aee763af`），本腿零字节改动，故不进本次 commit。
- `docs/evidence/s1/262-path-length-budget.md`：本件（新）。
- 票 262 面 Progress log：追加本腿一节（⛔ 未翻 AC 框）。
- commit：显式 pathspec 两条（本件＋票面），⛔ 无 push。hash 见回话。
