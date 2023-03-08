package gutil

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type SqlBuilder struct {
	BaseSql   string
	Temp      string
	Condition map[string]interface{}
	Total     int64
	pageSize  int
	pageNum   int64
	DB        *gorm.DB
}

// SqlSelect 创建sql合并
func SqlSelect(db *gorm.DB, s ...string) *SqlBuilder {
	tempStr := ""
	for _, item := range s {
		tempStr += item
	}
	return &SqlBuilder{
		BaseSql:   "",
		Temp:      tempStr,
		DB:        db,
		Condition: make(map[string]interface{}),
	}
}

func (s *SqlBuilder) Form(tableName string) *SqlBuilder {
	s.BaseSql = fmt.Sprintf("from %v ", tableName)
	return s
}

func (s *SqlBuilder) LeftJoin(v string) *SqlBuilder {
	s.BaseSql += fmt.Sprintf("left join %v ", v)
	return s
}

// Where 首次拼接
func (s *SqlBuilder) Where(w, k string, v interface{}) *SqlBuilder {
	if strings.Index(s.BaseSql, "where") == -1 {
		s.BaseSql = fmt.Sprintf("%v where %v", s.BaseSql, w)
		s.Condition[k] = v
		return s
	}
	s.BaseSql = fmt.Sprintf("%v and %v ", s.BaseSql, w)
	s.Condition[k] = v
	return s
}

func (s *SqlBuilder) Count() int64 {
	var count int64
	sql := fmt.Sprintf("select count(1) %v", s.BaseSql)
	if len(s.Condition) < 1 {
		err := s.DB.Raw(sql).Count(&count).Error
		if err != nil {
			fmt.Println(err.Error())
		}
		return count
	}
	err := s.DB.Raw(sql, s.Condition).Count(&count).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	return count
}

func (s *SqlBuilder) LimitNotPage(pageNum, pageSize int) *SqlBuilder {
	if pageNum == 0 {
		panic("页数不能为0")
	}
	pageNum--
	s.BaseSql += fmt.Sprintf(" limit %v,%v", pageNum*pageSize, pageSize)
	s.pageSize = pageSize
	return s
}

func (s *SqlBuilder) Limit(pageNum, pageSize int) *SqlBuilder {
	total := s.Count()
	if pageNum == 0 {
		panic("页数不能为0")
	}
	pageNum--
	totalPageNum := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPageNum++
	}
	s.BaseSql += fmt.Sprintf(" limit %v,%v", pageNum*pageSize, pageSize)
	s.pageNum = totalPageNum
	s.pageSize = pageSize
	s.Total = total
	return s
}

// ToSql 转换为sql
func (s *SqlBuilder) ToSql() (string, map[string]interface{}) {
	sql := fmt.Sprintf("select %v %v", s.Temp, s.BaseSql)
	return sql, s.Condition
}

// BuildNotPage 构建无需page信息 @returns 恐慌
func BuildNotPage[T any](s *SqlBuilder, v *T) error {
	sql := fmt.Sprintf("select %v %v", s.Temp, s.BaseSql)
	if len(s.Condition) < 1 {
		err := s.DB.Raw(sql).Find(&v).Error
		if err != nil {
			return err
		}
		return err
	}
	err := s.DB.Raw(sql, s.Condition).Find(&v).Error
	if err != nil {
		return err
	}
	return err
}

// Build 构建 r1 总页数 r2 总数据量 r3 恐慌
func Build[T any](s *SqlBuilder, v *T) (int64, int64, error) {
	sql := fmt.Sprintf("select %v %v", s.Temp, s.BaseSql)
	if len(s.Condition) < 1 {
		err := s.DB.Raw(sql).Find(&v).Error
		if err != nil {
			return 0, 0, err
		}
		return s.Total, s.pageNum, err
	}
	err := s.DB.Raw(sql, s.Condition).Find(&v).Error
	if err != nil {
		return 0, 0, err
	}
	return s.Total, s.pageNum, err
}
