# 259residual-1 起手锚（第一笔）

- 腿名：`259residual-1`（突变验证腿）。写面＝仅两枚产码文件的突变行（跑完逐发还原）＋仅本目录只新建 `.md`。
- 取锚时刻：`2026-10-08 19:23:14 +0800`
- HEAD：`623a4d81dde0c5de337f8f46f243b36a958e462d`　Thu Oct 8 19:22:56 2026 +0800　A742 落账＋代修一处被顶位引用（comment-fix-land-1 收回；零翻框零 push）
- `git status --porcelain | head -20`（命令原样登记；未打开其中任何文件的内容）：
```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 M .scratch/wisp/probes/242/r3/logs/probe-routed.txt
 M .scratch/wisp/probes/268/v1/evidence.md
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
```
（以上为命令原样输出：与他腿共享工作树的在飞改动，本腿不碰、不读其内容。）

## 料文与两发落点（现量复核）

- 料文＝`.scratch/wisp/probes/next-instruments-brief-1/02-ticket259-ac2-brief.md`（备料腿 `next-instruments-brief-1`，commit `a6ec48d3`，其取锚 `5f52f310`）。
- 发 A：`internal/agent/approval/approval.go:565`＝`\tif !equalSecret(v, nonce) {`（`spend` 内成员扫描：`:564` `for v, stored := range s.values`、`:566` `continue`、`:567` 收口）。
- 发 B：`internal/agent/approval/queue.go:390`＝`\tif it.state != statePending {`（`allowScoped` 内 state 卫）。
- 漂移复核：`git diff 5f52f310..HEAD -- internal/agent/approval/approval.go internal/agent/approval/queue.go` 空；两文件 `git status --porcelain` 空 ⇒ **两落点与料文逐字一致、零漂**（今日 12 处注释改动未触这两枚文件）。
- 种前 hash（本程现取，与料文同）：
  - A `internal/agent/approval/approval.go` ＝ `67fb146898dc7b235623bb8ad2fac0e0b2465fe1`
  - B `internal/agent/approval/queue.go` ＝ `66fec7ae375344d7bc741270b42c4ec5050cf41f`

## 派单模板 PATH 现核更正（具名报回）

- 派单原文 `PATH="$PWD/../../third_party/sherpa-onnx:$PATH"`（自 `internal/agent/approval` 数两级）解析为 `internal/third_party/sherpa-onnx`，**该路径不存在**（现核：`ls internal/third_party` 报 No such file or directory）。
- 正确级数＝**三级** `$PWD/../../../third_party/sherpa-onnx`；DLL 现量在 `third_party/sherpa-onnx/`（`onnxruntime.dll`、`sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll`）。
- 本腿按三级 PATH 跑 `go test . -run '259' -v -count=1`（靶向名册：`ticket259_denial_rulers_test.go` 6 枚 ＋ `ticket259_panel_capability_rulers_test.go` 6 枚 ＝ 12 枚，与料文基线口径一致）。

## 四件套纪律（本腿自用留档）

1. 每发：种前 `git hash-object`；盘上改行 `sed -n` 复量；跑测试取 rc＋红句 file:line；`git show HEAD:<文件> > <文件>` 还原后 hash 逐字等＋scoped porcelain 空。
2. 突变形状（料文）：
   - A `sed -i '565s/if !equalSecret(v, nonce) {/if false {/' internal/agent/approval/approval.go`
   - B `sed -i '390s/if it.state != statePending {/if false {/' internal/agent/approval/queue.go`
