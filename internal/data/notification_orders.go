package data

import "notifications/ent"

func OrderByCreatedAt() ent.OrderFunc {
	return ent.Asc(`created_at`)
}
