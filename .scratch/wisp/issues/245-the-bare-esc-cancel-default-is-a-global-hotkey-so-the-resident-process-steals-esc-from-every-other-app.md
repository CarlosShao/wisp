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

- [x] **AC#1 稳态名册里不含裸 Esc**：常驻腿起来之后、**没有卡挂着**的时候，实测注册集只能是三枚（`summon`／`mute`／`panel`），**`cancel` 不得在场**。尺＝`HotkeyReport().Bindings()`／`Live()` 现读，交回时带逐字读数（⚠ 别再拿"4/4 live"当好消息）。
- [x] **AC#2 `Confirming` 内借、离开还**：造一发真机——卡片挂起时按 Esc 真的否决（复用票 64 那一族），**卡片撤下后再按 Esc，别的程序要能收到**（判据形状＝"还"这一步有具名读数，不是只测"借"）。
- [x] **AC#3 ⛔ 撞钉预检的三枚钉要同批改期望，一枚都不许偷偷绕**（09-30 **17:0x 我再跑一遍把②的指认更正过**——我先前写「两处文案含 hotkeys live 4/4」是**词面指认写歪了**，盘上逐字没有那串，真实形状见下）：
  ① `internal/ball/live_windows_test.go:68-70` 现读逐字 `if got := b.HotkeyReport(); !got.AllLive() || len(got.Live()) != 4 {` ＋红句「the live ball did not register its four hotkeys」＝它把"恰好四枚在场"钉成了期望 ⇒ **本票落地后这一发必然红**，要改成**稳态三枚＋借用时四枚**（⚠ **收紧**不是放宽：判据要能同时抓"该借的时候没借上"和"稳态里多绑了一枚"）。
  ② `cmd/wisp/resident_ball_windows.go:111` 现读逐字 `hotkeys live %d/4: %s`（**枚数是 `len(rep.Live())` 现算的、不是写死的字面量**）＋`:114` 的 slog 字段 `"hotkeys_live", len(rep.Live())` ⇒ 落地后这行**自然变成 3/4**，不用改文案，但**必须确认它变的值与判据一致**；另有 `cmd/wisp/resident_ball_live_228_windows_test.go:144/:151` 那两把尺：`:144` 只取 `verdict` 里 "hotkeys live" 那段当日志、**不断枚数**，`:151` 只在 **"hotkeys live 0/4"** 时红 ⇒ **乙形（3/4）不会打红它**，但它 `:148-150` 那段注释写着"A ball whose **four** hot keys all failed to register" ⇒ **这句注释在本票之后会变成半假话**，一并改成带条件事实句。
  ③ `internal/ball/hotkey_live_test.go:33-36` 的 `liveHotkeys()` 注释（现逐字：生产默认是裸 Esc、真机测试刻意改绑 `Ctrl+Alt+V` 以免吞掉整台桌面的 Esc）⇒ **本票落地后"生产默认稳态也绑裸 Esc"这半句会变假**，要改成"稳态不绑、Confirming 期间借"的事实句。⚠ ⛔ **不许顺手把那几枚用例改回绑裸 Esc**——它们在稳态里躲开裸 Esc 仍然是对的（借期只有几秒，用例不该把整段跑批挂在借期内）。
  ⇒ 三处一律改成**带条件的事实句**；⛔ **不许靠删断言或改期望值蒙过去**；⛔ 不许新增词面型仪器。
- [ ] **AC#4 用户要能改**：`[hotkey] cancel` 真被常驻腿读进去（这条依赖票 228 的"常驻腿接 `config.toml`"那一格，**本票不许自己造第二条配置通路**——若 228 那格还没做，本格判不了就写成"阻塞于票 228"，不许假绿）。
- [x] **AC#5 残余登记**：乙形下"借 Esc 那 2–3 秒会吞别处的 Esc"这半条**今天不修**，按 `SPEC-12 §5` 五字段登记成推迟项（完成判据＋残缺表现），⛔ 不许只在票里写一句"以后再说"。

## 编排者收件 — 09-30 **17:41**（账 `A478`）：245-r1 交件已核过盘上三把尺；⛔ **一枚框都不翻**——那份表出自实现方自己，裁决要等另一程

