from typing import Optional, Dict, Any
from .vector import apply_decay, boost_field, update_last_updated, get_vector

# Веса действий (b_action из документа)
ACTION_BOOSTS = {
    "open_notebook": 0.3,
    "read_page": 0.5,
    "scroll": 0.3,
    "return_notebook": 0.7,
    "save_notebook": 1.0,
    "rate_notebook": 0.7,
    "subscribe_author": 1.5,
}


def process_event(event: Dict[str, Any]) -> bool:
    """
    Обрабатывает одно событие.
    
    event = {
        "user_id": "uuid",
        "type": "open_notebook",
        "target_id": "uuid",
        "tags": ["math", "algebra"],
        "author_id": "uuid",
        "timestamp": 1725960000
    }
    
    Возвращает True если обработано, False если пропущено.
    """
    user_id = event.get("user_id")
    event_type = event.get("type")
    
    if not user_id or not event_type:
        return False
    
    boost = ACTION_BOOSTS.get(event_type)
    if boost is None:
        return False
    
    # 1. затухание всего вектора
    apply_decay(user_id)
    
    # 2. буст по предметам (тегам тетради)
    tags = event.get("tags") or []
    if tags:
        per_tag = boost / len(tags)
        for tag in tags:
            boost_field(user_id, f"subject:{tag}", per_tag)
    
    # 3. буст по автору
    author_id = event.get("author_id")
    if author_id:
        boost_field(user_id, f"author:{author_id}", boost)
    
    # 4. буст по часу (из timestamp события)
    ts = event.get("timestamp")
    if ts:
        from datetime import datetime
        hour = datetime.fromtimestamp(ts).hour
        boost_field(user_id, f"hour:{hour}", boost)
    
    # 5. обновляем timestamp
    update_last_updated(user_id)
    
    return True


def process_batch(events: list) -> int:
    """Обрабатывает список событий. Возвращает количество обработанных."""
    count = 0
    for event in events:
        if process_event(event):
            count += 1
    return count