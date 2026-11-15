/* ============================================================================
   Settings modal shell (?harness=1 demo, 2026-09-26)
   ----------------------------------------------------------------------------
   遮罩 + 居中 680px 大卡（max-h 80vh）。THE SHELL IS ALL THIS FILE IS: the
   content is the real <ConfigScreen /> (src/components/config-screen.tsx)
   rendered verbatim - zero content copied, so its two live controls (the
   --panel-alpha slider and the 全开/关闭 motion switch) keep working exactly
   as they do on the standalone screen. The scrim rides Tailwind's built-in
   black at 35% (bg-black/35), the card speaks bg-surface / border-line /
   shadow-overlay.

   PROPS CONTRACT:

     open: boolean
         False renders nothing (return null), true mounts the overlay. The
         mounter owns the flag - there is no internal state, so re-opening
         replays the fade-up like a fresh mount.
     onClose: () => void
         Called by the X button and by a click on the scrim itself (guarded
         to the backdrop, clicks inside the card do not close). The Esc KEY
         is not wired here; the mounter owns the keyboard.

   Known cosmetic note, left for the assembler/owner: ConfigScreen renders
   its own 「设置」 h2 and intro paragraph, which sit under this shell's
   header title. Removing the duplicate would mean editing the real screen,
   which is out of this file's lane - the shell stays content-zero.

   ============================================================================ */

import { X } from "lucide-react";
import { ConfigScreen } from "@/components/config-screen";

export interface SettingsModalProps {
  open: boolean;
  onClose: () => void;
}

export function SettingsModal({ open, onClose }: SettingsModalProps) {
  if (!open) {
    return null;
  }

  return (
    <div
      aria-label="设置"
      aria-modal="true"
      className="fixed inset-0 z-40 flex items-center justify-center bg-black/35"
      onClick={(event) => {
        if (event.target === event.currentTarget) onClose();
      }}
      role="dialog"
      style={{ animation: "fade-in 200ms var(--ease-out-strong) both" }}
    >
      <div
        className="flex max-h-[80vh] w-[680px] flex-col overflow-hidden rounded-card border border-line bg-surface shadow-overlay"
        role="document"
        style={{ animation: "fade-up 300ms var(--ease-out-strong) both" }}
      >
        <header className="flex shrink-0 items-center justify-between gap-3 border-b border-line py-2.5 pl-5 pr-3">
          <p className="text-[13.5px] font-semibold text-ink">设置</p>
          <button
            aria-label="关闭设置"
            className="flex size-7 items-center justify-center rounded-full text-ink-3 transition-colors duration-100 hover:bg-hover hover:text-ink"
            onClick={onClose}
            type="button"
          >
            <X aria-hidden="true" className="size-4" />
          </button>
        </header>
        <div className="min-h-0 flex-1 overflow-y-auto px-5 py-4">
          <ConfigScreen />
        </div>
      </div>
    </div>
  );
}
