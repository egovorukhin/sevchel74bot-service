package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sevchel74botService/src/db"
	"strconv"
	"time"
)

type Model struct {
	Created  *time.Time `json:"created"`
	Modified *time.Time `json:"modified"`
	Enabled  bool       `json:"enabled"`
	Author   string     `json:"author"`
}

type Filter struct {
	Fields []string `json:"fields,omitempty"`
	Orders []string `json:"orders,omitempty"`
	Limit  int      `json:"limit,omitempty"`
	Offset int      `json:"offset,omitempty"`
}

type GroupMember struct {
	GroupName, MemberName string
	m                     map[string]int
}

type CrudResponse struct {
	Id interface{} `json:"id" jsonschema:"type=integer"`
	General
}

type GroupMemberResponse struct {
	GroupId  int `json:"group_id"`
	MemberId int `json:"member_id"`
	General
}

type General struct {
	State    string    `json:"state"`
	Datetime time.Time `json:"datetime"`
	Message  string    `json:"message"`
}

const (
	Read    = "Read"
	Created = "Created"
	Updated = "Updated"
	Deleted = "Deleted"

	Added   = "Added"
	Removed = "Removed"

	Failed = "Failed"
)

const (
	ModelCrudResponse        = "CrudResponse"
	ModelFilter              = "Filter"
	ModelGroupMember         = "GroupMember"
	ModelGroupMemberResponse = "GroupMemberResponse"
)

func GetValue(v interface{}, fieldName string) string {
	field := reflect.Indirect(reflect.ValueOf(v)).FieldByName(fieldName)
	if field.IsValid() {
		switch field.Kind() {
		case reflect.String:
			return field.Interface().(string)
		}
	}
	return ""
}

// NewCrudResponse Ответ для обычных таблиц
func NewCrudResponse(id interface{}, state, message string) CrudResponse {
	return CrudResponse{
		Id: id,
		General: General{
			State:    state,
			Datetime: time.Now(),
			Message:  message,
		},
	}
}

func NewCreated(id interface{}, message string) CrudResponse {
	return NewCrudResponse(id, Created, message)
}

func NewUpdated(id interface{}, message string) CrudResponse {
	return NewCrudResponse(id, Updated, message)
}

func NewDeleted(id interface{}, message string) CrudResponse {
	return NewCrudResponse(id, Deleted, message)
}

func NewCrudFailed(id, memberId interface{}, message string) CrudResponse {
	return NewCrudResponse(id, Failed, message)
}

// NewGroupMemberResponse Ответ для групповых таблиц
func NewGroupMemberResponse(groupId, memberId int, state, message string) GroupMemberResponse {
	return GroupMemberResponse{
		GroupId:  groupId,
		MemberId: memberId,
		General: General{
			State:    state,
			Datetime: time.Now(),
			Message:  message,
		},
	}
}

func NewAdded(groupId, memberId int, message string) GroupMemberResponse {
	return NewGroupMemberResponse(groupId, memberId, Added, message)
}

func NewRemoved(groupId, memberId int, message string) GroupMemberResponse {
	return NewGroupMemberResponse(groupId, memberId, Removed, message)
}

func NewGMFailed(groupId, memberId int, message string) GroupMemberResponse {
	return NewGroupMemberResponse(groupId, memberId, Failed, message)
}

// NewFilter Инициализация фильтра для таблиц
func NewFilter(data []byte) (f Filter, err error) {
	if len(data) > 0 {
		err = json.Unmarshal(data, &f)
	}
	return f, err
}

func (f Filter) GetRecord(tableName, query string) (map[string]interface{}, error) {
	return db.DB(tableName).Select(f.Fields...).OrderBy(f.Orders...).Limit(f.Limit).Offset(f.Offset).GetMapRecord(query)
}

func (f Filter) GetRecords(tableName, query string) ([]map[string]interface{}, error) {
	return db.DB(tableName).Select(f.Fields...).OrderBy(f.Orders...).Limit(f.Limit).Offset(f.Offset).GetMapRecords(query)
}

func NewGroupMember(groupName, memberName string) GroupMember {
	return GroupMember{
		GroupName:  groupName,
		MemberName: memberName,
		m:          map[string]int{},
	}
}

func (gm GroupMember) setValue(fieldName string, value interface{}) error {
	var s string
	switch value.(type) {
	case string:
		s = value.(string)
	case int, int8, int16, int32, int64:
		s = fmt.Sprintf("%d", value)
	default:
		return errors.New("Data type not valid. Use: string, int, int8, int16, int32, int64")
	}

	id, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	gm.m[fieldName] = id

	return nil
}

func (gm GroupMember) SetGroup(value interface{}) error {
	return gm.setValue(gm.GroupName, value)
}

func (gm GroupMember) SetMember(value interface{}) error {
	return gm.setValue(gm.MemberName, value)
}

func (gm GroupMember) Valid() bool {
	_, isGroup := gm.m[gm.GroupName]
	_, isMember := gm.m[gm.MemberName]
	if !isGroup || !isMember {
		return false
	}
	return true
}

func (gm GroupMember) Group() int {
	return gm.value(gm.GroupName)
}

func (gm GroupMember) Member() int {
	return gm.value(gm.MemberName)
}

func (gm GroupMember) value(name string) int {
	if value, ok := gm.m[name]; ok {
		return value
	}
	return -1
}

func (gm GroupMember) GroupString() string {
	return strconv.Itoa(gm.Group())
}

func (gm GroupMember) MemberString() string {
	return strconv.Itoa(gm.Member())
}

func (gm GroupMember) GroupNameValue() (string, int) {
	return gm.GroupName, gm.Group()
}

func (gm GroupMember) MemberNameValue() (string, int) {
	return gm.MemberName, gm.Member()
}

func nameValueString(name string, value int) string {
	return fmt.Sprintf("%s=%d", name, value)
}

func (gm GroupMember) GroupNameValueString() string {
	return nameValueString(gm.GroupName, gm.Group())
}

func (gm GroupMember) MemberNameValueString() string {
	return nameValueString(gm.MemberName, gm.Member())
}

func order(fieldName string, desc ...bool) string {
	if desc != nil {
		return fieldName + " desc"
	}
	return fieldName
}

func (gm GroupMember) GroupOrder(desc ...bool) string {
	return order(gm.GroupName, desc...)
}

func (gm GroupMember) MemberOrder(desc ...bool) string {
	return order(gm.MemberName, desc...)
}

func (gm *GroupMember) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &gm.m)
	if err != nil {
		return err
	}
	if !gm.Valid() {
		return errors.New(fmt.Sprintf("Required fields [%s, %s]", gm.GroupName, gm.MemberName))
	}
	return nil
}
