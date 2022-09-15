# Reisen

4chan Archive Koiwai-compatible frontend

## Features
- Support for multiple pure CSS themes
- Tegaki support
- Latex support
- Country flags support
- Board flags support (for both /pol/ and /mlp/)
- Colored IDs support
- EXIF tables support
- ShiftJIS support

## Configuration
See config.example.toml and follow along with the comments

Reisen will look for a configuration file

- wherever the environment variable `REISEN_CONFIG` says a file is
- at ./config.toml wherever you started reisen from

in that order

## Deployment
Reisen needs a postgres user with `SELECT` on 
the posts table, `SELECT`, `INSERT` and `UPDATE`
permissions on the reports table, network access
to your Lnx instance and just enough filesystem permissions
to write the robots.txt file at ./static/robots.txt

Strongly consider giving reisen its own user with negligible filesystem permissions,
preferrably just the robots.txt file

For DMCA reasons, Reisen expects to always be run behind a reverse proxy,
so configuring TLS in it directly will involve a (tiny) bit of code,
see the documentation for the echo framework
