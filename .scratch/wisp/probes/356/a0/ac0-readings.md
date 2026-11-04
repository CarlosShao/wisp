# 356 AC#0 · R3 A 档成对读数（2026-11-04，编排者现跑，另起尺不抄 174-v1）

探针＝`internal/tools/zz_356_ac0_probe_test.go`（build tag `probe356`，常规测试计数不受影响）。
命令＝`go test ./internal/tools/ -tags probe356 -run Test356AC0R3DenyPaired -v -count=1` ⇒ **PASS**。
造形＝`~/.ssh/id_ed25519` A-list 文件（`internal/risk/blacklist.go:244` 判 `ClassA`＝R3 绝对 Deny）；用**approving gateSpy**（`approveAns/windowAns=AnswerAllow`）在场，这样计数=0 才证明"R3 短路、根本不会出卡"，而不是"NoGate 拒"。notice 取 `announceOf(out.Text)`。

## 三形读数（逐字）

| 形 | 构造 | notice 原文（announceOf） | fs.read isError／level／class | approval 被问 |
|---|---|---|---|---|
| **i 正控：R2-only** | 普通文件在授权目录外（非 R3） | （省略，同 R2 两路文案） | `false`／`L2`／—（整读回 bytes=20000） | **1**（出卡且批准即读回） |
| **ii 越权 + R3 Deny** | `.ssh/id_ed25519` 在授权目录外 | 「…注意：这条路径现在读不到，它不在你被授权的目录范围内：fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡，全文见 <path>」 | `true`／`L2`／`permission_denied` 正文=「该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：R2: 目标路径在授权目录之外: <path>; R3: 命中 A 档敏感路径，拒绝（不可豁免）: <path>」 | **0** |
| **iii 授权内 + R3 Deny** | `.ssh/id_ed25519` 在授权目录内（home 整棵入 allowed_dirs） | **无告警句（warn 0 字节）**，只剩「…全文见 <path>…」指针 | `true`／`L2`／`permission_denied` 正文=「该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：R3: 命中 A 档敏感路径，拒绝（不可豁免）: <path>」 | **0** |

## 归因（这就是 356 要修的缺陷）
`pointerNotice`（`internal/tools/task.go:831`）**没有任何一支处理"被判 Deny"**：
- 形 ii 只命中 `!InAllowlist` 支 ⇒ 吐出"两路"，但 R3 Deny 下**加目录与批卡都无效、根本没有卡**（approval=0 实证）＝**点一张永不出示的卡**。
- 形 iii `InAllowlist==true`、`os.Stat` 成功、`IsRegular==true` ⇒ 全部支跳过 ⇒ **连"读不到"都不说**（warn 0 字节），而 fs.read 实为 `permission_denied`。
- 正控 i 证明尺分得清：同尺下 R2-only approval=1、整读回 20000B。

⇒ AC#0（只读现读＋归因）**达标**：notice 原文、fs.read 的 isError/level/class/正文、approval 次数（R3 两形=0／R2 正控=1）全逐字给出，另起尺。
