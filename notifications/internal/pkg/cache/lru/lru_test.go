package lru

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
 Название тестов:
 Test{структура}_{метод}_{данные}_{ожидание}
 Пример: TestBuilder_BuildFeneralStep_hasID_buildStep

*/

func TestLRUCache_Get_Item_OK(t *testing.T) {
	t.Parallel()
	// Arrange
	cache := NewLRUCache[string, string](5)
	key := "test_key"
	value := &Item[string, string]{key: key, value: "test_value"}
	elem := cache.queue.PushFront(value)
	cache.items[key] = elem

	// Act
	el, ok := cache.Get(key)

	// Assert
	require.True(t, ok)
	require.Equal(t, value.value, el)
}

func TestLRUCache_Set_Item_OK(t *testing.T) {
	t.Parallel()

	// Arrange
	cache := NewLRUCache[string, string](5)
	key := "test_key"
	value := "test_value"

	// Act
	err := cache.Set(key, value)

	// Assert
	el := cache.items[key].Value.(*Item[string, string]).value

	require.NoError(t, err)
	require.Equal(t, value, el)

}

func TestLRU_Set_ExistElementWithFULlQueueSync_MoveToFront(t *testing.T) {
	t.Parallel()
	// Arrange
	cache := NewLRUCache[string, int](3)
	cache.Set("Vasya", 10)
	cache.Set("Petya", 11)
	cache.Set("Kolya", 15)

	// Act
	cache.Set("Vasva", 15)

	resultFront, _ := cache.queue.Front().Value.(*Item[string, int])
	resultBack, _ := cache.queue.Back().Value.(*Item[string, int])

	require.Equal(t, resultFront.value, 15)
	require.Equal(t, resultBack.value, 11)
	require.Equal(t, cache.queue.Len(), 3)

}

func TestLRUCache_Set_MoreThanCap_MaxCap(t *testing.T) {
	t.Parallel()

	// Arrange
	cache := NewLRUCache[string, string](5)

	cache.Set("test_1", "value_1")
	cache.Set("test_2", "value_2")
	cache.Set("test_3", "value_3")
	cache.Set("test_4", "value_4")
	cache.Set("test_5", "value_5")
	cache.Set("test_6", "value_6")

	// Assert
	el := cache.items["test_1"]

	require.Nil(t, el)
}

func TestLRUCache_Delete_Item_OK(t *testing.T) {
	t.Parallel()
	// Arrange
	cache := NewLRUCache[string, string](5)
	key := "test_key"
	value := "test_value"
	elem := cache.queue.PushFront(value)
	cache.items[key] = elem

	// Act
	cache.Delete(key)

	// Assert
	_, ok := cache.items[key]

	require.False(t, ok)
}

func TestLRUCache_Clear_Cache_OK(t *testing.T) {
	t.Parallel()
	// Arrange
	cache := NewLRUCache[string, string](5)
	key := "test_key"
	value := "test_value"
	elem := cache.queue.PushFront(value)
	cache.items[key] = elem

	// Act
	cache.Clear()

	// Assert
	el := cache.queue.Front()
	require.Zero(t, len(cache.items))
	require.Nil(t, el)
}

func TestLRUCache_Count_Items_OK(t *testing.T) {
	t.Parallel()
	// Arrange
	cache := NewLRUCache[string, string](5)
	key := "test_key"
	value := "test_value"
	elem := cache.queue.PushFront(value)
	cache.items[key] = elem

	// Act
	el := cache.Count()

	// Assert
	require.Equal(t, el, 1)
}

func TestLRUCache_ConcurrentSafety_GetSetItems_OK(t *testing.T) {
	t.Parallel()
	// Arrange
	capacity := 3
	cache := NewLRUCache[string, string](capacity)
	var wg sync.WaitGroup
	wg.Add(capacity + capacity)
	// Act
	for i := 1; i < cache.capacity+1; i++ {
		li := i
		go func() {
			defer wg.Done()
			err := cache.Set(fmt.Sprintf("key%d", li), fmt.Sprintf("value%d", li))
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()
			cache.Get(fmt.Sprintf("key%d", li))
		}()
	}

	wg.Wait()

	size := cache.Count()
	require.Equal(t, 3, size)
	for i := 1; i < cache.capacity+1; i++ {
		value, ok := cache.Get(fmt.Sprintf("key%d", i))
		require.True(t, ok)
		require.Equal(t, fmt.Sprintf("value%d", i), value)
	}

}
