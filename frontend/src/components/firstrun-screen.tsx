/* ============================================================================
   First-run screen (showcase component, 2026-09-26)
   ----------------------------------------------------------------------------
   The ball state machine's FirstRun state (first boot: model download +
   directory authorization) as one guided screen, redrawn in the beautiful-ui
   vocabulary: numbered rail dots + connector lines, the ProgressRing download
   card, the verified-model list, EntityChip directory rows with a delete slot,
   and the masked-key row. The content data follows the demo's firstrun screen
   (tickets 14/63: four offline voice models with sizes and minisign states,
   the 68 percent ring, the DPAPI line); on this page it arrives exclusively
   through src/fixtures/harness.ts.

   PROPS CONTRACT - everything visible is a prop; the component holds no state
   and decides nothing:
     title / subtitle          heading pair
     steps                     rail labels in order (the firstrun spec fixes
                               three: choose models / authorize dirs / keys)
     currentStep               which step the screen shows, 0-based, clamped
     modelRows[]               name (mono) / desc / size / minisign, minisign
                               being "verified" or "verifying" (the two
                               verification states of the signed manifest)
     progress                  0..1; drives the ProgressRing arc only
     downloadStatus / note     the download card copy
     dirs / dirTitle / ...     the authorized-directory card; onRemoveDir(path)
                               and onAddDir(path) hand each change to whoever
                               mounts this - nothing is mutated here
     keyLabel / keyNote /      the masked key row (font-mono + an eye-off slot
     keyMasked / dpapiNote     that stays a slot) and the DPAPI line
     onSkipStep(step)          the per-step skip affordance
     onNext()                  advance / finish; the mounter owns step state

   生产路径由宿主触发（票 33/14 接线），本组件只负责观感：the callbacks are
   affordances, not behaviour - no route, no storage, no host call lives here.
   The three action labels (跳过 / 下一步 / 完成) are fixed UI vocabulary of
   the same class as the tasks screen's state names; every data string rides
   in through the props above.
   ============================================================================ */

import { Fragment, type FormEvent } from "react";
import {
  Check,
  ChevronRight,
  EyeOff,
  FolderLock,
  KeyRound,
  LoaderCircle,
  LockKeyhole,
  ShieldCheck,
  Trash,
} from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { EntityChip } from "@/components/ai-native/entity-chip";
import { ProgressRing } from "@/components/ai-native/progress-ring";
import { Shimmer } from "@/components/ai-native/shimmer";
import { StatusPill } from "@/components/ai-native/status-pill";

export type MinisignState = "verified" | "verifying";

export interface FirstrunModelRow {
  /** The model's identifier, set in mono. */
  name: string;
  /** The purpose line, verbatim from the manifest copy. */
  desc: string;
  /** The size label, verbatim ("245MB"). */
  size: string;
  /** The minisign verification state of the row. */
  minisign: MinisignState;
}

export interface FirstrunScreenProps {
  title: string;
  subtitle: string;
  steps: readonly string[];
  currentStep: number;
  modelRows: readonly FirstrunModelRow[];
  progress: number;
  downloadStatus: string;
  downloadNote: string;
  dirs: readonly string[];
  dirTitle: string;
  dirNote: string;
  dirPlaceholder: string;
  dirAddLabel: string;
  keyLabel: string;
  keyNote: string;
  keyMasked: string;
  dpapiNote: string;
  onRemoveDir?: (path: string) => void;
  onAddDir?: (path: string) => void;
  onSkipStep?: (step: number) => void;
  onNext?: () => void;
}

/** The two minisign states, named like the tasks screen names its six. */
const MINISIGN_LABELS: Record<MinisignState, string> = {
  verified: "已校验",
  verifying: "验签中",
};

function MinisignMark({ state }: { state: MinisignState }) {
  if (state === "verified") {
    return (
      <StatusPill className="shrink-0" dot={false} tone="green">
        <ShieldCheck aria-hidden="true" className="size-3" />
        {MINISIGN_LABELS[state]}
      </StatusPill>
    );
  }
  return (
    <StatusPill className="shrink-0" dot={false} tone="accent">
      <LoaderCircle
        aria-hidden="true"
        className="size-3"
        style={{ animation: "spin 1.1s linear infinite" }}
      />
      {MINISIGN_LABELS[state]}
    </StatusPill>
  );
}

/* The rail: solid accent discs for reached steps (check when done, the number
   while active, the halo marks the active one), hollow field discs for the
   steps not yet reached; the connector fills accent once passed. */
