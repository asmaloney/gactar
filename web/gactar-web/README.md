# gactar-web

This explains how to generate the website for development of the web interface, and how to build it for inclusion in the `gactar` binary.

It uses [vue.js](https://vuejs.org/), [vite](https://vitejs.dev/), [bulma](https://bulma.io/), and [buefy](https://buefy.org/) since they are the web tools I'm most familiar with.

It is set up to require recent versions of [node](https://nodejs.org/) & [npm](https://www.npmjs.com/). This is only because this is what I use. It may work with older versions.

## Setup

The first thing we need is to make sure the npm packages are installed:

```
npm install
```

## Development

This is set up so that the backend is served by running gactar and the frontend is served by [vite](https://vitejs.dev/) for live development.

- run `gactar -w` to serve the api endpoints
- run the vite server to see the site
  ```
  npm run dev
  ```

This will tell you what port it's running on - something like:

```
vite v2.7.3 dev server running at:

> Local: http://localhost:3000/
> Network: use `--host` to expose

ready in 416ms.
```

Navigate to this localhost site and you can work on the vue frontend with live updating.

## Release

To build & include the website in the gactar binary:

- compile the production site:
  ```
  npm run build
  ```
- copy the files from `dist/` to `gactar/web/build/`
- recompile gactar using `make`
