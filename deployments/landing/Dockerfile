FROM node:22-alpine as base

RUN npm install -g pnpm@10.13.1

WORKDIR /app

COPY internal/landing/package.json internal/landing/pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY internal/landing/ .

FROM base AS build

RUN pnpm run build

FROM nginx:alpine AS production

COPY deployments/landing/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
