FROM alpine:3.19.0

LABEL description="Stator the observable service"

COPY bin/* /usr/local/bin/

CMD stator
