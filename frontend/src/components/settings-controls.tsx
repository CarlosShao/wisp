/* ============================================================================
   Settings controls - the small pieces both settings pages (模型 / 语音) draw.
   ----------------------------------------------------------------------------
   One vocabulary, so the two pages cannot drift apart: a card, a status dot, a
   dropdown, a provider+model pair control, a feedback line and an empty state.
   Colours are token utilities only, zero literals; no emoji.

   The dropdown is a button + popover (no native <select>) because each option
   can carry a right-hand note ("已配置", "手动输入") that a <select> cannot
   draw. Dismissal is lib/use-dismiss: outside pointerdown, Escape, or its own
   trigger.
   ============================================================================ */

import { useRef, useState, type ReactNode } from "react";
import { Check, ChevronDown } from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { cn } from "@/lib/cn";
import { useDismiss } from "@/lib/use-dismiss";
import type { Feedback, Tone } from "@/lib/config-feedback";
import type { ConfigViewProvider } from "@/lib/config-view";

/* ---------------------------------------------------------------------------
   Card + labels.
   --------------------------------------------------------------------------- */

export function SettingsCard({
  title,
  description,
  children,
}: {
  title: string;
  description?: string;
  children: ReactNode;
}) {
  return (
    <section className="rounded-card border border-line bg-surface">
      <header className="border-b border-line px-4 py-2.5">
        <h3 className="text-[12.5px] font-semibold text-ink">{title}</h3>
        {description ? (
          <p className="mt-0.5 text-[11px] leading-[1.6] text-ink-3">{description}</p>
        ) : null}
      </header>
      <div className="px-4 py-3">{children}</div>
    </section>
  );
}

/** The dot that carries a tone without words where a word would crowd. */
export function StateDot({ tone, className }: { tone: Tone; className?: string }) {
  return (
    <span
      aria-hidden="true"
      className={cn(
        "inline-block size-2 shrink-0 rounded-full",
        tone === "ok" && "bg-green",
        tone === "warn" && "bg-orange",
        tone === "error" && "bg-red",
        tone === "neutral" && "bg-ink-3",
        className,
      )}
    />
  );
}

/** A compact tag (smaller than StatusPill): row-level states read as a tag,
    not as a banner. */
export function StateTag({ tone, text }: { tone: Tone; text: string }) {
  return (
    <span
      className={cn(
        "inline-flex h-5 shrink-0 items-center gap-1.5 rounded-full px-2 text-[11px] font-medium leading-none",
        tone === "ok" && "bg-green-tint text-green",
        tone === "warn" && "bg-orange-tint text-orange",
        tone === "error" && "bg-red-tint text-red",
        tone === "neutral" && "bg-inset text-ink-2",
      )}
    >
      <StateDot tone={tone} />
      {text}
    </span>
  );
}

export function FeedbackLine({ feedback }: { feedback: Feedback | null }) {
  if (feedback === null) return null;
  return (
    <p
      className={cn(
        "text-[11.5px] leading-[1.6]",
        feedback.tone === "ok" && "text-green",
        feedback.tone === "error" && "text-red",
        feedback.tone === "warn" && "text-orange",
        feedback.tone === "neutral" && "text-ink-3",
      )}
      role={feedback.tone === "error" ? "alert" : undefined}
    >
      {feedback.text}
    </p>
  );
}

export function EmptyState({
  text,
  actionLabel,
  onAction,
}: {
  text: string;
  actionLabel?: string;
  onAction?: () => void;
}) {
  return (
    <div className="flex flex-col items-start gap-3 py-2">
      <p className="max-w-[420px] text-[12px] leading-[1.7] text-ink-2">{text}</p>
      {actionLabel && onAction ? (
        <Button onClick={onAction} size="sm" variant="primary">
          {actionLabel}
        </Button>
      ) : null}
    </div>
  );
}

/* ---------------------------------------------------------------------------
   Dropdown.
   --------------------------------------------------------------------------- */

export interface PickerOption {
  value: string;
  /** Right-hand note; empty string draws nothing. */
  note?: string;
}

/** The sentinel a model list uses for its "手动输入" row. Not a legal id, so it
    can never collide with a real model name. */
export const MANUAL_ENTRY = "\u0000manual-entry";

