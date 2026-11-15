import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { Fragment, useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from "react";
import { Activity, ArrowLeft, ArrowUp, Check, ChevronDown, ChevronLeft, ChevronRight, CircleAlert, Command, Compass, FilePen, FileSearch, FileText, Folder, FolderCog, FolderInput, GitBranch, GitFork, GitPullRequest, Info, LoaderCircle, Mic, Moon, PanelLeft, PanelRight, PanelRightClose, Paperclip, PencilLine, Pin, Plus, RefreshCw, RotateCw, Search, Settings, ShieldAlert, ShieldCheck, Square, SquareTerminal, Sun, Terminal, Trash, X } from "lucide-react";
import { cva } from "class-variance-authority";
import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { Slot } from "radix-ui";
import { Fragment as Fragment$1, jsx, jsxs } from "react/jsx-runtime";
//#region src/lib/cn.ts
function cn(...inputs) {
	return twMerge(clsx(inputs));
}
//#endregion
//#region src/components/ui/badge.tsx
var badgeVariants = cva("inline-flex w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-full border border-transparent px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-[color,box-shadow] focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&>svg]:pointer-events-none [&>svg]:size-3", {
	variants: { variant: {
		default: "bg-primary text-primary-foreground [a&]:hover:bg-primary/90",
		secondary: "bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90",
		destructive: "bg-destructive text-on-solid focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40 [a&]:hover:bg-destructive/90",
		outline: "border-border text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
		ghost: "[a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
		link: "text-primary underline-offset-4 [a&]:hover:underline"
	} },
	defaultVariants: { variant: "default" }
});
function Badge({ className, variant = "default", asChild = false, ...props }) {
	const Comp = asChild ? Slot.Root : "span";
	return /* @__PURE__ */ jsx(Comp, {
		"data-slot": "badge",
		"data-variant": variant,
		className: cn(badgeVariants({ variant }), className),
		...props
	});
}
//#endregion
//#region src/components/ui/button.tsx
var buttonVariants$1 = cva("inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4", {
	variants: {
		variant: {
			default: "bg-primary text-primary-foreground hover:bg-primary/90",
			destructive: "bg-destructive text-on-solid hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40",
			outline: "border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50",
			secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80",
			ghost: "hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50",
			link: "text-primary underline-offset-4 hover:underline"
		},
		size: {
			default: "h-9 px-4 py-2 has-[>svg]:px-3",
			xs: "h-6 gap-1 rounded-md px-2 text-xs has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3",
			sm: "h-8 gap-1.5 rounded-md px-3 has-[>svg]:px-2.5",
			lg: "h-10 rounded-md px-6 has-[>svg]:px-4",
			icon: "size-9",
			"icon-xs": "size-6 rounded-md [&_svg:not([class*='size-'])]:size-3",
			"icon-sm": "size-8",
			"icon-lg": "size-10"
		}
	},
	defaultVariants: {
		variant: "default",
		size: "default"
	}
});
function Button$1({ className, variant = "default", size = "default", asChild = false, ...props }) {
	const Comp = asChild ? Slot.Root : "button";
	return /* @__PURE__ */ jsx(Comp, {
		"data-slot": "button",
		"data-variant": variant,
		"data-size": size,
		className: cn(buttonVariants$1({
			variant,
			size,
			className
		})),
		...props
	});
}
//#endregion
//#region src/components/ui/card.tsx
function Card({ className, ...props }) {
	return /* @__PURE__ */ jsx("div", {
		"data-slot": "card",
		className: cn("flex flex-col gap-6 rounded-xl border bg-card py-6 text-card-foreground shadow-sm", className),
		...props
	});
}
function CardHeader({ className, ...props }) {
	return /* @__PURE__ */ jsx("div", {
		"data-slot": "card-header",
		className: cn("@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-2 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6", className),
		...props
	});
}
function CardTitle({ className, ...props }) {
	return /* @__PURE__ */ jsx("div", {
		"data-slot": "card-title",
		className: cn("leading-none font-semibold", className),
		...props
	});
}
function CardContent({ className, ...props }) {
	return /* @__PURE__ */ jsx("div", {
		"data-slot": "card-content",
		className: cn("px-6", className),
		...props
	});
}
function CardFooter({ className, ...props }) {
	return /* @__PURE__ */ jsx("div", {
		"data-slot": "card-footer",
		className: cn("flex items-center px-6 [.border-t]:pt-6", className),
		...props
	});
}
//#endregion
//#region src/lib/panel.ts
function hostBridge() {
	if (typeof window === "undefined") return null;
	return window.wispBridge ?? window.chrome?.webview ?? null;
}
/**
* Send one bridge request and return nothing: responses arrive as a fresh
* PanelSnapshot push, never as a return value, because a panel that keeps a
* copy of the answer would be a second state holder (PLAN.md:1044).
* Throws when there is no host: a silent no-op would look like a UI that
* worked while nothing was ever asked.
*/
function requestApprovalResolution(correlationId, outcome) {
	const bridge = hostBridge();
	if (!bridge) throw new Error("panel: no native host attached - the request was not sent (ticket 33 owns the WebView2 host)");
	bridge.postMessage(JSON.stringify({
		method: "panel.approval.request",
		correlationId,
		outcome
	}));
}
/**
* One request on the panel's only envelope: the same postMessage pipe the
* approval card uses (C17). There is no second channel, and a composer that
* needed one would be a composer with authority the native side never gave it.
*
* Every composer request carries requestId + source: the id is what lets an
* audit line and a card point back at the user's own click, and the source is
* the native side's routing claim (panel.ParseComposerRequest refuses anything
* that does not name "panel-composer"). It is NOT authority - the payload still
* passes the same guards on the Go side.
*
* Throws when no host is attached: a silent no-op here would look like a UI
* that accepted the user's intent while nothing was ever asked (ticket 83's
* family again).
*/
var requestSeq = 0;
function nextRequestId() {
	requestSeq += 1;
	const uuid = globalThis.crypto?.randomUUID?.();
	return `pc-${requestSeq.toString()}-${uuid ?? "no-randomuuid"}`;
}
function sendRequest(method, payload) {
	const bridge = hostBridge();
	if (!bridge) throw new Error("panel: no native host attached - the request was not sent (ticket 33 owns the WebView2 host)");
	bridge.postMessage(JSON.stringify({
		method,
		requestId: nextRequestId(),
		source: "panel-composer",
		...payload
	}));
}
/**
* Ask the native side to switch the permission mode. This is intent only: the
* call that actually changes the mode is perm.Store.Set on the Go side, which
* still costs the R20/M4 L2 strong confirmation and writes the audit line.
* Switching INTO the档 that stops asking questions is the one that needs it,
* and a renderer that could do it silently would be approving on nobody's
* behalf - which is why nothing below can record an outcome.
*/
function requestModeSwitch(to) {
	sendRequest("panel.mode.request", { to });
}
/**
* Ask the native side to scope this session to a folder. The string is only a
* request: C26 resolves it, reparse/rewrite refusals come back as a refusal the
* panel displays, and the audit line is written natively.
*/
function requestWorkspaceChange(path) {
	sendRequest("panel.workspace.request", { path });
}
/**
* Send one attachment's BYTES to the native side.
*
* dataBase64 is the whole file, base64-encoded, because the envelope carries
* text. The native side decodes it and the decoded length is what the guards
* judge (panel.DecodeAttachmentPayload), so a truncated encode cannot pass as a
* complete file. Metadata alone would be a lie about the user's pick.
*/
function sendAttachmentBytes(name, declaredMime, sizeBytes, dataBase64) {
	sendRequest("panel.attachment.add", {
		name,
		declaredMime,
		sizeBytes,
		dataBase64
	});
}
/** Read a picked or pasted File and hand its bytes to the native side. */
function submitAttachment(file) {
	return new Promise((resolve, reject) => {
		const reader = new FileReader();
		reader.onerror = () => reject(/* @__PURE__ */ new Error(`panel: 读取 ${file.name} 失败`));
		reader.onload = () => {
			const buf = reader.result;
			if (!(buf instanceof ArrayBuffer)) {
				reject(/* @__PURE__ */ new Error(`panel: ${file.name} 的读取结果不是字节，未发送`));
				return;
			}
			sendAttachmentBytes(file.name, file.type, buf.byteLength, toBase64(buf));
			resolve();
		};
		reader.readAsArrayBuffer(file);
	});
}
/** Send one composed message: text plus the references native accepted. */
function sendMessage(text, attachments) {
	sendRequest("panel.message.send", {
		text,
		attachments
	});
}
function toBase64(buf) {
	const bytes = new Uint8Array(buf);
	let chunk = "";
	const step = 32768;
	for (let i = 0; i < bytes.length; i += step) chunk += String.fromCharCode(...bytes.subarray(i, i + step));
	return btoa(chunk);
}
async function requestHostAnswer(envelope) {
	if (typeof window === "undefined") return null;
	const door = window.wispDispatch;
	if (typeof door !== "function") return null;
	const answer = await Promise.resolve(door(JSON.stringify(envelope)));
	return typeof answer === "string" ? answer : null;
}
/** Ask the host for the settings view (config.get). The reply is ticket 248's
prose plus one line starting with the wisp-config-view marker; parse it
with parseConfigReply (lib/config-view.ts) and never by reading the prose. */
function requestConfigView() {
	return requestHostAnswer({
		method: "config.get",
		requestId: nextRequestId(),
		source: "panel-composer"
	});
}
/** One whitelisted settings write (config.set), value in configValue. The
credential field must NOT come through here - it has its own function so
a value can never be placed on the shared envelope's value key. */
function requestConfigWrite(write) {
	const envelope = {
		method: "config.set",
		requestId: nextRequestId(),
		source: "panel-composer",
		configField: write.field
	};
	if (write.provider !== void 0) envelope.configProvider = write.provider;
	if (write.model !== void 0) envelope.configModel = write.model;
	if (write.value !== void 0) envelope.configValue = write.value;
	return requestHostAnswer(envelope);
}
/** The credential write: the value rides credentialValue and only that key.
The reply is a status sentence and never echoes the value. */
function requestCredentialWrite(provider, value) {
	return requestHostAnswer({
		method: "config.set",
		requestId: nextRequestId(),
		source: "panel-composer",
		configField: "provider_credential",
		configProvider: provider,
		credentialValue: value
	});
}
/** Ask the host to fetch one provider's model catalog (config.models.list,
ticket 320's read port; the method name is registered on the Go side by
320-r2). The reply is prose plus one line starting with the
wisp-models-list marker; parse it with parseModelsListReply
(lib/config-view.ts). Until the Go half lands the door refuses the name
with its own sentence, which the page shows in place. */
function requestProviderModels(provider) {
	return requestHostAnswer({
		method: "config.models.list",
		requestId: nextRequestId(),
		source: "panel-composer",
		configProvider: provider
	});
}
//#endregion
//#region src/components/l2-approval-card.tsx
/** The header icon slot. A table keyed on the real gated tool names, falling
back to the name's first segment, and to FileSearch past that - a guessed
glyph per unknown tool would decorate the card with meaning nobody chose. */
var TOOL_ICONS = {
	"fs.delete": Trash,
	"fs.rename": PencilLine,
	"fs.write": FilePen,
	"fs.mv": FolderInput,
	"shell.exec": Terminal
};
var TOOL_PREFIX_ICONS = {
	fs: FileSearch,
	shell: Terminal
};
function toolIcon(tool) {
	const exact = TOOL_ICONS[tool];
	if (exact) return exact;
	return TOOL_PREFIX_ICONS[tool.split(".")[0] ?? ""] ?? FileSearch;
}
/**
* The level badge: PLAN's outline capsule - a 1px stroke in the level's hue,
* never a fill (the harness fails the card for the two fill spellings it
* knows). A level with no twin (Deny) takes the strong stroke and shows its
* raw string rather than an invented label.
*/
function LevelBadge({ level }) {
	const strong = level === "L2" || level === "Deny";
	if (level === "L1") return /* @__PURE__ */ jsxs(Badge, {
		className: "rounded-full border-accent/40 px-1.5 text-accent-ink",
		variant: "outline",
		children: [/* @__PURE__ */ jsx(ShieldCheck, {
			"aria-hidden": "true",
			className: "size-3"
		}), /* @__PURE__ */ jsx("span", { children: "L1 轻确认" })]
	});
	return /* @__PURE__ */ jsxs(Badge, {
		className: cn("rounded-full px-1.5", strong ? "border-red/40 text-red" : "border-line text-ink-2"),
		variant: "outline",
		children: [strong && /* @__PURE__ */ jsx(ShieldAlert, {
			"aria-hidden": "true",
			className: "size-3"
		}), /* @__PURE__ */ jsx("span", {
			className: strong ? "uppercase tracking-wide" : "",
			children: strong ? level === "L2" ? "L2 强确认" : level : level
		})]
	});
}
/** 12px semibold uppercase section head, the library's muted ink. */
function SectionTitle({ children }) {
	return /* @__PURE__ */ jsx("p", {
		className: "mb-2 text-[12px] font-semibold uppercase tracking-[0.06em] text-ink-3",
		children
	});
}
/**
* The argv, verbatim - SPEC-06 §7 / the frozen row want 完整参数, never a
* summary, so nothing here is elided. Collapsed: one line per element, tool
* first. Expanded: the same lines with an index and JSON quoting, because
* join() cannot show an empty argument, a leading space or a trailing
* newline - and those are exactly the shapes that change what a command
* deletes. The expand adds precision; it never unlocks withheld text.
*/
function ParamsBlock({ args, expanded, tool }) {
	return /* @__PURE__ */ jsx("div", {
		className: "overflow-x-auto rounded-control border border-line bg-inset px-3 py-2.5 shadow-hairline",
		children: expanded ? /* @__PURE__ */ jsx("ul", {
			className: "m-0 flex list-none flex-col gap-0.5 p-0",
			children: [tool, ...args].map((element, i) => /* @__PURE__ */ jsxs("li", {
				className: "whitespace-pre font-mono text-[11.5px] leading-[1.45] text-ink",
				children: [/* @__PURE__ */ jsxs("span", {
					className: "text-ink-3",
					children: [String(i).padStart(2, "0"), " "]
				}), JSON.stringify(element)]
			}, i))
		}) : /* @__PURE__ */ jsx("pre", {
			className: "m-0 whitespace-pre font-mono text-[11.5px] leading-[1.45] text-ink",
			children: [tool, ...args].join("\n")
		})
	});
}
/** Rules that fired, as C19 ids - the operator reads these against SPEC-06 §3.
The title repeats the id itself; inventing rule prose here would put words
in the assessor's mouth. Outline capsule, mono 10.5px, the id in the
title attribute for the audit hover. */
function RuleList({ rulesHit }) {
	if (rulesHit.length === 0) return /* @__PURE__ */ jsx("span", {
		className: "text-[11.5px] text-ink-3",
		children: "no rule id recorded"
	});
	return /* @__PURE__ */ jsx("ul", {
		"aria-label": "matched risk rules",
		className: "m-0 flex list-none flex-wrap gap-1 p-0",
		children: rulesHit.map((rule) => /* @__PURE__ */ jsx("li", { children: /* @__PURE__ */ jsx(Badge, {
			className: "rounded-full border-line px-1.5 text-ink-2",
			style: {
				fontFamily: "var(--font-mono-stack)",
				fontSize: "10.5px",
				height: "20px"
			},
			title: rule,
			variant: "outline",
			children: rule
		}) }, rule))
	});
}
/** The reason line, with the fail-closed wording when there is none: an absent
reason must never render as an empty slot beside a card that looks safe
(Q-23). Warn colour, not Tailwind amber - the token is the palette here. */
function ReasonLine({ view }) {
	if (!view.reasonKnown || view.reason.trim() === "") return /* @__PURE__ */ jsxs("p", {
		className: "text-[12.5px] leading-[1.6] text-ink-2",
		children: [/* @__PURE__ */ jsx("span", {
			className: "font-medium text-warn",
			children: "信息不足："
		}), "风险评估未给出原因文本，本卡不能推断为\"无风险\"。"]
	});
	return /* @__PURE__ */ jsx("p", {
		className: "text-[12.5px] leading-[1.6] text-ink-2",
		children: view.reason
	});
}
/**
* The standing row the frozen L2 row puts at the foot of the card: where the
* allow actually happens, drawn with the 20px ball glyph whose ring breathes
* on the 2s period the row's motion cell names. Loop animation is allowed
* here and only here: the non-idle state list names the approval wait, and
* this card does not exist while the panel is idle.
*/
function BallApproveHint() {
	return /* @__PURE__ */ jsxs("p", {
		className: "flex items-center gap-1.5 text-[11px] text-ink-2",
		children: [/* @__PURE__ */ jsxs("svg", {
			"aria-hidden": "true",
			className: "shrink-0 text-ink-2",
			height: "20",
			viewBox: "0 0 20 20",
			width: "20",
			children: [/* @__PURE__ */ jsx("circle", {
				className: "l2-ball-ring",
				cx: "10",
				cy: "10",
				fill: "none",
				r: "8",
				stroke: "currentColor",
				strokeWidth: "1"
			}), /* @__PURE__ */ jsx("circle", {
				cx: "10",
				cy: "10",
				fill: "currentColor",
				r: "5"
			})]
		}), "点击悬浮球以批准"]
	});
}
function L2ApprovalCard({ view, onIntent }) {
	const [expanded, setExpanded] = useState(false);
	const send = (outcome) => {
		requestApprovalResolution(view.correlationId, outcome);
		onIntent?.(view.correlationId, outcome);
	};
	const ToolIcon = toolIcon(view.tool);
	const strong = view.level === "L2" || view.level === "Deny";
	return /* @__PURE__ */ jsxs(Card, {
		className: "relative w-full gap-0 overflow-hidden rounded-card border-line bg-surface py-0 text-ink shadow-card",
		children: [
			/* @__PURE__ */ jsx("div", {
				"aria-hidden": "true",
				className: cn("h-[2px] w-full shrink-0", strong ? "bg-stop" : "bg-accent")
			}),
			/* @__PURE__ */ jsx(CardHeader, {
				className: "block px-4 pb-2 pt-3",
				children: /* @__PURE__ */ jsx("div", {
					className: "flex items-start justify-between gap-3",
					children: /* @__PURE__ */ jsxs("div", {
						className: "flex min-w-0 items-center gap-2.5",
						children: [/* @__PURE__ */ jsx("div", {
							className: "flex size-8 shrink-0 items-center justify-center rounded-control bg-inset text-ink-2 shadow-hairline",
							children: /* @__PURE__ */ jsx(ToolIcon, {
								"aria-hidden": "true",
								className: "size-4"
							})
						}), /* @__PURE__ */ jsxs("div", {
							className: "min-w-0",
							children: [/* @__PURE__ */ jsx("p", {
								className: "text-[13px] font-medium text-ink",
								children: "需要你的确认"
							}), /* @__PURE__ */ jsxs(CardTitle, {
								className: "mt-0.5 flex min-w-0 items-center gap-2 text-[11.5px] font-normal leading-snug text-ink-2",
								children: [/* @__PURE__ */ jsx("span", {
									className: "truncate font-mono",
									children: view.tool
								}), /* @__PURE__ */ jsx(LevelBadge, { level: view.level })]
							})]
						})]
					})
				})
			}),
			/* @__PURE__ */ jsxs(CardContent, {
				className: "flex flex-col gap-2.5 px-4 pb-3",
				children: [
					/* @__PURE__ */ jsxs("div", { children: [/* @__PURE__ */ jsx(SectionTitle, { children: "完整参数" }), /* @__PURE__ */ jsx(ParamsBlock, {
						args: view.args,
						expanded,
						tool: view.tool
					})] }),
					/* @__PURE__ */ jsxs("div", { children: [
						/* @__PURE__ */ jsx(SectionTitle, { children: "为什么需要确认" }),
						/* @__PURE__ */ jsx(ReasonLine, { view }),
						/* @__PURE__ */ jsx("div", {
							className: "mt-2",
							children: /* @__PURE__ */ jsx(RuleList, { rulesHit: view.rulesHit })
						})
					] }),
					view.callChain.length > 0 && /* @__PURE__ */ jsx("p", {
						className: "font-mono text-[11.5px] text-ink-2",
						children: `调用链：${view.callChain.join(" > ")}`
					}),
					view.sessionOverrideBlocked && /* @__PURE__ */ jsx("p", {
						className: "text-[11.5px] font-medium text-red",
						children: "本调用带 C25 污染标记，会话级授权对它无效。"
					}),
					/* @__PURE__ */ jsxs("div", {
						className: "flex items-start gap-1.5 text-[11px] leading-[1.6] text-ink-3",
						children: [/* @__PURE__ */ jsx(Info, {
							"aria-hidden": "true",
							className: "mt-px size-3.5 shrink-0"
						}), /* @__PURE__ */ jsx("span", { children: "批准权在原生侧 · 面板不提供允许按钮" })]
					})
				]
			}),
			/* @__PURE__ */ jsxs(CardFooter, {
				className: "flex-col items-stretch gap-2 border-t border-line px-4 py-2.5",
				children: [/* @__PURE__ */ jsx(BallApproveHint, {}), /* @__PURE__ */ jsxs("div", {
					className: "flex items-center gap-2",
					children: [
						/* @__PURE__ */ jsx("span", {
							className: "mr-auto text-[11px] text-ink-3",
							children: `判定来自 ${view.decidedBy === "native" ? "原生风险评估" : view.decidedBy}`
						}),
						/* @__PURE__ */ jsx(Button$1, {
							className: "rounded-control bg-red text-white hover:bg-red/90",
							onClick: () => send("refuse"),
							size: "sm",
							variant: "destructive",
							children: "拒绝"
						}),
						/* @__PURE__ */ jsx(Button$1, {
							className: "rounded-control border-line bg-surface text-ink hover:bg-hover",
							onClick: () => setExpanded((e) => !e),
							size: "sm",
							variant: "outline",
							children: expanded ? "收起参数" : "查看完整参数"
						})
					]
				})]
			})
		]
	});
}
//#endregion
//#region src/components/ai-native/prompt-bar.tsx
/** the last @word or /word being typed, if any - upstream logic verbatim. */
function parseToken(draft) {
	const match = /(^|\s)([@/])([\w-]*)$/.exec(draft);
	if (!match) return null;
	return {
		kind: match[2] === "@" ? "at" : "slash",
		query: match[3].toLowerCase(),
		start: match.index + match[1].length
	};
}
function PromptBar({ variant = "Rounded", placeholder = "输入消息…", sources = [], commands = [], models = [], model, onModelChange, permissions = [], permission, onPermissionChange, efforts = [], effort, onEffortChange, attachments = [], onRemoveAttachment, onAttach, onMicToggle, onSend, streaming = false, onStop, autoFocus = false, className }) {
	const pill = variant === "Pill";
	const [draft, setDraft] = useState("");
	const [menu, setMenu] = useState(null);
	const [active, setActive] = useState(0);
	const [modelOpen, setModelOpen] = useState(false);
	const [effortOpen, setEffortOpen] = useState(false);
	const [permissionOpen, setPermissionOpen] = useState(false);
	const [listening, setListening] = useState(false);
	const [expanded, setExpanded] = useState(false);
	const inputRef = useRef(null);
	const measureRef = useRef(null);
	const controlsRef = useRef(null);
	const modelRef = useRef(null);
	const token = parseToken(draft);
	const menuKind = token ? token.kind === "at" ? "at" : "slash" : menu === "attach" ? "attach" : null;
	const rows = (menuKind === "at" ? sources.filter((s) => s.name.toLowerCase().includes(token?.query ?? "")) : menuKind === "slash" ? commands.filter((c) => c.name.toLowerCase().includes(token?.query ?? "")) : sources.map((s) => ({
		key: s.key,
		name: s.name,
		desc: s.desc
	}))).map((r) => ({
		key: r.key,
		name: r.name,
		desc: "desc" in r ? r.desc ?? "" : ""
	}));
	const activeModel = model ? models.find((m) => m.key === model) ?? {
		key: model,
		name: model
	} : models[0];
	const activePermission = permissions.find((p) => p.key === permission);
	const activeEffort = effort ? efforts.find((e) => e.key === effort) ?? {
		key: effort,
		name: effort
	} : efforts[0];
	useLayoutEffect(() => {
		const input = inputRef.current;
		const controls = controlsRef.current;
		const measure = measureRef.current;
		const modelButton = modelRef.current;
		if (!input || !controls || !measure || !modelButton) return;
		const fixedControlsWidth = 84 + modelButton.offsetWidth;
		const inlineInputWidth = controls.clientWidth - fixedControlsWidth - 16;
		const needsFullWidth = draft.includes("\n") || measure.offsetWidth + 8 > inlineInputWidth;
		if (needsFullWidth !== expanded) setExpanded(needsFullWidth);
		const minHeight = 28;
		const maxHeight = 100;
		input.style.height = "0px";
		const contentHeight = input.scrollHeight;
		input.style.height = `${Math.min(Math.max(contentHeight, minHeight), maxHeight)}px`;
		input.style.overflowY = contentHeight > maxHeight ? "auto" : "hidden";
	}, [draft, expanded]);
	const closeMenus = () => {
		setMenu(null);
		setModelOpen(false);
		setEffortOpen(false);
		setPermissionOpen(false);
	};
	const pick = (row) => {
		if (menuKind === "attach") onAttach?.();
		else if (token) {
			const insert = menuKind === "at" ? `@${row.name} ` : `${row.name} `;
			setDraft(draft.slice(0, token.start) + insert);
		} else setDraft(`${draft}${row.name} `);
		setMenu(null);
		inputRef.current?.focus();
	};
	const canSend = draft.trim().length > 0 || attachments.length > 0;
	const send = () => {
		if (!canSend) return;
		onSend?.(draft);
		setDraft("");
		closeMenus();
	};
	return /* @__PURE__ */ jsxs("div", {
		className: cn("relative w-full", className),
		children: [
			menuKind && /* @__PURE__ */ jsxs("div", {
				className: "absolute inset-x-0 bottom-full z-10 mb-2 rounded-card bg-surface p-1 shadow-raised",
				style: {
					animation: "pop-in 180ms var(--ease-out-strong) both",
					transformOrigin: "bottom center"
				},
				children: [
					rows.map((row, i) => /* @__PURE__ */ jsxs("button", {
						type: "button",
						onMouseDown: (event) => event.preventDefault(),
						onClick: () => pick(row),
						className: cn("flex h-9 w-full items-center gap-2.5 rounded-chip px-2 text-left", i === active && "bg-hover"),
						children: [
							/* @__PURE__ */ jsx("span", {
								className: "flex size-5.5 shrink-0 items-center justify-center text-ink-2",
								children: menuKind === "slash" ? /* @__PURE__ */ jsx(Square, {
									"aria-hidden": "true",
									size: 11
								}) : /* @__PURE__ */ jsx(FileText, {
									"aria-hidden": "true",
									size: 14
								})
							}),
							/* @__PURE__ */ jsx("span", {
								className: "shrink-0 text-[12.5px] font-medium text-ink",
								children: row.name
							}),
							/* @__PURE__ */ jsx("span", {
								className: "min-w-0 flex-1 truncate text-[12px] text-ink-3",
								children: row.desc
							})
						]
					}, row.key)),
					rows.length === 0 && /* @__PURE__ */ jsx("div", {
						className: "flex h-9 items-center px-2 text-[12px] text-ink-3",
						children: "没有匹配项"
					}),
					/* @__PURE__ */ jsx("div", {
						className: "mt-1 border-t border-line px-2 pb-1 pt-1.5 text-[11px] text-ink-3",
						children: menuKind === "at" ? "输入以搜索来源与文件" : menuKind === "slash" ? "输入以搜索命令" : "附加来源与文件"
					})
				]
			}),
			modelOpen && /* @__PURE__ */ jsx("div", {
				className: "absolute right-0 bottom-full z-10 mb-2 w-44 rounded-card bg-surface p-1 shadow-raised",
				style: {
					animation: "pop-in 180ms var(--ease-out-strong) both",
					transformOrigin: "bottom right"
				},
				children: models.map((m) => /* @__PURE__ */ jsxs("button", {
					type: "button",
					onMouseDown: (event) => event.preventDefault(),
					onClick: () => {
						onModelChange?.(m.key);
						setModelOpen(false);
						inputRef.current?.focus();
					},
					className: "flex h-7.5 w-full items-center gap-2 rounded-chip px-2 text-left hover:bg-hover",
					children: [
						/* @__PURE__ */ jsx("span", {
							className: "min-w-0 flex-1 truncate text-[12.5px] font-medium text-ink",
							children: m.name
						}),
						/* @__PURE__ */ jsx("span", {
							className: "shrink-0 text-[11px] text-ink-3",
							children: m.tag
						}),
						/* @__PURE__ */ jsx("span", {
							className: cn("shrink-0 text-ink", m.key === activeModel?.key ? "" : "invisible"),
							children: /* @__PURE__ */ jsx(Check, {
								"aria-hidden": "true",
								size: 13,
								strokeWidth: 2.5
							})
						})
					]
				}, m.key))
			}),
			effortOpen && efforts.length > 0 && /* @__PURE__ */ jsxs("div", {
				className: "absolute right-0 bottom-full z-10 mb-2 w-40 rounded-card bg-surface p-1 shadow-raised",
				style: {
					animation: "pop-in 180ms var(--ease-out-strong) both",
					transformOrigin: "bottom right"
				},
				children: [/* @__PURE__ */ jsx("p", {
					className: "px-1.5 pb-0.5 pt-1 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3",
					children: "思考强度"
				}), efforts.map((e) => /* @__PURE__ */ jsxs("button", {
					type: "button",
					onMouseDown: (event) => event.preventDefault(),
					onClick: () => {
						onEffortChange?.(e.key);
						setEffortOpen(false);
						inputRef.current?.focus();
					},
					className: "flex h-7.5 w-full items-center gap-2 rounded-chip px-2 text-left hover:bg-hover",
					children: [/* @__PURE__ */ jsx("span", {
						className: "min-w-0 flex-1 truncate text-[12.5px] font-medium text-ink",
						children: e.name
					}), /* @__PURE__ */ jsx("span", {
						className: cn("shrink-0 text-ink", e.key === activeEffort?.key ? "" : "invisible"),
						children: /* @__PURE__ */ jsx(Check, {
							"aria-hidden": "true",
							size: 13,
							strokeWidth: 2.5
						})
					})]
				}, e.key))]
			}),
			permissionOpen && permissions.length > 0 && /* @__PURE__ */ jsxs("div", {
				className: "absolute left-0 bottom-full z-10 mb-2 w-56 rounded-card bg-surface p-1 shadow-raised",
				style: {
					animation: "pop-in 180ms var(--ease-out-strong) both",
					transformOrigin: "bottom left"
				},
				children: [/* @__PURE__ */ jsx("p", {
					className: "px-1.5 pb-0.5 pt-1 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3",
					children: "权限档"
				}), permissions.map((p) => /* @__PURE__ */ jsxs("button", {
					type: "button",
					onMouseDown: (event) => event.preventDefault(),
					onClick: () => {
						onPermissionChange?.(p.key);
						setPermissionOpen(false);
						inputRef.current?.focus();
					},
					className: "flex w-full items-center gap-2 rounded-chip px-2 py-1.5 text-left hover:bg-hover",
					children: [/* @__PURE__ */ jsxs("span", {
						className: "min-w-0 flex-1",
						children: [/* @__PURE__ */ jsx("span", {
							className: "block truncate text-[12.5px] font-medium text-ink",
							children: p.name
						}), p.gated ? /* @__PURE__ */ jsx("span", {
							className: "mt-0.5 block text-[10.5px] text-warn",
							children: "切过去要原生侧强确认"
						}) : null]
					}), /* @__PURE__ */ jsx("span", {
						className: cn("shrink-0 text-ink", p.key === activePermission?.key ? "" : "invisible"),
						children: /* @__PURE__ */ jsx(Check, {
							"aria-hidden": "true",
							size: 13,
							strokeWidth: 2.5
						})
					})]
				}, p.key))]
			}),
			/* @__PURE__ */ jsxs("div", {
				className: cn("relative isolate flex flex-col gap-1.5 overflow-hidden border border-line bg-surface p-1.5 shadow-card transition-[border-color] duration-150 focus-within:border-line-strong", pill ? attachments.length > 0 || expanded ? "rounded-[24px]" : "rounded-full" : "rounded-[14px]"),
				children: [
					/* @__PURE__ */ jsx("span", {
						ref: measureRef,
						"aria-hidden": "true",
						className: "pointer-events-none absolute invisible whitespace-pre text-[13px] leading-[18px]",
						children: draft
					}),
					attachments.length > 0 && /* @__PURE__ */ jsx("div", {
						className: cn("flex flex-wrap gap-1.5 pt-0.5", pill ? "px-1" : "px-0.5"),
						children: attachments.map((file) => /* @__PURE__ */ jsxs("span", {
							className: cn("flex h-6.5 items-center gap-1.5 bg-field py-1 pr-1 pl-1.5 text-[11.5px] text-ink-2 shadow-hairline", pill ? "rounded-full" : "rounded-chip"),
							style: { animation: "pop-in 200ms var(--ease-out-strong) both" },
							children: [
								/* @__PURE__ */ jsx(Paperclip, {
									"aria-hidden": "true",
									size: 12
								}),
								/* @__PURE__ */ jsx("span", {
									className: "max-w-36 truncate",
									children: file.name
								}),
								onRemoveAttachment ? /* @__PURE__ */ jsx("button", {
									type: "button",
									"aria-label": `移除 ${file.name}`,
									onClick: () => onRemoveAttachment?.(file.key),
									className: cn("flex size-4 items-center justify-center text-ink-3 transition-colors duration-100 hover:bg-line/70 hover:text-ink", pill ? "rounded-full" : "rounded-[4px]"),
									children: /* @__PURE__ */ jsx(X, {
										"aria-hidden": "true",
										size: 10,
										strokeWidth: 2.5
									})
								}) : null
							]
						}, file.key))
					}),
					/* @__PURE__ */ jsxs("div", {
						ref: controlsRef,
						className: cn("grid items-end gap-x-1 gap-y-1.5", expanded ? "grid-cols-[minmax(0,1fr)_auto_28px_28px]" : "grid-cols-[28px_minmax(0,1fr)_auto_28px_28px]"),
						children: [
							/* @__PURE__ */ jsx("button", {
								type: "button",
								"aria-label": "添加附件与来源",
								"aria-expanded": menu === "attach",
								onClick: () => {
									setModelOpen(false);
									setMenu((current) => current === "attach" ? null : "attach");
									setActive(0);
									inputRef.current?.focus();
								},
								className: cn("flex size-7 shrink-0 items-center justify-center justify-self-start text-ink-3 transition-[background-color,color,transform] duration-150 hover:bg-hover hover:text-ink active:scale-[0.94]", pill ? "rounded-full" : "rounded-[8px]", menu === "attach" && "bg-hover text-ink", expanded ? "col-start-1 row-start-2" : "col-start-1 row-start-1"),
								children: /* @__PURE__ */ jsx(Plus, {
									"aria-hidden": "true",
									size: 16,
									strokeWidth: 2
								})
							}),
							/* @__PURE__ */ jsx("textarea", {
								ref: inputRef,
								rows: 1,
								value: draft,
								autoFocus,
								onChange: (event) => {
									setDraft(event.target.value);
									setMenu(null);
								},
								onKeyDown: (event) => {
									if (menuKind && rows.length > 0) {
										if (event.key === "ArrowDown" || event.key === "ArrowUp") {
											event.preventDefault();
											setActive((current) => (current + (event.key === "ArrowDown" ? 1 : rows.length - 1)) % rows.length);
											return;
										}
										if (event.key === "Enter" && !event.shiftKey || event.key === "Tab") {
											event.preventDefault();
											pick(rows[active]);
											return;
										}
									}
									if (event.key === "Escape") {
										closeMenus();
										return;
									}
									if (event.key === "Enter" && !event.shiftKey && !event.nativeEvent.isComposing) {
										event.preventDefault();
										send();
									}
								},
								placeholder: listening ? "正在听写…" : placeholder,
								"aria-label": "Prompt",
								className: cn("min-h-7 min-w-0 w-full resize-none bg-transparent px-1 py-[5px] text-[13px] leading-[18px] text-ink outline-none [overflow-wrap:anywhere] placeholder:text-ink-3", expanded ? "col-span-full col-start-1 row-start-1" : "col-start-2 row-start-1")
							}),
							models.length > 0 && /* @__PURE__ */ jsxs("button", {
								ref: modelRef,
								type: "button",
								"aria-expanded": modelOpen,
								"aria-label": "选择模型",
								onClick: () => {
									setMenu(null);
									setEffortOpen(false);
									setPermissionOpen(false);
									setModelOpen((current) => !current);
								},
								className: cn("flex h-7 shrink-0 items-center gap-1 px-1.5 text-[12px] font-medium text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink", pill ? "rounded-full" : "rounded-[8px]", expanded ? "col-start-2 row-start-2" : "col-start-3 row-start-1"),
								children: [activeModel?.name ?? "未设置", /* @__PURE__ */ jsx("span", {
									className: "text-ink-3",
									children: /* @__PURE__ */ jsx(ChevronDown, {
										"aria-hidden": "true",
										size: 11,
										strokeWidth: 2.4
									})
								})]
							}),
							efforts.length > 0 && /* @__PURE__ */ jsxs("button", {
								type: "button",
								"aria-expanded": effortOpen,
								"aria-label": "思考强度",
								onClick: () => {
									setMenu(null);
									setModelOpen(false);
									setPermissionOpen(false);
									setEffortOpen((current) => !current);
								},
								className: cn("flex h-7 shrink-0 items-center gap-1 px-1.5 text-[12px] font-medium transition-colors duration-150", effortOpen ? "bg-hover text-ink" : "text-ink-2 hover:bg-hover hover:text-ink", pill ? "rounded-full" : "rounded-[8px]"),
								children: [activeEffort?.name ?? "未设置", /* @__PURE__ */ jsx("span", {
									className: "text-ink-3",
									children: /* @__PURE__ */ jsx(ChevronDown, {
										"aria-hidden": "true",
										size: 10,
										strokeWidth: 2.4
									})
								})]
							}),
							/* @__PURE__ */ jsx("button", {
								type: "button",
								"aria-label": listening ? "停止听写" : "开始听写",
								"aria-pressed": listening,
								onClick: () => {
									setListening((current) => {
										onMicToggle?.(!current);
										return !current;
									});
								},
								className: cn("flex size-7 shrink-0 items-center justify-center transition-[background-color,color,transform] duration-150 active:scale-[0.94]", pill ? "rounded-full" : "rounded-[8px]", listening ? "bg-accent-tint text-accent-ink" : "text-ink-3 hover:bg-hover hover:text-ink", expanded ? "col-start-3 row-start-2" : "col-start-4 row-start-1"),
								children: listening ? /* @__PURE__ */ jsx("span", {
									className: "flex h-3.5 items-center gap-[2.5px]",
									children: [
										0,
										1,
										2
									].map((i) => /* @__PURE__ */ jsx("span", {
										className: "w-[2.5px] rounded-full bg-current",
										style: {
											height: "100%",
											animation: `eq-bounce 900ms ease-in-out ${i * 150}ms infinite`
										}
									}, i))
								}) : /* @__PURE__ */ jsx(Mic, {
									"aria-hidden": "true",
									size: 15,
									strokeWidth: 2
								})
							}),
							/* @__PURE__ */ jsx("button", {
								type: "button",
								"aria-label": streaming ? "停止" : "发送",
								disabled: streaming ? false : !canSend,
								onClick: () => streaming ? onStop?.() : send(),
								className: cn("flex size-7 shrink-0 items-center justify-center transition-[background-color,color,transform] duration-200 enabled:active:scale-[0.94]", pill ? "rounded-full" : "rounded-[8px]", expanded ? "col-start-4 row-start-2" : "col-start-5 row-start-1"),
								style: streaming ? {
									background: "var(--red)",
									color: "var(--surface)"
								} : {
									background: canSend ? "var(--ink)" : "var(--line-strong)",
									color: canSend ? "var(--surface)" : "var(--ink-2)"
								},
								children: streaming ? /* @__PURE__ */ jsx(Square, {
									"aria-hidden": "true",
									size: 13,
									strokeWidth: 2.4
								}) : /* @__PURE__ */ jsx(ArrowUp, {
									"aria-hidden": "true",
									size: 16,
									strokeWidth: 2.4
								})
							})
						]
					})
				]
			}),
			permissions.length > 0 ? /* @__PURE__ */ jsx("div", {
				className: "mt-1.5 flex flex-wrap items-center gap-1.5 px-0.5",
				children: permissions.map((p) => {
					const on = p.key === activePermission?.key;
					return /* @__PURE__ */ jsxs("button", {
						type: "button",
						"aria-pressed": on,
						title: p.gated ? "切到这一档要原生侧强确认" : p.name,
						onClick: () => {
							setPermissionOpen(false);
							onPermissionChange?.(p.key);
						},
						className: cn("inline-flex h-6 items-center gap-1 rounded-full px-2.5 text-[11px] font-medium transition-colors duration-100", on ? p.gated ? "bg-warn-soft text-warn" : "bg-accent-tint text-accent-ink" : "bg-inset text-ink-2 hover:bg-hover"),
						children: [p.name, on && p.gated ? /* @__PURE__ */ jsx("span", {
							className: "text-[10px]",
							children: "需强确认"
						}) : null]
					}, p.key);
				})
			}) : null
		]
	});
}
//#endregion
//#region src/components/context-meter.tsx
/** `75.9M tok` - the dsh bottom-bar shape: a short count, not a digit
avalanche. Below 1000 stays exact because below that nobody wants a unit. */
function shortTokens(n) {
	if (n < 1e3) return String(n);
	if (n < 1e6) return `${(n / 1e3).toFixed(n < 1e4 ? 1 : 0)}k`;
	return `${(n / 1e6).toFixed(1)}M`;
}
function percent(fraction) {
	return `${Math.round(fraction * 100)}%`;
}
/** The dot between parts. Muted so it reads as punctuation, not content. */
function Dot() {
	return /* @__PURE__ */ jsx("span", {
		"aria-hidden": "true",
		className: "shrink-0 text-ink-3",
		children: "·"
	});
}
function ContextMeter({ usage, className }) {
	if (!(usage.tokens !== void 0 || usage.cacheHitRate !== void 0 || usage.occupancy !== void 0)) return /* @__PURE__ */ jsx("p", {
		className: cn("truncate text-[10.5px] text-ink-3", className),
		children: "上下文与用量还没有收到数据，面板不自己造这个数字。"
	});
	return /* @__PURE__ */ jsxs("div", {
		"aria-label": "上下文与用量",
		className: cn("flex items-center gap-1.5 text-[10.5px] tabular-nums text-ink-3", className),
		children: [
			usage.tokens !== void 0 ? /* @__PURE__ */ jsxs("span", {
				className: "shrink-0",
				children: [/* @__PURE__ */ jsx("span", {
					className: "text-ink-2",
					children: shortTokens(usage.tokens)
				}), " tok"]
			}) : null,
			usage.tokens !== void 0 && usage.cacheHitRate !== void 0 ? /* @__PURE__ */ jsx(Dot, {}) : null,
			usage.cacheHitRate !== void 0 ? /* @__PURE__ */ jsxs("span", {
				className: "shrink-0",
				children: ["缓存命中 ", percent(usage.cacheHitRate)]
			}) : null,
			usage.occupancy !== void 0 ? /* @__PURE__ */ jsxs(Fragment$1, { children: [usage.tokens !== void 0 || usage.cacheHitRate !== void 0 ? /* @__PURE__ */ jsx(Dot, {}) : null, /* @__PURE__ */ jsxs("span", {
				className: "shrink-0",
				children: ["占用 ", percent(usage.occupancy)]
			})] }) : null
		]
	});
}
//#endregion
//#region src/components/ai-native/button.tsx
var filledShadow = "shadow-btn";
var buttonVariants = cva(`inline-flex items-center justify-center font-medium select-none
   transition-[transform,background-color,opacity] duration-150 ease-out
   active:scale-[0.96] disabled:opacity-50 disabled:pointer-events-none`, {
	variants: {
		variant: {
			primary: `bg-ink text-canvas hover:opacity-90 dark:bg-ink dark:text-canvas ${filledShadow}`,
			secondary: "bg-surface text-ink shadow-btn hover:bg-inset aria-expanded:bg-hover",
			ghost: "bg-hover-2 text-ink hover:bg-line-strong",
			accent: `bg-accent text-white hover:bg-accent-ink ${filledShadow}`,
			success: `bg-green text-white hover:brightness-95 ${filledShadow}`,
			quiet: "text-ink hover:bg-hover"
		},
		size: {
			xs: "h-7 rounded-full px-2.5 text-[12px] font-normal leading-none gap-1",
			sm: "h-[27px] px-3 text-[13px] leading-none rounded-full gap-1.5",
			md: "px-4 py-[9px] text-sm leading-none rounded-full gap-2"
		}
	},
	defaultVariants: {
		variant: "secondary",
		size: "md"
	}
});
function Button({ variant, size, className, ...props }) {
	return /* @__PURE__ */ jsx("button", {
		className: cn(buttonVariants({
			variant,
			size
		}), className),
		...props
	});
}
//#endregion
//#region src/lib/use-dismiss.ts
function useDismiss(ref, open, onDismiss) {
	useEffect(() => {
		if (!open) return;
		function onPointerDown(event) {
			const node = ref.current;
			if (!node) return;
			const target = event.target;
			if (target instanceof Node && node.contains(target)) return;
			onDismiss();
		}
		function onKeyDown(event) {
			if (event.key === "Escape") onDismiss();
		}
		document.addEventListener("pointerdown", onPointerDown, true);
		document.addEventListener("keydown", onKeyDown);
		return () => {
			document.removeEventListener("pointerdown", onPointerDown, true);
			document.removeEventListener("keydown", onKeyDown);
		};
	}, [
		ref,
		open,
		onDismiss
	]);
}
//#endregion
//#region src/components/settings-controls.tsx
function SettingsCard({ title, description, children }) {
	return /* @__PURE__ */ jsxs("section", {
		className: "rounded-card border border-line bg-surface",
		children: [/* @__PURE__ */ jsxs("header", {
			className: "border-b border-line px-4 py-2.5",
			children: [/* @__PURE__ */ jsx("h3", {
				className: "text-[12.5px] font-semibold text-ink",
				children: title
			}), description ? /* @__PURE__ */ jsx("p", {
				className: "mt-0.5 text-[11px] leading-[1.6] text-ink-3",
				children: description
			}) : null]
		}), /* @__PURE__ */ jsx("div", {
			className: "px-4 py-3",
			children
		})]
	});
}
/** The dot that carries a tone without words where a word would crowd. */
function StateDot$1({ tone, className }) {
	return /* @__PURE__ */ jsx("span", {
		"aria-hidden": "true",
		className: cn("inline-block size-2 shrink-0 rounded-full", tone === "ok" && "bg-green", tone === "warn" && "bg-orange", tone === "error" && "bg-red", tone === "neutral" && "bg-ink-3", className)
	});
}
/** A compact tag (smaller than StatusPill): row-level states read as a tag,
not as a banner. */
function StateTag({ tone, text }) {
	return /* @__PURE__ */ jsxs("span", {
		className: cn("inline-flex h-5 shrink-0 items-center gap-1.5 rounded-full px-2 text-[11px] font-medium leading-none", tone === "ok" && "bg-green-tint text-green", tone === "warn" && "bg-orange-tint text-orange", tone === "error" && "bg-red-tint text-red", tone === "neutral" && "bg-inset text-ink-2"),
		children: [/* @__PURE__ */ jsx(StateDot$1, { tone }), text]
	});
}
function FeedbackLine({ feedback }) {
	if (feedback === null) return null;
	return /* @__PURE__ */ jsx("p", {
		className: cn("text-[11.5px] leading-[1.6]", feedback.tone === "ok" && "text-green", feedback.tone === "error" && "text-red", feedback.tone === "warn" && "text-orange", feedback.tone === "neutral" && "text-ink-3"),
		role: feedback.tone === "error" ? "alert" : void 0,
		children: feedback.text
	});
}
function EmptyState({ text, actionLabel, onAction }) {
	return /* @__PURE__ */ jsxs("div", {
		className: "flex flex-col items-start gap-3 py-2",
		children: [/* @__PURE__ */ jsx("p", {
			className: "max-w-[420px] text-[12px] leading-[1.7] text-ink-2",
			children: text
		}), actionLabel && onAction ? /* @__PURE__ */ jsx(Button, {
			onClick: onAction,
			size: "sm",
			variant: "primary",
			children: actionLabel
		}) : null]
	});
}
/** The sentinel a model list uses for its "手动输入" row. Not a legal id, so it
can never collide with a real model name. */
var MANUAL_ENTRY = "\0manual-entry";
function Picker({ ariaLabel, value, options, placeholder = "选择", emptyText = "没有可选项", onSelect, className }) {
	const [open, setOpen] = useState(false);
	const ref = useRef(null);
	useDismiss(ref, open, () => setOpen(false));
	const current = options.find((o) => o.value === value);
	return /* @__PURE__ */ jsxs("div", {
		className: cn("relative", className),
		ref,
		children: [/* @__PURE__ */ jsxs("button", {
			"aria-expanded": open,
			"aria-label": ariaLabel,
			className: "flex h-[30px] w-full items-center gap-1.5 rounded-control bg-inset px-2.5 text-left text-[12px] text-ink shadow-hairline transition-colors duration-150 hover:bg-hover",
			onClick: () => setOpen((v) => !v),
			type: "button",
			children: [/* @__PURE__ */ jsx("span", {
				className: cn("min-w-0 flex-1 truncate font-mono", current ? "text-ink" : "text-ink-3"),
				children: current?.value ?? placeholder
			}), /* @__PURE__ */ jsx(ChevronDown, {
				"aria-hidden": "true",
				className: "shrink-0 text-ink-3",
				size: 13,
				strokeWidth: 2
			})]
		}), open ? /* @__PURE__ */ jsx("div", {
			className: "absolute left-0 top-full z-40 mt-1 max-h-64 min-w-full overflow-y-auto rounded-card border border-line bg-surface p-1 shadow-overlay",
			style: { animation: "pop-in 160ms var(--ease-out-strong) both" },
			children: options.length === 0 ? /* @__PURE__ */ jsx("p", {
				className: "px-2 py-2 text-[11.5px] text-ink-3",
				children: emptyText
			}) : options.map((option) => /* @__PURE__ */ jsxs("button", {
				className: "flex w-full items-center gap-2 rounded-chip px-2 py-1.5 text-left transition-colors duration-100 hover:bg-hover",
				onClick: () => {
					onSelect(option.value);
					setOpen(false);
				},
				type: "button",
				children: [
					/* @__PURE__ */ jsx("span", {
						className: "min-w-0 flex-1 truncate font-mono text-[12px] text-ink",
						children: option.value
					}),
					option.note ? /* @__PURE__ */ jsx("span", {
						className: "shrink-0 text-[10.5px] text-ink-3",
						children: option.note
					}) : null,
					option.value === value ? /* @__PURE__ */ jsx(Check, {
						"aria-hidden": "true",
						className: "shrink-0 text-accent-ink",
						size: 12,
						strokeWidth: 2.4
					}) : null
				]
			}, option.value))
		}) : null]
	});
}
/** The appearance-family picker: one row per vocabulary entry the snapshot
carried, the current family checked, and the pick handed to the parent as
a value - this control performs no write itself, the same division the
rest of this file keeps (controls show and pass intent; lib/config-actions
is where a request goes out).

The vocabulary and the current value come from the snapshot verbatim; this
file spells neither family name, so the schema list stays the one source.
A snapshot with a vocabulary but an empty style renders the picker with
nothing checked (the page says "not configured" beside it); a null
section renders the honest no-reader sentence, never a guessed family. */
function BallStylePicker({ section, onSelect, className }) {
	if (section === null) return /* @__PURE__ */ jsx("span", {
		className: cn("flex h-[30px] items-center text-[11.5px] text-ink-3", className),
		children: "快照里还没有球风格这一节（本机没有装配球风格读口）"
	});
	const options = section.vocabulary.map((v) => ({ value: v }));
	const checked = options.some((o) => o.value === section.style);
	return /* @__PURE__ */ jsx(Picker, {
		ariaLabel: "球风格",
		className,
		emptyText: "配置里没有可切换的风格",
		onSelect,
		options,
		placeholder: checked ? "选择风格" : "未设置；选择一个",
		value: section.style
	});
}
/** The model half of a pair: the catalog list when the read has one, a typed id
when it does not (`manualFallback`), and the reason printed when neither can
produce a value. The typed id is the SECONDARY path - it exists because a
provider nobody has configured yet cannot have a catalog list. */
function ModelField({ row, provider, model, manual, onChange, onManual, manualFallback, className }) {
	const models = row !== null && row.models !== null && row.models.length > 0 ? row.models : null;
	if (models !== null && !manual) {
		const options = [...models.map((name) => ({ value: name })), ...manualFallback ? [{
			value: MANUAL_ENTRY,
			note: "手动输入"
		}] : []];
		return /* @__PURE__ */ jsx(Picker, {
			ariaLabel: "模型",
			className,
			onSelect: (next) => {
				if (next === "\0manual-entry") {
					onManual();
					onChange("");
				} else onChange(next);
			},
			options,
			placeholder: "选择模型",
			value: model
		});
	}
	if (manualFallback) return /* @__PURE__ */ jsx("input", {
		"aria-label": "模型 id",
		className: cn("h-[30px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3", className),
		onChange: (e) => onChange(e.currentTarget.value),
		placeholder: "模型 id",
		value: model
	});
	return /* @__PURE__ */ jsx("span", {
		className: cn("flex h-[30px] items-center text-[11.5px] text-ink-3", className),
		children: provider === "" ? "先选供应商" : row === null ? "这家供应商还没有配置" : "这家供应商暂无模型"
	});
}
/** Provider + model as one control. Owns only the manual/automatic flip; both
values are the parent's (a saved pair must stay visible after a reload). */
function PairFields({ providers, provider, model, onProvider, onModel, manualFallback, providerLabel = "供应商" }) {
	const [manualFor, setManualFor] = useState(null);
	const manual = manualFor === provider && provider !== "";
	const row = providers.find((p) => p.name === provider) ?? null;
	return /* @__PURE__ */ jsxs(Fragment$1, { children: [/* @__PURE__ */ jsx(Picker, {
		ariaLabel: providerLabel,
		className: "w-[180px] shrink-0",
		emptyText: "还没有已配置的供应商",
		onSelect: onProvider,
		options: providers.map((p) => ({
			value: p.name,
			note: p.credentialSet ? "密钥已保存" : "密钥未保存"
		})),
		placeholder: "选择供应商",
		value: provider
	}), /* @__PURE__ */ jsx(ModelField, {
		className: "w-[220px] shrink-0",
		manual,
		manualFallback,
		model,
		onChange: onModel,
		onManual: () => setManualFor(provider),
		provider,
		row
	})] });
}
//#endregion
//#region src/lib/config-feedback.ts
/** The page-level three states. "未读取" is a fact about the READ, not about
the config: nothing has been set AND nothing has been answered must not
look alike. */
function configStatus(view) {
	if (!view.readable) return {
		label: "未读取",
		tone: "neutral"
	};
	if (!view.configuredKnown) return {
		label: "未读取",
		tone: "neutral"
	};
	return view.configured ? {
		label: "已配置",
		tone: "ok"
	} : {
		label: "未配置",
		tone: "warn"
	};
}
/** One provider's key state as a dot's words. */
function keyLabel(row) {
	if (row.credentialError !== "") return {
		label: "密钥未读取",
		tone: "neutral"
	};
	return row.credentialSet ? {
		label: "密钥已保存",
		tone: "ok"
	} : {
		label: "密钥未保存",
		tone: "warn"
	};
}
/** The refusal matcher: sealed sentences, most specific first. Each maps to
the one line a person can act on. Order matters where two sentences share
a word (a model row's "already exists" must not be read as a provider's). */
var REFUSALS = [
	{
		match: "没有对应的写入腿",
		text: "这项暂不支持在界面上保存"
	},
	{
		match: "不是面板 composer 通路的能力入口",
		text: "暂不支持拉取模型列表"
	},
	{
		match: "建目录条目必须点名模型 id",
		text: "请先填入模型 id"
	},
	{
		match: "这一步是「建条目」",
		text: "这个模型已经在目录里了"
	},
	{
		match: "预设行请走",
		text: "这是内置供应商名，请用「内置供应商」入口添加"
	},
	{
		match: "不是内置预设名",
		text: "只能添加内置供应商"
	},
	{
		match: "billing",
		text: "计费模式只支持「套餐」或「按量」"
	},
	{
		match: "protocol",
		text: "协议只支持 openai-chat / openai-responses / anthropic"
	},
	{
		match: "两半",
		text: "自定义接入需要协议和基址两项"
	},
	{
		match: "already exists",
		text: "这家供应商已经存在（不会覆盖）"
	},
	{
		match: "必须一起点名",
		text: "供应商与模型要一起选好再保存"
	},
	{
		match: "空值不会被写成默认",
		text: "这项不能留空"
	},
	{
		match: "has no model",
		text: "这家供应商名下没有这个模型"
	},
	{
		match: "没有 model",
		text: "这家供应商名下没有这个模型"
	},
	{
		match: "凭据存储或配置对象未接入",
		text: "暂不支持在本机保存密钥"
	},
	{
		match: "凭据字段是空的",
		text: "请先填入密钥"
	},
	{
		match: "引用形状不合法",
		text: "密钥格式不对"
	},
	{
		match: "密钥没有写进凭据存储",
		text: "密钥没有保存成功"
	},
	{
		match: "密钥已进凭据存储",
		text: "密钥已保存，但供应商没有更新成功"
	}
];
/** The refusal half both write translators share: a matched sealed sentence
becomes its human line; anything else is the generic "保存失败", never a
guess dressed up as a diagnosis. */
function refusalFeedback(reply) {
	for (const rule of REFUSALS) if (reply.includes(rule.match)) return {
		tone: "error",
		text: rule.text
	};
	return {
		tone: "error",
		text: "保存失败"
	};
}
/** Translate one write reply into the single sentence shown to the user.
`null` means no host answered at all. */
function translateWriteReply(reply) {
	if (reply === null) return {
		tone: "error",
		text: "未连接到本机服务，未保存"
	};
	if (reply.includes("被拒绝")) return refusalFeedback(reply);
	if (reply.includes("密钥已录入")) return {
		tone: "ok",
		text: restartTier(reply) ? "密钥已保存，重启后生效" : "密钥已保存"
	};
	if (reply.includes("已写入")) return {
		tone: "ok",
		text: restartTier(reply) ? "已保存，重启后生效" : "已保存"
	};
	return {
		tone: "error",
		text: "保存失败"
	};
}
/** The composer row's model/thinking switches (ticket 187) translate through
this one instead of translateWriteReply. The Go receipt states the
effective tier itself ("这一项在下一次任务开始时生效" for a thinking
write, "这一项要重启进程…" for a model write - config_handlers.go's
tierSentence), and the two halves of the row have DIFFERENT tiers: a
translator that collapsed both into one generic "已保存" would let the
model half look in force before a restart and the tier half wait for one.
So the receipt's own tier wording is kept. Refusals match exactly like
translateWriteReply. */
function translateTieredWriteReply(reply) {
	if (reply === null) return {
		tone: "error",
		text: "未连接到本机服务，未保存"
	};
	if (reply.includes("被拒绝")) return refusalFeedback(reply);
	if (reply.includes("已写入")) {
		if (reply.includes("下一次任务")) return {
			tone: "ok",
			text: "已保存，下一次任务开始时生效"
		};
		if (reply.includes("重启")) return {
			tone: "ok",
			text: "已保存，重启进程后生效"
		};
		return {
			tone: "ok",
			text: "已保存"
		};
	}
	return {
		tone: "error",
		text: "保存失败"
	};
}
/** All settings writes are restart-tier today; the receipt carries the tier
sentence and this reads it rather than assuming. */
function restartTier(reply) {
	return reply.includes("重启");
}
/** The model-list fetch's answer when it did not carry a payload: a refusal
gets its matched human line; anything else (the host's D37 prose) is shown
verbatim - the fetch has no receipt to translate, and inventing a summary
for prose this page never parsed would be a guess. */
function translateFetchReply(reply) {
	if (reply === null) return {
		tone: "error",
		text: "未连接到本机服务，未拉取"
	};
	if (reply.includes("被拒绝")) {
		for (const rule of REFUSALS) if (reply.includes(rule.match)) return {
			tone: "error",
			text: rule.text
		};
	}
	const text = reply.trim();
	return {
		tone: "error",
		text: text === "" ? "拉取失败" : text
	};
}
/** Local validation lines (checked before anything is sent - the pair rule is
stated at save time, not by disabling the button, which reads as a black
box to the person looking at it). */
var COPY = {
	pairMissing: "供应商与模型都要选好再保存",
	chainSegmentMissing: "每一段都要同时选好供应商和模型",
	chainEmpty: "至少保留一段；链不能为空",
	chatPairMissing: "供应商与模型都要选好再保存",
	testConnectionUnavailable: "暂不支持连接测试",
	testKeyUnavailable: "暂不支持密钥测试",
	deleteUnavailable: "暂不支持删除供应商",
	editUnavailable: "暂不支持编辑",
	baseURLSaved: "已保存，重启后生效",
	noProvidersYet: "还没有供应商。添加一个供应商并填入密钥，就可以选择模型。",
	addProvider: "添加供应商",
	save: "保存",
	saving: "保存中…"
};
//#endregion
//#region src/lib/config-view.ts
/** The wire marker. Spelled once here; panel.SettingsPayloadMarker is the
other side of the same string. */
var CONFIG_VIEW_MARKER = "wisp-config-view:";
/** The model-catalog marker (ticket 320's config.models.list; the Go half is
320-r2). Same shape as the view marker: prose, then exactly one line. */
var MODELS_LIST_MARKER = "wisp-models-list:";
/** Parse one config.models.list reply. Only models[].id is read; unknown
fields are ignored on purpose (the wire shape is the Go side's to grow).
The prose half is never parsed - on failure it is returned verbatim for
display. */
function parseModelsListReply(reply) {
	const at = reply.indexOf(MODELS_LIST_MARKER);
	if (at < 0) return {
		kind: "no-payload",
		raw: reply
	};
	const line = reply.slice(at + 17).split("\n", 1)[0] ?? "";
	let data;
	try {
		data = JSON.parse(line);
	} catch {
		return {
			kind: "no-payload",
			raw: reply
		};
	}
	if (typeof data !== "object" || data === null) return {
		kind: "no-payload",
		raw: reply
	};
	const payload = data;
	const provider = typeof payload.provider === "string" ? payload.provider : "";
	const models = [];
	if (Array.isArray(payload.models)) for (const entry of payload.models) {
		if (typeof entry !== "object" || entry === null) continue;
		const id = entry.id;
		if (typeof id === "string" && id.trim() !== "") models.push(id.trim());
	}
	return {
		kind: "parsed",
		provider,
		models
	};
}
/** Parse one config.get reply. The prose half is returned as-is (display
only); the structured half is the line that starts with the marker. */
function parseConfigReply(reply) {
	const at = reply.indexOf(CONFIG_VIEW_MARKER);
	if (at < 0) return {
		kind: "no-payload",
		raw: reply
	};
	const line = reply.slice(at + 17).split("\n", 1)[0] ?? "";
	let data;
	try {
		data = JSON.parse(line);
	} catch {
		return {
			kind: "no-payload",
			raw: reply
		};
	}
	if (typeof data !== "object" || data === null) return {
		kind: "no-payload",
		raw: reply
	};
	const err = data.payloadError;
	if (typeof err === "string" && err !== "") return {
		kind: "payload-error",
		message: err,
		raw: reply
	};
	return {
		kind: "parsed",
		view: data,
		prose: reply.slice(0, at).trimEnd()
	};
}
/** Read the settings view once, and expose reload() for after a write. */
function useConfigView() {
	const [state, setState] = useState({ phase: "loading" });
	const seq = useRef(0);
	const reload = useCallback(() => {
		const mine = ++seq.current;
		requestConfigView().then((reply) => {
			if (mine !== seq.current) return;
			if (reply === null) {
				setState({ phase: "no-host" });
				return;
			}
			const parsed = parseConfigReply(reply);
			if (parsed.kind === "parsed") setState({
				phase: "ready",
				view: parsed.view,
				prose: parsed.prose,
				raw: reply
			});
			else if (parsed.kind === "payload-error") setState({
				phase: "payload-error",
				message: parsed.message,
				raw: reply
			});
			else setState({
				phase: "unreadable",
				raw: reply
			});
		});
	}, []);
	useEffect(() => {
		reload();
	}, [reload]);
	return {
		state,
		reload
	};
}
//#endregion
//#region src/lib/config-actions.ts
/** One whitelisted write (config.set), answered as the sentence to show. */
async function performConfigWrite(write) {
	return translateWriteReply(await requestConfigWrite(write));
}
/** The composer row's model switch (ticket 187): one offerable catalog row,
written as the role_chat_pair field the model page already uses, with the
receipt keeping the Go receipt's own tier sentence (model writes are
restart tier). Same door, same envelope, tiered translation. */
async function performComposerModelWrite(provider, model) {
	return translateTieredWriteReply(await requestConfigWrite({
		field: "role_chat_pair",
		provider,
		model
	}));
}
/** The composer row's thinking-tier switch (ticket 187): the value is one of
the snapshot's own thinkingLevels verbatim - never a fourth word invented
here - and the receipt keeps the next-task tier sentence the host stamps
on this field. */
async function performComposerThinkingWrite(level) {
	return translateTieredWriteReply(await requestConfigWrite({
		field: "role_chat_thinking",
		value: level
	}));
}
/** The credential write: the value rides credentialValue and only that key. */
async function performCredentialWrite(provider, value) {
	return translateWriteReply(await requestCredentialWrite(provider, value));
}
/** Fetch one provider's model catalog (config.models.list). */
async function performModelsFetch(provider) {
	const reply = await requestProviderModels(provider);
	const parsed = reply === null ? null : parseModelsListReply(reply);
	if (parsed !== null && parsed.kind === "parsed") return {
		ok: true,
		models: parsed.models
	};
	return {
		ok: false,
		feedback: translateFetchReply(reply)
	};
}
//#endregion
//#region src/components/composer-row.tsx
/** The one error sentence shape for a request the envelope refused to carry
(no host attached is the usual case - lib/panel's sendRequest throws).
Told to the user on the feedback line, never swallowed. */
function requestError(err) {
	return String(err instanceof Error ? err.message : err);
}
/** Thinking-tier enum words -> the Chinese the effort button shows. The write
value stays the enum word verbatim (PromptEffort.key carries it); a word
this table does not know shows its own text instead of being dropped. */
var TIER_LABELS = {
	off: "关",
	low: "低",
	medium: "中",
	high: "高"
};
function ComposerRow({ modelSwitch, composer, permission, permissions }) {
	const [feedback, setFeedback] = useState(null);
	const [workspaceDraft, setWorkspaceDraft] = useState("");
	const [attaching, setAttaching] = useState(false);
	const fileInput = useRef(null);
	const stored = (composer?.attachments ?? []).filter((a) => a.stored);
	const refused = (composer?.attachments ?? []).filter((a) => !a.stored);
	const models = (modelSwitch?.models ?? []).map((m) => ({
		key: `${m.provider}/${m.model}`,
		name: m.model,
		tag: m.provider
	}));
	const efforts = (modelSwitch?.thinkingLevels ?? []).map((level) => ({
		key: level,
		name: TIER_LABELS[level] ?? level
	}));
	const currentPair = modelSwitch !== void 0 && (modelSwitch.current.provider !== "" || modelSwitch.current.model !== "") ? `${modelSwitch.current.provider}/${modelSwitch.current.model}` : "";
	const currentThinking = modelSwitch !== void 0 && modelSwitch.currentThinking !== "" ? modelSwitch.currentThinking : "off";
	async function switchModel(pairKey) {
		setFeedback({
			tone: "neutral",
			text: "保存中…"
		});
		const cut = pairKey.indexOf("/");
		setFeedback(await performComposerModelWrite(cut < 0 ? pairKey : pairKey.slice(0, cut), cut < 0 ? "" : pairKey.slice(cut + 1)));
	}
	async function switchThinking(level) {
		setFeedback({
			tone: "neutral",
			text: "保存中…"
		});
		setFeedback(await performComposerThinkingWrite(level));
	}
	function askMode(to) {
		if (composer !== void 0 && to === composer.mode.current) return;
		try {
			requestModeSwitch(to);
		} catch (err) {
			setFeedback({
				tone: "error",
				text: `切档请求没有发出去：${requestError(err)}`
			});
		}
	}
	function switchWorkspace() {
		const path = workspaceDraft.trim();
		if (path === "") return;
		try {
			requestWorkspaceChange(path);
			setWorkspaceDraft("");
		} catch (err) {
			setFeedback({
				tone: "error",
				text: `工作区切换请求没有发出去：${requestError(err)}`
			});
		}
	}
	async function attach(files) {
		if (!files || files.length === 0) return;
		setAttaching(true);
		try {
			for (const f of Array.from(files)) await submitAttachment(f);
		} catch (err) {
			setFeedback({
				tone: "error",
				text: `附件没有送出去：${requestError(err)}`
			});
		} finally {
			setAttaching(false);
			if (fileInput.current) fileInput.current.value = "";
		}
	}
	function send(text) {
		const body = text.trim();
		if (body === "" && stored.length === 0) return;
		try {
			sendMessage(body, stored);
		} catch (err) {
			setFeedback({
				tone: "error",
				text: `消息没有发出去：${requestError(err)}`
			});
		}
	}
	return /* @__PURE__ */ jsxs(Fragment$1, { children: [
		modelSwitch !== void 0 && models.length === 0 && modelSwitch.reason !== "" ? /* @__PURE__ */ jsxs("p", {
			className: "mb-1.5 text-[11px] leading-[1.6] text-ink-3",
			role: "note",
			children: [modelSwitch.reason, "（模型选择器不画）"]
		}) : null,
		composer ? /* @__PURE__ */ jsxs("div", {
			className: "mb-1.5 flex items-center gap-1.5 text-[11.5px] text-ink-3",
			children: [
				/* @__PURE__ */ jsx(FolderCog, {
					"aria-hidden": "true",
					className: "size-3.5 shrink-0"
				}),
				/* @__PURE__ */ jsx("span", {
					className: "min-w-0 flex-1 truncate font-mono",
					title: composer.workspace.set ? composer.workspace.canonical : composer.workspace.reason,
					children: composer.workspace.set ? composer.workspace.canonical : composer.workspace.reason
				}),
				composer.workspace.reparse ? /* @__PURE__ */ jsx("span", {
					className: "shrink-0 text-warn",
					children: "经 reparse 点（例外名单）"
				}) : null,
				/* @__PURE__ */ jsx("input", {
					"aria-label": "工作区路径（由原生侧解析与校验）",
					className: "w-32 shrink-0 rounded-chip border border-line bg-surface px-1.5 py-0.5\n              font-mono text-[11px] text-ink outline-none transition-colors duration-100\n              focus:border-line-strong",
					onChange: (e) => setWorkspaceDraft(e.target.value),
					placeholder: composer.workspace.spelling || "D:\\work\\project",
					value: workspaceDraft
				}),
				/* @__PURE__ */ jsxs("button", {
					"aria-label": "请求切换工作区",
					className: "flex shrink-0 items-center gap-0.5 rounded-chip px-1 py-0.5 text-ink-3\n              transition-colors duration-100 hover:bg-hover hover:text-ink\n              disabled:pointer-events-none disabled:opacity-50",
					disabled: workspaceDraft.trim() === "",
					onClick: switchWorkspace,
					type: "button",
					children: [/* @__PURE__ */ jsx(ChevronRight, {
						"aria-hidden": "true",
						className: "size-3"
					}), "切换"]
				})
			]
		}) : null,
		/* @__PURE__ */ jsx(ContextMeter, {
			usage: {},
			className: "mb-1.5"
		}),
		/* @__PURE__ */ jsx(PromptBar, {
			attachments: stored.map((a) => ({
				key: a.id,
				name: a.name
			})),
			effort: modelSwitch !== void 0 ? currentThinking : void 0,
			efforts,
			model: modelSwitch !== void 0 ? currentPair : void 0,
			models,
			onAttach: composer ? () => fileInput.current?.click() : void 0,
			onEffortChange: (key) => void switchThinking(key),
			onModelChange: (key) => void switchModel(key),
			onPermissionChange: askMode,
			onSend: send,
			permission,
			permissions
		}),
		/* @__PURE__ */ jsx("input", {
			accept: composer && composer.acceptedAttachmentMimes.length > 0 ? composer.acceptedAttachmentMimes.join(",") : void 0,
			className: "hidden",
			multiple: true,
			onChange: (e) => void attach(e.target.files),
			ref: fileInput,
			type: "file"
		}),
		refused.length > 0 || (composer?.attachmentError ?? "") !== "" ? /* @__PURE__ */ jsxs("div", {
			className: "mt-1",
			children: [composer?.attachmentError ? /* @__PURE__ */ jsx("p", {
				className: "text-[11px] leading-[1.6] text-red",
				children: composer.attachmentError
			}) : null, refused.map((a) => /* @__PURE__ */ jsxs("p", {
				className: "text-[11px] leading-[1.6] text-red",
				children: [
					"未发送：",
					a.name,
					"——",
					a.reason
				]
			}, a.id))]
		}) : null,
		attaching ? /* @__PURE__ */ jsx("div", {
			className: "mt-1",
			children: /* @__PURE__ */ jsx(FeedbackLine, { feedback: {
				tone: "neutral",
				text: "附件读取中…"
			} })
		}) : null,
		feedback !== null ? /* @__PURE__ */ jsx("div", {
			className: "mt-1",
			children: /* @__PURE__ */ jsx(FeedbackLine, { feedback })
		}) : null
	] });
}
//#endregion
//#region src/components/ai-native/glide-menu.tsx
function GlideMenu({ children, className = "", highlightClassName = "inset-x-0 rounded-[8px] bg-hover", rowSelector = "[data-menu-row]" }) {
	const ref = useRef(null);
	const [box, setBox] = useState(null);
	const [visible, setVisible] = useState(false);
	const positioned = /\b(absolute|fixed|sticky)\b/.test(className) ? "" : "relative";
	const moveTo = (target) => {
		const container = ref.current;
		if (!(target instanceof Element) || !container) return;
		const row = target.closest(rowSelector);
		if (!(row instanceof HTMLElement) || !container.contains(row)) return;
		const containerRect = container.getBoundingClientRect();
		const rowRect = row.getBoundingClientRect();
		setBox({
			top: rowRect.top - containerRect.top,
			height: rowRect.height
		});
		setVisible(true);
	};
	return /* @__PURE__ */ jsxs("div", {
		className: `group/glide-menu ${positioned} ${className}`,
		onBlurCapture: (event) => {
			if (!ref.current?.contains(event.relatedTarget)) setVisible(false);
		},
		onFocusCapture: (event) => moveTo(event.target),
		onMouseLeave: () => setVisible(false),
		onMouseOver: (event) => moveTo(event.target),
		ref,
		children: [/* @__PURE__ */ jsx("span", {
			"aria-hidden": "true",
			className: `m-glide pointer-events-none absolute ${highlightClassName}`,
			style: {
				top: box?.top ?? 0,
				height: box?.height ?? 0,
				opacity: box && visible ? 1 : 0,
				transition: "top 220ms cubic-bezier(0.23,1,0.32,1), height 220ms cubic-bezier(0.23,1,0.32,1), opacity 150ms ease"
			}
		}), children]
	});
}
//#endregion
//#region src/components/git-branch.tsx
function GitBranchChip({ view, onCheckout, onNewWorktree }) {
	const [open, setOpen] = useState(false);
	const [creating, setCreating] = useState(false);
	const [draft, setDraft] = useState("");
	const popRef = useRef(null);
	useDismiss(popRef, open, () => setOpen(false));
	if (!view.isRepo) return null;
	return /* @__PURE__ */ jsxs("div", {
		className: "relative shrink-0",
		children: [/* @__PURE__ */ jsxs("button", {
			type: "button",
			"aria-expanded": open,
			"aria-label": `当前分支 ${view.current}`,
			onClick: () => setOpen((v) => !v),
			className: cn("flex h-6.5 items-center gap-1.5 rounded-chip px-2 text-[11.5px] transition-colors duration-150", open ? "bg-hover text-ink" : "text-ink-2 hover:bg-hover hover:text-ink"),
			children: [
				/* @__PURE__ */ jsx(GitBranch, {
					"aria-hidden": "true",
					size: 12,
					strokeWidth: 1.9
				}),
				/* @__PURE__ */ jsx("span", {
					className: "max-w-32 truncate font-mono",
					children: view.current
				}),
				/* @__PURE__ */ jsx(ChevronDown, {
					"aria-hidden": "true",
					size: 10,
					strokeWidth: 2.4
				})
			]
		}), open ? /* @__PURE__ */ jsx("div", {
			ref: popRef,
			className: "absolute left-0 top-full z-40 mt-1.5 w-64 rounded-card border border-line bg-surface p-1.5 shadow-overlay",
			style: { animation: "pop-in 160ms var(--ease-out-strong) both" },
			children: creating ? /* @__PURE__ */ jsxs("div", {
				className: "flex flex-col gap-1.5 p-1",
				children: [
					/* @__PURE__ */ jsx("p", {
						className: "text-[11px] text-ink-3",
						children: "新分支名"
					}),
					/* @__PURE__ */ jsx("input", {
						autoFocus: true,
						"aria-label": "新分支名",
						onChange: (e) => setDraft(e.target.value),
						onKeyDown: (e) => {
							if (e.key === "Escape") setCreating(false);
							if (e.key === "Enter" && draft.trim()) {
								onNewWorktree(draft.trim());
								setDraft("");
								setCreating(false);
								setOpen(false);
							}
						},
						placeholder: "feature/…",
						value: draft,
						className: "h-7 rounded-control border border-line bg-inset px-2 font-mono text-[12px] text-ink outline-none focus:border-line-strong"
					}),
					/* @__PURE__ */ jsx("p", {
						className: "text-[10.5px] leading-[1.6] text-ink-3",
						children: "新分支会开一个工作树，本地路径由原生侧定。请求已经递出去了，结果在下一份快照里。"
					}),
					/* @__PURE__ */ jsxs("div", {
						className: "flex items-center gap-1.5",
						children: [/* @__PURE__ */ jsx("button", {
							type: "button",
							disabled: !draft.trim(),
							onClick: () => {
								onNewWorktree(draft.trim());
								setDraft("");
								setCreating(false);
								setOpen(false);
							},
							className: "h-6.5 rounded-control bg-ink px-2.5 text-[11.5px] font-medium text-surface transition-opacity disabled:opacity-40",
							children: "建工作树"
						}), /* @__PURE__ */ jsx("button", {
							type: "button",
							onClick: () => setCreating(false),
							className: "h-6.5 rounded-control px-2 text-[11.5px] text-ink-2 transition-colors hover:bg-hover hover:text-ink",
							children: "取消"
						})]
					})
				]
			}) : /* @__PURE__ */ jsxs(GlideMenu, {
				className: "flex flex-col",
				highlightClassName: "inset-x-0 rounded-chip bg-hover",
				children: [
					/* @__PURE__ */ jsx("p", {
						className: "px-1.5 pb-0.5 pt-1 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3",
						children: "本地分支"
					}),
					view.local.length === 0 ? /* @__PURE__ */ jsx("p", {
						className: "px-1.5 py-1.5 text-[11.5px] text-ink-3",
						children: "只有当前这一个分支。"
					}) : view.local.map((branch) => /* @__PURE__ */ jsxs("button", {
						type: "button",
						"data-menu-row": true,
						onClick: () => {
							onCheckout(branch);
							setOpen(false);
						},
						className: "flex w-full items-center gap-2 rounded-chip px-1.5 py-1.5 text-left",
						children: [
							/* @__PURE__ */ jsx(GitBranch, {
								"aria-hidden": "true",
								className: "shrink-0 text-ink-3",
								size: 12,
								strokeWidth: 1.8
							}),
							/* @__PURE__ */ jsx("span", {
								className: "min-w-0 flex-1 truncate font-mono text-[12px] text-ink-2",
								children: branch
							}),
							/* @__PURE__ */ jsx("span", {
								className: cn("shrink-0 text-ink", branch === view.current ? "" : "invisible"),
								children: /* @__PURE__ */ jsx(Check, {
									"aria-hidden": "true",
									size: 12,
									strokeWidth: 2.5
								})
							})
						]
					}, branch)),
					view.worktrees.length > 0 ? /* @__PURE__ */ jsxs(Fragment$1, { children: [/* @__PURE__ */ jsx("p", {
						className: "px-1.5 pb-0.5 pt-2 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3",
						children: "工作树"
					}), view.worktrees.map((wt) => /* @__PURE__ */ jsxs("div", {
						className: "flex items-center gap-2 rounded-chip px-1.5 py-1.5",
						"data-menu-row": true,
						children: [/* @__PURE__ */ jsx(GitFork, {
							"aria-hidden": "true",
							className: "shrink-0 text-ink-3",
							size: 12,
							strokeWidth: 1.8
						}), /* @__PURE__ */ jsxs("span", {
							className: "min-w-0 flex-1 truncate text-[12px] text-ink-2",
							children: [/* @__PURE__ */ jsx("span", {
								className: "font-mono",
								children: wt.branch
							}), /* @__PURE__ */ jsx("span", {
								className: "ml-1.5 font-mono text-[10.5px] text-ink-3",
								children: wt.path
							})]
						})]
					}, wt.path))] }) : null,
					/* @__PURE__ */ jsx("div", {
						className: "mt-1 border-t border-line pt-1",
						children: /* @__PURE__ */ jsxs("button", {
							type: "button",
							"data-menu-row": true,
							onClick: () => setCreating(true),
							className: "flex w-full items-center gap-2 rounded-chip px-1.5 py-1.5 text-left text-[12px] text-ink-2 transition-colors hover:text-ink",
							children: [/* @__PURE__ */ jsx(Plus, {
								"aria-hidden": "true",
								size: 12,
								strokeWidth: 2.2
							}), "新建分支并开工作树"]
						})
					})
				]
			})
		}) : null]
	});
}
//#endregion
//#region src/components/harness/monitor-popover.tsx
/** D43 状态名 -> the Chinese the row shows (internal/statemachine/states.go is
the source of the names; App's TIER_NAMES is the precedent for spelling the
display map in the component that draws it). A name missing here falls to
the honest line below, never to a guessed word. */
var D43_LABELS = {
	FirstRun: "首次运行",
	Sleeping: "休眠",
	Armed: "已就绪",
	Muted: "已静音",
	Listening: "在听",
	Thinking: "思考中",
	Acting: "执行中",
	Speaking: "回话中",
	Warm: "预热",
	Conversation: "对话中",
	Confirming: "待确认（L1）",
	AwaitingApproval: "待审批（L2）",
	Settling: "收尾中",
	Downloading: "下载中",
	Unconfigured: "未配置",
	NoNetwork: "无网络",
	Error: "出错",
	Queued: "排队中",
	Stuck: "卡住",
	WatchdogAlert: "看门狗告警"
};
/** Only the failure family goes red; Settling is the one name that reads as
wrapping up. Every other known status is live work (accent); a status the
host never filed is idle grey. */
var D43_DOT = {
	Error: "failed",
	Stuck: "failed",
	WatchdogAlert: "failed",
	Settling: "done"
};
/** One row's displayed word and dot, derived from the wire fields and nothing
else. 等批准 outranks the status word: a task whose own card is on the
queue is waiting on the user, whatever its D43 name currently says. */
function agentDisplay(row) {
	if (row.statusKnown && D43_LABELS[row.status] !== void 0) return {
		label: D43_LABELS[row.status],
		dot: D43_DOT[row.status] ?? "running"
	};
	return {
		label: row.statusReason?.trim() || "宿主没有给出状态",
		dot: "idle"
	};
}
var DOT_CLASS = {
	running: "bg-accent",
	done: "bg-green",
	failed: "bg-red",
	blocked: "bg-warn",
	idle: "bg-ink-3"
};
/** The row's own truncation suffix, or null when this row's stream is whole.
Dropped outranks truncated: a dropped stream is not "partly shown", the
page for that row is GONE. Only an actual loss renders anything - a clean
row stays clean, because a permanent "0 字省略" chip on every row would be
alarm fatigue (the same argument App's header uses for the approval chip:
a badge that is always on trains you to ignore it). */
function rowBoundLabel(row) {
	if (row.streamDropped) return "流已丢";
	if (row.streamTruncated) return `省 ${row.streamElidedRunes ?? 0} 字`;
	return null;
}
/** The whole log's truncation sentence off the SECTION numbers (wire
TaskRosterSection's streamTruncated/streamElidedRunes/droppedCount, Go
:166-183), or null on a clean log. Same display rule as the row suffix:
abnormal states only. Both facts when both hold, in one line. */
function sectionBoundLine(truncated, elidedRunes, droppedCount) {
	const dropped = droppedCount ?? 0;
	if (truncated && dropped > 0) return `流已截断：省略 ${elidedRunes ?? 0} 字，已丢 ${dropped} 条`;
	if (truncated) return `流已截断：省略 ${elidedRunes ?? 0} 字`;
	if (dropped > 0) return `已丢 ${dropped} 条流`;
	return null;
}
/** 子智能体头像：conic 渐变圆片（token 色相，零字面量）。 */
function AgentAvatar() {
	return /* @__PURE__ */ jsx("span", {
		"aria-hidden": "true",
		className: "size-4 shrink-0 rounded-full",
		style: { background: "conic-gradient(var(--accent), var(--green), var(--orange), var(--accent))" }
	});
}
function StateDot({ state }) {
	return /* @__PURE__ */ jsx("span", {
		"aria-hidden": "true",
		className: cn("size-1.5 shrink-0 rounded-full", DOT_CLASS[state])
	});
}
function MonitorPopover({ open, env, agents, rosterPresent, poolInFlight, poolCap, streamTruncated, streamElidedRunes, streamDroppedCount, onPin }) {
	if (!open) return null;
	const boundLine = sectionBoundLine(streamTruncated ?? false, streamElidedRunes, streamDroppedCount);
	return /* @__PURE__ */ jsx("div", {
		className: "w-[300px] rounded-card border border-line bg-surface p-3 shadow-overlay",
		role: "dialog",
		"aria-label": "任务监控",
		style: { animation: "pop-in 180ms var(--ease-out-strong) both" },
		children: /* @__PURE__ */ jsxs(GlideMenu, {
			className: "flex flex-col",
			highlightClassName: "inset-x-0 rounded-chip bg-hover",
			children: [
				/* @__PURE__ */ jsxs("div", {
					className: "flex items-center justify-between gap-2 pb-1",
					children: [/* @__PURE__ */ jsx("p", {
						className: "text-[12.5px] font-semibold text-ink",
						children: "任务监控"
					}), /* @__PURE__ */ jsx("button", {
						"aria-label": "固定任务监控",
						className: "flex size-6 items-center justify-center rounded-control text-ink-3 transition-colors duration-150 hover:bg-hover hover:text-ink",
						onClick: onPin,
						title: "固定",
						type: "button",
						children: /* @__PURE__ */ jsx(Pin, {
							"aria-hidden": "true",
							size: 13,
							strokeWidth: 1.8
						})
					})]
				}),
				/* @__PURE__ */ jsx("p", {
					className: "px-1 pb-1 pt-1.5 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3",
					children: "环境信息"
				}),
				/* @__PURE__ */ jsx("div", {
					className: "flex flex-col gap-0.5",
					children: env.map((row) => /* @__PURE__ */ jsxs("div", {
						className: "relative flex items-center gap-2 rounded-chip px-1.5 py-1.5 text-[12px]",
						"data-menu-row": true,
						children: [/* @__PURE__ */ jsx("span", {
							className: "min-w-0 flex-1 truncate text-ink-2",
							children: row.label
						}), /* @__PURE__ */ jsx("span", {
							className: "shrink-0 font-mono text-[11.5px] tabular-nums text-ink",
							title: row.value,
							children: row.value
						})]
					}, row.label))
				}),
				/* @__PURE__ */ jsxs("div", {
					className: "flex items-center justify-between gap-2 px-1 pb-1 pt-2.5",
					children: [/* @__PURE__ */ jsx("p", {
						className: "text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3",
						children: "子智能体"
					}), poolInFlight !== void 0 && poolCap !== void 0 ? /* @__PURE__ */ jsxs("span", {
						className: "shrink-0 font-mono text-[11px] tabular-nums text-ink-3",
						title: "名册池：宿主正占着的槽位与上限",
						children: [
							"在跑 ",
							poolInFlight,
							" / 上限 ",
							poolCap
						]
					}) : null]
				}),
				/* @__PURE__ */ jsxs("div", {
					className: "flex flex-col gap-0.5",
					children: [agents.length === 0 ? /* @__PURE__ */ jsx("p", {
						className: "px-1.5 py-1.5 text-[11.5px] leading-relaxed text-ink-3",
						children: rosterPresent ? "还没有子任务在跑" : "还没收到名册（宿主没带任务这一节）"
					}) : agents.map((agent) => {
						const display = agentDisplay(agent);
						const bound = rowBoundLabel(agent);
						return /* @__PURE__ */ jsxs("div", {
							className: "relative flex items-center gap-2 rounded-chip px-1.5 py-1.5 text-[12px]",
							"data-menu-row": true,
							children: [
								/* @__PURE__ */ jsx(AgentAvatar, {}),
								/* @__PURE__ */ jsx("span", {
									className: "min-w-0 flex-1 truncate text-ink-2",
									title: agent.name,
									children: agent.name
								}),
								bound ? /* @__PURE__ */ jsx("span", {
									className: "shrink-0 font-mono text-[10.5px] tabular-nums text-warn",
									title: bound === "流已丢" ? "这一行的流页已被日志丢弃" : "这一行的流被日志折界，页面上少这么多字",
									children: bound
								}) : null,
								/* @__PURE__ */ jsx("span", {
									className: cn("shrink-0 text-[11px] tabular-nums", display.dot === "blocked" ? "text-warn" : "text-ink-3"),
									children: display.label
								}),
								/* @__PURE__ */ jsx(StateDot, { state: display.dot })
							]
						}, agent.id);
					}), boundLine ? /* @__PURE__ */ jsx("p", {
						className: "px-1.5 py-1.5 text-[11.5px] leading-relaxed text-warn",
						children: boundLine
					}) : null]
				})
			]
		})
	});
}
//#endregion
//#region src/components/palette-screen.tsx
/** A tiny name table for the icons the fixtures actually use; a name outside
it falls back to a neutral dot rather than a guessed glyph. */
var ITEM_ICONS = {
	"shield-check": ShieldCheck,
	settings: Settings,
	moon: Moon,
	terminal: Terminal,
	trash: Trash,
	"folder-input": FolderInput,
	"file-pen": FilePen
};
function ItemIcon({ icon }) {
	const Named = icon ? ITEM_ICONS[icon] : void 0;
	if (Named) return /* @__PURE__ */ jsx(Named, {
		"aria-hidden": "true",
		className: "size-3.5 shrink-0 text-ink-3"
	});
	return /* @__PURE__ */ jsx("span", {
		"aria-hidden": "true",
		className: "size-1.5 shrink-0 rounded-full bg-ink-3"
	});
}
function PaletteScreen({ query, onQueryChange, groups, onSelect }) {
	const q = query.trim().toLowerCase();
	const shown = groups.map((group) => ({
		name: group.name,
		items: q ? group.items.filter((item) => item.label.toLowerCase().includes(q)) : group.items
	})).filter((group) => group.items.length > 0);
	const empty = q.length > 0 && shown.length === 0;
	const flat = shown.flatMap((group) => group.items);
	const [active, setActive] = useState(0);
	const cursor = flat.length === 0 ? 0 : Math.min(active, flat.length - 1);
	useEffect(() => {
		setActive((current) => flat.length === 0 ? 0 : Math.min(current, flat.length - 1));
	}, [flat.length]);
	function onKey(e) {
		if (flat.length === 0) return;
		if (e.key === "ArrowDown") {
			e.preventDefault();
			setActive((current) => (current + 1) % flat.length);
			return;
		}
		if (e.key === "ArrowUp") {
			e.preventDefault();
			setActive((current) => (current - 1 + flat.length) % flat.length);
			return;
		}
		if (e.key === "Enter") {
			e.preventDefault();
			onSelect?.(flat[cursor]);
		}
	}
	let seen = 0;
	return /* @__PURE__ */ jsx("div", {
		className: "flex w-full flex-col items-stretch",
		children: /* @__PURE__ */ jsxs("div", {
			className: "w-full overflow-hidden rounded-card border border-line bg-surface shadow-card",
			children: [
				/* @__PURE__ */ jsxs("div", {
					className: "flex h-11 items-center gap-2.5 border-b border-line px-3.5 transition-colors duration-100 hover:bg-hover",
					children: [
						/* @__PURE__ */ jsxs("svg", {
							className: "shrink-0 text-ink-3",
							fill: "none",
							height: "15",
							stroke: "currentColor",
							strokeLinecap: "round",
							strokeWidth: "2",
							viewBox: "0 0 24 24",
							width: "15",
							children: [/* @__PURE__ */ jsx("circle", {
								cx: "11",
								cy: "11",
								r: "7"
							}), /* @__PURE__ */ jsx("path", { d: "M21 21l-4.3-4.3" })]
						}),
						/* @__PURE__ */ jsx("input", {
							"aria-label": "搜索命令",
							"aria-activedescendant": flat.length > 0 ? `palette-row-${cursor}` : void 0,
							className: "min-w-0 flex-1 bg-transparent text-[15px] text-ink outline-none placeholder:text-ink-3",
							onChange: (event) => onQueryChange(event.target.value),
							onKeyDown: onKey,
							placeholder: "输入命令或问题…",
							value: query
						}),
						query && /* @__PURE__ */ jsx("button", {
							"aria-label": "清空搜索",
							className: "flex size-6 items-center justify-center rounded-full text-ink-3 transition-colors duration-100 hover:bg-line/70 hover:text-ink",
							onClick: () => onQueryChange(""),
							style: { animation: "fade-in 150ms ease-out both" },
							type: "button",
							children: /* @__PURE__ */ jsx("svg", {
								fill: "none",
								height: "11",
								stroke: "currentColor",
								strokeLinecap: "round",
								strokeWidth: "2.2",
								viewBox: "0 0 24 24",
								width: "11",
								children: /* @__PURE__ */ jsx("path", { d: "M18 6L6 18M6 6l12 12" })
							})
						})
					]
				}),
				empty ? /* @__PURE__ */ jsxs("div", {
					className: "flex flex-col items-center justify-center gap-1 px-4 py-8",
					style: { animation: "fade-in 250ms ease-out both" },
					children: [
						/* @__PURE__ */ jsx("span", {
							className: "mb-1.5 flex size-8 items-center justify-center rounded-control bg-inset text-ink-3 shadow-hairline",
							children: /* @__PURE__ */ jsxs("svg", {
								fill: "none",
								height: "15",
								stroke: "currentColor",
								strokeWidth: "1.8",
								strokeLinecap: "round",
								viewBox: "0 0 24 24",
								children: [/* @__PURE__ */ jsx("circle", {
									cx: "11",
									cy: "11",
									r: "7"
								}), /* @__PURE__ */ jsx("path", { d: "M21 21l-4.3-4.3" })]
							})
						}),
						/* @__PURE__ */ jsx("span", {
							className: "text-[13px] font-medium text-ink",
							children: "没有匹配的命令"
						}),
						/* @__PURE__ */ jsx("span", {
							className: "text-[12px] text-ink-3",
							children: "换个说法再试一次"
						})
					]
				}) : /* @__PURE__ */ jsx(GlideMenu, {
					className: "p-1.5",
					highlightClassName: "inset-x-1.5 rounded-[6px] bg-hover",
					rowSelector: "[data-menu-row]",
					children: shown.map((group) => /* @__PURE__ */ jsxs("div", {
						className: "flex flex-col",
						children: [/* @__PURE__ */ jsx("p", {
							className: "px-2 pb-1 pt-2 text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3",
							children: group.name
						}), group.items.map((item) => {
							const index = seen++;
							const on = index === cursor;
							return /* @__PURE__ */ jsxs("button", {
								"aria-selected": on,
								className: cn("relative flex h-9 w-full items-center gap-2.5 rounded-[6px] py-0 pr-2 pl-2.5 text-left transition-colors duration-100", on ? "text-ink" : "text-ink-2 hover:text-ink"),
								"data-menu-row": true,
								id: `palette-row-${index}`,
								onClick: () => onSelect?.(item),
								onMouseEnter: () => setActive(index),
								role: "option",
								style: { animation: "fade-in 200ms ease-out both" },
								type: "button",
								children: [
									on ? /* @__PURE__ */ jsx("span", {
										"aria-hidden": "true",
										className: "absolute bottom-1.5 left-0.5 top-1.5 w-[2px] rounded-full bg-accent"
									}) : null,
									/* @__PURE__ */ jsx(ItemIcon, { icon: item.icon }),
									/* @__PURE__ */ jsx("span", {
										className: "min-w-0 flex-1 truncate text-[13px]",
										children: item.label
									}),
									item.hint && /* @__PURE__ */ jsx("span", {
										className: "shrink-0 text-[11px] text-ink-3",
										children: item.hint
									})
								]
							}, item.label);
						})]
					}, group.name))
				}),
				/* @__PURE__ */ jsxs("div", {
					className: "flex items-center gap-3 border-t border-line bg-inset px-3.5 py-2 text-[11px] text-ink-3",
					children: [
						/* @__PURE__ */ jsxs("span", {
							className: "flex items-center gap-1",
							children: [
								/* @__PURE__ */ jsx("span", {
									className: "kbd",
									children: "↑"
								}),
								/* @__PURE__ */ jsx("span", {
									className: "kbd",
									children: "↓"
								}),
								" 选择"
							]
						}),
						/* @__PURE__ */ jsxs("span", {
							className: "flex items-center gap-1",
							children: [/* @__PURE__ */ jsx("span", {
								className: "kbd",
								children: "Enter"
							}), " 执行"]
						}),
						/* @__PURE__ */ jsxs("span", {
							className: "flex items-center gap-1",
							children: [/* @__PURE__ */ jsx("span", {
								className: "kbd",
								children: "ESC"
							}), " 关闭"]
						})
					]
				})
			]
		})
	});
}
//#endregion
//#region src/components/react-bits/line-sidebar.tsx
var FALLOFF_CURVES = {
	linear: (p) => p,
	smooth: (p) => p * p * (3 - 2 * p),
	sharp: (p) => p * p * p
};
var DEFAULT_ITEMS = [
	"Overview",
	"Components",
	"Animations",
	"Backgrounds",
	"Showcase",
	"Playground",
	"Templates",
	"Changelog",
	"Community",
	"Resources",
	"Documentation",
	"Support"
];
function LineSidebar({ items = DEFAULT_ITEMS, accentColor = "var(--accent)", textColor = "var(--ink-3)", markerColor = "var(--line-strong)", showIndex = true, showMarker = true, proximityRadius = 100, maxShift = 30, falloff = "smooth", markerLength = 60, markerGap = 0, tickScale = .5, scaleTick = true, itemGap = 20, fontSize = 1.1, smoothing = 100, defaultActive = null, onItemClick, className = "" }) {
	const listRef = useRef(null);
	const itemRefs = useRef([]);
	const targetsRef = useRef([]);
	const currentRef = useRef([]);
	const rafRef = useRef(null);
	const lastRef = useRef(0);
	const activeRef = useRef(defaultActive);
	const smoothingRef = useRef(smoothing);
	const [activeIndex, setActiveIndex] = useState(defaultActive);
	activeRef.current = activeIndex;
	smoothingRef.current = smoothing;
	const runFrame = useCallback((now) => {
		const dt = Math.min((now - lastRef.current) / 1e3, .05);
		lastRef.current = now;
		const tau = Math.max(smoothingRef.current, 1) / 1e3;
		const k = 1 - Math.exp(-dt / tau);
		let moving = false;
		const els = itemRefs.current;
		for (let i = 0; i < els.length; i++) {
			const el = els[i];
			if (!el) continue;
			const target = Math.max(targetsRef.current[i] || 0, activeRef.current === i ? 1 : 0);
			const cur = currentRef.current[i] || 0;
			const next = cur + (target - cur) * k;
			const settled = Math.abs(target - next) < .0015;
			const value = settled ? target : next;
			currentRef.current[i] = value;
			el.style.setProperty("--effect", value.toFixed(4));
			if (!settled) moving = true;
		}
		rafRef.current = moving ? requestAnimationFrame(runFrame) : null;
	}, []);
	const startLoop = useCallback(() => {
		if (rafRef.current != null) cancelAnimationFrame(rafRef.current);
		lastRef.current = performance.now();
		rafRef.current = requestAnimationFrame(runFrame);
	}, [runFrame]);
	const handlePointerMove = useCallback((e) => {
		const list = listRef.current;
		if (!list) return;
		const rect = list.getBoundingClientRect();
		const pointerY = e.clientY - rect.top;
		const ease = FALLOFF_CURVES[falloff] ?? FALLOFF_CURVES.linear;
		const els = itemRefs.current;
		for (let i = 0; i < els.length; i++) {
			const el = els[i];
			if (!el) continue;
			const center = el.offsetTop + el.offsetHeight / 2;
			const distance = Math.abs(pointerY - center);
			targetsRef.current[i] = ease(Math.max(0, 1 - distance / proximityRadius));
		}
		startLoop();
	}, [
		falloff,
		proximityRadius,
		startLoop
	]);
	const handlePointerLeave = useCallback(() => {
		targetsRef.current = targetsRef.current.map(() => 0);
		startLoop();
	}, [startLoop]);
	const handleClick = useCallback((index, label) => {
		setActiveIndex(index);
		onItemClick?.(index, label);
	}, [onItemClick]);
	useEffect(() => {
		startLoop();
	}, [activeIndex, startLoop]);
	useEffect(() => () => {
		if (rafRef.current != null) cancelAnimationFrame(rafRef.current);
		rafRef.current = null;
	}, []);
	return /* @__PURE__ */ jsx("nav", {
		className: `line-sidebar${showMarker ? " line-sidebar--markers" : ""}${scaleTick ? " line-sidebar--scale-tick" : ""}${className ? ` ${className}` : ""}`,
		style: {
			"--accent-color": accentColor,
			"--text-color": textColor,
			"--marker-color": markerColor,
			"--marker-length": `${markerLength}px`,
			"--marker-gap": `${markerGap}px`,
			"--tick-scale": tickScale,
			"--max-shift": `${maxShift}px`,
			"--item-gap": `${itemGap}px`,
			"--font-size": `${fontSize}rem`,
			"--smoothing": `${smoothing}ms`
		},
		children: /* @__PURE__ */ jsx("ul", {
			ref: listRef,
			className: "line-sidebar__list",
			onPointerMove: handlePointerMove,
			onPointerLeave: handlePointerLeave,
			children: items.map((label, index) => /* @__PURE__ */ jsxs("li", {
				ref: (el) => {
					itemRefs.current[index] = el;
				},
				className: "line-sidebar__item",
				"aria-current": activeIndex === index ? "true" : void 0,
				onClick: () => handleClick(index, label),
				children: [showMarker ? /* @__PURE__ */ jsx("span", {
					className: "line-sidebar__marker",
					"aria-hidden": "true"
				}) : null, /* @__PURE__ */ jsxs("span", {
					className: "line-sidebar__label",
					children: [showIndex ? /* @__PURE__ */ jsx("span", {
						className: "line-sidebar__index",
						children: String(index + 1).padStart(2, "0")
					}) : null, /* @__PURE__ */ jsx("span", {
						className: "line-sidebar__text",
						children: label
					})]
				})]
			}, `${label}-${index}`))
		})
	});
}
//#endregion
//#region src/components/glide-highlight.tsx
/** Sub-pixel tolerance. Layout rects land on fractions, so an exact compare
would re-render on every frame of a scroll for no visible change. */
var EPS = .5;
/**
* The one way this tree says "you are here".
*
* Every list in the panel marks its selected row with THIS and nothing else:
* sidebar sessions, command palette rows, the review file list, the file tree,
* the task monitor's agent rows, the workspace picker, the launcher cards. One
* component rather than a className copied into nine files, because the first
* time it was copied into two they drifted and the panel stopped reading as one
* surface.
*
* It is a fixed accent bar on the left edge, deliberately NOT a background
* wash: the background is the hover channel (GlideMenu), and two things in one
* slot means moving the pointer erases where you are.
*/
function RowMark({ className }) {
	return /* @__PURE__ */ jsx("span", {
		"aria-hidden": "true",
		className: cn("pointer-events-none absolute bottom-1.5 left-0.5 top-1.5 w-[2px] rounded-full bg-accent", className)
	});
}
function sameRect(a, b) {
	return Math.abs(a.top - b.top) < EPS && Math.abs(a.height - b.height) < EPS && Math.abs(a.left - b.left) < EPS && Math.abs(a.width - b.width) < EPS;
}
function GlideHighlight({ activeKey, layoutKey, duration = 220, radius = "7px", blockClassName, className }) {
	const blockRef = useRef(null);
	const first = useRef(true);
	const [box, setBox] = useState(null);
	useLayoutEffect(() => {
		const parent = blockRef.current?.parentElement;
		if (!parent || activeKey === null) {
			setBox((current) => current === null ? current : null);
			return;
		}
		const target = parent.querySelector(`[data-glide-key="${CSS.escape(activeKey)}"]`);
		if (!target) {
			setBox((current) => current === null ? current : null);
			return;
		}
		const p = parent.getBoundingClientRect();
		const t = target.getBoundingClientRect();
		const next = {
			top: t.top - p.top,
			height: t.height,
			left: t.left - p.left,
			width: t.width
		};
		setBox((current) => current !== null && sameRect(current, next) ? current : next);
	}, [activeKey, layoutKey]);
	const instant = first.current && box === null;
	useLayoutEffect(() => {
		if (box !== null) first.current = false;
	}, [box]);
	return /* @__PURE__ */ jsx("span", {
		"aria-hidden": "true",
		className: cn("m-glide pointer-events-none absolute top-0 left-0 z-0 bg-hover", className, blockClassName),
		ref: blockRef,
		style: {
			top: box?.top ?? 0,
			height: box?.height ?? 0,
			left: box?.left ?? 0,
			width: box?.width ?? 0,
			opacity: box ? 1 : 0,
			borderRadius: radius,
			transition: instant ? "none" : `top ${duration}ms var(--ease-out-strong), left ${duration}ms var(--ease-out-strong), height ${duration}ms var(--ease-out-strong), width ${duration}ms var(--ease-out-strong), opacity 150ms ease`
		}
	});
}
//#endregion
//#region src/components/harness/right-rail.tsx
var TAB_DEFS = [
	{
		key: "review",
		label: "审查",
		closedAgo: "刚刚",
		blurb: "看这一轮改了哪些文件、每一处改动是什么"
	},
	{
		key: "files",
		label: "工作区文件",
		closedAgo: "1 小时",
		blurb: "浏览会话工作区的文件",
		keys: "Ctrl P"
	},
	{
		key: "terminal",
		label: "终端",
		closedAgo: "3 小时",
		blurb: "在会话工作区里跑命令",
		keys: "Ctrl `"
	},
	{
		key: "browser",
		label: "浏览器",
		closedAgo: "5 小时",
		blurb: "浏览网页、看生成物预览",
		keys: "Ctrl T"
	}
];
var WORKSPACE_ROOT = "D:\\work\\workspace\\projects plans\\Wisp";
var TREE = [
	{
		name: "frontend",
		kind: "dir",
		depth: 0
	},
	{
		name: "src",
		kind: "dir",
		depth: 1
	},
	{
		name: "components",
		kind: "dir",
		depth: 2
	},
	{
		name: "harness",
		kind: "dir",
		depth: 3,
		changed: true
	},
	{
		name: "right-rail.tsx",
		kind: "file",
		depth: 4,
		changed: true
	},
	{
		name: "sidebar.tsx",
		kind: "file",
		depth: 4,
		changed: true
	},
	{
		name: "palette-screen.tsx",
		kind: "file",
		depth: 4,
		changed: true
	},
	{
		name: "ai-native",
		kind: "dir",
		depth: 3
	},
	{
		name: "prompt-bar.tsx",
		kind: "file",
		depth: 4,
		locked: "fs.write"
	},
	{
		name: "styles",
		kind: "dir",
		depth: 2,
		changed: true
	},
	{
		name: "theme.css",
		kind: "file",
		depth: 3,
		changed: true
	},
	{
		name: "internal",
		kind: "dir",
		depth: 0
	},
	{
		name: "tools",
		kind: "dir",
		depth: 1
	},
	{
		name: "fs.go",
		kind: "file",
		depth: 2
	},
	{
		name: "fs_edit.go",
		kind: "file",
		depth: 2
	},
	{
		name: "docs",
		kind: "dir",
		depth: 0
	},
	{
		name: "specs",
		kind: "dir",
		depth: 1
	},
	{
		name: "tasks",
		kind: "dir",
		depth: 0
	},
	{
		name: "待办台账.json",
		kind: "file",
		depth: 1
	}
];
var REVIEW_SCOPES = [
	"未暂存",
	"已暂存",
	"全部分更改",
	"上一轮更改"
];
var REVIEW_FILES = [
	{
		name: "prompt-bar.tsx",
		dir: "frontend/src/components/ai-native/",
		additions: 214,
		removals: 0,
		tile: "bg-accent-tint",
		diff: [
			{
				sign: "+",
				text: "export function PromptBar({ sources, commands, models }: PromptBarProps) {"
			},
			{
				sign: "+",
				text: "  const token = parseToken(draft);"
			},
			{
				sign: "-",
				text: "  const draft = \"演示数据\";"
			}
		]
	},
	{
		name: "right-rail.tsx",
		dir: "frontend/src/components/harness/",
		additions: 331,
		removals: 12,
		tile: "bg-accent-tint",
		diff: [{
			sign: "+",
			text: "const [menuOpen, setMenuOpen] = useState(false);"
		}, {
			sign: "-",
			text: "  // 画布 = 给项目加个画布功能"
		}]
	},
	{
		name: "sidebar.tsx",
		dir: "frontend/src/components/harness/",
		additions: 180,
		removals: 96,
		tile: "bg-accent-tint",
		diff: [{
			sign: "+",
			text: "  // Qoder 三段式：自动化 / 工作区 / 最近任务"
		}]
	},
	{
		name: "theme.css",
		dir: "frontend/src/styles/",
		additions: 240,
		removals: 1204,
		tile: "bg-orange-tint",
		diff: [{
			sign: "-",
			text: "  /* demo 转写层全部退役 */"
		}]
	},
	{
		name: "tokens.generated.css",
		dir: "frontend/src/styles/",
		additions: 87,
		removals: 298,
		tile: "bg-orange-tint",
		diff: [{
			sign: "+",
			text: "  --page: var(--page);  /* beautiful-ui 主题收编 */"
		}]
	},
	{
		name: "gen-tokens.mjs",
		dir: "frontend/scripts/",
		additions: 190,
		removals: 467,
		tile: "bg-green-tint",
		diff: [{
			sign: "+",
			text: "  // 第四代：库主题逐字收编"
		}]
	},
	{
		name: "harness-app.ts",
		dir: "frontend/src/fixtures/",
		additions: 322,
		removals: 0,
		tile: "bg-green-tint",
		diff: [{
			sign: "+",
			text: "  export const SESSIONS = [ /* 4 条演示会话 */ ];"
		}]
	},
	{
		name: "VENDORED.md",
		dir: "frontend/",
		additions: 47,
		removals: 3,
		tile: "bg-field"
	}
];
var TERMINAL_LINES = [
	{
		kind: "cmd",
		text: "$ npm run typecheck"
	},
	{
		kind: "ok",
		text: "> tsc -b"
	},
	{
		kind: "cmd",
		text: "$ npm run build"
	},
	{
		kind: "ok",
		text: "OK built in 826ms"
	},
	{
		kind: "cmd",
		text: "$ wisp run --panel"
	},
	{
		kind: "out",
		text: "[panel] C17 桥接就绪，等待宿主推送快照…"
	},
	{
		kind: "out",
		text: "[ball] Direct2D 表面已挂载，20 态转移表加载完成"
	}
];
var BROWSER_URL = "wisp.local/preview/生成页面";
function RightRail({ openRequest, tabs = TAB_DEFS, reviewFiles = REVIEW_FILES, reviewScopes = REVIEW_SCOPES, fileTree = TREE, workspaceRoot = WORKSPACE_ROOT, terminalLines = TERMINAL_LINES, browserUrl = BROWSER_URL }) {
	const [openTabs, setOpenTabs] = useState([]);
	const [activeTab, setActiveTab] = useState(null);
	const [closedTabs, setClosedTabs] = useState([]);
	const [menuOpen, setMenuOpen] = useState(false);
	const [menuQuery, setMenuQuery] = useState("");
	const handled = useRef(void 0);
	useEffect(() => {
		if (!openRequest || openRequest.nonce === handled.current) return;
		handled.current = openRequest.nonce;
		setOpenTabs((current) => current.includes(openRequest.tab) ? current : [...current, openRequest.tab]);
		setClosedTabs((current) => current.filter((k) => k !== openRequest.tab));
		setActiveTab(openRequest.tab);
	}, [openRequest]);
	function closeTab(key) {
		setOpenTabs((current) => {
			const next = current.filter((k) => k !== key);
			if (activeTab === key) setActiveTab(next[next.length - 1] ?? null);
			return next;
		});
		setClosedTabs((current) => [key, ...current.filter((k) => k !== key)].slice(0, 3));
		setMenuOpen(false);
	}
	function reopenTab(key) {
		setOpenTabs((current) => current.includes(key) ? current : [...current, key]);
		setClosedTabs((current) => current.filter((k) => k !== key));
		setActiveTab(key);
		setMenuOpen(false);
	}
	const active = activeTab && openTabs.includes(activeTab) ? activeTab : null;
	const menuMatches = (def) => def.label.toLowerCase().includes(menuQuery.toLowerCase());
	return /* @__PURE__ */ jsxs("aside", {
		"aria-label": "上下文栏",
		className: "flex w-[340px] shrink-0 flex-col overflow-hidden border-l border-line bg-surface text-ink",
		children: [/* @__PURE__ */ jsxs("div", {
			className: "relative flex shrink-0 items-stretch gap-0.5 border-b border-line bg-inset p-1",
			children: [
				/* @__PURE__ */ jsx(GlideHighlight, {
					activeKey: active,
					layoutKey: openTabs.join("|"),
					radius: "6px"
				}),
				/* @__PURE__ */ jsx("button", {
					"aria-expanded": menuOpen,
					"aria-label": "标签页列表",
					className: cn("relative z-10 flex w-7 shrink-0 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:text-ink", menuOpen && "text-ink"),
					onClick: () => setMenuOpen((v) => !v),
					title: "标签页列表",
					type: "button",
					children: /* @__PURE__ */ jsx("svg", {
						fill: "none",
						height: "12",
						stroke: "currentColor",
						strokeLinecap: "round",
						strokeLinejoin: "round",
						strokeWidth: "2",
						viewBox: "0 0 24 24",
						width: "12",
						"aria-hidden": "true",
						children: /* @__PURE__ */ jsx("path", { d: "M7 15l5 5 5-5M7 9l5-5 5 5" })
					})
				}),
				openTabs.map((key) => {
					const def = tabs.find((t) => t.key === key);
					if (!def) return null;
					const on = key === active;
					return /* @__PURE__ */ jsxs("button", {
						"aria-selected": on,
						className: cn("group relative z-10 flex min-w-0 flex-1 items-center justify-center gap-1 rounded-control px-2 py-1.5 text-[12px] font-medium transition-colors duration-150", on ? "text-ink" : "text-ink-2 hover:text-ink"),
						"data-glide-key": key,
						onClick: () => setActiveTab(key),
						role: "tab",
						type: "button",
						children: [/* @__PURE__ */ jsx("span", {
							className: "truncate",
							children: def.label
						}), /* @__PURE__ */ jsx("span", {
							"aria-label": `关闭 ${def.label}`,
							className: cn("flex size-3.5 shrink-0 items-center justify-center rounded-[4px] text-ink-3 opacity-0 transition-[opacity,background-color,color] duration-100 hover:bg-line/70 hover:text-ink group-hover:opacity-100", on && "opacity-100"),
							onClick: (event) => {
								event.stopPropagation();
								closeTab(key);
							},
							role: "button",
							tabIndex: -1,
							children: /* @__PURE__ */ jsx(X, {
								"aria-hidden": "true",
								size: 9,
								strokeWidth: 2.6
							})
						})]
					}, key);
				}),
				/* @__PURE__ */ jsx("button", {
					"aria-label": "打开更多",
					className: "relative z-10 flex w-7 shrink-0 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:text-ink",
					onClick: () => setMenuOpen((v) => !v),
					title: "打开更多",
					type: "button",
					children: /* @__PURE__ */ jsx(Plus, {
						"aria-hidden": "true",
						size: 13,
						strokeWidth: 2.2
					})
				}),
				menuOpen && /* @__PURE__ */ jsxs(GlideMenu, {
					className: "absolute left-1 right-1 top-full z-40 mt-1 rounded-card border border-line bg-surface p-1.5 shadow-overlay",
					highlightClassName: "inset-x-0 rounded-[6px] bg-hover",
					children: [
						/* @__PURE__ */ jsx("label", {
							className: "mb-1 flex h-8 items-center gap-2 rounded-control bg-inset px-2.5 shadow-hairline",
							children: /* @__PURE__ */ jsx("input", {
								onChange: (event) => setMenuQuery(event.target.value),
								placeholder: "搜索标签页…",
								value: menuQuery,
								className: "min-w-0 flex-1 bg-transparent text-[12.5px] text-ink outline-none placeholder:text-ink-3"
							})
						}),
						/* @__PURE__ */ jsx("p", {
							className: "px-1.5 pb-0.5 pt-1.5 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3",
							children: "可打开"
						}),
						tabs.filter((t) => !openTabs.includes(t.key) && menuMatches(t)).map((def) => /* @__PURE__ */ jsxs("button", {
							className: "flex w-full items-center gap-2 rounded-chip px-1.5 py-1.5 text-left",
							"data-menu-row": true,
							onClick: () => reopenTab(def.key),
							title: def.blurb,
							type: "button",
							children: [
								/* @__PURE__ */ jsx(IconFor, { def }),
								/* @__PURE__ */ jsx("span", {
									className: "min-w-0 flex-1 truncate text-[12px] text-ink-2",
									children: def.label
								}),
								def.keys ? /* @__PURE__ */ jsx("kbd", {
									className: "kbd shrink-0",
									children: def.keys
								}) : null
							]
						}, def.key)),
						openTabs.filter((k) => menuMatches(tabs.find((t) => t.key === k))).length > 0 && /* @__PURE__ */ jsxs(Fragment$1, { children: [/* @__PURE__ */ jsx("p", {
							className: "px-1.5 pb-0.5 pt-2 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3",
							children: "已打开"
						}), openTabs.filter((k) => menuMatches(tabs.find((t) => t.key === k))).map((key) => {
							const def = tabs.find((t) => t.key === key);
							return /* @__PURE__ */ jsxs("div", {
								className: "flex items-center gap-2 rounded-chip px-1.5 py-1.5",
								"data-menu-row": true,
								children: [/* @__PURE__ */ jsx("span", {
									className: "min-w-0 flex-1 truncate text-left text-[12px] text-ink",
									children: def.label
								}), /* @__PURE__ */ jsx("span", {
									"aria-label": `关闭 ${def.label}`,
									className: "flex size-4 shrink-0 items-center justify-center rounded-[4px] text-ink-3 hover:text-ink",
									onClick: () => closeTab(key),
									role: "button",
									tabIndex: -1,
									children: /* @__PURE__ */ jsx(X, {
										"aria-hidden": "true",
										size: 9,
										strokeWidth: 2.6
									})
								})]
							}, key);
						})] }),
						closedTabs.filter((k) => menuMatches(tabs.find((t) => t.key === k))).length > 0 && /* @__PURE__ */ jsx("p", {
							className: "px-1.5 pb-0.5 pt-2 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3",
							children: "最近关闭的标签页"
						}),
						closedTabs.filter((k) => menuMatches(tabs.find((t) => t.key === k))).map((key) => {
							const def = tabs.find((t) => t.key === key);
							return /* @__PURE__ */ jsxs("button", {
								className: "flex w-full items-center gap-2 rounded-chip px-1.5 py-1.5 text-left",
								"data-menu-row": true,
								onClick: () => reopenTab(key),
								type: "button",
								children: [/* @__PURE__ */ jsx("span", {
									className: "min-w-0 flex-1 truncate text-[12px] text-ink-2",
									children: def.label
								}), /* @__PURE__ */ jsx("span", {
									className: "shrink-0 text-[10.5px] text-ink-3",
									children: def.closedAgo
								})]
							}, key);
						})
					]
				})
			]
		}), openTabs.length === 0 ? /* @__PURE__ */ jsxs("div", {
			className: "flex min-h-0 flex-1 flex-col items-center justify-center gap-3 px-4 py-6",
			style: { backgroundImage: "radial-gradient(circle at 50% 55%, var(--accent-tint) 0%, transparent 62%)" },
			children: [
				/* @__PURE__ */ jsx("p", {
					className: "text-[13px] font-semibold text-ink",
					children: "打开面板"
				}),
				/* @__PURE__ */ jsx("p", {
					className: "mb-1 text-[11.5px] text-ink-3",
					children: "选择要在侧边面板中查看的内容。"
				}),
				/* @__PURE__ */ jsx(GlideMenu, {
					className: "flex w-full flex-col gap-2.5",
					highlightClassName: "inset-x-0 rounded-card bg-hover",
					children: tabs.map((def) => /* @__PURE__ */ jsxs("button", {
						className: "group relative flex w-full items-center gap-3 rounded-card border border-line bg-surface p-4 text-left shadow-card transition-colors duration-150 hover:border-accent/40",
						"data-menu-row": true,
						onClick: () => reopenTab(def.key),
						type: "button",
						children: [
							/* @__PURE__ */ jsx("span", {
								className: "flex size-8 shrink-0 items-center justify-center rounded-control bg-field text-ink-2",
								children: /* @__PURE__ */ jsx(IconFor, { def })
							}),
							/* @__PURE__ */ jsxs("span", {
								className: "min-w-0 flex-1",
								children: [/* @__PURE__ */ jsx("span", {
									className: "block truncate text-[12.5px] font-medium text-ink",
									children: def.label
								}), /* @__PURE__ */ jsx("span", {
									className: "mt-0.5 block truncate text-[11px] text-ink-3",
									children: def.blurb
								})]
							}),
							def.keys ? /* @__PURE__ */ jsx("kbd", {
								className: "kbd shrink-0",
								children: def.keys
							}) : null
						]
					}, def.key))
				})
			]
		}) : /* @__PURE__ */ jsxs(Fragment$1, { children: [
			active === "review" && /* @__PURE__ */ jsx(ReviewPane, {
				files: reviewFiles,
				scopes: reviewScopes
			}),
			active === "files" && /* @__PURE__ */ jsx(FilesPane, {
				fileTree,
				workspaceRoot
			}),
			active === "terminal" && /* @__PURE__ */ jsx(TerminalPane, { lines: terminalLines }),
			active === "browser" && /* @__PURE__ */ jsx(BrowserPane, { url: browserUrl })
		] })]
	});
}
/** 启动页卡片上的图标，和标签页用的是同一套映射，避免两处各画一套。 */
var RAIL_ICONS = {
	review: GitPullRequest,
	files: Folder,
	terminal: SquareTerminal,
	browser: Compass
};
function IconFor({ def }) {
	const Icon = RAIL_ICONS[def.key];
	return /* @__PURE__ */ jsx(Icon, {
		"aria-hidden": "true",
		size: 15,
		strokeWidth: 1.8
	});
}
function FilesPane({ fileTree, workspaceRoot }) {
	const [collapsed, setCollapsed] = useState(/* @__PURE__ */ new Set());
	const [open, setOpen] = useState(null);
	const visible = fileTree.filter((row, i) => {
		if (row.kind !== "dir") return true;
		for (let j = i - 1; j >= 0; j--) if (fileTree[j].kind === "dir" && fileTree[j].depth < row.depth) return !collapsed.has(fileTree[j].name);
		return true;
	});
	return /* @__PURE__ */ jsxs("section", {
		"aria-label": "工作区文件",
		className: "flex min-h-0 flex-1 flex-col overflow-hidden",
		children: [
			/* @__PURE__ */ jsxs("div", {
				className: "flex shrink-0 items-center gap-2 border-b border-line px-3 py-2",
				children: [/* @__PURE__ */ jsx("span", {
					className: "min-w-0 flex-1 truncate font-mono text-[11px] text-ink-2",
					title: workspaceRoot,
					children: workspaceRoot
				}), /* @__PURE__ */ jsxs("button", {
					className: "flex shrink-0 items-center gap-1 rounded-control px-1.5 py-1 text-[11.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink",
					title: "刷新",
					type: "button",
					children: [/* @__PURE__ */ jsx(RotateCw, {
						"aria-hidden": "true",
						size: 12,
						strokeWidth: 2
					}), "刷新"]
				})]
			}),
			visible.length === 0 ? /* @__PURE__ */ jsx("p", {
				className: "px-3 py-4 text-[11.5px] leading-relaxed text-ink-3",
				children: "还没有收到工作区文件。面板不许自己造数据，这里只说明它缺 C17 的读口与票 35 的推送。"
			}) : /* @__PURE__ */ jsx(GlideMenu, {
				className: "min-h-0 flex-1 overflow-y-auto py-1",
				highlightClassName: "inset-x-1.5 rounded-chip bg-hover",
				children: visible.map((row) => {
					const isDir = row.kind === "dir";
					const shut = isDir && collapsed.has(row.name);
					return /* @__PURE__ */ jsxs("button", {
						"aria-expanded": isDir ? !shut : void 0,
						className: cn("relative flex w-full items-center gap-1.5 rounded-chip py-1 pr-2 text-left transition-colors duration-100", open === row.name ? "text-ink" : "text-ink-2"),
						"data-menu-row": true,
						onClick: () => {
							setOpen(row.name);
							if (isDir) setCollapsed((current) => {
								const next = new Set(current);
								if (next.has(row.name)) next.delete(row.name);
								else next.add(row.name);
								return next;
							});
						},
						style: { paddingLeft: 8 + row.depth * 12 },
						type: "button",
						children: [
							open === row.name ? /* @__PURE__ */ jsx(RowMark, {}) : null,
							isDir ? /* @__PURE__ */ jsx(ChevronRight, {
								"aria-hidden": "true",
								className: cn("shrink-0 text-ink-3 transition-transform duration-150", !shut && "rotate-90"),
								size: 11,
								strokeWidth: 2.2
							}) : /* @__PURE__ */ jsx("span", {
								"aria-hidden": "true",
								className: "w-[11px] shrink-0"
							}),
							isDir ? /* @__PURE__ */ jsx(Folder, {
								"aria-hidden": "true",
								className: "shrink-0 text-ink-3",
								size: 12,
								strokeWidth: 1.8
							}) : /* @__PURE__ */ jsx(FileText, {
								"aria-hidden": "true",
								className: "shrink-0 text-ink-3",
								size: 12,
								strokeWidth: 1.8
							}),
							/* @__PURE__ */ jsx("span", {
								className: cn("min-w-0 flex-1 truncate text-[12px]", row.changed ? "font-medium text-ink" : isDir ? "text-ink-2" : "text-ink-2"),
								children: row.name
							}),
							row.locked ? /* @__PURE__ */ jsxs("span", {
								className: "shrink-0 text-[10px] text-warn",
								children: [row.locked, " 正在写"]
							}) : row.changed ? /* @__PURE__ */ jsx("span", {
								"aria-label": "本轮已改动",
								className: "size-1.5 shrink-0 rounded-full bg-green"
							}) : null
						]
					}, `${row.depth}:${row.name}`);
				})
			}),
			/* @__PURE__ */ jsx("p", {
				className: "shrink-0 border-t border-line px-3 py-2 text-[10.5px] leading-[1.6] text-ink-3",
				children: "一缕只能读写你授权的目录；工作区必须落在其中（C26 拒绝符号链接外指）。 标着「正在写」的行被工具占着路径锁，这时点开只会拿到一句拒绝。"
			})
		]
	});
}
function ReviewPane({ files, scopes }) {
	const [scopeOpen, setScopeOpen] = useState(false);
	const [scope, setScope] = useState(scopes[0] ?? "未暂存");
	const [openFile, setOpenFile] = useState(null);
	return /* @__PURE__ */ jsxs("section", {
		"aria-label": "审查",
		className: "flex min-h-0 flex-1 flex-col overflow-hidden",
		children: [/* @__PURE__ */ jsxs("div", {
			className: "relative flex shrink-0 items-center justify-between gap-2 border-b border-line px-3 py-2",
			children: [
				/* @__PURE__ */ jsxs("button", {
					"aria-expanded": scopeOpen,
					className: "flex items-center gap-1.5 rounded-control px-1.5 py-1 text-[12px] font-medium text-ink transition-colors duration-150 hover:bg-hover",
					onClick: () => setScopeOpen((v) => !v),
					type: "button",
					children: [scope, /* @__PURE__ */ jsx(ChevronDown, {
						"aria-hidden": "true",
						size: 11,
						strokeWidth: 2.2
					})]
				}),
				/* @__PURE__ */ jsxs("button", {
					className: "flex items-center gap-1 rounded-control px-1.5 py-1 text-[11.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink",
					title: "刷新",
					type: "button",
					children: [/* @__PURE__ */ jsx(RotateCw, {
						"aria-hidden": "true",
						size: 12,
						strokeWidth: 2
					}), "刷新"]
				}),
				scopeOpen && /* @__PURE__ */ jsx("div", {
					className: "absolute left-3 top-full z-30 mt-1 w-44 rounded-card border border-line bg-surface p-1 shadow-overlay",
					style: { animation: "pop-in 160ms var(--ease-out-strong) both" },
					children: /* @__PURE__ */ jsx(GlideMenu, {
						className: "flex flex-col",
						highlightClassName: "inset-x-0 rounded-chip bg-hover",
						children: scopes.map((s) => /* @__PURE__ */ jsxs("button", {
							className: "relative flex w-full items-center gap-2 rounded-chip px-2 py-1.5 text-left text-[12px] text-ink",
							"data-menu-row": true,
							onClick: () => {
								setScope(s);
								setScopeOpen(false);
							},
							type: "button",
							children: [
								s === scope ? /* @__PURE__ */ jsx(RowMark, {}) : null,
								/* @__PURE__ */ jsx("span", {
									className: cn("w-3.5 shrink-0", s === scope ? "text-accent-ink" : "invisible"),
									children: /* @__PURE__ */ jsx("svg", {
										fill: "none",
										height: "11",
										stroke: "currentColor",
										strokeLinecap: "round",
										strokeWidth: "2.4",
										viewBox: "0 0 24 24",
										width: "11",
										children: /* @__PURE__ */ jsx("path", { d: "M20 6L9 17l-5-5" })
									})
								}),
								s
							]
						}, s))
					})
				})
			]
		}), files.length === 0 ? /* @__PURE__ */ jsx("p", {
			className: "px-3 py-4 text-[11.5px] leading-relaxed text-ink-3",
			children: "还没有收到审查数据。面板不许自己造数据，这里只说明它缺 C17 的读口与票 35 的推送。"
		}) : /* @__PURE__ */ jsx(GlideMenu, {
			className: "min-h-0 flex-1 overflow-y-auto",
			highlightClassName: "inset-x-1.5 rounded-chip bg-hover",
			children: files.map((file) => {
				const expanded = openFile === file.name;
				return /* @__PURE__ */ jsxs("div", {
					className: "border-b border-line last:border-b-0",
					children: [/* @__PURE__ */ jsxs("button", {
						"aria-expanded": expanded,
						className: "relative flex w-full items-center gap-2 px-3 py-2 pl-3.5 text-left transition-colors duration-100",
						"data-menu-row": true,
						onClick: () => setOpenFile(expanded ? null : file.name),
						type: "button",
						children: [
							/* @__PURE__ */ jsx("span", {
								"aria-hidden": "true",
								className: cn("size-3.5 shrink-0 rounded-[4px]", file.tile)
							}),
							/* @__PURE__ */ jsx("span", {
								className: "shrink-0 font-mono text-[11.5px] font-medium text-ink",
								children: file.name
							}),
							/* @__PURE__ */ jsx("span", {
								className: "min-w-0 flex-1 truncate text-[10.5px] text-ink-3",
								children: file.dir
							}),
							/* @__PURE__ */ jsxs("span", {
								className: "shrink-0 font-mono text-[10.5px] tabular-nums text-green",
								children: ["+", file.additions]
							}),
							/* @__PURE__ */ jsxs("span", {
								className: "shrink-0 font-mono text-[10.5px] tabular-nums text-red",
								children: ["-", file.removals]
							}),
							/* @__PURE__ */ jsx(ChevronDown, {
								"aria-hidden": "true",
								className: cn("shrink-0 text-ink-3 transition-transform duration-150", expanded && "rotate-180"),
								size: 11,
								strokeWidth: 2.2
							})
						]
					}), expanded && file.diff ? /* @__PURE__ */ jsx("div", {
						className: "mx-3 mb-2 overflow-hidden rounded-chip border border-line font-mono text-[10.5px] leading-[1.7]",
						style: { animation: "fade-in 150ms ease both" },
						children: file.diff.map((line, i) => /* @__PURE__ */ jsxs("div", {
							className: cn("px-2", line.sign === "+" ? "bg-green-tint text-green" : "bg-red-tint text-red"),
							children: [
								line.sign,
								" ",
								line.text
							]
						}, i))
					}) : null]
				}, file.name);
			})
		})]
	});
}
function TerminalPane({ lines }) {
	return /* @__PURE__ */ jsxs("section", {
		"aria-label": "终端",
		className: "flex min-h-0 flex-1 flex-col overflow-hidden",
		children: [
			/* @__PURE__ */ jsxs("div", {
				className: "flex items-center gap-1.5 border-b border-line px-3 py-2 text-[11.5px] text-ink-2",
				children: [
					/* @__PURE__ */ jsx("span", { className: "size-2 rounded-full bg-red" }),
					/* @__PURE__ */ jsx("span", { className: "size-2 rounded-full bg-orange" }),
					/* @__PURE__ */ jsx("span", { className: "size-2 rounded-full bg-green" }),
					/* @__PURE__ */ jsx("span", {
						className: "ml-1.5 font-mono",
						children: "wisp — 终端"
					})
				]
			}),
			/* @__PURE__ */ jsxs("div", {
				className: "min-h-0 flex-1 overflow-y-auto p-3 font-mono text-[11.5px] leading-[1.9]",
				style: {
					background: "var(--tooltip-bg)",
					color: "var(--tooltip-fg)"
				},
				children: [lines.length === 0 ? /* @__PURE__ */ jsx("p", {
					className: "text-[11.5px] leading-relaxed text-tooltip-muted",
					children: "还没有收到终端会话。面板不许自己造数据，这里只说明它缺 C17 的读口与票 35 的推送。"
				}) : lines.map((line, i) => /* @__PURE__ */ jsx("div", {
					className: cn("whitespace-pre-wrap", line.kind === "cmd" && "text-tooltip-fg", line.kind === "ok" && "text-tooltip-muted", line.kind === "out" && "text-tooltip-muted"),
					children: line.text
				}, i)), /* @__PURE__ */ jsxs("div", {
					className: "mt-1 flex items-center gap-1",
					children: [/* @__PURE__ */ jsx("span", {
						className: "text-green",
						children: "$"
					}), /* @__PURE__ */ jsx("span", {
						className: "inline-block h-3.5 w-[7px] bg-tooltip-fg",
						style: { animation: "caret-blink 1s step-end infinite" }
					})]
				})]
			}),
			/* @__PURE__ */ jsx("p", {
				className: "border-t border-line px-3 py-2 text-[10.5px] text-ink-3",
				children: "只读演示。真实终端会话由宿主进程托管（外部插件或内置集成待定案）。"
			})
		]
	});
}
function BrowserPane({ url }) {
	return /* @__PURE__ */ jsxs("section", {
		"aria-label": "浏览器",
		className: "flex min-h-0 flex-1 flex-col overflow-hidden",
		children: [/* @__PURE__ */ jsxs("div", {
			className: "flex shrink-0 items-center gap-1 border-b border-line px-2 py-1.5",
			children: [
				/* @__PURE__ */ jsx(ChevronLeft, {
					"aria-hidden": "true",
					className: "shrink-0 text-ink-3",
					size: 14,
					strokeWidth: 2
				}),
				/* @__PURE__ */ jsx(ChevronRight, {
					"aria-hidden": "true",
					className: "shrink-0 text-ink-3",
					size: 14,
					strokeWidth: 2
				}),
				/* @__PURE__ */ jsx(RotateCw, {
					"aria-hidden": "true",
					className: "shrink-0 text-ink-3",
					size: 12,
					strokeWidth: 2
				}),
				/* @__PURE__ */ jsx("span", {
					className: "ml-1 min-w-0 flex-1 truncate rounded-control bg-inset px-2 py-1 font-mono text-[10.5px] text-ink-2 shadow-hairline",
					children: url
				})
			]
		}), /* @__PURE__ */ jsx("div", {
			className: "min-h-0 flex-1 overflow-y-auto bg-page",
			children: /* @__PURE__ */ jsxs("div", {
				className: "mx-auto max-w-[420px] px-5 py-8",
				children: [/* @__PURE__ */ jsxs("div", {
					className: "rounded-card border border-line bg-surface p-4 shadow-card",
					children: [
						/* @__PURE__ */ jsx("p", {
							className: "text-[15px] font-semibold text-ink",
							children: "桌面截图归档报告"
						}),
						/* @__PURE__ */ jsx("p", {
							className: "mt-1 text-[11.5px] leading-relaxed text-ink-2",
							children: "本周共归档 47 张截图到 6 个月份目录，重复 12 张已移入待确认区。"
						}),
						/* @__PURE__ */ jsx("div", {
							className: "mt-3 flex flex-col gap-1.5",
							children: [
								"2026-09 · 18 张",
								"2026-08 · 11 张",
								"2026-07 · 9 张"
							].map((row) => /* @__PURE__ */ jsxs("div", {
								className: "flex items-center gap-2 rounded-control bg-inset px-2.5 py-1.5 text-[11.5px] text-ink-2",
								children: [/* @__PURE__ */ jsx("span", {
									className: "min-w-0 flex-1",
									children: row
								}), /* @__PURE__ */ jsx("span", {
									className: "h-1 w-16 overflow-hidden rounded-full bg-field",
									children: /* @__PURE__ */ jsx("span", {
										className: "block h-full rounded-full bg-accent",
										style: { width: "72%" }
									})
								})]
							}, row))
						})
					]
				}), /* @__PURE__ */ jsx("p", {
					className: "mt-4 text-center text-[10.5px] leading-relaxed text-ink-3",
					children: "演示页面 · 浏览器内核未接入（外部插件或内置集成待定案）"
				})]
			})
		})]
	});
}
//#endregion
//#region src/components/ai-native/chip.tsx
/** Monospace token chip — for code values like `updated_at`. */
function Chip({ children, tone = "neutral", className = "" }) {
	return /* @__PURE__ */ jsx("code", {
		className: `inline rounded-md px-1.5 py-0.5 font-mono text-[12px]
        leading-none align-[-1px] ${{
			neutral: "bg-inset text-ink-2",
			accent: "bg-accent-tint text-accent-ink",
			orange: "bg-orange-tint text-orange"
		}[tone]} ${className}`,
		children
	});
}
//#endregion
//#region src/components/ai-native/progress-ring.tsx
/** Small progress ring with content in the center — the task-badge from the refs. */
function ProgressRing({ progress, tone = "orange", children, size = 28 }) {
	const stroke = 2;
	const r = (size - stroke) / 2;
	const c = 2 * Math.PI * r;
	return /* @__PURE__ */ jsxs("span", {
		className: "relative inline-flex items-center justify-center",
		style: {
			width: size,
			height: size
		},
		children: [/* @__PURE__ */ jsxs("svg", {
			width: size,
			height: size,
			className: "-rotate-90 absolute inset-0",
			children: [/* @__PURE__ */ jsx("circle", {
				cx: size / 2,
				cy: size / 2,
				r,
				fill: "none",
				stroke: "var(--line)",
				strokeWidth: stroke
			}), /* @__PURE__ */ jsx("circle", {
				cx: size / 2,
				cy: size / 2,
				r,
				fill: "none",
				stroke: {
					orange: "var(--orange)",
					green: "var(--green)",
					red: "var(--red)",
					accent: "var(--accent)"
				}[tone],
				strokeWidth: stroke,
				strokeLinecap: "round",
				strokeDasharray: c,
				strokeDashoffset: c * (1 - Math.min(1, Math.max(0, progress))),
				style: { transition: "stroke-dashoffset 400ms cubic-bezier(0.23,1,0.32,1)" }
			})]
		}), /* @__PURE__ */ jsx("span", {
			className: "relative text-[12px] font-semibold tabular-nums",
			children
		})]
	});
}
//#endregion
//#region src/components/ai-native/segmented-control.tsx
/** Segmented control — equal-width segments, sliding thumb. */
function SegmentedControl({ options, value, onChange, className = "" }) {
	const index = options.indexOf(value);
	return /* @__PURE__ */ jsxs("div", {
		className: `relative inline-grid h-8 select-none rounded-full bg-line/60 p-0.5 ${className}`,
		style: { gridTemplateColumns: `repeat(${options.length}, 1fr)` },
		role: "tablist",
		children: [/* @__PURE__ */ jsx("span", {
			"aria-hidden": true,
			className: "absolute inset-y-0.5 rounded-full bg-surface shadow-hairline\n          transition-transform duration-200",
			style: {
				width: `calc((100% - 4px) / ${options.length})`,
				left: 2,
				transform: `translateX(${index * 100}%)`,
				transitionTimingFunction: "cubic-bezier(0.23, 1, 0.32, 1)"
			}
		}), options.map((opt) => /* @__PURE__ */ jsx("button", {
			role: "tab",
			"aria-selected": opt === value,
			onClick: () => onChange(opt),
			className: `relative z-10 rounded-full px-3 text-[13px] font-medium
            transition-colors duration-150
            ${opt === value ? "text-ink" : "text-ink-3 hover:text-ink-2"}`,
			children: opt
		}, opt))]
	});
}
//#endregion
//#region src/lib/motion.ts
/** What the settings row offers, in the order it offers them. */
var MOTION_TIERS = [
	"全开",
	"适中",
	"关闭"
];
/** The value each tier writes onto <html data-motion>. "full" is written as
the absence of nothing - it IS the value, because "full" has to be the
default and a document that forgot to set an attribute should still be
allowed to animate. */
var TIER_ATTR = {
	全开: "full",
	适中: "lite",
	关闭: "off"
};
var MOTION_ATTR = "data-motion";
/**
* Apply a tier to the current document. Safe to call before there is a document
* (the SSR render harnesses), where it does nothing.
*/
function applyMotionTier(tier) {
	if (typeof document === "undefined") return;
	document.documentElement.setAttribute(MOTION_ATTR, TIER_ATTR[tier]);
}
/** The tier this document is currently in, read back from the attribute rather
than from a mirror of it - so a tier set anywhere (including by a test) is
the tier the UI reports. */
function readMotionTier() {
	if (typeof document === "undefined") return "全开";
	const v = document.documentElement.getAttribute(MOTION_ATTR);
	if (v === "off") return "关闭";
	if (v === "lite") return "适中";
	return "全开";
}
//#endregion
//#region src/components/ai-native/switch.tsx
function Switch({ checked, onChange, label }) {
	return /* @__PURE__ */ jsx("button", {
		role: "switch",
		"aria-checked": checked,
		"aria-label": label,
		onClick: () => onChange(!checked),
		className: `relative h-6 w-10 shrink-0 rounded-full transition-colors duration-200
        ${checked ? "bg-ink" : "bg-line-strong"}`,
		children: /* @__PURE__ */ jsx("span", {
			className: "absolute top-0.5 left-0.5 size-5 rounded-full bg-white\n          shadow-btn transition-transform duration-200",
			style: {
				transform: checked ? "translateX(16px)" : "translateX(0)",
				transitionTimingFunction: "cubic-bezier(0.23, 1, 0.32, 1)"
			}
		})
	});
}
//#endregion
//#region src/components/config-screen.tsx
/** The knob the frost reads. Kept as one name so nothing else can drift. */
var ALPHA_VAR = "--panel-alpha";
/**
* The default comes from the generated token table, not from a number typed
* here - that is the whole point of having a table. Returns null when there
* is no document to ask (the render harnesses build this on the server) or
* when the property is missing, and a null default makes the row say so
* rather than guessing.
*/
function readPanelAlpha() {
	if (typeof document === "undefined") return null;
	const raw = getComputedStyle(document.documentElement).getPropertyValue(ALPHA_VAR).trim();
	const n = Number.parseFloat(raw);
	return Number.isFinite(n) ? Math.round(n) : null;
}
/** The appearance rows, in the demo owner's original order. */
var APPEARANCE_ROWS = [
	{
		group: "app",
		name: "theme",
		desc: "界面主题",
		blocked: "主题由原生宿主决定，面板只跟随"
	},
	{
		group: "panel",
		name: "opacity",
		desc: "面板背景不透明度",
		blocked: "",
		control: "alpha"
	},
	{
		group: "panel",
		name: "animations",
		desc: "界面动画（关掉后全部静帧，性能最省）",
		blocked: "",
		control: "motion"
	},
	{
		group: "panel",
		name: "font_size",
		desc: "面板正文字号（px）",
		blocked: "",
		control: "panelFontSize"
	},
	{
		group: "ball",
		name: "size",
		desc: "悬浮球直径（44-160）",
		blocked: "",
		control: "ballSize"
	},
	{
		group: "ball",
		name: "opacity_idle",
		desc: "空闲时悬浮球不透明度",
		blocked: "",
		control: "ballOpacity"
	},
	{
		group: "ball",
		name: "click_through",
		desc: "空闲时鼠标穿透悬浮球",
		blocked: "",
		control: "clickThrough"
	},
	{
		group: "ball",
		name: "style",
		desc: "悬浮球外观风格（玻璃 / Jarvis 粒子）",
		blocked: "",
		control: "ballStyle"
	},
	{
		group: "panel",
		name: "width",
		desc: "面板宽度（px）",
		blocked: "",
		control: "panelWidth"
	}
];
/** The effect badge: the two page-local rows say 即时 (they restyle this
document and nothing else), the config-write row says 能改 (its effect
tier is the receipt's to state - now when the ball is alive, restart when
it is not), dormant rows say 未接线. 即时 and 未接线 are this screen's
own needle words; 能改 joins them only because calling a restart-tier
write 即时 would be the lie the receipt exists to correct. */
function EffectBadge({ live, word }) {
	return /* @__PURE__ */ jsx("span", {
		className: cn("inline-flex h-[18px] shrink-0 items-center rounded-full px-1.5 text-[10px] font-medium leading-none", live ? "bg-green-tint text-green" : "bg-field text-ink-2"),
		children: live ? word ?? "即时" : "未接线"
	});
}
/**
* One px number field for the two panel rows (the appearance wave). Local
* draft state ONLY while the field is being edited - the shown value is what
* the host read, and an accepted write re-reads the view instead of patching
* a copy, so the field snaps back to the snapshot on the next push. The
* commit happens on blur / Enter (a keystroke-per-write number field would
* fire a config.set per digit), and an unparseable or non-positive draft is
* simply not sent - the field shows the snapshot's value again, it does not
* guess. The positive floor is the write route's own rule; the schema freezes
* no band for either key, so none is invented here.
*/
function NumberWriteField({ ariaLabel, value, onCommit, busy }) {
	const [editing, setEditing] = useState(false);
	const [draft, setDraft] = useState("");
	const commit = () => {
		const parsed = Number.parseInt(draft, 10);
		if (Number.isFinite(parsed) && parsed > 0 && parsed !== value) onCommit(parsed);
		setEditing(false);
	};
	return /* @__PURE__ */ jsx("input", {
		"aria-label": ariaLabel,
		className: "h-[30px] w-20 rounded-control bg-inset px-2 text-right font-mono text-[12px] text-ink shadow-hairline outline-none",
		disabled: busy,
		inputMode: "numeric",
		onBlur: () => {
			if (editing) commit();
		},
		onChange: (e) => {
			setEditing(true);
			setDraft(e.currentTarget.value);
		},
		onFocus: () => {
			setEditing(true);
			setDraft(String(value));
		},
		onKeyDown: (e) => {
			if (e.key === "Enter") e.currentTarget.blur();
		},
		type: "number",
		value: editing ? draft : String(value)
	});
}
function ConfigScreen({ heading = "设置", lede = "外观一节九行键。除了主题那一行由原生宿主决定之外，其余八行都能动：悬浮球风格、大小、空闲透明度、穿透，面板不透明度、宽度、字号和动画。", ballStyle = null, ballSize = null, ballOpacityIdle = null, ballClickThrough = null, panelWidth = null, panelFontSize = null }) {
	const [alpha, setAlpha] = useState(() => readPanelAlpha());
	const [motion, setMotion] = useState(() => readMotionTier());
	const [writeFeedback, setWriteFeedback] = useState(null);
	const [writeBusy, setWriteBusy] = useState(false);
	function moveAlpha(next) {
		document.documentElement.style.setProperty(ALPHA_VAR, `${next}%`);
		setAlpha(next);
	}
	function moveMotion(next) {
		applyMotionTier(next);
		setMotion(next);
	}
	/** The ball.style write: one config.set with the field's value, answered as
	the receipt sentence. Local validation first - an empty pick cannot be
	sent, because an empty value is refused by name on the Go side and the
	page can say so before anything leaves. */
	function moveBallStyle(next) {
		if (writeBusy || next === "") return;
		setWriteBusy(true);
		performConfigWrite({
			field: "ball_style",
			value: next
		}).then((fb) => {
			setWriteFeedback(fb);
			setWriteBusy(false);
		});
	}
	/** The ball.size write (tickets 413/422): same shape as moveBallStyle - one
	config.set with the field's value, answered as the receipt sentence.
	The size leaves as its decimal string; the 44..160 band is Go's own
	(the slider offers only the snapshot's min/max), so the page sends and
	lets the write route be the judge. */
	function moveBallSize(next) {
		if (writeBusy) return;
		setWriteBusy(true);
		performConfigWrite({
			field: "ball_size",
			value: String(next)
		}).then((fb) => {
			setWriteFeedback(fb);
			setWriteBusy(false);
		});
	}
	/** The ball.opacity_idle write (the appearance wave): the slider moves in
	percent so the row reads like the alpha row above it, but the value
	leaves as the 0..1 decimal the schema's own band validates - percent
	here is display arithmetic, never a second unit of truth. The band the
	slider offers is the snapshot's min/max, the setter's own range. */
	function moveBallOpacity(nextPercent) {
		if (writeBusy) return;
		setWriteBusy(true);
		performConfigWrite({
			field: "ball_opacity_idle",
			value: String(nextPercent / 100)
		}).then((fb) => {
			setWriteFeedback(fb);
			setWriteBusy(false);
		});
	}
	/** The ball.click_through write (the appearance wave): the Switch hands the
	next boolean over and it leaves as the plain value half the write route
	takes ("true" / "false"). The receipt states the tier - now when this
	process hosts the ball, restart when it does not. */
	function moveClickThrough(next) {
		if (writeBusy) return;
		setWriteBusy(true);
		performConfigWrite({
			field: "ball_click_through",
			value: next ? "true" : "false"
		}).then((fb) => {
			setWriteFeedback(fb);
			setWriteBusy(false);
		});
	}
	/** The panel.width write (the appearance wave): the px value leaves as its
	string and the write route judges it. No band is offered here - the
	schema freezes none, and the field's own positive-number floor is the
	write route's rule (see NumberWriteField). */
	function movePanelWidth(next) {
		if (writeBusy) return;
		setWriteBusy(true);
		performConfigWrite({
			field: "panel_width",
			value: String(next)
		}).then((fb) => {
			setWriteFeedback(fb);
			setWriteBusy(false);
		});
	}
	/** The panel.font_size write (the appearance wave): same shape as
	movePanelWidth, on the panel_font_size field. */
	function movePanelFontSize(next) {
		if (writeBusy) return;
		setWriteBusy(true);
		performConfigWrite({
			field: "panel_font_size",
			value: String(next)
		}).then((fb) => {
			setWriteFeedback(fb);
			setWriteBusy(false);
		});
	}
	return /* @__PURE__ */ jsxs("section", {
		"aria-label": "设置",
		className: "flex min-w-0 flex-col gap-1",
		children: [
			/* @__PURE__ */ jsx("h2", {
				className: "text-[15px] font-semibold text-ink",
				children: heading
			}),
			/* @__PURE__ */ jsx("p", {
				className: "text-[12.5px] text-ink-2",
				children: lede
			}),
			/* @__PURE__ */ jsx("div", {
				className: "mt-1 flex flex-col",
				children: APPEARANCE_ROWS.map((row) => {
					const live = row.blocked === "";
					return /* @__PURE__ */ jsxs("div", {
						className: "flex items-center justify-between gap-4 border-b border-line py-3 last:border-b-0",
						children: [/* @__PURE__ */ jsxs("div", {
							className: "min-w-0",
							children: [/* @__PURE__ */ jsxs(Chip, { children: [
								row.group,
								".",
								row.name
							] }), /* @__PURE__ */ jsx("div", {
								className: "mt-0.5 text-xs text-ink-2",
								children: row.desc
							})]
						}), /* @__PURE__ */ jsxs("div", {
							className: "flex shrink-0 items-center gap-2",
							children: [live && row.control === "motion" ? /* @__PURE__ */ jsx(SegmentedControl, {
								onChange: moveMotion,
								options: MOTION_TIERS,
								value: motion
							}) : live && row.control === "alpha" ? alpha === null ? /* @__PURE__ */ jsxs("span", {
								className: "text-xs text-ink-3",
								children: ["读不到 ", ALPHA_VAR]
							}) : /* @__PURE__ */ jsxs(Fragment$1, { children: [
								/* @__PURE__ */ jsx("input", {
									"aria-label": "面板背景不透明度",
									className: "w-28 accent-primary",
									max: 100,
									min: 30,
									onChange: (e) => moveAlpha(Number(e.currentTarget.value)),
									step: 1,
									type: "range",
									value: alpha
								}),
								/* @__PURE__ */ jsx(ProgressRing, {
									progress: alpha / 100,
									tone: "accent",
									children: alpha
								}),
								/* @__PURE__ */ jsx("span", {
									className: "text-[11px] text-ink-3",
									children: "%"
								})
							] }) : live && row.control === "ballSize" ? ballSize === null || ballSize === void 0 ? /* @__PURE__ */ jsx("span", {
								className: "text-xs text-ink-3",
								children: "等快照"
							}) : /* @__PURE__ */ jsxs(Fragment$1, { children: [/* @__PURE__ */ jsx("input", {
								"aria-label": "悬浮球大小",
								className: "w-28 accent-primary",
								max: ballSize.max,
								min: ballSize.min,
								onChange: (e) => moveBallSize(Number(e.currentTarget.value)),
								step: 1,
								type: "range",
								value: ballSize.size
							}), /* @__PURE__ */ jsxs("span", {
								className: "text-[11px] text-ink-3",
								children: [ballSize.size, "px"]
							})] }) : live && row.control === "ballOpacity" ? ballOpacityIdle === null || ballOpacityIdle === void 0 ? /* @__PURE__ */ jsx("span", {
								className: "text-xs text-ink-3",
								children: "等快照"
							}) : /* @__PURE__ */ jsxs(Fragment$1, { children: [/* @__PURE__ */ jsx("input", {
								"aria-label": "空闲球不透明度",
								className: "w-28 accent-primary",
								max: Math.round(ballOpacityIdle.max * 100),
								min: Math.round(ballOpacityIdle.min * 100),
								onChange: (e) => moveBallOpacity(Number(e.currentTarget.value)),
								step: 1,
								type: "range",
								value: Math.round(ballOpacityIdle.opacity * 100)
							}), /* @__PURE__ */ jsxs("span", {
								className: "text-[11px] text-ink-3",
								children: [Math.round(ballOpacityIdle.opacity * 100), "%"]
							})] }) : live && row.control === "clickThrough" ? ballClickThrough === null || ballClickThrough === void 0 ? /* @__PURE__ */ jsx("span", {
								className: "text-xs text-ink-3",
								children: "等快照"
							}) : /* @__PURE__ */ jsx(Switch, {
								checked: ballClickThrough.enabled,
								label: "空闲球鼠标穿透",
								onChange: moveClickThrough
							}) : live && row.control === "ballStyle" ? /* @__PURE__ */ jsx(BallStylePicker, {
								section: ballStyle,
								onSelect: moveBallStyle,
								className: "w-[180px] shrink-0"
							}) : live && row.control === "panelWidth" ? panelWidth === null || panelWidth === void 0 ? /* @__PURE__ */ jsx("span", {
								className: "text-xs text-ink-3",
								children: "等快照"
							}) : panelWidth.width <= 0 ? /* @__PURE__ */ jsx("span", {
								className: "text-xs text-ink-3",
								children: "快照里没有读到数值"
							}) : /* @__PURE__ */ jsx(NumberWriteField, {
								ariaLabel: "面板宽度",
								busy: writeBusy,
								onCommit: movePanelWidth,
								value: panelWidth.width
							}) : live && row.control === "panelFontSize" ? panelFontSize === null || panelFontSize === void 0 ? /* @__PURE__ */ jsx("span", {
								className: "text-xs text-ink-3",
								children: "等快照"
							}) : panelFontSize.fontSize <= 0 ? /* @__PURE__ */ jsx("span", {
								className: "text-xs text-ink-3",
								children: "快照里没有读到数值"
							}) : /* @__PURE__ */ jsx(NumberWriteField, {
								ariaLabel: "面板字号",
								busy: writeBusy,
								onCommit: movePanelFontSize,
								value: panelFontSize.fontSize
							}) : live ? /* @__PURE__ */ jsx("span", {
								className: "text-xs text-ink-3",
								children: "此行缺控件定义"
							}) : /* @__PURE__ */ jsx("span", {
								className: "text-xs text-ink-3",
								children: row.blocked
							}), /* @__PURE__ */ jsx(EffectBadge, {
								live,
								word: row.control === "ballStyle" || row.control === "ballSize" ? "能改" : row.control === "clickThrough" ? "立即/重启" : row.control === "ballOpacity" || row.control === "panelWidth" || row.control === "panelFontSize" ? "关窗/重启生效" : void 0
							})]
						})]
					}, `${row.group}-${row.name}`);
				})
			}),
			writeFeedback ? /* @__PURE__ */ jsx("div", {
				className: "mt-1 flex flex-col gap-1",
				children: /* @__PURE__ */ jsx(FeedbackLine, { feedback: writeFeedback })
			}) : null,
			/* @__PURE__ */ jsx("p", {
				className: "mt-3 text-[11px] leading-[1.6] text-ink-3",
				children: "「即时」那两行改的是这个页面自己的一层样式：关掉面板就回到 token 表里的默认值。 悬浮球风格、大小、空闲透明度、穿透和面板宽度、字号那六行不同：它们走设置写入通路真正落盘， 收据句会写明生效时机——穿透是有球当场就生效、没球要重启，宽度是关掉面板再打开就跟着变， 其余按收据照实说。"
			}),
			/* @__PURE__ */ jsxs("div", {
				className: "mt-2 flex flex-col gap-1 text-[11px] leading-[1.6] text-ink-3",
				children: [/* @__PURE__ */ jsx("span", { children: "其余各类在左边那一列上，各自写着它在等哪个数据源。" }), /* @__PURE__ */ jsx("span", { children: "设置项里没有一格能改变权限档位或工作区，那是面板的边界（P9）。" })]
			})
		]
	});
}
//#endregion
//#region src/components/config-read-state.tsx
function ConfigReadGate({ children }) {
	const { state, reload } = useConfigView();
	if (state.phase === "loading") return /* @__PURE__ */ jsx("p", {
		className: "mt-2 text-[12px] text-ink-3",
		children: "正在读取配置..."
	});
	if (state.phase === "no-host") return /* @__PURE__ */ jsx("div", {
		className: "mt-2 rounded-card border border-line bg-surface px-3.5 py-3",
		children: /* @__PURE__ */ jsx("p", {
			className: "text-[12px] leading-[1.7] text-ink-2",
			children: "暂未连接到本机服务，读不到配置。"
		})
	});
	if (state.phase === "unreadable" || state.phase === "payload-error") return /* @__PURE__ */ jsxs("div", {
		className: "mt-2 flex flex-col gap-2 rounded-card border border-line bg-surface px-3.5 py-3",
		children: [/* @__PURE__ */ jsx("p", {
			className: "text-[12px] leading-[1.7] text-ink-2",
			children: "暂时读不到配置，请稍后重试。"
		}), /* @__PURE__ */ jsx(Button, {
			className: "self-start",
			onClick: reload,
			size: "xs",
			variant: "secondary",
			children: "重试"
		})]
	});
	return /* @__PURE__ */ jsx(Fragment$1, { children: children(state.view, state.prose, reload) });
}
//#endregion
//#region src/components/model-settings.tsx
/** 票 313 AC#1（机主原话「可以把 stepfun 作为我们默认的模型来使用」；A877
裁定③＝`configuredKnown && !configured` 时默认高亮）：机主点名的默认
供应商。它只在两道门都过时才生效——用户还没配过任何东西
（`!view.configured`），且这个名字真的出现在 Go 侧送来的预设清单里
（`view.presets`）。它只是添加表单的初始预选：不点「添加」不写任何
东西，用户一旦有配置就永远以用户的为准，这里连预选都不再给。 */
var OWNER_DEFAULT_PROVIDER = "stepfun";
function ModelSettings({ heading = "模型", lede = "添加供应商、保存密钥，并选择默认聊天模型。" }) {
	return /* @__PURE__ */ jsxs("section", {
		"aria-label": "模型",
		className: "flex min-w-0 flex-col gap-1",
		children: [
			/* @__PURE__ */ jsx("h2", {
				className: "text-[15px] font-semibold text-ink",
				children: heading
			}),
			/* @__PURE__ */ jsx("p", {
				className: "text-[12.5px] text-ink-2",
				children: lede
			}),
			/* @__PURE__ */ jsx(ConfigReadGate, { children: (view, _prose, reload) => /* @__PURE__ */ jsx(ModelSettingsView, {
				onReload: reload,
				performCredential: performCredentialWrite,
				performFetchModels: performModelsFetch,
				performWrite: performConfigWrite,
				view
			}) })
		]
	});
}
function ModelSettingsView({ view, performWrite, performCredential, performFetchModels, onReload }) {
	const status = configStatus(view);
	const [busy, setBusy] = useState(null);
	const [providersFeedback, setProvidersFeedback] = useState(null);
	const [chatFeedback, setChatFeedback] = useState(null);
	const [addOpen, setAddOpen] = useState(false);
	const [editing, setEditing] = useState(null);
	const [chatProvider, setChatProvider] = useState(view.chatProvider);
	const [chatModel, setChatModel] = useState(view.chatModel);
	async function run(key, fn) {
		setBusy(key);
		const feedback = await fn();
		setBusy(null);
		onReload();
		return feedback;
	}
	return /* @__PURE__ */ jsxs("div", {
		className: "mt-2 flex min-w-0 flex-col gap-3",
		children: [/* @__PURE__ */ jsxs("div", {
			className: "flex items-center gap-2",
			children: [
				/* @__PURE__ */ jsx(StateTag, {
					text: status.label,
					tone: status.tone
				}),
				/* @__PURE__ */ jsx("span", { className: "min-w-0 flex-1" }),
				/* @__PURE__ */ jsxs(Button, {
					"aria-label": "刷新",
					onClick: onReload,
					size: "xs",
					variant: "quiet",
					children: [/* @__PURE__ */ jsx(RefreshCw, {
						"aria-hidden": "true",
						size: 11,
						strokeWidth: 2
					}), "刷新"]
				})
			]
		}), !view.readable ? /* @__PURE__ */ jsx(SettingsCard, {
			title: "配置",
			children: /* @__PURE__ */ jsx("p", {
				className: "text-[12px] leading-[1.7] text-ink-2",
				children: "暂时读不到配置，请稍后刷新。"
			})
		}) : /* @__PURE__ */ jsxs(Fragment$1, { children: [/* @__PURE__ */ jsx(ProvidersCard, {
			busy,
			editing,
			feedback: providersFeedback,
			onEditing: setEditing,
			onFeedback: setProvidersFeedback,
			performCredential,
			performFetchModels,
			performWrite,
			run,
			setAddOpen,
			addOpen,
			view
		}), /* @__PURE__ */ jsx(ChatModelCard, {
			busy,
			chatModel,
			chatProvider,
			feedback: chatFeedback,
			onFeedback: setChatFeedback,
			onModel: setChatModel,
			onProvider: (next) => {
				setChatProvider(next);
				setChatModel("");
			},
			performWrite,
			run,
			view
		})] })]
	});
}
function ProvidersCard({ view, busy, addOpen, setAddOpen, editing, onEditing, feedback, onFeedback, performWrite, performCredential, performFetchModels, run }) {
	const rows = view.providers;
	const [autoFetch, setAutoFetch] = useState(null);
	const availablePresets = view.presetsKnown && view.presets !== null ? view.presets.filter((name) => !rows.some((p) => p.name === name)) : null;
	const defaultProvider = view.configuredKnown && !view.configured && (availablePresets ?? []).includes("stepfun") ? OWNER_DEFAULT_PROVIDER : "";
	return /* @__PURE__ */ jsxs(SettingsCard, {
		title: "供应商",
		description: "密钥只保存在本机；界面上不会再显示密钥内容。",
		children: [rows.length === 0 && !addOpen ? /* @__PURE__ */ jsx(EmptyState, {
			actionLabel: COPY.addProvider,
			onAction: () => setAddOpen(true),
			text: COPY.noProvidersYet
		}) : /* @__PURE__ */ jsxs("div", {
			className: "flex flex-col",
			children: [rows.map((row) => /* @__PURE__ */ jsx(ProviderRow, {
				autoFetch: autoFetch === row.name,
				busy,
				editing: editing === row.name,
				onAutoFetchHandled: () => setAutoFetch(null),
				onEditing: () => onEditing(editing === row.name ? null : row.name),
				onFeedback,
				performCredential,
				performFetchModels,
				performWrite,
				row,
				run
			}, row.name)), addOpen ? /* @__PURE__ */ jsx(AddProviderForm, {
				available: availablePresets,
				busy,
				defaultProvider,
				onCancel: () => setAddOpen(false),
				onFeedback,
				performCredential,
				performWrite,
				presets: view.presets,
				run,
				onDone: (name) => {
					setAddOpen(false);
					onEditing(name);
					setAutoFetch(name);
				}
			}) : /* @__PURE__ */ jsx("div", {
				className: "pt-3",
				children: /* @__PURE__ */ jsx(Button, {
					onClick: () => setAddOpen(true),
					size: "sm",
					variant: "secondary",
					children: COPY.addProvider
				})
			})]
		}), feedback ? /* @__PURE__ */ jsx("div", {
			className: "pt-2",
			children: /* @__PURE__ */ jsx(FeedbackLine, { feedback })
		}) : null]
	});
}
/** Exported for the render instrument (scripts/render-model-settings.tsx):
the key row's markup (type="password", value never in state) is asserted
on this piece directly. */
function ProviderRow({ row, busy, editing, autoFetch, onAutoFetchHandled, onEditing, onFeedback, performWrite, performCredential, performFetchModels, run }) {
	const key = keyLabel(row);
	const [fetch, setFetch] = useState({ phase: "idle" });
	function startFetch() {
		if (fetch.phase === "loading" || busy !== null) return;
		if (!editing) onEditing();
		setFetch({ phase: "loading" });
		performFetchModels(row.name).then((result) => {
			setFetch(result.ok ? {
				phase: "ready",
				models: result.models
			} : {
				phase: "error",
				feedback: result.feedback
			});
		});
	}
	const [confirmRemove, setConfirmRemove] = useState(false);
	function removeRow() {
		run(`row:${row.name}`, () => performWrite({
			field: "provider_row_delete",
			provider: row.name
		})).then((fb) => {
			setConfirmRemove(false);
			onFeedback(fb);
			if (fb.tone === "ok" && editing) onEditing();
		});
	}
	const startFetchRef = useRef(startFetch);
	startFetchRef.current = startFetch;
	const autoFired = useRef(false);
	useEffect(() => {
		if (!autoFetch || autoFired.current) return;
		autoFired.current = true;
		onAutoFetchHandled();
		startFetchRef.current();
	}, [autoFetch, onAutoFetchHandled]);
	return /* @__PURE__ */ jsxs("div", {
		className: "border-b border-line last:border-b-0",
		children: [/* @__PURE__ */ jsxs("div", {
			className: "flex flex-wrap items-center gap-x-3 gap-y-1.5 py-2.5",
			children: [
				/* @__PURE__ */ jsx(StateDot$1, { tone: key.tone }),
				/* @__PURE__ */ jsx("span", {
					className: "shrink-0 text-[12.5px] font-medium text-ink",
					children: row.name
				}),
				row.billing === "plan" || row.billing === "pay-per-token" ? /* @__PURE__ */ jsx("span", {
					className: "shrink-0 rounded-chip bg-inset px-1.5 py-0.5 text-[10.5px] text-ink-2",
					children: row.billing === "plan" ? "套餐" : "按量"
				}) : null,
				row.protocol !== "" ? /* @__PURE__ */ jsx("span", {
					className: "shrink-0 font-mono text-[10.5px] text-ink-3",
					children: row.protocol
				}) : null,
				/* @__PURE__ */ jsx("span", {
					className: "min-w-0 flex-1 truncate font-mono text-[11px] text-ink-3",
					title: row.baseURL,
					children: row.baseURL !== "" ? row.baseURL : "默认基址（预设）"
				}),
				/* @__PURE__ */ jsx("span", {
					className: "shrink-0 text-[11px] text-ink-3",
					children: key.label
				}),
				/* @__PURE__ */ jsx("span", {
					className: "shrink-0 text-[11px] text-ink-3",
					children: row.modelCount > 0 ? `${row.modelCount} 个模型` : "暂无模型"
				}),
				/* @__PURE__ */ jsxs("div", {
					className: "flex shrink-0 items-center gap-0.5",
					children: [
						/* @__PURE__ */ jsx(Button, {
							onClick: () => onFeedback({
								tone: "neutral",
								text: COPY.testConnectionUnavailable
							}),
							size: "xs",
							variant: "quiet",
							children: "测试连接"
						}),
						/* @__PURE__ */ jsx(Button, {
							"aria-busy": fetch.phase === "loading",
							disabled: fetch.phase === "loading" || busy !== null,
							onClick: startFetch,
							size: "xs",
							variant: "quiet",
							children: fetch.phase === "loading" ? /* @__PURE__ */ jsxs("span", {
								className: "flex items-center gap-1",
								children: [/* @__PURE__ */ jsx(LoaderCircle, {
									"aria-hidden": "true",
									className: "animate-spin",
									size: 11,
									strokeWidth: 2
								}), "拉取中…"]
							}) : "拉取模型列表"
						}),
						/* @__PURE__ */ jsx(Button, {
							onClick: onEditing,
							size: "xs",
							variant: "quiet",
							children: editing ? "收起" : "编辑"
						}),
						confirmRemove ? /* @__PURE__ */ jsxs(Fragment$1, { children: [/* @__PURE__ */ jsx(Button, {
							disabled: busy !== null,
							onClick: removeRow,
							size: "xs",
							variant: "primary",
							children: "确认删除"
						}), /* @__PURE__ */ jsx(Button, {
							onClick: () => setConfirmRemove(false),
							size: "xs",
							variant: "quiet",
							children: "取消"
						})] }) : /* @__PURE__ */ jsx(Button, {
							onClick: () => setConfirmRemove(true),
							size: "xs",
							variant: "quiet",
							children: "删除"
						})
					]
				})
			]
		}), editing ? /* @__PURE__ */ jsx(ProviderEdit, {
			busy,
			fetch,
			onDismissFetch: () => setFetch({ phase: "idle" }),
			onFeedback,
			performCredential,
			performWrite,
			row,
			run
		}) : null]
	});
}
function ProviderEdit({ row, fetch, busy, onDismissFetch, onFeedback, performWrite, performCredential, run }) {
	const [baseURL, setBaseURL] = useState("");
	const [replacingKey, setReplacingKey] = useState(false);
	const [manualModel, setManualModel] = useState("");
	const [ticked, setTicked] = useState(/* @__PURE__ */ new Set());
	const keyInputRef = useRef(null);
	const existing = new Set(row.models ?? []);
	const fetchedNew = fetch.phase === "ready" ? fetch.models.filter((id) => !existing.has(id)) : [];
	function saveKey() {
		const input = keyInputRef.current;
		if (!input) return;
		const value = input.value;
		if (value === "") {
			onFeedback({
				tone: "error",
				text: "请先填入密钥"
			});
			return;
		}
		input.value = "";
		run(`credential:${row.name}`, () => performCredential(row.name, value)).then((fb) => {
			onFeedback(fb);
			if (fb.tone === "ok") setReplacingKey(false);
		});
	}
	function setBilling(value) {
		run(`billing:${row.name}`, () => performWrite({
			field: "provider_billing",
			provider: row.name,
			value
		})).then(onFeedback);
	}
	function addOneModel(id) {
		const model = id.trim();
		if (model === "") {
			onFeedback({
				tone: "error",
				text: "请先填入模型 id"
			});
			return;
		}
		run(`models:${row.name}`, () => performWrite({
			field: "model_row_create",
			provider: row.name,
			model
		})).then((fb) => {
			onFeedback(fb);
			if (fb.tone === "ok") setManualModel("");
		});
	}
	/** The taking-back half the operator asked for: one catalog entry, removed by
	name, with the setter's own refusal when it is already gone. */
	function removeOneModel(id) {
		run(`models:${row.name}`, () => performWrite({
			field: "model_row_delete",
			provider: row.name,
			model: id
		})).then(onFeedback);
	}
	/** One write per selected id - model_row_create is a single-entry step, and
	the answered receipt is per call. The summary line keeps the receipt's
	own meaning (saved + restart tier). */
	function addFetchedModels() {
		const ids = fetchedNew.filter((id) => ticked.has(id));
		if (ids.length === 0) {
			onFeedback({
				tone: "error",
				text: "先勾选要添加的模型"
			});
			return;
		}
		run(`models:${row.name}`, async () => {
			let saved = 0;
			let firstError = null;
			for (const id of ids) {
				const fb = await performWrite({
					field: "model_row_create",
					provider: row.name,
					model: id
				});
				if (fb.tone === "ok") saved += 1;
				else if (firstError === null) firstError = fb;
			}
			if (firstError === null) {
				onDismissFetch();
				return {
					tone: "ok",
					text: `已保存 ${saved} 个模型，重启后生效`
				};
			}
			return {
				tone: "error",
				text: saved > 0 ? `${saved} 个已保存；其余没有写成：${firstError.text}` : firstError.text
			};
		}).then(onFeedback);
	}
	return /* @__PURE__ */ jsxs("div", {
		className: "flex flex-col gap-3 pb-3 pl-5",
		children: [
			/* @__PURE__ */ jsxs("div", {
				className: "flex flex-wrap items-center gap-2",
				children: [
					/* @__PURE__ */ jsx("span", {
						className: "w-16 shrink-0 text-[11.5px] text-ink-3",
						children: "基址"
					}),
					/* @__PURE__ */ jsx("input", {
						"aria-label": `${row.name} 的接口基址`,
						className: "h-[30px] w-[320px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3",
						onChange: (e) => setBaseURL(e.currentTarget.value),
						placeholder: "https://…",
						value: baseURL
					}),
					/* @__PURE__ */ jsx(Button, {
						disabled: busy !== null,
						onClick: () => {
							if (baseURL.trim() === "") {
								onFeedback({
									tone: "error",
									text: "基址不能为空"
								});
								return;
							}
							run(`base:${row.name}`, () => performWrite({
								field: "provider_base_url",
								provider: row.name,
								value: baseURL.trim()
							})).then(onFeedback);
						},
						size: "xs",
						variant: "secondary",
						children: "保存基址"
					})
				]
			}),
			/* @__PURE__ */ jsxs("div", {
				className: "flex flex-wrap items-center gap-2",
				children: [
					/* @__PURE__ */ jsx("span", {
						className: "w-16 shrink-0 text-[11.5px] text-ink-3",
						children: "密钥"
					}),
					row.credentialSet && !replacingKey ? /* @__PURE__ */ jsxs(Fragment$1, { children: [/* @__PURE__ */ jsxs("span", {
						className: "text-[11.5px] text-ink-2",
						children: ["已保存", row.apiKeyRef !== "" ? `（${row.apiKeyRef}）` : ""]
					}), /* @__PURE__ */ jsx(Button, {
						onClick: () => setReplacingKey(true),
						size: "xs",
						variant: "secondary",
						children: "更换"
					})] }) : /* @__PURE__ */ jsxs(Fragment$1, { children: [
						/* @__PURE__ */ jsx("input", {
							"aria-label": `${row.name} 的密钥`,
							className: "h-[30px] w-[260px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3",
							onKeyDown: (e) => {
								if (e.key === "Enter") saveKey();
							},
							placeholder: "粘贴密钥（保存后不再显示）",
							ref: keyInputRef,
							type: "password"
						}),
						/* @__PURE__ */ jsx(Button, {
							disabled: busy !== null,
							onClick: saveKey,
							size: "xs",
							variant: "secondary",
							children: "保存密钥"
						}),
						replacingKey ? /* @__PURE__ */ jsx(Button, {
							onClick: () => setReplacingKey(false),
							size: "xs",
							variant: "quiet",
							children: "取消"
						}) : null
					] }),
					/* @__PURE__ */ jsx(Button, {
						onClick: () => onFeedback({
							tone: "neutral",
							text: COPY.testKeyUnavailable
						}),
						size: "xs",
						variant: "quiet",
						children: "测试密钥"
					})
				]
			}),
			/* @__PURE__ */ jsxs("div", {
				className: "flex flex-wrap items-center gap-2",
				children: [
					/* @__PURE__ */ jsx("span", {
						className: "w-16 shrink-0 text-[11.5px] text-ink-3",
						children: "计费模式"
					}),
					/* @__PURE__ */ jsx(Button, {
						disabled: busy !== null,
						onClick: () => setBilling("plan"),
						size: "xs",
						variant: row.billing === "plan" ? "primary" : "secondary",
						children: "套餐"
					}),
					/* @__PURE__ */ jsx(Button, {
						disabled: busy !== null,
						onClick: () => setBilling("pay-per-token"),
						size: "xs",
						variant: row.billing === "pay-per-token" ? "primary" : "secondary",
						children: "按量"
					}),
					/* @__PURE__ */ jsx("span", {
						className: "text-[11px] text-ink-3",
						children: row.billing === "plan" ? "当前：套餐（点击「按量」即改）" : row.billing === "pay-per-token" ? "当前：按量（点击「套餐」即改）" : "未设置；点击即保存（同一厂商两个计费口各记各的）"
					})
				]
			}),
			/* @__PURE__ */ jsxs("div", {
				className: "flex flex-col gap-2",
				children: [
					/* @__PURE__ */ jsx("span", {
						className: "text-[11.5px] text-ink-3",
						children: "模型目录"
					}),
					/* @__PURE__ */ jsx("div", {
						className: "flex flex-wrap items-center gap-1.5",
						children: row.models !== null && row.models.length > 0 ? row.models.map((id) => /* @__PURE__ */ jsxs("span", {
							className: "inline-flex h-5 items-center gap-1 rounded-chip bg-inset pl-2 pr-1 font-mono text-[11px] text-ink-2",
							children: [id, /* @__PURE__ */ jsx("button", {
								"aria-label": `删除模型 ${id}`,
								className: "rounded-full px-1 text-[13px] leading-none text-ink-3 hover:text-ink disabled:opacity-40",
								disabled: busy !== null,
								onClick: () => removeOneModel(id),
								type: "button",
								children: "×"
							})]
						}, id)) : /* @__PURE__ */ jsx("span", {
							className: "text-[11.5px] text-ink-3",
							children: row.models === null ? "暂时读不到模型清单" : "暂无模型"
						})
					}),
					fetch.phase === "loading" ? /* @__PURE__ */ jsxs("p", {
						className: "flex items-center gap-1.5 text-[11.5px] text-ink-3",
						children: [/* @__PURE__ */ jsx(LoaderCircle, {
							"aria-hidden": "true",
							className: "animate-spin",
							size: 11,
							strokeWidth: 2
						}), "正在向上游拉取模型列表…"]
					}) : fetch.phase === "error" ? /* @__PURE__ */ jsxs("div", {
						className: "flex flex-col gap-1",
						children: [/* @__PURE__ */ jsx(FeedbackLine, { feedback: fetch.feedback }), /* @__PURE__ */ jsx("div", { children: /* @__PURE__ */ jsx(Button, {
							onClick: onDismissFetch,
							size: "xs",
							variant: "quiet",
							children: "知道了"
						}) })]
					}) : fetch.phase === "ready" ? /* @__PURE__ */ jsxs("div", {
						className: "flex flex-col gap-2 rounded-control bg-inset p-2.5",
						children: [fetch.models.length === 0 ? /* @__PURE__ */ jsx("p", {
							className: "text-[11.5px] text-ink-3",
							children: "上游返回的清单是空的。"
						}) : fetchedNew.length === 0 ? /* @__PURE__ */ jsx("p", {
							className: "text-[11.5px] text-ink-3",
							children: "上游清单里的模型都已经在目录里了。"
						}) : /* @__PURE__ */ jsxs(Fragment$1, { children: [
							/* @__PURE__ */ jsxs("p", {
								className: "text-[11.5px] text-ink-3",
								children: [
									"拉到 ",
									fetchedNew.length,
									" 个新模型，勾选后写入："
								]
							}),
							/* @__PURE__ */ jsx("div", {
								className: "flex max-h-44 flex-col gap-0.5 overflow-y-auto",
								children: fetchedNew.map((id) => /* @__PURE__ */ jsxs("label", {
									className: "flex items-center gap-2 rounded-chip px-1.5 py-1 hover:bg-hover",
									children: [/* @__PURE__ */ jsx("input", {
										checked: ticked.has(id),
										className: "size-3.5 accent-primary",
										onChange: (e) => {
											const isChecked = e.currentTarget.checked;
											setTicked((current) => {
												const next = new Set(current);
												if (isChecked) next.add(id);
												else next.delete(id);
												return next;
											});
										},
										type: "checkbox"
									}), /* @__PURE__ */ jsx("span", {
										className: "font-mono text-[11.5px] text-ink",
										children: id
									})]
								}, id))
							}),
							/* @__PURE__ */ jsxs("div", {
								className: "flex items-center gap-2",
								children: [
									/* @__PURE__ */ jsx(Button, {
										disabled: busy !== null,
										onClick: addFetchedModels,
										size: "xs",
										variant: "primary",
										children: "写入所选"
									}),
									/* @__PURE__ */ jsx(Button, {
										onClick: () => setTicked((current) => {
											const next = new Set(current);
											for (const id of fetchedNew) next.add(id);
											return next;
										}),
										size: "xs",
										variant: "quiet",
										children: "全选"
									}),
									/* @__PURE__ */ jsx(Button, {
										onClick: onDismissFetch,
										size: "xs",
										variant: "quiet",
										children: "取消"
									})
								]
							})
						] }), fetch.models.length > 0 && fetchedNew.length === 0 ? /* @__PURE__ */ jsx("div", { children: /* @__PURE__ */ jsx(Button, {
							onClick: onDismissFetch,
							size: "xs",
							variant: "quiet",
							children: "取消"
						}) }) : null]
					}) : null,
					/* @__PURE__ */ jsxs("div", {
						className: "flex flex-wrap items-center gap-2",
						children: [/* @__PURE__ */ jsx("input", {
							"aria-label": "模型 id",
							className: "h-[30px] w-[220px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3",
							onChange: (e) => setManualModel(e.currentTarget.value),
							onKeyDown: (e) => {
								if (e.key === "Enter") addOneModel(manualModel);
							},
							placeholder: "模型 id",
							value: manualModel
						}), /* @__PURE__ */ jsx(Button, {
							disabled: busy !== null,
							onClick: () => addOneModel(manualModel),
							size: "xs",
							variant: "secondary",
							children: "添加"
						})]
					})
				]
			})
		]
	});
}
/** Exported for the render instrument (scripts/render-model-settings.tsx):
static markup cannot click the 添加供应商 button, so the instrument mounts
the pure form directly. */
function AddProviderForm({ available, presets, busy, defaultProvider, onCancel, onDone, onFeedback, performWrite, performCredential, run }) {
	const MODES = ["内置供应商", "自定义接入"];
	const BILLING = ["套餐（plan）", "按量（pay-per-token）"];
	const [selectedMode, setSelectedMode] = useState("内置供应商");
	const [provider, setProvider] = useState(defaultProvider);
	const [customName, setCustomName] = useState("");
	const [protocol, setProtocol] = useState("");
	const [customBase, setCustomBase] = useState("");
	const [billing, setBilling] = useState("");
	const [model, setModel] = useState("");
	const keyInputRef = useRef(null);
	const KNOWN_ALT_ENDPOINTS = { stepfun: {
		label: "按量（API）口",
		base: "https://api.stepfun.com/v1"
	} };
	const vendorKey = customName.trim().split("-")[0]?.toLowerCase() ?? "";
	const altEndpoint = KNOWN_ALT_ENDPOINTS[vendorKey];
	const takenAll = (presets ?? []).filter((name) => !(available ?? []).includes(name));
	const secondRows = takenAll.filter((name) => KNOWN_ALT_ENDPOINTS[name.toLowerCase()]).map((name) => {
		const alt = KNOWN_ALT_ENDPOINTS[name.toLowerCase()];
		return {
			name: `${name}-api`,
			label: alt.label,
			base: alt.base,
			vendor: name
		};
	});
	function prefillSecondRow(row) {
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
				onFeedback({
					tone: "error",
					text: "请先选择供应商"
				});
				return;
			}
		} else {
			if (customName.trim() === "") {
				onFeedback({
					tone: "error",
					text: "请先填行名"
				});
				return;
			}
			if (presets !== null && presets.includes(customName.trim())) {
				onFeedback({
					tone: "error",
					text: "这是内置供应商名，请用「内置供应商」添加"
				});
				return;
			}
			if (protocol === "") {
				onFeedback({
					tone: "error",
					text: "请先选择协议"
				});
				return;
			}
			if (customBase.trim() === "") {
				onFeedback({
					tone: "error",
					text: "请先填基址"
				});
				return;
			}
		}
		if (keyInputRef.current) keyInputRef.current.value = "";
		const name = selectedMode === "内置供应商" ? provider : customName.trim();
		run("add-provider", async () => {
			const created = await performWrite(selectedMode === "内置供应商" ? {
				field: "provider_row_create",
				provider: name,
				...model.trim() !== "" ? { model: model.trim() } : {}
			} : {
				field: "provider_row_create_custom",
				provider: name,
				value: `${protocol},${customBase.trim()}`,
				...model.trim() !== "" ? { model: model.trim() } : {}
			});
			if (created.tone !== "ok") return created;
			if (billingValue !== "") {
				const billed = await performWrite({
					field: "provider_billing",
					provider: name,
					value: billingValue
				});
				if (billed.tone !== "ok") return {
					tone: "error",
					text: `供应商已添加，但计费模式没有写入：${billed.text}`
				};
			}
			if (keyValue === "") return created;
			const keyed = await performCredential(name, keyValue);
			return keyed.tone === "ok" ? created : {
				tone: "error",
				text: `供应商已添加，但密钥没有保存：${keyed.text}`
			};
		}).then((fb) => {
			onFeedback(fb);
			if (fb.tone === "ok") onDone(name);
		});
	}
	const rowLabel = "w-16 shrink-0 text-[11.5px] text-ink-3";
	const inputCls = "h-[30px] rounded-control bg-surface px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3";
	return /* @__PURE__ */ jsxs("div", {
		className: "mt-1 flex flex-col gap-2.5 rounded-control bg-inset p-3",
		style: { animation: "fade-up 200ms var(--ease-out-strong) both" },
		children: [
			/* @__PURE__ */ jsx("div", {
				className: "flex items-center gap-2",
				children: /* @__PURE__ */ jsx(SegmentedControl, {
					onChange: setSelectedMode,
					options: MODES,
					value: selectedMode
				})
			}),
			selectedMode === "内置供应商" ? /* @__PURE__ */ jsxs("div", {
				className: "flex flex-col gap-2",
				children: [
					/* @__PURE__ */ jsxs("div", {
						className: "flex flex-wrap items-center gap-2",
						children: [
							/* @__PURE__ */ jsx("span", {
								className: rowLabel,
								children: "供应商"
							}),
							available === null ? /* @__PURE__ */ jsx("span", {
								className: "text-[11.5px] text-ink-3",
								children: "暂未读取到供应商列表，请稍后刷新"
							}) : /* @__PURE__ */ jsx(Picker, {
								ariaLabel: "供应商",
								className: "w-[220px]",
								emptyText: "内置供应商都已添加",
								onSelect: setProvider,
								options: available.map((name) => ({
									value: name,
									note: name === "stepfun" ? "推荐" : ""
								})),
								placeholder: "选择供应商",
								value: provider
							}),
							/* @__PURE__ */ jsx("span", {
								className: "text-[11px] text-ink-3",
								children: "这里只列还没添加过的内置供应商；同一家要第二行请用下面那条"
							})
						]
					}),
					takenAll.length > 0 ? /* @__PURE__ */ jsx("div", {
						className: "flex flex-wrap items-center gap-2",
						children: /* @__PURE__ */ jsxs("span", {
							className: "text-[11px] text-ink-3",
							children: ["已添加：", takenAll.join("、")]
						})
					}) : null,
					secondRows.length > 0 ? /* @__PURE__ */ jsxs("div", {
						className: "flex flex-wrap items-center gap-2 rounded-control bg-surface/60 px-2 py-1.5",
						children: [
							/* @__PURE__ */ jsx("span", {
								className: rowLabel,
								children: "常用第二行"
							}),
							secondRows.map((row) => /* @__PURE__ */ jsxs(Button, {
								onClick: () => prefillSecondRow(row),
								size: "xs",
								variant: "quiet",
								children: [
									"加一行 ",
									row.name,
									"（",
									row.label,
									"）"
								]
							}, row.name)),
							/* @__PURE__ */ jsx("span", {
								className: "text-[11px] text-ink-3",
								children: "同一厂商两行并存：各用各的基址、各拉各的模型清单"
							})
						]
					}) : null,
					provider !== "" && KNOWN_ALT_ENDPOINTS[provider.toLowerCase()] ? /* @__PURE__ */ jsxs("div", {
						className: "flex flex-wrap items-center gap-2",
						children: [/* @__PURE__ */ jsxs("span", {
							className: "text-[11px] text-ink-3",
							children: [
								provider,
								" 还有另一个计费口（",
								KNOWN_ALT_ENDPOINTS[provider.toLowerCase()].label,
								"） ——同一厂商两行并存时，两行各有各的基址与模型目录："
							]
						}), /* @__PURE__ */ jsx(Button, {
							onClick: () => prefillSecondRow({
								name: `${provider}-api`,
								base: KNOWN_ALT_ENDPOINTS[provider.toLowerCase()].base
							}),
							size: "xs",
							variant: "quiet",
							children: "切到自定义并预填这一行"
						})]
					}) : null
				]
			}) : /* @__PURE__ */ jsxs(Fragment$1, { children: [
				/* @__PURE__ */ jsxs("div", {
					className: "flex flex-wrap items-center gap-2",
					children: [
						/* @__PURE__ */ jsx("span", {
							className: rowLabel,
							children: "行名"
						}),
						/* @__PURE__ */ jsx("input", {
							"aria-label": "行名",
							className: `${inputCls} w-[220px]`,
							onChange: (e) => setCustomName(e.currentTarget.value),
							placeholder: "例如 stepfun-api",
							value: customName
						}),
						/* @__PURE__ */ jsx("span", {
							className: "text-[11px] text-ink-3",
							children: "同一厂商可以两行并存，各用各的基址与密钥"
						})
					]
				}),
				/* @__PURE__ */ jsxs("div", {
					className: "flex flex-wrap items-center gap-2",
					children: [/* @__PURE__ */ jsx("span", {
						className: rowLabel,
						children: "协议"
					}), /* @__PURE__ */ jsx(Picker, {
						ariaLabel: "协议",
						className: "w-[220px]",
						onSelect: setProtocol,
						options: [
							{
								value: "openai-chat",
								note: "OpenAI 兼容"
							},
							{
								value: "openai-responses",
								note: "OpenAI 响应式"
							},
							{
								value: "anthropic",
								note: "Anthropic"
							}
						],
						placeholder: "选择协议",
						value: protocol
					})]
				}),
				/* @__PURE__ */ jsxs("div", {
					className: "flex flex-wrap items-center gap-2",
					children: [/* @__PURE__ */ jsx("span", {
						className: rowLabel,
						children: "基址"
					}), /* @__PURE__ */ jsx("input", {
						"aria-label": "基址",
						className: `${inputCls} w-[320px]`,
						onChange: (e) => setCustomBase(e.currentTarget.value),
						placeholder: "https://…",
						value: customBase
					})]
				}),
				altEndpoint ? /* @__PURE__ */ jsxs("div", {
					className: "flex flex-wrap items-center gap-2",
					children: [/* @__PURE__ */ jsxs("span", {
						className: "text-[11px] text-ink-3",
						children: [
							vendorKey,
							" 的",
							altEndpoint.label,
							"已知基址：",
							altEndpoint.base
						]
					}), /* @__PURE__ */ jsx(Button, {
						onClick: () => {
							setCustomBase(altEndpoint.base);
							if (billing === "") setBilling(BILLING[1]);
						},
						size: "xs",
						variant: "quiet",
						children: "填入"
					})]
				}) : null,
				/* @__PURE__ */ jsxs("div", {
					className: "flex flex-wrap items-center gap-2",
					children: [
						/* @__PURE__ */ jsx("span", {
							className: rowLabel,
							children: "计费模式"
						}),
						/* @__PURE__ */ jsx(SegmentedControl, {
							onChange: (next) => setBilling(next === BILLING[0] ? BILLING[0] : next === BILLING[1] ? BILLING[1] : ""),
							options: BILLING,
							value: billing === BILLING[0] ? BILLING[0] : billing === BILLING[1] ? BILLING[1] : ""
						}),
						/* @__PURE__ */ jsx("span", {
							className: "text-[11px] text-ink-3",
							children: "选填；同一厂商两个计费口各记各的"
						})
					]
				})
			] }),
			/* @__PURE__ */ jsxs("div", {
				className: "flex flex-wrap items-center gap-2",
				children: [
					/* @__PURE__ */ jsx("span", {
						className: rowLabel,
						children: "模型 id"
					}),
					/* @__PURE__ */ jsx("input", {
						"aria-label": "模型 id（选填）",
						className: `${inputCls} w-[220px]`,
						onChange: (e) => setModel(e.currentTarget.value),
						placeholder: "选填",
						value: model
					}),
					/* @__PURE__ */ jsx("span", {
						className: "text-[11px] text-ink-3",
						children: "选填；添加后会自动拉一次这行的模型清单"
					})
				]
			}),
			/* @__PURE__ */ jsxs("div", {
				className: "flex flex-wrap items-center gap-2",
				children: [/* @__PURE__ */ jsx("span", {
					className: rowLabel,
					children: "密钥"
				}), /* @__PURE__ */ jsx("input", {
					"aria-label": "密钥（选填）",
					className: `${inputCls} w-[260px]`,
					placeholder: "粘贴密钥（保存后不再显示）",
					ref: keyInputRef,
					type: "password"
				})]
			}),
			/* @__PURE__ */ jsxs("div", {
				className: "flex items-center gap-2 pl-18",
				children: [/* @__PURE__ */ jsx(Button, {
					disabled: busy !== null,
					onClick: submit,
					size: "sm",
					variant: "primary",
					children: busy === "add-provider" ? COPY.saving : "添加"
				}), /* @__PURE__ */ jsx(Button, {
					onClick: onCancel,
					size: "sm",
					variant: "quiet",
					children: "取消"
				})]
			})
		]
	});
}
function ChatModelCard({ view, busy, chatProvider, chatModel, onProvider, onModel, feedback, onFeedback, performWrite, run }) {
	const current = view.chatModelKnown ? view.chatProvider !== "" && view.chatModel !== "" ? `${view.chatProvider} / ${view.chatModel}` : "未设置" : "暂未读取";
	function save() {
		if (chatProvider === "" || chatModel === "") {
			const row = view.providers.find((p) => p.name === chatProvider);
			const noModels = row != null && (row.models === null || row.models.length === 0);
			onFeedback({
				tone: "error",
				text: chatProvider === "" ? COPY.chatPairMissing : noModels ? "这家供应商还没有模型；在供应商行的「编辑」里用「拉取模型列表」或手动添加" : COPY.pairMissing
			});
			return;
		}
		run("chat-pair", () => performWrite({
			field: "role_chat_pair",
			provider: chatProvider,
			model: chatModel
		})).then(onFeedback);
	}
	return /* @__PURE__ */ jsx(SettingsCard, {
		title: "默认聊天模型",
		description: "对话默认使用的供应商与模型。",
		children: view.providers.length === 0 ? /* @__PURE__ */ jsx(EmptyState, { text: "先添加一个供应商，然后就能选择聊天模型。" }) : /* @__PURE__ */ jsxs("div", {
			className: "flex flex-col gap-2.5",
			children: [
				/* @__PURE__ */ jsxs("div", {
					className: "flex flex-wrap items-center gap-2",
					children: [/* @__PURE__ */ jsx(PairFields, {
						manualFallback: false,
						model: chatModel,
						onModel,
						onProvider,
						provider: chatProvider,
						providers: view.providers
					}), /* @__PURE__ */ jsx(Button, {
						disabled: busy !== null,
						onClick: save,
						size: "sm",
						variant: "primary",
						children: busy === "chat-pair" ? COPY.saving : COPY.save
					})]
				}),
				/* @__PURE__ */ jsxs("p", {
					className: "text-[11.5px] text-ink-3",
					children: ["当前：", current]
				}),
				/* @__PURE__ */ jsx(FeedbackLine, { feedback })
			]
		})
	});
}
//#endregion
//#region src/components/voice-settings.tsx
/** The census ruling these three cards carry: the save is real, the effect is
not. The generic translated receipt ("已保存，重启后生效") promises an
effect no runtime path can deliver today, so every ok answer on this page
is restated as a warn - and the tone downgrade is the point: a green
"saved" line on a card that cannot take effect reads as "working". */
function placeholderReceipt(feedback) {
	return feedback.tone === "ok" ? {
		tone: "warn",
		text: "已写入配置文件；这一项还没有接到运行路径，当前不生效"
	} : feedback;
}
function VoiceSettings({ heading = "语音", lede = "云端听写链、云端朗读链与实时语音的供应商和模型（三张卡都是占位：写入会保存，但云端链路与实时引擎还没有接到运行路径，当前不生效）；本机音频的后端读者已接真，面板先等快照。" }) {
	return /* @__PURE__ */ jsxs("section", {
		"aria-label": "语音",
		className: "flex min-w-0 flex-col gap-1",
		children: [
			/* @__PURE__ */ jsx("h2", {
				className: "text-[15px] font-semibold text-ink",
				children: heading
			}),
			/* @__PURE__ */ jsx("p", {
				className: "text-[12.5px] text-ink-2",
				children: lede
			}),
			/* @__PURE__ */ jsx(ConfigReadGate, { children: (view, _prose, reload) => /* @__PURE__ */ jsx(VoiceSettingsView, {
				onReload: reload,
				performWrite: performConfigWrite,
				view
			}) })
		]
	});
}
function VoiceSettingsView({ view, performWrite, onReload }) {
	const status = configStatus(view);
	const [busy, setBusy] = useState(null);
	async function run(key, fn) {
		setBusy(key);
		const feedback = await fn();
		setBusy(null);
		onReload();
		return feedback;
	}
	return /* @__PURE__ */ jsxs("div", {
		className: "mt-2 flex min-w-0 flex-col gap-3",
		children: [/* @__PURE__ */ jsxs("div", {
			className: "flex items-center gap-2",
			children: [
				/* @__PURE__ */ jsx(StateTag, {
					text: status.label,
					tone: status.tone
				}),
				/* @__PURE__ */ jsx("span", { className: "min-w-0 flex-1" }),
				/* @__PURE__ */ jsxs(Button, {
					"aria-label": "刷新",
					onClick: onReload,
					size: "xs",
					variant: "quiet",
					children: [/* @__PURE__ */ jsx(RefreshCw, {
						"aria-hidden": "true",
						size: 11,
						strokeWidth: 2
					}), "刷新"]
				})
			]
		}), !view.readable ? /* @__PURE__ */ jsx(SettingsCard, {
			title: "配置",
			children: /* @__PURE__ */ jsx("p", {
				className: "text-[12px] leading-[1.7] text-ink-2",
				children: "暂时读不到配置，请稍后刷新。"
			})
		}) : view.providers.length === 0 ? /* @__PURE__ */ jsx(SettingsCard, {
			title: "语音",
			children: /* @__PURE__ */ jsx(EmptyState, { text: "先在「模型」页添加一个供应商，再回来配置语音。" })
		}) : /* @__PURE__ */ jsxs(Fragment$1, { children: [
			/* @__PURE__ */ jsx(ChainEditor, {
				busy,
				description: "占位卡：常驻语音的听写目前只走本地引擎。这里列出的云端顺序可以写入配置，但还没有接到听写运行路径，接入前改它不会生效。",
				field: "voice_cloud_asr_chain",
				initial: view.voice.cloudASRChain,
				performWrite,
				providers: view.providers,
				run,
				title: "听写（ASR）"
			}),
			/* @__PURE__ */ jsx(ChainEditor, {
				busy,
				description: "占位卡：朗读的云端引擎尚未接入。这里列出的顺序可以写入配置，但当前没有运行路径按它取朗读，接入前改它不会生效。",
				field: "voice_cloud_tts_chain",
				initial: view.voice.cloudTTSChain,
				performWrite,
				providers: view.providers,
				run,
				title: "朗读（TTS）"
			}),
			/* @__PURE__ */ jsx(RealtimeCard, {
				busy,
				performWrite,
				providers: view.providers,
				run,
				view
			}),
			/* @__PURE__ */ jsx(LocalAudioCard, {})
		] })]
	});
}
/** "provider/model" as the config spells a chain element; a hand-written
element without a slash keeps its whole text in the provider half so the
editor never silently drops it. */
function parseChain(chain) {
	if (chain === null) return [];
	return chain.map((el) => {
		const cut = el.indexOf("/");
		return cut < 0 ? {
			provider: el,
			model: ""
		} : {
			provider: el.slice(0, cut),
			model: el.slice(cut + 1)
		};
	});
}
function chainText(chain) {
	if (chain === null || chain.length === 0) return "未设置";
	return chain.join(", ");
}
function ChainEditor({ title, description, field, initial, providers, busy, run, performWrite }) {
	const [segments, setSegments] = useState(() => {
		const parsed = parseChain(initial);
		return parsed.length > 0 ? parsed : [{
			provider: "",
			model: ""
		}];
	});
	const [feedback, setFeedback] = useState(null);
	function update(index, next) {
		setSegments((all) => all.map((seg, i) => i === index ? {
			...seg,
			...next
		} : seg));
	}
	function save() {
		if (segments.length === 0) {
			setFeedback({
				tone: "error",
				text: COPY.chainEmpty
			});
			return;
		}
		for (const seg of segments) if (seg.provider === "" || seg.model === "") {
			const row = providers.find((p) => p.name === seg.provider);
			const noModels = seg.provider !== "" && row !== void 0 && (row.models === null || row.models.length === 0);
			setFeedback({
				tone: "error",
				text: noModels ? "这家供应商还没有模型" : COPY.chainSegmentMissing
			});
			return;
		}
		run(`${field}`, () => performWrite({
			field,
			value: segments.map((seg) => `${seg.provider}/${seg.model}`).join(",")
		})).then((feedback) => setFeedback(placeholderReceipt(feedback)));
	}
	return /* @__PURE__ */ jsx(SettingsCard, {
		title,
		description,
		children: /* @__PURE__ */ jsxs("div", {
			className: "flex flex-col gap-2",
			children: [
				/* @__PURE__ */ jsx(StateTag, {
					text: "占位·当前不生效",
					tone: "warn"
				}),
				/* @__PURE__ */ jsx("div", {
					className: "flex flex-col gap-1.5",
					children: segments.map((seg, index) => /* @__PURE__ */ jsxs("div", {
						className: "flex flex-wrap items-center gap-2",
						children: [
							/* @__PURE__ */ jsx("span", {
								className: "w-8 shrink-0 font-mono text-[11px] text-ink-3",
								children: index + 1
							}),
							/* @__PURE__ */ jsx(PairFields, {
								manualFallback: false,
								model: seg.model,
								onModel: (next) => update(index, { model: next }),
								onProvider: (next) => update(index, {
									provider: next,
									model: ""
								}),
								provider: seg.provider,
								providers
							}),
							/* @__PURE__ */ jsx(Button, {
								"aria-label": `移除第 ${index + 1} 段`,
								onClick: () => setSegments((all) => all.filter((_, i) => i !== index)),
								size: "xs",
								variant: "quiet",
								children: /* @__PURE__ */ jsx(X, {
									"aria-hidden": "true",
									size: 12,
									strokeWidth: 2
								})
							})
						]
					}, `${index}-${seg.provider}`))
				}),
				/* @__PURE__ */ jsxs("div", {
					className: "flex flex-wrap items-center gap-2",
					children: [/* @__PURE__ */ jsxs(Button, {
						onClick: () => setSegments((all) => [...all, {
							provider: "",
							model: ""
						}]),
						size: "xs",
						variant: "secondary",
						children: [/* @__PURE__ */ jsx(Plus, {
							"aria-hidden": "true",
							size: 12,
							strokeWidth: 2
						}), "添加一段"]
					}), /* @__PURE__ */ jsx(Button, {
						disabled: busy !== null,
						onClick: save,
						size: "sm",
						variant: "primary",
						children: busy === field ? COPY.saving : COPY.save
					})]
				}),
				/* @__PURE__ */ jsxs("p", {
					className: "text-[11.5px] text-ink-3",
					children: ["当前：", chainText(initial)]
				}),
				/* @__PURE__ */ jsx(FeedbackLine, { feedback })
			]
		})
	});
}
/** One row of the local-audio card: the key as the config spells it, the
backend fact (why it is real today), and the waiting state. */
var LOCAL_AUDIO_ROWS = [
	{
		key: "voice.aec.enabled",
		name: "回声消除",
		fact: "后端已接：喂给打断监测的回声消除开关（重启后生效）"
	},
	{
		key: "audio.input_device",
		name: "输入设备",
		fact: "后端已接：每次开麦与热插拔重开都按这个名字解析设备"
	},
	{
		key: "audio.sample_rate",
		name: "采样率",
		fact: "后端守卫：采集链只跑 16000，配置里写别的值会被拒"
	}
];
/** The three local-audio keys the backend already reads for real, but the
panel cannot show or change yet: the config.get snapshot's voice payload
carries five keys only (pinned in config_313_test.go), none of them these
three, and the config.set whitelist has no write field for them either.
So this card is the honest WAITING state - the backend fact as words plus
"等快照" - never a guessed value, never a control that looks switchable.
When Go adds the snapshot keys and the write fields, this card graduates
to real controls; drawing them now would re-open the C-class lie. */
function LocalAudioCard() {
	return /* @__PURE__ */ jsx(SettingsCard, {
		title: "本机音频",
		description: "这三项的后端读者/守卫已接（真），但配置读口的快照还没有它们的当前值，写入路由也未接。快照补上之前这里显示等待态，不画能点但改不动的控件。",
		children: /* @__PURE__ */ jsx("div", {
			className: "flex flex-col gap-1.5",
			children: LOCAL_AUDIO_ROWS.map((row) => /* @__PURE__ */ jsxs("div", {
				className: "flex flex-wrap items-baseline gap-x-2 gap-y-0.5",
				children: [
					/* @__PURE__ */ jsx("span", {
						className: "text-[12px] font-medium text-ink",
						children: row.name
					}),
					/* @__PURE__ */ jsx("span", {
						className: "font-mono text-[10.5px] text-ink-3",
						children: row.key
					}),
					/* @__PURE__ */ jsx("span", {
						className: "min-w-0 flex-1 basis-40 text-[11.5px] leading-[1.6] text-ink-2",
						children: row.fact
					}),
					/* @__PURE__ */ jsx(StateTag, {
						text: "等快照",
						tone: "neutral"
					})
				]
			}, row.key))
		})
	});
}
function RealtimeCard({ view, providers, busy, run, performWrite }) {
	const [provider, setProvider] = useState(view.voice.realtimeProvider);
	const [model, setModel] = useState(view.voice.realtimeModel);
	const [feedback, setFeedback] = useState(null);
	const current = view.voice.realtimeProvider !== "" || view.voice.realtimeModel !== "" ? `${view.voice.realtimeProvider} / ${view.voice.realtimeModel}` : "未设置";
	function save() {
		if (provider === "" || model === "") {
			setFeedback({
				tone: "error",
				text: COPY.pairMissing
			});
			return;
		}
		run("realtime", async () => {
			const head = await performWrite({
				field: "voice_realtime_provider",
				value: provider
			});
			if (head.tone !== "ok") return head;
			const tail = await performWrite({
				field: "voice_realtime_model",
				value: model
			});
			return tail.tone === "ok" ? placeholderReceipt(tail) : tail;
		}).then(setFeedback);
	}
	return /* @__PURE__ */ jsx(SettingsCard, {
		title: "实时语音",
		description: "占位卡：实时通话引擎尚未实现，这里配置的供应商与模型暂时没有通话链使用；写入会保存，接入前不生效。",
		children: /* @__PURE__ */ jsxs("div", {
			className: "flex flex-col gap-2",
			children: [
				/* @__PURE__ */ jsx(StateTag, {
					text: "占位·当前不生效",
					tone: "warn"
				}),
				/* @__PURE__ */ jsxs("div", {
					className: "flex flex-wrap items-center gap-2",
					children: [/* @__PURE__ */ jsx(PairFields, {
						manualFallback: true,
						model,
						onModel: setModel,
						onProvider: (next) => {
							setProvider(next);
							setModel("");
						},
						provider,
						providers
					}), /* @__PURE__ */ jsx(Button, {
						disabled: busy !== null,
						onClick: save,
						size: "sm",
						variant: "primary",
						children: busy === "realtime" ? COPY.saving : COPY.save
					})]
				}),
				/* @__PURE__ */ jsxs("p", {
					className: "text-[11.5px] text-ink-3",
					children: ["当前：", current]
				}),
				/* @__PURE__ */ jsx(FeedbackLine, { feedback })
			]
		})
	});
}
//#endregion
//#region src/components/harness/settings-page.tsx
var CATEGORIES = [
	{
		name: "通用",
		purpose: "启动行为、默认工作区、面板与球的分工。",
		waiting: "配置读口（config.get）已经落地：模型、语音两节走读口回执，悬浮球与面板各节随快照推送。这一类自己的键（启动行为、默认工作区）还没有节推上来，等它们上读口再画。"
	},
	{
		name: "外观",
		purpose: "主题、面板底色深浅、动画档位、字号。",
		waiting: "",
		control: "config"
	},
	{
		name: "模型",
		purpose: "添加供应商、保存密钥，并选择默认聊天模型。",
		waiting: "",
		control: "model"
	},
	{
		name: "语音",
		purpose: "配置听写、朗读与实时语音使用的供应商和模型，以及本机音频的回声消除与输入设备。",
		waiting: "",
		control: "voice"
	},
	{
		name: "记忆",
		purpose: "长期记忆的开关、范围与清空。",
		waiting: "还没收到记忆数据。D20 定义了记忆系统，但没有任何一条推送告诉面板「记了什么」。"
	},
	{
		name: "技能",
		purpose: "装了哪些技能、每个技能的权限范围。",
		waiting: "还没收到技能清单。技能要能被列出来，前提是它们从插件运行时注册上来。"
	},
	{
		name: "MCP",
		purpose: "MCP 服务端、每个服务提供哪些工具。",
		waiting: "还没收到 MCP 状态。D13 定了 MCP 的取舍，但没有字段把已连的服务端报给面板。"
	},
	{
		name: "安全",
		purpose: "权限档位、门控阈值、批量授权的边界。",
		waiting: "还没收到策略快照。这一节最敏感：档位和阈值属于原生侧，面板在这里显示的任何数字都必须是读来的，不能是猜的。"
	},
	{
		name: "隐私",
		purpose: "哪些数据留在本机、录音留存多久、遥测开关。",
		waiting: "还没收到隐私配置。本节的所有项都要落到 config.toml 的读口上，没有读口就一个字也不能显示。"
	},
	{
		name: "成本",
		purpose: "每轮花费、本月累计、预算与超支提醒。",
		waiting: "还没收到用量数据。成本口径由 C23 定，但面板快照里没有 token 数，所以这一格连一个空壳数字都不画。"
	},
	{
		name: "引导",
		purpose: "首次使用的几步，以及随时能重看的说明。",
		waiting: "还没收到引导进度。首启那套画面在首启屏上，不在这里重画一遍。"
	},
	{
		name: "悬浮球",
		purpose: "球的大小、空闲透明度、鼠标穿透、唤醒词、热键。",
		waiting: "球的大小、外观风格、空闲透明度、穿透都已能切：去「外观」类的对应行，走设置写入通路落盘。唤醒词与热键还没有数据源。"
	}
];
function SettingsPage({ activeSection, onSelectSection, onClose, ballStyle = null, ballSize = null, ballOpacityIdle = null, ballClickThrough = null, panelWidth = null, panelFontSize = null }) {
	const cat = CATEGORIES.find((c) => c.name === activeSection) ?? CATEGORIES[0];
	return /* @__PURE__ */ jsxs("div", {
		className: "flex h-full min-h-0 flex-1 flex-col bg-page",
		children: [/* @__PURE__ */ jsxs("header", {
			className: "flex h-10 shrink-0 items-center gap-2 border-b border-line bg-surface px-3",
			children: [
				/* @__PURE__ */ jsxs("button", {
					type: "button",
					onClick: onClose,
					className: "flex shrink-0 items-center gap-1 rounded-control px-1.5 py-1 text-[12px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink",
					children: [/* @__PURE__ */ jsx(ArrowLeft, {
						"aria-hidden": "true",
						size: 13,
						strokeWidth: 2
					}), "返回"]
				}),
				/* @__PURE__ */ jsx("span", {
					"aria-hidden": "true",
					className: "h-3.5 w-px shrink-0 bg-line"
				}),
				/* @__PURE__ */ jsx("span", {
					className: "min-w-0 flex-1 truncate text-[12.5px] font-medium text-ink",
					children: cat.name
				})
			]
		}), /* @__PURE__ */ jsxs("div", {
			className: "flex min-h-0 flex-1 overflow-hidden",
			children: [/* @__PURE__ */ jsx("nav", {
				"aria-label": "设置分类",
				className: "flex w-48 shrink-0 flex-col overflow-y-auto border-r border-line bg-surface p-2",
				children: /* @__PURE__ */ jsx(GlideMenu, {
					className: "flex flex-col gap-px",
					highlightClassName: "inset-x-0 rounded-[7px] bg-hover",
					children: CATEGORIES.map((c) => {
						const on = c.name === cat.name;
						const name = c.waiting ? `${c.name}（还没有数据）` : c.name;
						return /* @__PURE__ */ jsxs("button", {
							type: "button",
							"aria-current": on ? "page" : void 0,
							"aria-label": name,
							"data-menu-row": true,
							title: c.purpose,
							onClick: () => onSelectSection(c.name),
							className: cn("relative z-10 flex w-full items-center rounded-[7px] py-1.5 pr-2 text-left text-[12.5px] transition-colors duration-100", on ? "pl-3.5 font-medium text-ink" : "pl-3.5 text-ink-2"),
							children: [
								on ? /* @__PURE__ */ jsx("span", {
									"aria-hidden": "true",
									className: "absolute bottom-1.5 left-1 top-1.5 w-[2px] rounded-full bg-accent"
								}) : null,
								/* @__PURE__ */ jsx("span", {
									"aria-hidden": "true",
									className: "min-w-0 flex-1 truncate",
									children: c.name
								}),
								c.waiting ? /* @__PURE__ */ jsx("span", {
									"aria-hidden": "true",
									className: "size-1.5 shrink-0 rounded-full bg-line-strong"
								}) : null
							]
						}, c.name);
					})
				})
			}), /* @__PURE__ */ jsx("div", {
				className: "min-h-0 flex-1 overflow-y-auto px-6 py-5",
				children: cat.control === "config" ? /* @__PURE__ */ jsx(ConfigScreen, {
					heading: cat.name,
					lede: cat.purpose,
					ballStyle,
					ballSize,
					ballOpacityIdle,
					ballClickThrough,
					panelWidth,
					panelFontSize
				}) : cat.control === "model" ? /* @__PURE__ */ jsx(ModelSettings, {
					heading: cat.name,
					lede: cat.purpose
				}) : cat.control === "voice" ? /* @__PURE__ */ jsx(VoiceSettings, {
					heading: cat.name,
					lede: cat.purpose
				}) : /* @__PURE__ */ jsxs("section", {
					"aria-label": cat.name,
					className: "flex min-w-0 flex-col gap-1",
					children: [
						/* @__PURE__ */ jsx("h2", {
							className: "text-[15px] font-semibold text-ink",
							children: cat.name
						}),
						/* @__PURE__ */ jsx("p", {
							className: "text-[12.5px] text-ink-2",
							children: cat.purpose
						}),
						/* @__PURE__ */ jsx("div", {
							className: "mt-4 rounded-card border border-line bg-surface px-3.5 py-3",
							children: /* @__PURE__ */ jsx("p", {
								className: "text-[12px] leading-[1.7] text-ink-3",
								children: cat.waiting
							})
						}),
						/* @__PURE__ */ jsx("p", {
							className: "mt-2.5 text-[11px] leading-[1.6] text-ink-3",
							children: "这一类不画控件。画一个空壳开关，等于让一个改不动的旋钮看起来像能改。"
						})
					]
				})
			})]
		})]
	});
}
//#endregion
//#region src/components/panel-sidebar.tsx
function SectionHeader({ label }) {
	return /* @__PURE__ */ jsx("div", {
		className: "px-2 pb-1 pt-2 text-[10.5px] font-medium tracking-[0.08em] text-ink-3",
		children: label
	});
}
/** The leading state glyph on a session row. A leaf slot is kept for `done` so
every row's text starts at the same x - a column that jitters when a
session changes state is unreadable. */
function StateGlyph({ state }) {
	return /* @__PURE__ */ jsx("span", {
		className: "flex w-3.5 shrink-0 justify-center",
		children: state === "failed" ? /* @__PURE__ */ jsx(CircleAlert, {
			"aria-hidden": "true",
			className: "shrink-0 text-red",
			size: 12,
			strokeWidth: 2
		}) : state === "streaming" ? /* @__PURE__ */ jsx("span", {
			"aria-hidden": "true",
			className: "block size-1.5 rounded-full bg-accent"
		}) : null
	});
}
function SessionRow({ session, active, onSelect, indent }) {
	return /* @__PURE__ */ jsxs("button", {
		type: "button",
		"aria-current": active ? "page" : void 0,
		"aria-label": session.title,
		"data-menu-row": true,
		onClick: onSelect,
		style: { paddingLeft: 8 + indent * 12 },
		className: cn("relative z-10 flex w-full items-center gap-1.5 rounded-[7px] py-1.5 pr-2 text-left transition-[color,transform] duration-150 active:scale-[0.96]", active ? "font-medium text-ink" : "text-ink-2"),
		children: [
			active ? /* @__PURE__ */ jsx("span", {
				"aria-hidden": "true",
				className: "absolute bottom-1.5 left-1 top-1.5 w-[2px] rounded-full bg-accent"
			}) : null,
			/* @__PURE__ */ jsx(StateGlyph, { state: session.state }),
			/* @__PURE__ */ jsx("span", {
				className: "min-w-0 flex-1 truncate text-[13px]",
				children: session.title
			}),
			session.meta ? /* @__PURE__ */ jsx("span", {
				className: "shrink-0 text-[10.5px] tabular-nums text-ink-3",
				children: session.meta
			}) : null
		]
	});
}
function PanelSidebar({ activeSessionId, onSelectSession, onNewTask, onOpenSettings, onToggleTheme, theme, settingsOpen, workspaceName, workspacePath, workspaceSessions, historyGroups }) {
	const [query, setQuery] = useState("");
	const [workspaceOpen, setWorkspaceOpen] = useState(true);
	const [collapsed, setCollapsed] = useState(/* @__PURE__ */ new Set());
	const q = query.trim().toLowerCase();
	const shownWorkspace = useMemo(() => q ? workspaceSessions.filter((s) => s.title.toLowerCase().includes(q)) : workspaceSessions, [workspaceSessions, q]);
	const shownHistory = useMemo(() => historyGroups.map((g) => ({
		...g,
		rows: q ? g.rows.filter((s) => s.title.toLowerCase().includes(q)) : g.rows
	})).filter((g) => g.rows.length > 0), [historyGroups, q]);
	function toggleGroup(label) {
		setCollapsed((current) => {
			const next = new Set(current);
			if (next.has(label)) next.delete(label);
			else next.add(label);
			return next;
		});
	}
	const nothingToShow = shownWorkspace.length === 0 && shownHistory.length === 0;
	return /* @__PURE__ */ jsxs("aside", {
		className: "flex h-full w-60 min-w-60 shrink-0 flex-col border-r border-line bg-surface",
		children: [/* @__PURE__ */ jsxs("nav", {
			"aria-label": "工作区会话",
			className: "flex min-h-0 flex-1 flex-col overflow-y-auto p-2",
			children: [
				/* @__PURE__ */ jsxs("button", {
					type: "button",
					title: workspacePath,
					className: "mb-1.5 flex w-full items-center gap-2.5 rounded-control p-1.5 text-left transition-[background-color,transform] duration-100 hover:bg-hover active:scale-[0.96]",
					children: [
						/* @__PURE__ */ jsx("span", {
							className: "flex size-8 shrink-0 items-center justify-center rounded-[8px] bg-ink text-[13px] font-semibold text-surface",
							children: "W"
						}),
						/* @__PURE__ */ jsxs("span", {
							className: "min-w-0 flex-1",
							children: [/* @__PURE__ */ jsx("span", {
								className: "block truncate text-[13px] font-medium leading-tight text-ink",
								children: workspaceName
							}), /* @__PURE__ */ jsx("span", {
								className: "block truncate font-mono text-[10.5px] leading-tight text-ink-3",
								children: workspacePath
							})]
						}),
						/* @__PURE__ */ jsx(ChevronDown, {
							"aria-hidden": "true",
							className: "shrink-0 rotate-180 text-ink-3",
							size: 12,
							strokeWidth: 2
						})
					]
				}),
				/* @__PURE__ */ jsxs("label", {
					className: "mb-1 flex h-8 items-center gap-2 rounded-control bg-inset px-2.5 shadow-hairline",
					children: [
						/* @__PURE__ */ jsx(Search, {
							"aria-hidden": "true",
							className: "shrink-0 text-ink-3",
							size: 12,
							strokeWidth: 2
						}),
						/* @__PURE__ */ jsx("input", {
							"aria-label": "搜索会话",
							onChange: (event) => setQuery(event.target.value),
							placeholder: "搜索会话",
							value: query,
							className: "min-w-0 flex-1 bg-transparent text-[12.5px] text-ink outline-none placeholder:text-ink-3"
						}),
						query ? /* @__PURE__ */ jsx("button", {
							type: "button",
							"aria-label": "清空搜索",
							onClick: () => setQuery(""),
							className: "flex size-4.5 shrink-0 items-center justify-center rounded-[5px] text-ink-3 transition-colors duration-100 hover:bg-hover hover:text-ink",
							children: /* @__PURE__ */ jsx(X, {
								"aria-hidden": "true",
								size: 11,
								strokeWidth: 2.4
							})
						}) : /* @__PURE__ */ jsx("kbd", {
							className: "flex size-4.5 shrink-0 items-center justify-center rounded-[5px] bg-surface text-[10px] text-ink-3 shadow-hairline",
							children: "/"
						})
					]
				}),
				/* @__PURE__ */ jsxs("button", {
					type: "button",
					onClick: onNewTask,
					className: "mb-1 flex w-full items-center gap-2 rounded-control px-2 py-1.5 text-[13px] font-medium text-accent transition-[background-color,transform] duration-100 hover:bg-accent-tint active:scale-[0.96]",
					children: [/* @__PURE__ */ jsx("span", {
						className: "min-w-0 flex-1 truncate text-left",
						children: "新建任务"
					}), /* @__PURE__ */ jsx("span", {
						className: "flex size-4 shrink-0 items-center justify-center rounded-full bg-accent text-white",
						children: /* @__PURE__ */ jsx(Plus, {
							"aria-hidden": "true",
							size: 9,
							strokeWidth: 3
						})
					})]
				}),
				/* @__PURE__ */ jsxs(GlideMenu, {
					className: "flex flex-col gap-0.5",
					highlightClassName: "inset-x-0 rounded-[7px] bg-hover",
					children: [/* @__PURE__ */ jsxs("div", { children: [
						/* @__PURE__ */ jsx(SectionHeader, { label: "工作区" }),
						/* @__PURE__ */ jsxs("button", {
							type: "button",
							"aria-expanded": workspaceOpen,
							"aria-label": workspaceName,
							"data-menu-row": true,
							onClick: () => setWorkspaceOpen((v) => !v),
							className: "relative z-10 flex w-full items-center gap-1.5 rounded-[7px] px-2 py-1.5 text-left transition-[color,transform] duration-150 active:scale-[0.96]",
							children: [
								/* @__PURE__ */ jsx(ChevronRight, {
									"aria-hidden": "true",
									className: cn("shrink-0 text-ink-3 transition-transform duration-150", workspaceOpen && "rotate-90"),
									size: 10,
									strokeWidth: 2.2
								}),
								/* @__PURE__ */ jsx(Folder, {
									"aria-hidden": "true",
									className: "shrink-0 text-ink-2",
									size: 13,
									strokeWidth: 1.8
								}),
								/* @__PURE__ */ jsx("span", {
									className: "min-w-0 flex-1 truncate text-[13px] text-ink-2",
									children: workspaceName
								}),
								shownWorkspace.length > 0 ? /* @__PURE__ */ jsx("span", {
									className: "shrink-0 text-[10.5px] tabular-nums text-ink-3",
									children: shownWorkspace.length
								}) : null
							]
						}),
						workspaceOpen && shownWorkspace.length > 0 ? /* @__PURE__ */ jsx("div", {
							className: "flex flex-col gap-px",
							children: shownWorkspace.map((session) => /* @__PURE__ */ jsx(SessionRow, {
								active: session.id === activeSessionId,
								indent: 1,
								onSelect: () => onSelectSession(session.id),
								session
							}, session.id))
						}) : null,
						workspaceOpen && shownWorkspace.length === 0 ? /* @__PURE__ */ jsx("p", {
							className: "py-1 pl-6 pr-2 text-[11.5px] leading-[1.6] text-ink-3",
							children: q ? "没有匹配的会话。" : "这个工作区下还没有会话。"
						}) : null
					] }), shownHistory.map((group) => {
						const shut = collapsed.has(group.label);
						return /* @__PURE__ */ jsxs(Fragment, { children: [
							/* @__PURE__ */ jsx(SectionHeader, { label: group.label }),
							/* @__PURE__ */ jsxs("button", {
								type: "button",
								"aria-expanded": !shut,
								onClick: () => toggleGroup(group.label),
								className: "flex w-full items-center gap-1 px-2 pb-0.5 text-[10.5px] font-medium tracking-[0.06em] text-ink-3",
								children: [
									/* @__PURE__ */ jsx(ChevronDown, {
										"aria-hidden": "true",
										className: cn("shrink-0 transition-transform duration-150", shut && "-rotate-90"),
										size: 10,
										strokeWidth: 2.2
									}),
									/* @__PURE__ */ jsx("span", {
										className: "min-w-0 flex-1 truncate text-left",
										children: group.label
									}),
									/* @__PURE__ */ jsx("span", {
										className: "shrink-0 tabular-nums",
										children: group.rows.length
									})
								]
							}),
							!shut && /* @__PURE__ */ jsx("div", {
								className: "flex flex-col gap-px",
								children: group.rows.map((session) => /* @__PURE__ */ jsx(SessionRow, {
									active: session.id === activeSessionId,
									indent: 1,
									onSelect: () => onSelectSession(session.id),
									session
								}, session.id))
							})
						] }, group.label);
					})]
				}),
				nothingToShow && q ? /* @__PURE__ */ jsxs("p", {
					className: "px-3 py-2 text-[12px] leading-[1.6] text-ink-3",
					children: [
						"没有匹配「",
						query.trim(),
						"」的会话。"
					]
				}) : null
			]
		}), /* @__PURE__ */ jsxs(GlideMenu, {
			className: "flex flex-col gap-1 border-t border-line p-2",
			highlightClassName: "inset-x-2 rounded-[7px] bg-hover",
			children: [/* @__PURE__ */ jsxs("button", {
				type: "button",
				"aria-label": "设置",
				"aria-current": settingsOpen ? "page" : void 0,
				"data-menu-row": true,
				onClick: onOpenSettings,
				className: cn("relative z-10 flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left text-[12.5px] transition-colors duration-150", settingsOpen ? "font-medium text-ink" : "text-ink-2 hover:text-ink"),
				children: [
					settingsOpen ? /* @__PURE__ */ jsx("span", {
						"aria-hidden": "true",
						className: "absolute bottom-1.5 left-0.5 top-1.5 w-[2px] rounded-full bg-accent"
					}) : null,
					/* @__PURE__ */ jsx(Settings, {
						"aria-hidden": "true",
						size: 13,
						strokeWidth: 1.8
					}),
					"设置"
				]
			}), /* @__PURE__ */ jsxs("button", {
				type: "button",
				"data-menu-row": true,
				onClick: onToggleTheme,
				className: "relative z-10 flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left text-[12.5px] text-ink-2 transition-colors duration-150 hover:text-ink",
				children: [theme === "dark" ? /* @__PURE__ */ jsx(Sun, {
					"aria-hidden": "true",
					size: 13,
					strokeWidth: 1.8
				}) : /* @__PURE__ */ jsx(Moon, {
					"aria-hidden": "true",
					size: 13,
					strokeWidth: 1.8
				}), theme === "dark" ? "浅色" : "暗色"]
			})]
		})]
	});
}
//#endregion
//#region src/lib/panel-views.ts
var PANEL_VIEWS = [
	{
		id: "chat",
		label: "对话",
		icon: "message-square-text",
		demoIcon: "message-square-text",
		fed: true
	},
	{
		id: "approval",
		label: "审批",
		icon: "shield-check",
		demoIcon: "shield-check",
		fed: true
	},
	{
		id: "palette",
		label: "命令",
		icon: "command",
		demoIcon: "command",
		fed: false
	},
	{
		id: "tasks",
		label: "任务",
		icon: "list-checks",
		demoIcon: "list-checks",
		fed: false
	},
	{
		id: "ball",
		label: "球状态",
		icon: "orbit",
		demoIcon: "orbit",
		fed: false
	},
	{
		id: "config",
		label: "设置",
		icon: "settings",
		demoIcon: "settings",
		fed: false,
		selfFed: true
	},
	{
		id: "security",
		label: "安全",
		icon: "shield-alert",
		demoIcon: "shield-alert",
		fed: false
	},
	{
		id: "privacy",
		label: "隐私",
		icon: "lock",
		demoIcon: "lock",
		fed: false
	},
	{
		id: "cost",
		label: "成本",
		icon: "chart-column",
		demoIcon: "chart-column",
		fed: false
	}
];
/** The screen the panel opens on when the snapshot names none. Chosen because
it is the only row with anything behind it today; it is a fallback for a
missing field, not a memory of where the user last was. */
var DEFAULT_VIEW = "chat";
/**
* The one honest answer to "which screen is this". The field does not exist on
* PanelSnapshot yet - ticket 35's pump owns it, and `panel.resync` (the push
* that would re-establish it on every show) is still zero-hit repo-wide. Until
* then a screen change is asked for and refused, so what is on screen stays
* whatever the caller passed.
*
* Reading it through this one function is what keeps the swap a one-line change
* later: no component holds a view, and none derives one either.
*/
function currentView(snapshot) {
	const named = snapshot.view;
	if (named === void 0) return DEFAULT_VIEW;
	if (PANEL_VIEWS.some((v) => v.id === named)) return named;
	console.error(`panel-views: snapshot names unknown view ${JSON.stringify(named)}`);
	return DEFAULT_VIEW;
}
//#endregion
//#region src/App.tsx
var EMPTY_COMPOSER = {
	mode: {
		current: "unknown",
		names: [],
		l2ConfirmNames: []
	},
	workspace: {
		set: false,
		spelling: "",
		canonical: "",
		reparse: false,
		rewritten: false,
		reason: "尚未收到原生侧的状态快照"
	},
	attachments: [],
	acceptedAttachmentMimes: [],
	maxAttachmentBytes: 0,
	attachmentError: ""
};
var EMPTY = {
	pending: [],
	results: [],
	composer: EMPTY_COMPOSER,
	generatedAt: ""
};
/** Native mode names -> the Chinese the tier row shows. Spelled once, in the
component that draws them, so a fourth tier cannot be invented here. */
var TIER_NAMES = {
	ask_every_step: "每步都问",
	ask_high_risk: "高风险才问",
	auto_approve: "全自动",
	unknown: "档位未知"
};
/** The task monitor's roster rows, straight off the wire (ticket 197/220: the
fields are copied verbatim - blockedOnApproval is the host's own join,
status is a D43 name this file neither judges nor fills in, and the row's
stream triple rides as-is (the popover draws it as a suffix, never a
status); the display words live in the popover that draws them). Until
this derivation existed the monitor was fed a bare empty array and every
pushed row died unseen. An absent tasks key stays distinguishable from an
empty roster: the popover is told which one it got (lib/panel.ts's
optionality rule). */
function monitorAgents(tasks) {
	return (tasks?.rows ?? []).map((row) => ({
		id: row.taskId,
		name: row.label || row.taskId,
		blockedOnApproval: row.blockedOnApproval,
		status: row.status,
		statusKnown: row.statusKnown,
		statusReason: row.statusReason,
		streamTruncated: row.streamTruncated,
		streamElidedRunes: row.streamElidedRunes,
		streamDropped: row.streamDropped
	}));
}
/** Narrower than this and the sidebar stops fitting a session title, so it
collapses to its icon rail. 640px is the DEFAULT panel width, not an edge
case - see the header. */
var SIDEBAR_COLLAPSE_BELOW = 860;
function App({ snapshot = EMPTY }) {
	const composer = snapshot.composer ?? EMPTY_COMPOSER;
	const [local, setLocal] = useState({
		view: null,
		railOpen: false,
		settingsOpen: false,
		settingsSection: "通用",
		dark: false,
		monitorOpen: false,
		paletteOpen: false,
		sessionId: null,
		activeTurn: null
	});
	const view = snapshot.view !== void 0 ? currentView(snapshot) : local.view ?? currentView(snapshot);
	const categoryOf = {
		config: "外观",
		security: "安全",
		privacy: "隐私",
		cost: "成本"
	}[view];
	const narrow = useNarrow();
	const [sidebarOpen, setSidebarOpen] = useState(false);
	useEffect(() => {
		setSidebarOpen(false);
	}, [narrow]);
	const sidebarVisible = !narrow || sidebarOpen;
	const monitorRef = useRef(null);
	useDismiss(monitorRef, local.monitorOpen, useCallback(() => setLocal((c) => ({
		...c,
		monitorOpen: false
	})), []));
	useEffect(() => {
		function onKey(event) {
			if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k") {
				event.preventDefault();
				setLocal((c) => ({
					...c,
					paletteOpen: !c.paletteOpen
				}));
				return;
			}
			if (event.key === "Escape") {
				setLocal((c) => ({
					...c,
					paletteOpen: false,
					monitorOpen: false
				}));
				if (narrow) setSidebarOpen(false);
			}
		}
		document.addEventListener("keydown", onKey);
		return () => document.removeEventListener("keydown", onKey);
	}, [narrow]);
	function toggleTheme() {
		const next = !local.dark;
		setLocal((c) => ({
			...c,
			dark: next
		}));
		document.documentElement.dataset.theme = next ? "dark" : "light";
	}
	const tiers = composer.mode.names.length > 0 ? composer.mode.names : ["unknown"];
	return /* @__PURE__ */ jsxs("div", {
		className: "flex h-screen w-full flex-col overflow-hidden bg-page text-ink",
		children: [
			/* @__PURE__ */ jsxs("header", {
				className: "flex h-11 shrink-0 select-none items-center gap-2 border-b border-line bg-surface pl-4 pr-2.5",
				children: [
					/* @__PURE__ */ jsxs("div", {
						className: "flex min-w-0 shrink-0 items-center gap-2.5",
						children: [
							narrow ? /* @__PURE__ */ jsx("button", {
								type: "button",
								"aria-label": sidebarOpen ? "收起会话树" : "展开会话树",
								"aria-expanded": sidebarOpen,
								onClick: () => setSidebarOpen((v) => !v),
								className: "flex size-7 shrink-0 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink",
								children: /* @__PURE__ */ jsx(PanelLeft, {
									"aria-hidden": "true",
									size: 15
								})
							}) : null,
							/* @__PURE__ */ jsx("span", {
								"aria-hidden": "true",
								className: "size-2.5 rounded-full bg-accent"
							}),
							/* @__PURE__ */ jsx("span", {
								className: "shrink-0 text-[13px] font-medium text-ink",
								children: "一缕"
							}),
							snapshot.pending.length > 0 ? /* @__PURE__ */ jsxs("span", {
								className: "shrink-0 rounded-full bg-warn-soft px-2 py-0.5 text-[10.5px] font-medium tabular-nums text-warn",
								children: ["待审批 ", snapshot.pending.length]
							}) : null
						]
					}),
					/* @__PURE__ */ jsx("div", {
						className: "flex min-w-0 flex-1 items-center gap-1.5",
						children: /* @__PURE__ */ jsx(GitBranchChip, {
							view: {
								current: "",
								isRepo: false,
								local: [],
								worktrees: []
							},
							onCheckout: () => void 0,
							onNewWorktree: () => void 0
						})
					}),
					/* @__PURE__ */ jsxs("div", {
						className: "flex shrink-0 items-center gap-0.5",
						children: [
							/* @__PURE__ */ jsxs("div", {
								className: "relative",
								ref: monitorRef,
								children: [/* @__PURE__ */ jsx("button", {
									type: "button",
									"aria-expanded": local.monitorOpen,
									"aria-label": "任务监控",
									onClick: () => setLocal((c) => ({
										...c,
										monitorOpen: !c.monitorOpen
									})),
									className: `flex h-7 size-7 shrink-0 items-center justify-center rounded-control transition-colors duration-150 ${local.monitorOpen ? "bg-hover text-ink" : "text-ink-2 hover:bg-hover hover:text-ink"}`,
									children: /* @__PURE__ */ jsx(Activity, {
										"aria-hidden": "true",
										size: 14
									})
								}), local.monitorOpen ? /* @__PURE__ */ jsx("div", {
									className: "absolute right-0 top-full z-40 mt-1.5",
									children: /* @__PURE__ */ jsx(MonitorPopover, {
										open: true,
										env: [{
											label: "工作区",
											value: composer.workspace.set ? composer.workspace.canonical : "尚未收到"
										}, {
											label: "分支",
											value: "尚未读到"
										}],
										agents: monitorAgents(snapshot.tasks),
										rosterPresent: snapshot.tasks !== void 0,
										poolInFlight: snapshot.tasks?.inFlightSlots,
										poolCap: snapshot.tasks?.poolCap,
										streamTruncated: snapshot.tasks?.streamTruncated,
										streamElidedRunes: snapshot.tasks?.streamElidedRunes,
										streamDroppedCount: snapshot.tasks?.droppedCount,
										onPin: () => void 0
									})
								}) : null]
							}),
							/* @__PURE__ */ jsx("kbd", {
								className: "kbd hidden sm:flex",
								children: "Ctrl K"
							}),
							/* @__PURE__ */ jsx("button", {
								type: "button",
								"aria-label": "命令面板",
								className: "flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink",
								onClick: () => setLocal((c) => ({
									...c,
									paletteOpen: !c.paletteOpen
								})),
								title: "命令面板",
								children: /* @__PURE__ */ jsx(Command, {
									"aria-hidden": "true",
									size: 14,
									strokeWidth: 1.9
								})
							}),
							/* @__PURE__ */ jsx("button", {
								type: "button",
								"aria-label": local.railOpen ? "收起右侧栏" : "展开右侧栏",
								className: "flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink",
								onClick: () => setLocal((c) => ({
									...c,
									railOpen: !c.railOpen
								})),
								title: local.railOpen ? "收起右侧栏" : "展开右侧栏",
								children: local.railOpen ? /* @__PURE__ */ jsx(PanelRightClose, {
									"aria-hidden": "true",
									size: 15
								}) : /* @__PURE__ */ jsx(PanelRight, {
									"aria-hidden": "true",
									size: 15
								})
							}),
							/* @__PURE__ */ jsx("button", {
								type: "button",
								"aria-label": "切换主题",
								className: "flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink",
								onClick: toggleTheme,
								title: "切换主题",
								children: local.dark ? /* @__PURE__ */ jsx(Sun, {
									"aria-hidden": "true",
									size: 15
								}) : /* @__PURE__ */ jsx(Moon, {
									"aria-hidden": "true",
									size: 15
								})
							})
						]
					})
				]
			}),
			/* @__PURE__ */ jsxs("div", {
				className: "relative flex min-h-0 flex-1 overflow-hidden",
				children: [
					sidebarVisible ? /* @__PURE__ */ jsx("div", {
						className: narrow ? "absolute inset-y-0 left-0 z-30 shadow-overlay" : "contents",
						children: /* @__PURE__ */ jsx(PanelSidebar, {
							activeSessionId: local.sessionId,
							historyGroups: [],
							onNewTask: () => setLocal((c) => ({
								...c,
								sessionId: null,
								view: "chat"
							})),
							onOpenSettings: () => setLocal((c) => ({
								...c,
								settingsOpen: true
							})),
							onSelectSession: (id) => setLocal((c) => ({
								...c,
								sessionId: id,
								view: "chat"
							})),
							onToggleTheme: toggleTheme,
							settingsOpen: local.settingsOpen,
							theme: local.dark ? "dark" : "light",
							workspaceName: "一缕",
							workspacePath: composer.workspace.set ? composer.workspace.canonical : "尚未收到工作区",
							workspaceSessions: []
						})
					}) : null,
					/* @__PURE__ */ jsx("div", {
						className: "flex min-w-0 flex-1 flex-col",
						children: local.settingsOpen || categoryOf ? /* @__PURE__ */ jsx(SettingsPage, {
							activeSection: categoryOf ?? local.settingsSection,
							onClose: () => setLocal((c) => ({
								...c,
								settingsOpen: false
							})),
							onSelectSection: (next) => setLocal((c) => ({
								...c,
								settingsSection: next
							})),
							ballStyle: snapshot.ballStyle ?? null,
							ballSize: snapshot.ballSize ?? null,
							ballOpacityIdle: snapshot.ballOpacityIdle ?? null,
							ballClickThrough: snapshot.ballClickThrough ?? null,
							panelWidth: snapshot.panelWidth ?? null,
							panelFontSize: snapshot.panelFontSize ?? null
						}) : /* @__PURE__ */ jsxs(Fragment$1, { children: [snapshot.pending.map((v) => /* @__PURE__ */ jsx(L2ApprovalCard, { view: v }, v.correlationId)), /* @__PURE__ */ jsx("div", {
							className: "flex min-h-0 flex-1 flex-col",
							children: view === "approval" || view === "tasks" || view === "ball" ? /* @__PURE__ */ jsx(RetiredView, { id: view }) : /* @__PURE__ */ jsx(Conversation, {
								composer,
								onJump: (id) => setLocal((c) => ({
									...c,
									activeTurn: id
								})),
								permission: composer.mode.current,
								snapshot,
								tiers
							})
						})] })
					}),
					local.railOpen ? /* @__PURE__ */ jsx("div", {
						className: narrow ? "absolute inset-y-0 right-0 z-30 shadow-overlay" : "",
						children: /* @__PURE__ */ jsx(RightRail, {
							reviewFiles: [],
							fileTree: [],
							terminalLines: [],
							browserUrl: ""
						})
					}) : null
				]
			}),
			local.paletteOpen ? /* @__PURE__ */ jsx("div", {
				"aria-label": "命令面板",
				"aria-modal": "true",
				className: "fixed inset-0 z-40 flex items-start justify-center bg-black/35 pt-[12vh]",
				onClick: (e) => {
					if (e.target === e.currentTarget) setLocal((c) => ({
						...c,
						paletteOpen: false
					}));
				},
				role: "dialog",
				children: /* @__PURE__ */ jsx("div", {
					className: "w-[560px] overflow-hidden rounded-card border border-line bg-surface shadow-overlay",
					style: { animation: "fade-up 200ms var(--ease-out-strong) both" },
					children: /* @__PURE__ */ jsx(PaletteScreen, {
						groups: [],
						onQueryChange: () => void 0,
						onSelect: () => setLocal((c) => ({
							...c,
							paletteOpen: false
						})),
						query: ""
					})
				})
			}) : null
		]
	});
}
/** Track the window width. A resize listener, not a media query: the threshold
is about whether a session title fits beside the icons, not about a device
class. 640px is the DEFAULT panel width (internal/config/schema.go), so
this path is the common one, not the edge case. */
function useNarrow() {
	const read = useCallback(() => typeof window !== "undefined" && window.innerWidth < SIDEBAR_COLLAPSE_BELOW, []);
	const [narrow, setNarrow] = useState(read);
	useEffect(() => {
		const onResize = () => setNarrow(read());
		window.addEventListener("resize", onResize);
		onResize();
		return () => window.removeEventListener("resize", onResize);
	}, [read]);
	return narrow;
}
function Conversation({ snapshot, tiers, permission, composer, onJump }) {
	const turns = snapshot.results.map((chunk) => ({
		id: chunk.correlationId,
		label: chunk.text.slice(0, 18) || "（空）"
	}));
	return /* @__PURE__ */ jsxs(Fragment$1, { children: [/* @__PURE__ */ jsxs("div", {
		className: "flex min-h-0 flex-1 overflow-hidden",
		children: [turns.length > 1 ? /* @__PURE__ */ jsx("div", {
			className: "flex w-44 shrink-0 items-center border-r border-line px-2",
			children: /* @__PURE__ */ jsx(LineSidebar, {
				fontSize: .8,
				itemGap: 10,
				items: turns.map((t) => t.label),
				markerLength: 26,
				maxShift: 14,
				onItemClick: (index) => onJump(turns[index].id),
				showIndex: true
			})
		}) : null, /* @__PURE__ */ jsx("div", {
			className: "flex min-h-0 flex-1 flex-col overflow-y-auto",
			children: /* @__PURE__ */ jsx("div", {
				className: "mx-auto flex w-full max-w-[720px] flex-1 flex-col px-6",
				children: snapshot.results.length === 0 ? /* @__PURE__ */ jsx(ConversationPlaceholder, {}) : /* @__PURE__ */ jsx("div", {
					className: "flex flex-col gap-5 py-5",
					children: snapshot.results.map((chunk) => /* @__PURE__ */ jsx("p", {
						className: "text-[13px] leading-relaxed text-ink",
						children: chunk.text
					}, chunk.correlationId))
				})
			})
		})]
	}), /* @__PURE__ */ jsx("div", {
		className: "shrink-0 border-t border-line bg-surface px-6 pb-4 pt-2.5",
		children: /* @__PURE__ */ jsx("div", {
			className: "mx-auto w-full max-w-[720px]",
			children: /* @__PURE__ */ jsx(ComposerRow, {
				composer,
				modelSwitch: snapshot.modelSwitch,
				permission,
				permissions: tiers.map((t) => ({
					key: t,
					name: TIER_NAMES[t] ?? t,
					gated: composer.mode.l2ConfirmNames.includes(t)
				}))
			})
		})
	})] });
}
function RetiredView({ id }) {
	return /* @__PURE__ */ jsx("div", {
		className: "flex min-h-0 flex-1 items-start justify-center px-6 py-10",
		children: /* @__PURE__ */ jsx("div", {
			className: "w-full max-w-[520px]",
			style: { animation: "fade-up 320ms var(--ease-out-strong) both" },
			children: /* @__PURE__ */ jsx("p", {
				className: "text-[13px] leading-[1.8] text-ink-2",
				children: {
					approval: "审批不是一屏，是一个卡片：真有高危请求时它会自己出现在面板顶部，右上角的计数告诉你有几个在等你。",
					tasks: "任务不是一个地方，是一个浮层：点顶栏的「任务监控」看正在跑的活。",
					ball: "悬浮球是原生画的（D29），不在面板里。球的设置在设置页。"
				}[id] ?? "这一屏已经不存在了。"
			})
		})
	});
}
var ARRIVES = [
	"流式回复",
	"工具调用",
	"思考过程"
];
function ConversationPlaceholder() {
	return /* @__PURE__ */ jsx("div", {
		className: "flex flex-1 items-start justify-center pt-[13vh]",
		children: /* @__PURE__ */ jsxs("div", {
			className: "w-full max-w-[440px] text-center",
			style: { animation: "fade-up 320ms var(--ease-out-strong) both" },
			children: [
				/* @__PURE__ */ jsx("p", {
					className: "text-[15px] font-medium text-ink",
					children: "一缕在这儿"
				}),
				/* @__PURE__ */ jsx("p", {
					className: "mx-auto mt-1.5 max-w-[360px] text-[12.5px] leading-[1.7] text-ink-3",
					children: "面板不自己造内容。在下面说一句话，或者从悬浮球开口，原生侧会把结果一条条推上来。"
				}),
				/* @__PURE__ */ jsx("div", {
					className: "mt-4 flex items-center justify-center gap-1.5",
					children: ARRIVES.map((label) => /* @__PURE__ */ jsx("span", {
						className: "rounded-full border border-line bg-surface px-2.5 py-1 text-[11px] text-ink-3",
						children: label
					}, label))
				}),
				/* @__PURE__ */ jsx("p", {
					className: "mt-3.5 text-[10.5px] text-ink-3",
					children: "还没收到任何一条"
				})
			]
		})
	});
}
//#endregion
//#region scripts/render-task-roster.tsx
var failures = [];
function check(label, ok, detail = "") {
	if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}
