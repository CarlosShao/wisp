# 票 133 AC#2 第二格复算 —— 非实现者验收方的第二次裁决（同一格第二次退回即终止）

**执行方**：`acceptor-ticket133-ac2-r2`（验收位，非那枚尺的作者，也非回修方）
**被验的说话**：回修方 `worker-ticket133-ac2-fix` 在票面末尾那 52 行里的那一句 —— "它把退回件的五条账修完了"
**判据物**：票面 AC#2 那一格 ＋ `docs/evidence/s1/133-adversarial-acceptance.md` §3.1／§6／§8
**回修方证据**：`docs/evidence/s1/133-ac2-fix.md`
**本表只裁 AC#2**。AC#1／AC#3／AC#4／AC#5／AC#6 一律未裁。**本表不翻任何勾**（AC#2 由编排者按本表翻）。

## 0. 锚点与树

### 0.1 锚点（开工第一步自己量）

```
$ git rev-parse --short HEAD            048a9e4
$ git rev-parse --abbrev-ref HEAD       dev
$ git show -s --format="%H %ad" --date=iso 048a9e4  048a9e4c0239aaa80722abc20841bad8a395ade3 2026-09-23 22:58:09 +0800
$ git merge-base --is-ancestor fbf420c HEAD && echo yes      yes
$ date -u                               Wed Sep 23 15:00:09 UTC 2026  （开工；+8 ＝ 23:00）
```

派单写"HEAD 约 `048a9e4`"，本程自量也是 `048a9e4`。回修方 `next=` 要求"在 `fbf420c`（或其后代）上复跑"，
`fbf420c` 是 `048a9e4` 的祖先 ⇒ 条件成立。

### 0.2 锚点之后的提交：逐枚判断有没有动 `cmd/wisp/**`

```
$ git log --oneline 048a9e4 ^fbf420c      共 11 枚，全部列出（见下）
048a9e4 docs(122,AC#7 锚点精度)          docs/evidence + .scratch 票面
6e027e4 docs(133,AC#2 回修收尾落盘)      .scratch 票面 133 + docs/evidence/s1/133-ac2-fix.md
21c8def docs(133,AC#2 回修证据 §5/§3.6)  docs/evidence/s1/133-ac2-fix.md
2f291d0 test(136,AC#9①)                  internal/observe/sampler_settle_zerosample_136_test.go
595abd3 docs(136,票面 log)               .scratch 票面 136
36443f2 docs(136,AC#8②③ 实现方证据)      docs/evidence + .scratch 票面 136
ca2b34a docs(133,AC#2 回修证据 §3.3/§3.4) docs/evidence/s1/133-ac2-fix.md
79ddd49 test(136,AC#8②)                  internal/observe/sampler_test.go
1d38206 docs(A134,A135;122,136 翻勾)      .scratch 票面 ×3 + docs/evidence ×2 + docs/reports
7bdfbbb docs(133,AC#2 回修证据 §2.5/…)    docs/evidence/s1/133-ac2-fix.md
ae6a01e docs(136,AC#1 验收方 三处措辞)    docs/evidence/s1/136-ac1-adversarial-acceptance.md
```

机器核过，不靠上面的目测：

```
$ git diff --name-only fbf420c..048a9e4 | grep -c "^cmd/wisp/"     0
$ git diff --stat  fbf420c..048a9e4 -- cmd/wisp/                   （空输出）
```

⇒ **`cmd/wisp/**` 自 `fbf420c` 起一字未动**。本表所有读数取于 `048a9e4` 的归档树，
其中被测那枚尺与 `fbf420c` 逐字节同（见 §0.3），回修方在 `fbf420c` 上取的 §3.3／§3.6 读数
因此与本程同树可比。锚点上唯一未提交面是 `docs/evidence/s1/136-ac8-ac9-impl.md`（兄弟在飞）
与一枚来源未明的未跟踪件（§6 末），**都不在 `cmd/wisp/**`，本程未读、未改、未提交**。

### 0.3 快照树（只在仓外建，**只建不删**；仓内零 worktree／零 checkout）

