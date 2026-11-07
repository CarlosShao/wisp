# 330-ac5 台面解析尺 v2（纯 awk，无产码）
# 口径（与票 329 AC#7/227-v2 同族）：RUN=全部 ^=== RUN 含子测试；PASS/FAIL/SKIP_TOP=仅零缩进判决行；
# 顶层=名字不带 "/" 的零缩进 RUN；死点=顶层 RUN 无判决行。
LOG="$1"
LABEL="${2:-log}"
awk -v lbl="$LABEL" '
  /^=== RUN /   { run++; name=$3;
                  if (name !~ /\//) { top++; topname[top]=name; topdead[name]=1 }
                  next }
  /^--- PASS/   { passt++; delete topdead[$3]; next }
  /^--- FAIL/   { failt++; red[nred++]=$3; delete topdead[$3]; next }
  /^--- SKIP/   { skipt++; delete topdead[$3]; next }
  /(^|[^[:alpha:]])panic:/ { pan++ }
  END {
    printf "[%s] RUN=%d PASS_TOP=%d FAIL_TOP=%d SKIP_TOP=%d PANIC=%d\n", lbl, run, passt, failt, skipt, pan;
    nd=0; s="";
    for (i=1; i<=top; i++) if (topname[i] in topdead) { nd++; s=s (nd>1?" ; ":"") topname[i] }
    printf "[%s] DEADPOINTS=%d %s\n", lbl, nd, s;
    printf "[%s] RED_ROSTER=%d\n", lbl, nred;
    for (i=0; i<nred; i++) printf "[%s] RED: %s\n", lbl, red[i];
  }' "$LOG"
