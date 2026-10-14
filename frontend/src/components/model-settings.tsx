/* ============================================================================
   模型 settings page - v2 (owner review 2026-10-14: the v1 page was the
   shortest path that made the wiring work, not a product surface).
   ----------------------------------------------------------------------------
   WHAT CHANGED FROM v1. The interface is untouched (313-panel-ask.md); this is
   the product cut of the page:

     - copy: every visible word is plain language. The spec's internal
       vocabulary is gone from the screen, statuses read 未读取 / 未配置 /
       已配置, and a write's answer is one translated sentence instead of a
       receipt with key paths (lib/config-feedback.ts owns all of it).
     - structure: 供应商 is a LIST (name, base URL, key state, actions) with
       「添加供应商」, not a three-step form; the chat role is one provider +
       model pair with 「保存」.
     - validation: the pair rule is stated at SAVE time (a pre-disabled button
       is a black box to the person looking at it).
     - the add form carries an optional model id because a model can only be
       created together with a provider row (config.go's CreateProviderRow is
       the only creator, and an existing row is refused) - a provider added
       without one cannot have a model added later through the panel.

   THE VIEW IS PURE (ModelSettingsView takes the view plus two write callbacks)
   so scripts/render-settings.tsx can paint every state from a fixture and the
   owner can review the design without a backend. The container below binds it
   to the read gate and the settings door.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useRef, useState } from "react";
import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { ConfigReadGate } from "@/components/config-read-state";
import {
  COPY,
  configStatus,
  keyLabel,
  translateWriteReply,
  type Feedback,
} from "@/lib/config-feedback";
import type { ConfigView, ConfigViewProvider } from "@/lib/config-view";
import { requestConfigWrite, requestCredentialWrite } from "@/lib/panel";
import {
  EmptyState,
  FeedbackLine,
  PairFields,
  Picker,
  SettingsCard,
  StateDot,
  StateTag,
} from "@/components/settings-controls";

/** One settings write, in the shape the door takes. */
export interface ConfigWireWrite {
  field: string;
  provider?: string;
  model?: string;
  value?: string;
}

export interface ModelSettingsViewProps {
  view: ConfigView;
  /** One whitelisted write, answered as the sentence to show. */
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  /** The credential write: value rides its own key, never state. */
  performCredential: (provider: string, value: string) => Promise<Feedback>;
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
            performCredential={async (provider, value) =>
              translateWriteReply(await requestCredentialWrite(provider, value))
            }
            performWrite={async (write) => translateWriteReply(await requestConfigWrite(write))}
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
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
}) {
  const rows = view.providers;
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
              busy={busy}
              editing={editing === row.name}
              key={row.name}
              onEditing={() => onEditing(editing === row.name ? null : row.name)}
              onFeedback={onFeedback}
              performCredential={performCredential}
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
              run={run}
              onDone={() => setAddOpen(false)}
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

function ProviderRow({
  row,
  busy,
  editing,
  onEditing,
  onFeedback,
  performWrite,
  performCredential,
  run,
}: {
  row: ConfigViewProvider;
  busy: string | null;
  editing: boolean;
  onEditing: () => void;
  onFeedback: (next: Feedback | null) => void;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  performCredential: (provider: string, value: string) => Promise<Feedback>;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
}) {
  const key = keyLabel(row);

  return (
    <div className="border-b border-line last:border-b-0">
      <div className="flex items-center gap-3 py-2.5">
        <StateDot tone={key.tone} />
        <span className="shrink-0 text-[12.5px] font-medium text-ink">{row.name}</span>
        <span className="min-w-0 flex-1 truncate font-mono text-[11px] text-ink-3">
          {row.hasBaseURL ? "已设置基址" : "默认基址"}
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
  busy,
  onFeedback,
  performWrite,
  performCredential,
  run,
}: {
  row: ConfigViewProvider;
  busy: string | null;
  onFeedback: (next: Feedback | null) => void;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  performCredential: (provider: string, value: string) => Promise<Feedback>;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
}) {
  const [baseURL, setBaseURL] = useState("");
  const [replacingKey, setReplacingKey] = useState(false);
  const keyInputRef = useRef<HTMLInputElement>(null);

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
    </div>
  );
}

function AddProviderForm({
  available,
  busy,
  onCancel,
  onDone,
  onFeedback,
  performWrite,
  performCredential,
  run,
}: {
  available: string[] | null;
  busy: string | null;
  onCancel: () => void;
  onDone: () => void;
  onFeedback: (next: Feedback | null) => void;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  performCredential: (provider: string, value: string) => Promise<Feedback>;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
}) {
  const [provider, setProvider] = useState("");
  const [model, setModel] = useState("");
  const keyInputRef = useRef<HTMLInputElement>(null);

  function submit() {
    if (provider === "") {
      onFeedback({ tone: "error", text: "请先选择供应商" });
      return;
    }
    const keyValue = keyInputRef.current?.value ?? "";
    if (keyInputRef.current) keyInputRef.current.value = "";
    void run("add-provider", async () => {
      const created = await performWrite({
        field: "provider_row_create",
        provider,
        ...(model.trim() !== "" ? { model: model.trim() } : {}),
      });
      if (created.tone !== "ok") return created;
      if (keyValue === "") return created;
      const keyed = await performCredential(provider, keyValue);
      return keyed.tone === "ok"
        ? created
        : { tone: "error" as const, text: `供应商已添加，但密钥没有保存：${keyed.text}` };
    }).then((fb) => {
      onFeedback(fb);
      if (fb.tone === "ok") onDone();
    });
  }

  return (
    <div
      className="mt-1 flex flex-col gap-2.5 rounded-control bg-inset p-3"
      style={{ animation: "fade-up 200ms var(--ease-out-strong) both" }}
    >
      <div className="flex flex-wrap items-center gap-2">
        <span className="w-16 shrink-0 text-[11.5px] text-ink-3">供应商</span>
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
      </div>
      <div className="flex flex-wrap items-center gap-2">
        <span className="w-16 shrink-0 text-[11.5px] text-ink-3">模型 id</span>
        <input
          aria-label="模型 id（选填）"
          className="h-[30px] w-[220px] rounded-control bg-surface px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
          onChange={(e) => setModel(e.currentTarget.value)}
          placeholder="选填"
          value={model}
        />
        <span className="text-[11px] text-ink-3">现在填上，之后就能直接选</span>
      </div>
      <div className="flex flex-wrap items-center gap-2">
        <span className="w-16 shrink-0 text-[11.5px] text-ink-3">密钥</span>
        <input
          aria-label="密钥（选填）"
          className="h-[30px] w-[260px] rounded-control bg-surface px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
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
        text: chatProvider === "" ? COPY.chatPairMissing : noModels ? "这家供应商还没有模型" : COPY.pairMissing,
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

