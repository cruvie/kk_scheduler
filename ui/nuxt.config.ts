// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    ssr: false,
    app: {
        head: {
            title: 'kk_scheduler'
        }
    },
    compatibilityDate: '2025-09-05',
    devtools: {enabled: true},
    modules: [
        '@nuxt/ui'
    ],
    css: ['~/assets/css/main.css'],
    build: {
    },
    ui: {
        fonts: false,
    },
    vite: {
        optimizeDeps: {
            include: [
                '@vue/devtools-core',
                '@vue/devtools-kit',
                '@bufbuild/protobuf',
                '@bufbuild/protobuf/codegenv2',
                '@connectrpc/connect',
                '@connectrpc/connect-web',
            ]
        }
    }
})
