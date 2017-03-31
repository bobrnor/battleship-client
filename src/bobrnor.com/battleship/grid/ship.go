package grid

import (
	"math"
	"math/rand"
)

// ship ...
type ship struct {
	Left   int
	Right  int
	Top    int
	Bottom int
	Size   int

	indexOfX     int
	indexOfY     int
	indexOfAngle int

	randomizedXs     []int
	randomizedYs     []int
	randomizedAngels []int
}

func NewShip(size int) *ship {
	s := ship{
		Size: size,

		indexOfX:     0,
		indexOfY:     0,
		indexOfAngle: 0,

		randomizedXs:     rand.Perm(10),
		randomizedYs:     rand.Perm(10),
		randomizedAngels: rand.Perm(4),
	}
	s.updateCoordinates()
	return &s
}

func (s *ship) FindNextPosition() bool {
	for true {
		s.indexOfAngle++
		if s.indexOfAngle > 3 {
			s.indexOfAngle = 0
			s.indexOfY++
		}

		if s.indexOfY > 9 {
			s.indexOfY = 0
			s.indexOfX++
		}

		if s.indexOfX > 9 {
			// no more position
			s.Reset()
			s.updateCoordinates()
			return false
		}

		s.updateCoordinates()
		return true
	}
	return false
}

func (s *ship) Reset() {
	s.indexOfX = 0
	s.indexOfY = 0
	s.indexOfAngle = 0
}

func (s *ship) updateCoordinates() {
	angle := float64(s.randomizedAngels[s.indexOfAngle]) * 0.5 * math.Pi

	x0 := float64(s.randomizedXs[s.indexOfX])
	x1 := float64(x0 + float64((s.Size-1)*int(math.Cos(angle))))

	y0 := float64(s.randomizedYs[s.indexOfY])
	y1 := float64(y0 + float64((s.Size-1)*int(math.Sin(angle))))

	s.Left = int(math.Min(x0, x1))
	s.Right = int(math.Max(x0, x1))
	s.Top = int(math.Min(y0, y1))
	s.Bottom = int(math.Max(y0, y1))
}

func (s *ship) Intersect(ship *ship) bool {
	left := s.Left - 1
	right := s.Right + 1
	top := s.Top - 1
	bottom := s.Bottom + 1

	leftCollision := ship.Left >= left && ship.Left <= right
	rightCollision := ship.Right >= left && ship.Right <= right
	topCollision := ship.Top >= top && ship.Top <= bottom
	bottomCollision := ship.Bottom >= top && ship.Bottom <= bottom
	return (leftCollision || rightCollision) && (topCollision || bottomCollision)
}

func (s *ship) CheckHit(x, y int) bool {
	return s.Left <= x && s.Right >= x && s.Top <= y && s.Bottom >= y
}

func (s *ship) IsValid() bool {
	return s.Left >= 0 && s.Right < 10 && s.Top >= 0 && s.Bottom < 10
}
