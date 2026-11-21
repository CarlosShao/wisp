# 58 — S8: i18n / English locale + English voice models — DEFERRED (voice-chain caveat)

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 54-s7-acceptance + user go-ahead + English voice model selection (P-new)
**Parallel slots:** ≤2 sub-agents
**Spec refs:** D23 corrected i18n row, SPEC-12 §5 i18n caveat, D9, S8

## What to build
Internationalization done right per the corrected DEFERRED row: this is NOT just UI strings —
the whole voice chain is Chinese models, so non-Chinese users are voice-dead. Deliver UI string
externalization + full English locale + English ASR/TTS/KWS model selections passing
English baseline sets.

## Key constraints
- UI strings externalized from S1 onward (ticket 07/34 kept them out of code); locale files +
  `app.language`; RTL not in scope.
- English voice: sherpa-compatible English KWS/ASR/TTS candidates evaluated (license per P3
  discipline); baseline wav sets (en near_clean/noisy_far) + CER gates re-established (6%/15%
  analogs); punctuation model for English.
- Locale switch = hot; voice-locale = reload tier; mixed-language utterances documented
  (best-effort).
- Docs/README bilingual minimum (EN + ZH).

## Acceptance criteria
- [ ] All UI via locale files; zero hardcoded Chinese in frontend/native render paths (scan).
- [ ] English voice chain passes its committed baseline gates; latency budgets re-measured.
- [ ] Locale switch e2e (UI + voice reload) without restart.
- [ ] DEFERRED row closure evidence filed.

## Progress log (append-only, newest last)

