FROM nginx:1.27.0

RUN echo "Asia/Shanghai" > /etc/timezone

COPY nginx.conf /etc/nginx/nginx.conf
COPY dist /usr/share/nginx/dist/

# 启动 Nginx 服务
CMD ["nginx", "-g", "daemon off;"]
