# Use Node.js 22 on Alpine
FROM node:22-alpine

WORKDIR /app

# Enable pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files first (layer caching)
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install

# Copy the source code
COPY . .

# Expose the dev server port (default 5173)
EXPOSE 5173

# Command to run Vite in dev mode
CMD ["pnpm", "dev", "--host", "0.0.0.0"]
