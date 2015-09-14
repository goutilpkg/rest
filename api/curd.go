package api

import (
	"github.com/gin-gonic/gin"
	. "pkg.deepin.io/server/utils/logger"
	"reflect"
)

type CURD interface {
	GetBy(key string, value ...interface{}) error
	Create() error
	Check() error
	Delete() error
	Data() interface{}
}

func Create(v interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		rr := NewRsp(c)
		defer rr.Render()

		curd := reflect.New(reflect.TypeOf(v).Elem()).Interface().(CURD)
		if err := c.Bind(curd); nil != err {
			Logger.Warning("Check Data %v Failed: %v", curd, err)
			rr.Error(400, NewError(ErrIllegalDataFormat, err.Error()))
			return
		}

		if err := curd.Check(); nil != err {
			Logger.Warning("Check Data %v Failed: %v", curd, err)
			rr.Error(400, NewError(err, ""))
			return
		}

		if err := curd.Create(); nil != err {
			Logger.Error("Create Data %v Failed: %v", curd, err)
			rr.Error(400, NewError(err, ""))
			return
		}
		rr.Data = curd.Data()
	}
}

func Delete(v interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		rr := NewRsp(c)
		defer rr.Render()

		curd := reflect.New(reflect.TypeOf(v).Elem()).Interface().(CURD)
		id := c.Params.ByName("id")
		if err := curd.GetBy("`id`=?", id); nil != err {
			Logger.Error("%v", err)
			rr.Error(404, NewError(err, ""))
			return
		}

		if err := curd.Delete(); nil != err {
			Logger.Error("%v", err)
			rr.Error(400, NewError(err, ""))
			return
		}

		v, _ := c.Get("Callback")
		if callback, _ := v.(func(interface{})); callback != nil {
			callback(curd)
		}
		rr.Data = curd.Data()
	}
}
