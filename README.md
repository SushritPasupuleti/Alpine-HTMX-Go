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

- [x] Fully CSP compliant.

- [x] Authentication and Authorization flows with JWT.

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

## Design Considerations and Advice

- Create helper functions to return specific `error` types. This will help you to handle errors in a more granular way. Example: `NotFound`, `Unauthorized`, `Forbidden`, `BadRequest`, `InternalServerError` could trigger rendering of a toast notification or a banner.

## Pros and Cons

This approach is not a replacement for React and co. But it can be used when you are looking for the following set of pros, while being aware of the cons.

### Pros

- Potentially higher RPS (Requests per second) if coupled with a strong backend, with caching/CDNs and other optimizations.

- Server Driven UIs. (SEO, Easier A/B testing, etc.)

- Co-opt (incrementally) with existing applications.

- Simpler DevOps. (No need to deploy a separate frontend app)

- Skip expensive JSON serialization and deserialization on both ends.

### Cons

- CSP (Content Security Policy) can be a nightmare to configure.

    - Refer to [this](https://alpinejs.dev/advanced/csp) for Alpine.js. (TL;DR: Strict CSP takes away most of the ease of Alpine.js)

    - Refer to [this](https://htmx.org/docs/#security) for HTMX. (TL;DR: HTMX with tight configuration can be secure, but the lack of Alpine.js's ease makes the overall experience less pleasant)

- More data transferred per Request (HTML is more verbose than JSON)

- Slower/Constrained Developer Experience (Hard to Debug, Poor IDE support, syntax highlighting, autocomplete, etc.)

- Code Sharing and Monorepo advantages from a stack like React+React Native are lost.
