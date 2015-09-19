package rule110

import (
  "fmt"
  "image"
  "image/color"
  "image/png"
  "io/ioutil"
  "net/http"
  "regexp"
  "strconv"

  "github.com/go-zoo/bone"
  "google.golang.org/cloud/storage"
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
  mux.GetFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "http://github.com/matthewdu/rule110-go", 301)
  }) 

  http.Handle("/", mux)
}

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

var rulesUint8 = [2][2][2]uint8{
  [2][2]uint8{
    [2]uint8{0, 1},
    [2]uint8{1, 1},
  },
  [2][2]uint8{
    [2]uint8{0, 1},
    [2]uint8{1, 0},
  },
}

func rule(i, j, k bool) bool {
  var x, y, z int
  if i { x = 1 } else { x = 0 }
  if j { y = 1 } else { y = 0 }
  if k { z = 1 } else { z = 0 }

  return rules[x][y][z]
}

// Black and white palette
var bwPalette []color.Color = []color.Color{
  color.RGBA{0xff, 0xff, 0xff, 0xff},
  color.RGBA{0x00, 0x00, 0x00, 0xff},
}

func HtmlHandler(w http.ResponseWriter, r *http.Request) {
  val := bone.GetValue(r, "rows")
  re := regexp.MustCompile("\\d*")
  i, err := strconv.Atoi(re.FindString(val))
  if err != nil {
    fmt.Fprint(w, "Please pass in number of rows")
    return
  }
  if i > 500 {
    fmt.Fprint(w, "Please pass in a row equal or smaller than 500")
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
  rowsStr := re.FindString(val)
  rows, err := strconv.Atoi(rowsStr)
  if err != nil {
    fmt.Fprint(w, "Please pass in number of row")
    return
  }
  if rows > 5000 {
    fmt.Fprint(w, "Please pass in a row equal or smaller than 5000")
    return
  }

  ctx, err := cloudAuthContext(r)
  if err == nil {
    rc, err := storage.NewReader(ctx, bucket, rowsStr + ".png")
    if err == nil {
      image, err := ioutil.ReadAll(rc)
      rc.Close()
      if err == nil {
        w.Header().Set("Content-type", "image/png")
        w.Header().Set("Cache-control", "public, max-age=259200")
        w.Write(image)
        return
      }
    }
  }

  m := image.NewPaletted(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{rows, rows}}, bwPalette)

  m.Pix[m.PixOffset(rows-1, 0)] = 1
  for row := 1; row < rows; row++ {
    m.Pix[m.PixOffset(rows-1, row)] = 1
    for j := rows-row; j < rows-1; j++ {
      mid := m.PixOffset(j, row-1)
      left, right := mid-1, mid+1
      m.Pix[m.PixOffset(j, row)] = rulesUint8[m.Pix[left]][m.Pix[mid]][m.Pix[right]]
    }
  }
  w.Header().Set("Content-type", "image/png")
  w.Header().Set("Cache-control", "public, max-age=259200")
  png.Encode(w, m)
}

