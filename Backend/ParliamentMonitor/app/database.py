import os
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, Session
from typing import Generator

def get_connection_string() -> str:
    server = os.environ.get("DBSERVERADDRESS", "localhost")
    name = os.environ.get("DBNAME", "parlimentdb")
    user = os.environ.get("DBUSER", "bobo")
    password = os.environ.get("DBPASSWORD", "password123")
    port = os.environ.get("DBPORT", "5432")
    return f"postgresql://{user}:{password}@{server}:{port}/{name}"

engine = create_engine(get_connection_string(), pool_pre_ping=True)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

def get_db() -> Generator[Session, None, None]:
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
