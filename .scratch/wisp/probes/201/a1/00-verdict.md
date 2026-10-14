# 201-a1 探针：junction 卡面标记串 → `WideningRule`/`dirOf` 推出半截「目录」—— **复现成立**

- 身份：探针腿 `201-a1`（**只量不修**）。发送 2026-10-14 16:4x +08；面对：`internal/agent/approval`（唯一写面＝新增一枚 zz 探针测试件）＋`cmd/wisp` **只读**（`114-r6` 在飞）。
- 交件：本文件 ＋ `logs/run3-full.txt`（探针真跑，全部读数逐字）＋ `logs/run4-package-full.txt`（同包全量回归绿）＋ 探针件 `internal/agent/approval/zz_201_a1_junction_marker_windows_test.go`。票 201 尾追加 1 条 bullet。
- **一句话结论**：复现成立、逐字。junction 被拒 ⇒ 卡面 `Paths[0]` 带标记串 ⇒ `ReplyCard.WideningRule`（`replies.go:503`）经 `dirOf`（`batch.go:91`）在**标记里的 `/`** 处切开，产出的「目录」尾部是半截 **`(junction`**；该串两道 `AddAllowedDir` 闸都过、**实测真写进了 config.toml**，并被逐字抬进第二张 L2 卡。
- **最坏后果一句**：连着两次人工答复（`always` ＋ 第二张卡 allow）后，`[fs] allowed_dirs` 多出一条**不是目录**的持久垃圾行；按 `rootsContain` 前缀匹配它命中不了任何真实路径（**惰性行、非直接越权**——此最后半跳是读码推演，见「没试的形」），但它正是 AC#4 要防的「存了一条读不懂的规则」，且其来源是**半截标记**。

## 链路（逐跳 file:line，全部现读＋现跑）

1. `internal/tools/bridge.go:296` `dec.Paths = b.displayPaths(rawPaths)`；`:1020-1034` 对 C26 拒收的路径返回 `p+" (无法规范化: "+err.Error()+")"`——**原样路径＋标记**，不是 canonical。
2. 今天 C26（`internal/tools/paths.go:131-138` → `risk.Resolve`）能返回的错误文本只有两枚：`ErrReparseDenied`（`internal/risk/pathresolver.go:36`，逐字含 `(junction/symlink)`，**带一个 `/`**）与 `tools: empty path`。junction 卡面的标记**只有这一种文本** ⇒ 这条粉碎不是「可能」，是**确定**（下节真跑）。
3. `Decision.Paths` 原样进 `Prompt.Paths`（`gate.go:597-614` `promptFor` 的 `Paths: append(...)`），再原样进 `ReplyCard.Paths`（`replies.go:170-195` `Record`，逐字 `append([]string(nil), p.Paths...)`）；生产落点＝`cmd/wisp/run.go:1386` `u.live.record(p)`（console UI 的 Prompt，无任何剥标记步骤）。
4. `cmd/wisp/approval_always.go:84` `card.WideningRule()` → `replies.go:503-518`：`dirOf(p)` → `batch.go:91-100` `LastIndexAny(t, "/\\")`。标记里的 `/` 是整串**最后一个**斜杠 ⇒ 在 `(junction` 后截断；再 `ReplaceAll("\\","/")` 出成品。`isBareDrive`/`"."`/空 三道拒收均不沾边 ⇒ `ok=true`。
5. 写口：`approval_always.go:134` `s.rt.mgr.AddAllowedDir(dir)`（生产已装配：`run.go:428` `rt.mgr = mgr` ⇒ `:80` 的「配置写面未装配」提前拒不会触发）。`internal/config/allowdirs.go:62-93` 只有两道闸（空、含 `..`），碎片两道都过。

## 逐字读数（真跑，`logs/run3-full.txt`；临时目录名为随机噪声）

