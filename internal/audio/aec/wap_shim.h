//go:build windows && cgo && webrtcap

/*
 * wap_shim.h - C ABI over the webrtc-audio-processing APM (AEC3) for the Go
 * binding in aec_webrtc_windows.go. All webrtc C++ types stay behind this
 * shim; the Go side sees void* sessions, float frames and int return codes.
 */
#ifndef WAP_SHIM_H
#define WAP_SHIM_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Creates a mono APM with only the echo canceller enabled. Returns NULL on
 * failure and sets *out_err to a shim-level code (1 = APM create failed,
 * 2 = allocation/exception). */
void* wap_new(int sample_rate_hz, int* out_err);

/* Releases a session created by wap_new. NULL is allowed. */
void wap_free(void* session);

/* Feeds one ~10ms far (render) frame; the APM buffers the reference
 * internally. Returns 0 or an AudioProcessing error code. */
int wap_process_reverse(void* session, const float* far_frame, int frames);

/* Feeds one ~10ms near (capture) frame with the measured stream delay and
 * writes the echo-cancelled frame to out_frame. Returns 0 or an
 * AudioProcessing error code. */
int wap_process_capture(void* session, const float* near_frame, int frames,
                        int delay_ms, float* out_frame);

/* Human-readable reason for the last non-zero return code on this session. */
const char* wap_last_error(void* session);

#ifdef __cplusplus
}
#endif

#endif /* WAP_SHIM_H */
