FROM goofygoof/phi-pnpm

WORKDIR /app

COPY . .

RUN pnpm install --force; exit 0

EXPOSE 3000
CMD ["pnpm", "start" ]

