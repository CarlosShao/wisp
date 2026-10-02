# 票 198 · 落地写腿 `198-r2` · 证据件（四小格：AC#2 仪器／AC#5 失败支用例／AC#4 "去哪儿补"／落点仪器收回本票）

> 本件是**实现腿**（写腿）的证据。射程＝票面 §6「10-02 10:1x 验收后裁定」第 5 条排程的四小格：
> ① 给"改一枚默认值"装上牙（AC#2）；② 给"首建失败那句"补一枚用例（AC#5）；
> ③ AC#4 的"去哪儿补"要先回答 DPAPI 那一层；④ 把"落点错了"从别人家的钉收回到本票自己的仪器。
> 上一枚非实现者裁决表＝`docs/evidence/s1/198-firstrun-config-v1.md`（本件要补的洞由它 §1／§3／§6-甲点名）。
> ⛔ 本腿不碰任何 AC 勾选框；⛔ 一字未动 `docs/PLAN.md`／`docs/specs/**`／阈值／golden／`allowlist.txt`／
> `internal/panel` 两枚钉／`internal/perm/ticket90_persist_test.go`／`.github/workflows/ci.yml`。
> ⛔ `frontend/**`／`design/**` 不读、不写、不引。

---

## §0 起手锚（同发取 date／HEAD／porcelain）

**〔本节由编排者代提——腿死于 15:17 前后的模型服务当日额度（通知 15:4x 到、只当未验证；盘上现量为凭），产码它已提交（`fa5593c0`）但证据件还是骨架；下面每一行是编排者 15:46–15:48 现跑，非本腿自跑〕**

- `date`：`2026-10-02 15:46:57 +0800`。
- HEAD：**`f2e5f72a`**（＝腿自己的骨架发，15:18）。
- `git status --porcelain -- cmd internal`＝**0 行**⇒ 本腿产码已全部入库、无突变残留（腿死前把 `mutate.sh` 落在台件目录但**没有执行残留**：写面干净）。
- 本腿两发产码提交：`1d8106d8`（证据件骨架，15:18 前落）＋ **`fa5593c0`**（15:12，`cmd/wisp/firstrun.go` **23/1**＋新增 `cmd/wisp/firstrun_198r2_test.go` **473/0**，去重路径＝2）。

---

## §1 四格各自的修法与判据（每格带 file:line）

**〔本节同上，编排者从 `git show fa5593c0` 的 diff 里逐行读出来代提；判语一格不写〕**

1. **AC#2（改一枚默认值必须响）**：新增 `cmd/wisp/firstrun_198r2_test.go:221` `TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags` 与 `:274` `TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags`。形状＝**期望侧不抄任何字面默认值**，只读 `schema.go` 各字段的 `default` 标签、与首建写出的文件**逐字节比**（提交标题原话"自带八枚种下必红的正控"——哪八枚、种在哪，属 §3 变异自证表，本节不代填）。
2. **AC#5（失败句要有用例守）**：新增 `:299` `TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile`。形状＝把 secrets 存储落点用**普通文件占位**（winsec 拒密封那族）⇒ 钉住三样：失败句响亮、盘上不留半份文件、进程仍退码 2。
3. **AC#4（"去哪儿补"说真话）**：产码只动 `cmd/wisp/firstrun.go:92-115`（diff `+23/-1`）——回执从一句涨到三段，**只写现读出来的两条真走得通的路**：key 走 `wisp secret set <blob 名>`（隐藏输入／`--from-stdin`，明文进 DPAPI 存储、文件里补的是 `api_key_ref = "dpapi:<blob 名>"` 那个**名字**，不是明文；不用 DPAPI 可写 `env:<环境变量名>`）；模型没有 CLI 写入者（`wisp providers discover/probe` 只读该文件去问端点、不代写）⇒ 文案照实说"只有改这份文件一条路"。⚠ 三段里只有命令名／字段名／占位名，**无凭据值**（编排者只读 diff 复核过）。
4. **J1 落点仪器收回本票**：新增 `:432` `TestTicket198R2J1TheAssemblyRootStillCreatesNothing`（两枚同样干净的空目录：装配根不建、`wisp run` 入口建）⇒ "落点错了"从此不再只靠别人票的两枚钉。

---

## §2 门禁四数（build／`gofumpt -l` 本腿动过的目录／`go test` 相关包 `-count=1`／`d22scan` exe）

**〔本节同上，编排者 15:46–15:48 现跑代提，HEAD `f2e5f72a`、写面 porcelain＝0；数字与署名各归各，⛔ 不是本腿自跑〕**

- **尺一 `GOFLAGS= go build ./...`** ⇒ **rc=0、合并输出 0 字节**（15:46:57）。
- **尺二 `gofumpt -l cmd/wisp`** ⇒ **rc=0、list-count=0**（15:47:08；真身 `$(go env GOPATH)/bin/gofumpt.exe`，不在本 shell PATH）。⚠ 口径＝本腿动过的目录（`cmd/wisp`）。
- **尺三 `go test ./cmd/wisp ./internal/config -count=1`**（带 sherpa PATH 前缀）⇒ **合并 rc=0**：`ok cmd/wisp 336.478s`／`ok internal/config 3.759s`——**整包零红**，含本腿新增那五枚用例；那枚既有间歇红（`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`）这一发**没赶上**（它在册读数现为 5 发整包里 1 红，见台账 `A533` §3）。⚠ 336.478s 按【同机可能有争用】读。
- **尺四 `./tools/d22scan/d22scan.exe`** ⇒ **rc=0**，逐字 `clean - no D22 ban violations`＋`examined 262 production Go files under internal/ and cmd/`；分母变化具名：`ban #8 cmd/` 从前值 86 涨到 **87**（＝本腿新增的那枚测试件）。

---

## §3 变异自证表（每格一行：种什么形 → 哪枚必须红 → 复跑终值）

（取数中——⛔ 编排者不代填：腿死前没跑，种哪些形只有它自己的台件 `.scratch/wisp/probes/198/r2/mutate.sh` 记着；这一节归下一枚非实现者验收腿自己取数，我代填就把"谁做的判"洗混）

---

## §4 我可能写错的条目

（取数中——归验收腿）

---

## §5 判不动的地方

（取数中——归验收腿）

---

## §6 交件判语

（取数中——归验收腿；编排者只代提了 §0–§2 的读数，四格有没有牙由验收腿判）
