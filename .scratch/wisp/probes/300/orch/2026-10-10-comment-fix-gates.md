# 票 300｜编排者代笔：三处"注释／文案里的假话或机器专属路径"修完，门禁七把全落 `rc=`

跑＝编排者本人（⛔ 腿）。`date_start=2026-10-10 14:15:32 +0800`／`date_end=14:15:58`（尺＝`date` stdout），
第二批发件段 `14:2x`。动的那枚文件＝`internal/audio/parse_wave_format_300_windows_test.go`（**⛔ 断言、⛔ 夹具字节、⛔ 偏移、⛔ 生产码**——三处全是注释与 `why:` 里的说明串）。

## ① 改了哪三处（每处给"来路＋为什么"）

1. **头顶注释那句"没有一枚单独是决定性的"**（原文逐字：`No single face here is the decisive one; the four together only accept 24.`）
   ⇒ 前半句⛔ 成立：按节实测 **S2 与 S3 各自只绿 `offset=24`**（在 20／22／26／28 全红），所以单形即已钉死偏移；
   后半句（四枚一起只放行 24）⛔ 变。
   来路＝非实现者验收腿 `300-v2`（commit `21c9abe3`）第 6 条顶回＋我自己按节复算（见 `2026-10-10-orch-five-offset-scan.txt` 第 112–116 行 tally）。
   ⇒ 新句逐字写明"S2／S3 只绿 24＝单形有牙；S4 绿 22／24／28、S1 绿 20／24＝这两枚单形钉⛔ 住"。
2. **`why:` 里的权威指针 `mmreg.h:2474` ⇒ `:2475`**（来路＝`300-v2` 必答第 2 条自报"腿的指路偏 1"，我自己 `sed -n '2470,2482p'` 现量对上：
   `:2474` 那行是 `DEFINE_WAVEFORMATEX_GUID(WAVE_FORMAT_PCM)`， carrying `Data1=0x00000001` 字面量的是 `:2475` 那行 `DEFINE_GUIDSTRUCT("00000001-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_PCM);`）。
   ⚠ 同批核过⛔ 改的三枚：`:2376`／`:2540`／`:2550`／`:2483`（`300-v2` 与我两把都判✅）。
3. **注释里那枚绝对 SDK 路径**（`C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h`）⇒ 换成 `shared/mmreg.h` ＋ 一句"绝对路径按机器与 Kit 版本变，故意⛔ 写死在这里"。
   理由＝机器专属＋版本专属，且产码注释自己用的就是相对名（`wasapi_windows.go` 那三行布局注释）；⛔ 把某台开发机的目录结构喂给下一个读代码的人。

## ② 门禁七把（每把自己落一行 `rc=`，⛔ 0 字节＝那格没交）

```
gofumpt -l internal/audio/parse_wave_format_300_windows_test.go   名册＝空        rc-gofumpt=0
gofmt   -l 同上                                                    名册＝空        rc-gofmt=0
改动区非 ASCII 自取（注释里不许出现票面符号 ⛔★⚠）                 命中 0 行        rc-ascii=0
go vet ./internal/audio/                                                          rc-vet=0
go test ./internal/audio/ -run TestParseWaveFormatSubFormatOffset300 -count=1 -v   rc-test=0
  尺＝grep -cE '^--- PASS|^=== RUN' 的**合计数**＝6（⛔ 一把尺：其中 `--- PASS` 顶层 1 枚＋`=== RUN` 5 枚＝1 父＋4 子）
  第二发（⛔ -v）：ok github.com/CarlosShao/wisp/internal/audio 0.026s
sh scripts/d22scan.sh                                                              rc-d22scan=0
GOFLAGS= go build ./...                                                            rc-build=0
```

⚠ 本机 PATH ⛔ `gofumpt`（第一发我因此只取了 `gofmt`）⇒ 第二发改用 `$(go env GOPATH)/bin/gofumpt` 才拿到上面那行；
**定式＝仪器⛔ 在 PATH 时先定位再判"这一把没跑"**，⛔ 把"命令找⛔ 到"写成"检查通过"。

## ③ 越界自查

- 写面＝`internal/audio/parse_wave_format_300_windows_test.go`（注释＋一条说明串）＋本件＋票 300＋台账＋停车点＋那份按节件（追加更正一节）。
- ⛔ 碰 `internal/audio/wasapi_windows.go`（生产那枚偏移 24⛔ 动）；⛔ 动任何断言／夹具字节；⛔ 动 `level.go`／`liquid.go`／`thresholds.go`／golden／`allowlist.txt`／D43／C1–C32／D1–D47；⛔ 动三枚冻结件；⛔ 动 `.gitattributes` 或任何 git 配置；⛔ 动 `frontend/**`／`design/**`。
- ⛔ 零 push（机主从未授权）；commit 带显式 pathspec。

## ④ 顺带记下的两枚交付形状缺陷（⛔ 阻断 `AC#2`，但必须具名）

- `probes/300/r1/logs/post-commit-recheck.txt`（817 B，mtime 13:40:44）**从未入库**，而 `40` 件 §1／§2 两处写着"逐字＝该件"
  ⇒ 一枚**blob 层拿不到的悬空引用**（读数本身在 `.md` 里有，所以那格⛔ 空）。⇒ 并回＝**文中指名的件一律当待验断言**，
  与"注释里的测试名一律当待验断言"（`A` 账第五形）同级；尺＝`git ls-files --error-unmatch <路径>`（本次回 `Did you forget to 'git add'?`）。
- `probes/300/r1/30-gates.md:116` 残留未填的 `<!--ROSTER-->` 占位符（尺＝`grep -rn ROSTER probes/300/r1/`）。
  ⚠ `300-v2` 自己在 `2328934e` 里也抓到并更正了它**同一形**的脑补号——两枚都记在盘上，⛔ 抹。
