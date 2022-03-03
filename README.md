# hcmut-co1027

Build

```bash
docker build -t hcmut-co1027 .
```

Stop and remove

```bash
docker stop hcmut-co1027
docker rm hcmut-co1027
```

Start

```bash
docker run --name hcmut-co1027 -v /root/hcmut-co1027/cases:/usr/src/app/cases -v /root/hcmut-co1027/archive:/usr/src/app/archive -d -p 80:8080 --env APP_URI=http://co1027.hoangvvo.com hcmut-co1027
```
