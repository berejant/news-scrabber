# Whisper server using the official OpenAI whisper package and a simple FastAPI wrapper
FROM python:3.14-slim

ENV PIP_NO_CACHE_DIR=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    TORCH_HOME=/models \
    WHISPER_MODEL=base \
    WHISPER_DEVICE=cpu

# ffmpeg is required by whisper for audio decoding
RUN apt-get update \
    && apt-get install -y --no-install-recommends ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# Install whisper plus API/runtime deps (includes uvicorn entrypoint)
RUN pip install --break-system-packages openai-whisper fastapi uvicorn[standard] python-multipart

WORKDIR /app
COPY server.py /app/server.py

VOLUME ["/models"]
EXPOSE 10300

CMD ["uvicorn", "server:app", "--host", "0.0.0.0", "--port", "10300"]
