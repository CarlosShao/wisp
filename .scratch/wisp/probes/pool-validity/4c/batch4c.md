# pool-validity-4c / 批4c —— 零勾开放票有效性普查（号段 15–26，8 枚最老首轮切片票）

> 只读普查腿 `pool-validity-4c`。没有前文，本件自证。
> 起手 HEAD `3d9b8374`（2026-10-06 13:12:18 +0800，分支 `dev`），量尺时刻 10-06 13:1x +08。⚠ 换 HEAD 要重量。
> 本腿**零 go 命令**（并行写腿 `236-r2` 此刻在 `internal/tools` 取突变读数，抢一把就洗掉它）。
> 枚数／行数一律 `git ls-tree`／`git ls-files`／`wc`；看码一律 `git show HEAD:<path>`／`git grep … HEAD`。
> 票面零改动：零翻勾、零改名、零 `Status:` 改写；台账 `docs/reports/pending-and-issues.md` 与 `docs/**` 一字不碰。
> 唯一写入路径＝`.scratch/wisp/probes/pool-validity/4c/**`（新建；起手 `ls` 验过 `4c` 未占用）。
> ⛔ 零读零引 `frontend/**`、`design/**`、`.gitignore`、`.scratch/wisp/probes/236/**`、`.scratch/wisp/probes/pool-validity/4b/**`。
> ⛔ `cmd/wisp/**` 与 `internal/**` 的工作树内容面不读（写腿地界），一律走 `HEAD:` 尺。
> `docs/evidence/s1/**` 只按**文件名**取、不读内容。
> 本段与 §0 之外的所有读数都是**本腿现量**，不采信任何旧腿结论；号段 17–20、27 以上不归本腿。

## §0 四档尺（本腿判法，学 `.scratch/wisp/probes/pool-validity/3a/batch3a.md` §0 的表形，口径自立）

| 档 | 判据 | 复法（必须留命令原文＋读数） |
|---|---|---|
| **仍成立** | 票面点名的缺陷／缺口今天还在树上 | 把票面 AC 里**最硬的一条**自己跑尺：`git --no-pager grep -n '<符号或那句字面>' HEAD -- internal/ cmd/ scripts/` 的命中数＋file:line；或 `git ls-files '<那个包>/*_test.go'` 之类存在性尺 |
| **已失效／已被别人做掉** | 票面要的东西已经在树上 | ⛔ 硬门：一枚 commit 号＋`git log --name-only` 命中，**或**今树读到实现件在场的复算尺；拿不出就不许写这一档 |
| **差翻勾** | 东西已在，但凭据是"票内注记／台账 `A##`／裁决表" | 凭据种类写进复法列；⚠ 不等于本腿裁它可结案，翻勾归编排者 |
| **量不到** | 判据落运行期行为／真麦克风／真 TTS 播放／真开窗／真机双击／CI 侧读数／`frontend/**`·`design/**`（禁读） | 归口写清**缺哪一行读数**；⛔ 不许由 grep 命中外推成"这功能能用" |

★ 这批是 9 月开的首轮切片票，写的时候产品还不存在 ⇒ **"票面点名的东西今天存不存在"只由尺回答，不由名字像不像回答**。
★ 切片级验收票（票 16／票 25）的射程盖的是**别的票的产物**：本腿不替它判"整片通过"，只回答"它点名的那几枚产物今天在不在树上"，其余进 §4。
★ `上次动过` 列＝`git log -1 --format='%h %ad' --date=format:'%m-%d' -- <票文件>`，本腿逐枚实跑。

## §1 名册（8 枚，票面标题原文）

分母尺（本腿 10-06 13:1x 现量，逐枚精确 glob，不含同前缀他号）：

```
for f in .scratch/wisp/issues/15-*.md .scratch/wisp/issues/16-*.md .scratch/wisp/issues/21-*.md .scratch/wisp/issues/22-*.md .scratch/wisp/issues/23-*.md .scratch/wisp/issues/24-*.md .scratch/wisp/issues/25-*.md .scratch/wisp/issues/26-*.md; do
  b=$(grep -c '^- \[ \]' "$f"); x=$(grep -c '^- \[x\]' "$f");
  echo "== $(basename $f) 未勾=$b 已勾=$x"; done
```

读数原文：

