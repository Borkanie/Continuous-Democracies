import logging

from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from app.config import settings
from app.middleware import RequestLoggingMiddleware
from app.routers import health, party, politicians, voting

logging.basicConfig(level=getattr(logging, settings.LOG_LEVEL, logging.WARNING))

app = FastAPI(title="Parliament Monitor API", version="2.0.0")

app.add_middleware(RequestLoggingMiddleware)

app.include_router(health.router)
app.include_router(party.router)
app.include_router(politicians.router)
app.include_router(voting.router)


@app.exception_handler(Exception)
async def unhandled_exception_handler(request: Request, exc: Exception):
    logging.getLogger(__name__).exception("Unhandled error: %s", exc)
    return JSONResponse(status_code=500, content="Internal server error")
