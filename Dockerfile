FROM alpine:latest
LABEL maintainer="Remi Philippe <remi@cisco.com>"

# Copy tetviz
RUN mkdir -p /demo-policy-pipeline/web
WORKDIR /demo-policy-pipeline

COPY web /demo-policy-pipeline/web

EXPOSE 1323

ENTRYPOINT ["./demo-policy-pipeline"]
