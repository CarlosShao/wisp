# 417 — `tools/d22scan` 的 blob 正控在 CI 上必红：它 `git show` 一个老 commit，而 `actions/checkout` 默认是**浅克隆**

**Status**: open（2026-11-09 23:5x 编排者建，从 CI 真红反推）

## 现量（⛔ 引用前先重跑）

**1. CI 上的真红（`gh api` 拉托管侧日志，`lint` job `129250565029`，run `42643523161`）**：
```
2026-11-09T15:20:14.4992476Z ##[error]    riskwordlist_test.go:118: git show 29d258f9:internal/risk/blacklist.go: exit status 128
2026-11-09T15:20:14.5001117Z --- FAIL: TestRiskWordlistStartAnchorControlIsExactlyAEquiv (0.00s)
```

**2. 本机为什么不红**：
```
$ git merge-base --is-ancestor 29d258f9 origin/dev  →  YES（是祖先）
$ git show 29d258f9:internal/risk/blacklist.go     →  package risk ...（成功）
```
⇒ 本机是**全克隆**，那枚老 commit 在本地。

**3. CI 为什么红**：
```
$ grep -n "fetch-depth" .github/workflows/ci.yml   →  （空：一处都没设）
$ grep -c "actions/checkout" .github/workflows/ci.yml → 7
```
⇒ **7 处 `actions/checkout` 全用默认 `fetch-depth: 1`（浅克隆）** ⇒ `29d258f9` 不在 CI 的本地库 ⇒ `git show` 返 128。

**4. 肇事代码（现读）**：
`tools/d22scan/riskwordlist_test.go:52-57`：
```go
func writeBlobAt(t *testing.T, repoRoot, rev, rel, base string) {
	t.Helper()
	out, err := exec.Command("git", "-C", repoRoot, "show", rev+":"+rel).Output()
	if err != nil {
		t.Fatalf("git show %s:%s: %v", rev, rel, err)   // ← 浅克隆时死在这儿
	}
```
调用点 `:118`：`writeBlobAt(t, root, "29d258f9", "internal/risk/blacklist.go", fixture)`

## ★为什么这条最有价值（铁律 ⓪⓪⓪）

**它不是"CI 偶发红"，是"一条控制尺在 CI 上永远量不到东西"。** `TestRiskWordlistStartAnchorControlIsExactlyAEquiv` 是票 350 的**正控**（证明 phantom-symbol 扫描器能精确命中唯一一枚 `aEquiv`）。它红了 ⇒ **票 350 那道闸在 CI 上等于不存在**。这比"修个 flake"价值高一个档。

## AC

- [ ] **AC#1 判形状（⛔ 直接抄"加 fetch-depth: 0"）**：三种主流解法，腿**先判后改**并写明理由：
  - **ⓐ `fetch-depth: 0`**：最简单，但把整个历史拉进每个 job（7 处 checkout × 全历史 ⇒ CI 时间与磁盘涨）。
  - **ⓑ `fetch-depth: <N>`（够到 `29d258f9` 即可）**：仍需算 N，且**下次再引更老的 commit 又会红** —— 是把定时炸弹往后挪。
  - **ⓒ 测试自足（最主流）**：`writeBlobAt` 不再依赖 git 历史，改成**把那个 blob 的内容作为 fixture 常量内嵌在测试里**（或用 `git hash-object -w` 现场造）。★**这是主流 harness 的做法**：测试不该依赖版本库历史的可得性。
  - ★**我倾向 ⓒ**，理由：①CI 与本地同形（不靠"本地恰好有历史"）②不受 checkout 深度影响 ③省 CI 时间。但**腿自己判**，⛔ 照抄我。
- [ ] **AC#2 反向判据**：造突变「把内嵌 fixture 改一个字符」⇒ 正控必须**变红**（证明它还在量真东西）；还原 ⇒ 变绿。**两个读数都贴**。
- [ ] **AC#3 别处同族**：`git grep -n 'exec.Command("git"' -- '*_test.go'` 找出**所有**依赖 git 历史的测试，逐枚判是否同病（已知一枚：`cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go:46` 只是**注释里**引 commit，⛔ 同病）。
- [ ] **AC#4 门禁**：`gofumpt`（绝对路径）空／`go vet` rc0／`go build` rc0／`go test ./tools/d22scan` 绿／`d22scan` clean／path-length GREEN。
- [ ] **AC#5 CI 侧验证**：推一发后 `gh api repos/CarlosShao/wisp/actions/jobs/<id>/logs` 拉 `lint` job 日志，**`TestRiskWordlistStartAnchorControlIsExactlyAEquiv` 必须是 PASS**，并贴那几行。

## 禁区
- ⛔ 改 `.github/workflows/ci.yml` 除非 AC#1 判为 ⓐ/ⓑ 并写明理由（改 CI 属**契约轴**，需机主批）
- ⛔ 删掉或跳过那条正控（那是**把闸拆了**，不是修）
- ⛔ 动 `tools/d22scan/allowlist.txt` / thresholds / golden
- ⛔ 碰 `frontend/**`、`design/**`

## 撤销口令
- 「**417 撤回**」⇒ 本票作废（无产码改动，随时可撤）。