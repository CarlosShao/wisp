/* ============================================================================
   Config editor section (ticket 39 AC#1 / AC#4, readface consume wave).
   ----------------------------------------------------------------------------
   WHAT THIS IS. The config editor as a pure VIEW over config.toml: the FULL
   D36 section tree the config.get leg answers (`sections[]{key,value,tier}`,
   internal/config/sections_456.go -> internal/panel/config_handlers.go:805)
   drawn as ONE TABLE PER SECTION, and the line-numbered load errors
   (`errors[]{line,message}`, internal/config/loaderrors_456.go -> :806)
   listed under their lines. The GUI owns nothing: every value on screen is
   the struct's own spelling handed through the payload, every tier badge is
   the registry's own word VERBATIM (hot / reload / restart / locked) - a
   second vocabulary here would be exactly the second source of truth D36
   rule 4 bans.

   THE TIER BADGE READS THE WIRE, NEVER A MAP: the badge word IS row.tier.
   The registry's words are hot/reload/restart and the locked marker; a row
   whose tier came back "" (an unregistered path - tierForPath answers "" for
   a path neither the per-key nor the section table carries) shows 未登记档位
   as a named absence, never a tier word it was not handed. Tones follow the
   words: hot = applies now, reload = subsystem reload note, restart = the
   restart banner, locked = the danger tint ticket 39's lock affordance
   spells (risk/fs/net/plugins loosening queues an L2 card, planLocked).

   THE ERRORS LIST IS HONEST ABOUT LINES: line is the 1-based line the strict
   decoder pointed at; 0 = the problem names no line (a semantic validation
   failure or a read failure), rendered as （无行号） - never folded into
   第 1 行, because the file did not point there. Wire order is kept: the
   pipeline answers first failure stage first (loaderrors_456.go's stage
   order), so the list shows what the decoder hit first, first.

   THE RELOAD EVENT (AC#1's "after save, applied-tier feedback surfaced from
   the reload events"). An ACCEPTED config.set carries one
   `wisp-config-reload:` event line on its reply (config_handlers.go:579); a
   refused write carries none. The write helpers (lib/config-actions.ts)
   detect the line once per reply and fan it out; this section subscribes and
   RE-READS the whole view on the event instead of patching any local copy
   (PLAN.md:1044). A refusal notifies nobody - nothing changed, so nothing
   reloads.

   THE NAMED ABSENCES (omitempty = absent = the leg did not answer):
     sections absent  "the leg did not answer the tree" - never rendered as
                      an empty file (the schema always has sections; an empty
                      tree here would be a lie about the file).
     errors absent    "loads clean, or the leg did not check" - the page says
                      it cannot tell the two apart from THIS wire; it does
                      not invent an error or a clean bill.
   readable=false    the leg answered and said the config is unreadable - a
                      state of its own, never "nothing configured".

   MOUNT POINT. This section hangs at the foot of the settings page's 外观
   category (components/config-screen.tsx) for this wave; it is
   self-reading (one config.get on mount, like lib/cost-view's hook shape)
   so no snapshot prop plumbing is needed. Re-slotting it into its own rail
   category is a page-layout move, not a wire move.
   ============================================================================ */

import { useEffect } from "react";
import { Button } from "@/components/ai-native/button";
import { Chip } from "@/components/ai-native/chip";
import { SettingsCard } from "@/components/settings-controls";
import { onConfigReloadEvent } from "@/lib/config-actions";
import {
  useConfigView,
  type ConfigLoadErrorView,
  type ConfigSectionRowView,
} from "@/lib/config-view";
import { cn } from "@/lib/cn";

/** The badge tone for one tier word. Unknown words (including "") take the
    neutral tone and are shown verbatim - a tone for a word nobody registered
    would be a judgement the registry never made. */
function tierBadgeTone(tier: string): "ok" | "warn" | "neutral" | "danger" {
  if (tier === "hot") return "ok";
  if (tier === "reload") return "warn";
  if (tier === "locked") return "danger";
  return "neutral";
}

/** One tier badge: the word IS the wire's word. hot/reload/restart/locked
    逐字; "" shows the named absence instead of a tier word. */
