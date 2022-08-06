// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"notifications/ent/notification"
	"notifications/ent/schema"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// NotificationCreate is the builder for creating a Notification entity.
type NotificationCreate struct {
	config
	mutation *NotificationMutation
	hooks    []Hook
}

// SetSenderID sets the "sender_id" field.
func (nc *NotificationCreate) SetSenderID(i int) *NotificationCreate {
	nc.mutation.SetSenderID(i)
	return nc
}

// SetType sets the "type" field.
func (nc *NotificationCreate) SetType(st schema.NotificationType) *NotificationCreate {
	nc.mutation.SetType(st)
	return nc
}

// SetNillableType sets the "type" field if the given value is not nil.
func (nc *NotificationCreate) SetNillableType(st *schema.NotificationType) *NotificationCreate {
	if st != nil {
		nc.SetType(*st)
	}
	return nc
}

// SetPayload sets the "payload" field.
func (nc *NotificationCreate) SetPayload(s schema.Payload) *NotificationCreate {
	nc.mutation.SetPayload(s)
	return nc
}

// SetTTL sets the "ttl" field.
func (nc *NotificationCreate) SetTTL(i int) *NotificationCreate {
	nc.mutation.SetTTL(i)
	return nc
}

// SetStatus sets the "status" field.
func (nc *NotificationCreate) SetStatus(ss schema.NotificationStatus) *NotificationCreate {
	nc.mutation.SetStatus(ss)
	return nc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (nc *NotificationCreate) SetNillableStatus(ss *schema.NotificationStatus) *NotificationCreate {
	if ss != nil {
		nc.SetStatus(*ss)
	}
	return nc
}

// SetCreatedAt sets the "created_at" field.
func (nc *NotificationCreate) SetCreatedAt(t time.Time) *NotificationCreate {
	nc.mutation.SetCreatedAt(t)
	return nc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (nc *NotificationCreate) SetNillableCreatedAt(t *time.Time) *NotificationCreate {
	if t != nil {
		nc.SetCreatedAt(*t)
	}
	return nc
}

// SetUpdatedAt sets the "updated_at" field.
func (nc *NotificationCreate) SetUpdatedAt(t time.Time) *NotificationCreate {
	nc.mutation.SetUpdatedAt(t)
	return nc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (nc *NotificationCreate) SetNillableUpdatedAt(t *time.Time) *NotificationCreate {
	if t != nil {
		nc.SetUpdatedAt(*t)
	}
	return nc
}

// SetPlannedAt sets the "planned_at" field.
func (nc *NotificationCreate) SetPlannedAt(t time.Time) *NotificationCreate {
	nc.mutation.SetPlannedAt(t)
	return nc
}

// SetNillablePlannedAt sets the "planned_at" field if the given value is not nil.
func (nc *NotificationCreate) SetNillablePlannedAt(t *time.Time) *NotificationCreate {
	if t != nil {
		nc.SetPlannedAt(*t)
	}
	return nc
}

// SetRetries sets the "retries" field.
func (nc *NotificationCreate) SetRetries(i int) *NotificationCreate {
	nc.mutation.SetRetries(i)
	return nc
}

// SetNillableRetries sets the "retries" field if the given value is not nil.
func (nc *NotificationCreate) SetNillableRetries(i *int) *NotificationCreate {
	if i != nil {
		nc.SetRetries(*i)
	}
	return nc
}

// SetSentAt sets the "sent_at" field.
func (nc *NotificationCreate) SetSentAt(t time.Time) *NotificationCreate {
	nc.mutation.SetSentAt(t)
	return nc
}

// SetNillableSentAt sets the "sent_at" field if the given value is not nil.
func (nc *NotificationCreate) SetNillableSentAt(t *time.Time) *NotificationCreate {
	if t != nil {
		nc.SetSentAt(*t)
	}
	return nc
}

// Mutation returns the NotificationMutation object of the builder.
func (nc *NotificationCreate) Mutation() *NotificationMutation {
	return nc.mutation
}

// Save creates the Notification in the database.
func (nc *NotificationCreate) Save(ctx context.Context) (*Notification, error) {
	var (
		err  error
		node *Notification
	)
	nc.defaults()
	if len(nc.hooks) == 0 {
		if err = nc.check(); err != nil {
			return nil, err
		}
		node, err = nc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*NotificationMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = nc.check(); err != nil {
				return nil, err
			}
			nc.mutation = mutation
			if node, err = nc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(nc.hooks) - 1; i >= 0; i-- {
			if nc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = nc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, nc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Notification)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from NotificationMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (nc *NotificationCreate) SaveX(ctx context.Context) *Notification {
	v, err := nc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (nc *NotificationCreate) Exec(ctx context.Context) error {
	_, err := nc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nc *NotificationCreate) ExecX(ctx context.Context) {
	if err := nc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (nc *NotificationCreate) defaults() {
	if _, ok := nc.mutation.GetType(); !ok {
		v := notification.DefaultType
		nc.mutation.SetType(v)
	}
	if _, ok := nc.mutation.Status(); !ok {
		v := notification.DefaultStatus
		nc.mutation.SetStatus(v)
	}
	if _, ok := nc.mutation.CreatedAt(); !ok {
		v := notification.DefaultCreatedAt()
		nc.mutation.SetCreatedAt(v)
	}
	if _, ok := nc.mutation.UpdatedAt(); !ok {
		v := notification.DefaultUpdatedAt()
		nc.mutation.SetUpdatedAt(v)
	}
	if _, ok := nc.mutation.PlannedAt(); !ok {
		v := notification.DefaultPlannedAt()
		nc.mutation.SetPlannedAt(v)
	}
	if _, ok := nc.mutation.Retries(); !ok {
		v := notification.DefaultRetries
		nc.mutation.SetRetries(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nc *NotificationCreate) check() error {
	if _, ok := nc.mutation.SenderID(); !ok {
		return &ValidationError{Name: "sender_id", err: errors.New(`ent: missing required field "Notification.sender_id"`)}
	}
	if _, ok := nc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "Notification.type"`)}
	}
	if v, ok := nc.mutation.GetType(); ok {
		if err := notification.TypeValidator(string(v)); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Notification.type": %w`, err)}
		}
	}
	if _, ok := nc.mutation.Payload(); !ok {
		return &ValidationError{Name: "payload", err: errors.New(`ent: missing required field "Notification.payload"`)}
	}
	if _, ok := nc.mutation.TTL(); !ok {
		return &ValidationError{Name: "ttl", err: errors.New(`ent: missing required field "Notification.ttl"`)}
	}
	if _, ok := nc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "Notification.status"`)}
	}
	if v, ok := nc.mutation.Status(); ok {
		if err := notification.StatusValidator(string(v)); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Notification.status": %w`, err)}
		}
	}
	if _, ok := nc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Notification.created_at"`)}
	}
	if _, ok := nc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Notification.updated_at"`)}
	}
	if _, ok := nc.mutation.PlannedAt(); !ok {
		return &ValidationError{Name: "planned_at", err: errors.New(`ent: missing required field "Notification.planned_at"`)}
	}
	if _, ok := nc.mutation.Retries(); !ok {
		return &ValidationError{Name: "retries", err: errors.New(`ent: missing required field "Notification.retries"`)}
	}
	return nil
}

