# 票 256 `256-r1` 落地写码腿证据件 —— 常驻那条腿的审批门改吃 `[risk]` 两枚配置（形ⓐ：签名不加参）

**腿**＝`256-r1`（只做编排者已裁的 `[risk]` 那一半，`Grants` 那一半归口待立票 265，本轮不碰）
**任务书**＝`.scratch/wisp/issues/256-resident-gate-built-before-session-grants.md` §8（编排者 10-04 09:5x 裁定，账 `A591`）＋只读普查件 `.scratch/wisp/probes/256/a2/census.md`
**写面**（逐字照任务书）＝`cmd/wisp/resident_approval_windows.go` ＋ `cmd/wisp/resident_windows.go` 的 `:126` 一处调用 ＋ 新增测试件 ＋ 本证据件目录
**仓库**＝`D:\work\workspace\projects plans\Wisp`，分支 `dev`
**起手锚点**＝HEAD `d0847aa1`（`264-a1 先写满尾部两节…`）

---

## §1 起手撞钉预检读数

### 1.1 环境读数（本轮新量，具名记着）

（表已就位，读数见 §1.2；这一节记 `cmd/wisp` 测试二进制在本机的加载前提）

### 1.2 今天绿的用例名册（起手一发，逐字抄）

（名册见下，字节非 0）

### 1.3 P1／P3／P4 会不会被本轮顶到

（逐枚判定见下）

---

## §2 落地形状（file:line 逐处）

（三处落点与形状，见下）

---

## §3 四枚常驻判据（① 默认档钉 ② 吃配置的正控 ③ 构造期定值的限制 ④ 落点自证）

（四枚判据的用例名、尺面、终值见下）

---

## §4 突变记录（红句逐字）

（摘掉新逻辑 ⇒ 指名用例必须红；红句逐字进这一节，见下）

---

## §5 `Grants` 那一格：本轮钉没钉、为什么、归口给谁

（结论见下）

---

## §6 门禁读数（终）

（`go build ./...`／`go vet ./cmd/wisp/`／定向 `go test`／`sh scripts/d22scan.sh`／`bash scripts/check-path-length-budget.sh`／`gofumpt -l` 逐条读数见下）

---

## §7 判不动的地方／没做成的（具名，不许用"应该没问题"填空）

（逐条见下）

---

## §8 Git 流水

（本腿每一发 commit 的 pathspec 与 `git show --numstat` 枚数见下）
