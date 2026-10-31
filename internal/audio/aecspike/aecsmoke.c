/* aecsmoke.c - the ticket 59 AEC spike's Voice Capture DSP (CWMAudioAEC)
 * driver. Windows-only, compiled by the pinned mingw-w64 toolchain through
 * cgo. See aecsmoke.h for the contract.
 *
 * Why a C shim at all: the DMO is a classic COM/DMO object; driving it from
 * Go over raw vtables would duplicate the mediaobj.h surface, while C gets
 * the headers directly. The spike deliberately does NOT link any import
 * library: every GUID/property key is a local constant, so the build depends
 * only on ole32 and the headers this MSYS2 toolchain ships.
 */

#define COBJMACROS
#define WIN32_LEAN_AND_MEAN
#include <stdint.h>
#include <windows.h>
#include <mmreg.h>
#include <mmdeviceapi.h>
#include <audioclient.h>
#include <audiopolicy.h>
#include <mediaobj.h>
#include "aecsmoke.h"

#if defined(__has_include)
#if __has_include(<propsys.h>)
#include <propsys.h>
#define HAVE_PROPSYS_H 1
#endif
#endif

#ifndef HAVE_PROPSYS_H
/* Minimal IPropertyStore declaration for toolchains whose propsys.h is
 * absent. PROPERTYKEY/PROPVARIANT come from propidl.h via windows.h. */
#define __IPropertyStore_FWD_DEFINED__ 1
typedef struct IPropertyStore IPropertyStore;
typedef struct IPropertyStoreVtbl {
    HRESULT(STDMETHODCALLTYPE *QueryInterface)(IPropertyStore *, REFIID, void **);
    ULONG(STDMETHODCALLTYPE *AddRef)(IPropertyStore *);
    ULONG(STDMETHODCALLTYPE *Release)(IPropertyStore *);
    HRESULT(STDMETHODCALLTYPE *GetCount)(IPropertyStore *, DWORD *);
    HRESULT(STDMETHODCALLTYPE *GetAt)(IPropertyStore *, DWORD, PROPERTYKEY *);
    HRESULT(STDMETHODCALLTYPE *GetValue)(IPropertyStore *, REFPROPERTYKEY, PROPVARIANT *);
    HRESULT(STDMETHODCALLTYPE *SetValue)(IPropertyStore *, REFPROPERTYKEY, const PROPVARIANT *);
    HRESULT(STDMETHODCALLTYPE *Commit)(IPropertyStore *);
} IPropertyStoreVtbl;
struct IPropertyStore {
    IPropertyStoreVtbl *lpVtbl;
};
static const IID wisp_IID_IPropertyStore_real = {
    0x886d8eeb, 0x8cf2, 0x4446, {0x8d, 0x02, 0xcd, 0xba, 0x1d, 0xbd, 0xcf, 0x99}};
#endif /* HAVE_PROPSYS_H */

/* IUnknown pinned locally: the shim links no import library, so the
 * uuid-provided IID_IUnknown symbol is not available. */
static const GUID wisp_IID_IUnknown = {
    0x00000000, 0x0000, 0x0000, {0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}};

/* The Voice Capture DSP coclass (wmcodecdsp.h documents the class; the GUID
 * is pinned here so the spike does not depend on the uuid import lib). */
static const CLSID wisp_CLSID_CWMAudioAEC = {
    0x745057C7, 0xF353, 0x4F2D, {0xA7, 0xEE, 0x58, 0x43, 0x42, 0x7A, 0x0E, 0xF4}};
static const IID wisp_IID_IMediaObject = {
    0xD8AD0F58, 0x5494, 0x4102, {0xAB, 0x73, 0xE7, 0x80, 0x68, 0x56, 0x77, 0xEE}};
static const IID wisp_IID_IMediaBuffer = {
    0x59BAFF3C, 0x795A, 0x4ADD, {0x9D, 0xB4, 0x4E, 0x31, 0xFF, 0x83, 0xD1, 0xC6}};

