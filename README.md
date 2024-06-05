# qrcode-generation-service


### To run the container you should write
docker build -t qrcode-generation-service .
docker run -p 7070:7070 --env-file .env -ti qrcode-generation-service