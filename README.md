# Chủ đề
Viết một chương trình về crawler/scrapy với bất kì ngôn ngữ nào mà bạn đã biết, và đảm bảo các yêu cầu sau:
# Yêu cầu
1. Bạn tự chọn một đường dẫn của bất kì của một trang web tin tức (VNExpress, TuoiTre, TheSaiGonTime,...)
2. Các bạn phải download được nội dung (crawl) và phân tích (parse) bài viết để lấy được: tiêu đề, tác giả, ngày xuất bản.
3. Từ nội dung bài viết đó các bạn phải tìm kiếm xem nó còn có các đường dẫn đến bài viết nào nữa hay không, với mỗi đường dẫn tìm được, tiếp tục thực hiện bước 1, và 2.

# Run project
go run main.go <link_crawler>
ex: go run main.go https://www.thesaigontimes.vn/121624/Cuoc-cach-mang-dau-khi-da-phien
# Link Refer
https://github.com/jackdanger/collectlinks
https://godoc.org/github.com/PuerkitoBio/goquery
https://www.reddit.com/r/golang/comments/3fcabt/question_read_value_from_html_input_tag/
