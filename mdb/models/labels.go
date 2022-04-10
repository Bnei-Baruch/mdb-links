// Code generated by SQLBoiler 4.8.6 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package mdbmodels

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Label is an object representing the database table.
type Label struct {
	ID            int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
	UID           string    `boil:"uid" json:"uid" toml:"uid" yaml:"uid"`
	ContentUnitID int64     `boil:"content_unit_id" json:"content_unit_id" toml:"content_unit_id" yaml:"content_unit_id"`
	MediaType     string    `boil:"media_type" json:"media_type" toml:"media_type" yaml:"media_type"`
	Properties    null.JSON `boil:"properties" json:"properties,omitempty" toml:"properties" yaml:"properties,omitempty"`
	ApproveState  int16     `boil:"approve_state" json:"approve_state" toml:"approve_state" yaml:"approve_state"`
	CreatedAt     time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *labelR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L labelL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var LabelColumns = struct {
	ID            string
	UID           string
	ContentUnitID string
	MediaType     string
	Properties    string
	ApproveState  string
	CreatedAt     string
}{
	ID:            "id",
	UID:           "uid",
	ContentUnitID: "content_unit_id",
	MediaType:     "media_type",
	Properties:    "properties",
	ApproveState:  "approve_state",
	CreatedAt:     "created_at",
}

var LabelTableColumns = struct {
	ID            string
	UID           string
	ContentUnitID string
	MediaType     string
	Properties    string
	ApproveState  string
	CreatedAt     string
}{
	ID:            "labels.id",
	UID:           "labels.uid",
	ContentUnitID: "labels.content_unit_id",
	MediaType:     "labels.media_type",
	Properties:    "labels.properties",
	ApproveState:  "labels.approve_state",
	CreatedAt:     "labels.created_at",
}

// Generated where

var LabelWhere = struct {
	ID            whereHelperint64
	UID           whereHelperstring
	ContentUnitID whereHelperint64
	MediaType     whereHelperstring
	Properties    whereHelpernull_JSON
	ApproveState  whereHelperint16
	CreatedAt     whereHelpertime_Time
}{
	ID:            whereHelperint64{field: "\"labels\".\"id\""},
	UID:           whereHelperstring{field: "\"labels\".\"uid\""},
	ContentUnitID: whereHelperint64{field: "\"labels\".\"content_unit_id\""},
	MediaType:     whereHelperstring{field: "\"labels\".\"media_type\""},
	Properties:    whereHelpernull_JSON{field: "\"labels\".\"properties\""},
	ApproveState:  whereHelperint16{field: "\"labels\".\"approve_state\""},
	CreatedAt:     whereHelpertime_Time{field: "\"labels\".\"created_at\""},
}

// LabelRels is where relationship names are stored.
var LabelRels = struct {
	ContentUnit string
	LabelI18ns  string
	Tags        string
}{
	ContentUnit: "ContentUnit",
	LabelI18ns:  "LabelI18ns",
	Tags:        "Tags",
}

// labelR is where relationships are stored.
type labelR struct {
	ContentUnit *ContentUnit   `boil:"ContentUnit" json:"ContentUnit" toml:"ContentUnit" yaml:"ContentUnit"`
	LabelI18ns  LabelI18nSlice `boil:"LabelI18ns" json:"LabelI18ns" toml:"LabelI18ns" yaml:"LabelI18ns"`
	Tags        TagSlice       `boil:"Tags" json:"Tags" toml:"Tags" yaml:"Tags"`
}

// NewStruct creates a new relationship struct
func (*labelR) NewStruct() *labelR {
	return &labelR{}
}

// labelL is where Load methods for each relationship are stored.
type labelL struct{}

var (
	labelAllColumns            = []string{"id", "uid", "content_unit_id", "media_type", "properties", "approve_state", "created_at"}
	labelColumnsWithoutDefault = []string{"uid", "content_unit_id", "media_type"}
	labelColumnsWithDefault    = []string{"id", "properties", "approve_state", "created_at"}
	labelPrimaryKeyColumns     = []string{"id"}
	labelGeneratedColumns      = []string{}
)

type (
	// LabelSlice is an alias for a slice of pointers to Label.
	// This should almost always be used instead of []Label.
	LabelSlice []*Label

	labelQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	labelType                 = reflect.TypeOf(&Label{})
	labelMapping              = queries.MakeStructMapping(labelType)
	labelPrimaryKeyMapping, _ = queries.BindMapping(labelType, labelMapping, labelPrimaryKeyColumns)
	labelInsertCacheMut       sync.RWMutex
	labelInsertCache          = make(map[string]insertCache)
	labelUpdateCacheMut       sync.RWMutex
	labelUpdateCache          = make(map[string]updateCache)
	labelUpsertCacheMut       sync.RWMutex
	labelUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single label record from the query.
func (q labelQuery) One(exec boil.Executor) (*Label, error) {
	o := &Label{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for labels")
	}

	return o, nil
}

// All returns all Label records from the query.
func (q labelQuery) All(exec boil.Executor) (LabelSlice, error) {
	var o []*Label

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Label slice")
	}

	return o, nil
}

// Count returns the count of all Label records in the query.
func (q labelQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count labels rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q labelQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if labels exists")
	}

	return count > 0, nil
}

