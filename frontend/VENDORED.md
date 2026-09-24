# Vendored code ledger (ticket 77 AC#7)

Every third-party file that enters `frontend/` is listed here, per file, with its
source repository, the component it is, its license, the exact upstream revision
or registry URL, and what this repo changed. Each such file also carries the same
facts in its own header so `git grep` on any file answers the provenance question
without opening this one.

Delivery model: all three sources below are copy-paste libraries (Q-20). **None of
them is an npm dependency.** `package.json` names only real runtime/build packages
(react, react-dom, tailwindcss, lucide-react, clsx, tailwind-merge,
class-variance-authority, radix-ui, vite, typescript, oxlint).

## How a file here got here

```
git clone --depth 1 https://github.com/TurboKach/ai-native-react-components.git /tmp/beautifului
cd frontend && node scripts/vendor.mjs /tmp/beautifului

curl -o /tmp/sh-button.json https://ui.shadcn.com/r/styles/new-york-v4/button.json   # card, badge likewise
cd frontend && node scripts/vendor-shadcn.mjs /tmp
```

Both scripts are in the repo, both are idempotent, and both refuse to run if the
transformation they claim to perform does not apply (a surviving `use client`
directive, or a dingbat with no ASCII mapping, exits 2). Re-running them is the
way to check this ledger against the tree; hand-editing a vendored file is what
this file exists to make visible.

## Beautiful UI - `TurboKach/ai-native-react-components` (MIT)

Upstream commit vendored: `05dab2d2b5f1f3e40029776e339a486d70491079`
License text: upstream `LICENSE`, "MIT License, Copyright (c) 2026 Turbo".
Owner ruling **R19**: this library is the base of the first version.

