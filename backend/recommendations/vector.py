import time
from typing import Dict, List, Tuple

from .redis_client import get_client

DECAY_RATE = 0.98  # затухание в час
VECTOR_TTL = 60 * 60 * 24 * 90  # 90 дней неактивности → удаляем


def _key(user_id: str) -> str:
    return f"user:vector:{user_id}"


def get_vector(user_id: str) -> Dict[str, float]:
    """Читает весь вектор пользователя из Redis."""
    client = get_client()
    data = client.hgetall(_key(user_id))

    if not data:
        return {}

    vector = {}
    for field, value in data.items():
        if field == "_last_updated":
            continue
        try:
            vector[field] = float(value)
        except ValueError:
            continue

    return vector


def get_last_updated(user_id: str) -> float:
    """Возвращает timestamp последнего обновления вектора."""
    client = get_client()
    value = client.hget(_key(user_id), "_last_updated")

    if not value:
        return time.time()

    try:
        return float(value)
    except ValueError:
        return time.time()


def apply_decay(user_id: str) -> Dict[str, float]:
    """
    Применяет затухание ко всем полям вектора.
    Возвращает вектор после затухания.
    """
    client = get_client()
    vector = get_vector(user_id)
    last_updated = get_last_updated(user_id)

    now = time.time()
    hours_passed = (now - last_updated) / 3600

    if hours_passed < 0.01:
        return vector

    decay_factor = DECAY_RATE ** hours_passed

    if vector:
        updated = {k: v * decay_factor for k, v in vector.items()}
        client.hset(_key(user_id), mapping=updated)
        return updated

    return {}


def boost_field(user_id: str, field: str, amount: float) -> None:
    """Добавляет boost к одному полю вектора."""
    client = get_client()
    client.hincrbyfloat(_key(user_id), field, amount)


def update_last_updated(user_id: str) -> None:
    """Обновляет timestamp последнего обновления."""
    client = get_client()
    client.hset(_key(user_id), "_last_updated", time.time())
    client.expire(_key(user_id), VECTOR_TTL)


def clear_vector(user_id: str) -> None:
    """Полностью удаляет вектор пользователя."""
    client = get_client()
    client.delete(_key(user_id))


def get_top_fields(user_id: str, prefix: str, n: int = 10) -> List[Tuple[str, float]]:
    """
    Возвращает топ-N полей с указанным префиксом.
    Например, get_top_fields(uid, "subject:", 5) → [("subject:math", 4.2), ...]
    """
    vector = get_vector(user_id)
    filtered = [(k, v) for k, v in vector.items() if k.startswith(prefix)]
    filtered.sort(key=lambda x: x[1], reverse=True)
    return filtered[:n]
