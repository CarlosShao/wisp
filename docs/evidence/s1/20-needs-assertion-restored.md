# 票 20 复验：`Needs` 配对断言被删（MAJOR）— 编排者独立复核 + 已修 + 变异检验

**时间：** 2026-09-20 21:29–21:31（本地）· **HEAD 起点：** `5c6e661` · **执行者：** 编排者（非票 20 实现代理）
**触发：** 只读审计代理对票 20 的预审报出一条 MAJOR——"删掉的判据与它判的代码在同一个 commit"。
按纪律我**不采信代理结论**，先自己核，核完再修。

## 1. 复核：属实，且比代理说的更值得修

代理给的证据是 commit `0986d63` 把 `TestFSRegistrationIsTheL0Pair` 改名成
`TestFSRegistrationIsTheD34Roster` 时删掉了逐条 `Needs == [fs.read]` 断言。我自己看 diff：

```
-func TestFSRegistrationIsTheL0Pair(t *testing.T) {
+func TestFSRegistrationIsTheD34Roster(t *testing.T) {
-		if len(e.Decl.Needs) != 1 || e.Decl.Needs[0] != CapFSRead {
-			t.Errorf("%s needs %v, want exactly [fs.read] (D34 capability column)", n, e.Decl.Needs)
+func TestDeleteEnabledAddsFSDelete(t *testing.T) {
+	if len(del.Needs) != 1 || del.Needs[0] != CapFSWrite {
```

同 commit `--stat`：`internal/tools/fs_write.go | 742 +++++`、`fs_test.go | 75 ++-`。
**判据与被判的代码在同一个 commit 里被改小**，这是本仓最优先攻击面。

**一处需要替代理说句公道话（也是我给后面所有人定的口径）**：这次改名**不是**为了让测试闭眼——
老断言写的是"五个条目全都 needs `[fs.read]`"，而 `0986d63` 新增了 write/trash/move，
它们**本来就该**是 `[fs.write]`，所以老断言**必须**被改。**错的是把它换成"什么都没有"**
（只留下 `fs.delete` 一条），而不是换成一张逐条表。删守卫必须与替代同批；这次替代缺席。

**为什么 `Needs` 是有牙的字段**（我读 `internal/tools/bridge.go:276-296` 确认，不是猜的）：

```go
Capabilities: dedupe(entry.Decl.Needs),      // :279 —— 决策里记录的能力集直接取自 Needs
need := entry.Decl.Needs                      // :284
if miss := declared.missing(need); ...        // :286 C3：未声明即拒
if miss := b.authz.missing(need); ...         // :291 主机授权集对照的也是 Needs
```

⇒ **C3 的授权判定拿 `Needs` 当需求集**。若 `fs.write` 的 `Needs` 被误写成 `[fs.read]`，
一台只授权 `fs.read` 的机器会**放行一个真会写盘的调用**，而且审计行里记的能力集也是错的。
`registry.go:140` 那条 `Capabilities ⊇ Needs` 只挡"需求报大了"，**挡不住报小了**。

**当前无可利用缺陷**：今天六个条目的 `Needs` 全部声明正确（`fs.go:297/310` read、
`fs_write.go:675/690/703/717` write），丢的是**回归检测力**，不是活漏洞 ⇒ 严重度 MAJOR，非 BLOCKER。

## 2. 修法：把逐条表补回来（`internal/tools/fs_test.go`，测试文件，零生产码改动）

```go
wantNeeds := map[string][]Capability{
    "fs.read": {CapFSRead}, "fs.list": {CapFSRead},
    "fs.write": {CapFSWrite}, "fs.trash": {CapFSWrite}, "fs.move": {CapFSWrite},
}
...
if len(e.Decl.Needs) != len(wantN) || e.Decl.Needs[0] != wantN[0] {
    t.Errorf("%s needs %v, want exactly %v (C3 requirement set)", n, e.Decl.Needs, wantN)
}
```

## 3. 变异检验（先证变异落地，再跑测试 —— 这条是我自己的规矩，今天救过一次冤案）

| 步 | 动作 | 实际输出 |
|---|---|---|
| 前置 | `git status --short internal/tools/fs_write.go` | 空（干净，变异不会踩到别人 WIP） |
| 真码 | `go test ./internal/tools/ -run TestFSRegistrationIsTheD34Roster` | `ok ... 0.035s` **PASS** |
| 变异 | `sed -i '675s/CapFSWrite/CapFSRead/' internal/tools/fs_write.go` | — |
| **验变异落地** | `grep -n "Needs:" fs_write.go` | `675: Needs: []Capability{CapFSRead}` ✅ 真换了 |
| 红？ | 同一条测试 | `fs_test.go:216: fs.write needs [fs.read], want exactly [fs.write] (C3 requirement set)` → **FAIL** ✅ |
| 还原 | `git checkout -- internal/tools/fs_write.go` + `sed -n '675p'` | `Needs: []Capability{CapFSWrite}` ✅ |
| 还原后全套件 | `go test -count=2 ./internal/tools/` | **`ok ... 20.600s`**（含 `checkCaps`/`bridge`/`fs_write` 全套，其余测试仍全绿 ⇒ 旧套件确实从没覆盖这个形状） |
| 门禁 | `gofmt -l internal/tools/`、`go vet ./internal/tools/` | 均无输出（clean） |

**为什么最后一条"其余测试仍全绿"是关键证据**：变异只让**新加的**那条断言转红，
说明删掉老断言之后的整个套件对"五个 fs 工具的 C3 需求集"是**完全盲的** ——
这正是 MAJOR 的量化形式。

## 4. 顺带从该审计确认的三条好消息（我逐条自己复跑/复看）

- **票 20 票面那句 "R12 ANSWER: NOBODY YET" 已经过期**：`cd011b8 feat(12)` 之后
  `cmd/wisp/run.go` 真的 `tools.New`（:262）+ `approval.New`（:255）+ 注册全部
  `BuiltinFSEntries`，且 `loop.go` 的 `decideRisk` 已改为"有 `AdmitTask` 时 L1/L2 放行、无 gate 时仍拒"。
  ⇒ 票 20 的能力**今天在生产装配根里可达**，A12/A13 那条"done 但永不生效"的系统性缺陷对票 20 已解除。
  （这条我要在票 20 票面把过期表述改掉，见 A17 的收尾。）
- **first-match 同族（C-3/M-7/A16）在这条路径上是干净的**：`bridge.go:850-869` 的 `pathArgs`
  遍历**全部**声明键（`fs.move` 是 `{"from","to"}`），`inScope`（:780-789）对所有被judge路径循环；
  值为列表时 `paramString` 返回 `""` → `{Irreversible:["delete"]}` → **L2，fail-closed**。
- 冻结边界未被违反：`c9c3a6f 64c5fea 0986d63 94827e4 e669567` 逐个 `--name-only`，
  **无一命中** `internal/risk`、`docs/PLAN.md`、`docs/specs`、`rules_gateway.go`。

## 5. 审计里我**不采纳**的部分

代理建议把 **AC#3 从 `[x]` 降为 PARTIAL**（理由：`Hooks.Kill` 是**进程内**返回 error，
真 `taskkill` 不跑 Go 的清理，`.wisp-tmp-*` 会留在用户目录）。
它这个技术观察**可能对**，但我没有自己复现 ⇒ **不动框**（纪律：不复现不改勾），
改为登记 **A18**：真·外部 kill 下的暂存文件残留，附完成判据，等票 20 收尾或票 39 一起做。
