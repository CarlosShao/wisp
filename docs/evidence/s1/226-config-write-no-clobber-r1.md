# 票 226 — 交付证据件 r1（写腿 `226-w1`）：`always` 不再把整份内存快照写回 `config.toml`

- 工单：`.scratch/wisp/issues/226-always-answer-rewrites-the-whole-config-snapshot-and-hides-hand-edits.md`
- 台账出处：`docs/reports/pending-and-issues.md` 的 `A428` 第 ② 条（缺陷怎么被发现）＋ `A424` 第 1 条边界（"持久、可审计、可撤销"）
- 治理：根目录 `AGENTS.md` §1.1（`D/C/R` 契约与 SLO/golden/`thresholds.go` 一字不动）、§1.2（`tools/d22scan` 禁形）、§1.4（只 commit 不 push、commit 带显式 pathspec）
- 本件作者＝**实现者本人**，不是裁决者；按 `AGENTS.md` §0 第 3 句与 `SPEC-12 §4.3` #1/#3，**缺口审计与对抗验收必须由另一个 agent 做**，本件不构成自证。
- 定锚（起手 HEAD）：`4540d5b11abce0522f811068b10c902647c06846`（committer 时间 `2026-09-29T12:25:21+08:00`，现跑 `git log -1`）
- 起手时刻：`2026-09-29 12:26 +0800`（现跑 `date`）

## 0. 一句话结论

写路径改成**写前重读磁盘＋只替换这一节的那一枚键**（票面候选**甲**），并把 `statOwnWrite()` 的认领范围收窄到"这次写产生的内容与内存一致"那一支；`SetPermissionMode` 同形那半（AC#4）**同批改掉**，没有留残缺。

## 1. 选了哪一支，为什么（票面把这一支交给自己定并写明理由）

〔待填〕

## 2. 改了哪几枚文件（包级路径＋现跑 `grep -n` 行号）

〔待填〕

## 3. 门与读数

〔待填〕

## 4. 逐格判据

〔待填〕

## 5. 没做完的、卡在哪、要新增却没敢造的名字

〔待填〕

## 6. Git 纪律

〔待填〕
