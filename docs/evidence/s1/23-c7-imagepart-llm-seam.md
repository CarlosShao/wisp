# 票 23 / AC#2 — C7 ImagePart 到 LLM 的接缝（实现腿交付记录）

- 票面：`.scratch/wisp/issues/23-system-window-input-tools.md`（Status: ready-for-agent，Blocked by 20-host-bridge-fs-tools = `20-host-bridge-fs-tools-done.md`，已 done）
- 基线 HEAD：`2f7f9de16`
- 本腿写面：`internal/agent/**`、`internal/llm/**`（含聚焦 `*_test.go`）+ 本文件
- 本文件是**实现证据**，不是裁决表：AC 的对抗验收须由非实现者另做（SPEC-12 §4.3）

## 1. 三核前提（动手前现跑）

| 核 | 问题 | 读数 |
|---|---|---|
| (a) 票面点名的字符串/行还在不在 | `llm.ImagePart{MimeType,BytesRef,Alt}`、`Request.ResolveBytes`、D34 的 screen.capture 元数据渲染 | **在**：`internal/llm/content.go:43-47`、`internal/llm/provider.go:83-86`、`internal/tools/system.go:242-264` |
| (b) 票面引用的 file:line 对不对 | AC#2 指「C7 ImagePart→LLM 接缝」；票面进度日志称「需 internal/agent loop + llm.ImagePart 集成」 | **对**。HEAD 上三协议适配器都**显式拒绝**工具结果里的 ImagePart（报 `only TextPart`），agent loop 也从不产生 ImagePart；`internal/tools/system.go:262-263` 逐字写着「C7 ImagePart 到 LLM 的接缝本仓当前不存在」 |
| (c) 有没有吃别人的在飞改动 | `git status` 在写面内 | 只有一枚**未跟踪**的 `internal/agent/ticket23_image_seam_test.go`（09:59:35 落盘，早于本腿接手约 20 秒）。它不是本腿写的，且**当时不可编译**（缺 `internal/llm` import、引用尚不存在的 `ImageAttachment`/`newImageStore`/`maxImageBytes`）。⇒ 按「临时件只建不删」保留，本腿先补 import 再把它当 RED 用，见 §3。其余写面干净（无并发腿写 internal/agent、internal/llm） |

Blocked-by 复核：`20-host-bridge-fs-tools-done.md` 存在 ⇒ 不阻塞。Status 与派单一致（ready-for-agent）。

本腿只做 **AC#2 的 loop+llm 半边**（C7 ImagePart 进历史、进线上字节、内存句柄不落盘/库/日志）。AC#2 的 tools 半边（`tools.Result` 图像通道 + bridge 映射 + `system.go` 那句「接缝不存在」）在写面外，见 §7 STOP-note。AC#3 需真机 fixture app + 前端确认卡（写面外/机主外部 agent），AC#1/4/5 已由前序腿结清。

## 2. 设计（为什么是这条）

C7 已经把协议差异挡在适配器之外，`ToolResultPart.Content` 注释本来就允许「richer media when the protocol allows」。所以：

- **agent 侧**：图像放进 `ToolResultPart.Content` 的 `llm.ImagePart`（协议中立）；`BytesRef` 是 `mem:<n>` 内存句柄，`Request.ResolveBytes` 只在**编码上线那一刻**把句柄解成字节。
- **llm 侧**：三种协议各按自己的线格式落：
  - anthropic：`tool_result.content` 原生接受 image block ⇒ 图像留在**同一个工具结果里**；
  - openai-responses：`function_call_output.output` 接受 input 数组 ⇒ 图像留在**同一个调用结果里**；
  - openai-chat：协议只允许 role:"tool" 里放文本 ⇒ 主流做法是文本留在 tool 消息、像素挂在**紧随其后的 user 轮**并标注是哪个调用返回的。这条不是妥协：三家都会拒 role:"tool" 里的图像，旧适配器拒的方式是把整个请求判失败，等于「截一张图 = 这段对话在该 provider 上死掉」。
