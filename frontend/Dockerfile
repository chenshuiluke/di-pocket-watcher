# Use an official Node runtime as the base image
FROM node:20-alpine

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json (if available)
COPY package*.json ./

# Install dependencies
RUN npm install --include dev

# Copy the rest of the application code
COPY . .

# Expose the port Vite runs on by default
EXPOSE 5173

# Start the Vite dev server
COPY start.sh .
RUN chmod +x start.sh