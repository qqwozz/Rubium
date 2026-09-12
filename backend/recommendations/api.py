import random
from typing import Dict, List

from .scoring import score, score_cold_start
from .similarity import jaccard

LAMBDA_MMR = 0.7
EPSILON = 0.1


def mmr_rank(
    user_id: str, candidates: List[Dict], k: int = 10, min_pop: float = 0.0, max_pop: float = 10.0
) -> List[Dict]:
    """
    Ранжирование с MMR (Maximal Marginal Relevance).
    Возвращает топ-K тетрадей.
    """
    if not candidates:
        return []

    # первый этап: считаем score
    scored = [(nb, score(user_id, nb, min_pop, max_pop)) for nb in candidates]
    scored.sort(key=lambda x: x[1], reverse=True)

    selected = []
    selected_tags: List[set] = []

    # первый элемент — лучший по score
    first_nb, _ = scored.pop(0)
    selected.append(first_nb)
    selected_tags.append(set(first_nb.get("tags") or []))

    # остальные — через MMR
    while len(selected) < k and scored:
        best_idx = -1
        best_mmr = -float("inf")

        for idx, (nb, s) in enumerate(scored):
            tags = set(nb.get("tags") or [])
            max_sim = max((jaccard(tags, st) for st in selected_tags), default=0.0)
            mmr = LAMBDA_MMR * s - (1 - LAMBDA_MMR) * max_sim

            if mmr > best_mmr:
                best_mmr = mmr
                best_idx = idx

        if best_idx < 0:
            break

        nb, _ = scored.pop(best_idx)
        selected.append(nb)
        selected_tags.append(set(nb.get("tags") or []))

    # exploration: заменяем последний элемент на случайный
    if random.random() < EPSILON and len(candidates) > len(selected):
        pool = [c for c in candidates if c not in selected]
        if pool:
            selected[-1] = random.choice(pool)

    return selected


def recommend_cold_start(
    candidates: List[Dict], k: int = 10, min_pop: float = 0.0, max_pop: float = 10.0
) -> List[Dict]:
    """Рекомендации для нового пользователя."""
    scored = [(nb, score_cold_start(nb, min_pop, max_pop)) for nb in candidates]
    scored.sort(key=lambda x: x[1], reverse=True)
    return [nb for nb, _ in scored[:k]]


def recommend(
    user_id: str,
    candidates: List[Dict],
    k: int = 10,
    is_cold: bool = False,
    min_pop: float = 0.0,
    max_pop: float = 10.0,
) -> List[Dict]:
    """Главная функция рекомендаций."""
    if is_cold:
        return recommend_cold_start(candidates, k, min_pop, max_pop)
    return mmr_rank(user_id, candidates, k, min_pop, max_pop)
