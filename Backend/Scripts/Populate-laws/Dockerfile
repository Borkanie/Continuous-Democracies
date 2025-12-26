FROM python:3.11-slim

# Build-time args (can be set from CI using --build-arg)
ARG DBName
ARG DBUser
ARG DBPassword
ARG DBServerAddress
ARG DBPort=5432

# Runtime environment variables (populated from args by default)
ENV DBName=${DBName}
ENV DBUser=${DBUser}
ENV DBPassword=${DBPassword}
ENV DBServerAddress=${DBServerAddress}
ENV DBPort=${DBPort}

# Install system dependencies for Poppler and Tesseract + fonts
RUN apt-get update \
     && apt-get install -y --no-install-recommends \
         curl \
         poppler-utils \
         tesseract-ocr \
         tesseract-ocr-ron \
         libgl1 \
         libsm6 \
         libxrender1 \
         libxext6 \
         libpq-dev \
         build-essential \
         python3-dev \
     && rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /app

# Copy requirements and install
COPY requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt

# Copy app files
COPY . /app

# Ensure temp and logs exist
RUN mkdir -p /app/temp

# Default command
CMD ["python", "-u", "orchestrator.py"]
