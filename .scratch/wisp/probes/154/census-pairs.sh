# 票 154 AC#3 —— 同族普查的读数生成器（可复算）。
# 用法：bash .scratch/wisp/probes/154/census-pairs.sh [锚点]  > 读数文件
# 每一对给三行：开那侧的生产出现点／关那侧的生产出现点／差值判读由证据件写。
set -u
cd "$(git rev-parse --show-toplevel)" || exit 9
A="${1:-$(git rev-parse HEAD)}"
P='internal' ; P2='cmd'
EX=':!*_test.go'

pair() { # $1 标题 $2 开侧 pattern $3 关闭侧 pattern [$4 额外 pathspec 排除]
	local title="$1" open="$2" close="$3" extra="${4:-}"
	echo
	echo "########## $title"
	echo "--- 开那侧（生产码）---"
	# shellcheck disable=SC2086
	git grep -nE "$open" "$A" -- "$P" "$P2" "$EX" $extra || echo "(零命中)"
	echo "--- 关那侧（生产码）---"
	# shellcheck disable=SC2086
	git grep -nE "$close" "$A" -- "$P" "$P2" "$EX" $extra || echo "(零命中)"
}

echo "# 票 154 AC#3 同族普查 —— 锚点 $A / 生成 $(date -Iseconds)"
echo "# 尺：git grep -nE（大小写敏感、不整词、pathspec=internal+cmd、排除 _test.go；每对另给的排除项写在标题里）"

pair "1 Bridge.OpenTask / Bridge.CloseTask（本票主题）" 'OpenTask' 'CloseTask' ':!internal/tools/bridge.go'
pair "2 Provenance.OpenScope / CloseScope（internal/risk 冻结面，本程只读）" 'OpenScope\(' 'CloseScope\('
pair "3 Gate.AdmitTextTask 的 revoke（关那侧是闭包，方法名 grep 抓不到）" 'AdmitTextTask' 'revoke\(\)|defer revoke'
pair "4 plugin.DisposalScope：NewDisposalScope/Defer vs Dispose" 'NewDisposalScope|\.Defer' '\.Dispose\('
pair "5 Store.StartRetentionJob（需要一枚 DisposalScope）" 'StartRetentionJob' 'StartRetentionJob'
pair "6 tools.Registry.Register / Registry.RegisterProvider，有无 Unregister" '\.Register\(|RegisterProvider' 'Unregister'
pair "7 RiskAssessor.RegisterSubAssessor（只有注册、没有反注册）" 'RegisterSubAssessor' 'UnregisterSubAssessor|removeSubAssessor'
pair "8 proc.JobScope：StartInJob vs Close" 'StartInJob' 'job\.Close\(|Job\.Close\(|CloseJob'
pair "9 proc.SingleInstance：Acquire 与 Release" 'AcquireSingleInstance' '\.Release\('
pair "10 audio：麦克风腿的 Start/Stop 成对（整条腿今天接没接宿主）" 'NewHalfDuplexGate|WASAPIMicrophone\{|NewWavInjector' '\.Stop\(\)'
pair "11 panel.StreamLog：Append 隐式开 key vs Close(key)" 'stream\.Append|\.Append\(' 'stream\.Close\(|func \(s \*StreamLog\) Close'
pair "12 memory.Store：InsertGrant vs RevokeGrant/DeleteGrant" 'InsertGrant\(' 'RevokeGrant\(|DeleteGrant\('
pair "13 ball：RegisteredHotkeys 注册 vs UnregisterHotKey（Win32）" 'RegisteredHotkeys|RegisterHotKey' 'UnregisterHotKey|ReleaseEscAfterSession|\.Close\(\)'
pair "14 observe.LogPipeline.Close / rollingWriter.Close" 'NewLogPipeline' 'pipeline\.Close\(|\.Close\(\)'

echo
echo "########## 15 可见性那一族（被上游码挡住）"
echo "--- fs.read 的声明档 ---"
git grep -n -A 4 "func FSReadDecl" "$A" -- internal/tools/fs.go | sed -n '1,12p'
echo "--- loop 那侧的 default 分支（L0/未分级走哪条）---"
git grep -nE "case RiskL1, RiskL2:|default:|PassThroughUnclassifiedRisk" "$A" -- internal/agent/loop.go
echo "--- 生产组合根的 agent.Config 字面量里有没有 PassThroughUnclassifiedRisk ---"
git grep -n "PassThroughUnclassifiedRisk" "$A" -- cmd/wisp/run.go || echo "(零命中)"
echo "--- 测试那侧谁打开了它 ---"
git grep -n "PassThroughUnclassifiedRisk:" "$A" -- '*.go'
echo "--- 能 Mark 的工具名交集（注册名 ∩ sensitiveSourceTools）---"
git grep -hE 'Name\(\) string \{ return "' "$A" -- internal/tools | grep -v _test | sed -E 's/.*return "([^"]+)".*/\1/' | sort > /d/tmp/wisp154-registered-names.txt
git grep -hE "Src[A-Za-z]+ += \"" "$A" -- internal/risk/provenance.go | sed -E 's/.*= "([^"]+)".*/\1/' | sort > /d/tmp/wisp154-sensitive-names.txt
echo "注册名：$(tr '\n' ' ' < /d/tmp/wisp154-registered-names.txt)"
echo "sensitive：$(tr '\n' ' ' < /d/tmp/wisp154-sensitive-names.txt)"
echo "交集：$(comm -12 /d/tmp/wisp154-registered-names.txt /d/tmp/wisp154-sensitive-names.txt | tr '\n' ' ')"