```
$ mkdir -p /d/tmp/wisp133-r2-tree1 && git archive 048a9e4 | tar -x -C /d/tmp/wisp133-r2-tree1
$ git show 048a9e4:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum   906201f4a3995d10b0a65910aa4cc68e1f745a9d
$ sha1sum /d/tmp/wisp133-r2-tree1/cmd/wisp/leg_dispatch_gate_133_test.go 906201f4a3995d10b0a65910aa4cc68e1f745a9d
$ ls -l cmd/wisp/ | grep -iE "131x|probe|\.off|\.bak"                （空 ⇒ 归档里没有种件残留）
$ cp <repo>/third_party/sherpa-onnx/*.dll <树>/third_party/sherpa-onnx/
    onnxruntime.dll 17799168 / sherpa-onnx-c-api.dll 4605952 / sherpa-onnx-cxx-api.dll 259584（缺即 0xc0000135）
$ go build ./cmd/wisp/   rc=0
```

`906201f4a3995d10b0a65910aa4cc68e1f745a9d` 与回修方 §3.3 自报的 `ARCHIVE-MODULE` 行**逐字同**
⇒ 本程量的就是它交回的那一枚尺。

**每发读数一棵新树**：本程不复用回修方／前任验收方的树，`/d/tmp/wisp133-r2-run.sh` 在每次读数前
`cp -r tree1 wisp133-r2-T<n>`，树已存在则拒跑（`TREE-EXISTS-REFUSING`）。
这样做是为了绕开回修方 §3.2 自己踩到的那个坑 —— `probe.py` 的 `clean()` 会把上一轮留下的
`leg_dispatch_gate_133_test.go.p1.bak` 盖回当前树，让"改后的树"读出"改前的读数"。
本程每发的树都是干净归档的副本，**没有可被盖回的 `.bak`**，所以那类作废读数在本程结构上不可能出现。

### 0.4 台件（复用的只有"种法定义"，没有复用任何人的读数）

前任验收方的两枚台件与本程逐字节同（本程自己 `sha1sum` 过，不是照抄回修方的自述）：

```
$ sha1sum /d/tmp/wisp133-acc-r1-mutate.py /d/tmp/wisp133-ac2fix-mutate.py
  cd88168c05cc14db3410913524ec97992029dc73  （两枚相同）
$ sha1sum /d/tmp/wisp133-acc-r1-probe.py   /d/tmp/wisp133-ac2fix-probe.py
  c44704ab47a89e66325407c4c8104932053e5e63  （两枚相同）
```

⇒ 回修方 §0 那句"复跑台是验收方台件的逐字节副本"**成立**，本程 `p1`／`p2`／`p3`／`p4`／`p5` 与原
裁决表 §3.1 是**同一发件**。本程自己的副本：`/d/tmp/wisp133-r2-{mutate,probe}.py`（`sha1sum` 同上）。
`p6` 是回修方自造的（`/d/tmp/wisp133-ac2fix-p6.sh`），本程按它交接段里写死的形状
**另写一份** `/d/tmp/wisp133-r2-p6.sh` 复种（§1.6 判它造得对不对）。
`p7`／`p8` 是本程自造（§1.7／§1.8），**前任验收方与回修方都没造过** ⇒ 按本仓规矩不记成别人已测。

日志与逐发红名全在 `/d/tmp/wisp133-r2-out/`；树 `/d/tmp/wisp133-r2-T01…T12`（＋后续），**保留不删**。

### 0.5 本程的硬约束（照派单，写给下一位读者）

- `cmd/wisp/**` 只读；被测那枚尺一个字没改（§1 每发的落地证明都是"改树"，改的是快照副本）。
- 在飞兄弟不碰：`worker-ticket136-ac8-ac9` 的 `internal/observe/**`、`worker-ticket137-ac1` 的
  `internal/winsec/**` —— 本程**未读它们未提交的半成品**，`internal/observe/**` 里 `2f291d0`／`79ddd49`
  两枚已提交的件本程也不需要（不在 AC#2 判据物上）。
- 禁改清单（`docs/PLAN.md`／`docs/specs/**`／`internal/risk/**`／阈值／golden／`frontend/**`…）零接触。
- **不跑任何计时类断言**（本机另有代理在飞、`slo-full` 会随 push 自启抢 CPU）：本程所有读数都是
  四数／红名／rc／字节数，无 D32 那两枚数。
- 凭据值零接触：本表只可能出现变量名。
