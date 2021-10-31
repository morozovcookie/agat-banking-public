# migrate/migrate:v4.15.1
FROM migrate/migrate@sha256:254d67343a2389383b530d9049610f21f4d302557288c649ebbb761ea96b10af

COPY ./migrations/percona /migrations

COPY ./scripts/docker/bankingd-mysql-schema/docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/bin/sh", "/docker-entrypoint.sh"]
