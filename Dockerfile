FROM scratch
COPY phi phi
COPY /sql/schema /sql/schema
ENTRYPOINT ["./phi"]