**我现跑复认的（不采信通知正文）**：四枚 commit 逐枚 `git log -1` 都在（`f91feddd` 骨架 17:06 → `e6d6426f` 产码 17:25 → `08886b8d` 表填满 17:38 → `60f888c8` §6 名册 17:39）；表「docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r1.md」＝**114 行／26,518 字节、占位符 0 枚**；写面已空（`git status --porcelain cmd/wisp internal/ball`＝**0 枚**）。
**我自己重跑的那把能力尺**（不是背它的读数）：`grep -rn "standby|Confirming" internal/ball/ball_windows.go` ⇒ 稳态那一支真的存在且注释具名指向本票（`:243` 逐字「an idle ball may hold and leaves the cancel slot standby (ticket 245)」、`:808`「set is the idle set, and the idle set does not include cancel」、`:868`「this is the ONLY path that registers the cancel id」）。

- [ ] **AC#6（新补，我自己抓到的两笔）"四枚"这句用户可见文案与稳态枚数不一致 ＋ 判据锚在"枚数"而不是"裸键"**：
      ① `cmd/wisp/main.go:24`（`wisp -h` 打给用户的正文，反引号原始串）现读逐字仍写「its four global hot keys in this same process」，而 245-r1 的实测稳态读数＝**3/4** ⇒ **本票落地把这句变成了半假话**（它上一轮刚被 228-r1 改过，那一轮改的是"票 07"，枚数没人管）。**完成判据**＝改成带条件的事实句（说清"稳态三枚、进入 Confirming 借第四枚"），凭据要现跑那发 `winlive` 读数。
      ② **J#2 由我裁（这是我重切自己的判据，不是盘上现成的答案）**：本格字面"注册集只能是三枚"写窄了——真缺陷的形状是「**没有任何修饰键的裸键被全局抢走**」，不是"稳态必须恰好三枚"。⇒ 用户把 `[hotkey] cancel` 配成**带修饰键**的组合时**应当常绑**（那不会吞别的程序，且与 `ApplyHotkeyDefaults` 逐字写着的方向一致：「a live hotkey the user can see, never a dead silent one」）；配成**裸键**时仍走借用。⚠ **今天没有行为差别**（`[hotkey]` 还没接进常驻腿＝AC#4 阻塞于票 228，默认档只有一枚裸键）⇒ 本格是**判据精度与钉子形状**的账：AC#3① 那枚"稳态三枚"的钉要改成**判裸键**（判据＝稳态注册集里不存在"无修饰键"的那一枚），否则等 AC#4 通了，用户配的组合键会被这枚钉打红、或被实现方为了过钉而悄悄不绑。**完成判据**＝钉子改成按"有无修饰键"判＋正控两形各一发（稳态绑裸 Esc 必红／稳态绑用户配的组合键必不红）。
- ⚠ **SPEC-12 §5 那一行我不搬**：`docs/specs/**` 是**一字不改**的禁地（对腿、对我都成立；那张表最后一次被动是 09-19，早于编队开工）。五字段登记落在**本票末尾＋台账 A478**，⛔ 未加 `DEFERRED(D-xx)` 代码标记（加了就会打破"标记 ↔ 表 1:1 双向"那条硬规矩——245-r1 这一步判断正确，记它一笔）。**欠账具名**：要真进那张总表，需要 owner 一句「搬」。
      ✅ **09-30 17:5x 他回了「搬呗，我都不知道你说的啥玩意，看都看不懂」⇒ 已搬**（上面那句"我不搬"**作废、原句不抹**，账 `A479`）。落点＝`docs/specs/SPEC-12-roadmap-governance.md` §5 表内**新增一枚 DEFERRED 行**（五字段齐全，逐字把"要连那 2–3 秒都不吞别处＝换键名＝触碰 D43/C12 冻结、须人工批准"写进"前置"栏）；尺＝`git diff --numstat` **1 增 0 删**，其余一字未动。⛔ 仍**未加** `DEFERRED(D-xx)` 代码标记（现况表 13 行／代码 0 枚标记，加一枚会造出"只有这行有标记"的第二形状）。**撤销口令「撤 245 那行」**。⚠ 归我自己：这类账目归位**本就不该摆给他做选择题**，下次只做＋告知。

