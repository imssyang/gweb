FROM ghcr.io/imssyang/formatui:v1.0.2 AS formatui
FROM golang:1.20-buster

WORKDIR /opt/gweb

RUN apt-get update && apt-get install -y curl tree rsync python3 libpython3-dev

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY public ./public
COPY templates ./templates
COPY tests ./tests
COPY third_party ./third_party
COPY --from=formatui /opt/formatui/dist/plugins/ public/plugins/
COPY --from=formatui /opt/formatui/dist/index.min.css public/css/formatify.min.css
COPY --from=formatui /opt/formatui/dist/index.min.js public/js/formatify.min.js
COPY --from=formatui /opt/formatui/src/img/formatui.svg public/img/formatify.svg

ENV CGO_CFLAGS -I/usr/include/python3.7m
ENV CGO_CXXFLAGS -I/usr/include/python3.7m -I/opt/gweb/third_party
ENV CGO_LDFLAGS -L/usr/lib/x86_64-linux-gnu -lpython3.7m

RUN python3 -m compileall -b ./internal/api/formatify
RUN rsync -av --include="*/" --include="*.pyc" --exclude="*" \
    internal/api/formatify ./lib
RUN go mod download && go mod verify
RUN go build -v -o ./bin/gweb ./cmd/gweb.go
RUN rm -rf go.mod go.sum \
    cmd \
    configs \
    internal \
    public \
    templates \
    tests \
    third_party

ENV PATH ${PATH}:/opt/gweb/bin
ENV PYTHONPATH /opt/gweb/lib
CMD ["gweb"]
