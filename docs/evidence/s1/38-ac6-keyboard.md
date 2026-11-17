# 38 AC#6 — 键盘全流（keyboard-only palette flow）裁决表

- **日期**：2026-11-17 10:4x +08
- **实现腿**：agent-106（GLM 舰队，铁律 4 前端归本编队；写面按派单：`palette-screen.tsx`、`App.tsx`、`render-palette-kbd-38.tsx`、本文档）
- **票**：`.scratch/wisp/issues/38-command-palette-tasks.md` AC#6「Keyboard-only: full palette flow reachable without mouse (tab/arrow/enter audit)」
- **基线**：`1d5d686733eb08b359772a5ed8a32546e33e9c95`（AC#1/AC#2 已在 `817785260` 收件勾框；本腿不动那两格）

## 改了什么（机制，三处）

1. `frontend/src/components/palette-screen.tsx`
   - 键盘词汇收成**一个纯决策**并导出：`paletteKeyDecision(key, blocked, rowCount)` ＋
     `paletteCursorStep(cursor, action, rowCount)`。组件的 `onKey`（输入框）与新增的
     `onSurfaceKey`（行列表）跑的是**同一个**决策——仪器驱动的就是组件跑的代码路径，
     不是仪器里另抄一份。
   - `blocked = composing || e.nativeEvent.isComposing`（AC#2 守卫原样保留，一字未放宽）。
   - **Tab 不被面板吞**：决策对 Tab 返回 `none` 且不 preventDefault ⇒ 浏览器原生 DOM 次序
     Tab 走：搜索框 → 清空按钮（有查询时）→ 行按钮（行本身就是行内操作，逐行是
     `<button>`）按 `palette-row-N` 的 N 次序推进。
   - **焦点落行、高亮随行**：`onFocusCapture` 把落到某行上的焦点同步成键盘游标（`setActive`），
     所以 Tab 到某行后按 Enter 选的是**那一行**，不是指针留下的旧游标。
   - **↑↓ 在行上也能用**：`onSurfaceKey` 接住落在行按钮上的方向键，移动高亮**并把焦点
     送回输入框**（aria-activedescendant 住在输入框上，播报行恒等于 Enter 落点行）；
     输入框里按键由输入框自己的处理器负责（`e.target === inputRef.current` 即返回），
     不会一步走两行。
   - **aria**：行列表包 `role="listbox"`（`aria-label="命令列表"`），分组包 `role="group"`，
     行原有 `role="option"`＋`aria-selected` 不动；输入框 `aria-activedescendant` 随游标移
     （空结果集不输出该属性）。
2. `frontend/src/App.tsx`
   - **ESC/Ctrl+K 补 composing 守卫**：document 级 handler 最前面 `if (event.isComposing) return;`
     ——组合期 ESC/Ctrl+K 全部不触发（AC#2 的守卫从输入框扩到对话框动词，没有削弱）。
   - **Tab 焦点陷阱（focus trap）**：面板开着时 Tab 在 `[role="dialog"][aria-label="命令面板"]`
     内取 `input, button:not([disabled])` 按 DOM 次序走，两端环绕（末行 Tab 回搜索框、
     搜索框 Shift+Tab 回末行），跑出对话框外的游离焦点被拉回第一个可聚焦元素。
   - **关闭后焦点回原处**：打开的两条路（Ctrl+K、头部按钮 onClick）都先
     `rememberPaletteReturnFocus()` 记下 `document.activeElement`；关闭（ESC / Ctrl+K 再按 /
     点遮罩 / 行派发任何一条关闭路径）统一走 effect 把焦点还给原元素。
3. `frontend/package.json` ＋ 新仪器 `frontend/scripts/render-palette-kbd-38.tsx`：
   新脚本 `render:palette38kbd`（与既有 `render:palette38` 分工——那只管线形状，这只只管键盘流）。

## 五支牙与读数（`npm --prefix frontend run render:palette38kbd`，rc=0）