function TierBadge({ tier }: { tier: string }) {
  return (
    <span
      className={cn(
        "inline-flex h-[18px] shrink-0 items-center rounded-full px-1.5 font-mono text-[10px] font-medium leading-none",
        tierBadgeTone(tier) === "ok" && "bg-green-tint text-green",
        tierBadgeTone(tier) === "warn" && "bg-orange-tint text-orange",
        tierBadgeTone(tier) === "danger" && "bg-red-tint text-red",
        tierBadgeTone(tier) === "neutral" && "bg-field text-ink-2",
      )}
      title={
        tier === ""
          ? "这一路径没有登记生效档位（TierRegistry 无此键、也无其所在节的整节行）"
          : `生效档位：${tier}（读自 config.get 的 sections[].tier，逐字）`
      }
    >
      {tier === "" ? "未登记档位" : tier}
    </span>
  );
}

/** Group the section rows by their top-level segment, keeping the wire's own
    order of first appearance (the Go walk runs in schema field order; the
    page re-orders nothing). One group = one table. */
function groupBySection(sections: readonly ConfigSectionRowView[]): { section: string; rows: ConfigSectionRowView[] }[] {
  const groups: { section: string; rows: ConfigSectionRowView[] }[] = [];
  const index = new Map<string, number>();
  for (const row of sections) {
    const dot = row.key.indexOf(".");
    const section = dot > 0 ? row.key.slice(0, dot) : row.key;
    const at = index.get(section);
    if (at === undefined) {
      index.set(section, groups.length);
      groups.push({ section, rows: [row] });
    } else {
      groups[at]!.rows.push(row);
    }
  }
  return groups;
}

/** The per-section tables. Fed rows come back out verbatim - the behavioural
    tooth the render instrument feeds (scripts/render-config-screen.tsx): N
    sections in, N tables out, one row per key, tier badge verbatim; a
    section the next read no longer carries leaves the screen entirely, with
    no stale value left behind (the tables render from THIS read, never from
    a remembered one). */
