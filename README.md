# reisen

4chan Koiwai-compatible frontend

## Features
- Not foolfuuka
- Looks better (thanks Tibix)
- Tegaki support (unprecedented)

## Tailwind

We're using tailwind for CSS. To generate tailwind.css
install tailwind with `npm i -g tailwind` and run
`tailwind -o static/tailwind.css`.

## Configuration
* PostgresConfig
  * ConnectionString - Postgres connection string
* LnxConfig
  * Host - Lnx host
  * Port - Lnx port
* Hosts - array indicating hosts Reisen should try to get certificates for from Let's Encrypt for no proxy deployments
* Production - tells Reisen to try to obtain certificates from Let's Encrypt
* Boards
  * Name - board name (without slashes)
  * Description - board description (as it appears in the index page)
* CspConfig - CSP configuration
* TemplateConfig
  * ImagesUrl - url from which media is served (media is assumed to be avaiable in ${ImagesUrl}/{mediaSha256InBase64}
  * ThumbnailsUrl - url from which thumbnails are served (thumbnails are assumed to be avaiable in ${ImagesUrl}/{mediaSha256InBase64}
  * FaviconUrl - url with favicon
  * JsUrl - url with reisen.js
  * TailwindCssUrl - url with tailwind.css
  * TegakiJsUrl - url with tegaki.js or tegaki.min.js (available in [the actual correct github](https://github.com/desuwa/tegaki) or copied and pasted into the static folder in this thingy)
  * TegakiCssUrl - url with tegaki.css available above
  * CssUrl - url with reisen.css
