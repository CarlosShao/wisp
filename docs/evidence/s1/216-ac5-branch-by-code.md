# 票 216 AC#5 · 错误提示按错误码分支（写腿 `216-ac5`，2026-11-14 17:3x–18:1x）

AC#5 原文（票面逐字）：**错误提示按错误码分支，不许按英文错误串里有没有某个数字。**
反面教材：openchamber `ChatInput.tsx:2205-2206`。

## §0 一句话判

- 全仓量到 **3 处**「用户可见句子由英文错误文本选」：2 处已改成按码、各配一枚有牙的尺（突变红／还原绿）；1 处
  （`cmd/wisp/config_reload.go` 的 `describeReloadFailure`，5 分支 6 句）**未改**——它的修法要动 D37 错误模型，
  属「改契约文字」，摆编排者/机主那一格（§5-A 给具名修法与阻塞）。
- panel 那层（`216-r3` 落的 `internal/panel/error_copy_216.go`）我逐字复核并**重打编排者那发突变**：两尺同红
  （5 条具名断言），还原后 SHA256 与还原前同串。
- 另有 **1 处射程外同形**（`internal/tools` 两行英文全串相等，判 errno 存在性，不是文案选择）具名上报，见 §5-B。
- **AC#5 的框不翻**，理由见 §6。

## §1 反面教材：逐字读到的原文（本票起手要求）

上游 openchamber 仓不在本树；盘上唯一副本是调研件 `docs/reports/survey-2026-09-28-composer-plus-menu.md:378-382`
的逐字转录。我把它逐字抄进 `probes/216/ac5/control/openchamber-chatinput-2205.txt` 并读了三遍。原文：

```ts
ChatInput.tsx:2205-2206
if (normalized.includes('payload too large') || normalized.includes('413') || normalized.includes('entity too large')) {
    toast.error(t('chat.chatInput.toast.attachmentsTooLarge'));
```

**票面转述有两处偏差，逐条纠正**（这是「只读到腿的转述、未逐字复核」那句要防的）：

1. 不是「靠 413 子串决定弹哪句」，是**三串 OR**（`payload too large`／`413`／`entity too large`），`413` 只是其中一支。
2. `413` 在这里是**十进制子串**，不是状态码值：响应正文里任何位置出现 `413`（id、offset、模型名）都会命中。
3. 反面不在「数字」这个字眼，在「**判据是散文**」：服务端换个措辞（"the uploaded file exceeds the maximum size"）
   这三支全不中，提示就永远不弹。数字只是这个错误的**显眼症状**，不是病根。

同仓正面对照（调研件 `:385`）：`ExtensionsPage.tsx:58` `if (code === 'too-large') return '...zipTooLarge';`
——**码精确相等**，这是 AC#5 要的形状。

## §2 现量：尺与读数（可复跑）

脚本 `probes/216/ac5/logs/census-ac5.sh`（带 root 参数时走 `git grep --no-index`，可对纯导出树复跑）。
**两发并排**是本节的判据：

- **改前** `logs/census-ac5-before.log`：树＝`git archive 81ba68874` 的纯导出（`cmd/wisp/resident_voice_engines_windows.go`
  与 `internal/llm/errors.go` 的改前字节）。锚差核过：`git diff --stat 81ba68874..HEAD -- internal cmd tools` 空读数。
- **改后** `logs/census-ac5-after.log`（＝`census-ac5.log`）：锚 `HEAD=7150a61f6`。
- 两发之间**只有** A1 从 1 行→0 行、B3 从 3 行→0 行；其余尺同值（B2 5 分支、A2 2 行照旧＝未改项，§5）。

