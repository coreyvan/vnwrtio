FROM golang:1.13-alpine as builder

WORKDIR /vnwrtio

# Git needed for the go get
# RUN apk update && apk add --no-cache git ca-certificates tzdata openssh-client && update-ca-certificates

COPY go.mod ./
RUN go mod download

# Copy in source
COPY ./go.* ./
COPY ./*.go ./
COPY ./src/* ./src/
# RUN ls ./src

# CGO_ENABLED=0 to avoid using C routines for the dns resolution.
# These routines are not available in the scratch docker image
# see https://golang.org/pkg/net/#hdr-Name_Resolution
RUN CGO_ENABLED=0 go build -o server

FROM scratch

# Copy across the binary
COPY --from=builder /vnwrtio/server .
COPY --from=builder /vnwrtio/src/*.css ./src/css/
COPY --from=builder /vnwrtio/src/*.js ./src/js/
COPY --from=builder /vnwrtio/src/*.html ./src/templates/
COPY --from=builder /vnwrtio/src/*.jpg ./src/assets/
COPY --from=builder /vnwrtio/src/*.ico ./src/assets/

EXPOSE 80
ENTRYPOINT ["./server"]