package imagecreator

import (
  "image"
  "image/png"
  "os"
  "log"
)

func Show(f func(int, int) [][]uint8) {
  const (
    dx = 256
    dy = 256
  )
  data := f(dx, dy)
  m := image.NewNRGBA(image.Rect(0, 0, dx, dy))
  for y := 0; y < dy; y++ {
    for x := 0; x < dx; x++ {
      v := data[y][x]
      i := y*m.Stride + x*4
      m.Pix[i] = v
      m.Pix[i+1] = v
      m.Pix[i+2] = 255
      m.Pix[i+3] = 255
    }
  }
  StoreImage(m, "D:/go-image-tri/image.png")
}

func StoreImage(img image.Image, filename string) {
  f, err := os.Create(filename)
  if err != nil {
    log.Fatal(err)
  }

  if err := png.Encode(f, img); err != nil {
    f.Close()
    log.Fatal(err)
  }

  if err := f.Close(); err != nil {
    log.Fatal(err)
  }
}
