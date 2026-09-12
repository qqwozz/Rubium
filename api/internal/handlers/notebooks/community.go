package notebooks

import (
	"fmt"
	"net/http"
	"strings"
	"encoding/json"

	"api/internal/validation"

	"github.com/gin-gonic/gin"
)

func (h *NotebooksHandler) GetCommunityNotebooks(c *gin.Context) {
	sortVal := c.DefaultQuery("sort", "rating")
	if _, ok := validation.ValidSort(sortVal); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный sort"})
		return
	}

	search := c.Query("search")
	if search != "" {
		if _, ok := validation.SafeSearchString(search); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный поисковый запрос"})
			return
		}
	}

	limitVal := c.DefaultQuery("limit", "50")
	limit, ok := validation.ValidLimit(limitVal)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit должен быть от 1 до 100"})
		return
	}

	selectFields := "id,user_id,title,description,color,tags,is_public,average_rating,ratings_count,views_count,copies_count,created_at,updated_at"

	filters := []string{
		"select=" + selectFields,
		"is_public=eq.true",
	}

	if search != "" {
		safeSearch, _ := validation.SafeSearchString(search)
		filters = append(filters, fmt.Sprintf("or=(title.ilike.*%s*,description.ilike.*%s*)", safeSearch, safeSearch))
	}

	switch sortVal {
	case "newest":
		filters = append(filters, "order=created_at.desc")
	case "popular":
		filters = append(filters, "order=views_count.desc")
	default:
		filters = append(filters, "order=average_rating.desc")
	}

	filters = append(filters, fmt.Sprintf("limit=%d", limit))

	endpoint := "notebooks?" + strings.Join(filters, "&")

	rawNotebooks, err := h.client.RawQuery(c.Request.Context(), endpoint, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(rawNotebooks) > 0 && rawNotebooks[0] == '{' {
		rawNotebooks = append([]byte("["), append(rawNotebooks, []byte("]")...)...)
	}

	var rows []notebookRow
	if err := json.Unmarshal(rawNotebooks, &rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userIDs := make([]string, 0, len(rows))
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
	}

	authors := make(map[string]gin.H)
	if len(userIDs) > 0 {
		usersEndpoint := fmt.Sprintf("rubium_users?select=id,first_name,last_name,email,avatar_url&id=in.(%s)", strings.Join(userIDs, ","))
		rawUsers, err := h.client.RawQuery(c.Request.Context(), usersEndpoint, false)
		if err == nil {
			if len(rawUsers) > 0 && rawUsers[0] == '{' {
				rawUsers = append([]byte("["), append(rawUsers, []byte("]")...)...)
			}
			var users []struct {
				ID        string `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Email     string `json:"email"`
				AvatarURL string `json:"avatar_url"`
			}
			if json.Unmarshal(rawUsers, &users) == nil {
				for _, u := range users {
					authors[u.ID] = gin.H{
						"id":         u.ID,
						"first_name": u.FirstName,
						"last_name":  u.LastName,
						"email":      u.Email,
						"avatar_url": u.AvatarURL,
					}
				}
			}
		}
	}

	notebooks := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		resp := r.toResponse(false)
		if author, ok := authors[r.UserID]; ok {
			resp["author"] = author
		} else {
			resp["author"] = gin.H{
				"id":         "",
				"first_name": "Автор",
				"email":      "",
				"avatar_url": "",
			}
		}
		notebooks = append(notebooks, resp)
	}

	c.JSON(http.StatusOK, gin.H{"notebooks": notebooks})
}