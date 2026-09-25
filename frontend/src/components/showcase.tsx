/* ============================================================================
   Showcase - the harness walkthrough page (main.tsx renders it under
   ?harness=1; owner 2026-09-25: "起码要有主流 harness 的样子").
   ----------------------------------------------------------------------------
   The shape is the component-library showcase: a sticky header (ball dot +
   title, light/dark toggle on the right), a sticky section index on the left,
   and numbered sections on the right - each with one line of intent and one
   to three framed demo cards, every card fed by src/fixtures/harness.ts.
   The footer carries the honesty line: the showcase is demo data, the product
   path reads the snapshot Go pushes.

   What is real here: the components. L2ApprovalCard, ApprovalScreen,
   ConfigScreen, the palette / tasks / ball screens and the vendored atoms are
   the same files the panel mounts; they eat fixtures only because this page
   runs without a host. What is NOT real: every value on the page - hence the
   demo-data chip on each card and the banner main.tsx adds.

   State policy for this file only: it holds the walkthrough's own view
   preferences (theme toggle, palette query, the two demo control values).
   None of it is persisted, none of it reaches the bridge, and none of it
   belongs to any component a Go-fed screen would mount - those stay
   stateless and prop-driven.
   ============================================================================ */

import { useState, type ReactNode } from "react";
import { Moon, Sun } from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { CodeBlock } from "@/components/ai-native/code-block";
import { ContextCards } from "@/components/ai-native/context-cards";
import { DiffTable } from "@/components/ai-native/diff-table";
import { EntityChip } from "@/components/ai-native/entity-chip";
import { LoadingState } from "@/components/ai-native/loading-state";
import { SegmentedControl } from "@/components/ai-native/segmented-control";
import { StreamingText } from "@/components/ai-native/streaming-text";
import { Switch } from "@/components/ai-native/switch";
import { TextRow } from "@/components/ai-native/text-row";
import { Thinking } from "@/components/ai-native/thinking";
import { ToolChips } from "@/components/ai-native/tool-chips";
import { ValuePill } from "@/components/ai-native/value-pill";
import { ApprovalScreen } from "@/components/approval-screen";
import { BallScreen } from "@/components/ball-screen";
import { ConfigScreen } from "@/components/config-screen";
import { L2ApprovalCard } from "@/components/l2-approval-card";
import { PaletteScreen } from "@/components/palette-screen";
import { TasksScreen } from "@/components/tasks-screen";
import {
  SHOWCASE_APPROVALS,
  SHOWCASE_APPROVAL_SNAPSHOT,
  SHOWCASE_AUDIT_ROWS,
  SHOWCASE_BALL_RINGS,
  SHOWCASE_BALL_STATES,
  SHOWCASE_CODE,
  SHOWCASE_CONTEXT,
  SHOWCASE_COST,
  SHOWCASE_DIFF,
  SHOWCASE_PALETTE_GROUPS,
  SHOWCASE_PERMISSION_TIERS,
  SHOWCASE_PRIVACY_MEMORY,
  SHOWCASE_PRIVACY_PROFILE,
  SHOWCASE_RETENTION_DEFAULT,
  SHOWCASE_RETENTION_OPTIONS,
  SHOWCASE_STREAM,
  SHOWCASE_TASKS,
  SHOWCASE_THINKING,
  SHOWCASE_TOOLS,
} from "@/fixtures/harness";

interface SectionMeta {
  id: string;
  num: string;
  name: string;
  desc: string;
}

const SECTIONS: readonly SectionMeta[] = [
  { id: "sec-chat", num: "01", name: "对话", desc: "流式回复、思考块与工具条：beautiful-ui 的 agent 语言，字段一到就画。" },
  { id: "sec-approval", num: "02", name: "审批", desc: "L2 强确认卡与审批中心：面板只递请求，批准权在原生侧。" },
  { id: "sec-palette", num: "03", name: "命令", desc: "命令面板：大输入框、分组过滤与键盘词汇。" },
  { id: "sec-tasks", num: "04", name: "任务", desc: "任务列表：六态色标、运行弧与 80ms 入场节奏。" },
  { id: "sec-config", num: "05", name: "设置", desc: "检查器风格：mono 键名 + 控件 + 生效徽标，唯一活控制是面板不透明度。" },
  { id: "sec-ball", num: "06", name: "球状态", desc: "悬浮球状态机的态与视觉（PLAN.md D43 冻结的转移表不在此复述）。" },
  { id: "sec-security", num: "07", name: "安全", desc: "权限档三卡与授权记录：原生侧做决定，面板转述。" },
  { id: "sec-privacy", num: "08", name: "隐私", desc: "L1 画像 / L2 记忆两卡，外加保留期与本地保留开关。" },
  { id: "sec-cost", num: "09", name: "成本", desc: "成本三数与一周用量；数字口径是 C23 的 mono tabular 纯文字。" },
];

