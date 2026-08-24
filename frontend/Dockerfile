FROM node:22-slim AS builder
WORKDIR /app


COPY . .


RUN rm -rf node_modules package-lock.json && \
    npm install && \
    npm run build


FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]