| 尺 | 范围 | 改前 rc／命中 | 改后 rc／命中 | 判 |
|---|---|---|---|---|
| A1 `.Error()` 上的字符串谓词（`Contains/HasPrefix/HasSuffix/Index/MatchString/EqualFold`） | 生产 Go | 0／1 行：`resident_voice_engines_windows.go:1240` | 1／0 | 不合格→已改 |
| A2 `.Error()` 英文全串相等 | 生产 Go | 0／2 行：`system_backend_windows.go:286`／`:426` | 0／2 行 | 射程外同形（§5-B） |
| A3 `switch *.Error()` | 生产 Go | 1／0 | 1／0 | 无人整串 switch |
| B1 `.Detail/.Message/.ProviderCode` 别名 | 生产 Go | 0／1 行：`config_reload.go:431 d := oe.Detail` | 0／1 行 | 下游 B2 |
| B2 在 `d` 上的谓词 | 生产 Go | 0／`config_reload.go:433/437/441/442/460/488`（5 分支）＋`:471` 注释＋`config/allowdirs.go:67`（目录串，非错误） | 0／同 | 不合格→未改（§5-A） |
| B3 `d.Message` 散文匹配 | 生产 Go | 0／`llm/errors.go:181/184/196` | 1／0 | 不合格→已改 |
| B4 上游 `code` 子串匹配 | 生产 Go | 0／`llm/errors.go:205/213/214/215` | 0／同 | **合格**：键是机器令牌，不是散文 |
| C1 数字当字符串（`"413"` 等 9 个） | 生产 Go | 0／**0 生产命中**；2 条命中是**注释里点名反面教材**（`panel/bridge.go:204`、`panel/error_copy_216.go:5`） | 0／同 | 合格 |
| C2 同数字当 int case | 生产 Go | 0／6 行（`llm/errors.go:178/188/195/204/209/219`） | 0／同 | **正控**：这正是按码分支 |
| D1 前端散文嗅探（`includes('413')`／`payload too large`） | `frontend/src` | 1／0 | 1／0 | 合格（只读，零写） |
| D2 前端 `includes(` 正控 | `frontend/src` | 0／38 | 0／38 | 尺不瞎 |
| D3 前端状态比较 | `frontend/src` | 0／7 行 | 0／同 | 比的是**本方枚举**（task/tool status），不是错误散文 |
| E1 panel 文案表键 | `error_copy_216.go` | —／17 键＝D37 17 类 | —／同 | 合格侧 |
| E3 `observe.ClassOf(` 生产消费 | 生产 Go | —／`panel/error_copy_216.go:71` 等 6 处 | —／同 | 码源在场 |
| F1/F2 两句话 | `instructions_200.go`／`workspace.go` | —／五状态 5 常量；workspace 2 常量 | —／同 | 见 §6 |
| 分母 | 生产 Go | —／`.Error()` 274 处、`errors.Is/As` 138 处 | —／同 | 0 不是尺瞎 |

控制：C1 同一把尺在存档转录件上命中（`grep` rc=0）；C2 证明数字以 int 形态在本仓大量存在。⇒「0 个状态码子串」不是盲零。

## §3 逐条判（点名表）

