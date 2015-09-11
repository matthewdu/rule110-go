package rule110

import (
  "fmt"
  "image"
  "image/color"
  "image/png"
  "net/http"
  "regexp"
  "strconv"
)

var rules = [2][2][2]bool{
  [2][2]bool{
    [2]bool{false, true},
    [2]bool{true, true},
  },
  [2][2]bool{
    [2]bool{false, true},
    [2]bool{true, false},
  },
}

func main() {
  http.HandleFunc("/image", ImageHandler)
  http.HandleFunc("/html", Rule110Handler)
  http.ListenAndServe(":3000", nil)
}

func init() {
  http.HandleFunc("/image/", ImageHandler)
  http.HandleFunc("/html/", Rule110Handler)
}

func Rule110Handler(w http.ResponseWriter, r *http.Request) {
  re := regexp.MustCompile("/html/\\d*")
  i, err := strconv.Atoi(re.FindString(r.RequestURI)[6:])
  if err != nil {
    fmt.Fprint(w, "Please pass in number of rows")
    return
  }

  fmt.Fprint(w,
    `<html>
        <head>
          <style>
            .b {
              display: inline-block;
              position: relative;
              height: 4px;
              width: 4px;
              background: black;
            }
            .w {
              display: inline-block;
              position: relative;
              height: 4px;
              width: 4px;
              background: white;
            }
          </style>
        </head>
        <body>
        </body>
    </html>`)

  arr := [2][]bool{
    make([]bool, i),
    make([]bool, i),
  }
  arr[0][i-1] = true
  arr[1][i-1] = true
  for row := 1; row < i; row++ {
    arr[0][i-row-1] = true
    if row != 0 {
      for j := i-row; j < i-1; j++ {
        if j == 0 { continue }
        arr[1][j] = Rule(arr[0][j-1], arr[0][j], arr[0][j+1])
      }
    }
    fmt.Fprint(w, "<div>")
    for j := range arr[1] {
      arr[0][j] = arr[1][j]
      if arr[1][j] {
        fmt.Fprint(w, "<div class=\"b\"></div>")
      } else {
        fmt.Fprint(w, "<div class=\"w\"></div>")
      }
    }
    fmt.Fprint(w, "</div>")
  }

  fmt.Fprintln(w, `
        </body>
    </html>`)
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
  re := regexp.MustCompile("/image/\\d*")
  i, err := strconv.Atoi(re.FindString(r.RequestURI)[7:])
  if err != nil {
    fmt.Fprint(w, "Please pass in number of row")
    return
  }

  m := image.NewRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{i, i}})
//  for y := 0; y < i; y++ {
//    for x := 0; x < i; x++ {
//      m.SetRGBA(x, y, color.RGBA{uint8(x), uint8((x + y) / 2), uint8(y), 255})
//    }
//  }

  arr := [2][]bool{
    make([]bool, i),
    make([]bool, i),
  }
  arr[0][i-1] = true
  arr[1][i-1] = true
  for row := 1; row < i; row++ {
    arr[0][i-row-1] = true
    if row != 0 {
      for j := i-row; j < i-1; j++ {
        if j == 0 { continue }
        arr[1][j] = Rule(arr[0][j-1], arr[0][j], arr[0][j+1])
      }
    }
    for j := range arr[1] {
      arr[0][j] = arr[1][j]
      if arr[1][j] {
        m.SetRGBA(j, row, color.RGBA{0, 0, 0, 255})
      } else {
        m.SetRGBA(j, row, color.RGBA{255, 255, 255, 255})
      }
    }
  }
  w.Header().Set("Content-type", "image/png")
  w.Header().Set("Cache-control", "public, max-age=259200")
  png.Encode(w, m)
}

func Rule(i, j, k bool) bool {
  var x, y, z int
  if i { x = 1 } else { x = 0 }
  if j { y = 1 } else { y = 0 }
  if k { z = 1 } else { z = 0 }

  return rules[x][y][z]
}

