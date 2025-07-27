FROM scratch
WORKDIR /
COPY fluffly .
ENTRYPOINT ["/fluffly"]