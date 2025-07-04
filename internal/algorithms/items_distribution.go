package algorithms

import (
	"errors"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/roulette"
)

type ItemsDistribution struct {
	avgCost   int64
	intervals []ItemInterval
}

func NewIntervalArray(items []roulette.Item, avgCost int64) ItemsDistribution {
	result := make([]ItemInterval, len(items))
	for i, item := range items {
		result[i] = NewItemInterval(item)
	}
	return ItemsDistribution{
		intervals: result,
		avgCost:   avgCost,
	}
}

func (d ItemsDistribution) TotalIntervalLength() float64 {
	if len(d.intervals) == 0 {
		return 0
	}
	return d.intervals[len(d.intervals)-1].End - d.intervals[0].Begin
}

func (d ItemsDistribution) Roll() (*ItemInterval, error) {
	if len(d.intervals) == 0 {
		return nil, errors.New("can't roll empty intervals")
	}

	if len(d.intervals) == 1 {
		return &d.intervals[0], nil
	}

	randFloat, err := CryptoRandFloat64()
	if err != nil {
		return nil, err
	}

	p := randFloat * d.TotalIntervalLength()
	return d.findInverval(p), nil
}

func (d ItemsDistribution) findInverval(point float64) *ItemInterval {
	if len(d.intervals) == 0 {
		return nil
	}

	if len(d.intervals) == 1 {
		interval := d.intervals[0]
		if point >= interval.Begin && point <= interval.End {
			return &interval
		}
		return nil
	}

	lastIndx := len(d.intervals) - 1
	for i, interval := range d.intervals {
		if i == lastIndx { // aslo include end at last interval
			if point >= interval.Begin && point <= interval.End {
				return &d.intervals[i]
			}
			continue
		}
		if point >= interval.Begin && point < interval.End {
			return &d.intervals[i]
		}
	}

	return nil
}

func calculateWeigth(avgCost int64, itemCost int64) float64 {
	avgDiff := avgCost - itemCost
	if avgDiff == 0 {
		avgDiff = 1
	}
	if avgDiff < 0 {
		return 1 - float64(-avgDiff)/float64(itemCost)
	} else {
		return 1 + float64(avgDiff)/float64(itemCost)
	}
}

func (d ItemsDistribution) Split() {
	start := float64(0)
	for i := range d.intervals {
		d.intervals[i].Begin = start
		d.intervals[i].End = start + calculateWeigth(d.avgCost, d.intervals[i].Item.TotalCost)
		start = d.intervals[i].End
	}
	d.normalize()
}

func (d ItemsDistribution) EstimateProbabilities() []float64 {
	// edge case - no intervals, no chance to win
	if len(d.intervals) == 0 {
		return nil
	}

	// edge case - for single interval probability equals 1
	if len(d.intervals) == 1 {
		return []float64{1}
	}

	// len(invervals) > 1
	totalLength := d.TotalIntervalLength()
	var res = make([]float64, len(d.intervals))
	for indx, i := range d.intervals {
		res[indx] = 100 * i.Length() / totalLength
	}
	return res
}

func (d ItemsDistribution) Intervals() []ItemInterval {
	return d.intervals
}

func (d ItemsDistribution) normalize() {
	totalLength := d.TotalIntervalLength()
	for i := range d.intervals {
		d.intervals[i].Begin /= totalLength
		d.intervals[i].End /= totalLength
	}
}
