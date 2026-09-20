# 69 — C21 token 表的 61 行新内容没有任何机器检查（A24-D4）

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-20
**Blocked by:** 68（它正在改 `internal/ball/`，同包并行=假并行）⇒ **本票在票 68 收尾前不得开工**
**Parallel slots:** ≤1 sub-agent（`internal/ball/` + `docs/evidence/s1/c21-native-tokens.md`）
**Spec refs:** C21、SPEC-08 §2（**冻结，不得编辑**）、D29
**登记项:** A24-D4

## 为什么要这张票

票 12 AC#7 的代理把 `docs/evidence/s1/c21-native-tokens.md` 从 **133 行补到 227 行**
（新增 **25 条几何/动效常量 + 36 个 look 色**），并把 78 条配色行与
`design/assets/tokens.css`、`internal/ball/tokens.go` **逐值对账到零漂移**。这是真交付。

**但它同时暴露了一个空档**：今天只有 **20 条**配色受
`TestTokenGoldenValues` 保护，`TestNoHardcodedColorsInBallPackage` 管的是"不许在绘制码里写字面量"，
**不是**"表与码一致"。⇒ **新加的 61 行没有任何机器检查**：
改了 `tokens.go` 不改表、或改了表不改码，都不会被任何东西抓到。
表头现在也写明了"36 个 look 色在 `tokens.css` 里无出处、只有人工表这一道保证，
**不得当契约级事实引用**"——这句话今天是**真的**，本票的任务就是让它变成不需要这句话的东西。

## 验收标准

- [ ] **AC#1 反向检查落地**：一条机器断言，遍历本票负责的两类 token
  （几何/动效常量、look 色），证明**表里每一行都能在码里找到同名同值的对应物，且反之亦然**。
  实现方式自选（testdata 里的表 + 解析，或从 Go 侧生成再 diff），但**判据必须是双向**：
  只查"码里的都在表里"是不够的（那正是漂移的方向）。
- [ ] **AC#2 变异检验**：故意把 `tokens.go` 里一个几何常量改一个数（或删一行 look 色），
  新用例**必须转红**；然后还原并证明全套件回到绿。⚠ **先 grep 证明变异真的落地**再跑测试，
  跑完立刻 `git checkout --` 还原；原始输出留在 `docs/evidence/s1/69-mutation-*.md`。
- [ ] **AC#3 不许弱化**：`TestTokenGoldenValues` 现有 20 条断言**一条都不许删或放宽**；
  只能是超集。若某条 look 色因票 65（玻璃质感返工）注定要变，**登记成显式豁免行**而不是删断言。
- [ ] **AC#4 门禁**：`gofmt -l` 触及包为空、`go vet ./internal/ball/`、
  `go test -count=2 ./internal/ball/`（**非 winlive**；若本票不需要桌面就别碰桌面）。贴原始输出。

## 边界（不要越界）

- **不碰** `docs/specs/*`（SPEC-08 冻结）。若发现表与契约冲突，**报出来**由我裁定，不改文本。
- **不在本票改球面默认值**（那是票 68 AC#2 的活）；两票若都动 `internal/ball` 必须**串行**。
- 若票 65 的质感返工尚未开始，本票**只锁现有值**，不预判新值。

## Progress log（append-only；每个 commit 一行 `- [ISO-UTC] agent=... did=...`）
