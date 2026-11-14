# 274-c1 · 读数件（逐字摘录＋时间戳对拉；2026-11-14）

本件只摘录，不判断（判断在 `docs/evidence/s1/274-c1-three-cells-verdict.md`）。全部行由 `gh run view <run> --job <job> --log` 现拉，完整日志落在 `logs/`（未入库，7.4 MB；重跑 §4 命令可复得）。

## 1. 两发绿 run 的逐 job 时长

| run | head | job | 起→止 (UTC) | 整 job |
|---|---|---|---|---|
| 43559883923 | b52365b71 | lint-frontend | 05:06:21→05:06:48 | 27s |
| 43559883923 | b52365b71 | test-windows | 05:06:21→05:27:28 | 1267s |
| 43559883923 | b52365b71 | slo-smoke | 05:06:20→05:08:39 | 139s |
| 43559883923 | b52365b71 | slo-full | 05:10:22→05:14:38 | 256s |
| 43561054718 | df6a8e035 | lint-frontend | 05:20:31→05:21:06 | 35s |
| 43561054718 | df6a8e035 | test-windows | 05:20:31→05:52:39 | 1928s |
| 43561054718 | df6a8e035 | slo-smoke | 05:20:31→05:23:25 | 174s |
| 43561054718 | df6a8e035 | slo-full | 05:20:31→05:24:54 | 263s |
| 37703959747（pre-甲 对照，整发红于无关原因） | cc31526 | lint-frontend | 23:44:32→23:44:56 | 24s |
| 同上 | cc31526 | test-windows | 23:44:32→23:54:39 | 607s |
| 同上 | cc31526 | slo-smoke | 23:44:32→23:46:46 | 134s |
| 同上 | cc31526 | slo-full | 00:13:53→00:15:09 | 76s |

## 2. `build.ps1` 步内前端段（三枚产 exe job，两发 run）

| run | job | node 行 | 快照行 | `frontend ok` 行 | npm ci | npm run build | 前端段合计 |
|---|---|---|---|---|---|---|---|
| 43559883923 | test-windows | 05:08:26.127 | 05:08:34.181 | 05:08:39.563 | 8.05s | 5.38s | **13.44s** |
| 43559883923 | slo-smoke | 05:07:11.413 | 05:07:22.162 | 05:07:31.898 | 10.75s | 9.74s | **20.49s** |
| 43559883923 | slo-full | 05:12:33.496 | 05:13:35.483 | 05:13:56.008 | 61.99s | 20.53s | **82.51s** |
| 43561054718 | test-windows | 05:23:42.772 | 05:23:56.235 | 05:24:04.467 | 13.46s | 8.23s | **21.70s** |
| 43561054718 | slo-smoke | 05:21:39.258 | 05:21:54.869 | 05:22:03.641 | 15.61s | 8.77s | **24.38s** |
| 43561054718 | slo-full | 05:23:03.776 | 05:23:53.406 | 05:24:13.486 | 49.63s | 20.08s | **69.71s** |

（`npm ci` ＝ node 行→快照行；`npm run build` ＝ 快照行→`frontend ok` 行。）

## 3. 逐字摘录

### 3.1 run 43559883923 · lint-frontend（job 132340969677）

```
lint-frontend	build (vite build -> frontend/dist, the bytes go:embed carries)	2026-11-14T05:06:42.4561626Z vite v8.3.0 building client environment for production...
lint-frontend	build (vite build -> frontend/dist, the bytes go:embed carries)	2026-11-14T05:06:42.7067362Z ✓ 2454 modules transformed.
lint-frontend	build (vite build -> frontend/dist, the bytes go:embed carries)	2026-11-14T05:06:42.8247869Z dist/index.html                   1.04 kB │ gzip:   0.63 kB
lint-frontend	build (vite build -> frontend/dist, the bytes go:embed carries)	2026-11-14T05:06:42.8248681Z dist/assets/index-CZoXrb9N.css   65.26 kB │ gzip:  12.28 kB
lint-frontend	build (vite build -> frontend/dist, the bytes go:embed carries)	2026-11-14T05:06:42.8249440Z dist/assets/index-zFGXlVs5.js   640.29 kB │ gzip: 193.92 kB
lint-frontend	build (vite build -> frontend/dist, the bytes go:embed carries)	2026-11-14T05:06:42.8327867Z ✓ built in 375ms
```

