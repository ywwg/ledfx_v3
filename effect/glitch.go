package effect

import (
	"ledfx/audio"
	"ledfx/color"
	"ledfx/logger"
	"math"
)

/*
Audio reactive port of PixelBlaze effect "Glitch Bands"
Credit to Ben Henke for original concept & fantastic LED art
https://electromage.com/patterns
*/

type Glitch struct{}

// movement speed based on beat
// modify saturation by vocals/mids
// make the bars/legs grow and shrink
// modulate t1 with highs

// Apply new pixels to an existing pixel array.
func (e *Glitch) assembleFrame(base *Effect, p color.Pixels) {
	mel, err := audio.Analyzer.GetMelbank(base.ID)
	if err != nil {
		logger.Logger.WithField("context", "Effect").Error(err)
		return
	}

	t1 := base.time(0.1) * math.Pi * 2
	t2 := base.time(0.1)
	t3 := base.time(0.5)
	t4 := base.time(0.2) * math.Pi * 2
	t5 := base.time(0.05)
	t6 := base.time(0.02)

	for i := 0; i < len(p); i++ {
		fi := float64(i)
		m := (0.3 + base.triangle(t2)*0.2)
		h := math.Sin(t1) + math.Mod(((fi-base.pixelScaler/2)/base.pixelScaler)*(base.triangle(t3)*10+4*math.Sin(t4)), m)
		s1 := base.triangle(math.Mod(t5+fi/base.pixelScaler*5, 1))
		s1 = math.Pow(s1, 2)
		s2 := base.triangle(math.Mod(t6-(fi-base.pixelScaler)/base.pixelScaler, 1))
		s2 = math.Pow(s2, 4)
		s := 1 - base.triangle(s1*s2)
		v := 0.5
		if s1 > s2 {
			v = 1 - s1
		} else {
			v = 0.5 + s2
		}

		p[i][0] = h + mel.HighAmplitude()
		p[i][1] = s - mel.LowsAmplitude()
		p[i][2] = v + mel.MidsAmplitude()
	}
}
