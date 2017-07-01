FROM alpine:3.6
MAINTAINER danil@nulana.com

COPY battleship /internal/bin/battleship_client
COPY entrypoint.sh /entrypoint.sh

RUN ls -l
RUN ls -l /internal/bin/

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/internal/bin/battleship_client"]