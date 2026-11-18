//go:build windows

package ball

import (
	"reflect"
	"testing"
)

// 票 441 AC#2 的渲染器入参牙：SetTheme 必须真换 ACTIVE palette（pal）——
// renderer 的每一笔都读它（renderer_windows.go 的 pal 注释），装配根的
// SetTheme 调用（cmd/wisp/resident_windows.go）值不值得接线就看这一枚。
// 纯函数、不开窗：pal 与 PaletteFor 都是包内值，双向各证一次
// （light 真换、dark 换回），收尾恢复 dark 默认，不给同包其他测试留痕。
func TestSetTheme441SwitchesTheActivePalette(t *testing.T) {
	t.Cleanup(func() { SetTheme(ThemeDark) })

	SetTheme(ThemeLight)
	if !reflect.DeepEqual(pal, PaletteFor(ThemeLight)) {
		t.Fatalf("after SetTheme(ThemeLight) the active palette is not PaletteFor(ThemeLight): the switch is a lie")
	}
	if reflect.DeepEqual(pal, DarkPalette()) {
		t.Fatalf("after SetTheme(ThemeLight) the active palette still equals DarkPalette(): the tooth cannot tell the two palettes apart")
	}

	SetTheme(ThemeDark)
	if !reflect.DeepEqual(pal, DarkPalette()) {
		t.Fatalf("after SetTheme(ThemeDark) the active palette is not DarkPalette(): the switch back is a lie")
	}
	if reflect.DeepEqual(pal, LightPalette()) {
		t.Fatalf("after SetTheme(ThemeDark) the active palette still equals LightPalette(): the tooth cannot tell the two palettes apart")
	}
}
