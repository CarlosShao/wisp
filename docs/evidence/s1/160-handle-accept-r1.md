# 票 160 · 160-v1 验收程（非实现者）—— r1 三格独立复算判语

- 派单：`.scratch/wisp/dispatches/2026-09-27-115x-accept-160-r1-v1-independent-verdict.md`
- 本程时刻：09-27 11:51 起手
- 被验对象：`docs/evidence/s1/160-handle-r1.md`（＝待复算断言，非取证来源）＋ 三枚提交 `34c0b4e4`／`f576cf08`／`027ff82e`
- 本程台件：`.scratch/wisp/probes/160/v1/**`（**只建不删**；两枚自有样本＋六枚 34c0b4e4 blob＋三枚变异体＋读数。实现者 `probes/160/r1/**` 与 `c1/**` 一枚字节未动）
- **总判语：三格全部「成立」。** 逐格判据与读数见 §3；承重式与残项见 §5／§12。本格不勾票面、不改 Status、不动 `probes/**`、不修 G5pos。

---

## 1. step-0 四件

1. `date` → `Sun Sep 27 11:51:48 CST 2026`。
2. `git rev-parse --abbrev-ref HEAD` → **`dev`**（全程未切分支；末次复看见 §7 提交时的 `git show --name-only HEAD`）。
3. `git rev-parse HEAD` → `027ff82e6872e6b0d7600245d09b5b142c5b3b4c`（＝被验第③格提交本身；编排者给的三枚号**未当锚抄**，本程的锚＝工作树现量）。
4. `git status --porcelain -- internal/risk/ internal/tools/ cmd/wisp/` → **空**（工作树确在被验三枚号上 ⇒ "读工作树"这一路合法）。

票面 AC 连**现量行号**（`.scratch/wisp/issues/160-...-closer.md`，本程自己 `Read` 全份）：

| 行 | 条目 |
|---|---|
| `:26` | **AC#1** 改前"只开不关"证成读数：编译器不拦、门不响、只有普查尺点名；量不到就停手 |
| `:27`（更正块 `:28-33`） | **AC#2** 换形状：句柄自带关闭＋认身份＋幂等；后半句 `DisposalScope` 被 `:32` 降格为"由调用方在已有 CloseTask 缝上关"；`:33` 把票面 `:12`"写不出来"降格为"关不掉别人的、别人的也关不了你的；忘关由门响亮地点名" |
| `:34` | **AC#3** 反向判据两问：摘掉认身份是否存在跨包误关打不红？有无外部可见读数变过？ |
| `:35`／`:36` | AC#4／AC#5 —— **不在本单射程**（派单 §1 只交三格） |

## 2. 本程没测什么（一条不许省）

- **没跑全仓 `go test ./...`**：派单 §4 明令逐包单跑，本程只跑了 `./internal/risk/`、`./internal/tools/`、`./cmd/wisp/`（健康 bench）与 `tools/d22scan`。
- **没复跑 gofumpt**：该尺本程由编排者代跑（`gofumpt v0.12.0 (go1.27.1)`、548 枚已跟踪 `.go`、0 行红）——**档位登记**：实现者自报只跑了 `gofmt`（gofumpt 不在其 PATH），编排者代跑的才是 gofumpt 档；本程按派单不重复，引用时带这个档位差别。
- **没跑 golden／SLO／`thresholds.go` 任何尺**：非本票写面、非派单要求。
- **没有把实现者 `probes/160/r1/ac1/` 再跑一遍当我的正控**：改前两向读数一律由本程自有样本产出（`probes/160/v1/ac1|ac1b`）；它的存档读数文件只作文字对照。
- **未复现 `0xc0000135` 档**：本程 `cmd/wisp` 直接走健康 bench（`PATH=$PWD/third_party/sherpa-onnx`），一次成功（见 §7），没有复算"不接 PATH 整包量不到"那一档（那是既有环境事实，非本票改动）。
- **未复现计时噪声红**：`TestResolvePerCallBudget` 在本程全部五趟 `internal/risk` 全量跑（基线／pre-full 改前／m1／m2／m3）里**一枚都没红**；阈值一字未动、未 Skip、未改名，归类＝**本程未复现的既有计时噪声**（真红 0 枚）。
- **AC#4／AC#5 两格未裁**（不在派单）；票面一格未勾。
- `probes/154/gate-clauses.sh` 本程**只执行、零写入**（输出重定向进本程 `probes/160/v1/out/`）。

