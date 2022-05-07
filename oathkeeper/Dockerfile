FROM oryd/oathkeeper:v0.38.23

WORKDIR /home/ory

ADD config.yaml .oathkeeper.yaml
ADD rules.yaml rules.yaml

ENTRYPOINT [ "sh" ]
CMD ["-c", "oathkeeper credentials generate --alg RS256 > jwks.json && oathkeeper serve"]
