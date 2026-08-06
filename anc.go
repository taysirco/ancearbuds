// Package ancearbuds models active noise cancellation in true wireless earbuds:
// what it can physically remove, what it cannot, and what it costs in battery.
//
// ANC works by inverting an incoming waveform and playing the inverse back, so
// the two cancel. That works well for predictable low-frequency sound and poorly
// for sudden high-frequency sound, because the system cannot model what it
// cannot anticipate. Passive isolation from the ear tip covers the opposite end
// of the spectrum, which is why fit matters as much as the electronics.
//
// Reference model data used to validate this package is published at
// https://cairovolt.com/en/soundcore/audio
package ancearbuds

import "math"

// ANCType is the microphone arrangement used for cancellation.
type ANCType int

const (
	// None means no active cancellation, passive isolation only.
	None ANCType = iota
	// Feedforward places the microphone outside the earbud.
	Feedforward
	// Feedback places the microphone inside the ear canal.
	Feedback
	// Hybrid combines both and is what most premium models ship.
	Hybrid
)

func (a ANCType) String() string {
	switch a {
	case Feedforward:
		return "feedforward"
	case Feedback:
		return "feedback"
	case Hybrid:
		return "hybrid"
	default:
		return "none"
	}
}

// MaxReductionDB is a representative peak attenuation for the arrangement,
// measured against steady low-frequency noise.
func (a ANCType) MaxReductionDB() float64 {
	switch a {
	case Feedforward:
		return 20
	case Feedback:
		return 22
	case Hybrid:
		return 30
	default:
		return 0
	}
}

// BatteryPenalty is the fraction of playback time lost when ANC runs.
func (a ANCType) BatteryPenalty() float64 {
	if a == None {
		return 0
	}
	if a == Hybrid {
		return 0.30
	}
	return 0.22
}

// Fit describes how well the ear tip seals, which governs passive isolation.
type Fit int

const (
	// PoorSeal leaks audibly and undermines both passive and active performance.
	PoorSeal Fit = iota
	// GoodSeal is a correctly sized tip.
	GoodSeal
	// FoamSeal is memory foam, the best passive isolator.
	FoamSeal
)

// PassiveReductionDB is the attenuation the seal provides without electronics.
func (f Fit) PassiveReductionDB() float64 {
	switch f {
	case GoodSeal:
		return 18
	case FoamSeal:
		return 25
	default:
		return 6
	}
}

// Earbuds combines the electronics and the physical seal.
type Earbuds struct {
	ANC          ANCType
	Fit          Fit
	RatedHoursNoANC float64
}

// ReductionAtDB estimates total attenuation at a given frequency in hertz.
//
// ANC effectiveness falls off sharply above a few hundred hertz because the
// system cannot invert fast, unpredictable waveforms in time. Passive isolation
// does the opposite: it is weak in the bass and strong in the treble.
func (e Earbuds) ReductionAtDB(hz float64) float64 {
	active := e.ANC.MaxReductionDB()
	if hz > 300 {
		// roll off roughly 6 dB per octave above 300 Hz
		octaves := math.Log2(hz / 300)
		active = math.Max(0, active-6*octaves)
	}
	passive := e.Fit.PassiveReductionDB()
	if hz < 500 {
		passive *= hz / 500 // seals do little against deep bass
	}
	return math.Round((active+passive)*10) / 10
}

// PlaybackHours returns expected playback time with ANC enabled.
func (e Earbuds) PlaybackHours() float64 {
	h := e.RatedHoursNoANC * (1 - e.ANC.BatteryPenalty())
	return math.Round(h*10) / 10
}

// BetterAgainst reports which of two configurations attenuates more at a
// frequency, returning "a", "b" or "equal".
func BetterAgainst(hz float64, a, b Earbuds) string {
	x, y := a.ReductionAtDB(hz), b.ReductionAtDB(hz)
	switch {
	case x > y:
		return "a"
	case y > x:
		return "b"
	default:
		return "equal"
	}
}
