# 票 266 CI 步腿 `266-ci` —— 形 ①-A 的另一半接线（`.github/workflows/ci.yml` 加一步）

**腿名**：`266-ci`（CI-step 腿，写面＝`.github/workflows/ci.yml` ＋ `docs/evidence`）
**派单钉的 HEAD**：`0fa1527f9`；**本腿实际开工 HEAD**：`0af3b17aa`（＝钉的下一笔，编排者的收件 commit
"收wave17:230-AC4 seam1 6a46140ea+266 PASS;登记overrides.go:78 C26红;A1553"，只动 `docs/reports/pending-and-issues.md` 6 行）
⇒ **漂移 +1 笔**，本腿在 `0af3b17aa` 上做，未 reset / 未 rebase。
**交件 commit**：`6e5957756`（ci.yml，14 行纯插入，父 `c6d1d187e`）＋ `23cfd5725`（本文件）＋ 本条补正（sha 由
`git log --oneline -- docs/evidence/s1/ticket266-266-ci-step-v1.md` 现取；本文件不写自己未知的 sha）。
**窗口内 HEAD 漂移（现读，非推断）**：`0af3b17aa` → **`c6d1d187e`**（兄弟腿 "284 rename 收尾:删旧名"，`2026-11-05 18:45:44`）
→ `6e5957756`（本腿 ci.yml，`18:48:15`）。`0af3b17aa` 与 `6e5957756^` 的 `.github/workflows/ci.yml` **blob 逐字节相同**（`git rev-parse` 两枚 sha 相等），
且 `git log --oneline 0af3b17aa..HEAD -- .github/workflows/ci.yml` **只有本腿一笔** ⇒ 漂移没夹带别人对 ci.yml 的改动，本腿的 14 行仍是唯一增量。
**来源**：`docs/evidence/s1/ticket266-266-r2-impl-v1.md` §6（逐字补丁）——本腿是"无脑照抄"，**未增删一字**（§2 有逐行字节比对读数）。

---

## §1 落点：前/后逐字片段

### BEFORE（`0af3b17aa` 的 `.github/workflows/ci.yml:203-211`）

```yaml
203	        # forbidden action. The step runs the seeded-path positive control first (in
204	        # a throwaway repository, so the control cannot be red because of this tree)
205	        # and refuses to run the scan at all if that control cannot be made red.
206	        if: ${{ !cancelled() }}
207	        run: sh scripts/check-path-length-budget.sh --with-self-test
208	
209	      # CALL SURFACE CHANGED (ticket 362, AC#1 form (A'), owner-approved edit of
210	      # .github/**; the approval is docs/reports/pending-and-issues.md A1161 (1),
211	      # and the dispatch that carries it is commit 76ae3bbb). Four claims:
```

### AFTER（`6e5957756` 的同一处：206-224）

```yaml
206	        if: ${{ !cancelled() }}
207	        run: sh scripts/check-path-length-budget.sh --with-self-test
208	
209	      - name: "SLO nail census (ticket 266 - the outer tooth over slo-check.ps1's nails)"
210	        # ADDED STEP ONLY. This is a text reader with no toolchain: it declares in
211	        # scripts/slo-nail-census.sh that the D32 gate carries exactly two
212	        # self-locking nails, and it goes red when one of those call lines is
213	        # deleted or commented out (ticket 263 mutant delete-nails measured an
214	        # all-green screen with both calls removed; ticket 266 is that hole).
215	        # It refuses to print a verdict at all - exit 2 - when the target is
216	        # missing/empty or its sampling anchor is gone, so "did not look" cannot
217	        # read as green. Runs its own 32-shot positive control first and refuses
218	        # to print a reading if the control cannot be made red. No threshold, no
219	        # golden, no allowlist entry, and no existing step touched.
220	        if: ${{ !cancelled() }}
221	        run: sh scripts/slo-nail-census.sh --with-self-test
222	
223	      # CALL SURFACE CHANGED (ticket 362, AC#1 form (A'), owner-approved edit of
224	      # .github/**; the approval is docs/reports/pending-and-issues.md A1161 (1),
```

**形**：ADDED STEP ONLY。`git diff --numstat -- .github/workflows/ci.yml` = **`14  0`**（纯插入、零删除、零修改行）。
`on:` / `paths:` / `concurrency:` / `runs-on:` / 既有 14 步：一字未动。

---

## §2 校验读数（本腿现跑，带取数时刻）

| 校验 | 命令 | 读数 | 时刻 |
|---|---|---|---|
| YAML 解析 | `python -c 'import yaml;yaml.safe_load(open(".github/workflows/ci.yml"))'` | **rc=0**；`jobs` 仍是 `['lint','test-core','test-windows','slo-smoke','slo-full','lint-frontend']`（6 个没多没少） | 18:45:5x |
| 步数与位置 | 同上解析 | `lint.steps` **14 → 15**；新步 `index=6`，`keys=['if','name','run']`（**无** `continue-on-error`）；前一步 `'Tracked path-length budget (ticket 262)'`，后一步 `'gofmt (gofumpt)'` | 18:45:5x |
| 自动步名唯一 | 同上 | 除 `[None]`（checkout / setup-go 两步本无名，既有）外无重名 | 18:45:5x |
| 与 §6 逐字一致 | 抽证据 §6 围栏内 13 行 vs ci.yml:209-221 逐行 | **byte-exact = True**（13/13 行完全相同，无尾随空白、无 tab） | 18:46:1x |
| 行尾一致 | 全文件字节统计 | 1376 行 **全 LF、CRLF=0**（`0af3b17aa` 的 blob：1362 行、CRLF=0）⇒ 未混入 CRLF | 18:46:4x |
| 空白 | `git diff --check` | 干净（无空白错误） | 18:46:4x |
| 语法（脚本） | `bash -n scripts/slo-nail-census.sh` / `scripts/d22scan.sh` / `scripts/check-path-length-budget.sh` | **rc=0 / rc=0 / rc=0**（本腿**未改任何脚本**，这三发只是证"步里点的那把尺本身语法成立"） | 18:42:0x |

