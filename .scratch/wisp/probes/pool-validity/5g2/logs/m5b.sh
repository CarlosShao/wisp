#!/bin/bash
# M5b: join m4 role census + m3 filename list + passD ticket census -> classification table.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
awk -F'|' -v LS="$L/passD.tsv" '
NR==FNR && FILENAME ~ /m4\.tsv$/ {
  f=$1; ha=0; hi=0;
  for(i=2;i<=NF;i++){ if($i ~ /^hacc=/){split($i,x,"=");ha=x[2]} if($i ~ /^himp=/){split($i,y,"=");hi=y[2]} }
  role[f]=(ha>0?"ACC":(hi>0?"IMP":"UNK"));
  next
}
FILENAME ~ /m3-list\.txt$/ {
  n=$1; f=$2; if(f=="")next;
  r=role[f]; if(r=="")r="MISSING";
  cnt[n]++; acc[n]+=(r=="ACC"); imp[n]+=(r=="IMP"); unk[n]+=(r=="UNK");
  s1[n]+=($2 ~ /^docs\/evidence\/s1\//);
  files[n]=files[n]" "substr(f,1,52)"("r")";
  next
}
END{
  while((getline line < LS) > 0){
    split(line, b, "|"); if(b[1]=="")continue;
    t=b[1];
    printf "%s|ACC=%d|IMP=%d|UNK=%d|FILES=%d|S1=%d|%s|%s\n", t, acc[t]+0, imp[t]+0, unk[t]+0, cnt[t]+0, s1[t]+0, b[2], files[t];
  }
}' "$L/m4.tsv" "$L/m3-list.txt" > "$L/m5.tsv" 2>&1
echo "M5=$(wc -l < "$L/m5.tsv")"
echo "ZEROACC: "; awk -F'|' '$2=="ACC=0"{printf "%s ",$1}' "$L/m5.tsv"; echo
echo "UNKONLY(no ACC,no IMP): "; awk -F'|' '$2=="ACC=0" && $3=="IMP=0"{printf "%s ",$1}' "$L/m5.tsv"; echo
echo "IMPCOUNT>0: "; awk -F'|' '$3!="IMP=0"{printf "%s:%s ",$1,$3}' "$L/m5.tsv"; echo
echo "FILES=0: "; awk -F'|' '$5=="FILES=0"{printf "%s ",$1}' "$L/m5.tsv"; echo
