# hcmut-co1027

Bộ chấm bài lớp CO1027 - KTLT. Chấm bài | Upload test case

## Workflow

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