/** The framed demo card every section drops its components into. The chip in
    the corner is the honesty mark: everything inside eats fixtures. */
function DemoCard({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="relative overflow-hidden rounded-card border border-line bg-surface p-4 shadow-card">
      <span className="absolute right-2.5 top-2.5 rounded-chip bg-field px-1.5 text-[10px] leading-[18px] text-ink-3">
        demo 数据
      </span>
      <p className="mb-3 pr-16 text-[13px] font-medium text-ink">{title}</p>
      {children}
    </div>
  );
}

function Section({ meta, children }: { meta: SectionMeta; children: ReactNode }) {
  return (
    <section className="scroll-mt-20" id={meta.id}>
      <h2 className="text-[14px] font-semibold text-ink">
        <span className="mr-2 font-mono font-normal text-ink-3">{meta.num}</span>
        {meta.name}
      </h2>
      <p className="mt-1 text-[12.5px] leading-[1.6] text-ink-2">{meta.desc}</p>
      <div className="mt-3 flex flex-col gap-4">{children}</div>
    </section>
  );
}

/** Authorization-record actor mark: the tool name on a tier-coloured disc.
    The disc colour is a token reference, never a literal. */
const TIER_COLORS: Record<string, string> = {
  L0: "var(--ink-3)",
  L1: "var(--accent)",
  L2: "var(--red)",
};

const TIER_OPTIONS = ["L0 只读", "L1 可逆写", "L2 不可逆"] as const;

