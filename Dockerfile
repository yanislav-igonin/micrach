FROM busybox
WORKDIR /app
COPY templates/ templates/
COPY micrach ./
RUN chmod +x /app/micrach
ENTRYPOINT ["/app/micrach"]