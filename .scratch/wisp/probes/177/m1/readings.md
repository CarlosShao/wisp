# 177-m1 台件与读数（变异台件，产码一枚不入提交；本目录只存台件文本与原始日志）

起手锚（`date` 现量）：`Sun Sep 27 22:29:27 CST 2026`｜branch `dev`｜`git rev-parse HEAD` = `9b4107530df050a2a8d6a918b2b38c1ccd4fbce0`
还原后同一棵树复跑，`git status --porcelain -- internal/ cmd/` 为空。

## 1. 甲形变异的本体（三枚被改文件的改法；全部已用 `git cat-file blob HEAD:<path> > <path>` 还原）

### 1.1 `internal/risk/taintmatch.go`（建索引那一圈）

```go
// newFragmentIndex 形参从 (norm string, n int) 改为 (norm string, n int, skip ...[2]int)
// （S1/S2/S3/S5 用变长参数，S4 单独测"必选第三枚参数"那枚载具）
func newFragmentIndex(norm string, n int, skip ...[2]int) *fragmentIndex {
	rp := []rune(norm)
	// ...原有 n 钳制与 len(rp)<n 分支不动...
	for i := 0; i+n <= len(rp); i++ {
		if m1WindowSkipped(i, n, skip) { continue }   // 只跳过被声明的窗口
		copy(win, rp[i:i+n])
		hashes = append(hashes, hashWindow(win))
	}
```

- `f.src` 全文与 `contains()` 里的 `strings.Contains(f.src, w)` 复核**一字未动**。
- 新增件：`m1WindowSkipped(i, n, skip)`（窗口与 `[lo,hi)` 是否相交）、`runeIndexOf(hay, needle []rune)`（归一化坐标里定位 P）、
  `m1GlobalOn`/`m1GlobalSkip`/`m1GlobalVals`/`m1ValueSkips`（只在 S3 系用来模拟"豁免跨出它自己那枚 mark"）。

### 1.2 `internal/risk/provenance.go`（Mark 侧；导出面在 S1/S2/S3/S5 全部未动）

```go
func (p *Provenance) Mark(scopeID, tool, origin, content string) bool {
	return p.markWithHostPath(scopeID, tool, origin, content, "")   // 载具＝一枚未导出参数
}

func (p *Provenance) markWithHostPath(scopeID, tool, origin, content, hostPath string) bool {
	norm := normalizeTaint(content)
	if norm == "" { logf(...); return false }        // 空判定排在排除之前（原 :475-479 位置不动）
	rp := []rune(norm)
	// 原有 short-source / maxSrc 截断 / IsSensitiveSource 日志三支原样保留
	var skip [][2]int
	if hp := normalizeTaint(hostPath); hp != "" {
		nhp := []rune(hp)
		if lo := runeIndexOf(rp, nhp); lo >= 0 {      // 跨度定位在归一化之后
			span := [2]int{lo, lo + len(nhp)}
			skip = append(skip, span)
			m1RecordSpan(span)
			logf("risk/C25: 177-m1 shape A excludes normalized span [%d,%d) ...", span[0], span[1])
		} else { logf("... host-minted path absent ... nothing excluded") }
	}
	m := &taintMark{tool: tool, origin: origin,
		idx: newFragmentIndex(string(rp), p.minChars, skip...), at: observe.NowWallUTC().Unix()}
	// 其余（入 scope / maxMarks 日志 / truncated 日志）原样
}
```

S5（行 2 的反例）在 span 之后插了一枚 `if lo == 0 && span[1] == len(rp) { return false }`。
S4（载具对照）把 `newFragmentIndex` 的第三枚参数改成**必选** `skip [][2]int`。

### 1.3 `internal/risk/provenance_test.go`（六发台件，追加在文件末尾）

```go
const (
	m1Path         = `C:/Users/x/AppData/Roaming/wisp/artifacts/tool-output-call-177-m1.txt`
	m1Ghost        = `C:/Users/x/AppData/Roaming/wisp/artifacts-ghost/tool-output-call-177-m1.txt`
	m1Stub         = "head SUPERSECRET-177M1 tail […输出已落文件，全文见 " + m1Path + "…]"
	m1Foreign      = "请帮我复核这份清单：" + m1Path + " 里是原始输出"
	m1ForeignGhost = "请帮我复核这份清单：" + m1Ghost + " 里是原始输出"
)
```

