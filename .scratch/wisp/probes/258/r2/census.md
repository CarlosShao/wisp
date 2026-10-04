# 票 258 · 实现腿 `258-r2` 证据件（写码腿，非验收者）

起手时刻 `2026-10-04 09:0x +0800`（自取 `date`）；起手锚 HEAD `fd269de1`（`git log --oneline -1` 现取）。
本件只建不删；读数逐字；AC 框一枚未碰。

**本腿射程（派单三件事）**：
① AC#1 条件①：`[hotkey]` 节缺失时档位词必须印 "defaults"（措辞说真话），并配一枚不依赖真窗、
   对「把 defaults 改回 config」这一发有红的用例。
② AC#2 条件：两枚 winlive 尺子修形状（rebind 枚的 JSON/等号格式失配、occupied 枚的环境赌博），
   改成"读得到就判、读不到就具名失败/诚实 skip"。⛔ 本轮不跑 `-tags winlive`。
③ build 销账：全量 `go build ./...` 逐字输出＋rc 落 §3。

## §0 起手基线与绿名册

- 起手 HEAD：`fd269de1`。
- 基线整包 `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v ./cmd/wisp/`：读数档
  `.scratch/wisp/probes/258/r2/baseline-gotest.txt`（跑完后把三枚数与含 hotkey/Hotkey/ball 的绿名册抄录如下）。

| 项 | 读数 |
|---|---|
| rc | （基线跑完后填，非占位交付） |
| PASS/FAIL/SKIP 三枚数 | （基线跑完后填） |
| 含 hotkey/ball 的绿用例名册 | （基线跑完后填） |

## §1 ①：档位词裁定形状（defaults 说真话）

（写码前把现状三点抄录：终值来源、印词、缺失退回；改后逐处 `file:line` 落点与用例名。）

## §2 ②：两枚 winlive 尺子的形状修复

（rebind 枚：sink 是 JSON 冒号格式、判据曾 grep 等号格式；occupied 枚：曾赌 Ctrl+Alt+U 未占。
改后的判据形状逐条写，含真窗禁跑欠账的具名登记。）

## §3 门禁读数

| 尺 | 结果 |
|---|---|
| `go test -count=1 ./cmd/wisp/`（带 sherpa PATH） | （改后复跑填） |
| `go vet -tags winlive ./cmd/wisp/` | （填） |
| `go build ./...`（仓根全量） | （填，逐字输出） |
| `scripts/d22scan.sh`（若可用） | （填） |

## §4 变异自证

- ①那枚：真跑（`go test` 可跑）——把 defaults 换回 config 的发必须红，换回必须绿；红句逐字贴这里。
- ②两枚：真窗不许跑 ⇒ 静态反证（"永远满足"形为什么在真窗下也会绿）；写清楚，不喊"有牙"。

## §5 判不动的地方（欠账具名）

- （此处登记：修完的 winlive 尺子有没有牙＝本轮量不到，归 258-v2 之后与 owner 同意的窗口。）

## §6 越界与禁区自查

（git diff 名册逐枚过：写面只允许 cmd/wisp/** 与 .scratch/wisp/probes/258/r2/**；
tools/d22scan、internal/** 零触碰；票 260 射程（cancel 键注册）零顺手修；票面 AC 框零触碰。）

## §7 Progress log

- 09:0x 起手：date/锚现取；票面翻勾节＋258-v1 判决书＋五枚落点现读。基线整包后台跑。
- （后续追加，不删。）
