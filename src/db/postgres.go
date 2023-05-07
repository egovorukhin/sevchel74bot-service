package db

import (
	"encoding/base64"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
	"time"
)

type Config struct {
	Server     Server     `yaml:"server"`
	Username   string     `yaml:"username"`
	Password   string     `yaml:"password"`
	Name       string     `yaml:"name"`
	SSL        bool       `yaml:"ssl"`
	ConnConfig ConnConfig `yaml:"connConfig"`
}

type Server struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

type ConnConfig struct {
	MaxIdle     int `yaml:"maxIdle"`
	MaxOpen     int `yaml:"maxOpen"`
	MaxLifetime int `yaml:"maxLifetime"`
}

type Table struct {
	Name    string
	db      *gorm.DB
	selects []string
	orders  []string
	limit   int
	offset  int
}

var db *gorm.DB

const InformationSchemaColumns = "information_schema.columns"

func DB(table string) *Table {
	return &Table{
		Name:   table,
		db:     db,
		limit:  -1,
		offset: -1,
	}
}

func Init(config Config) error {

	password, err := base64.StdEncoding.DecodeString(config.Password)
	if err != nil {
		return err
	}

	dataSource := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		config.Server.Hostname,
		config.Username,
		password,
		config.Name,
		config.Server.Port,
	)
	if !config.SSL {
		dataSource += " sslmode=disable"
	}

	db, err = gorm.Open(postgres.Open(dataSource), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(config.ConnConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(config.ConnConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnConfig.MaxLifetime) * time.Minute)

	return nil
}

func Close() error {
	own, err := db.DB()
	if err != nil {
		return err
	}
	return own.Close()
}

// OrderBy Указываются поля для сортировки
func (t *Table) OrderBy(orders ...string) *Table {
	t.orders = append(t.orders, orders...)
	return t
}

// Select Выбор отображаемых полей при загрузке данных
func (t *Table) Select(selects ...string) *Table {
	t.selects = append(t.selects, selects...)
	return t
}

// Limit Установка лимита возвращаемых записей
func (t *Table) Limit(l int) *Table {
	t.limit = l
	return t
}

// Offset Установка строки с какой возвращать записи
func (t *Table) Offset(o int) *Table {
	t.offset = o
	return t
}

// GetRecord Возвращаем одну запись
func (t *Table) GetRecord(v interface{}, query interface{}, args ...interface{}) error {
	return t.db.Table(t.Name).Select(t.selects).Where(query, args...).Order(strings.Join(t.orders, ",")).Limit(t.limit).Offset(t.offset).Take(v).Error
}

// GetMapRecord Возвращаем одну запись
func (t *Table) GetMapRecord(query interface{}, args ...interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := t.db.Table(t.Name).Select(t.selects).Where(query, args...).Order(strings.Join(t.orders, ",")).Limit(t.limit).Offset(t.offset).Take(&m).Error
	if err == gorm.ErrRecordNotFound {
		m = nil
	}
	return m, err
}

// GetRecords Возвращаем всю таблицу
func (t *Table) GetRecords(v interface{}, query interface{}, args ...interface{}) error {
	return t.db.Table(t.Name).Select(t.selects).Where(query, args...).Order(strings.Join(t.orders, ",")).Limit(t.limit).Offset(t.offset).Find(v).Error
}

// GetMapRecords Возвращаем всю таблицу
func (t *Table) GetMapRecords(query interface{}, args ...interface{}) ([]map[string]interface{}, error) {
	var m []map[string]interface{}
	err := t.db.Table(t.Name).Select(t.selects).Where(query, args...).Order(strings.Join(t.orders, ",")).Limit(t.limit).Offset(t.offset).Find(&m).Error
	if err == gorm.ErrRecordNotFound {
		m = nil
	}
	return m, err
}

// Create Создать запись
func (t *Table) Create(v interface{}) error {
	now := time.Now().UTC()
	setFieldValue(v, "Created", &now)
	setFieldValue(v, "Modified", &now)
	return t.db.Table(t.Name).Create(v).Error
}

