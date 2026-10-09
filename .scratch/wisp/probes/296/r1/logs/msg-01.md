296-r1 起手锚（只读复核 + 落点计划，⛔ 零产码改动）

票 296 落地腿开工第一笔：date + git log --oneline -1 + git status --porcelain -- cmd/wisp/
逐字落盘（cmd/wisp 开工前干净；wisp.exe / balldebug.exe 各 0 枚）。

现跑复核裁定节记的四处行号（锚 = 字面量）：
- resident_windows.go:213 = hotReload258 读不到配置文件那支的 return ball.HotkeyConfig{}（相符）
- resident_windows.go:215 = hotReload258 裸映射那支（相符；票面原句 :214 已由裁定节更正）
- resident_windows.go:208 = 「The bridge keeps the current bindings on an empty answer」那句过期承诺（相符）
- internal/ball/hotkey_reload.go:93 = Check() 里唯一的跳过条件 reflect.DeepEqual；全文件零条 all-empty 守卫（相符）

落点 = 甲-全（两枚 return 都过 ApplyHotkeyDefaults）+ 同批改那句注释，⛔ 不动 internal/ball。

具名上报一条实现约束（不自选形）：测试拿不到 runResident 的局部闭包本体 ⇒「重打一遍同形」=
对产码修没修都不敏感的仪器，满足不了 AC#1 的「改前必红」。本腿把 hotReload258 的本体提为
cmd/wisp/resident_windows.go 里的包级函数（同文件、同射程），产码与测试共用它 —— 先落「提取但未套
合并」这一笔让夹具红，再落「套上」那一笔转绿。

件：.scratch/wisp/probes/296/r1/00-anchor.md
