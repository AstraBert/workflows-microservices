FROM ghcr.io/astral-sh/uv:python3.13-alpine

WORKDIR /orders/

COPY ./ /orders/

RUN uv sync

ENV PYTHONUNBUFFERED=1
ENTRYPOINT ["/bin/sh", "-c", "source .venv/bin/activate && process-orders"]