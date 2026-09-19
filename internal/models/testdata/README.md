# testdata fixtures (internal/models)

## tiny-model.tar.bz2

Deterministic 3-file "model" archive in the k2-fsa release layout (single top
directory, flattened on extraction by `ExtractTarBz2`). Used by the archive
installation tests; the pinned SHA256 values live in `downloader_test.go`
(`tinyArchive*` constants) and must match this file byte for byte.

Regenerate (Python 3, deterministic mtime/content):

    py - <<'EOF'
    import io, tarfile, bz2, hashlib
    buf = io.BytesIO()
    with tarfile.open(fileobj=buf, mode='w') as tf:
        def add(name, data):
            ti = tarfile.TarInfo(name); ti.size = len(data); ti.mtime = 0
            tf.addfile(ti, io.BytesIO(data))
        add('tiny-model-2024-01-01/model.onnx', bytes(range(256))*4)
        add('tiny-model-2024-01-01/tokens.txt', b'token one two three\n')
        add('tiny-model-2024-01-01/dict/inner.txt', b'inner dict payload\n')
        add('tiny-model-2024-01-01/README.md', b'# tiny fixture model\n')
    comp = bz2.compress(buf.getvalue(), 9)
    open('tiny-model.tar.bz2','wb').write(comp)
    print(hashlib.sha256(comp).hexdigest(), len(comp))
    EOF

Then copy the printed hash/size into `downloader_test.go`.
