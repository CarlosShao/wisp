# model-mirror fixtures (SPEC-11 §4, ticket 14)

Three tiers, one per C29 acceptance case (SPEC-04 §9). The mirror serves
BYTES ONLY (D33/F3): every hash lives in the in-repo signed manifest
(`models/manifest.json`), never here.

    good/     18081  correct bytes               -> downloads succeed
    corrupt/  18082  same paths, flipped bytes   -> sha256 must reject+delete
    missing/  18083  empty root                  -> 404, exercises failover

Layout under each tier follows the downloader's mirror convention
(`<mirror-base>/<model-id>/<file-name>`, see internal/models.BuildCandidates):

    fixtures/good/vad-fixture/vad.onnx          sha256 fb043a49e248bba9881648d9362abeeb78e129a2b1fcb4f5c73340317536c54d
    fixtures/good/kws-fixture/tiny-model.tar.bz2 sha256 fb043a49e248bba9881648d9362abeeb78e129a2b1fcb4f5c73340317536c54d
    (corrupt/ holds the same paths with one byte flipped; hashes in SHA256SUMS)

These tiny fixtures are the deterministic `internal/models/testdata` archive
and a 1KB prefix of it; unit tests use httptest servers, compose tiers exist
for dev envs and the S3 adversarial pass.

## Using real models in dev

Copy real artifacts into the matching manifest paths, e.g.:

    fixtures/good/vad-silero/silero_vad.onnx
    fixtures/good/punc-ct-transformer-zh-en-vocab272727/sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8.tar.bz2

Local copies already verified by T14 live in `third_party/model-fetch/`
(git-ignored). Never place a file here whose sha256 is not pinned in
`models/manifest.json` - a mirror file that fails the manifest hash is
rejected by design.
