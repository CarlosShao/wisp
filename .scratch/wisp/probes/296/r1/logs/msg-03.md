296-r1 AC#2 落地（甲-全：热加载源两枚 return 都过 ApplyHotkeyDefaults）+ 同批改那句过期承诺

根因＝常驻腿热加载那一半的配置闭包把 [hotkey] 四格裸映射给桥，链路里没有任何 ApplyHotkeyDefaults，
而桥照规矩把空串读成「配置里关掉了」⇒ 机主那份 summon=''／mute=''／panel='' 的配置文件在开机 1 秒后
把三枚全局键拆成 live=0（回归引入点 5e8748b3，账在编排者名下）。

改的只有 cmd/wisp：hotkeyReloadSource296（＝hotReload258 的本体）两枚 return 全过合并 —— 裁定节点名的那枚
裸映射，以及「配置文件读不到」那支 return ball.HotkeyConfig{}（四格全空同样进桥、同样 live=0）。
只套前者会漏掉后者，形二的夹具已把这一支钉住。

同批把 resident_windows.go 里那句「The bridge keeps the current bindings on an empty answer」删掉并改成按行为写的实话：
桥**不**保持当前绑定 —— internal/ball/hotkey_reload.go 的 Check() 里唯一的跳过条件是 reflect.DeepEqual，
全文件零条 all-empty 守卫 ⇒ 全空一旦与已应用集不同就真去重绑。调用点那条「goes to the bridge as-is」一并更正。

零动：internal/ball（乙案未裁）、internal/config/schema.go 的四枚默认值、cmd/balldebug、HotkeyDisabled 与
none/off 一族（Q-82 默认动作＝不做 ⇒ 零 normalize）、机主的 config.toml（未读未写）、hotCfg258 的 Warn 时机。

目标用例改前红（rc=1，红句逐字含 live=0，进程读数与真机 :27/:28/:30/:31 同形）→ 改后绿（rc=0，
且名册里 hotkeys rebound after config change 0 命中＝DeepEqual 跳过、零次 Win32 重绑）。
件：.scratch/wisp/probes/296/r1/20-ac2-fix-and-comment.md ＋ logs/green-after-fix.md
AC#3 真机那一发归编排者跑；AC#4 门禁的整包改前／改后各两发在下一笔。
