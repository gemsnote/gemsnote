FROM python:3.12-slim
WORKDIR /worker
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && useradd --uid 10001 worker
COPY tasks.py .
USER worker
CMD ["celery", "-A", "tasks", "worker", "--loglevel=INFO", "--concurrency=2", "--queues=gemsnote.mail"]
