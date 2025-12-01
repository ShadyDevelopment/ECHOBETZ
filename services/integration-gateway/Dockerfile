# Using multi-stage build for development speed
FROM node:20-alpine AS dev
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
# Compile TypeScript
RUN npm run build 

FROM node:20-alpine AS prod
WORKDIR /app
COPY --from=dev /app/node_modules ./node_modules
COPY --from=dev /app/dist ./dist # Compiled code
COPY index.html ./ # Static HTML file

EXPOSE 8080
CMD ["node", "dist/server.js"]