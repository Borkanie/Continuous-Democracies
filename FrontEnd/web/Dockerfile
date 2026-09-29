# dist/ is built on the host (`pnpm install && pnpm build`), not inside this
# image: this network TLS-intercepts outbound HTTPS, and corepack/pnpm
# inside a container can't fetch from registry.npmjs.org here the way the
# host can. This Dockerfile just serves the pre-built static output.
FROM nginx:alpine

COPY dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
