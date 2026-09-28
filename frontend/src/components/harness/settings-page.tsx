/* ============================================================================
   Settings page - a full page with a category rail, not a modal.
   ----------------------------------------------------------------------------
   Owner 2026-09-28, second pass. The first pass had the right eleven
   categories and one flaw worth recording: every unfed category rendered the
   SAME sentence, so the rail read as eleven identical rows - the user could
   not tell which one was which, only that all of them were empty. Each
   category now carries its own one-liner saying what it is for and what it is
   waiting on, so the rail itself becomes the honest status board.

   Where the 9 legacy rail rows went, since the sidebar used to show them:

     命令  -> not a view. It is the Ctrl+K palette, opened from the header.
     设置  -> this page, reached once from the sidebar bottom.
     安全 / 隐私 / 成本 -> categories on the rail below, which is where a
              setting belongs. Putting them in the sidebar made a settings
              list look like navigation.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { ArrowLeft } from "lucide-react";
import { ConfigScreen } from "@/components/config-screen";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { cn } from "@/lib/cn";

/** One category: what it is, and what it is still waiting for. The second
    sentence is the reason this table exists rather than a bare label list. */
interface Category {
  name: string;
  /** One line: what a person would expect to change here. */
  purpose: string;
  /** What the panel is waiting on. "" means the category has real controls. */
  waiting: string;
}

const CATEGORIES: readonly Category[] = [
  {
    name: "通用",
    purpose: "启动行为、默认工作区、面板与球的分工。",
    waiting: "还没收到配置读口。config.toml 的全量 section 树在 D36 里定，但 C17 还没有把任何一节推给面板。",
  },
  {
    name: "外观",
    purpose: "主题、面板底色深浅、动画档位、字号。",
    waiting: "", // ConfigScreen renders here; two of its rows are live today.
  },
  {
    name: "模型",
    purpose: "默认模型、温度、上下文长度、每家的密钥。",
    waiting: "还没收到模型清单。C5 抽象了 provider，但面板没有一个字段能列出「有哪些模型可选」。",
  },
  {
    name: "语音",
    purpose: "听写与朗读、唤醒词、打断策略。",
    waiting: "还没收到语音链路的状态。D47 定了半双工干活 / 全双工陪聊两条路，面板看不到它们走到哪一步。",
  },
  {
    name: "记忆",
    purpose: "长期记忆的开关、范围与清空。",
    waiting: "还没收到记忆数据。D20 定义了记忆系统，但没有任何一条推送告诉面板「记了什么」。",
  },
  {
    name: "技能",
    purpose: "装了哪些技能、每个技能的权限范围。",
    waiting: "还没收到技能清单。技能要能被列出来，前提是它们从插件运行时注册上来。",
  },
  {
    name: "MCP",
    purpose: "MCP 服务端、每个服务提供哪些工具。",
    waiting: "还没收到 MCP 状态。D13 定了 MCP 的取舍，但没有字段把已连的服务端报给面板。",
  },
  {
    name: "安全",
    purpose: "权限档位、门控阈值、批量授权的边界。",
    waiting: "还没收到策略快照。这一节最敏感：档位和阈值属于原生侧，面板在这里显示的任何数字都必须是读来的，不能是猜的。",
  },
  {
    name: "隐私",
    purpose: "哪些数据留在本机、录音留存多久、遥测开关。",
    waiting: "还没收到隐私配置。本节的所有项都要落到 config.toml 的读口上，没有读口就一个字也不能显示。",
  },
  {
    name: "成本",
    purpose: "每轮花费、本月累计、预算与超支提醒。",
    waiting: "还没收到用量数据。成本口径由 C23 定，但面板快照里没有 token 数，所以这一格连一个空壳数字都不画。",
  },
  {
    name: "引导",
    purpose: "首次使用的几步，以及随时能重看的说明。",
    waiting: "还没收到引导进度。首启那套画面在首启屏上，不在这里重画一遍。",
  },
];

