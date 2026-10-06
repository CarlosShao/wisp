#!/bin/bash
# M9c: robust final join -> one compact row per roster ticket with tier + evidence counts.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
awk -F'|' '{print $1"|"$2}' "$L/m8.tsv" > "$L/m9-role2.txt"
awk -F'|' '$2=="OFFINDEX"{ if(!(($1"|"$3) in seen)){ seen[$1"|"$3]=1; tot[$1]++;
  f=$3; if(f ~ /adversarial|accept|verdict|recheck|audit|裁决|验收|终裁|归口|backfill|gap|evidence|s1/) vf[$1]++ } }
END{ for(k in tot) print k" offv="vf[k]+0" offtot="tot[k] }' "$L/m7-offindex.txt" > "$L/m9-offsum.txt"
awk -F'|' -v OUT="$L/m9-rows.txt" '
function gnum(s, key,   v){ v=s; sub(".*" key "=", "", v); sub("[^0-9].*", "", v); return v+0 }
FILENAME ~ /m9-role2/ { role[$1]=$2; next }
FILENAME ~ /m3-list/  { n=$1; f=$2; r=(f in role)?role[f]:"MISSING"; cnt[n]++;
                        if(r ~ /^ACC|^VLEG/) acc[n]++; else if(r ~ /^IMP|^RLEG/) imp[n]++; else unk[n]++; next }
FILENAME ~ /m9-offsum/{ split($0,a," "); offv[a[1]]=substr(a[2],5)+0; offtot[a[1]]=substr(a[3],9)+0; next }
FILENAME ~ /passD/ {
  n=$1; F[n]=$2;
  un[n]=gnum($0,"unind"); any[n]=gnum($0,"anyind"); dn[n]=gnum($0,"done"); bx[n]=gnum($0,"boxes");
  pv[n]=gnum($0,"pl_verdict"); pi[n]=gnum($0,"pl_impl"); rf[n]=gnum($0,"refmd");
  order[++M]=n; next }
END{ for(i=1;i<=M;i++){ n=order[i];
  a=acc[n]+0; im=imp[n]+0; u=unk[n]+0; c=cnt[n]+0; ov=offv[n]+0; ot=offtot[n]+0;
  if(a>0) t="A"; else if(ov>0) t="B"; else if(im>0||u>0) t="C"; else if(pv[n]>0||rf[n]>0) t="D"; else t="E";
  printf "%s|%s|tier=%s|nA=%d|nI=%d|nU=%d|nT=%d|offv=%d|offtot=%d|un=%d|any=%d|done=%d|boxes=%d|pv=%d|pi=%d|refmd=%d\n",
    n, F[n], t, a, im, u, c, ov, ot, un[n], any[n], dn[n], bx[n], pv[n], pi[n], rf[n] > OUT
}}' "$L/m9-role2.txt" "$L/m3-list.txt" "$L/m9-offsum.txt" "$L/passD.tsv"
echo "M9ROWS=$(wc -l < "$L/m9-rows.txt")"
awk -F'|' '{for(i=1;i<=NF;i++) if($i ~ /^tier=/) print substr($i,6)}' "$L/m9-rows.txt" | sort | uniq -c
echo "TIER_A=$(awk -F'|' '$3=="tier=A"{printf "%s ",$1}' $L/m9-rows.txt)"
echo "TIER_B=$(awk -F'|' '$3=="tier=B"{printf "%s ",$1}' $L/m9-rows.txt)"
echo "TIER_C=$(awk -F'|' '$3=="tier=C"{printf "%s ",$1}' $L/m9-rows.txt)"
echo "TIER_D=$(awk -F'|' '$3=="tier=D"{printf "%s ",$1}' $L/m9-rows.txt)"
echo "TIER_E=$(awk -F'|' '$3=="tier=E"{printf "%s ",$1}' $L/m9-rows.txt)"
