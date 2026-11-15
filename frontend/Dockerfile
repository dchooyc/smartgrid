# Build stage: installs the dependencies and builds the production app.
FROM node:26-slim AS build
WORKDIR /app

# Dependencies first, so they are only installed again when they change.
COPY package.json package-lock.json ./
RUN npm ci

COPY . .
# The backend's address is fixed into the app when it is built.
ARG NEXT_PUBLIC_API_URL=http://localhost:8080
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
RUN npm run build

# Run stage: serves the built app.
FROM node:26-slim
WORKDIR /app
ENV NODE_ENV=production

COPY --from=build --chown=node:node /app/package.json /app/package-lock.json /app/next.config.ts ./
COPY --from=build --chown=node:node /app/node_modules ./node_modules
COPY --from=build --chown=node:node /app/.next ./.next
COPY --from=build --chown=node:node /app/public ./public

USER node
EXPOSE 3000
CMD ["npm", "run", "start"]
