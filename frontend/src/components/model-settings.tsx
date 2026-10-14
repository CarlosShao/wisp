/* ============================================================================
   模型 settings page - ticket 313's page half (2026-10-13).
   ----------------------------------------------------------------------------
   WHAT THIS IS. The 模型 category of the settings page, wired to the settings
   door the Go side landed (internal/panel/config_handlers.go, config_handlers
   read + config_write_313.go). Everything rendered here comes from a
   config.get reply the page just asked for - the prose half is shown as the
   human sentence, and the pickers are built ONLY from the structured line
   (wisp-config-view:), never by parsing the prose (313-panel-ask §1).

   THE THREE WRITES, in the order the spec recommends (§3):
     1. provider_row_create - the provider's section does not exist yet; this
        is the one field that creates it. It accepts a built-in preset name
        only, and a row that already exists is refused. So the page shows the
        step as DONE (from the read payload, not by trying) when providers[]
        already names the provider.
     2. provider_credential - the value rides credentialValue and is cleared
        from the input the instant it is sent; it is never stored in state and
        never rendered again. What stays on screen is the status sentence and
        the reference name the store recorded.
     3. role_chat_pair - llm.roles.chat's provider and model, both halves in
        ONE write, because validateRoles refuses both half states - so the
        pair guard (set both or neither) is enforced HERE, on the page.

   THE DEFAULTS. While configuredKnown && configured === false, the provider
   picker default-highlights stepfun (A877 ruling 3) - and that default writes
   NOTHING. "The user changed something, the user wins" holds by construction:
   no request leaves this page until one of the buttons is pressed.

   THE MODEL PICKER has a manual id fallback by requirement: stepfun has zero
   real model ids in this repository (census c-default.md C4), so a dropdown
   of "the models the config names" is empty until the user wrote one. When
   the provider's models list is null (not answered) or empty, the control is
   a hand-entry field with the reason printed next to it. No model name is
   ever invented here.

   RECEIPTS ARE FACTS (§3): a write's reply is shown verbatim - the receipt
   sentence for an accepted write, the refusal sentence for a refused one -
   and this page does not add any "already in force" wording of its own; the
   tier sentence inside the reply is the only place that answer lives.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useMemo, useRef, useState, type ReactNode } from "react";
import { ChevronDown, KeyRound, RefreshCw } from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { Chip } from "@/components/ai-native/chip";
import { ConfigReadGate } from "@/components/config-read-state";
import { cn } from "@/lib/cn";
import { useDismiss } from "@/lib/use-dismiss";
import { CREDENTIAL_LABELS, type ConfigView, type ConfigViewProvider } from "@/lib/config-view";
import { requestConfigWrite, requestCredentialWrite } from "@/lib/panel";

/* ---------------------------------------------------------------------------
   The page shell: the shared read gate owns loading / no-host / unreadable /
   payload-error, and hands the ready view to the body below.
   --------------------------------------------------------------------------- */

export function ModelSettings({
  heading = "模型",
  lede = "选服务商、录密钥、把聊天角色指过去。这一页的每一次写入都要重启进程才用得上，回执里会说。",
}: {
  /** The settings page already names this category in its own bar, so it
      passes its own words in here - two stacked headings for one section is
      a layout bug, not a title (same convention as config-screen.tsx). */
  heading?: string;
  lede?: string;
}) {
  return (
    <section aria-label="模型" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>
      <ConfigReadGate>
        {(view, prose, reload) => <ModelBody onChanged={reload} prose={prose} view={view} />}
      </ConfigReadGate>
    </section>
  );
}

/* ---------------------------------------------------------------------------
   The body: status, then the three steps.
   --------------------------------------------------------------------------- */

/** One write's outcome, kept as the sentence the host sent - never rewrapped
    into "success"/"saved" wording of this page's own. */
interface Receipt {
  step: string;
  text: string;
  refused: boolean;
}

