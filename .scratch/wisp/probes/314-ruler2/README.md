# 票 314 · 形状尺补洞（`w314-ruler2`）读数根

本目录是 `docs/evidence/s1/314-shape-ruler.md` §7 的原始读数。**所有突变都在仓外导出树做**，
跟踪件只改了 `cmd/wisp/panel_profile_leaf_shape_windows_test.go`（SHA256 见下）。

| 件 | 是什么 |
|---|---|
| `zzz_ruler2_probe_windows_test.go` | 能编译的探针：审计那 10 形（A2a–A2e、A3a/A3b/A3c、X1/X2）＋1 条反向控制 |
| `probe.ps1` | 仓外两棵树（`old`＝HEAD 的尺、`new`＝补洞后的尺）建 + 跑同一枚探针 |
| `teeth.ps1` | 仓外 `teeth` 树里逐条打断匹配臂，看常驻控制是否当场红，再还原 |
| `mutants2.log` | 探针读数：旧尺 **275 文件/2 命中/0 红 PASS**；新尺 **275/14/12 红 FAIL** |
| `mutants.log` | 第一次跑歪的读数（PowerShell 把 `-test.run=…` 拆成 `-test` + `.run=…`，rc=2）——见下面「两个坑」第 1 条 |
| `baseline.log` | 补洞后**第一次**现跑：尺子把自己控制表里两条裸绝对字面量判红（274 文件／4 命中／**2 实红**，`:1374`/`:1387`）——证据 §7.6 的处置由此而来 |
| `baseline2.log` | 处置完（裸绝对字面量改走 SINK）后的现跑：**274/2/2/0，PASS** |
| `final-green.log` | 定稿现跑：**274/2/2/0，PASS**（与 `baseline2.log` 同读数、同字节数 10356 B） |
| `teeth.log` | T1–T4 四条自残突变各让对应控制红；还原后 SHA 与仓内同串 |
| `d22scan.log` · `path-budget.log` | 两条门禁的完整读数 |

## 复跑（逐字）

```powershell
# 1) 仓外导出树：old = HEAD 的尺；new = 补洞后的尺；两边都放同一枚探针
$root='C:\Temp\wisp-ruler2-20261114'
New-Item -ItemType Directory -Force "$root\stage-old","$root\stage-new" | Out-Null
git archive HEAD -o "$root\src.zip"
Expand-Archive -Path "$root\src.zip" -DestinationPath "$root\old" -Force
Expand-Archive -Path "$root\src.zip" -DestinationPath "$root\new" -Force
Copy-Item -Force "cmd\wisp\panel_profile_leaf_shape_windows_test.go" "$root\new\cmd\wisp\"
Copy-Item -Force ".scratch\wisp\probes\314-ruler2\zzz_ruler2_probe_windows_test.go" "$root\old\cmd\wisp\"
Copy-Item -Force ".scratch\wisp\probes\314-ruler2\zzz_ruler2_probe_windows_test.go" "$root\new\cmd\wisp\"
Copy-Item -Force "third_party\sherpa-onnx\*.dll" "$root\stage-old\"
Copy-Item -Force "third_party\sherpa-onnx\*.dll" "$root\stage-new\"

# 2) 两棵树各自建 + 跑（DLL 必须挨着 exe，cwd 必须是 cmd/wisp）
pwsh -NoProfile -File .scratch\wisp\probes\314-ruler2\probe.ps1

# 3) 常驻控制有没有牙：四条自残突变
pwsh -NoProfile -File .scratch\wisp\probes\314-ruler2\teeth.ps1

# 4) 本仓现跑（只许这一条 -test.run，⛔ 跑整包）
sh .scratch/wisp/probes/314-ruler/run-ruler.sh '-test.run=^Test314'
```

## 两个 PowerShell 坑（把读数跑歪过一次，记下来）

1. `-test.v` 不引号会连带把前一个 `-test.run=…` 拆成 `-test` + `.run=…`（`flag provided but not defined: -test`）。
   **两个都加单引号**：`'-test.run=^Test314' '-test.v'`。
2. 测试二进制要求 cwd＝`cmd/wisp`：在仓根跑会读到「扫描集为空」并 Fatal（这是尺子故意的，不是故障）。

## 尺子 SHA256

- 补洞前（＝审计件钉的那一枚）：`7D0CB395303D09212DE920A9E24F3658F06E93FD48371B373A5A86C03F89A6E7`
- 补洞后：`6E54A0B34D6A3CD79888F5CD33535EB0754F4549DDD3D2FE1FD8FFE3CAA723F2`
