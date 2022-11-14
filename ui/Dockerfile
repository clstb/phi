FROM node:19-slim

WORKDIR /app

COPY . .

RUN npm i --force && npm run build && npm install http-server -g

EXPOSE 3000

CMD ["http-server", "build/", "--port=3000" ]

