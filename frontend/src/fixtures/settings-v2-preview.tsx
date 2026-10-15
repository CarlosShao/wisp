/* ============================================================================
   Settings v2 preview - fixture data, real components (2026-10-14).
   ----------------------------------------------------------------------------
   WHY THIS EXISTS. The model / voice settings pages only render their real
   shape when a host answers config.get; a browser has no host, so the product
   page shows its empty state and there is nothing to look at. The v2 rework is
   a DESIGN deliverable ("附改前/改后截图"), so the two views are mounted here
   against fixture READS - and the WRITE path is the real one, which without a
   host answers "未连接到本机服务，未保存" (the same sentence the product shows
   in a browser). Nothing on this page pretends a save succeeded.

   NOT PART OF THE PRODUCT BUILD: vite.config.ts builds index.html only; this
   file is reachable through the dev server alone
   (npm run dev → /fixtures/settings-v2-preview.html).

   The model ids below are PLACEHOLDERS on purpose - the repository has zero
   real stepfun ids (owner will supply the list), and a screenshot is not a
   reason to invent one. The page header says so.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { createRoot } from "react-dom/client";
import type { ReactNode } from "react";
import { ModelSettingsView } from "@/components/model-settings";
import { VoiceSettingsView } from "@/components/voice-settings";
import type { ConfigView } from "@/lib/config-view";
import { performConfigWrite, performCredentialWrite } from "@/lib/config-actions";
import "../styles/theme.css";

/* The real path, bound the way the containers bind it (lib/config-actions.ts).
   A browser has no host, so every save and the fetch report 未连接到本机服务 -
   which is the truth here. */
const performWrite = performConfigWrite;
const performCredential = performCredentialWrite;
const noReload = () => undefined;
/* The fetch is the one call the preview replaces with fixture data: without a
   host the real one can only answer 未连接到本机服务, and the fetched-list UI
   (tick, write selected) is exactly what needs eyes on it. The ids are
   placeholders, said so in the page header. */
const performFetchModels = async () => ({
  ok: true as const,
  models: ["示例模型-1", "示例模型-2", "示例模型-3"],
});

/* ---------------------------------------------------------------------------
   Fixture reads. Key sets and value shapes follow config_313_test.go's envelope.
   --------------------------------------------------------------------------- */

const UNCONFIGURED: ConfigView = {
  readable: true,
  credential: "no_ref_declared",
  presetsKnown: true,
  presets: ["anthropic", "deepseek", "openai", "stepfun"],
  configuredKnown: true,
  configured: false,
  chatProvider: "",
  chatModel: "",
  chatModelKnown: true,
  providers: [],
  voice: {
    cloudASRChain: null,
    cloudTTSChain: null,
    realtimeEnabled: false,
    realtimeProvider: "",
    realtimeModel: "",
  },
};

const CONFIGURED: ConfigView = {
  readable: true,
  credential: "all_recorded",
  presetsKnown: true,
  presets: ["anthropic", "deepseek", "openai", "stepfun"],
  configuredKnown: true,
  configured: true,
  chatProvider: "stepfun",
  chatModel: "示例模型-1",
  chatModelKnown: true,
  providers: [
    {
      name: "stepfun",
      apiKeyRef: "dpapi:9f2c41d0",
      hasBaseURL: false,
      modelCount: 2,
      models: ["示例模型-1", "示例模型-2"],
      credentialSet: true,
      credentialError: "",
    },
    {
      name: "deepseek",
      apiKeyRef: "",
      hasBaseURL: true,
      modelCount: 0,
      models: [],
      credentialSet: false,
      credentialError: "",
    },
  ],
  voice: {
    cloudASRChain: ["stepfun/示例模型-1"],
    cloudTTSChain: ["stepfun/示例模型-1", "stepfun/示例模型-2"],
    realtimeEnabled: false,
    realtimeProvider: "stepfun",
    realtimeModel: "示例模型-1",
  },
};

const VOICE_NO_PROVIDERS: ConfigView = {
  ...UNCONFIGURED,
  voice: UNCONFIGURED.voice,
};

const UNREADABLE: ConfigView = {
  ...UNCONFIGURED,
  readable: false,
  configuredKnown: false,
  chatModelKnown: false,
};

/* ---------------------------------------------------------------------------
   Page.
   --------------------------------------------------------------------------- */

export function Block({ title, note, children }: { title: string; note: string; children: ReactNode }) {
  return (
    <section className="flex flex-col gap-2">
      <header>
        <h2 className="text-[14px] font-semibold text-ink">{title}</h2>
        <p className="text-[11.5px] leading-[1.6] text-ink-3">{note}</p>
      </header>
      <div className="rounded-card border border-line bg-page p-5">{children}</div>
    </section>
  );
}

export function Preview() {
  return (
    <div className="flex min-h-screen flex-col gap-8 bg-page px-8 py-8 text-ink">
      <header className="max-w-[820px]">
        <h1 className="text-[17px] font-semibold text-ink">模型 / 语音设置页 · v2 夹具预览</h1>
        <p className="mt-1 text-[12px] leading-[1.7] text-ink-2">
          这一页用夹具数据挂真组件；写入走真路径，浏览器里没有本机服务，所以点「保存」会得到真实的
          「未连接到本机服务」；「拉取模型列表」是本页唯一替换成夹具的调用（否则看不到勾选写入那一块）。
          模型名是占位示例（仓里没有真实清单，由配置读取后显示真名单）。
        </p>
      </header>

      <div className="mx-auto flex w-full max-w-[820px] flex-col gap-8">
        <Block note="还没有供应商：引导文案 + 主按钮。" title="模型页 · 未配置">
          <ModelSettingsView
            onReload={noReload}
            performCredential={performCredential}
            performFetchModels={performFetchModels}
            performWrite={performWrite}
            view={UNCONFIGURED}
          />
        </Block>

        <Block note="两家供应商（一家密钥已保存并有两个模型、一家只设了基址）、聊天模型已选择。" title="模型页 · 已配置">
          <ModelSettingsView
            onReload={noReload}
            performCredential={performCredential}
            performFetchModels={performFetchModels}
            performWrite={performWrite}
            view={CONFIGURED}
          />
        </Block>

        <Block note="同上的深色主题，验证状态点与卡片在暗色下的可读性。" title="模型页 · 已配置（深色）">
          <div data-theme="dark" className="rounded-card bg-page p-1">
            <ModelSettingsView
              onReload={noReload}
              performCredential={performCredential}
              performFetchModels={performFetchModels}
              performWrite={performWrite}
              view={CONFIGURED}
            />
          </div>
        </Block>

        <Block note="读不到配置时的状态（可刷新重试）。" title="模型页 · 未读取">
          <ModelSettingsView
            onReload={noReload}
            performCredential={performCredential}
            performFetchModels={performFetchModels}
            performWrite={performWrite}
            view={UNREADABLE}
          />
        </Block>

        <Block note="听写一段、朗读两段的回退链；实时语音已指向一家供应商与模型。" title="语音页 · 已配置">
          <VoiceSettingsView onReload={noReload} performWrite={performWrite} view={CONFIGURED} />
        </Block>

        <Block note="供应商还没有添加时的引导。" title="语音页 · 无供应商">
          <VoiceSettingsView onReload={noReload} performWrite={performWrite} view={VOICE_NO_PROVIDERS} />
        </Block>
      </div>
    </div>
  );
}

createRoot(document.getElementById("root")!).render(<Preview />);
