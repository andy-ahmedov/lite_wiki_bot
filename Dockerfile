FROM alpine
ENV LANGUAGE="en"
COPY telegrambot.go .
#RUN apk add --no-cache ca-certificates && chmod +x code
EXPOSE 80/tcp
#CMD ["./bash"]