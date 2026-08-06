package ancearbuds

import "testing"

func TestANCDominatesLowFrequencies(t *testing.T) {
	e := Earbuds{ANC: Hybrid, Fit: GoodSeal, RatedHoursNoANC: 8}
	engine := e.ReductionAtDB(100) // aircraft drone
	voice := e.ReductionAtDB(2000) // speech
	if engine <= voice {
		t.Errorf("ANC should do more against 100Hz (%v) than 2kHz (%v)", engine, voice)
	}
}

func TestSealCarriesTheTreble(t *testing.T) {
	good := Earbuds{ANC: None, Fit: FoamSeal, RatedHoursNoANC: 8}
	poor := Earbuds{ANC: None, Fit: PoorSeal, RatedHoursNoANC: 8}
	if good.ReductionAtDB(4000) <= poor.ReductionAtDB(4000) {
		t.Error("a foam seal must outperform a poor seal against treble")
	}
}

func TestPoorSealUnderminesGoodElectronics(t *testing.T) {
	flagship := Earbuds{ANC: Hybrid, Fit: PoorSeal, RatedHoursNoANC: 8}
	budget := Earbuds{ANC: Feedforward, Fit: FoamSeal, RatedHoursNoANC: 8}
	if got := BetterAgainst(3000, flagship, budget); got != "b" {
		t.Errorf("a well-sealed budget pair should win at 3kHz, got %q", got)
	}
}

func TestANCCostsBatteryLife(t *testing.T) {
	e := Earbuds{ANC: Hybrid, Fit: GoodSeal, RatedHoursNoANC: 10}
	if got := e.PlaybackHours(); got != 7 {
		t.Errorf("10h rated with hybrid ANC: got %vh, want 7", got)
	}
	off := Earbuds{ANC: None, Fit: GoodSeal, RatedHoursNoANC: 10}
	if off.PlaybackHours() != 10 {
		t.Error("no ANC means no battery penalty")
	}
}
