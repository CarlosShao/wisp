# 308-r1 · 10-guard-landing（`AC#1` 守卫落点；编排者裁语 ⓐ–ⓒ 落地）

按裁语 ⓐ：产码 `pathresolver.go:263` 的 `drive[:1]` 前一行 `if len(drive)==2 && drive[1]=='$' && …` 已有守卫（尺＝`git show HEAD:internal/risk/pathresolver.go | sed -n '262,263p'`），⛔ 动产码；两枚 `//go:build windows` 测试件是唯一写面。

按裁语 ⓑ：本机 `os.TempDir()` 走 `GetTempPath` ⇒ `filepath.VolumeName` 恒非空；三把台面（default / posix / alias）`go test` 全部 `PASS`、零 panic（尺＝`probes/308/a1/90-final.md` §1 表格＋本腿 `logs/pre-*-full.txt` 各 `panic hits=0`）。⇒ 守卫牙⛔ 挂在"机器台面"上，挂**输入层**：直接把无卷标字面量喂给 helper。

按裁语 ⓒ：三件一起做，下面 §1 §2 §3 逐件落。

---

## 1. 共享 helper（新件 `internal/risk/pathresolver_volumeless_spelling_308_test.go`，⛔ build tag ⇒ 所有档位编译）

两个函数，尺＝`grep -n '^func ' internal/risk/pathresolver_volumeless_spelling_308_test.go`：

- `driveShare(vol string) (share string, ok bool)`：`len(vol)==0` ⇒ 返回 `("", false)`；否则 `strings.ToLower(vol[:1])+"$"`。★ 唯一"取盘符"动作面，`[:1]` 只出现在这条 `len` 守卫之后（尺＝`grep -n 'if len(vol)' …`＋`grep -n 'vol\[:1\]' …`，行 24 / 行 27）。
- `volumeInvariantSpellings(home, long string) (keep map[string]string, dropped []string)`：7 枚⛔ 依赖卷标的拼法（long / upper / lower / extended / slashes / dot-and-dotdot / not-yet-there）恒在 `keep`；`driveShare` 报 ok ⇒ 追加 `unc-admin`、`unc-extended`；报 ¬ok ⇒ `dropped=["unc-admin","unc-extended"]`，`keep` 保持 7 项。⛔ panic、⛔ Skip、⛔ 整枚去掉。

---

## 2. 两枚崩点件按内容锚改为经 helper（尺＝`git diff -U0` 两枚文件，非 0 行只此两处）

- `pathresolver_anchor_spelling_windows_test.go` 内 `TestClassifySpellingInvariance`：
  删去 `vol := filepath.VolumeName(long)` / `share := strings.ToLower(vol[:1]) + "$"` / `rest := strings.TrimPrefix(long, vol)` 三行与 9 枚 inline 拼法 map；换 `spellings, dropped := volumeInvariantSpellings(home, long)`；`len(dropped) > 0` 时 `t.Logf` 具名"去掉的是哪两项＋还剩哪几项"（分母显式，⛔ 让下一位以为全跑）。⛔ 动后续 `if s := getShortPath(...)` 那枚 8.3-short 折叠、⛔ 动主循环、⛔ 动 controls 段。
- `pathresolver_junction_windows_test.go` 内 `TestPathResolverUNCAListDenied`：
  删去 `share := strings.ToLower(vol[:1]) + "$"`；换 `share, haveShare := driveShare(vol)`；`!haveShare` ⇒ `t.Fatalf(...)` 具名"⛔ 卷标 ⇒ 三枚 UNC 拼法无法构造"。⛔ Skip（票面 AC#1 逐字禁）、⛔ panic；三枚 UNC 拼法构造式与断言体本身⛔ 变（`\\?\UNC\localhost\`+share+rest / `\\localhost\`+share+rest / `//localhost/`+share+rest 逐字仍在；尺＝`git show HEAD:…junction… | grep -c localhost`＝4，其中 3 枚＝拼法字面量、1 枚＝用例顶部 `Case 3` 注释里同名出现）。

