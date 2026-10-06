#!/bin/bash
# M6: aggregate passD + m5 into bucket candidates, print only short summaries.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
awk -F'|' '
{ n=$1; for(i=3;i<=NF;i++){ v=$i; sub(/^[a-z_]+=/,"",v); lab=$i; sub(/=.*/,"",lab); val[lab]=v }
  un[n]=val["unind"]; an[n]=val["anyind"]; dn[n]=val["done"]; bx[n]=val["boxes"]; pv[n]=val["pl_verdict"]; pi[n]=val["pl_impl"]; rf[n]=val["refmd"];
  name[n]=$2; }
END{
  printf "BOXES0: "; for(k in bx) if(bx[k]+0==0) printf "%s ",k; print "";
  printf "UNIND_GT0: "; for(k in un) if(un[k]+0>0) printf "%s(%s) ",k,un[k]; print "";
  printf "ANYIND_GT_UNIND: "; for(k in an) if(an[k]+0>un[k]+0) printf "%s(%s/%s) ",k,un[k],an[k]; print "";
  printf "PLVERDICT0: "; for(k in pv) if(pv[k]+0==0) printf "%s ",k; print "";
  printf "PLIMPL_GT5: "; for(k in pi) if(pi[k]+0>5) printf "%s(%s) ",k,pi[k]; print "";
  printf "REFMD0: "; for(k in rf) if(rf[k]+0==0) printf "%s ",k; print "";
  print "TOTALROWS="NR;
}' "$L/passD.tsv"
echo "== per-ticket evidence-name rollup (ACC/IMP/UNK/FILES/S1) =="
awk -F'|' '{printf "%s %s %s %s %s %s | %s\n",$1,substr($2,5),substr($3,5),substr($4,5),substr($5,6),substr($6,4),substr($7,1,26)}' "$L/m5.tsv"
