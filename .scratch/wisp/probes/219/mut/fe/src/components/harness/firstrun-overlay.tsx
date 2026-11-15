/* ============================================================================
   First-run guide overlay (?harness=1 demo, 2026-09-26)
   ----------------------------------------------------------------------------
   遮罩 + 居中 560px 卡，内容是真正的 <FirstrunScreen />（src/components/
   firstrun-screen.tsx）逐 props 透传 - 零复制：三步轨、ProgressRing 下载卡、
   目录授权行、掩码密钥行全部原样工作。The shell adds exactly one thing the
   screen does not have: a 「跳过全部」 text button in the card's top-right
   corner, wired to onClose. The scrim rides Tailwind's built-in black at 35%
   (bg-black/35).

   PROPS CONTRACT:

     FirstrunOverlayProps = FirstrunScreenProps & { onClose: () => void }

       - every FirstrunScreen field is forwarded untouched (title/subtitle/
         steps/currentStep/modelRows/progress/downloadStatus/downloadNote/
         dirs/dirTitle/dirNote/dirPlaceholder/dirAddLabel/keyLabel/keyNote/
         keyMasked/dpapiNote + the four callbacks onRemoveDir/onAddDir/
         onSkipStep/onNext). The step STATE stays with the mounter, exactly
         as the screen's own header contract demands.
       - onClose: the 「跳过全部」 button. There is deliberately NO
         backdrop-click close: a stray click must not abandon a first-run
         guide; leaving is an explicit act (跳过全部).

   DEMO DATA: RB_FIRSTRUN at the bottom of this file, typed as
   FirstrunScreenProps so it spreads directly:
     <FirstrunOverlay {...RB_FIRSTRUN} onClose={...} />
   Same content as the fixtures' SHOWCASE_FIRSTRUN (this file must not touch
   src/fixtures/harness.ts, so the demo values live here).
   ============================================================================ */

import { FirstrunScreen, type FirstrunScreenProps } from "@/components/firstrun-screen";

export type FirstrunOverlayProps = FirstrunScreenProps & {
  /** The 「跳过全部」 exit. The only way out of the overlay. */
  onClose: () => void;
};

export function FirstrunOverlay({ onClose, ...screen }: FirstrunOverlayProps) {
  return (
    <div
      aria-label="首次启动引导"
      aria-modal="true"
      className="fixed inset-0 z-40 flex items-center justify-center bg-black/35"
      role="dialog"
      style={{ animation: "fade-in 200ms var(--ease-out-strong) both" }}
    >
      <div
        className="w-[560px] rounded-card border border-line bg-surface p-6 shadow-overlay"
        role="document"
        style={{ animation: "fade-up 300ms var(--ease-out-strong) both" }}
      >
        <div className="flex justify-end">
          <button
            className="rounded-full px-2 py-1 text-[11.5px] text-ink-3 transition-colors duration-100 hover:bg-hover hover:text-ink"
            onClick={onClose}
            type="button"
          >
            跳过全部
          </button>
        </div>
        <FirstrunScreen {...screen} />
      </div>
    </div>
  );
}

/* ============================================================================
   RB_FIRSTRUN - the harness fixture for <FirstrunOverlay />, typed as
   FirstrunScreenProps so the assembler can spread it and add onClose.
   Demo only: the production first run is host-triggered (tickets 33/14).
   ============================================================================ */

export const RB_FIRSTRUN: FirstrunScreenProps = {
  title: "欢迎使用一缕",
  subtitle: "三步完成首次配置 · 每步均可跳过，稍后在设置里补",
  steps: ["选择模型", "目录授权", "凭据录入"],
  currentStep: 0,
  modelRows: [
    { name: "asr-streaming-paraformer", desc: "流式识别 · int8", size: "245MB", minisign: "verifying" },
    { name: "asr-offline-sensevoice", desc: "离线识别 · int8", size: "118MB", minisign: "verified" },
    { name: "punc-ct-transformer", desc: "标点恢复 · int8", size: "88MB", minisign: "verified" },
    { name: "kws-zipformer", desc: "唤醒词 · int8", size: "42MB", minisign: "verified" },
  ],
  progress: 0.68,
  downloadStatus: "在装模型",
  downloadNote:
    "哈希一律取自已验签的清单，不从镜像取；下载支持断点续传，可随时跳过，后台继续。",
  dirs: ["C:\\Users\\swq\\Desktop", "D:\\work\\workspace"],
  dirTitle: "授权目录",
  dirNote: "一缕只能读写你授权的目录；工作区必须落在其中（C26 拒绝符号链接外指）。",
  dirPlaceholder: "粘贴或输入目录绝对路径…",
  dirAddLabel: "添加",
  keyLabel: "LLM 凭据",
  keyNote: "用于文本与工具调用；不填也能用离线语音，但对话能力不可用。",
  keyMasked: "sk-************3f9a",
  dpapiNote: "密钥由 Windows DPAPI 按用户加密落盘（SPEC-03）",
};
