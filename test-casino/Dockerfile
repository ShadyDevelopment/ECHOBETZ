# Use a minimal Node.js base image
FROM node:20-alpine

# Set the working directory inside the container
WORKDIR /app

# Corrected COPY command:
# Source: test-casino/package*.json (relative to the ECHOBETZ root context)
# Destination: ./ (which is /app inside the container)
COPY test-casino/package*.json ./

# Install dependencies (will now find package.json in /app)
RUN npm install

# Copy the server logic file
COPY test-casino/server.js ./

# Expose the API port
EXPOSE 3000

# Command to run the test casino server
CMD ["node", "server.js"]