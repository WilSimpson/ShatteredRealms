##########
# LEGACY #
##########
FROM debian:latest

RUN apt-get update \
    && apt-get upgrade -y
    
RUN adduser server
COPY ./Package/1.0.0/Shipping/LinuxServer /home/server/server
RUN chown -R server:server /home/server
RUN chmod o+x /home/server/server

USER server

EXPOSE 7777/udp

ENTRYPOINT ["/home/server/server/SROServer.sh"]