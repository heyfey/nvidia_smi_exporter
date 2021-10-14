FROM nvidia/cuda:11.2.2-base-ubuntu18.04
MAINTAINER Tsung-Tso Hsieh <tsungtsohsieh@gmail.com>

COPY nvidia_smi_exporter.go /nvidia_smi_exporter.go
RUN apt-get update && apt-get install -y golang-go && go build /nvidia_smi_exporter.go && apt-get remove -y golang-go

EXPOSE 9101:9101

ENTRYPOINT ["/nvidia_smi_exporter"]
