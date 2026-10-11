# 307-a1 · 30-test-nails —— 改完可能颜色变哪些既有钉（逐枚名＋尺）

工作树；`internal/`＋`cmd/` 对 HEAD 干净（`git status --porcelain -- internal cmd`＝0 行，见 `00-anchor.md`）。
凡"会红／会绿"＝〔读码推到〕，本腿⛔ 跑过任何 `go test`；本文件只交**名册＋归属＋依赖哪一枚折法**。

## 尺（三条，可逐字重跑）

- 尺①（票面点名的 5 枚断言行现跑）：`grep -n "write under fixture root must be sync\|under-profile fallback must stay armed\|suspect fallback must stay armed\|under-profile must be suspect\|under-profile path must be sync-suspect" internal/risk/syncdirs_test.go` ⇒ 命中 **162 / 213 / 233 / 355 / 402**，**与票面同号**（这 5 枚没腐烂，腐烂的是产码那两枚）。
- 尺②（`home` 折法的消费面）：`grep -rn --include=*.go "isUnder(cand" internal/risk/*.go` ⇒ `syncdirs.go:422`（`e.canon`）、`syncdirs.go:428`（`s.home`）。
- 尺③（负向断言名册）：`grep -rn --include=*_test.go -e "not be suspect" -e "must NOT be sync" -e "if st.Sync" -e ".Sync {" internal/risk`（输出全量在 `logs/g-i-negatives.txt`）。

---

## A. 票面点名的 5 枚（＝乙族），逐枚名＋所属函数＋依赖

| 断言行（快照，尺①） | 逐字消息 | 所属函数（func 起始行现跑） | 走哪枚 `home` 赋值点 | 依赖的折法 |
|---|---|---|---|---|
| `syncdirs_test.go:162` | `t.Fatal("write under fixture root must be sync")` | `TestSyncFixtureFallbackAndMatch`（`:152`） | `provenance.go:345`（`NoProbe:true, HomeDir:home, SyncFixturePath:fx`） | **两枚都吃**：`:422`（fixture 根 `home\Cloud\NutstoreSync` 从未创建 ⇒ `add()` 落 `:160` 短形）＋ `:428` |
| `:213` | `t.Errorf("source %q: under-profile fallback must stay armed", src)` | `TestSyncFallbackNotDisarmableByWeakRoot`（`:198`） | `provenance.go:345` | `:428`（`home` 侧） |
| `:233` | `t.Fatal("suspect fallback must stay armed")` | `TestSyncUnverifiedRootKeepsFallback`（`:224`） | `provenance.go:345` | `:428` |
| `:355` | `t.Fatal("under-profile must be suspect")` | `TestSyncSuspectFallbackIsComponentBounded`（`:344`） | `provenance.go:345` | `:428` |
| `:402` | `t.Fatalf("under-profile path must be sync-suspect, got %+v", st)` | `TestSyncSuspectFallbackWhenUndetectable`（`:390`） | `provenance.go:345` | `:428` |

⇒ **读数（决定落点形状的那一条）**：5 枚**全部**经 `NoProbe:true` ⇒ 全部走 `provenance.go:345` 那一枚复制件，**没有一枚**走 `syncdirs.go:88`。若落点只折 `syncdirs.go:88`，这 5 枚在别名台面上颜色不动。〔现跑：尺＝`grep -rn "NewProvenance(ProvOptions{" internal/risk/syncdirs_test.go`，输出在 `logs/g-e-nails.txt`／`logs/g-h-prod.txt`〕

## B. 同族但票面名册没点名的 `home` 依赖位（现跑）

这些也吃 `:428`，改完可能一起动，⛔ 让落地腿"以为只有那 5 枚"：

| 位 | 原文锚 | 依赖 | 备注 |
|---|---|---|---|
| `syncdirs_test.go:174` | `if !p.IsSyncPath(filepath.Join(home, "LocalOnly", "a.txt")).Sync` | `:428` | 在 A 表第一枚同一个函数体内（M-3 那半边） |
| `syncdirs_test.go:501` | `if !p.IsSyncPath(filepath.Join(home, "Documents", "x.txt")).Sync` | `:428` ＋ `syncdirs.go:88` | **唯一一枚吃 `:88` 复制件的 fold 断言**（`ProvOptions{}`＝live env，`home := userHomeDir()` 于 `:496`）；它前置 `if !p.SyncDetectionComplete()`（`:500`）⇒ 生产/台面若检出确认根就整块跳过 ⇒〔读码推到〕**这枚天然是"可能静悄悄不执行"的那一族** |
| `syncdirs_test.go:506` | `if !p.IsSyncPath(filepath.Join(od, "anything.txt")).Sync`（`od := home\OneDrive`） | 视台面：根命中或 `:428` | 前置 `dirExists(od)` ⇒ 本机无 OneDrive 时不执行 |
| `syncdirs_windows_test.go:58` | `if !p.IsSyncPath(filepath.Join(root, "notes.md")).Sync` | `:422`（confirmed 根） | `TestSyncConfirmedGradesDisarmFallbackWindows`（`:40`）体内 |
| `syncdirs_windows_test.go:54` | `if p.IsSyncPath(filepath.Join(home, "Documents", "exfil.md")).Sync` | **反向**：靠 `!s.complete` 短路 | ⚠ 别名台面上这枚今天**为错误的原因绿**：`complete=true` 本该让 `:428` 不参与，而短形 `home` 又让它必然不命中 ⇒〔读码推到〕落点若同时动 `finalize()`/`canonical` 会先在这里现形 |
| `syncdirs_windows_test.go:92` | `if st := p.IsSyncPath(inside); st.Sync` | `!s.complete` 短路 | `TestWriteGatePlainTargetInsideProfileWindows`（`:75`），前置 `:88` `SyncDetectionComplete()` |
| `syncdirs_redteam_windows_test.go:179` | `if !p.IsSyncPath(filepath.Join(root, "Notes", "x.md")).Sync` | `:422` | Windows 档，体内用 `shortPathOf`/`tryShortPath`（`:29-:54`，无别名时 `t.Skipf`） |
| `syncdirs_other_test.go:99` | `if st := e.p.IsSyncPath(filepath.Join(e.home, "Documents", "x.txt")); !st.Sync \|\| st.Root.Source != "suspect-fallback"` | `:428` **且锁 `Root.Source`** | `//go:build !windows` 档 ⇒ Windows 台面不编译；**它是唯一一枚断言 fallback 出口 `Root.Source` 的用例**（本腿⛔ 找到断言 `Root.Path` 等于 `home` 的用例） |

