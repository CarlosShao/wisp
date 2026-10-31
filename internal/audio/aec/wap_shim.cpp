//go:build windows && cgo && webrtcap

/*
 * wap_shim.cpp - the webrtc side of the aec package's C ABI. Compiled by
 * cgo together with the package when the webrtcap tag is set; include dirs
 * for the fetched webrtc-audio-processing tree (and its pinned abseil)
 * come in through CGO_CXXFLAGS, the built static archive through
 * CGO_LDFLAGS. See deps.toml [webrtc-audio-processing] and
 * docs/evidence/s1/aec3r1-remeasure.md for the pinned provenance.
 *
 * Measurement note: only the echo canceller is enabled here. The upstream
 * offline example also turns on AGC1/AGC2/high-pass, but gain control
 * rescales the near signal - and with it the residual echo - so the rig's
 * ERLE reading must run with those stages off to measure the engine alone.
 */
#include "wap_shim.h"

#include <string>
#include <vector>

#include "api/scoped_refptr.h"
#include "modules/audio_processing/include/audio_processing.h"

namespace {

struct WapSession {
    rtc::scoped_refptr<webrtc::AudioProcessing> apm;
    webrtc::StreamConfig config{0, 1};
    std::vector<float> scratch;
    std::string last_error;
};

void note_error(WapSession* s, const char* what, int code) {
    s->last_error = std::string(what) + " code=" + std::to_string(code);
}

}  // namespace

extern "C" void* wap_new(int sample_rate_hz, int* out_err) {
    *out_err = 0;
    WapSession* s = nullptr;
    try {
        s = new WapSession();
    } catch (...) {
        *out_err = 2;
        return nullptr;
    }
    s->apm = webrtc::AudioProcessingBuilder().Create();
    if (!s->apm) {
        delete s;
        *out_err = 1;
        return nullptr;
    }
    webrtc::AudioProcessing::Config config;
    config.echo_canceller.enabled = true;
    config.echo_canceller.mobile_mode = false;
    s->apm->ApplyConfig(config);
    s->config = webrtc::StreamConfig(sample_rate_hz, 1);
    s->scratch.resize(s->config.num_frames());
    return s;
}

extern "C" void wap_free(void* session) { delete static_cast<WapSession*>(session); }

extern "C" int wap_process_reverse(void* session, const float* far_frame, int frames) {
    WapSession* s = static_cast<WapSession*>(session);
    const float* src[1] = {far_frame};
    float* dest[1] = {s->scratch.data()};
    int rc = s->apm->ProcessReverseStream(src, s->config, s->config, dest);
    if (rc != 0) {
        note_error(s, "ProcessReverseStream", rc);
    }
    return rc;
}

extern "C" int wap_process_capture(void* session, const float* near_frame, int frames,
                                   int delay_ms, float* out_frame) {
    WapSession* s = static_cast<WapSession*>(session);
    int drc = s->apm->set_stream_delay_ms(delay_ms);
    if (drc != 0) {
        note_error(s, "set_stream_delay_ms", drc);
        return drc;
    }
    const float* src[1] = {near_frame};
    float* dest[1] = {out_frame};
    int rc = s->apm->ProcessStream(src, s->config, s->config, dest);
    if (rc != 0) {
        note_error(s, "ProcessStream", rc);
    }
    return rc;
}

extern "C" const char* wap_last_error(void* session) {
    WapSession* s = static_cast<WapSession*>(session);
    return s->last_error.c_str();
}
