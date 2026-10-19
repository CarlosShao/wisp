#!/usr/bin/env bash
# 342-v1: build the mandatory (i) scatter table as a single named artefact, every number
# recomputed from the RAW logs this leg wrote (⛔ copying a leg's reading, ⛔ copying the
# orchestrator's ledger). Each row: shape / bench / repeats / per-repeat numerators.
set -u
S=/d/tmp/wisp/20261018/342v1
L=$S/logs
R="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/342/v1/logs"
mkdir -p "$R"
grab () { # $1=file $2=shape label -> prints damage/truncated/deleted/flips
  grep -E "V342V1 $2: attempts" "$1" | sed -E 's/.*attempts=([0-9]+) (linkState)?flips=([0-9]+) flipErrs=([0-9]+) damage=([0-9]+) truncated=([0-9]+) deleted=([0-9]+).*/att=\1 flips=\3 damage=\5 trunc=\6 del=\7/'
}
{
echo "342-v1 scatter table (mandatory question (i)) -- every number recomputed from the RAW log named on its own line; the bench of this leg, built from git archive (⛔ the census leg's trees)."
echo
echo "== POSIX bench=golang:1.27 container, CGO_ENABLED=0 GOPROXY=off, host module cache ro ; A-config = GOTMPDIR bound to a D: dir (this leg's own bench defect, see verdict 7); B-config = matched to the leg's driver env (no bind) =="
echo
echo "-- VBASE (shipped by-name open, A-config, RAW 10-race-smoke.log / 10-race-r1..r3.log):"
for n in smoke r1 r2 r3; do printf "   %-6s " "$n"; grab "$L/10-race-$n.log" VBASE; done
echo "-- VPROD (production PrivateFile, A-config, same RAWs):"
for n in smoke r1 r2 r3; do printf "   %-6s " "$n"; grab "$L/10-race-$n.log" VPROD; done
echo "-- VJIA (甲 simulated landing, A-config):"
for n in smoke r1 r2 r3; do printf "   %-6s " "$n"; grab "$L/10-race-$n.log" VJIA; done
echo "-- VYI (乙 simulated landing, A-config):"
for n in smoke r1 r2 r3; do printf "   %-6s " "$n"; grab "$L/10-race-$n.log" VYI; done
echo "-- VBING (丙 report-after-the-act, A-config):"
for n in smoke r1 r2 r3; do printf "   %-6s " "$n"; grab "$L/10-race-$n.log" VBING; done
echo
echo "-- VBING re-run B-config, 2 repeats (RAW 20-vbingB-r1.log / 20-vbingB-r2.log):"
for n in 1 2; do printf "   r%-6s " "$n"; grab "$L/20-vbingB-r$n.log" VBING; done
echo
echo "-- VPROD racer-pacing sweep (attribution of the magnitude miss), B-config, attempts=1500 each (RAW 15-vprod-pace-{0,100,300,1000}.log):"
for us in 0 100 300 1000; do printf "   pace=%-5sus " "$us"; grab "$L/15-vprod-pace-$us.log" VPROD; done
echo
echo "-- landed production rows (甲/乙 IN the shipped route), B-config (RAW 16-landed-{jia,yi}-r{1,2}.log):"
for form in jia yi; do for n in 1 2; do printf "   %-4s r%-2s " "$form" "$n"; grab "$L/16-landed-$form-r$n.log" VPROD-LANDED; done; done
echo
echo "== Windows bench = this host, go1.27.1 windows/amd64, GOTMPDIR/TMPDIR/TEMP/TMP unset; linkStateFlips = flips proved reparseAt(ancestor)=true =="
echo "-- WBASE (by-name open, oNoFollowPlatform==0) RAW 31-win-race-r1..r3.log:"
for n in 1 2 3; do printf "   r%-3s " "$n"; grab "$L/31-win-race-r$n.log" WBASE; done
echo "-- WPROD (production PrivateFile):"
for n in 1 2 3; do printf "   r%-3s " "$n"; grab "$L/31-win-race-r$n.log" WPROD; done
echo "-- WJIA (甲's Windows face, widened reparse probe + by-name open):"
for n in 1 2 3; do printf "   r%-3s " "$n"; grab "$L/31-win-race-r$n.log" WJIA; done
echo "-- WYI (乙's Windows face, os.Root descent):"
for n in 1 2 3; do printf "   r%-3s " "$n"; grab "$L/31-win-race-r$n.log" WYI; done
echo
echo "== shape verdicts (量级/形状, ⛔ 逐字相同 requirement) =="
echo "reproduced: VBASE open (484/458/479 vs leg by-name 242 = same open shape, higher rate);"
echo "            VJIA/VYI damage=0 (truncating branch closed); VBING damage>0 (reported, did not prevent);"
echo "            WBASE 975/871/878 vs leg 897 = magnitude reproduced; WPROD 7/6/4 vs leg 3 = same tier;"
echo "            WYI 0/0/0 with root-refuse ~1341-1390 = reproduced; A3 deterministic delete-shape = reproduced 3/3."
echo "NOT reproduced, stated as a finding, no smoothing:"
echo "  (1) VPROD magnitude: leg 106 vs this leg 2/2/1 (A-config) and 56/18/32/28 (B-config pacing sweep)."
echo "      Attribution: racer tempo. This leg's racer lands ~4x the flips (1581 vs 358) and leaves the ancestor"
echo "      ABSENT on 1353/1500 attempts (ENOENT), so the truncating open itself rarely runs; under throttling"
echo "      (flips back to the leg's tier) damage returns to 18-56 and is then almost entirely DELETED, not truncated."
echo "  (2) WJIA damage: leg reported 0, this leg measures 3/5/5 with flips 59/56/63."
echo "      Reading: on Windows 甲's criterion does not close the truncating branch, it only shrinks it."
echo "  (3) leg's yinofb red-sentence shape: it claimed 4 reds naming openleaf_322_other.go:200 as a DAMAGED red;"
echo "      this leg's same control makes that nail red on the LIVENESS line instead (:206, RED instrument)."
} > "$R/22-scatter-table.txt"
wc -c "$R/22-scatter-table.txt"; tail -3 "$R/22-scatter-table.txt"
