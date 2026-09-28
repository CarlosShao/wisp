# 183-v2 — 非实现者对抗验收（收窄版）：宿主指针豁免的三发变异恒真性 / 安全退让面进攻 / 三条边界 / AC#2 措辞判读

- 验收代理：`183-v2`（角色＝非实现者对抗验收；**只裁不改**：产码零改动，AC 框一枚不勾）
- 派单：`.scratch/wisp/dispatches/2026-09-28-132x-accept-183-v2-narrowed-incremental-commit-after-the-previous-leg-died.md`
- 被验收对象：`183-r1` 交件 `docs/evidence/s1/183-pointer-exemption-r1.md` ＋ 产码 commit `0662a35a`
- 本表是**增量交付**：第 1 节＝step-0 与起手现场；§2＝恒真性矩阵；§3＝安全退让面进攻；§4＝三条边界；§5＝AC 判读；§6＝本程没裁什么；§7＝现场与护栏自证；§8＝next。
- ⚠ 本节以下标 `待填` 的格，是"写表时还没跑"，**不是**"判定通过"。

## 1. step-0 五件（本程现量）

```
$ date "+%Y-%m-%d %H:%M:%S %z"
2026-09-28 13:23:30 +0800

$ git rev-parse --abbrev-ref HEAD ; git rev-parse --short HEAD
dev
7ec24d04

$ git merge-base --is-ancestor d61aee1b HEAD
（真＝派单起手锚 d61aee1b 是当前 HEAD 7ec24d04 的祖先；HEAD 多出的那一枚正是派单 183-v2 自己）

$ git status --porcelain -- internal/ cmd/
（空）

# 三向 md5（逐行管道到 md5sum，与本程起手一致）
$ sed -n '468,473p' internal/risk/provenance.go | md5sum
858e45116383caa3e7c1dd4b0924fad1 *-      ← 派单预期值，一致
$ sed -n '11,15p'   internal/risk/taintmatch.go | md5sum
5680ddd18e2d2ec2a85e485b54f4c12e *-      ← 一致
$ sed -n '482,506p' internal/risk/provenance.go | md5sum
f89e891e5eee3f3ea2b4f89d921072c4 *-      ← 一致

$ git show --numstat --format="%h" 0662a35a
```

`0662a35a` 的 numstat（本程现量，删除列**逐枚全 0**）：名下产码件 = `internal/risk/provenance.go` 48/0、`internal/risk/taintmatch.go` 59/0、`internal/risk/pointer_183_test.go` 331/0；其余为 `docs/evidence/s1/183-pointer-exemption-r1.md` 178/0、`probes/183/r1/**`（mut-f1/f2/freq/head 副本＋logs＋overlay＋zz183r1_e2e_test.go）、票 183 面 1/0。没有一枚落在 `frontend/**`／`design/**`。

判据件顶层 `func Test*` 枚数（本程现量，尺＝`grep -n "^func Test" internal/risk/pointer_183_test.go`）：**7 枚**，派单与被验收表写"顶层 6 枚" ⇒ 数目差一枚，§2 逐名列出这 7 枚。

## 2. 恒真性矩阵（三发变异 M1/M2/M3 × 逐名判据）

待填（本程自取变异，还原只用 `git cat-file blob HEAD:<path> > <path>`）。

## 3. 安全退让面的进攻（自造样本）＋ 频率读数口径复量

待填。

## 4. 三条边界

待填：(a) 参数侧按值放行有没有换张脸回来；(b) `declaredPath` 跨 mark 可见性；(c) `taintmatch.go:11-15` "逐 token 追踪已被否决"算不算被复活。

## 5. AC 判读（只给判读，不勾框）

待填：票 183 AC#2 档位；票 177 AC#3 条件格／票 175 AC#5／票 176 AC#3-5 可否翻；票 185 独立还是并进 183。

## 6. 本程没裁什么（逐枚具名＋为什么）

待填。

## 7. 现场与护栏自证

待填（含被拒/没成功的调用、有没有跑过删除命令、工具调用枚数 vs 硬顶 40、伪授权两栏）。

## 8. next=

待填。
