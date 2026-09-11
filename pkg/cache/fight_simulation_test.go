package cache

import (
	"testing"

	"github.com/br-lemes/golem/pkg/models"
)

func TestFightSimulationCacheSavesImmediately(t *testing.T) {
	initializeTestCache(t)
	simulation := models.FightSimulation{Key: "immediate", Version: 1, XP: 42}
	SaveFightSimulation(simulation)
	got, found := GetFightSimulation("immediate", 1)
	if !found || got.XP != 42 {
		t.Fatalf("GetFightSimulation() = %#v, %t", got, found)
	}
	_, found = GetFightSimulation("immediate", 2)
	if found {
		t.Fatal("simulation with wrong version was found")
	}
}

func TestFightSimulationCacheBatch(t *testing.T) {
	initializeTestCache(t)
	batchFightSimulations = false
	FlushFightSimulationBatch()
	BeginFightSimulationBatch()
	SaveFightSimulation(models.FightSimulation{Key: "batch", Version: 1, XP: 7})
	_, found := GetFightSimulation("batch", 1)
	if found {
		t.Fatal("batched simulation was saved before flush")
	}
	FlushFightSimulationBatch()
	got, found := GetFightSimulation("batch", 1)
	if !found || got.XP != 7 {
		t.Fatalf("flushed simulation = %#v, %t", got, found)
	}
	FlushFightSimulationBatch()
}

func TestFightSimulationCacheFlushesEmptyBatch(t *testing.T) {
	initializeTestCache(t)
	BeginFightSimulationBatch()
	FlushFightSimulationBatch()
	if batchFightSimulations || pendingFightSimulations != nil {
		t.Fatal("empty batch was not reset")
	}
}
