# routes/opportunities.py
import json
import os
from datetime import datetime
from flask import Blueprint, request, jsonify

opportunities_bp = Blueprint("opportunities", __name__)

CACHE = {
    "data": [],
    "updated_at": None
}

def parse_universities():
    all_opportunities = []
    
    json_path = os.path.join(os.path.dirname(__file__), '..', '..', 'assets', 'uni.json')
    
    try:
        with open(json_path, 'r', encoding='utf-8') as f:
            universities = json.load(f)
        
        for uni in universities:
            all_opportunities.append({
                "id": f"uni-{uni['name'].lower().replace(' ', '-').replace('.', '')}",
                "title": f"Поступление в {uni['name']}",
                "company": uni['name'],
                "city": uni.get('city', 'Москва'),
                "type": "university",
                "direction": uni.get('direction', 'other'),
                "format": "office",
                "salary": "стипендия / бюджет",
                "description": uni.get('description', ''),
                "url": uni['url'],
                "date_posted": datetime.now().strftime("%Y-%m-%d"),
                "tags": uni.get('programs', [])
            })
            
        print(f"[Uni] Загружено {len(all_opportunities)} университетов из JSON")
        
    except Exception as e:
        print(f"[Uni] Ошибка загрузки JSON: {e}")
    
    return all_opportunities


@opportunities_bp.route("/ai/v1/opportunities", methods=["GET"])
def get_opportunities():
    global CACHE
    
    if not CACHE["data"]:
        print("[OPPORTUNITIES] Загружаем данные...")
        CACHE["data"] = parse_universities()
        CACHE["updated_at"] = datetime.now().isoformat()
        print(f"[OPPORTUNITIES] Загружено {len(CACHE['data'])}")
    
    result = CACHE["data"]
    
    opp_type = request.args.get("type")
    if opp_type:
        result = [o for o in result if o["type"] == opp_type]
    
    direction = request.args.get("direction")
    if direction:
        result = [o for o in result if o["direction"] == direction]
    
    city = request.args.get("city")
    if city:
        city_lower = city.lower()
        result = [o for o in result if city_lower in o["city"].lower()]
    
    page = int(request.args.get("page", 1))
    limit = int(request.args.get("limit", 20))
    start = (page - 1) * limit
    end = start + limit
    
    return jsonify({
        "opportunities": result[start:end],
        "total": len(result),
        "page": page,
        "limit": limit,
        "updated_at": CACHE["updated_at"]
    })