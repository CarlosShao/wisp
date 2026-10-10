build.ps1: repo root D:\work\workspace\projects plans\Wisp (Env=dev)
build.ps1: toolchain: go version go1.27.1 windows/amd64; CC=E:\work\base\msys64\mingw64\bin\gcc.exe (gcc.exe (Rev3, Built by MSYS2 project) 16.2.0)
fetch-deps: cache hit - third_party/sherpa-onnx matches deps.toml (sherpa-onnx 1.13.8)
build.ps1: frontend: node = D:\work\server\node14\node.exe (v24.9.0); npm = D:\work\server\node14\npm.cmd (11.6.0)
build.ps1: frontend: node_modules EXISTS -> taking the npm-ci-SKIPPED branch (npm run build only). This branch does not verify the installed tree against package-lock.json; delete frontend/node_modules to force a clean npm ci.
build.ps1: frontend: pre-build dist snapshot = 4 file(s) [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html] (this run's npm must produce an artifact name outside it).
build.ps1: frontend: running npm run build (tsc -b && vite build, output goes to frontend/dist).

> frontend@0.0.0 build
> tsc -b && vite build

[36mvite v8.3.0 [32mbuilding client environment for production...[36m[39m
transforming...
✓ 2439 modules transformed.
rendering chunks...
computing gzip size...
dist/index.html                   1.04 kB │ gzip:   0.63 kB
dist/assets/index-yy8KMgdf.css   49.54 kB │ gzip:   9.84 kB
dist/assets/index-B8yINMF1.js   551.98 kB │ gzip: 170.73 kB

[33m[plugin builtin:vite-reporter] 
(!) Some chunks are larger than 500 kB after minification. Consider:
- Using dynamic import() to code-split the application
- Use build.rolldownOptions.output.codeSplitting to improve chunking: https://rolldown.rs/reference/OutputOptions.codeSplitting
- Adjust chunk size limit for this warning via build.chunkSizeWarningLimit.[39m
[32m✓ built in 885ms[39m
build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1044 bytes): .gitkeep=0 assets\index-B8yINMF1.js=551989 assets\index-yy8KMgdf.css=49540 index.html=1044 || provenance: this run's npm added new artifact name(s) [assets\index-B8yINMF1.js assets\index-yy8KMgdf.css] over the pre-build snapshot 4 file(s) [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html]; gone since snapshot: [assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js]
build.ps1: go build ok (cgo linked against sherpa-onnx C API)
build.ps1: DLLs colocated into D:\work\workspace\projects plans\Wisp\build
build.ps1: signed model manifest colocated into D:\work\workspace\projects plans\Wisp\build\models
build.ps1: wrote D:\work\workspace\projects plans\Wisp\build\SHA256SUMS
build.ps1: smoke test - running wisp.exe doctor
build.ps1: done. Artifacts in build\: wisp.exe, onnxruntime.dll, sherpa-onnx-c-api.dll, sherpa-onnx-cxx-api.dll, SHA256SUMS
time=2026-10-10T08:51:18.694+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
wisp doctor - build chain self-check
------------------------------------------------------------------------
[INFO] wisp build                         version=0.0.0-dev commit=29081a13 built=2026-10-10T00:51:10Z WISP_ENV=dev
[INFO] Go runtime                         go1.27.1 (toolchain pinned by go.mod)
[INFO] C29 minisign public key            untrusted comment: wisp models signing key (dev)
RWSjyHlPP9lPxdEQRvWj3zFLMbc1tTEkKMwTDuVXXQDxsWRpA/m5jk9j (placeholder until C29 lands; hardcoded into buildinfo per SPEC-11 §7.3)
[PASS] gcc (build-time)                   gcc.exe (Rev3, Built by MSYS2 project) 16.2.0
[PASS] sherpa-onnx C API                  runtime 1.13.8 matches build pin (exe dir D:\work\workspace\projects plans\Wisp\build)
[PASS] onnxruntime.dll colocated          file version 1.28.2.0 matches build pin 1.28.2
[INFO] sherpa-onnx built against onnxruntime 1.28.2
[PASS] DLL colocated: sherpa-onnx-c-api.dll D:\work\workspace\projects plans\Wisp\build\sherpa-onnx-c-api.dll
[PASS] DLL colocated: sherpa-onnx-cxx-api.dll D:\work\workspace\projects plans\Wisp\build\sherpa-onnx-cxx-api.dll
[PASS] deps.toml sherpa-onnx pin          1.13.8
[PASS] deps.toml onnxruntime pin          1.28.2
[PASS] deps.toml sherpa-onnx-go pin       v1.13.8
[PASS] data dir writable (dev)            C:\Users\swq\AppData\Roaming\wisp-dev
------------------------------------------------------------------------
wisp doctor: PASS
