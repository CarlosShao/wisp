/* ============================================================================
   模型 settings page - v3 (2026-10-15: ticket 320's write capabilities landed
   on the Go side; see docs/evidence/s1/model-page-v3-brief.md §6/§7).
   ----------------------------------------------------------------------------
   WHAT v3 ADDS, on top of v2's product cut (see the v2 notes in git history):

     - 拉取模型列表: reads one provider's catalog through config.models.list
       (panel.ts's requestProviderModels; the Go half is 320-r2, so until it
       lands the door refuses the name and the page shows that sentence in
       place). The fetched ids are listed for ticking and written ONE PER CALL
       through model_row_create - one key per guarded write.
     - 计费模式: plan / pay-per-token, written through provider_billing. The
       read payload carries no billing field, so the page says 暂未读取 for the
       current value instead of guessing one.
     - 模型目录: the row's models plus a manual "add one id" box
       (model_row_create).
     - 同一厂商两行: the add form has a 自定义接入 mode
       (provider_row_create_custom, value "<protocol>,<base_url>") so
       stepfun-plan and stepfun-api can coexist - the preset path stays for
       built-in names.

   THE VIEW IS PURE (ModelSettingsView takes the view plus three callbacks), so
   the fixture preview (src/fixtures/settings-v2-preview.tsx) can paint every
   state without a backend.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useEffect, useRef, useState } from "react";
import { LoaderCircle, RefreshCw } from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { ConfigReadGate } from "@/components/config-read-state";
import { SegmentedControl } from "@/components/ai-native/segmented-control";
import { COPY, configStatus, keyLabel, type Feedback } from "@/lib/config-feedback";
import type { ConfigView, ConfigViewProvider } from "@/lib/config-view";
import {
  performConfigWrite,
  performCredentialWrite,
  performModelsFetch,
  type ConfigWireWrite,
  type ModelsFetchResult,
} from "@/lib/config-actions";
import {
  EmptyState,
  FeedbackLine,
  PairFields,
  Picker,
  SettingsCard,
  StateDot,
  StateTag,
} from "@/components/settings-controls";

/* voice-settings.tsx still imports ConfigWireWrite from here; the definition
   itself lives in lib/config-actions.ts with the performers that take it. */
export type { ConfigWireWrite, ModelsFetchResult } from "@/lib/config-actions";

export interface ModelSettingsViewProps {
  view: ConfigView;
  /** One whitelisted write, answered as the sentence to show. */
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  /** The credential write: value rides its own key, never state. */
  performCredential: (provider: string, value: string) => Promise<Feedback>;
  /** Fetch one provider's model catalog (config.models.list). */
  performFetchModels: (provider: string) => Promise<ModelsFetchResult>;
  /** Re-read the config (after every accepted write). */
  onReload: () => void;
}

export function ModelSettings({
  heading = "模型",
  lede = "添加供应商、保存密钥，并选择默认聊天模型。",
}: {
  heading?: string;
  lede?: string;
}) {
  return (
    <section aria-label="模型" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>
      <ConfigReadGate>
        {(view, _prose, reload) => (
          <ModelSettingsView
            onReload={reload}
            performCredential={performCredentialWrite}
            performFetchModels={performModelsFetch}
            performWrite={performConfigWrite}
            view={view}
          />
        )}
      </ConfigReadGate>
    </section>
  );
}

/* ---------------------------------------------------------------------------
   The pure view.
   --------------------------------------------------------------------------- */

