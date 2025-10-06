FROM ghcr.io/astral-sh/uv:python3.13-alpine

WORKDIR /stock/

COPY ./ /stock/

RUN uv sync

ENTRYPOINT [ "uvx", "--from", ".", "process-stock" ]