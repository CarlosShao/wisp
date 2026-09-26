#!/usr/bin/env python
"""Mutation driver for the 153 acceptance: pristine copies + named mutations,
each one prints the landing line so a mutation can never silently no-op."""
import os
import re
import shutil
import subprocess
import sys

SNAP = 'D:/tmp/wisp153-accept-r1/snap-delivered'
PRISTINE = 'D:/tmp/wisp153-accept-r1/pristine'
FILES = ['internal/agent/compress.go', 'internal/agent/loop.go',
         'internal/agent/compress_trace_test.go']
OUT = 'D:/tmp/wisp153-accept-r1/out'


def rd(p):
    return open(p, encoding='utf-8', newline='').read()


def wr(p, s):
    open(p, 'w', encoding='utf-8', newline='').write(s)


def restore():
    for f in FILES:
        shutil.copyfile(os.path.join(PRISTINE, f.split('/')[-1]), os.path.join(SNAP, f))
    bad = []
    for f in FILES:
        if rd(os.path.join(SNAP, f)) != rd(os.path.join(PRISTINE, f.split('/')[-1])):
            bad.append(f)
    print('RESTORE %s' % ('OK (3/3 byte-identical to pristine)' if not bad else 'MISMATCH ' + str(bad)))
    return not bad


def drop_func(path, name):
    """Delete a top-level func (plus its immediately preceding contiguous comment
    block) line-wise from a Go file."""
    s = rd(path)
    lines = s.split('\n')
    sig = next((i for i, l in enumerate(lines)
                if l.startswith('func ' + name + '(t *testing.T) {')), None)
    if sig is None:
        raise SystemExit('ANCHOR NOT FOUND: func %s in %s' % (name, path))
    end = next(i for i in range(sig, len(lines)) if lines[i] == '}')
    start = sig
    while start > 0 and lines[start - 1].startswith('//'):
        start -= 1
    before = len(lines)
    kept = lines[:start] + lines[end + 1:]
    # collapse the doubled blank line the cut leaves behind
    out = '\n'.join(kept).replace('\n\n\n', '\n\n')
    wr(path, out)
    gone = ('func ' + name) not in out
    print('DROPPED %s from %s (lines %d -> %d, comment lines taken %d, func gone=%s, name mentions left=%d)'
          % (name, os.path.basename(path), before, len(out.split('\n')), sig - start, gone,
             out.count(name)))
    if not gone:
        raise SystemExit('drop failed: %s still present' % name)
    return True


def sub_once(path, old, new, label):
    s = rd(path)
    n = s.count(old)
    if n != 1:
        raise SystemExit('MUT %s: anchor %r occurs %d times (want exactly 1)' % (label, old, n))
    wr(path, s.replace(old, new))
    chk = rd(path)
    print('MUT %s landed: %s  (now %d occurrence(s) of the new text; file %d -> %d bytes)'
          % (label, path.replace(SNAP, ''), chk.count(new), len(s), len(chk)))
    return new in chk and old not in chk


def run(tag):
    p = subprocess.run([sys.executable, 'D:/tmp/wisp153-accept-r1/ruler153.py',
                        './internal/agent/', '1', OUT + '/' + tag],
                       cwd=SNAP, capture_output=True, text=True, encoding='utf-8', errors='replace')
    print('---- %s ----' % tag)
    print(p.stdout.strip())
    if p.stderr.strip():
        print('STDERR: ' + p.stderr.strip()[:600])


ACTIONS = {
    'A0-baseline-delivered': lambda: None,
    'A1-m5': lambda: sub_once(SNAP + '/internal/agent/compress.go',
                              '\tif rep.Ran {\n', '\tif c.Need(hist) {\n', 'M5'),
    'A2-m5-dropAC1test': lambda: (sub_once(SNAP + '/internal/agent/compress.go',
                                            '\tif rep.Ran {\n', '\tif c.Need(hist) {\n', 'M5'),
                                  drop_func(SNAP + '/internal/agent/compress_trace_test.go',
                                            'TestCompressionTraceSilentWhenNothingFoldableOverThreshold')),
    'A3-dropAC1test-only': lambda: drop_func(SNAP + '/internal/agent/compress_trace_test.go',
                                             'TestCompressionTraceSilentWhenNothingFoldableOverThreshold'),
    'A4-dropN2-looptag': lambda: sub_once(SNAP + '/internal/agent/loop.go',
                                           'l.comp.Compress(withTraceTask(ctx, taskID), hist)',
                                           'l.comp.Compress(ctx, hist)', 'N2'),
    'A5-dropN2-dropowntaskidtest': lambda: (sub_once(SNAP + '/internal/agent/loop.go',
                                                     'l.comp.Compress(withTraceTask(ctx, taskID), hist)',
                                                     'l.comp.Compress(ctx, hist)', 'N2'),
                                            drop_func(SNAP + '/internal/agent/compress_trace_test.go',
                                                      'TestCompressionTraceCarriesTheOwningTaskID')),
    'A6-dropN3-attrwrite': lambda: sub_once(
        SNAP + '/internal/agent/compress.go',
        '\t\tif task := traceTaskID(ctx); task != "" {\n\t\t\tattrs = append(attrs, "task", task)\n\t\t}\n',
        '', 'N3'),
    'A7-dropN3-dropbothtasktests': lambda: (sub_once(
        SNAP + '/internal/agent/compress.go',
        '\t\tif task := traceTaskID(ctx); task != "" {\n\t\t\tattrs = append(attrs, "task", task)\n\t\t}\n',
        '', 'N3'),
        drop_func(SNAP + '/internal/agent/compress_trace_test.go',
                  'TestCompressionTraceCarriesTheOwningTaskID'),
        drop_func(SNAP + '/internal/agent/compress_trace_test.go',
                  'TestCompressionTraceNeverInventsATaskID')),
    'A8-m6-unconditional-attr': lambda: sub_once(
        SNAP + '/internal/agent/compress.go',
        '\t\tif task := traceTaskID(ctx); task != "" {\n\t\t\tattrs = append(attrs, "task", task)\n\t\t}\n',
        '\t\tattrs = append(attrs, "task", traceTaskID(ctx))\n', 'M6-unconditional'),
    'A9-m7-constant-fake-id': lambda: sub_once(
        SNAP + '/internal/agent/compress.go',
        '\ts, _ := ctx.Value(traceTaskKey{}).(string)\n\treturn s\n',
        '\t_ = ctx.Value(traceTaskKey{})\n\treturn "00000000-0000-4000-8000-000000000153"\n',
        'M7-constant-id'),
}

if __name__ == '__main__':
    tag = sys.argv[1]
    if not restore():
        raise SystemExit('restore failed')
    ACTIONS[tag]()
    run(tag)
