# 票 272 对抗验收裁决件 — 腿 `272-v1`（非实现者）

> 本件是**对抗验收**，不是复述实现件 `272-r2`。实现者是另一枚腿 `272-r2`（已停），
> 那 52 行不是我写的。我的职责是攻它。任何一格我**不**因为它自陈"做了"就判成立。
> 结构：§0 起手锚｜§1 六格逐格判语｜§2 恒真性进攻｜§3 我推翻实现件的哪几句｜
> §4 我推翻编排者的哪几句｜§5 未做完的格｜§6 终态自证。
> ⛔ 本腿零产码改动、零 AC 框翻动。

---

## §0 起手锚（第 1 笔，任何长跑命令之前落盘）

原文读数（本腿 2026-10-07 自跑）：

```
$ date
Wed Oct  7 11:02:56 CST 2026

$ git rev-parse --short HEAD
e18e32da            # 分支 dev

$ git status --porcelain -- cmd internal scripts tools .github docs frontend
（空 —— 六族路径起手全干净，无别人的未提交 Go/脚本/文档改动）

$ git log --oneline -3 -- cmd/wisp/config_reload_223_test.go
8d30a862 验 5g2-v1（非实现者只读腿·第 1 笔）：... [顺手提走那 52 行的编队事故笔]
3d9b8374 probes(232-r2 收尾代提): ...
a16d1ff7 票 231 AC#3＋同批名册 · 常驻钉加一行，cause=newer-build 进两张互斥名册
```

被审对象落点核对（本腿现量，非抄实现件）：

```
$ git show --stat --oneline 8d30a862 -- cmd/wisp/config_reload_223_test.go
 cmd/wisp/config_reload_223_test.go | 52 ++++++++++++++++++++++++++++++++++++++
 1 file changed, 52 insertions(+)

$ git diff --numstat 8d30a862^ 8d30a862 -- cmd/wisp/config_reload_223_test.go
52      0       cmd/wisp/config_reload_223_test.go        # 只增 52／0 删，与票面/实现件一致

$ git diff --numstat 8d30a862 HEAD -- cmd/wisp/config_reload_223_test.go
（空 —— 8d30a862 之后测试文件再无改动，故 HEAD 版 == 8d30a862 版 == 被审版）

wc -l HEAD 版测试文件 = 797
wc -l 8d30a862^（改前尺）版 = 745      # 差 52，对得上
$ git show HEAD:cmd/wisp/config_reload.go | md5sum
5ce441ca5e72b64d18a6c26f1c066882       # 起手时产码基线 md5，§6 还原自证以此为拉
```

结论性锚点：**"改前尺" = `8d30a862^:cmd/wisp/config_reload_223_test.go`（745 行）**；
"被审尺" = `HEAD:cmd/wisp/config_reload_223_test.go`（797 行）。二者差恰为那 52 行。
