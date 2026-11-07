# 票 316 · AC#1 三选一裁定（腿 `316-r2` · 2026-11-07）

- 腿＝写腿（编排者派单，四步拆步落盘）；锚点＝`git rev-parse --short HEAD` ＝ **`ea602a61b`**（2026-11-07 15:5x 现读）。
- 产码 md5 `cmd/wisp/config_reload.go` ＝ `0934ee7a88cd1c6e7cbd3dd0ce123971`；尺 `cmd/wisp/config_reload_223_test.go` ＝ `11d42757a64f167dbd60d8d44b9128a6`。
- ★先交代身份事实：票面 Status 已是 **withdrawn**（编排者 2026-10-14 13:38 撤，判与票 272 重复，台账 `A924`；实现腿 `316-r1` 全套读数在 `.scratch/wisp/probes/316/r1/`，交件 `8cc028f8`）。本腿按派单四步走，**独立现跑重核撤票的事实基础今天是否仍成立**；⛔ 不翻 Status、不翻 AC 勾。
- 起手三核③：`git status --porcelain -- cmd internal` 有**他腿在飞未跟踪件 3 枚**（`cmd/wisp/panel_resident_323_windows_test.go`、`cmd/wisp/panel_resident_assembly_323_test.go`、`internal/plugin/manifestslot/command_138_cancel_test.go`）——本腿零碰；家族测试编译包含它们且全绿 ⇒ 本件全部红绿读数非在飞件造成。

## 前半（步1）：三候选今天各自的真实形状（全部本腿现跑／现读）

### (a) `strings.Count(win, 该句)==1`（窗口内恰一次）—— **已在盘（重复落地＝零增量）**

- 载体：`cmd/wisp/config_reload_223_test.go:748`：`if win := r.stdoutSince232(t, mark); strings.Count(win, restartSentence272) != 1 {`（**arm2**）。
- 来路：票 272 落地（`8d30a862`，2026-10-07 09:33）。
- 对 M 形的牙：**零**——M 的窗与正常一发逐字节相同（横幅在 mark 之前、窗外），窗内怎么数都是 1（r1 引 272-r2 甲-alone 尖读数：M PASS 2.15s）。它保留的牙是对"窗内 ≥2 份"。

### (b) 标记取在 `state=armed` 横幅之后 —— **空操作**

- `mark` 现址 `:654`（`mark := r.h.out.String()`），取在测试体内；测试体跑在 `rtHook`（assembleRuntime 返回）之后；启动横幅＋`state=armed` 都由 `startConfigReload`（`config_reload.go:118`；调用点 `run.go:1149`）**同步**写进 syncWriter ⇒ mark 本来就在全部启动输出（含 armed 行）之后，再往后挪＝一句可挪出去的字节都没有。
- r1 引 272-r2 尖读数：乙在 M PASS 2.88s / cur PASS 2.15s ＝ **既不变严也不变松**。

### (c) 要求该句偏移晚于 `state=applied` —— **结构性误伤（落下去必假红）**

- 今日产码次序实测（`cmd/wisp/config_reload.go`）：`rt.mgr.OnRestartPending = rt.reportRestartPending`（`:128`）→ reload 轮询 `CheckAndReload()` 内经 OnRestartPending 发操作员句（`reportRestartPending` `:346`）→ `rt.reportReload(rep)`（`:193`）→ `auditf("config: HOT-RELOAD state=applied …")`（`:211`）。
- ⇒ **操作员句天然早于 applied 行**；"要求晚于"＝正常一发确定性 FAIL（r1 引 272-r2 丙：cur FAIL 42.20s 超时）＝把真绿判红。

### M 形今天现跑（本腿以 HEAD 为 base 重建突变；全走 overlay，共享树零编辑）

- 手法：`probes/316/r2/config_reload.M.go` ＝ HEAD 产码副本＋在启动横幅（`config_reload.go:148` 的 Fprintf）后追加一行含 `本次运行不会生效` 的横幅句（注入文本逐字承 `232-v1`/`316-r1`）；overlay 件 `probes/316/r2/prod-M.json`。
- 读数：`go test -count=1 -run TestTicket223RestartTierSaysItWillNotApply -overlay .scratch/wisp/probes/316/r2/prod-M.json ./cmd/wisp` → **rc=1，`--- FAIL: TestTicket223RestartTierSaysItWillNotApply (2.85s)`，红＝`config_reload_223_test.go:676`（arm1）**，红句 dump 第三条＝M 注入的横幅同名句逐字（log：`probes/316/r2/logs/M-rerun-ea602a61b.txt`，窗 15:59:0x→15:59:22）。
- 正控 cur（无突变）同发：**PASS，ok 3.928s，rc=0**（`probes/316/r2/logs/cur-rerun-ea602a61b.txt`）。
- ⇒ **判据①（复造 M ⇒ 必须红，红句逮住横幅同名句）今天在 HEAD `ea602a61b` 上现跑成立；让 M 红的是 arm1（票 272 的枚举外臂），不在本票三候选菜单内。**

### 行号漂移备注（三核①）

- 票面引 arm1 `:595/:596`、arm2 `:668` 是 10-07 那版尺的行号；其后 351-r2（`7c5074f91`）在同文件 +48 行 ⇒ 今日实读：arm1＝`:675-:680`（红句在 `:676`）、arm2＝`:748`、mark＝`:654`。**串都在，行号过期**——引用以今日 sed/grep 现读为准。

（后半部分：步2 裁定、步3 落地结果、步4 门禁，见下方追加节。）
