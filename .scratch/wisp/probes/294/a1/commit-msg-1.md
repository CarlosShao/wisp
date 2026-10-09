票 294 AC#0 只读普查（腿 294-a1，只读、零 Go 命令）

静音"生产形状"名册：
- startResidentBall 测试调用点 10 枚全为真窗、手搓 &residentBall{} 13 枚；attachMuteGate 3 枚测试直调全挂在手搓球/nilBall 上，无一枚真窗后再挂门。
- attachMuteGate 锚调用形状计数：产码 1（resident_windows.go，现随 296-r1 并发从 :325 漂到 :336）、测试直调 3、定义 1、注释 1、AST 走查 0。
- 装配根"直调有走查零"欠账名册除静音外至少 4 枚：bindBallHost、RegisterShutdownHook(StepCancelTasks)、startResidentTaskSource 调用点、startResidentAudio。
- AC#1 模板=resident_hotkey_258_test.go（Ident 形状，attachMuteGate 是 SelectorExpr 方法调用，需加一层不能纯照抄）。
- AC#2 IsLive 方法存在（internal/ball 已对 hkMute 读），但 hkMute 在 cmd/wisp 非导出、无本地镜像常量，cmd/wisp 从未对 mute 读 IsLive。

零源码改动、零翻框、未动 docs。go build/vet/test 一把没跑（296-r1 独占 Go 编译面）。
