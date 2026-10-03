#!/usr/bin/env python3
# 212-a3 fixed-ruler harvest v2 (ticket 212 AC#1). Zero Go commands.
# Fixes carried:
#   (F1) token sits AFTER the comment marker [a2 postscript 2], and the marker finder is
#        STRING-AWARE: // or # inside double-quoted/raw/backtick strings (e.g. https:// in
#        string literals) does NOT open a comment; /* open also string-aware.
#   (F2) ext alternation longest-first (jsonl before json, tsv before ts) [a2 postscript 1].
#   (F3) go block comments /* ... */ harvested, line-granular.
#   (F4) placeholders <> captured whole -> class 3.
#   (F5) leading '/' split: POSIX system dirs = 4 posix-abs; root-relative resolved on disk.
#   (F6) class 2 split in-index (worktree-deleted) vs true-absent; ../ compare slash-normalized.
#   (F7) tokens may start with a dot (.scratch/...), and '.' cannot precede a token start
#        (kills the scratch/ vs .scratch/ split).
#   (F8) URL shapes (scheme://...) stripped from comment regions before tokenizing (URLs are
#        out of AC#1 scope per a2 SS1; prevents https://x/y.md harvesting y.md).
#   (F9) resolution ladder for relative cites on disk miss: (a) repo-root, (b) .scratch/wisp/+token,
#        (c) unique tracked path-suffix -> 1 (suffix-hit); multiple -> 3 suffix-ambiguous.
#   (F10) stemless tokens like html/.jsonl -> 3 stemless (tokenizer artifact, not a real cite).
import subprocess, os, re, json
from collections import defaultdict

ROOT = os.getcwd()
WORK = os.path.join(ROOT, '.scratch', 'wisp', 'probes', '212', 'a3', 'work')

EXTS = ['md','go','sh','ps1','py','json','jsonl','sse','txt','log','toml','sql','yml','yaml',
        'csv','tsv','html','css','js','ts','tsx','jsx','mjs','png','jpg','svg','exe','db','mod','sum','golden']
EXTS = sorted(set(EXTS), key=len, reverse=True)   # F2
EXTALT = '|'.join(EXTS)
R_MAIN = re.compile(r'(?<![A-Za-z0-9_\-/.])((?:[A-Za-z]:[\\/])?(?:\.{1,2}/)?/?\.?[A-Za-z0-9_][A-Za-z0-9_./+\-<>]*\.(?:'+EXTALT+r'))')  # F7
R_BS = re.compile(r'(?<![A-Za-z0-9_.\\\-/])((?:[A-Za-z]:)?(?:\\\\)?(?:%[A-Za-z][A-Za-z0-9_]*%)?[A-Za-z0-9_.\\\-]+\.(?:'+EXTALT+r'))', re.IGNORECASE)
R_ELL = re.compile(r'(?<![A-Za-z0-9_\-/.])([A-Za-z0-9_][A-Za-z0-9_./+\-]*(?:\.\.\.|…)[A-Za-z0-9_./+\-]*\.(?:'+EXTALT+r'))')
R_WILD = re.compile(r'(?<![A-Za-z0-9_\-/.])([A-Za-z0-9_][A-Za-z0-9_./+\-]*\*[A-Za-z0-9_./+\-]*\.(?:'+EXTALT+r'))')
R_URL = re.compile(r'[A-Za-z][A-Za-z0-9+.\-]*://\S+')          # F8
PUNCT = '.,;:!?)]}>"\'、，。；：）」』》”’`*|'
DEN_EXT = ('.go', '.md', '.sh', '.ps1', '.py')
POSIX_SYS = ('tmp', 'bin', 'usr', 'sbin', 'etc', 'var', 'opt', 'home', 'Users', 'dev', 'proc', 'sys',
             'private', 'Library', 'Applications', 'mnt', 'media', 'root', 'srv', 'run')
STATS = {'url_stripped': 0, 'stemless_dropped': 0}

def git(args):
    return subprocess.run(['git'] + args, capture_output=True, text=True, encoding='utf-8',
                          errors='replace', cwd=ROOT).stdout

