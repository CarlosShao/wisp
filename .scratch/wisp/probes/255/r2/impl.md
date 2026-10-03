# 票 255 · 实现腿 `255-r2` · 证据件

写面＝`cmd/wisp/`＋`internal/config/`（本腿只做 AC#1 主格＋AC#3 残余补尺）。
本件由实现腿自己写，**判语不归自己**：AC 勾选框一枚未碰。

---

## 0. 起手锚（同一发命令取数）

```
$ date
Sat Oct  3 09:10:04 CST 2026
$ git log -1 --format=%h
8a3790f0
$ git status --porcelain -- cmd internal tools docs .scratch | wc -l
352
$ git status --porcelain -- cmd internal
（空＝本腿写面起手干净）
```

- 起手时刻 `2026-10-03 09:10 +08`，锚点 HEAD `8a3790f0`。
- `cmd internal` 两枚目录 porcelain＝**0 行** ⇒ 工作树在写面上与 HEAD 同字节（AC 判据里"工作树＝HEAD"的那把尺成立）。
- 全仓 porcelain 352 行主要来自 `.scratch/`（别的腿的临时件，规则＝只建不删），与写面无冲突。
- 开工前已完整读票面：`.scratch/wisp/issues/255-config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers.md`（55 行，含编排者 2026-10-02 11:00 裁定小节与 13:1x 行号更正）。

**AC#1 判据原文（逐字抄自票面第 17 行）**：

> AC#1 那句"已立即生效"只许说真话：判据＝stdout 里"这些段已立即生效"那一段**只列真有生产读者的段**；**正控**＝手改 `[panel] width`  ⇒ 那句话**不许出现 `[panel]`**（要么不出现、要么换成诚实形如"值已换、但当前无组件应用它"）；**反控**＝真有读者的段（逐名给出证据）**仍要出现**。⛔ 不许把整句删掉——那是把热加载的可见性摘了，本仓定式：⛔ 摘尺不许当修尺。

**AC#3 残余补尺判据（逐字抄自票面第 19 行编排者翻勾段末句）**：

> 范围披露：`[app]` 新增键不在键级尺射程（固定名单非反射走查）——我裁＝记范围不算欠账，补尺（一行循环）归 `255-r2` 顺手做

---

## 1. 现量（回执那几行 · `rep.Hot` 的来源 · 逐名"谁读它"）

（取数中）

---

## 2. 两格各自的修法与代码

（取数中）

---

## 3. 门禁四数（真实读数）

（取数中）

---

## 4. 变异自证表（种什么形 → 哪枚必须红 → 复跑终值）

（取数中）

---

## 5. 我可能写错的条目（自我对抗）

（取数中）

---

## 6. 判不动的地方

（取数中）

---

## 7. 交件判语

（取数中）
