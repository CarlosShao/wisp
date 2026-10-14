# 票 316 · 腿 `316-r1`（实现者）——M 形收口：**现读裁定**与四形读数

裁决对象＝`cmd/wisp/config_reload_223_test.go`（现行 HEAD `35da8bc9`，md5 `2ea71578…`）＋票面现量四条。
本腿全部突变走 `go test -overlay`（以**现行 HEAD 产码**为 base 重建；手法/注入文本逐字照抄
`232-v1` 的 `gen_mut{,2,3}.py` 与 `232-r2` 的 `mutG_banner_only`），**共享工作树零编辑**。
⛔ 本腿不勾 `- [ ]`、不动 `Status:`/票名、不碰别人脏件。

**结论先行（四行）**

1. **M 形今天已经红**（红在 `:596` arm1，红句 dump 里横幅同名句逐字可见）⇒ 票面判据①今天已满足，
   **无需新写任何断言**。票面"现量"里那句"终态仍 PASS 2.62s"是 `232-v1` 在 **10-06 19:0x** 量
   `3d9b8374` 那版尺的读数（现量：那版尺 `grep -c restartSentence272`＝**0**），**早于票 272 落地**
   （`8d30a862`，10-07 09:33 入库 arm1＋arm2）⇒ **该现量在现行 HEAD 上是过期读数**。
2. **三选一逐支现读裁定**：(a) **已在盘**（arm2 `:668`；对 M 零牙，272-r2 甲-alone 读数 M PASS 2.15s）；
   (b) **空操作**（`mark` `:574` 本来就在全部启动输出之后；272-r2 乙 M PASS 2.88s）；
   (c) **结构性误伤**（产码次序＝操作员句先、`state=applied` 后；272-r2 丙在**正常一发**上 FAIL 42.20s）
   ⇒ **无一支可落**；能同时让 M 红又不误伤的**只有 arm1（272 的枚举外臂），它已在盘**。
3. **停下上报的具名触发**（派单 §6）：**"正例 G 仍绿"与判据①互斥**——`232-r2` 的 G（横幅带全字样、
   通知完好）与 `232-v1` 的 M（横幅带同样字样、重启句照打）**是同一个形**（都是"种前横幅含
   `本次运行不会生效`＋真句照打"）；今天 G 红在同一 arm1（§3 读数）。要让 G 绿只能**放宽 arm1**＝
   违判据④与票 272 AC#3 ⇒ **本腿零写面、停手上报**，请编排者裁（§7）。
4. **"0ms"读法（本票归我裁）**：按**"不许假绿"**读（承 232 翻勾记录 10-06、"票 272"记录 10-07 的同一裁定）。
   本腿 6 发读数里 **0ms 从不出现**：红＝4.38s／4.38s／44.94s／44.24s，绿＝4.81s／5.64s。

---

## §1 起手锚（本腿 13:0x–13:2x 现量）

| 尺 | 命令 | 读数 |
|---|---|---|
| HEAD | `git rev-parse --short HEAD` | `35da8bc9` |
| 写面 | `git status --porcelain -- cmd internal`（起手） | **0 行** ⇒ `cmd` 写面干净（落盘时刻 `internal/tools` 出现他腿 2 枚 M：`task_output_pointer_notice_test.go`、`taskid_source_236r4_teeth_test.go`——**本腿未碰**，在飞的 `174-r6`/`236-r6` 地界） |
| 产码 | `md5sum cmd/wisp/config_reload.go` vs `git show HEAD:… \| md5sum` | 同值 `348f5a67a73442d45f751079e3b80a6c`（⚠ 与 10-06 的 `5ce441ca…` 不同：`diff` 全量＝**24 行、两处 hunk，全在文件头注释**（`c5040ea7`/`8dbae00c`），产码语义零动 ⇒ 本腿突变以现行 HEAD 为 base 重建，不直接用 10-06 的副本） |
| 被审尺 | `md5sum cmd/wisp/config_reload_223_test.go` vs `git show HEAD:` | 同值 `2ea71578ba9ca015c535673bc9792d09`（＝`272-v1` 验收记录里的同值） |
| 票面框 | `grep -c '^- \[ \]' <票面>` | **1**（AC#1 未勾）⇒ ⛔ 本腿不翻 |

## §2 现读：M 已红（票面现量的时间线对拉）

- `3d9b8374`（232-r2 尺，`232-v1` 量它）与 `d6243757`（`232-v1` 起手 HEAD）：`grep -c restartSentence272`＝**0／0** ⇒ 无臂。
- `8d30a862`（10-07 09:33，272-r2 落地）：arm1 `:595`（`strings.Count(prePlant, restartSentence272) != 0`）、arm2 `:668`（窗内 `!= 1`）。
- 现行 HEAD `35da8bc9` 尺：双臂在盘（`:594-598`／`:667-671`）。

