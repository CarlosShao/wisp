#!/bin/sh
gotmp_run="${WISP_TEST_GOTMP_ROOT:-"${LOCALAPPDATA//\\//}/wisp-test-gotmp"}/PROBE-$(date +%Y%m%d-%H%M%S)-$$"
echo "PROBE_LINE=[$gotmp_run]"
if mkdir -p "$gotmp_run" 2>/dev/null; then echo MKDIR_RC=0; else echo MKDIR_RC=1; fi
[ -d "$gotmp_run" ] && echo EXISTS=yes || echo EXISTS=no