**框为什么一枚都不翻**：`docs/evidence/s1/` 里的裁决表按 `AGENTS.md §1.5`／`SPEC-12 §4.3` **必须出自非实现者**，而这份出自 `245-r1` 自己（它同时是写手）。所以它现在只是〔实现方自述＋台件〕，不是判语。**下一程＝`245-v1`（非实现者对抗验收）**，交回后我据它的表定勾或退回。
**它的自述里我最想让验收腿攻的几处**：① 新探针 `escBorrowProbe`（`hotkey_live_test.go:57`）拿**另一枚 id** 向 Win32 注册裸 Esc 来判"有没有人占"——这条**是不是真能证到"别的应用此刻收得到 Esc"**（它自己列了 J#3：证不到）；② `SKIP-LOUD 0 枚`那发的起跑名册与终态名册差集是否只落在自己落点；③ M1/M2a/M2b 三发突变的还原是否逐枚 `git diff --quiet` 复认；④ 它揪出的**两枚票面未点名**的同类计数钉（`hotkey_live_test.go:74/:101`、`interaction_live_test.go:234`）改得对不对、有没有把别人的票 64 判据顺手放宽。

## 编排者收表 — 09-30 **18:23**（账 `A480`）：非实现者裁决表 `245-v1` 交回 ⇒ **AC#1／AC#2／AC#3／AC#5 翻勾（AC#2 附两条具名条件）**，AC#4／AC#6 不翻，另补三格

裁决表：`docs/evidence/s1/245-esc-not-a-standby-global-hotkey-v1.md`＝**256 行／74,002 字节、占位符 0 枚**；四枚 commit 逐枚 `git log -1` 复认（「2f2d398e 骨架 17:48 → 6f5595fd 判语 18:15 → ec2dcd74 攻不动 18:19 → 7a3d5a3c 名册 18:21」）；`git diff --quiet -- internal cmd` ⇒ **rc=0**（三发突变全还原）；票面六枚框它一枚未动（我复跑 `grep -c`）。
**我自己重跑的两把尺（不采信它的转述）**：① `grep -rn "TakeEscForCancel|ReleaseEscAfterSession" --include=*.go internal cmd | grep -v _test` ⇒ **产品进程里零调用者，唯一调用者是旁支程序 `cmd/balldebug/main.go:635`**；② `grep -n "Confirming|approval" cmd/wisp/resident_ball_windows.go` ⇒ 该文件 `:22`／`:69-70`／`:163` 逐字自陈"这条腿没有审批门、没有任务管线、卡片进不来、手势没有执行者"。**⇒ 这两把尺撑起了 AC#2 的条件②，也是本票之后最该干的那件后端活（转立票 246）。**

- [x] **AC#1 成立**（凭据＝它自己那发出厂真机逐字「hotkeys live 3/4: … cancel=Esc not bound while idle」＋默认键名一字未动：`hotkey_windows.go:69` 仍是 `Cancel: "Esc"`，且被反向钉住 ⇒ **未触碰 C12／D43**）。
- [x] **AC#2 附两条具名条件成立**（⚠ 不是"放行"）：**条件①**＝票面字面那句"别的程序要能收到 Esc"在仓里**零判据**，它另造了第二枚真 GUI 进程补测（四发读数），**那套台件不进仓就守不住回归** ⇒ 落成下面 AC#9；**条件②**＝**出厂常驻进程既没有卡片、也就永远走不到"借"那一支**（上面尺①）⇒ 本格只到"机制对了"，**不许读成"Esc 否决已经能用"**（那是票 246 的活）。
- [x] **AC#3 成立**：三枚钉逐枚字面对撞判"收紧非放宽"（`!=4`→`!=3` 且新增双向 `requireIdleRoster`；`0/4` 那一支一字未动、**另加一条让 `4/4` 判红**；注释改准而用例仍躲开裸键＝票面禁令守住）。**"收紧后仍能红"它实跑了**：V-M1 之下红在 `TestBallLiveLifecycle`，红句逐字「the live ball did not register its three idle hotkeys (summon/mute/panel, cancel is borrowed only during Confirming)」。**外加一枚它新抓的缺陷＝下面 AC#8。**
- [ ] **AC#4 仍不翻**（它判"阻塞于票 228 如实、合格、不许翻勾"：`resident_ball_windows.go:86` 逐字 `Hotkeys: ball.DefaultHotkeys(),`）。
- [x] **AC#5 成立**，且比我落笔时更全：`SPEC-12 §5` 那一行已在（`docs/specs/SPEC-12-roadmap-governance.md:80`，就是我 17:50 那枚「dc694430」搬进去的），代码里 `DEFERRED(D-` 标记**仍 0 枚** ⇒ 双向核对未被我这一步打破（它替我复认了这件事）。
- [ ] **AC#6 两半都不翻**：①`main.go:25` 那句「its four global hot keys」**不成立**——它跑了这一版二进制 `-h`，打出的就是这句，而同一枚二进制稳态实测 3/4；②它**独立裁了我那句自裁**：判"编排者把锚改到'裸键'方向成立"，但**四枚钉全还锚在枚数**（`hotkey_live_test.go:435`／`live_windows_test.go:74`／`hotkey_status_test.go:181`/`:245`）⇒ **顺序必须"先改钉、再改行为"**，且它实测到"用户配成组合键也不绑"这一形**已经在读配置那条腿上真发生**（E5 逐字「cancel hotkey left unbound while idle … binding=Ctrl+Alt+V」＋`live=3`）⇒ 我那句"今天零行为差别"**说轻了，就地更正**：差别在另一条腿上今天就能量到。

