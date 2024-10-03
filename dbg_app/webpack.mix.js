const mix = require('laravel-mix')
const webpack = require('webpack')
mix.js('app/app.js', 'debugger.js')
    .setPublicPath('server')
    .browserSync({
        online: false,
        host: 'localhost',
        port: 3666,
        ui: {port: 3667},
        server: {
            baseDir: "./server",
            index: "index.html",
        },
        single: true,
        open: false,
        watch: true
    })
    .webpackConfig({
        plugins: [new webpack.DefinePlugin({
            __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
        })]
    })
    .vue()
    .version()
    .then(() => {
        generateIndexHtml('./index.html.template', './server/mix-manifest.json', './server/index.html')
    })


const fs = require('fs');

function generateIndexHtml(templatePath, manifestPath, outputPath) {
    const templateContent = fs.readFileSync(templatePath, 'utf8')
    const manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf8'))

    let updatedContent = templateContent
    for (const asset in manifest) {
        updatedContent = updatedContent.replace(new RegExp(asset, 'g'), manifest[asset])
    }

    fs.writeFileSync(outputPath, updatedContent, 'utf8')

    console.log(`Generated ${outputPath}!`)
}