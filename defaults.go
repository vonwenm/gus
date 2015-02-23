package main

// This is the default configuration file for the GUS system. You should select any
// database packages here, alter the constants that are used in the code and put any
// system-wide configurations here.
import (
	/*  Database support. This is compiled in and available for selection.
	    To remove the support (and overhead) you should comment it out. */
	_ "github.com/cgentry/gus/storage/sqlite"

	/* Encryption support */
	_ "github.com/cgentry/gus/encryption/drivers/bcrypt"
	_ "github.com/cgentry/gus/encryption/drivers/sha512"
	/* REMOVE WHEN IN PRODUCTION */
	_ "github.com/cgentry/gus/encryption/drivers/plaintext"
)

const (
	DEFAULT_CONFIG_FILENAME    = "/etc/gus/config.json"
	DEFAULT_CONFIG_PERMISSIONS = 0600
)
