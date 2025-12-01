# Stage 1: Development/Build Stage
FROM node:20-alpine AS dev

# Set the working directory
WORKDIR /app

# Copy package files from the service directory to install dependencies
COPY services/integration-gateway/package*.json ./
RUN npm install

# Copy all project files (source code, client files)
COPY . .

# Compile TypeScript (Requires 'build' script defined in package.json)
RUN npm run build 

# Stage 2: Production Stage (Minimal Image)
FROM node:20-alpine AS prod
WORKDIR /app

# Copy only necessary runtime files
COPY --from=dev /app/node_modules ./node_modules
COPY --from=dev /app/dist ./dist 

# Copy the static HTML wrapper file that loads the game client
COPY services/integration-gateway/index.html ./ 

# Expose the HTTP/WebSocket port
EXPOSE 8080

# Command to run the compiled JavaScript server
CMD ["node", "dist/server.js"]