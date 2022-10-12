FROM alpine:3.16.2

COPY zed_entrypoint.sh /zed_entrypoint.sh

RUN echo "Download zed client" \
    && cd /tmp \
    && wget https://github.com/authzed/zed/releases/download/v0.7.3/zed_0.7.3_linux_amd64_musl.tar.gz  \
    && tar xzf zed_0.7.3_linux_amd64_musl.tar.gz \
    && chmod u+x ./zed \
    && mv ./zed /usr/local/bin/zed

RUN chmod u+x /zed_entrypoint.sh

CMD ["/zed_entrypoint.sh"]