# kumo 🕷️

Crawl channels for new videos and queues them up into
[`tasq`](https://github.com/ragtag-archive/tasq).

## Configuration

See `config.example.toml`. A `tsv` file hosted on a public URL is used as the
source of the videos. The format of the `tsv` is as follows:

| Agency name      | Batch name | Channel ID | Channel Name | Cron Preset |
| ---------------- | ---------- | ---------- | ------------ | ----------- |
| VTuber Group LLC | Gen -1     | something  | some vtuber  | Low         |

## Deployment

You can use this `docker-compose.yaml` to run a pre-built image of kumo.

```yaml
version: '3'
services:
  kumo:
    image: ghcr.io/ragtag-archive/kumo:main
    restart: unless-stopped
    volumes:
      - ./config.toml:/app/config.toml:ro
```
