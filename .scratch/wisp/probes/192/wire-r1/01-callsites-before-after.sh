#!/usr/bin/env bash
# Ticket 192's wiring leg - the three call-site readings, reproducible.
#
#   sh .scratch/wisp/probes/192/wire-r1/01-callsites-before-after.sh
#
# "before" is HEAD: the commit this leg started from. "after" is the worktree,
# with --untracked so the new production file is counted (a plain `git grep`
# skips it and would report the same zero as before - the trap this script
# exists to avoid).
set -u
cd "$(dirname "$0")/../../../../.." || exit 2

echo "== A) QuoteForModel CALL sites ('.QuoteForModel('), non-test .go"
echo "-- before (HEAD):"
git grep -n '\.QuoteForModel(' HEAD -- '*.go' ':!*_test.go'
echo "   count=$(git grep -n '\.QuoteForModel(' HEAD -- '*.go' ':!*_test.go' | grep -c . )"
echo "-- after (worktree, untracked included):"
git grep -n --untracked '\.QuoteForModel(' -- '*.go' ':!*_test.go'
echo "   count=$(git grep -n --untracked '\.QuoteForModel(' -- '*.go' ':!*_test.go' | grep -c . )"

echo
echo "== B) panel.NewPreviewResolver CALL sites, cmd/ + internal/, non-test"
echo "-- before (HEAD):"
git grep -n 'panel\.NewPreviewResolver(' HEAD -- 'cmd/*.go' 'internal/*.go' ':!*_test.go'
echo "   count=$(git grep -n 'panel\.NewPreviewResolver(' HEAD -- 'cmd/*.go' 'internal/*.go' ':!*_test.go' | grep -c . )"
echo "-- after (worktree, untracked included):"
git grep -n --untracked 'panel\.NewPreviewResolver(' -- 'cmd/*.go' 'internal/*.go' ':!*_test.go'
echo "   count=$(git grep -n --untracked 'panel\.NewPreviewResolver(' -- 'cmd/*.go' 'internal/*.go' ':!*_test.go' | grep -c . )"

echo
echo "== C) the preview BRANCH in cmd/wisp (panel.PreviewPathPrefix), non-test"
echo "-- before (HEAD):"
git grep -n 'PreviewPathPrefix' HEAD -- 'cmd/*.go' ':!*_test.go'
echo "   count=$(git grep -n 'PreviewPathPrefix' HEAD -- 'cmd/*.go' ':!*_test.go' | grep -c . )"
echo "-- after (worktree, untracked included):"
git grep -n --untracked 'PreviewPathPrefix' -- 'cmd/*.go' ':!*_test.go'
echo "   count=$(git grep -n --untracked 'PreviewPathPrefix' -- 'cmd/*.go' ':!*_test.go' | grep -c . )"

echo
echo "== D) the one production file this leg adds"
git status --porcelain -- cmd/wisp/panel_preview_192.go
