# 票 111 AC#11 — 非实现者对抗验收腿 `111-v1` · 锚件（第 1 笔）

本目录 = 验收腿 `111-v1` 的证据面。**只攻不润色**；写腿 = `111-r6`。

## 起手读数（本腿自己现跑）

- `date` = `Thu Oct  8 09:11:37 CST 2026`
- `git rev-parse HEAD` = `a809ab74a65544b8a608cb0f4b60b2d049e8da3b`（`a809ab74`，编排者的台账笔）
- 分支 = `dev`
- `git status --porcelain` = **753 行**（`??` 722 / ` D` 16 / ` M` 15）——
  别人的在飞改动：` M .gitignore`、`design/**` 16 枚删除＋`design/doubao/**` 4 枚修改。
  ⛔ 本腿一字未动、未 stage、未恢复。

## 被验的东西（票面原文，本腿 `sed -n '85,100p'` 自取）

- 出处：`.scratch/wisp/issues/111-ci-tests-20-of-33-packages.md:91-94` = **AC#11**，框态 `- [ ]`（未勾）。
- 判据三分：ⓐ 先取当前读数（rc=?，若红则名册入台账，⛔ 不许改那 12 枚断言凑绿）；
  ⓑ 反形自证（塞语法错 ⇒ 那一步必须红，还原并证名册逐枚同）；
  ⓒ **牙的口径＝只买"这 12 枚改坏会被发现"**，⛔ 不许被任何件读成"票 35 `:52`／`:75(c)` 的真窗读数有了 CI 载体"。
- 票面还明写：⚠ 与 AC#7 不同轴（⛔ 不许互相抵账）；⚠ 与 AC#9 的"不许只编译不执行"相反方向——
  本条**本来就只主张编译覆盖**，所以必须**具名标成编译覆盖**。
- 12 枚名册（尺＝`grep -rlE '//go:build.*winlive' --include=*.go .`，本腿自跑＝**12**，与票面一致）：
  `cmd/wisp/` 7 枚（`panel_geometry_255_winlive_test.go`、`panel_host_windows_live_test.go`、
  `panel_transport_live_35v2_windows_test.go`、`resident_approval_live_246_windows_test.go`、
  `resident_ball_live_228_windows_test.go`、`resident_hotkey_live_258_windows_test.go`、
  `resident_task_source_live_246_windows_test.go`）
  ＋`internal/ball/` 5 枚（`hotkey_cancel_borrow_live_260_test.go`、`hotkey_live_test.go`、
  `interaction_live_test.go`、`live_guard_windows_test.go`、`live_windows_test.go`）。

## 写腿 `111-r6` 自报（待攻，非本腿结论）

六笔：`02505442` 锚／`779714f8` ⓐ／`351e5a5e` 加步／`04b01c6c` ⓑ／`e6981776` 门／`b437d82a` 交付件。
自报：ⓐ rc=0；ⓑ 带 tag rc=1（`internal\ball\live_windows_test.go:691:1 expected declaration, found this`）
／同 overlay 摘 winlive tag rc=0；还原 12/12 `git hash-object` 同基线；d22scan rc=0；yaml.safe_load rc=0；步数 9→10。

## 本腿纪律声明

- ⛔ 零 push／零 `add -A`／零 `--amend`／零 `reset`/`rebase`/`stash`/`checkout .`/`clean`。
- ⛔ 未翻任何 `[ ]`→`[x]`；⛔ 未改产码、未改跟踪测试、未改 `ci.yml`。
- ⛔ 未跑 `-tags winlive` 的 `go test`（零开窗）；只跑 `go vet -overlay` 与门。
- 写面只有 `.scratch/wisp/probes/111/v1/**`。
- ⚠ 落件规避：根 `.gitignore:8` 的 `*.out` 会静默跳过（本腿 `git check-ignore -v` 现量证实
  `.scratch/wisp/probes/111/v1/logs/a.out` 被 `.gitignore:8:*.out` 命中），故本腿全部落 `.txt`/`.md`。

## 票 111 框数现量（⛔ 只量不改）

- 尺＝`grep -cE '^[[:space:]]*- \[ \]'` / `'- \[x\]'` 同一文件。
- 未勾 = **5**，已勾 = **6**。与编排者转述的 5/6 **一致**。

## 四问的判（正文见 verdict.md；此处只记每问的落件计划）

- Q1 归属＋前置链：`logs/q1-*.txt`；真实 run 读数＝**判不动（零 push）**。
- Q2 ⓑ 反形复做（自造坏法＋overlay 三连先证台件）：`logs/q2-*.txt`。
- Q3 连带面（`slo-freshness.sh` P1 是对 ci.yml 做断言的仪器，须验）＋两发门：`logs/q3-*.txt`。
- Q4 交付件与结论一致性逐件对拉：`logs/q4-*.txt`。
