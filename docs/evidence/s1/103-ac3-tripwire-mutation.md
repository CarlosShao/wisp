# 票 103 AC#3③ 复验 —— "tripwire 改下降"的变异让责任用例红在下降路径上；还原位同；两包 `-count=2` 绿

**腿：** 复验腿（票 103 派单；本腿不实现、不改断言）
**交件时的树：** HEAD `071541fb5`（派单写 `1aa96f8ac`；`1aa96f8ac` 是 `071541fb5` 的父提交，已用 `git merge-base --is-ancestor 1aa96f8ac HEAD` 量到 `ANC_YES`；两枚之间 `git diff --stat 1aa96f8ac HEAD -- internal/memory internal/winsec` 为**空** ⇒ 本腿读数对 `1aa96f8ac` 的写面同样成立）
**写面核对（开工前实测）：** `git status --porcelain internal/memory internal/winsec` 空。
**写面：** 只写 `docs/evidence/s1/103-ac3-tripwire-mutation.md` 一份；快照与日志全部在仓外/`.scratch`（未提交）。
**环境：** `export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP`；`go version go1.27.1 windows/amd64`；全程**未用** `CGO_ENABLED=0`。
**快照：** `C:/Users/swq/AppData/Local/Temp/wisp103-ac3-2008`（`git archive HEAD | tar -x` 出来的**仓外**副本，非 git 仓库；仓库目录内**未建 worktree、未 checkout**）。

---

## 0. 一句话结论

**PASS。** 在 HEAD 的仓外快照里做 AC#3③ 的唯一一发变异——把 `internal/memory` 回收遍历从
`filepath.WalkDir`（Lstat 基、不下降）换成 `os.Stat`+`os.ReadDir` 的递归（**跟链接＝下降**）——
`go build` rc=0 先量到，责任用例 `TestAC2ReclaimWalkNeverDescendsIntoAJunction` **FAIL**，红行点名
`removeStray(junklink)` 与它构造出的 `junklink\sub\keep-me.txt`（即"遍历确实下降了"）。
同一发里票 79 的旧用例 `TestStrayRemovalDoesNotFollowLinks` 仍 **PASS**（新 tripwire 是真增量）。
还原后 `internal/memory` + `internal/winsec` 共 **119 个文件 sha256 全同、`diff -rq` 零差异**，
同一快照复跑该用例 **PASS**。真实树上四道门（build / 两包 `-count=2` / `gofumpt -l` 空 / d22）全过。

**没有发现生产 bug，不触发 STOP**：MUT-3 下 `winsec.RemoveUnlinked` 的祖先链接守卫仍在位，
外来文件**没有被删**（leg 3 未红），所以这是"下降被响亮拒掉"，不是"别人树里的文件静默消失"。

---

## 1. 树、写面、环境（可核）

| 项 | 命令 | 读数 |
|---|---|---|
| 当前 HEAD | `git rev-parse HEAD` | `071541fb5dc385acf39aa3ed83426072a2a8ab3d` |
| 派单 HEAD 是祖先 | `git merge-base --is-ancestor 1aa96f8ac HEAD` | `ANC_YES` |
| 两枚之间写面无漂移 | `git diff --stat 1aa96f8ac HEAD -- internal/memory internal/winsec` | 空 |
| 开工前写面干净 | `git status --porcelain internal/memory internal/winsec` | 空 |
| Go | `go version` | `go1.27.1 windows/amd64` |

---

## 2. 变异（AC#3③ 唯一一发：tripwire 的"不下降"改成"下降"）

**落点（同链 `grep -n` 证的落地行）：** `internal/memory/artifacts.go:188`

```
142:		walkErr := filepath.WalkDir(filepath.Join(dir, e.Name()),        ← 测量走法，未动
188:	err := descendFollowingLinks103(filepath.Join(s.artifactsDir, name), ← 回收走法，被换
```

`diff -u`（对 `git show HEAD:internal/memory/artifacts.go`，**只有这一行**）：

