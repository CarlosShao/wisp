# 361-a1 普查腿（只读）part 1／2：前提新鲜度＋35 枚逐名归类

- 腿：`361-a1`（只读腿，只做 AC#0 现读＋归因；不落 `.go`、不改票面、不 push）
- 时刻：2026-10-22 16:4x +08
- HEAD：`21f05dd4`；本机 `core.autocrlf=true`；gofumpt `v0.12.0 (go1.27.1)`
- 票面：`.scratch/wisp/issues/361-attrib-tracked-only-is-red-locally-but-green-in-ci-the-third-cross-shape-case.md`
- NOTE 读数原件在仓外 `D:\tmp\wisp\20261022\361a1\`（`local-tracked-only.txt` / `archive-tracked-only.txt` / `archive-with-poscontrol.txt` / `selftest.txt` / `classified.txt` / `d22scan.txt` / `pathlen.txt`）
- 续件：`docs/evidence/s1/361-a1-census-part2.md`（sec 4 起）

**结论先行：票面的两个前提都不成立，已触发停手条件 (1) 与 (2)，逐字读数见 part 2 的 sec 9。**
票面说「31 枚台件（归属规则允许红）＋ 4 枚 CRLF」；实测的**原因**分裂是 **20 枚入库 blob 本身就脏（真红）＋ 15 枚 CRLF 假红**。

## sec 1 前提新鲜度（逐条现跑，对 HEAD blob 复量）

(1) `ci.yml:240` 逐字（`git show HEAD:.github/workflows/ci.yml`，第 240 行，行首 8 空格）：

```
        run: sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only
