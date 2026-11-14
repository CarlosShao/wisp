# 票 23 AC#3 — 提权目标拒绝：**修复后端到端真机读数**

> 执行：编排者。时间 `2026-11-14 10:57–10:59 +08`。基线：`331ecf2a9`（修 `realProcessElevated` 的句柄语义）之上。
> 本件补的是**修复腿自己声明欠着的那一笔**：对真实 High 窗调用 `input.type` 必须返回显式 UIPI 错误。

## 1. 修复前的读数（腿 23-AC#3 原始，`ad50ee24e`）

```
elevated_target      : taskmgr.exe（IL 0x3000 = 12288），self 0x2000（8192）
elevated_is_error    : false
elevated_text        : 已向 taskmgr.exe 注入 15 个字符        ← 该拒绝，没拒绝
```
根因：`internal/tools/system_backend_windows.go:741` 把 `OpenProcess` 的**进程句柄**当**令牌句柄**喂 `GetTokenInformation` ⇒ `handle is invalid` ⇒ 返 0 ⇒ `0 > self` 恒假。

## 2. 修复后的读数（本件，`331ecf2a9` 之上）

```
probe=probe23 pid=65920  injector_integrity=0x2000  roster_entries=11
windows_above_our_integrity=2
elevated_candidate_0=hwnd=0x1706be pid=6904 name=taskmgr.exe il=0x3000 title="Task Manager"
elevated_candidate_1=hwnd=0x10346 pid=12780 name=cua-driver.exe il=0x3000 title="Cua.AgentCursorOverlay.default"
elevated_target=hwnd=0x1706be pid=6904 name=taskmgr.exe il=0x3000 title="Task Manager"
elevated_is_error=true
elevated_text=目标程序以管理员身份运行，Wisp 无法向其输入（Windows UIPI 平台硬限制：普通权限进程不能向更高完整性级别的窗口发送输入，本工具不会尝试提权绕过）。请把目标应用以普通权限运行后重试
elevated_fixture_chars_unchanged=true
```

**同一个目标窗（`il=0x3000`）：修复前 `IsError=false`＋“已注入 15 个字符”，修复后 `IsError=true`＋冻结拒绝句＋fixture 零变化。**
产品还额外量到：`window.focus` 对该提权目标返回**无 host 错误**，但 `elevated_foreground_after_focus` 仍是提权窗自身 —— 与 `AttachThreadInput` 被拒一致，这条**如实留着不夸大成“聚焦成功”**。

## 3. 我改过探针一处（编排者行为，留痕）

`probe23` 原本在「抢不到前台」时 `os.Exit(3)`，于是**提权那一半根本没跑**。实测本机桌面被一个 High 窗（Taskmgr）占前台，Medium 进程抢不回来 ⇒ 探针每轮都 abort，读数空白。

我加了一条分叉（`.scratch/wisp/probes/23/probe23/main.go`，`PROBE23-UIPI-SPLIT` 注释）：**前台抢不到时先把提权形跑完再退出**。理由是这两半前提不同——
- 300 字那一半**需要**前台；
- 提权那一半**不需要**：拒绝发生在完整性门、发生在打第一个字之前，所以它没有前台也测得到。

这不是放宽判据，是把「本来可测却被前提顺带掐掉」的读数捞回来。`ABORTED` 记录**保留**，没有伪装成通过。

## 4. 仍未满足的那一格（不糊）

票 23 AC#3 还含**第三半**：确认卡三选项、且「所有应用」带 danger 样式。本件**没有**量它；此外 300 字那半的**修复后复跑**也还欠（被 Taskmgr 占前台挡住，需机主让出前台）。故 **AC#3 框仍未翻**。

## 5. 复跑

```powershell
go build -o .scratch\wisp\probes\23\logs\probe23-fixed.exe ./.scratch/wisp/probes/23/probe23
.\.scratch\wisp\probes\23\logs\probe23-fixed.exe `
  .\.scratch\wisp\probes\23\logs\wisp23fixture.exe .scratch\wisp\probes\23\logs
# 提权形在 summary 里，键为 elevated_*
```