- [ ] **AC#7（新补，它的 V-F1）借被拒那一支的记账陈旧**：`takeEsc` 失败时不设 `escTakenOver`，于是 `ReleaseEscAfterSession` 首行早退 ⇒ 报告永久停在 `HotkeyError`、`AllLive()` 恒 false、`Problems()` 一直说"没注册"，**而桌面其实什么都没绑**（`ball_windows.go:878-899`）。完成判据＝失败那一支也要把状态归位＋一枚正控（借被拒之后报告必须说"稳态三枚、cancel 未绑"，不是"注册失败"）。
- [ ] **AC#8（新补，它的 V-F2）别的程序已经占了裸 Esc 时，正确的实现也过不了**：它起真外来进程占住 Esc（E6）⇒ 四枚真机用例**全硬红**（逐字「ticket 245 RED: bare Esc is claimed on this desktop while the ball is idle」），`SKIP-LOUD` 0 枚 ⇒ 在这种机器上**测不绿**，而写腿表里"被拒＝具名绿"那一支结构性走不到。完成判据＝把"桌面已有占位者"与"我们自己没绑上"这两种红分开判（分开两条红句），⛔ 不许改成跳过。
- [ ] **AC#9（新补，它的条件①）把那套"第二个进程真收到 Esc"的台件搬进仓里当常驻尺**：现在它只在验收侧（E1 正常 `keydown_esc=1`／E2 自占裸 Esc `0`＋`wm_hotkey=1`／E3 出厂进程挂着仍 `1`／E4 反回去的 `0`）。完成判据＝仓内一条能跑的用例做同样的四形＋**不依赖人工桌面**。⚠ 改文案时注意它抓到的另一枚钉：`leg_dispatch_gate_133_test.go:156` 的 `usageBare133` 钉住 `  wisp   ` 行首形状 ⇒ AC#6① 那处改动必须保住它。

**它的门禁（我不转述成自己的）**：build rc=0；vet 四档全 rc=0（含 `-tags winlive` 两档，比写腿多两档）；`scripts/d22scan.sh` 正控先吃读数（`PASS=34 FAIL=0 SKIP=0`）后真扫 clean；`internal/ball` 默认层唯一红＝`TestC21TableColourRowsMatchTokensCSS`（起手即在、末次改动属票 74，非本票）。⚠⛔ **一枚挂着的不定性读数**：它自己说 `cmd/wisp` 默认层读到 **1 红 2 绿**，但首发用例名丢了（只留 `tail -20`）⇒ **欠一次归因**，我不替它编原因，挂在这一行，下一枚动 `cmd/wisp` 的腿**必须先把这枚红是谁归出来**再动手。
**它顺手发现、不属于本票的（我一枚没改，全部分了归口）**：① `ball_windows.go:809-812` 那句"没有调用者同时轮询配置与卡片"对 `cmd/balldebug` **不成立**（同进程两件事都在）⇒ 归票 228 后续片，注释先改准；② 票 64 已勾框里那句「断 `len(live)==4`」今天码里是 `==3`（实质未变、引文过期）⇒ **归我补记**，见 `A480`；③ `cmd/wisp` 测试输出里出现 `goroutine outside the D38 roster` × `subagent-finish-*`（`internal/tools/subagent_197.go:380-381`）⇒ **归票 197 名下**，排进 197-r3 的射程。
**本票状态**：勾 4／未勾 5（AC#4＋AC#6..AC#9）⇒ **不加 `-done`**。

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

## 残余登记（乙形那半条代价，`SPEC-12 §5` 五字段；AC#5 的落点）

