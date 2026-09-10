
## Go-эндпоинт `GET /sitemap.xml`, генерирующий XML со всеми публичными URL.

Удали файл как сделаешь

**Что включить:**

- Статика: `/`, `/community`, `/courses`, `/rubium_tech`
- Все публичные тетради: `/notebook/{id}`
- Все публичные профили: `/user/{id}` (если применимо)

**Формат:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://rubium.tech/</loc>
    <changefreq>weekly</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>https://rubium.tech/notebook/{uuid}</loc>
    <lastmod>2026-09-10</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.7</priority>
  </url>
</urlset>
```

**Поля:**

- `loc` — URL
- `lastmod` — `updated_at` тетради
- `changefreq` — `monthly` для тетрадей, `weekly` для главной
- `priority` — `1.0` главная, `0.9` каталог, `0.7` тетради, `0.6` профили

**Content-Type:** `application/xml`

**Nginx:** проксировать `/sitemap.xml` → Go API

**Лимит:** максимум 50 000 URL в одном файле. Если больше — разбить на `sitemap-1.xml`, `sitemap-2.xml` + индексный `sitemap.xml`. Пока не актуально, но пусть знает.