**本腿自造 M（现行 base 重建）**：`rc=1`，`--- FAIL: TestTicket223RestartTierSaysItWillNotApply (4.38s)`，
**只有 1 条红**（无超时）⇒ 让 M 红的就是 arm1；红句逐字：

```
config_reload_223_test.go:596: the operator stream carried "本次运行不会生效" 1 time(s) BEFORE this run planted the edit; a copy nobody planted cannot be evidence that this run will not use the new values. stdout before the plant:
    …（答复监听横幅）…
    …（配置热加载已被接管的启动横幅）…
    wisp run: 这些段的改动本次运行不会生效（D36 重启档，需要重启进程）。原因：上一版横幅。涉及：app.autostart / app.language / app.single_instance。注意两件事都没发生：文件里的新值没有被丢掉。
```

⇒ 判据①的"红句体现**横幅同名句被逮住**"**逐字成立**（第三条 dump 行就是 M 注入的横幅句）。

## §3 四形＋三枚正控：全部本腿自跑（逐发窗口起止＝wall 时刻；全部 `logs/AC1-*.txt`）

| 形 | overlay | rc | 结果（test 时长） | 红句 | 窗口起止 |
|---|---|---|---|---|---|
| `cur`（无突变） | — | 0 | **PASS 5.64s**（ok 7.882s） | — | 13:20:56→13:21:10 |
| `MIRROR`（逐字节副本） | `prod-MIRROR-finaltest.json` | 0 | **PASS 4.81s**（ok 7.078s） | — | 13:21:21→13:21:34 |
| **`M`**（横幅＋同样字样，真句照打） | `prod-M-finaltest.json` | 1 | **FAIL 4.38s**（ok 6.811s） | `:596` arm1（§2 逐字） | 13:21:34→13:21:49 |
| **`G`**（＝232-r2 正例：横幅带全字样，通知未动） | `prod-G-finaltest.json` | 1 | **FAIL 4.38s**（ok 6.600s） | `:596` arm1（dump 含 G 注入行逐字） | 13:21:49→13:22:04 |
| `MK`（M 横幅＋真句删光） | `prod-MK-finaltest.json` | 1 | **FAIL 44.94s** | `:596` ＋ `:623` | 13:22:38→13:23:36 |
| `SA`（审计两行挪到启动处） | `prod-SA-finaltest.json` | 1 | **FAIL 44.24s** | `:600` | 13:23:36→13:24:32 |
| 家族最小集 `-run 'TestTicket223'` | — | 0 | **PASS=10 FAIL=0**（62s） | — | 13:28:06→13:29:08 |

红句逐字（第二、三行形）：

```
MK :596 the operator stream carried "本次运行不会生效" 1 time(s) BEFORE this run planted the edit; …（同上 arm1）
MK :623 stdout never carried "本次运行不会生效" AFTER the plant within 40s; window since the mark:
      （窗口空行＝真没打印；full stdout 里红句前三行带横幅句＝种前字节）
SA :600 the audit trail never carried "config: HOT-RELOAD state=restart-pending sections=[app]" AFTER the plant within 40s; window since the mark:
      [audit] config: HOT-RELOAD state=applied hot=[] reload=[] restart=[app] locked=0   ← 窗内只剩它
```

⇒ 对照判据：**M/MK/SA 全红**（不回归）；`cur`/`MIRROR` 两枚真绿不误伤；**G 红**（§5 要编排者裁）。

## §4 三选一逐支判定（具名理由）

- **(a) 窗口内恰一次**：**已在盘**（arm2 `:667-671`，272-r2 落地）。272-r2 的尖读数：**甲-alone 在 M 上 PASS 2.15s**
  （`probes/272/r2/logs/AC2-candidate-JIA-M.txt`）＝M 的窗与正常一发逐字节相同（横幅在窗外），窗内怎么数都是 1
  ⇒ **(a) 对 M 零牙**；它保留的价值是对"窗内 ≥2 份"有牙（甲-alone REPEAT FAIL 2.12s）。**落无可落（重复）。**
- **(b) 标记取在 `state=armed` 之后**：**空操作**。`mark`（`:574`）由测试体取，而测试体跑在 `rtHook`
  （`assembleRuntime` 返回之后）；启动横幅与 `state=armed` 都由 `startConfigReload()`（`run.go:813`）
  **同步**写进 `syncWriter` ⇒ mark 本来就在全部启动输出之后。272-r2 尖读数：乙 在 M 上 **PASS 2.88s**、
  在 cur 上 PASS 2.15s（`AC2-candidate-B-{M,cur}.txt`）＝**既不变严也不变松**。
