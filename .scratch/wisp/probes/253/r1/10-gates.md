# 253-r1 · 门禁四把读数（10-gates）

时刻 `2026-10-08 16:53 +0800`（`date` 现跑）；HEAD 此刻 `89443d0a`（别人 ticket-181 的笔，我的第 1 笔 `65f4c968` 用 `git merge-base --is-ancestor 65f4c968 HEAD` 现验＝**是祖先，rc=0**，未被埋）。

⚠ 本程所有 `go test` 都要 `PATH` 里带 `third_party/sherpa-onnx`（`cmd/wisp` 的测试二进制链 sherpa cgo 导入库，缺 DLL 会在**加载期**给 `exit status 0xc0000135`＝仪器没跑到，不是红）。权威原文＝现读 `scripts/wisp-cli-tests.sh:99-113`：`pwd -W` 那种盘符＋正斜杠形会打回 `0xc0000135`，本程用的是 shell 自己的路径形 `export PATH="$PWD/third_party/sherpa-onnx:$PATH"`（`$PWD=/d/work/workspace/projects plans/Wisp`）。
我在没注 PATH 的第一发起过一次，读数＝`exit status 0xc0000135` ＋ **0 条 `=== RUN`** ⇒ 判"根本没跑到"，不当红也不当绿（`scripts/wisp-cli-tests.sh:13-19` 那条负控同形）。

## 1) 行尾符尺（先跑这把再看 gofmt）

| 对象 | 工作树 CR 字节 | `git show HEAD:` CR 字节 |
|---|---|---|
| `cmd/wisp/panel_host_windows.go` | **0** | **0** |
| `cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go`（我的新件） | **0** | （新件，无 HEAD 版本；入库后 `git hash-object`＝`f8dadd15…`＝`git rev-parse HEAD:` 同值） |

⇒ `core.autocrlf=true` 这条设置是真的，但**这两枚文件盘上就是 LF**，`git hash-object <工作树>` 与 `git rev-parse HEAD:<path>` **逐字相等**（`58e2b155dbcc8d380793e90bffc5514ab00743f1`）⇒ 本程"没被 autocrlf 洗过"是量出来的，不是猜的；也因此 `git show HEAD:<file> > <file>` 这种还原法是字节安全的（六发突变全部这样还原，见 `20-mutations.md`）。

## 2) gofmt

```
gofmt -l cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go   → 空     rc=0
```
⚠ 第一版**被点过名**（我手写的 struct 字段对齐不合 gofmt），`gofmt -w` 之后清空；CR 字节前后都是 0＝`gofmt -w` 在这台机上没换行尾。

## 3) go vet

```
go vet ./cmd/wisp/   → 零输出   rc=0
```
（六发突变每发的最终态都是"还原后的干净树"，上面那枚 `rc=0` 是在干净树上跑的。突变态里我只在 M5 那次顺手跑过一次 `go vet ./cmd/wisp/ 2>&1 | tail -3`＝**零输出**，⚠ 但那发尾接了管道 ⇒ **`rc` 没取到**，只能算"vet 没吐诊断"，⛔ 不许被我这段话引成"种了突变 vet 仍 rc=0"的证据。）

## 4) 靶向 go test（⛔ 没跑整包）

尺＝
```
go test -count=1 ./cmd/wisp/ -run 'TestTransportDoorBindingMatchesRoster253r1|TestBindingRosterBitesItsOwnFixtures253r1|TestPagePostMessageEnvelopeReachesDispatchRawViaTransport|TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge|TestInboundSourceGuardRefusesForeignSourceOnPageEdge|TestInboundRequestIDGuardRefusesMissingIDOnPageEdge' -v
```
读数＝**6 枚顶层全 PASS／0 FAIL／0 SKIP，rc=0**（`ok github.com/CarlosShao/wisp/cmd/wisp 0.091s`）：
- `TestTransportDoorBindingMatchesRoster253r1` PASS 0.01s ← 本腿那枚盘上尺，日志逐字含 `installPanelTransport at panel_host_windows.go:801 made 1 Bind call(s) and 1 Init call(s); roster holds 1 name(s) over **34 production file(s)**`（整包射程的实证，不是一枚文件）
- `TestBindingRosterBitesItsOwnFixtures253r1` PASS（4 枚子用例 good／drifted-bind／empty-bind／untied-init 全 PASS）
- `TestPagePostMessageEnvelopeReachesDispatchRawViaTransport` PASS 0.00s ← 票 253 `AC#1` 收窄后**能力形那一半**（`panel_transport_35r1_test.go:195→:203` 真调产码）
- `TestInboundRoster{Guard,SourceGuard,RequestIDGuard}...OnPageEdge` 三枚 PASS 0.00s ← 同族既有尺（`panel_inbound_guards_35r3_test.go:150/:158/:167`）

`internal/panel` 那一半（票 253 `AC#2` 的在册尺，`inbound_roster_253_test.go:419/:597/:696/:769/:809`）我也**只按名靶向**跑了一遍：`go test -count=1 ./internal/panel/ -run '<那五枚名>'` → **rc=0**（`ok github.com/CarlosShao/wisp/internal/panel 0.062s`）⇒ 本腿没把这族既有尺弄红。⛔ 没跑 `./internal/panel/` 整包（`A716` §4 现量那包今天 6 枚红＝2 本程有意保留＋4 既有，与我不相干，跑整包只会把别人的红算到我头上）。

⛔ **没跑 `./cmd/wisp/` 整包**——按派单由编排者自己跑；理由具名＝`A718` §3 那发真窗读数（`cmd/wisp/panel_resident_windows_test.go:314 TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` FAIL 20.03s／rc=1，窗口依赖＋本机过期 `frontend/dist`），那格红在我名下没修也不该由我修。

## 5) sh scripts/d22scan.sh

```
sh scripts/d22scan.sh   → rc=0   "clean - no D22 ban violations"
```
分母现量（同一发输出里逐字）：`bans #1-5 internal/=228`、`bans #1-5 cmd/=38`、`ban #6 frontend/=85`、`ban #7 internal/tools/=23`、`ban #8 design/=39`、`ban #8 frontend/=85`、`ban #8 internal/=516`、**`ban #8 cmd/=113 Go files, comments and _test.go included`**（⚠ 按台账记过的口径：ban #8 **含** `_test.go`，所以我这枚新测试文件真在分母里；它零 emoji＝那 4 处非 ASCII 全是 CJK `票`，不在仪器射程的 `U+1F000–1FAFF`／`2200–22FF`／`2600–27BF`／`2B00–2BFF`／`FE0F`／`1F1E6–1F1FF` 任何一段）。另：`skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`。
⛔ 没从仓库根 `go run ./tools/d22scan`（那是独立 module）。

## 6) frontend/dist 依赖声明

本腿**没读也没引** `frontend/**`（派单禁区＋票 253 `:23`），这枚尺只读 `cmd/wisp/*.go` 的 Go 源文本，所以 `A718` §3 那枚"现量 `frontend/dist/index.html` mtime＝09-27 10:59、HEAD 只跟踪 `dist/.gitkeep`"的过期产物限制**对本尺不适用**；我这边任何"绑定成立"的结论都建立在工作树里的 `.go` 字节上，不是建立在那份产物上。
