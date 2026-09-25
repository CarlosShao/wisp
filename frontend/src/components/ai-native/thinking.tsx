/* ============================================================================
   Adapted component.
   ----------------------------------------------------------------------------
   来源 = derived from TurboKach/ai-native-react-components
          components/thinking.tsx (MIT, Copyright (c) 2026 Turbo,
          upstream commit 05dab2d2).
   本地改动：
   - props 化: { seconds?, steps?, reasoning?, sources? } - 上游是定时器驱动
     的四变体演示（STAGES 分幕自动展开收起），这里全部由 props 决定；全部
     缺省时不渲染任何主体（宁缺毋造）。working = steps 里还有未完成步。
   - 数据移除: VARIANTS 演示数据（Steps/Reasoning/Search/Coding 四套）、
     Search 的外链 href、"+7 more" 演示尾巴、Coding 变体全部删除。
   - 图标替换: 上游折叠头的 sparkles 四角星按 D23 图标规则禁用，换成
     loading-state 点阵 Drive 图案的缩微版（2px 单元 + pixel-on 关键帧）。
   - 字面量换 token: 全部走 bg-hover-2 / bg-line / text-ink / text-ink-2 /
     text-ink-3 / bg-accent / bg-orange / bg-green 等 token 工具类，无颜色
     字面量；上游用 JS 量高画的展开竖线改为纯 CSS 定位线（SSR 安全）。
   - 三个子页签（步骤/推理/检索）按 bu-pill 词汇用 Tailwind 手写：
     rounded-full、bg-hover、active 为 bg-surface + shadow-btn。
   ============================================================================ */

import { useState } from "react";
import { ChevronDown } from "lucide-react";

/** One step row: done shows a muted check, not-done the rotating arc. */
export interface ThinkingStep {
  text: string;
  done: boolean;
}

/** One search-source row: title plus a body snippet, no URL in the model. */
export interface ThinkingSource {
  title: string;
  body: string;
}

type TabId = "steps" | "reasoning" | "search";

const TAB_LABELS: Record<TabId, string> = {
  steps: "步骤",
  reasoning: "推理",
  search: "检索",
};

/* loading-state Drive 点阵的延迟表（(列 + 与中行的行距) * 90ms），缩微复用。 */
const PIXEL_DELAYS: number[] = Array.from({ length: 9 }, (_, i) => {
  const r = Math.floor(i / 3);
  const c = i % 3;
  return (c + Math.abs(r - 1)) * 90;
});

function PixelGlyph({ working }: { working: boolean }) {
  return (
    <span aria-hidden="true" className="grid shrink-0 grid-cols-[repeat(3,2px)] gap-[1px]">
      {PIXEL_DELAYS.map((d, i) => (
        <span
          className="size-[2px] rounded-[0.5px] bg-ink"
          key={i}
          style={{
            opacity: 0.15,
            animation: working ? `pixel-on 650ms ease-in-out ${d}ms infinite` : "none",
          }}
        />
      ))}
    </span>
  );
}

/* 检索行的圆点三色，照上游 Search 变体的 TONES。 */
const DOT_TONES = ["bg-accent", "bg-orange", "bg-green"];