// ContentUnit pointed to by the foreign key.
func (o *Label) ContentUnit(mods ...qm.QueryMod) contentUnitQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ContentUnitID),
	}

	queryMods = append(queryMods, mods...)

	query := ContentUnits(queryMods...)
	queries.SetFrom(query.Query, "\"content_units\"")

	return query
}

// LabelI18ns retrieves all the label_i18n's LabelI18ns with an executor.
func (o *Label) LabelI18ns(mods ...qm.QueryMod) labelI18nQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"label_i18n\".\"label_id\"=?", o.ID),
	)

	query := LabelI18ns(queryMods...)
	queries.SetFrom(query.Query, "\"label_i18n\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"label_i18n\".*"})
	}

	return query
}

// Tags retrieves all the tag's Tags with an executor.
func (o *Label) Tags(mods ...qm.QueryMod) tagQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.InnerJoin("\"label_tag\" on \"tags\".\"id\" = \"label_tag\".\"tag_id\""),
		qm.Where("\"label_tag\".\"label_id\"=?", o.ID),
	)

	query := Tags(queryMods...)
	queries.SetFrom(query.Query, "\"tags\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"tags\".*"})
	}

	return query
}

// LoadContentUnit allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (labelL) LoadContentUnit(e boil.Executor, singular bool, maybeLabel interface{}, mods queries.Applicator) error {
	var slice []*Label
	var object *Label

	if singular {
		object = maybeLabel.(*Label)
	} else {
		slice = *maybeLabel.(*[]*Label)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &labelR{}
		}
		args = append(args, object.ContentUnitID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &labelR{}
			}

			for _, a := range args {
				if a == obj.ContentUnitID {
					continue Outer
				}
			}

			args = append(args, obj.ContentUnitID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`content_units`),
		qm.WhereIn(`content_units.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ContentUnit")
	}

	var resultSlice []*ContentUnit
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ContentUnit")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for content_units")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for content_units")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.ContentUnit = foreign
		if foreign.R == nil {
			foreign.R = &contentUnitR{}
		}
		foreign.R.Labels = append(foreign.R.Labels, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ContentUnitID == foreign.ID {
				local.R.ContentUnit = foreign
				if foreign.R == nil {
					foreign.R = &contentUnitR{}
				}
				foreign.R.Labels = append(foreign.R.Labels, local)
				break
			}
		}
	}

	return nil
}

// LoadLabelI18ns allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (labelL) LoadLabelI18ns(e boil.Executor, singular bool, maybeLabel interface{}, mods queries.Applicator) error {
	var slice []*Label
	var object *Label

	if singular {
		object = maybeLabel.(*Label)
	} else {
		slice = *maybeLabel.(*[]*Label)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &labelR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &labelR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`label_i18n`),
		qm.WhereIn(`label_i18n.label_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load label_i18n")
	}

	var resultSlice []*LabelI18n
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice label_i18n")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on label_i18n")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for label_i18n")
	}

	if singular {
		object.R.LabelI18ns = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &labelI18nR{}
			}
			foreign.R.Label = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.LabelID {
				local.R.LabelI18ns = append(local.R.LabelI18ns, foreign)
				if foreign.R == nil {
					foreign.R = &labelI18nR{}
				}
				foreign.R.Label = local
				break
			}
		}
	}

	return nil
}

