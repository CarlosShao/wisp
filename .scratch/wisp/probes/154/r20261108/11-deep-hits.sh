# 2026-11-08 修腿：深层文件命中普查（引号修法是否动基线；只建不删）
cd "$(git rev-parse --show-toplevel)" || exit 9
echo "== 深层文件(NF>=4)上的腿命中 =="
git grep -nE 'ToolRequest\{|\.Execute\(|OpenTask|CloseTask|OpenScope|CloseScope|Defer(Named)?|DisposalScope' HEAD -- 'internal/**/*.go' 'cmd/**/*.go' | awk -F: '{n=split($1,a,"/"); if(n>=4) print}' | tee /dev/stderr | wc -l
echo "== 测试文件命中 =="
git grep -nE 'ToolRequest\{|\.Execute\(|OpenTask|CloseTask|OpenScope|CloseScope|Defer(Named)?|DisposalScope' HEAD -- 'internal/**/*.go' 'cmd/**/*.go' | grep -c '_test\.go:' || true
