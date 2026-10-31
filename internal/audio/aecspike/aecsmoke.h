/* aecsmoke.h - the C surface of the ticket 59 AEC spike's Voice Capture DSP
 * probe (CLSID_CWMAudioAEC, the Windows built-in voice AEC DMO).
 *
 * Every call must run on ONE pinned OS thread (the rig's captured thread):
 * the session state below is process-global on purpose, and COM is
 * initialized once in aecsmoke_begin and balanced in aecsmoke_end.
 *
 * Return values are HRESULTs (0 = S_OK) except where noted:
 *   aecsmoke_process: 0 = output chunk produced, 1 = no output yet
 *     (S_FALSE, feed more), otherwise the failing HRESULT.
 */

#ifndef WISP_AECSMOKE_H
#define WISP_AECSMOKE_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/* One-shot: CoCreateInstance the voice AEC DMO, ask for its stream layout,
 * release it. Writes the input/output stream counts on success. */
int64_t aecsmoke_probe(int32_t *in_streams, int32_t *out_streams);

/* Begin a processing session: COM init, create the DMO in feed mode
 * (MFPKEY_WMAAECMA_DMO_SOURCE_MODE = FALSE, AEC system mode), set 16k/mono/
 * s16 on both inputs and the output. 0 on success. */
int64_t aecsmoke_begin(void);

/* Feed one aligned chunk of the far reference (stream 1) and the near/mic
 * signal (stream 0); on output, copy up to n processed samples. */
int64_t aecsmoke_process(const int16_t *far_ref, const int16_t *near_in,
                         int16_t *out, int32_t n, int32_t *produced);

/* End the session and balance CoUninitialize. */
void aecsmoke_end(void);

/* Set the rig render stream's session master volume through a real C call
 * (float arguments ride XMM registers Go's syscall cannot populate) and
 * undo any earlier integer-path ISimpleAudioVolume damage. Also unmutes. */
int64_t aecsmoke_set_client_volume(void *client, float v);

/* The session volume read back by the last set_client_volume call. */
float aecsmoke_last_volume(void);

#ifdef __cplusplus
}
#endif

#endif /* WISP_AECSMOKE_H */