### §2A 步里那条命令本机现跑（＝该步在 CI 上会做的事）

```
$ sh scripts/slo-nail-census.sh --with-self-test
slo-nail-census: positive control PASSED - 32 control(s) ran, 0 failed
slo-nail-census: sampling anchor: first line carrying [-FilePath $WispExe] is scripts/slo-check.ps1:533
slo-nail-census: PASS - 2/2 nails declared here are defined once, carry their own ok sentence,
                 are called once at column 0, and run before the first sample at scripts/slo-check.ps1:533
RC=0
```

- **接步前**（18:40:3x）：`control ... ok` × 32、`0 failed`、`PASS - 2/2`、**rc=0**
- **交件后**（18:48:30）：同上，`control ... ok` 计数 **32**、**rc=0**
- 失败数 **0/32**；`cannot-decide` 未被触发（被检件在、非空、采样锚 `:533` 在）。

### §2B D22 扫描（`sh scripts/d22scan.sh`）

三条读数必须一起看，否则会把别人的红算到这一步头上：

| 树 | 命令 | 读数 | 时刻 |
|---|---|---|---|
| **在飞工作树**（含兄弟腿 230-ac4 未提交的 `internal/config/overrides.go`） | `sh scripts/d22scan.sh` | **rc=1**，1 条 finding：`internal/config/overrides.go:57: [phantom-citation] comment cites repo path "docs/evidence/s1/230-ac4-config-rewrite-account.md" which does not exist on disk` | 18:41:5x |
| 同一在飞树（兄弟腿把那个文件写出来之后） | 同上 | **rc=0**，`clean - no D22 ban violations`（scopes：`internal/`363、`cmd/`58、`frontend/`116、`internal/tools/`48、`design/`39、ban #8 `internal/`904、`cmd/`211、`frontend/`116） | 18:49:4x |
| **`0af3b17aa` 的纯 HEAD 内容**（`git archive HEAD` 解到 `%LOCALAPPDATA%\wisp-tmp\266ci-head-export`，**不碰工作树**） | 该目录内 `sh scripts/d22scan.sh` | **rc=0**，`clean`（无 git ⇒ ignore 规则不生效，分母略不同：`design/`14、`internal/`902） | 18:47:0x |
| **同一 HEAD 内容 ＋ 本腿的新 ci.yml** | 同上 | **rc=0**，`clean`；`grep -c slo-nail-census .github/workflows/ci.yml` = **2**（name 行＋run 行） | 18:48:0x |

**天归（这条最重要）**：18:41 那条红**不是本步造成的**——①D22 的 scope 表只有 `internal/ cmd/ frontend/ internal/tools/ design/` 五个，**`.github/**` 根本不在扫描面内**；
②那条 finding 落在 `internal/config/overrides.go`，本腿一个字节都没碰（`git diff --stat -- internal/config/overrides.go` 的改动属兄弟腿 230-ac4 seam1）；
③把本腿的新 ci.yml 叠在**纯 HEAD 内容**上单跑 ⇒ rc=0。18:49 兄弟腿补齐被引文件后，在飞树回到 rc=0。

---

## §3 写面与交件干净度

- 本腿**只**写：`.github/workflows/ci.yml`（commit `6e5957756`）、本文件。
- `git status --porcelain -- .github docs/evidence`（commit `6e5957756` 后）：`.github` **空**；`docs/evidence` 只有别的腿的未跟踪件
  （`195-r1-direction-verdict.md`、`230-ac4-config-rewrite-account.md`、`351-ac2-ac3-verdict.md`、`378-r1-impl.md`），**无本腿遗留的 M**。
- 未碰（逐条现读核过）：`.github/workflows/` 下的其它文件、任何 `scripts/**`、`docs/PLAN.md`、`docs/specs/**`、
  `internal/observe/thresholds.go`、任何 golden、`tools/d22scan/allowlist.txt`、`internal/config/**`、`cmd/wisp/**`、`frontend/**`、`design/**`。
- **无 push**（推送只由编排者做）；`git add`/`commit` 均带显式 pathspec，带 `index.lock` 重试循环。

---

## §4 这个步的射程（诚实档：它买到什么、买不到什么）

1. **买到的**：从 `6e5957756` 起，"有人把 `scripts/slo-check.ps1` 里那两枚钉的调用删掉/注释掉"这件事，**每次 push/PR 的 lint 作业都会红**，
   而且红的名字是具名的（`slo-nail-census-nail-call-missing` / `...-ordering` / `...-declaration-drift` / `...-cannot-decide`）。
   此前这把尺**零 CI 调用点**，全靠人手跑（票 266-r2 证据 §6 "未接线的后果"）。
2. **买不到的**：它只钉**形状**（定义 1 次、自带 ok 句、列 0 未注释地调用 1 次、位置在首枚采样 `:533` 之前）。
   "那两行 ok 句**真的被打印过**"仍无人断言——那是票 266 形 ②（`slo-report.json` 记账字段 ＋ `slo-smoke` 尾部新步）。
3. **没动门禁口径**：无阈值、无 golden、无 allowlist、无 `continue-on-error`、无 `if: false`；只加了一步。

**撤销口令（一句话复原）**：`git revert 6e5957756` —— 只回退这一个步的 14 行插入，`.github/workflows/ci.yml` 回到 `0af3b17aa` 的形状。
