# 票 159 落地腿 159-r1 实现件（part 2/2）：防再漂的尺 ＋ 反向判据三读数 ＋ 门禁

接 `159-r1-impl.md`（part 1：三处原文、盘上依据、更正块、插入代价）。

---

## 3. "防再漂"的尺：最小可见单位 ＋ 逐名 ＋ 为什么今天响

`cmd/wisp/slo_comment_tense_159_windows_test.go`（新增，`//go:build windows`，`package main`）。
**最小可见单位 ＝ 一枚 (pointer, fragment) 对**，不是"某关键词在文件里出现几次"。

派单要求的"上一轮经验两问"，我先答再动手：

- **最小可见单位是什么？** 一个"某行仍载着某片段"的断言。证据拆到两个单位
  （原文行／更正块的引文行）仍各自可读；摘掉任一半，尺子逐名报出**哪一半**没了。
- **拆到两个单位里还看得见吗？** 看得见，而且**必须看得见**：更正块引用原文，
  所以同一句必然在盘上出现两次；我的判据因此不是"全文件唯一"，而是
  "**更正块之外恰有一处未加标记的家**"（块内的家必须是 `|...|` 标记形）。
  NOTE 这一条是我写尺时**先写错、被尺自己抓出来**才改对的：第一版判"全文件唯一"，
  落地即红（引文本身构成第二处命中）。这正是"恒真/恒红判据"该被排掉的地方。

12 枚子用例逐名（全部 PASS，读数见 sec 4）：

```
the-block-itself-is-still-here
the-corrected-sentences-are-still-on-their-lines
the-correction-sits-after-what-it-corrects
the-block-quotes-its-subjects-and-does-not-paraphrase-them
a-stale-sentence-has-exactly-one-unmarked-home
the-numbers-still-match-their-logs
the-block-states-the-numbers-its-logs-carry
the-correction-still-says-these-sentences-are-past
the-pointer-table-still-designates-its-lines
the-published-md5-still-hashes-the-line-it-names
the-old-numbers-record-a-real-move
the-four-historicised-lines-are-four-lines
```

**为什么它今天响（而不是恒绿）**：sec 4 七发变异逐枚把它打红、逐枚报红名；
其中 M3 只改一个数字、M5 只改一个时态词，都**单独**打红了对应的那一枚子用例。
反向也测过：还原后台面回到绿（sec 4 的 `restore=identical` 与 `ok`）。

---

## 4. 反向判据：七发变异（红）＋还原（绿）＋原样（绿）逐字

台面 ＝ 工作树 `@ff21530f` ＋ 本腿改动；每发**只改被测文件**、跑**只这一枚用例**、跑完 `cp` 还原。

```
M1 drop the (1) :742 quotation row
--- FAIL: TestSLO159Case14CommentKeepsItsTenseAndItsNumbers (0.01s)
    --- FAIL: .../the-block-quotes-its-subjects-and-does-not-paraphrase-them
    --- FAIL: .../the-pointer-table-still-designates-its-lines
    --- FAIL: .../the-published-md5-still-hashes-the-line-it-names

M2 rewrite the historicised sentence (word drift)
--- FAIL: .../the-corrected-sentences-are-still-on-their-lines
    --- FAIL: .../a-stale-sentence-has-exactly-one-unmarked-home

M3 edit a number in the block without re-measuring
--- FAIL: .../the-block-states-the-numbers-its-logs-carry

M4 move the Case 13 block down 3 lines
--- FAIL: .../the-pointer-table-still-designates-its-lines
    --- FAIL: .../the-published-md5-still-hashes-the-line-it-names

M5 restore the tense paragraph to present tense
--- FAIL: .../the-correction-still-says-these-sentences-are-past

M6 copy a stale sentence to a second unmarked line
--- FAIL: .../a-stale-sentence-has-exactly-one-unmarked-home
    --- FAIL: .../the-pointer-table-still-designates-its-lines
    --- FAIL: .../the-published-md5-still-hashes-the-line-it-names

M7 edit the Case 13 header so its published md5 no longer matches
--- FAIL: .../the-pointer-table-still-designates-its-lines
    --- FAIL: .../the-published-md5-still-hashes-the-line-it-names

RESTORE CHECK
restore=identical
PRISTINE (unmutated)
ok  github.com/CarlosShao/wisp/cmd/wisp  0.217s
```

三读数：**变异红 7/7**、**还原绿**（`restore=identical` 后 `ok`）、**原样绿**（`ok`）。
M4／M6／M7 附带打红 md5 子用例是**正确行为**：三发都移动或改动了 `:872` 那一行，md5 断言本就该响。

### 4.1 sha256 三处（提交版／还原台／`git show HEAD:` blob）

```
工作树（＝提交版内容）  a34d1a24a3f7ef2ed58307e2f0f5d497b2333957b261514d113a1ca3b8373df6
还原台 pristine-144.go  a34d1a24a3f7ef2ed58307e2f0f5d497b2333957b261514d113a1ca3b8373df6
cmp                     IDENTICAL
```

NOTE **第三处（`git show HEAD:` blob）与上面两处不同，这是设计如此**：
`06861e1f8a8a6ccab6299b23e7d841739db4051c7eb700c3d43685197d1e903d` 是**本腿改动之前**的 blob
（写这部分时 HEAD 还是 `ff21530f`）。提交之后的三处读数补在 sec 6.3。

---

接 `159-r1-impl-part3.md`（sec 5 门禁、sec 6 整包读数与红名册、提交后 sha256）。
