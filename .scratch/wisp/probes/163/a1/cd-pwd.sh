#!/bin/sh
# 票 163 a1 台件（只读设计核，量现状，不修码）
# 这一发量的是「一次性 spawn」的形状：两次完全独立的进程生发，
# 第一次 cd、第二次取 pwd —— 中间没有任何存活载体。
# 这就是今天 shell.exec 名义下唯一可量的真实基底（注册表里其实一枚 shell 工具都没有，
# 见同目录 gotest 台件的 Lookup 读数）。
# ⚠ logdir 取脚本自身目录，不继承 CWD。临时件只建不删。
set -u
PROBE_DIR="$(cd "$(dirname "$0")" && pwd)"
LOGDIR="$PROBE_DIR/logs"
mkdir -p "$LOGDIR"
LOG="$LOGDIR/cd-pwd-2nd.log"  # 第一枚 cd-pwd.log 里是横幅坏读数，按规矩留存不删

{
	echo "# 台件 cd-pwd.sh 读数"
	date -Iseconds
	echo "# 生发器＝cmd.exe（Windows 上 shell.exec 未来会生发的那一枚真 shell；非镜像、非假进程）"
	echo
# ⚠ MSYS 会把 /c 当路径改写，这里逐发禁掉路径转换（MSYS_NO_PATHCONV=1），
# 第一版没加，cmd.exe 落成了交互横幅（读数无效，旧日志原样留在 logs/ 不删）。
	echo "## 发 1: cmd.exe /c \"cd /d C:\\Windows & cd\""
	MSYS_NO_PATHCONV=1 cmd.exe /c 'cd /d C:\Windows & cd' 2>&1
	rc1=$?
	echo "# rc=$rc1"
	echo
	echo "## 发 2（另起一发全新进程，只问当前目录）: cmd.exe /c \"cd\""
	MSYS_NO_PATHCONV=1 cmd.exe /c 'cd' 2>&1
	rc2=$?
	echo "# rc=$rc2"
	echo
	echo "## 读数判语：发 2 打印的目录若不是 C:\\Windows，则两步之间没有任何状态存活。"
	echo "# 另外：exit 码取回验证 cmd.exe /c exit 7"
	MSYS_NO_PATHCONV=1 cmd.exe /c 'exit 7' 2>&1
	echo "# rc=$?（这是操作系统 Wait 状态取回的字段，不是从输出文本里猜的）"
} > "$LOG" 2>&1
cat "$LOG"
echo "# log=$LOG"
