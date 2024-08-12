package cache

import (
    "github.com/bradfitz/gomemcache/memcache"
    "encoding/json"
)

type MemcachedClient struct {
    client *memcache.Client
}

func NewMemcachedClient(serverAddress string) *MemcachedClient {
    return &MemcachedClient{
        client: memcache.New(serverAddress),
    }
}

func (m *MemcachedClient) Set(key string, value interface{}, expiration int32) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return m.client.Set(&memcache.Item{Key: key, Value: data, Expiration: expiration})
}

func (m *MemcachedClient) Get(key string, value interface{}) error {
    item, err := m.client.Get(key)
    if err != nil {
        return err
    }
    return json.Unmarshal(item.Value, value)
}

func (m *MemcachedClient) Delete(key string) error {
    return m.client.Delete(key)
}