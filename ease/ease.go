package ease

func EaseOutQuad(t float64) float64 {
	return t * (2 - t)
}

func EaseOutCubic(t float64) float64 {
	t = t - 1
	return t*t*t + 1
}

func EaseOutQuart(t float64) float64 {
	t = t - 1
	return 1 - t*t*t*t
}

func EaseOutQuint(t float64) float64 {
	t = t - 1
	return 1 + t*t*t*t*t
}

func EaseInQuad(t float64) float64 {
	return t * t
}

func EaseInCubic(t float64) float64 {
	return t * t * t
}

func EaseInQuart(t float64) float64 {
	return t * t * t * t
}

func EaseInQuint(t float64) float64 {
	return t * t * t * t * t
}