- **(c) 该句偏移晚于 `state=applied`**：**结构性不可能**。现行产码次序（`cmd/wisp/config_reload.go`）：
  `:157 reloadOnce → rt.mgr.CheckAndReload()`（操作员句 `:315-330` 经 `OnRestartPending` 在**其内**发出）
  → `:164 rt.reportReload(rep)` → `:183 auditf("…state=applied…")`。272-r2 尖读数：**丙 在正常一发 `cur` 上
  FAIL 42.20s 超时**（`AC2-candidate-C-cur.txt`）＝确定性地把真绿判红 ⇒ 直接违"不许误伤真绿"。

⇒ **三支联合也过不了 {① M 红 ∧ ② 不误伤}**：能对 M 有牙的闭口形状只有 arm1（272 的枚举外臂，"票 272 ☑合法落地"），
**它已在盘**。⇒ 本票"收口"的实体已在 10-07 由票 272 交付，**票 316 是同一缺口的重复票**（来路、措辞、三候选逐字同源）。

## §5 G 的冲突（本腿停下来交给编排者的那一格）

- `232-r2` 的 G（`mutG_banner_only`，正例对照"窗口化不许误伤好文案"）在**232-r2 的尺**上 **PASS 2.26s**。
- **今天**：同发 G **FAIL 4.38s 红 `:596`**（同一 arm1；dump 第三条＝G 注入行逐字，`logs/AC1-G.txt`）。
- 机制：**M（232-v1）与 G（232-r2）是同一个形**——都是"种前横幅含 `本次运行不会生效` ＋ 真句照打"，
  只有横幅句的外壳文字不同。任何以"该句在种前出现过"为判据的臂（arm1 就是）**必然两形同判**；
  要"红 M 而绿 G"只能按**横幅句的具体措辞**取判据——不在票面菜单、非原理性，本腿不落。
- 两票标签相反（232-v1：M＝必须红；232-r2：G＝必须绿），票 272 的 arm1 把它裁成"必须红"（编排者 `A661` 翻勾时在案）。
- 今天真绿的正例＝`cur`（无突变）＋`MIRROR`（overlay 落地正控）；"横幅只含部分字样"边界＝272-r2 的 PART **PASS 2.02s**。

## §6 门禁（本腿亲跑，读数在 `logs/`）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/wisp/` | **rc=0** |
| gofmt／gofumpt | `gofmt -l`／`gofumpt v0.12.0 -l`（两枚 `.go`） | **空**（本腿零写面，两枚文件本来就干净） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**；正控先绿 `top-level: PASS=35 FAIL=0 SKIP=0`；`clean`；`ban #8 cmd/ examined 125 Go files`（`logs/G1-d22scan.txt`） |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **rc=0 / VERDICT GREEN**（12675 tracked paths，over-budget 59 全在册，`logs/G2-path-length-budget.txt`） |
| 家族 | `go test -count=1 -run 'TestTicket223' -v ./cmd/wisp` | **rc=0，PASS=10 FAIL=0**（`logs/AC3-family-TestTicket223.txt`） |
| 写面终态 | `git status --porcelain -- cmd`＋`md5sum cmd/wisp/config_reload.go` | `cmd` **0 行**；产码 md5 `348f5a67…` 复量同值（6 发突变全 overlay，工作树零编辑） |

## §7 具名没做／交编排者裁

1. **零写面**：本腿**未动 `cmd/wisp/config_reload_223_test.go` 一个字节**（无新增臂、无注释级改动）——
   理由＝§2＋§4＋§5：闭合已存在、菜单三支皆不可落、唯一剩余可满足形状要放宽既有断言（违④）。
2. **(b)／(c) 未落地**（理由 §4）；**(a) 未重复落地**（已在盘）。
3. **"审计侧种前同名臂"（SA 的 arm1 类比）未做**：不在票面菜单，属票 272 判据面；今天正常路径打印 0 份 ⇒
   属"前向形状"加严，需编排者裁后另派（⛔ 本腿不自选加戏）。
4. **未跑整包**（派单禁令：314 死区）；⛔ 未碰 `internal/tools` 两枚他腿脏件；零 push。
5. **建议编排者裁**：票 316 与票 272 同一缺口（票 272 现只剩 AC#5/AC#7 未勾）——
   收编／撤票／改验证票由你定；本件全部读数可作那张裁决的底座。
