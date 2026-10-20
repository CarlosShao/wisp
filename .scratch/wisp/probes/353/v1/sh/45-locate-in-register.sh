#!/bin/sh
S7=$(printf '\xc2\xa7')
cd '/d/work/workspace/projects plans/Wisp' || exit 9
grep -n "$S7" .scratch/wisp/probes/353/v1/00-ruler-register.md | cut -c1-200