function ModelBody({
  view,
  prose,
  onChanged,
}: {
  view: ConfigView;
  prose: string;
  onChanged: () => void;
}) {
  const presets = view.presetsKnown && view.presets !== null ? view.presets : null;
  // Default-highlight stepfun while nothing is configured (A877 ruling 3).
  // Computed, not written: this is a selection the picker shows, and nothing
  // leaves the page until a button below is pressed.
  const suggested =
    view.configuredKnown && !view.configured && presets?.includes("stepfun") ? "stepfun" : null;
  const [pickedProvider, setPickedProvider] = useState<string | null>(null);
  const provider = pickedProvider ?? suggested ?? (view.chatProvider !== "" ? view.chatProvider : (presets?.[0] ?? ""));
  const providerRow = useMemo(
    () => view.providers.find((p) => p.name === provider) ?? null,
    [view.providers, provider],
  );

  const [model, setModel] = useState<string>(view.chatModelKnown ? view.chatModel : "");
  const [manual, setManual] = useState(false);
  const [receipts, setReceipts] = useState<Receipt[]>([]);
  const [busy, setBusy] = useState<string | null>(null);

  async function run(step: string, send: () => Promise<string | null>) {
    setBusy(step);
    const reply = await send();
    setBusy(null);
    if (reply === null) {
      setReceipts((all) => [{ step, text: "没接到宿主，请求没有发出去。", refused: true }, ...all].slice(0, 6));
      return;
    }
    // The refusal sentence is sealed on the Go side ("面板请求被拒绝 [...]");
    // classifying it colours the card, and the text is shown verbatim either
    // way, so a wrong guess here cannot turn a refusal into a success.
    setReceipts((all) => [{ step, text: reply, refused: reply.includes("被拒绝") }, ...all].slice(0, 6));
    onChanged();
  }

  const rowExists = providerRow !== null;
  // The pair write needs both halves. When the leg answered a catalog and the
  // picker is showing its first entry without the user having touched it, that
  // shown entry IS the value the write will carry - a control that displays a
  // model while the button next to it acts on "" would be the classic
  // picker/save mismatch.
  const catalogModels =
    providerRow !== null && providerRow.models !== null && providerRow.models.length > 0
      ? providerRow.models
      : null;
  const effectiveModel =
    model.trim() !== "" ? model.trim() : catalogModels !== null && !manual ? catalogModels[0]! : "";

  return (
    <div className="mt-2 flex min-w-0 flex-col gap-3">
      {/* 读到的现状：散文照抄 + 结构化事实。 */}
      <div className="flex flex-col gap-1.5 rounded-card border border-line bg-surface px-3.5 py-3">
        <div className="flex items-center justify-between gap-2">
          <p className="text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3">配置读回</p>
          <Button aria-label="重读配置" onClick={onChanged} size="xs" variant="quiet">
            <RefreshCw aria-hidden="true" size={11} strokeWidth={2} />
            重读
          </Button>
        </div>
        <pre className="whitespace-pre-wrap text-[12px] leading-[1.7] text-ink-2">{prose}</pre>
        <div className="flex flex-wrap items-center gap-1.5 pt-1">
          <Chip>凭据</Chip>
          <span className="text-[11.5px] text-ink-2">
            {CREDENTIAL_LABELS[view.credential] ?? `未知状态（${view.credential}）`}
          </span>
        </div>
        {view.configuredKnown ? (
          <p className="text-[11.5px] leading-[1.6] text-ink-3">
            {view.configured
              ? "已经配置过：下面的选择器按你已有的配置显示。"
              : "还没有配置过：服务商默认先落在 stepfun 上，不点下面的按钮不写任何东西。"}
          </p>
        ) : (
          <p className="text-[11.5px] leading-[1.6] text-ink-3">
            是否配置过没读到（腿还没答），所以没有默认高亮。
          </p>
        )}
      </div>

      {!view.readable ? (
        <div className="rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">
            配置当前不可读（这不是未配置）：上面那句原文就是宿主的话，先把它修好，下面的控件才谈得上有意义。
          </p>
        </div>
      ) : (
        <>
          {/* 第一步：服务商。 */}
          <StepCard n={1} title="选服务商" done={rowExists}>
            {presets === null ? (
              <p className="text-[12px] leading-[1.7] text-ink-3">
                {view.presetsKnown
                  ? "宿主答了预设清单，但清单是空的。"
                  : "预设清单没读到（腿还没答）。没有清单就不列服务商名，也不编名字。"}
              </p>
            ) : presets.length === 0 ? (
              <p className="text-[12px] leading-[1.7] text-ink-3">宿主答了：内置预设清单是空的。</p>
            ) : (
              <div className="flex flex-wrap items-center gap-2">
                <Dropdown
                  label="服务商"
                  onSelect={(next) => {
                    setPickedProvider(next);
                    setModel("");
                    setManual(false);
                  }}
                  options={presets.map((name) => ({
                    value: name,
                    note:
                      name === suggested
                        ? "未配置时的默认"
                        : view.providers.some((p) => p.name === name)
                          ? "已建行"
                          : "",
                  }))}
                  value={provider}
                />
                {rowExists ? (
                  <span className="text-[11.5px] text-ink-3">
                    这一行已经在配置里了：建行的口只建不改，所以本步不需要动作。
                  </span>
                ) : (
                  <Button
                    disabled={busy !== null}
                    onClick={() =>
                      run("建服务商行", () =>
                        requestConfigWrite({
                          field: "provider_row_create",
                          provider,
                          ...(model.trim() !== "" ? { model: model.trim() } : {}),
                        }),
                      )
                    }
                    size="sm"
                    variant="primary"
                  >
                    {busy === "建服务商行" ? "写入中..." : `建立 [${provider}] 行`}
                  </Button>
                )}
              </div>
            )}
          </StepCard>

          {/* 第二步：密钥。 */}
          <StepCard n={2} title="录密钥" done={providerRow?.credentialSet === true}>
            <CredentialRow
              busy={busy !== null}
              onSend={(value) => run("录密钥", () => requestCredentialWrite(provider, value))}
              provider={provider}
              row={providerRow}
            />
          </StepCard>

          {/* 第三步：成对写聊天角色。 */}
          <StepCard n={3} title="把聊天角色指过去（成对写）">
            <div className="flex flex-col gap-2">
              <div className="flex flex-wrap items-center gap-2">
                <span className="text-[11.5px] text-ink-3">模型</span>
                <ModelControl
                  manual={manual}
                  model={model}
                  onChange={setModel}
                  onManual={() => setManual(true)}
                  row={providerRow}
                />
              </div>
              <div className="flex flex-wrap items-center gap-2">
                <Button
                  disabled={busy !== null || provider === "" || effectiveModel === ""}
                  onClick={() =>
                    run("成对写", () =>
                      requestConfigWrite({ field: "role_chat_pair", provider, model: effectiveModel }),
                    )
                  }
                  size="sm"
                  variant="primary"
                >
                  {busy === "成对写" ? "写入中..." : "写入 provider + model（一笔）"}
                </Button>
                {(provider === "" || effectiveModel === "") && (
                  <span className="text-[11px] text-ink-3">
                    provider 与 model 两半都要有：只有一半的写入会被校验拒掉，所以这里按钮就不开。
                  </span>
                )}
              </div>
              {view.chatModelKnown ? (
                <p className="text-[11.5px] text-ink-3">
                  当前聊天角色指向：{view.chatProvider === "" ? "（provider 未设）" : view.chatProvider} /{" "}
                  {view.chatModel === "" ? "（model 未设）" : view.chatModel}
                </p>
              ) : (
                <p className="text-[11.5px] text-ink-3">当前聊天角色的两半没读到（腿还没答）。</p>
              )}
            </div>
          </StepCard>

          {/* 回执：每一句都是宿主发的原文。 */}
          {receipts.length > 0 ? (
            <div className="flex flex-col gap-2">
              {receipts.map((r, i) => (
                <div
                  className={cn(
                    "rounded-card border px-3.5 py-2.5",
                    r.refused ? "border-red/40 bg-red-tint" : "border-line bg-surface",
                  )}
                  key={`${r.step}-${i}`}
                >
                  <p className="text-[11px] font-medium text-ink-3">{r.step}</p>
                  <pre className="whitespace-pre-wrap text-[11.5px] leading-[1.7] text-ink-2">{r.text}</pre>
                </div>
              ))}
            </div>
          ) : null}

          <p className="text-[11px] leading-[1.6] text-ink-3">
            写入的生效级别由回执里那句话说，页面不另说。密钥只写不回显：录完之后这里只显示已录入和引用名。
          </p>
        </>
      )}
    </div>
  );
}