| File here | Upstream source | Component | Changes | Mounted? |
|---|---|---|---|---|
| `src/components/ai-native/stream-text.tsx` | `components/atoms/StreamText.tsx` | StreamText | header, dropped `use client`, LF endings | **yes** - `src/components/result-stream.tsx` |
| `src/components/ai-native/shimmer.tsx` | `components/atoms/Shimmer.tsx` | Shimmer | header, dropped `use client`, LF endings | **yes** - `src/components/panel-skeleton.tsx` |
| `src/components/ai-native/approval-card.tsx` | `components/approval-card.tsx` | ApprovalCard | header, dropped `use client`, LF endings | no - upstream demo questionnaire; it is the visual reference `src/components/l2-approval-card.tsx` was cut from |
| `src/components/ai-native/loading-state.tsx` | `components/loading-state.tsx` | LoadingState | header, dropped `use client`, LF endings | no - demo strings intact |
| `src/components/ai-native/thinking.tsx` | `components/thinking.tsx` | Thinking | header, dropped `use client`, LF endings, **1 glyph ASCII-ized by hand** (U+2212 in the diff row; `vendor.mjs`'s `EMOJI_RE` does not carry the math band so it cannot catch or re-apply this) | no - demo strings intact |
| `src/components/ai-native/tool-chips.tsx` | `components/tool-chips.tsx` | ToolChips | header, dropped `use client`, LF endings, **3 glyphs ASCII-ized** (2x U+2713 via `scripts/vendor.mjs` + 1x U+2212 by hand in the diff row, see the ban #8 note) | no - tool traces arrive with ticket 35 |
| `src/components/ai-native/task-rows.tsx` | `components/task-rows.tsx` | TaskRows | header, dropped `use client`, LF endings | no - history rows are ticket 35+ |
| `src/components/ai-native/streaming-text.tsx` | `components/streaming-text.tsx` | StreamingText | header, dropped `use client`, LF endings | no - carries upstream's demo source list with external URLs, which must not render in a desktop panel |

What was **not** taken from Beautiful UI: its `app/globals.css` `:root` block. Those
40-odd declarations are that library's own palette (`--surface: #fff`, `--ink:
#1f2124`, ...), and copying them would have put a second style truth source into the
tree, which C21 forbids. Instead `src/styles/theme.css` maps their *names* onto Wisp
tokens with `var(...)` only - `--surface: var(--bg-raised-color)` - so the vendored
components render in C21's palette. That mapping is asserted literal-free by
`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`.

The 6 unmounted files are kept in the tree on purpose (R19 makes this library the
base, and ticket 35+ mounts them), and `TestVendoredDemoComponentsAreNotMounted`
proves which ones are unreachable, so "unmounted" is a fact the gate maintains
rather than a claim in this table. Vite tree-shakes an unimported module out of the
bundle entirely, so no demo string reaches the binary.

## shadcn/ui (MIT) - base primitives, registry style `new-york-v4`

License: "MIT, Copyright (c) 2023 shadcn" (upstream `LICENSE` in
`github.com/shadcn-ui/ui`). Retrieved from the public registry on 2026-09-21; the
registry items carry no commit id, so the URL and date are recorded per file header.

| File here | Registry item | Changes |
|---|---|---|
| `src/components/ui/button.tsx` | `r/styles/new-york-v4/button.json` | header; `from "cn"` -> `from "@/lib/cn"`; `text-white` -> `text-on-solid` |
| `src/components/ui/card.tsx` | `r/styles/new-york-v4/card.json` | header; `from "cn"` -> `from "@/lib/cn"` |
| `src/components/ui/badge.tsx` | `r/styles/new-york-v4/badge.json` | header; `from "cn"` -> `from "@/lib/cn"`; `text-white` -> `text-on-solid` |
| `src/lib/cn.ts` | `r/styles/new-york-v4/utils.json` (the `cn` helper) | header only |

shadcn's colour vocabulary (`primary`, `muted`, `destructive`, `card`, `ring`, ...)
is aliased onto C21 in `src/styles/theme.css`, and its `rounded-md`/`rounded-xl`
utilities are re-pointed at C21's `--r-*` scale, so even a radius written by
upstream resolves to a Wisp token.

## React Bits - **phase two, zero code in this tree**

| Item | Value |
|---|---|
| Repository | `dillionverma/react-bits` (reactbits.dev), 47.7k stars |
| License | **MIT + Commons Clause** - not a plain MIT |
| Status in this repo | **not vendored, one line of it is in the tree.** R19: the owner picked 12 animated components for a later phase and ordered phase one to ship on Beautiful UI alone |
| Gate before it may enter | **owner re-review of the Commons Clause terms is required before any React Bits code is vendored** (the clause restricts selling/Redistribution of the components themselves, which is a product-shape question, not a code question) |
| Phase-two constraints already agreed | at most one full-screen WebGL ambience layer on screen at a time (`Glass Flow` / `Aura Blob` / `Neural Float` / `Fog Sphere` are mutually exclusive), and it must be destroyed when the panel hides, not merely paused - WebView2 CPU counts against the D32 budget. `Agentic Ball` may only ever be a secondary indicator inside the panel; the desktop ball stays Win32/Direct2D (ticket 77's first hard constraint) |

Verification that the exclusion holds, runnable at any time:

```
git grep -in "react-bits\|reactbits" -- frontend/ | grep -v VENDORED.md
```

Expected: no hits. This row exists so React Bits cannot disappear from the books
merely because phase one does not use it.

## Not vendored, and why it is named here anyway

- `design/` (11 prototype screens + `design/assets/tokens.css`) is this repo's own
  work and, per ruling **R18**, a *reference* rather than a blueprint. Only
  `tokens.css` is load-bearing: it is the C21 source the frontend theme generates
  from (`scripts/gen-tokens.mjs`).
- `lucide-react` (ISC) is a genuine npm dependency, used for the D23 icon rule
  (SVG icons, never emoji). It is not copy-paste code, so it is not in this ledger.

## Gate notes for this tree (AC#4)

- **ban #6 (the panel-side approval-decision identifier, `frontend/`)** - the very
  first hit this gate produced was **this file**: writing the forbidden
  identifier into the ledger to explain the rule tripped both the scanner and
  `TestFrontendNeverNamesAnApprovalDecision`. It is written around here instead,
  because the ban is on the identifier appearing in this tree, not on
  understanding it - `tools/d22scan/main.go` and PLAN.md carry the wording.
  `tools/d22scan` walks
  `frontend/` as text files and the planted-violation run named the file
  (`frontend/src/__ban6_probe.ts:1: [panel-approval]`). The panel tree never
  contains that identifier; the callback that reports what a user asked for is
  `requestApprovalResolution` and the bridge method is `panel.approval.request`,
  because an allow decision is made natively and only natively (D33/F2). The
  scanner's `frontend/` entry is still registered `absentOK`, so the drift guard
  exits 2 while `frontend/` exists - arming it (`live:true`) is a
  `tools/d22scan` change, which is not this ticket's to make.
- **ban #8 (zero emoji)** - its declared scope is `design/`, `internal/`, `cmd/`;
  `frontend/` is not in it. Choice recorded per the ticket: **enforce it ourselves**
  with `TestFrontendHasNoEmoji`, using the scanner's own character ranges, rather
  than register `frontend/` as a permitted-emoji tree. 24 files scanned, 0 hits;
  the two glyphs that upstream shipped in `tool-chips.tsx` were ASCII-ized during
  vendoring, not after.
  > **2026-09-25 追加更正（前端会话，随本轮那 6 处字符的改动一起落）：上面这段的三句已不成立，原句保留只为旧引用可核。**
  > ① **`frontend/` 现在在扫描射程里** - 票 96 把它补进了 `emojiScopes()`，且声明为 `everyFile: true`
  > （`tools/d22scan/main.go:491-498`、`:867-871`），所以"门没照到我们、我们自行执行"这个前提已经作废；
  > ② **"0 hits" 当时是假的** - 同一把尺在 `composer.tsx:170`、`fixtures/composer-states.html:2/5/8` 的
  > U+2264 与 `thinking.tsx:213`、`tool-chips.tsx:186` 的 U+2212 上报 6 行，`internal/panel` 那份
  > `TestFrontendHasNoEmoji` 报 0，因为它的字符类（`frontend_hygiene_test.go:66`）**自称逐字抄自扫描器却没有
  > `\x{2200}-\x{22FF}`** 那一段（票 141 只加宽了扫描器那一份）。**这一处副本归编排者修（台账 `U1`），不是本文件能自结的。**
  > ③ **"during vendoring, not after" 现在只对 2 枚** - 本轮这 6 处是**手工在 vendoring 之后**改的，
  > 且 `scripts/vendor.mjs` 的 `EMOJI_RE` 同样缺数学段 ⇒ **它既看不见 U+2212、也就无法在重新拉取时把它挡回来**：
  > 一次 `npm run vendor:beautifului` 会把 `thinking.tsx` / `tool-chips.tsx` 的 U+2212 恢复原状，
  > 而那句"unmapped glyph is fatal rather than silently shipped"对它不生效。
  > ⇒ 本轮之后**盘上确实 0 命中**，但**这一格的持久性依赖 `U1` 那批把三份副本一起加宽**，不依赖本文件的自述。
