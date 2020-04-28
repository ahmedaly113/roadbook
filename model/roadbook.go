package model

import "image/color"

// Waypoint represents an individual item within the roadbook
type Waypoint struct {
	Distance   float32
	Tulip      []byte
	Notes      []byte
	Background color.Color
	DZ			bool
	FZ 			bool
}

// Model contains the present Book and sensor state
type Model struct {
	Book     []Waypoint
	Idx      int
	Speed    float32
	SpeedLimit float32
	IsSpeedZone bool
	Heading  float32
	Distance float32
}
