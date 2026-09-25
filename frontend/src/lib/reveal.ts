/* ============================================================================
   Reveal segmentation (frontend session, 2026-09-25).

   PLAN.md:3476 (the frozen "SSE 流式输出" row of the D29 state-visual table)
   asks for three things at once:
     视觉  文字逐段出现                        <- text arrives segment by segment
     动效  光标 opacity 1->0.2->1 呼吸，1s      <- see src/styles/theme.css
     不用  逐字符动画（CPU 与可读性都差，按 SSE 分块追加即可）

   The vendored Beautiful UI atom this replaces (src/components/ai-native/
   stream-text.tsx, NOT mounted) produced its segments with text.split(" ").
   A Chinese answer has no spaces, so that call returns ONE element: the whole
   paragraph blurred in as a single unit and the caret was dropped 46 ms later.
   "逐段" was therefore satisfied for Latin text and structurally unsatisfiable
   for the product's main language.

   This file holds the rule on its own so it can be checked without a browser:
   the reveal itself is timed by React and cannot be observed in static markup.
   ============================================================================ */

/** True for one code point of CJK text (Han incl. Ext A, kana, Hangul). */
const CJK = /[\u3400-\u4DBF\u4E00-\u9FFF\uF900-\uFAFF\u3040-\u30FF\uAC00-\uD7AF]/;

/** Clause-final punctuation. Deliberately only what ENDS a clause: an opening
    「 or （ must stay with what follows it, or the reveal stutters. */
const CLAUSE_END = /[，。、；：！？…]/u;

/** A run of CJK with no punctuation in it is still cut, so one unpunctuated
    answer cannot become one 300-character segment. A Latin word is never
    cut, which is why the test is on the segment's first code point. */
export const MAX_SEGMENT = 12;

function cutAtBoundaries(text: string): string[] {
  const out: string[] = [];
  let start = 0;
  let i = 0;
  while (i < text.length) {
    if (/\s/.test(text[i])) {
      // Whitespace attaches to what precedes it (upstream's split(" ") shape),
      // except a run at the very start of a segment, which belongs to the
      // text after it -- otherwise an indented code line reveals as three
      // empty-looking spaces first.
      let j = i;
      while (j < text.length && /\s/.test(text[j])) j += 1;
      if (j > start) {
        out.push(text.slice(start, j));
        start = j;
      }
      i = j;
      continue;
    }
    if (CLAUSE_END.test(text[i])) {
      out.push(text.slice(start, i + 1));
      start = i + 1;
    }
    i += 1;
  }
  if (start < text.length) out.push(text.slice(start));
  return out;
}

function cap(seg: string, out: string[]): void {
  // Array.from, not slice: chunking by UTF-16 unit would cut an astral code
  // point in half and render a replacement glyph mid-answer.
  const chars = Array.from(seg);
  if (chars.length <= MAX_SEGMENT || !CJK.test(chars[0] as string)) {
    out.push(seg);
    return;
  }
  const chunks: string[] = [];
  for (let i = 0; i < chars.length; i += MAX_SEGMENT) {
    chunks.push(chars.slice(i, i + MAX_SEGMENT).join(""));
  }
  const last = chunks[chunks.length - 1] as string;
  if (/^[\s，。、；：！？…]+$/u.test(last)) {
    // A trailing punctuation-only segment blurs in alone and reads as a glitch.
    chunks[chunks.length - 2] += last;
    chunks.pop();
  }
  out.push(...chunks);
}

/** Split assistant text into reveal units. join() of the result is always the
    input: no character is inserted, dropped or reordered. */
export function revealSegments(text: string): string[] {
  const out: string[] = [];
  for (const seg of cutAtBoundaries(text)) cap(seg, out);
  return out;
}
