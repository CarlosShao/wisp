#!/usr/bin/env python
"""Adversarial-acceptance ruler for ticket 153: parse `go test -v` output into
four numbers + roster names. Only anchored '--- FAIL' lines count as red."""
import re
import subprocess
import sys
import collections
import os

pkg = sys.argv[1]
count = sys.argv[2]
prefix = sys.argv[3]
extra = sys.argv[4:]

cmd = ['go', 'test', pkg, '-count=' + count, '-v'] + extra
p = subprocess.run(cmd, capture_output=True, text=True, encoding='utf-8', errors='replace')
out = (p.stdout or '') + (p.stderr or '')

run_all = re.findall(r'^\s*=== RUN\s+(\S+)', out, re.M)
top_pass = re.findall(r'^--- PASS: (\S+)', out, re.M)
sub_pass = re.findall(r'^\s+--- PASS: (\S+)', out, re.M)
fail = re.findall(r'^--- FAIL: (\S+)', out, re.M)
sub_fail = re.findall(r'^\s+--- FAIL: (\S+)', out, re.M)
skip = re.findall(r'^--- SKIP: (\S+)', out, re.M)
sub_skip = re.findall(r'^\s+--- SKIP: (\S+)', out, re.M)
panic = re.findall(r'^panic:', out, re.M)

names = collections.OrderedDict()
for n in top_pass + fail + skip:
    names[n.split('/')[0]] = None

print('CMD: ' + ' '.join(cmd))
print('rc=%d' % p.returncode)
print('RUN(all)=%d PASS_top=%d PASS_sub=%d FAIL_top=%d FAIL_sub=%d SKIP_top=%d SKIP_sub=%d panic=%d unique=%d'
      % (len(run_all), len(top_pass), len(sub_pass), len(fail), len(sub_fail),
         len(skip), len(sub_skip), len(panic), len(names)))
print('FAIL_NAMES: ' + ' '.join(fail + sub_fail))
print('SKIP_NAMES: ' + ' '.join(skip + sub_skip))
d = os.path.dirname(os.path.abspath(prefix))
os.makedirs(d, exist_ok=True)
with open(prefix + '.roster.txt', 'w', encoding='utf-8') as f:
    for n in sorted(names):
        f.write(n + '\n')
with open(prefix + '.log.txt', 'w', encoding='utf-8') as f:
    f.write(out)
