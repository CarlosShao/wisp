# pool-validity-2 / 批4 —— 零勾开放票有效性普查（号段 150 以上，46 枚）

> 只读普查腿 `pool-validity-2`，接 `pool-validity-1`（死腿）的活。
> 起手 HEAD `949a5b92`（2026-10-06 09:03 +0800）。⚠ 换 HEAD 要重量。
> 本腿**零 go 命令**（并行腿 `231-r1` 此刻真在 `cmd/wisp` 跑定向用例与整包读数）。
> 工具全集＝`git log`／`ls`／`grep`／`wc`／`find`／Read。四档尺与本批名册的表头见 `batch1.md` §0。
> 票面零改动；唯一写入路径＝`.scratch/wisp/probes/pool-validity/2/**`；`pool-validity/1/**` 只读只引。
> ⛔ 零读零引 `frontend/**`、`design/**`、`cmd/wisp/**`、`scripts/**`、`.github/**`、
> `docs/reports/HANDOVER.md`、`docs/evidence/s1/**`、`.scratch/wisp/probes/{231,268,111,evidence-close}/**`。

## §1 名册与判定（46 枚）

| 票号 | 票文件 | 上次动过 | 判定 | 复法（跑的那把尺） | 读数 |
|---|---|---|---|---|---|
| 150 | `150-four-deferred-markers-have-no-row-in-the-spec-12-registry-so-the-1-to-1-rule-is-violated-today-and-no-instrument-scans-it.md` | c6c6a499 09-26 | 未判 | — | — |
| 159 | `159-record-level-cleanup-c1-c5-left-by-ticket-156-acceptance.md` | d0b72775 09-26 | 未判 | — | — |
| 165 | `165-plan-first-then-act-a-brake-the-voice-entry-needs.md` | a651fe80 09-27 | 未判 | — | — |
| 166 | `166-no-session-layer-list-search-rename-archive-fork-on-the-go-side.md` | f1ce2f2b 09-26 | 未判 | — | — |
| 168 | `168-built-in-browser-as-a-disable-able-plugin-two-rules-to-set-first.md` | a651fe80 09-27 | 未判 | — | — |
| 169 | `169-a-single-decorative-glyph-in-frontend-makes-the-repo-wide-gate-red-and-kills-every-ci-step-after-it.md` | 82476e85 09-26 | 未判 | — | — |
| 170 | `170-the-cli-legs-l1-block-window-has-nobody-to-veto-it-asks-and-then-executes.md` | 27f798ec 09-26 | 未判 | — | — |
| 172 | `172-no-ruler-links-the-frozen-tier-in-plan-md-to-the-tier-the-gate-actually-gives-a-zero-tests-go-red.md` | 9b5f9e9d 09-27 | 未判 | — | — |
| 173 | `173-the-audit-line-prints-in-allowlist-scope-true-when-nothing-was-judged-blank-path-both-meanings-one-field.md` | d0c164a1 09-27 | 未判 | — | — |
| 178 | `178-pair-census-ruler-g6-negative-leg-cannot-express-closing-a-scope-opened-elsewhere-and-rings-every-honest-test.md` | 386716a6 09-28 | 未判 | — | — |
| 182 | `182-the-task-monitor-rail-is-outside-ticket-145s-fourteen-row-table-so-nobody-has-counted-which-of-its-stacks-have-a-go-side-source.md` | fccaf3e3 09-28 | 未判 | — | — |
| 185 | `185-every-successful-reread-mints-the-blocker-for-the-next-one-because-the-fs-read-result-becomes-an-undeclared-c25-mark.md` | 1f22f1aa 09-28 | 未判 | — | — |
| 186 | `186-the-composer-area-should-detect-git-and-let-the-user-switch-between-local-worktree-and-branch.md` | 38747655 09-28 | 未判 | — | — |
| 187 | `187-the-model-and-thinking-tier-in-the-composer-row-must-be-changeable-not-just-displayed.md` | 03c01d71 09-28 | 未判 | — | — |
| 189 | `189-the-review-stack-needs-uncommitted-count-and-per-file-diff-but-go-side-has-zero-source.md` | fccaf3e3 09-28 | 未判 | — | — |
| 194 | `194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md` | 74524fe5 09-28 | 未判 | — | — |
| 195 | `195-config-set-needs-a-per-section-write-foundation-before-the-panel-can-touch-it-and-the-only-writer-in-the-repo-is-permmode.md` | f25f9572 10-02 | 未判 | — | — |
| 196 | `196-two-task-state-vocabularies-already-coexist-and-the-schema-has-no-check-so-the-d43-names-are-not-the-ones-in-use.md` | fccaf3e3 09-28 | 未判 | — | — |
| 200 | `200-nothing-reads-the-projects-instructions-for-ai-although-all-five-harnesses-do.md` | 85da3060 09-28 | 未判 | — | — |
| 211 | `211-subagent-pool-cannot-exceed-the-d38d-tool-ceiling.md` | 8e309d22 09-29 | 未判 | — | — |
| 213 | `213-composer-plus-menu-slash-command-catalog.md` | 8e309d22 09-29 | 未判 | — | — |
| 214 | `214-composer-plus-menu-attachments-and-context-wired.md` | 485d6d25 09-29 | 未判 | — | — |
| 215 | `215-composer-plus-menu-skills-plugins-inventory.md` | c1fa2e1d 09-28 | 未判 | — | — |
| 216 | `216-menu-display-base-scrub-control-chars-and-two-sentences.md` | c1fa2e1d 09-28 | 未判 | — | — |
| 219 | `219-approval-card-three-reply-buttons-and-a-reason-box.md` | c45d9b6f 09-29 | 未判 | — | — |
| 220 | `220-roster-blockedOnApproval-only-sees-L2-and-second-card-silently-drops.md` | 6cf2cc60 09-29 | 未判 | — | — |
| 225 | `225-deferred-markers-do-not-match-the-spec12-registry-and-nothing-checks-it.md` | 326fbf78 10-04 | 未判 | — | — |
| 227 | `227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md` | fac60ad4 09-29 | 未判 | — | — |
| 229 | `229-subagent-spawns-an-unregistered-goroutine-name-and-drowns-the-leak-alarm.md` | c774f8da 09-29 | 未判 | — | — |
| 230 | `230-four-cells-left-unfinished-inside-closed-tickets.md` | 5759d7fb 09-29 | 未判 | — | — |
| 231 | `231-a-config-written-by-a-newer-build-is-reported-to-the-operator-as-a-validation-failure.md` | fac60ad4 09-29 | 未判 | — | — |
| 232 | `232-the-restart-tier-test-pins-that-the-sentence-exists-not-that-it-carries-the-why.md` | fac60ad4 09-29 | 未判 | — | — |
| 233 | `233-d36-three-tier-enum-is-tagged-on-no-key-and-the-reload-tier-has-neither-producer-nor-consumer.md` | fac60ad4 09-29 | 未判 | — | — |
| 234 | `234-the-twelve-cells-that-only-need-a-reading-mutation-and-gate-rerun-on-current-head.md` | 6d22dd6c 09-29 | 未判 | — | — |
| 236 | `236-six-cells-that-only-surface-at-the-reading-layer.md` | c308be29 10-05 | 未判 | — | — |
| 237 | `237-race-run-exposes-a-parentheft-inference-that-never-held-plus-two-things-only-a-human-reads.md` | d5a59d66 09-29 | 未判 | — | — |
| 238 | `238-first-run-on-a-fresh-machine-creates-the-whole-data-root-with-no-seal-to-inherit.md` | cf03fa31 09-30 | 未判 | — | — |
| 242 | `242-the-grant-binding-layer-and-the-panel-facing-read-surface-both-have-zero-rulers.md` | 745173ad 10-03 | 未判 | — | — |
| 244 | `244-spec-11-wants-the-gui-subsystem-build.md` | 378e9ca6 10-05 | 未判 | — | — |
| 247 | `247-the-capture-stack-has-zero-importers-so-no-real-microphone-level-ever-reaches-the-ball.md` | 8d17e8cf 09-30 | 未判 | — | — |
| 249 | `249-the-resident-process-installs-no-panic-sink-so-a-real-crash-leaves-no-stack-on-disk.md` | 05f91799 10-01 | 未判 | — | — |
| 253 | `253-panel-inbound-three-ruler-holes.md` | 46079fcc 10-04 | 未判 | — | — |
| 262 | `262-tracked-path-length-gate-for-ci-checkout.md` | da929f2b 10-04 | 未判 | — | — |
| 264 | `264-error-text-carries-absolute-paths-into-model-visible-text.md` | 275864f1 10-04 | 未判 | — | — |
| 266 | `266-no-external-tooth-over-the-slo-check-self-locking-nails.md` | a32f30fb 10-04 | 未判 | — | — |
| 269 | `269-gofmt-step-is-held-red-by-its-own-committed-mutation-sample-and-eats-three-unprotected-steps.md` | c308be29 10-05 | 未判 | — | — |

## §2 本批计数

仍成立 **未判** ／ 已失效 **未判** ／ 差翻勾 **未判** ／ 量不到 **未判** ／ 合计 46
