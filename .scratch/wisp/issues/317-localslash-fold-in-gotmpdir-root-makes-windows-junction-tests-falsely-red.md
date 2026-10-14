# 317 — 310/311 的 `GOTMPDIR` 转换把 `\` 主动折成 `/`：windows 腿上 junction/路径字面量族用例**假红**（`mklink … Invalid switch - "Users"`；`want` 混合分隔符）

- **Status**：**ready-for-agent**（编排者 2026-10-14 2026-10-14 13:18 立；起手锚点＝本票现读 `git rev-parse --short HEAD`）。
- 来路＝**本次推送 CI（run `38508311376`）vs 上一发（`38484004521`）的差集**＋编排者现读定位：新增红共 5 枚，其中 4 枚＝`TestSyncAncestorActableLegFailsClosedOnMovedAncestor116`／`TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116`／`TestSyncAncestorLegReResolvesThroughJunctionToRealTree`／`TestSyncAncestorLegNotExemptedStaysFailClosed`（`internal/risk`），**本机同树全绿**（`go test -run … ./internal/risk/` 0.8s 全 PASS）⇒ 环境形。
- **因果（我在 CI 日志里读到逐字）**：
  - `precondition: deepestExistingAncestor("C:\Users\runneradmin\…\%WISP116ANC\B%\new.md") = "C:\Users\…\B%" want "C:/Users/runneradmin/…\…\B%"` ⇒ **want 是正斜杠拼法**；
  - `mklink /J C:\Users\… TestSyncAncestorLegReResolves…\001 failed: exit status 1: Invalid switch - "Users".` ⇒ **正斜杠目标被 cmd 当开关**。
- **元凶行（现读）**：`scripts/portable-tests.sh:144` `gotmp_base=${WISP_TEST_GOTMP_ROOT:-"${LOCALAPPDATA//\//}/wisp-test-gotmp"}`——`${LOCALAPPDATA//\//}` 把 `\` 折成 `/`；`winsec-tests.sh:63`／`wisp-cli-tests.sh:76` 同形（另三处 `wisp_311_ref_path=${LOCALAPPDATA//\//}` 是 311 的对照位，⛔ 未必同一问题，落地腿逐处判）。
- **出身核（编排者现跑）**：该转换**不在旧 tip**（`git show 2c764c04:scripts/portable-tests.sh | grep LOCALAPPDATA` 空）；引入笔＝`6a8f537a`（票 310 AC#1"send the test temp family to a dedicated GOTMPDIR root"）＋`f8b1f446`（票 311）——**都在未推送存量里，本发 CI 首次暴露**。旧发同族用例跑在默认 `…\Temp\`（反斜杠形）下不红。
- **交叉先例**：`236-v2` 2026-10-14 本机实测同一效应（posix 形 GOTMPDIR ⇒ 5 枚 junction 用例红）；本机 `GOTMPDIR` 一律用反斜杠形是运行纪律（已入铁律）。

## 要建什么

- [ ] **AC#1 复现（先读 310／311 的原意，再动手）**：先读票 310／311 与 `6a8f537a`／`f8b1f446` 的 rationale——**为什么当时要折成正斜杠**（若为 linux 腿形状/若是"只许 C: 卷"的顺手写法/若 tests 依赖它 ⇒ 逐条具名）；再在本机用 `WISP_TEST_GOTMP_ROOT='C:/tmp/xxx'`（posix 形）复现 4 枚红（236-v2 形状），反斜杠形复现绿。
- [ ] **AC#2 修（形状由落地腿现读后裁，须同时保住 310／311 的原意）**：候选 ⓐ windows 腿用反斜杠形（只对非 windows/linux 腿保留折法）；ⓑ 全腿都不折（核 310 原意是否允许）；ⓒ 只修"当 GOTMPDIR 值"的那两处、`wisp_311_ref_path` 那些对照位逐处判。**判据**：① 修后以 posix 形 `WISP_TEST_GOTMP_ROOT` 注入时，4 枚**仍须能触发同类红**（说明尺没被拆——等等，反了：修后 harness 自己产出的 GOTMPDIR 应是反斜杠形 ⇒ 4 枚在 harness 下**绿**）；② 本地能以 harness 的 windows 腿跑通（若本机可跑；不可跑则给出下一发 CI 的预期读数与判法）；③ `bash -n` rc=0；⛔ 不许为绿放开任何测试/删 skip。
- ⚠ **若修必须动 `.github/**` 或触发表 ⇒ 停手上报**（契约邻接）；写面＝`scripts/**`（＋本票证据件）。

## 禁区

⛔ 三枚冻结件／`thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`specs/**`；⛔ `frontend/**`／`design/**`；⛔ 不动 174-r5 已入册的 ledger 行（除非本修实证其表述需改——那就停手上报）；⛔ push。

## Progress log (append-only, newest last)

- [2026-10-14 13:18 +08] 编排者立票（来路＝两发 CI 差集＋现读定位；元凶行三处同形；出身笔 `6a8f537a`／`f8b1f446` 在存量里）。
