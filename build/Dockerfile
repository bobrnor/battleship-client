FROM alpine:3.6
MAINTAINER danil@nulana.com

COPY battleship /internal/bin/battleship-client
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/internal/bin/battleship-client"]