```diff
@@ -185,7 +185,7 @@
 // symlink standing where a child was listed gets unlinked, not followed.
 func (s *Store) removeStray(name string) error {
 	var dirs []string // artifacts-relative, walk (pre-)order; reversed = children first
-	err := filepath.WalkDir(filepath.Join(s.artifactsDir, name),
+	err := descendFollowingLinks103(filepath.Join(s.artifactsDir, name),
 		func(p string, d fs.DirEntry, err error) error {
 			if err != nil {
 				return err
```

新增的变异体（**只存在于仓外快照**，不进仓库）`internal/memory/mut3_descend_windows.go` 全文：

```go
//go:build windows

package memory

import (
	"io/fs"
	"os"
	"path/filepath"
)

// descendFollowingLinks103 is ticket 103 AC#3(3) MUT-3 shape, kept out of the
// repository: a hand-rolled pre-order walk that stats THROUGH links (os.Stat,
// not Lstat), so a junction reports as a directory and is descended into. It
// exists only to prove that TestAC2ReclaimWalkNeverDescendsIntoAJunction bites
// when the reclaim walk starts following links.
func descendFollowingLinks103(root string, fn fs.WalkDirFunc) error {
	return walk103(root, fn)
}

func walk103(p string, fn fs.WalkDirFunc) error {
	fi, err := os.Stat(p) // the mutation: follows the link instead of lstat-ing it
	if err != nil {
		return fn(p, nil, err)
	}
	d := fs.FileInfoToDirEntry(fi)
	if err := fn(p, d, nil); err != nil {
		return err
	}
	if !fi.IsDir() {
		return nil
	}
	entries, err := os.ReadDir(p)
	if err != nil {
		return fn(p, d, err)
	}
	for _, e := range entries {
		if err := walk103(filepath.Join(p, e.Name()), fn); err != nil {
			return err
		}
	}
	return nil
}
```

**编译先量到（编译失败不算变异）：** `go build ./internal/memory/` ⇒ `BUILD_RC=0`。

> 为什么不用 `filepath.EvalSymlinks` 版：票面 137 行已登记"Windows 上 EvalSymlinks 对 junction 不解引用，第一发没落地成下降"；
> 本腿直接用 `os.Stat`（跟链接）+`os.ReadDir`，与 `103-adversarial-acceptance.md` §MUT-3 的形状一致，可复现。

---

## 3. RED 读数（责任用例红在"下降"上）

`go test ./internal/memory/ -run "TestAC2ReclaimWalkNeverDescendsIntoAJunction|TestStrayRemovalDoesNotFollowLinks" -count=1 -timeout 420s -v` ⇒ `TEST_RC=1`

```
=== RUN   TestAC2ReclaimWalkNeverDescendsIntoAJunction
    artifacts_junction_tripwire_windows_test.go:40: mklink /J: Junction created for ...\data\artifacts\junklink <<===>> ...\data\someone-elses-tree
    artifacts_junction_tripwire_windows_test.go:60: walk over the junction entry reported [...\data\artifacts\junklink]
    artifacts_junction_tripwire_windows_test.go:68: removeStray(junklink): memory: reclaim artifacts stray "junklink":
      memory: remove artifacts stray file "junklink\\sub\\keep-me.txt":
      winsec: entry is a link to something else, not private data ...\data\artifacts\junklink\sub\keep-me.txt:
      the spelling reaches it through the link at ...\data\artifacts\junklink,
      and whatever lives behind that link is not this tree's data to delete
--- FAIL: TestAC2ReclaimWalkNeverDescendsIntoAJunction (0.24s)
=== RUN   TestStrayRemovalDoesNotFollowLinks
--- PASS: TestStrayRemovalDoesNotFollowLinks (0.17s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/memory	0.471s
TEST_RC=1
```

读法（先说得准确，免得被读成"tripwire 每条腿都咬住了"）：

1. **责任用例确实红了，且红行点名了被削弱的行为**：`removeStray(junklink)` 报出的 failed path 是
   `junklink\sub\keep-me.txt`——这是**只有遍历下降了才会构造出来的路径**；红行同时点名了 `removeStray`（回收遍历的入口）。