### 3.2 run 43561054718 · lint-frontend（job 132344591193）

```
✓ 2454 modules transformed.
dist/index.html                   1.04 kB │ gzip:   0.63 kB
dist/assets/index-CZoXrb9N.css   65.26 kB │ gzip:  12.28 kB
dist/assets/index-zFGXlVs5.js   640.29 kB │ gzip: 193.92 kB
✓ built in 442ms
```

### 3.3 `build.ps1` 的成功行（run 43559883923；三枚 job 与两发 run 逐字同）

```
build.ps1: frontend: node = C:\Program Files\nodejs\node.exe (v22.23.4); npm = C:\Program Files\nodejs\npm.cmd (10.9.9)
build.ps1: frontend: node_modules ABSENT -> running npm ci (clean install pinned by package-lock.json).
build.ps1: frontend: pre-build dist snapshot = 1 file(s) [.gitkeep] (this run's npm must produce an artifact name outside it).
build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1068 bytes): .gitkeep=0 assets\index-CZoXrb9N.css=65264 assets\index-DmXz6VZ6.js=640338 index.html=1068 || provenance: this run's npm added new artifact name(s) [assets\index-CZoXrb9N.css assets\index-DmXz6VZ6.js index.html] over the pre-build snapshot 1 file(s) [.gitkeep]; gone since snapshot: []
```

自托管那台（slo-full）的 node 行不同，其余同：

```
build.ps1: frontend: node = D:\work\server\node24\node.exe (v24.21.0); npm = D:\work\server\node24\npm.cmd (11.19.0)
```

### 3.4 `AC#12` 读数（test-windows，`panel_host_gate_test.go:118`）

改后（两发 run 各一行）：

```
    panel_host_gate_test.go:118: AC#12 reading (head b52365b7): shape=page-bundle built=true entry-bytes=1068 entry-ctype="text/html; charset=utf-8" entry-err=<nil> refs=2 check-err=<nil> manifest-entries=4 manifest-err=<nil> | git-metadata=true tracked=1 tracked-beyond-anchor=0 tracked-has-entry=false ignored-or-untracked=2
    panel_host_gate_test.go:118: AC#12 reading (head df6a8e03): shape=page-bundle built=true entry-bytes=1068 entry-ctype="text/html; charset=utf-8" entry-err=<nil> refs=2 check-err=<nil> manifest-entries=4 manifest-err=<nil> | git-metadata=true tracked=1 tracked-beyond-anchor=0 tracked-has-entry=false ignored-or-untracked=2
```

改前（pre-甲，run 37703959747，行号为当时那一版 `:109`）：

```
build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 Â§2.2)
    panel_host_gate_test.go:109: AC#12 reading (head cc31526): shape=anchor-only built=false entry-bytes=0 entry-ctype="" entry-err=panel: embedded assets are not built (run npm run build in frontend/) refs=0 check-err=panel: embedded assets are not built (run npm run build in frontend/) manifest-entries=0 manifest-err=panel: embedded assets are not built (run npm run build in frontend/) | git-metadata=true tracked=1 tracked-beyond-anchor=0 tracked-has-entry=false ignored-or-untracked=0
```

### 3.5 本机 `panel-assets`（只读；exe mtime 2026-11-02 17:16:22／32,966,840 B）

```
exit=0 197ms  panel assets embedded: 4 files, entry=index.html built=true
exit=0 150ms  panel assets check: entry=index.html built=true 2 asset refs resolve [./assets/index-DpRdJ-7r.js ./assets/index-Dxt2y5oW.css]
```

## 4. 复跑

见 `docs/evidence/s1/274-c1-three-cells-verdict.md` §4（同一串命令）。本目录下另有 `run-43559883923-jobs.json`／`run-43561054718-jobs.json`／`base-37703959747-jobs.json` 三枚机器读数（`gh run view --json headSha,conclusion,jobs` 原样）。
