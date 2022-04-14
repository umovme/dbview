FROM postgres:14-bullseye

LABEL maintainer="PostGIS Project - https://postgis.net"

ENV DBVIEW_VERSION "0.8.1"
ENV POSTGIS_MAJOR 3
ENV POSTGIS_VERSION 3.2.1+dfsg-1.pgdg110+1

RUN apt-get update \
      && apt-cache showpkg postgresql-$PG_MAJOR-postgis-$POSTGIS_MAJOR \
      && apt-get install -y --no-install-recommends \
           postgresql-$PG_MAJOR-postgis-$POSTGIS_MAJOR=$POSTGIS_VERSION \
           postgresql-$PG_MAJOR-postgis-$POSTGIS_MAJOR-scripts \
      && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y \
      build-essential \
	devscripts \
      curl \
      vim \
      && apt-get clean \
      && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /docker-entrypoint-initdb.d
COPY ./initdb-postgis.sh /docker-entrypoint-initdb.d/10_postgis.sh
COPY ./update-postgis.sh /usr/local/bin

COPY ./config.toml ./  
#COPY ./dbview.pgbkp ./

RUN curl -L https://github.com/umovme/dbview/releases/download/v${DBVIEW_VERSION}/dbview_${DBVIEW_VERSION}_linux_amd64.deb -o ./tmp/dbview.deb

RUN dpkg -i ./tmp/dbview.deb   && \
    rm -rf ./tmp/dbview.deb
ENV TZ=America/Sao_Paulo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
