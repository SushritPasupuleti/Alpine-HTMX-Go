# Alpine HTMX Go Web App

Serving a web experience similar to React, through Alpine.js, HTMX served by a Go-chi server.

Server bootstrapped [from](https://github.com/SushritPasupuleti/Go-Chi-Boilerplate).

## Stack

- [alpine.js](https://alpinejs.dev/)

- [htmx](https://htmx.org/)

- [Tailwind CSS](https://tailwindcss.com/)

- [Hyper UI Components](https://www.hyperui.dev/)

- [Go](https://golang.org/)

- [Go Chi](https://github.com/go-chi/chi)

## Features

- [x] Switch between `json` and `html` response types based on `Accept` header. Keeps the API interoperable with other clients.

- [x] Caching for `html` responses in addition to `json` responses.

Base features are carried over from the [boilerplate](https://github.com/SushritPasupuleti/Go-Chi-Boilerplate).

## Setup

Run `make` to see all available commands.

### Install dependencies

```bash
cd server
make packages_install
```

### Run

Before running, make sure generate the `tailwindcss` styles.

```bash
cd server/templates
yarn install
yarn build
```

> Alternatively, you can run `yarn watch` to watch for changes during development.

```bash
cd server
make run
```
