# Base image for building the Next.js app
FROM node:18-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json first (for efficient caching)
COPY package.json ./
COPY package-lock.json ./

# Install dependencies
RUN npm install --frozen-lockfile

# Copy the rest of the application files
COPY . .

# Build the Next.js application
RUN npm run build:prod

# Remove development dependencies (optional)
RUN npm prune --production

# Production image
FROM node:18-alpine AS runner

# Set working directory
WORKDIR /app

# Copy only necessary files from builder
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/package.json ./package.json

# Expose the Next.js default port
EXPOSE 3000

# Set environment variables (optional)
ENV NODE_ENV=production

# Start the Next.js application
CMD ["npm", "run", "start:prod"]
