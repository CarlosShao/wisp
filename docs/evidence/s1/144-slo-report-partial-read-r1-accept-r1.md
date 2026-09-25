# 144 — 对抗验收 r1（非实现者）：`collectReport` 分"没写完 / 坏了"这一味

- 验收程：本件作者，**不是** `95885fb` / `f63f0e3` / `a91d7c2` 的作者，进场前对本票无先验上下文。
- 被验锚点：**`1e94672`**（本程自己 `git cat-file -t 1e94672` = `commit`；`git diff --stat 95885fb 1e94672 -- cmd/wisp/slo_windows.go` 现量为空 ⇒ 收笔那句"本票零改动 `slo_windows.go`"**成立**）。
- 本程进场时刻仓库 `HEAD` = `b694378`（分支 `dev`，本程自己 `rev-parse` 量的）——**HEAD 已被别家程推走**，
  所以本件所有读数出自 `1e94672` / `95885fb^` 的**仓外纯净副本**，不出自工作树。
- 派单里编排者写的每一句（含"读数是成立的"那几句）在本件里一律按**未验证断言**处理；
  下面每一格给"现量命令 + 读数 + 判定"，判定档位只有 **成立 / 退回 / 退回（附条件入账、文件保留、不许勾）**。

## 0　副本与两件自证

### 0.1　纯净副本配方（两个 `-c` 都带，并自证生效）

```
$ cd "D:\work\workspace\projects plans\Wisp"
$ git -c core.autocrlf=false -c core.eol=lf archive 1e94672 | tar -x -C D:\tmp\wisp-144acc1\snap
$ cd /d/tmp/wisp-144acc1/snap && gofumpt -l . ; echo rc=$?
rc=0            ← 输出为空 ⇒ 副本没有被烙上 CRLF，两个 -c 确实都必要
```

