# 300-a2 · 90 — 终态锚＋名册差集＋欠读数（交回件）

- 起手锚＝4e00357f 2026-10-10T17:34:47+08:00（00-anchor.md）
- 终态锚与逐笔名册、差集＝**由尺现跑追加在本行下面**（同形＝git log -1 / git status --porcelain 计数 / tasklist 三枚计数；名册＝对 起手锚..HEAD 逐笔 git show 型 --name-only）
- 授权名册＝probes/300/a2/** 之内；⛔ 越界一枚、⛔ 动产码、⛔ push
- ⚠ 自我引用一处：本件的**最后一笔 commit**（追加终态读数那一笔）⛔ 能把⛔ 自己的号写进⛔ 自己 ⇒ 那一枚号在腿的回报正文里给，⛔ 编排者按 git log 现跑核对。

## 三件交付物（各带尺，见文件）

| 件 | 是什么 | 尺面 |
|---|---|---|
| 10-candidate-landing-points.md | `AC#6` 候选落点 **5 枚**（C1 既有 tagged 用例加段／C2 新枚同包 tagged 文件／C3 权威表读盘／C4 跨实现一致性／C5 生产侧抽谓词），逐枚写"动哪枚文件／判据形状／3→1 与 3→0 红⛔ 红的**依据**／买到⛔ 买到" | 五枚硬事实 F1–F5，全 blob 层 |
| 20-near-miss-nails.md | 擦边钉 **7 枚**（N1–N7）逐枚"断的是什么／撞⛔ 撞"，含"会撞的只有 N4（动容差/token 值）与 N6（新增未登记 SKIP 要改 scripts）两形" | 7 把名册级 grep，跨 internal/** ＋ cmd/wisp/** ＋ scripts ＋ tools/d22scan |
| 30-ci-visibility.md | 答复＝**落在 CI 上**，档＝test-windows 的 --scope=windows 那一发，OS＝windows-latest（托管），step 级 if 逐字为 !cancelled() 故**有求值**；ubuntu 档永远⛔ 求值它（GUARD A 注释自己写着） | portable-tests.sh 名册 × ci.yml job/runs-on/调用点，调用点一律内容锚 |

## 欠读数（具名，⛔ 静默；本腿⛔ 可取）

1. **新判据的成对两发颜色**（改⇒红／逐字节还原⇒绿）＝**归落地腿**（本腿⛔ 编译面授权：`go build/test/vet/run/list` 全部 0 发；`go env` 亦 **0 发**，此处自报）。
2. **CI 那一发色**（test-windows 档跑起来是 PASS/FAIL/SKIP）＝**归编排者**，票面 next 已钉"跟票 303 修复后那批推送一起取、⛔ 为它单推"。
3. **托管 windows 镜像上 SDK 头文件在不在**（候选 C3 的生死判）＝**归那一发 2**；仓内名册级尺对 Windows Kits/mmreg.h 在 scripts/tools/.github 的命中＝**0 枚**（本腿现跑）。
4. **internal/audio 在 windows 档的当前基线色**（票 301 AC#2 那发前后名册作差）＝**归编排者／301 的腿**；probes/301 在我授权名册外，本腿只按票面引、⛔ 查。
5. **TestLiveWasapiSmoke 在托管 runner 上到底是登记过的 SKIP 还是红**＝**归那一发 2**（ledger 那行的理由句自称托管无音频端点，⛔ 实测凭据在盘上）。
6. **cmd/wisp 那族 255 名册仪器会不会去看 audio 的字节面**＝**归 cmd/wisp 台面持有腿（303-r1）或编排者**；本腿只量到"⛔ 人按行号引 wasapi_windows.go／parse_wave_format_300*"（0 命中），⛔ 证明那族仪器⛔ 看字节面。
7. **顶回一处枚数口径**：派单/票面那句"整包 **41** 枚顶层"与本腿源码尺（blob 名册逐枚统计 func Test 行）数到的 **42 枚**差 1，而该包 func TestMain 命中 **0** ⇒ 两把尺⛔ 同形（PASS 行 vs 源码行），⛔ 我判谁对，落笔请带尺名。
8. **顶回一处派单措辞**：票面 AC#6 原文还有一半边派单⛔ 写——"且红句具名指向那枚常量而⛔ 指向 tag 的等值比较"；本腿按**原文**办（原文优先于转述），差异具名报回。
9. ⛔ 派单一处**⛔ 我拍**：C4（跨实现一致性）撞票面 AC#0 那句"⛔ 拿本仓另一枚实现当裁判"的**字面射程**——那一格钉偏移、本候选钉**值**，算⛔ 同一条禁令＝编排者裁；C3 撞既有注释声明＋要改 scripts（越出 AC#4 名册）＝编排者裁；C5 改产码语义面＝⛔ 本票授权范围，需编排者明示。
