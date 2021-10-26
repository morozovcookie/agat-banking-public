module.exports = (config, env) => {
    const path = require('path');

    let alias = config.resolve.alias;

    alias = {
        ...alias,

        '@': path.resolve(__dirname, 'src'),
        '@assets': path.resolve(__dirname, 'src/assets'),
        '@components': path.resolve(__dirname, 'src/components'),
        '@constants': path.resolve(__dirname, 'src/constants'),
        '@contexts': path.resolve(__dirname, 'src/contexts'),
        '@hooks': path.resolve(__dirname, 'src/hooks'),
        '@i18n': path.resolve(__dirname, 'src/i18n'),
        '@pages': path.resolve(__dirname, 'src/pages'),
        '@routes': path.resolve(__dirname, 'src/routes'),
        '@services': path.resolve(__dirname, 'src/services'),
        '@stores': path.resolve(__dirname, 'src/stores'),
        '@utils': path.resolve(__dirname, 'src/utils')
    }

    config.resolve.alias = alias;

    return config;
}
