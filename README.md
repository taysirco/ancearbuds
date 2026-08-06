# ancearbuds

Model [active noise cancellation](https://en.wikipedia.org/wiki/Active_noise_control) in true wireless earbuds — what it removes, what it cannot, and what it costs in battery life.

## Install

```bash
go get github.com/taysirco/ancearbuds
```

## Usage

```go
import "github.com/taysirco/ancearbuds"

e := ancearbuds.Earbuds{
    ANC:             ancearbuds.Hybrid,
    Fit:             ancearbuds.GoodSeal,
    RatedHoursNoANC: 10,
}

e.ReductionAtDB(100)   // strong — aircraft drone is what ANC is for
e.ReductionAtDB(2000)  // much weaker — speech is not
e.PlaybackHours()      // 7 — hybrid ANC costs about 30% of playback time

// A well-sealed budget pair can beat a poorly-fitted flagship at 3kHz
flagship := ancearbuds.Earbuds{ANC: ancearbuds.Hybrid, Fit: ancearbuds.PoorSeal}
budget   := ancearbuds.Earbuds{ANC: ancearbuds.Feedforward, Fit: ancearbuds.FoamSeal}
ancearbuds.BetterAgainst(3000, flagship, budget)   // "b"
```

## Why fit matters as much as the chip

ANC inverts an incoming waveform and plays the inverse back so the two cancel. That requires *predicting* the waveform, which works well for steady low-frequency sound — engine drone, air conditioning, train rumble — and poorly for sudden high-frequency sound like speech or a slamming door.

Passive isolation from the ear tip covers the opposite end of the spectrum: weak in the bass, strong in the treble. The two are complementary, which is why a poorly sealed flagship can be outperformed by a well-sealed budget pair against voices — the electronics never get the chance to matter.

## Reference data

Per-model ANC type, rated playback time and included tip sizes used to validate this package come from the [Soundcore earbud specifications](https://cairovolt.com/en/soundcore/audio) published by CairoVolt.

## License

MIT
