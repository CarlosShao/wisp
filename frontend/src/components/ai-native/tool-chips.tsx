/* ============================================================================
   Adapted component.
   ----------------------------------------------------------------------------
   来源 = derived from TurboKach/ai-native-react-components
          components/tool-chips.tsx (MIT, Copyright (c) 2026 Turbo,
          upstream commit 05dab2d2).
   本地改动：
   - props 化: { tools: { name, status, summary?, detail? }[] } - 上游是定时
     分幕的演示（ROWS/step 计数），行的出现顺序与展开态全部由 props 决定；
     tools 为空时不渲染任何东西（宁缺毋造）。
   - 四态: pending = 旋转弧（border-line-strong 底 + border-t-ink-2 的弧）；
     running = 同弧加 shimmer 文案「执行中」；done = lucide Check 绿；
     denied = lucide Ban 红 + 整行红描边（border-red/40，无新颜色字面量）。
   - 数据移除: 上游演示行（think/write/run/read）、DIFFS 文件差分 chip 区、
     "+2 more" 演示尾巴全部删除；折叠头计数改由 tools.length 推导，不造数。
   - 图标映射: fs.read/fs.list 取 FileText，其余 fs.* 取 FolderInput；
     shell.* 取 Terminal；web.* 取 Globe；回退 Wrench（lucide-react 已逐一
     核实存在）。上游 think/write/run/read 内联 SVG 删除（sparkles 形状的
     think 星形按 D23 图标规则不再使用）。
   - 字面量换 token: 上游 chip 的硬编码文字灰（一个 # 号十六进制灰）换成
     text-ink-2；dark: 前缀
     变体删除（token 表随 data-theme 自动翻转，不需要组件层再写一份）。
   - 详情展开行保持蓝本形状：左边框线 + mono 行；detail 是工具输出，统一
     mono。不可展开的行用 div 而不是无路由的假按钮。
   ============================================================================ */

import { useState, type ReactNode } from "react";
import {
  Ban,
  Check,
  ChevronDown,
  FileText,
  FolderInput,
  Globe,
  Terminal,
  Wrench,
  type LucideIcon,
} from "lucide-react";

export type ToolStatus = "pending" | "running" | "done" | "denied";

export interface ToolCall {
  name: string;
  status: ToolStatus;
  summary?: string;
  detail?: string[];
}

/* 图标名映射表：GatedTool 的点分名到 lucide 图标。 */
function toolIcon(name: string): LucideIcon {
  if (name.startsWith("fs.")) {
    /* 读文件画纸面，写入/删除/建目录这类改动性调用画进目录的文件夹。 */
    return name === "fs.read" || name === "fs.list" ? FileText : FolderInput;
  }
  if (name.startsWith("shell.")) return Terminal;
  if (name.startsWith("web.")) return Globe;
  return Wrench;
}

/* pending 与 running 共用旋转弧；running 另有 shimmer 的「执行中」字样。 */
function StatusGlyph({ status }: { status: ToolStatus }) {
  if (status === "done") {
    return <Check aria-hidden="true" className="shrink-0 text-green" size={14} strokeWidth={2.5} />;
  }
  if (status === "denied") {
    return <Ban aria-hidden="true" className="shrink-0 text-red" size={14} strokeWidth={2.5} />;
  }
  return (
    <span
      aria-hidden="true"
      className="size-3 shrink-0 rounded-full border-[1.5px] border-line-strong border-t-ink-2"
      style={{ animation: "spin 700ms linear infinite" }}
    />
  );
}

function RunningLabel() {
  return (
    <span
      className="shrink-0 bg-clip-text text-[11.5px] font-medium text-transparent"
      style={{
        backgroundImage:
          "linear-gradient(90deg, var(--ink-3) 35%, var(--ink) 50%, var(--ink-3) 65%)",
        backgroundSize: "200% 100%",
        animation: "shimmer-text 1.4s linear infinite",
      }}
    >
      执行中
    </span>
  );
}

