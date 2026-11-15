# 454 诚实的 QuickTime 文件被误拒（`.mov` 走不通）

> 立票 2026-11-15 15:4x，由编排者立。来源＝票 214 AC#3 收件（`f563a3c0e`）真跑 10 例探针发现。

**Status**: ready-for-agent
**Blocked by**: 无

## 现象（腿的真跑读数，不是推断）

一个**诚实的** QuickTime 文件（品牌 `qt  `）**声明** `video/quicktime` ⇒ **被拒**，拒的理由逐字是"文件头表明它是 video/mp4"。
同字节内容**声明** `video/mp4` ⇒ **收**。

## 机制

`mp4` 与 `mov` **共用同一个 ISO 谓词**（都走 `ftyp` 那一族的嗅探），而判定表的**顺序让 mp4 先命中** ⇒ 声明为 `mov` 的文件被按 mp4 判，然后与自己的声明对不上 ⇒ 拒。

⇒ 这是**误拒**（false negative），不是安全洞：拒的方向是安全的一侧，但**诚实用户的 `.mov` 用不了**。

## AC

- [x] **AC#1 先复现**：写一枚测试，用**真** QuickTime 头（品牌 `qt  `）＋声明 `video/quicktime` ⇒ 断言"被收"。**它必须先红**，红句原文抄进证据。
- [x] **AC#2 修在一个地方**：要么让 mov 与 mp4 **各自的谓词分开**、要么让"声明与嗅探一致"那一步**接受同一族内的互换**。⛔ 不要靠"把 mov 从允许表里删掉"或"放宽成谁都可以"来变绿。
- [x] **AC#3 不许放宽到危险**：同一条改动后，一个**声明 mp4 但实际是别的族**的文件**仍须被拒**（给会红的测试名）。
- [x] **AC#4 门禁**：gofumpt／vet／d22scan／path-length；`./internal/panel` 定向跑绿。

## 约束

- ⛔ 动 `thresholds.go`／golden／`allowlist.txt`／契约；⛔ 改 `frontend/**`。
---

## 收件（2026-11-15 18:0x，`834aae917` ＋ 证据 `f8ffbf5`／`a5706b7bc`）：**4/4，结案**

**尺1**：三枚提交在；`834aae917` **两件**（`internal/panel/attachments.go` ＋ `attachments_454_test.go`）。
**尺2**：`docs/evidence/s1/454-honest-quicktime-stored.md` **11745 B**。
**尺3 我自己复跑**：`go test ./internal/panel/ -run Test454 -count=1 -v` ⇒ **2 枚 PASS**，其中 AC#3 的 **4 枚子测试全 PASS**（声明 mp4 遇 PNG 拒／遇可执行拒为 unsupported／声明 quicktime 遇 PNG 拒／**声明 mp4 遇 `qt  ` 字节拒**）。

### AC#1 红句原文（修前，逐字）
> `attachments_454_test.go:51: AC#1 一枚诚实的 QuickTime 文件（主品牌 "qt  "）声明 video/quicktime 被拒了：panel: attachment rejected: "holiday.mov" 声明为 video/quicktime，但文件头表明它是 video/mp4`

夹具是**真 QuickTime 头**（20 字节 `ftyp` box、主品牌 `qt  `、minor version、compatible brand），⛔ 不是造字符串。

### 修在一处（现读复核）
`internal/panel/attachments.go:131-133` 两条谓词分开 ＋ `:157 func matchesISOBaseMedia(h []byte, brands ...string)`：
| 谓词 | 只认 |
|---|---|
| `matchesMP4` | `isom`／`avc1`／`mp42`／`M4V `／`mmp4`（**五枚**） |
| `matchesQuickTime` | `qt  ` |

⇒ **判语由品牌自己说，表序不再决定**（这就是原缺陷的根）。

### AC#3 没放宽到危险
品牌并集逐字未变（6 枚）、`mimeAgrees` 与"声明-嗅探一致"那一步一字未动。**两种"放宽"变异都打红**：
- **M2a**（`mimeAgrees` 首行 `return true`）⇒ AC#3 红 3 枚，红句 `refused nothing: ref={… MIME:image/png … Stored:true}`
- **M2b**（mp4 谓词 `return true`）⇒ AC#3 红，红句 `refused nothing: … Name:payload.mp4 MIME:video/mp4`
⇒ **修对了**，不是把门开大。

### A1880 行为牙
**M1**＝把 `attachments.go` 整件还原成修前版 ⇒ AC#1 红，红句与修前**逐字同**；断言对象是 `Ingest` 返回的 ref 与盘上字节 ⇒ **行为红**，⛔ 字面钉。两次变异后还原 `restore-identical=True`（SHA256 `5768CDE8…6240F`）。

### AC#4 门禁（我另核过两份）
`go test ./internal/panel/` 终态 `ok 33.907s` rc=0（基线 42.807s）／`go vet` rc=0／gofumpt 空／`d22scan` clean rc=0／path-length **GREEN**（tracked 20691、over=68＝名册）。
⚠ 它登记了一条**裸 `bash` 命中 WSL 假红**的坑 —— 两个 sh 门必须用 `D:\work\soft\Git\bin\bash.exe`，我认（这条该进台账）。

### ★它主动登记的一处"收紧"
同一枚 `qt  ` 字节**声明 `video/mp4`**：修前 **收**（票 214 探针的对照例）⇒ 修后 **拒**（`声明为 video/mp4，但文件头表明它是 video/quicktime`）。
**我裁：这是对的。** 声明与嗅探对不上就该拒 —— 与"诚实 mov 被误拒"是同一个谓词的**两面**；本票修的是**误拒**那一面，收紧那一面由 AC#3 子测试 4 钉住。若哪天要"两个方向都收"，那是另一张票，⛔ 不在本票场景里。

### 本票 4/4 结案
撤销口令：`git revert 834aae917`
