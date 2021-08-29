FROM busybox
WORKDIR /app
COPY templates/ templates/
COPY micrach ./
RUN chmod +x /app/micrach
RUN ls
ENTRYPOINT ["/app/micrach"]