def strip_tok(t):
    t = t.strip()
    while t and t[-1] in PUNCT:
        t = t[:-1]
    return t

def first_comment_span(ln, hashes=False):
    """Return index of first comment marker outside string literals, or -1.
    Go: '//' or '/*'. sh/ps1/py (hashes=True): '#'. String-aware (F1)."""
    inDq = inBt = inSq = False
    i, n = 0, len(ln)
    while i < n:
        c = ln[i]
        if inDq:
            if c == '\\' and i + 1 < n:
                i += 2; continue
            if c == '"':
                inDq = False
        elif inBt:
            if c == '`':
                inBt = False
        elif inSq:
            if c == '\\' and i + 1 < n:
                i += 2; continue
            if c == "'":
                inSq = False
        else:
            if c == '"':
                inDq = True
            elif c == '`':
                inBt = True
            elif c == "'" :
                inSq = True
            elif hashes:
                if c == '#':
                    return i
            else:
                if c == '/' and i + 1 < n and ln[i+1] == '/':
                    return i
                if c == '/' and i + 1 < n and ln[i+1] == '*':
                    return -i - 2          # negative encoding: block open at ~(i+2)
        i += 1
    return -1

def comment_regions(path, text):
    """Yield (line, comment_text). Go block comments (F3): state machine, line-granular."""
    ext = path.rsplit('.', 1)[-1].lower()
    in_block = False
    for i, ln in enumerate(text.splitlines(), 1):
        if ext == 'md':
            yield i, ln
            continue
        if ext == 'go':
            if in_block:
                close = ln.find('*/')
                if close < 0:
                    yield i, ln
                    continue
                yield i, ln[:close]                 # tail inside block until */
                tail = ln[close+2:]
                if tail:
                    p = first_comment_span(tail)
                    if p >= 0:
                        yield i, tail[p+2:]
                    elif p <= -2:
                        j = -p - 2
                        yield i, tail[j+2:]
                        if '*/' not in tail[j+2:]:
                            in_block = True
                in_block = False
                continue
            p = first_comment_span(ln)
            if p >= 0:
                yield i, ln[p+2:]
            elif p <= -2:
                j = -p - 2
                seg = ln[j+2:]
                close = seg.find('*/')
                if close < 0:
                    yield i, seg
                    in_block = True
                else:
                    yield i, seg[:close]
                    tail = seg[close+2:]
                    if tail:
                        q = first_comment_span(tail)
                        if q >= 0:
                            yield i, tail[q+2:]
            continue
        # sh / ps1 / py
        if i == 1 and ln.startswith('#!'):
            continue
        p = first_comment_span(ln, hashes=True)
        if p >= 0:
            yield i, ln[p+1:]

