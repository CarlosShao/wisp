# 262-close-r1 — 票 262「能否收口」核验（裁决者轮，判＝不可收口）

- **腿名**：Wisp 写腿 `262-close-r1`（核票/裁决轮，非实现者）。**时刻**：2026-11-09 14:1x–14:3x +08。
- **派单**：pool-audit-r6 §1 第 6 条称票 262 要补的「跟踪路径长度尺」已存在（被 230-r1／216-r2／230-v1 多腿引用 GREEN）⇒ 建议编排者收口。本腿按三核现跑后**具名判＝不可收口**：AC 框一枚未翻、票名未加 `-done`。
- **写面**：本件＋票 262 Progress log 追加。**零产码**；`frontend/**`·`design/**`·`cmd/wisp`·`internal/**`·`scripts/**`·`.github/**` 零写入。

## §1 三核①：票面未勾 AC 原文（现读，非转述）

- 票文件＝`.scratch/wisp/issues/262-tracked-path-length-gate-for-ci-checkout.md`（**未** `-done`，40,524 字节＋本轮追加）。
- `- [ ]` 现读＝**恰 2 枚**（核尺 `grep -c "^- \[ \]"`＝2，本轮前后同值＝一枚没翻）：
  - **AC#5**（`:64`）：「CI 侧真出颜色：一次推送（或 run 重跑）后，该步在**本机 self-hosted 那一侧**跑过且有日志行可指；托管 Windows／Linux job 不误加……」
  - **AC#5ⓑ**（`:73`）：「（＝原文那句"该步在**本机 self-hosted 那一侧**跑过"）：**现形下结构上不可能满足**，⛔ 本格不勾、也不许拿 ⓐ 的读数替它填。」——**禁勾写在 AC 原文里**。
- 已勾 8 枚里，**AC#1（脚本存在且被真调用）/ AC#2（HEAD 跑＝绿＋分母读数）** 早在 2026-10-07 `A674` 验收轮翻勾；AC#5ⓐ（CI 侧颜色＋托管侧不误加那支）同轮已勾。
- ⇒ **audit 前提指认错格**：它说的「要补的那格」＝"尺存在"＝AC#1/#2，**早已是勾**；今天未勾的两格要的是 **CI 侧 self-hosted 真跑证据**，与「尺存不存在」无关。

## §2 三核②：尺在仓、但 HEAD 上是 RED（不是 GREEN）

- 脚本在仓：`scripts/check-path-length-budget.sh`（48,586 字节）——「存在」这半成立。
- 本树裸跑（工作目录 HEAD `cc5b297e2`）：`bash scripts/check-path-length-budget.sh` ⇒ **`VERDICT RED`／rc=1**。
  读数：`denominator: tracked paths=19218  over-budget=67  covered by roster=63  not in roster=4`；
  `bands: over the hat=67 … roster entries=65`；`RED - count guard: the roster holds 65 entries, the tree has 67 over-budget tracked paths`。
  ⚠ 首轮读数经管道尾取码得 0＝`tail` 的退出码，**已作废**；改直读脚本退出码。
- ★**shared clone 导出树复跑**（`git clone --shared` → `/tmp/w262c`，HEAD 同 `cc5b297e2`）：**rc=1，VERDICT RED，逐行同** ⇒ 红不是脏树、不是别人在飞的未提交件，**是 HEAD 真状态**。
- **红因（逐枚点名）**：`not in roster` 4 枚全是近期 `-done` 改名（已推）把末段名推过 100 帽：
  1. `138-…-terminatejobobject-done.md`（名 105 字符，full 170 > 165）
  2. `140-…-windows-runner-done.md`（名 113，full 178）
  3. `148-…-is-unwritten-done.md`（名 157，full 222）
  4. `301-…-is-package-level-done.md`（名 103，full 168）
  另 `roster drift` 2 枚＝140/148 的**旧名**（`-done` 前的名）已不是跟踪路径；count guard 65≠67。
- ⇒ **「尺真 GREEN」这半不成立**：HEAD 今日读 RED。且分母历史 14721（216-v3 时）→17219（230-v1 时）→19218（今日），凡引用必带时刻。

## §3 三核③：audit 点名的 GREEN 引文真在、但都是「当时」的读数

- `docs/evidence/s1/216-v3-verdict.md:49`：`check-path-length-budget.sh rc=0 GREEN（tracked 14721／over 60／在册 60／not in roster 0）`——真在。
- `docs/evidence/s1/230-v1-verdict.md:73-75`：`rc=0 VERDICT GREEN，denominator=17219、over-budget=60＝roster=60、not-in-roster=0`——真在。
- 另多枚同类：`109-v2-verdict.md:124`／`114-r2-impl.md:130`／`114-v2-verdict.md:126`／`115-…corrections.md:127`／`120-sealhandle-fixture.md:130`／`124-r1/r2`——全部 rc=0 GREEN，**各自当时**。
- ⇒ 引文层核过＝audit 没造假；**判断层不成立**＝"别腿引用过 GREEN"推不出"HEAD 今天 GREEN"。今天引用这把尺只会读到 §2 那发 RED。

## §4 判语与处置

**判＝票 262 不可收口。** 四条理由：
1. 未勾两格（AC#5／AC#5ⓑ）与「尺存在」无关——audit 前提过期，指认错格。
2. AC#5ⓑ 原文自带「⛔ 本格不勾」；AC#5 父格被 ⓑ 卡死——谁翻勾谁违 AC 原文（本腿禁改 AC 原文，只许翻勾，而这两格**不许翻**）。
3. 2026-10-07 编排者终态裁定（票面 Progress log「⑥ 状态尺（终态，⛔ 不加 `-done`）」）未被废止，仍是现行裁定。
4. 门在 HEAD 此刻是 RED——在门红时刻宣告"收口"更不成立。

⇒ 本腿**零翻勾、零改名**；写面仅＝票面 Progress log 追加＋本件。

## §5 ★上报编排者的新账（⛔ 本腿不修，超出本腿写面）

**HEAD 门红是真发现**：票 138/140/148/301 的 `-done` 名越帽（§2 点名）。处置二选一：
- **循先例改短**：`f67af953` 曾把 273/274/275 的名改短后门回绿（名册 57＝树 57，正文 R100 一字未动）；对这四枚同法可正控回绿，不动名册、不动帽。
- 或对名册／帽出新裁定（抬帽与通配名册均为票 262 自身禁区，需编排者解冻）。

**修好前，每条跑这把尺的腿都会读到 RED（rc=1）**——包括把该尺列入门禁引用的腿（AC#6 式「读数不扩大」比对会被这发红干扰）。建议尽快裁。

## §6 本轮门禁

- 零产码 ⇒ `gofumpt`／`go vet`／`go build` 无对象（未跑＝无对象可跑，非跳过）。
- `sh scripts/d22scan.sh` ⇒ **rc=0 clean**（日志 `/tmp/w262-d22.log`：ban #8 internal/=971、cmd/=245 枚等，全 clean）。
- Git 纪律：只 commit、⛔ 零 push；add/commit 两处显式 pathspec（票 262＋本件两枚）；无 amend/reset/rebase/stash/checkout ./clean；临时件 `.scratch/262-close-r1-append.md`（Progress log 追加片段）只建不删、不入库。
