FROM node:18-alpine

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

EXPOSE 8001

CMD ["npm", "run", "dev", "--", "--host"]