// LoadTags allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (labelL) LoadTags(e boil.Executor, singular bool, maybeLabel interface{}, mods queries.Applicator) error {
	var slice []*Label
	var object *Label

	if singular {
		object = maybeLabel.(*Label)
	} else {
		slice = *maybeLabel.(*[]*Label)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &labelR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &labelR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.Select("\"tags\".id, \"tags\".description, \"tags\".parent_id, \"tags\".uid, \"tags\".pattern, \"a\".\"label_id\""),
		qm.From("\"tags\""),
		qm.InnerJoin("\"label_tag\" as \"a\" on \"tags\".\"id\" = \"a\".\"tag_id\""),
		qm.WhereIn("\"a\".\"label_id\" in ?", args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load tags")
	}

	var resultSlice []*Tag

	var localJoinCols []int64
	for results.Next() {
		one := new(Tag)
		var localJoinCol int64

		err = results.Scan(&one.ID, &one.Description, &one.ParentID, &one.UID, &one.Pattern, &localJoinCol)
		if err != nil {
			return errors.Wrap(err, "failed to scan eager loaded results for tags")
		}
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice tags")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on tags")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for tags")
	}

	if singular {
		object.R.Tags = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &tagR{}
			}
			foreign.R.Labels = append(foreign.R.Labels, object)
		}
		return nil
	}

	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			if local.ID == localJoinCol {
				local.R.Tags = append(local.R.Tags, foreign)
				if foreign.R == nil {
					foreign.R = &tagR{}
				}
				foreign.R.Labels = append(foreign.R.Labels, local)
				break
			}
		}
	}

	return nil
}

// SetContentUnit of the label to the related item.
// Sets o.R.ContentUnit to related.
// Adds o to related.R.Labels.
func (o *Label) SetContentUnit(exec boil.Executor, insert bool, related *ContentUnit) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"labels\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"content_unit_id"}),
		strmangle.WhereClause("\"", "\"", 2, labelPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ContentUnitID = related.ID
	if o.R == nil {
		o.R = &labelR{
			ContentUnit: related,
		}
	} else {
		o.R.ContentUnit = related
	}

	if related.R == nil {
		related.R = &contentUnitR{
			Labels: LabelSlice{o},
		}
	} else {
		related.R.Labels = append(related.R.Labels, o)
	}

	return nil
}

// AddLabelI18ns adds the given related objects to the existing relationships
// of the label, optionally inserting them as new records.
// Appends related to o.R.LabelI18ns.
// Sets related.R.Label appropriately.
func (o *Label) AddLabelI18ns(exec boil.Executor, insert bool, related ...*LabelI18n) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.LabelID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"label_i18n\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"label_id"}),
				strmangle.WhereClause("\"", "\"", 2, labelI18nPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.LabelID, rel.Language}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.LabelID = o.ID
		}
	}

	if o.R == nil {
		o.R = &labelR{
			LabelI18ns: related,
		}
	} else {
		o.R.LabelI18ns = append(o.R.LabelI18ns, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &labelI18nR{
				Label: o,
			}
		} else {
			rel.R.Label = o
		}
	}
	return nil
}

// AddTags adds the given related objects to the existing relationships
// of the label, optionally inserting them as new records.
// Appends related to o.R.Tags.
// Sets related.R.Labels appropriately.
func (o *Label) AddTags(exec boil.Executor, insert bool, related ...*Tag) error {
	var err error
	for _, rel := range related {
		if insert {
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}
	}

	for _, rel := range related {
		query := "insert into \"label_tag\" (\"label_id\", \"tag_id\") values ($1, $2)"
		values := []interface{}{o.ID, rel.ID}

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, query)
			fmt.Fprintln(boil.DebugWriter, values)
		}
		_, err = exec.Exec(query, values...)
		if err != nil {
			return errors.Wrap(err, "failed to insert into join table")
		}
	}
	if o.R == nil {
		o.R = &labelR{
			Tags: related,
		}
	} else {
		o.R.Tags = append(o.R.Tags, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &tagR{
				Labels: LabelSlice{o},
			}
		} else {
			rel.R.Labels = append(rel.R.Labels, o)
		}
	}
	return nil
}

