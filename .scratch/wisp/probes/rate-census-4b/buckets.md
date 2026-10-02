# rate-census-4b — 号段 50–99 净勾选格普查

## §0 起手锚

- 分支：`dev`
- HEAD sha：`6fe300c5035a60586e681f586a4ead5e2b56137b`
- 取数时刻：`2026-10-02 08:46:45 +0800`
- 本腿角色：**只读普查**。全程未跑 `go test` / `go build` / `go vet`，未执行任何 exe，未跑任何 CPU 负载工具；只读文件与 git 历史。
- 票面零改动：本腿不碰任何票的 AC 勾选框、不改任何票名、不写 `docs/reports/pending-and-issues.md`、不动任何冻结件。
- 唯一写入路径：`.scratch/wisp/probes/rate-census-4b/`。
- 零读取目录：`frontend/**`、`design/**`（本段 62/64/65/68/77/96 那几枚涉及界面与球视觉的票，判据**只看票面正文**，未去界面侧核对做没做）。
- 跨号段纪律：50 以下与 100 以上的票**未读票面正文**，§4 只按编排者题面具名点出前置，不代其计数。

## §1 分母实测

复跑命令（在 `.scratch/wisp/issues/` 内）：

```
ls | grep -E "^[5-9][0-9]-" | grep -v -- "-done" | wc -l   ->  24
ls | grep -E "^[5-9][0-9]-" | grep -- "-done"    | wc -l   ->  26
ls | grep -E "^[5-9][0-9]-"                      | wc -l   ->  50
```

- **号段 50–99 共 50 枚**，`seq 50 99` 逐号扫无缺号。
- **已 `-done` 26 枚**，**开放 24 枚**。
- **与编排者题面的 24 一致。**
- `voided/` 之类子目录**不存在**（`ls -d */` 返回空），故分母无混入风险。
- 判完成**只认文件名 `-done` 后缀**；票面里的 `Status:` 行一律未采信。

开放名册（24 枚，文件名照抄）：

```
50-tier1-manifest-plugins.md
51-tier2-goja.md
52-d46-command-plugins.md
53-larkcli-plugin-e2e.md
54-s7-acceptance.md
55-s8-macos-port.md
56-s8-signing-distribution-naming.md
57-s8-plugin-sdk-registry.md
58-s8-i18n-english.md
59-p15-aec-spike.md
60-c32-realtime-engine.md
61-cloud-voice-providers-c9.md
62-liquid-glass-ball-visuals.md
64-ball-defects-hotkey-interactive.md
65-ball-glass-quality-rework.md
68-ball-default-visuals-parity.md
70-ci-actually-green.md
71-gates-must-self-report.md
72-atble-classification-runner.md
77-frontend-scaffold-reactbits-beautifului-shadcn.md
85-lint-tools-never-produced-a-verdict.md
86-resolvepercallbudget-wallclock-fragility.md
91-os-isolation-spike-restricted-token-vs-appcontainer.md
98-cmd-wisp-tests-never-run-on-this-host.md
```

## §2 四桶加总

（待填）

## §3 名册

（待填）

## §4 主链表

（待填）

## §5 本号段到能用 ≈ N 枚

（待填）

## §6 我可能分错的条目

（待填）

## §7 占位符自查

（待填）
