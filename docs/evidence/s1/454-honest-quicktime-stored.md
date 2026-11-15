# 454 诚实的 QuickTime 文件被误拒（`.mov` 走不通）—— 修好判据表

> 实现腿交件件。票：`.scratch/wisp/issues/454-an-honest-quicktime-file-is-refused-because-mov-and-mp4-share-one-sniff-predicate.md`
> 修前锚＝`6c60148e6`（本腿起手 `git log -1` 现量；交件时本地 HEAD 已由别腿前移，本腿修复落在其上）。
> **修复提交＝`834aae917`**（`834aae9177d44eeed34c5beefccb91ab8143abfc`，两枚件：`internal/panel/attachments.go`＋`internal/panel/attachments_454_test.go`）；本证据件是其后的一笔（中间夹着别腿的 `31b5a6e9`，与本腿无关），所以修复 sha 能在这里逐字写死。

**一句话结论**：mp4 与 QuickTime 共用同一个嗅探谓词时，表序让 mp4 先命中，诚实的 `.mov`（主品牌 `qt  `）声明 `video/quicktime` 被按 mp4 判、再被"声明与文件头不符"拒掉。修法＝**按品牌把两枚谓词分开**（唯一一处改动，`internal/panel/attachments.go`）：`qt  ` 只归 `video/quicktime`，MP4 品牌只归 `video/mp4`。品牌集合**一个都没加**，声明-嗅探一致性那一步（`mimeAgrees`）**一字未动**。

---

## 1. AC#1 修前红句（逐字，真跑）

**测试**：`Test454HonestQuickTimeFileIsStored`（`internal/panel/attachments_454_test.go:42`）。
夹具＝真 QuickTime 头：20 字节 `ftyp` box、主品牌 `qt  `、minor version、compatible brand（`attachments_454_test.go:30`）；声明 `video/quicktime`；断言 `Ingest` 必须**收**（`Stored`、`MIME=video/quicktime`、`Kind=video`、产物扩展名 `.mov`、盘上字节等于源字节）。

**修前命令**（日志 `.scratch/wisp/probes/454/01-ac1-red-before-fix.txt`）：

```powershell
go test ./internal/panel/ -run 'Test454' -count=1 -v
```

**修前输出原文**：

```
=== RUN   Test454HonestQuickTimeFileIsStored
    attachments_454_test.go:51: AC#1 一枚诚实的 QuickTime 文件（主品牌 "qt  "）声明 video/quicktime 被拒了：panel: attachment rejected: "holiday.mov" 声明为 video/quicktime，但文件头表明它是 video/mp4
--- FAIL: Test454HonestQuickTimeFileIsStored (0.01s)
```

⇒ 红句与票 214 收件腿真跑读到的现象**逐字一致**（`声明为 video/quicktime，但文件头表明它是 video/mp4`），不是本腿编的。

**修后**（日志 `.scratch/wisp/probes/454/02-ac1-ac3-green-after-fix.txt`）：`--- PASS: Test454HonestQuickTimeFileIsStored (0.02s)`。

---

## 2. 修在哪一处（file:line ＋ 一句话原理）

**唯一改动面**：`internal/panel/attachments.go`（提交 `834aae917` 的 numstat ＝ `attachments.go` 29+/10−、新测试件 117+/0−；生产件只有这一枚）。

| 位置 | 修前 | 修后 | 一句话原理 |
|---|---|---|---|
| `attachments.go:132` | `match: matchesISOBaseMedia` | `match: matchesMP4` | mp4 行只认 MP4 主品牌 `{isom, avc1, mp42, M4V , mmp4}`（`attachments.go:139`） |
| `attachments.go:133` | `match: matchesISOBaseMedia` | `match: matchesQuickTime` | QuickTime 行只认 `qt  `（`attachments.go:149`）——品牌是两枚白名单类型**唯一**的分界 |
| `attachments.go:157` | `func matchesISOBaseMedia(h []byte) bool` 内含品牌 `switch` | `func matchesISOBaseMedia(h []byte, brands ...string) bool` | 容器形状（4 字节 size＋`ftyp`＋主品牌）与"品牌归谁"拆开：品牌表由两枚谓词各自持有，不再有一枚共用谓词让表序决定判语 |

**为什么选"谓词分开"而不是"同一族内互换"**（票面 AC#2 给的两支）：
- 判定结果就是字节的判语（本文件开头的成文约束：`MIME` 是**嗅探出来的**、声明只是 claim）。品牌分开后 `qt  ` 文件存成 `video/quicktime`＋`.mov`、MP4 文件存成 `video/mp4`＋`.mp4`，**MIME 与扩展名都跟字节走**；"族内互换"那一支会把 QuickTime 字节记成 mp4、命名成 `.mp4`，是另一种"按声明改标签"。
- 表序不再参与判定：任何品牌只命中一行（或都不命中⇒按不受支持拒），原缺陷的成因（顺序让 mp4 先命中）被**移除**而不是被绕过。
- 代价（如实记）：品牌集合与修前**完全相同**，没新增可接受字节（`git diff` 可核：品牌字面量只是从一枚 `switch` 搬到两个调用点）。反向那一格的行为翻转见 §7。

