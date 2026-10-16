#!/usr/bin/env python3
# 220-v1 突变器（裁决腿自用）：只打在仓外导出树（$WISP220ROOT）的物理字节上。
# 每发：先备份原字节 ⇒ 打一处精确串替换（命中数必须＝1，否则拒打）⇒ 交回 bash 跑用例
# ⇒ restore 回原字节并逐枚 sha256 对回。⛔ -overlay，⛔ 动仓内（共享工作树）任何字节。
import hashlib
import os
import sys

ROOT = os.environ["WISP220ROOT"]

PANEL = "internal/panel/subagent_roster_197.go"
PUMP = "cmd/wisp/panel_pump.go"
RUN = "cmd/wisp/run.go"

MUTS = {
    # MU-A：AC#1 的牙——把卡片那一半的"corr＋#序号 去牙双索引"退成 220 之前的单查形
    # （只按卡片自己的 corr 建索引，不再认 queue.go 的 "<id>#<seq>" 改写）。
    "MU-A": (PANEL,
             "\tfor _, card := range pending {\n\t\tindexWaitingKey(waiting, card.CorrelationID)\n\t}\n",
             "\tfor _, card := range pending {\n\t\tif card.CorrelationID != \"\" {\n\t\t\twaiting[card.CorrelationID] = true\n\t\t}\n\t}\n"),
    # MU-B：AC#1③ 的反形——L1 那一半只按 corr 索引，丢掉 taskID 那一半。
    "MU-B": (PANEL,
             "\t\tindexWaitingKey(waiting, w.CorrelationID)\n\t\tindexWaitingKey(waiting, w.TaskID)\n",
             "\t\tindexWaitingKey(waiting, w.CorrelationID)\n"),
    # MU-C：AC#1③ 的第二形反形——查键改成行自己的 corr 维（StreamKey），不再按 taskID 查。
    "MU-C": (PANEL,
             "BlockedOnApproval: row.TaskID != \"\" && waiting[row.TaskID],",
             "BlockedOnApproval: row.StreamKey != \"\" && waiting[row.StreamKey],"),
    # MU-D：AC#2③ 的牙——摘掉装配根那一行供数（生产者＝0 枚那一形）。
    "MU-D": (RUN, "\t\tL1Windows: rt.liveL1WindowWaits,\n", ""),
    # MU-E：AC#2④ 的反形——供一枚常量非空列表（名册那一行恒 true，与窗口在不在无关）。
    "MU-E": (PUMP,
             "\twins := rt.gate.LiveL1Windows()\n\tif len(wins) == 0 {\n\t\treturn nil\n\t}\n",
             "\twins := rt.gate.LiveL1Windows()\n\t_ = wins\n\tst := rt.taskRosterState()\n\tconstant := make([]panel.L1WindowWait, 0, len(st.Rows))\n\tfor _, r := range st.Rows {\n\t\tconstant = append(constant, panel.L1WindowWait{CorrelationID: \"const-220v1\", TaskID: r.TaskID})\n\t}\n\tif len(constant) > 0 {\n\t\treturn constant\n\t}\n"),
}


def path(rel):
    return os.path.join(ROOT, rel.replace("/", os.sep))


def sha(p):
    with open(p, "rb") as f:
        return hashlib.sha256(f.read()).hexdigest()


def main():
    action = sys.argv[1]
    label = sys.argv[2]
    rel, old, new = MUTS[label][0], MUTS[label][1], MUTS[label][2]
    p = path(rel)
    if action == "pre":
        print("%s BASE %s sha=%s" % (label, p, sha(p)))
        with open(p, "rb") as f:
            data = f.read()
        n = data.count(old.encode())
        print("%s anchor_hits=%d" % (label, n))
        if n != 1:
            sys.exit(3)
    elif action == "apply":
        with open(p + ".orig-220v1", "wb") as f:
            f.write(open(p, "rb").read())
        data = open(p, "rb").read()
        assert data.count(old.encode()) == 1, label
        with open(p, "wb") as f:
            f.write(data.replace(old.encode(), new.encode(), 1))
        print("%s APPLIED sha=%s" % (label, sha(p)))
    elif action == "restore":
        with open(p, "wb") as f:
            f.write(open(p + ".orig-220v1", "rb").read())
        base = sha(p + ".orig-220v1")
        now = sha(p)
        print("%s RESTORED sha=%s base=%s match=%s" % (label, now, base, now == base))
        if now != base:
            sys.exit(4)
    else:
        print("unknown action")
        sys.exit(2)


main()
