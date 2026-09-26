/* ============================================================================
   Toast stack (?harness=1 demo, 2026-09-26)
   ----------------------------------------------------------------------------
   右下角 fixed 展示栈：tooltip 皮肤（bg-tooltip-bg 底 / text-tooltip-fg 字 /
   border-tooltip-border 描边 / shadow-overlay 投影），pop-in 入场，tone 色点
   （info accent / ok green / warn warn / error red）。纯展示：pointer-events
   关闭，无状态、无定时、无关闭钮 - 消失的时机由持有 items 的人决定。

   PROPS CONTRACT:

     items: readonly ToastItem[]
         One entry per visible toast, rendered bottom-up in array order.
         id      stable React key (and nothing else - no routing on it)
         text    the message, verbatim
         tone?   "info" | "ok" | "warn" | "error"; default "info"

   The empty array renders nothing (return null) - an empty stack must not
   leave an invisible fixed container behind.

   DEMO DATA: RB_TOASTS at the bottom of this file:
     <ToastStack items={RB_TOASTS} />
   ============================================================================ */

import { cn } from "@/lib/cn";

export type ToastTone = "info" | "ok" | "warn" | "error";

export interface ToastItem {
  /** Stable key. Display identity only. */
  id: string;
  /** The message, verbatim. */
  text: string;
  /** Colour of the leading dot. Default "info". */
  tone?: ToastTone;
}

/** The tone dot's fill; the toast body itself stays tooltip-neutral. */
const DOT_CLS: Record<ToastTone, string> = {
  info: "bg-accent",
  ok: "bg-green",
  warn: "bg-warn",
  error: "bg-red",
};

export function ToastStack({ items }: { items: readonly ToastItem[] }) {
  if (items.length === 0) {
    return null;
  }

  return (
    <div
      aria-live="polite"
      className="pointer-events-none fixed bottom-4 right-4 z-50 flex flex-col items-end gap-2"
      role="status"
    >
      {items.map((item) => (
        <div
          className="flex max-w-[360px] items-center gap-2 rounded-control border border-tooltip-border bg-tooltip-bg px-3 py-2 text-[12.5px] leading-snug text-tooltip-fg shadow-overlay"
          key={item.id}
          style={{ animation: "pop-in 200ms var(--ease-out-strong) both" }}
        >
          <span
            aria-hidden="true"
            className={cn("size-1.5 shrink-0 rounded-full", DOT_CLS[item.tone ?? "info"])}
          />
          <span className="min-w-0">{item.text}</span>
        </div>
      ))}
    </div>
  );
}

/* ============================================================================
   RB_TOASTS - the harness fixture for <ToastStack />. Demo only. Four tones
   plus one entry without a tone so the info default is visible.
   ============================================================================ */

export const RB_TOASTS: ToastItem[] = [
  { id: "rb-toast-1", text: "已归档 9 张截图到 Desktop\\截图", tone: "ok" },
  { id: "rb-toast-2", text: "下载模型文件失败，可重试", tone: "error" },
  { id: "rb-toast-3", text: "月预算已用 117%", tone: "warn" },
  { id: "rb-toast-4", text: "配置已重读" },
];