export function SettingsPage({
  activeSection,
  onSelectSection,
  onClose,
}: {
  activeSection: string;
  onSelectSection: (section: string) => void;
  onClose: () => void;
}) {
  // An unknown name is not guessed at: it lands on the first category, the
  // same way an unknown view does. A rail that lights nothing is worse.
  const active =
    CATEGORIES.find((c) => c.name === activeSection) ?? CATEGORIES[0]!;
  const cat = active;

  return (
    <div className="flex h-full min-h-0 flex-1 flex-col bg-page">
      <header className="flex h-10 shrink-0 items-center gap-2 border-b border-line bg-surface px-3">
        <button
          type="button"
          onClick={onClose}
          className="flex shrink-0 items-center gap-1 rounded-control px-1.5 py-1 text-[12px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
        >
          <ArrowLeft aria-hidden="true" size={13} strokeWidth={2} />
          返回
        </button>
        <span aria-hidden="true" className="h-3.5 w-px shrink-0 bg-line" />
        <span className="min-w-0 flex-1 truncate text-[12.5px] font-medium text-ink">
          {cat.name}
        </span>
      </header>

      <div className="flex min-h-0 flex-1 overflow-hidden">
        <nav
          aria-label="设置分类"
          className="flex w-48 shrink-0 flex-col overflow-y-auto border-r border-line bg-surface p-2"
        >
          <GlideMenu className="flex flex-col gap-px" highlightClassName="inset-x-0 rounded-[7px] bg-hover">
            {CATEGORIES.map((c) => {
              const on = c.name === cat.name;
              // The status dot is decorative, so it is hidden from assistive
              // tech and the state is folded into the button's own name
              // instead. An aria-label on the DOT leaked into the button's
              // accessible name (「成本还没有数据」), which broke every by-name
              // lookup of the rail - including a smoke test's own click.
              const name = c.waiting ? `${c.name}（还没有数据）` : c.name;
              return (
                <button
                  key={c.name}
                  type="button"
                  aria-current={on ? "page" : undefined}
                  aria-label={name}
                  data-menu-row
                  title={c.purpose}
                  onClick={() => onSelectSection(c.name)}
                  className={cn(
                    "relative z-10 flex w-full items-center rounded-[7px] py-1.5 pr-2 text-left text-[12.5px] transition-colors duration-100",
                    on ? "pl-3.5 font-medium text-ink" : "pl-3.5 text-ink-2",
                  )}
                >
                  {on ? (
                    <span
                      aria-hidden="true"
                      className="absolute bottom-1.5 left-1 top-1.5 w-[2px] rounded-full bg-accent"
                    />
                  ) : null}
                  <span aria-hidden="true" className="min-w-0 flex-1 truncate">
                    {c.name}
                  </span>
                  {/* The unfed ones are marked, so the rail doubles as a
                      status board: you can see at a glance which half of
                      settings is real. */}
                  {c.waiting ? (
                    <span
                      aria-hidden="true"
                      className="size-1.5 shrink-0 rounded-full bg-line-strong"
                    />
                  ) : null}
                </button>
              );
            })}
          </GlideMenu>
        </nav>

        <div className="min-h-0 flex-1 overflow-y-auto px-6 py-5">
          {cat.waiting === "" ? (
            <ConfigScreen
              heading={cat.name}
              lede={cat.purpose}
            />
          ) : (
            <section aria-label={cat.name} className="flex min-w-0 flex-col gap-1">
              <h2 className="text-[15px] font-semibold text-ink">{cat.name}</h2>
              <p className="text-[12.5px] text-ink-2">{cat.purpose}</p>
              <div className="mt-4 rounded-card border border-line bg-surface px-3.5 py-3">
                <p className="text-[12px] leading-[1.7] text-ink-3">{cat.waiting}</p>
              </div>
              <p className="mt-2.5 text-[11px] leading-[1.6] text-ink-3">
                这一类不画控件。画一个空壳开关，等于让一个改不动的旋钮看起来像能改。
              </p>
            </section>
          )}
        </div>
      </div>
    </div>
  );
}
