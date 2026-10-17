#!/usr/bin/env python3
# 330-r2: append EXACTLY ONE bullet at the end of ticket 330. The existing bytes
# are read, hash-checked against the pre-append baseline and written back
# unchanged; only the new bullet is added (no AC box touched).
import hashlib
import io
import sys

P = r"D:\work\workspace\projects plans\Wisp\.scratch\wisp\issues\330-doorbell-instrument-mirror-roster-never-reconciled.md"
PRE = "e5f60b34038acebf458b23ba95f99a5e50805a565a6ac395ee3d50250b362cb3"
BULLET = (
    "- **330-r2 交件（2026-10-17 09:3x +08·只追加这一条、⛔ 翻任何 AC 框、⛔ 改上面一字）**："
    "G2 那格的\"缺世界\"**已造出**——新建 `internal/agent/plan_present_crowded_world_330r2_test.go`（272 行，⛔ 产码字节）"
    "把真镜像条目（组合 `planPresentDirectory165r3()`、⛔ 第二份抄本）与**一枚**不驻留的长尾候选放进同一目录，"
    "检索名额按 `internal/agent/inject.go:232-243` 的 300 token 预算削：长尾那条索引行占 291、镜像那条占 20，"
    "⇒ 镜像条目一旦从驻留支拨到检索支就抢不过名额、整枚不上线（档位这时才决定谁被注入）。"
    "★成对两发现跑（只打仓外导出树 `D:\\tmp\\w330pst` 的字节、逐发当场还原核回提交内 sha256）："
    "去牙发＝长尾候选 `Resident: false→true` ⇒ 唯一红在 `plan_present_crowded_world_330r2_test.go:247`"
    "（`tier leak: resident arm = [plan.present ledger.sweep330r2], want \"ledger.sweep330r2\" kept out of it …`）；"
    "反控发＝查询换成与长尾无共同字 `\"frobnicate widget quietly\"` ⇒ 红在另一行 `:254`"
    "（`retrieval arm = [], want \"ledger.sweep330r2\" in it …`）——红句与行号都不同 ⇒ 这枚用例响的是档位、⛔ 恒响。"
    "G2 那一发重放（`plan_present_doorbell_165r3_test.go:69` 拨 `false`）⇒ 红在 `:188`（第一轮请求体里那扇门整个不见了）"
    "＋`:240`（驻留臂空），而 165 那枚用例本身同发仍 `--- PASS`＝它的盲区属 G1／330-r1 的 `cmd/wisp` 面。"
    "★必答题：把新增那枚候选删掉（只删本腿自己加的夹具）再重放档位突变，"
    "`TestCrowdedWorldStillOffersTheBrake330r2` **变回绿＝答\"是\"**（牙确实挂在多候选世界上）；"
    "另一枚 `TestResidencyTierDecidesTheInjectionArm330r2` 仍红，接住它的是**枚数前提门 `:221`**"
    "（那扇门按名字提名了被删的候选、与档位无关），具名交裁、⛔ 算本腿的 G2 牙。"
    "台面与颜色：仓外树八件套＋仓根 `docs`＋手工三件＋`GOTMPDIR` 指 `D:` 卷，"
    "改前（锚 `feab45e3`）`rc=0／RUN 110／PASS_TOP 86／FAIL_TOP 0／SKIP 0／PANIC 0` ⇄ "
    "改后（含本腿两枚文件）`rc=0／RUN 112／PASS_TOP 88／FAIL_TOP 0／SKIP 0／PANIC 0`，红名册作差＝0、"
    "同包⛔ 一枚别枚用例被打红（故无需改任何既有期望值）。"
    "★零产码尺：`git diff --name-only feab45e3 HEAD -- internal/agent | grep -vE '(_test\\.go|^internal/agent/testdata/)'` "
    "＝空、rc=1；同一把尺的正控跑在 `git clone --local` 的副本里（把 `internal/agent/inject.go` 一处非注释字面量 "
    "`MaxRetrievedTools` 改成 `MaxRetrievedTools+1`）＝非空、rc=0 ⇒ 这把尺不失明。"
    "门禁四数只在共享工作树：`gofumpt -l internal/agent/` 列出 `approval/pending_read.go`／`tools.go` 两枚"
    "＝本锚点既存漂移（`git status` 对它们为空、⛔ 本腿碰过，本腿两枚文件不在名单）；`sh scripts/d22scan.sh` rc=1，"
    "唯一 finding＝`cmd/wisp/panel_snapshot_deliver_329.go:14 [phantom-citation]`（并行腿 329-r1 的在飞字节，"
    "本腿新文件被提及 0 次，具名交编排者）；长度门 rc=0 VERDICT GREEN；`go vet ./internal/agent/` rc=0 输出 0 字节。"
    "G3 随批＝`plan_present_doorbell_165r3_test.go:131` 那句注释改成\"冗余非零门，牙在 :145 的 `rowOfTool165r3` 独立遍历\""
    "（1 删 3 加、221→223 行、新串 `6542773d3d2edcc6…`；⛔ 删那扇门；\"把 sink 计数喂成常量\"那一发本腿未复测，"
    "⛔ 记本腿读数）。基线件一字未动：`internal/agent/testdata/golden/plan-present.sse` 仍 `7d26fa889b59fc3a…`。"
    "证据件＝`docs/evidence/s1/330-r2-g2-teeth.md`（102 行／9647 字节），五发原始读数＝`.scratch/wisp/probes/330/r2/logs/m1..m5.log`。"
    "AC#1/#2 仍空（等 330-r1 在 `cmd/wisp/**` 落对接钉）；AC#0/#3/#4/#5/#6/#7 本腿⛔ 自翻、交裁决腿。\n")


def main():
    with io.open(P, "r", encoding="utf-8", newline="") as f:
        before = f.read()
    h = hashlib.sha256(before.encode("utf-8")).hexdigest()
    if h != PRE:
        print("ABORT: ticket file sha256=%s != pre-append baseline %s" % (h, PRE))
        return 5
    if not before.endswith("\n"):
        print("ABORT: ticket file does not end with a newline")
        return 6
    if "330-r2" in before:
        print("ABORT: a 330-r2 bullet is already present, refusing to append twice")
        return 7
    with io.open(P, "a", encoding="utf-8", newline="") as f:
        f.write(BULLET)
    with io.open(P, "r", encoding="utf-8", newline="") as f:
        after = f.read()
    if not after.startswith(before) or after != before + BULLET:
        print("ABORT: the append is not a pure tail append")
        return 8
    print("appended one bullet: bytes %d -> %d, lines %d -> %d"
          % (len(before.encode("utf-8")), len(after.encode("utf-8")),
             before.count("\n"), after.count("\n")))
    print("sha256_after=%s" % hashlib.sha256(after.encode("utf-8")).hexdigest())
    return 0


if __name__ == "__main__":
    sys.exit(main())
