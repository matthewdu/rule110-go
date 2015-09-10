package rule110

import (
  "fmt"
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
  http.HandleFunc("/", Rule110Handler)
  http.ListenAndServe(":3000", nil)
}

func init() {
  http.HandleFunc("/", Rule110Handler)
}

func Rule110Handler(w http.ResponseWriter, r *http.Request) {
  re := regexp.MustCompile("/\\d*")
  i, err := strconv.Atoi(re.FindString(r.RequestURI)[1:])
  if err != nil {
    fmt.Fprint(w, "Please pass in number of rows")
    return
  }

  fmt.Fprint(w,
    `<html>
        <head>
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
        fmt.Fprint(w, "<div style=\"display: inline-block; position:relative; height:4px; width:4px; background:black;\"></div>")
      } else {
        fmt.Fprint(w, "<div style=\"display: inline-block; position:relative; height:4px; width:4px; background:white;\"></div>")
      }
    }
    fmt.Fprint(w, "</div>")
  }

  fmt.Fprintln(w, `
        </body>
    </html>`)
}

func Rule(i, j, k bool) bool {
  var x, y, z int
  if i { x = 1 } else { x = 0 }
  if j { y = 1 } else { y = 0 }
  if k { z = 1 } else { z = 0 }

  return rules[x][y][z]
}

