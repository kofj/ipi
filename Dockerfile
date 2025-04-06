FROM alpine:3.18

COPY _output/linux/amd64 /_build/
CMD [ipi]