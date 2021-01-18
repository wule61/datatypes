package datatypes

import (
	"context"
	"fmt"
	"github.com/twpayne/go-geom/encoding/wkb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
)

type Polygon []Point

func (p Polygon) GeomFromText() string {

	var buf strings.Builder

	buf.WriteString("GeomFromText('Polygon((")

	for k,v := range p {
		buf.WriteString(v.Lon + " " +v.Lat)
		if k < len(p) -1 {
			buf.WriteString(",")
		}
	}

	buf.WriteString("))')")

	return buf.String()
}

// 通过自定义类型创建记录
type Point struct {
	Lon string
	Lat string
}

func (loc *Point) GeomFromText() string {
	return fmt.Sprintf("GeomFromText('Point(%v %v)')", loc.Lon, loc.Lat)
}

func (loc *Point) Scan(v interface{}) error {

	if b, ok := v.([]byte); ok {
		p, err := wkb.Unmarshal(b[4:])
		if err != nil {
			return err
		}

		loc.Lon = strconv.FormatFloat(p.FlatCoords()[0], 'E', -1, 64)
		loc.Lat = strconv.FormatFloat(p.FlatCoords()[1], 'E', -1, 64)
	}

	return nil
}

func (loc Point) GormDataType() string {
	return "point"
}

func (loc Point) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%v %v)", loc.Lon, loc.Lat)},
	}
}

type PointQueryExpression struct {
	column   string
	polygon Polygon
}

func PointQuery(column string) *PointQueryExpression {
	return &PointQueryExpression{column: column}
}

func (pointQuery *PointQueryExpression) MBRContains(polygon Polygon) *PointQueryExpression {

	pointQuery.polygon = polygon

	return pointQuery
}

func (pointQuery *PointQueryExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		switch stmt.Dialector.Name() {
		case "mysql":

			_, _ = builder.WriteString(fmt.Sprintf("MBRContains(%s,%s)", pointQuery.polygon.GeomFromText(), pointQuery.column))

		case "postgres":

		}
	}
}
