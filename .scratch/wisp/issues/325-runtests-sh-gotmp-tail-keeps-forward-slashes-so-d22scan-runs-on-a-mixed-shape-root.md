# 325 — `tools/d22scan/runtests.sh` 的测试临时根＝**混形活坑**（基座反斜杠＋尾接正斜杠）：同一族分隔符病在票 317 射程外的那一处

- **Status**：已立（2026-10-16 09:2x，编排者），**派 `325-r1`**（写面＝`tools/d22scan/**`，起手时该面无写腿）
- **Type**：测量完整性缺陷（与票 317 同族、同一种写法），⛔ 不是功能缺陷、⛔ 不是为了变绿
- **来源**：票 317 的非实现者验收腿 `317-v1`（表＝`docs/evidence/s1/317-gotmp-separator.md`，笔 `ceaa50ac`）§4 具名上报；编排者 09:2x 现读复核并附自己的成对复跑（`.scratch/wisp/probes/orch/20261016/logs/10-nat.txt`／`10-psx.txt`）。

## 这是什么（人话）

跑测试时，Go 要把临时文件放在哪，是由一个环境变量（`GOTMPDIR`）说定的。我们这台 Windows 上，这个值**必须写成带反斜杠的样子**（`C:\Users\…`）；一旦里面混进正斜杠（`C:/Users/…`），有几条专门测"目录链接（junction）"的用例就会**假红**——不是代码坏了，是量它的尺自己把路径写坏了（`mklink` 会把 `C:/…` 当成命令行开关，报 `Invalid switch`）。

票 317 已经修好了三处这种写法（都在 `scripts/` 里）。**同一族写法还剩一处没修**，就在 `tools/d22scan/runtests.sh:92`：它的开头用了反斜杠、结尾接名字时用了正斜杠，所以生成的路径是"半对半错"的混形。今天没炸，只因为这条脚本的射程里没有那几条会炸的用例——**这是"坑还在、只是暂时没人踩"，不是"没坑"**。

## 现量（编排者 2026-10-16 09:2x，全部自己跑）

1. 字节形（`sed -n '92p' tools/d22scan/runtests.sh`）＝`gotmp_run="${WISP_TEST_GOTMP_ROOT:-"$LOCALAPPDATA\\wisp-test-gotmp"}/d22scan-$(date +%Y%m%d-%H%M%S)-$$"` ⇒ 基座 `\\`（原生）＋尾接 `/`（正斜杠）＝**混形**。
2. 已修的三处对照（同尺）＝`scripts/portable-tests.sh:155-156`／`scripts/winsec-tests.sh:69`／`scripts/wisp-cli-tests.sh:82` 三处**尾接都是 `\\`**＝原生形。
3. 病本身可复现（我这把尺，仓外导出树、同一个目录只差分隔符）：`GOTMPDIR='C:\tmp\wisp317orch\pair'` ⇒ `RUN=4 PASS=4 FAIL=0` rc=0 ⇄ `GOTMPDIR='C:/tmp/wisp317orch/pair'` ⇒ `RUN=4 PASS=0 FAIL=4` rc=1，红句逐字同行同号（`internal/risk/syncdirs_ancestor_actable_leg_116_test.go:182`／`:286`＋`…reparse_ticket105_windows_test.go:59`／`:116` 的 `Invalid switch - "tmp".`）。

## 要建什么

- [ ] **AC#1 改形**：把 `tools/d22scan/runtests.sh:92` 的尾接 `/d22scan-` 改成 `\\d22scan-`（与 `scripts/**` 那三处逐字同形），旁边留一条说明"为什么这里不许折成正斜杠"的注释（⛔ 只写"风格"、要写**会炸的那条用例**）。判据＝改前/改后各跑一发 `sh scripts/d22scan.sh`（干净台面：`GOTMPDIR`／`WISP_TEST_GOTMP_ROOT` 抽干），两发都 rc=0 clean，并把那一行 `GOTMPDIR=` 的逐字形状落在读数件里。⚠ **⛔ 顺手把 `wisp_311_ref_path` 那四处对照位（`portable-tests.sh:818`／`winsec-tests.sh:174`／`wisp-cli-tests.sh:186`／`runtests.sh:155`）一起"改一致"**——票 317 已逐处判过它们是对照位显示值、⛔ 流向 `GOTMPDIR`；动了＝越界，要具名上报。
- [ ] **AC#2 成对反控（这族尺的牙）**：在**仓外导出树**上把这一行退回混形（物理改字节，⛔ `-overlay`——这族尺运行时读盘），跑一次能让 `./internal/risk/` 那四枚用例红的构造 ⇒ 证明"混形会炸"今天仍然真；然后还原、复跑绿、记录还原哈希。⛔ 只做"改前必红"就交，还要写清"射程里今天有没有那四枚"——如果 `d22scan` 那档根本不调 `./internal/risk/`，就具名写成**缺世界**（坑活着但今天看不见），⛔ 假装测到了。
- [ ] **AC#3 顺带交出票 317 AC#3 那一发**（同一台面、同一档，省一枚腿）：修完后跑一发 `bash scripts/portable-tests.sh --scope=windows`，三样齐＝① `GOTMPDIR=` 行逐字原生形 ② 八包四数**逐包**落盘 ③ 红名册与 10-14 那两发 CI 的 `Invalid switch` 名册作差＝新增 0。
- [ ] **AC#4 门禁**：`bash -n tools/d22scan/runtests.sh`／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh --with-self-test`／`go build ./...` 逐枚自落 `rc=`。⚠ 两枚台面雷：d22scan 会把 `GOTMPDIR` 指到非 `%LOCALAPPDATA%` 卷时回 **rc=2 的票 311 GUARD 射程拒绝**（⛔ 当红）；`cmd/wisp` 测试缺 `third_party/sherpa-onnx` 那三枚 DLL 会回 `0xc0000135`＋`RUN=0`（＝台面故障，⛔ 读数）。
- [ ] **AC#5 账**：票 317 面上那句"同形、今天未爆"要**重述**成"混形，射程遮蔽中的活坑"——由编排者在 317 的 Progress log 追加，⛔ 本腿改 317 原文。

## 禁区

⛔ 动 `scripts/**`（那三处已由票 317 修毕并验过）／`wisp_311_ref_path` 四处对照位／`.github/**`／`ci.yml`／`thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件／`frontend/**`／`design/**`／`PLAN.md`／`specs/**`；⛔ 为变绿放宽任何断言或删用例或加 `t.Skip`；⛔ push（只 commit，显式 pathspec）；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 派单前置条件

无（`tools/**` 面空）。写面三处＝`.scratch/wisp/probes/325/r1/`＋`docs/evidence/s1/325-runtests-gotmp.md`＋票 325 末尾 1 条 bullet；⛔ 勾 AC 框（翻框归编排者）。

## Progress log (append-only, newest last)

- [2026-10-16 09:2x +08] 编排者立票（来路＝`317-v1` 具名上报＋我自己现读三把尺：字节形／已修三处对照／成对复跑 4 红⇄4 绿）。⚠ 出身笔⛔ 在本腿现场反查：`317-v1` 记的是 `581dd252`（编排者），立票时未复跑 ⇒ 留给 `325-r1` 第 0 步 `git log -1 581dd252` 自证。
