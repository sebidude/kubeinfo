FROM scratch

COPY build/linux/kubeinfo /usr/bin/kubeinfo
ENTRYPOINT ["/usr/bin/kubeinfo"]
