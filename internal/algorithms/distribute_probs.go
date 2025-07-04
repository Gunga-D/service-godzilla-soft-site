package algorithms

func (i ItemInterval) Length() float64 {
	return i.End - i.Begin
}

func SplitIntervals(weights []float64) ([]ItemInterval, error) {
	start := float64(0)
	intervals := make([]ItemInterval, len(weights))
	for indx, weight := range weights {
		end := weight
		intervals[indx] = ItemInterval{
			Begin: start,
			End:   start + end,
		}
		start += end
	}
	return intervals, nil
}
