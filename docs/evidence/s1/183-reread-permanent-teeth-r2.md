# 183-r2 — 测试面装牙：AC#9 的第 8 枚常驻腿 / AC#2b 的 CLI 接缝常驻判据 / AC#4 那一形 / AC#6 一句归口 / AC#8 门禁重取

- 角色：**实现者（写码腿·测试面）**。预期**零产码改动**。
- 派单：`.scratch/wisp/dispatches/2026-09-28-135x-impl-183-r2-grow-the-two-missing-teeth-and-re-take-the-gates-after-your-own-last-commit.md`
- 依据表：`docs/evidence/s1/183-pointer-exemption-accept-v2.md`（§2 恒真性矩阵 M3 列零枚红＝AC#9 的来由；§3 四发进攻）· `docs/evidence/s1/183-pointer-exemption-r1.md`
- 交付是**增量**的：本文件先落骨架＋step-0，没做的格标 `待填`——`待填` ＝ **还没跑**，不是**判定通过**。
- 本节以下凡自指文件行号，一律写"符号名＋`grep -n` 尺"，不写死号。

## 1. step-0 五件（本程现量）

**(1) 时刻／分支／起手锚**

```
$ date
Mon Sep 28 13:47:39 CST 2026        # 2026-09-28 13:47 +0800（本机 CST）

$ git rev-parse --abbrev-ref HEAD
dev

$ git rev-parse --short HEAD | sed 's/./& /g'
5 3 0 5 5 0 3 e                      # HEAD = 5305503e
```

⚠ **起手锚差异（按派单要求登记，不据此改判）**：派单写 `98af6535`，盘上 HEAD＝`5305503e`。
`git log --oneline 98af6535..HEAD` 现量＝**恰好多出 1 枚**，且那一枚正是派单 183-r2 自己（`ledger(A366)＋票 183 翻四格补两格＋派 183-r2`）。
⇒ 派单写的是"派单前"的锚，HEAD 靠后一枚是编排者记账＋派单 commit；**按实际 HEAD `5305503e` 做**，终态尺也以此为比较点。

**(2) 在飞脏件闸门**

```
$ git status --porcelain -- internal/ cmd/
（空）
```

⇒ **空**＝没有别人在飞的产码脏件，我不需要 overlay；后续每一次 commit 前复量。

**(3) 三向 md5（起手向＝本程第 3 枚调用现量；交件向见 §6）**

| 尺 | 派单预期值 | 起手实测 | 判 |
|---|---|---|---|
| `sed -n '468,473p' internal/risk/provenance.go \| md5sum` | `858e45116383caa3e7c1dd4b0924fad1` | `858e45116383caa3e7c1dd4b0924fad1` | 一致 |
| `sed -n '11,15p' internal/risk/taintmatch.go \| md5sum` | `5680ddd18e2d2ec2a85e485b54f4c12e` | `5680ddd18e2d2ec2a85e485b54f4c12e` | 一致 |
| `sed -n '482,506p' internal/risk/provenance.go \| md5sum` | `f89e891e5eee3f3ea2b4f89d921072c4` | `f89e891e5eee3f3ea2b4f89d921072c4` | 一致 |

**(4) 判据枚数现量（尺＝`grep -c "^func Test" internal/risk/pointer_183_test.go`）**

起手实测＝**7 枚** ＝ 派单 §0(4) 的本程量值，一致。票面"八枚"与我派单里曾写的"六枚"**都不是现量值**。
⇒ 本程要补的第 8 枚（AC#9）若落地，这一把尺应当从 7 变 8；终态 §6 复量并贴两向。

**(5) 门禁基线（在册）**

```
$ sh scripts/d22scan.sh
…
d22scan: scope ban #8 internal/         examined 432 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations
rc=0
```

⇒ `ban #8 internal/` ＝ **432**，与派单在册基线一致（本程新增判据件会让它 +1，具名登记在 §6）。

```
$ bash .scratch/wisp/probes/154/gate-clauses.sh
…
# BAD  腿=G6neg 声明=ring 基线=1枚 实测=2枚 因=新增未成对（票 171 AC#2：实测 > 基线）
# 腿数＝14 声明与实测不符＝1
# 聚合退码＝1
```

尺＝`grep "BAD" <输出> | grep -oE "腿=G[0-9a-z]+" | sort -u` ⇒ 名册＝**只有 `腿=G6neg`**（与派单在册一致；**比名册不比退码**，退码 1 是它在册）。
本程没跑 `probes/161/r6/flip-declaration.sh`（禁面）。

## 2. AC#9 — 第 8 枚常驻腿的"牙"（M3 变异下必须红）

待填。

## 3. AC#2b — CLI 接缝形状落在 `internal/tools/` 的常驻判据（含回填产生的 `ArtifactPath`）

待填。含**未修码红读数逐字**。

## 4. AC#4 那一形（同名／兄弟路径 × CLI 形状）＋ AC#6 秒数读数与归口

待填。

## 5. 本程没测什么（逐枚具名＋为什么）

待填。

## 6. 门禁终态（取在本程最后一枚 commit 之后，逐枚贴数）

待填。

## 7. 被拒／没成功的调用

待填。⚠ `build failed` 之类 rc=1 **不算判据红**，分开写。

## 8. 有没有跑过删除命令

待填。

## 9. 工具调用枚数 vs 硬顶 40

待填。

## 10. 伪授权两栏

待填。

## 11. 凭据值零抄录

待填。

## 12. next=

待填。