export function ConfigSectionTables({ sections }: { sections: readonly ConfigSectionRowView[] }) {
  const groups = groupBySection(sections);
  return (
    <div className="flex flex-col gap-3">
      <p className="text-[11.5px] leading-[1.6] text-ink-3">
        每节一张表：key 与 value 是 config.get 的 sections[] 逐字（复合值按 Go 侧压成一行 JSON），
        徽章是 TierRegistry 自己的档位词（hot/reload/restart/locked），页面不造第二个词表。
      </p>
      {groups.map((group) => (
        <div aria-label={`配置节 ${group.section}`} className="rounded-card border border-line bg-surface px-3.5 py-3" key={group.section}>
          <div className="mb-1.5 flex items-center gap-2">
            <Chip>{group.section}</Chip>
            <span className="text-[10.5px] text-ink-3">{group.rows.length} 行</span>
          </div>
          <table className="w-full border-collapse text-left">
            <thead>
              <tr className="border-b border-line">
                <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">键</th>
                <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">值</th>
                <th className="py-1.5 text-[11px] font-medium text-ink-3">生效档位</th>
              </tr>
            </thead>
            <tbody>
              {group.rows.map((row) => {
                const dot = row.key.indexOf(".");
                const rest = dot > 0 ? row.key.slice(dot + 1) : row.key;
                return (
                  <tr className="border-b border-line last:border-b-0" key={row.key}>
                    <td className="py-1.5 pr-3 align-top font-mono text-[12px] text-ink">{rest}</td>
                    <td className="py-1.5 pr-3 align-top font-mono text-[12px] break-all text-ink-2">{row.value}</td>
                    <td className="py-1.5 align-top">
                      <TierBadge tier={row.tier} />
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      ))}
    </div>
  );
}

/** The load-error list: one row per problem, under its line. Line 0 renders
    （无行号）- the wire's own "names no line" marker - never 第 0 行 and
    never 第 1 行. The message rides verbatim (it is the decoder's own words,
    loaderrors_456.go builds it from the decoder's error shapes). */
export function ConfigErrorList({ errors }: { errors: readonly ConfigLoadErrorView[] }) {
  return (
    <div className="flex flex-col gap-1.5">
      <p className="text-[11.5px] leading-[1.6] text-orange" role="note">
        读口在这份 config.toml 里查到了 {errors.length} 个问题（{"errors[]{line,message}"} 逐字；
        {"line=0"} 表示这个问题不指向某一行，照实说「无行号」）：
      </p>
      <ol className="m-0 flex list-none flex-col gap-1 p-0">
        {errors.map((problem, i) => (
          <li className="flex items-baseline gap-2" key={`${problem.line}-${i}`}>
            <span className="shrink-0 font-mono text-[11.5px] tabular-nums text-orange">
              {problem.line === 0 ? "（无行号）" : `第 ${problem.line} 行`}
            </span>
            <span className="min-w-0 font-mono text-[11.5px] break-all text-ink">{problem.message}</span>
          </li>
        ))}
      </ol>
    </div>
  );
}

/** The editor section itself: one config.get read, the reload-event
    subscription, and the four honest states around the two faces. The
    subscription re-reads on every accepted write's event - the tables are
    always THIS read's, never a patch. */
export function ConfigEditorSection({ heading = "配置文件" }: { heading?: string }) {
  const { state, reload } = useConfigView();

  // Ticket 39 AC#1: an accepted write fires the wisp-config-reload: event
  // (the write helpers detect the line once); this page answers it with a
  // full re-read. A refused write notifies nobody - zero re-reads, zero
  // flicker - which is the wire's own semantics (nothing changed).
  useEffect(() => onConfigReloadEvent(() => reload()), [reload]);

  return (
    <SettingsCard
      title={heading}
      description="config.get 的全节树读面（sections[]{key,value,tier}）与行号错误面（errors[]{line,message}）：页面只画读口给的，不造第二个配置树（D36 规则 4）。写入被接受时宿主带 wisp-config-reload 事件，这里整页重读。"
    >
      {state.phase === "loading" ? <p className="text-[12px] text-ink-3">正在读取配置文件...</p> : null}
      {state.phase === "no-host" ? (
        <p className="text-[12px] leading-[1.7] text-ink-2">
          暂未连接到本机服务，读不到配置文件全树（config.get 门在本机侧，浏览器页没有宿主）。
        </p>
      ) : null}
      {state.phase === "unreadable" || state.phase === "payload-error" ? (
        <div className="flex flex-col gap-2">
          <p className="text-[12px] leading-[1.7] text-ink-2">
            {state.phase === "unreadable" ? "读口拒绝了这个请求，原话逐字如下：" : "配置载荷没排出来（payloadError），原话如下："}
          </p>
          <p className="whitespace-pre-wrap font-mono text-[11px] leading-[1.6] text-ink-3">
            {state.phase === "unreadable" ? state.raw : state.message}
          </p>
          <Button className="self-start" onClick={reload} size="xs" variant="secondary">
            重试
          </Button>
        </div>
      ) : null}
      {state.phase === "ready" ? (
        <div className="flex flex-col gap-2">
          {!state.view.readable ? (
            <p className="text-[12px] leading-[1.7] text-orange" role="note">
              读口答了：这份配置不可读（readable=false）——这是读到的状态，不是「没有配置」。
            </p>
          ) : null}
          {state.view.sections === undefined ? (
            <p className="text-[12px] leading-[1.7] text-ink-2">
              这次读口没带全节树（sections 键不在载荷上＝本腿没有答这一面）——不把「没答」画成「文件是空的」。
            </p>
          ) : (
            <ConfigSectionTables sections={state.view.sections} />
          )}
          {state.view.errors === undefined ? (
            <p className="text-[11px] leading-[1.6] text-ink-3">
              错误面不在载荷上：文件装载干净，或本腿没有检查——这一侧从线上分不出这两种，照实说，不编一个错误出来。
            </p>
          ) : state.view.errors.length > 0 ? (
            <ConfigErrorList errors={state.view.errors} />
          ) : (
            <p className="text-[11px] leading-[1.6] text-ink-3">读口带了错误面且为空：这次检查没有查出问题。</p>
          )}
        </div>
      ) : null}
    </SettingsCard>
  );
}
