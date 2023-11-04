# Go htmx server

## Dependencies

- docker (docker-desktop if you are using windows)
- docker-compose (comes with docker-desktop, but can install [here](https://docs.docker.com/compose/install/standalone/) if you are not on windows)
- [Node LTS v18](https://nodejs.org/en/download)
- [Go v1.21](https://go.dev/doc/install)
- [Air](https://github.com/cosmtrek/air)
- [templ](https://github.com/a-h/templ)

## Features

https://github.com/Jason-CKY/go-htmx-example/assets/27609953/9519a6ea-a5e4-407d-8a29-dcdd76bc2857

- [Gin](https://gin-gonic.com/) web server that serves html on htmx endpoints
- [HTMX](https://htmx.org/) for interactivity, minimal js needed
- Lazy loading with HTMX
- [DaisyUI](daisyui.com/) with [theme-changing library](https://github.com/saadeghi/theme-change) for CSS styling and themes
- [SortableJS](https://github.com/SortableJS/Sortable) for drag and drop of tasks (sorting and updates)
- [Directus](https://directus.io/) for headless CMS and API routes for CRUD operations

## Quickstart (development mode)

You can either start up using `docker-compose`:

```sh
# Run install-deps once to install all dev dependencies
make install-deps
```

```sh
make build-dev
# make sure directus is up on http://localhost:8055 before running migrations for directus
make initialize-db
```

Or you can run locally with:

```sh
# start directus
docker-compose -f docker-compose.dev.yml start directus
# make sure directus is up on http://localhost:8055 before running migrations for directus
make initialize-db
# install air
go install github.com/cosmtrek/air@latest
# install templ
go install github.com/a-h/templ/cmd/templ@latest
# start golang server with code reload using air
air
```

## Format on save

Refer to this [link](https://www.digitalocean.com/community/tutorials/how-to-format-code-with-prettier-in-visual-studio-code) on how to install and set prettier to format on save.

## VS-code extensions for good developer experience

- [Prettier - Code formatter](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode)
- [Tailwind CSS IntelliSense](https://marketplace.visualstudio.com/items?itemName=bradlc.vscode-tailwindcss)

### Syntax highlighting of golang template files on vscode

- Download [templ-vscode](https://marketplace.visualstudio.com/items?itemName=a-h.templ) vscode extension for go-templ syntax highlighting
- Add the following into your vscode `settings.json` to allow for tailwind syntax highlighting in your go `templ` files:

```json
"tailwindCSS.includeLanguages": {
"templ": "html"
}
```