## 3. 三格逐格判语

### 3.1 AC#1 —— 判语：**成立**

改前读数一律用本程自造样本，打在 `34c0b4e4` 的真 `internal/risk/provenance.go` 上（`-overlay` 映回，工作树 `internal/**` 零写入）。

| 腿 | 本程命令（可复制） | 本程读数 | 对照 r1 断言 |
|---|---|---|---|
| LEG 1 只开不关·**改前** | `cd .scratch/wisp/probes/160/v1/ac1 && go vet -overlay="$PF" ./...`（`$PF=…/prebase/overlay-pre-full.json`） | **rc=0** | 一致（它：rc=0） |
| | `go run -overlay="$PF" .` | `V1-OPENONLY leak-inspect hit=true src="web.fetch" taints=1 guessed-id-taints=0` ← **污点活着、scope 开着** | 一致（形状同、marker 各造） |
| LEG 1·**改后（工作树，无 overlay）** | `cd …/v1/ac1 && go vet ./... && go run .` | vet rc=0；运行读数**与改前逐字节同式**（hit=true taints=1） | ⇒ r1 §4.3 第 3 枚"statement 形照样写得出来"的降格**属实** |
| LEG 3 外人按 id 关·**改前两向** | `cd …/v1/ac1b && go vet -overlay="$PF" ./...`（rc=0）；`go run -overlay="$PF" . 2>&1` | `V1-FOREIGN-A no-other-taints before=1 after=0 inspectHit=false`（**静默消失，零读数**）／`V1-FOREIGN-B witness-tainted inspectHit=true src="unbound-scope"`（登记里有别的污点时判决变响） | **比它多做了一向**：r1 只量了 A 静默；本程 B 向钉住"日志/判决只在『别的 scope 带污点』时才存在"这个条件分支确实按它说的形状走 |
| `:453` fail-closed 日志打没打 | `grep -c 'fail-closed R4' out/ac1-pre-foreign.txt` | **0**（A 向分支不走 ⇒ 没得打；stderr 捕获里除 `winsec` 一行 INFO 外零行） | **它说没打 ⇒ 复算属实**；〔推断→本机读数〕的升档**配得上**（本程另加 B 向控制，档位不降反升） |
| 门（改前） | `bash .scratch/wisp/probes/154/gate-clauses.sh 34c0b4e4`（`out/gate-34c0b4e4.txt`） | **gate_rc=0**；G5 主尺名册＝**1 枚、粒度是文件名**（`UNPAIRED cmd/wisp/panel_assets.go`）；G5 正控 quiet；G5-负一负按设计响 | 一致；本台件在 `.scratch/` 下、G5 射程（`internal/**/*.go cmd/**/*.go`）**结构上吃不到本程样本** ⇒ r1"换个目录就看不见"一句在本程同样成立 |
| 事后观察口 | 本程样本里 `ScopeTaints(已知 id)=1`／`ScopeTaints(猜的 id)=0`（同一行输出） | 无导出"列出开着哪些 scope"的口（`git grep -n "^func (p \*Provenance)"` 复核：按-scope 出口只有 `ScopeTaints`/`scopeMarks`，前者要 id、后者非导出） | 一致（三态塌成 0） |

**判据本体核对**（票面 `:26`：编译器不拦✓／门不响（只点文件名，且不吃台件目录）✓／只有普查尺点名✓／事后查不出✓）。四把尺本程全部现量、未借它的数据。

### 3.2 AC#2 —— 判语：**成立**

**(i)／(i′) 认身份两枚用例＋(ii) 幂等分格——本程读码核了三件事**（`internal/risk/provenance.go:383-466`、`provenance_test.go:906-997`）：

- **分格没有**：`Close()` 四味各一条出口、各一枚 sentinel——`s==nil||s.p==nil`→`ErrScopeNotOpen`（`:426-428`）、`s.closed`→`ErrScopeAlreadyClosed`（`:432-434`）、`reg==nil`→`ErrScopeNotOpen`（`:437-439`）、`reg.owner!=s`→`ErrScopeNotOwner`（`:440-441`）；测试面用 `errors.Is` **逐枚钉**（`:912`／`:940` NotOwner、`:966`／`:972` AlreadyClosed、`:978` NotOpen），**没有一格靠 `err != nil` 混打**。⇒ 不是"一条 err!=nil 吃四味"。
  ⚠ 一枚**观察**（不阻断判语）：`reg==nil` 且句柄未 spent 那一格（stale 句柄碰上 fresh 关完又没了登记）今天**没有专测**——唯一被钉到的 `ErrScopeNotOpen` 是 nil-receiver 形。可达序列存在（`stale→fresh→fresh.Close()` 后 `stale.Close()`＝NotOpen）；本程判定该缺口属 AC#4"明写治不到什么"的射程（未派），记进 §12 next。
