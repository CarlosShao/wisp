# 23 写腿（POSIX/vet 收口）—— `platformSystemBackend` 的 POSIX 面 ＋ 两处 `unsafe.Pointer` 真修

- 腿：`23-posix-vet`（写腿；只动 `internal/tools/**` ＋ 一支 `_test.go` ＋ 本篇）
- 锚：HEAD `232b0846c`
- 时刻：2026-11-05 10:5x–11:2x +08
- 面：`internal/tools/system_ops_other.go`（新）· `internal/tools/system_backend_windows.go` · `internal/tools/system_bridge.go`（注释 1 行）· `internal/tools/system_23_clipboard_memory_windows_test.go`（新）
- ⛔ 改 `Status:`、⛔ 改台账、⛔ 动 `system_23_test.go`/AC#2 那批件、⛔ push

## 0 这两件是谁点名的（不是我挑的）

票 23 自己是 WIP：`119c3fa91` 落了后端六文件，把 `internal/tools` 的 POSIX 编译打断；两处指针伪造是同一批代码的 vet 红。点名处：

- `docs/evidence/s1/107-verdict.md` §2.2「POSIX —— HEAD 字面口径：跑不起来」：`GOOS=linux go test -c ./internal/tools/` rc=1，唯一未定义符号 `platformSystemBackend`（`system_bridge.go:129`）；
- 同上 §2.3 的补充口径就是"HEAD + 一份 7 行 POSIX 桩"——本腿把那枚桩落进仓（桩的语义逐字相同，只是补上了"为什么"）；
- 同上 §4「归属与建议」第 1 条直接写了文件名与函数体；
- 同上 §AC#5 的 vet 红：`system_backend_windows.go:337:44` 与 `:355:32` `possible misuse of unsafe.Pointer`（本腿现量复现，见 §3）。

## 1 改了什么

| 文件 | 改动 | 是什么 |
|---|---|---|
| `internal/tools/system_ops_other.go` | 新，`//go:build !windows` | `platformSystemBackend() (systemFuncs, error)` ＝ `return nil, ErrNoSystemBackend`。补的是**唯一**一枚符号：其余文件都是平台中立的 |
| `internal/tools/system_backend_windows.go` | +2 proc、+1 helper、两个剪贴板腿重写 | `procGlobalSize` / `procRtlMoveMemory` ＋ `moveMemory`；`realClipboardRead` / `realClipboardWrite` 不再把 OS 地址做成 Go 指针 |
| `internal/tools/system_bridge.go` | 注释 1 行 | `:37` 原写 `system_backend_other.go`，与 `:110` 自己写的 `system_ops_other.go` 冲突（前者在 HEAD 上不存在）；按 `:110`（也是票 107 复判引用的那一行）改正 |
| `internal/tools/system_23_clipboard_memory_windows_test.go` | 新，`//go:build windows` | 两支：GMEM 块「Go→OS→Go」双向真搬字节；零长度块＝空串而非错误 |

`system_ops_other.go` 的语义选择（写进文件头注释）：**不接 xdotool / AppleScript 这类 shim**。POSIX 上"音量已设为 40%"这种读数没人能手验，硬凑等于把 D34 风险表变成小说；`nil + ErrNoSystemBackend` 让每一个动词落到家族已有的 fail-closed 路径（`IsError` 文本、任务继续，D37），与零值 `SystemBridge` 逐字同答。

## 2 指针那两处：为什么这样修（而不是压诊断 / 加 build tag）

vet 的 `unsafeptr` 判据是"`uintptr` 实参转 `unsafe.Pointer` 且不是 ptr→uintptr 往返/算术"（`hasBasicType` 走 `Underlying()`，所以**给 uintptr 起个新类型名也拦**）。`GlobalLock` 返回的是地址，它不是 Go 堆上的对象，GC 也追不到它——所以这不是"工具太严"，而是"这里确实没有可证的安全形状"。

修法＝**不造 Go 指针**，让 OS 搬字节：

- 读：`GlobalSize(h)` 拿块长 → `make([]uint16, size/2)` → `RtlMoveMemory(&buf[0], p, size)`；
- 写：`RtlMoveMemory(p, &utf16[0], size)`。

两侧都是 `uintptr(unsafe.Pointer(...))` 的安全方向（就是 `syscall.Syscall` 自己的形状），OS 那一侧全程 `uintptr`。附带两枚真收益：

