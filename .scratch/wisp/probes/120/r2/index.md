# 120-r2 探针目录（落地腿自留：跑了什么、怎么复跑、盘上留了什么）

起手锚 `58de465a`（dev）。写面只有两枚：`internal/winsec/winsec_other.go`（产码）、
`internal/winsec/posix_seal_window_120_other_test.go`（仪器）。证据件＝`docs/evidence/s1/120-fd-anchor.md`。
本目录不放 `.go`（规矩），只放脚本与读数件；临时件全部**只建不删**。

## ① 三发测量各是什么

| 发 | 代码字节 | 仪器 | 档 | 日志 |
|---|---|---|---|---|
| 复现 r1 | 58de465a 原样（`os.Chmod(path, mode)`） | r1 那份读者件 | 读者 | `120-prefix-reader-reading.log` |
| 成对·前半 | 58de465a 原样产码 ＋ **本轮**测试文件 | 本轮仪器 | `WISP_120_BASELINE=1` | `120-prefix-newinstrument-reader.log` |
| 成对·后半 | 本轮产码（`pinPlacement` 钉 fd） | 本轮仪器 | 默认严格档 | `120-postfix-gate-count2-v.log` |
| 反控 | 本轮产码，`sameObjectPinned` 两枚判据各 `&& false`（盘上物理改写） | 本轮仪器 | 默认严格档 | `120-control-gate-neutered.log` ＋ `120-control-driver.log` |

快照都在仓外 `/d/tmp/probe120/r2/`：`base-head`（改前原样）、`base-head-newinstr`（改前产码＋本轮测试件）、
`base-fixed`（本轮两枚文件）、`permface`（权限面台件，独立 module，未入库）。

## ② 复跑命令（按顺序；全部现量口径）

```sh
cd 'D:/work/workspace/projects plans/Wisp'
export GOTMPDIR='C:\Users\swq\AppData\Local\wisp-test-gotmp\120r2gotmp'
S=/d/tmp/probe120/r2

# 成对·前半的快照怎么造的（改前产码 + 本轮测试文件）
rm -rf "$S/base-head" "$S/base-head-newinstr"; mkdir -p "$S/base-head" "$S/base-head-newinstr"
git archive 58de465a | tar -x -C "$S/base-head"
git archive 58de465a | tar -x -C "$S/base-head-newinstr"
cp internal/winsec/posix_seal_window_120_other_test.go "$S/base-head-newinstr/internal/winsec/"
find "$S/base-head-newinstr" -name '*.go' | wc -l        # 先证非空再测：1072

# 半发读数（容器；⛔ 用 -overlay，这枚仪器运行时读盘上产码）
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e WISP_120_BASELINE=1 \
  -v "$S/base-head-newinstr":/work -v probe120build:/root/.cache/go-build -w /work golang:1.27 \
  sh -c 'go test -count=1 -v -run "TestPOSIX120" ./internal/winsec/'

# 成对·后半 + 门禁：直接用本目录脚本
sh .scratch/wisp/probes/120/r2/gate2.sh   # 建 base-fixed、跑 AC#3 名册（后半 uid 四发死于已知 bug）
sh .scratch/wisp/probes/120/r2/gate3.sh   # uid 0 / uid 1000 的 permface 四发 + 整包两发
sh .scratch/wisp/probes/120/r2/control2.sh # 反控：去牙→红→还原→证 sha
sh .scratch/wisp/probes/120/r2/gate-helper.sh # 从 count=2 日志里取具名计数
```

## ③ 读数速记（细节与判读口径见证据件）

- 成对·前半（改前产码）：三形 `hit-foreign` 7 / 8 / 11，分母各 2000 ⇒ 合 **26/6000 > 0**，凭据第一半到手。
- 成对·后半（改后默认档，`-count=2`）：六发全 **0/2000**；活性同发非零
  `refused-at-check=1003/114/138（第 2 发 993/133/164）`、`refused-at-act=11/8/23（7/5/51）`、
  `sealed-own-object=864/128/181（860/109/191）`；flip 错误 0；整包 `RUN=110 PASS=110 FAIL=0 SKIP=0`。
- 反控：形 3 四枚样本全红（8+4+9+4 ⇒ `hit-foreign=25/2000`），`control rc=1`，
  还原后 sha 与 clean 同值（`RESTORE MATCHES`），仓内文件 sha 全程未变（`137442cd…`）。
  已知射程缺口：`sealHandle` 自己那半 `SameFile` 没被去牙，也没有具名单元台件。
- 权限面（r 位收紧的代价）：uid 0 改前改后九形全 `SEALED`（失真档）；uid 1000 改后 5 形转 `REFUSED`
  （`file-0200`／`file-0000`／`dir-0300`／`handle-0200`／`dir-0300-parent`），整包两发仍 `ok`。

## ④ 本目录件清册

脚本：`gate-helper.sh`、`gate2.sh`（后半四发有已知 bug，`MSYS_NO_PATHCONV=1: command not found`）、
`gate3.sh`（补跑）、`control.sh`（首版，死于 `cp` 到不存在的 `/work/.scratch`）、`control2.sh`。
读数：`120-prefix-reader-reading.log`、`120-prefix-newinstrument-reader.log`、
`120-postfix-gate-count2-v.log`、`120-ac3-roster-113-names.log`、`120-control-gate-neutered.log`、
`120-control-driver.log`、`120-host-rulers.log`、`120-d22scan-postfix.log`、
`120-root-permface-base-{head,fixed}.txt`、`120-uid1000-permface-base-{head,fixed}.txt`、
`120-uid1000-package-base-{head,fixed}.log`、`exp-readout.log`、`exp2-readout.log`、`tmp-postfix-smoke.log`、
失败件 `120-uid{0,1000}-permface-{prefix,postfix}.txt`、`120-uid1000-package-{prefix,postfix}.log`（gate2 bug 所致）。

`exp-readout.log`／`exp2-readout.log`＝写码前对 `os.Root` 语义的现量（兄弟相对链接会被跟、绝对/`..` 目标才拒、
`O_NOFOLLOW` 经 `Root.OpenFile` 失效、`.` 指根本身、`rename(2)` 换祖先的时序面、FIFO 上纯 `O_RDONLY` 阻塞 1.5s
vs 带 `O_NONBLOCK` 115µs）。这两发决定了产码里"⛔ 把 `O_NOFOLLOW`／Root 的链接规则当防线"那句与
`O_NONBLOCK` 的取舍。

## ⑤ 没做完的（逐枚具名）

1. `next=2 补集`：`sealHandle` 的同一对象复验无具名台件，反控也未去它的牙。
2. `next=4 补集`：改后默认档没把"每形每发 `refused-at-act>0`"升成硬判据（今天只 `t.Logf`）。
3. `next=5`（r1 转下）：票 120 票面那句 `154/4000` 仍是改前口径，改措写要人批 ⇒ 本腿只追加 1 条 bullet。
4. 射程外读数：`internal/risk`、`cmd/wisp` 等消费者在 Linux CI 上会不会因 r 位收紧变色，本腿写面之外，未量。
5. `privateFile` 开那一跳（先 `Create` 再 `Open` 再封）＝**票 322、未修**，本腿⛔动。
