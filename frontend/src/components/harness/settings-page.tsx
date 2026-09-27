/* ============================================================================
   Settings page - full page, not a modal (owner 2026-09-27).
   ----------------------------------------------------------------------------
   Replaces the old SettingsModal. A left category sidebar (11 sections) and a
   content area on the right. The 外观 section renders the real ConfigScreen
   verbatim (its two live controls keep working); every other section is an
   honest "还没有接到数据" state, never a fabricated control.

   The categories are owner-approved (2026-09-27): 通用 / 外观 / 模型 / 语音 /
   记忆 / 技能 / MCP / 安全 / 隐私 / 成本 / 引导.
   ============================================================================ */

import { X } from "lucide-react";
import { ConfigScreen } from "@/components/config-screen";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { cn } from "@/lib/cn";

const SECTIONS = [
  "通用",
  "外观",
  "模型",
  "语音",
  "记忆",
  "技能",
  "MCP",
  "安全",
  "隐私",
  "成本",
  "引导",
] as const;

export type SettingsSection = (typeof SECTIONS)[number];

export interface SettingsPageProps {
  activeSection: string;
  onSelectSection: (section: SettingsSection) => void;
  onClose: () => void;
}

/** The honest empty state for a section that has no data source yet. */
function UnfedSection({ name }: { name: string }) {
  return (
    <div className="flex flex-col gap-2 px-6 py-8">
      <h2 className="text-[14px] font-semibold text-ink">{name}</h2>
      <p className="text-[12.5px] leading-relaxed text-ink-3">
        这一节还没有接到数据。面板不许自己造内容，所以这里只说明它缺
        C17 的读口与票 35 的推送。能动的设置在「外观」一节。
      </p>
    </div>
  );
}

export function SettingsPage({ activeSection, onSelectSection, onClose }: SettingsPageProps) {
  const active = SECTIONS.includes(activeSection as SettingsSection)
    ? (activeSection as SettingsSection)
    : "通用";

  return (
    <div className="flex h-full min-h-0 flex-1 flex-col bg-page">
      {/* header */}
      <header className="flex h-11 shrink-0 items-center justify-between border-b border-line bg-surface px-4">
        <p className="text-[13px] font-semibold text-ink">设置</p>
        <button
          aria-label="关闭设置"
          className="flex size-7 items-center justify-center rounded-full text-ink-3 transition-colors duration-100 hover:bg-hover hover:text-ink"
          onClick={onClose}
          type="button"
        >
          <X aria-hidden="true" className="size-4" />
        </button>
      </header>

      <div className="flex min-h-0 flex-1 overflow-hidden">
        {/* left category sidebar */}
        <nav
          aria-label="设置分类"
          className="flex w-44 shrink-0 flex-col gap-px overflow-y-auto border-r border-line bg-surface p-2"
        >
          <GlideMenu highlightClassName="inset-x-0 rounded-[7px] bg-hover">
            {SECTIONS.map((section) => {
              const on = section === active;
              return (
                <button
                  aria-current={on ? "page" : undefined}
                  className={cn(
                    "relative z-10 flex w-full items-center rounded-[7px] px-2.5 py-1.5 text-left text-[12.5px] transition-colors duration-100",
                    on ? "font-medium text-ink" : "text-ink-2 hover:text-ink",
                  )}
                  data-menu-row
                  key={section}
                  onClick={() => onSelectSection(section)}
                  type="button"
                >
                  {on ? (
                    <span
                      aria-hidden="true"
                      className="absolute bottom-1.5 left-0.5 top-1.5 w-[2px] rounded-full bg-accent"
                    />
                  ) : null}
                  {section}
                </button>
              );
            })}
          </GlideMenu>
        </nav>

        {/* content */}
        <div className="min-h-0 flex-1 overflow-y-auto px-6 py-5">
          {active === "外观" ? (
            <ConfigScreen />
          ) : (
            <UnfedSection name={active} />
          )}
        </div>
      </div>
    </div>
  );
}