def main():
    den_subset = [f for f in git(['ls-files', '--', 'cmd', 'internal', 'tools', 'scripts', 'docs']).splitlines() if f]
    all_tracked = [f for f in git(['ls-files']).splitlines() if f]
    TRACKED = set(all_tracked)
    TRACKED_LIST = all_tracked
    den = [f for f in den_subset if f.lower().endswith(DEN_EXT)]
    go_files  = [f for f in den if f.lower().endswith('.go')]
    md_files  = [f for f in den if f.lower().endswith('.md')]
    scr_files = [f for f in den if f.lower().endswith(('.sh', '.ps1', '.py'))]
    bmap = defaultdict(list)
    for f in all_tracked:
        bmap[os.path.basename(f)].append(f)
    bmap_allowed = {b: [f for f in v if not f.startswith(('design/', 'frontend/'))] for b, v in bmap.items()}

    occ = []
    seen = set()
    def emit(cls, citing, line, token, sub=''):
        key = (citing, line, token)
        if key in seen:
            return
        seen.add(key)
        occ.append({'cls': cls, 'citing': citing, 'line': line, 'token': token, 'sub': sub})

    def suffix_resolve(t2):
        cands = [f for f in TRACKED_LIST if f == t2 or f.endswith('/' + t2)]
        return cands

    def fixture_sub(t):
        if 'OneDrive' in t:
            return 'external-fixture'
        if re.match(r'^cwd-[A-Za-z]/', t):
            return 'runtime-fixture'
        return ''

    def resolve_relative(citing, line, t, p):
        """p = repo-relative path form of the token (no leading ./ or ../)."""
        if os.path.isfile(p):
            emit('1', citing, line, t); return
        alt = os.path.join('.scratch', 'wisp', p)
        if os.path.isfile(alt):
            emit('1', citing, line, t, 'scratch-prefix-hit'); return
        cands = suffix_resolve(p)                                  # F9c
        if len(cands) == 1:
            if os.path.isfile(cands[0]):
                emit('1', citing, line, t, 'suffix-hit'); return
            emit('2', citing, line, t, 'in-index' if cands[0] in TRACKED else 'true-absent'); return
        if len(cands) > 1:
            emit('3', citing, line, t, 'suffix-ambiguous'); return
        fs = fixture_sub(p)
        emit('2', citing, line, t, fs if fs else ('in-index' if p in TRACKED else 'true-absent'))

    def judge(citing, line, raw):
        t = strip_tok(raw)
        if not t or '.' not in t:
            return
        if re.search(r'/\.[A-Za-z0-9]+$', t):                      # F10 stemless artifact
            STATS['stemless_dropped'] += 1
            emit('3', citing, line, t, 'stemless'); return
        if re.match(r'(^|/)[A-Za-z0-9_][A-Za-z0-9_+\-]*\.[A-Za-z0-9]+/.+\.[A-Za-z0-9]+$', t, re.I):  # F11 dual-ext (go.mod/go.sum); hidden dirs like .scratch/ do NOT match
            emit('3', citing, line, t, 'multi-ext'); return
        if re.match(r'^[A-Za-z]:[\\/]', t) or t.startswith('\\\\') or t.startswith('%') or t.startswith('/'):
            if re.match(r'^[A-Za-z]:', t):
                emit('4', citing, line, t, 'drive'); return
            if t.startswith('\\\\'):
                emit('4', citing, line, t, 'unc'); return
            if t.startswith('%'):
                emit('4', citing, line, t, 'envvar'); return
            rt = t.lstrip('/')                                     # F5
            if any(rt == p or rt.startswith(p + '/') for p in POSIX_SYS):
                emit('4', citing, line, t, 'posix-abs'); return
            if os.path.isfile(rt):
                emit('1', citing, line, t, 'root-rel-hit'); return
            cands = suffix_resolve(rt)
            if len(cands) == 1:
                emit('1', citing, line, t, 'suffix-hit'); return
            emit('3', citing, line, t, 'root-ambiguous'); return
        if '...' in t or '…' in t:
            emit('3', citing, line, t, 'ellipsis'); return
        if '<' in t or '>' in t:
            emit('3', citing, line, t, 'placeholder'); return      # F4
        if '*' in t:
            emit('3', citing, line, t, 'wildcard'); return
        if '/' in t:
            if t.startswith('../'):
                res = os.path.normpath(os.path.join(os.path.dirname(citing), t))
                if res.startswith('..'):
                    emit('4', citing, line, t, 'escapes-root'); return
                res = res.replace(os.sep, '/')
                if os.path.isfile(res):
                    emit('1', citing, line, t); return
                p2 = t
                while p2.startswith('../'):
                    p2 = p2[3:]
                cands = suffix_resolve(p2)
                if len(cands) == 1:
                    emit('1', citing, line, t, 'suffix-hit'); return
                if len(cands) > 1:
                    emit('3', citing, line, t, 'suffix-ambiguous'); return
                fs = fixture_sub(p2)
                emit('2', citing, line, t, fs if fs else 'true-absent')
                return
            p = t[2:] if t.startswith('./') else t
            resolve_relative(citing, line, t, p)                   # F9
            return
        cands = bmap_allowed.get(t, [])
        if not cands:
            if bmap.get(t):
                emit('3', citing, line, t, 'bare-only-design-frontend')
            else:
                alt = os.path.join('.scratch', 'wisp', t)
                if os.path.isfile(alt):
                    emit('1', citing, line, t, 'scratch-prefix-hit')
                else:
                    emit('3', citing, line, t, 'bare-unresolvable')
            return
        if len(cands) > 1:
            emit('3', citing, line, t, 'bare-ambiguous'); return
        if os.path.isfile(cands[0]):
            emit('1', citing, line, t)
        else:
            emit('2', citing, line, t, 'in-index' if cands[0] in TRACKED else 'true-absent')

    for f in den:
        try:
            with open(f, 'r', encoding='utf-8', errors='replace') as fh:
                text = fh.read()
        except OSError:
            continue
        ext = f.rsplit('.', 1)[-1].lower()
        for line, region in comment_regions(f, text):
            if R_URL.search(region):                               # F8
                STATS['url_stripped'] += 1
                region = R_URL.sub(' ', region)
            if ext != 'md':
                for m in R_BS.findall(region):
                    t = strip_tok(m)
                    if t and ('\\' in t or t.startswith('%') or re.match(r'^[A-Za-z]:', t)):
                        if re.match(r'^[A-Za-z]:[\\/]', t) or t.startswith('\\\\') or t.startswith('%'):
                            sub = 'drive' if re.match(r'^[A-Za-z]:', t) else ('unc' if t.startswith('\\\\') else 'envvar')
                            emit('4', f, line, t, sub)
                for m in R_ELL.findall(region) + R_WILD.findall(region):
                    t = strip_tok(m)
                    if t:
                        judge(f, line, t)
            for m in R_MAIN.findall(region):
                judge(f, line, m)

    with open(os.path.join(WORK, 'classified.tsv'), 'w', encoding='utf-8') as fh:
        fh.write('class\tsub\tciting\tline\ttoken\n')
        for o in occ:
            fh.write(f"{o['cls']}\t{o['sub']}\t{o['citing']}\t{o['line']}\t{o['token']}\n")
    for cls in ('2', '3', '4'):
        with open(os.path.join(WORK, f'roster{cls}.tsv'), 'w', encoding='utf-8') as fh:
            for o in occ:
                if o['cls'] == cls:
                    fh.write(f"{o['citing']}:{o['line']}\t{o['token']}\t{o['sub']}\n")

    def cnt(pred):
        return sum(1 for o in occ if pred(o))
    narrow = [o for o in occ if o['cls'] in ('2', '3')
              and o['citing'].endswith('.go')
              and (o['citing'].startswith('internal/') or o['citing'].startswith('cmd/'))
              and o['token'].endswith('.md')]
    counts = {
        'denominator': {'go': len(go_files), 'md': len(md_files), 'script': len(scr_files),
                        'den_total': len(den), 'all_tracked': len(all_tracked)},
        'citations': {'1': cnt(lambda o: o['cls'] == '1'),
                      '2': cnt(lambda o: o['cls'] == '2'),
                      '3': cnt(lambda o: o['cls'] == '3'),
                      '4': cnt(lambda o: o['cls'] == '4'),
                      'unique_1_tokens': len({o['token'] for o in occ if o['cls'] == '1'}),
                      'code_side': cnt(lambda o: not o['citing'].endswith('.md')),
                      'md_side': cnt(lambda o: o['citing'].endswith('.md')),
                      'citing_files_2': len({o['citing'] for o in occ if o['cls'] == '2'}),
                      'citing_files_3': len({o['citing'] for o in occ if o['cls'] == '3'}),
                      'citing_files_4': len({o['citing'] for o in occ if o['cls'] == '4'})},
        'sub1': {s: cnt(lambda o, s=s: o['cls'] == '1' and o['sub'] == s) for s in ('suffix-hit', 'scratch-prefix-hit', 'root-rel-hit')},
        'sub2': {s: cnt(lambda o, s=s: o['cls'] == '2' and o['sub'] == s) for s in
                 ('in-index', 'true-absent', 'external-fixture', 'runtime-fixture')},
        'sub3': {s: cnt(lambda o, s=s: o['cls'] == '3' and o['sub'] == s) for s in
                 ('ellipsis', 'wildcard', 'placeholder', 'stemless', 'multi-ext', 'root-ambiguous', 'suffix-ambiguous',
                  'bare-unresolvable', 'bare-ambiguous', 'bare-only-design-frontend')},
        'sub4': {s: cnt(lambda o, s=s: o['cls'] == '4' and o['sub'] == s) for s in
                 ('drive', 'unc', 'envvar', 'posix-abs', 'escapes-root')},
        'per_class_side': {c: {'code': cnt(lambda o, c=c: o['cls'] == c and not o['citing'].endswith('.md')),
                               'md': cnt(lambda o, c=c: o['cls'] == c and o['citing'].endswith('.md'))}
                           for c in ('1', '2', '3', '4')},
        'ruler_stats': dict(STATS),
        'ticket_caliber_slice': {'go_internal_cmd_md_token_2': sum(1 for o in narrow if o['cls'] == '2'),
                                 'go_internal_cmd_md_token_3': sum(1 for o in narrow if o['cls'] == '3'),
                                 'go_internal_cmd_docs_md_2': sum(1 for o in narrow if o['cls'] == '2' and o['token'].startswith('docs/')),
                                 'go_internal_cmd_docs_md_3': sum(1 for o in narrow if o['cls'] == '3' and o['token'].startswith('docs/'))},
    }
    with open(os.path.join(WORK, 'counts.json'), 'w', encoding='utf-8') as fh:
        json.dump(counts, fh, ensure_ascii=False, indent=1)
    print(json.dumps(counts, ensure_ascii=False, indent=1))

    # ---- calibration on the 4 known samples + 2 artifact checks -------
    cal = []
    def find(citing, tokfrag, expect):
        hits = [(o['cls'], o['sub']) for o in occ if o['citing'] == citing and tokfrag in o['token']]
        cal.append((citing, tokfrag, expect, hits))
    find('internal/tools/subagent_197_test.go', '197-subagent-entity-r1b.md', 'A:1-exists')
    find('cmd/wisp/subagent_stream_key_197_test.go', '197-subagent-entity-r1b.md', 'B:1-exists')
    find('cmd/wisp/slo_report_144_windows_test.go', '152-', 'C:3-ellipsis(+full-name 1)')
    find('cmd/wisp/config_reload.go', 'restart_tier_keys_255r2_test.go', 'D:1-exists')
    jl = sorted({o['token'] for o in occ if o['token'].endswith('.jsonl')})
    art1 = []
    try:
        with open('tools/d22scan/scan_test.go', 'r', encoding='utf-8', errors='replace') as fh:
            for i, ln in enumerate(fh.read().splitlines(), 1):
                if 'docs/' in ln:
                    j = ln.find('docs/')
                    k = ln.find('//')
                    if k < 0 or j < k:
                        art1.append((i, ln.strip()[:110]))
    except OSError:
        pass
    lit_lines = {i for i, _ in art1}
    leaked = [o for o in occ if o['citing'] == 'tools/d22scan/scan_test.go' and o['line'] in lit_lines]
    # string-aware spot check: URL in code must not be harvested as comment
    url_leak = [(o['citing'], o['line'], o['token']) for o in occ
                if o['citing'] == 'internal/risk/provenance_test.go' and 'example.com' in o['token']]
    with open(os.path.join(WORK, 'calibration.txt'), 'w', encoding='utf-8') as fh:
        for c in cal:
            fh.write(f"{c[1]}\texpect={c[2]}\tgot={c[3]}\n")
        fh.write(f"jsonl-unique-tokens={jl}\n")
        fh.write(f"artifact1-literal-lines(path-before-marker)={art1}\n")
        fh.write(f"artifact1-leak-assertion={'PASS' if not leaked else 'FAIL: '+str(leaked)}\n")
        fh.write(f"url-in-code-leak={'PASS' if not url_leak else 'FAIL: '+str(url_leak)}\n")
        fh.write(f"ruler-stats={STATS}\n")
    print('CALIBRATION:')
    for c in cal:
        print(f"  {c[1]}  expect={c[2]}  got={c[3]}")
    print(f"  jsonl unique tokens: {len(jl)} (first 5: {jl[:5]})")
    print(f"  artifact1 assertion: {'PASS' if not leaked else 'FAIL ' + str(leaked)}")
    print(f"  url-in-code leak assertion: {'PASS' if not url_leak else 'FAIL ' + str(url_leak)}")

if __name__ == '__main__':
    main()
