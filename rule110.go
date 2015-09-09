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

  arr := make([][]int, i)
  for row := range arr {
    arr[row] = make([]int, i)
    arr[row][i-row-1] = 1
    arr[row][i-1] = 1
    if row != 0 {
      for j := i-row-1; j < i-1; j++ {
        if j == 0 {
          continue
        }
        arr[row][j] = Rule(arr[row-1][j-1], arr[row-1][j], arr[row-1][j+1])
      }
    }
  }

  for row := range arr {
    for col := range arr[row] {
      if arr[row][col] == 1 {
        fmt.Fprint(w, "█")
      } else {
        fmt.Fprint(w, " ")
      }
    }
    fmt.Fprintln(w, "")
  }
  fmt.Fprintln(w, i)
}

func Rule(i, j, k int) int {
  arr := [2][2][2]int{
    [2][2]int{
      [2]int{0, 1},
      [2]int{1, 1},
    },
    [2][2]int{
      [2]int{0, 1},
      [2]int{1, 0},
    },
  }
  return arr[i][j][k]
}

