FROM node:23-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
COPY .env.staging .env
RUN npm run build

FROM node:23-alpine AS runner
WORKDIR /app
COPY --from=builder /app/dist /app/dist
RUN npm i -g serve
EXPOSE 3000
CMD ["serve", "-p", "3000", "-s", "dist"]

