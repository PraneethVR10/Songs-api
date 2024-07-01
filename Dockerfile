FROM golang:1.22

WORKDIR /Songs-api

COPY . .


RUN cd /Songs-api/pkg && go build -o /Songs-api/Songs-api

CMD ["/Songs-api/Songs-api"]