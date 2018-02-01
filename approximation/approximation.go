package approximation

import (
  "image/color"
  "image"
  "time"
  "log"
  "math"
  "math/rand"
)

type Approximation struct {
  Triangles []Triangle
}

type Triangle struct {
  V1, V2, V3 Point
  Color      color.Color
}

type Point struct {
  X, Y int
}

func sign(p1, p2, p3 Point) float64 {
  return float64((p1.X - p3.X) * (p2.Y - p3.Y) - (p2.X - p3.X) * (p1.Y - p3.Y))
}

func (t Triangle) Contains(pt Point) bool {
  var b1, b2, b3 bool;
  b1 = sign(pt, t.V1, t.V2) < 0;
  b2 = sign(pt, t.V2, t.V3) < 0;
  b3 = sign(pt, t.V3, t.V1) < 0;
  return ((b1 == b2) && (b2 == b3));
}

func (p Point) move(x, y int) Point {
  return Point{p.X + x, p.Y + y}
}

func Apply(img *image.NRGBA, approx Approximation) {
  defer timeTrack(time.Now(), "approximation.Apply")
  for _, v := range approx.Triangles {
    applySingle(img, v)
  }
}

func applySingle(img *image.NRGBA, triang Triangle) {
  width, height := img.Bounds().Dx(), img.Bounds().Dy()
  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
      if triang.Contains(Point{x, y}) {
        img.Set(x, y, triang.Color)
      }
    }
  }
}

func AdjustColors(img *image.NRGBA, approx Approximation) Approximation {
  defer timeTrack(time.Now(), "approximation.AdjustColors")
  for i, v := range approx.Triangles {
    approx.Triangles[i].Color = avgColor(img, v)
  }
  return approx
}

func avgColor(img *image.NRGBA, triang Triangle) color.Color {
  width, height := img.Bounds().Dx(), img.Bounds().Dy()
  var rSum, gSum, bSum float64
  var count float64
  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
      if triang.Contains(Point{x, y}) {
        var r, g, b uint8
        var point color.NRGBA = img.At(x, y).(color.NRGBA)
        r, g, b = point.R, point.G, point.B
        rSum += float64(r)
        gSum += float64(g)
        bSum += float64(b)
        count++
      }
    }
  }
  if count > 0 {
    return color.NRGBA{uint8(rSum / count), uint8(gSum / count), uint8(bSum / count), 255}
  } else {
    return color.Black
  }
}

func Diff(base, img *image.NRGBA) float64 {
  defer timeTrack(time.Now(), "approximation.Diff")
  width, height := img.Bounds().Dx(), img.Bounds().Dy()
  var sum float64
  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
      var r1, g1, b1 uint32
      var r2, g2, b2 uint32
      r1, g1, b1, _ = img.At(x, y).RGBA()
      r2, g2, b2, _ = base.At(x, y).RGBA()
      sum += math.Abs(float64(r1 - r2)) + math.Abs(float64(g1 - g2)) + math.Abs(float64(b1 - b2))
    }
  }
  return sum
}

func Initial(img *image.NRGBA) Approximation {
  width, height := img.Bounds().Dx(), img.Bounds().Dy()
  xCells, yCells := 16, 16
  xStep, yStep := width / xCells, height / yCells
  triangles := make([]Triangle, 0)
  for x := 0; x < width; x += xStep {
    for y := 0; y < height; y += yStep {
      p1 := Point{x, y}
      p2 := Point{x + xStep, y}
      p3 := Point{x, y + yStep}
      triag := Triangle{p1, p2, p3, color.Black}
      triangles = append(triangles, triag)

      pp1 := Point{x + xStep, y}
      pp2 := Point{x + xStep, y + yStep}
      pp3 := Point{x, y + yStep}
      triagp := Triangle{pp1, pp2, pp3, color.Black}
      triangles = append(triangles, triagp)
    }
  }
  return AdjustColors(img, Approximation{triangles})
}

func Shake(img *image.NRGBA, approx Approximation) Approximation {
  defer timeTrack(time.Now(), "approximation.Shake")
  width, height := img.Bounds().Dx(), img.Bounds().Dy()
  for i, v := range approx.Triangles {
    approx.Triangles[i] = randomMoveTriangle(width, height, v)
  }
  return approx
}

func randomMoveTriangle(width, height int, t Triangle) Triangle {
  return Triangle{
    randomMovePoint(width, height, t.V1),
    randomMovePoint(width, height, t.V2),
    randomMovePoint(width, height, t.V3),
    t.Color,
  }
}

func randomMovePoint(width, height int, p Point) Point {
  return p.move(rand.Intn(width/64)-width/128, rand.Intn(height/64)-width/128)
}

func timeTrack(start time.Time, name string) {
  elapsed := time.Since(start)
  log.Printf("%s took %s", name, elapsed)
}
