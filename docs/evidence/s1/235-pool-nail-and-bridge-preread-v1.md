# 票 235 裁决表 —— 235-v1（非实现者验收腿）

- **验收人**：编排者之外的非实现代理（腿号 `235-v1`）。**本表一行实现码都没写**：零产码、零测试码改动，
  所有突变只经 `go test -overlay` 施加，每发之后 `git status --porcelain -- internal cmd` 的读数逐条记在
  `.scratch/wisp/probes/235/v1/readings.md` §5。
- **被验对象**：票
  `.scratch/wisp/issues/235-the-pool-nail-comment-reason-was-falsified-by-ticket-222-and-only-one-of-three-cases-sees-the-real-bridge.md`
  的四格判据，交件腿＝`235-r1`。
- **起手锚点（本腿现跑 `git rev-parse --short HEAD`）**：`98df640a`。
  被验实现腿的五枚 commit＝`795ed767`／`1d2ad737`／`a71f0be9`／`feb88b91`／`9c96f103`；
  改动面（`git diff --name-only 25556879..HEAD -- internal cmd`）＝
  `internal/tools/subagent_197_test.go`＋`internal/tools/subagent_222_test.go`，**零产码**（本腿现跑对上）。
- **行号口径**：票面／台件里的旧行号一律只当线索。下表"票面行号"＝**本腿在锚点 `98df640a` 现读的该 AC 判据所在行**。
- **读数标记**：〔我现跑〕＝本腿自己执行的命令；〔读台件，未复跑〕＝只看 `235-r1` 的现场件。两者不混。

## 0. 逐格判据（本表正文见 §1-§4；骨架先落，逐格裁完一格追加一节）

| 格 | 票面行号（现读） | 判语 | 凭据（本腿读数所在节） |
|---|---|---|---|
| AC#1 改注释、结论与断言零字符变更；M6 抬池仍是有效检测力证明 | `:19` | 待裁 | §1 |
| AC#2 三枚 222 用例的前置读数；M3 下三枚全红且红不来自 30s 护栏 | `:20` | 待裁 | §2 |
| AC#3 两枚自证（① 摘掉前置读数退回 1/3；② 注释改回旧句是否有仪器保护） | `:21` | 待裁 | §3 |
| AC#4 门禁与终态读数（四数只从 `-v` 量／gofumpt／vet／d22scan／逐名双向差集） | `:22` | 待裁 | §4 |

## 1. AC#1 —— 待裁

## 2. AC#2 —— 待裁

## 3. AC#3 —— 待裁

## 4. AC#4 —— 待裁

## 5. 我攻不动的地方

（逐格裁完后回填，不许为空）

## 6. 没做完／留给编排者

（逐格裁完后回填，不许为空）
