#!/usr/bin/env bash
cd /tmp/w239r1-export || exit 1
for x in $(git ls-files | grep -Ei 'orb.*\.html|jarvis.*\.html'); do
  wt=$(tr -cd '\r' < "$x" | wc -c)
  bl=$(git cat-file blob "HEAD:$x" | tr -cd '\r' | wc -c)
  echo "$x worktree_CR=$wt blob_CR=$bl"
done
echo "export_autocrlf=$(git config core.autocrlf || echo unset)"
cd "/d/work/workspace/projects plans/Wisp"
echo "main_autocrlf=$(git config core.autocrlf || echo unset)"
