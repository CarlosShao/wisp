# 227-r1 落地腿（写代理）·票 227 AC#6 · part2/2

本文件对：part1（`227-r1-impl-part1.md`）交改前读数／新用例／真轮询证据／反向判据三读数，
本件交门禁／写面自查／交回。

## sec 5 门禁（本腿现跑，逐条 rc）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| gofumpt | `gofumpt.exe -l/-d` 两枚测试件（绝对路径） | 0 | 两模式都空输出（`logs/50-gofumpt.txt`） |
| build | `go build ./...` | 0 | 0 行（`logs/20-build-precommit.txt`） |
| vet | `go vet ./cmd/wisp/` | 0 | 0 行（`logs/21-vet-precommit.txt`） |
| d22scan | `sh scripts/d22scan.sh` | 0 | `d22scan: clean - no D22 ban violations` |
| path-length | `sh scripts/check-path-length-budget.sh` | 0 | `VERDICT GREEN`，17558 tracked、61 在册、0 不在册 |
| 显式选择器 | `-run '^(TestAC6A...|TestAC1Always...)$'` | 0 | 两枚全 PASS（part1 sec 4） |

整包：主检出会中途死，故另立对照树。主检出在 `.../projects plans/Wisp`（路径带空格），整包跑到
`TestAC14GoSideEvalPushReachesThePage` 时测试二进制直接死（无结果行、无 panic 行，
`logs/15-fullpkg-v-final.log` 尾部停在 `MARK299 before Embed`），221/348 名到不了。这与本腿改动
无关，且对照树（同 HEAD `527b7377` 的 `git archive` 导出、无空格路径）跑到 472 个 RUN：HEAD 不含
本腿文件为 `RUN 472 / PASS 322 / FAIL 20 / SKIP 5`；加本腿文件为
`RUN 473 / PASS 324 / FAIL 19 / SKIP 5`。逐名作差（`logs/86-redset-diff.txt`）：

```
=== reached: control 347  after 348 ===
=== names reached in after but NOT in control (this leg adds one) ===
TestAC6AHandEditIsStillSeenUnderTheRealPoll
=== NEW RED: red in after, not red in control ===
  (nothing listed above means zero new reds)
=== FLIPPED GREEN: red in control, not in after === TestAC4FocusReturnToPriorWindowGap33r5
=== shared red count: 19 ===   --- PASS: TestAC6AHandEditIsStillSeenUnderTheRealPoll (3.27s)
```

新增红为 0；19 枚共同红在对照树（不含本腿任何字节）上就红，属存量；唯一翻绿那枚是桌面窗口族的
争用型，本腿未碰其目标文件。主树红名册：`--- FAIL:` 0 枚，但只跑到 221/348 名就死，故主树的
「整包绿」不能当任何结论用（正是派单警告那族）。

## sec 6 写面与越界自查

| 文件 | 动作 | 是否在写面内 |
|---|---|---|
| `cmd/wisp/always_write_no_clobber_227_test.go` | 新建（测试件） | 是 |
| `cmd/wisp/always_write_no_clobber_226_test.go` | 只追加注释（`git diff` 无 `-` 行） | 是（就地追加） |
| `docs/evidence/s1/227-r1-impl-part1.md` 与本件 | 交付件 | 是 |
| `cmd/wisp/config_reload.go` | 变异用、已还原（sha 三处同串） | 未落地 |
| `cmd/wisp/run.go` | 未动 | 无 |

`git status --porcelain -- cmd/wisp` 起手 0 行（无别人的未提交改动）。禁码位自查
（`logs/91-banned-added.sh`）：226 的新增侧、227 整件、交付件全清；正控
（`docs/reports/HANDOVER.md`）当场命中，尺灵。

## sec 7 具名交回编排者

1. 主检出的整包门禁在这台机器上不可判：路径带空格，panel-host 族在
   `TestAC14GoSideEvalPushReachesThePage` 处把二进制打死。本腿用同 HEAD 的导出树取整包读数，
   不拿主树「0 枚 FAIL」当整包绿。
2. 票 227 的框本腿一枚没翻（现跑尺：判据区 `- [ ] **AC#` 为 6、`- [x]` 为 0，与起手相同）。
3. `226-v2` 第 6 节第 8 条点的那枚「226 写回不 Adopt 内存」确定性用例本腿没做：票面 AC#6 只点名
   「真轮询之下手改仍看得见」这一跳，那枚若仍要属另一格。
