package main

import (
	"Sadeem-RestAPI/internal/models"
	"Sadeem-RestAPI/internal/server"
	"Sadeem-RestAPI/internal/translation"
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml"
	"golang.org/x/text/language"
)

//go:generate goi18n extract -sourceLanguage en
//go:generate goi18n extract -sourceLanguage en

func main() {
	// Initilize goi18n
	i18nInit()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("could not connect to database")
	}
	server := server.NewServer()

	models.Models = &models.ModelStruct{
		User: &models.UserModel{
			DB: pool,
		},
	}

	print("starting server at http://localhost", server.Addr)

	err = server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}

func i18nInit() {
	translation.Bundle = *i18n.NewBundle(language.English)
	translation.Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	//_, _ = translation.Bundle.LoadMessageFileFS(LocalFS, "active.en.toml")
	//_, _ = translation.Bundle.LoadMessageFileFS(LocalFS, "active.ar.toml")
	_, _ = translation.Bundle.LoadMessageFile("active.en.toml")
	_, _ = translation.Bundle.LoadMessageFile("active.ar.toml")
}