var EMOJI_RE = /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;
/** Same escaping React applies to a text node - so a match here is a match on
screen, not a match on an unescaped substring that never renders. */
function esc(text) {
	return text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#x27;");
}
var UNFILED_REASON = "宿主没有给出这一维的读数（fail-closed：不替任务编一个状态）";
function row(over) {
	return {
		kind: "subagent",
		parentTaskId: "root-1",
		status: "",
		statusKnown: false,
		streamKey: over.taskId,
		blockedOnApproval: false,
		streamTruncated: false,
		streamElidedRunes: 0,
		streamDropped: false,
		...over
	};
}
function roster(rows) {
	return {
		rows,
		inFlightSlots: 1,
		poolCap: 8,
		streamTruncated: false,
		streamElidedRunes: 0,
		droppedStreamKeys: [],
		droppedCount: 0
	};
}
var COMPOSER = {
	mode: {
		current: "ask_high_risk",
		names: [
			"ask_every_step",
			"ask_high_risk",
			"auto_approve"
		],
		l2ConfirmNames: ["auto_approve"]
	},
	workspace: {
		set: true,
		spelling: "D:/work/wisp",
		canonical: "D:/work/wisp",
		reparse: false,
		rewritten: false,
		reason: ""
	},
	attachments: [],
	acceptedAttachmentMimes: [],
	maxAttachmentBytes: 0,
	attachmentError: ""
};
function snapshotWith(tasks) {
	return {
		pending: [],
		results: [],
		composer: COMPOSER,
		generatedAt: "2026-11-08T14:00:00Z",
		tasks
	};
}
var ENV = [{
	label: "工作区",
	value: "D:/work/wisp"
}, {
	label: "分支",
	value: "dev"
}];
check("an absent tasks key must yield zero rows without inventing any", monitorAgents(void 0).length === 0);
var BLOCKED_ROW = row({
	taskId: "child-9",
	label: "等批准的子任务",
	status: "Acting",
	statusKnown: true,
	blockedOnApproval: true
});
var fromPresent = monitorAgents(roster([BLOCKED_ROW]));
check("a present tasks key must yield one derived row per wire row", fromPresent.length === 1, `rows=${fromPresent.length}`);
check("the derived row must carry blockedOnApproval verbatim from the wire", fromPresent[0]?.blockedOnApproval === true);
check("the derived row must carry the D43 status verbatim (this layer never judges it)", fromPresent[0]?.status === "Acting" && fromPresent[0]?.statusKnown === true);
check("the derived row must carry the host's id and label verbatim", fromPresent[0]?.id === "child-9" && fromPresent[0]?.name === "等批准的子任务");
var blockedHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: monitorAgents(roster([BLOCKED_ROW])),
	rosterPresent: true
}));
check("the blocked row's own label must reach the screen", blockedHtml.includes(esc("等批准的子任务")));
check("the blocked row must show 等批准", blockedHtml.includes(esc("等批准")));
check("the blocked row must NOT show the D43 status word (等批准 outranks it)", !blockedHtml.includes(esc("执行中")));
check("the blocked row's dot must be the warn tier", blockedHtml.includes("bg-warn"));
check("the blocked markup must hold zero emoji", !EMOJI_RE.test(blockedHtml));
var freeRow = row({
	taskId: "child-4",
	label: "没被卡住的子任务",
	status: "Acting",
	statusKnown: true,
	blockedOnApproval: false
});
var freeHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: monitorAgents(roster([freeRow])),
	rosterPresent: true
}));
check("the free row must show the D43 status word", freeHtml.includes(esc("执行中")));
check("the free row must NOT show 等批准", !freeHtml.includes(esc("等批准")));
check("the free row's dot must be the live-work tier", freeHtml.includes("bg-accent"));
check("the free markup must hold zero emoji", !EMOJI_RE.test(freeHtml));
var unknownRow = row({
	taskId: "child-7",
	label: "状态没登记的子任务",
	status: "",
	statusKnown: false,
	blockedOnApproval: false,
	statusReason: UNFILED_REASON
});
var unknownHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: monitorAgents(roster([unknownRow])),
	rosterPresent: true
}));
check("an unfiled status must show the HOST's own sentence, never a guessed word", unknownHtml.includes(esc(UNFILED_REASON)), unknownHtml.slice(unknownHtml.indexOf("状态没登记") - 10, unknownHtml.indexOf("状态没登记") + 220));
check("an unfiled status must not render any D43 word", ![
	"执行中",
	"休眠",
	"收尾中",
	"出错"
].some((w) => unknownHtml.includes(esc(w))));
check("the unfiled-status dot must be the idle tier", unknownHtml.includes("bg-ink-3"));
var emptyRosterHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: [],
	rosterPresent: false
}));
check("with no tasks key the popover must say the roster never arrived", emptyRosterHtml.includes(esc("宿主没带任务这一节")));
check("the no-roster markup must hold zero emoji", !EMOJI_RE.test(emptyRosterHtml));
var spawnedEmptyHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: [],
	rosterPresent: true
}));
check("with a present-but-empty roster the popover must say nothing is spawned (a different sentence)", spawnedEmptyHtml.includes(esc("还没有子任务在跑")) && !spawnedEmptyHtml.includes(esc("宿主没带任务这一节")));
var appHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(App, { snapshot: snapshotWith(roster([BLOCKED_ROW])) }));
check("App with a pushed tasks section must render without crashing", appHtml.includes("一缕"));
check("the App markup must hold zero emoji", !EMOJI_RE.test(appHtml));
check("App without a tasks section must render without crashing", renderToStaticMarkup(/* @__PURE__ */ jsx(App, { snapshot: snapshotWith(void 0) })).includes("一缕"));
var appSrc = readFileSync("src/App.tsx", "utf8");
check("the monitor must be fed the derivation over snapshot.tasks, never a bare empty array", appSrc.includes("agents={monitorAgents(snapshot.tasks)}") && !appSrc.includes("agents={[]}"));
check("the popover must be told whether the snapshot carried a tasks section at all", appSrc.includes("rosterPresent={snapshot.tasks !== undefined}"));
check("等批准 must be derived from the wire's blockedOnApproval, not from any status word", readFileSync("src/components/harness/monitor-popover.tsx", "utf8").includes("if (row.blockedOnApproval) return { label: \"等批准\""));
var pooledHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: [],
	rosterPresent: true,
	poolInFlight: 3,
	poolCap: 8
}));
check("the pool numbers must read 在跑 3 / 上限 8", pooledHtml.includes(esc("在跑 3 / 上限 8")));
check("the pool markup must hold zero emoji", !EMOJI_RE.test(pooledHtml));
var bothHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: [],
	rosterPresent: true,
	streamTruncated: true,
	streamElidedRunes: 421,
	streamDroppedCount: 2
}));
check("a truncated log with drops must name both losses in one line", bothHtml.includes(esc("流已截断：省略 421 字，已丢 2 条")));
check("the truncation markup must hold zero emoji", !EMOJI_RE.test(bothHtml));
var truncatedOnlyHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: [],
	rosterPresent: true,
	streamTruncated: true,
	streamElidedRunes: 421
}));
check("a truncated log must show the elision count", truncatedOnlyHtml.includes(esc("流已截断：省略 421 字")));
check("truncated-only must not claim dropped streams it does not have", !truncatedOnlyHtml.includes(esc("已丢")));
var droppedOnlyHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: [],
	rosterPresent: true,
	streamDroppedCount: 1
}));
check("dropped streams must warn even when nothing was elided", droppedOnlyHtml.includes(esc("已丢 1 条流")) && !droppedOnlyHtml.includes(esc("省略")));
var cleanBoundHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: [],
	rosterPresent: true,
	streamTruncated: false,
	streamElidedRunes: 0,
	streamDroppedCount: 0
}));
check("a clean log must render no truncation warning (no fake zero, no noise)", !cleanBoundHtml.includes(esc("流已截断")) && !cleanBoundHtml.includes(esc("已丢")) && !cleanBoundHtml.includes(esc("省略")));
var truncatedRowHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: monitorAgents(roster([row({
		taskId: "child-11",
		label: "流被折界的子任务",
		streamTruncated: true,
		streamElidedRunes: 900
	})])),
	rosterPresent: true
}));
check("a folded row must show its elided count beside the status word", truncatedRowHtml.includes(esc("省 900 字")) && truncatedRowHtml.includes(esc("流被折界的子任务")));
var droppedRowHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: monitorAgents(roster([row({
		taskId: "child-12",
		label: "流页被丢的子任务",
		streamDropped: true
	})])),
	rosterPresent: true
}));
check("a dropped row must say its page is gone, not a fake elision count", droppedRowHtml.includes(esc("流已丢")) && !droppedRowHtml.includes(esc("省 ")));
check("the row-suffix markup must hold zero emoji", !EMOJI_RE.test(droppedRowHtml));
check("a whole row must show no truncation suffix", !freeHtml.includes(esc("省 ")) && !freeHtml.includes(esc("流已丢")));
var noNumbersHtml = renderToStaticMarkup(/* @__PURE__ */ jsx(MonitorPopover, {
	open: true,
	env: ENV,
	agents: [],
	rosterPresent: false
}));
check("with no tasks key nothing pool-shaped may render", !noNumbersHtml.includes(esc("在跑")));
check("with no tasks key nothing truncation-shaped may render", !noNumbersHtml.includes(esc("流已截断")) && !noNumbersHtml.includes(esc("已丢")) && !noNumbersHtml.includes(esc("省 ")));
var truncatedSection = {
	...roster([BLOCKED_ROW]),
	streamTruncated: true,
	streamElidedRunes: 421,
	droppedCount: 2
};
check("App with a truncated roster must render without crashing", renderToStaticMarkup(/* @__PURE__ */ jsx(App, { snapshot: snapshotWith(truncatedSection) })).includes("一缕"));
check("the monitor must be fed the pool numbers straight off snapshot.tasks", appSrc.includes("poolInFlight={snapshot.tasks?.inFlightSlots}") && appSrc.includes("poolCap={snapshot.tasks?.poolCap}"));
check("the monitor must be fed the whole-log truncation triple straight off snapshot.tasks", appSrc.includes("streamTruncated={snapshot.tasks?.streamTruncated}") && appSrc.includes("streamElidedRunes={snapshot.tasks?.streamElidedRunes}") && appSrc.includes("streamDroppedCount={snapshot.tasks?.droppedCount}"));
check("the derivation must carry the row stream triple verbatim", appSrc.includes("streamTruncated: row.streamTruncated") && appSrc.includes("streamElidedRunes: row.streamElidedRunes") && appSrc.includes("streamDropped: row.streamDropped"));
if (failures.length > 0) {
	for (const f of failures) console.error("render-roster: FAIL " + f);
	console.error(`render-roster: FAIL, ${failures.length} check(s) failed`);
	process.exit(1);
}
console.error("render-roster: OK, the task monitor consumes snapshot.tasks (等批准 is its own tier from blockedOnApproval, D43 words verbatim, unfiled status shows the host's sentence, absent key is an honest empty state, pool numbers and stream truncation land on the page and warn only when a bound bit)");
//#endregion
export {};