| # | 牙名（仪器 check） | 读数 |
|---|---|---|
| 1 | Tab 次序（正控，可判读 DOM 次序） | `<input` 在全部行 id 之前；`palette-row-0 < 1 < 2` 按升序出现；`paletteKeyDecision("Tab",false,3)==="none"`（面板不吞 Tab，原生走序保留） |
| 2 | ↑↓ 移行（正控＝焦点 id 序列） | Down Down Up Up 走出 `palette-row-0 → 1 → 2 → 1 → 0`；`paletteCursorStep(0,up,3)==2`、`(2,down,3)==0` 双向环绕不脱 |
| 3 | Enter 选行（正控） | 决策返回 `select`；组件源含 `onSelect?.(flat[cursor])`；`aria-activedescendant="palette-row-0"` 与 `aria-selected="true"` 是**同一行** |
| 4 | ESC 关闭 | 面板决策对 Escape 返回 `none`（不吞）；App 源 `event.key === "Escape"` → `paletteOpen: false`；负控：`paletteKeyDecision("Escape",true,3)==="none"` |
| 5 | Ctrl+K 开关 | App 源 `(event.ctrlKey||event.metaKey) && key.toLowerCase()==="k"` 拨 `paletteOpen`；守卫 `if (event.isComposing) return;` 在它**之前**；trap/回焦源级三件（dialog 查询、双端环绕、`back.focus()`）齐 |

**负控（总）**：`composing=true` 时 `ArrowDown/ArrowUp/Enter/Escape/Tab/k` 六键决策**全部 `none`**
（AC#2 守卫未弄坏）；空结果集渲染**无** listbox、**无** activedescendant、**无**行 id。
**正控（总）**：三行开着时上面整表逐条可判读通过。

## 门禁读数（串行，node=D:\work\server\node24 v24.21.0）

| 门 | rc |
|---|---|
| `npm --prefix frontend run build`（tsc -b + vite build） | **0** |
| `npm --prefix frontend run typecheck` | **0** |
| `npm --prefix frontend run lint`（oxlint，40 warnings 0 errors，既有基线） | **0** |
| `render:palette38kbd`（新牙） | **0** |
| `render:palette38`（既有牙回归） | **0** |
| render 回归批：nav / approval37 / history36 / shell77 / receiver / stream / composer / model-settings / roster | 全 **0** |

## 既有红（非本腿写面，如实转报）

- `render:security40` rc=1（1 check：`the grants absence must name the nil grant_id writer`，
  断言 security 页 markup 含 `internal/agent/journal.go`）。输入面是 settings/security 组件，
  本腿零触碰——基线就红。
- `tools/d22scan` 5 枚红全在 Go 面（`cmd/wisp/approval_reply.go:379/:464`、ticket49 测试三处），
  与上一腿转报的是同一批；frontend/ 零命中，本腿 markup 仪器自带零 emoji 断言过。
- 本波回归跑 `render:composer` 把 HEAD 上**过期**的 `frontend/fixtures/composer-states.html`
  重生成出 4 行 diff（上一 composer 波没回写 fixture）；已按 HEAD 字节还原、未提交。
  修复归 composer 波：重跑 `render:composer` 把再生成的 fixture 一并 commit。
  （还原后该文件 status 上的 ` M` 是 autocrlf 行尾幻影，`git diff` 为空。）
- `render:l2` 未跑（其裸跑 rc=2 是既有 usage 退出，见票 38 收件批注）。

## 边界与不造值

- 名册行仍是**只读**（`commands_roster_213.go:21-26` 只承载不执行）——键盘全流不改变这一条：
  名册区没有可聚焦元素，Tab 次序里的「行内操作」就是行按钮本身，App 的
  `runPaletteItem` 是唯一执行面（`sendMessage(body, [])`，同一封套，无旁路）。
- 焦点陷阱与回焦运行在 `document.activeElement` / keydown 监听里，静态 markup 不可点，
  仪器按本树既有 provenance 在 App **源级**断言（机制串逐字），行为读数由纯决策＋markup 承担。
- dist 构建产物未入库（`git ls-files frontend/dist` 仍只有 `.gitkeep`，构建输出走 gitignore）。
- AC#6 勾框留给非实现者（收件裁决者勾）。
