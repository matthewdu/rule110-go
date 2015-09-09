package rule110

import (
  "fmt"
  "net/http"
  "regexp"
  "strconv"
)

func main() {
  http.HandleFunc("/", Rule110Handler)
  http.ListenAndServe(":3000", nil)
}

func init() {
  http.HandleFunc("/", Rule110Handler)
}

func Rule110Handler(w http.ResponseWriter, r *http.Request) {
//  fmt.Fprint(w, r.RequestURI)
  re := regexp.MustCompile("/\\d*")
  i, err := strconv.Atoi(re.FindString(r.RequestURI)[1:])
  if err != nil {
    fmt.Fprint(w, "Please pass in number of rows")
    return
  }

//  arr := make([][]bool, i)
//  for row := range arr {
//    arr[row] = make([]bool, i)
//    arr[row][i-row-1] = true
//    arr[row][i-1] = true
//    if row != 0 {
//      for j := i-row-1; j < i-1; j++ {
//        if j == 0 {
//          continue
//        }
//        arr[row][j] = Rule(arr[row-1][j-1], arr[row-1][j], arr[row-1][j+1])
//      }
//    }
//  }
//
//  for row := range arr {
//    for col := range arr[row] {
//      if arr[row][col] {
//        fmt.Fprint(w, "█")
//      } else {
//        fmt.Fprint(w, " ")
//      }
//    }
//    fmt.Fprintln(w, "")
//  }

  arr := [2][]bool{
    make([]bool, i),
    make([]bool, i),
  }
  arr[0][i-1] = true
  arr[1][i-1] = true
  for row := 0; row < i; row++ {
    arr[0][i-row-1] = true
    if row != 0 {
      for j := i-row-1; j < i-1; j++ {
        if j == 0 { continue }
        arr[1][j] = Rule(arr[0][j-1], arr[0][j], arr[0][j+1])
      }
    }
    for j := range arr[1] {
      arr[0][j] = arr[1][j]
      if arr[1][j] {
        fmt.Fprint(w, "█")
      } else {
        fmt.Fprint(w, " ")
      }
    }
    fmt.Fprintln(w, "")
  }

  fmt.Fprintln(w, i)
}

func Rule(i, j, k bool) bool {
  arr := [2][2][2]bool{
    [2][2]bool{
      [2]bool{false, true},
      [2]bool{true, true},
    },
    [2][2]bool{
      [2]bool{false, true},
      [2]bool{true, false},
    },
  }
  var x, y, z int
  if i { x = 1 } else { x = 0 }
  if j { y = 1 } else { y = 0 }
  if k { z = 1 } else { z = 0 }

  return arr[x][y][z]
}

