from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    DBSERVERADDRESS: str = "localhost"
    DBNAME: str = "parliamentdb"
    DBUSER: str = ""
    DBPASSWORD: str = ""
    DBPORT: int = 5432
    LOG_LEVEL: str = "WARNING"


settings = Settings()