> 由产码腿 `245-r1` 于 09-30 17:3x 追加。同一份内容逐字在裁决支持表 `docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r1.md` §4.1。
> ⚠ **`docs/specs/SPEC-12-roadmap-governance.md` §5 那张表本身是产码腿禁地**（"不改 `docs/specs/**` 一字"），所以这一行**要由编排者搬进那张表**；本腿也**没有**在代码里加 `DEFERRED(D-xx)` 标记，因为 `issues/README` 硬约束要求标记与登记表 **1:1 双向**对得上——造一枚没人登记的标记＝当场把那道核对打破。

| 字段 | 内容 |
|---|---|
| **类型／项** | DEFERRED：`Confirming` 借 `Esc` 那 2–3 秒仍吞别处的 `Esc`（票 245 乙形残余） |
| **为什么现在不做（依据）** | 本票「裁定」三支理由走的都是"稳态不绑"，没裁"借期也不抢"；要把借期也消掉得换新机制（低级键盘钩子／只在面板内取消），属新契约面。借期长度不是自由参数：`internal/statemachine/timeouts.go:48` 现读 `{state: StateConfirming}: {3 * time.Second, EvConfirmExpired}`＝L1 窗口上限 3 秒（`SPEC-06` 的 2–3 秒） |
| **完成判据** | 卡片挂着期间按 `Esc`：Wisp 的否决生效（D43 #22）**且**同一发里焦点在别的应用时那个应用也收到 `Esc`（现读形状的判据，不是"我们不再记账了"）；若结论是"做不到不抢"，就要给出**不接管全局 `Esc`** 的第二条否决通路并过 `Confirming` 那族真机用例 |
| **当前残缺表现** | 等批的那 2–3 秒内别的应用收不到 `Esc`（记事本关窗、资源管理器改名、对话框取消、游戏/编辑器退全屏那一段会失灵）；反向同样存在：想否决 Wisp 时若焦点在别处，那发 `Esc` 归 Wisp |
| **前置／落哪个模块·切片／谁做／触发条件** | 落 `internal/ball`（`takeEscWith`／`releaseEscWith` 这一对）＋ `internal/agent/approval` 的卡片生命周期；切片＝票 228 后续片把卡片接进常驻腿之后（S6 之后的球／门控联动面）；谁做＝**编排者立后续票**（产码腿不许自己开票），形状需 owner 定（钩子 vs 只在面板取消）；触发条件＝常驻腿出现真卡片，或 owner 实机撞到"那两三秒按 `Esc` 没反应"任一条 |

### 本腿另交的四条具名残余（不是推迟项，是归属清楚的账，详见表 §0.1／§4.2）

1. **AC#4＝阻塞于票 228**：常驻腿整棵 `config.toml` 没接（`cmd/wisp/resident_ball_windows.go:86` 逐字 `Hotkeys:  ball.DefaultHotkeys(),`），本腿未新增第二条配置通路。
2. `cmd/wisp/main.go:25-26` 那句 usage `its four global hot keys` 在本票之后成半假话 ⇒ 归票 228 的 AC#7/AC#8 文案面（台账 `A477` 已把 `main.go` 那族句子判归 228），本腿不跨票改。
3. `RebindHotkeys` 会丢掉一次在飞的借用（`internal/ball/ball_windows.go:822`，既有形状，本腿只补注释未改行为）⇒ 票 228 后续片。
4. **待人拍板（本腿按票面字面取了严的一支）**：用户把 `[hotkey] cancel` 配成带修饰键的组合时，稳态**也不绑**它——票面 AC#1 写的是无条件"只能是三枚"。另一种读法（只有裸键才不绑）会让 AC#3① 那枚钉的"稳态三枚"不成立。⇒ 见表 §5 J#2，需编排者／owner 一句定案。

## Progress log

- [2026-09-30T17:37:19+08:00] agent=245-r1 claimed=票 245（产码腿，非裁决腿） did=乙形落地：稳态不绑 cancel（`HotkeyStandby`）＋`takeEscWith`/`releaseEscWith` 接成唯一借用路径；三枚 AC#3 钉同批收紧（另加两枚票面未点名的同类计数钉）；真机探针 `escBorrowProbe` 证借还两头；M1/M2a/M2b 三发突变逐发红过并逐发还原成绿；本表末尾追加五字段残余登记与这节 log evidence=docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r1.md（commits f91feddd／e6d6426f／本枚） skipped=none（本票无前端面） next=交非实现者对抗验收；AC 框一枚未勾（5 未勾／0 已勾，勾归编排者）；Status 那一行按编排者原文保留未改，⛔ 未新增 `DEFERRED(D-xx)` 代码标记（理由见上）
