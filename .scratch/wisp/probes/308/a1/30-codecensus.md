# 308-a1 · 30-codecensus（Ⓑ：产码⛔ 在射程内的全仓 `[:1]` 名册表）

尺＝`grep -rn '\[:1\]' internal cmd --include=*.go` ⇒ 7 枚命中（全仓含 `tools` 亦 7 枚，⛔ 更多）。逐枚看切片前面有没有守卫、属不属于本票射程（射程＝两枚 `internal/risk` 测试文件内的 `vol[:1]`）。

| # | 件:行 | 切片表达式 | 被切的是 | 前面有守卫？ | 属本票 308 射程？ |
|---|---|---|---|---|---|
| 1 | `internal/risk/pathresolver.go:263` | `drive[:1]` | `drive`（来自 share 切段） | **有**：`:262` `if len(drive)==2 && drive[1]=='$' && ('a'-'z' 或 'A'-'Z')` | ⛔ 是（产码、且守卫后、恒 len==2 才切） |
| 2 | `internal/risk/pathresolver_anchor_spelling_windows_test.go:103` | `vol[:1]` | `vol := filepath.VolumeName(long)`（`:102`） | **⛔ 无** | **是＝崩点 #1**（`TestClassifySpellingInvariance`；`vol==""` ⇒ panic） |
| 3 | `internal/risk/pathresolver_junction_windows_test.go:152` | `vol[:1]` | `vol := filepath.VolumeName(long)`（`:151`，注 `// e.g. C:`） | **⛔ 无** | **是＝崩点 #2**（`TestPathResolverUNCAListDenied`；与 #2 同形） |
| 4 | `internal/agent/spill_path_invariant_test.go:389` | `abs[:1]` | `abs := filepath.Join(userDir,"diary.txt")`（`:386`） | 末元素是字面量 `diary.txt` ⇒ `abs⛔` 空串（`Join` 非空）；但**⛔ 是 `VolumeName` 形** | ⛔ 是（另一枚文件、另一族机制：切路径首字符、⛔ 切卷标） |
| 5 | `internal/agent/spill_path_invariant_test.go:406` | `produced[:1]` | `produced`（`[]string`，`:379` 循环 append） | ⛔ 无显式 `len` 守卫（靠 fixture 非空）；⛔ 是字符串切片形 | ⛔ 是（非 `vol[:1]` 崩形；别件） |
| 6 | `internal/llm/openaichat/mockllm_integ_test.go:190` | `rs[:1]` | `rs`（切片，喂 `golden.NewReplayer`） | n/a（切片首元素；与 VolumeName 无关） | ⛔ 是（非路径卷标崩形；别件别包） |
| 7 | `cmd/wisp/firstrun_198r2_test.go:401` | `word[:1]` | `word := m[1]`（`:396` 正则捕获组） | 靠 `seen[word]` 去重、**⛔ 长度守卫**（正则捕获空串理论上会崩，属另一枚潜在形） | ⛔ 是（`cmd/wisp` 包＝本腿⛔ 打的面；非 risk 射程） |

## 产码面两条（票面 §现量 已给，本腿当场复跑坐实）

- `internal/risk/pathresolver.go` 的 `drive[:1]`（#1）：切片前 `:262` 那条 `len(drive)==2 && drive[1]=='$'` 守卫**在场** ⇒ 只有恰两字符且第二位 `$` 的段才进去切 ⇒ ⛔ 会 `[:1] with length 0`。
- `internal/risk/pathresolver_windows.go:60 reparseComponents()`：`:61 vol := filepath.VolumeName(p)` 取了 `vol`，但**体内⛔ 做 `[:1]` 切片**（尺＝`grep -n 'reparseComponents\|VolumeName\|\[:1\]' …pathresolver_windows.go`，`[:1]` 在该文件 **0 命中**）⇒ 产码面⛔ 有本票要防的崩形。

## 射程结论

- **本票 308 射程 ＝ #2 与 #3**（两枚 `internal/risk` 测试文件里、`vol := filepath.VolumeName(…)` 后⛔ 守卫的 `vol[:1]`）。
- **#1／#4／#5／#6／#7 ⛔ 属本票射程**：#1 是守卫后的产码；#4 末元素字面量⛔ 会空、且切的是路径首字符⛔ 是卷标；#5／#6 是切片首元素、⛔ 是 `VolumeName` 崩形；#7 在 `cmd/wisp`（本腿⛔ 打那面）。
- ⇒ 印证票面 §现量"产码⛔ 在射程内"：**崩只在测试面两枚同形（#2/#3），产码（#1）有守卫、`reparseComponents` ⛔ 切片**。