## C. 负向钉（把 `home`/根折长＝判级变严 ⇒ 这些是会被弄红的位）

尺③全量在 `logs/g-i-negatives.txt`。名册（13 枚）：

`syncdirs_test.go:89`、`:178`、`:357`、`:375`、`:405`、`:447`；
`syncdirs_other_test.go:91`（`!windows` 档）；
`syncdirs_windows_test.go:54`、`:92`；
`syncdirs_ancestor_actable_leg_116_test.go:209`、`:212`；
`syncdirs_redteam_windows_test.go:163`；
`provenance_test.go:487`、`:627`。

判据形状（现跑读到的，逐字）：`:178/:405` 的消息是 `path outside the user profile must not be suspect`，路径由 `outOfProfile(t, home)` 造（`syncdirs_test.go:105-116`，体内自带 `isUnder(normPath(p), normPath(home))` 前提断言 ⇒ **⚠ 测试代码自己在用"短形对短形"的那把尺**，落点若改 `home` 折法而⛔ 改这枚前提，它仍按词汇形判"不在 profile 内"）；`:357` 是 N-5 那枚 `profileevil` 兄弟目录钉（component-bounded，**任何新折法⛔ 让它变成前缀命中**）。

⇒〔读码推到〕形 A／形 C 对 C 组应当无影响（只增命中、且并集⛔ 减命中）；形 B 因为改 `:422`，**风险正集中在 C 组**。⛔ 现跑证据在此，必须由编排者或落地腿两台面各跑一发。

## D. `SyncDetectionComplete(` 非测试调用者＝0 这一枚事实对影响面的意义

- 尺＝`grep -rn --include=*.go "SyncDetectionComplete(" internal cmd` ⇒ **22 行**（工作树，含注释行）。逐层拆（同一条尺的两个过滤）：定义 1 枚（`syncdirs.go:458`）、**非注释调用 14 枚**（13 枚在 `_test.go` ＋ 定义 1 枚）、纯注释提及 8 枚（`provenance_syncdirs_windows_test.go:37`、`provenance_test.go:457`、`syncdirs_ancestor_actable_leg_116_test.go:128`、`syncdirs_other_test.go:20,43`、`syncdirs_test.go:367,434`、`syncdirs_windows_test.go:18`）。
- 尺＝`grep -rn --include=*.go "SyncDetectionComplete(" internal cmd | grep -v "_test.go"` ⇒ **只剩定义那一行** ⇒ **非测试面调用者＝0**（`cmd/` 零命中）。
- 意义（读数，⛔ 结论）：
  1. `complete` 这枚标志**对包外不可观测**——它的唯一生产效果是 `syncdirs.go:428` 里那个 `!s.complete` 短路。⇒ 落点若只改"判定形"（`home`/`canon` 的折法），**不会**改变任何包外可见行为，包外面只有 `IsSyncPath(...).Sync/Root/Why`（生产消费者＝`provenance.go:928` 一枚，票面已现量）。
  2. 反过来：`complete` 是**测试可见的契约**，13 枚测试调用点里 8 枚是"必须为真／必须为假"的前置闸（上表 B 组 `syncdirs_windows_test.go:51,88`、`syncdirs_test.go:171,209,229,397,489,500` 等）。⇒ **任何顺手改 `finalize()`／`canonical` 判定（＝把某支根升/降级成 canonical）的落点，红会先出现在这批闸上，而不是出现在 B/C 那批判定上**。这是"最小诚实落点"该⛔ 碰 `finalize()` 的具体理由。
  3. 因 `cmd/` 零调用 ⇒ 本票⛔ 需要面板/CLI 侧改动，`AC#1`④ 那三枚冻结件与本腿读到的任何一位都无耦合〔现跑：尺＝`grep -rn --include=*.go "SyncDetectionComplete(\|IsSyncPath(" internal/panel internal/perm cmd` ⇒ 0 命中〕。

## E. 需要编排者车道跑的具名发数（本腿⛔ 跑）

- `go test ./internal/risk/ -count=1 -timeout 300s -v -run 'TestSyncFixtureFallbackAndMatch|TestSyncFallbackNotDisarmableByWeakRoot|TestSyncUnverifiedRootKeepsFallback|TestSyncSuspectFallbackIsComponentBounded|TestSyncSuspectFallbackWhenUndetectable'` ——本机台面一发、`TEMP` 指到 `C:\Users\swq\tmp\RISKAT~1\RUNNER~1`（长形 `...\riskattrib1\RunnerTempAreaForAttribution1`，本腿现跑复核它**仍在盘上**，见 `40-machine-alias.md`）一发。
- 上面 C 组那 13 枚名册的正控（同一台面两发，取改前/改后交集）。
