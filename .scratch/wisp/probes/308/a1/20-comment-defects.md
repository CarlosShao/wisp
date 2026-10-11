# 308-a1 · 20-comment-defects（Ⓐ：独立复跑 `A843` §2 那两枚注释缺陷的尺）

尺一律现量；件 `internal/risk/pathresolver_anchor_spelling_windows_test.go`（崩点文件，本票射程内）。

## Ⓐ① 指向一枚⛔ 存在的证据件（指针抄错）

- 崩点文件里的句子（现量行号 **`:157`**，非 `A843` 记的 `:158`——同句，行号随改动漂移）：

      // exactly as it was. See probes/72/r1/20-canonical-form-choice.md for why the

  尺＝`grep -n 'canonical-form-choice' internal/risk/pathresolver_anchor_spelling_windows_test.go` ⇒ 命中 `:157`。
- `ls .scratch/wisp/probes/72/r1/20-canonical-form-choice.md` ⇒ **No such file**（⛔ 存在）。
- `ls .scratch/wisp/probes/72/r1/` 全名册（真件）＝`00-anchor.md` / `10-class-a-expected-side.md` /
  `20-class-b-negative-nails.md` / `30-evidence-readings-gates-mutations.md` / `logs/` / `run-two-surface.sh`。
  ⇒ **没有 `20-canonical-form-choice.md` 这一枚**；`20-` 前缀真实占位是 `20-class-b-negative-nails.md`（乙组，⛔ 是本句要的选型）。
- `find .scratch/wisp/probes -name '*canonical-form-choice*'` ⇒ 0 命中（全仓⛔ 此件）。
- 那段"为什么单形计数⛔ 再可断言"的选型真身＝在 `10-class-a-expected-side.md` 的 **`## #2 TestCanonicalInputGainsNoSecondForm`** 一节（尺＝`grep -nE '^#' 10-class-a-expected-side.md`，`#2` 在 `:27`）。
- ⇒ 与 `A843` §2 判断一致：**指针抄错（该指 `10-…#2`），⛔ 是编造凭据**（区别要紧：判语照认，但这枚件名必改）。归口＝票 308（同一枚文件、注释写成现在时声称那族病）。

## Ⓐ② "deliberately left byte-for-byte untouched" 是现在时的假话

- 崩点文件 `:12-16`（现量命中行 **`:15`**）逐字：

      // pathresolver_junction_windows_test.go, which is the failing instrument and
      // is deliberately left byte-for-byte untouched: its expected values, its

  尺＝`grep -n 'byte-for-byte untouched' …anchor_spelling_windows_test.go` ⇒ 命中 `:15`。它指的对象＝`pathresolver_junction_windows_test.go`。
- `git log --oneline -- internal/risk/pathresolver_junction_windows_test.go` ⇒ 顶端一枚就是 **`f278c57f 72-r1 class A: move 7 expectation sides off hardcoded temp spellings onto verdicts`**（即"把 junction 文件改了"的那笔）。
- `git show --numstat --format= f278c57f -- internal/risk/pathresolver_junction_windows_test.go` ⇒ **`36  7`**（`+36 −7`，⛔ 逐字节不变）。
  （同笔还改了 `pathresolver_anchor_spelling_windows_test.go +92 −18`，`10-class-a-expected-side.md +142`、`run-two-surface.sh +49`。）
- ⇒ 与 `A843` §2 判断一致：**引入这半句的正是同一票 72 的 `f278c57f`，它当场把那枚"号称逐字节不动"的文件改了 `+36 −7`** ⇒ 注释把已作废的做法写成现在时。归口＝票 308。
- 本腿⛔ 改这两处（同两枚件⛔ 是裁决者射程，规矩：注释缺陷登记、⛔ 顺手改对；`A843` 已具名"⛔ 由我改、归口票 308"）。
