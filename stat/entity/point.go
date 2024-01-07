package entity

// Todo: namespace entity/metric(stat)

import (
	"fmt"
	"time"
)

// Value specifies Stringer.
type Value interface {
	fmt.Stringer
}

// Uint holds uint64 values.
type Uint struct {
	Data uint64
}

// String implements Stringer.
func (val Uint) String() string {
	return fmt.Sprintf("%d", val.Data)
}

// Float holds float64 values.
type Float struct {
	Data float64
}

// String implements Stringer.
func (val Float) String() string {
	return fmt.Sprintf("%.2f", val.Data)
}

// Label is a key/val pair associated with a point or points.
type Label struct {
	Key string
	Val string
}

// Labels is a multiplicity of Label.
type Labels []Label

// Point represents a collected stat.
type Point struct {
	Name   string
	Desc   string
	Unit   string
	Type   string
	Labels Labels
	Value  Value
}

// PointsAt are points with a common root name and labels collected at the same time.
type PointsAt struct {
	Name   string
	Stamp  time.Time
	Labels Labels
	Points []Point
}

// Stats are a collection unrelated PointsAt.
type Stats []PointsAt

// Todo: wring hands over name and ref
