# blockchain-platform

blockchain network platform

## Tech Stacks

- Go

---

## Development

- Copy `node/config/config.example.yml` to `node/config/config.yml`, then adjust your env configuration

- Adjust nodes service in `compose.yml`

- Start network (use `--build` tag to apply changes):

```bash
docker-compose up -d
```

- Stop network:

```bash
docker-compose down
```
