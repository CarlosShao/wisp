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

#include <cmath>
#include <memory>
#include <string>
#include <vector>

#include "api/scoped_refptr.h"
#include "api/audio/echo_canceller3_config.h"
#include "modules/audio_processing/aec3/echo_canceller3.h"
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

// The ticket-392 tuning factory: it hands APM an EchoCanceller3 built from
// the caller's config - the same class the APM-internal default path would
// construct, so a tuned session differs from a default one only in the
// config fields the caller overrode. This is the binding layer's parameter
// face; the webrtc tree itself is pinned and untouched.
class TunedEchoControlFactory : public webrtc::EchoControlFactory {
 public:
  explicit TunedEchoControlFactory(webrtc::EchoCanceller3Config config)
      : config_(config) {}

  std::unique_ptr<webrtc::EchoControl> Create(int sample_rate_hz,
                                              int num_render_channels,
                                              int num_capture_channels) override {
    return std::make_unique<webrtc::EchoCanceller3>(
        config_, std::nullopt, sample_rate_hz,
        static_cast<size_t>(num_render_channels),
        static_cast<size_t>(num_capture_channels));
  }

 private:
  webrtc::EchoCanceller3Config config_;
};

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

extern "C" void* wap_new_tuned(int sample_rate_hz, const wap_aec3_tuning* tuning,
                               int* out_err) {
    *out_err = 0;
    if (tuning == nullptr) {
        return wap_new(sample_rate_hz, out_err);
    }
    WapSession* s = nullptr;
    try {
        webrtc::EchoCanceller3Config cfg;
        if (tuning->dominant_nearend_enr_threshold > 0) {
            cfg.suppressor.dominant_nearend_detection.enr_threshold =
                static_cast<float>(tuning->dominant_nearend_enr_threshold);
        }
        if (tuning->dominant_nearend_trigger_threshold > 0) {
            cfg.suppressor.dominant_nearend_detection.trigger_threshold =
                tuning->dominant_nearend_trigger_threshold;
        }
        if (tuning->nearend_max_dec_factor_lf > 0) {
            cfg.suppressor.nearend_tuning.max_dec_factor_lf =
                static_cast<float>(tuning->nearend_max_dec_factor_lf);
        }
        webrtc::EchoCanceller3Config::Validate(&cfg);
        webrtc::AudioProcessing::Config config;
        config.echo_canceller.enabled = true;
        config.echo_canceller.mobile_mode = false;
        auto apm = webrtc::AudioProcessingBuilder()
                       .SetConfig(config)
                       .SetEchoControlFactory(
                           std::make_unique<TunedEchoControlFactory>(cfg))
                       .Create();
        if (!apm) {
            *out_err = 1;
            return nullptr;
        }
        s = new WapSession();
        s->apm = apm;
        s->config = webrtc::StreamConfig(sample_rate_hz, 1);
        s->scratch.resize(s->config.num_frames());
        return s;
    } catch (...) {
        delete s;
        *out_err = 2;
        return nullptr;
    }
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
