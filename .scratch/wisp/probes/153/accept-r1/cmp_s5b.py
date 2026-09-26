import io, re, difflib

def rd(p):
    return io.open(p, 'rb').read().decode('utf-8')

base = 'D:/tmp/wisp153-accept-r1/dumps/'
raw = {}
for k in ('bec1727', '9a8766e', '6de3d1c'):
    t = rd(base + k + '.md')
    raw[k] = t
    print('%s: CRLF=%d  LF=%d  bare_CR=%d' % (k, t.count('\r\n'), t.count('\n'), t.count('\r')))

def lines(k):
    return raw[k].replace('\r\n', '\n').split('\n')

for k in ('bec1727', '9a8766e', '6de3d1c'):
    L = lines(k)
    blank = sum(1 for l in L if not l.strip())
    print('%s lines=%d blank=%d nonblank=%d' % (k, len(L), blank, len(L) - blank))

a = lines('9a8766e')
b = lines('6de3d1c')
A = [l for l in a[690:830] if l.strip()]
B = [l for l in b[345:419] if l.strip()]


def squash(ls):
    s = ' '.join(ls).replace('**', '')
    return re.sub(r'\s+', ' ', s).strip()


wa, wb = squash(A).split(' '), squash(B).split(' ')
sm = difflib.SequenceMatcher(None, wa, wb, autojunk=False)
print('\n=== section 5 word-level diff: 9a8766e(<) vs 6de3d1c(>) ===')
for tag, i1, i2, j1, j2 in sm.get_opcodes():
    if tag == 'equal':
        continue
    print('[-OLD] ' + ' '.join(wa[i1:i2]))
    print('[+NEW] ' + ' '.join(wb[j1:j2]))
print('old_tokens=%d new_tokens=%d' % (len(wa), len(wb)))
