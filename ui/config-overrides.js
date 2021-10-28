const { override, addLessLoader, fixBabelImports, addWebpackAlias } = require('customize-cra');
const path = require('path');

module.exports = override(
    addWebpackAlias({
        ['@']: path.resolve(__dirname, 'src'),
        ['@assets']: path.resolve(__dirname, 'src/assets'),
        ['@components']: path.resolve(__dirname, 'src/components'),
        ['@constants']: path.resolve(__dirname, 'src/constants'),
        ['@contexts']: path.resolve(__dirname, 'src/contexts'),
        ['@hooks']: path.resolve(__dirname, 'src/hooks'),
        ['@i18n']: path.resolve(__dirname, 'src/i18n'),
        ['@pages']: path.resolve(__dirname, 'src/pages'),
        ['@routes']: path.resolve(__dirname, 'src/routes'),
        ['@services']: path.resolve(__dirname, 'src/services'),
        ['@store']: path.resolve(__dirname, 'src/store'),
        ['@utils']: path.resolve(__dirname, 'src/utils')
    }),
    fixBabelImports('antd', {
        libraryDirectory: 'es',
        style: true,
    }),
    addLessLoader({
        javascriptEnabled: true,
        modifyVars: {
            '@primary-color': '#344468',
        },
    })
)