```

紧邻三行（238/239/241）＝空行、`        if: ${{ !cancelled() }}`、空行。**该步在盘、逐字与票面一致。**

(2) `attrib.sh` 在盘：`.scratch/wisp/probes/161/r5/attrib.sh`，662 行、36828 字节，可执行。

(3) 同一 job（`lint`）上方 `ci.yml:172` 那步是 `gofumpt -l . tools/d22scan tools/mockllm`（非 tracked-only）。
`ci.yml:191-193` 自己写着「On a checkout every file is tracked, so the step above and this one read the same bytes and must agree」——这是本次判决的关键伏笔。

(4) 分母：`git ls-files '*.go'` = **1155**；其中 `.scratch/**` = **360**、非 `.scratch` = **795**。

## sec 2 35 枚逐名＋逐枚归类

`sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` 本机 rc=**1**，`lines=35 files=35`。
归类判据（逐枚现跑）：`blob_CR` ＝ `git cat-file -p HEAD:<f> | tr -dc '\r' | wc -c`，`wt_CR` 同法读工作树。

- **B ＝ 入库 blob 本身就脏**（`wt_CR=0` 且 `blob_CR=0`，纯 LF 仍然脏）＝> 真红，CI 同形
- **C ＝ CRLF 假红**（`blob_CR=0` 但 `wt_CR>0`）＝> 读入库 blob 时干净

| # | 路径（仓库相对） | 台件? | 归类 |
|---|---|---|---|
| 1 | `.scratch/wisp/probes/163/a1/main.go` | 是 | B |
| 2 | `.scratch/wisp/probes/174/c2/zz174c2_wiring_pair_windows_test.go` | 是 | B |
| 3 | `.scratch/wisp/probes/183/r2/mut/task-boxset-off.go` | 是 | B |
| 4 | `.scratch/wisp/probes/185/r1/mut-m1/hostpath_185.go` | 是 | B |
| 5 | `.scratch/wisp/probes/197/r1c/pre/subagent_197.go` | 是 | B |
| 6 | `.scratch/wisp/probes/197/r1c/pre/subagent_197_test.go` | 是 | B |
| 7 | `.scratch/wisp/probes/212/v1/mut/main-noq9.go` | 是 | B |
| 8 | `.scratch/wisp/probes/220/r1/prechange/prechange220_probe_test.go` | 是 | B |
| 9 | `.scratch/wisp/probes/222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go` | 是 | B |
| 10 | `.scratch/wisp/probes/224/v2/dialect_probe_test.go` | 是 | B |
| 11 | `.scratch/wisp/probes/231/v1/mutation/config_reload_branch_after_prefix.go` | 是 | B |
| 12 | `.scratch/wisp/probes/235/v1/mut/prereadoff-m3_222_test.go` | 是 | B |
| 13 | `.scratch/wisp/probes/241/v1/posctl/badly_formatted.go` | 是 | B |
| 14 | `.scratch/wisp/probes/241/v1/probe_v1_readings_test.go` | 是 | B |
| 15 | `.scratch/wisp/probes/253/v2/denoms253v2.go` | 是 | B |
| 16 | `.scratch/wisp/probes/259/r1/gofumpt-negctl/probe.go` | 是 | B |
| 17 | `.scratch/wisp/probes/263/v1/src/main.go` | 是 | B |
| 18 | `.scratch/wisp/probes/33/p1/q1/main.go` | 是 | B |
| 19 | `.scratch/wisp/probes/33/p1/q2/main.go` | 是 | B |
| 20 | `.scratch/wisp/probes/33/p1/q3/main.go` | 是 | B |
| 21 | `.scratch/wisp/probes/257/r2/mut/M1/firstrun.go` | 是 | C (wt_CR=185) |
| 22 | `.scratch/wisp/probes/257/r2/mut/M2/firstrun.go` | 是 | C (wt_CR=176) |
| 23 | `.scratch/wisp/probes/257/r2/mut/M3/firstrun.go` | 是 | C (wt_CR=174) |
| 24 | `.scratch/wisp/probes/257/r2/mut/M4/firstrun.go` | 是 | C (wt_CR=188) |
| 25 | `.scratch/wisp/probes/257/r2/mut/M5/firstrun.go` | 是 | C (wt_CR=185) |
| 26 | `.scratch/wisp/probes/257/r2/mut/M6/firstrun.go` | 是 | C (wt_CR=185) |
| 27 | `.scratch/wisp/probes/257/r2/mut/M7/firstrun.go` | 是 | C (wt_CR=185) |
| 28 | `.scratch/wisp/probes/257/r2/mut/M8/firstrun_257_nonpreset_test.go` | 是 | C (wt_CR=176) |
| 29 | `.scratch/wisp/probes/287/r1/mutants/stale-crlf/M1-disable-branch.go` | 是 | C (wt_CR=612) |
| 30 | `.scratch/wisp/probes/287/r1/mutants/stale-crlf/M2-identity-to-corr.go` | 是 | C (wt_CR=612) |
| 31 | `.scratch/wisp/probes/287/r1/mutants/stale-crlf/M3-overbroad-any-row.go` | 是 | C (wt_CR=612) |
| 32 | `cmd/wisp/models.go` | 否 | C (wt_CR=389) |
| 33 | `internal/agent/approval/pending_read.go` | 否 | C (wt_CR=131) |
| 34 | `internal/agent/tools.go` | 否 | C (wt_CR=225) |
| 35 | `internal/risk/provenance.go` | 否 | C (wt_CR=1115) |

**票面的 31/4 是路径分裂（台件 31、非台件 4），这一步是对的；但把它当成原因分裂是错的。**
真正的原因分裂是 **B=20（全是台件）＋ C=15（台件 11 ＋ 非台件 4）**。
票面把 11 枚 CRLF 假红的台件混进了「归属规则允许」那一堆。

## sec 3 `attrib.sh` 的归属规则：逐字、哪一行、怎么判

归属规则由**两个函数**构成，只在 `classify_line()` 里合流：

`attrib.sh:185-193`（`ticket_of`，判路径形状）：

```sh
ticket_of() {
    case $1 in
    .scratch/wisp/probes/[0-9]*/*)
        t=${1#.scratch/wisp/probes/}
        printf '%s' "${t%%/*}"
        ;;
    *) printf '' ;;
    esac
}
```

`attrib.sh:196-203`（`ticket_known`，判票是否真实存在）：

```sh
ticket_known() {
    [ -n "$1" ] || { printf '0'; return; }
    if ls "$root"/.scratch/wisp/issues/"$1"-*.md >/dev/null 2>&1; then
        printf '1'
    else
        printf '0'
    fi
}
```

规则本体写在 `attrib.sh:60-63`（文件头自述）：

```
# A line is attributed iff its repo-relative path is
#   .scratch/wisp/probes/<NN>/<...>               AND
#   .scratch/wisp/issues/<NN>-*.md exists        (the ticket is real)
```

**怎么判**：`classify_line()` 先调 `classify()`（`attrib.sh:172-182`，问 git 这路径是否 tracked）。
若 tracked 则直接判 `TRACKED<CI-RED>` 并 `CL_RC=1`（`attrib.sh:225-229`）；只有**非 tracked** 时才走 `ticket_known` 那半条归属规则（`attrib.sh:231-235`）。
**归属规则只对工作树（B 尺）的非 tracked 行生效。**

**关键：这条规则在 `--tracked-only` 模式下根本不执行。**
`--tracked-only` 走 `attrib.sh:627-640`，只调 `ruler_a()` 和 `summary()`；`ruler_a()`（`attrib.sh:539-583`）**从不调用 `classify_line`**：它把 gofumpt 吐的每一行无条件记成 `TRACKED-DIRTY`、`TRACKED_DIRTY` 加一，然后 `attrib.sh:633` `[ "$TRACKED_DIRTY" -gt 0 ] && RC=1`。
`attrib.sh:52` 自己把这条写死了：`rc=1  (A) printed anything at all  -> tracked bytes are dirty; CI red`。

＝> **票面「`attrib.sh` 自己的归属规则允许 `.scratch` 台件红」这句话，对 `--tracked-only` 是假的。**
该模式下 `.scratch` 台件与 `internal/` 产码**同权**：只要 tracked 且不 gofumpt-clean，一律红、一律点名。

另记一条副作用：`attrib.sh:631` 在 tracked-only 模式下仍打印 `attribution rule = ...`，而该模式 (B) 根本没跑——这行字会误导读者以为归属规则参与了判决。
