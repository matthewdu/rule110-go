package rule110

import (
  "fmt"
  "image"
  "image/color/palette"
  "image/png"
  "net/http"
  "regexp"
  "strconv"

  "github.com/go-zoo/bone"
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

const (
  blackIndex uint8 = 0
  whiteIndex uint8 = 215
)

func main() {
  mux := bone.New()

  mux.GetFunc("/image/:rows", ImageHandler)
  mux.GetFunc("/html/:rows", HtmlHandler)

  http.ListenAndServe(":3000", mux)
}

func init() {
  mux := bone.New()

  mux.GetFunc("/image/:rows", ImageHandler)
  mux.GetFunc("/html/:rows", HtmlHandler)

  http.Handle("/", mux)
}

func rule(i, j, k bool) bool {
  var x, y, z int
  if i { x = 1 } else { x = 0 }
  if j { y = 1 } else { y = 0 }
  if k { z = 1 } else { z = 0 }

  return rules[x][y][z]
}

func HtmlHandler(w http.ResponseWriter, r *http.Request) {
  val := bone.GetValue(r, "rows")
  re := regexp.MustCompile("\\d*")
  i, err := strconv.Atoi(re.FindString(val))
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
        arr[1][j] = rule(arr[0][j-1], arr[0][j], arr[0][j+1])
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
  val := bone.GetValue(r, "rows")
  re := regexp.MustCompile("\\d*")
  i, err := strconv.Atoi(re.FindString(val))
  if err != nil {
    fmt.Fprint(w, "Please pass in number of row")
    return
  }

  m := image.NewPaletted(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{i, i}}, palette.WebSafe)

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
        arr[1][j] = rule(arr[0][j-1], arr[0][j], arr[0][j+1])
      }
    }
    for j := range arr[1] {
      arr[0][j] = arr[1][j]
      if arr[1][j] {
        m.SetColorIndex(j, row, blackIndex)
      } else {
        m.SetColorIndex(j, row, whiteIndex)
      }
    }
  }
  w.Header().Set("Content-type", "image/png")
  w.Header().Set("Cache-control", "public, max-age=259200")
  png.Encode(w, m)
}


