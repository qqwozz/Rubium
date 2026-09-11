# Rubium — Рекомендательная система

> Математическая модель персонализированных рекомендаций тетрадей.
> Не ML. Эвристики на основе статистики действий.

---

## 1. Обозначения

| Символ      | Значение                                                                        |
| ----------------- | --------------------------------------------------------------------------------------- |
| $u$             | пользователь                                                                |
| $i$             | тетрадь (материал)                                                       |
| $\mathcal{I}$   | множество всех публичных тетрадей                         |
| $\mathcal{I}_u$ | тетради, которые$u$ уже видел                                   |
| $p_u$           | вектор предпочтений пользователя                          |
| $q_i$           | вектор признаков тетради                                          |
| $\Delta t$      | время с последнего обновления (в часах или днях) |

---

## 2. Вектор пользователя $p_u$

Три оси:

$$
p_u = \left( p_u^{\text{subject}},\ p_u^{\text{author}},\ p_u^{\text{hour}} \right)
$$

### 2.1. Ось предметов

$p_u^{\text{subject}} \in \mathbb{R}^{S}$, где $S$ — количество предметов (тегов).

Обновляется при каждом действии:

$$
p_u^{\text{subject}}[s] \leftarrow p_u^{\text{subject}}[s] \cdot \gamma^{\Delta t} + b_{\text{action}} \cdot \mathbb{1}[s \in \text{tags}(i)]
$$

Где:

- $\gamma = 0.98$ — коэффициент затухания в час
- $b_{\text{action}}$ — вес действия (см. таблицу)
- $\mathbb{1}[\cdot]$ — индикатор

### 2.2. Ось авторов

$p_u^{\text{author}} \in \mathbb{R}^{A}$, где $A$ — количество авторов.

$$
p_u^{\text{author}}[a] \leftarrow p_u^{\text{author}}[a] \cdot \gamma^{\Delta t} + b_{\text{action}} \cdot \mathbb{1}[a = \text{author}(i)]
$$

### 2.3. Ось часов

$p_u^{\text{hour}} \in \mathbb{R}^{24}$

$$
p_u^{\text{hour}}[h] \leftarrow p_u^{\text{hour}}[h] \cdot \gamma^{\Delta t} + b_{\text{action}} \cdot \mathbb{1}[h = \text{hour}(\text{now})]
$$

### 2.4. Веса действий $b_{\text{action}}$

| Действие                              | $b_{\text{action}}$ |
| --------------------------------------------- | --------------------- |
| Открыл тетрадь                   | 0.3                   |
| Провёл >30 сек на странице | 0.5                   |
| Проскроллил >50%                   | 0.3                   |
| Вернулся на тетрадь          | 0.7                   |
| Сохранил тетрадь               | 1.0                   |
| Оценил тетрадь                   | 0.7                   |
| Подписался на автора        | 1.5                   |

---

## 3. Softmax по каждой оси

Вектор предпочтений превращается в распределение вероятностей:

$$
\sigma(p_u^{\text{subject}})_s = \frac{\exp(p_u^{\text{subject}}[s] / T)}{\sum_{s'} \exp(p_u^{\text{subject}}[s'] / T)}
$$

Где:

- $T$ — температура
- $T = 1$ — стандарт
- $T < 1$ — резкие пики (более сфокусировано)
- $T > 1$ — сглажено

Аналогично для авторов и часов.

---

## 4. Вектор тетради $q_i$

### 4.1. Ось предметов

Равномерное распределение по тегам тетради:

$$
q_i^{\text{subject}}[s] = \frac{1}{|\text{tags}(i)|} \cdot \mathbb{1}[s \in \text{tags}(i)]
$$

### 4.2. Ось авторов

One-hot на автора тетради:

$$
q_i^{\text{author}}[a] = \mathbb{1}[a = \text{author}(i)]
$$

### 4.3. Ось часов

Не применима к тетради. Выпадает из similarity.

---

## 5. Similarity

### 5.1. По предметам

Косинусное сходство между softmax-вектором пользователя и вектором тетради:

$$
\text{sim}_{\text{subject}}(u, i) = \frac{\sigma(p_u^{\text{subject}}) \cdot q_i^{\text{subject}}}{\|\sigma(p_u^{\text{subject}})\| \cdot \|q_i^{\text{subject}}\|}
$$

### 5.2. По авторам

$$
\text{sim}_{\text{author}}(u, i) = \frac{\sigma(p_u^{\text{author}}) \cdot q_i^{\text{author}}}{\|\sigma(p_u^{\text{author}})\| \cdot \|q_i^{\text{author}}\|}
$$

### 5.3. Итоговое

$$
\text{sim}(u, i) = \alpha \cdot \text{sim}_{\text{subject}} + \beta \cdot \text{sim}_{\text{author}}
$$

Где $\alpha + \beta = 1$. Начальные значения: $\alpha = 0.7$, $\beta = 0.3$.

---

## 6. Остальные компоненты

### 6.1. Freshness (свежесть)

Экспоненциальное затухание от времени обновления:

$$
\text{fresh}(i) = 0.5^{\Delta t_i / T_{1/2}}
$$

Где:

- $\Delta t_i$ — дней с момента `updated_at`
- $T_{1/2} = 14$ — период полураспада в днях

### 6.2. Quality (качество)

Байесовское среднее (как у IMDb):

$$
\text{qual}(i) = \frac{R_i \cdot v_i + C \cdot m}{v_i + m}
$$

Где:

- $R_i$ — средняя оценка тетради
- $v_i$ — число оценок
- $C = 3.5$ — глобальная средняя оценка
- $m = 10$ — порог доверия

**Пример:**

Материал с 1 оценкой 5.0:

$$
\text{qual} = \frac{5 \cdot 1 + 3.5 \cdot 10}{1 + 10} = \frac{40}{11} \approx 3.64
$$

Материал с 100 оценками 5.0:

$$
\text{qual} = \frac{5 \cdot 100 + 3.5 \cdot 10}{100 + 10} = \frac{535}{110} \approx 4.86
$$

### 6.3. Popularity (популярность)

Логарифм от действий (чтобы топы не доминировали):

$$
\text{pop}(i) = \log(1 + \text{views}_i + 2 \cdot \text{saves}_i)
$$

---

## 7. Итоговый скор

$$
\text{score}(u, i) = w_1 \cdot \text{sim}(u,i) + w_2 \cdot \text{fresh}(i) + w_3 \cdot \text{qual}(i) + w_4 \cdot \text{pop}(i)
$$

### 7.1. Нормировка компонентов

Каждый компонент нормируется в $[0, 1]$ перед суммой:

$$
\text{norm}(x) = \frac{x - x_{\min}}{x_{\max} - x_{\min}}
$$

Или сигмоида:

$$
\text{norm}(x) = \frac{1}{1 + e^{-x}}
$$

### 7.2. Веса

Начальные значения: $w_1 = 0.4$, $w_2 = 0.2$, $w_3 = 0.2$, $w_4 = 0.2$.

Сумма $\sum w_k = 1$.

---

## 8. Diversity (MMR)

Maximal Marginal Relevance — чтобы топ не был однообразным.

Первый элемент — с максимальным score.

Каждый следующий:

$$
\text{MMR}(i) = \lambda \cdot \text{score}(u, i) - (1 - \lambda) \cdot \max_{j \in \mathcal{S}} \text{sim}(i, j)
$$

Где:

- $\mathcal{S}$ — уже выбранные тетради
- $\lambda = 0.7$ — баланс между релевантностью и разнообразием

### 8.1. Similarity между тетрадями

Jaccard по тегам:

$$
\text{sim}(i, j) = \frac{|\text{tags}(i) \cap \text{tags}(j)|}{|\text{tags}(i) \cup \text{tags}(j)|}
$$

---

## 9. Exploration vs Exploitation

ε-greedy:

$$
i^* = \begin{cases}
\arg\max_i \text{score}(u, i) & \text{с вероятностью } 1 - \varepsilon \\
\text{random}(\mathcal{I} \setminus \mathcal{I}_u) & \text{с вероятностью } \varepsilon
\end{cases}
$$

$\varepsilon = 0.1$ — 10% показываем случайное.

**Зачем:** собрать данные о том, что юзеру может понравиться, но он ещё не видел.

---

## 10. Cold Start

Новый пользователь ($\|p_u\| = 0$):

$$
\text{score}(u, i) = w_2 \cdot \text{fresh}(i) + w_3 \cdot \text{qual}(i) + w_4 \cdot \text{pop}(i)
$$

Без similarity. Показываем общий топ.

---

## 11. Псевдокод

```
function recommend(u, K=10):
    if is_cold_start(u):
        candidates = all_public_notebooks()
        ranked = sort_by(candidates, key=score_cold)
        return top_K(ranked, K)
  
    candidates = all_public_notebooks() - seen(u)
  
    scored = []
    for i in candidates:
        s = w1 * similarity(p_u, q_i)
          + w2 * freshness(i)
          + w3 * quality(i)
          + w4 * popularity(i)
        scored.append((i, s))
  
    # MMR
    selected = []
    selected.append(argmax(scored, key=score))
    scored.remove(selected[0])
  
    while len(selected) < K:
        for (i, s) in scored:
            mmr = lambda * s - (1 - lambda) * max_similarity(i, selected)
            ...
        selected.append(argmax(scored, key=mmr))
        scored.remove(...)
  
    # exploration
    if random() < epsilon:
        selected[-1] = random_choice(candidates - selected)
  
    return selected
```

---

## 12. Параметры для настройки

| Параметр | Значение | Что значит                                      |
| ---------------- | ---------------- | -------------------------------------------------------- |
| $\gamma$       | 0.98             | Коэффициент затухания в час      |
| $T$            | 1.0              | Температура softmax                           |
| $\alpha$       | 0.7              | Вес similarity по предметам                |
| $\beta$        | 0.3              | Вес similarity по авторам                    |
| $T_{1/2}$      | 14               | Период полураспада freshness (дней) |
| $C$            | 3.5              | Глобальная средняя оценка         |
| $m$            | 10               | Порог доверия для quality                 |
| $w_1$          | 0.4              | Вес similarity в score                               |
| $w_2$          | 0.2              | Вес freshness в score                                |
| $w_3$          | 0.2              | Вес quality в score                                  |
| $w_4$          | 0.2              | Вес popularity в score                               |
| $\lambda$      | 0.7              | Баланс MMR                                         |
| $\varepsilon$  | 0.1              | Exploration rate                                         |

---

## 13. Этапы реализации

1. **Сбор событий** — фронт шлёт, Python принимает, Redis хранит
2. **Обновление вектора** — decay + boost при каждом событии
3. **API рекомендаций** — top-K с MMR
4. **Cold start** — топ по quality + freshness + popularity
5. **Логирование показов** — для будущего ML
