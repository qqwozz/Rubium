package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	maxRetries     = 2
	retryBaseDelay = 200 * time.Millisecond
)

type Client struct {
	baseURL    string
	anonKey    string
	serviceKey string
	http       *http.Client
}

func NewClient(baseURL, anonKey, serviceKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		anonKey:    anonKey,
		serviceKey: serviceKey,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) headers() map[string]string {
	return map[string]string{
		"apikey":       c.anonKey,
		"Authorization": "Bearer " + c.anonKey,
		"Content-Type":  "application/json",
	}
}

func (c *Client) headersService() map[string]string {
	return map[string]string{
		"apikey":       c.serviceKey,
		"Authorization": "Bearer " + c.serviceKey,
		"Content-Type":  "application/json",
	}
}

// do — выполнение запроса с retry на сетевые ошибки и 5 раз
// TODO: написать комментарии к функции
func (c *Client) do(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryBaseDelay * time.Duration(1<<(attempt-1))
			time.Sleep(delay)
		}

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("executing request: %w", err)
			continue
		}

		if resp.StatusCode < 500 {
			return resp, nil
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		lastErr = fmt.Errorf("supabase returned %d: %s", resp.StatusCode, string(body))
	}

	return nil, lastErr
}

// RawQuery — GET-запрос к PostgREST
// TODO: написать комментарии к функции
func (c *Client) RawQuery(endpoint string, useServiceRole bool) ([]byte, error) {
	url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return body, nil
}

// Query — запрос с распаковкой JSON
// TODO: написать комментарии к функции
func (c *Client) Query(endpoint string, useServiceRole bool, result interface{}) error {
	body, err := c.RawQuery(endpoint, useServiceRole)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, result)
}

// Patch — PATCH-запрос к PostgREST (обновление записей)
// TODO: написать комментарии к функции
func (c *Client) Patch(endpoint string, useServiceRole bool, payload interface{}) error {
	url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	h["Prefer"] = "return=minimal"
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

//! POST-запрос к PostgREST (создание записей)
//! Post отправляет POST-запрос к PostgREST API для создания новых записей в базе данных.
func (c *Client) Post(endpoint string, useServiceRole bool, payload interface{}, result interface{}) error {
	// Параметры:
	//   - endpoint: путь к таблице или RPC-функции (например, "users" или "rpc/my_function")
	//   - useServiceRole: если true, использует service role key (обход RLS),
	//     если false - использует anon key (с учетом RLS)
	//   - payload: данные для вставки (структура, map или слайс структур)
	//   - result: указатель на переменную для сохранения созданной записи (может быть nil)
	//
	// Возвращает:
	//   - error: ошибка при выполнении запроса или парсинге ответа


    //* Формируем полный URL для запроса к PostgREST API
    url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

    //* Сериализуем payload в JSON
    // NOTE: json.Marshal может вернуть ошибку при несериализуемых данных
    data, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("marshaling payload: %w", err)
    }

    //* Создаем HTTP-запрос с телом в виде байтового буфера
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    if err != nil {
        return fmt.Errorf("creating request: %w", err)
    }

    //* Выбираем заголовки в зависимости от необходимости service role
    h := c.headers()
    if useServiceRole {
        h = c.headersService()
    }

    // Устанавливаем Prefer header для получения созданной записи
    // representation - возвращает созданные данные, включая дефолтные значения
    // Возможные альтернативы: "return=minimal" - не возвращает данные
    h["Prefer"] = "return=representation"

    //* Применяем все заголовки к запросу
    for k, v := range h {
        req.Header.Set(k, v)
    }

    //* Выполняем запрос с автоматическими повторными попытками
    resp, err := c.do(req)
    if err != nil {
        return err
    }
    // NOTE: закрываем тело ответа после обработки
    defer resp.Body.Close()

    //* Читаем тело ответа
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("reading response: %w", err)
    }

    //* Проверяем статус ответа
    //! POST может возвращать 200 OK или 201 Created
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("supabase returned %d: %s", resp.StatusCode, string(body))
    }

    //* Если result не nil, десериализуем ответ в указанную структуру
    //NOTE: Это позволяет получить созданную запись (с ID, created_at и т.д.)
    if result != nil {
        return json.Unmarshal(body, result)
    }

    return nil
}

//! AuthUser — проверка Bearer-токена через Supabase Auth
// возвращает user_id

//? AuthUser проверяет валидность Bearer-токена через Supabase Auth API
//? и возвращает идентификатор пользователя.

func (c *Client) AuthUser(token string) (string, error) {
	// Параметры:
	//   - token: JWT-токен пользователя (Bearer token из Authorization header)
	//
	// Возвращает:
	//   - string: ID пользователя (UUID)
	//   - error: ошибка при проверке токена или получении данных пользователя

	//* создать запрос к Auth API для получения данных пользователя
	//* GET /auth/v1/user - эндпойнт для получения текцщего пользователя
	req, err := http.NewRequest("GET", c.baseURL+"/auth/v1/user", nil)
	if err != nil {
		return "", err
	}

	// Устанавливаем необходимые заголовки для аутинфекации
	// используем anon key как apikey а сам токен в authtorization
	req.Header.Set("apikey", c.anonKey)
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполняем запрос
	resp, err := c.do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //! по окончанию функции закрываем (для отсуствия траты ресурсов)


	// читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("невалидный токен")
	}

	var u struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &u); err != nil {
		return "", err
	}

	// проверяем что ID не пустой
	if u.ID == "" {
		return "", fmt.Errorf("не удалось получить user_id")
	}
	return u.ID, nil
}

// Delete — DELETE-запрос к PostgREST (удаление записей)
// TODO: написать комментарии к функции
func (c *Client) Delete(endpoint string, useServiceRole bool) error {
	url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	h["Prefer"] = "return=minimal"
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
