import httpx

from config import RUBIUM_API_URL, INTERNAL_API_KEY


def get_notebooks_by_tag(tag: str, limit: int = 50) -> list[dict]:
    response = httpx.get(
        f"{RUBIUM_API_URL}/internal/v1/notebooks/by-tag",
        params={
            "tag": tag,
            "limit": limit,
        },
        headers={
            "X-Internal-Key": INTERNAL_API_KEY,
        },
        timeout=10.0,
    )

    if response.status_code == 401:
        raise RuntimeError("Неверный INTERNAL_API_KEY")

    if response.status_code == 400:
        data = response.json()
        raise ValueError(data.get("error", "Невалидный запрос"))

    response.raise_for_status()

    data = response.json()

    '''
    response example:
    {
        "notebooks": [
            {
            "id": "uuid",
            "title": "Математика ЕГЭ",
            "description": "Подготовка к ЕГЭ",
            "color": "#ffffff",
            "tags": [
                "математика",
                "егэ"
            ],
            "is_public": true,
            "sections_count": 5,
            "pages_count": 24,
            "views_count": 120,
            "copies_count": 15,
            "average_rating": 4.8,
            "ratings_count": 20,
            "created_at": "2026-09-01T12:00:00Z",
            "updated_at": "2026-09-05T15:30:00Z"
            }
        ]
    }

    example using:
    python: notebooks = get_notebooks_by_tag("математика")

    get:
    [
        {
            "id": "...",
            "title": "Математика ЕГЭ",
            "description": "Подготовка к ЕГЭ",
            "color": "#ffffff",
            "tags": ["математика", "егэ"],
            "is_public": True,
            "sections_count": 5,
            "pages_count": 24,
            "views_count": 120,
            "copies_count": 15,
            "average_rating": 4.8,
            "ratings_count": 20,
            "created_at": "...",
            "updated_at": "..."
        }
    ]
    '''
    return data.get("notebooks", [])