# hcmut-co1027

Build

```bash
docker build -t hcmut-co1027 .
docker stop hcmut-co1027

```

Stop and remove

```bash
docker stop hcmut-co1027
docker rm hcmut-co1027
```

Start

```bash
docker run --name hcmut-co1027 -v /root/hcmut-co1027/cases:/usr/src/app/cases -p 80:8080 hcmut-co1027
```