func (nc *NotificationCreate) sqlSave(ctx context.Context) (*Notification, error) {
	_node, _spec := nc.createSpec()
	if err := sqlgraph.CreateNode(ctx, nc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (nc *NotificationCreate) createSpec() (*Notification, *sqlgraph.CreateSpec) {
	var (
		_node = &Notification{config: nc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: notification.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: notification.FieldID,
			},
		}
	)
	if value, ok := nc.mutation.SenderID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: notification.FieldSenderID,
		})
		_node.SenderID = value
	}
	if value, ok := nc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: notification.FieldType,
		})
		_node.Type = value
	}
	if value, ok := nc.mutation.Payload(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: notification.FieldPayload,
		})
		_node.Payload = value
	}
	if value, ok := nc.mutation.TTL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: notification.FieldTTL,
		})
		_node.TTL = value
	}
	if value, ok := nc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: notification.FieldStatus,
		})
		_node.Status = value
	}
	if value, ok := nc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: notification.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := nc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: notification.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := nc.mutation.PlannedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: notification.FieldPlannedAt,
		})
		_node.PlannedAt = value
	}
	if value, ok := nc.mutation.Retries(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: notification.FieldRetries,
		})
		_node.Retries = value
	}
	if value, ok := nc.mutation.SentAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: notification.FieldSentAt,
		})
		_node.SentAt = &value
	}
	return _node, _spec
}

// NotificationCreateBulk is the builder for creating many Notification entities in bulk.
type NotificationCreateBulk struct {
	config
	builders []*NotificationCreate
}

// Save creates the Notification entities in the database.
func (ncb *NotificationCreateBulk) Save(ctx context.Context) ([]*Notification, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ncb.builders))
	nodes := make([]*Notification, len(ncb.builders))
	mutators := make([]Mutator, len(ncb.builders))
	for i := range ncb.builders {
		func(i int, root context.Context) {
			builder := ncb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*NotificationMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ncb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ncb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ncb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ncb *NotificationCreateBulk) SaveX(ctx context.Context) []*Notification {
	v, err := ncb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ncb *NotificationCreateBulk) Exec(ctx context.Context) error {
	_, err := ncb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ncb *NotificationCreateBulk) ExecX(ctx context.Context) {
	if err := ncb.Exec(ctx); err != nil {
		panic(err)
	}
}
