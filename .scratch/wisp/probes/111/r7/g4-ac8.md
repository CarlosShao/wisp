# 111-r7 / g4 — 票 111 AC#8（包匹配前缀锚定 + 阳性自证）逐格现读

**本格判（第 0 步现读）**：**半做**——两处匹配式的锚定/转义/时长尾巴**已在盘上**（§1 逐字钉），
缺的是票面逐字要的那发"**种一个只在字符串里出现的包名 ⇒ 不许计入分母**"的**当前 HEAD 阳性自证**
（旧的两发只在票面 Progress log 的散句里，200-111 的非实现者件点名"台件/CI 里没有"）⇒ 本腿补 §2；
⛔ 未改 `scripts/**` 一字节（两处匹配式与所有测试文案零改动）。

## §1 现读（逐字钉"已做的那半"）

- 站点 1 `scripts/portable-tests.sh`（GUARD B，ticket 111 AC#2+AC#8）：
  `:874` `escape_re() { printf '%s' "$1" | sed 's/[][\\^$.*+?(){}|]/\\&/g'; }`；
  `:884` `result_tail='[[:space:]]+([0-9]+\.[0-9]+s|\[build failed\])$'`；
  `:890` `line=$(grep -E "^(ok|FAIL)[[:space:]]+$(escape_re "$pkg")${result_tail}" "$capture" | tail -1 || true)`。
- 站点 2 `scripts/winsec-tests.sh`（guard 2，票面 R-110-3 的病就地）：`:234` `esc_pkgpath=...`（同一条 sed 转义串）；
  `:235` `result_line=$(grep -E "^(ok|FAIL)[[:space:]]+${esc_pkgpath}[[:space:]]+([0-9]+\.[0-9]+s|\[build failed\])$" ...)`。
- 病与三枚钉的注释在 `scripts/winsec-tests.sh:216-233` 逐字（"the original form of the audit was a loose
  substring search for "winsec" over the log, which counted 18 hits in a run where the package was tested 0 times"）。
- CI 侧：这两个站点**每一 run 都在跑**（windows 腿 portable 步／winsec 步的 GUARD B 表与 guard 2 行都在日志里，
  例：run `38676514264` 的 `g4` 同族读数见 g1/g3 件）——但**它们只证明"分母算得对"，不证明"匹配式有牙"**；
  有牙这件事由 §2 的自证交付。

## §2 阳性自证（本腿现取，脚本可复跑）

跑法：`bash .scratch/wisp/probes/111/r7/logs/g4-ac8-proof.sh <scratch>`（脚本内从**提交字节**抽取两条匹配式，
`go` 进程换成仓内既有 shim `scripts/testdata/portable-tests/go`——与 `portable-tests-selftest.sh` 同一约定；
⛔ 不改 `scripts/**`、⛔ 不跑真测试二进制）。逐字输出 `logs/g4-ac8-proof-output.txt`。

**B 面（合成 capture，匹配式本体）**——capture 里 `PKG=.../internal/perm` 只以两种"字符串"出现
（一行 golden fixture 文案、一行把点写成 `X` 的文档串），另有一枚**真** `ok` 行（`.../internal/plugin`，对照）：

| 尺 | 读数 | 判 |
|---|---|---|
| 宽松 `grep -c '<PKG>'` | **2** | 旧仪器（R-110-3 形状）会把"词在日志里"读成"它跑过" |
| 宽松含未转义点（`. = 通配`） | **2** | 连 `githubXcom/...` 那枚也算 ⇒ 转义是承重的 |
| 提交的 GUARD B 匹配式（:874+:884+:890 逐字） | **0** | 字符串提及**不**给分母贡献 |
| 提交的匹配式对真 `ok` 行（对照） | **命中**（打回 `ok .../internal/plugin 0.012s`） | 不是"永不匹配"的死尺 |
| 提交的 winsec guard-2 匹配式（:234+:235 逐字） | **[]（0）** | 两处站点行为同形 |

**A 面（端到端，真 `portable-tests.sh` ＋ 真 `runtests.sh`，只把 `go` 换 shim）**：
- 对照（不种任何东西）：**rc=0**，`portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/internal/perm`
  ⇒ 真结果行照常计入。
- 种坏（把 perm 的结果行抽掉＋把它的包名只种进一条字符串/日志行，`FAKEGO_DROP_PKG` + `FAKEGO_STDERR_FILE`）：
  **rc=1**；全日志里该包名命中 3 行，逐行分类——`:2` 是种进去的字符串行、`:113` 是 GUARD B 自己点的
  `<NO TOP-LEVEL RESULT LINE> missing`、`:125` 是同段名册行；**执行形状行（`=== RUN/--- PASS/FAIL/SKIP`）命名它的 = 0**；
  GUARD B 逐字 `GUARD B - 1 of 27 packages in scope mode=core printed`。
  ⇒ **"词在日志里"≠"它跑过"**：字符串不救分母，缺结果行照红。

## §3 本格判语

- 匹配式=已做（§1 两站点逐字）；票面要的阳性自证=**本腿补齐**（§2 两发，含正反方向与两处站点）。
  ⇒ 建议非实现者按 §2 的脚本复跑后可翻；⛔ 本腿不翻框。
- 诚实边界（不伪造，与 `scripts/portable-tests.sh:877-883` 的自述一致）：这是**输出行形状匹配**，
  不是 go test 协议性质——在被搜包自身往 stdout 写一条完整伪造结果行的那一层，由 GUARD A + `=== RUN` 计数兜，
  本件不假称已绝。旧两发（票面 Progress log 的合成 capture 与 `/tmp` 快照诱饵）**不重做**，只在此具名。
