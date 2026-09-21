# 73 — Sweep orphan `.wisp-tmp-*` staging files left by a real kill

**Status:** ready-for-review — AC#1–AC#5 全部已证并勾（证据 `docs/evidence/s1/73-orphan-staging-sweep.md`）；本票收尾，未 push
**Type:** defect-fix (bookkeeping of a proven leak)
**Blocks:** nothing · **Blocked by:** nothing (the characterization test already exists)
**Spec refs:** D31 atomic-rename, SPEC-07 §2–§3, registry **A18**
**Packages:** `internal/tools` only (do not touch `internal/ball`, `internal/observe`, `tools/d22scan` — other agents are in flight there)

## Proven premise (do not re-derive, do not re-litigate)
`761447f` added `internal/tools/bridge_a18_kill_windows_test.go`: a **real child process** + **real
`taskkill /F`** mid-write. It proves three things, all already measured, not guessed:
1. The destination file is byte-for-byte intact or absent ⇒ **D31 holds under a real fault**.
2. **Each kill leaves exactly one `.wisp-tmp-*` behind** (4 bytes each in the test).
3. **Nothing sweeps it** — a later successful write to the same directory does not remove the orphan.

That test deliberately pins the assertion to "the orphan is still there". **This ticket flips that
one assertion to "the orphan is gone"** and must not weaken any other assertion in the file.

## What to build
Recovery of interrupted atomic writes at the place that already owns the target path, so a killed
session does not accumulate litter in the user's allowed dirs:
- On the next write (or on bridge start-up, your choice — say which in the commit message and why),
  reclaim `.wisp-tmp-*` files in the directory being written **that this bridge itself created**.
