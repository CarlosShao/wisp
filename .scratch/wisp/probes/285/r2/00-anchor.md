# 285-r2 · 起手锚

## 取数（`date` / `git` 现量，非名义）

- 时刻：`2026-10-09 11:04 +0800`（`date '+%F %H:%M %z'`）
- 锚 sha：`ce18b3d6cc3b6cc03066d65f108053ae68fcc1d7`
- 锚标题：`A760 收 285-r1＋reconcile-8：票 285 三格翻勾（AC#4 按住）＋台账七笔承重数逐笔复跑更正`
- 分支：`dev`

## `git log --oneline -6`（sha＋subject 前段；两枚长句在 `：` 处截断，全文以 `git log` 现量为准）

```
ce18b3d6 A760 收 285-r1＋reconcile-8：票 285 三格翻勾（AC#4 按住）＋台账七笔承重数逐笔复跑更正
1d6bcfc3 reconcile-8 自纠笔数：汇总表 A755 行与合计（对上 61→60，复跑不了列补上脚本那一发）
df778832 reconcile-8 只读对账：A751-A756 承重数复跑 71 处（对上 61／对不上 7／复跑不了 4）
fec4f5e7 票 285 腿 285-r1 时钟订正：本腿两处名义时刻（11:2x／11:5x）改成 date 与 git log 现量（锚 commit 97392ab＝10:32，读数 commit f5452c9＝10:48，收工 10:52）；票面 Progress log 标题同批，⛔ 判据句与四格框未动
23b29215 A759 落账＋票 279 两格翻勾（AC#2/AC#3，AC#1 整格按住撞 285-r1 独占面）
6378718c 票 285 腿 285-r1 就地订正 00-anchor：唯一用例名口径 61（未量下笔）→ 现量 72 枚，附取数命令
```

## 工作树快照（在飞件＝别人的，一律不碰）

`git status --porcelain | wc -l` = **759**，按状态码分布：

```
    16  D
    17  M
   726 ??
```

已知在飞面（本腿零读写）：`M .gitignore`、`design/**` 的 16 枚删除与 4 枚修改、大量
`.scratch/wisp/probes/**` 与 `.scratch/ci-logs/**`、仓根一枚名叫 `-` 的未跟踪文件。

## 产码面在起手的状态（本腿收工要与之逐字等值）

`git status --porcelain -- internal cmd` 在起手＝**空**（原文无输出）。

起手 hash（`md5sum`）：

```
8eb37e9f59e563b58325d2ebc898b843 *internal/agent/loop.go
018dc1480d64c6825eb7e0abe1a63ced *internal/tools/ticket283_corr_identity_rulers_test.go
```

## 本腿范围

票 `.scratch/wisp/issues/285-cross-card-denial-name-unasserted-at-route-level-plus-four-green-shapes.md`
的 **AC#4**（票末「追加一格」节 · 第五形：per-call corr 的「每枚调用各不相同」今天零尺）。
新建会响的断言 + 四件套突变自证（两发，全种在盘上，零 `-overlay`）。
⛔ 不碰 `internal/agent/approval/gate.go`（284 随批搭载格不触发）。
⛔ 不翻勾、不改判据句、不 push。
