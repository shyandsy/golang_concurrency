package cache

import (
	"testing"
	"time"
)

func checkKeyExpired(t *testing.T, c Cache, key string) {
	_, ok := c.Get(key)
	if ok {
		t.Fatalf("expetect key `%s` expired", key)
	}
}

func checkKeyExist(t *testing.T, c Cache, key string, value string) {
	v, ok := c.Get(key)
	if !ok {
		t.Fatalf("expetect key `%s` exist", key)
	}
	if v != value {
		t.Fatalf("expetect key `%s` value `%s`, got `%s`", key, value, v)
	}
}

func TestSetGet(t *testing.T) {
	c := NewMapCahe()

	c.Set("AAA", "111", 3*time.Second)
	c.Set("B", "", 3*time.Second)
	c.Set("", "", 3*time.Second)

	checkKeyExist(t, c, "AAA", "111")
	checkKeyExist(t, c, "B", "")
	checkKeyExist(t, c, "", "")

	c.Close()
}

func TestExpiration(t *testing.T) {
	c := NewMapCahe()

	c.Set("AAA", "111", 3*time.Second)

	checkKeyExist(t, c, "AAA", "111")

	time.Sleep(3 * time.Second)

	checkKeyExpired(t, c, "AAA")

	c.Close()
}

func TestPartialExpiration(t *testing.T) {
	c := NewMapCahe()

	c.Set("AAA", "111", 1*time.Second)
	c.Set("B", "", 3*time.Second)
	c.Set("", "", 6*time.Second)

	checkKeyExist(t, c, "AAA", "111")
	checkKeyExist(t, c, "B", "")
	checkKeyExist(t, c, "", "")

	time.Sleep(1 * time.Second)

	checkKeyExpired(t, c, "AAA")

	time.Sleep(3 * time.Second)

	checkKeyExpired(t, c, "B")

	time.Sleep(6 * time.Second)

	checkKeyExpired(t, c, "")

	c.Close()
}
