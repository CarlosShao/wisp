# 35-v6｜编排者本人现跑的"本波真窗"一发（票 35 `:75` (c) 支的载体）

时刻＝2026-10-08 10:1x＋0800；HEAD 起点＝`601c2b18`（跑完后我又落了 `93147902`，动的只有 `docs/reports/**`，产码与夹具一字未动）。
执行者＝编排者（⛔ 不是非实现者腿；这台尺今天没有独立复跑人）。

## 1. 配方（harness，逐字可复制）

```
mkdir -p /d/tmp/wisp35v6 .scratch/wisp/probes/35/v6/logs
cp -f third_party/sherpa-onnx/*.dll /d/tmp/wisp35v6/          # onnxruntime / sherpa-onnx-c-api / sherpa-onnx-cxx-api
go test -c -tags winlive -o /d/tmp/wisp35v6/live.test.exe ./cmd/wisp/   # build rc=0  → logs/build-winlive.txt
cd cmd/wisp                                                    # CWD 必须是包目录
/d/tmp/wisp35v6/live.test.exe -test.run '^TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor$' -test.v -test.count=1
```

## 2. 读数（本波真窗）

- `logs/live-run.txt` 末行＝`LIVE rc=0`，`--- PASS: TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor (8.30s)`。
- 页面自己报回的形状（判据由页面报，不是 Go 侧猜）：
  - `canary arrival=true after 1 post(s); door has recorded [pc-35v2-canary]`
  - `DELIVERED TO GO: method="panel.mode.request" requestId="pc-35v2-pagepost-1" source="panel-composer" to="ask_every_step" (attempt 1)`
  - 门最终录制＝`[pc-35v2-canary pc-35v2-pagepost-1 pc-35v2-desc-ping pc-35v2-DESC-OWN-WR-CF-HOOK-NOERR]`（四枚，含"包裹器在位且没报错"那一枚 `…-HOOK-NOERR`）
- 进程残留尺：日志内 `webview processes at start: tree=0 machineNamed=18 pids=0` → `after teardown: tree=0 machineNamed=18 pids=0`（**同色＝我这发没留窗**）；跑完后我另外在 shell 里量 `tasklist` 的 `wisp.test.exe`／`live.test.exe`＝**0／0**。
- 产码未动尺＝`git status --porcelain -- cmd/ internal/ frontend/ docs/PLAN.md` 现量 **0 行**（本发全程只跑，不改）。

## 3. 同波控制组（两支，⛔ 第三支没做）

- **控制①（tag 缺席）**：`go test -c -o /d/tmp/wisp35v6/plain.test.exe ./cmd/wisp/`（不带 tag，`build rc=0`）→ `logs/list-plain.txt` 里 `grep -c 'TestLive35v2'`＝**0**。⇒ 上面那发读数只能来自 `-tags winlive` 那半边，普通构建里根本没有这枚用例（这正是票 111 `AC#11` 那步存在的理由）。
- **控制②（同波解释器名册）**：`plain.test.exe -test.run 'Forwarding'` ⇒ `logs/control-interpreter-roster.txt`＝`--- PASS` **5** 枚／`--- FAIL` **0** 枚／`INTERPRETER-ROSTER rc=0`。⚠ 这 **5** 枚是"`-test.run 'Forwarding'` 这一把尺的分母"，⛔ 与本票别处引的 10／12／13 枚名册**不是同一分母**，不许互数。
- **控制③（本波反形＝没做，具名）**：⛔ 这一波我**没有**用 overlay 把产码包裹器改坏再跑真窗，所以"这枚真窗用例今天能转红"**没有本波凭据**。历史上它红过那一发（`:52` 记的 `fb2fb802`：录制 `[]`、58.39s FAIL）属**另一波**，我没复跑。⇒ 下一波若再碰传输，控制③要一起做。

## 4. 这一发买到什么／买不到什么

- **买到**：票 35 `:75` 的 (c) 那句"要么给它一个 CI 可达替身，要么在它自己脸上写〔仅本机可量〕并且**每一波碰到传输就复跑**"里"复跑"这一支——本波确实复跑了，且件在盘上（`:75` 因此翻勾，翻勾凭据逐格见工单末节与台账 `A696`）。
- **买不到（⛔ 谁都不许往这里推）**：
  1. 它**不是** CI 载体。winlive 在 CI 里仍只有"编译查"（票 111 `AC#11` 那步 `go vet -tags winlive`），**零执行**；
  2. 它**不是**"这台机器以外也成立"的记录——`machineNamed=18` 那批 webview 进程是本机既有状态，我这发没增没减；
  3. 它**不覆盖** `frontend/**` 那个真页面：用例注入的是宿主自己的脚本，页面分支还没合、构建带页面那格也没落地（票 274 族）；
  4. 它**不证明**时延——见下面第 5 条那枚坏读数。

## 5. ★抓到的仪器缺陷（具名，本程未修）

日志里这一行是**坏读数**：`35v2 cold bring-up measured on this box: -1.000 ms, hwnd=0x1d0d7a`。
负数毫秒不可能是时长，形状像是"哨兵值 −1 没被替换"。⇒ **`cmd/wisp/panel_transport_live_35v2_windows_test.go` 那枚冷启动时延 `t.Logf` 今天不可引用**（它不是断言，所以用例仍 PASS，⛔ 但任何人拿这行去答"面板冷启动多少毫秒"就是拿一个未初始化常量当实测）。本程不改（改它要再跑一发真窗，归下一枚碰传输的写腿；先例排程＝票 33 `AC#13` 那一族同波做）。

## 6. ★§2 那句"产码未动"当场作废（我自己写歪了，附本波第二发读数）

原文 §2 写着"产码未动尺＝`git status --porcelain -- cmd/ …`＝0 行（本发全程只跑，不改）"、§起手写着"HEAD 起点＝`601c2b18`"。**两句都不成立**：并行写腿 `33-r10` 在 **10:16** 落了 `70b00885`，把 `cmd/wisp/panel_host_windows.go` 的 `bringUp` 尾段抽成新函数 `coldStartPageHandover`（＋30／−11）。我那发 `go test -c -tags winlive` 建在 10:1x 与 10:2x 之间，**不能断定编进去的是搬动前还是搬动后的源码**。⚠ 我原文那句"0 行"本身是真的，但它证的是"我没改产码"，⛔ 被我写成了"这一波的产码没动"——这是我造的一句超出证据的话（第 130 条那一族，这次是我自己）。

**补救（同法重跑，读数与 HEAD 绑定）**：

```
HEAD at rerun=2dcef1a6898ac08e54dc54312494c954aeb7ce83      # 含 33-r10 的产码改动
git hash-object cmd/wisp/panel_host_windows.go = 1f9060dfff33cffab1317e5f653e1a95f9678e97
go test -c -tags winlive -o /d/tmp/wisp35v6/live.test.exe ./cmd/wisp/    # rebuild rc=0
/d/tmp/wisp35v6/live.test.exe -test.run '^TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor$' -test.count=1
→ logs/live-rerun-at-33r10.txt：PASS ／ RERUN rc=0 ／ 件里带 HEAD= 行
```

⇒ 现在对**含本波产码改动的树**有一份绑 HEAD 的真窗读数；上面 §2／§3 那些读数**归属仍按本条收紧**：只到"编排者那一发跑过、颜色如此"，不再被读成"整波产码未动"。
