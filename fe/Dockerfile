# Base image Bun/Node
FROM oven/bun

# Set working directory
WORKDIR /app

# Copy dependency files dulu supaya caching layer optimal
COPY package.json bun.lockb ./

# Install dependencies
RUN bun install

# Copy seluruh source code
COPY . .

# Expose port container (container port tetap 3000)
EXPOSE 3000

# Jalankan frontend
CMD ["bun", "run", "start"]
