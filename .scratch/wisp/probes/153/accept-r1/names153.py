#!/usr/bin/env python
"""Re-parse saved go-test logs into the same roster shape the 153 impl used:
every anchored '--- PASS|FAIL|SKIP' name (subtests included), de-duplicated."""
import re
import sys

for path in sys.argv[1:]:
    txt = open(path, encoding='utf-8', errors='replace').read()
    top_pass = re.findall(r'^--- PASS: (\S+)', txt, re.M)
    sub_pass = re.findall(r'^\s+--- PASS: (\S+)', txt, re.M)
    fail = re.findall(r'^--- FAIL: (\S+)', txt, re.M) + re.findall(r'^\s+--- FAIL: (\S+)', txt, re.M)
    skip = re.findall(r'^--- SKIP: (\S+)', txt, re.M) + re.findall(r'^\s+--- SKIP: (\S+)', txt, re.M)
    names = sorted(set(top_pass + sub_pass + fail + skip))
    out = path.replace('.log.txt', '.names.txt')
    with open(out, 'w', encoding='utf-8') as f:
        f.write('\n'.join(names) + '\n')
    print('%s -> total=%d unique=%d (top_pass=%d sub_pass=%d fail=%d skip=%d)  names=%s'
          % (path, len(top_pass) + len(sub_pass) + len(fail) + len(skip), len(names),
             len(top_pass), len(sub_pass), len(fail), len(skip), out))
