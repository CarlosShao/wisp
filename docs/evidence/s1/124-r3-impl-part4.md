# 票 124 落地腿 124-r3 交件 part4（sec 9：提交号、逐枚 `.go`、越界与正控）

> 前接 `docs/evidence/s1/124-r3-impl.md`（sec 0..4）、`124-r3-impl-part2.md`（sec 5..6）、
> `124-r3-impl-part3.md`（sec 7..8）。

本腿**两笔**，主题逐字以 `票 124 落地腿 124-r3` 开头；只 commit，NOT push。

| 提交 | 主题 | 路径 | `.go` | 越界 |
|---|---|---|---|---|
| `6b712571` | cmd/wisp 的 37 枚可转用例接上票 119 那条纪律 | 9 | **9** | **0** |
| `<见下一笔>` | 交件（cmd/wisp 批 3 的读数与三态自证） | 4 | **0** | 4（全是本腿证据 `.md`） |

```
$ git show --name-only --format='' 6b712571 | grep -c '\.go$'                    => 9
$ git show --name-only --format='' 6b712571 | grep -v '^cmd/wisp/[^/]*_test\.go$' | grep -c .
                                                                                 => 0
$ git show --name-only --format='' 6b712571 | grep -cE '^(internal/|tools/|scripts/|frontend/|design/|docs/(PLAN|specs|reports)/|go\.mod|go\.sum|allowlist\.txt|\.gitignore|\.github/)'
                                                                                 => 0
$ git show --numstat --format='' 6b712571                                        => 15 增 / 15 删，9 枚逐枚对称
$ git show --name-only --format='' <交件那笔> | grep -c '\.go$'                    => 0
$ git ls-tree -r --name-only HEAD | grep -c 'docs/evidence/s1/124-r3-impl'       => 4
```

**正控（同一把尺打真含 `.go` 的提交，证明尺是活的、NOT 恒 0）**：
`fafe2b4a` => 9/9、`dfa3dc4` => 18/18、`62dda11` => 8/8、`5265c3a` => 3/3、`99f1f9c7` => 15/15
（都是"路径数 = `.go` 数"）。

**负控（同一把尺打纯文档提交）**：`<本件那笔>` => 4 路径 / 0 `.go`、`d055326e` => 3/0、`c50218bf` => 2/0。

**新文件检查**：本腿**零新增 `.go`**（9 枚全是既有文件的等量替换）；证据 4 枚是新件但 NOT `.go`：

```
$ git ls-tree -r --name-only HEAD | grep -c .docs.evidence.s1.124-r3-impl.   => 4
      docs/evidence/s1/124-r3-impl.md
      docs/evidence/s1/124-r3-impl-part2.md
      docs/evidence/s1/124-r3-impl-part3.md
      docs/evidence/s1/124-r3-impl-part4.md
$ （无 cmd/wisp/**124-r3* 之类的新 .go —— 本腿复用批 4 已入树的那一枚 helper）
```

**共树漂移登记**（NOT 本腿的路径）：`design/doubao/**` 25 枚未跟踪件（机主外部 agent 写面，d22scan `ban #8 design/`
由 14 涨到 39 全是它们）、以及 `.scratch/**` 下数千条未跟踪件。逐枚核过这 25 枚在 `19a4f2c9` 里各 0 命中 => 本腿越界仍为 **0**。

**本腿写面（逐条对派单的天花板）**：只 `cmd/wisp/**/*_test.go`（9 枚）＋ `docs/evidence/s1/124-r3-impl*.md`（4 枚）。
`cmd/wisp` **生产码 0 hunk**（`git status --porcelain -- cmd/wisp | grep '\.go$' | grep -vc '_test\.go$'` => 0）；
`internal/{agent,tools,config,winsec,proc,risk}/**`、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、任何 golden、
`go.mod`/`go.sum`/`.gitignore`/`.github/**`、`frontend/**`、`design/**`、`docs/{PLAN.md,specs,reports}/**` **一律零 hunk**。

**临时件（只建不删）**：`D:\tmp\wisp\20261021\124r3\**` —— `snap/`（HEAD 纯净快照 17503 枚）、`post/`、`mut/`（变异台）、
`restored/`（还原台）、`probe3/`（三处"有意不改"的反证台）、`headcheck/`（HEAD 码 + DLL 的真窗族基线台）、
`harness/`（9 支脚本）、`logs/`（8 发 `-v` ＋ 3 发宿主 ＋ 2 发 d22scan ＋ 1 发 pathlen ＋ 2 发非测试容器）、
`rosters/`（逐台名册与 `comm` 差集）。
零 `.out`、未建 worktree/checkout、未用 `-overlay`、未删任何既有临时件。