2. **红在 leg 2（`artifacts_junction_tripwire_windows_test.go:68`），不在 leg 1（:60）**：leg 1 是对
   **标准库 `filepath.WalkDir` 的直接调用**，它当然还只报一个条目（:60 的读数就是证据）。
   票面 136 行与 `103-adversarial-acceptance.md` §MUT-3（197-203 行）**已各自登记过这一点**：
   leg 1 测的是 stdlib 行为，MUT-3 咬住的是**结果面那条腿**。本腿复现读数与它们一致，不是新缺口。
3. **不是"红在噪声"**：同一次运行只有这一条 FAIL；leg 3（外来 `keep-me.txt` 与 `someone-elses-tree` 仍在）
   与 leg 4（阳性对照：普通 `stray-dir` 照样被清）都**没有**报错 ⇒ 守卫挡住了删除，红是"下降被拒"而不是"删了别人的文件"。

---

## 4. 还原位同（`diff` 干净 + 119 文件 sha256 全同 + 同快照复绿）

还原动作（快照内）：`rm internal/memory/mut3_descend_windows.go`；
`git show HEAD:internal/memory/artifacts.go > internal/memory/artifacts.go`。

| 判据 | 命令 | 读数 |
|---|---|---|
| 被改文件逐字节 | `diff -q <(git show HEAD:internal/memory/artifacts.go) $SNAP/internal/memory/artifacts.go` | 无输出（`ARTIFACTS_CLEAN`） |
| 整个写面 vs 新 `git archive HEAD` | `diff -rq $SNAP/internal/memory $VER/internal/memory`、`diff -rq $SNAP/internal/winsec $VER/internal/winsec` | 无输出（`SURFACE_CLEAN`） |
| 逐文件哈希 | 两份树各 `find ... -type f \| sort \| xargs sha256sum` 后 `diff` | 无输出（`SHA256_IDENTICAL`），**119** 个文件 |
| 同快照复跑 | `go test ./internal/memory/ -run TestAC2ReclaimWalkNeverDescendsIntoAJunction -count=1 -v` | `--- PASS (0.40s)`，`ok memory 0.484s`，`RESTORED_RC=0` |

⇒ 红是变异造成的（同一快照红→还原→绿），快照里没有残留文件。

---

## 5. 门禁（真实树，逐条现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | `BUILD_RC=0` |
| 回归（任务指定的两包） | `go test ./internal/winsec ./internal/memory -count=2 -timeout 420s` | `TEST_RC=0`；`ok winsec 72.518s` / `ok memory 84.235s`；`--- FAIL` **0** |
| 格式 | `/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/ internal/memory/` | **空**（`GOFUMPT_RC=0`，`GOFUMPT_EMPTY`） |
| D22 | `sh scripts/d22scan.sh` | `D22_RC=0`，`clean - no D22 ban violations`；scope 计数与复核腿一致（bans #1-5 internal/=363、cmd/=58，ban #6 frontend/=116，ban #7 internal/tools/=48，ban #8 design/=39、frontend/=116、internal/=904、cmd/=211，不降） |

门禁跑完复查：`git status --porcelain internal/memory internal/winsec` 仍空 ⇒ 本腿没有在生产写面留下任何改动。

---

## 6. 本腿做了什么、没做什么

**做了：** AC#3③ 的唯一一发变异（§2）、`go build` rc=0 先量（§2）、责任用例的 RED 读数与对照组（§3）、
还原与逐文件哈希证明（§4）、四道门（§5）。

**没做 / 没碰：**
- 未新增/修改任何生产码、任何既有用例、任何断言、任何 SKIP 条件（写面 `git status` 空）。
- 未动 `internal/risk/**`、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、golden、`docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`。
- 未勾 AC 框、未改票面、未改名 `-done`、**未 push**；未在仓库目录内建 worktree / checkout。
- 只复验 AC#3③；AC#3①/②不在本腿射程（读数见 `103-head-reappraisal.md` §5 与 `103-adversarial-acceptance.md` §MUT-1/2）。

