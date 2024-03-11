docker rm netcatContainer
docker rmi netcat_image
docker build -t netcat_image .
docker run -it --name netcatContainer -p 8989:8989 -p 8990:8990 -p 8991:8991 -p 8992:8992 -p 8993:8993 -p 8994:8994 -p 8995:8995 -p 8996:8996 -p 8997:8997 -p 8998:8998 -p 8999:8999 -p 9000:9000 netcat_image

