// Package builddeps keeps dependencies imported only by Revel's generated
// application entrypoint visible to Go module maintenance commands.
package builddeps

import (
	_ "github.com/bradfitz/gomemcache/memcache"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/patrickmn/go-cache"
)
