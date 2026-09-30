# 245 — 球一进常驻进程，**裸 `Esc` 就被当成全局热键抢走了**：仓里自己的测试注释写着"这会吞掉整台桌面上别的程序的 Esc"，而生产默认值恰恰就是裸 Esc

- Status: **已立，未派，排在队列最前**（09-30 16:3x，编排者立；触发＝实现腿 `228-r1` 把球装进了常驻那条腿，**这件副作用从今天起在正常跑的产品进程里生效**，不再只在调试旁支程序里）。
- 来源：`228-r1` 交件清单第⑥节第 4 条（它具名上交的事实："默认 `cancel=Esc` 是**全局** `RegisterHotKey`，既有设计、但**今天第一次跑在常驻进程里**"）＋**我自己现读复认**（下面每行都是我自己跑过的，不是转抄）。
- ⚠ **这不是"本机被入侵"，照三行读**：① **现象在哪**＝**你的键盘**：只要 `wisp` 在跑，**Esc 按下去可能别的程序收不到**（记事本关窗、资源管理器改名、游戏/编辑器退全屏、对话框取消都靠 Esc）；② **有没有本机被入侵的证据**＝**没有**；③ **最坏后果是什么形状**＝**你在用别的程序时 Esc 失灵**，以及反过来"**想否决 Wisp 时按 Esc 不管用**"（两个方向都是真代价）。⇒ 这是**产品行为缺陷**，不是安全攻破。