export function ModelSettingsView({
  view,
  performWrite,
  performCredential,
  performFetchModels,
  onReload,
}: ModelSettingsViewProps) {
  const status = configStatus(view);
  const [busy, setBusy] = useState<string | null>(null);
  const [providersFeedback, setProvidersFeedback] = useState<Feedback | null>(null);
  const [chatFeedback, setChatFeedback] = useState<Feedback | null>(null);
  const [addOpen, setAddOpen] = useState(false);
  const [editing, setEditing] = useState<string | null>(null);
  // Drafts. chatProvider/chatModel start from the read so the current pair is
  // visible and savable as-is; the reload after a save brings back what landed.
  const [chatProvider, setChatProvider] = useState(view.chatProvider);
  const [chatModel, setChatModel] = useState(view.chatModel);

  async function run(key: string, fn: () => Promise<Feedback>): Promise<Feedback> {
    setBusy(key);
    const feedback = await fn();
    setBusy(null);
    onReload();
    return feedback;
  }

  return (
    <div className="mt-2 flex min-w-0 flex-col gap-3">
      {/* 页头一行：状态 + 刷新。 */}
      <div className="flex items-center gap-2">
        <StateTag text={status.label} tone={status.tone} />
        <span className="min-w-0 flex-1" />
        <Button aria-label="刷新" onClick={onReload} size="xs" variant="quiet">
          <RefreshCw aria-hidden="true" size={11} strokeWidth={2} />
          刷新
        </Button>
      </div>

      {!view.readable ? (
        <SettingsCard title="配置">
          <p className="text-[12px] leading-[1.7] text-ink-2">
            暂时读不到配置，请稍后刷新。
          </p>
        </SettingsCard>
      ) : (
        <>
          <ProvidersCard
            busy={busy}
            editing={editing}
            feedback={providersFeedback}
            onEditing={setEditing}
            onFeedback={setProvidersFeedback}
            performCredential={performCredential}
            performFetchModels={performFetchModels}
            performWrite={performWrite}
            run={run}
            setAddOpen={setAddOpen}
            addOpen={addOpen}
            view={view}
          />

          <ChatModelCard
            busy={busy}
            chatModel={chatModel}
            chatProvider={chatProvider}
            feedback={chatFeedback}
            onFeedback={setChatFeedback}
            onModel={setChatModel}
            onProvider={(next) => {
              setChatProvider(next);
              setChatModel("");
            }}
            performWrite={performWrite}
            run={run}
            view={view}
          />
        </>
      )}
    </div>
  );
}

/* ---------------------------------------------------------------------------
   供应商.
   --------------------------------------------------------------------------- */

function ProvidersCard({
  view,
  busy,
  addOpen,
  setAddOpen,
  editing,
  onEditing,
  feedback,
  onFeedback,
  performWrite,
  performCredential,
  performFetchModels,
  run,
}: {
  view: ConfigView;
  busy: string | null;
  addOpen: boolean;
  setAddOpen: (next: boolean) => void;
  editing: string | null;
  onEditing: (next: string | null) => void;
  feedback: Feedback | null;
  onFeedback: (next: Feedback | null) => void;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  performCredential: (provider: string, value: string) => Promise<Feedback>;
  performFetchModels: (provider: string) => Promise<ModelsFetchResult>;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
}) {
  const rows = view.providers;
  const [autoFetch, setAutoFetch] = useState<string | null>(null);
  const availablePresets =
    view.presetsKnown && view.presets !== null
      ? view.presets.filter((name) => !rows.some((p) => p.name === name))
      : null;

  return (
    <SettingsCard title="供应商" description="密钥只保存在本机；界面上不会再显示密钥内容。">
      {rows.length === 0 && !addOpen ? (
        <EmptyState actionLabel={COPY.addProvider} onAction={() => setAddOpen(true)} text={COPY.noProvidersYet} />
      ) : (
        <div className="flex flex-col">
          {rows.map((row) => (
            <ProviderRow
              autoFetch={autoFetch === row.name}
              busy={busy}
              editing={editing === row.name}
              key={row.name}
              onAutoFetchHandled={() => setAutoFetch(null)}
              onEditing={() => onEditing(editing === row.name ? null : row.name)}
              onFeedback={onFeedback}
              performCredential={performCredential}
              performFetchModels={performFetchModels}
              performWrite={performWrite}
              row={row}
              run={run}
            />
          ))}

          {addOpen ? (
            <AddProviderForm
              available={availablePresets}
              busy={busy}
              onCancel={() => setAddOpen(false)}
              onFeedback={onFeedback}
              performCredential={performCredential}
              performWrite={performWrite}
              presets={view.presets}
              run={run}
              onDone={(name) => {
                // ADD, THEN FETCH. A brand-new row has nothing to pick from, and
                // making the operator open the row and press 拉取 themselves is
                // the two-step shape this form exists to remove: the fetch runs
                // for the row that was just created and its result lands in that
                // row's panel.
                setAddOpen(false);
                onEditing(name);
                setAutoFetch(name);
              }}
            />
          ) : (
            <div className="pt-3">
              <Button onClick={() => setAddOpen(true)} size="sm" variant="secondary">
                {COPY.addProvider}
              </Button>
            </div>
          )}
        </div>
      )}
      {feedback ? (
        <div className="pt-2">
          <FeedbackLine feedback={feedback} />
        </div>
      ) : null}
    </SettingsCard>
  );
}

/** The fetch's four states; `ready` holds the ids the host listed. */
type FetchState =
  | { phase: "idle" }
  | { phase: "loading" }
  | { phase: "error"; feedback: Feedback }
  | { phase: "ready"; models: string[] };

