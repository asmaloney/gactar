# gactar-web

This is the setup to generate the website which gets included in `gactar` to edit and run models.

It uses [vue.js](https://vuejs.org/), [bulma](https://bulma.io/), and [buefy](https://buefy.org/) since they are the web tools I'm most familiar with.

It is set up to require recent versions of node & npm. This is only because this is what I use. It may work with older versions.

To include the website in gactar:

- install modules
  ```
  npm install
  ```
- (if doing development):
  - run `gactar -w` to serve the endpoints
  - run the npm server to see the site
    ```
    npm run serve
    ```
- compile the production site:
  ```
  npm run build
  ```
- finally, copy the files from `dist/` to `gactar/web/build/` and recompile gactar
