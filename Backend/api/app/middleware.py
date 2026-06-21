import time
import logging

from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request

logger = logging.getLogger(__name__)


class RequestLoggingMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        start = time.perf_counter()
        logger.info("→ %s %s from %s", request.method, request.url.path, request.client)
        response = await call_next(request)
        elapsed = (time.perf_counter() - start) * 1000
        logger.info("← %s %s %.1fms", response.status_code, request.url.path, elapsed)
        return response