/* ---------------------------------------------------------------------------
   Pieces.
   --------------------------------------------------------------------------- */

function StepCard({
  n,
  title,
  done,
  children,
}: {
  n: number;
  title: string;
  done?: boolean;
  children: ReactNode;
}) {
  return (
    <div className="rounded-card border border-line bg-surface px-3.5 py-3">
      <div className="mb-2 flex items-center gap-2">
        <span
          aria-hidden="true"
          className="flex size-4.5 shrink-0 items-center justify-center rounded-full bg-inset font-mono text-[10.5px] text-ink-2"
        >
          {n}
        </span>
        <span className="text-[12.5px] font-medium text-ink">{title}</span>
        {done ? (
          <span className="rounded-full bg-green-tint px-1.5 py-0.5 text-[10px] font-medium text-green">
            已就位
          </span>
        ) : null}
      </div>
      {children}
    </div>
  );
}

/** The credential input. The value is read off the DOM at submit time, the
    field is emptied in the same tick, and nothing about it is ever put into
    React state or rendered again - not even a tail. */
function CredentialRow({
  provider,
  row,
  busy,
  onSend,
}: {
  provider: string;
  row: ConfigViewProvider | null;
  busy: boolean;
  onSend: (value: string) => void;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [empty, setEmpty] = useState(true);

  function submit() {
    const input = inputRef.current;
    if (!input) return;
    const value = input.value;
    if (value === "") return;
    input.value = "";
    setEmpty(true);
    onSend(value);
  }

  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-wrap items-center gap-2">
        <span className="flex h-[27px] items-center gap-1.5 rounded-control bg-inset px-2 shadow-hairline">
          <KeyRound aria-hidden="true" size={12} strokeWidth={1.9} className="text-ink-3" />
          <input
            aria-label={`${provider} 的密钥`}
            className="w-56 min-w-0 bg-transparent font-mono text-[12px] text-ink outline-none placeholder:font-sans placeholder:text-ink-3"
            onChange={(e) => setEmpty(e.currentTarget.value === "")}
            onKeyDown={(e) => {
              if (e.key === "Enter") submit();
            }}
            placeholder={provider === "" ? "先选服务商" : "粘贴密钥（只写不回显）"}
            ref={inputRef}
            type="password"
          />
        </span>
        <Button disabled={busy || empty || provider === ""} onClick={submit} size="sm" variant="primary">
          录入密钥
        </Button>
      </div>
      {row !== null ? (
        <p className="text-[11.5px] leading-[1.6] text-ink-3">
          {row.credentialSet ? "已录入。" : "未录入。"}
          {row.apiKeyRef !== "" ? ` 记录的是引用 ${row.apiKeyRef}。` : ""}
          {row.credentialError !== "" ? ` 读取这条引用时：${row.credentialError}` : ""}
        </p>
      ) : (
        <p className="text-[11.5px] leading-[1.6] text-ink-3">
          这个服务商还没有配置行。先做第 1 步再录密钥（密钥写的是这一行的引用）。
        </p>
      )}
    </div>
  );
}

