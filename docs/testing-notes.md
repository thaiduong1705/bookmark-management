# Ghi chú test (hexagonal)

## Mỗi tầng test bằng gì

| Tầng | Loại | Dependency xử lý |
|---|---|---|
| Repo (adapter ra) | Unit/integration nhẹ | miniredis thay Redis thật |
| Service | Unit | Mock repo (mockery) |
| Handler (adapter vào) | Unit | Mock service + `httptest` + `gin.CreateTestContext` |
| Cả luồng | Integration | `api.NewEngine` thật + miniredis, gửi HTTP qua `ServeHTTP` |

## Hỏi gì khi nghĩ testcase

1. Happy path
2. Input sai (rỗng, sai định dạng, âm)
3. Biên (0, 1, N, N+1)
4. Mỗi dependency lỗi thì sao?
5. Trạng thái dữ liệu (có / không có / hết hạn / trùng)
6. Side-effect: có ghi đúng không? Có gọi thứ không nên gọi không?

## Mock (mockery + testify)

- `.On("Method", args...).Return(...)` khai báo mong đợi.
- `mock.Anything`: giá trị nào cũng được. `mock.MatchedBy(func(x T) bool {...})`: tự viết điều kiện (dùng cho giá trị ngẫu nhiên).
- `NewXxx(t)` tự kiểm tra mọi `.On` đã được gọi. Case không được gọi service (vd 400) thì **không** khai báo `.On`.
- `AssertCalled(...)` nối đầu vào với kết quả trả về (vd code trả ra = code đã lưu).

## Lỗi đã gặp

- Mock thiếu tham số so với chữ ký thật -> `Unexpected Method Call`.
- `errors.New("x")` tạo 2 lần là 2 lỗi khác nhau -> `ErrorIs` fail. Dùng 1 biến chung.
- `mock.Anything` chỉ có nghĩa trong `.On`, không dùng trong `assert.Equal`.
- Handler test: quên gán `ginCtx.Request` -> bind luôn lỗi, trả 400.
- Handler dùng `c.Request.Context()` nhưng mock chờ `ginCtx` -> ctx không khớp. Dùng `mock.Anything`.
- `ShouldBindQuery` đọc tag `form`, không phải `json`.
- `expectedResp` phải là JSON hợp lệ (có `{}`) khi dùng `JSONEq`.
- Integration: gõ sai route (`link` vs `links`) -> 404 mặc định của gin. Kiểm tra cả body, đừng chỉ kiểm tra status.
- miniredis không tự hết hạn theo giờ thật -> dùng `mr.FastForward(...)`.
- Gọi `rdb.Close()` rồi dùng -> `redis.ErrClosed` (cách giả lập "Redis chết").

## Thiết kế rút ra

- `redis.Nil` không được rò ra ngoài repo -> repo đổi thành `domain.ErrNotFound`, handler map sang 404.
- Interface nhỏ -> mock dễ. Dependency truyền qua constructor.
- Phá code thử xem test có đỏ không: cách nhanh nhất biết test có bảo vệ được gì.

## Việc còn lại

- [ ] Integration: kịch bản tạo rồi lấy lại (POST lấy code, GET đúng code đó)
- [ ] Tách `pkg/redis/mock.go` ra package `redistest`
- [ ] `go mod tidy`

## Ghi chú của tôi

(tự viết thêm ở đây)
