FROM python:3.10-slim

ENV DEBIAN_FRONTEND="noninteractive"

RUN apt-get update && apt-get install -y build-essential nodejs npm

WORKDIR /app

COPY . .

RUN pip install -r requirements.txt

RUN make frontend

RUN make proto

EXPOSE 8083

RUN mkdir .data

WORKDIR /app/src/

CMD [ "python3", "main.py" ]
