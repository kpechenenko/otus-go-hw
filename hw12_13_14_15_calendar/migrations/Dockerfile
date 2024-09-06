FROM gomicro/goose

ADD *.sql /migrations/
COPY --chmod=755 entrypoint.sh /migrations/entrypoint.sh

ENTRYPOINT ["/migrations/entrypoint.sh"]