**改动没碰的**：`mimeAgrees`（`attachments.go:365`）、声明-嗅探一致那一步（`attachments.go:271-273`）、白名单种类数、扩展名映射、上限/阈值/golden/allowlist/PLAN/specs/frontend/design、`cmd/wisp/**`，全部零字节。

---

## 3. AC#3 不放宽到危险

**会红的测试名**：`Test454DeclaredMP4WithForeignBytesIsStillRefused`（`attachments_454_test.go:76`），四枚子测试：

1. `declared video/mp4 over a PNG is refused` ⇒ 必须是 `声明为 video/mp4，但文件头表明它是 image/png`；
2. `declared video/mp4 over an executable is refused as unsupported` ⇒ `不受支持`（`MZ` 头）；
3. `declared video/quicktime over a PNG is refused`；
4. `declared video/mp4 over QuickTime bytes is refused, not re-labelled a mov` ⇒ `声明为 video/mp4，但文件头表明它是 video/quicktime`。

修后四枚全 PASS（日志 `02-ac1-ac3-green-after-fix.txt`），且断言 `f.puts == 0`、artifacts 目录零条目（"拒了但落了盘"也抓）。

### 3.1 变异证明 AC#3 有牙（不是字面钉）

| 变异 | 改哪 | 结果 | 日志 |
|---|---|---|---|
| M2a 放宽"声明与嗅探一致" | `mimeAgrees` 首行 `return true` | **AC#3 红 3 枚**（PNG×2、qt-as-mp4）；红句 `attachments_454_test.go:86: refused nothing: ref={... MIME:image/png ... Stored:true ...}`；AC#1 仍 PASS（说明这类放宽对 AC#1 没用、只买来危险） | `mut/m2a-agreement-always-true.txt` |
| M2b 放宽 mp4 谓词 | mp4 行 `match: func([]byte) bool { return true }` | **AC#3 红**（exe 子测试＋qt-as-mp4 子测试）；红句 `attachments_454_test.go:92: refused nothing: ref={... Name:payload.mp4 MIME:video/mp4 ... Stored:true ...}`；AC#1 也红（万能谓词抢在 quicktime 前命中） | `mut/m2b-mp4-predicate-any-bytes.txt` |

两次变异都在**同一棵工作树**上做，改前 `attachments.go` SHA256＝`5768CDE8C6B89B3DC8B6BEFCD22C545E665C3832401C8DEB070ED536C146240F`，每次跑完 `Copy-Item` 还原并**逐字节比 hash 相同**（`restore-identical=True`）。备份与变异文本留在 `.scratch/wisp/probes/454/mut/`。

---

## 4. A1880：行为牙（删掉改动 ⇒ 行为测试红）

**判据**：删掉 §2 那处改动，响的必须是**行为**（`Ingest` 的返回/落盘），不是源码字符串断言。

- 变异 M1＝把 `internal/panel/attachments.go` **整件还原成修前版**（当时本腿还没提交，HEAD 仍是 `6c60148e6`，所以 `git show HEAD:internal/panel/attachments.go` 覆盖工作树就是"删掉本腿改动"的最强形；今天复跑要写 `git show 6c60148e6:internal/panel/attachments.go`），只跑 AC#1：
  - 结果：`--- FAIL: Test454HonestQuickTimeFileIsStored (0.01s)`，红句与 §1 逐字同（`AC#1 ... 被拒了：panel: attachment rejected: "holiday.mov" 声明为 video/quicktime，但文件头表明它是 video/mp4`）。
  - 断言对象＝`Ingest` 的返回 `AttachmentRef` 与 artifacts 目录里的字节 ⇒ **行为红**，不是字面钉。
  - 还原后 `restore-identical=True`（日志 `mut/m1-revert-ac1-red.txt`）。

---

## 5. AC#4 门禁（本腿现跑，读数照抄）

