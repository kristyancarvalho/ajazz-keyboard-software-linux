package keyboard

type LightingMode uint8

const (
	LED_OFF LightingMode = iota
	STATIC
	SINGLE_ON
	SINGLE_OFF
	GLITTERING
	FALLING
	COLOURFUL
	BREATH
	SPECTRUM
	OUTWARD
	SCROLLING
	ROLLING
	ROTATING
	EXPLODE
	LAUNCH
	RIPPLES
	FLOWING
	PULSATING
	TILT
	SHUTTLE
)

type Direction uint8

const (
	LEFT Direction = iota
	DOWN
	UP
	RIGHT
)

var allModes = []LightingMode{
	LED_OFF,
	STATIC,
	SINGLE_ON,
	SINGLE_OFF,
	GLITTERING,
	FALLING,
	COLOURFUL,
	BREATH,
	SPECTRUM,
	OUTWARD,
	SCROLLING,
	ROLLING,
	ROTATING,
	EXPLODE,
	LAUNCH,
	RIPPLES,
	FLOWING,
	PULSATING,
	TILT,
	SHUTTLE,
}

var modeNames = map[LightingMode]string{
	LED_OFF:    "LED Off",
	STATIC:     "Static",
	SINGLE_ON:  "Single On",
	SINGLE_OFF: "Single Off",
	GLITTERING: "Glittering",
	FALLING:    "Falling",
	COLOURFUL:  "Colourful",
	BREATH:     "Breath",
	SPECTRUM:   "Spectrum",
	OUTWARD:    "Outward",
	SCROLLING:  "Scrolling",
	ROLLING:    "Rolling",
	ROTATING:   "Rotating",
	EXPLODE:    "Explode",
	LAUNCH:     "Launch",
	RIPPLES:    "Ripples",
	FLOWING:    "Flowing",
	PULSATING:  "Pulsating",
	TILT:       "Tilt",
	SHUTTLE:    "Shuttle",
}

var directionNames = map[Direction]string{
	LEFT:  "Left",
	DOWN:  "Down",
	UP:    "Up",
	RIGHT: "Right",
}

const (
	LEDOff    = LED_OFF
	Static    = STATIC
	SingleOn  = SINGLE_ON
	SingleOff = SINGLE_OFF
	Breath    = BREATH
	Spectrum  = SPECTRUM
	Outward   = OUTWARD
	Scrolling = SCROLLING
	Rolling   = ROLLING
	Rotating  = ROTATING
	Explode   = EXPLODE
	Launch    = LAUNCH
	Ripples   = RIPPLES
	Flowing   = FLOWING
	Pulsating = PULSATING
	Tilt      = TILT
	Shuttle   = SHUTTLE
	Left      = LEFT
	Down      = DOWN
	Up        = UP
	Right     = RIGHT
)

func (m LightingMode) String() string {
	if name, ok := modeNames[m]; ok {
		return name
	}
	return "Unknown"
}

func (LightingMode) All() []LightingMode {
	return append([]LightingMode(nil), allModes...)
}

func (d Direction) String() string {
	if name, ok := directionNames[d]; ok {
		return name
	}
	return "Unknown"
}

func DirectionsForMode(mode LightingMode) []Direction {
	switch mode {
	case SCROLLING:
		return []Direction{UP, DOWN}
	case ROLLING, FLOWING, TILT:
		return []Direction{LEFT, RIGHT}
	default:
		return []Direction{LEFT, DOWN, UP, RIGHT}
	}
}

func DefaultDirection(mode LightingMode) Direction {
	options := DirectionsForMode(mode)
	if len(options) == 0 {
		return LEFT
	}
	return options[0]
}
