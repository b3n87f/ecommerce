FROM golang:1.21

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . ./
EXPOSE 8080
RUN go build -v -o itm_project
RUN ["chmod", "+x", "/app/itm_project"]
CMD ["/app/itm_project"]
