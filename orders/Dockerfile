FROM ghcr.io/astral-sh/uv:python3.13-alpine

WORKDIR /orders/

COPY ./ /orders/

RUN uv sync

ENTRYPOINT [ "uvx", "--from", ".", "process-orders" ]