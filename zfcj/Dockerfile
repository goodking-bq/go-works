FROM alpine
# 设置时区为上海
ENV TZ "Asia/Shanghai"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache --update tzdata \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone
COPY bin/zfcj /opt/zfcj
#RUN echo "0 * * * * /opt/zfcj crawl >> /var/log/zfcj.log" >> /var/spool/cron/crontabs/root
CMD /opt/zfcj exporter