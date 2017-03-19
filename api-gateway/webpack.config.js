var path = require('path');

module.exports = {
    entry: [
        path.resolve(__dirname, 'client/src/index.js')
    ],
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'client/public/dist')
    },
    module: {
        loaders: [{
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                loader: 'babel-loader',
                query: {
                    presets: ['es2015', 'react']
                }
            },
            // CSS handling
            // * style-loader: Embeds referenced CSS code using a <style>-element in your index.html file
            // * css-loader: Parses the actual CSS files referenced from your code. Modifies url()-statements in your
            //   CSS files to match images handled by url loader (see below)
            { test: /\.css$/, loader: 'style-loader!css-loader' },

            // Image Handling
            // * url-loader: Returns all referenced png/jpg files up to the specified limit as inline Data Url
            //   or - if above that limit - copies the file to your output directory and returns the url to that copied file
            //   Both values can be used for example for the 'src' attribute on an <img> element
            { test: /\.(png|jpg)$/, loader: 'url-loader?limit=25000' },

            // JSon file handling
            // * Enables you to 'require'/'import' json files from your JS files
            { test: /\.json$/, loader: 'json-loader' }
        ]
    },
    devtool: 'source-map'
};