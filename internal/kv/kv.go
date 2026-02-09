package kv

// Package kv provides a local alias to the gofiber Storage interface
// to make it explicit that we use it as a Key-Value store abstraction.
//
// This helps avoid confusion with other kinds of storage (e.g., filesystems).
//
// Under the hood, it maps directly to github.com/gofiber/storage.Storage,
// so all existing Fiber storage drivers (redis, badger, etc.) are compatible.

import (
	"github.com/gofiber/storage"
)

// KVStore is a key-value storage abstraction used across the project.
// It is an alias to gofiber's storage.Storage interface.
//
// Use implementations from github.com/gofiber/storage/* packages
// such as redis or badger when constructing a KVStore.
//
// Since this is a type alias (not a new interface), there is no runtime
// overhead and all methods/implementations remain identical.
// You can continue using storage drivers exactly as before.
//
// Example:
//   import (
//     "github.com/gofiber/storage/redis/v2"
//     "news-scrabber/internal/kv"
//   )
//   var store kv.KVStore = redis.New(redis.Config{ URL: "redis://..." })
//
// Methods available are the same as storage.Storage.
// See https://pkg.go.dev/github.com/gofiber/storage for details.
 type KVStore = storage.Storage
