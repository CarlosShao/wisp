#!/bin/sh
# 120-r2 gate round 3: the uid-1000 permission faces, pre-fix and post-fix on the
# same probe bytes, plus the whole package at uid 1000 on both. (gate2.sh carried a
# bug: a shell variable holding "MSYS_NO_PATHCONV=1 docker ..." is executed as one
# command word, so every run below it died before reaching docker. Kept, not deleted.)
set -e
OUT="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/120/r2"
S=/d/tmp/probe120/r2

runin() {
	# runin <snapshot> <ignored> <docker-user-flag (two words, split on purpose)> <command>
	MSYS_NO_PATHCONV=1 docker run --rm $3 -e CGO_ENABLED=0 -e GOCACHE=/tmp/gocache \
		-v "$1":/work -v probe120build:/root/.cache/go-build -w /work golang:1.27 sh -c "$4"
}

for code in base-head base-fixed; do
	for who in root uid1000; do
		if [ "$who" = root ]; then
			flag="--user 0:0"
		else
			flag="--user 1000:1000"
		fi
		echo "== permface code=$code as=$who =="
		runin "/d/tmp/probe120/r2/$code" x "$flag" 'id -u; go run ./permface' > "$OUT/120-$who-permface-$code.txt" 2>&1 || true
		sed -n '1p;$p' "$OUT/120-$who-permface-$code.txt"
	done
done

echo "== whole package at uid 1000, pre-fix vs post-fix =="
for code in base-head base-fixed; do
	echo "-- code=$code --"
	runin "/d/tmp/probe120/r2/$code" x "--user 1000:1000" 'id -u; go test -count=1 ./internal/winsec/ 2>&1 | tail -40' > "$OUT/120-uid1000-package-$code.log" 2>&1 || true
	tail -6 "$OUT/120-uid1000-package-$code.log"
done
