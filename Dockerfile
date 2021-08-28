FROM busybox
WORKDIR /app
COPY templates/ templates/
COPY micrach ./
ENTRYPOINT ["/app/micrach"]