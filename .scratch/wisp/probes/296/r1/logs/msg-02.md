296-r1 AC#1 夹具先行（两形，改产码前红；产码面＝只提取、零行为变化）

新增 cmd/wisp/resident_hotkey_296_windows_test.go 两枚在程用例，跑真装配
startResidentBall(reg, nil, hotCfg, hotReload)，两半是两枚不同闭包（⛔ 非 258 那种 src, src）：
- 形一＝机主 config.toml 的 [hotkey] 四行逐字（summon=''／mute=''／cancel='Esc'／panel=''）
- 形二＝热加载那一半读到的目录里没有 config.toml（裁定节新增射程：return ball.HotkeyConfig{} 那一支）
两形各自带前提自证（形一读回必须只有 cancel 一枚；形二必须走 LoadFile 报错那一支），
断言 ⓐ 桥第一次 Check() 之后 Live() 仍含 summon/mute/panel（live=3）ⓑ ConfiguredHotkeys() 是合并后的值。

产码面只把 hotReload258 的本体提为包级 hotkeyReloadSource296(dataDir)，闭包与测试共用同一枚本体
（同仓先例＝票 258 把构造那一半提为 residentBallHotkeyChain258）；
行为零变化，GOFLAGS= go build ./... rc=0。不提取就只能在测试里重打一遍裸映射，那枚仪器对产码修没修不敏感，
给不出票面 AC#1 要的「改前红」——已在 00-anchor.md §3 具名上报。

改前红（尺＝PATH=sherpa-onnx:build go test ./cmd/wisp/ -count=1 -run Test296BridgeKeeps -v，rc=1，
起手 wisp.exe／balldebug.exe 各 0 枚）：
形一 live=0 [] 并逐字复现真机那组 hotkey disabled (unset in config) + live=0 读数；形二 live=0 []（连 cancel 都 disabled）。
件：.scratch/wisp/probes/296/r1/10-ac1-red-before-fix.md（红句逐字）＋ logs/red-before-fix.md（原始 stdout）
下一笔＝甲-全：两枚 return 都过 ApplyHotkeyDefaults ＋同批改 resident_windows.go 那句过期承诺。