- **(iii) 旧口子的档位**：本程**自造**外人按 id 关的样本（`probes/160/v1/ac1b`，与 r1 那枚同谓词、不同文件）打在改后工作树：`go vet ./...` → `vet.exe: .\closebyid.go:18:54: p.CloseScope undefined (type *risk.Provenance has no field or method CloseScope)` rc=1；`git grep CloseScope -- internal cmd` 非注释命中＝0（只剩 `provenance.go:361/:562` 两处**讲历史的注释**与两枚测试注释）⇒ 导出面**整枚消失**属实。
  它买到的是**"写不出来"**（编译器读数，可信）——但只覆盖**外人关**这一形；**忘关**那一形本程实测改前改后**同一份样本照样编译、照样漏**（§3.1 LEG 1 两档同读数）。票面 `:12` 那句在 `:33` 被降格成"关不掉别人的、别人的也关不了你的；忘关由门响亮地点名"——**逐字对 `provenance.go:366-372` 的新注释块**：写的正是这个降格、还自带"measured rather than claimed"指回台件。承诺文字与读数**一致**，无oversell。
- **生产三枚调用点改口属实**：`bridge.go:648` 开腿把句柄**存进桥账本**（`:143` `map[string]*risk.Scope`）、`:706` 关腿 `scope.Close()`＋错误并进审计行 `:708`（追加 `close_err=%v`）；`panel_assets.go:243` `_ = prov.OpenScope(...)` **写明故意的丢弃**。`git grep "risk\.Scope\b"` 排包内 ⇒ **只有 `bridge.go:143/:186` 两枚＋一枚测试注释**，无第四枚地界。
- **DEFERRED 枚数尺**：`git grep -c "DEFERRED(C25-loop-wiring)" -- internal cmd tools` ⇒ 三文件各 1、**恰 3 枚**（枚数未变；行号 55→58 的位移与台账 `:7660` 过期两笔，编排者已在台账 09-27 11:5x 补账条目里自落——本程复核**该补账在盘**）。

## 4. 恒真检查那一发（两向读数都贴）

问：AC#2 那两枚新用例（实为四枚：`TestScopeCloseRefusesAnotherOwnersRegistration`／`TestScopeCloseRefusedOnEmptyRegistrationIsStillVisible`／`TestScopeCloseIsIdempotent`／`TestScopeIDGrantsNoPower`）在**改前**码上响不响？

| 向 | 命令 | 读数 |
|---|---|---|
| 新测试 × 旧实现（只映回 `provenance.go`） | `go test -count=1 -overlay="$V1/prebase/overlay-pre-provonly.json" -run 'TestScopeClose|TestScopeID' ./internal/risk/` | **build failed**：`provenance_test.go:908:11: p.OpenScope("task-victim") (no value) used as value`、`:912:43/:940:43: undefined: ErrScopeNotOwner`、`:341:33: undefined: Scope` 等 ⇒ **旧码上根本站不起来，谈不上响** |
| 全 pre-base 六文件映回 `34c0b4e4` | `go test -count=1 -overlay="$V1/prebase/overlay-pre-full.json" -run 'TestScopeClose|TestScopeID' -v ./internal/risk/` | `testing: warning: no tests to run`／`ok … [no tests to run]` ⇒ **改前这四枚不存在**，旧全量 100 枚（99 PASS＋1 SKIP）对着旧码全绿（`out/roster-pre.txt`） |
| **真检的反证（比"改前不响"更硬的一向）** | 三枚自有变异打改后码（§5） | m2 摘认身份 ⇒ **恰好两枚认身份用例红**、其余 101＋SKIP 照旧；m1 摘幂等分格 ⇒ **恰好幂等用例红**；m3 摘 nil 守卫 ⇒ 该用例 panic 红。**每一味摘掉都有且只有点名它的用例红** ⇒ 这四枚不是恒真，是**承重探针**。 |

⇒ 本仓否过四次的形状（用例在坏码上从不响）这一发**没有蒙过去**：两向＋变异反证三路全在案。
