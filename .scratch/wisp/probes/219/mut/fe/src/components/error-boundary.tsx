/* ============================================================================
   ErrorBoundary - the panel's last line, so one throw is not a blank window.
   ----------------------------------------------------------------------------
   Why this exists: a render-time throw unmounts the whole React tree and
   leaves an empty WebView. There are real ones on this path - panel-views'
   viewOf() throws on an unknown view id, nav-rail throws on an icon name it
   has no component for, and any component reading a snapshot shape other than
   the one it declared throws inside render. Before this file, the first such
   throw left the user staring at an empty window with no way to tell whether
   the panel was broken or simply had no data.

   What it does NOT catch, stated plainly because it is the more likely one:
   an error boundary does not see throws from event handlers. Composer.send()
   and the L2 card's 拒绝 both call lib/panel.ts, which throws on purpose when
   no native host is attached - that is an event-handler throw, so it surfaces
   as an unhandled error and the tree survives. Those two call sites need their
   own try/catch and a visible message (App.tsx still has to pass onUserError,
   which nothing does today); this file is the other half, not a substitute.

   What it shows, and what it refuses to show:
     - a human sentence first, in the panel's own vocabulary;
     - the error's own message in mono, because an operator needs it and a
       user does not have to read it;
     - NO stack. A stack is a development artefact; the panel ships to a user.
     - a 重试 button that remounts the subtree, which is the one action that
       is always correct: whatever threw may be transient (a missing host, a
       half-initialised snapshot), and remounting re-runs the same read.

   It holds no state beyond "did something throw" and it never sends anything
   anywhere - the panel stays a second state holder-free (PLAN.md:1044). The
   error is NOT remembered anywhere: a WebView restart re-runs the same code,
   which is exactly the behaviour that makes a deterministic bug reproducible
   instead of a ghost.

   Scope note: it catches render/lifecycle throws below itself. It cannot catch
   an async rejection nobody awaited - that still reaches window.onerror and
   shows nothing. The retry button is there so the user is never stuck.
   ============================================================================ */

import { Component, type ErrorInfo, type ReactNode } from "react";
import { Button } from "@/components/ai-native/button";

interface Props {
  children: ReactNode;
  /** Short label for which part of the tree died; the panel has three roots
      (the product panel, the app harness, the component showcase) and "出错了"
      alone would not say which. */
  area: string;
}

interface State {
  failed: boolean;
  message: string;
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { failed: false, message: "" };

  static getDerivedStateFromError(error: unknown): State {
    return {
      failed: true,
      message: error instanceof Error ? error.message : String(error),
    };
  }

  componentDidCatch(error: Error, info: ErrorInfo): void {
    // Console only, and deliberately not the component stack: the panel has no
    // telemetry route and inventing one is not this file's call.
    console.error(`panel: ${this.props.area} 渲染时抛出`, error, info.componentStack);
  }

  private retry = (): void => {
    this.setState({ failed: false, message: "" });
  };

  render(): ReactNode {
    if (!this.state.failed) {
      return this.props.children;
    }
    return (
      <div className="flex h-full w-full flex-col items-center justify-center gap-3 bg-page px-6 text-center">
        <p className="text-[14px] font-semibold text-ink">{this.props.area}出了点问题</p>
        <p className="max-w-[420px] text-[12.5px] leading-[1.7] text-ink-2">
          这一块没有画出来。面板不会拿假内容补上它——你看到的空白就是真实状态。
          点重试会重新读一次宿主状态；如果每次都这样，多半是宿主那一侧还没把
          快照推过来。
        </p>
        <p className="max-w-[520px] break-words font-mono text-[11px] leading-[1.6] text-ink-3">
          {this.state.message}
        </p>
        <Button className="mt-1" onClick={this.retry} size="sm" variant="secondary">
          重试
        </Button>
      </div>
    );
  }
}
