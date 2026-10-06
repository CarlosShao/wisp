#!/bin/bash
# M9b: pure-awk final join -> one compact row per roster ticket with tier + evidence counts.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
awk -F'|' -v OUT="$L/m9-rows.txt" -v LDIR="$L" '
function flushline(){ if(NR2=="")return }
FNR==1 { fileidx++ }
fileidx==1 { role[$2]=$3; next }                 # m9-role2.txt : file|role  (built below)
{ }
' /dev/null
# build simple role file first
awk -F'|' '{print $1"|"$2}' "$L/m8.tsv" | sed 's/|ROLE=/|/' > "$L/m9-role2.txt"
# build off-index per-ticket counts
awk -F'|' '$2=="OFFINDEX"{ if(!(($1"|"$3) in seen)){ seen[$1"|"$3]=1; tot[$1]++;
  f=$3; if(f ~ /adversarial|accept|verdict|recheck|audit|裁决|验收|终裁|归口|backfill|gap|evidence/) vf[$1]++ } }
END{ for(k in tot) print k" offv="(vf[k]+0)" offtot="tot[k] }' "$L/m7-offindex.txt" > "$L/m9-offsum.txt"
awk -F'|' -v OUT="$L/m9-rows.txt" '
FILENAME ~ /m9-role2/ { role[$1]=$2; next }
FILENAME ~ /m3-list/  { n=$1; f=$2; r=(f in role)?role[f]:"MISSING"; cnt[n]++;
                        if(r ~ /^ACC|^VLEG/) acc[n]++; else if(r ~ /^IMP|^RLEG/) imp[n]++; else unk[n]++; next }
FILENAME ~ /m9-offsum/{ split($0,a," "); offv[a[1]]=substr(a[2],5)+0; offtot[a[1]]=substr(a[3],8)+0; next }
FILENAME ~ /passD/ {
  n=$1; F[n]=$2; un[n]=substr($3,7); any[n]=substr($4,7); dn[n]=substr($5,6); bx[n]=substr($6,7);
  pv[n]=substr($7,11); pi[n]=substr($8,8); rf[n]=substr($9,6); order[++M]=n; next }
END{ for(i=1;i<=M;i++){ n=order[i];
  a=acc[n]+0; im=imp[n]+0; u=unk[n]+0; c=cnt[n]+0; ov=offv[n]+0; ot=offtot[n]+0;
  t="A"; if(a==0) t=(ov>0)?"B":"C"; if(a==0 && ov==0 && im==0 && u==0) t=(pi[n]>0||rf[n]>0)?"D":"E";
  printf "%s|%s|tier=%s|nA=%s|nI=%s|nU=%s|nT=%s|offv=%s|offtot=%s|un=%s|any=%s|done=%s|boxes=%s|pv=%s|pi=%s|refmd=%s\n",
    n, F[n], t, a, im, u, c, ov, ot, un[n], any[n], dn[n], bx[n], pv[n], pi[n], rf[n] > OUT
}}' "$L/m9-role2.txt" "$L/m3-list.txt" "$L/m9-offsum.txt" "$L/passD.tsv"
echo "M9ROWS=$(wc -l < "$L/m9-rows.txt")"
echo "== tier counts =="; awk -F'|' '{for(i=1;i<=NF;i++) if($i ~ /^tier=/) print substr($i,6)}' "$L/m9-rows.txt" | sort | uniq -c
echo "== tier B rows (evidence only under other ticket numbers) =="; grep '|tier=B|' "$L/m9-rows.txt" | cut -c1-120
echo "== tier C/D/E rows =="; grep -E '\|tier=[CDE]\|' "$L/m9-rows.txt" | cut -c1-160
