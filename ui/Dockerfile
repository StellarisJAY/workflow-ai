FROM nginx:1.25.3-alpine
COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./dist /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
EXPOSE 8088

CMD ["nginx", "-g", "daemon off;"]