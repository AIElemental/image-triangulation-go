package main

import (
  "image"
  "image/color"
  "github.com/AIElemental/image-triangulation-go/imagecreator"
  "github.com/AIElemental/image-triangulation-go/approximation"
  "fmt"
  "time"
  "log"
)

func main() {
  defer timeTrack(time.Now(), "main.main")
  const width, height = 8, 8

  // Create a colored image of the given width and height.
  img := image.NewNRGBA(image.Rect(0, 0, width, height))

  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
      seed := 255
      if (x - width / 2) * (x - width / 2) + (y - height / 2) * (y - height / 2) < 128*128 {
        seed = x + y
      }
      img.Set(x, y, color.NRGBA{
        R: uint8(seed & 255),
        G: uint8(seed << 1 & 255),
        B: uint8(seed << 2 & 255),
        A: 255,
      })
    }
  }
  imagecreator.StoreImage(img, "D:/go-image-tri/circle.png")

  bestApprox := approximation.Initial(img)

  approxImg := image.NewNRGBA(image.Rect(0, 0, width, height))
  approximation.Apply(approxImg, bestApprox)
  imagecreator.StoreImage(approxImg, fmt.Sprint("D:/go-image-tri/circle-approx-0.png"))

  version := 0
  shakes := make([]approximation.Approximation, 4)
  bestMinDiff := approximation.Diff(img, approxImg)

  const iterations = 1
  for range make([]bool, iterations) {
    start := time.Now()
    version++
    fmt.Println("iteration ", version)
    minDiffIdx := -1
    for i := range shakes {
      shakes[i] = approximation.Shake(img, bestApprox)
      shakes[i] = approximation.AdjustColors(img, shakes[i])
      shakesImage := image.NewNRGBA(image.Rect(0, 0, width, height))
      approximation.Apply(shakesImage, shakes[i])
      diffI := approximation.Diff(img, shakesImage)
      fmt.Println("Approx with diff", diffI)
      if diffI < bestMinDiff {
        bestMinDiff = diffI
        minDiffIdx = i
      }
    }
    if minDiffIdx >= 0 {
      fmt.Println("New best approx with diff", bestMinDiff)
      bestApprox = shakes[minDiffIdx]
    }
    updatedVariant := image.NewNRGBA(image.Rect(0, 0, width, height))
    approximation.Apply(updatedVariant, bestApprox)
    imagecreator.StoreImage(updatedVariant, fmt.Sprintf("D:/go-image-tri/circle-approx-%d.png", version))
    timeTrack(start, fmt.Sprintf("iteration %d", version))
  }

  lastVariant := image.NewNRGBA(image.Rect(0, 0, width, height))
  approximation.Apply(lastVariant, bestApprox)
  imagecreator.StoreImage(lastVariant, "D:/go-image-tri/circle-approx-final.png")
}

func timeTrack(start time.Time, name string) {
  elapsed := time.Since(start)
  log.Printf("%s took %s", name, elapsed)
}