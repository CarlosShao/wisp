# 72-r1 乙组 2 枚负向钉（票 307 AC#2 欠账，同批归本票）

写面＝`internal/risk/syncdirs_test.go`。产码⛔ 动。

归口原文＝`307-…-done.md:218`（"那枚『乙族 2 枚⛔ 负向钉』的欠账（体内只有正向断言，恒真形下它们绿）**归票 72 下一程**"，撤销口令「307 负向钉另立一票」）＋票 72 追加节 §4（`72-…-runner.md:288-290`）。⚠ 台账写的行号是 `syncdirs_test.go:209/:212/:216`、`:229/:232`，本腿 `sed` 现量红句在 `:213`/`:217` 与 `:233` ⇒ **认内容锚**（`"under-profile fallback must stay armed"`／`"writes under the injected root must be sync"`／`"suspect fallback must stay armed"`），⛔ 只信行号。

成因＝这两枚体内**只有** `Sync:true` 那一族断言：把 `match()` 的下述分支改成恒真（`syncdirs.go:442` 那枚 `if !s.complete && s.home != "" && aUnder(candForms, anchorForms(s.home)) {` → `if true {`），两枚**全绿**＝仪器哑。

---

## 乙-1 `TestSyncFallbackNotDisarmableByWeakRoot`

新增夹具（在函数体开头，正向断言之所以⛔ 受影响）：
```go
	// A tree that is neither under the fake profile nor under the injected
	// root. It lives one level up (sandbox's home is <base>\profile), so no
	// membership test and no under-profile fallback may reach it.
	outside := filepath.Join(filepath.Dir(home), "not-a-sync-tree", "notes.md")
```
新增负向断言（逐字，`for _, src := range []string{"default", "fixture", "options", ""}` 循环内、三枚正向断言之后，四档 src 各问一次）：
```go
		// NEGATIVE NAIL (ticket 307 AC#2 debt, handed to ticket 72 乙组): the
		// three assertions above are all of the shape "this must be Sync:true",
		// so a detector that answered Sync:true for every write on the volume
		// passed all three. The fallback is a BOUNDED net under the profile,
		// not a blanket deny -- pinning the miss is what makes "not disarmable"
		// a claim about a rule instead of a claim about a constant.
		if st := p.IsSyncPath(outside); st.Sync {
			t.Errorf("source %q: %q is outside the profile and outside the injected root, must NOT be sync: %+v", src, outside, st)
		}
```
它响的那一发＝`50-mutations.md` 反形 ⓑ（恒真形，两台面上 `TestSyncFallbackNotDisarmableByWeakRoot` 转红、正向三枚仍绿＝只有新负向那枚在响）。原三枚正向断言**逐字未动**。

## 乙-2 `TestSyncUnverifiedRootKeepsFallback`

逐字：
```go
	// NEGATIVE NAIL (ticket 307 AC#2 debt, ticket 72 乙组): armed is not the
	// same word as everywhere. A target outside the profile and outside the
	// unverified root must still read Sync:false, or this case would also be
	// satisfied by a detector that flags every write it is asked about -- the
	// shape that keeps both assertions above green.
	if st := p.IsSyncPath(filepath.Join(filepath.Dir(home), "not-a-sync-tree", "x.md")); st.Sync {
		t.Fatalf("a write outside the profile and outside the unverified root must not be flagged: %+v", st)
	}
```
为什么这一枚⛔ 撞上那枚"根"：注入的根是 `home\OneDrive\..\..\definitely-missing-9f3a`（C26 ⛔ 能验证 ⇒ `SyncDetectionComplete()` 必须 false），它在盘上⛔ 存在；负向靶落在 `sandbox` 基目录的另一个兄弟树上，与该根既⛔ 同源也⛔ 同前缀（`not-a-sync-tree` 对比 `definitely-missing-9f3a`），所以"⛔ 被恒真形绿"与"⛔ 被那条根绿"两件事分开钉住。`..` 只在**注入根**的字符串里（那是被测的那件事），负向靶自身⛔ 含 `..`（含 `..` 会命中 `syncdirs.go:395` 的 N-10 fail-closed，那不是我要问的）。

两台面上这两枚**都绿**（读数＝`30-evidence.md` §1 的 sync 族发），恒真形下**都红**（`50-mutations.md` ⓑ）。⛔ 新增 `t.Skip`／build tag；⛔ 动产码；⛔ 拔掉正向断言（三枚正向 + `SyncDetectionComplete` 那两枚**逐字保留**）。
