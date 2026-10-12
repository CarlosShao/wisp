#!/usr/bin/env bash
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
D=.scratch/wisp/probes/310/v1/logs
: > "$D/10-block-bytes.txt"
say() { echo "$@" >> "$D/10-block-bytes.txt"; }
say "go_env_GOOS=$(go env GOOS)"
say "block_bounds_as_read_by_me (start=the TICKET 310 comment line, end=the closing fi of that block):"
for spec in "scripts/portable-tests.sh 99 166" "scripts/winsec-tests.sh 46 66" "scripts/wisp-cli-tests.sh 60 79"; do
  set -- $spec
  f=$1; s=$2; e=$3
  bytes=$(awk -v s="$s" -v e="$e" 'NR>=s && NR<=e {b+=length($0)+1} END{print b+0}' "$f")
  lines=$((e - s + 1))
  left=$(awk -v s="$s" -v e="$e" 'NR>=s && NR<=e' "$f" | grep -c 'left as-is' || true)
  exp=$(awk -v s="$s" -v e="$e" 'NR>=s && NR<=e' "$f" | grep -c 'export GOTMPDIR' || true)
  deg=$(awk -v s="$s" -v e="$e" 'NR>=s && NR<=e' "$f" | grep -c 'DEGRADED' || true)
  say "block=$f lines=$lines bytes=$bytes export_GOTMPDIR_lines=$exp DEGRADED_lines=$deg left_as_is_lines=$left"
done
say "non_windows_branch_reads:"
awk 'NR>=160 && NR<=166' scripts/portable-tests.sh >> "$D/10-block-bytes.txt"
say "winsec_block_condition_line=$(awk 'NR==54' scripts/winsec-tests.sh | cut -c1-140)"
say "cli_block_condition_line=$(awk 'NR==67' scripts/wisp-cli-tests.sh | cut -c1-140)"
say "cli_pre_existing_non_windows_guard_at=$(awk 'NR==80' scripts/wisp-cli-tests.sh | cut -c1-140) (present in the pre-leg blob too: grep count was 1)"
say "ubuntu_leg_named_on_disk=$(grep -n 'ubuntu leg is left exactly as it is' scripts/portable-tests.sh | cut -c1-140)"
say "block_bytes_done=1"
