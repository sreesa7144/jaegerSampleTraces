FROM golang:1.13

# Create test directory
RUN mkdir -p /$GOPATH/src/github.com/sreesa7144/jaegerSampleTraces
WORKDIR /$GOPATH/src/github.com/sreesa7144/jaegerSampleTraces

COPY . /$GOPATH/src/github.com/sreesa7144/jaegerSampleTraces
RUN go build .
CMD ["./jaegerSampleTraces"]