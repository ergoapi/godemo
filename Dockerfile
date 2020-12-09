FROM registry.cn-beijing.aliyuncs.com/k7scn/goa as build

COPY . /go/src

WORKDIR /go/src

RUN go build

FROM registry.cn-beijing.aliyuncs.com/k7scn/alpine

COPY --from=build /go/src/godemo /usr/bin/godemo

RUN chmod +x /usr/bin/godemo

CMD [ "/usr/bin/godemo" ]