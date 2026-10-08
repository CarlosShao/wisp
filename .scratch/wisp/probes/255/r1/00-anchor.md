# 255-r1 起手锚（写码腿）

生成时刻 `2026-10-08 11:08 +08`。**本文件在任何 `go` 命令之前落盘并单独 commit。**

## 1. 起手 HEAD（现量，非编排者转述）

`git log -1 --format=%H` ⇒

```
e6c3be1725bb772fce4b94582a02ff00f3664695
```

分支：`dev`。与编排者给的 `e6c3be17` **恰好相同**，但本程不把它当假设：
判据写作「不早于 `e6c3be17`」，共享树随时前进。

## 2. 五件 blob 名册（`git rev-parse HEAD:<path>` 现量）

| 路径 | blob |
|---|---|
| `cmd/wisp/panel_host_windows.go` | `1f9060dfff33cffab1317e5f653e1a95f9678e97` |
| `cmd/wisp/resident_windows.go` | `02a9a0b01ea6133ac8641340cf61185bd6c4db06` |
| `internal/config/tiers.go` | `5ef31255d97a9919cc3a74223c1f5a3f8714819a` |
| `cmd/wisp/config_readers_255.go` | `aa53f254dacdb3d24aede7bbe1dd19252092afdc` |
| `cmd/wisp/config_receipt_255_test.go` | `f6e733c4d18ed231c05abeddf2370f88e6190076` |

## 3. 禁区三方差现量

命令：

```
git diff --name-only e6c3be1725bb772fce4b94582a02ff00f3664695 -- frontend/ docs/specs/ docs/PLAN.md internal/panel/bridge.go
```

输出：**空**（`rc=0`，零行）⇒ 起手时禁区四面对起手 HEAD 无位移。

## 4. 起手工作区备注

`git status --porcelain` 起手即含**他人未提交改动**（`design/**`  deletion、`.gitignore`、
`.scratch/wisp/probes/161/**`、`268/v1/evidence.md` 等）。本程 **一律不 stage、不 commit**，
每笔 commit 带显式 pathspec，⛔ `git add -A` / `git add .`。

## 5. 本程范围（编排者已裁形＝甲形，⛔ 不重开）

1. 甲形落地：`Show`/重新显示那一路取一次当前配置并经 `Dispatch` 应用 `SetSize`。
2. 清掉 `cmd/wisp/config_readers_255.go:17` 与 `:141` 对 `panel_host_windows.go:304` 的已漂行号引用（改**内容锚**）。
3. 无窗仪器：钉「发没发／发的是不是几何源此刻那一对数／hint 是哪枚」；真宽度〔仅本机可量〕。
