# Собираем в гошке
FROM golang:1.22 AS build

ENV BIN_FILE=/opt/calendar/calendar-app
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

COPY . .

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} -mod vendor ${CODE_DIR}/cmd/calendar/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL SERVICE="calendar"

ENV BIN_FILE="/opt/calendar/calendar-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE=/etc/calendar/config.yaml
COPY configs/calendar/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} --config ${CONFIG_FILE}