export function Showcase() {
  // The walkthrough's own view preferences. All of them die with the page;
  // nothing here is persisted and nothing reaches the bridge.
  const [dark, setDark] = useState(false);
  const [query, setQuery] = useState("");
  const [tier, setTier] = useState<(typeof TIER_OPTIONS)[number]>("L1 可逆写");
  const [retention, setRetention] = useState<(typeof SHOWCASE_RETENTION_OPTIONS)[number]>(
    SHOWCASE_RETENTION_DEFAULT,
  );
  const [keepAudio, setKeepAudio] = useState(false);

  function toggleTheme() {
    const next = !dark;
    setDark(next);
    document.documentElement.dataset.theme = next ? "dark" : "light";
  }

  const maxPct = Math.max(...SHOWCASE_COST.week.map((d) => d.pct));

  return (
    <div className="min-h-screen bg-page text-ink">
      {/* header */}
      <header className="sticky top-0 z-20 border-b border-line bg-page/95 backdrop-blur-sm">
        <div className="mx-auto flex h-14 w-full max-w-[960px] items-center justify-between gap-3 px-8">
          <div className="flex min-w-0 items-center gap-2.5">
            <span aria-hidden="true" className="size-3.5 shrink-0 rounded-full bg-accent" />
            <span className="truncate text-[13px] font-medium text-ink">一缕 · 面板组件总览</span>
            <span className="hidden truncate text-[11.5px] text-ink-3 sm:inline">
              九个章节，组件吃假数据
            </span>
          </div>
          <Button onClick={toggleTheme} size="xs" variant="secondary">
            {dark ? <Sun aria-hidden="true" /> : <Moon aria-hidden="true" />}
            <span>{dark ? "浅色" : "暗色"}</span>
          </Button>
        </div>
      </header>

      {/* body: sticky index left, numbered sections right */}
      <div className="mx-auto grid w-full max-w-[960px] grid-cols-1 gap-10 px-8 py-10 md:grid-cols-[150px_minmax(0,1fr)]">
        <nav aria-label="showcase 章节" className="sticky top-24 hidden self-start md:block">
          <ul className="m-0 flex list-none flex-col gap-2 p-0">
            {SECTIONS.map((meta) => (
              <li key={meta.id}>
                <a
                  className="text-[12px] text-ink-2 transition-colors duration-100 hover:text-ink"
                  href={`#${meta.id}`}
                >
                  <span className="mr-1.5 font-mono text-[11px] text-ink-3">{meta.num}</span>
                  {meta.name}
                </a>
              </li>
            ))}
          </ul>
        </nav>

        <div className="flex min-w-0 flex-col gap-12">
          {/* 01 对话 */}
          <Section meta={SECTIONS[0]}>
            <DemoCard title="流式回复（进行中 / 已完成）">
              <div className="flex max-w-[620px] flex-col gap-5">
                <StreamingText
                  streaming
                  sources={SHOWCASE_STREAM.sources}
                  text={SHOWCASE_STREAM.streaming}
                />
                <StreamingText followUps={SHOWCASE_STREAM.followUps} text={SHOWCASE_STREAM.done} streaming={false} />
              </div>
            </DemoCard>
            <DemoCard title="思考块（可展开）">
              <div className="max-w-[620px]">
                <Thinking
                  reasoning={SHOWCASE_THINKING.reasoning}
                  seconds={SHOWCASE_THINKING.seconds}
                  sources={SHOWCASE_THINKING.sources}
                  steps={SHOWCASE_THINKING.steps}
                />
              </div>
            </DemoCard>
            <DemoCard title="工具条 · 四态">
              <div className="max-w-[620px]">
                <ToolChips tools={SHOWCASE_TOOLS} />
              </div>
            </DemoCard>
            <DemoCard title="上下文来源卡 / 差异行 / 代码块 / 全局加载">
              <div className="flex max-w-[620px] flex-col gap-4">
                <ContextCards items={SHOWCASE_CONTEXT} />
                <DiffTable rows={SHOWCASE_DIFF} />
                <CodeBlock lines={SHOWCASE_CODE} />
                <LoadingState label="正在唤醒" />
              </div>
            </DemoCard>
          </Section>

          {/* 02 审批 */}
          <Section meta={SECTIONS[1]}>
            <DemoCard title="L2 强确认卡 · 三个档位各一张">
              <div className="flex max-w-[560px] flex-col gap-3">
                {SHOWCASE_APPROVALS.map((view) => (
                  <L2ApprovalCard key={view.correlationId} view={view} />
                ))}
              </div>
              <p className="mt-3 text-[11px] leading-[1.6] text-ink-3">
                拒绝按钮把请求递给原生宿主；在无宿主的浏览器里会按规矩把「未送达」抛进控制台，而不是假装成功。倒计时与推荐卡没有字段支撑，不画。
              </p>
            </DemoCard>
            <DemoCard title="审批中心（审批屏）">
              <ApprovalScreen snapshot={SHOWCASE_APPROVAL_SNAPSHOT} />
            </DemoCard>
          </Section>

          {/* 03 命令 */}
          <Section meta={SECTIONS[2]}>
            <DemoCard title="命令面板">
              <PaletteScreen
                groups={SHOWCASE_PALETTE_GROUPS}
                onQueryChange={setQuery}
                query={query}
              />
            </DemoCard>
          </Section>

          {/* 04 任务 */}
          <Section meta={SECTIONS[3]}>
            <DemoCard title="任务列表 · 六态">
              <TasksScreen tasks={SHOWCASE_TASKS} />
            </DemoCard>
          </Section>

          {/* 05 设置 */}
          <Section meta={SECTIONS[4]}>
            <DemoCard title="设置屏（真实组件，唯一活控制）">
              <ConfigScreen />
            </DemoCard>
          </Section>

          {/* 06 球状态 */}
          <Section meta={SECTIONS[5]}>
            <DemoCard title="悬浮球状态机 · 态与视觉">
              <BallScreen rings={SHOWCASE_BALL_RINGS} states={SHOWCASE_BALL_STATES} />
            </DemoCard>
          </Section>

          {/* 07 安全 */}
          <Section meta={SECTIONS[6]}>
            <DemoCard title="权限档（C2 的三档风险级）">
              <SegmentedControl
                onChange={setTier}
                options={TIER_OPTIONS}
                value={tier}
              />
              <div className="mt-3 grid grid-cols-1 gap-2 sm:grid-cols-3">
                {SHOWCASE_PERMISSION_TIERS.map((t) => {
                  const active = tier.startsWith(t.tier);
                  return (
                    <div
                      className={
                        active
                          ? "rounded-control border border-accent/40 bg-inset p-2.5"
                          : "rounded-control border border-line bg-inset p-2.5"
                      }
                      key={t.tier}
                    >
                      <ValuePill tone={t.tier === "L2" ? "red" : t.tier === "L1" ? "accent" : "neutral"}>
                        {t.tier}
                      </ValuePill>
                      <p className="mt-1.5 text-[12.5px] font-medium text-ink">{t.name}</p>
                      <p className="mt-1 text-[11.5px] leading-[1.6] text-ink-2">{t.scope}</p>
                      <p className="mt-1 text-[11px] leading-[1.6] text-ink-3">{t.confirm}</p>
                    </div>
                  );
                })}
              </div>
            </DemoCard>
            <DemoCard title="授权记录">
              <div className="overflow-hidden rounded-control border border-line">
                <div className="grid grid-cols-[52px_minmax(0,1fr)_36px_minmax(0,1.4fr)] gap-2 border-b border-line bg-inset px-2.5 py-1.5 text-[11px] font-medium text-ink-3">
                  <span>时间</span>
                  <span>工具</span>
                  <span>档</span>
                  <span>结果</span>
                </div>
                {SHOWCASE_AUDIT_ROWS.map((row) => (
                  <div
                    className="grid grid-cols-[52px_minmax(0,1fr)_36px_minmax(0,1.4fr)] items-center gap-2 border-b border-line px-2.5 py-2 text-[12px] transition-colors duration-100 last:border-b-0 hover:bg-hover"
                    key={`${row.time}-${row.tool}`}
                  >
                    <span className="font-mono text-[11px] tabular-nums text-ink-3">{row.time}</span>
                    <EntityChip color={TIER_COLORS[row.tier] ?? TIER_COLORS.L0} name={row.tool} />
                    <span className="font-mono text-[11px] text-ink-2">{row.tier}</span>
                    <span className="text-ink-2">{row.result}</span>
                  </div>
                ))}
              </div>
            </DemoCard>
          </Section>

          {/* 08 隐私 */}
          <Section meta={SECTIONS[7]}>
            <DemoCard title="两卡：L1 画像 / L2 记忆">
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <div className="rounded-control border border-line p-3">
                  <ValuePill tone="accent">L1 画像</ValuePill>
                  <div className="mt-1 flex flex-col divide-y divide-line">
                    {SHOWCASE_PRIVACY_PROFILE.map((row) => (
                      <TextRow key={row.label} label={row.label} meta={row.meta} value={row.value} />
                    ))}
                  </div>
                  <p className="mt-1 text-[11px] leading-[1.6] text-ink-3">
                    只在本机，跟着会话走，不出进程。
                  </p>
                </div>
                <div className="rounded-control border border-line p-3">
                  <ValuePill tone="orange">L2 记忆</ValuePill>
                  <div className="mt-1 flex flex-col divide-y divide-line">
                    {SHOWCASE_PRIVACY_MEMORY.map((row) => (
                      <TextRow key={row.label} label={row.label} meta={row.meta} value={row.value} />
                    ))}
                  </div>
                  <div className="mt-3 flex flex-col gap-2.5">
                    <div className="flex items-center justify-between gap-3">
                      <span className="text-[12.5px] text-ink-2">保留期</span>
                      <SegmentedControl
                        onChange={setRetention}
                        options={SHOWCASE_RETENTION_OPTIONS}
                        value={retention}
                      />
                    </div>
                    <div className="flex items-center justify-between gap-3">
                      <span className="text-[12.5px] text-ink-2">保留原始音频（默认关闭，D16）</span>
                      <Switch
                        checked={keepAudio}
                        label="保留原始音频"
                        onChange={setKeepAudio}
                      />
                    </div>
                  </div>
                </div>
              </div>
            </DemoCard>
          </Section>

          {/* 09 成本 */}
          <Section meta={SECTIONS[8]}>
            <DemoCard title="成本三数">
              <div className="grid grid-cols-1 gap-2 sm:grid-cols-3">
                <div className="rounded-control border border-line p-3">
                  <ValuePill>今日 tokens</ValuePill>
                  <p className="mt-1.5 font-mono text-[17px] font-semibold tabular-nums text-ink">
                    {SHOWCASE_COST.tokens}
                  </p>
                </div>
                <div className="rounded-control border border-line p-3">
                  <ValuePill>今日花费</ValuePill>
                  <p className="mt-1.5 font-mono text-[17px] font-semibold tabular-nums text-ink">
                    {SHOWCASE_COST.money}
                  </p>
                </div>
                <div className="rounded-control border border-line p-3">
                  <ValuePill tone="red">月预算占用</ValuePill>
                  <p className="mt-1.5 font-mono text-[17px] font-semibold tabular-nums text-ink">
                    {SHOWCASE_COST.budget}
                  </p>
                </div>
              </div>
            </DemoCard>
            <DemoCard title="一周用量（纯 CSS 柱状，无图表库）">
              <div className="flex h-28 items-end gap-1.5">
                {SHOWCASE_COST.week.map((d) => (
                  <div className="flex h-full min-w-0 flex-1 flex-col items-center justify-end gap-1" key={d.day}>
                    <div
                      className={d.pct === maxPct ? "w-full rounded-t-[4px] bg-accent" : "w-full rounded-t-[4px] bg-accent-tint"}
                      style={{ height: `${d.pct}%` }}
                    />
                    <span className="text-[10.5px] text-ink-3">{d.day}</span>
                  </div>
                ))}
              </div>
            </DemoCard>
          </Section>
        </div>
      </div>

      <footer className="border-t border-line py-6 text-center text-[11px] text-ink-3">
        HARNESS 假数据 · 生产路径读 Go 推送的快照
      </footer>
    </div>
  );
}
