FROM node:18.17.1-alpine AS base

WORKDIR /app

COPY package*.json ./

FROM base AS dev

RUN npm install

COPY . .

CMD ["yarn", "start:dev"]


FROM base AS production

RUN npm install

COPY . .

RUN npm run build

CMD ["yarn", "start:prod"]



