# 35-r7 — 起手锚（commit-first，先于任何 `go` 命令）

腿：`35-r7`，写面限定 `internal/panel/**` ＋ `.scratch/wisp/probes/35/r7/**`。
⛔ 不动 `cmd/wisp/**`（`255-r1` 的编译面）。本件落地时**尚未跑过任何 go 命令**。

## 1. HEAD 现量

- `git log -1 --format=%H` ＝ `d1dddb77dcc6a8415282e8ea3142b1e17064f3d0`（短号 `d1dddb77`），分支 `dev`。
- 锚点尺：`git merge-base --is-ancestor bfcb23e4 HEAD` → `ancestor_rc=0`（不早于 `bfcb23e4`，成立）。
- ⚠ 共享树在取锚期间自行前进过一枚（首读 `052b393f` → 交锚 `d1dddb77`），以下为交锚时刻复量。

## 2. `internal/panel/` blob 名册（`git ls-files -s internal/panel/`，34 枚）

```
100644 61f6db5e0f73e9bc9996f4d3fd7a9fae3ebf4149 0	internal/panel/approval.go
100644 070e3c84fadedc825c300fc781fd7d52d3908fe8 0	internal/panel/approval_test.go
100644 18458ff0c264d52f2fbaa6fead5e70a9be891cd1 0	internal/panel/assets.go
100644 6b0dd86df628e0288fd1dcc0866b2f3ea6601bca 0	internal/panel/assets_test.go
100644 745912faaf1805969713fed424dd32d421d2e727 0	internal/panel/attachments.go
100644 97d40c1cf07f56a662a8db430f0f2ae4a7f9a5a9 0	internal/panel/attachments_test.go
100644 bebe8e702a85c640641551e93dd09936530950f4 0	internal/panel/bridge.go
100644 582dc4d0113b2d206884b15b2baeddc58e122be2 0	internal/panel/bridge_test.go
100644 f0f638861416904086ec12cd00e83f699ec4e154 0	internal/panel/composer.go
100644 1757a487d672897a8e59d221fc613912a48790cf 0	internal/panel/composer_dispatch.go
100644 3e69cdf4fd00d086e34625d7b8011bad374d504e 0	internal/panel/composer_dispatch_test.go
100644 e8fdf8015806d027f3356b0c73ca641542a19895 0	internal/panel/composer_handlers.go
100644 74ad9969405f4708a24c9f4cdebe514cbc532bd6 0	internal/panel/composer_handlers_test.go
100644 e247b38d87d892254cc2e8c9bc22934f1d34501c 0	internal/panel/composer_test.go
100644 19445546ea1267b15101018d17d464f91a265472 0	internal/panel/config_handlers.go
100644 2a0c7c22c0c986204647587859b1168b47f4314e 0	internal/panel/config_route_248_test.go
100644 b934c270ca5243b856bf1c18e5f7f1367c8b200a 0	internal/panel/doc.go
100644 31e55fdfd5a07913a76f53b036750a6a320f4b23 0	internal/panel/frontend_hygiene_test.go
100644 a537c5e797fd65de79fe5d4418a6b649e980bd2f 0	internal/panel/git.go
100644 96cf2ec3410319fa8a66d1b5a36b8e4771977484 0	internal/panel/git_test.go
100644 774c765c973177cf9c8e05c6d496ad540c10efbf 0	internal/panel/inbound_roster_253_test.go
100644 9563e640dc8d79506fea2750490191a59c6fb3fa 0	internal/panel/instructions_200.go
100644 87060c5276ec1cd481fc3a4f028b79a83f692b17 0	internal/panel/instructions_200_test.go
100644 6668f19bc14f3391fe9d5169186403db8cf82087 0	internal/panel/l2_grant_boundary_test.go
100644 645e1f02bbf97dd35cdc1ff2563c9c212a6f9a18 0	internal/panel/pump.go
100644 3f84c3360063bb2ae26dd6cc5ccb819ac510e7a4 0	internal/panel/pump_test.go
100644 a6b552e5985641024821a939be0209f87f178f6a 0	internal/panel/subagent_blocked_220_test.go
100644 341878c86d63477499c3a704dd5eb662038bf329 0	internal/panel/subagent_roster_197.go
100644 f803179df6465eaa7e5ed4933bc1210283023b1d 0	internal/panel/subagent_roster_197_test.go
100644 0d124985470738763a7313ed95517c4c9d265306 0	internal/panel/subagent_stream_197_test.go
100644 b60acc6145ac2cee222bd3d77758e8087fa585f0 0	internal/panel/tokens_fourway_test.go
100644 4db510ce5594e24e31c2f0ef275fff0e6f8bf4bc 0	internal/panel/workspace.go
100644 f9e4a60dd512509a299db68bd4ec434a290e7bca 0	internal/panel/workspace_account_181r3_test.go
100644 8bef67e00bddb1ee8106562cea0ec5492a3b1adc 0	internal/panel/workspace_test.go
```

## 3. 脏面现量（别人在飞的件，⛔ 一个字节都不许动）

`git status --porcelain -- cmd internal docs .github scripts` ＝ **4 行**：

```
 M cmd/wisp/config_readers_255.go
 M cmd/wisp/panel_host_windows.go
 M cmd/wisp/panel_resident_windows.go
?? cmd/wisp/panel_reshow_255r1_windows_test.go
```

全部属 `255-r1`（在飞写面）。`internal/` 范围 0 行＝干净。
（另：仓内其它路径 `.gitignore`、`.scratch/wisp/probes/**`、`design/**` 等有别的腿/程挂脏，与本腿无关，本腿一枚没 stage。）

## 4. 起手读数（只读，非 go）

- `git grep -n 'hits248'`：命中全部在 `cmd/wisp/panel_config_248_test.go`（`:52` 定义、`:187-192`/`:235`/`:253`/`:256` 使用）。
  ⇒ **具名冲突预告（顶回候选）**：`hits248`/`secretShape248` 是 `cmd/wisp` 包的**包内测试符号**，
  Go 不允许 import 测试文件，`internal/panel` 的测试**够不到它**；而派单又限定本腿写面 ⛔ `cmd/wisp/**`。
  处置见 `impl.md`（镜像钉源方案），冲突具名顶回，不静默改写派单。
- `git grep -ni 'redact' -- internal/panel`：仅 `composer.go:260-261`、`config_handlers.go:140-141` 两处**注释**引用
  `internal/secret` 的终端 redactor；包内**没有**任何入向原文 redact 的产码支。
- `internal/panel/config_route_248_test.go:349`（`const canary = "canary-248-not-a-real-key-0f2a"`，
  `:375` `strings.Contains` 判）＝票面 :51 残差①在**包内**最近的一格既有覆盖，行号为本腿现量。
