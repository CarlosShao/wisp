package audio

// In-process linear-interpolation resampler (SPEC-04 sec 3: device-native
// rate -> 16k, no new dependencies). Streams with a 64-bit fixed-point
// position accumulator in ABSOLUTE input-stream coordinates, so chunk
// boundaries never drift and the output for a given total input is exactly
// the one-shot output regardless of chunking (asserted by test).
//
// Latency bound: linear interpolation needs exactly one lookahead sample,
// so a converted sample is at most one input sample behind the source
// (about 21us at 48kHz) - bounded by construction, asserted by test.
//
// Deterministic length: an output sample i sits at input position
// i*step (32.32 fixed point) and is emitted only while both interpolation
// endpoints exist (floor(i*step) <= n-2). The trailing partial step
// (<= 1 input sample, ~21us at 48k) is omitted by design; the capture path
// is endless anyway and Flush exists only to reset the state.

// Resampler converts mono int16 samples from FromRate to ToRate, streaming.
// Use NewResampler; the zero value is not usable.
type Resampler struct {
	from, to int
	step     int64 // input-sample advance per output sample, 32.32 fixed point
	pos      int64 // next output position in ABSOLUTE input coordinates, 32.32

	consumed   int   // input samples fully absorbed into the window
	carry      int16 // last consumed sample (in[consumed-1] of the stream)
	carryValid bool
}

// NewResampler creates a streaming rate converter. from == to is legal and
// makes Process pass samples through untouched.
func NewResampler(from, to int) *Resampler {
	if from <= 0 || to <= 0 {
		return nil
	}
	return &Resampler{
		from: from,
		to:   to,
		step: (int64(from) << 32) / int64(to),
	}
}

// Rates reports the configured conversion.
func (r *Resampler) Rates() (from, to int) { return r.from, r.to }

// Process consumes one chunk and returns the output samples derivable from
// all input seen so far. Chunks may be any size (including empty).
func (r *Resampler) Process(in []int16) []int16 {
	if r.from == r.to {
		return in
	}
	var out []int16
	end := int64(r.consumed + len(in)) // absolute count of available input
	for r.pos < (end-1)<<32 {
		ip := r.pos >> 32
		frac := uint64(r.pos & 0xffffffff)
		a := int32(r.abs(in, ip))
		b := int32(r.abs(in, ip+1))
		delta := (int64(b-a) * int64(frac)) >> 32
		out = append(out, clampI16(a+int32(delta)))
		r.pos += r.step
	}
	if len(in) > 0 {
		r.carry = in[len(in)-1]
		r.carryValid = true
		r.consumed += len(in)
	}
	return out
}

// abs returns absolute input sample i; i is guaranteed to lie inside the
// window [consumed-1, consumed+len(in)-1] by the emit condition.
func (r *Resampler) abs(in []int16, i int64) int16 {
	if i == int64(r.consumed)-1 {
		return r.carry
	}
	return in[i-int64(r.consumed)]
}

// Flush resets the stream state (see the file comment for the tail policy).
// It returns nothing; output lengths are deterministic by construction.
func (r *Resampler) Flush() {
	r.pos = 0
	r.consumed = 0
	r.carry = 0
	r.carryValid = false
}

// ResampleLinear is the one-shot form: stream the whole signal.
// outLen = ceil((n-1)*toRate/fromRate) (the i >= 0 with floor(i*step) <= n-2).
func ResampleLinear(in []int16, fromRate, toRate int) []int16 {
	if fromRate == toRate {
		return in
	}
	r := NewResampler(fromRate, toRate)
	return r.Process(in)
}

// MonoDownmix averages interleaved multi-channel samples to mono. chans <= 1
// returns the input unchanged.
func MonoDownmix(in []int16, chans int) []int16 {
	if chans <= 1 || len(in) == 0 {
		return in
	}
	n := len(in) / chans
	out := make([]int16, n)
	for i := 0; i < n; i++ {
		var acc int32
		for c := 0; c < chans; c++ {
			acc += int32(in[i*chans+c])
		}
		out[i] = clampI16(acc / int32(chans))
	}
	return out
}

// FloatToPCM16 converts interleaved float32 samples [-1,1] to int16.
func FloatToPCM16(in []float32) []int16 {
	out := make([]int16, len(in))
	for i, v := range in {
		s := int64(v*32767 + 0.5)
		if v < 0 {
			s = int64(v*32768 - 0.5)
		}
		out[i] = clampI16(int32(s))
	}
	return out
}

func clampI16(v int32) int16 {
	if v > 32767 {
		return 32767
	}
	if v < -32768 {
		return -32768
	}
	return int16(v)
}
