# 252-r3（AC#4／AC#5 两格）锚、方法与 runner 形状复刻

- 本腿派单原名 252-r2 撞 10-03 旧腿；编排者 10-14 12:1x 更正为 252-r3（存档 `.scratch/wisp/probes/orch/20261013/txt/252r3-correction.txt`）。证据件全部落 `probes/252/r3/`，r2/ 零触碰（本腿起手误建的 3 件已即场移出 r2/、未追加任何字节）。
- 临时件路径沿用派单原名拼法 `C:\tmp\wisp252r2\`（更名前已建；只建不删；日志与字节数见 01/02 两件）。
- 起手读数：HEAD `e437db14`（2026-10-14 12:00:40 +08 现读；派单发时近 `d5895f7f`）；`git status --porcelain cmd/wisp internal/tools internal/risk` ＝ **0 行**（起手、终态两次现读）。
- 改前锚＝本票自己的修法 `2b1a3071`（fix(tools/paths): ticket 252 AC#2）的父：`2b1a3071^` ＝ **`62ca6c34`**。
- 交付锚＝载具 commit `d3cccfda`（只动 `cmd/wisp/ticket224_assembly_test.go`，+10/−1；交付文件 `sha256sum` ＝ `366da8f5f64488be4975a2b67f7a7e2f2e40bcbb86166c8859958936487180ff`）。

## 改前树（仓外复刻，不是现读工作树）

```
git archive 2b1a3071^ | tar -x -C /c/tmp/wisp252r2/tree-before     # 落 62ca6c34 全树
```

三枚 witness 在该树逐枚在位：`cmd/wisp/ticket224_assembly_test.go:233`（W1）、`run_test.go:330`（W2）、`run_mode101_test.go:389`（W3）。
载具常数逐字＝W1/W2 载具 `newReplyHost(t, 20*time.Second)` 与 `runFixture` 不写 `confirm_timeout_sec`（default 300）、W3 载具 `confirm_timeout_sec = 1` ⇒ 对上 CI 三发 41.72s／301.16s／4.53s 的量级。

## runner 形状的本机复刻（沿用 AC#3 已验的"真短名当临时根"）

- 建 `C:\tmp\wisp252r2\short-tmp-root-252r2`；`dir /x` 现读其 8.3 别名＝`SHORT-~1`。
- 所有跑发 `TMP=TEMP=C:\tmp\wisp252r2\SHORT-~1`（Git Bash 到 Windows 进程的传递用 `cmd //c echo` 现验过）。
  效果逐字（`dir=` 行）：`C:\tmp\wisp252r2\SHORT-~1\Test…\002\logs` ⇒ `t.TempDir()` 的返回串带短名，与 runner 的 `C:\Users\RUNNER~1\…` 同形。
- **仪器坑（本腿实测）**：设了 `GOTMPDIR` 后 go 工具会把它透传给 test 二进制 ⇒ 临时根落回 `GOTMPDIR`（长名），形状复刻失效（`w3-before-longtmp.log` 即此形：pre-fix 也 PASS 5.47s）。⇒ 短名发一律 `unset GOTMPDIR`（go 工具与 test 二进制同用 `TMP` 短名根，仍在 C: 盘）。
- 形状对照（同机、同树、只换 TMP 拼法）：pre-fix 长名 W3 PASS 5.47s／W2 PASS 5.80s ↔ pre-fix 短名三枚全 FAIL。

## 三枚点名行号现验

- `cmd/wisp/ticket224_assembly_test.go:206-217 grantLinesOf`（`02-ac5-level-dump.md` 的靶）＝**真**（`:208` 起函数体、`:217` 收）。
- `cmd/wisp/run_test.go:378`＝`TestComposedGateBlocksAWriteForTwoSeconds` 体内；pre-fix 树里红句正报 `:378`。
- `cmd/wisp/run_mode101_test.go:506`＝`TestTicket101SessionGrantDoesNotCrossRestart` 体内；红句报 `:506`。
