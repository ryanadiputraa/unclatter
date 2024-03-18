# UnClatter 📑

UnClatter is an article bookmarking tool that also removes distractions like ads and popups. It offers a clean reading experience, allowing users to focus solely on the essential content.

## Development 🧑🏻‍💻

Setup development database or use [docker compose](https://docs.docker.com/compose/) and start services

```bash
make compose-up
```

Adjust your `config/config.yml` for development server

```bash
make env
```

Start server

```bash
make server
```

Run test cases

```bash
make test
```

Use [mockery](https://github.com/vektra/mockery) for mocking
```bash
mockery --dir <interface_dir> --name=<interface_name> --filename=<out_mock_fil> --output=<out_mock_dir> --outpkg=<mock_pkg>
```

## API Spec 📝

API Specification can be seen on `/docs` endpoint.
