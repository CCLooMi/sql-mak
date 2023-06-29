package utils

import (
	"fmt"
	"go/importer"
	"reflect"
)

// 包结构体
type Package struct {
	Module reflect.Value
}

// 加载模块
func LoadModule(pkgName string) (reflect.Value, error) {
	pkg, err := ImportPackage(pkgName)
	if err != nil {
		return reflect.Value{}, err
	}
	module := pkg.Module
	if module.IsValid() {
		return module, nil
	}
	return reflect.Value{}, fmt.Errorf("Module not found")
}

// 导入包
func ImportPackage(pkgName string) (*Package, error) {
	pkg, err := importer.Default().Import(pkgName)
	if err != nil {
		return nil, err
	}
	return &Package{Module: reflect.ValueOf(pkg)}, nil
}