export function Picker({
  ariaLabel,
  value,
  options,
  placeholder = "选择",
  emptyText = "没有可选项",
  onSelect,
  className,
}: {
  ariaLabel: string;
  value: string;
  options: readonly PickerOption[];
  placeholder?: string;
  emptyText?: string;
  onSelect: (next: string) => void;
  className?: string;
}) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  useDismiss(ref, open, () => setOpen(false));
  const current = options.find((o) => o.value === value);

  return (
    <div className={cn("relative", className)} ref={ref}>
      <button
        aria-expanded={open}
        aria-label={ariaLabel}
        className="flex h-[30px] w-full items-center gap-1.5 rounded-control bg-inset px-2.5 text-left text-[12px] text-ink shadow-hairline transition-colors duration-150 hover:bg-hover"
        onClick={() => setOpen((v) => !v)}
        type="button"
      >
        <span className={cn("min-w-0 flex-1 truncate font-mono", current ? "text-ink" : "text-ink-3")}>
          {current?.value ?? placeholder}
        </span>
        <ChevronDown aria-hidden="true" className="shrink-0 text-ink-3" size={13} strokeWidth={2} />
      </button>
      {open ? (
        <div
          className="absolute left-0 top-full z-40 mt-1 max-h-64 min-w-full overflow-y-auto rounded-card border border-line bg-surface p-1 shadow-overlay"
          style={{ animation: "pop-in 160ms var(--ease-out-strong) both" }}
        >
          {options.length === 0 ? (
            <p className="px-2 py-2 text-[11.5px] text-ink-3">{emptyText}</p>
          ) : (
            options.map((option) => (
              <button
                className="flex w-full items-center gap-2 rounded-chip px-2 py-1.5 text-left transition-colors duration-100 hover:bg-hover"
                key={option.value}
                onClick={() => {
                  onSelect(option.value);
                  setOpen(false);
                }}
                type="button"
              >
                <span className="min-w-0 flex-1 truncate font-mono text-[12px] text-ink">
                  {option.value}
                </span>
                {option.note ? (
                  <span className="shrink-0 text-[10.5px] text-ink-3">{option.note}</span>
                ) : null}
                {option.value === value ? (
                  <Check aria-hidden="true" className="shrink-0 text-accent-ink" size={12} strokeWidth={2.4} />
                ) : null}
              </button>
            ))
          )}
        </div>
      ) : null}
    </div>
  );
}

/* ---------------------------------------------------------------------------
   Provider + model pair.
   --------------------------------------------------------------------------- */

/** The model half of a pair: the catalog list when the read has one, a typed id
    when it does not (`manualFallback`), and the reason printed when neither can
    produce a value. The typed id is the SECONDARY path - it exists because a
    provider nobody has configured yet cannot have a catalog list. */
export function ModelField({
  row,
  provider,
  model,
  manual,
  onChange,
  onManual,
  manualFallback,
  className,
}: {
  row: ConfigViewProvider | null;
  provider: string;
  model: string;
  manual: boolean;
  onChange: (next: string) => void;
  onManual: () => void;
  /** Whether a typed id is offered at all (an existing provider's catalog is
      authoritative, so its pair pickers do not offer one). */
  manualFallback: boolean;
  className?: string;
}) {
  const models = row !== null && row.models !== null && row.models.length > 0 ? row.models : null;

  if (models !== null && !manual) {
    const options: PickerOption[] = [
      ...models.map((name) => ({ value: name })),
      ...(manualFallback ? [{ value: MANUAL_ENTRY, note: "手动输入" }] : []),
    ];
    return (
      <Picker
        ariaLabel="模型"
        className={className}
        onSelect={(next) => {
          if (next === MANUAL_ENTRY) {
            onManual();
            onChange("");
          } else {
            onChange(next);
          }
        }}
        options={options}
        placeholder="选择模型"
        value={model}
      />
    );
  }

  if (manualFallback) {
    return (
      <input
        aria-label="模型 id"
        className={cn(
          "h-[30px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3",
          className,
        )}
        onChange={(e) => onChange(e.currentTarget.value)}
        placeholder="模型 id"
        value={model}
      />
    );
  }

  // No list and no manual path: the honest reason, as words, not a dead input.
  return (
    <span className={cn("flex h-[30px] items-center text-[11.5px] text-ink-3", className)}>
      {provider === ""
        ? "先选供应商"
        : row === null
          ? "这家供应商还没有配置"
          : "这家供应商暂无模型"}
    </span>
  );
}

/** Provider + model as one control. Owns only the manual/automatic flip; both
    values are the parent's (a saved pair must stay visible after a reload). */
export function PairFields({
  providers,
  provider,
  model,
  onProvider,
  onModel,
  manualFallback,
  providerLabel = "供应商",
}: {
  providers: readonly ConfigViewProvider[];
  provider: string;
  model: string;
  onProvider: (next: string) => void;
  onModel: (next: string) => void;
  manualFallback: boolean;
  providerLabel?: string;
}) {
  // Which provider the manual flag belongs to. Switching providers therefore
  // drops out of manual mode by derivation (manualFor !== provider) instead of
  // by an effect that resets state - a control whose modes leak across a
  // provider switch is exactly the bug this shape prevents.
  const [manualFor, setManualFor] = useState<string | null>(null);
  const manual = manualFor === provider && provider !== "";
  const row = providers.find((p) => p.name === provider) ?? null;
  return (
    <>
      <Picker
        ariaLabel={providerLabel}
        className="w-[180px] shrink-0"
        emptyText="还没有已配置的供应商"
        onSelect={onProvider}
        options={providers.map((p) => ({
          value: p.name,
          note: p.credentialSet ? "密钥已保存" : "密钥未保存",
        }))}
        placeholder="选择供应商"
        value={provider}
      />
      <ModelField
        className="w-[220px] shrink-0"
        manual={manual}
        manualFallback={manualFallback}
        model={model}
        onChange={onModel}
        onManual={() => setManualFor(provider)}
        provider={provider}
        row={row}
      />
    </>
  );
}
