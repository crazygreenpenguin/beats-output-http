package resolver

import (
	"context"
	"net"
	"sync"
	"time"
)

const MaxCacheSize = 512 //Must be greater then 16

func NewDNSResolver() *DNSResolver {
	tmp := DNSResolver{
		resolver: net.DefaultResolver,
		cache:    make(map[string]DNSRecord),
	}

	go func() {
		for range time.Tick(60 * time.Second) {
			tmp.refresh()
		}
	}()
	return &tmp
}

type DNSRecord struct {
	addr []string
	err  error
}

type DNSResolver struct {
	resolver *net.Resolver

	cache map[string]DNSRecord
	sync.Mutex
}

func (dr *DNSResolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	dr.Lock()
	defer dr.Unlock()

	if val, ok := dr.cache[host]; ok {
		return val.addr, val.err
	}

	addr, err := dr.resolver.LookupHost(ctx, host)
	dr.cache[host] = DNSRecord{
		addr: addr,
		err:  err,
	}
	return addr, err
}

func (dr *DNSResolver) refresh() {
	dr.Lock()
	defer dr.Unlock()

	//free cache space
	if len(dr.cache) > MaxCacheSize {
		for key := range dr.cache {
			delete(dr.cache, key)
			if len(dr.cache) <= MaxCacheSize-16 {
				break
			}
		}
	}

	for key, val := range dr.cache {
		addr, err := net.LookupHost(key)
		val.addr = addr
		val.err = err
	}
}
