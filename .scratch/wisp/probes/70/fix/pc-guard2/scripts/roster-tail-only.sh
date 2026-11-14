#!/usr/bin/env bash
# POSITIVE CONTROL input for ticket 70 / winsec-tests.sh GUARD 2.
#
# This is the exact byte shape CI run 43549905387 test-windows step 4 fed GUARD 2
# (docs/evidence/s1/70-ci-red-triage.md step 4): the `ok  <pkg>  <duration>` line
# that `go test -list` prints as its trailer, and NOTHING that only a real -v run
# can print - no `--- PASS|FAIL|SKIP:` verdict line, no `portable-tests.sh: four
# numbers` line. It exits 0 on purpose: the point of the control is that GUARD 2
# must red this on its own, without help from a non-zero child exit code.
printf 'ok  \tgithub.com/CarlosShao/wisp/internal/winsec\t0.069s\n'
printf 'PASS\n'
exit 0
