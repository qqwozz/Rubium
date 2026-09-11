
Что нужно

Добавить в `GET /notebooks/community` новый параметр `sort=for_you`, который сортирует тетради через Python-сервис.

## Пайплайн

```
Фронт → GET /notebooks/community?sort=for_you
              │
              ▼
        Go берёт тетради из Supabase
              │
              ▼
        Go фильтрует кандидатов:
          - is_public = true
          - не свои (user_id != current)
          - топ-100 по (rating + freshness)
              │
              ▼
        POST http://backend:5080/recommendations/{user_id}
        Body: {
          "candidates": [...],
          "k": 20,
          "is_cold": false,
          "min_pop": 0.0,
          "max_pop": 10.0
        }
              │
              ▼
        Python возвращает отсортированные тетради
              │
              ▼
        Go отдаёт фронту (тот же формат что обычный /community)
```

## API-эндпоинт Python

```
POST /recommendations/{user_id}
Content-Type: application/json

{
  "candidates": [
    {
      "id": "uuid",
      "tags": ["math", "algebra"],
      "author_id": "uuid",
      "updated_at": "2026-09-05T00:00:00+00:00",
      "average_rating": 4.5,
      "ratings_count": 10,
      "views_count": 100,
      "saves_count": 5
    }
  ],
  "k": 20,
  "is_cold": false,
  "min_pop": 0.0,
  "max_pop": 10.0
}

Response:
{
  "recommendations": [ ...тот же объект + поле score ]
}
```

## Что важно

1. **Fallback:** если Python недоступен (`curl` timeout 500ms, connection error) → сортировать по рейтингу. Не ломать каталог.
2. **Таймаут:** максимум 500 мс. Python должен успевать.
3. **Cold start:** если у юзера нет вектора (Python вернёт всё как есть) → тоже fallback на рейтинг.
4. **`min_pop`/`max_pop`:** можно сначала посчитать в Go по всем кандидатам или захардкодить `0` и `10`.
5. **URL сервиса:** `http://backend:5080` (имя контейнера в docker-compose).

## Тест

```bash
curl "http://localhost:8080/notebooks/community?sort=for_you" \
  -H "Authorization: Bearer <token>"
```

Должен вернуть тетради в порядке, отличном от `sort=rating`.