**读法边界（写给裁决者）：**
- 若把"tripwire 咬住变异"理解为"leg 1 那行断言必须红"，则**它不会红**，因为 leg 1 调的是 stdlib `filepath.WalkDir`，
  不经过生产码；这在票面 136 行已自述。真正咬住变异的是 leg 2/3（结果面），本腿的 RED 读数就在 leg 2。
  这**不构成 STOP**：生产守卫（`removeStray` 用 `filepath.WalkDir`）在 HEAD 上仍在位，变异只存在于快照；
  且变异下外来文件没有被删（leg 3 绿），不是生产 bug。
- 本腿**不建议**为了"让 leg 1 也红"去改生产 API 或加新断言：leg 2（`removeStray(junklink)` 必须 nil）已经把
  "遍历下降"变成可观测红，再加一层只会与现有断言重复；若裁决者坚持要 leg 1 直接测生产走法，应另开票，不在本腿私自改被测件。

---

## 7. 复跑最短路径（复制即用）

```bash
cd "/d/work/workspace/projects plans/Wisp"
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP

# 0) 基线（应 PASS）
go test ./internal/memory/ -run "TestAC2ReclaimWalkNeverDescendsIntoAJunction" -count=1 -timeout 420s -v

# 1) 仓外快照 + MUT-3（把 WalkDir 换成跟链接的递归）
SNAP="/c/Users/swq/AppData/Local/Temp/wisp103-ac3-replay-$$"; mkdir -p "$SNAP"
git archive HEAD | tar -x -C "$SNAP"
cat > "$SNAP/internal/memory/mut3_descend_windows.go" <<'EOF'
//go:build windows

package memory

import (
	"io/fs"
	"os"
	"path/filepath"
)

func descendFollowingLinks103(root string, fn fs.WalkDirFunc) error { return walk103(root, fn) }

func walk103(p string, fn fs.WalkDirFunc) error {
	fi, err := os.Stat(p) // follows the link: the mutation
	if err != nil {
		return fn(p, nil, err)
	}
	d := fs.FileInfoToDirEntry(fi)
	if err := fn(p, d, nil); err != nil {
		return err
	}
	if !fi.IsDir() {
		return nil
	}
	entries, err := os.ReadDir(p)
	if err != nil {
		return fn(p, d, err)
	}
	for _, e := range entries {
		if err := walk103(filepath.Join(p, e.Name()), fn); err != nil {
			return err
		}
	}
	return nil
}
EOF
sed -i 's|err := filepath.WalkDir(filepath.Join(s.artifactsDir, name),|err := descendFollowingLinks103(filepath.Join(s.artifactsDir, name),|' \
  "$SNAP/internal/memory/artifacts.go"
grep -n "descendFollowingLinks103\|filepath.WalkDir" "$SNAP/internal/memory/artifacts.go"
( cd "$SNAP" && go build ./internal/memory/ && echo BUILD_RC=$? )   # 必须先 rc=0

# 2) RED（应 FAIL，红行含 junklink\sub\keep-me.txt）
( cd "$SNAP" && go test ./internal/memory/ \
    -run "TestAC2ReclaimWalkNeverDescendsIntoAJunction|TestStrayRemovalDoesNotFollowLinks" \
    -count=1 -timeout 420s -v )

# 3) 还原 + 位同证明
git show HEAD:internal/memory/artifacts.go > "$SNAP/internal/memory/artifacts.go"
rm -f "$SNAP/internal/memory/mut3_descend_windows.go"
diff -q <(git show HEAD:internal/memory/artifacts.go) "$SNAP/internal/memory/artifacts.go" && echo CLEAN

# 4) 真实树门禁
go build ./cmd/... ./internal/... ./tools/... ; echo BUILD_RC=$?
go test ./internal/winsec ./internal/memory -count=2 -timeout 420s ; echo TEST_RC=$?
/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/ internal/memory/ ; echo GOFUMPT_RC=$?
sh scripts/d22scan.sh ; echo D22_RC=$?
```
