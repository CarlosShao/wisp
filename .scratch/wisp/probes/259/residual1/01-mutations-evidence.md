# 259residual-1 两发突变实测证据（发 A 拆 store 成员扫描／发 B state 卫）

- 料文＝`.scratch/wisp/probes/next-instruments-brief-1/02-ticket259-ac2-brief.md`（备料腿 `next-instruments-brief-1`，commit `a6ec48d3`）。
- 跑法（三跑同）：`cd internal/agent/approval && PATH="$PWD/../../../third_party/sherpa-onnx:$PATH" go test . -run '259' -v -count=1`
- PATH 具名更正：派单模板两级（`../../`）解析为不存在的 `internal/third_party`（现核 No such file or directory）；本腿按三级 `$PWD/../../../third_party/sherpa-onnx` 跑。三跑均实打 12 条 `=== RUN`（非 0 条的 0xc0000135 死法）。
- 三跑均先打一行 `winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2` INFO 日志（装尺噪声，非红）。
- 基线（未突变，对照组，先于两发）：rc=0，12 枚全 `--- PASS`，末行 `ok  github.com/CarlosShao/wisp/internal/agent/approval  0.049s`。

## 发 A｜拆 store 成员扫描（`approval.go:565`）

1. 种前 hash：`67fb146898dc7b235623bb8ad2fac0e0b2465fe1`（与料文同；种前 scoped porcelain 空，种前 `sed -n '565p'`＝`\tif !equalSecret(v, nonce) {`）。
2. 盘上改行复量（种法 `sed -i '565s/if !equalSecret(v, nonce) {/if false {/'`）：
```
== post 564-567 ==
	for v, stored := range s.values {
		if false {
			continue
		}
== diff numstat ==
1	1	internal/agent/approval/approval.go
```
3. 跑测试：rc=1（`MUTA_RC=1`）。红句逐字（file:line 与短语同一发现取）：
   - `ticket259_denial_rulers_test.go:139: AC#2 RED: the never-issued-proof route returned <nil>, want ErrBadGrant`（`--- FAIL: TestTicket259R1DenialNamesSpentOrNeverLiveNonce`）
   - `ticket259_denial_rulers_test.go:259: AC#2 RED: never-live proof returned <nil>, want ErrBadGrant`（`--- FAIL: TestTicket259R1OutwardAnswerStaysMergedAcrossDenials`）
   - 连带第三枚（料文未预判；其 §4.3 第 3 条自认"未核两发是否连带打红别的在册项"）：`--- FAIL: TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer`，两行：
     - `ticket259_panel_capability_rulers_test.go:301: AC#3 RED: a value read out of the panel face spent as a native grant and allowed the card; the outward read face and the answer face are not the same channel, and this is the case that says so`
     - `ticket259_panel_capability_rulers_test.go:308: AC#3 RED: the panel-observable values were refused, but so was the live native grant (approval: correlation_id 无对应待审批项) - the card stopped being answerable before the control ran, which makes the pass above meaningless`
   - 其余 9 枚 `--- PASS`。
   - 料文预判核对：`:139`／`:259` 命中＝**对**（逐字）；"两枚"口径按测试枚数成立、按红行枚数不成立（红行 4 条、涉 3 测试）⇒ 具名报回。
4. 还原（`git show HEAD:internal/agent/approval/approval.go > internal/agent/approval/approval.go`）：hash＝`67fb146898dc7b235623bb8ad2fac0e0b2465fe1`（逐字等）；`git status --porcelain -- internal/agent/approval/approval.go` 空。

## 发 B｜state 卫（`queue.go:390`）

1. 种前 hash：`66fec7ae375344d7bc741270b42c4ec5050cf41f`（与料文同；种前 scoped porcelain 空，种前 `sed -n '390p'`＝`\tif it.state != statePending {`）。
2. 盘上改行复量（种法 `sed -i '390s/if it.state != statePending {/if false {/'`）：
```
== post 389-399 ==
	}
	if false {
		// Ticket 259 AC#2: the denial is named in the audit even on the cause
		// the store never sees. Read the tool under the lock - this item is
		// settled, and a settled item is the one thing this queue must never
		// re-describe from a stale pointer after the lock is gone.
		settled := it.Dec.Tool
		q.mu.Unlock()
		q.logf("approval: GRANT-DENY corr=%s tool=%s denial=%s", corr, settled, denialNotPending.label())
		return ErrNotPending
	}
== diff numstat ==
1	1	internal/agent/approval/queue.go
```
3. 跑测试：rc=1（`MUTB_RC=1`）。红句逐字：
   - `ticket259_denial_rulers_test.go:208: AC#2 RED: the settled-card refusal burned the live grant; that branch never reaches the store`（`--- FAIL: TestTicket259R1DenialNamesCardNotPending`）
   - 其余 11 枚 `--- PASS`；无连带。
   - 料文两点核对：红枚位置 `:208` 命中＝**对**；"其余断言不红"＝**对**（测试走到 `:208` 才 Fatal，说明 `:200-201` err＝ErrNotPending 与 `:203-205` 审计行 `denial=card-not-pending` 都过了——与料文 §2 的"牙只长在 live() 断言上"一致）。
4. 还原（`git show HEAD:internal/agent/approval/queue.go > internal/agent/approval/queue.go`）：hash＝`66fec7ae375344d7bc741270b42c4ec5050cf41f`（逐字等）；scoped porcelain 空。

## 终态与登记

- 两发还原后联合复核：`git status --porcelain -- internal/agent/approval/approval.go internal/agent/approval/queue.go` 空；`git diff --numstat --` 同两文件空。
- 还原期间他腿（242-corrcensus-1）在共享树落了新 commit `0e0ac7b9`（本腿取锚后）；两次还原均按当时 HEAD 取 blob 且 hash 逐字相等 ⇒ 该两文件 blobs 全程未被第三方改动（README 规则 8：临时件只建不删，本腿无删除动作）。
- 逐发判：**发 A＝料文预判对（另报连带 1 枚，料文自认未预判）；发 B＝料文预判对（无连带）**。
- 新恒真面：**无**（两发今天都真红；发 A 的连带红是真实牙齿）。
- 写面声明：本腿只改过 `approval.go:565` 与 `queue.go:390` 两行且已逐字还原；三枚冻结件未碰（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`——射程判断，非内容引用）；未读未引 `frontend/**`／`design/**`；零翻框；零 push。