★ 两枚件现量⛔ 一处 `[:1]` 未守卫（尺＝`grep -rn '\[:1\]' internal/risk`＝8 行：`pathresolver.go:263` 有 `len(drive)==2 && drive[1]=='$'` 守卫；`pathresolver_volumeless_spelling_308_test.go:27` 是 helper 里唯一"取盘符"动作、行 24 `if len(vol)==0 {return "",false}` 之后；其余 6 行⛔ 我件⛔ 释旧崩形的注释文本，⛔ 一行可执行码）。

---

## 3. 直接喂无卷标字面量的单元用例（`TestVolumeInvariantSpellingsVolumeLess308`，⛔ build tag ⇒ 所有档位跑）

- 输入＝`home := "/tmp/wisp-308-no-volume"` 与 `long := filepath.Join(home, ".ssh", "id_testkey")`。★ 边界：字面量是**喂给 helper 的输入**（`§1` 裁语 ⓒ ⚠），⛔ 换成既有用例里的 `t.TempDir()`（票面硬约束＋票 72 `:201` 同形禁线）。
- 断言三件事（票面 `AC#1` 逐字）：
  ①⛔ panic ＝ 用例自己以 `PASS/FAIL` 结束（尺＝`grep -E '^--- (PASS|FAIL).*TestVolumeInvariantSpellingsVolumeLess308' logs/post-restored-default-full.txt`，命中 `PASS`）；
  ② `dropped` 逐名＝ `[unc-admin unc-extended]`，`keep`⛔ 含二者；用例 `t.Logf` 具名"去掉哪两项＋还剩哪几项"（分母可见：7 kept＋2 dropped；尺＝`grep -A1 'TestVolumeInvariantSpellingsVolumeLess308' logs/post-default-full.txt`）；
  ③ 长/大写/小写/扩展那四种⛔ 消失（连同 slashes/dot-and-dotdot/not-yet-there 全 7 项），逐项 `keep[name]` 必 present（尺＝`grep -n '"long", "upper"' internal/risk/*_308_test.go`）。
- 前置在场尺：`filepath.VolumeName(long)` 必为 `""`，否则 `t.Fatalf` 具名"premise broken"，⛔ 用绿色冒充"跑过守卫"。

实测三枚颜色（尺＝`logs/post-{default,alias}-full.txt`，`grep -E '^--- (PASS|FAIL).*…SpellingInvariance|…UNCAListDenied|…VolumeLess308'`）：default 台面三枚全 PASS、alias 台面三枚全 PASS、`panic hits=0`、顶层 `PASS=132`（前＝131，＋1＝本用例）⛔ 少任何一枚、⛔ 多任何一枚⛔ 我名。

---

## 4. 覆盖⛔ 搬家式消失（票 §3 硬约束）

`TestClassifySpellingInvariance` 在有卷标台面下求值 9 枚（＋可能 8.3-short），helper 输出与旧版逐字等值（`unc-admin`/`unc-extended` 表达式⛔ 变）；⛔ 台面无卷标那形求值 7 枚。★ 分母在日志里，⛔ 让下一位误当"全跑"（尺＝`grep -c 'still evaluated' logs/post-restored-default-full.txt`＝1，命中即 helper 走 ¬ok 分支那次；本机⛔ 到那分支，本机 `dropped=∅`、那行⛔ 触发）。

`TestPathResolverUNCAListDenied` 三枚拼法在有卷标台面下逐字⛔ 变（改前后同集合，尺＝`logs/post-{default,alias}-names.txt` ⇄ `logs/pre-{default,alias}-names.txt` 名⇄色一致）；无卷标那形⛔ 覆盖搬家式消失＝用例自己红（`t.Fatalf`），⛔ Skip。

---

## 5. ⛔ 做到什么（具名）

- ⛔ 让无卷标台面在本机现量成立（裁语 ⓑ：`GetTempPath` 恒给 drive-qualified）。守卫的牙只在 §3 那枚直喂字面量的单元用例里现量。
- ⛔ 改任何判级逻辑（`Classify`/`Gate`/`Resolve` 函数体本身）。
- ⛔ 动产码：`pathresolver.go`/`pathresolver_windows.go` 及任何⛔ `_test.go` 文件。
- ⛔ `t.Skip`（票面 AC#1 逐字禁），⛔ `t.TempDir()` 换成手写常量。
- ⛔ 为变绿放宽断言（`dropped` 长度、`share` 值、`keep` 键集都是**加**的强断言，⛔ 减）。
