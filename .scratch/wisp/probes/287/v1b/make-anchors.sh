#!/usr/bin/env bash
# 287-v1b：写出八枚突变台的锚文件（逐字节锚，⛔ 行号锚；正文全是 ASCII，无 shell 元字符风险）
set -eu
A="$(cd "$(dirname "$0")" && pwd)/anchors"
mkdir -p "$A"

w() { # name, content
  printf '%s' "$2" >"$A/$1"
}

# V1 深度那一支摘掉（恒真检查的主发；＝"判定永假"）
w v1.old 'if rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {'
w v1.new 'if rec, ok := t.d.Roster.Look(parentID); false && ok && rec.Kind == TaskKindSubagent {'
# V3 判定放宽成"有名册行就拒"（恒真检测器：只有反形该红）
w v3.old 'if rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {'
w v3.new 'if rec, ok := t.d.Roster.Look(parentID); ok {'
# V2 身份那一维改读 corr（判定读错身份）
w v2.old 'parentID := TaskID(ctx)'
w v2.new 'parentID := CorrelationID(ctx)'
# V4 判定仍按 taskID、但句子把派生者写成 corr ⇒ 钉夹具 :180-182 那一发（前任具名"无人钉"的那枚）
w v4.old 'MaxSubagentDepth, parentID, rec.Kind), IsError: true}, nil'
w v4.new 'MaxSubagentDepth, CorrelationID(ctx), rec.Kind), IsError: true}, nil'
# V6 派发不走真 loop 的那一形：请求侧 corr 塌到 taskID ⇒ 钉夹具 corrOf287 :140-141 那一发
w v6.old 'TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),'
w v6.new 'TaskID: taskID, CorrelationID: taskID,'
# V7 Kind 由产码 PublishSubagent 写进名册那一支 ⇒ 钉夹具 :159-161 那一发（⛔ 夹具自填字段）
w v7.old 'ParentTaskID: parentTaskID, Label: label, Kind: TaskKindSubagent,'
w v7.new 'ParentTaskID: parentTaskID, Label: label, Kind: TaskKindRoot,'
# V8 派发经真 loop：AdmitTask 钩子不再被 loop 调用 ⇒ 钉夹具 :95-97 那一发
w v8.old 'if l.opt.AdmitTask != nil {'
w v8.new 'if false && l.opt.AdmitTask != nil {'

ls -1 "$A" | sed 's/^/anchor: /'
wc -c "$A"/* | tail -1
