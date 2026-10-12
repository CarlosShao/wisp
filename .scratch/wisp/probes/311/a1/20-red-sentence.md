# 311-a1 `20-red-sentence` — 那两枚红用例的红句逐字＋那枚 SID（`AC#0` 问②）

读数件＝`logs/21-gotest-rootB-dtmp.txt`（根 B，换卷那发，`gotest_rc=1`）⇄ `logs/20-gotest-rootA-localappdata.txt`（根 A，默认那发，`gotest_rc=0`）。

## 1. 红名册逐字（根 B，`--- FAIL` 那四行＋它们各自的红句）

```
136:--- FAIL: TestGateJudgesThePrivateSetByResolvedSID (5.45s)
137:    --- FAIL: TestGateJudgesThePrivateSetByResolvedSID/runner_descriptor_is_judged_by_what_it_resolves_to (1.55s)
139:    --- FAIL: TestGateJudgesThePrivateSetByResolvedSID/a_name_spelling_cannot_change_the_verdict (1.93s)
143:--- FAIL: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (1.57s)
```

四枚红的**全部**消息内容（逐字，⛔ 截短；行号同件）：

```
130:    testing.go:1617: TempDir RemoveAll cleanup: openfdat D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSIDrunner_descriptor_is_ju1386323932\001\runner: Access is denied.
133:    testing.go:1617: TempDir RemoveAll cleanup: openfdat D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSIDa_name_spelling_cannot_2711238538\001\named: Access is denied.
135:    testing.go:1617: TempDir RemoveAll cleanup: openfdat D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSID3972594964\001\la: Access is denied.
142:    testing.go:1617: TempDir RemoveAll cleanup: openfdat D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594\TestSealNarrowsAndNamesThePrincipalItRemovedBySID1250499198\002\la: Access is denied.
```

★**票面那条前提要更正**（这是问②最要紧的一半答案）：这两枚用例的红句里**没有一句是判据说的**——
四行全是 `testing.go:1617` 的 `TempDir` 自我收尾错误（`c.Errorf`），⛔ 一条 `t.Fatalf` 触发过。
所以严格讲，"红句里那枚'解析出来却⛔ 在集合内'的 SID"**在那四行里⛔ 存在**（那四行只带路径＋`Access is denied.`）。
两枚用例内那几处会打印 trustee 名册的判据句（`:154`／`:158`／`:204`／`:207`／`:340`／`:348`／`:352`／`:359`）**一枚都没响**，尺＝对两发全文各数一遍关键字，逐字读数：
```
grep -c 'gate accepted a descriptor'      A=0 B=0
grep -c 'gate refused a descriptor'       A=0 B=0
grep -c 'decided differently by spelling' A=0 B=0
grep -c 'did not land on the same DACL'   A=0 B=0
grep -c 'seal left a foreign grant'       A=0 B=0
grep -c 'private_set_sid_windows_test.go:' A=1 B=1   ← 两发都只有 :314 那行 t.Logf
```
〔真机量的〕

## 2. 判据侧的读数在两枚根上逐字一致（＝"颜色变了"⛔ 来自判定）

同一条夹具自报（`private_set_sid_windows_test.go:314`），根 A ⇄ 根 B 逐字：

```
A(128): private_set_sid_windows_test.go:314: planted as strangers: [S-1-1-0 S-1-5-32-546 S-1-5-21-1228170099-895614386-1166154857-500]; planted as a private-set member: S-1-5-32-544 (LA=S-1-5-21-1228170099-895614386-1166154857-500, in set: false, set=[S-1-5-21-1228170099-895614386-1166154857-1001 S-1-5-18 S-1-5-32-544])
B(141): private_set_sid_windows_test.go:314: planted as strangers: [S-1-1-0 S-1-5-32-546 S-1-5-21-1228170099-895614386-1166154857-500]; planted as a private-set member: S-1-5-32-544 (LA=S-1-5-21-1228170099-895614386-1166154857-500, in set: false, set=[S-1-5-18 S-1-5-32-544 S-1-5-21-1228170099-895614386-1166154857-1001])
```
（唯一差别＝`set=[…]` 那枚 map 的**打印顺序**，⛔ 内容差别。）

成对时长（同一枚用例，A ⇄ B）：
```
A: --- PASS: TestGateJudgesThePrivateSetByResolvedSID (0.01s)   B: --- FAIL: … (5.45s)
A:     --- PASS: …/runner_descriptor_is_judged_by_what_it_resolves_to (0.00s)   B:     --- FAIL: … (1.55s)
A:     --- PASS: …/a_name_spelling_cannot_change_the_verdict (0.00s)            B:     --- FAIL: … (1.93s)
A: --- PASS: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (0.08s)          B: --- FAIL: … (1.57s)
```
那 1.5–2.0 秒＝`testing.go` 里 `removeAll` 对 `Access is denied` 的**重试窗**（`const arbitraryTimeout = 2 * time.Second`，GOROOT 逐字见 `30-cleanup-failure.md` §2）〔读码推的：时长与窗长的对应是本腿推的；窗与重试本身是 GOROOT 源码逐字〕。

