# Reisen

4chan archive Koiwai-compatible frontend

## Features
- Not foolfuuka
- Looks better (thanks Tibix)
- Multiple default themes
- Tegaki support (unprecedented)
- Minimal JS
- Simple, customizable CSS
- Reasonably fast
- Reasonably secure
- Reasonably maintainable
- Reasonable SEO

## Configuration
See config.example.toml and follow along with the comments

Reisen will look for a configuration file
- Wherever the environment variable `REISEN_CONFIG` says a file is
- At ./config.toml wherever you started reisen from

## Deployment
Reisen needs a postgres user with `SELECT` on the posts table and network access to your Lnx instance

Strongly consider giving reisen its own user with negligible filesystem permissions

For DMCA reasons, reisen expects to always be run behind a reverse proxy, so configuring TLS in it will involve a bit of code
