FROM node:17

WORKDIR /app

COPY package*.json ./

RUN npm install --force
# RUN npm ci --only=production

# Bundle app source
COPY . .

EXPOSE 3000
CMD [ "npm", "start" ]