function ProviderRow({
  row,
  busy,
  editing,
  autoFetch,
  onAutoFetchHandled,
  onEditing,
  onFeedback,
  performWrite,
  performCredential,
  performFetchModels,
  run,
}: {
  row: ConfigViewProvider;
  busy: string | null;
  editing: boolean;
  autoFetch: boolean;
  onAutoFetchHandled: () => void;
  onEditing: () => void;
  onFeedback: (next: Feedback | null) => void;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  performCredential: (provider: string, value: string) => Promise<Feedback>;
  performFetchModels: (provider: string) => Promise<ModelsFetchResult>;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
}) {
  const key = keyLabel(row);
  const [fetch, setFetch] = useState<FetchState>({ phase: "idle" });

  // The row's own 拉取模型列表 button: opens the detail panel so the result is
  // visible where it will be written from, then runs the fetch. The loading
  // state is the button's (Go answers synchronously, budget <= 10s).
  function startFetch() {
    if (fetch.phase === "loading" || busy !== null) return;
    if (!editing) onEditing();
    setFetch({ phase: "loading" });
    void performFetchModels(row.name).then((result) => {
      setFetch(result.ok ? { phase: "ready", models: result.models } : { phase: "error", feedback: result.feedback });
    });
  }

  // AUTO-FETCH AFTER ADD: the form hands a freshly created row's name up, and
  // this row is the one that then asks the host for its catalog - so "add a
  // provider, then fetch its models" is one flow, not two. The ref keeps the
  // effect to a single firing per handed-over name (the row re-renders as the
  // fetch progresses).
  const startFetchRef = useRef(startFetch);
  startFetchRef.current = startFetch;
  const autoFired = useRef(false);
  useEffect(() => {
    if (!autoFetch || autoFired.current) return;
    autoFired.current = true;
    onAutoFetchHandled();
    startFetchRef.current();
  }, [autoFetch, onAutoFetchHandled]);

  return (
    <div className="border-b border-line last:border-b-0">
      <div className="flex flex-wrap items-center gap-x-3 gap-y-1.5 py-2.5">
        <StateDot tone={key.tone} />
        <span className="shrink-0 text-[12.5px] font-medium text-ink">{row.name}</span>
        {row.billing === "plan" || row.billing === "pay-per-token" ? (
          <span className="shrink-0 rounded-chip bg-inset px-1.5 py-0.5 text-[10.5px] text-ink-2">
            {row.billing === "plan" ? "套餐" : "按量"}
          </span>
        ) : null}
        {row.protocol !== "" ? (
          <span className="shrink-0 font-mono text-[10.5px] text-ink-3">{row.protocol}</span>
        ) : null}
        <span className="min-w-0 flex-1 truncate font-mono text-[11px] text-ink-3" title={row.baseURL}>
          {row.baseURL !== "" ? row.baseURL : "默认基址（预设）"}
        </span>
        <span className="shrink-0 text-[11px] text-ink-3">{key.label}</span>
        <span className="shrink-0 text-[11px] text-ink-3">
          {row.modelCount > 0 ? `${row.modelCount} 个模型` : "暂无模型"}
        </span>
        <div className="flex shrink-0 items-center gap-0.5">
          <Button
            onClick={() => onFeedback({ tone: "neutral", text: COPY.testConnectionUnavailable })}
            size="xs"
            variant="quiet"
          >
            测试连接
          </Button>
          <Button
            aria-busy={fetch.phase === "loading"}
            disabled={fetch.phase === "loading" || busy !== null}
            onClick={startFetch}
            size="xs"
            variant="quiet"
          >
            {fetch.phase === "loading" ? (
              <span className="flex items-center gap-1">
                <LoaderCircle aria-hidden="true" className="animate-spin" size={11} strokeWidth={2} />
                拉取中…
              </span>
            ) : (
              "拉取模型列表"
            )}
          </Button>
          <Button onClick={onEditing} size="xs" variant="quiet">
            {editing ? "收起" : "编辑"}
          </Button>
          <Button
            onClick={() => onFeedback({ tone: "neutral", text: COPY.deleteUnavailable })}
            size="xs"
            variant="quiet"
          >
            删除
          </Button>
        </div>
      </div>

      {editing ? (
        <ProviderEdit
          busy={busy}
          fetch={fetch}
          onDismissFetch={() => setFetch({ phase: "idle" })}
          onFeedback={onFeedback}
          performCredential={performCredential}
          performWrite={performWrite}
          row={row}
          run={run}
        />
      ) : null}
    </div>
  );
}

