# 100 — AppContainer 落地票：**今天不做**，但把"什么时候必须回头做"写成一条**可判定的门**（来源=票 91 两会话实测）

**Status:** **RESERVED（不是 open，也不是 blocked-on-owner）**：AC#0 那条触发门**没满足之前，本票不许开工**，
且**不许用任何"先在主进程里做点什么"的折中绕过去**。
· **Blocks:** nothing · **Blocked by:** **AC#0（就是本票第一条判据，它是门不是任务）**
**Type:** OS 级隔离（owner 2026-09-21 那句"别连这种最基础的安全的东西都没有就搞笑了"的**唯一有牙的答复**）
**Packages:** 触发后才定。预估 `internal/proc/`（起进程处）+ `internal/winsec/`（授权与撤销）+ 新探针用例。
              **禁改**：`docs/PLAN.md`、`docs/specs/*.md`（AC#1/AC#6 若要动契约措辞走 **D22**，要 owner 批）、
              `internal/risk/assessor.go`、`internal/risk/pathresolver*.go`、`rules_gateway.go`、
              `tools/d22scan/**` 与 `allowlist.txt`。

## 为什么它是"待重评"而不是"已解决"或"以后再说"

票 91 两个会话的实测（`docs/evidence/s1/91-isolation-options.md` + `91-os-isolation-memo.md`，档位：**代理实测**，
编排者未逐格复现）：

- **受限令牌：不解决问题**（第一会话量的，第二会话横向复测**未推翻**）。四种降权令牌配方读
  `%USERPROFILE%` 靶文件与列目录**全部 OK**，只挡得住写。
  第二会话还加强了一刀（第一会话没测）：**`OpenProcess(宿主, PROCESS_VM_READ)` 在全部四种令牌下都 OK**
  ⇒ 只要秘密还在主进程内存里（今天工具的读、DPAPI 解密、SQLite 页、模型权重都在，
  `internal/tools/fs.go:144` 是主进程内 `os.Open`），降权子进程就是**同房间另一个进程**，
  **降权形同装饰**——这条比"能不能读目录"更致命。
- **独立低权限账户：不推荐**。`net user /add` = **error 5**；`LogonUserW` 对内置账户全部 **1326（密码不对）
  而不是 1314（privilege 不够）** ⇒ 挡路的是**凭据**不是 API；跨用户起进程**未证**；
  而且产品要保管一个能登录本机的口令 ⇒ **新增 D33 攻击面**（自己给自己造一个要保护的secret）。
- **AppContainer：唯一"能干活 + 真挡住"同时成立的路线**（读 / 写 / 宿主内存 / HKCU **全 DENIED**，
  显式授权的路径可读写），**非管理员就能建 profile**。代价也量到了：
  官方文档常量 `PROC_THREAD_ATTRIBUTE_SECURITY_CAPABILITY`（`0x00020000`）在本机
  `UpdateProcThreadAttribute` 报 **errno 24 `ERROR_BAD_LENGTH`**（`Length` 填 32/28 都试过）；
  **真正起得来的是被标废弃的 `0x00020009` + 无 `Length` 字段的载荷布局**
  （子进程自报 `isAppContainer=1`、完整性级别 `S-1-16-4096`、AC SID 与 `CreateAppContainerProfile` 一致）。
  ⚠ **"能用的恰好是被标废弃的那个"这条矛盾必须跟着进 AC#1，不许写成"按 MSDN 做"**——
  那是一句会害死实现者的话（第二会话原话）。

## AC（1:1；**AC#0 之前的一切都不许动**）

- [ ] **AC#0（触发门，先判这条；不满足 ⇒ 本票保持 RESERVED）**：
      `shell.exec` 或任一 D46 Tier-1 插件**在生产路径真的会起子进程**。
      **判据（可机器跑）**：`grep -rn "exec.Command" internal/tools/ internal/plugin/` 出现**非测试**命中
      **且被 bridge 调到**。
      今天的全仓非测试 `exec.Command` 只有 `mockllm` / SLO 自测 / `doctor` ⇒ **门未开，符合预期**。
      ⚠ **不许**用"先在主进程里做 ACL"来绕过这道门：那不是 AppContainer，也拿不到"同房间另一个进程"的防线。
- [ ] **AC#1 能力实测可复现**：把备忘录 §4.1 的**两条否定 + 一条肯定**做成 CI 用例
      （`0x00020000` → 必须证成 `ERROR_BAD_LENGTH` 或更好；`0x00020009` → 必须起得来且子进程自报
      `isAppContainer=1`）。**"按 MSDN"四个字不算判据。**
