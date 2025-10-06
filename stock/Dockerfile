FROM ghcr.io/astral-sh/uv:python3.13-alpine

WORKDIR /stock/

COPY ./ /stock/

RUN uv sync

ENV PYTHONUNBUFFERED=1
ENTRYPOINT ["/bin/sh", "-c", "source .venv/bin/activate && process-stock"]