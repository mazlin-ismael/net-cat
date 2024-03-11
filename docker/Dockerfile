# taking golang from docker library
FROM golang:1.22.0

# META DATA 
LABEL version="1.0" maintainer="imazlin"

# COPY NEEDED FILES
COPY dockerNC /

# work repertory
WORKDIR /newNET-CAT

# LAUNCH PORTS start 8989 to 9000
EXPOSE 8989-9000

CMD ["./NETCAT"]