- [ ] **AC#2 读边界**：AC 子进程读 `%USERPROFILE%` 靶文件、读 `%APPDATA%`、`OpenProcess(宿主, VM_READ)`
      **三条都必须 DENIED**，探针输出用 `PROBE k=v` 那种**可被日志判定**的格式，不靠人眼。
- [ ] **AC#3 授权与撤销**：给定 `allowed_dirs`，子进程**只**能读写清单内路径；进程退出（含被 kill）后
      ACL **必须**回到快照态——用例必须包含 `TerminateProcess` 与主进程被 Job 连带清理
      （C30 `KILL_ON_JOB_CLOSE`）**两条崩溃路径**。
- [ ] **AC#4 双向变异**：① 去掉撤销 ⇒ 有用例红；② 把"授权失败"改成"当作没降权继续跑" ⇒ 有用例红；
      ③ 把 DENIED 断言改成"跳过不存在的路径" ⇒ 有用例红（防 §4.2 那类假绿）。
- [ ] **AC#5 体验面**：`%TEMP%`、`%LOCALAPPDATA%`、HKCU 三处各有**明确错误文案**用例，**禁止"看起来成功了"**。
- [ ] **AC#6 D32 无回归**：`Sleeping` 的句柄/进程数采样与 `PLAN.md:2253` 两档数字**一字不动**地跑；
      另测"按需拉起 + 用完退出"下 `Conversation` 峰值不越过同表工作态上限。
      ⚠ 一个**常驻**的降权 executor 会直接违反 `Sleeping…无子进程` ⇒ **形态只能是"按需拉起、用完退出"**。

## 已知未证（当结论用就是自欺；开工时先补）

1. **AC 子进程能否解密本用户的 DPAPI blob**：第二会话那套装置**自证作废**
   （`CryptProtectData` 产物不以 `DPAPI` 魔数开头、连控制组自解都报 errno 13）⇒ 只有推断
   （AC 读 `%APPDATA%` 已实测 DENIED，主密钥在 `%APPDATA%\Microsoft\Protect\<SID>\`）。
   要证就用第一会话 `wisp91-sess91-a7c3\spike2\` 那套**跑通过的**装置重跑。
2. 跨用户 `CreateProcessAsUserW`（路 3 的起进程步骤）：本机无第二账户令牌 ⇒ 未证。
3. `TOKEN_MANDATORY_POLICY_NO_READ_UP`（连读都挡掉的 Low IL 配方）：本机非提权，只拿到"需要特权未持有"的负向读数。
4. **Windows 10 / 其他 Win11 build 上 `0x00020009` 是否还起得来**：未证
   ⇒ **AC#1 必须是跑出来的门，不能是一次性读数**。

## Rules（本仓固定，开工时逐条生效）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/
`checkout .`（A34）；**不 push**；不在仓内建 worktree（A38④，快照目录带会话后缀）；
票面 append-only（改行前先读；标题前插段落要重抄标题，`git diff --numstat` 删除列必须 0）；
四种假绿逐条点名；**变异只在 `/tmp` 仓外快照里做**；收尾前必跑 `sh scripts/d22scan.sh`（A64②）；
数字不达标写 FAIL 附数字，不许调阈值、不许挑运气那次、多样本全报；15 次工具调用内交回第一枚 checkpoint。
⚠ 引用要自己打开核对（本仓已抓到 5 起引用腐坏：A30、A51⑥、票 87 的"票 41"、票 96 的"清单没同步"、
以及**本票：备忘录把落地票写成"建议编号 95"，而 95 已被"其余 `0o600` 落点"占用 ⇒ 实际归口就是本票 100**）。

## Progress log（append-only）

- 2026-09-21 17:0x（编排者）：建票并登记为 **RESERVED**（来源：票 91 两会话实测 + 裁定 **R21**）。
  我把它单独立案而不是塞进票 95，理由：**票 95 是"把没接的落点接上"（现在就能干），
  本票是"要不要引入一种新的进程形态"（要等门开）**——两件事混在一张票里，
  前者会被后者拖成"整票不动"，而后者会拿到前者"已经绿了"的读数当自己的证据。
  ⚠ **owner 的可选决策挂在 Q-30**（发布/售卖前要不要排一期）；不答**不影响任何在飞的事**，
  因为今天**没有可降权的对象**（工具的读全部发生在主进程内）。
  next= **不派单**。门（AC#0）一开——即 `shell.exec` 或 Tier-1 插件在生产路径起子进程——立刻派，
  且派单简报的第一句就是"先跑 AC#0 的 grep 证明门开着"。
