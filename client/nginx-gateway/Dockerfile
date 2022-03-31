FROM nginx:1.15.9-alpine
ARG FILE
COPY $FILE /etc/nginx/nginx.conf
WORKDIR /usr/share/nginx/html
EXPOSE 80
