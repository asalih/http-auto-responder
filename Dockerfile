FROM ubuntu:18.04
RUN wget https://github.com/asalih/http-auto-responder/releases/download/v1.0.0/http-auto-responder-linux.tar.gz
RUN tar -xf http-auto-responder-linux.tar.gz && cd http-auto-responder-linux
CMD ./http-auto-responder

EXPOSE 80