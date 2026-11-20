# 票 193 · 界面半（AC#5/AC#4 UI 侧）证据 · r1

- 腿：193-ui-finish（接管 2026-11-19 被打断腿的盘上成果并交件；本腿只做自查、四发变异取证、证据件与提交）
- 日期：2026-11-20（会话时间 Asia/Shanghai 20:0x–20:2x）
- 执行：接管写腿（agent-158）；**本证据件出自实现腿本身，按 SPEC-12 §4.3 仍需另一 agent 做对抗验收**

## 0. 基线现跑（接管后、变异前，全部本腿亲手复跑）

| 门 | 命令 | 读数 |
|---|---|---|
| typecheck | `npm run typecheck`（frontend/） | rc=0（tsc -b 无输出） |
| build | `npm run build`（frontend/） | rc=0（✓ built in 1.54s，仅 chunk 体积警告） |
| 仪器 | `npm run render:animation193`（frontend/） | rc=0，OK 句含 "off 真停: getAnimations()=0 + zero active rAF handles, in-flight cancelled on the flip, starts declined" |
| go test | `go test ./internal/panel -count=1` | ok 20.722s（本腿复跑；编排者先前复跑 17.630s） |
| d22scan | `sh scripts/d22scan.sh` | `d22scan.sh: OK - control rc=0, scan rc=0`，0 findings（末行 OK 句） |

## 1. AC#5 逐条自查（派单第 1 条）

1. **tag 逐字一致**：`internal/panel/panel_animation_193.go` 的 `PanelAnimationSection`/`Tier string \`json:"tier"\``＋`composer.go` 的 `PanelAnimation *PanelAnimationSection \`json:"panelAnimation,omitempty"\``，对照 `frontend/src/lib/panel.ts` 的 `panelAnimation?: PanelAnimationSection`＋`export interface PanelAnimationSection { tier: string }`。键名与形状逐字一致（见 §3 对照表）。仪器另钉：`ANIMATION_TIERS` 与 `internal/config/validate.go` 的 `panelAnimationVocabulary` 逐字同序同词（仪器直接读 Go 源比对）。
2. **三档 UI＋默认档标**：`config-screen.tsx` 的 panel.animation 行挂三档 SegmentedControl（全开／减少／关闭，label 表 `ANIMATION_TIER_LABELS` 与词表一一对应，无第四档）；行描述逐字含「「减少」档为默认（默认，待定）——默认档位待 owner 拍板」。仪器判读钉死 `「减少」档为默认（默认，待定）` 出现在 unknown 渲染里。
3. **写走 config.set、字段与 Go 同名**：写半只有一条路——`performPanelAnimationWrite`（`config-actions.ts`）经 `requestConfigWrite({ field: "panel_animation", ... })`，`panel_animation` 即 Go 侧第九名册拼写（`FieldPanelAnimation = "panel_animation"`），无第二条通道；`config-screen.tsx` 的 RowSource 配对 `field: "panel_animation", section: "panelAnimation"` 被两台仪器（render-config-screen、render-animation-193）源码判读钉死。
4. **读面没发键 ⇒ UI 显未知**：`config-screen.tsx` 三分支——`panelAnimation == null` ⇒ 显「档位未知（快照未携带动画档）」且不挂任何控件；`tier === ""` ⇒ 「档位未设置（读面为空档位）」；词表外词 ⇒ 「档位未知（读面词不在三档词表内）」。`lib/motion.ts` 的 `applyPanelAnimationTier` 对未知形状撤 `data-motion` 属性回出厂默认。绝不默认成 off／full。

## 2. AC#4 四枚牙变异取证（逐发实际改→红→还原→绿）

变异底本存 `.scratch/wisp/probes/orch/20261120/193-anim-mutations/`（pristine 副本），每发跑完从底本逐字节还原并复跑仪器确认回绿。仪器＝`npm run render:animation193`（frontend/）。

