# 票 407 审计收尾记录 — 可疑①rows 行级守卫 ＋ 可疑②拒绝 fixture 从 Go 现读

**日期**：2026-11-19 18:0x–18:2x+08。**来路**：派单 `.scratch/wisp/probes/orch/20261119/txt/407-audit-followup-dispatch.txt`（audit407 对 `1d08513ec` 的 7 通过/1 可疑/0 不通过收尾，本编队前端写腿）。
**性质**：实现腿自录证据（裁决/对抗验收由非实现者另出）。
**裁定执行**：只动 `frontend/**`；`internal/panel/**` 与 `cmd/wisp/**` 产码**零触碰**（下文 §3：fixture 现读方案**不需要**新增 Go 常量导出）。

## 1. 可疑①：rows 行级守卫＋未知分支（plugins-screen.tsx）

审计读数：`plugins-screen.tsx:146` 直读 `section.rows.length`，读面缺 `rows` 键 ⇒ `Cannot read properties of undefined (reading 'length')`（生产被 `internal/panel/pump.go:703-705` 的 nil→`[]PluginRow{}` 归一掩住；页面吃 parsed JSON，不能押注宿主永远良形）。

修法（`frontend/src/components/plugins-screen.tsx`）：

- 新增 `rowsGiven(rows: unknown): rows is PluginRowView[]`（`Array.isArray`）——rows 与格字段同一条未知规则：TS 接口是良形描述，页面消费的是 parsed JSON。
- 分支三态（⛔ 不改任何契约文字，行为即语义）：
  - `rows` **缺席／null／非数组** ⇒ 该页显示「**未知（读面没有给清单）**」＋解释句（这不是「还没有插件」：没有清单可读不等于清单为空），**零表格**，不再抛异常；
  - `rows === []`（给了的空清单）⇒ 仍显「**还没有插件**」＋原实测说明句（正控，不许退化成未知）；
  - `rows` 非空数组 ⇒ 清单表，行为不变。
- 组件头注释补 THE MISSING ROWS LIST 段（注释非契约，只是把规则写清）。

## 2. 可疑②：拒绝 fixture 从 Go 载体现读 ＋ 折保存失败断言（render-pluginpage407.tsx）

审计读数：旧 fixture 手抄（`:155-156`），缺 `panel: 配置字段不在可写名册内:` 前缀（`ErrConfigFieldRefused` 的 `%w` 段）；且翻译器下游若把拒绝折成「保存失败」，旧 fixture 上的仪器仍全绿——两处都修：

1. **fixture 现读装配（tooth 3b）**：`REFUSAL_REPLY` 不再手抄，仪器启动时从三个 Go 文件行现读字面量、按 wire 的拼装顺序组装：
   - 闸门句：`internal/panel/config_handlers.go` 的 `fmt.Errorf` 字面量（`必须指名插件 id（configPlugin＝[plugins.<id>] 的 <id>；requestId=%q）`），两枚 `%q` 按字段序填 `plugin_enabled`／`r-407-1`；
   - `%w` 载体：同文件 `ErrConfigFieldRefused = errors.New("panel: 配置字段不在可写名册内")` 的字面量；
   - 信封前缀：`internal/panel/bridge.go` `RefusedEnvelopeForUser` 的 `fmt.Sprintf("面板请求被拒绝 [%s %s]：%s", ...)`，按 wire 同序填 method/requestId/原因。
   - 装配后自带一枚 check（闸门句＋configPlugin 选择器词＋前缀＋名册句四要素齐才过）；抽行工具（`goLine`/`goStringLiteral`）找不到 Go 行**直接 throw**——缺行＝仪器崩，不是静默绿。
   - **未新增任何 Go 侧常量导出**（readFileSync 现读文件字面量已足），internal/panel 产码零触碰——派单允许的「新增一枚只读导出」未动用。
2. **折保存失败断言（tooth 3c）**：在等价返回点（`performPluginToggleWrite` 的返回值＝`translatePluginToggleReply` 的产物，仪器牙 2/牙 3(source) 源码尺钉住这条链）断言：**`tone==="error"` 时 `text ≠ 「保存失败」`**——检测器 `foldsIntoGenericSaveFailure(f)` 吃 Feedback，对拒绝/已接受/重启三个回执全过；并自带一发**自证红**（喂 `{tone:"error",text:"保存失败"}` 必须命中），保证检测器不是恒过。

## 3. 新牙读数（派单要的四发＋变异演示）

