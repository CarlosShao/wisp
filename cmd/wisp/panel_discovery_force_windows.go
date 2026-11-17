//go:build windows

package main

import (
	"context"
	"fmt"

	"github.com/CarlosShao/wisp/internal/panel"
)

// Readface wave 3 (2026-11-17, ticket 38 AC#3/AC#5): the late binds of the
// discovery/force doors and the launch exit. The panel chain is assembled
// before the task pipeline exists (resident_windows.go's ordering, the gap
// ticket 453's relay and the roster-cancel seam already crossed), so the
// doors leave their seams nil at construction and this call - on the first
// line where both objects exist - hands them the runtime's own machinery:
//
//	tools.list  -> agentRuntime.toolsRosterRows (the registry's own entries,
//	               the pill rendered by internal/risk's vocabulary);
//	tool.launch -> the resident source's SubmitPanelToolLaunch (compose
//	               against the registry, submit through the ONE task slot;
//	               the loop's dispatch is the bridge gate - no bypass);
//	rules.force -> agentRuntime.forceRulesView (the ONE resolver the prompt
//	               injection also consults);
//	the message door's force chip stage -> agentRuntime.stageForceOverride.
//
// A boot with no panel, no dispatch chain or no task source leaves the seams
// nil. That is the honest state of the CLI seam and of a refused pipeline:
// the doors then refuse by name instead of dropping.
func bindPanelDiscoveryForce(m *PanelManager, src *residentTaskSource) {
	if m == nil || m.disp == nil || src == nil || src.run == nil {
		return
	}
	if doors, ok := m.disp.Launch.(*discoveryDoors); ok && doors != nil {
		// One struct carries all three sockets; the two read doors are the
		// same object the Launch field holds (panel_inbound.go wires all
		// three to one discoveryDoors), so one bind covers the wave.
		doors.bind(
			src.run.toolsRosterRows,
			func(ctx context.Context, tool, args string) error {
				return src.SubmitPanelToolLaunch(ctx, tool, args)
			},
			src.run.forceRulesView,
		)
	}
	if door, ok := m.disp.Message.(*panelMessageDoor453); ok && door != nil {
		door.bindForceStageW3(src.run.stageForceOverride)
	}
}

// SubmitPanelToolLaunch is the resident source's launch exit (ticket 38
// AC#3): compose the task book against the registry the run actually
// assembled, then submit through the ONE task slot. A refused submission
// clears the staged force chip in the same call (the 453 rule).
func (src *residentTaskSource) SubmitPanelToolLaunch(ctx context.Context, tool, args string) error {
	if src == nil || src.run == nil {
		return fmt.Errorf("%w: 本进程没有任务管线（任务入口没有起来），工具没有启动", panel.ErrNoHandlerAttached)
	}
	if ctx != nil && ctx.Err() != nil {
		return fmt.Errorf("%w: 这次启动在提交前已被取消", panel.ErrComposerRequest)
	}
	text, err := src.run.composeToolLaunch(tool, args)
	if err != nil {
		return fmt.Errorf("%w: %v", panel.ErrComposerRequest, err)
	}
	return src.submitTask(text)
}
