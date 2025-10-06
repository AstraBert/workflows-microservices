FROM ghcr.io/astral-sh/uv:python3.13-alpine

WORKDIR /payments/

COPY ./ /payments/

RUN uv sync

ENTRYPOINT [ "uvx", "--from", ".", "process-payments" ]