# pool-validity-4a / 批4a —— 零勾开放票有效性普查（号段 200–238，11 枚）

> 只读普查腿 `pool-validity-4a`，接 `pool-validity-1`／`pool-validity-2`（死在"扫全池"这个形状）、
> `pool-validity/3a`（100–149，13 枚）、`pool-validity/3b`（150–199，18 枚）。
> **本件只判 200–238 号段那 11 枚零勾票。**
> 起手 HEAD `1a93cc7b`（2026-10-06 11:53:51 +0800，分支 `dev`）。⚠ 换 HEAD 要重量。
> 量尺时刻 `10-06 11:5x`。
> 本腿**零 go 命令**（并行写腿 `232-r2` 此刻真在 `cmd/wisp` 取整包终态）。
> 工具全集＝`git log`／`git grep HEAD`／`git show HEAD:<path>`／`git ls-tree`／`git ls-files`／`ls`／`wc`／`grep -c`／Read。
> 票面零改动：零翻勾、零改名、零 `Status:` 改写、零产码改动、零台账写入、`docs/**` 一字不碰。
> 唯一写入路径＝`.scratch/wisp/probes/pool-validity/4a/**`（新建；代号 `4a` 起手经 `ls` 确认未被占用）。
> ⛔ 不读不引 `frontend/**`、`design/**`、`.gitignore`、`.scratch/wisp/probes/232/**`、`.scratch/wisp/probes/111/ci1/**`。
> ⛔ `docs/evidence/s1/**` 只按**文件名**取（`git ls-tree`），不读内容面。
> ⛔ `cmd/wisp/**` 的工作树内容面不读，一律 `git show HEAD:<path>`。
> ⛔ 前腿件（`batch1`／`batch2`／`3a`／`3b`）只读只学表形，**不照抄其口径、不采信其判档**；
>   `pool-validity/1/**` 与 `pool-validity/2/batch3.md`·`batch4.md`（死腿空骨架）一字不动。
> ⚠ 三枚不判：**232**（写腿正在修它）、**269**／**249**（已撤销／已撤回，判它们＝把已并案的缺陷立案成第二次）。
>   232 出现在我的分母尺原始输出里（它确实零勾、6 格），**是我按指令剔除的，不是量漏**。

## §0 四档尺（本腿判法，口径自立；表形学 3a）

| 档 | 判据 | 复法（命令原文＋读数必须留在表里） |
|---|---|---|
| **仍成立** | 票面点名的那枚缺陷今天还在树上 | 把票面 §现场／AC 里**最硬的一条**自己跑尺：`git --no-pager grep -n '<那句字面或符号>' HEAD -- internal/ cmd/ scripts/ .github/` 的命中数＋file:line |
| **已失效／已被别人做掉** | 缺陷不在了 | ⛔ 硬门：必须给出**一枚 commit 号**＋`git log --name-only` 命中，**或**今树 0 命中的复算尺；拿不出就不许写这一档 |
| **差翻勾** | 缺陷不在了，但凭据是"票内注记／台账 `A##`／一枚非实现者裁决表" | 凭据种类写进复法列；⛔ **不等于本腿裁它可结案**，翻勾归编排者 |
| **量不到** | 判据落在运行期行为／真开窗／真机双击／DPAPI／多账号登录／CI 侧读数／`frontend/**`·`design/**`（本腿禁读） | 归口写清**缺哪一行读数**；⛔ 不许由 grep 命中外推成"这功能能用" |

★ `上次动过` 列＝`git --no-pager log -1 --format='%h %ad' --date=format:'%m-%d' -- <票文件>`（本腿逐枚实跑，见 §1）。
★ 本段预登记的**一处族群陷阱**：234／236 两枚的判据本身就是"CI 侧读数／变异复跑"，
  按 §0 那一档它们的核心格天然落〔量不到〕；本腿要判的是**"票面点名的静态那半边今天还在不在"**，
  两半分开写，⛔ 不许把"我 grep 到了形状"外推成"那一格能结"。

## §1 名册（11 枚，逐枚带票面标题原文）

> 标题＝`head -1` 逐字读，⛔ 未改写。"上次动过"＝票文件自己的最近一笔 commit。