function StepRail({ steps, current }: { steps: readonly string[]; current: number }) {
  return (
    <ol
      className="flex w-full items-start"
      style={{ animation: "fade-up 300ms cubic-bezier(0.23,1,0.32,1) 40ms both" }}
    >
      {steps.map((label, i) => {
        const done = i < current;
        const active = i === current;
        const dotCls = active
          ? "bg-accent text-white ring-4 ring-accent-tint"
          : done
            ? "bg-accent text-white"
            : "border border-line-strong bg-field text-ink-3";
        const labelCls = active
          ? "font-medium text-accent-ink"
          : done
            ? "font-medium text-ink"
            : "text-ink-3";
        return (
          <Fragment key={label}>
            {i > 0 && (
              <span
                aria-hidden="true"
                className={`mt-[11px] h-0.5 min-w-4 flex-1 rounded-full ${i <= current ? "bg-accent" : "bg-line-strong"}`}
              />
            )}
            <li className="flex w-16 shrink-0 flex-col items-center gap-1.5">
              <span
                aria-hidden="true"
                className={`flex size-6 items-center justify-center rounded-full text-[11px] font-medium leading-none ${dotCls}`}
              >
                {done ? <Check className="size-3.5" /> : i + 1}
              </span>
              <span className={`text-[11px] leading-none ${labelCls}`}>{label}</span>
            </li>
          </Fragment>
        );
      })}
    </ol>
  );
}

function ModelStepPanel({
  rows,
  progress,
  status,
  note,
}: {
  rows: readonly FirstrunModelRow[];
  progress: number;
  status: string;
  note: string;
}) {
  return (
    <div
      className="flex flex-col gap-2.5"
      style={{ animation: "fade-up 300ms cubic-bezier(0.23,1,0.32,1) 80ms both" }}
    >
      <div className="flex items-center gap-3.5 rounded-control border border-line bg-inset p-3.5">
        <ProgressRing progress={progress} size={56} tone="accent">
          {Math.round(progress * 100)}%
        </ProgressRing>
        <div className="min-w-0">
          <Shimmer className="text-[13px] font-semibold">{status}</Shimmer>
          <p className="mt-1 text-[11.5px] leading-[1.6] text-ink-3">{note}</p>
        </div>
      </div>
      <div className="overflow-hidden rounded-control border border-line bg-surface">
        {rows.map((row, i) => (
          <div
            className="flex items-center justify-between gap-3 border-b border-line px-3 py-2.5 transition-colors duration-100 last:border-b-0 hover:bg-hover"
            key={row.name}
            style={{ animation: `fade-up 300ms cubic-bezier(0.23,1,0.32,1) ${80 + i * 40}ms both` }}
          >
            <div className="min-w-0">
              <p className="truncate font-mono text-[12px] font-medium text-ink">{row.name}</p>
              <p className="mt-0.5 text-[11px] text-ink-3">
                {row.desc} · {row.size}
              </p>
            </div>
            <MinisignMark state={row.minisign} />
          </div>
        ))}
      </div>
    </div>
  );
}

function DirsStepPanel({
  paths,
  title,
  note,
  placeholder,
  addLabel,
  onRemoveDir,
  onAddDir,
}: {
  paths: readonly string[];
  title: string;
  note: string;
  placeholder: string;
  addLabel: string;
  onRemoveDir?: (path: string) => void;
  onAddDir?: (path: string) => void;
}) {
  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const input = event.currentTarget.elements.namedItem("dir");
    if (!(input instanceof HTMLInputElement) || !onAddDir) {
      return;
    }
    const value = input.value.trim();
    if (value) {
      onAddDir(value);
    }
    event.currentTarget.reset();
  }

  return (
    <div
      className="rounded-control border border-line bg-inset p-3.5"
      style={{ animation: "fade-up 300ms cubic-bezier(0.23,1,0.32,1) 80ms both" }}
    >
      <p className="text-[13px] font-medium text-ink">{title}</p>
      <p className="mt-1 text-[11.5px] leading-[1.6] text-ink-3">{note}</p>
      <div className="mt-2.5 overflow-hidden rounded-control border border-line bg-surface">
        {paths.map((path) => (
          <div
            className="flex items-center justify-between gap-2 border-b border-line px-2.5 py-2 transition-colors duration-100 last:border-b-0 hover:bg-hover"
            key={path}
          >
            <EntityChip
              className="min-w-0"
              color="var(--accent)"
              monogram={<FolderLock aria-hidden="true" className="size-2.5" />}
              name={path}
            />
            <button
              aria-label={`删除目录 ${path}`}
              className="flex size-6 shrink-0 items-center justify-center rounded-full text-ink-3 transition-colors duration-100 hover:bg-hover-2 hover:text-red"
              onClick={() => onRemoveDir?.(path)}
              type="button"
            >
              <Trash aria-hidden="true" className="size-3.5" />
            </button>
          </div>
        ))}
      </div>
      <form className="mt-2.5 flex gap-2" onSubmit={submit}>
        <input
          className="h-7 min-w-0 flex-1 rounded-control border border-line bg-surface px-2.5 text-[12px] text-ink outline-none placeholder:text-ink-3 focus:border-accent"
          name="dir"
          placeholder={placeholder}
          type="text"
        />
        <Button size="xs" type="submit" variant="secondary">
          {addLabel}
        </Button>
      </form>
    </div>
  );
}

