# 2026-11-08 argv-overflow 修腿量化件：bash 未引号展开规模 vs git glob 射程（只建不删）
cd "$(git rev-parse --show-toplevel)" || exit 9
GO='internal/**/*.go'
GO2='cmd/**/*.go'
cnt=0; chars=0
for w in $GO $GO2; do cnt=$((cnt + 1)); chars=$((chars + ${#w} + 1)); done
echo "bash-unquoted-expansion: words=$cnt chars=$chars"
echo "git-ls-files-recursive(**): $(git ls-files 'internal/**/*.go' 'cmd/**/*.go' | wc -l)"
echo "git-ls-files-onedep(*/*): $(git ls-files 'internal/*/*.go' 'cmd/*/*.go' | wc -l)"
echo "depth>=4 files: $(git ls-files 'internal/**/*.go' 'cmd/**/*.go' | awk -F/ 'NF>=4' | wc -l)"
echo "quoted-glob git grep G6neg probe:"
git grep -nEw 'OpenTask|CloseTask' HEAD -- 'internal/**/*.go' 'cmd/**/*.go' >/dev/null; echo "  rc=$? lines=$(git grep -nEw 'OpenTask|CloseTask' HEAD -- 'internal/**/*.go' 'cmd/**/*.go' | wc -l)"
