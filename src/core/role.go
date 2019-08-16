package core

// TeamForce is
type TeamForce string

// BeatForce is beat team force
const BeatForce = TeamForce("beat")

// ThinkForce is think team force
const ThinkForce = TeamForce("think")

// String returns the string of it
func (f TeamForce) String() string {
	return string(f)
}

// Reverse return the opposite of it
func (f TeamForce) Reverse() TeamForce {
	if f == BeatForce {
		return ThinkForce
	}
	return BeatForce
}

// Equal compares the string expression of a team force
func (f TeamForce) Equal(o string) bool {
	return string(f) == o
}

// BattleMode is
type BattleMode string

// FireMode is fire battle mode
const FireMode = BattleMode("fire")

// WaterMode is water battle mode
const WaterMode = BattleMode("water")

// String returns the string
func (m BattleMode) String() string {
	return string(m)
}

// PowerForce returns the powerful team force
func (m BattleMode) PowerForce() TeamForce {
	if m == FireMode {
		return BeatForce
	}

	return ThinkForce
}
