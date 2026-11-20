# 票 191（界面半）收件 · terminal 区上屏 + 账本收口 + 仪器

- 日期：2026-11-20
- 写腿：本腿（frontend/** ＋ internal/panel/composer_test.go 账本收口）
- 载体来源：`8513bc2e3`（wave 2a 面板侧，`internal/panel/terminal_carrier_191.go` → `Snapshot.Terminal`，json `"terminal"`：`input`/`output` 各一枚 `{text, truncated, droppedRunes}`；没到过就缺席，空态＝`"terminal": {}`）

## 交付

1. **panel.ts 声明**（`frontend/src/lib/panel.ts`）：`PanelSnapshot.terminal?: TerminalSection`；`TerminalSection{input?; output?}` 两键各指 `TerminalWindow{text; truncated; droppedRunes}`，与 Go 侧 json tag 逐字一致（行内三键无 omitempty，页面三个键全必填；两枚 Go 行结构 TerminalInputRow / TerminalOutputRow 共用这一页形）。
2. **账本收口**（`internal/panel/composer_test.go`，同一提交）：`composerTSGapLedger` 的 `"PanelSnapshot": {"terminal"}` 登记删除（收口注释按 193/4017962b3 先例留档）；`composerTSNameMap` 加三条翻译：`TerminalSection→TerminalSection`、`TerminalInputRow→TerminalWindow`、`TerminalOutputRow→TerminalWindow`。⛔ 未放宽任何断言——两向对账、stale/dead 判读全部原样生效。
3. **UI**（`frontend/src/components/terminal-region.tsx`，挂点 `frontend/src/App.tsx` 对话列、composer 上方）：一个可折叠「终端」区；上方只读命令行（input）、下方滚动尾窗（output，追尾随最新行）；`truncated=true` 显示「已截断，另有 N 字符未显示」（N＝droppedRunes 原数）；`terminal` 键缺席显示「终端未连接」＋缺键原因，绝不空白；空态（`"terminal":{}`）显示 尚无输入/尚无输出，与缺席两态不共用一张脸。区内有且仅有一个按钮＝折叠钮（视图偏好，随文档消亡）；零桥接动词、零 emoji、零颜色字面量。
4. **仪器** `render:terminal191`（`frontend/scripts/render-terminal-191.tsx`，`npm run render:terminal191`）：三态渲染判读＋三牙（见下读数），另钉：panel.ts 声明逐字、App 喂 `snapshot.terminal` 本尊（喂而不读＝假接线）、单按钮＋零桥接动词、零 emoji 零颜色字面量。

## 三发牙读数（仪器内判读名，全绿）

1. 牙1（缺席显「终端未连接」⇒改空白⇒红）：`a pump without the terminal key must read 终端未连接, never blank (牙1: 删掉这句话 = 红)` 绿；反面 `the EMPTY section ... must NOT read 终端未连接` 绿。
2. 牙2（truncated=true 不显截断提示⇒红）：`truncated=true must carry the truncation notice verbatim` 绿（整句「已截断，另有 123 字符未显示」逐字在屏）；反面：未截断节零提示 绿（live 渲染「已截断」恰 1 次、clean 节 0 次）。
3. 牙3（droppedRunes 显示错⇒红）：`a second window with droppedRunes=7 must show 7` 绿 ＋ `must NOT show the other window's 123` 绿（数字来自行，非常量）；输入行自己的 `4097` 同样逐字上屏 绿。

## 门禁读数（全部本腿现跑，2026-11-20 22:0x–22:2x +08）

| 门禁 | 读数 | rc |
|---|---|---|
| `npm run typecheck` | tsc -b 无错 | 0 |
| `npm run build` | vite build ok（2733 modules，产物在 frontend/dist，**未提交**） | 0 |
| `npm run render:terminal191`（新仪器） | OK，全判读绿（读数见上） | 0 |
| `go test ./internal/panel -count=1` | RUN=306 PASS=306 FAIL=0 SKIP=0（ok 23.726s） | 0 |
| `sh scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations`；末行 `d22scan.sh: OK - control rc=0, scan rc=0` | 0 |

仪器回归扫描（渲染整版 App 的既有仪器，防新挂点破形）：render:nav / render:snapshot-receiver / render:task-roster / render:history36 / render:palette38 / render:approval37 / render:shell77 / render:hiddenpause36 / render:animation193 全部 rc=0（本腿收尾时附最终清单）。

## 边界与不动项

- ⛔ 未动：`cmd/wisp/**`、`internal/tools/**`、`internal/risk/**`、`internal/panel` 产码（本腿只改 `composer_test.go` 账本）、thresholds/golden/allowlist；⛔ 未 push；⛔ 未提交 frontend/dist 字节（构建产物只落盘，d22scan 亦按 gitignore 跳过）。
- ⚠ 他人在飞件：`frontend/fixtures/composer-states.html`（M，工作树未提交改动）——**只报不动**，本腿未读写该文件。
- 终端区本波只读：不画任何可执行按钮；命令本身走与一切命令相同的门控工具腿（AC#0 甲裁定），面板不新增执行通路。泵源接线（wave 2b，cmd/wisp）与泄漏半（Q-457 的 C25 管道）不在本腿。