| 门 | 命令 | 读数 | 日志 |
|---|---|---|---|
| 定向用例（全包） | `go test ./internal/panel/ -count=1` | `ok github.com/CarlosShao/wisp/internal/panel 39.882s` rc=0（起手基线 `42.807s` rc=0） | `03-full-panel-after-fix.txt`、`00-baseline-panel.txt` |
| `go vet` | `go vet ./internal/panel/` | rc=0（无输出） | `04-vet-panel.txt` |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/panel` | **空**，rc=0；两枚改动件单跑同样空 | `05-gofumpt-panel.txt` |
| d22scan | `sh scripts/d22scan.sh`（Git Bash：`D:\work\soft\Git\bin\bash.exe`） | rc=0，`d22scan: clean - no D22 ban violations`；自测 `PASS=46 FAIL=0 SKIP=0`；各 scope：bans #1-5 internal/=405、cmd/=86；#6 frontend/=123；#7 internal/tools/=51；#10 internal/risk/=21；#8 design/=39、frontend/=123、internal/=1021、cmd/=302 | `06-d22scan.txt` |
| 路径长度 | `sh scripts/check-path-length-budget.sh` | rc=0 **VERDICT GREEN**（tracked 20691、over-budget=68＝名册 68、not-in-roster=0） | `07-path-length.txt` |

⚠ 机器口径（照实记）：本机 `bash` 裸名先命中 `C:\Windows\system32\bash.exe`（WSL，`/bin/bash` 不存在 ⇒ rc=1、不是门红），上表 d22scan/长度门是用 `D:\work\soft\Git\bin\bash.exe` 跑的；跑前 `Remove-Item Env:GOTMPDIR/TMPDIR`。

---

## 6. 复跑命令（一条条可抄）

```powershell
go test ./internal/panel/ -run 'Test454' -count=1 -v          # AC#1 + AC#3 两枚测试名逐条
go test ./internal/panel/ -count=1                            # 全包
go vet ./internal/panel/
& "$(go env GOPATH)\bin\gofumpt.exe" -l internal/panel
& 'D:\work\soft\Git\bin\bash.exe' -lc 'cd "/d/work/workspace/projects plans/Wisp" && sh scripts/d22scan.sh'
& 'D:\work\soft\Git\bin\bash.exe' -lc 'cd "/d/work/workspace/projects plans/Wisp" && sh scripts/check-path-length-budget.sh'
```

---

## 7. 行为翻转与残余（如实登记，不粉饰）

1. **同一枚 `qt  ` 字节、声明 `video/mp4`：修前"收"（票 214 探针的对照例）⇒ 修后"拒"**，理由逐字 `声明为 video/mp4，但文件头表明它是 video/quicktime`。这是本改动的**收紧**、不是放宽：品牌说了它是 QuickTime，按 mp4 入库会把 `.mov` 字节记成 `video/mp4`＋`.mp4`。已用 AC#3 子测试 4 钉住（`attachments_454_test.go:101-108`）。若编排者认为该方向也该收（族内两向互换），那是 AC#2 的另一支（改 `mimeAgrees` 接受同族互换），**本腿没有做**，按票面"修在一个地方"只选了一支。
2. **反向残余（修前就存在、本腿未改）**：主品牌是 `isom` 等 MP4 品牌、却声明 `video/quicktime` 的文件，修前修后都拒（修前：mp4 先命中＋声明不符；修后：嗅探 mp4＋声明不符）。属"声明错"一侧，方向安全。
3. 未放宽带：品牌集合逐字未增（修前 `{isom, avc1, mp42, M4V , qt  , mmp4}`，修后 mp4 五枚＋quicktime 一枚，并集相同）。
4. 票 214 的 10 例探针里其余 9 例不在本腿写面，未重跑；本腿只对 454 的 `.mov` 一格出读数。

---

## 8. 边界自查（硬约束逐条）

- ⛔ `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／`frontend/**`／`design/**`：零改动。核法＝本腿两笔提交（`git show --name-only 834aae917` 加上证据笔：`git log --format=%H -- docs/evidence/s1/454-honest-quicktime-stored.md | head -1`）只列本腿三枚件（两枚 Go＋本件）；`docs/evidence/s1/` 里另有别腿未跟踪的 evidence 件，与本腿无关。
- ⛔ 附件上限（D36 的配置面）与 `MaxAttachmentBytes`：零改动。
- ⛔ 放宽既有断言：既有测试一枚未改；本腿只**新增**测试件 `internal/panel/attachments_454_test.go`（5244 字节 / 117 行）。
- ⛔ 碰 `cmd/wisp/**`：零改动（没读也没写那边的在飞件）。
- ⛔ push：未推。

---

## 9. 撤销口令

「**撤 454**」＝`git revert 834aae917`（本件 §2 那两行谓词＋`matchesISOBaseMedia` 的 variadic 形一起回退；测试件可留可删，删掉后 AC 读数回到 §1 的红）。若只撤行为、保留证据：把 `attachments.go:132-133` 两行改回 `match: matchesISOBaseMedia` 并把 `matchesISOBaseMedia` 签名改回单参＋`switch`。
