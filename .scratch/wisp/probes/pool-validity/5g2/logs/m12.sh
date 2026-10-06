#!/bin/bash
# M12: final join with robust numeric parsing (fixes substr truncation on 2-digit counts).
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
awk -F'|' '{print $1"|"$2}' "$L/m8.tsv" > "$L/m12-role.txt"
awk -F'|' '{ k=$1"|"$2; if(k in s) next; s[k]=1; tot[$1]++;
  f=$2; if(f ~ /adversarial|accept|verdict|recheck|裁决|验收|终裁|归口|backfill|gap-audit|缺口/) vf[$1]++;
  if(f ~ /s1/) s1[$1]++;
  if(f ~ /^docs\/evidence/) d[$1]++ }
END{ for(k in tot) print k" offv="vf[k]+0" offtot="tot[k]+0" offs1="s1[k]+0" oevd="d[k]+0 }' "$L/m7-offindex.txt" > "$L/m12-offsum.txt"
awk -F'|' '
function num(line, key,   s){ if(match(line, key"=[0-9]+")){ s=substr(line,RSTART,RLENGTH); sub(/.*=/,"",s); return s+0 } return -1 }
function num2(line, key,   s){ if(match(line, " "key"=[0-9]+")){ s=substr(line,RSTART+1,RLENGTH-1); sub(/.*=/,"",s); return s+0 } return -1 }
FILENAME ~ /m12-role/  { role[$1]=$2; next }
FILENAME ~ /m3-list/   { n=$1; f=$2; r=(f in role)?role[f]:"MISSING"; cnt[n]++;
                         if(r ~ /^ACC|^VLEG/) acc[n]++; else if(r ~ /^IMP|^RLEG/) imp[n]++; else unk[n]++;
                         if(r ~ /^ACC|^VLEG/) files[n]=files[n]" "substr(f,1,46); next }
FILENAME ~ /m12-offsum/{ split($0,z," "); offv[z[1]]=num2($0,"offv"); offtot[z[1]]=num2($0,"offtot"); offs1[z[1]]=num2($0,"offs1"); oevd[z[1]]=num2($0,"oevd"); next }
FILENAME ~ /passD/     { n=$1; F[n]=$2;
                         un[n]=num($0,"unind"); any[n]=num($0,"anyind"); dn[n]=num($0,"done"); bx[n]=num($0,"boxes");
                         pv[n]=num($0,"pl_verdict"); pi[n]=num($0,"pl_impl"); rf[n]=num($0,"refmd");
                         order[++M]=n; next }
END{ for(i=1;i<=M;i++){ n=order[i];
  a=acc[n]+0; im=imp[n]+0; u=unk[n]+0; c=cnt[n]+0; ov=offv[n]+0; ot=offtot[n]+0;
  if(a>0) t="A"; else if(ov>0) t="B"; else if(im>0||u>0) t="C"; else if(pv[n]>0||rf[n]>0) t="D"; else t="E";
  printf "%s|%s|tier=%s|nA=%d|nI=%d|nU=%d|nT=%d|offv=%d|offtot=%d|offs1=%d|un=%d|any=%d|done=%d|boxes=%d|pv=%d|pi=%d|ref=%d|%s\n",
    n, F[n], t, a, im, u, c, ov, ot, offs1[n]+0, un[n], any[n], dn[n], bx[n], pv[n], pi[n], rf[n], files[n] > OUT
}}' -v OUT="$L/m12-rows.txt" "$L/m12-role.txt" "$L/m3-list.txt" "$L/m12-offsum.txt" "$L/passD.tsv" 2>/dev/null
# awk -v must come before file list: rerun properly
rm -f "$L/m12-rows.txt" 2>/dev/null
awk -v OUT="$L/m12-rows.txt" -F'|' '
function num(line, key,   s){ if(match(line, key"=[0-9]+")){ s=substr(line,RSTART,RLENGTH); sub(/.*=/,"",s); return s+0 } return -1 }
function num2(line, key,   s){ if(match(line, " "key"=[0-9]+")){ s=substr(line,RSTART+1,RLENGTH-1); sub(/.*=/,"",s); return s+0 } return -1 }
FILENAME ~ /m12-role/  { role[$1]=$2; next }
FILENAME ~ /m3-list/   { n=$1; f=$2; r=(f in role)?role[f]:"MISSING"; cnt[n]++;
                         if(r ~ /^ACC|^VLEG/) acc[n]++; else if(r ~ /^IMP|^RLEG/) imp[n]++; else unk[n]++;
                         if(r ~ /^ACC|^VLEG/) files[n]=files[n]" "substr(f,1,44); next }
FILENAME ~ /m12-offsum/{ split($0,z," "); offv[z[1]]=num2($0,"offv"); offtot[z[1]]=num2($0,"offtot"); offs1[z[1]]=num2($0,"offs1"); next }
FILENAME ~ /passD/     { n=$1; F[n]=$2;
                         un[n]=num($0,"unind"); any[n]=num($0,"anyind"); dn[n]=num($0,"done"); bx[n]=num($0,"boxes");
                         pv[n]=num($0,"pl_verdict"); pi[n]=num($0,"pl_impl"); rf[n]=num($0,"refmd");
                         order[++M]=n; next }
END{ for(i=1;i<=M;i++){ n=order[i];
  a=acc[n]+0; im=imp[n]+0; u=unk[n]+0; c=cnt[n]+0; ov=offv[n]+0; ot=offtot[n]+0;
  if(a>0) t="A"; else if(ov>0) t="B"; else if(im>0||u>0) t="C"; else if(pv[n]>0||rf[n]>0) t="D"; else t="E";
  printf "%s|%s|tier=%s|nA=%d|nI=%d|nU=%d|nT=%d|offv=%d|offtot=%d|offs1=%d|un=%d|any=%d|done=%d|boxes=%d|pv=%d|pi=%d|ref=%d|%s\n",
    n, F[n], t, a, im, u, c, ov, ot, offs1[n]+0, un[n], any[n], dn[n], bx[n], pv[n], pi[n], rf[n], files[n] > OUT
}}' "$L/m12-role.txt" "$L/m3-list.txt" "$L/m12-offsum.txt" "$L/passD.tsv"
echo "M12=$(wc -l < "$L/m12-rows.txt")"
for T in A B C D E; do echo "TIER$T: $(awk -F'|' -v t="tier=$T" '$3==t{printf "%s ",$1}' $L/m12-rows.txt)"; done
echo "== offv>0 sorted =="; awk -F'|' '{for(i=1;i<=NF;i++) if($i ~ /^offv=/){split($i,x,"="); if(x[2]+0>0) printf "%s:%s ",$1,x[2]}}' $L/m12-rows.txt
