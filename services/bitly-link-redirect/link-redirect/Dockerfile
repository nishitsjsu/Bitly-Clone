FROM golang:latest 
EXPOSE 3002
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
RUN cd /app ; go install goapi
CMD ["/app/bin/goapi"]