| # | 位置 | 判 | 依据 |
|---|---|---|---|
| 1 | `internal/panel/error_copy_216.go:30-48`（表）＋`:67-92`（分支点）＋`bridge.go:206-216`（消费） | **合格** | 键＝17 个 `observe.Class*` 常量；分支点 `observe.ClassOf`（errors.As）；`err.Error()` 只在 `:91` 当原因后缀，⛔ 不参与选句 |
| 2 | `internal/panel/instructions_200.go:27-40` | 合格 | 五状态是常量枚举，分支键是状态不是句子 |
| 3 | `internal/panel/git.go:365/372/413/468/470` | 合格 | `errors.Is(err, fs.ErrNotExist/ErrPermission)`＝码判 |
| 4 | `cmd/wisp/config_reload.go:433/437/441-442/460/488` | **不合格** | 6 句操作员文案由 `oe.Detail` 的英文词面选；与反面教材同形（只是"散文"由错误正文换成 Detail） |
| 5 | `cmd/wisp/resident_voice_engines_windows.go:1240`（改前） | **不合格** | `strings.Contains(speakErr.Error(), "interrupted")` 决定「中断 vs 失败」、决定是否记 problem（用户可见） |
| 6 | `internal/llm/errors.go:181/184/196`（改前） | **不合格** | `Contains(strings.ToLower(d.Message), ...)` 选错误类 ⇒ 决定 panel 那句（`errorCopy216` 按类选人话） |
| 7 | `internal/llm/errors.go:205/213-215` | 合格 | 键是上游 `code` 机器令牌（`insufficient_quota` 等）；码不是散文 |
| 8 | `internal/tools/system_backend_windows.go:286/426` | 射程外同形 | `e.Error() != "The operation completed successfully."` 判 errno 是否存在，不选文案；但**随系统语言变**＝真 latent bug（§5-B） |
| 9 | `internal/config/parse.go:153` | 合格 | `TrimPrefix(de.Error(),"toml: ")` 是造文案，不是分支 |
| 10 | `config_readers_255.go:507`／`slo_windows.go:960`／`resident_voice_engines_windows.go:483` | 射程外 | 分别是名册标记、文件标记、审计行首格式；都不是错误文本 |
| 11 | `frontend/**` | 合格（只读） | 零散文嗅探；D3 的 7 行只比本方枚举；本腿零写 |

## §4 改了什么（两处，逐枚点名）

### A `cmd/wisp/resident_voice_engines_windows.go`

- 新增包级 `errSpeakInterrupted`（文本仍是 `interrupted mid-speak`，只为日志不动）＋ `speakFailureIsInterrupt(err) bool`
  ＝`errors.Is(err, errSpeakInterrupted)`，**唯一**的 interrupted/failed 决策点（`+33 −3`）。
- `:1212` 改为 `return errSpeakInterrupted`；`:1240` 改为 `if speakFailureIsInterrupt(speakErr)`。
- 语义保证：`MatchaTts.Speak` 把 emit 错误**逐字**返回（`internal/speech/tts.go:212-214`），`CloudTTS.Speak` 文档同样承诺
  （`cloudtts.go:110-116`），`CloudTTSChain` 用 `%w` 包（`cloudchain.go:655`）⇒ 码在三条链上都到得了。
- **牙** `cmd/wisp/speak_interrupt_code_216_windows_test.go` 五臂：码本身真；`%w` 包过的码真；**同文本不同错误必须假**
  （这条是牙）；nil 假；邻居码 `errSpeakerAborted` 假。
- **突变**（盘上现做）：helper 换回 `err != nil && strings.Contains(err.Error(), "interrupted")` ⇒
  `speak_interrupt_code_216_windows_test.go:49` 红：
  `an unrelated error with the text "interrupted mid-speak" classified as a user interrupt - the branch is reading the English text again`。
  还原 ⇒ `ok 0.120s`，文件 SHA256 `A966E8231702DCE1B92A8293438BA6569ED29358705966EE5E86D159155C9580`（突变前后同串）。

### B `internal/llm/errors.go`

- `ClassifyHTTPStatus` 的 400/422 与 404 去掉 `d.Message` 散文匹配（3 个条件）；键只剩**HTTP 数字状态**与上游 `code` 令牌。
  文件头注释写明「message 只作为人类可读原因，永远不是分支键」＋反面教材出处（`+17 −3`）。
- **行为差（如实报）**：只靠散文给「model / context」提示、又不带 code 的上游，今天归 `provider` 而不是 `model`。
  有 code 的上游不受影响（`openaichat/wire.go:146` 填 `codeString()`；`anthropic/stream.go:175-178` 用 code-or-type 兜底）。
  这换来的是：类不再由散文决定，且 `d.Message` 里出现 "413"/"model" 不再改判。
- **牙** `internal/llm/errors_216_test.go` 两枚：①同状态 5 档 × 5 句散文（其中三句正是旧实现匹配过的词）必须同类同可重试；
  ②正控：带 code 时仍然按 code 出类（`model_not_found`→model、`context_length_exceeded`→model、`insufficient_quota`→budget），
  且伴一句「unrelated prose」不能翻案。