| 发 | 测试名 | 断言 | 对应上游件 |
|---|---|---|---|
| A | `TestM1RigALegReReadOfHostPathIsClean` | 宿主自产 mark 声明 P 后，`Inspect("notify",{text:P})` 必须干净 | 夹具 R-2 正向对照 |
| B | `TestM1RigBRestOfHostMarkStillIndexed` | 同那一枚 mark 的**非 P 片段** `SUPERSECRET-177M1` 仍须命中 | 防"恒不变"的正控 |
| C | `TestM1RigCLegForeignMarkStillHits` | 外来 mark 逐字带同一条 P ⇒ `Inspect(P)` 仍须命中 | 夹具 R-1 反向判据 |
| D | `TestM1RigDLegNeverMintedPathStillHits` | 外来 mark 带的是宿主**没写下过**的路径 ⇒ 仍须命中 | 夹具 R-3 过界哨 |
| E | `TestM1RigEPathOnlyMarkStillRecorded` | 正文恰好只剩 P：mark 仍须被记录，且枚数＝1 | 上游 §4.1 行 2＋行 3 的前提检查 |
| F | `TestM1RigFExemptionDiesWithScope` | `Scope.Close` 之后 mark 与排除一起消失（`ScopeTaints==0`、无命中） | 上游 §8 d 生命周期 |

每发起手 `m1Reset()`＋`defer m1Reset()` 清 `m1GlobalSkip`/`m1GlobalVals`，免得跨测试污染别人的读数。

## 2. 状态与读数（每条都可复制）

| 态 | 改法 | 命令 | 读数 |
|---|---|---|---|
| 基线 | 未动 | `go test -count=1 ./internal/risk/ ./internal/tools/` | `ok risk 5.616s` / `ok tools 17.980s`（顶层 2 枚包 PASS；`-v` 口径 risk `=== RUN` 172＝171 PASS＋1 SKIP(`TestSyncRegistryProbeLive`)，tools 173） |
| S1 | 甲形·精确（未导出载具） | `go test -count=1 ./internal/risk/ ./internal/tools/`；`go test -count=1 -v -run 'TestFragmentMatchExactBoundary\|TestFragmentIndexShortSourceNeverMatches\|TestFragmentHashCollisionCannotFakeHit\|TestMarkEmptyAfterNormalization\|TestMarkUnknownSourceToolFailClosed\|TestConcurrentMarkInspect\|TestScopeCloseRefusesAnotherOwnersRegistration\|TestM1Rig' ./internal/risk/` | 两包 rc=0；13 枚点名用例逐名 PASS（含六发台件全绿） |
| S2 | 排除放宽成"整段 mark 不入索引" | 同 S1 的第二枚命令 | 仅 `TestM1RigBRestOfHostMarkStillIndexed` FAIL（`rc=1`；包级 `FAIL risk 5.450s`，`ok tools 15.881s`） |
| S3a | 只把排除**按位置**推给全树所有 mark | 同 S1 的第二枚命令 | 六发全 PASS ⇒ 这放宽法打不动夹具（读数本身是条否定结论） |
| S3b | 排除**按值**推给所有 mark（名册形；已关掉位置记录以做归因） | `go test -count=1 -v -run 'TestM1Rig\|TestMarkUnknownSourceToolFailClosed\|TestMarkEmptyAfterNormalization\|TestConcurrentMarkInspect' ./internal/risk/` | `TestM1RigCLegForeignMarkStillHits` FAIL；A/B/D/F PASS；E FAIL 是 S5 那一支混在同一树里 |
| S3c | 名册里放的是**目录前缀**（同一棵树里的那一发 D） | 同上（合并跑） | `TestM1RigDLegNeverMintedPathStillHits` FAIL，但关掉位置记录后转绿 ⇒ D 的红**不能**记在"前缀豁免"账上，见证据件 §5 末注 |
| S5 | 空判定排在排除**之后** | 同 S1 的第二枚命令 | `TestM1RigEPathOnlyMarkStillRecorded` FAIL（`provenance_test.go:1096`）；`TestMarkEmptyAfterNormalization` 仍 PASS |
| S4 | 载具换成"必选第三枚参数" | `go test -count=1 ./internal/risk/ ./internal/tools/` | `FAIL risk [build failed]`；`ok tools 15.925s`；编译错 7 处：`taintmatch_test.go:31/:51/:66/:108/:124/:134/:148`（尺：`grep -cE 'not enough arguments in call to newFragmentIndex'` → 7） |
| 还原 | 三枚文件逐枚 `git cat-file blob HEAD:<path> > <path>` | `git status --porcelain -- internal/ cmd/`；`git diff --numstat -- internal/ cmd/`；`go test -count=1 ./internal/risk/ ./internal/tools/` | 两条 git 命令均无输出；两包 `ok risk 5.513s` / `ok tools 15.410s` |

原始日志：本目录 `s1-pkg.log`／`s2-pkg.log`／`s3b-targeted.log`／`s3b-isolated.log`／`s3s5-targeted.log`／`s4-pkg.log`／`s4-vet.log`／`s1-d22scan.log`／`s1-gates.log`／`s1-runtests.log`／`r-d22scan.log`／`r-gates.log`／`r-runtests.log`／`r-pass-mut.txt`／`r-pass-res.txt`。
`logdir` 取脚本自身目录：本程没有新脚本，日志一律用相对仓库根的显式路径写进本目录（未依赖 CWD 之外的默认落点）。
