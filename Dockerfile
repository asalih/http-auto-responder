FROM ubuntu:18.04
RUN  apt-get update \
  && apt-get install -y wget \
  && rm -rf /var/lib/apt/lists/*
RUN wget https://github.com/asalih/http-auto-responder/releases/download/v1.0.0/http-auto-responder-linux.tar.gz
RUN tar -xf http-auto-responder-linux.tar.gz
WORKDIR "/http-auto-responder-linux"
CMD ["./http-auto-responder"]

EXPOSE 80