- The reclaimer must be **conservative**: never delete a `.wisp-tmp-*` it cannot attribute to itself
  (attribution mechanism is your design — record the naming scheme or a marker; explain the choice),
  and never follow a junction/reparse point or a symlink while sweeping (that would turn a cleanup
  into a delete primitive outside the allowed dirs — this repo's C26 exists precisely for that).
- A temp file that is **currently held by another live process** must survive; a retry-then-skip is
  acceptable, a hard failure is not.

## AC (1:1 verdict table required, one row per box)
- [x] **AC#1** The flipped A18 assertion: after a real `taskkill /F` mid-write, a subsequent write
  through the bridge leaves the target directory with **zero** `.wisp-tmp-*` entries; the destination
  file is still byte-exact. Mutation check required: revert the reclaimer to no-op ⇒ this test goes
  red (grep-prove the mutation landed before running, restore and grep-prove it is gone).
- [x] **AC#2** Attribution is real, not "delete every `.wisp-tmp-*` in the dir": a `*.wisp-tmp-*`
  file **the bridge did not create** (e.g. a foreign file with the same prefix, created directly by
  the test) survives the sweep, with a test that names the surviving file.
- [x] **AC#3** Reparse-point safety: an orphan-shaped name reached through a **real NTFS junction**
  pointing outside the allowed dirs is not deleted. `mklink /J` for real, `t.Fatalf` (not `t.Skip`)
  if the fixture cannot be built, and a positive control proving the target file exists first.
- [x] **AC#4** A live temp file held by another process is not deleted and does not fail the write:
  spawn a child that opens the temp and blocks, then write. `//go:build windows`, zero `t.Skip`.
- [x] **AC#5** Gates: `gofmt -l internal/tools` empty, `go vet ./internal/tools/...` rc=0,
  `go test ./internal/tools/... -count=2`, `-race`, and the repo's `d22scan` invoked **exactly as
  the CI line invokes it** (`cd tools/d22scan && go run . -root <abs>`), with a seeded-violation
  self-test shown red at least once so a green scan means something.

## Explicit prohibitions for this ticket
- No `git commit --amend`, `git reset`, `git rebase`, `git stash`, `git checkout .` (shared worktree,
  registry A34). Commit with `git commit -F - -- <explicit paths>` and a **quoted** heredoc
  (`<<'MSGEOF'`) — an unquoted one made the shell execute backticks on `9bf05cb` (A31).
- After `git add`, before `git commit`, run `git diff --cached --name-only` and confirm only your
  paths are staged.
- Do not tick an AC you did not reproduce. Do not add `t.Skip`, do not lower a threshold, do not
  touch frozen contracts (`docs/PLAN.md`, `docs/specs/*`, `internal/risk/assessor.go`,
  `internal/risk/pathresolver*.go`, `rules_gateway.go`).
- If an AC cannot be proven, **leave the box empty and say why** in the Progress log.
- Commit your first checkpoint **within your first 15 tool calls** (turn cap is ~150; a ticket that
  dies with nothing committed is a total loss). Sync the ticket face (Status + boxes + one Progress
  log line ending in `next=`) at **every** commit.

## Progress log (append-only, newest last)

- 2026-09-21（实现代理，checkpoint 1）：开工前只做了设计拍板，代码未落，AC 框一个都没勾。
  **清扫时机 = 下一次写盘**（`stageAndRename`/`crossVolume` 在 `os.CreateTemp` 之前对目标目录扫一次），
  不是桥启动时——理由：本仓 `internal/tools` 里没有"桥启动"这个可挂钩的单点（`New()` 只是构造，
  `cmd/wisp` 不许我碰），而证据 §5 已经写明"下一次启动"在这张票的可达面上只能用"新建一套桥 + 一次成功写盘"代表，
  正好就是 A18 用例第三条在做的事。**归属机制 = 暂存文件名自带身份**：
  `.wisp-tmp-<owner8>-<pid>-<rand>`，owner8 = sha256(可执行文件路径 + 用户名 + tempPrefix) 前 8 位十六进制
  （跨进程稳定 ⇒ 新进程能认死进程的文件；不同程序/不同用户 ⇒ 不同 token，前缀相同也不会互删）；
  pid 用来判"创建者还活着吗"（活人持有的暂存文件不是孤儿）。扫到名字不合式、是 reparse point、
  不是普通文件、创建者 pid 还活着、或删不动（重试后仍失败）的一律**跳过不删**。
  **next=落 `internal/tools/fs_staging.go` + 平台 liveness 文件 + 翻 A18 第三条断言，然后按 AC#1 做变异检验**。

- 2026-09-21（实现代理，checkpoint 2 · **AC#1 已勾**）：清扫器落地 `internal/tools/fs_staging.go`
  （命名/归属/扫描）+ `staging_live_windows.go`（`OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` 判活
  + `FILE_ATTRIBUTE_REPARSE_POINT` 判 reparse）+ `staging_live_other.go`（非 Windows 一律"判不了 ⇒ 不删"）。
  挂点在 `fs_write.go:282`（fs.write 的 `stageAndRename`）与 `:505`（fs.move 跨卷），都在
  `os.CreateTemp` **之前**，所以清扫永远看不到本次写盘自己的暂存文件。A18 第三条断言已翻成
  "下一次写盘扫干净 + 目标 `existing.txt` 逐字未动"，用例改名 `TestA18…TheNextWriteReclaimsTheStagingFile`
  （常量 `a18ChildTest` 同批改，子进程 re-exec 指的就是它）。
  **②的阳性对照从 2 个改成 1 个**：现在第二个子进程自己写盘前就把第一个的孤儿扫掉了 ⇒ 每次 kill 仍留 1 个、
  但可见残留只剩最后一个。这是行为变化不是放宽：0 个会让"扫掉了"变成空跑。
  **AC#1 变异检验**：`fs_staging.go` 里插 `return got // MUTATION-AC1` ⇒ grep 命中 210 行 ⇒
  `go test -run TestA18` **FAIL**（"残留 = [两个 .wisp-tmp-*]，want 1 个"）；删掉该行 ⇒
  `grep -rn MUTATION-AC1 internal/tools/` rc=1（零命中）⇒ 复跑 `ok 1.046s`。
  **AC#2/#3/#4/#5 未勾**：判据用例还没写；`-count=2`/`-race`/d22scan 门禁还没跑。
  **next=写 `internal/tools/fs_staging_windows_test.go`：外人不删（AC#2）、真 junction 不穿（AC#3）、
  活人持有的不删且写盘照常成功（AC#4），然后跑全套门禁 + d22scan 种子阳性。**

- 2026-09-21（实现代理，checkpoint 3 · **AC#2/#3/#4 已勾**）：新增 `internal/tools/fs_staging_windows_test.go`
  （`//go:build windows`，零 `t.Skip`）。AC#2 `TestSweepReclaimsOnlyItsOwnStagingFiles`：6 个撞前缀的外人文件
  （逐个点名，含"形状全对但 owner 不是本程序"和"pid 段不是数字"）+ 1 个"本程序活 pid"的孤儿全部存活，
  同目录那个死 pid 孤儿被扫掉；另有 `TestSweepNamingSchemeRejectsNearMisses` 钉字符串层
  （发现并修掉一个真洞：`pid == "0"` 曾被当成合法创建者）。AC#3 `TestSweepNeverDeletesThroughARealJunction`：
  `mklink /J` 真 junction，链名**就是可归因的孤儿名**，两道阳性对照（读得到对面字节 + 同目录普通孤儿确实被扫掉），
  然后断言对面文件字节未动**且 junction 本身没被删**。AC#4 `TestSweepSparesATempFileHeldByAnotherProcess`：
  子进程 `CreateFile` 共享位**故意不给 FILE_SHARE_DELETE** 后阻塞；A 形状（名字里是活 pid）与 B 形状
  （名字里是死 pid、但句柄在活人手里）都存活且写盘成功；B 再真 taskkill 释放后下一次写盘必须扫掉它
  ——证明上一条的幸存是句柄而不是名字。子进程 re-exec 自己（`-test.run=^TestSweepSpares…$`），三处 `t.Fatalf`
  写清缺什么前提。
  **变异检验（超出票面要求，因为 AC#2/#3 是这张票的危险面）**：把 `sweepStagingOrphans` 退回"前缀即我的"
  （grep 命中 `fs_staging.go:224` 的 `MUTATION-73`）⇒ AC#2/AC#3/AC#4 **三条全红**（AC#3 报"清扫器连 junction
  本身都删了"）；删掉该行 ⇒ `grep -rn MUTATION-73 internal/tools/` rc=1 零命中 ⇒ `TestSweep|TestA18` 全绿。
  门禁部分完成：`gofmt -l internal/tools` 空、`go vet ./internal/tools/...` rc=0、
  `go test -count=2 ./internal/tools/...` `ok 27.617s`、`-race` `ok 18.041s`。
  **next=按 CI 原样跑 `scripts/d22scan.sh`（含 d22scan 自身 seeded 阳性）后勾 AC#5，并写证据文件。**

- 2026-09-21（实现代理，checkpoint 4 · **AC#5 已勾，本票五框全绿**）：门禁按票面逐条量过——
  `gofmt -l internal/tools` 空、`go vet ./internal/tools/...` rc=0、
  `go test ./internal/tools/... -count=2` `ok 27.617s`、`-race -count=2` `ok 36.283s`；
  d22scan **按 CI 原样**跑：`sh scripts/d22scan.sh`（= 先 `cd tools/d22scan && go test ./...` 的
  seeded 阳性对照 `ok`，再 `go run . -root "D:/work/workspace/projects plans/Wisp"` → `clean`，
  ban #7 scope `internal/tools/` examined 16）与票面写的手动式都跑过。
  **green 之前先证红**：临时在 `internal/tools/` 播一个 `filepath.Clean(filepath.Abs(p))` 种子文件
  ⇒ 扫描 **rc=1** 并打两条 `[pathresolver-bypass]` + `2 finding(s)`；删掉种子文件 ⇒ 复扫 rc=0 clean，
  `ls` 确认树上没留。证据 `docs/evidence/s1/73-orphan-staging-sweep.md`（§1 结果表、§2 两次变异、
  §4 门禁数字、§5 四条没证明的事）。
  登记两条不属于本票、我没碰的观察：①第一次跑 `sh scripts/d22scan.sh` 时 `tools/d22scan/main.go:818`
  正被别的代理改到编译不过（`undefined: io`），脚本把失败吞成了 `script-rc=0`（我读的是 `tail` 的 rc，
  但同一形态在 CI 里用 `${PIPESTATUS}` 之外也会溜）；②本票改名后的 A18 用例函数名与
  `docs/evidence/s1/20-a18-taskkill-residue-characterization.md` 里引用的旧名不再一致（历史证据不改写）。
  **next=owner/复核代理按 AC#1–AC#5 逐条复现；registry A18 判据②与 Q-16 的收口话术归 orchestrator，不在我职权内。**
