//go:build windows

package main

// 票 441 AC#2 的 boot 半边：app.theme 的球侧种子。外观波（7f3177f93）给
// [ball] 四枚键装了 boot 种子，但 [app].theme 的球侧消费一直悬着：
// internal/ball/renderer_windows.go:106 的 SetTheme 全仓零调用者，渲染器的
// ACTIVE C21 palette 永远停在 DarkPalette()（census-unwired 2026-11-16）。
// 本件把读的那半补上——调用的那半在 cmd/wisp/resident_windows.go 的装配根
// （与 ballStyleSeed65 → ball.SetStyle 同一座位、同一时序：都在 ball.New
// 之前，所以第一帧就画新 palette，"caller must trigger a repaint afterwards"
// 的那条注释由预建时序免掉）。
//
// WHY A FRESH READ INSTEAD OF THE MANAGER. 同 ballIdleAppearanceSeed 的理由
// （该文件头部那段原话）：球建在任务管线的 config.Manager 之前，种子按
// 每次 boot 用 config.LoadFile 现读。
//
// ok 的语义（与其他种子一致，fallback 说出口、不冒充）：
//   - ok=true  ⇒ 文件读到且 app.theme 是 dark/light 之一，theme 是文件的
//     自己的回答（默认值也算回答）；
//   - ok=false ⇒ 文件读不了（renderer 的 dark 默认原地站着，Warn 具名），
//     或 app.theme = auto——本进程在这里没有"系统主题"这个答案可解析
//     （D29 的 system-theme follow 是另一笔登记的活），auto 冒充 dark
//     就是谎报，所以原样不调，Warn 把它说出口。
//
// NO BAND AND NO DEFAULT OF ITS OWN. 词表住在 internal/config（schema 的
// themes 切片 + Validate），映射只做 dark/light 两枚枚举的直译，本件不造
// 第三档、不造第三份默认。

import (
	"log/slog"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
)

// ballThemeSeed441 reads [app].theme the way ballIdleAppearanceSeed reads
// [ball].opacity_idle (one LoadFile per boot, the file's own answer or a
// named fallback). ok=true says "call SetTheme with this"; ok=false says
// "nothing to call with - the renderer's dark palette stands, and the Warn
// named why".
func ballThemeSeed441(dataDir string) (theme ball.Theme, ok bool) {
	cfgPath := filepath.Join(dataDir, configFileName)
	c, _, err := config.LoadFile(cfgPath, nil)
	if err != nil || c == nil {
		slog.Warn("ball: [app].theme source unreadable at boot; the renderer's dark palette stands",
			"path", cfgPath, "err", err)
		return ball.ThemeDark, false
	}
	switch c.App.Theme {
	case config.ThemeLight:
		return ball.ThemeLight, true
	case config.ThemeDark:
		return ball.ThemeDark, true
	default:
		// auto (以及任何校验漏进来的词)：这里没有系统主题可解析。
		slog.Warn("ball: [app].theme has no resolvable answer at boot (system-theme follow is a separate D29 item); the renderer's dark palette stands",
			"theme", c.App.Theme, "path", cfgPath)
		return ball.ThemeDark, false
	}
}
