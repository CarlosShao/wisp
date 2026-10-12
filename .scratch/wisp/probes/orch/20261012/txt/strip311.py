import sys

START = "# ---- TICKET 311 AC#1"

def strip(path):
    with open(path, encoding="utf-8") as fh:
        text = fh.read()
    crlf = "\r\n" in text
    lines = text.split("\n")
    out = []
    i = 0
    removed = 0
    blocks = 0
    while i < len(lines):
        if lines[i].rstrip("\r").startswith(START):
            j = i
            while j < len(lines) and lines[j].rstrip("\r") != "fi":
                j += 1
            if j >= len(lines):
                raise SystemExit("UNBOUNDED BLOCK in %s" % path)
            blocks += 1
            removed += (j - i + 1)
            i = j + 1
            continue
        out.append(lines[i])
        i += 1
    body = "\r\n".join(out) if crlf else "\n".join(out)
    with open(path, "w", encoding="utf-8", newline="") as fh:
        fh.write(body)
    print("path=%s blocks=%d lines_removed=%d crlf=%s" % (path, blocks, removed, crlf))

for p in sys.argv[1:]:
    strip(p)