/** The model control: the provider's known model names when the leg answered
    them, and a hand-entry id otherwise - never an invented name. */
function ModelControl({
  row,
  model,
  onChange,
  manual,
  onManual,
}: {
  row: ConfigViewProvider | null;
  model: string;
  onChange: (next: string) => void;
  manual: boolean;
  onManual: () => void;
}) {
  const models = row !== null && row.models !== null && row.models.length > 0 ? row.models : null;
  if (models !== null && !manual) {
    return (
      <div className="flex flex-wrap items-center gap-2">
        <Dropdown
          label="模型"
          onSelect={(next) => {
            if (next === MANUAL) {
              onManual();
              onChange("");
            } else {
              onChange(next);
            }
          }}
          options={[
            ...models.map((name) => ({ value: name, note: "" })),
            { value: MANUAL, note: "自己写 id" },
          ]}
          value={model === "" ? models[0]! : model}
        />
        <span className="text-[11px] text-ink-3">清单来自配置里这个服务商名下的目录模型。</span>
      </div>
    );
  }
  return (
    <div className="flex flex-wrap items-center gap-2">
      <input
        aria-label="手填模型 id"
        className="h-[27px] w-56 rounded-control bg-inset px-2 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
        onChange={(e) => onChange(e.currentTarget.value)}
        placeholder="手填模型 id"
        value={model}
      />
      <span className="text-[11px] leading-[1.6] text-ink-3">
        {row === null
          ? "这一行还不存在，所以没有目录模型可列；先建行，或者直接手填 id 随建行一起写。"
          : row.models === null
            ? "这个服务商的模型清单没读到（腿还没答），只能手填。"
            : "这个服务商名下还没有目录模型，只能手填 id。"}
      </span>
    </div>
  );
}

