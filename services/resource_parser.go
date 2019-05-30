package services

type ResourceParser func(data string) (interface{}, int64, error)
