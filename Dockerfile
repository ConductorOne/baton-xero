FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-xero"]
COPY baton-xero /