# 表① per-package tagged census @ HEAD blob 86d89478
# ruler: git ls-tree -r --name-only 86d89478 | grep '_test$' ; per-file git show 86d89478:path | head -3 | grep -c 'go:build windows' ; cases = git show 86d89478:path | grep -c '^func Test' (tagged files only)
PACKAGE                                    TAGGED_FILES TAGGED_TOPLEVEL_CASES
.scratch/wisp/probes/147                        1      3
.scratch/wisp/probes/149                        1      1
.scratch/wisp/probes/152                        2      2
.scratch/wisp/probes/156                        1      1
.scratch/wisp/probes/156-accept/anchor          1      5
.scratch/wisp/probes/174/c1                     1      3
.scratch/wisp/probes/174/r1                     2      2
.scratch/wisp/probes/252/v1/overlay             1      2
.scratch/wisp/probes/255/r6/mutations           3     10
.scratch/wisp/probes/33/r11/mutations           2      2
cmd/wisp                                       42    168
internal/agent                                  1      1
internal/audio                                  2      9
internal/ball                                  11     45
internal/config                                 1      3
internal/memory                                 1      1
internal/models                                 3      3
internal/proc                                   7     23
internal/risk                                   6     26
internal/secret                                 3     17
internal/tools                                 10     27
internal/winsec                                17     54