#ifndef HAVE_PROPSYS_H
#define wisp_IID_IPropertyStore_real wisp_IID_IPropertyStore
#else
static const IID wisp_IID_IPropertyStore_real = {
    0x886d8eeb, 0x8cf2, 0x4446, {0x8d, 0x02, 0xcd, 0xba, 0x1d, 0xbd, 0xcf, 0x99}};
#endif

/* MFPKEY_WMAAECMA_* property set (Windows SDK wmcodecdsp.h; mingw's copy of
 * that header does not carry these, so they are pinned verbatim here). */
static const PROPERTYKEY wisp_PKEY_WMAAECMA_FEATURE_MODE = {
    {0x8b3fd9df, 0x4456, 0x4f09, {0xa4, 0x87, 0x7f, 0xac, 0x16, 0xa5, 0x4f, 0x0e}}, 0};
static const PROPERTYKEY wisp_PKEY_WMAAECMA_MIC_ARRAY_MODE = {
    {0x8b3fd9df, 0x4456, 0x4f09, {0xa4, 0x87, 0x7f, 0xac, 0x16, 0xa5, 0x4f, 0x0e}}, 2};
static const PROPERTYKEY wisp_PKEY_WMAAECMA_SYSTEM_MODE = {
    {0x8b3fd9df, 0x4456, 0x4f09, {0xa4, 0x87, 0x7f, 0xac, 0x16, 0xa5, 0x4f, 0x0e}}, 6};
static const PROPERTYKEY wisp_PKEY_WMAAECMA_DMO_SOURCE_MODE = {
    {0x8b3fd9df, 0x4456, 0x4f09, {0xa4, 0x87, 0x7f, 0xac, 0x16, 0xa5, 0x4f, 0x0e}}, 7};

/* DMO media type GUIDs (uuid.h/uuid.lib territory; pinned locally because
 * the shim links no import library beyond ole32). */
static const GUID wisp_MEDIATYPE_Audio = {
    0x73647561, 0x0000, 0x0010, {0x80, 0x00, 0x00, 0xAA, 0x00, 0x38, 0x9B, 0x71}};
static const GUID wisp_MEDIASUBTYPE_PCM = {
    0x00000001, 0x0000, 0x0010, {0x80, 0x00, 0x00, 0xAA, 0x00, 0x38, 0x9B, 0x71}};
static const GUID wisp_FORMAT_WaveFormatEx = {
    0x05589f81, 0xc356, 0x11ce, {0xbf, 0x01, 0x00, 0xaa, 0x00, 0x55, 0x59, 0x5a}};
static const IID wisp_IID_ISimpleAudioVolume = {
    0x87CE5498, 0x68D6, 0x44E5, {0x92, 0x15, 0x6D, 0xA4, 0x7E, 0xF8, 0x83, 0xD8}};

/* AEC_SYSTEM_MODE: SINGLE_CHANNEL_AEC = 0 (wmcodecdsp.h). */
#define WISP_SINGLE_CHANNEL_AEC 0

#define WISP_CHUNK 1600 /* 100ms at 16k */

static IMediaObject *g_mo = NULL;
static IPropertyStore *g_ps = NULL;
static LONG g_comInit = 0;

/* App-provided IMediaBuffer implementations (the DMO takes buffers, not
 * memory pointers). One each for near-in (stream 0), far-in (stream 1) and
 * out (stream 0). */
typedef struct WispBuffer {
    IMediaBuffer vt;
    LONG refs;
    BYTE *data;
    DWORD maxLen;
    DWORD len;
} WispBuffer;

