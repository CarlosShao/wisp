//go:build windows

package ball

// Direct2D / DirectWrite COM plumbing for the ball renderer (no cgo, same
// syscall discipline as the S0 spike). Vtable slot indices were taken from
// the mingw-w64 headers (d2d1.h / dwrite.h, MIDL order) and match the spike
// anchors: CreateHwndRenderTarget=14, CreateSolidColorBrush=8, FillRectangle
// =17, DrawText=27, Clear=47, BeginDraw=48, EndDraw=49.
//
// Float-argument ABI (spike finding): syscall args are integer words. On
// win-x64 a float argument in a STACK slot (arg position >= 4) is passed as
// its 4-byte IEEE bits in the low half of the 8-byte slot. Floats in
// REGISTER positions (0..3) cannot be passed this way (the callee reads
// XMM registers) - such calls are avoided by design here:
//   - DrawEllipse/DrawLine strokeWidth sits at position 3/4: DrawLine is
//     safe (stack), DrawEllipse is NOT -> every stroked circle/ring is drawn
//     as a polyline of DrawLine segments (polylineCircle below).
//   - SetOpacity/SetDpi/SetTransform-with-floats are never called; opacity
//     is baked into colors (SetColor) or gradient stops, DPI is fixed at
//     96 on the render target so 1 DIP == 1 physical pixel.
//   - All gradient/brush geometry lives in structs passed by pointer.

import (
	"fmt"
	"math"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modD2D1   = windows.NewLazySystemDLL("d2d1.dll")
	modDWrite = windows.NewLazySystemDLL("dwrite.dll")

	pD2D1CreateFactory   = modD2D1.NewProc("D2D1CreateFactory")
	pDWriteCreateFactory = modDWrite.NewProc("DWriteCreateFactory")
)

// COM vtable slots (0-based; IUnknown::QueryInterface = 0).
const (
	slotFactoryCreateDCRT = 16 // ID2D1Factory::CreateDCRenderTarget
)

const (
	// ID2D1RenderTarget (Base = IUnknown 0-2 + ID2D1Resource::GetFactory 3)
	slotRTCreateSolidColorBrush = 8
	slotRTCreateGradientStops   = 9
	slotRTCreateRadialGradient  = 11
	slotRTDrawLine              = 15
	slotRTFillRectangle         = 17
	slotRTFillEllipse           = 21
	slotRTDrawText              = 27
	slotRTSetTransform          = 30
	slotRTClear                 = 47
	slotRTBeginDraw             = 48
	slotRTEndDraw               = 49
)

const (
	// ID2D1DCRenderTarget (Base = full ID2D1RenderTarget vtable 0..56)
	slotDCRTBindDC = 57
)

const (
	// ID2D1SolidColorBrush (Base = IUnknown 0-2 + ID2D1Resource::GetFactory 3)
	slotSolidSetColor = 8 // after ID2D1Brush (SetOpacity=4..GetTransform=7)
)

const (
	// IDWriteFactory
	slotDWriteCreateTextFormat = 15
	// IDWriteTextFormat
	slotTFSetTextAlignment      = 3
	slotTFSetParagraphAlignment = 4
)

const (
	slotIUnknownRelease = 2
)

// D2D enums.
const (
	d2d1FactoryTypeSingleThreaded = 0
	d2d1Gamma22                   = 0
	d2d1ExtendModeClamp           = 0
	dxgiFormatB8G8R8A8            = 87
	d2d1AlphaModePremultiplied    = 1
)

// GUIDs (mingw-w64 headers; also in docs/evidence/s0/02-spike-report.md).
var (
	iidID2D1Factory   = guid([8]byte{0x92, 0x45, 0x11, 0x8b, 0xfd, 0x3b, 0x60, 0x07}, 0x06152247, 0x6f50, 0x465a)
	iidIDWriteFactory = guid([8]byte{0xa2, 0xe8, 0x1a, 0xdc, 0x7d, 0x93, 0xdb, 0x48}, 0xb859ee5a, 0xd838, 0x4b5b)
)

