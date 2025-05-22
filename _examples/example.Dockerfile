# Stolen from my repo
# https://github.com/coffeemakingtoaster/dotfile-loader/tree/main

FROM golang:1.22.1-alpine3.18 AS build 

WORKDIR /build

COPY . .

RUN go build -o dotfile-server 

## -- RUNTIME STAGE --
FROM alpine:3.18.4 AS runtime

WORKDIR /app

ARG USER=docker
ARG UID=5432
ARG GID=5433

# Create user for execution

#User group has same name as user
RUN addgroup -g $GID $USER 


RUN adduser \
	--disabled-password \
	--gecos "" \
	--ingroup "$USER" \
	--no-create-home \
	--uid "$UID" \
	"$USER"

# Copy build with permissions
COPY --from=build --chown=$USER:$USER /build/dotfile-server /app/dotfile-server

# Ensure that backend can be run
RUN chmod +x /app/dotfile-server

USER $USER 

CMD ["/app/dotfile-server"]