static HRESULT STDMETHODCALLTYPE wbQueryInterface(IMediaBuffer *self, REFIID riid, void **out) {
    if (IsEqualGUID(riid, &wisp_IID_IUnknown) || IsEqualGUID(riid, &wisp_IID_IMediaBuffer)) {
        *out = self;
        self->lpVtbl->AddRef(self);
        return S_OK;
    }
    *out = NULL;
    return E_NOINTERFACE;
}

static ULONG STDMETHODCALLTYPE wbAddRef(IMediaBuffer *self) {
    WispBuffer *b = (WispBuffer *)self;
    return InterlockedIncrement(&b->refs);
}

static ULONG STDMETHODCALLTYPE wbRelease(IMediaBuffer *self) {
    WispBuffer *b = (WispBuffer *)self;
    return InterlockedDecrement(&b->refs);
}

static HRESULT STDMETHODCALLTYPE wbSetLength(IMediaBuffer *self, DWORD len) {
    WispBuffer *b = (WispBuffer *)self;
    if (len > b->maxLen) {
        return E_INVALIDARG;
    }
    b->len = len;
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE wbGetBufferAndLength(IMediaBuffer *self, BYTE **buf, DWORD *len) {
    WispBuffer *b = (WispBuffer *)self;
    if (buf) {
        *buf = b->data;
    }
    if (len) {
        *len = b->len;
    }
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE wbGetMaxLength(IMediaBuffer *self, DWORD *maxLen) {
    WispBuffer *b = (WispBuffer *)self;
    if (maxLen) {
        *maxLen = b->maxLen;
    }
    return S_OK;
}

static IMediaBufferVtbl wbVtbl = {
    wbQueryInterface, wbAddRef, wbRelease, wbSetLength, wbGetMaxLength, wbGetBufferAndLength};

static WispBuffer g_nearBuf, g_farBuf, g_outBuf;

static void wbInit(WispBuffer *b) {
    b->vt.lpVtbl = &wbVtbl;
    b->refs = 1;
    b->maxLen = WISP_CHUNK * 2;
    b->len = 0;
    b->data = (BYTE *)HeapAlloc(GetProcessHeap(), 0, b->maxLen);
}

static REFERENCE_TIME g_ts = 0;

static HRESULT wispSetProp(IPropertyStore *ps, const PROPERTYKEY *key, PROPVARIANT *pv) {
    return IPropertyStore_SetValue(ps, key, pv);
}

int64_t aecsmoke_probe(int32_t *in_streams, int32_t *out_streams) {
    HRESULT hr = CoInitializeEx(NULL, COINIT_MULTITHREADED);
    if (FAILED(hr)) {
        return (int64_t)hr;
    }
    IMediaObject *mo = NULL;
    hr = CoCreateInstance(&wisp_CLSID_CWMAudioAEC, NULL, CLSCTX_INPROC_SERVER,
                          &wisp_IID_IMediaObject, (void **)&mo);
    if (SUCCEEDED(hr)) {
        DWORD cin = 0, cout = 0;
        HRESULT hrStreams = IMediaObject_GetStreamCount(mo, &cin, &cout);
        if (SUCCEEDED(hrStreams)) {
            if (in_streams) {
                *in_streams = (int32_t)cin;
            }
            if (out_streams) {
                *out_streams = (int32_t)cout;
            }
            hr = S_OK;
        } else {
            hr = hrStreams;
        }
        IMediaObject_Release(mo);
    }
    CoUninitialize();
    return (int64_t)hr;
}

int64_t aecsmoke_begin(void) {
    HRESULT hr = CoInitializeEx(NULL, COINIT_MULTITHREADED);
    if (FAILED(hr)) {
        return (int64_t)hr;
    }
    g_comInit = 1;
    hr = CoCreateInstance(&wisp_CLSID_CWMAudioAEC, NULL, CLSCTX_INPROC_SERVER,
                          &wisp_IID_IMediaObject, (void **)&g_mo);
    if (FAILED(hr)) {
        aecsmoke_end();
        return (int64_t)hr;
    }
    if (SUCCEEDED(g_mo->lpVtbl->QueryInterface(g_mo, &wisp_IID_IPropertyStore_real,
                                               (void **)&g_ps))) {
        PROPVARIANT pv;
        pv.vt = VT_BOOL;
        pv.boolVal = VARIANT_FALSE; /* feed mode: the caller owns the devices */
        wispSetProp(g_ps, &wisp_PKEY_WMAAECMA_DMO_SOURCE_MODE, &pv);
        pv.vt = VT_I4;
        pv.lVal = WISP_SINGLE_CHANNEL_AEC;
        wispSetProp(g_ps, &wisp_PKEY_WMAAECMA_SYSTEM_MODE, &pv);
        wispSetProp(g_ps, &wisp_PKEY_WMAAECMA_MIC_ARRAY_MODE, &pv);
        IPropertyStore_Release(g_ps);
        g_ps = NULL;
    }
    /* 16k/mono/s16 on both inputs and the output, as DMO_MEDIA_TYPE with a
     * FORMAT_WaveFormatEx format blob (SetInputType does not take a bare
     * WAVEFORMATEX). */
    static WAVEFORMATEX wfx;
    ZeroMemory(&wfx, sizeof(wfx));
    wfx.wFormatTag = WAVE_FORMAT_PCM;
    wfx.nChannels = 1;
    wfx.nSamplesPerSec = 16000;
    wfx.wBitsPerSample = 16;
    wfx.nBlockAlign = 2;
    wfx.nAvgBytesPerSec = 32000;
    DMO_MEDIA_TYPE mt;
    ZeroMemory(&mt, sizeof(mt));
    mt.majortype = wisp_MEDIATYPE_Audio;
    mt.subtype = wisp_MEDIASUBTYPE_PCM;
    mt.bFixedSizeSamples = TRUE;
    mt.lSampleSize = 2;
    mt.formattype = wisp_FORMAT_WaveFormatEx;
    mt.cbFormat = sizeof(wfx);
    mt.pbFormat = (BYTE *)&wfx;
    hr = IMediaObject_SetInputType(g_mo, 0, &mt, 0);
    if (FAILED(hr)) {
        aecsmoke_end();
        return (int64_t)hr;
    }
    hr = IMediaObject_SetInputType(g_mo, 1, &mt, 0);
    if (FAILED(hr)) {
        aecsmoke_end();
        return (int64_t)hr;
    }
    hr = IMediaObject_SetOutputType(g_mo, 0, &mt, 0);
    if (FAILED(hr)) {
        aecsmoke_end();
        return (int64_t)hr;
    }
    IMediaObject_AllocateStreamingResources(g_mo);
    wbInit(&g_nearBuf);
    wbInit(&g_farBuf);
    wbInit(&g_outBuf);
    g_ts = 0;
    return (int64_t)S_OK;
}

int64_t aecsmoke_process(const int16_t *far_ref, const int16_t *near_in,
                         int16_t *out, int32_t n, int32_t *produced) {
    if (produced) {
        *produced = 0;
    }
    if (!g_mo || n <= 0 || n > WISP_CHUNK) {
        return (int64_t)E_INVALIDARG;
    }
    DWORD status = 0;
    if (SUCCEEDED(IMediaObject_GetInputStatus(g_mo, 0, &status)) &&
        (status & DMO_INPUT_STATUSF_ACCEPT_DATA)) {
        memcpy(g_nearBuf.data, near_in, (size_t)n * 2);
        g_nearBuf.len = (DWORD)n * 2;
        IMediaObject_ProcessInput(g_mo, 0, &g_nearBuf.vt, 0, g_ts, (REFERENCE_TIME)n * 625);
    }
    if (SUCCEEDED(IMediaObject_GetInputStatus(g_mo, 1, &status)) &&
        (status & DMO_INPUT_STATUSF_ACCEPT_DATA)) {
        memcpy(g_farBuf.data, far_ref, (size_t)n * 2);
        g_farBuf.len = (DWORD)n * 2;
        IMediaObject_ProcessInput(g_mo, 1, &g_farBuf.vt, 0, g_ts, (REFERENCE_TIME)n * 625);
    }
    DMO_OUTPUT_DATA_BUFFER odb;
    ZeroMemory(&odb, sizeof(odb));
    odb.pBuffer = &g_outBuf.vt;
    DWORD reserved = 0;
    HRESULT hr = IMediaObject_ProcessOutput(g_mo, 0, 1, &odb, &reserved);
    if (hr == S_FALSE) {
        g_ts += (REFERENCE_TIME)n * 625;
        return 1; /* engine still buffering; feed the next chunk */
    }
    if (FAILED(hr)) {
        return (int64_t)hr;
    }
    DWORD made = g_outBuf.len / 2;
    if (made > (DWORD)n) {
        made = (DWORD)n;
    }
    if (made > 0) {
        memcpy(out, g_outBuf.data, (size_t)made * 2);
    }
    if (produced) {
        *produced = (int32_t)made;
    }
    g_ts += (REFERENCE_TIME)n * 625;
    return 0;
}

void aecsmoke_end(void) {
    if (g_mo) {
        IMediaObject_FreeStreamingResources(g_mo);
        IMediaObject_Release(g_mo);
        g_mo = NULL;
    }
    if (g_ps) {
        IPropertyStore_Release(g_ps);
        g_ps = NULL;
    }
    if (g_nearBuf.data) {
        HeapFree(GetProcessHeap(), 0, g_nearBuf.data);
        ZeroMemory(&g_nearBuf, sizeof(g_nearBuf));
    }
    if (g_farBuf.data) {
        HeapFree(GetProcessHeap(), 0, g_farBuf.data);
        ZeroMemory(&g_farBuf, sizeof(g_farBuf));
    }
    if (g_outBuf.data) {
        HeapFree(GetProcessHeap(), 0, g_outBuf.data);
        ZeroMemory(&g_outBuf, sizeof(g_outBuf));
    }
    if (g_comInit) {
        CoUninitialize();
        g_comInit = 0;
    }
}

/* --- session volume -----------------------------------------------------
 * ISimpleAudioVolume::SetMasterVolume takes a FLOAT parameter. On the x64
 * Windows ABI float arguments ride the XMM registers, which Go's
 * syscall.SyscallN cannot populate - an integer-path call made the engine
 * read garbage (observed: a near-zero session volume that Windows then
 * PERSISTED for this process's session, silencing every later render run).
 * The fix is a real C call on the rig's own audio client (GetService hands
 * out the stream's session volume - no session enumeration needed). */
static float g_lastVolume = -1.0f;

int64_t aecsmoke_set_client_volume(void *client, float v) {
    if (!client) {
        return (int64_t)E_POINTER;
    }
    HRESULT hrInit = CoInitializeEx(NULL, COINIT_MULTITHREADED);
    if (FAILED(hrInit)) {
        return (int64_t)hrInit;
    }
    IAudioClient *ac = (IAudioClient *)client;
    ISimpleAudioVolume *sv = NULL;
    HRESULT hr = IAudioClient_GetService(ac, &wisp_IID_ISimpleAudioVolume, (void **)&sv);
    if (SUCCEEDED(hr)) {
        hr = ISimpleAudioVolume_SetMasterVolume(sv, v, NULL);
        ISimpleAudioVolume_SetMute(sv, FALSE, NULL);
        ISimpleAudioVolume_GetMasterVolume(sv, &g_lastVolume);
        ISimpleAudioVolume_Release(sv);
    }
    CoUninitialize();
    return (int64_t)hr;
}

/* Read back the session volume the last aecsmoke_set_client_volume saw. */
float aecsmoke_last_volume(void) {
    return g_lastVolume;
}