export function Thinking({
  seconds,
  steps,
  reasoning,
  sources,
}: {
  seconds?: number;
  steps?: ThinkingStep[];
  reasoning?: string;
  sources?: ThinkingSource[];
}) {
  const [manualExpanded, setManualExpanded] = useState<boolean | null>(null);
  const [picked, setPicked] = useState<TabId | null>(null);

  const tabs: TabId[] = [];
  if (steps && steps.length > 0) tabs.push("steps");
  if (reasoning && reasoning.trim() !== "") tabs.push("reasoning");
  if (sources && sources.length > 0) tabs.push("search");

  const working = steps ? steps.some((s) => !s.done) : false;
  const expanded = manualExpanded ?? true;
  const active: TabId | null = picked !== null && tabs.includes(picked) ? picked : (tabs[0] ?? null);

  if (tabs.length === 0) return null;

  const headerLabel = working
    ? "思考中"
    : seconds !== undefined
      ? `已思考 ${seconds} 秒`
      : "已完成";

  return (
    <div className="flex w-full flex-col">
      {/* 折叠头：点阵 + shimmer 文案 + 旋转箭头，蓝本的头部骨架 */}
      <button
        aria-expanded={expanded}
        className="-mx-1.5 flex w-fit items-center gap-2 rounded-control px-1.5 py-1
          transition-colors duration-100 hover:bg-hover-2"
        onClick={() => setManualExpanded((current) => !(current ?? true))}
        type="button"
      >
        <PixelGlyph working={working} />
        {working ? (
          <span
            className="bg-clip-text whitespace-nowrap text-[13px] font-medium text-transparent"
            style={{
              backgroundImage:
                "linear-gradient(90deg, var(--ink-3) 35%, var(--ink) 50%, var(--ink-3) 65%)",
              backgroundSize: "200% 100%",
              animation: "shimmer-text 1.4s linear infinite",
            }}
          >
            {headerLabel}
          </span>
        ) : (
          <span
            className="whitespace-nowrap text-[13px] font-medium text-ink-2"
            style={{ animation: "fade-in 350ms ease-out both" }}
          >
            {headerLabel}
          </span>
        )}
        <ChevronDown
          aria-hidden="true"
          className="text-ink-3 transition-transform duration-300"
          size={14}
          strokeWidth={2.2}
          style={{ transform: expanded ? "rotate(180deg)" : "rotate(0deg)" }}
        />
      </button>

      {/* 展开体：子页签 + 竖线走带，grid-rows 收展动画照蓝本 */}
      <div
        className="grid transition-[grid-template-rows,opacity] duration-400"
        style={{
          gridTemplateRows: expanded ? "1fr" : "0fr",
          opacity: expanded ? 1 : 0,
          transitionTimingFunction: "var(--ease-out-strong)",
        }}
      >
        <div className="overflow-hidden">
          <div className="flex items-center gap-1 pt-1.5">
            {tabs.map((id) => (
              <button
                aria-pressed={active === id}
                className={`inline-flex h-6 items-center rounded-full px-2.5 text-[11px] font-medium
                  transition-colors duration-100 ${
                    active === id
                      ? "bg-surface text-ink shadow-btn"
                      : "bg-hover text-ink-2 hover:bg-hover-2"
                  }`}
                key={id}
                onClick={() => setPicked(id)}
                type="button"
              >
                {TAB_LABELS[id]}
              </button>
            ))}
          </div>

          <div className="relative mt-1 ml-[5px] pl-4">
            <span aria-hidden="true" className="absolute bottom-1 left-[3px] top-1 w-px bg-line" />
            <div className="flex flex-col gap-1 py-1">
              {active === "steps" && steps
                ? steps.map((step, i) => (
                    <div
                      className="flex min-h-7 w-full items-center gap-2 rounded-[6px] px-1.5 py-0.5 text-left"
                      key={`${i}-${step.text}`}
                      style={{ animation: `fade-up 320ms var(--ease-out-strong) ${i * 120}ms both` }}
                    >
                      {step.done ? (
                        <svg
                          className="shrink-0"
                          fill="none"
                          height="14"
                          stroke="var(--ink-3)"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2.5"
                          viewBox="0 0 24 24"
                          width="14"
                        >
                          <path d="M20 6L9 17l-5-5" />
                        </svg>
                      ) : (
                        <span
                          aria-hidden="true"
                          className="size-3 shrink-0 rounded-full border-[1.5px] border-line-strong border-t-ink-2"
                          style={{ animation: "spin 700ms linear infinite" }}
                        />
                      )}
                      <span className="min-w-0 truncate text-[12.5px] font-medium text-ink">
                        {step.text}
                      </span>
                    </div>
                  ))
                : null}

              {active === "reasoning" && reasoning ? (
                <p className="whitespace-normal px-1.5 py-1 text-[12.5px] leading-relaxed text-ink-2">
                  {reasoning}
                </p>
              ) : null}

              {active === "search" && sources
                ? sources.map((source, i) => (
                    <div
                      className="flex min-h-7 w-full items-center gap-2 rounded-[6px] px-1.5 py-0.5 text-left"
                      key={`${i}-${source.title}`}
                      style={{ animation: `fade-up 320ms var(--ease-out-strong) ${i * 120}ms both` }}
                    >
                      <span
                        aria-hidden="true"
                        className={`size-2 shrink-0 rounded-full ${DOT_TONES[i % DOT_TONES.length]}`}
                      />
                      <span className="min-w-0 truncate text-[12.5px] font-medium text-ink">
                        {source.title}
                      </span>
                      <span className="max-w-[40%] shrink-0 truncate text-right text-[11.5px] text-ink-3">
                        {source.body}
                      </span>
                    </div>
                  ))
                : null}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