- **突变**（盘上现做）：把 `|| strings.Contains(strings.ToLower(d.Message), "model")` 加回 400/422 ⇒ **6 条具名红**
  （status 400/422 × "model not found: gpt-9"／"the model `gpt-9` does not exist"／"context length exceeded..."）。
  还原 ⇒ `ok 0.090s`，SHA256 `8AAA000DD331652923843418785B03DD7F235DDD24FB899695964605E5E9C279`（突变前后同串）。

## §5 没改的两处（具名，附修法与阻塞）

### A `cmd/wisp/config_reload.go:418-501 describeReloadFailure`（5 分支 / 6 句操作员文案）

- **判**：不合格。分支键是 `oe.Detail` 的英文词面（`HasPrefix "config.toml parse"`、`Contains "unknown key"`、
  两串 "cannot migrate..."/"no migration registered"、`Contains "was written by a newer build"`、`HasPrefix "config.toml:"`），
  输出是操作员看到的那一句 `cause=... detail="..."`。形与反面教材同构。
- **为什么这次不改**：把它改成按码，需要错误在**生成处**就带一个 `cmd/wisp` 能读的码。三条路今天都越线：
  1. `observe.Error.ProviderCode`：它是已导出字段，但非空会让 `Error()` 渲染出 `config (<code>): ...`（`observe/errors.go:194-199`）
     ⇒ 改既有 `err.Error()` 全文比对的钉句（钉手 `loader_test.go:273`、`migrate_test.go:160`）。票 231 §4 已把这条路具名登记为「有副作用」。
  2. 在 `internal/config` 加 sentinel + `%w`：同样会改 `Error()` 的尾巴（`observe.Error` 渲染 cause）。
  3. 给 `observe.Error` 加一个**不参与渲染**的 code 字段：最干净，但那是 **D37/C6 错误模型**的字段＝改契约文字。
- **两条路摆编排者**：①走 3（加不渲染字段）＝要机主批契约，收益：全仓错误可机读分类，票 231/232 那族一起解；
  ②走 2（sentinel + `%w`，同批改钉句测试）＝不改契约但改渲染，属票 232 面。
  ★现状的「没有 vs 读不到」两句在这函数里**已经是码判**（`errors.Is(err, fs.ErrNotExist)` vs `fs.ErrPermission`，`:420/:424`），
  改的时候这两支不用动。
- 登记出处：`.scratch/wisp/probes/231/r1/evidence.md:128-167`（票 231 AC#4「只登记不实现」）＋账 `A1010`。

### B `internal/tools/system_backend_windows.go:286/426`

- 形：`e.Error() != "The operation completed successfully."` —— 拿**英文系统消息**判「这次 Win32 调用到底有没有错」。
- **不在 AC#5 射程**（它不选文案，只判错误存在），但这是真 latent bug：非英文 Windows 上这句不等于，成功的调用会被当失败。
- 修法（另一个票的射程，一行级）：`errors.Is(e, syscall.Errno(0))` 或 `if errno, ok := e.(syscall.Errno); ok && errno == 0 { return nil }`。
- 本次**不动**：`internal/tools` 是 D22 ban #7 的扫描面，改动要连带它那套门；且与本票文案无关，混进来会污染这条腿的边界。

## §6 两句话（「没有」vs「读不到」）与为什么框不翻

盘上「两件事两句话」已经成立的三处（都可复跑）：

1. `cmd/wisp/config_reload.go:420-428`：`fs.ErrNotExist`→`cause=missing`、`fs.ErrPermission`→`cause=permission`，**这两支是码判**；
2. `internal/panel/workspace.go:77-80`：`workspaceNotChosenReason`（未选择）vs `workspaceNoReadingReason`（没有读数）；
3. `internal/panel/instructions_200.go:27-40`：`none_found` 与 `not_run` 是两个常量，不是同一个缺省。

