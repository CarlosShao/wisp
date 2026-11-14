#!/usr/bin/env bash
# POSITIVE CONTROL input for ticket 70 / winsec-tests.sh GUARD 2, the GREEN arm.
#
# Replays bytes captured from a REAL `bash scripts/portable-tests.sh ./internal/winsec/`
# invocation (PC70_REAL_CHILD names the file). No line in the replayed stream is
# invented here: every `--- PASS` verdict and every `portable-tests.sh: four
# numbers` line came out of the delegated runner's own run. The exit code is 0,
# matching the real run this capture was taken from.
cat "${PC70_REAL_CHILD:?PC70_REAL_CHILD must name the captured real child output}" || exit 2
exit 0