## 3. 问②的答复：那枚 SID 是 `S-1-5-21-1228170099-895614386-1166154857-500`

- **它是什么**：本机 `BUILTIN` 的**本地 Administrator 账号**（RID 500），SDDL 里渲染成 `LA`。逐字名册尺＝`logs/05-gate-start.txt`：
  `LOCALACCOUNT name=Administrator sid=S-1-5-21-1228170099-895614386-1166154857-500 rid=500 disabled=True`；
  本 token 自己＝`token_user_sid=S-1-5-21-1228170099-895614386-1166154857-1001`（`DESKTOP-LVS7839\swq`）⇒ 这台机上 `LA` ⛔ 是干活的人（＝票面 §"现量"第 2 条那句"这一枚在⛔ 集合内"的期望形状）。
- **"解析出来却⛔ 在集合内"逐字凭据**：上面 §2 那行的 `LA=S-1-5-21-…-500, in set: false`，集合＝`{S-1-5-18, S-1-5-32-544, S-1-5-21-…-1001}`。
- **它同时就是那四行失败的锁门人**：`\001\la`、`\002\la` 两枚残留叶的 DACL 逐字＝`D:PAI(A;;FA;;;LA)`（唯一 trustee 就是 `…-500`）；`\001\runner`、`\001\named`、`\001\numeric` 三枚叶的 DACL 逐字＝`D:PAI(A;;FA;;;SY)(A;OICIIO;GA;;;SY)(A;;FA;;;BA)(A;OICIIO;GA;;;BA)(A;;FA;;;LA)(A;OICIIO;GA;;;LA)`（`logs/16`／`logs/18`）——本 token⛔ 在其中一枚都没有（`S-1-5-32-544` 在 token 里是 `"Group used for deny only"`）。
  ⚠ 第四枚残留叶 `\003\la` 的锁门人⛔ 是它而是 `S-1-5-32-546`（`D:PAI(A;;FA;;;BG)`，BUILTIN\Guests，本 token⛔ 持），它⛔ 出现在那四行里（同一枚 tempDir 只报第一枚错，见 `30-cleanup-failure.md` §3）。
- ⇒ 一句话：**判据没换答案，换答案的是"谁能删掉那枚种子目录"**；`…-500` 是那枚"解析出来⛔ 在集合内、于是被正确拒绝、也正确地把自己锁在门里"的 SID。

## 4. 一条与 `AC#1` 裁形直接相关的读数（⛔ 由本腿裁，只落数）

同一批两发里，"封之后清了谁"的**通告**（`noticeNarrowed` 那条 WARN）枚数随卷变：

```
尺＝grep -c 'cleared_inherited=' ：根 A 那发 = 17（其中有内容 12）⇄ 根 B 那发 = 58（其中有内容 53）
尺＝grep -c 'TempDir RemoveAll cleanup' ：根 A = 0 ⇄ 根 B = 4
```
同一枚夹具（`TestAC1DefaultLogSaysInherited` 的 `\001\store`）两侧逐字：
```
B(39): … path=D:\tmp\…\TestAC1DefaultLogSaysInherited1148699846\001\store kind=inherited cleared="" cleared_inherited=S-1-5-11(A;ID;0x1301bf;;;AU),S-1-5-11(A;OICIIOID;SDGXGWGR;;;AU),S-1-5-32-545(A;ID;0x1200a9;;;BU),S-1-5-32-545(A;OICIIOID;GXGR;;;BU)
A(39): … path=C:\Users\swq\AppData\Local\…\TestAC1DefaultLogSaysInherited1372533800\001\store\log-leg.txt kind=inherited cleared="" cleared_inherited=S-1-1-0(A;ID;0x1200a9;;;WD)
```
⇒ 换卷之后 winsec 的**封**真的多清了 `S-1-5-11`／`S-1-5-32-545` 两组继承 ACE（这正是 §1 那张"差的 trustee"表在码里的自证），但**⛔ 改变任何一枚通过／拒绝的结论**。
票 311 `AC#1` 甲形那句"让判据⛔ 依赖继承"在本票现量里**够不着那两枚红**——红⛔ 吃判据，红吃的是收尾。（这一句是读数＋成因的分工陈述，裁形⛔ 归编排者。）
