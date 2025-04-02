package main

import "fileupbackendv2/internal/storage"

func migrate() {
	storage.Migrate()
}
func checkMigration() {
	storage.CheckMigrationSuccessFull()
}
