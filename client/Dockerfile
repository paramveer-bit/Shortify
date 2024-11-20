# Use official Node.js image to build the Next.js app
FROM node:18-alpine as builder

# Set the working directory inside the container
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install


# Copy the frontend files
COPY . .

# Install dependencies and build the Next.js app
RUN npm run build

# Install a production-ready Node.js server to serve the app
FROM node:18-alpine

WORKDIR /app

# Copy the built Next.js app from the builder image
COPY --from=builder /app ./

# Expose the port that Next.js will run on
EXPOSE 3000

# Start the Next.js app
CMD ["npm", "start"]