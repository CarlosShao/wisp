# 180-a1（10-02 第二程）—— 配置面全名册现量：哑键普查 · `unwiredKeys` 覆盖率 · 生效级别三档 · "此项不生效"的出口 · `[panel] width` 断链点

派单：只读取证腿 `180-a1`，10-02 09:5x 起手。**产码零字节改动**；票面 AC 框一枚未碰；票面正文未改（只在末尾**追加**一节）。

---

## §0 起手锚（逐字）

- `date -Iseconds` ⇒ `2026-10-02T09:57:47+08:00`
- `git log -1` ⇒ commit `00e7efeb0ca22e1ce3cbe62763d0a3088dfe4330`（"ledger(A521)：167-a1 结档…"，Author CarlosShao，Date Fri Oct 2 09:57:39 2026 +0800）
- `git rev-parse --abbrev-ref HEAD` ⇒ `dev`
- `git status --porcelain internal/config internal/panel cmd/wisp` ⇒ **空输出**（三目录起手干净）
- 全树 `git status --porcelain`：**不干净**——别人的脏件与未推临时件大量在场（`design/**` 一批 D/M、`.scratch/**` 一批 ??、`scripts/testdata/portable-tests/go` M）。本程**不还原、不提交**它们的任何东西。
- 本程读到的**所有行号均取自 HEAD `00e7efe`**（除注明）。票面锚是 `039efb47`、上一程 `180-a1` 锚是 `38fc7c0e`——**两处行号都已过期**，差集见 §8。

## §1 全名册：叶子键 → 读者 → 非测试／仅测试／无

（表头先行；枚数在 §1 正文，尺与射程逐条附。）

| 段 | 叶子键（TOML 路径） | Go 字段 | 生产读者 file:line | 判 |
|---|---|---|---|---|
| — | — | — | — | — |

## §2 `unwiredKeys` 名册覆盖率

## §3 "生效级别"三档的现量

## §4 "此项今天不生效"有没有出口 ＋ 最近可扩展点

## §5 `[panel] width` 那一环断在哪

## §6 我可能判错的条目

## §7 判不动的地方

## §8 我推翻票面与派单哪一句
