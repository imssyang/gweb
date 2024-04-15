FROM golang:1.21

ARG NVM_VERSION=0.39.7
ARG NODE_VERSION=18.18.2

WORKDIR /opt/gweb

RUN apt-get update && apt-get upgrade -y
RUN apt-get update && apt-get install -y curl gnupg2 rsync python3 libpython3-dev
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource.gpg.key | gpg --dearmor -o /usr/share/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/usr/share/keyrings/nodesource.gpg] https://deb.nodesource.com/node_18.x bullseye main" > /etc/apt/sources.list.d/nodesource.list
RUN apt-get update && apt-get install -y nodejs


RUN echo $(pkg-config --cflags python3)
RUN echo $(pkg-config --libs python3)
RUN echo $(pkg-config --cflags --libs python3)

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN python3 -m compileall -b internal/api/formatify
RUN rsync -av --include="*/" --include="*.pyc" --exclude="*" \
    internal/api/formatify lib
#RUN go build -v -o bin/gweb cmd/gweb

CMD ["/opt/gweb/bin/gweb"]