```
== 15-speech-engines-cer-harness.md 未勾=6 已勾=0
== 16-s2-acceptance.md 未勾=5 已勾=0
== 21-approval-gates-minimal.md 未勾=5 已勾=0
== 22-web-tools-d30.md 未勾=6 已勾=0
== 23-system-window-input-tools.md 未勾=5 已勾=0
== 24-doc-search-tools.md 未勾=5 已勾=0
== 25-s3-acceptance.md 未勾=5 已勾=0
== 26-tts-output.md 未勾=7 已勾=0
```

⇒ 8 枚**已勾全部＝0**，与编排者 10-06 13:1x 的现量一致；本腿无"已被别人翻过勾"要具名跳过的票。

| 票号 | 票文件（`.scratch/wisp/issues/`） | 票面标题原文 | 未勾格数 | 已勾 |
|---|---|---|---|---|
| 15 | `15-speech-engines-cer-harness.md` | 15 — Speech engines: VAD + streaming ASR via sherpa, engine slot mutex, CER harness | 6 | 0 |
| 16 | `16-s2-acceptance.md` | 16 — S2 acceptance: voice end-to-end gate, latency segments, hotplug pre-mortems | 5 | 0 |
| 21 | `21-approval-gates-minimal.md` | 21 — Approval gates minimal: L1 block window, L2 native card, queue trivial, batch D45-1 | 5 | 0 |
| 22 | `22-web-tools-d30.md` | 22 — Web & open tools: web.search/fetch/open, file.open, app.launch, D30 five layers | 6 | 0 |
| 23 | `23-system-window-input-tools.md` | 23 — System/window/input tool family: system.*, window.*, clipboard, media, screen, input.type | 5 | 0 |
| 24 | `24-doc-search-tools.md` | 24 — doc.read + local search tools: PDF/docx extraction, search.files/content/apps | 5 | 0 |
| 25 | `25-s3-acceptance.md` | 25 — S3 acceptance: capability scenarios ①② + full security red-team suite | 5 | 0 |
| 26 | `26-tts-output.md` | 26 — TTS output: sherpa matcha engine, serial slot, Speaking pipeline, P7 gate | 7 | 0 |

## §2 判档表（8 行 × 4 列）

| 票号 | 判定 | 复法（命令原文＋读数） | 备注／归口 |
|---|---|---|---|
| 15 | 尚未判，本腿下一步判它（票号升序第 1 枚） | 预取事实已在 §4 前置尺里，判档时重跑 | — |
| 16 | 尚未判，本腿下一步判它（票号升序第 2 枚，切片级验收票，射程受限见 §0★） | 判档时只答"它点名的产物在不在树上" | — |
| 21 | 尚未判，本腿下一步判它（票号升序第 3 枚；本枚有一整半截在界面侧，按 §3 边界会落"量不到"） | 判档时逐 AC 跑尺 | — |
| 22 | 尚未判，本腿下一步判它（票号升序第 4 枚） | 判档时逐 AC 跑尺 | — |
| 23 | 尚未判，本腿下一步判它（票号升序第 5 枚） | 判档时逐 AC 跑尺 | — |
| 24 | 尚未判，本腿下一步判它（票号升序第 6 枚） | 判档时逐 AC 跑尺 | — |
| 25 | 尚未判，本腿下一步判它（票号升序第 7 枚，切片级验收票） | 判档时只答产物在不在 | — |
| 26 | 尚未判，本腿下一步判它（票号升序第 8 枚） | 判档时逐 AC 跑尺 | — |

## §3 判得心虚的枚数与具名理由

本节当前**尚未写**，本腿判完 §2 之后按票号逐枚写满；⛔ 不留空。

## §4 一处尺有歧义／判不动的地方

本节当前**尚未写**，与 §3 同步补满；已预登记两处，届时展开并补其余。

- 预登记①：票 15/26 的 `internal/speech` 存在性尺——本腿 13:2x 现量 `git ls-tree -r --name-only HEAD internal/speech/ | wc -l` ＝ `1`，唯一文件 `internal/speech/doc.go`；同尺 `internal/audio/` ＝ `16` 枚。旧探针读数一律不采信。
- 预登记②：票 21/23/26 的"界面半边"判据落 `frontend/**`·`design/**`（本腿禁读）＋真机行为，按边界必须落〔量不到〕，不许绕。