export function ToolChips({ tools }: { tools: ToolCall[] }) {
  const [open, setOpen] = useState(true);
  const [openRows, setOpenRows] = useState<Set<number>>(new Set());

  if (tools.length === 0) return null;

  const toggleRow = (index: number) =>
    setOpenRows((current) => {
      const next = new Set(current);
      if (next.has(index)) {
        next.delete(index);
      } else {
        next.add(index);
      }
      return next;
    });

  return (
    <div className="w-full pb-1">
      {/* 折叠头：计数由 props 推导 */}
      <button
        aria-expanded={open}
        className="-mx-1.5 flex w-fit items-center gap-1.5 rounded-control px-1.5 py-1
          text-[12.5px] text-ink-2 transition-colors duration-100 hover:bg-hover-2"
        onClick={() => setOpen((current) => !current)}
        type="button"
      >
        <ChevronDown
          aria-hidden="true"
          className="transition-transform duration-200"
          size={12}
          strokeWidth={2.2}
          style={{ transform: open ? "rotate(0deg)" : "rotate(-90deg)" }}
        />
        <span className="tabular-nums">{tools.length} 次工具调用</span>
      </button>

      {/* 行列表：-mx-1 + px-1.5 让行悬停药丸有地方站（照蓝本的裁剪盒写法） */}
      <div
        className="grid transition-[grid-template-rows,opacity] duration-300"
        style={{ gridTemplateRows: open ? "1fr" : "0fr", opacity: open ? 1 : 0 }}
      >
        <div className="-mx-1 overflow-hidden px-1.5 pb-1">
          <div className="mt-1.5 flex flex-col gap-1">
            {tools.map((tool, i) => {
              const Icon = toolIcon(tool.name);
              const detail = tool.detail ?? [];
              const expandable = detail.length > 0;
              const rowOpen = openRows.has(i);
              const denied = tool.status === "denied";
              const rowClass = `-mx-[3px] flex h-7 w-[calc(100%+6px)] min-w-0 items-center gap-2
                rounded-control px-[3px] text-left transition-colors duration-100
                ${denied ? "border border-red/40" : "border border-transparent"}
                ${expandable ? "hover:bg-hover-2" : ""}`;
              const rowBody = (
                <RowBody
                  denied={denied}
                  icon={<Icon aria-hidden="true" size={13} strokeWidth={2} />}
                  name={tool.name}
                  open={rowOpen}
                  status={tool.status}
                  summary={tool.summary}
                />
              );
              return (
                <div
                  key={`${i}-${tool.name}`}
                  style={{ animation: "fade-up 300ms var(--ease-out-strong) both" }}
                >
                  {expandable ? (
                    <button
                      aria-expanded={rowOpen}
                      className={`group/row ${rowClass}`}
                      onClick={() => toggleRow(i)}
                      type="button"
                    >
                      {rowBody}
                    </button>
                  ) : (
                    <div className={rowClass}>{rowBody}</div>
                  )}

                  {expandable ? (
                    <div
                      className="grid transition-[grid-template-rows,opacity] duration-300"
                      style={{
                        gridTemplateRows: rowOpen ? "1fr" : "0fr",
                        opacity: rowOpen ? 1 : 0,
                        transitionTimingFunction: "var(--ease-out-strong)",
                      }}
                    >
                      <div className="min-h-0 overflow-hidden">
                        <div className="mt-0.5 mb-1 ml-2 flex flex-col gap-0.5 border-l border-line py-0.5 pl-3.5">
                          {detail.map((line, j) => (
                            <span
                              className="truncate font-mono text-[11.5px] leading-[1.6] text-ink-2"
                              key={`${j}-${line}`}
                            >
                              {line}
                            </span>
                          ))}
                        </div>
                      </div>
                    </div>
                  ) : null}
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}

/* 行内容：图标位（悬停换箭头）+ 工具名 + 状态符 + 执行中字样 + 摘要 chip。 */
function RowBody({
  icon,
  name,
  open,
  status,
  summary,
  denied,
}: {
  icon: ReactNode;
  name: string;
  open: boolean;
  status: ToolStatus;
  summary?: string;
  denied: boolean;
}) {
  return (
    <>
      <span className="relative flex size-4 shrink-0 items-center justify-center text-ink-3">
        <span
          className={`flex transition-opacity duration-100 group-hover/row:opacity-0 ${open ? "opacity-0" : ""}`}
        >
          {icon}
        </span>
        <ChevronDown
          aria-hidden="true"
          className={`absolute transition-[opacity,transform] duration-150 group-hover/row:opacity-100 ${
            open ? "rotate-0 opacity-100" : "-rotate-90 opacity-0"
          }`}
          size={12}
          strokeWidth={2.2}
        />
      </span>
      <span className="shrink-0 text-[12.5px] font-medium text-ink">{name}</span>
      <span className="flex size-4 shrink-0 items-center justify-center">
        <StatusGlyph status={status} />
      </span>
      {status === "running" ? <RunningLabel /> : null}
      {summary ? (
        <span
          className={`inline-flex h-5.5 min-w-0 flex-1 items-center truncate rounded-chip px-1.5
            text-[11.5px] shadow-hairline transition-colors duration-100 ${
              denied ? "bg-red-tint text-red" : "bg-hover-2 text-ink-2 hover:bg-line-strong"
            }`}
        >
          {summary}
        </span>
      ) : (
        <span className="min-w-0 flex-1" />
      )}
    </>
  );
}
