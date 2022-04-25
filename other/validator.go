package other

import (
	"MFile/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"os"
	"reflect"
)

func RegisterValidation(e *gin.Engine) {
	validation, ok := binding.Validator.Engine().(*validator.Validate)

	if !ok {
		os.Exit(1)
	}

	err:=validation.RegisterValidation("need", needleValidator)
	if err!=nil{
		logger.MLogger.Error("register validation err",zap.Error(err))
	}
}

func needleValidator(f validator.FieldLevel) bool {
	fl := f.Field()
	if !fl.IsValid() || fl.Kind() == reflect.Ptr && fl.IsNil() {
		return false
	}
	return true
}
