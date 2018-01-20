FROM scratch

MAINTAINER Luigi Riefolo <luigi.riefolo@gmail.com>

ARG CONFIG_FILE
ARG SERVICE
ARG VERSION

LABEL service=${SERVICE} \
    version=${VERSION}

#HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD curl -f http://localhost:9090/health || exit 1

#ADD ca-certificates.crt /ssl/certs/

#USER 751

ADD $CONFIG_FILE /

ADD src/${SERVICE}/cmd/${SERVICE} /

EXPOSE 8080

CMD ["/gateway", "service", "-config", "global_conf.toml"]