建了四棵仓外树（全部在 `D:\tmp\wisp-144acc1\` 下，**只建不删**）：

| 树 | 来源 | 用途 |
|---|---|---|
| `snap/` | `1e94672` archive + 本程那枚验收驱动 | 修后形状，跑本程自己的差分 |
| `pre/` | `95885fb^` archive + 同一枚驱动 | 修前形状，跑**同一发**差分 |
| `prer/` | `95885fb^` archive，**无驱动** | AC#4 改前名册（`ls prer/cmd/wisp \| grep -c slo_report_144` = **0** ⇒ 真改前树） |
| `post/` | `1e94672` archive，**无驱动** | AC#4 改后名册 |

DLL 不在 git 里（现量 `git ls-files third_party` 为空）⇒ 四棵树各从工作树 `cp third_party/sherpa-onnx/*.dll`，
`PATH` 用 `/d/…` 形（前一程 §3.4 那发 127 的坑本程照它躲过了，未"顺手优化"那枚脚本）。

### 0.2　五枚 AC 框此刻的勾态与本程之前的表

```
$ grep -nE '^\- \[[ x]\] \*\*AC#' .scratch/wisp/issues/144-*.md
52:- [x] **AC#1 …    58:- [x] **AC#2 …    61:- [x] **AC#3 …    64:- [x] **AC#4 …    70:- [x] **AC#5 …
$ ls docs/evidence/s1/ | grep -i 144
144-slo-report-partial-read-r1.md          ← 只有实现方自证，无 accept 件
$ ls docs/evidence/s1/144-slo-report-partial-read-r1-accept-r1.md
No such file or directory
```

⇒ 编排者那句断言**复算相符**：五格全 `[x]`、非实现者的表在本程之前**不存在**。
本件就是那枚表；五格在下面**逐格出档**，不做整体判断。票面的勾本程**不动**（地界外）。

### 0.3　先量最承重的那一发：真修前树 vs 真修后树，同一枚真半截文件

前一程 §3.4 的那张差分表用的是 **667/668 字节**（它自己那枚夹具）。本程**不拿它的表当凭据**，
另写一枚只调**生产入口** `sloSubject.collectReport()` 的驱动（`acc144driver_windows_test.go`，
同文件同内容放进两棵树，仓内一字未改），自己造真半截文件：

**这份半截文件"是不是真的半截"——不是本程挑的形状，是主体自己的序列化器产的**：

```
ACC1 full doc via writeSLO = 1950 bytes; last byte = 0x7d (0x0a would mean trailing newline)
ACC1 0-byte file exists and is readable (0 bytes) => a reader CAN see a file before its bytes arrived
```

`writeSLO` 现量（`95885fb^` 树 `cmd/wisp/slo_windows.go:628-638`）= `json.MarshalIndent` 整块进内存 →
**一次** `os.WriteFile(out, data, 0o644)`，无临时文件、无 rename、尾字节 `0x7d`（不是换行）
⇒ 读者期间可见的**只有那份文档的前缀**，本程取 `full[:len(full)-1]` 就是**它能出现的最真的一形**。
前置自检：`bytes.HasPrefix(full, pre)` 为真、`json.Valid(pre)` 为**假**（否则这发不携带信息）。

**同一发在两棵树上的读数（逐字）**：

```
### pre/  (95885fb^)
ACC2 raw json.Unmarshal on the 1949-byte prefix = unexpected end of JSON input
ACC2 tree prefix: file=1949 bytes of a 1950 byte document
ACC2 collectReport error = wisp slo: subject report: unexpected end of JSON input
ACC2 elapsed=0s (budget constant … =33s; a first-read death returns far under it)

### post/ snap (1e94672)
ACC2 collectReport error = wisp slo: subject 4242 never wrote a complete report within 33s (last read: 1949 bytes read, document still open at offset 0: the tail had not arrived)
ACC2 elapsed=33.02s
```

⇒ §3.4 那张表的**句子形状**本程复算相符（字节数不同只因夹具不同），且**"修前第一发即死"是真的**
（`elapsed=0s` vs `33.02s`，两形用的是同一枚预算常量）。

⚠ 本程自己那发 `--- FAIL: TestACC2StaticPrefix` 是**本程驱动的断言写坏了**
（`elapsed > budget` 漏算最后一枚 `time.Sleep(poll)`，33.018s > 33s），**不是被验对象的红**；
它反过来是一枚正控：**预算常量确实没被挪**（否则那一发不会正好落在 33s 那一格上）。

**本程补的一发（前一程没取、且它才是"值不值"那一问的答案）—— ACC3**：
半截文件先在盘上，尾巴 **400ms 后**到达；只调生产入口，问"这份报告收不收得回来"：

```
pre/   ACC3 VERDICT=error     : wisp slo: subject report: unexpected end of JSON input
snap/  ACC3 VERDICT=collected  mode=subject-in-tree pass=true
```

⇒ **摘掉这一味，外部可见的不只是"句子变短"，是"一份真会到达的报告被丢掉、整发红"**。
这条是本票"承重"的操作定义下**唯一需要的外部读数**，见 §3（AC#3 那一格）。

### 0.4　本程撞到的、两程都没量的一件事：那一发竞态在本机**撞得出来**

前一程 §3.3 的结论是"无争抢时读者落不进那次写"，续程据此把外部读数判成"两版相同"。
本程拿同一份真报告做 reader-vs-writer 硬撞（写侧循环 `os.WriteFile` + `os.Remove`，读侧循环 `os.ReadFile`，
各 3 秒），并带**已知正控**（一枚手截 100 字节的文件确实可被 `os.ReadFile` 读到）：

```
pre/   ACC4 positive control: a hand-truncated 100-byte file IS readable … json.Unmarshal says unexpected end of JSON input
pre/   ACC4 partial read: *json.SyntaxError len=0 x998
pre/   ACC4 partial reads observed during 3s of concurrent writeSLO = 998
snap/  ACC4 partial read: *json.SyntaxError len=0 x871
snap/  ACC4 partial reads observed during 3s of concurrent writeSLO = 871
```

⇒ 三点，前两点是**新增读数**、不是复算：
① 那一形**在这台机器上可重放**（871–998 次/3 秒），所以票面"AC#1 自证作废"那一支**不必走**，
   前一程 §6 第 1 条"从没在真进程上把那一发竞态召出来"是**没做**、不是**做不到**；
② 观察到的**唯一**部分形状是 **`len=0`**（文件已建、字节还没灌），本程**没有**撞到 1..1949 的中间前缀
   ⇒ 生产里真正救回来的是 `io.EOF`-on-0-bytes 那一支（`reportUnwritten` 的 0 字节分支），
   `io.ErrUnexpectedEOF` 那一支在本机今天走不到；两支同归 `reportUnwritten`，所以修法仍然对得上，
   但**"半截文件"这个票面措辞与实测量到的形状不完全是一回事**，登记给归因；
③ 这两发**没拿"多少秒"当判据**：871/998 是**次数**，且两版数量同级 ⇒ 它证的是"这形造得出"，不是"哪版快"。
