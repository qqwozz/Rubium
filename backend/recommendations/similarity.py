import math
from typing import Dict, List, Set
from .vector import get_vector

TEMPERATURE = 1.0
ALPHA = 0.7  # вес similarity по предметам
BETA = 0.3   # вес similarity по авторам


def softmax(values: Dict[str, float], temperature: float = TEMPERATURE) -> Dict[str, float]:
    """
    Softmax по словарю {key: value} → {key: probability}.
    """
    if not values:
        return {}
    
    max_v = max(values.values())
    exps = {k: math.exp((v - max_v) / temperature) for k, v in values.items()}
    total = sum(exps.values())
    
    if total == 0:
        return {k: 0.0 for k in values}
    
    return {k: v / total for k, v in exps.items()}


def extract_axis(vector: Dict[str, float], prefix: str) -> Dict[str, float]:
    """Достаёт ось из вектора: subject:, author:, hour:."""
    return {k[len(prefix):]: v for k, v in vector.items() if k.startswith(prefix)}


def cosine_similarity(a: Dict[str, float], b: Dict[str, float]) -> float:
    """
    Косинусное сходство между двумя словарями.
    Ключи не совпадающие — игнорируются.
    """
    if not a or not b:
        return 0.0
    
    common = set(a.keys()) & set(b.keys())
    if not common:
        return 0.0
    
    dot = sum(a[k] * b[k] for k in common)
    norm_a = math.sqrt(sum(v * v for v in a.values()))
    norm_b = math.sqrt(sum(v * v for v in b.values()))
    
    if norm_a == 0 or norm_b == 0:
        return 0.0
    
    return dot / (norm_a * norm_b)


def similarity_subject(user_id: str, notebook_tags: List[str]) -> float:
    """Similarity между предпочтениями юзера и тегами тетради."""
    vector = get_vector(user_id)
    subjects = extract_axis(vector, "subject:")
    
    if not subjects or not notebook_tags:
        return 0.0
    
    # softmax по предпочтениям
    probs = softmax(subjects)
    
    # one-hot / равномерное распределение по тегам тетради
    n = len(notebook_tags)
    q = {tag: 1.0 / n for tag in notebook_tags}
    
    # нормализуем p (softmax уже даёт сумму 1)
    return cosine_similarity(probs, q)


def similarity_author(user_id: str, author_id: str) -> float:
    """Similarity между предпочтениями юзера и автором тетради."""
    if not author_id:
        return 0.0
    
    vector = get_vector(user_id)
    authors = extract_axis(vector, "author:")
    
    if not authors or author_id not in authors:
        return 0.0
    
    probs = softmax(authors)
    return probs.get(author_id, 0.0)


def similarity(user_id: str, notebook: Dict) -> float:
    """
    Итоговая similarity между юзером и тетрадью.
    notebook = {"tags": [...], "author_id": "uuid"}
    """
    tags = notebook.get("tags") or []
    author_id = notebook.get("author_id")
    
    sim_subj = similarity_subject(user_id, tags)
    sim_auth = similarity_author(user_id, author_id)
    
    return ALPHA * sim_subj + BETA * sim_auth


def jaccard(a: Set[str], b: Set[str]) -> float:
    """Jaccard similarity для тегов."""
    if not a or not b:
        return 0.0
    inter = len(a & b)
    union = len(a | b)
    if union == 0:
        return 0.0
    return inter / union