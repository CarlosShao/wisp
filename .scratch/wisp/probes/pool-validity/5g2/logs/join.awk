BEGIN{ FS="|" }
FNR==1{ f=FILENAME }
f ~ /k-tier/    { tier[$1]=$2; next }
f ~ /k-nameidx/ { ni[$1]=$0; next }
f ~ /k-boxes/   { bx[$1]=$0; next }
f ~ /k-strict/  { st[$1]=$2; next }
f ~ /k-verdict/ { vd[$1]=$0; next }
f ~ /k-files/   { fn[$1]=$2; ord[++M]=$1; next }
f ~ /k-accfiles/{ acc[$1]=substr($0, index($0,"|")+1); next }
f ~ /k-impfiles/{ imp[$1]=substr($0, index($0,"|")+1); next }
END{ for(i=1;i<=M;i++){ n=ord[i];
  printf "%s|%s|tier=%s|%s|strict_unchecked=%s|%s|%s|accfiles=%s|impfiles=%s\n",
    n, fn[n], tier[n], ni[n], st[n], bx[n], vd[n], acc[n], imp[n] > OUT
} }
