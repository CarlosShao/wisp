set -u
F=internal/tools/bridge.go
LOGS=.scratch/wisp/probes/309/r1/logs
ORIG="$LOGS/bridge.go.orig"
cp "$F" "$ORIG"
restore() { cp "$ORIG" "$F"; }
trap restore EXIT

runrig() { # $1=rig name
  echo "########## RIG $1 ##########"
  PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1 \
    -run 'Test309AssessorForEmptyTaskIDDropsTaintDetector' -v 2>&1
  echo "RIG_$1_rc=$?"
  restore
}

# --- before (unmutated): must be GREEN ---
echo "########## BEFORE (unmutated) ##########"
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1 \
  -run 'Test309AssessorForEmptyTaskIDDropsTaintDetector' -v 2>&1
echo "BEFORE_rc=$?"

# --- RIG 1 (shape A): drop the `|| taskID == ""` limb, content-anchored ---
cp "$ORIG" "$F"
perl -0pi -e 's/\tif b\.injected \|\| b\.prov == nil \|\| taskID == "" \{/\tif b.injected || b.prov == nil {/' "$F"
grep -n 'if b.injected || b.prov == nil' "$F"
runrig A_drop_limb

# --- RIG 2 (shape B): make the attach branch attach nothing (detector not attached) ---
cp "$ORIG" "$F"
perl -0pi -e 's/WithTaintDetector\(b\.prov\.Detector\(taskID\)\)/WithTaintDetector(nil)/' "$F"
grep -n 'WithTaintDetector(nil)' "$F"
runrig B_no_detector

# --- prove byte-identical restore ---
restore
git diff --exit-code -- "$F" && echo "RESTORE_BYTE_IDENTICAL=yes rc=$?" || echo "RESTORE_DIFF rc=$?"
diff -q "$ORIG" "$F" && echo "IDENTICAL_TO_ORIG=yes"