## 现量（09-30 16:3x 编排者自己跑，别信行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| **生产默认把 cancel 定成裸 Esc** | `internal/ball/hotkey_windows.go:62` 逐字 `return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", `**`Cancel: "Esc"`**`, Panel: "Ctrl+Alt+P"}` | `sed -n '60,64p'` 我现读 |
| **开窗那一刻就把四枚全绑上**（含那枚裸 Esc） | `internal/ball/ball_windows.go:250` `b.hotkeyReport = registerAll(b.hwnd, b.opts.Hotkeys)`，注释逐字「Tray + hotkeys live on the same window」 | `sed -n '240,252p'` 我现读 |
| ⛔ **本仓自己写着这件事会吞键** | `internal/ball/hotkey_live_test.go:33-36` 逐字「the bare "Esc" cancel default: **registering Esc as a GLOBAL hotkey swallows Esc from every other app on the desktop for the whole test run**, so the live tests bind a modifier combination instead (**see the ticket's found-defect note**)」⇒ **真机测试刻意躲开这枚默认值**（它们用 `Ctrl+Alt+V`），**而生产路径用的是它** | `sed -n '28,45p'` 我现读 |
| **设计上 Esc 只在"等人批"那几秒被借走** | `hotkey_windows.go:434-436` takeEsc 注释逐字「binds VK_ESCAPE as the cancel hotkey (**B1 Confirming takeover**)」＋`:448-451` releaseEsc「hands Esc back (**B1**): …the Esc key **must be handed back when the session ends**」 | `sed -n '430,470p'` 我现读 |
| ⇒ **两枚机制互相矛盾** | 稳态默认已经是 Esc ⇒ "进 Confirming 才抢、离开就还"这套 `takeEsc`／`releaseEsc` **今天等于永久租期**：从开机那一刻就抢，还无可还 | 上面两行合读（这一列是我自己的判断，非现读） |
| 常驻腿**没读 `[hotkey]` 配置** | `228-r1` 第⑥节：`[hotkey]`／`[ball]` 未接 ⇒ 常驻那条腿里热键＝**编译默认**（用户改 `config.toml` 无效）；`ApplyHotkeyDefaults`（`hotkey_windows.go:76-90`）还会把空值填回默认，逐字理由「a live hotkey the user can see, never a dead silent one」 | 〔腿报＋我自己复读了那枚函数体〕 |
| 卡片文案与状态表都写着 Esc | `internal/agent/approval/approval.go:90` `ChannelEsc: "按 Esc 键"`；`docs/PLAN.md:3082`（**D43 转移表＝C12 冻结那张**）第 22 行逐字「否决（**单击球 / `Esc` / KWS 否决词 / 面板拒绝**，B1）」 | `sed` 我现读 |

## 裁定（我自己定方向，理由摆在盘上）

**走乙：稳态不绑 cancel 那一枚，只在进入 `Confirming` 时 `takeEsc`、离开时 `releaseEsc`。** 三条理由：
1. **不动冻结文字**：`PLAN.md:3082` 那行属于 **D43 转移表＝C12 冻结**，把默认键名改掉（甲形）＝**契约变更、要他一句批准**；乙形保留"`Esc` 是否决键"这句话**一字不动**，只把它的作用区间改回它本来被写成的样子（借—还）。
2. **复用现成机制**：`takeEsc`／`releaseEsc` 已经在包里、已被票 64 的真机用例钉过（`interaction_live_test.go:207/:223`），乙形是**接上**它，不是新造一套。
3. **丙形（让用户在 `[hotkey]` 里关掉）与既有设计相反**：`ApplyHotkeyDefaults` 逐字写着"空＝重新启用默认，方向要是安全的那一边：看得见的热键，绝不是一枚悄悄死掉的键"⇒ 靠"留空"关不掉，硬关要改那函数＝改另一处定案。

**⚠ 乙形留下的那半代价，具名不藏**：**确实在"等人批"那几秒里，别的程序的 Esc 仍会被抢走**（窗口只有 2–3 秒；`SPEC-06` 的 L1 定义）。这半条**今天不解决**，登记在本票末尾当残余。

## 验收判据（逐格要 `file:line` 与正控）

- [ ] **AC#1 稳态名册里不含裸 Esc**：常驻腿起来之后、**没有卡挂着**的时候，实测注册集只能是三枚（`summon`／`mute`／`panel`），**`cancel` 不得在场**。尺＝`HotkeyReport().Bindings()`／`Live()` 现读，交回时带逐字读数（⚠ 别再拿"4/4 live"当好消息）。
- [ ] **AC#2 `Confirming` 内借、离开还**：造一发真机——卡片挂起时按 Esc 真的否决（复用票 64 那一族），**卡片撤下后再按 Esc，别的程序要能收到**（判据形状＝"还"这一步有具名读数，不是只测"借"）。
- [ ] **AC#3 ⛔ 撞钉预检的三枚钉要同批改期望，一枚都不许偷偷绕**（我 09-30 16:3x 现读出来的射程）：① `internal/ball/live_windows_test.go:69` 逐字红句「**the live ball did not register its four hotkeys**」＝它把"四枚都注册"钉成了期望；② `cmd/wisp/resident_ball_228_test.go` 与 `resident_windows.go` 里那两处**新写的用户可见文案**含「hotkeys live 4/4」；③ `internal/ball/hotkey_live_test.go` 的 `liveHotkeys()` 注释（它现在写着"生产默认是裸 Esc、我们躲开"，本票落地后**这句会变成假话**）。⇒ 三处一律改成**带条件的事实句**＋把①那枚钉**收紧**（不是放宽：从"恰好四枚"改成"稳态三枚＋借用时四枚"）。⛔ **不许靠删断言或改期望值蒙过去**。
- [ ] **AC#4 用户要能改**：`[hotkey] cancel` 真被常驻腿读进去（这条依赖票 228 的"常驻腿接 `config.toml`"那一格，**本票不许自己造第二条配置通路**——若 228 那格还没做，本格判不了就写成"阻塞于票 228"，不许假绿）。
- [ ] **AC#5 残余登记**：乙形下"借 Esc 那 2–3 秒会吞别处的 Esc"这半条**今天不修**，按 `SPEC-12 §5` 五字段登记成推迟项（完成判据＋残缺表现），⛔ 不许只在票里写一句"以后再说"。

## 禁区

- ⛔ **不改 `PLAN.md` 与 `docs/specs/**` 一字**：D43 转移表（C12 冻结）里"`Esc` 是否决键"这句话本票**原样保留**；若某位认为要走甲形（换默认键名），那**必须先落一枚人工批准的 `A##`**，本票不许顺手换。
- ⛔ 不动 `internal/observe/thresholds.go`／golden／`allowlist.txt`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；不许为变绿放宽任何断言。
- `frontend/**`／`design/**` 零读零写零转述。
- ⛔ **不新增词面型仪器**（不许造"扫注释里热键名"的门）；判据一律问能力（注册集里有没有那一枚），不问文案。
- git：只 commit 不 push；commit 必带显式 pathspec；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件。

## 排程与串行

- ⛔ **排在队列最前**（先于 `197-r3`／`224-r3`／`242`）：理由是**它影响的是他这台机器现在的键盘**，不是仓里的账。
- 动 `internal/ball`＋`cmd/wisp` ⇒ ⛔ 与票 228 的后续片串行；⚠ **跑真机热键用例前 `tasklist //FI "IMAGENAME eq balldebug.exe"` 与 `wisp.exe` 都必须为 0**（本仓实测：自己的调试进程占着全局热键会造出 **8 枚假红**）。
- ⚠ 与票 244（GUI 子系统）同属"这一轮之后才看得见的外观/行为"，两枚**互不阻塞**，但**别同一批改**：244 改构建链、245 改热键稳态，读数面不同，混在一枚腿里出事无法归因。