1. 读回的字节在 `GlobalUnlock` 之后仍然活着（旧写法是解锁后继续指 OS 内存）；
2. **拷不动就报错**：`RtlMoveMemory` / `GlobalSize` 取不到时返回错误，而不是"静默返回空串假装剪贴板是空的"——空串是读数，假空串是谎。

## 3 前后读数（逐条命令 ＋ rc，同一台机、同一棵工作树）

改前用 `git show HEAD:<path>` 把两枚被改文件临时还原、把两枚新文件移出包，量完立刻还原（还原后 sha256 与备份逐枚相同）：

| 闸门 | 改前 | 改后 |
|---|---|---|
| `go vet ./internal/tools` | **rc=1**：`system_backend_windows.go:337:44` / `:355:32 possible misuse of unsafe.Pointer` | **rc=0** |
| `GOOS=linux go build ./internal/tools` | **rc=1**：`system_bridge.go:129:9: undefined: platformSystemBackend` | **rc=0** |
| `GOOS=linux CGO_ENABLED=0 go vet ./internal/tools/` | rc=1（同上未定义符号；107-verdict §2.2 逐字） | **rc=0** |
| `GOOS=linux CGO_ENABLED=0 go test -c -o /tmp/tools_linux.test ./internal/tools/` | rc=1（同上） | **rc=0**（产出 39,069,508 字节） |
| `go build ./internal/tools` | rc=0 | **rc=0** |
| `go test ./internal/tools -run '23\|System' -count=1 -timeout 200s` | ok | **ok**（含新增两支） |
| `gofumpt -l internal/tools/` | 空 | **空**（0 行） |
| `sh scripts/d22scan.sh` | rc=0 clean | **rc=0 clean**（`internal/tools=48` 生产文件，`bans#1-5 internal/=356 cmd/=56`） |

## 4 POSIX 不是"能编"而已——容器里真跑

`GOOS=linux` 只证明可编。按票 107 AC#3 的口径（POSIX 侧必须真执行）在容器里跑了一次（镜像本地已有，⛔ 不拉新镜像；仓 `:ro` 挂载，`-w /src/internal/tools` 与 `go test` 的 cwd 一致）：

```
docker run --rm -v D:/tmp/23posix:/t:ro -v "<repo>:/src:ro" -w /src/internal/tools \
  golang:1.27 /t/tools_linux.test -test.run "23|System" -test.v
```

读数：**RUN=28 / PASS=28 / FAIL=0 / SKIP=0，容器 rc=0**（日志 `.scratch/23posix/posix-23-system.log`）。`Test23WindowsZeroBridgeReachesABackend` 按构建标签不出现在 POSIX 侧（正确：它断言的正是"本平台有后端"）。

第一次跑把 `-w` 留在 `/src`，两支 236 的用例报 `open task.go: no such file` / `../agent/testdata/golden/... 读不回来`——那是**cwd 不对**（`go test` 平时把 cwd 设在包目录），不是代码红；`-w /src/internal/tools` 后全绿。

## 5 如实登记（没做到 / 不在本腿）

1. **票 23 的残余 AC 一件没闭**：AC#2 的 DRM 黑帧/`bytes_ref` 内存 only 已在别腿、AC#3 的真 fixture app ＋ UIPI 真机 ＋ 确认卡三选项 UI（前端）、AC#4 的真剪贴板 round-trip 仍是残余。本腿只清"挡别人门"的两枚红。
2. **全包 `go test ./internal/tools` 仍 rc=1**，唯一红是 `Test229SpawnNamesAreRosterShaped`（扫 `cmd/**` 的 goroutine 名册普查，与 `internal/tools` 无关；HEAD 提交 `232b0846c`/`111bfdb04` 的信息也记着"Test229红=票29既有非本改动"）。本腿**没有**为它放水。
3. **树里有一枚不属于本腿的未提交件**：`internal/tools/ticket23_image_channel_test.go`（mtime 10:31，AC#2 的 tools 半边验收尺）。本腿没改它、没提交它，但它参与本腿的门禁编译（因此本腿的 `-run '23|System'` 读数里含它那五支，且是绿的）。
4. **POSIX 真跑只在 `-run '23|System'` 这个子集上**：`internal/tools` 全量的 POSIX 跑法与读数不在本腿（票 107 的裁决者要做的那一格）。
5. 探针与日志：`.scratch/23posix/`（改前改后备份、`posix-23-system.log`、`d22scan-final.log`）。仓库内**没有**建 worktree/checkout，没有 `git add -A`，没有 push。
