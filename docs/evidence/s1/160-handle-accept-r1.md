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
