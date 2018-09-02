package lgimage

import (
	"image/color"

	"github.com/pkg/errors"
)

// ColorMapFunc is wrapper interface of ColorMap
type ColorMapFunc func(v Value) color.Color

func (c ColorMapFunc) Color(v Value) color.Color {
	return c(v)
}
func (c ColorMapFunc) List() []ExceptionColor {
	return []ExceptionColor{}
}

// ColorMap is encoder of to Color from Value
type ColorMap interface {
	Color(v Value) color.Color
	List() []ExceptionColor
}

// ExceptionColor declare the exception value to map to color
type ExceptionColor struct {
	String string
	Color  color.Color
}

// ValueMapWithFunc wrap the input value and to map to index[0->1]
type ValueMapWithFunc struct {
	Vmin, Vmax    float64
	ColorFunc     func(d float64) color.Color // 0 -> 1
	ExceptionList map[string]color.Color
}

func (m ValueMapWithFunc) Color(v Value) color.Color {
	switch x := v.Value().(type) {
	case float64:
		// over range
		if x > m.Vmax {
			if c, ok := m.ExceptionList["over"]; ok {
				return c
			}
			return color.NRGBA{0, 0, 0, 255}
		}
		if x < m.Vmin {
			if c, ok := m.ExceptionList["under"]; ok {
				return c
			}
			return color.NRGBA{255, 255, 255, 255}
		}

		vd := m.Vmax - m.Vmin
		d := (x - m.Vmin) / vd
		return m.ColorFunc(d)
	case string:
		if c, ok := m.ExceptionList[x]; ok {
			return c
		}
		return color.NRGBA{0, 0, 0, 0}
	}
	panic(errors.Errorf("Invalid type in ValueMapWithFunc [%v]", v))
}

func (m ValueMapWithFunc) List() []ExceptionColor {
	var list []ExceptionColor
	for k, v := range m.ExceptionList {
		list = append(list, ExceptionColor{k, v})
	}
	return list
}

// invalid data
// N/A, nil, Overflow

type Value struct {
	Float  float64
	String string
	Valid  bool
}

func NewValues(i ...interface{}) (values []Value, err error) {
	for _, v := range i {
		nv, err := NewValue(v)
		if err != nil {
			return values, err
		}
		values = append(values, nv)
	}
	return
}

func NewValue(i interface{}) (Value, error) {
	if i == nil {
		return Value{
			String: "nil",
		}, nil
	}
	switch v := i.(type) {
	case float64:
		return Value{
			Float: v,
			Valid: true,
		}, nil
	case float32:
		return Value{
			Float: float64(v),
			Valid: true,
		}, nil
	case int:
		return Value{
			Float: float64(v),
			Valid: true,
		}, nil
	case uint:
		return Value{
			Float: float64(v),
			Valid: true,
		}, nil
	case int32:
		return Value{
			Float: float64(v),
			Valid: true,
		}, nil
	case uint32:
		return Value{
			Float: float64(v),
			Valid: true,
		}, nil
	case int64:
		return Value{
			Float: float64(v),
			Valid: true,
		}, nil
	case uint64:
		return Value{
			Float: float64(v),
			Valid: true,
		}, nil
	case string:
		return Value{
			String: v,
		}, nil
	default:
		return Value{}, errors.Errorf("Value is can not parse input [%v]", v)
	}
}

func (v Value) Value() interface{} {
	if v.Valid {
		return v.Float
	}
	return v.String
}
