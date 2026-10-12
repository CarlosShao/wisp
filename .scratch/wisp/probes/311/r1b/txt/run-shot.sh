#!/usr/bin/env bash
# run-shot.sh - 311-r1b shot runner: one reading per call, self-labelled.
# usage: sh run-shot.sh OUT WORKDIR [ENV=V ...] -- CMD [ARGS...]
# Writes OUT as: header (workdir + env assignments + verbatim command line),
# then the run's stdout+stderr interleaved (2>&1), then a trailing "rc=N" line.
# Deliberately: no pipe into head (that would take head's exit code), no
# 2>/dev/null (a ruler's stderr must stay in the same file as its hit count).
out=$1
workdir=$2
shift 2
# Absolute-ise BEFORE any cd: the capture path used to be relative and the very
# first trip died with "No such file or directory" after the cd below (no reading
# was taken, nothing landed). Fixed here, not silently.
case $out in
    /* | [A-Za-z]:/* ) ;;
    *) out="$PWD/$out" ;;
esac
case $workdir in
    /* | [A-Za-z]:/* ) ;;
    *) workdir="$PWD/$workdir" ;;
esac
envs=()
while [ $# -gt 0 ] && [ "$1" != "--" ]; do
    envs+=("$1")
    shift
done
if [ "${1:-}" = "--" ]; then shift; fi
# Guard: the FIRST trip of this instrument lost its arguments to word splitting
# (the repo path contains a space) and no reading landed. Fail loudly here rather
# than write a capture into the wrong place or run the command in the wrong tree.
if [ ! -d "$workdir" ]; then
    echo "REFUSED workdir is not a directory: [$workdir]" >&2
    exit 96
fi
if [ ! -d "$(dirname "$out")" ]; then
    echo "REFUSED out dir does not exist: [$(dirname "$out")]" >&2
    exit 95
fi
part="$out.part"
{
    echo "# ruler=probes/311/r1b/txt/run-shot.sh ; shot taken by leg 311-r1b"
    echo "workdir=$workdir"
    if [ "${#envs[@]}" -gt 0 ]; then
        printf 'env_assignment=%s\n' "${envs[@]}"
    else
        echo "env_assignment=<none: factory default shape>"
    fi
    echo "command_line=$*"
    echo "---- interleaved stdout+stderr follows ----"
} > "$part"
cd "$workdir" || exit 97
if [ "${#envs[@]}" -gt 0 ]; then
    env "${envs[@]}" "$@" >> "$part" 2>&1
else
    env "$@" >> "$part" 2>&1
fi
rc=$?
echo "---- capture ends ----" >> "$part"
echo "rc=$rc" >> "$part"
mv "$part" "$out"
exit 0
