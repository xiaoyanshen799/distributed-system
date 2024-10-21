FROM ubuntu:latest
LABEL authors="zhanchengbo"

ENTRYPOINT ["top", "-b"]

FROM python:3.9-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python", "server.py"]