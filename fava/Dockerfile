FROM node:19-slim as builder

WORKDIR /app

COPY . .

RUN cd frontend && npm i && npm run build


FROM python:3.10-slim

COPY --from=builder /app .

RUN cd .. & pip install -r requirements.txt

RUN mkdir .data

EXPOSE 8083

WORKDIR /app/src/

CMD [ "python3", "main.py" ]