| 牙 | 变异（实际改了什么） | 红读数（逐字） | rc | 还原后 |
|---|---|---|---|---|
| ① off 档仍注册 rAF ⇒ 红 | `frontend/src/lib/panel-visibility.ts`：删 `createPanelAnimationLoop.start()` 里的档位门 `if (!motionHeavyAllowedOn(doc)) return;` | RED, 3 failed: `AC#4 off 真停: the loop factory must gate start() on the tier`；`牙1: a start() under off must register NOTHING and stay stopped`；`AC#4 reduce 边界: a start() under reduce must decline too` | 1 | OK, rc=0 |
| ② 把 off 当 prefers-reduced-motion ⇒ 红 | `frontend/src/styles/theme.css`：把 `[data-motion="off"]` 规则嵌进 `@media (prefers-reduced-motion: reduce)` 块 | RED, 1 failed: `牙2 静态半: the off rules must sit at the stylesheet's top level (not inside any @media block)` `[brace depth at the off rules = 1]` | 1 | OK, rc=0 |
| ③ 读面缺键时默认成 off ⇒ 红 | `frontend/src/lib/motion.ts`：`applyPanelAnimationTier` 的未知分支由 `root.removeAttribute(MOTION_ATTR)` 改为 `applyMotionTier("off")` | RED, 11 failed: `AC#4 未知: ... must withdraw the attribute`＋牙3 页面半逐形状红（undefined/null/""/loud/OFF 各 2 条：撤属性＋零取消零注册） | 1 | OK, rc=0 |
| ④ getAnimations() 非 0 时仍判停 ⇒ 红 | `frontend/src/lib/motion.ts`：删翻 off 档时的在跑取消 `for (const a of live) a.cancel();`（保留 getAnimations 读取） | RED, 3 failed: `AC#4 off 真停: the flip to off must cancel what is already in flight`；`AC#4 off 真停: getAnimations() must be EMPTY after the flip`；`幂等: a second push of the same tier must be byte-identical on the root` | 1 | OK, rc=0 |

四发全部：改→红（rc=1）→还原→绿（rc=0）。还原后三个文件 `git diff --stat` 与变异前完全一致（panel-visibility.ts +30/-2 行档、motion.ts +161/…、theme.css +66/…，即上一腿原 diff，无残余改动）。

## 3. 与 Go 载体的 tag 对照表

| Go（internal/panel，已提交载体） | json tag | TS（frontend/src/lib/panel.ts，本腿盘上成果） | 一致 |
|---|---|---|---|
| `composer.go`: `PanelAnimation *PanelAnimationSection` | `panelAnimation,omitempty` | `PanelSnapshot.panelAnimation?: PanelAnimationSection` | 逐字（optionality 语义同向：Go omitempty 不发键 ⇔ TS 可选缺省 ⇒ UI 显未知） |
| `panel_animation_193.go`: `Tier string` | `tier` | `PanelAnimationSection.tier: string` | 逐字 |
| `validate.go`: `panelAnimationVocabulary` = full/reduce/off | — | `ANIMATION_TIERS = ["full","reduce","off"]`（仪器读 Go 源逐字比对） | 同词同序 |
| 第九名册 `FieldPanelAnimation` | — | 写半 `field: "panel_animation"`（config.set 唯一门） | 逐字 |

## 4. 默认档位：待 owner 拍

**动画默认档位待 owner 拍板。** schema 占位默认为 `reduce`，设置页仅标「「减少」档为默认（默认，待定）」占位，不当作定案；owner 拍板后改 schema 默认值＋该行文案即可，UI 机制不变。

## 5. 边界与不动项

- `frontend/fixtures/composer-states.html`：进场前即为 M（他人在飞件），本腿未触碰、未提交。
- `frontend/dist`：未提交任何字节（本腿未产出构建产物入库；dist 仍只有既有内容）。
- `cmd/wisp/**`、`internal/**` 产码：零改动（仅 `internal/panel/composer_test.go` 按派单保持上一腿的账本收口登记，无断言放宽）。
- thresholds/golden/allowlist/PLAN/specs/D43/design：零改动。
- 提交用显式 pathspec；无 add -A／amend／reset／rebase／stash／checkout；未 push。
