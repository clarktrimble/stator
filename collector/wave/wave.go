// Package wave produces periodic stats from series of sinusoids.
package wave

import (
	"math"
	"math/rand"
	"time"

	"stator/entity"
)

const (
	name            = "wave"
	fFactor float64 = 20
)

// Wave generates a sine wave stat.
//
// count is used to drive the sine function.
// Which is a little odd, maybe use time with a possibly different fFactor.
type Wave struct {
	count  int
	series []series
}

// New creates a Wave generator.
func New() *Wave {

	return &Wave{
		series: []series{
			simple(),
			threeRandom(),
			square(),
		},
	}

}

// Collect collects stats.
func (wv *Wave) Collect(ts time.Time) (pa entity.PointsAt, err error) {

	// Todo: wring hands about ref to collector elsewheres

	points := make([]entity.Point, len(wv.series))
	for i, srs := range wv.series {

		var val float64
		for _, wave := range srs.waves {
			val += wave.amplitude * math.Sin((wave.frequency/fFactor)*float64(wv.count)+wave.phase)
			// https://en.wikipedia.org/wiki/Phase_(waves)#Formula_for_phase_of_an_oscillation_or_a_periodic_signal
		}

		points[i] = entity.Point{
			Name:   "sine",
			Desc:   "Sine wave(s)",
			Type:   "gauge",
			Labels: entity.Labels{{Key: "name", Val: srs.name}},
			Value:  entity.Float{Data: val},
		}
	}

	pa = entity.PointsAt{
		Name:   name,
		Stamp:  ts,
		Points: points,
	}

	wv.count++
	return
}

// unexported

type series struct {
	name  string
	waves []wave
}

type wave struct {
	frequency float64
	amplitude float64
	phase     float64
}

func simple() series {
	return series{
		name: "simple",
		waves: []wave{
			{
				frequency: 1,
				amplitude: 1,
				phase:     0,
			},
		},
	}
}

func threeRandom() series {

	waves := make([]wave, 3)
	for i := 0; i < 3; i++ {

		n := float64(i) + 1
		waves[i] = wave{
			frequency: 1 + rand.Float64()*n,
			amplitude: rand.Float64() / n,
			phase:     rand.Float64() * 2 * math.Pi,
		}
	}

	return series{
		name:  "three_random",
		waves: waves,
	}
}

func square() series {

	waves := make([]wave, 4)
	for i := 0; i < 4; i++ {

		// https://mathworld.wolfram.com/FourierSeriesSquareWave.html
		n := 2*float64(i) + 1
		waves[i] = wave{
			frequency: n,
			amplitude: 1 / n,
			phase:     0,
		}
	}

	return series{
		name:  "square",
		waves: waves,
	}
}
