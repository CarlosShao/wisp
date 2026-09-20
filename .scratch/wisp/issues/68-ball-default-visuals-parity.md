# 68 — 让默认构建画的就是 owner 签收的那个球（`prototypeVisuals` 默认值 + Sleeping 尺寸三方不一致）

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-20
**Blocked by:** —（只碰 `internal/ball/` + `cmd/balldebug/`；票 66 在 `internal/proc`/`internal/observe`/`cmd/wisp`，票 67 在 `internal/llm/adaptertest`/`tools/d22scan`，三包不相交）
**Parallel slots:** ≤1 sub-agent；**AC#2/AC#3 需真桌面** ⇒ 若桌面被占，**只做 AC#1 与 AC#4 的静态半，其余保持未勾并写明**
**Spec refs:** SPEC-08 §2 / §2.1（**冻结，本票不得编辑**；含 2026-09-20 的 INTERIM 标记）、D29、D32（不放宽）
**登记项:** 票 12 AC#7 代理报出的争议 **D1 / D2**（registry 待并入 A24）

## 事实（我自己读码确认，非转述）

`internal/ball/statevisual.go:76-84`：

```go
// prototypeVisuals selects the ticket 62 liquid-glass rendering. It is OFF by
// default: SPEC-08 §2.1 (12px/0.35 Sleeping dot etc.) is frozen until the
// owner signs the new look, and the ball's own tests assert that table.
var prototypeVisuals bool
func EnablePrototypeVisuals(on bool) { prototypeVisuals = on }
```

唯一把它打开的生产侧调用者：`cmd/balldebug/main.go:104` → `ball.EnablePrototypeVisuals(!*frozen)`。
被它门控的代码面：`dock_windows.go:90/114/176/206/229/245`、`liquid.go:313`、`liquid_windows.go:52/154`（票 62/64 的吸附与液面全在这里）。

⇒ **owner 2026-09-20 用肉眼签收的那个 44px 静态玻璃体，只在 `balldebug` 里存在；库默认画的是旧的 12px 微点。**
而我今天把 SPEC-08 §2 的 `Sleeping` 行改成了"44px 静态玻璃体"——**契约文本现在描述的是一个非默认配置**。
这个不一致是**我造成的**，本票就是还这笔账。

⚠ 注释里那句"until the owner signs the new look"的**前提今天部分成立**（R13：临时通过、先完成核心功能），
所以翻转默认值是**与已改契约一致**的动作，不是擅自扩大签收范围：**owner 批的只有"尺寸 / 可见 / 零定时器"三项，
质感仍是"赝品"**（见 SPEC-08 §2 INTERIM 与票 65）。

## 第二处不一致（D1，比默认值更要紧）

同一态下三个数对不上：SPEC-08 §2 INTERIM 写 **44px**；
`statevisual.go::stateSize` 在默认 56 基准下按 `SleepRestRatio = 0.62` 算出 **34.72px**（下限 `SleepingRestMinPx = 30`）；
而 `docs/evidence/s1/c21-native-tokens.md` 的旧行写 **12px**。
但 `docs/SLO.md` 附录 A.2 实测到的差分化像框是 **46×46、2103 像素变化 ≥8/255** ⇒ 那个跑法里 Sleeping **确实是 44px 级**。
**所以要先把"到底是哪个尺寸在什么条件下成立"查清楚**，再决定改码还是改文档——**不许**为了让三方一致就去编辑 SPEC-08。

## 验收标准

- [ ] **AC#1 尺寸真相**：以代码 + 一次实测（若桌面可得）钉死 `Sleeping` 在
  ①`prototypeVisuals=false`、②`=true` 且未靠边、③`=true` 且已吸附（dock ramp 之后）三种情形下的**实际像素尺寸**，
  逐情形给出 `stateSize` 的输入与输出、以及差分化像框。**完成判据是一张三列对照表**，不是一个结论句。
  与 SPEC-08 §2 的 44px 不符的那一格**如实报不符**，由我裁定。
- [ ] **AC#2 默认值翻转（需桌面复测）**：`prototypeVisuals` 默认改为 `true`，
  使**默认构建 == owner 签收的样子**；`EnablePrototypeVisuals(false)` 与 `balldebug -frozen` 保留为"对照旧冻结规格"的逃生门。
  连带把 `internal/ball` 里**断言旧默认**的测试改到位：⚠ **不得**为了让测试变绿而删除断言或放宽阈值——
  必须**逐条**把断言迁移到"新默认下应有的行为"，并在 log 里列出"这条原来断什么、现在断什么、为什么等价"。
  若某条测试的存在意义就是"钉住旧的冻结规格"，**保留它**并显式 `EnablePrototypeVisuals(false)`，别删。
- [ ] **AC#3 签收面复测（需桌面）**：翻转后跑一次 `cmd/balldebug` 的差分化像，
  证明 `Sleeping` 的 px≥8/255 与成像框**不低于** A.2 已录的 2103 像素 / 46×46，且 `timers=no` 仍成立（D32 零定时器）。
  数字不过就报 FAIL 并附样本，**不许调阈值、不许重测到运气好的那次**。
- [ ] **AC#4 无桌面的半格**：`gofmt -l` 空、`go vet ./internal/ball/ ./cmd/balldebug/`、
  `go test -count=2 ./internal/ball/`、`go test -count=2 ./cmd/balldebug/`（非 winlive 部分）。贴原始输出。

## 编排者已裁定

1. **SPEC-08 与一切 `docs/specs/*` 冻结**：本票**零编辑**契约文件。尺寸若与契约不符，**改代码或报不符**，不许改文本。
2. **D32 的零定时器与 ≤0.5% CPU 不因"只是默认值"而放松**。
3. 若桌面被票 66/67 占着：AC#2/AC#3 **保持未勾**并写明"等桌面"，**先交付 AC#1 + AC#4**，不要为了推进而偷跑测量。

## Progress log（append-only）
