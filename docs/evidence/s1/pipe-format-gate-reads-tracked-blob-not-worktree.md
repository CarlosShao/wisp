# `gofumpt` 判据口径裁定：按「入库 blob」读，⛔ 按工作树副本读

- 裁定人：编排者（按 `AGENTS.md` ⓪⓪⓪ 铁律第 4 条：**技术选型 agent 自己拍**）
- 时间：2026-10-21 22:2x +08
- 适用范围：票 **321 `AC#2`**／票 **298 `AC#5`**／以及**今后任何**"格式门禁"类判据
- 性质：**判据口径裁定**（⛔ 改契约文字、⛔ 改任何阈值、⛔ 改任何 `.go` 字节）

---

## 一、问题（一句话）

本机跑 `gofumpt -l` 会**恒红一枚**（`cmd/wisp/models.go`），而**入库的那份字节是干净的**。
⇒ 于是两张票的门禁格**在本机永远跑不出绿**，而它们**本该是绿的**。

---

## 二、铁证（两条读数，各自可复跑）

```
$ git cat-file blob HEAD:cmd/wisp/models.go > /tmp/x.go
$ /d/work/base/gopath/bin/gofumpt.exe -l /tmp/x.go     # 对「入库 blob」
0                                                       # ← 干净
$ /d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/models.go   # 对「工作树副本」
cmd/wisp/models.go                                      # ← 红一枚
```

**同一份内容、两处落点、两种答案。** 差别只有一个：

| 事实 | 读数 |
|---|---|
| 本机 `git config --get core.autocrlf` | **`true`** |
| `.gitattributes` 对 `*.go` 的规定 | **`*.go text eol=lf`** |
| 该枚**工作树**副本里的 CR 数 | **334** |
| 该枚**入库 blob** 里的 CR 数 | **0** |
| `git hash-object <工作树>` ⇄ `git rev-parse HEAD:<path>` | **同串（`4259a20f…`）** |

⇒ **`autocrlf=true` 在检出时把 LF 转成了 CRLF**，而 `gofumpt` 读的是工作树那个副本 ⇒ **红的是"本机检出的副本"，⛔ "我们入库的字节"**。
★这也解释了为什么它是**孤例**：`cmd/wisp` 140 枚 `.go` 里只有它一枚带 CRLF。

---

## 三、为什么按 blob 读是**主流 harness 的做法**（三条，均可核）

1. **★CI 自己就是这么定的，而且写在注释里。** `.github/workflows/ci.yml:189` 那一步的名字逐字是
   `"gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)"`，注释里逐字写着
   **"the denominator of the formatting gate is the TRACKED set, not the working tree"**（`ci.yml` 那段的原话）。
   ⇒ **裁定不是我发明的口径，是把仓库已经声明过的口径落到本机尺上。**
2. **主流格式化门禁读的是"要合进去的东西"。** 任何 code-review 门禁关心的是 **PR 里那棵树**（＝入库字节），
   ⛔ "某台机器磁盘上检出后的副本"。否则**本机 `autocrlf` 配置会改变门禁颜色**——那是**环境决定判据**，不是代码决定判据。
3. **反向也成立**：若按工作树读，则任何 `eol` 转换（本机配置、编辑器、`core.autocrlf` 三种都能改）都会让门禁**假红**，
   而**真问题**（有人提交了未格式化的字节）反而可能被这枚假红**淹没在噪音里**。

---

## 四、裁定（硬）

1. **今后凡"格式门禁"类判据**，一律跑在**入库 blob** 上：
   `git cat-file blob HEAD:<path> > <临时文件>` 再对临时文件跑 `gofumpt -l`。
2. ★**必须同时报两把**（⛔ 只报一把）：
   - **(甲) blob 尺**（**判据用这把**）；
   - **(乙) 工作树尺**（**只作诊断**：若它红而甲绿 ⇒ **本机 `autocrlf` 假红**，登记一句、⛔ 判该票失败）。
3. **⛔ 因此判任何票"格式门禁不通过"**，除非 **blob 尺**确实列出该文件。
4. **⛔ 修那枚 CRLF 副本**（`git add --renormalize` 之类）当成本票的活——那是**本机检出状态**，不是仓库内容；
   要不要统一本机 `autocrlf` 是**另一件事**（属"动本机配置"，我另记、⛔ 顺手做）。
5. ★**这条口径不改 CI**：CI 上每个文件都是 tracked、检出即 LF（Linux runner），两把尺本就同值 ⇒ **裁定只影响本机读数**。

---

## 五、这条裁定解开的两张票

| 票 | 格 | 之前为什么卡 | 裁定后 |
|---|---|---|---|
| **321** | `AC#2` 门禁五把 | `gofumpt` 红 1 枚（`models.go`） | **blob 尺 0 枚 ⇒ 这一把过**；其余四把本已 `rc=0` |
| **298** | `AC#5` gofumpt 残留 1 枚 | 同上 | **同上** |

★两票**共用同一枚环境事实**（这也是裁定要一次做、⛔ 一票一次的原因）。

---

## 六、我**没有**做什么（边界，⛔ 读大）

- ⛔ 改 `cmd/wisp/models.go` 一个字节（入库 blob 本来就干净）。
- ⛔ 改 `.gitattributes`（`*.go text eol=lf` 是对的；问题在本机 `autocrlf=true` 覆盖了它）。
- ⛔ 改 CI 那一步。
- ⛔ 动任何阈值／golden／断言。
- ★**⛔ 声称"本机工作树是干净的"**——它确实带 CRLF，那是**事实**；裁定说的是"**判据读 blob**"，⛔ "工作树没问题"。