function ProviderEdit({
  row,
  fetch,
  busy,
  onDismissFetch,
  onFeedback,
  performWrite,
  performCredential,
  run,
}: {
  row: ConfigViewProvider;
  fetch: FetchState;
  busy: string | null;
  onDismissFetch: () => void;
  onFeedback: (next: Feedback | null) => void;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  performCredential: (provider: string, value: string) => Promise<Feedback>;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
}) {
  const [baseURL, setBaseURL] = useState("");
  const [replacingKey, setReplacingKey] = useState(false);
  const [manualModel, setManualModel] = useState("");
  // Ticked ids from the last fetch. New ids start ticked; ids the row already
  // lists are shown as 已有 and cannot be ticked again.
  const [ticked, setTicked] = useState<ReadonlySet<string>>(new Set());
  const keyInputRef = useRef<HTMLInputElement>(null);

  const existing = new Set(row.models ?? []);
  const fetchedNew = fetch.phase === "ready" ? fetch.models.filter((id) => !existing.has(id)) : [];

  function saveKey() {
    const input = keyInputRef.current;
    if (!input) return;
    const value = input.value;
    if (value === "") {
      onFeedback({ tone: "error", text: "请先填入密钥" });
      return;
    }
    // The value leaves the DOM the instant it is sent; nothing puts it into
    // state, and nothing renders it again.
    input.value = "";
    void run(`credential:${row.name}`, () => performCredential(row.name, value)).then((fb) => {
      onFeedback(fb);
      if (fb.tone === "ok") setReplacingKey(false);
    });
  }

  function setBilling(value: "plan" | "pay-per-token") {
    void run(`billing:${row.name}`, () =>
      performWrite({ field: "provider_billing", provider: row.name, value }),
    ).then(onFeedback);
  }

  function addOneModel(id: string) {
    const model = id.trim();
    if (model === "") {
      onFeedback({ tone: "error", text: "请先填入模型 id" });
      return;
    }
    void run(`models:${row.name}`, () =>
      performWrite({ field: "model_row_create", provider: row.name, model }),
    ).then((fb) => {
      onFeedback(fb);
      if (fb.tone === "ok") setManualModel("");
    });
  }

  /** One write per selected id - model_row_create is a single-entry step, and
      the answered receipt is per call. The summary line keeps the receipt's
      own meaning (saved + restart tier). */
  function addFetchedModels() {
    const ids = fetchedNew.filter((id) => ticked.has(id));
    if (ids.length === 0) {
      onFeedback({ tone: "error", text: "先勾选要添加的模型" });
      return;
    }
    void run(`models:${row.name}`, async () => {
      let saved = 0;
      let firstError: Feedback | null = null;
      for (const id of ids) {
        const fb = await performWrite({ field: "model_row_create", provider: row.name, model: id });
        if (fb.tone === "ok") saved += 1;
        else if (firstError === null) firstError = fb;
      }
      if (firstError === null) {
        onDismissFetch();
        return { tone: "ok" as const, text: `已保存 ${saved} 个模型，重启后生效` };
      }
      return {
        tone: "error" as const,
        text: saved > 0 ? `${saved} 个已保存；其余没有写成：${firstError.text}` : firstError.text,
      };
    }).then(onFeedback);
  }

  return (
    <div className="flex flex-col gap-3 pb-3 pl-5">
      <div className="flex flex-wrap items-center gap-2">
        <span className="w-16 shrink-0 text-[11.5px] text-ink-3">基址</span>
        <input
          aria-label={`${row.name} 的接口基址`}
          className="h-[30px] w-[320px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
          onChange={(e) => setBaseURL(e.currentTarget.value)}
          placeholder="https://…"
          value={baseURL}
        />
        <Button
          disabled={busy !== null}
          onClick={() => {
            if (baseURL.trim() === "") {
              onFeedback({ tone: "error", text: "基址不能为空" });
              return;
            }
            void run(`base:${row.name}`, () =>
              performWrite({ field: "provider_base_url", provider: row.name, value: baseURL.trim() }),
            ).then(onFeedback);
          }}
          size="xs"
          variant="secondary"
        >
          保存基址
        </Button>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <span className="w-16 shrink-0 text-[11.5px] text-ink-3">密钥</span>
        {row.credentialSet && !replacingKey ? (
          <>
            <span className="text-[11.5px] text-ink-2">
              已保存{row.apiKeyRef !== "" ? `（${row.apiKeyRef}）` : ""}
            </span>
            <Button onClick={() => setReplacingKey(true)} size="xs" variant="secondary">
              更换
            </Button>
          </>
        ) : (
          <>
            <input
              aria-label={`${row.name} 的密钥`}
              className="h-[30px] w-[260px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
              onKeyDown={(e) => {
                if (e.key === "Enter") saveKey();
              }}
              placeholder="粘贴密钥（保存后不再显示）"
              ref={keyInputRef}
              type="password"
            />
            <Button disabled={busy !== null} onClick={saveKey} size="xs" variant="secondary">
              保存密钥
            </Button>
            {replacingKey ? (
              <Button onClick={() => setReplacingKey(false)} size="xs" variant="quiet">
                取消
              </Button>
            ) : null}
          </>
        )}
        <Button
          onClick={() => onFeedback({ tone: "neutral", text: COPY.testKeyUnavailable })}
          size="xs"
          variant="quiet"
        >
          测试密钥
        </Button>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <span className="w-16 shrink-0 text-[11.5px] text-ink-3">计费模式</span>
        <Button
          disabled={busy !== null}
          onClick={() => setBilling("plan")}
          size="xs"
          variant={row.billing === "plan" ? "primary" : "secondary"}
        >
          套餐
        </Button>
        <Button
          disabled={busy !== null}
          onClick={() => setBilling("pay-per-token")}
          size="xs"
          variant={row.billing === "pay-per-token" ? "primary" : "secondary"}
        >
          按量
        </Button>
        <span className="text-[11px] text-ink-3">
          {row.billing === "plan"
            ? "当前：套餐（点击「按量」即改）"
            : row.billing === "pay-per-token"
              ? "当前：按量（点击「套餐」即改）"
              : "未设置；点击即保存（同一厂商两个计费口各记各的）"}
        </span>
      </div>

      <div className="flex flex-col gap-2">
        <span className="text-[11.5px] text-ink-3">模型目录</span>
        <div className="flex flex-wrap items-center gap-1.5">
          {row.models !== null && row.models.length > 0 ? (
            row.models.map((id) => (
              <span
                className="inline-flex h-5 items-center rounded-chip bg-inset px-2 font-mono text-[11px] text-ink-2"
                key={id}
              >
                {id}
              </span>
            ))
          ) : (
            <span className="text-[11.5px] text-ink-3">
              {row.models === null ? "暂时读不到模型清单" : "暂无模型"}
            </span>
          )}
        </div>

        {fetch.phase === "loading" ? (
          <p className="flex items-center gap-1.5 text-[11.5px] text-ink-3">
            <LoaderCircle aria-hidden="true" className="animate-spin" size={11} strokeWidth={2} />
            正在向上游拉取模型列表…
          </p>
        ) : fetch.phase === "error" ? (
          <div className="flex flex-col gap-1">
            <FeedbackLine feedback={fetch.feedback} />
            <div>
              <Button onClick={onDismissFetch} size="xs" variant="quiet">
                知道了
              </Button>
            </div>
          </div>
        ) : fetch.phase === "ready" ? (
          <div className="flex flex-col gap-2 rounded-control bg-inset p-2.5">
            {fetch.models.length === 0 ? (
              <p className="text-[11.5px] text-ink-3">上游返回的清单是空的。</p>
            ) : fetchedNew.length === 0 ? (
              <p className="text-[11.5px] text-ink-3">上游清单里的模型都已经在目录里了。</p>
            ) : (
              <>
                <p className="text-[11.5px] text-ink-3">
                  拉到 {fetchedNew.length} 个新模型，勾选后写入：
                </p>
                <div className="flex max-h-44 flex-col gap-0.5 overflow-y-auto">
                  {fetchedNew.map((id) => (
                    <label className="flex items-center gap-2 rounded-chip px-1.5 py-1 hover:bg-hover" key={id}>
                      <input
                        checked={ticked.has(id)}
                        className="size-3.5 accent-primary"
                        onChange={(e) => {
                          // Read the DOM flag DURING dispatch. The updater below
                          // is not promised to run here: React only computes it
                          // eagerly when the queue is empty, and by the time a
                          // deferred updater runs, React has already nulled
                          // e.currentTarget - which is exactly the
                          // "Cannot read properties of null (reading 'checked')"
                          // crash the operator hit on the SECOND tick.
                          const isChecked = e.currentTarget.checked;
                          setTicked((current) => {
                            const next = new Set(current);
                            if (isChecked) next.add(id);
                            else next.delete(id);
                            return next;
                          });
                        }}
                        type="checkbox"
                      />
                      <span className="font-mono text-[11.5px] text-ink">{id}</span>
                    </label>
                  ))}
                </div>
                <div className="flex items-center gap-2">
                  <Button disabled={busy !== null} onClick={addFetchedModels} size="xs" variant="primary">
                    写入所选
                  </Button>
                  <Button
                    onClick={() =>
                      setTicked((current) => {
                        const next = new Set(current);
                        for (const id of fetchedNew) next.add(id);
                        return next;
                      })
                    }
                    size="xs"
                    variant="quiet"
                  >
                    全选
                  </Button>
                  <Button onClick={onDismissFetch} size="xs" variant="quiet">
                    取消
                  </Button>
                </div>
              </>
            )}
            {fetch.models.length > 0 && fetchedNew.length === 0 ? (
              <div>
                <Button onClick={onDismissFetch} size="xs" variant="quiet">
                  取消
                </Button>
              </div>
            ) : null}
          </div>
        ) : null}

        <div className="flex flex-wrap items-center gap-2">
          <input
            aria-label="模型 id"
            className="h-[30px] w-[220px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
            onChange={(e) => setManualModel(e.currentTarget.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") addOneModel(manualModel);
            }}
            placeholder="模型 id"
            value={manualModel}
          />
          <Button disabled={busy !== null} onClick={() => addOneModel(manualModel)} size="xs" variant="secondary">
            添加
          </Button>
        </div>
      </div>
    </div>
  );
}