四发新牙（仪器 `npm run render:pluginpage407`，rc=0，全部绿）：

| 发 | 臂 | 读数 |
|---|---|---|
| ① | `rows` 键删（JSON round-trip＋delete，真缺席） | 显示「未知（读面没有给清单）」、零 `<table>`、零行泄漏、无空清单句，**不抛** |
| ② | `rows: null` | 同①（null 不是空清单） |
| ③ | `rows: []`（正控） | 仍显「还没有插件」＋零表格，**不退化成未知**（断言同时禁未知句出现在空态 markup） |
| ④ | 折保存失败检测（tooth 3c） | 拒绝回执 `tone=error` 且 text＝宿主原话（非「保存失败」）过；检测器自证可红 |

变异演示（临时改源→跑仪器→抓红→还原，盘上字节已复原；每次还原后复跑仪器回绿 rc=0）：

| 变异 | 改动 | 读数 |
|---|---|---|
| M1 | 页面 `rowsGiven` 改恒真（守卫被拆） | 仪器**崩红** rc=1，`TypeError: Cannot read properties of undefined (reading 'length')`——正是审计报的崩形，rows 缺/null 两臂咬住 |
| M2 | 拒绝 fixture 改回手抄折断（少前缀＋缺闸门后缀，旧形状） | `FAIL tooth 3: the missing-configPlugin refusal rides VERBATIM ... [tone=error]`，rc=1 |
| M3 | 翻译器下游折断：`config-feedback.ts` 对带名册前缀的拒绝折成 `{tone:"error",text:"保存失败"}`（审计说的躲绿形状） | **3 处 FAIL**：牙 3 逐字＋牙 3 上屏＋**新牙 3c** `[refusal text head: 保存失败]`，rc=1——旧 fixture 缺前缀时 M3 躲绿，新 fixture 现读全文后必咬 |

M3 单独验证了审计的可疑点本身：折断在旧 fixture 下仪器全绿，现在三枚牙（3/3b/3c）同时咬。

## 4. 现跑读数（2026-11-19 18:0x–18:2x+08，全部本腿现跑）

| 闸 | 命令 | 读数 |
|---|---|---|
| 仪器（含新牙） | `npm run render:pluginpage407` | **rc=0**，OK 句末行含 rows 未知规则＋fixture 现读＋折断检测三段新文 |
| typecheck | `npm run typecheck`（tsc -b） | **rc=0** |
| build | `npm run build` | **rc=0**（1.87s；chunk>500kB 警告＝存量） |
| Go 门 | `go test ./internal/panel -run '^TestFrontendHasNoEmoji$'` | **ok ... 0.552s，rc=0** |
| 扫描 | `sh scripts/d22scan.sh` | 末行照录：`d22scan.sh: OK - control rc=0, scan rc=0`（clean 行：`d22scan: clean - no D22 ban violations; ...`，**0 findings**） |

本腿不跑 `go test ./cmd/wisp`（未碰 cmd/wisp，遵派单串行）。

## 5. 交付面（全部在派单写面内）

| 件 | 变更 |
|---|---|
| `frontend/src/components/plugins-screen.tsx` | rows 行级守卫＋未知分支（§1） |
| `frontend/scripts/render-pluginpage407.tsx` | tooth 3b fixture 现读装配＋tooth 3c 折保存失败断言＋rows 三臂牙＋OK 句扩写（§2/§3） |
| `docs/evidence/s1/407-plugin-panel-markup-empty.html` / `-with-plugins.html` | 仪器重产出（证据快照，非 golden，无闸门引用） |
| `docs/evidence/s1/407-plugin-panel-audit-followup.md` | 本件 |

**没动的**：`cmd/wisp/**`、`internal/panel/**`（含 pump.go 的掩蔽点——那是读口的既有归一，页面侧守卫按裁定补在消费端）、`frontend/src/lib/config-feedback.ts`／`config-actions.ts` 产码（M3 变异只在演示中临时改过、已还原，还原后仪器/typecheck 双绿坐实）、thresholds/golden/allowlist/PLAN/specs/D43/design、他人在飞件（`frontend/fixtures/composer-states.html` 等）。
**放宽断言：无**（新检查全为收紧；唯一一处分档是把「禁『还没有插件』字样出现」精确到空态分支的独有形状——`>还没有插件</p>` 标题元素＋其解释句——因为未知分支的**解释文字**本身合法地提及「还没有插件」，首版断言被自己的解释句咬中，属检查精度修正，非放宽）。
