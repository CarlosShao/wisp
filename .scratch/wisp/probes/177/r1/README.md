# 177-r1 台件（常驻判据腿的变异读数与门禁日志；临时件只建不删）

- 判据本体＝`internal/risk/shape_a_exemption_test.go`（提交 `a19fa2e7` 及其后的成对更正一枚）。
- 逐枚归因表＝`docs/evidence/s1/177-shape-a-impl-r1.md` §2。
- 变异复跑口径（每次：临时改产码→`grep -n MUT` 证落地→下面这条尺→还原→`git diff --exit-code -- internal/risk/provenance.go` 证还原）：

```
go test -count=1 -run 'TestShapeA' -v ./internal/risk/ 2>&1 | grep -E '^--- '
```

- 各变异的一次性配方（本目录**不存产码 diff**，防被照抄落地；配方即证据）：
  - `MUT-W1`＝`MarkWithHostPath` 里 span 换成 `{{0, len(rp)}}`（整段不索引）
  - `MUT-W2`＝包级 `[]string` 值名册＋每枚 mark 用 `runeIndexOf` 按值开窗（按值放行）
  - `MUT-W3`＝同 W2 但名册存 `hp[:strings.LastIndex(hp,"/")]`，位置记录全程为零（目录前缀豁免，归因干净版）
  - `MUT-W4`＝包级 `[][2]int` 跨度名册，`m := &taintMark{` 前 `skip = append(skip, 名册...)`（按位置推给全树＝177-m1 的 S3a，旧全仓 rc=0 那形）
  - `MUT-W5`＝排除后 `kept==0 → return false`（空判定排后＝S5）
  - `MUT-ORDER`＝跨度定位改 `runeIndexOf([]rune(content), []rune(hostPath))`（原始下标，违反顺序(i)）
  - `MUT-NONE`＝`skip = nil`（什么都不豁免，A 正控腿的反证）
- `logs/` 逐名读数：`newtests-green.txt`（未变异全绿）、`mut-*.txt`（七发）、`restored-green-w5.txt`、`final-d22scan.log`、`final-gates.log`（G5 漂移修复后的那发）、`final-runtests.log`、`final-pass.txt`（与 `probes/177/m1/r-pass-mut.txt` `comm -3` 空差集的右侧）。
- `logdir` 口径：本程无新脚本；日志一律以仓库根为基准的显式相对路径写入本目录 `logs/`。