function AddProviderForm({
  available,
  presets,
  busy,
  onCancel,
  onDone,
  onFeedback,
  performWrite,
  performCredential,
  run,
}: {
  available: string[] | null;
  presets: string[] | null;
  busy: string | null;
  onCancel: () => void;
  onDone: (name: string) => void;
  onFeedback: (next: Feedback | null) => void;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  performCredential: (provider: string, value: string) => Promise<Feedback>;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
}) {
  const MODES = ["内置供应商", "自定义接入"] as const;
  const BILLING = ["套餐（plan）", "按量（pay-per-token）"] as const;
  const [selectedMode, setSelectedMode] = useState<(typeof MODES)[number]>("内置供应商");
  const [provider, setProvider] = useState("");
  const [customName, setCustomName] = useState("");
  const [protocol, setProtocol] = useState("");
  const [customBase, setCustomBase] = useState("");
  const [billing, setBilling] = useState("");
  const [model, setModel] = useState("");
  const keyInputRef = useRef<HTMLInputElement>(null);

  // KNOWN ALTERNATE ENDPOINTS. A vendor's plan 口 and its API 口 are different
  // base URLs, and the operator's StepFun case (audio models are only free on the
  // API-billing side) is the reason this list exists: the same key answers a
  // different catalog per endpoint. One entry today, spelled out rather than
  // guessed; the plan-side address is the operator's own and is never invented
  // here.
  const KNOWN_ALT_ENDPOINTS: Record<string, { label: string; base: string }> = {
    stepfun: { label: "按量（API）口", base: "https://api.stepfun.com/v1" },
  };
  const vendorKey = customName.trim().split("-")[0]?.toLowerCase() ?? "";
  const altEndpoint = KNOWN_ALT_ENDPOINTS[vendorKey];
  // "COMMON SECOND ROWS": the vendors that are already configured AND have a
  // known second endpoint. The preset picker above is filtered to vendors the
  // operator has NOT added yet (a preset row cannot be added twice), which reads
  // as "stepfun is missing" to anyone who already has a stepfun row - so the
  // second-billing-entry case is offered here by name instead of hiding behind a
  // picker entry that is deliberately absent.
  const takenAll = (presets ?? []).filter((name) => !(available ?? []).includes(name));
  const takenVendors = takenAll.filter((name) => KNOWN_ALT_ENDPOINTS[name.toLowerCase()]);
  const secondRows = takenVendors.map((name) => {
    const alt = KNOWN_ALT_ENDPOINTS[name.toLowerCase()];
    return { name: `${name}-api`, label: alt.label, base: alt.base, vendor: name };
  });

  function prefillSecondRow(row: { name: string; base: string }) {
    setSelectedMode("自定义接入");
    setCustomName(row.name);
    setProtocol("openai-chat");
    setCustomBase(row.base);
    setBilling(BILLING[1]);
  }

  function submit() {
    const keyValue = keyInputRef.current?.value ?? "";
    const billingValue = billing === "" ? "" : billing === BILLING[0] ? "plan" : "pay-per-token";
    if (selectedMode === "内置供应商") {
      if (provider === "") {
        onFeedback({ tone: "error", text: "请先选择供应商" });
        return;
      }
    } else {
      if (customName.trim() === "") {
        onFeedback({ tone: "error", text: "请先填行名" });
        return;
      }
      if (presets !== null && presets.includes(customName.trim())) {
        onFeedback({ tone: "error", text: "这是内置供应商名，请用「内置供应商」添加" });
        return;
      }
      if (protocol === "") {
        onFeedback({ tone: "error", text: "请先选择协议" });
        return;
      }
      if (customBase.trim() === "") {
        onFeedback({ tone: "error", text: "请先填基址" });
        return;
      }
    }
    if (keyInputRef.current) keyInputRef.current.value = "";
    const name = selectedMode === "内置供应商" ? provider : customName.trim();
    void run("add-provider", async () => {
      const created = await performWrite(
        selectedMode === "内置供应商"
          ? {
              field: "provider_row_create",
              provider: name,
              ...(model.trim() !== "" ? { model: model.trim() } : {}),
            }
          : {
              field: "provider_row_create_custom",
              provider: name,
              value: `${protocol},${customBase.trim()}`,
              ...(model.trim() !== "" ? { model: model.trim() } : {}),
            },
      );
      if (created.tone !== "ok") return created;
      // Billing is its own one-key write (one key per call, settings.go), so a row
      // born as 按量 is two writes: the row, then the leaf that names how it is
      // billed. Failure here still reports the row that did land.
      if (billingValue !== "") {
        const billed = await performWrite({ field: "provider_billing", provider: name, value: billingValue });
        if (billed.tone !== "ok") {
          return { tone: "error" as const, text: `供应商已添加，但计费模式没有写入：${billed.text}` };
        }
      }
      if (keyValue === "") return created;
      const keyed = await performCredential(name, keyValue);
      return keyed.tone === "ok"
        ? created
        : { tone: "error" as const, text: `供应商已添加，但密钥没有保存：${keyed.text}` };
    }).then((fb) => {
      onFeedback(fb);
      if (fb.tone === "ok") onDone(name);
    });
  }

  const rowLabel = "w-16 shrink-0 text-[11.5px] text-ink-3";
  const inputCls =
    "h-[30px] rounded-control bg-surface px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3";

  return (
    <div
      className="mt-1 flex flex-col gap-2.5 rounded-control bg-inset p-3"
      style={{ animation: "fade-up 200ms var(--ease-out-strong) both" }}
    >
      <div className="flex items-center gap-2">
        <SegmentedControl onChange={setSelectedMode} options={MODES} value={selectedMode} />
      </div>

      {selectedMode === "内置供应商" ? (
        <div className="flex flex-col gap-2">
          <div className="flex flex-wrap items-center gap-2">
            <span className={rowLabel}>供应商</span>
            {available === null ? (
              <span className="text-[11.5px] text-ink-3">暂未读取到供应商列表，请稍后刷新</span>
            ) : (
              <Picker
                ariaLabel="供应商"
                className="w-[220px]"
                emptyText="内置供应商都已添加"
                onSelect={setProvider}
                options={available.map((name) => ({ value: name, note: name === "stepfun" ? "推荐" : "" }))}
                placeholder="选择供应商"
                value={provider}
              />
            )}
            <span className="text-[11px] text-ink-3">
              这里只列还没添加过的内置供应商；同一家要第二行请用下面那条
            </span>
          </div>
          {takenAll.length > 0 ? (
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-[11px] text-ink-3">已添加：{takenAll.join("、")}</span>
            </div>
          ) : null}
          {secondRows.length > 0 ? (
            <div className="flex flex-wrap items-center gap-2 rounded-control bg-surface/60 px-2 py-1.5">
              <span className={rowLabel}>常用第二行</span>
              {secondRows.map((row) => (
                <Button key={row.name} onClick={() => prefillSecondRow(row)} size="xs" variant="quiet">
                  加一行 {row.name}（{row.label}）
                </Button>
              ))}
              <span className="text-[11px] text-ink-3">
                同一厂商两行并存：各用各的基址、各拉各的模型清单
              </span>
            </div>
          ) : null}
          {provider !== "" && KNOWN_ALT_ENDPOINTS[provider.toLowerCase()] ? (
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-[11px] text-ink-3">
                {provider} 还有另一个计费口（{KNOWN_ALT_ENDPOINTS[provider.toLowerCase()].label}）
                ——同一厂商两行并存时，两行各有各的基址与模型目录：
              </span>
              <Button
                onClick={() =>
                  prefillSecondRow({
                    name: `${provider}-api`,
                    base: KNOWN_ALT_ENDPOINTS[provider.toLowerCase()].base,
                  })
                }
                size="xs"
                variant="quiet"
              >
                切到自定义并预填这一行
              </Button>
            </div>
          ) : null}
        </div>
      ) : (
        <>
          <div className="flex flex-wrap items-center gap-2">
            <span className={rowLabel}>行名</span>
            <input
              aria-label="行名"
              className={`${inputCls} w-[220px]`}
              onChange={(e) => setCustomName(e.currentTarget.value)}
              placeholder="例如 stepfun-api"
              value={customName}
            />
            <span className="text-[11px] text-ink-3">同一厂商可以两行并存，各用各的基址与密钥</span>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <span className={rowLabel}>协议</span>
            <Picker
              ariaLabel="协议"
              className="w-[220px]"
              onSelect={setProtocol}
              options={[
                { value: "openai-chat", note: "OpenAI 兼容" },
                { value: "openai-responses", note: "OpenAI 响应式" },
                { value: "anthropic", note: "Anthropic" },
              ]}
              placeholder="选择协议"
              value={protocol}
            />
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <span className={rowLabel}>基址</span>
            <input
              aria-label="基址"
              className={`${inputCls} w-[320px]`}
              onChange={(e) => setCustomBase(e.currentTarget.value)}
              placeholder="https://…"
              value={customBase}
            />
          </div>
          {altEndpoint ? (
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-[11px] text-ink-3">
                {vendorKey} 的{altEndpoint.label}已知基址：{altEndpoint.base}
              </span>
              <Button
                onClick={() => {
                  setCustomBase(altEndpoint.base);
                  if (billing === "") setBilling(BILLING[1]);
                }}
                size="xs"
                variant="quiet"
              >
                填入
              </Button>
            </div>
          ) : null}
          <div className="flex flex-wrap items-center gap-2">
            <span className={rowLabel}>计费模式</span>
            <SegmentedControl
              onChange={(next) =>
                setBilling(next === BILLING[0] ? BILLING[0] : next === BILLING[1] ? BILLING[1] : "")
              }
              options={BILLING}
              value={billing === BILLING[0] ? BILLING[0] : billing === BILLING[1] ? BILLING[1] : ""}
            />
            <span className="text-[11px] text-ink-3">选填；同一厂商两个计费口各记各的</span>
          </div>
        </>
      )}

      <div className="flex flex-wrap items-center gap-2">
        <span className={rowLabel}>模型 id</span>
        <input
          aria-label="模型 id（选填）"
          className={`${inputCls} w-[220px]`}
          onChange={(e) => setModel(e.currentTarget.value)}
          placeholder="选填"
          value={model}
        />
        <span className="text-[11px] text-ink-3">选填；添加后会自动拉一次这行的模型清单</span>
      </div>
      <div className="flex flex-wrap items-center gap-2">
        <span className={rowLabel}>密钥</span>
        <input
          aria-label="密钥（选填）"
          className={`${inputCls} w-[260px]`}
          placeholder="粘贴密钥（保存后不再显示）"
          ref={keyInputRef}
          type="password"
        />
      </div>
      <div className="flex items-center gap-2 pl-18">
        <Button disabled={busy !== null} onClick={submit} size="sm" variant="primary">
          {busy === "add-provider" ? COPY.saving : "添加"}
        </Button>
        <Button onClick={onCancel} size="sm" variant="quiet">
          取消
        </Button>
      </div>
    </div>
  );
}