// SetTags removes all previously related items of the
// label replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Labels's Tags accordingly.
// Replaces o.R.Tags with related.
// Sets related.R.Labels's Tags accordingly.
func (o *Label) SetTags(exec boil.Executor, insert bool, related ...*Tag) error {
	query := "delete from \"label_tag\" where \"label_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	removeTagsFromLabelsSlice(o, related)
	if o.R != nil {
		o.R.Tags = nil
	}
	return o.AddTags(exec, insert, related...)
}

// RemoveTags relationships from objects passed in.
// Removes related items from R.Tags (uses pointer comparison, removal does not keep order)
// Sets related.R.Labels.
func (o *Label) RemoveTags(exec boil.Executor, related ...*Tag) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	query := fmt.Sprintf(
		"delete from \"label_tag\" where \"label_id\" = $1 and \"tag_id\" in (%s)",
		strmangle.Placeholders(dialect.UseIndexPlaceholders, len(related), 2, 1),
	)
	values := []interface{}{o.ID}
	for _, rel := range related {
		values = append(values, rel.ID)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err = exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}
	removeTagsFromLabelsSlice(o, related)
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Tags {
			if rel != ri {
				continue
			}

			ln := len(o.R.Tags)
			if ln > 1 && i < ln-1 {
				o.R.Tags[i] = o.R.Tags[ln-1]
			}
			o.R.Tags = o.R.Tags[:ln-1]
			break
		}
	}

	return nil
}

func removeTagsFromLabelsSlice(o *Label, related []*Tag) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.Labels {
			if o.ID != ri.ID {
				continue
			}

			ln := len(rel.R.Labels)
			if ln > 1 && i < ln-1 {
				rel.R.Labels[i] = rel.R.Labels[ln-1]
			}
			rel.R.Labels = rel.R.Labels[:ln-1]
			break
		}
	}
}

// Labels retrieves all the records using an executor.
func Labels(mods ...qm.QueryMod) labelQuery {
	mods = append(mods, qm.From("\"labels\""))
	return labelQuery{NewQuery(mods...)}
}

// FindLabel retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindLabel(exec boil.Executor, iD int64, selectCols ...string) (*Label, error) {
	labelObj := &Label{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"labels\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, labelObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from labels")
	}

	return labelObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Label) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no labels provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(labelColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	labelInsertCacheMut.RLock()
	cache, cached := labelInsertCache[key]
	labelInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			labelAllColumns,
			labelColumnsWithDefault,
			labelColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(labelType, labelMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(labelType, labelMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"labels\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"labels\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into labels")
	}

	if !cached {
		labelInsertCacheMut.Lock()
		labelInsertCache[key] = cache
		labelInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Label.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Label) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	labelUpdateCacheMut.RLock()
	cache, cached := labelUpdateCache[key]
	labelUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			labelAllColumns,
			labelPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update labels, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"labels\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, labelPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(labelType, labelMapping, append(wl, labelPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	var result sql.Result
	result, err = exec.Exec(cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update labels row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for labels")
	}

	if !cached {
		labelUpdateCacheMut.Lock()
		labelUpdateCache[key] = cache
		labelUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q labelQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for labels")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for labels")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o LabelSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), labelPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"labels\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, labelPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in label slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all label")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Label) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no labels provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(labelColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	labelUpsertCacheMut.RLock()
	cache, cached := labelUpsertCache[key]
	labelUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			labelAllColumns,
			labelColumnsWithDefault,
			labelColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			labelAllColumns,
			labelPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert labels, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(labelPrimaryKeyColumns))
			copy(conflict, labelPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"labels\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(labelType, labelMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(labelType, labelMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert labels")
	}

	if !cached {
		labelUpsertCacheMut.Lock()
		labelUpsertCache[key] = cache
		labelUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Label record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Label) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Label provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), labelPrimaryKeyMapping)
	sql := "DELETE FROM \"labels\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from labels")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for labels")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q labelQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no labelQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from labels")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for labels")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o LabelSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), labelPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"labels\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, labelPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from label slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for labels")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Label) Reload(exec boil.Executor) error {
	ret, err := FindLabel(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *LabelSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := LabelSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), labelPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"labels\".* FROM \"labels\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, labelPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in LabelSlice")
	}

	*o = slice

	return nil
}

// LabelExists checks if the Label row exists.
func LabelExists(exec boil.Executor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"labels\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if labels exists")
	}

	return exists, nil
}
