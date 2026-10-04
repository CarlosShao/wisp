# 260-r3 证据件 — 票 260 条件②剩下的 10 枚「Esc」人看文案（具名解冻 `A595`）

> 本节标题就是骨架时刻（2026-10-04 10:5x +08）落的第一枚 commit；§1／§2 的读数是同发现跑抄进来的，
> §3–§9 在下面第 60 轮之前填满（门禁读数与判不动的地方两节优先）。
> 授权来源＝`docs/reports/pending-and-issues.md` **A595** §1 的五样（文件／行／理由／边界／撤销口令）。
> 票面＝`.scratch/wisp/issues/260-configured-cancel-hotkey-never-registers.md` §编排者选形裁定（A588）条件②。

## 1. 起手读数与锚点（现跑）

- 时刻：`2026-10-04 10:44 +0800`（`date "+%Y-%m-%d %H:%M %z"`）；起手 HEAD＝`57a33804`（`git rev-parse --short HEAD`，分支 `dev`）。
- 尺①（写面一）：`grep -nE '"[^"]*Esc[^"]*"' internal/agent/approval/approval.go` ⇒ **2 枚**，逐字：
  - `:90` `	ChannelEsc:   "按 Esc 键",`
  - `:112` `		return "Esc 取消不可用"`
- 尺②（写面二）：`grep -nE '"[^"]*Esc[^"]*"' cmd/wisp/resident_approval_windows.go` ⇒ **8 枚**，行号 `271/279/287/309/311/495/584/586`，逐字见 §4 表。
  ⇒ 复跑行号与 `A595` §1 记的号**一致**（编排者那一趟是 `10:19:36`／锚 `84dbda52`；本腿这一趟是 `10:44`／锚 `57a33804`，中间落的是 `260-r2` 的 ball 侧件，未动这两枚文件的行距）。
- 口径：只数「字符串字面量内部含 `Esc`」的行；注释与标识符（`ChannelEsc`／`vetoByEsc`／`TakeEscForCancel`／`escLoad`）不计。
- 编队起手（`git status --porcelain internal cmd tools scripts .github docs/evidence`）：见 §9 的并发自查。

## 2. 授权面与射程核对

- `A595` §1 点名的写面＝上面两枚文件；`A595` §1 边界①②③④逐条抄录并逐条对账（本节在下面补对账结论）：
  - ① 只改指代那枚键的名词，不动派发／加载／借还的任何一行逻辑；
  - ② 不碰 `internal/agent/approval/gate.go`、不碰 `internal/ball/**`；
  - ③ 票 245 的时机与存在性一字不动；
  - ④ `:271`／`:279`／`:112` 三枚的诚实形是「取消键无处可借／不可用」，不是换成另一枚键名。
- 撤销口令：「260 文案撤回」。

## 3. 取值链（单一真相源怎么走）

本节写：`ConfiguredHotkeys().Cancel` → `resolveCancelBorrow` → `cancelBorrowedLineFor`/`cancelIdleLine` → `HotkeyReport` 的 `cancel` 那一行 `Binding` → `cmd/wisp` 的读口；
以及为什么 `cmd/wisp` 不自己拼组合键、为什么空串落点是 `ball.DefaultHotkeys().Cancel`。逐处 `file:line`。

## 4. 十枚逐句分类（三档）＋改前／改后逐字

本节写：三档各几枚、每枚的旧句逐字、新句逐字、以及"这句在说什么事实"。⛔ 不是一句"统一替换"。

## 5. 判据（两格）与红句逐字

本节写：① 默认档零漂移（逐字等旧句）；② 种一发别的组合键 ⇒ 句中出现那枚新键、且不再出现 `Esc`；两格各自的用例名、断言体、以及②的红句逐字。

## 6. winlive 那枚钉的同步（`vetoDoneClaim`）

本节写：`cmd/wisp/resident_task_source_live_246_windows_test.go:69` 的改前／改后逐字，台架 `prepareResidentHarness` 写的 `config.toml` **不含 `[hotkey]`**（⇒ 台架跑默认档），
以及为什么这一格只能写成**判据**〔预测红/预测绿〕而不是"我验证过了"（CI 无 `winlive` tag、本机 `0xc0000135` 无 `--- FAIL`）。

## 7. 门禁终值读数（逐名）

本节写（第 60 轮前先填满）：`go build ./...`／`go vet ./cmd/wisp/`／`go vet -tags winlive ./cmd/wisp/`／带 PATH 的定向 `go test -v -run …  ./cmd/wisp/`／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh`／`gofumpt -l`。

## 8. 判不动的地方／量不到的地方（含 `:90` 的停手上报）

本节写：`:90` 为什么在本腿射程内做不到"默认档零漂移＋吃配置"两者同时成立（`gate.go:314`／`gate.go:477`／`report.go:142` 直接读 `channelNames` map ⇒ 不改这三枚就没有线程安全的注入面；改这三枚越出 `A595` §1 的写面），三形候选与代价，⛔ 本腿不自行选形。

## 9. 我越界了什么＋遗留与欠账

本节写：逐枚 diff 尺（`git diff --numstat`）、有没有第 11 枚、有没有动逻辑行、并发卷尘自查、以及交给 `260-v1` 的格子。