/* ---------------------------------------------------------------------------
   默认聊天模型.
   --------------------------------------------------------------------------- */

function ChatModelCard({
  view,
  busy,
  chatProvider,
  chatModel,
  onProvider,
  onModel,
  feedback,
  onFeedback,
  performWrite,
  run,
}: {
  view: ConfigView;
  busy: string | null;
  chatProvider: string;
  chatModel: string;
  onProvider: (next: string) => void;
  onModel: (next: string) => void;
  feedback: Feedback | null;
  onFeedback: (next: Feedback | null) => void;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
}) {
  const current = view.chatModelKnown
    ? view.chatProvider !== "" && view.chatModel !== ""
      ? `${view.chatProvider} / ${view.chatModel}`
      : "未设置"
    : "暂未读取";

  function save() {
    if (chatProvider === "" || chatModel === "") {
      const row = view.providers.find((p) => p.name === chatProvider);
      const noModels = row != null && (row.models === null || row.models.length === 0);
      onFeedback({
        tone: "error",
        text: chatProvider === ""
          ? COPY.chatPairMissing
          : noModels
            ? "这家供应商还没有模型；在供应商行的「编辑」里用「拉取模型列表」或手动添加"
            : COPY.pairMissing,
      });
      return;
    }
    void run("chat-pair", () =>
      performWrite({ field: "role_chat_pair", provider: chatProvider, model: chatModel }),
    ).then(onFeedback);
  }

  return (
    <SettingsCard title="默认聊天模型" description="对话默认使用的供应商与模型。">
      {view.providers.length === 0 ? (
        <EmptyState text="先添加一个供应商，然后就能选择聊天模型。" />
      ) : (
        <div className="flex flex-col gap-2.5">
          <div className="flex flex-wrap items-center gap-2">
            <PairFields
              manualFallback={false}
              model={chatModel}
              onModel={onModel}
              onProvider={onProvider}
              provider={chatProvider}
              providers={view.providers}
            />
            <Button disabled={busy !== null} onClick={save} size="sm" variant="primary">
              {busy === "chat-pair" ? COPY.saving : COPY.save}
            </Button>
          </div>
          <p className="text-[11.5px] text-ink-3">当前：{current}</p>
          <FeedbackLine feedback={feedback} />
        </div>
      )}
    </SettingsCard>
  );
}