- **内存 only**：`imageStore` 不碰文件系统，句柄无路径、无分隔符；单图上限 4 MiB、会话预算 32 MiB；超限**明说**给模型（「未送入模型」），不静默丢。清空点＝`Loop.Reset`（会话结束），因为任务结束时历史里仍有旧截图，早清会让下一轮 resolve 失败。

## 3. RED-before（两条，均为改动前盘上读数）

**R1 适配器：工具结果带图 = 整请求失败（三协议同时）**

```
$ go test ./internal/llm/openaichat ./internal/llm/anthropic ./internal/llm/openairesponses -run 23 -count=1
--- FAIL: TestToolResultImageReachesChatWire23
    buildWire refused a tool result carrying an image:
      internal: openai-chat: tool result part 1 is llm.ImagePart (only TextPart)
--- FAIL: TestToolResultImageReachesAnthropicWire23
      internal: anthropic: tool result part 1 is llm.ImagePart (only TextPart)
--- FAIL: TestToolResultImageReachesResponsesWire23
      internal: openai-responses: tool result part 1 is llm.ImagePart (only TextPart)
FAIL
```

**R2 agent：接缝类型根本不存在（编译期）**

```
$ go vet ./internal/agent
internal\agent\ticket23_image_seam_test.go:95:34: undefined: ImageAttachment
```

补上类型/仓库之后、装配 loop 前的中间红（证明「历史里没有 ImagePart」这一半也真实红过）：

```
$ go test ./internal/agent -run 23 -count=1
--- FAIL: TestScreenCaptureImageReachesMockLLM23
    fixture did not run: status=failed msg=internal: openai-chat: tool result part 1 is llm.ImagePart (only TextPart)
```

## 4. 改动

| 文件 | 改动 |
|---|---|
| `internal/agent/images.go`（新） | `ImageAttachment` 的落点仓库：`newImageStore/put/resolve/clear`、`maxImageBytes=4MiB`、`imageSessionBudget=32MiB`、`imageDropNotice`；无文件系统调用 |
| `internal/agent/tools.go` | `ToolOutcome.Images []ImageAttachment` + `ImageAttachment{MimeType,Data,Alt}`（无 Path/Ref 字段） |
| `internal/agent/loop.go` | `Loop.images`；`New` 建仓库；`buildRequest` 挂 `Request.ResolveBytes`；`executeCalls` 把工具图放进 `ToolResultPart.Content`（超限→一句「未送入模型」）；`Reset` 同时清仓库 |
| `internal/llm/openaichat/request.go` | `encodeMessage` 改多消息；工具结果带图 ⇒ tool 文本消息 + 紧随的 user 图像轮（data URL） |
| `internal/llm/anthropic/request.go` | `tool_result.content` 接受 image block（原生形状） |
| `internal/llm/openairesponses/request.go` | `function_call_output.output` 接受 `input_image`（原生形状） |
| `internal/agent/ticket23_image_seam_test.go`（接手时已存在、未跟踪） | 补 `internal/llm` import 使其可编译；5 条测试未改断言 |
| `internal/llm/{openaichat,anthropic,openairesponses}/tool_result_image_23_test.go`（新） | 各自线上形状 + 「句柄不上线」断言；openaichat 另钉「纯文本工具结果不得多长出 user 轮」 |

## 5. GREEN-after

```
$ go build ./internal/llm/... ./internal/agent/...        # rc=0
$ go test ./internal/agent -run 23 -count=1 -timeout 200s
ok   github.com/CarlosShao/wisp/internal/agent
    TestScreenCaptureImageReachesMockLLM23        PASS
    TestImageBytesRefNeverPersistedInJournal23    PASS
    TestImageBudgetRefusalIsStated23              PASS
    TestImageStoreIsMemoryOnly23                  PASS
    TestImageStoreRejectsOverBudget23             PASS
$ go test ./internal/llm/openaichat ./internal/llm/anthropic ./internal/llm/openairesponses -run 23 -count=1
ok  .../openaichat   TestToolResultImageReachesChatWire23        PASS
ok  .../anthropic    TestToolResultImageReachesAnthropicWire23   PASS
ok  .../openairesponses TestToolResultImageReachesResponsesWire23 PASS
```

