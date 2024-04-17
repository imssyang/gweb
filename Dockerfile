FROM ghcr.io/imssyang/formatui:latest AS formatui
FROM golang:1.20-buster

WORKDIR /opt/gweb

COPY --from=formatui /opt/formatui/dist/plugins/ public/plugins/
COPY --from=formatui /opt/formatui/dist/index.min.css public/css/formatify.min.css
COPY --from=formatui /opt/formatui/dist/index.min.js public/js/formatify.min.js
COPY --from=formatui /opt/formatui/src/img/formatui.svg public/img/formatify.svg

RUN apt-get update && apt-get install -y curl tree rsync python3 libpython3-dev

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY public ./public
COPY templates ./templates
COPY tests ./tests
COPY third_party ./third_party

RUN python3 -m compileall -b internal/api/formatify
ENV CGO_CFLAGS -I/usr/include/python3.7m
ENV CGO_CXXFLAGS -I/usr/include/python3.7m -I/opt/gweb/third_party
ENV CGO_LDFLAGS -L/usr/lib/x86_64-linux-gnu -lpython3.7m
RUN rsync -av --include="*/" --include="*.pyc" --exclude="*" \
    internal/api/formatify deploy
RUN go build -v -o deploy/gweb cmd/gweb.go

CMD ["/opt/gweb/bin/gweb"]
