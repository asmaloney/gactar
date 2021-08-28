module.exports = {
  devServer: {
    // This lets us run gactar to serve the endpoints, but run the UI through npm for testing
    proxy: 'http://localhost:8181',
  },

  // don't include maps in the production output
  productionSourceMap: false,
}
