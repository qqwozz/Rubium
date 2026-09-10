import os
from pathlib import Path
from dotenv import load_dotenv
import redis

load_dotenv(Path(__file__).parent.parent.parent / '.env')

REDIS_HOST = os.getenv("REDIS_HOST", "localhost")
REDIS_PORT = int(os.getenv("REDIS_PORT", 6379))
REDIS_DB = int(os.getenv("REDIS_DB", 0))

_client = None


def get_client() -> redis.Redis:
    global _client
    if _client is None:
        _client = redis.Redis(
            host=REDIS_HOST,
            port=REDIS_PORT,
            db=REDIS_DB,
            decode_responses=True
        )
    return _client


def ping() -> bool:
    try:
        return get_client().ping()
    except Exception as e:
        print(f"Redis connection failed: {e}")
        return False