FROM node:18-alpine

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

# Указываем порт, который будет слушать dev-сервер Vite
EXPOSE 8001

# Запускаем dev-сервер Vite
CMD ["npm", "run", "dev", "--", "--host"]