```
READING-1  bridge 路由的 Decision: level=L2 rules=[R2 R3] refusal=""
READING-2  卡面逐字 Paths[0] = "C:/tmp/wisp201a1gotmp/TestZZ…/001/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions)"
READING-2b 反斜杠拼法：同一碎片（dir 尾 = "(junction"）也复现 ⇒ 与拼法无关
READING-3  ledger 卡: level=L2 tool=fs.read paths=[同上逐字]
READING-4  WideningRule: ok=true
READING-5  dir  = "C:/tmp/wisp201a1gotmp/TestZZ…/001/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction"
READING-6  rule = "[fs] allowed_dirs += \"C:/tmp/wisp201a1gotmp/TestZZ…/001/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction\""
READING-7  真父目录（对照）= "C:/tmp/wisp201a1gotmp/TestZZ…/001/allowed/jn"
READING-8  allowdirs 两道闸预探: trimspace非空=true 含".."=false
READING-9  AddAllowedDir 接受了碎片，config.toml 现为:
             allowed_dirs = ['C:/tmp/wisp201a1gotmp/TestZZ…/001/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction']
READING-10 第二张卡逐字 Paths[0] = 与 READING-5 的 dir 完全一致
READING-11 对照（无标记的越界读）：paths=["C:\\tmp\\…\\outside\\victim.txt"] dir="C:/tmp/…/outside" ok=true（干净卡面推出真父目录）
READING-12 结论：碎片只在带标记串的卡面上出现
```

（`READING-2/5/6/9/10` 的完整未删节行见 log；`TestZZ…` 为 `t.TempDir()` 随机名。）

## 复现前提（缺一条就哑＝半复现边界，逐条具名）

- 卡面必须发 **L2** 且 `dirOf` 结果**恰好一个** distinct 值：单路径 fs.read 过一个未豁免 junction ＝ 今天是 L2（R2+R3 fail-closed，READING-1）。⇒ 单路径是**最小形**。
- 多路径混排若推出 >1 个「目录」 ⇒ `WideningRule` 返回 `ok=false`（不产碎片，改为响亮 refuse）——本腿未单测此分支（读码）。
- `always` 动词可达另需 `rt.mgr != nil`（生产满足）＋ 人工真答 `always`；写进 config.toml 还需第二张卡被 allow（不答＝C18 拒绝⇒不写）。
- 与分隔符拼法**无关**：前后两种拼法同碎（READING-2/2b）。

## 我没试的形（⛔ 不许当已量）

- **没真跑 `wisp run` 端到端**（cmd/wisp 本腿只读）：`always` 的动词解析/audit 打印行未执行；第二张卡与其后的 `AddAllowedDir` 是**同语句复刻**（`approval_always.go:115-123/134` 字段逐一同值）在审批包内真跑。
- 8.3 短名 / UNC 的卡面：它们不是「拒收」（短名会被 C26 折叠成功）⇒ 今天产不出标记，未试。
- 空路径标记 `" (无法规范化: tools: empty path)"`：该串无斜杠，`dirOf` 会把**整条标记**当「目录」返回（读码推演，未跑；且 fs 工具大概率在参数校验先拒、产不出卡——两处都未实量）。
- 「惰性行」下游：下一启动 `NewPathCanonicalizer` 把它收成 root 后是否真命中不了任何路径 ＝ 按 `rootsContain` 前缀匹配＋该串名字含 `:`（NTFS 不可造）推演，未跑。
- 收件面两处对照（读码，非跑）：`internal/session/grants.go:264` 与 `internal/tools/bridge.go:1273` **都**对 `无法规范化` 特判（会话授权拒存 / 判为 out-of-scope）；**「放宽」这条链是唯一不认标记的消费方**。

## 下一位要点

- 修点候选（只量不裁）：或 `WideningRule` 入口拒收含 `无法规范化` 的 path（对齐 `grants.go:264` 既有口径），或让标记不进 `Paths`（`displayPaths` 改形＝动 `tools.Decision.Paths` 语义面）。归口与排修由编排者定。
- 注释根因：`tools/gate.go:23` 的 `Paths` doc 写「canonical, as judged by R2/R3」，而 `displayPaths` 的拒收支与它**不符**；`batch.go:86-90` 的「already-canonical 所以纯字符串切分无害」前提在标记串上不成立。
- 复跑尺：`GOTMPDIR='C:\tmp\wisp201a1gotmp' go test ./internal/agent/approval -run TestZZ201A1 -count=1 -v`（全包：同命令去掉 `-run`；本腿收件时全包绿 = `ok 1.592s`，见 `logs/run4-package-full.txt`）。
- ⏰ 量取时刻：2026-10-14 16:3x +08（`run3` 16:34、`run4` 16:35 的机器时间）；锚点提交见票 201 尾那条 bullet。
