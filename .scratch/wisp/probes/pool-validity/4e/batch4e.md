# pool-validity 4e — 死腿 4a 名册收尾（233／234／237）＋ 六批重叠普查

> 本腿代号 `pool-validity-4e`。只读普查腿，档位口径四档：**仍成立／已失效（已被别人做掉）／差翻勾／量不到**。
> 本腿不写票面、不动台账、不碰 `docs/**`、零 go 命令、只 commit 不 push。

---

## 0. 起手锚

- HEAD 短号：`95cb3d7a`
- 分支：`dev`
- 时刻：`10-06 15:13`（`date '+%m-%d %H:%M'`）
- `git status --porcelain -- internal cmd` 原文（起手第一次读数）：

```
?? internal/tools/tasklist_deferred_236r3_teeth_test.go
```

⇒ 与"另一枚写腿正在做 236 AC#2"的告知对得上（`236-r3` 在 `internal/tools` 落了未跟踪的突变件）。本腿因此**不读 `internal/**` 工作树内容面**，一律 `git --no-pager grep ... HEAD`。

- 落盘位置检查：起手 `ls .scratch/wisp/probes/pool-validity/` 原文：

```
1
2
3a
3b
4a
4b
4c
4d
msg-a633.txt
```

⇒ `4e` 此前**不存在**（编排者那句"4e 未用过"成立），本腿按指定用 `4e`，不启用备选名 `4e2`。
★附带事实具名：`pool-validity/` 下除编排者列的六批（`2/`、`3a/`、`3b/`、`4a/`、`4b/`、`4c/`）之外，还有 **`1/`** 与 **`4d/`** 两个目录。`4d` 是本腿的禁读前缀（在飞腿地界），不读、不采信；`1/` 不在六批名册里，本腿在 §3 里把它当**补充对账面**读票号列（见 §3.4），⛔ 不据此改任何一批的档。

---

## 1. 三枚名册（本腿实际要判的枚）与分母现量

命令原文（编排者给的尺，本腿自跑）：

```
for n in 233 234 236 237 238; do f=$(ls .scratch/wisp/issues/${n}-*.md 2>/dev/null | head -1);
  printf "%s 文件=%s 未勾=%s 已勾=%s\n" "$n" "$(basename $f)" "$(grep -c '^- \[ \]' $f)" "$(grep -c '^- \[x\]' $f)"; done
```

读数原文（15:13 现量）：

```
233 文件=233-d36-three-tier-enum-is-tagged-on-no-key-and-the-reload-tier-has-neither-producer-nor-consumer.md 未勾=6 已勾=0
234 文件=234-the-twelve-cells-that-only-need-a-reading-mutation-and-gate-rerun-on-current-head.md 未勾=6 已勾=0
236 文件=236-six-cells-that-only-surface-at-the-reading-layer.md 未勾=7 已勾=0
237 文件=237-race-run-exposes-a-parentheft-inference-that-never-held-plus-two-things-only-a-human-reads.md 未勾=3 已勾=0
238 文件=238-first-run-on-a-fresh-machine-creates-the-whole-data-root-with-no-seal-to-inherit.md 未勾=4 已勾=0
```

⇒ 与编排者"236 有写腿在做 AC#2、238 已被 4b 判过"的口径**没有冲突**：五枚**全部 0 已勾**，即没有任何一枚被勾掉过。
⇒ 编排者那句"读数与我上面那句不一致时以你的为准"：上面读数**未见冲突可言**（编排者没给具体枚数，只给了"5 枚尚未判"的名册），本腿照原样上报。

三枚的判与不判：

| 号 | 本腿动作 | 具名理由 |
|---|---|---|
| 233 | **判** | 4a 表里逐行写着"尚未判" |
| 234 | **判** | 同上 |
| 236 | ⛔ **不判** | 写腿正在做它的 AC#2（`236-r3`，`internal/tools` overlay 突变，起手锚里那枚未跟踪件即其证据）。判它会与它互相污染。 |
| 237 | **判** | 4a 表里"尚未判" |
| 238 | ⛔ **不判，只当参照样本读** | `4b/batch4b.md` 已判（仍成立＋状态半边已被机主关死，台账 `A476`）。本腿只读它的**口径形状**，不重判。 |

---

## 2. 判档表（233／234／237）

尚未算，本腿下一步算它。名册与四列形状已就位：号｜档｜复法命令原文｜读数。

---

## 3. 六批重叠普查

尚未算，本腿下一步算它。射程：`2/`(46)、`3a/`(13)、`3b/`(18)、`4a/`(11 名册/6 判)、`4b/`(8)、`4c/`(8)。只抽**票号列与合计行**，⛔ 不采信各批判档，⛔ 不重判、不改档。

---

## 4. 判得心虚的枚数与具名理由

尚未写，本腿到约第 100 次调用前先补这一节。

---

## 5. 尺有歧义／判不动的地方

尚未写，本腿到约第 100 次调用前先补这一节。已定形的一条：`pool-validity/1/` 与 `4d/` 不在编排者六批名册内（见 §0）。
