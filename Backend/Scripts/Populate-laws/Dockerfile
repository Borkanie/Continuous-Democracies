FROM python:3.11-slim

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
CMD ["python", "-u", "pupulate_laws.py"]
