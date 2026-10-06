import io, sys
lines = io.open(sys.argv[1], encoding='utf-8').read().split('\n')
start = int(sys.argv[2]); end = int(sys.argv[3])
for idx in range(start-1, end):
    s = lines[idx]
    if not s.startswith('|'):
        continue
    total = s.count('|')
    esc = s.count('\\|')
    print('line', idx+1, 'total', total, 'escaped', esc, 'cols', total-esc-1)
