#!/bin/bash
L="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/pool-validity/5g2/logs"
awk -F'|' '
function cf(f,k,  n,c){ n=0; while((getline line < f)>0){ if(line ~ "^"k"[|]"){ c=line; close(f); sub(/.*acc=/,"",c); sub(/[|].*/,"",c); return c } } close(f); return "NA" }
FILENAME ~ /m28-final/ { slug[$1]=$2; un[$1]=$3; any[$1]=$4; chk[$1]=$5; tot[$1]=$6; selfc[$1]=$8; order[++n]=$1; next }
' "$L/m28-final.txt"
awk -F'|' '
NR==FNR{ if($0 ~ /^m19/){} ; next }
' /dev/null /dev/null
# 1) tier counts
echo "TIERCOUNT:"; awk -F'|' '{print $3}' "$L/m19-table.tsv" | sort | uniq -c
# 2) accfiles counts from m19 col12
echo "ACCFILES:"; awk -F'|' '{split($12,a,"="); print a[2]}' "$L/m19-table.tsv" | sort -n | uniq -c
# 3) selfcells>0 set
echo "SELFCELLS_GT0:"; awk -F'|' '{split($8,a,"="); if(a[2]+0>0) printf "%s(%s) ",$1,a[2]}' "$L/m28-final.txt"; echo
# 4) tot != chk+un
echo "TOT_NE_CHK_PLUS_UN:"; awk -F'|' '{split($3,u,"=");split($4,v,"=");split($5,c,"=");split($6,t,"="); if(t[2]+0 != c[2]+0+u[2]+0) printf "%s(tot=%s,chk=%s,un=%s) ",$1,t[2],c[2],u[2]}' "$L/m28-final.txt"; echo
# 5) rows where tot == chk+un exactly
echo "TOT_EQ_COUNT:"; awk -F'|' '{split($3,u,"=");split($5,c,"=");split($6,t,"="); if(t[2]+0 == c[2]+u[2]+0) n++} END{print n}' "$L/m28-final.txt"
# 6) m31 NOSIG/NOACC set
echo "M31_NOSIG:"; awk -F'|' '$3 ~ /NOSIG|NOACC/{printf "%s ",$1}' "$L/m31-sign.txt"; echo
echo "M31_SIGGED:"; awk -F'|' '$3 !~ /NOSIG|NOACC/{printf "%s(%s) ",$1,substr($3,5,14)}' "$L/m31-sign.txt"; echo
# 7) m33 NO_CRED
echo "M33_NOCRED:"; awk -F'|' '$2 ~ /NO_CRED/{print $1}' "$L/m33-ident.tsv"
# 8) disagreement rows: build joined derive table
awk -F'|' '
FILENAME ~ /m33-ident/ { ident[$1]=($2=="NO_CRED"?"NO_CRED":"ident="$4); acc33[$1]=substr($3,5); next }
FILENAME ~ /m31-sign/  { s=substr($3,5); sig[$1]=(s=="NOSIG"?"NOSIG":(s=="NOACC"?"NOACC":s)); sigfile[$1]=$2; next }
FILENAME ~ /m19-table/ { tier[$1]=substr($3,6); acc[$1]=substr($12,10); sub(/ .*/,"",acc[$1]); rej[$1]=substr($14,10); retk[$1]=substr($13,6); next }
FILENAME ~ /m28-final/ {
  n++; k=$1; slug[k]=$2; un[k]=substr($3,4); any[k]=substr($4,5); chk[k]=substr($5,5); tot[k]=substr($6,5); sc[k]=substr($8,10);
  order[n]=k; next }
END{
  for(i=1;i<=n;i++){ k=order[i];
    d1 = (sig[k]=="NOACC"||sig[k]=="NOSIG") ? "no" : "yes";
    d2 = (ident[k] ~ /ident=[1-9]/) ? "yes" : "no";
    conf = (d1==d2) ? "" : "  ***CONFLICT";
    printf "%s|tier=%s|acc=%s|un=%s|any=%s|chk=%s|tot=%s|self=%s|m31=%s|m33=%s|%s%s\n", k, tier[k], acc[k], un[k], any[k], chk[k], tot[k], sc[k], sig[k], ident[k], sigfile[k], conf }
}' "$L/m33-ident.tsv" "$L/m31-sign.txt" "$L/m19-table.tsv" "$L/m28-final.txt" > "$L/5g3-derive.tsv"
echo "DERIVE_ROWS=$(wc -l < "$L/5g3-derive.tsv")"
echo "CONFLICT_SET:"; awk -F'|' '/CONFLICT/{printf "%s ",$1}' "$L/5g3-derive.tsv"; echo
echo "CONFLICT_COUNT:"; grep -c CONFLICT "$L/5g3-derive.tsv"
