package contract

import "time"

type CheckFavoriteResponse struct {
	UserId     string `json:"userId"`
	ItemId     string `json:"itemId"`
	IsFavorite bool   `json:"isFavorite"`
}

type ListFavoriteResponseItem struct {
	UserId    string    `json:"userId"`
	ItemId    string    `json:"itemId"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListFavoriteResponse struct {
	Items     []ListFavoriteResponseItem `json:"items"`
	NextToken string                     `json:"nextToken"`
}
