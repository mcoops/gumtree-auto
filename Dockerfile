FROM golang:1.20 as build

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 go build -o gumtree-auto 

FROM chromedp/headless-shell:latest
ENV USERNAME=username
ENV PASSWORD=password

COPY --from=build /app/gumtree-auto /app/gumtree-auto

RUN chmod +x /app/gumtree-auto
ENTRYPOINT ["/app/gumtree-auto"]