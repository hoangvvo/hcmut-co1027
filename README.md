# hcmut-co1027

Bộ chấm bài lớp CO1027 - KTLT. Bộ test tự upload bởi người dùng / hỗ trợ chấm các bài tập lớn khác nhau (theo tên file) nên có thể dùng cho nhiều năm sau :)

Website: http://co1027.hoangvvo.com/ (mình sẽ đóng server nếu không có bài tập lớn để tiết kiệm chi phí, liên hệ mình để mở lại nếu cần. Ngoài ra nếu được các bạn có thể [donate](DONATE.md))

## Workflow

Chương trình nên được deploy bằng [docker](https://www.docker.com/). Nếu chạy trực tiếp trên host có thể dẫn tới vấn đề bảo mật vì chương trình cho phép thực thi code bất kì!

### Development

```bash
go run main.go
```

### Build

```bash
docker build -t hcmut-co1027 .
```

### Stop and remove

```bash
docker stop hcmut-co1027
docker rm hcmut-co1027
```

It is a good idea to save the logs somewhere:

```bash
docker logs hcmut-co1027 > /root/hcmut-co1027/logs
```

### Start

```bash
docker run --name hcmut-co1027 -v /root/hcmut-co1027/cases:/usr/src/app/cases -v /root/hcmut-co1027/archive:/usr/src/app/archive -d -p 80:8080 --env APP_URI=http://co1027.hoangvvo.com hcmut-co1027
```

## Contribute

Tạo PR để fix bug/thêm feature. Hoặc tạo issue để báo lỗi.

## LICENSE

[MIT](LICENSE)

Nếu các bạn fork và upload lên trang cá nhân. Vui lòng để credit và license gốc, cảm ơn!
