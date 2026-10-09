票 294 AC#1..3 落点判定 + Progress log（腿 294-a1，只读、零 Go 命令）

AC#1 走查落点：新建 untagged resident_mute_294_test.go 照 258 模板但需加 SelectorExpr 一层（attachMuteGate 是方法调用），锚形状非行号（resident_windows.go 正被 296-r1 改）。三件事静态可满足。
AC#2 真窗落点：新建 resident_mute_294_windows_test.go 复用 258 真窗形状；IsLive 方法存在但 hkMute 非导出、cmd/wisp 无镜像常量 → 缺 instrument=补 hkMute 本地镜像常量 + 一次现读；SKIP-LOUD 本机非退路。
AC#3 executed 断言名册：唯一 executed 断言 resident_mute_290:197 只断无门⇒false，无一人断 happy 的 executed==true；toggleMute :172/:174 两 true；happy 世界（真门）已在 290 夹具；M2 反形翻 :172/:174 撞票 295 同写面须串行、验后复原。
我没跑的尺归 294-r1/296-r1（走查红绿/真窗 IsLive 实读/executed==true/门禁四数）。
票面不符：:303 生产挂载现漂到 :336（两把 grep 间被并发编辑），:217/:278 亦过期，改锚形状。

零源码改动、零翻框、未动 docs。本把 go build/vet/test 一把没跑。