func guid(data4 [8]byte, d1 uint32, d2, d3 uint16) windows.GUID {
	return windows.GUID{Data1: d1, Data2: d2, Data3: d3, Data4: data4}
}

// comCall invokes slot `slot` (0-based) of the COM object at `this`. COM
// object pointers are kept as unsafe.Pointer throughout the package: the
// vtable deref below then needs no uintptr->pointer conversions (vet-clean
// provenance). Args remain uintptr words (the Win64 syscall ABI).
func comCall(this unsafe.Pointer, slot uintptr, args ...uintptr) (uintptr, uintptr, error) {
	vt := *(*unsafe.Pointer)(this)
	fn := *(*unsafe.Pointer)(unsafe.Add(vt, slot*8))
	return syscall.SyscallN(uintptr(fn), append([]uintptr{uintptr(this)}, args...)...)
}

func comRelease(obj *unsafe.Pointer) {
	if obj != nil && *obj != nil {
		_, _, _ = comCall(*obj, slotIUnknownRelease)
		*obj = nil
	}
}

func f32bits(f float32) uintptr { return uintptr(math.Float32bits(f)) }

// ---------------------------------------------------------------- D2D types

type d2d1ColorF struct{ r, g, b, a float32 }

func colorF(c Color) d2d1ColorF { return d2d1ColorF{c.R, c.G, c.B, c.A} }

type d2d1Point2F struct{ x, y float32 }

type d2d1RectF struct{ l, t, r, b float32 }

type d2d1Ellipse struct {
	center  d2d1Point2F
	radiusX float32
	radiusY float32
}

type d2d1RenderTargetProps struct {
	typ       uint32
	format    uint32
	alphaMode uint32
	dpiX      float32 // 96 -> 1 DIP == 1 physical pixel
	dpiY      float32
	usage     uint32
	minLevel  uint32
}

type d2d1BrushProps struct {
	opacity float32
	m       [6]float32 // identity D2D1_MATRIX_3X2_F
}

func identityBrushProps() d2d1BrushProps {
	return d2d1BrushProps{opacity: 1.0, m: [6]float32{1, 0, 0, 1, 0, 0}}
}

type d2d1GradientStop struct {
	position float32
	color    d2d1ColorF
}

type d2d1RadialGradientProps struct {
	center       d2d1Point2F
	originOffset d2d1Point2F
	radiusX      float32
	radiusY      float32
}

// d2dFactories holds the process-wide COM factories. ONE D2D factory per
// process (D38a): the future panel renderer reuses d2dFactory via the
// accessors; nothing else may call D2D1CreateFactory.
var (
	d2dFactory    unsafe.Pointer
	dwriteFactory unsafe.Pointer
	factoryOnce   sync.Once
	factoryErr    error
)

// ensureFactories creates the shared D2D/DWrite factories. Must run on the
// ui-sta thread (D2D1_FACTORY_TYPE_SINGLE_THREADED).
func ensureFactories() error {
	factoryOnce.Do(func() {
		var opts uint32 // D2D1_DEBUG_LEVEL_NONE
		r1, _, _ := pD2D1CreateFactory.Call(
			d2d1FactoryTypeSingleThreaded,
			unsafePtr(&iidID2D1Factory),
			unsafePtr(&opts),
			unsafePtr(&d2dFactory))
		if r1 != 0 {
			factoryErr = fmt.Errorf("D2D1CreateFactory failed hr=0x%x", r1)
			return
		}
		r1, _, _ = pDWriteCreateFactory.Call(
			0, // DWRITE_FACTORY_TYPE_SHARED
			unsafePtr(&iidIDWriteFactory),
			unsafePtr(&dwriteFactory))
		if r1 != 0 {
			factoryErr = fmt.Errorf("DWriteCreateFactory failed hr=0x%x", r1)
		}
	})
	return factoryErr
}

// FactoriesForPanel exposes the shared factories to the future panel renderer
// (D38a: one D2D factory only). Returns (0, 0, err) before ball init.
func FactoriesForPanel() (d2d, dwrite unsafe.Pointer, err error) {
	if e := ensureFactories(); e != nil {
		return nil, nil, e
	}
	return d2dFactory, dwriteFactory, nil
}
