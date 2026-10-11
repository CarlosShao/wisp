#!/bin/sh
# A846 读数 第四发：票 294 裁定节里那五条承重数字，我自己现跑一遍（writeface-1 报四条过期＋一枚颜色写反；
# ⛔ 转抄它的读数。尺一律＝git grep HEAD 对象层，⛔ 工作树＝树会呼吸那条仍成立）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A846-ticket294-rulers.txt

{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')  HEAD=$(git rev-parse --short HEAD)"
  echo "### 1) 测试直调 attachMuteGate 的枚数（票面裁定节写 3／腿报 5）"
  git grep -n "attachMuteGate" HEAD -- cmd/wisp | sed 's/^HEAD://'
  echo "  计数（含产码与注释）=$(git grep -c "attachMuteGate" HEAD -- cmd/wisp | awk -F: '{s+=$NF} END {print s}')"
  echo "  --- 分行名册（下一位可逐枚点）：直调那一把＝剥掉定义/注释后的 _test.go 命中"
  echo "  测试件内命中=$(git grep -n "attachMuteGate" HEAD -- 'cmd/wisp/*_test.go' | grep -c '')"
  echo "  测试件内剥注释行=$(git grep -n "attachMuteGate" HEAD -- 'cmd/wisp/*_test.go' | grep -v '//' | grep -c '')"
  echo "  产码命中=$(git grep -n "attachMuteGate" HEAD -- cmd/wisp | grep -v '_test.go' | grep -c '')"
} > $LOG 2>&1

{ echo "### 2) 手搓 &residentBall{ 的枚数（票面写 13／腿报 16）"
  echo "  尺=git grep -n '&residentBall{' HEAD -- cmd/wisp 命中=$(git grep -n "&residentBall{" HEAD -- cmd/wisp | grep -c '')"
  echo "  其中 _test.go=$(git grep -n "&residentBall{" HEAD -- 'cmd/wisp/*_test.go' | grep -c '')  剥注释=$(git grep -n "&residentBall{" HEAD -- 'cmd/wisp/*_test.go' | grep -v '//' | grep -c '')"
  echo "  逐枚件名计数=$(git grep -c "&residentBall{" HEAD -- cmd/wisp | tr '\n' ' ')"
} >> $LOG 2>&1

{ echo "### 3) 两枚行号锚（票面 :336 生产挂载 / :472 定义）现在在哪"
  echo "  attachMuteGate 产码行=$(git grep -n "rb.attachMuteGate" HEAD -- cmd/wisp | sed 's/^HEAD://')"
  echo "  attachMuteGate 定义行=$(git grep -n "func .*attachMuteGate" HEAD -- cmd/wisp | sed 's/^HEAD://')"
} >> $LOG 2>&1

{ echo "### 4) ★那枚'恒绿假钉'的颜色：票面 AC#1 ⓐ 逐字要的是什么"
  grep -n 'AC#1' .scratch/wisp/issues/294-*.md | sed -n '1,3p' | cut -c1-900
} >> $LOG 2>&1

{ echo "### 5) SelectorExpr 那把尺（我说'必须加一层'，同包是否早有范本）"
  echo "  git grep -c 'ast.SelectorExpr' HEAD -- cmd/wisp：命中总数=$(git grep -n 'ast.SelectorExpr' HEAD -- cmd/wisp | grep -c '')  件数=$(git grep -c 'ast.SelectorExpr' HEAD -- cmd/wisp | grep -c '')"
  git grep -n 'ast.SelectorExpr' HEAD -- cmd/wisp | sed 's/^HEAD://' | cut -c1-160
} >> $LOG 2>&1

{ echo "### 6) executed ⇄ started 零耦合那把尺（串行约束的真身）"
  echo "  executed 在 cmd/wisp 命中=$(git grep -n 'executed' HEAD -- cmd/wisp | grep -c '')"
  echo "  executed 与 started 同现于一枚行？=$(git grep -nE 'executed.*started|started.*executed' HEAD -- cmd/wisp | grep -c '')"
  echo "  唯一断言那一行=$(git grep -n 'if executed {' HEAD -- cmd/wisp | sed 's/^HEAD://')"
  echo "  started 在 internal/audio 命中（各件）=$(git grep -c '\bstarted\b' HEAD -- internal/audio | tr '\n' ' ')"
  echo "  started 在 cmd/wisp 命中=$(git grep -n '\.started' HEAD -- cmd/wisp | grep -c '')"
} >> $LOG 2>&1
cat $LOG
