//go:build windows

// ticket 448 -- 行为那一半的常驻闸：errnoPresent 的判据必须是码，不是句子。
//
// 446 的证明是一枚 -overlay 临时探针（probe446_errno_presence_windows_test.go），
// 跑完就没了；本文件是把它留下来的版本，并多做了一件那枚探针没做的事：把「哪些
// 行在任何语言/区域的 Windows 上都会红」写成机器可查的字段，而不是散文。
//
// 为什么这事需要一个字段：本机（en-US 消息表拿得到）上，旧写法
// `e != nil && e.Error() != "The operation completed successfully."` 在
// errno-zero-raw / errno-zero-via-windows-alias 这两行上**碰巧是绿的** —— 那正是
// 446 追到 syscall_windows.go:169-171 的那个坑（Go 先显式要 en-US，失败才退回系统
// 语言）。所以这两行不能当牙。真正与语言无关的牙是三处：
//   - 本地化文本：errno 0 带着非英文文本（真机在没有 en-US 表的机器上就是这个状态）
//     旧写法判它「有 errno」⇒ 红；
//   - 包了一层：wrapped 文本是前缀加句子，句子比较只看得见最外层 ⇒ 旧写法判「有
//     errno」⇒ 红；
//   - 长得像句子的错误：一个只是**读起来**像成功句子的真错误，句子比较把它当成功
//     吞掉 ⇒ 红。
//
// 这三行是 AC#3「任何语言的 Windows 上都成立」的可执行形式，门里断言至少 3 处。
//
// 输入全是 syscall.Errno(0) 的字面值／本地化样例，不读 GetUserDefaultUILanguage，
// 也不读 GetSystemDefaultLangID：闸门的答案不许随本机语言变。
package tools

import (
	"errors"
	"fmt"
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

// errnoZeroLocalized448 是 errno 0 带着非英文文本，也就是拿不到 en-US 消息表的
// Windows 交回来的那个状态。真机不需要任何假类型就能到这个状态；假类型只为了让这
// 台有英文 UI 语言包的机器也能把那条路走一遍。Is 是 errors.Is 拆包后咨询的方法，
// 所以判据在这里必须答「没有 errno」。
type errnoZeroLocalized448 struct{}

func (errnoZeroLocalized448) Error() string        { return "操作成功完成。" }
func (errnoZeroLocalized448) Is(target error) bool { return target == syscall.Errno(0) }

// errnoCase448 是一行判据。localeIndependentTooth 是「446 的旧写法在任何语言的
// Windows 上都会在这行上答错」的钉子，不是注解。
type errnoCase448 struct {
	name                   string
	err                    error
	want                   bool
	localeIndependentTooth bool
}

// Test448ErrnoPresenceIsDecidedByCode 走 errno 0 的四种形状，加上真错误与控制行。
func Test448ErrnoPresenceIsDecidedByCode(t *testing.T) {
	wrappedZero := fmt.Errorf("SendInput wrapper: %w", syscall.Errno(0))
	cases := []errnoCase448{
		// errno=0 那条路，也就是旧比较以为自己认得的那条。
		{name: "nil", err: nil, want: false},
		// 这两行旧写法在本机（en-US）碰巧答对，所以刻意不标成牙：它们只在拿不到
		// en-US 表的机器上才红，而那台机器恰好是 446 要修的那种。
		{name: "errno-zero-raw", err: syscall.Errno(0), want: false},
		{name: "errno-zero-via-windows-alias", err: windows.Errno(0), want: false},
		{name: "errno-zero-localized-text", err: errnoZeroLocalized448{}, want: false, localeIndependentTooth: true},
		{name: "errno-zero-wrapped", err: wrappedZero, want: false, localeIndependentTooth: true},
		{name: "errno-zero-wrapped-twice", err: fmt.Errorf("outer: %w", wrappedZero), want: false, localeIndependentTooth: true},
		// 真失败必须照样被报出来。
		{name: "errno-5-access-denied", err: syscall.Errno(5), want: true},
		{name: "errno-87-bad-parameter", err: syscall.Errno(87), want: true},
		{name: "errno-87-wrapped", err: fmt.Errorf("ctx: %w", syscall.Errno(87)), want: true},
		// 句子不是身份：一个只是读起来像成功文本的真错误，旧比较把它当「没有 errno」
		// 吞掉。这一行就是下面那个英文字面量存在的唯一理由，也是与语言无关的牙。
		{
			name: "lookalike-text-is-not-errno-zero",
			err:  errors.New("The operation completed successfully."), want: true, localeIndependentTooth: true,
		},
	}

	teeth := 0
	teethRed := 0
	red := 0
	for _, c := range cases {
		if c.localeIndependentTooth {
			teeth++
		}
		got := errnoPresent(c.err)
		t.Logf("448 case %-32s errnoPresent=%-5v want=%-5v tooth=%-5v err=%q",
			c.name, got, c.want, c.localeIndependentTooth, errnoText448(c.err))
		if got == c.want {
			continue
		}
		red++
		if c.localeIndependentTooth {
			teethRed++
		}
		// 逐行 Errorf 而不是 Fatalf：一次失败要把**所有**读文本的行都点出来，
		// 否则证据里只会留第一条，读者看不出牙有几颗。
		t.Errorf("448 红在 %q：errnoPresent(%q) = %v，want %v -- 判定又在读消息文本，而不是读码",
			c.name, errnoText448(c.err), got, c.want)
	}
	if teeth < 3 {
		t.Fatalf("与语言无关的牙只剩 %d 处（至少 3 处）：这个闸门又变成「只在本机这种语言上成立」了", teeth)
	}
	if red != 0 {
		t.Fatalf("448 行为闸红：%d/%d 行答错，其中 %d 行是与语言无关的牙（这些行在任何语言的 Windows 上都会红，所以这不是本机语言的问题）",
			red, len(cases), teethRed)
	}
	t.Logf("448 行为闸：%d 行全绿，其中 %d 行在任何语言的 Windows 上都会红", len(cases), teeth)
}

func errnoText448(e error) string {
	if e == nil {
		return "<nil>"
	}
	return e.Error()
}
