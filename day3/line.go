package main

type Point struct {
	X int64
	Y int64
}

type Line struct {
	A     Point
	B     Point
	Sp    Point
	Dir   string
	Len   int64
	Steps int64
}

func NewPoint(sp Point, direction string, steps int64) Point {
	var ep Point
	ep = sp

	switch direction {
	case "R":
		ep.X += steps
	case "L":
		ep.X -= steps
	case "U":
		ep.Y += steps
	case "D":
		ep.Y -= steps
	}

	return ep
}

func NewLine(sp Point, ep Point, steps int64) Line {
	var dir string
	var a, b Point
	var len int64

	a = sp
	b = ep

	if sp.X == ep.X {
		dir = "vertical"
		if sp.Y > ep.Y {
			a = ep
			b = sp
		}
		len = b.Y - a.Y
	} else if sp.Y == ep.Y {
		dir = "horizontal"
		if sp.X > ep.X {
			a = ep
			b = sp
		}
		len = b.X - a.X
	} else {
		dir = "normal"
	}
	return Line{
		A:     a,
		B:     b,
		Dir:   dir,
		Sp:    sp,
		Len:   len,
		Steps: steps,
	}
}

func (l Line) Contains(p Point) bool {
	if l.Dir == "vertical" {
		if p.X == l.A.X &&
			p.Y >= l.A.Y && p.Y <= l.B.Y {
			return true
		}
	} else if l.Dir == "horizontal" {
		if p.Y == l.A.Y &&
			p.X >= l.A.X && p.X <= l.B.X {
			return true
		}
	}
	return false
}

func (l Line) StepsTake(p Point) int64 {
	var s int64
	if l.Dir == "vertical" {
		if p.Y > l.Sp.Y {
			s = p.Y - l.Sp.Y
		} else {
			s = l.Sp.Y - p.Y
		}
	} else if l.Dir == "horizontal" {
		if p.X > l.Sp.X {
			s = p.X - l.Sp.X
		} else {
			s = l.Sp.X - p.X
		}
	}
	return l.Steps + s
}
