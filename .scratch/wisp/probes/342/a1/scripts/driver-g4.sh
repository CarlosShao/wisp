set -u
cd /src
echo "== cross-vet BUILD matrix for ./internal/winsec (CGO off, GOPROXY off) =="
for g in linux darwin freebsd openbsd netbsd dragonfly solaris windows; do
  out=$(GOOS=$g go build ./internal/winsec/ 2>&1); rc=$?
  echo "GOOS=$g rc=$rc $(echo "$out" | head -2 | tr '\n' ' ')"
done
outjs=$(GOOS=js go build ./internal/winsec/ 2>&1); echo "GOOS=js rc=$? $(echo "$outjs" | head -2 | tr '\n' ' ')"
outp9=$(GOOS=plan9 go build ./internal/winsec/ 2>&1); echo "GOOS=plan9 rc=$? $(echo "$outp9" | head -2 | tr '\n' ' ')"
echo "== config full-package 3x resample per face (BASE treeB vs mut-jia vs mut-yi) =="
for t in /src /jia /yi; do
  echo "face=$t"
  for i in 1 2 3; do
    r=$(cd $t && GOMODCACHE=/gomod GOCACHE=/gobuild go test -count=1 -timeout 10m ./internal/config/ 2>&1 | grep -E '^(ok|FAIL|--- FAIL)' | tr '\n' ' ')
    echo "  run$i: $r"
  done
done
echo "== risk C26 name 3x resample per face =="
for t in /src /jia /yi; do
  echo "face=$t"
  for i in 1 2 3; do
    r=$(cd $t && GOMODCACHE=/gomod GOCACHE=/gobuild go test -count=1 -timeout 20m -run TestC26RewriteAccountIsConsumedAtEverySecurityLeg ./internal/risk/ 2>&1 | grep -E '^(ok|FAIL|--- FAIL)' | tr '\n' ' ')
    echo "  run$i: $r"
  done
done
echo DONE-G4