const paths = require('react-scripts/config/paths')
const path = require('path')

paths.appPublic = path.resolve(__dirname, 'webapp');
paths.appHtml = path.resolve(__dirname, 'webapp/index.html');
paths.appSrc = path.resolve(__dirname, 'webapp')
paths.appIndexJs = path.resolve(__dirname, 'webapp/index.js')

module.exports = function override(config, env) {

    return config;
}

