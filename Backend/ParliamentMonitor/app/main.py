import os
import logging
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from app.middleware import RequestLoggingMiddleware
from app.routers import party, politicians, voting

log_level = os.environ.get("LOG_LEVEL", "INFO").upper()
logging.basicConfig(
    level=getattr(logging, log_level, logging.INFO),
    format="%(asctime)s %(levelname)s %(name)s: %(message)s",
)

app = FastAPI(
    title="API democratia continua",
    version="v1",
)

app.add_middleware(RequestLoggingMiddleware)

app.include_router(party.router)
app.include_router(politicians.router)
app.include_router(voting.router)


@app.get("/health")
def health_check():
    return "Healthy"


@app.exception_handler(Exception)
async def global_exception_handler(request: Request, exc: Exception):
    logging.getLogger(__name__).error("Unhandled exception: %s", exc, exc_info=True)
    return JSONResponse(status_code=500, content={"error": str(exc)})