function CredentialsStepPanel({
  label,
  note,
  keyMasked,
  dpapiNote,
}: {
  label: string;
  note: string;
  keyMasked: string;
  dpapiNote: string;
}) {
  return (
    <div
      className="rounded-control border border-line bg-inset p-3.5"
      style={{ animation: "fade-up 300ms cubic-bezier(0.23,1,0.32,1) 80ms both" }}
    >
      <p className="text-[13px] font-medium text-ink">{label}</p>
      <p className="mt-1 text-[11.5px] leading-[1.6] text-ink-3">{note}</p>
      <div className="mt-2.5 flex h-9 items-center gap-2.5 rounded-control border border-line bg-surface px-3">
        <KeyRound aria-hidden="true" className="size-3.5 shrink-0 text-ink-3" />
        <span className="min-w-0 flex-1 truncate font-mono text-[12.5px] text-ink">{keyMasked}</span>
        <button
          aria-label="显示密钥"
          className="flex size-6 shrink-0 items-center justify-center rounded-full text-ink-3 transition-colors duration-100 hover:bg-hover-2 hover:text-ink"
          type="button"
        >
          <EyeOff aria-hidden="true" className="size-3.5" />
        </button>
      </div>
      <p className="mt-2.5 flex items-start gap-1.5 text-[11px] leading-[1.6] text-ink-3">
        <LockKeyhole aria-hidden="true" className="mt-0.5 size-3.5 shrink-0" />
        <span>{dpapiNote}</span>
      </p>
    </div>
  );
}

export function FirstrunScreen({
  title,
  subtitle,
  steps,
  currentStep,
  modelRows,
  progress,
  downloadStatus,
  downloadNote,
  dirs,
  dirTitle,
  dirNote,
  dirPlaceholder,
  dirAddLabel,
  keyLabel,
  keyNote,
  keyMasked,
  dpapiNote,
  onRemoveDir,
  onAddDir,
  onSkipStep,
  onNext,
}: FirstrunScreenProps) {
  const step = Math.min(Math.max(currentStep, 0), steps.length - 1);
  const last = step === steps.length - 1;

  return (
    <div className="mx-auto flex w-full max-w-[440px] flex-col gap-4">
      <div style={{ animation: "fade-up 300ms cubic-bezier(0.23,1,0.32,1) both" }}>
        <p className="text-[15px] font-semibold text-ink">{title}</p>
        <p className="mt-1 text-[12px] text-ink-2">{subtitle}</p>
      </div>

      <StepRail current={step} steps={steps} />

      {step === 0 && (
        <ModelStepPanel note={downloadNote} progress={progress} rows={modelRows} status={downloadStatus} />
      )}
      {step === 1 && (
        <DirsStepPanel
          addLabel={dirAddLabel}
          note={dirNote}
          onAddDir={onAddDir}
          onRemoveDir={onRemoveDir}
          paths={dirs}
          placeholder={dirPlaceholder}
          title={dirTitle}
        />
      )}
      {step === 2 && (
        <CredentialsStepPanel dpapiNote={dpapiNote} keyMasked={keyMasked} label={keyLabel} note={keyNote} />
      )}

      <div
        className="flex items-center justify-end gap-2"
        style={{ animation: "fade-in 250ms ease-out 120ms both" }}
      >
        <Button size="xs" onClick={() => onSkipStep?.(step)} variant="ghost">
          跳过
        </Button>
        <Button size="xs" onClick={onNext} variant="primary">
          {last ? "完成" : "下一步"}
          {!last && <ChevronRight aria-hidden="true" className="size-3.5" />}
        </Button>
      </div>
    </div>
  );
}
