FROM scratch
WORKDIR /app
COPY templates/ templates/
COPY static/ static/
COPY micrach ./
ENTRYPOINT ["/app/micrach"]