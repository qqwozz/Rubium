import math
from datetime import datetime
from typing import Dict, List

from .similarity import similarity

# Веса компонентов
W_SIM = 0.4
W_FRESH = 0.2
W_QUAL = 0.2
W_POP = 0.2

# Freshness
T_HALF = 14  # дней

# Quality
C_GLOBAL = 3.5  # глобальная средняя оценка
M_TRUST = 10    # порог доверия


def freshness(notebook: Dict) -> float:
    """Свежесть тетради. 0..1, где 1 — только что обновлена."""
    updated_at = notebook.get("updated_at")
    if not updated_at:
        return 0.0
    
    try:
        dt = datetime.fromisoformat(updated_at.replace("Z", "+00:00"))
    except Exception:
        return 0.0
    
    now = datetime.now(dt.tzinfo)
    days_passed = (now - dt).total_seconds() / 86400
    
    return 0.5 ** (days_passed / T_HALF)


def quality(notebook: Dict) -> float:
    """Байесовское среднее оценок, нормированное в 0..1."""
    R = float(notebook.get("average_rating") or 0)
    v = int(notebook.get("ratings_count") or 0)
    
    bayes = (R * v + C_GLOBAL * M_TRUST) / (v + M_TRUST)
    
    return bayes / 5.0


def popularity(notebook: Dict) -> float:
    """Логарифм от действий."""
    views = int(notebook.get("views_count") or 0)
    saves = int(notebook.get("saves_count") or 0)
    
    return math.log(1 + views + 2 * saves)


def normalize(value: float, min_v: float, max_v: float) -> float:
    """Min-max нормировка."""
    if max_v == min_v:
        return 0.0
    return (value - min_v) / (max_v - min_v)


def score(user_id: str, notebook: Dict, 
          min_pop: float = 0.0, max_pop: float = 10.0) -> float:
    """
    Итоговый скор тетради для пользователя.
    min_pop/max_pop нужны для нормировки popularity.
    """
    sim = similarity(user_id, notebook)
    fresh = freshness(notebook)
    qual = quality(notebook)
    pop = popularity(notebook)
    pop_norm = normalize(pop, min_pop, max_pop)
    
    return (
        W_SIM * sim +
        W_FRESH * fresh +
        W_QUAL * qual +
        W_POP * pop_norm
    )


def score_cold_start(notebook: Dict, min_pop: float = 0.0, max_pop: float = 10.0) -> float:
    """Скор для нового пользователя (без similarity)."""
    fresh = freshness(notebook)
    qual = quality(notebook)
    pop = popularity(notebook)
    pop_norm = normalize(pop, min_pop, max_pop)
    
    total_w = W_FRESH + W_QUAL + W_POP
    return (
        W_FRESH * fresh +
        W_QUAL * qual +
        W_POP * pop_norm
    ) / total_w