错误文案层（`error_copy_216.go`）：17 条头句**没有一条**声称「没有」或「读不到」——6 条是中性的「见下面这句具体原因」，
其余点名子系统（配置/网络/模型/…），原因句逐字后缀 ⇒ 两种事实不可能被同一句吞掉。
**没有为不可达的 fs 分支新加文案**：panel 拒绝路径的错误源是 `ErrAttachmentRejected` 哨兵与解析错误，不带 `fs.ErrNotExist`
（`attachments.go:53/181/290`），加两句话就是恒绿装饰——AC#6 的教训。

**框不翻的理由**：3 处不合格里 1 处（§5-A）未改，且它的修法要动契约（D37 错误模型）⇒ AC#5 不能算「全部落地」。
这一步是编排者的裁决格，不是这条腿能自己盖章的。

## §7 门禁（逐条命令＋读数）

| 门 | 命令 | 读数 |
|---|---|---|
| 构建 | `go build ./internal/... ./cmd/... ./tools/...` | rc=0 |
| vet | `go vet ./internal/llm ./cmd/wisp` | rc=0 |
| llm 全家 | `go test -count=1 ./internal/llm/...` | 6 包全 ok：`122.990s / 10.869s / 34.168s / 0.377s / 58.290s / 52.166s` |
| panel | `go test -count=1 ./internal/panel` | `ok 29.912s` |
| panel AC#5 | `go test -count=1 ./internal/panel -run 216` | `ok 0.083s`（4 枚） |
| panel 突变重打 | `ClassCancelled` 句→`ClassConfig` 句 | 两尺同红，5 条具名：`error_copy_216_test.go:80/85/134/137/140`；还原后 SHA256 `90E6436F…3ADF` 同串、复绿 |
| cmd/wisp 定向 | `go test -count=1 ./cmd/wisp -run "Test381r2\|Test390r1\|Test180K\|Test395r1\|Test26R3\|TestWireR1\|TestSpeakInterrupt"` | 改前 30 项 vs 改后 31 项，**PASS→FAIL 0**；两枚 FAIL→PASS 是基线导出树缺未跟踪 `third_party/spike-models`（基线日志逐字：`config_wire_180k_windows_test.go:256 cannot find third_party\spike-models\silero_vad.onnx`） |
| cmd/wisp 环境红 | 同上 | `WISP_SPIKE_MODELS is not set`／onnxruntime 1.17.1 影子——改前改后同色，非本腿引入 |
| 整包 cmd/wisp | **没跑** | 本机在册 hang（今天新立的票）＋历史整包超时；定向集以外的包级结论本腿不冒充 |
| 禁用形状 | `sh scripts/d22scan.sh` | `clean` rc=0 |
| 路径预算 | `sh scripts/check-path-length-budget.sh` | `VERDICT GREEN` rc=0（20110 条、超 69 全在名册） |
| 格式 | `gofmt -l <四个改动 .go>` | 空 |
| gofumpt | 本机**未安装**（`gofumpt: not recognized`） | 如实记：没有拿 gofmt 冒充 gofumpt |

## §8 可复跑清单

```sh
# 改后（本仓）
sh .scratch/wisp/probes/216/ac5/logs/census-ac5.sh \
  > .scratch/wisp/probes/216/ac5/logs/census-ac5-after.log 2>&1
# 改前（纯导出树：先 git archive 81ba68874 -- . ':(exclude).scratch' | tar -x -C <dir>）
sh .scratch/wisp/probes/216/ac5/logs/census-ac5.sh <dir> \
  > .scratch/wisp/probes/216/ac5/logs/census-ac5-before.log 2>&1
go test -count=1 ./internal/panel -run 216                                  # panel 两尺
go test -count=1 ./internal/llm -run 216                                    # llm 两尺
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test -count=1 ./cmd/wisp -run TestSpeakInterruptIsDecidedByCode216     # 中断/失败码判
```
