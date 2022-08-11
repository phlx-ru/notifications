// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// NotificationsColumns holds the columns for the "notifications" table.
	NotificationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "sender_id", Type: field.TypeInt},
		{Name: "type", Type: field.TypeString, Default: "email"},
		{Name: "payload", Type: field.TypeJSON},
		{Name: "ttl", Type: field.TypeInt},
		{Name: "status", Type: field.TypeString, Default: "draft"},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "planned_at", Type: field.TypeTime},
		{Name: "retry_at", Type: field.TypeTime, Nullable: true},
		{Name: "retries", Type: field.TypeInt, Default: 0},
		{Name: "sent_at", Type: field.TypeTime, Nullable: true},
	}
	// NotificationsTable holds the schema information for the "notifications" table.
	NotificationsTable = &schema.Table{
		Name:       "notifications",
		Columns:    NotificationsColumns,
		PrimaryKey: []*schema.Column{NotificationsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "notification_type",
				Unique:  false,
				Columns: []*schema.Column{NotificationsColumns[2]},
			},
			{
				Name:    "notification_status",
				Unique:  false,
				Columns: []*schema.Column{NotificationsColumns[5]},
			},
			{
				Name:    "notification_planned_at",
				Unique:  false,
				Columns: []*schema.Column{NotificationsColumns[8]},
			},
			{
				Name:    "notification_sent_at",
				Unique:  false,
				Columns: []*schema.Column{NotificationsColumns[11]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		NotificationsTable,
	}
)

func init() {
}