// CreateAndSelect Создать и вернуть последнюю запись
func (t *Table) CreateAndSelect(input, output interface{}) error {
	err := t.Create(input)
	if err != nil {
		return err
	}
	return t.OrderBy("id desc").GetRecord("", output)
}

// CreateAndGetRecord Создать и вернуть последнюю запись
func (t *Table) CreateAndGetRecord(input interface{}, orders ...string) (map[string]interface{}, error) {
	err := t.Create(input)
	if err != nil {
		return nil, err
	}
	return t.OrderBy(orders...).GetMapRecord("")
}

// Update Обновить запись, можно указать Select для полей, которые хотим обновить, * - все поля
func (t *Table) Update(v interface{}, query interface{}, args ...interface{}) error {
	now := time.Now().UTC()
	setFieldValue(v, "Modified", &now)
	return t.db.Table(t.Name).Where(query, args...).Updates(v).Error
}

// UpdateAndSelect Обновить и вернуть запись по критериям обновления
func (t *Table) UpdateAndSelect(input, output interface{}, query interface{}, args ...interface{}) error {
	err := t.Update(input, query, args...)
	if err != nil {
		return err
	}
	return t.GetRecord(output, query, args...)
}

// UpdateAndGetRecord Обновить и вернуть запись по критериям обновления
func (t *Table) UpdateAndGetRecord(input interface{}, where string, orders ...string) (map[string]interface{}, error) {
	err := t.Update(input, where)
	if err != nil {
		return nil, err
	}
	return t.OrderBy(orders...).GetMapRecord(where)
}

// Delete Удалить запись
func (t *Table) Delete(v interface{}, query interface{}, args ...interface{}) error {
	return t.db.Table(t.Name).Where(query, args...).Delete(v).Error
}

// LastId Возвращаем id последней записи
func (t *Table) LastId(fieldName string, query interface{}, args ...interface{}) interface{} {

	result := make(map[string]interface{})
	t.db.Table(t.Name).Where(query, args...).Limit(1).Order(clause.OrderByColumn{
		Column: clause.Column{Table: clause.CurrentTable, Name: fieldName},
		Desc:   true,
	}).Find(&result)

	if id, ok := result[fieldName]; ok {
		return id
	}

	return nil
}

// UpdateOrCreate Обновить или создать запись на основе проверки наличия записи
func (t *Table) UpdateOrCreate(v interface{}, query interface{}, args ...interface{} /*, selects ...interface{}*/) error {
	now := time.Now()
	setFieldValue(v, "Modified", &now)
	tx := t.db.Table(t.Name).Where(query, args...).Updates(v)
	if tx.RowsAffected == 0 {
		return t.Create(v)
	}
	return tx.Error
}

func (t *Table) Columns(fields ...string) ([]map[string]interface{}, error) {
	tableName := t.Name
	t.Name = InformationSchemaColumns
	dotIndex := strings.Index(tableName, ".")
	schema := tableName[:dotIndex]
	tableName = tableName[dotIndex+1:]
	return t.Select(fields...).GetMapRecords(fmt.Sprintf("table_schema='%s' and table_name='%s'", schema, tableName))
}

// Установка значения в поле структуры
func setFieldValue(i interface{}, fieldName string, value interface{}) {

	if i == nil {
		return
	}

	v := reflect.Indirect(reflect.ValueOf(i))
	switch v.Kind() {
	case reflect.Map, reflect.Slice:
		return
	}
	/*if v.Kind() == reflect.Map {
		v.SetMapIndex(reflect.ValueOf(strings.ToLower(fieldName)), reflect.ValueOf(value))
		return
	}*/
	f := v.FieldByName(fieldName)
	if f.IsValid() {
		if f.CanSet() {
			switch f.Kind() {
			case reflect.Struct:
				f.Set(reflect.ValueOf(value))
				return
			case reflect.Ptr:
				if f.IsNil() {
					f.Set(reflect.ValueOf(value))
				}
				return
			}
		}
	}
}