- **2026-11-21 16:5x +08 — AC#1 wave 1（界面文案外置＋渲染路径 CJK 扫描牙），实现腿（GLM 5.3 Flash）。只做 AC#1；AC#2/AC#3（语音链）仍阻塞于票 54＋机主 go-ahead＋英文语音模型选择，本波未碰。**
  - **普查（改动前，frontend/src/\*\* 渲染路径＝main.tsx import 闭包，注释不计，仅数字符串字面量＋JSX 文本行）**：闭包 92 文件 **1602 枚**（字面量 1189＋JSX 文本行 413）；闭包外（未挂载演示/夹具）22 文件 57 枚。普查前快照存 `.scratch/wisp/probes/orch/20261121/58-ac1-census-pre.txt`。口径与工具：`frontend/scripts/cjk-census.mjs`（`npm run` 不挂，`node scripts/cjk-census.mjs`；bands＝汉字 U+3400-4DBF/4E00-9FFF/F900-FAFF＋假名＋谚文，全角标点不计）。
  - **Go 侧 native render paths（AST 字符串字面量、注释天然不计；工具＝临时 go 程序，读数存 `.scratch/wisp/probes/orch/20261121/58-ac1-census-go.txt`）**：`internal/panel` 277 枚/43 文件、`internal/ball` 7 枚/29 文件、`internal/present` 3 枚/4 文件、`cmd/wisp` 992 枚/103 文件（控制台回执/文案＝CLI 渲染路径）。`internal/tools` 811 枚是面向 LLM 的工具描述文案，不算 UI 渲染路径，另列。**Go 侧本波零改动**（不在写面），全部记欠账。
  - **locale 基础**：新增 `frontend/src/lib/locale/`（`zh.ts` 240+ 键中文、`en.ts` 同键 fluent 英文、`index.ts` 手写 `t(key, params?)`＋`setLanguage()`，默认 zh，零 npm 依赖、零 i18n 库；en.ts 由 `Record<LocaleKey,string>` 类型强制与 zh 键一一对应，缺键/多键 tsc 直接红）。无状态纪律：语言选择不落 localStorage，将来由 config wire 的 app.language 推送驱动 setLanguage（票 58 的切换 AC 另行接线）。
  - **外置（4 个产品面文件，普查枚数最多的非夹具/非演示文件）**：`model-settings.tsx`（169→0）、`lib/config-feedback.ts`（93→0）、`voice-settings.tsx`（71→0）、`App.tsx`（64→0），共 397 枚搬进 locale，全部改 `t(key)`；新增 `lib/config-protocol.ts` 承载 Go 封句 match 令牌（被拒绝/盘上已变等，只匹配不上屏）。默认 zh，`t()` 返回与原字面量逐字相同，渲染输出不变（`npm run render:model-settings` 仪器 rc=0 复核）。
  - **扫描牙（AC#1 判据）**：新增 `frontend/scripts/lint-cjk.mjs`（package.json：`cjk` / `cjk:self-test`）。扫 main.tsx import 闭包全部 .ts/.tsx（注释先剥，同一把尺来自 cjk-census.mjs，普查与门禁不分歧）：locale 目录豁免（zh.ts 是中文的唯一合法居所）；`config-protocol.ts` 白名单（理由：Go 封句 match 令牌，匹配用不渲染，外置即断契约接缝）；其余闭包文件一一登记 **55 条白名单（cap＝落牙时计数，只许降不许涨；归零须摘条；未登记文件必须 0）**。正控五种全红：清洁样本绿、种一句中文的字面量红（逐字红句：`src/components/seed-cjk-sample.tsx:1 literal "扫描牙正控：这句中文必须红"`）、种 JSX 文本红、cap 涨红、白名单过期红；另做**真盘临时树正控**（OS 临时目录种 main.tsx+seed.tsx，闭包走查命中后红，同句逐字复现）。`npm run cjk:self-test` rc=0；`npm run cjk` rc=0 findings=0。
  - **自证**：`npx tsc -b` rc=0；`npm run build` rc=0；`npm run render:model-settings` rc=0；`npm run seeds` / `seeds:self-test` rc=0（未动 77 号牙的任何族）；`sh scripts/d22scan.sh` **OK - control rc=0, scan rc=0（clean - no D22 ban violations）**；⛔ 真窗用例未跑；⛔ PLAN/specs/thresholds/golden/allowlist/design/** 一字未动；⛔ 无 emoji；⛔ 未 push。
  - **欠账清单（诚实账本，已固化进 lint-cjk.mjs 的 REGISTRY，55 条 cap）**：产品面 backlog **40 文件 587 枚**（config-screen 49、cost-settings 49、plugins-screen 40、privacy-settings 40、palette-screen 37、history-page 31、diagnostics-settings 30、config-editor 25、security-settings 25、prompt-bar 23、settings-controls 22、l2-approval-card 21、composer-row 18、result-history 17、session-rules-settings 15、approval-screen 14、panel-sidebar 13、git-branch 11、terminal-region 10、panel-views 9、tasks-screen 8、firstrun-screen 7、thinking 6、error-boundary 6、skills-settings 6、history-view 6、approval-queue 5、ball-screen 5、diagnostics-view 5、tools-discovery-view 5、config-read-state 4、context-meter 4、tier-names 4、tool-chips 2、cost-view 2、panel.ts 2、session-rules-view 2、config-actions 3、privacy-view 3、motion 3）；演示面（harness/showcase/fixtures，?harness 路由专用非产品文案）**14 文件 618 枚**；协议令牌 `config-protocol.ts` **1 文件 23 枚**（永久豁免理由见上）。zh.ts 286 枚＝中文合法居所。后续波次每清一个文件即从 REGISTRY 摘条，cap 只降不涨。
  - **留给编排者的两件事**：① CI 接线——`.github/workflows/ci.yml` 不在本波写面，需在 lint-frontend 作业底部按 77 号先例追加两步（`npm run cjk:self-test`、`npm run cjk`，均带 `!cancelled()`），接上之前这枚牙只在本地有牙。② Go 侧（internal/panel 277＋cmd/wisp 992＋internal/ball 7＋internal/present 3）未建扫描牙，属 AC#1 原生半边的后续波次。
  - AC#1 未勾（勾归非实现者）。

- **2026-11-21 18:3x +08 — AC#1 wave 2（产品面第二波外置：5 文件清零摘条＋1 文件降余 2 枚；registry 55→50），实现腿（GLM 5.3 Flash）。只做外置；扫描器逻辑、既有 cap 数字、AC#1 勾选框均未动。**
  - **选样（按派单规则"产品面枚数最多 4-6 个"）**：config-screen 49、cost-settings 49、plugins-screen 40、privacy-settings 40、palette-screen 37、history-page 31（六文件合计 246 枚，占产品面 backlog 587 的 42%）。逐一核对全部 render 仪器的**源码级钉**（readFileSync(...).includes）后才动手：config-screen 的 BADGE_WORDS 三词、cost 的 5 词、privacy 的 3 词、history 的 1 词被仪器钉在源码——第一波的既有先例是"copy 迁走后源码钉由注释里的针词承载"（config-feedback.ts 的 `需要重启` 钉即如此），本波沿用，仪器零放宽。
  - **清零摘条 5 个（k=5）**：`config-screen.tsx` 49→**0**、`cost-settings.tsx` 49→**0**、`privacy-settings.tsx` 40→**0**、`palette-screen.tsx` 37→**0**、`history-page.tsx` 31→**0**，全部改 `t(key, params?)`，中文搬进 `lib/locale/zh.ts`（wave2 段新增 ~184 键）、英文搬进 `en.ts`（`Record<LocaleKey,string>` 类型强制同键）。REGISTRY 每清一个整条摘掉（摘条注释留痕），stale-entry 牙全程没有绕过——每归零一个就先吃一条"stale registry entry"红再摘。
  - **降不清零 1 个**：`plugins-screen.tsx` 40→**2**，条目保留、cap 40 未动。余 2 枚＝`stateTone()` 里对 Go 名册三态词的**匹配令牌**（`state === "本轮生效"` / `=== "需重启"`，匹配用不上屏——上屏的 state 文本始终走 wire 原样），与 `config-protocol.ts` 同一性质；protocol 白名单 cap=23 已满、⛔ 不许给任何文件加/涨上限，故令牌留在源文件、理由写进函数注释。
  - **仪器源码钉随迁 4 台（第一波遗留账，非本波新红）**：本波首次基线即发现 `render:nav` / `render:pluginpage407` / `render:palette38` / `render:palette38kbd` 四台 **rc=1**——根因是第一波把 copy 迁进 locale 时这四台的源码级钉没跟着真值走（`aria-label="命令面板"` 已变 `aria-label={PALETTE_LABEL}`＝t()、config-feedback 翻译函数体已改用 GO_REFUSED 常量等）。仪器自己写着迁移规则「the consume waves move these pins with the truth, never loosen them」——按此把四台钉迁到新居所，且**每条改成都断言两半**（App 接线半 `t("app.*")` / GO_REFUSED ＋ zh.ts/config-protocol.ts 句子半逐字中文），牙比原钉更咬、零放宽。另：本波险些二次踩钉——history-page 的仪器有负向钉 `!src.includes('"history.`，而我最初起的键名前缀 `t("history.` 会撞上它，键名改 `hist.*` 避开（不是绕牙：钉禁的是历史方法名硬编码，键名撞车是假阳性）。
  - **自证（全部现跑）**：`npx tsc -b` rc=0；`npm run build` rc=0；`node scripts/lint-cjk.mjs` ⇒ closure files=100 scanned=93 **registry entries=50 findings: 0**；`node scripts/lint-cjk.mjs --self-test` ⇒ 6 路全红照旧（清洁样本绿＋种字面量/种 JSX/cap 涨/stale 条/真盘临时树全红）rc=0；`npm run lint`（oxlint）0 errors；**12 台涉事 render 仪器全 rc=0**（animation193/nav/security40/privacy40/pluginpage407/palette38/palette38kbd/forcechip38/discovery38/historypage36/composer/config-screen）。
  - **普查读数（cjk-census 现跑）**：closure **1514→1454 枚**（字面量 1243＋JSX 文本行 211）；**zh.ts 470 枚＝中文唯一合法居所**（第一波 286→+184）。registry 余 50 条（演示面 14 文件 618 枚未动、协议 23 枚未动）合计 **961 枚**，其中产品面 backlog 34 文件 **343 枚**（587−244）。下波最大余量：diagnostics-settings 30、config-editor 25、security-settings 25、prompt-bar 23。
  - AC#1 未勾（勾归非实现者）。
