FROM ysicing/god as build

COPY . /go/src

WORKDIR /go/src

RUN GOOS=linux GOARCH=amd64 go build -o godemo

FROM ysicing/debian

COPY --from=build /go/src/godemo /usr/bin/godemo

RUN chmod +x /usr/bin/godemo

CMD [ "/usr/bin/godemo" ]
