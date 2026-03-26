package out

// AVRational is a rational number (pair of numerator and denominator).
// Matches the C struct layout for pass-by-value through purego.
type AVRational struct {
	Num int32
	Den int32
}
