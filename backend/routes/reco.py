from flask import Blueprint, jsonify, request

from recommendations.api import recommend
from recommendations.events import process_batch
from recommendations.redis_client import ping
from recommendations.vector import clear_vector, get_vector

recommendations_bp = Blueprint("recommendations", __name__)


@recommendations_bp.route("/health", methods=["GET"])
def health():
    return jsonify({"redis": "ok" if ping() else "fail"})


@recommendations_bp.route("/events", methods=["POST"])
def events():
    """
    Принимает батч событий.
    Body: {"events": [{...}, {...}]}
    """
    data = request.get_json(silent=True)
    if not data:
        return jsonify({"error": "no data"}), 400

    events_list = data.get("events", [])
    if not isinstance(events_list, list):
        return jsonify({"error": "events must be a list"}), 400

    if len(events_list) > 100:
        return jsonify({"error": "too many events (max 100)"}), 400

    processed = process_batch(events_list)
    return jsonify({"processed": processed}), 200


@recommendations_bp.route("/recommendations/<user_id>", methods=["POST"])
def get_recommendations(user_id):
    """
    Возвращает топ-K рекомендаций.
    Body: {
        "candidates": [{notebook}, ...],
        "k": 10,
        "is_cold": false,
        "min_pop": 0.0,
        "max_pop": 10.0
    }
    """
    data = request.get_json(silent=True) or {}

    candidates = data.get("candidates", [])
    if not candidates:
        return jsonify({"error": "no candidates"}), 400

    k = int(data.get("k", 10))
    k = max(1, min(k, 50))

    is_cold = bool(data.get("is_cold", False))
    min_pop = float(data.get("min_pop", 0.0))
    max_pop = float(data.get("max_pop", 10.0))

    result = recommend(
        user_id=user_id,
        candidates=candidates,
        k=k,
        is_cold=is_cold,
        min_pop=min_pop,
        max_pop=max_pop,
    )

    return jsonify({"recommendations": result}), 200


@recommendations_bp.route("/vector/<user_id>", methods=["GET"])
def get_user_vector(user_id):
    """Отладочный эндпоинт — показывает вектор пользователя."""
    return jsonify({"vector": get_vector(user_id)}), 200


@recommendations_bp.route("/vector/<user_id>", methods=["DELETE"])
def delete_user_vector(user_id):
    """Удаление вектора (для тестов)."""
    clear_vector(user_id)
    return jsonify({"status": "deleted"}), 200
