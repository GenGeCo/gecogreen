# ============================================
# GecoGreen Frontend - Production Dockerfile
# ============================================

# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build arguments for API URL
ARG PUBLIC_API_URL=http://localhost:8080
ENV PUBLIC_API_URL=$PUBLIC_API_URL

# Build the application
RUN npm run build

# Production stage
FROM node:20-alpine AS runner

WORKDIR /app

# Create non-root user
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 sveltekit

# Copy built application
COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules

# Set ownership
RUN chown -R sveltekit:nodejs /app

USER sveltekit

# Expose port
EXPOSE 3000

ENV NODE_ENV=production
ENV PORT=3000
ENV ORIGIN=http://localhost:3000

# Start the server
CMD ["node", "build"]