| 票号 | 票文件（`.scratch/wisp/issues/`） | 上次动过 | 票面标题原文 |
|---|---|---|---|
| 200 | `200-nothing-reads-the-projects-instructions-for-ai-although-all-five-harnesses-do.md` | `85da3060` 09-28 | 它进到一个项目里，**完全不知道这个项目的规矩**：四家 harness 全都在读的"给 AI 看的项目说明"，我们一枚代码都没有 |
| 220 | `220-roster-blockedOnApproval-only-sees-L2-and-second-card-silently-drops.md` | `6cf2cc60` 09-29 | **名册上"卡在等批准"那一格今天只看得见 L2，而同一枚任务的第二张卡会静默读成"没被卡住"**（Go 侧送出去的布尔本身就是假的，不是界面没画） |
| 225 | `225-deferred-markers-do-not-match-the-spec12-registry-and-nothing-checks-it.md` | `326fbf78` 10-04 | **`DEFERRED(...)` 代码标记与 `SPEC-12 §5` 登记表今天对不上，而且没有任何仪器在查这一条**：`AGENTS.md` §1.1 那条"1:1 双向"目前是**只靠人**的规矩 |
| 227 | `227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md` | `fac60ad4` 09-29 | 受守卫写还剩两种"手改形状"没人钉：**删掉的行会被补成 schema 默认值**，而 **schema 落后时那句"什么都不写"会说谎**（迁移-on-读照样把文件重排、注释丢光） |
| 229 | `229-subagent-spawns-an-unregistered-goroutine-name-and-drowns-the-leak-alarm.md` | `c774f8da` 09-29 | 子代理收尾时注册了一枚**不在 D38 名册里的协程名**（`subagent-finish-*`）：每 spawn 一次就打一行"疑似泄漏"警告，还被算进 SLO 的协程总数；而"名字不在名册"本该是发现真泄漏时唯一的信号 |
| 230 | `230-four-cells-left-unfinished-inside-closed-tickets.md` | `5759d7fb` 09-29 | 结案票里那 **4 枚真残缺**：票 115 判了七格却**一张裁决表都没出**、票 97 的注释自称"唯一的读者"而同一文件里有三处、票 80 把活交给票 21 而票 21 一字不知 |
| 233 | `233-d36-three-tier-enum-is-tagged-on-no-key-and-the-reload-tier-has-neither-producer-nor-consumer.md` | `fac60ad4` 09-29 | D36 那三档生效级别**一枚键都没标上**：`internal/config.Tier` 三枚常量在盘上只活在 `schema.go` 自己的注释与声明里；`reload` 档**既没有生产者键、也没有消费者**；真正在分档的是 `cmd/wisp/config_reload.go:299` 那枚硬编码 `restartTierKeys` |
| 234 | `234-the-twelve-cells-that-only-need-a-reading-mutation-and-gate-rerun-on-current-head.md` | `6d22dd6c` 09-29 | 结案票残余那 **12 枚"只差一发读数"** 的判据：一次「变异＋门禁」复跑把它们从〔仅自述〕升成〔有读数〕，顺带销掉 10 枚挂在开放票名下的重复分母 |
| 236 | `236-six-cells-that-only-surface-at-the-reading-layer.md` | `c308be29` 10-05 | **六枚只在"读数层"暴露的仪器缺口**（`221-v1` 与 `ci-red-1` 两枚非实现者腿交回）：两支 fail-closed 无尺／那把 DEFERRED 尺的 (b) 支扫词面无牙／"任务 id 必须宿主铸造"零钉／`ci.yml` 那句"今天会红"已过期／**在册红名册不是稳定集合**／**一枚刻意坏的夹具被入库之后污染了格式门的 tracked 分母、并吃掉两道 `go vet`** |
| 237 | `237-race-run-exposes-a-parentheft-inference-that-never-held-plus-two-things-only-a-human-reads.md` | `d5a59d66` 09-29 | `-race` 那一跑撞出的红**不是票 235 造的**：`parkParents` 里"起孩子"排在"父发派生回执"之前、而测试等满孩子信号后立刻要求回执数已满 ＋ 守池注释"只能靠人读"这件事**代码旁边一个字都没写** ＋ 新前置读数只要求"至少一枚"、**部分绕过仍是瞎的** |
| 238 | `238-first-run-on-a-fresh-machine-creates-the-whole-data-root-with-no-seal-to-inherit.md` | `cf03fa31` 09-30 | **新机器上先跑常驻（GUI）那条腿时，整个数据根是"没有任何封条可继承"造出来的**（`internal/proc` 对 `winsec` 导入数＝0，连 `MkdirAll` 的父档都不是窄的） |

## §2 判档表（11 枚 × 4 列）

> ⚠ 本节的**每一行此刻都尚未判**。这不是占位符，这是实话：骨架是先落盘的交件形状，
> 判档由本腿在接下来的调用里逐枚跑尺填入，每判完 3 枚立刻 commit（不攒）。
> 若本腿在填完之前死掉，本节读出来的就是"0 枚有判语"，与死腿 `pool-validity/1` 可分辨——
> 分辨凭据＝上表 §1 的 11 枚名册与下面已实跑的分母尺读数都是**真读数**，不是空壳。

| 号 | 档 | 复法（命令原文） | 读数 |
|---|---|---|---|
| 200 | 尚未判，本腿按票号升序第一枚判它 | 尚未跑尺 | 尚未跑尺 |
| 220 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |
| 225 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |
| 227 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |
| 229 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |
| 230 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |
| 233 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |
| 234 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |
| 236 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |
| 237 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |
| 238 | 尚未判，本腿下一步判它 | 尚未跑尺 | 尚未跑尺 |

### §2-0 分母尺（本腿自己现量，不照抄派单给的数）

命令原文：
```
for f in .scratch/wisp/issues/2*.md; do case "$f" in *-done.md) continue;; esac; \
  b=$(grep -c '^- \[ \]' "$f"); x=$(grep -c '^- \[x\]' "$f"); \
  if [ "$b" -gt 0 ] && [ "$x" -eq 0 ]; then echo "$(basename $f .md) :: 未勾$b"; fi; done
```
读数：`2*` 通配把 20–29／2xx 全吞进来，原始输出 30 行；按 `2[0-3][0-9]` 数值段收窄到 200–238 后得 **12 行**：
200(8)／220(5)／225(4)／227(6)／229(6)／230(5)／**232(6)**／233(6)／234(6)／236(7)／237(3)／238(4)。
剔除 **232**（派单 §3 明令不判：写腿 `232-r2` 正在修它）⇒ **本腿分母＝11 枚**，
与派单 10-06 11:4x 的现量（200／220／225／227／229／230／233／234／236／237／238）**逐枚对上，零差**。
⚠ 号段 239 以上不归本腿（另派）；269／249 不在本尺输出里（已改名或已非零勾，本腿不判、不追）。

## §3 判得心虚的枚数与具名理由

尚未开始填（本腿在判档进行中回填，⛔ 不由编排者代填）。

## §4 一处尺有歧义／判不动的地方

尚未开始填（本腿在判档进行中回填，⛔ 不由编排者代填）。
