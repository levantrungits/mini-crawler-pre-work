package crawler

import (
  "github.com/jackdanger/collectlinks"
  "net/http"
  "fmt"
)

func OkGaden() {
  resp, _ := http.Get("https://www.thesaigontimes.vn/td/288229/hoang-mang-mua-hang-truc-tuyen-nhung-khong-duoc-kiem-hang.html")
  links := collectlinks.All(resp.Body)
  fmt.Println(links)
}