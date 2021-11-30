//go:build wireinject
// +build wireinject

package main

import (
	"crud-with-auth/api"
	"crud-with-auth/db"

	"github.com/google/wire"
)

func InitializeApp() *api.Api {
	wire.Build(db.ProvideDB, api.ProviderAPI)
	return &api.Api{}
}
