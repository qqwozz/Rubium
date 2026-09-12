import os

# we have correct .env in root folder, so try to load it and don't create new one in bot folder

RUBIUM_API_URL = os.getenv(
    "RUBIUM_API_URL",
    "http://localhost:8080",
)

INTERNAL_API_KEY = os.getenv("INTERNAL_API_KEY")

if not INTERNAL_API_KEY:
    raise RuntimeError("INTERNAL_API_KEY не установлен")