AC#2 逐句对应：

- 「screenshot reaches mock LLM as ImagePart（C7 wire format asserted）」：`TestScreenCaptureImageReachesMockLLM23` 同时钉 C7 历史（`ToolResultPart` 里的 `llm.ImagePart{MimeType,BytesRef,Alt}`）与第二枚请求的线体（`data:image/png;base64,...` 的 `image_url` 部件）；用的是真 openai-chat 适配器 + golden SSE 回放，不是函数 mock。`ResolveBytes` 未接上时适配器会拒绝并让任务失败，所以绿不可能是假绿。
- 「bytes_ref absent from disk/db/logs」：同一测试的隐私三断言（slog 全记录、artifacts 目录全盘走、journal 的 tool_call 行逐列）；`TestImageBytesRefNeverPersistedInJournal23` 专钉库侧。
- 「DRM window → black frame + doc note」：DRM 黑帧在 `internal/tools/screen_capture_windows.go` 的后端（写面外），本腿不碰；接缝对 `drm=true` 的元数据不特殊对待，截图照常进模型（黑帧就是它该看到的画面）。

## 6. 门禁读数

```
go build ./internal/llm/... ./internal/agent/...            rc=0
go test  ./internal/llm/... -count=1 -timeout 200s          ok（llm, adaptertest, anthropic, golden, openaichat, openairesponses 全 ok）
go test  ./internal/agent/... -count=1 -timeout 200s        仅 Test31DisarmStopsTheWait 红（见下）
gofumpt -l <10 个改动 .go>                                   空（rc=0）
sh scripts/d22scan.sh                                        rc=0，clean（410 生产文件 + 871 internal Go 文件 emoji 扫）
go build ./internal/... ./cmd/...                            rc=0
```

`Test31DisarmStopsTheWait` 归因：**既有、与本节无关**。它单独重跑 5 次 = 3 红 2 绿（同一棵树、只改了图像接缝文件），从不构造 `Loop`；台账 `docs/reports/pending-and-issues.md:30794` 已把它记为「未在任何账册/票面出现，疑计时敏感/待核，⛔ 未证为回归」。本腿不改 `reminder_timer*`。

`go build ./...` 另有 `build/mut_rvw.go`（变异测试脚手架，未跟踪、非本腿文件）的 undefined 符号红，属既有；`./internal/... ./cmd/...` 全绿。

## 7. STOP-note：本腿写面外、AC#2 尚未闭环的地方

接缝的 **tools 半边**不在 `internal/agent/**` / `internal/llm/**` 里，本腿不能碰，必须另派一枚 `internal/tools` 腿（三处，都很小）：

1. `internal/tools/tool.go:39` `Result` 没有图像通道 ⇒ 加 `ImageBytes []byte` / `ImageMime string` / `ImageAlt string`（或一个 `Image *ImagePayload`），由 `renderSystemResult`（`internal/tools/system.go:246`）在 `len(res.Bytes)>0` 分支填充；
2. `internal/tools/bridge.go:630` 的 `agent.ToolOutcome{...}` 只搬 `Text/IsError/...`，需补 `Images: []agent.ImageAttachment{...}`（字段名见 `internal/agent/tools.go`）；
3. `internal/tools/system.go:261-263` 那段「C7 ImagePart 到 LLM 的接缝本仓当前不存在，所以模型看不到画面」在本次改动后**已不成立**，必须改写成「截图已作为图像送入模型；像素只在内存，未落盘」——否则票面/产品会对模型说假话（D22 诚实规则）。

在这三处补上之前：真实 `screen.capture` 仍只把元数据交给模型（`internal/tools` 里那枚工具没往 `ToolOutcome.Images` 填东西），但**接缝本身已就绪**：任何 `ToolProvider` 返回 `Images`，图像就会以 ImagePart 进历史、以协议原生形状上线，且句柄不落盘/库/日志。

AC#3 残余（input.type 真 fixture app + UIPI 真机 + 确认卡三选项 UI）需真机与前端，写面外，本腿不做。
