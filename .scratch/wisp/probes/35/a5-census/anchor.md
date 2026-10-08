# 35-a5-census 起手锚（只读腿，零 go 命令 / 零开窗 / 零 push）

## 锚
- `git log -1` 完整号：`dffb8062b765a32230d983b82b3c02d425a0c2b3`
  （Thu Oct 8 08:26:10 2026 +0800，标题「票 35 追加一节（10-08 08:2x 编排者收 35-r4 盘上遗产…）」）
- branch = `dev`
- 落笔前 index 状态：`git diff --cached --name-only | wc -l` = `0`（别人没有 staged 件；本腿全部 commit 带显式 pathspec）
- 尺日志：`logs/a0-anchor-git.txt`（含 `rc=0`）

## 编排者转述的三条读数，逐条复量
1. 「`.github/workflows/ci.yml` 里 `grep -n winlive` 零命中；`grep -n -- '-tags'` 也零命中」
   ⇒ 复量：`grep -n winlive .github/workflows/ci.yml` → 0 命中（`rc=1`，件 `logs/a1-ci-winlive.txt`）；
   `grep -n -- '-tags' .github/workflows/ci.yml` → 0 命中（`rc=1`，件 `logs/a1-ci-tags.txt`）。
   **〔与我转述同〕**（两条都同）。
2. 「`git log --oneline -S winlive -- .github/workflows/` 为空 ⇒ CI 历史上从未有过这一步」
   ⇒ 复量：输出 0 行（`rc=0`，件 `logs/a2-ci-history.txt`）。**〔与我转述同〕**。
3. 「go test 步骤走 `scripts/runtests.sh` 那形」
   ⇒ 复量：**`scripts/` 目录里根本没有 `runtests.sh`**（`ls scripts/` 全名册见 `logs/q2-scripts-ls.txt`：
   build.ps1 / check-path-length-budget.sh / d22scan.sh / dev / fetch-deps.ps1 / portable-tests-selftest.sh /
   portable-tests.sh / sign-models.ps1 / slo-check.ps1 / slo-freshness.sh / spike / testdata / winsec-tests.sh / wisp-cli-tests.sh）。
   **〔与我转述冲突〕** —— 这一条待本腿在 Q2 里用 ci.yml 现量行号补正 CI 真实的 go test 调用形状。

## 其它已量的地基
- `.gitignore:8` = `*.out`（全仓）⇒ 本腿证据件一律 `.txt`/`.tsv`/`.md`，不用 `.out`。
- `.github/workflows/` 只有两枚：`ci.yml`、`slo-fresh.yml`。
