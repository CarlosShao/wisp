296-r1 补第三发改后读数（隔离跑 hotkey 家族：Test296BridgeKeeps* + Test258*）

尺＝PATH=sherpa-onnx:build go test ./cmd/wisp/ -count=1 -run 'Test296BridgeKeeps|Test258' -v
起手 tasklist wisp.exe／balldebug.exe 均 0 枚，date＝Fri Oct 9 21:38:38 CST 2026，锚＝HEAD 2746ae95（产码面＝fc4aedaf，此后 cmd/wisp 零改动）。

rc=0；--- FAIL 顶层 0 枚、--- SKIP 顶层 0 枚；^--- PASS 顶层＝13 枚、--- PASS 含子＝13 枚（两把尺同数＝这一族无子测试）。
⇒ 本票两枚夹具绿、票 258 那三枚桥用例（含 :153 那句"源可以返回全空"的正控与 AC#2/AC#3 两枚）一字未放宽、同批仍绿。

另具名一条读数是算出来的还是量出来的：cmd/wisp 自锚点 7285ef83 起只有本腿两笔 commit
（尺＝git log --oneline 7285ef83..HEAD -- cmd/wisp/ ＝ fc4aedaf ＋ 28a2ff4e 两枚）⇒ 前面那四发整包名册没被别的腿的
cmd/wisp 改动污染，这条是被证成的、不是假设的。
