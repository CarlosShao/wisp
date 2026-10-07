# Extract a per-step roster from a GitHub Actions workflow YAML.
# Output columns: job-name | step-pos | line-no | if-guard | kind | name/run digest
# Pure text processing: NO go commands, NO CI calls.
BEGIN {
    job = ""
    injobs = 0
    pos = 0
    stepline = 0
    stepif = ""
    stepname = ""
    stepbody = ""
    stepkind = ""
}
/^jobs:/ { injobs = 1; next }
/^  [A-Za-z0-9_-]+:$/ && injobs {
    if (stepline) { printf "%s|%d|%d|%s|%s|%s\n", job, pos, stepline, (stepif ? stepif : "NO-GUARD"), stepkind, digest(stepname " || " stepbody) }
    stepline = 0
    job = $1
    sub(/:$/, "", job)
    pos = 0
    next
}
injobs && /^      - name: / {
    if (stepline) { printf "%s|%d|%d|%s|%s|%s\n", job, pos, stepline, (stepif ? stepif : "NO-GUARD"), stepkind, digest(stepname " || " stepbody) }
    pos++
    stepline = NR
    stepif = ""
    stepname = $0
    sub(/^      - name: /, "", stepname)
    gsub(/\|/, "/", stepname)
    stepbody = ""
    stepkind = ""
    next
}
injobs && /^        if: / {
    s = $0
    sub(/^        if: /, "", s)
    gsub(/\|/, "/", s)
    stepif = "IF@" NR ":" s
    next
}
injobs && /^        (run|uses|shell): / {
    s = $0
    sub(/^        /, "", s)
    gsub(/\|/, "/", s)
    if (s ~ /^uses:/) stepkind = "uses"
    stepbody = stepbody " " s
    next
}
END {
    if (stepline) { printf "%s|%d|%d|%s|%s|%s\n", job, pos, stepline, (stepif ? stepif : "NO-GUARD"), stepkind, digest(stepname " || " stepbody) }
}
function digest(t) {
    t = substr(t, 1, 300)
    gsub(/\|/, "/", t)
    return t
}