/** The hand-entry sentinel: a dropdown row that flips the control to a text
    field instead of selecting anything. Not a value any config may hold. */
const MANUAL = "\u0000manual-model-entry";

/** A small dropdown: button + popover list, dismissed on outside pointerdown
    and on Escape (lib/use-dismiss). No native <select>, because the list
    carries a per-option note column. */
function Dropdown({
  label,
  value,
  options,
  onSelect,
}: {
  label: string;
  value: string;
  options: readonly { value: string; note: string }[];
  onSelect: (next: string) => void;
}) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  useDismiss(ref, open, () => setOpen(false));
  const current = options.find((o) => o.value === value);
  return (
    <div className="relative" ref={ref}>
      <button
        aria-expanded={open}
        aria-label={label}
        className="flex h-[27px] items-center gap-1.5 rounded-control bg-inset px-2.5 text-[12px] text-ink shadow-hairline transition-colors duration-150 hover:bg-hover"
        onClick={() => setOpen((v) => !v)}
        type="button"
      >
        <span className="font-mono">{current?.value ?? value}</span>
        <ChevronDown aria-hidden="true" size={12} strokeWidth={2} className="text-ink-3" />
      </button>
      {open ? (
        <div
          className="absolute left-0 top-full z-40 mt-1 min-w-[220px] rounded-card border border-line bg-surface p-1 shadow-overlay"
          style={{ animation: "pop-in 160ms var(--ease-out-strong) both" }}
        >
          {options.map((option) => (
            <button
              className="flex w-full items-center gap-2 rounded-chip px-2 py-1.5 text-left transition-colors duration-100 hover:bg-hover"
              key={option.value}
              onClick={() => {
                onSelect(option.value);
                setOpen(false);
              }}
              type="button"
            >
              <span className="min-w-0 flex-1 truncate font-mono text-[12px] text-ink">{option.value}</span>
              {option.note !== "" ? (
                <span className="shrink-0 text-[10.5px] text-ink-3">{option.note}</span>
              ) : null}
            </button>
          ))}
        </div>
      ) : null}
    </div>
  );
}
