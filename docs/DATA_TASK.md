
ЗАДАЧА: Собрать 50k+ пар (input, output) для обучения модели text → LaTeX.

## Формат

JSONL, одна пара на строку:
{"input": "...", "output": "...", "patterns": ["power", "fraction"]}

Где:

- input — как студент пишет формулу текстом (НЕ LaTeX).
- output — валидный LaTeX.
- patterns — список типов конструкций в формуле.

## Что такое patterns

Разложи каждую формулу на атомарные типы. Пример:

Input:  "x^2 + y^2"
Output: "x^{2} + y^{2}"
Patterns: ["power", "sum", "power"]

Input:  "MSE = 1/n * Σ(y_i - ŷ_i)²"
Output: "\\text{MSE} = \\frac{1}{n} \\sum (y_i - \\hat{y}_i)^{2}"
Patterns: ["equation", "fraction", "sum", "indexed", "hatted", "subtraction", "power"]

## Список patterns (используй только эти)

- power (x^2)
- index (x_i)
- fraction (a/b)
- sqrt (sqrt(x))
- sum (Σ)
- integral (∫)
- limit (lim)
- derivative (deriv, f')
- log (log, ln)
- trig (sin, cos, tan)
- equation (=)
- sum (a+b)
- subtraction (a-b)
- multiplication (a*b)
- hatted (ŷ)
- barred (x̄)
- greek (α, β)
- indexed (x_i)
- nested (скобки внутри скобок)

## Распределение (обязательно)

Собери ровно:

- 5000 пар — только один pattern (например, только power).
- 15000 пар — два pattern (power + sum).
- 20000 пар — три-четыре pattern.
- 10000 пар — пять+ pattern (сложные).

Итого: 50000.

## Источники

1. ChatGPT — генерируй по 200 пар на промпт, для каждой категории:

   - Степени (x^N, y^N, a^N).
   - Дроби (числа, переменные, смешанные).
   - Корни (sqrt, cbrt).
   - Тригонометрия (sin, cos, tan + аргументы).
   - Логарифмы (log, ln + основания).
   - Производные (deriv, f', d/dx).
   - Интегралы (с пределами и без).
   - Суммы (с пределами).
   - Пределы.
   - Уравнения (=).
   - ML-формулы (MSE, MAE, R², precision, recall).
   - Физика (E=mc², F=ma, etc).
2. Русский + английский — 30% примеров на русском:

   - "корень из x"
   - "x в квадрате"
   - "сумма от i до n"
   - "производная f"
   - "предел при x→0"
3. Разные стили input для одной формулы:

   - "x^2", "x squared", "x в квадрате", "X^2".

## Требования

- Каждый output — валидный LaTeX (проверь через KaTeX).
- Каждый input — реалистичный (как пишет студент).
- Никаких дубликатов (проверяй уникальность input).
- Каждая строка содержит поле patterns.
- 30% примеров — с опечатками (derivativ, производня, sqrtt).
- 20% примеров — с пробелами внутри формул (x ^ 2, a / b).

## Формат сдачи

- data/train_part1.jsonl — 25k пар.
- data/train_part2.jsonl — 25k пар.
- Проверь каждую строку:
  - Валидный JSON.
  - Есть поля input, output, patterns.
  - Output — валидный LaTeX (можно отрендерить).
  - Input уникален.

## Сроки

5 дней.
