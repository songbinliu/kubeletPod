FROM alpine:3.3

COPY ./_output/kclient.linux /bin/kclient

ENTRYPOINT ["/bin/kclient"]
