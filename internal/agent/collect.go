package agent

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/pkg/helper"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

const (
	// Alloc - текущий объем выделенной памяти в байтах, который еще не освобожден.
	// Показывает "живую" память программы на данный момент.
	alloc = "Alloc"

	// BuckHashSys - память в байтах, используемая для хеш-таблиц профилирования.
	// Внутренние структуры для сбора статистики и профилирования Go runtime.
	buckHashSys = "BuckHashSys"

	// Frees - общее количество вызовов освобождения памяти за время работы программы.
	// Счетчик увеличивается при каждом освобождении объекта сборщиком мусора.
	frees = "Frees"

	// GCCPUFraction - доля времени процессора, потраченного на сборку мусора.
	// Значение от 0.0 до 1.0, где 0.05 означает 5% времени CPU на GC.
	gCCPUFraction = "GCCPUFraction"

	// GCSys - память в байтах, используемая для метаданных сборщика мусора.
	// Память, выделенная под структуры данных, необходимые для работы GC.
	gCSys = "GCSys"

	// HeapAlloc - байты, выделенные в куче и все еще используемые.
	// Аналогично Alloc, но более точно отражает использование именно кучи.
	heapAlloc = "HeapAlloc"

	// HeapIdle - байты в свободных (неиспользуемых) частях кучи.
	// Память, зарезервированная для кучи, но в данный момент свободная.
	heapIdle = "HeapIdle"

	// HeapInuse - байты в активных (используемых) частях кучи.
	// Память кучи, которая сейчас содержит данные программы.
	heapInuse = "HeapInuse"

	// HeapObjects - количество объектов, находящихся в куче в данный момент.
	// Разница между Mallocs и Frees - количество "живых" объектов.
	heapObjects = "HeapObjects"

	// HeapReleased - байты физической памяти, возвращенные операционной системе.
	// Go может освобождать неиспользуемую память, чтобы уменьшить RSS (Resident Set Size).
	heapReleased = "HeapReleased"

	// HeapSys - общий объем памяти в байтах, запрошенный у ОС для кучи.
	// Включает HeapInuse + HeapIdle + служебные структуры управления памятью.
	heapSys = "HeapSys"

	// LastGC - время последней сборки мусора в наносекундах с Unix эпохи.
	// Для получения времени в секундах: (time.Now().UnixNano() - LastGC) / 1e9
	lastGC = "LastGC"

	// Lookups - количество поисков указателей в runtime за время работы программы.
	// Используется для отладки и анализа производительности.
	lookups = "Lookups"

	// MCacheInuse - байты, используемые в локальных кэшах потоков (mcache).
	// Кэши для быстрого выделения мелких объектов без блокировок.
	mCacheInuse = "MCacheInuse"

	// MCacheSys - общий объем памяти для локальных кэшей потоков.
	// Включает как активную, так и зарезервированную память для mcache.
	mCacheSys = "MCacheSys"

	// MSpanInuse - байты, используемые структурами mspan (управление страницами памяти).
	// Структуры, которые отслеживают состояние страниц памяти в куче.
	mSpanInuse = "MSpanInuse"

	// MSpanSys - общий объем памяти для структур mspan.
	// Вся память, выделенная под управление страницами памяти.
	mSpanSys = "MSpanSys"

	// Mallocs - общее количество выделений памяти за время работы программы.
	// Счетчик увеличивается при каждом создании объекта в куче.
	mallocs = "Mallocs"

	// NextGC - целевое значение HeapAlloc для следующей сборки мусора в байтах.
	// Когда HeapAlloc достигнет этого значения, будет запущен GC.
	nextGC = "NextGC"

	// NumForcedGC - количество сборок мусора, явно вызванных через runtime.GC().
	// Показывает количество принудительных сборок, в отличие от автоматических.
	numForcedGC = "NumForcedGC"

	// NumGC - общее количество завершенных циклов сборки мусора.
	// Инкрементируется каждый раз после успешного завершения GC.
	numGC = "NumGC"

	// OtherSys - байты памяти, используемые для прочих нужд runtime.
	// Включает память для системных структур, не вошедших в другие категории.
	otherSys = "OtherSys"

	// PauseTotalNs - суммарное время в наносекундах всех пауз сборки мусора.
	// Показывает общее время, в течение которого программа была остановлена из-за GC.
	pauseTotalNs = "PauseTotalNs"

	// StackInuse - байты, используемые стеками горутин.
	// Активно используемая память для хранения стековых фреймов.
	stackInuse = "StackInuse"

	// StackSys - общий объем памяти для стеков горутин.
	// Вся память, выделенная под стеки, включая зарезервированную.
	stackSys = "StackSys"

	// Sys - общий объем виртуальной памяти в байтах, запрошенной у ОС.
	// Сумма всей памяти, которую Go runtime запросил у операционной системы.
	sys = "Sys"

	// TotalAlloc - общий объем выделенной памяти за все время работы программы.
	// Кумулятивная сумма всех выделений, включая уже освобожденные.
	totalAlloc = "TotalAlloc"
)

var gaugesMemStatNames = []string{
	alloc,
	buckHashSys,
	frees,
	gCCPUFraction,
	gCSys,
	heapAlloc,
	heapIdle,
	heapInuse,
	heapObjects,
	heapReleased,
	heapSys,
	lastGC,
	lookups,
	mCacheInuse,
	mCacheSys,
	mSpanInuse,
	mSpanSys,
	mallocs,
	nextGC,
	numForcedGC,
	numGC,
	otherSys,
	pauseTotalNs,
	stackInuse,
	stackSys,
	sys,
	totalAlloc,
}

const (
	pollCount   = "PollCount"
	randomValue = "RandomValue"
)

func (a *Agent) Collect(ctx context.Context, cnf *config.AgentCnf) error {
	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			if err := a.addGopsutil(ctx); err != nil {
				return err
			}

			time.Sleep(time.Duration(cnf.PollInterval) * time.Second)
		}
	})

	errGroup.Go(func() error {
		var memStats runtime.MemStats

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			runtime.ReadMemStats(&memStats)

			a.addPollCount()
			a.addStatsMetric(memStats)
			a.addRandomValue()

			time.Sleep(time.Duration(cnf.PollInterval) * time.Second)
		}
	})

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("agent collect: %w", err)
	}

	close(a.Metrics)
	return nil
}

func (a *Agent) addGopsutil(ctx context.Context) error {
	v, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return fmt.Errorf("mem.VirtualMemoryWithContext: %w", err)
	}
	a.pushMetric(models.Metrics{
		ID:    "TotalMemory",
		MType: models.Gauge,
		Value: helper.NewRefFloat64(float64(v.Total)),
	})
	a.pushMetric(models.Metrics{
		ID:    "FreeMemory",
		MType: models.Gauge,
		Value: helper.NewRefFloat64(float64(v.Free)),
	})

	percentages, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		return fmt.Errorf("cpu.PercentWithContext: %w", err)
	}
	for i, percent := range percentages {
		a.pushMetric(models.Metrics{
			ID:    fmt.Sprintf("CPUutilization%d", i+1),
			MType: models.Gauge,
			Value: &percent,
		})
	}

	return nil
}

func (a *Agent) addPollCount() {
	a.pollCount.Add(1)
	delta := a.pollCount.Load()
	a.pushMetric(models.Metrics{
		ID:    pollCount,
		MType: models.Counter,
		Delta: &delta,
	})
}

func (a *Agent) addRandomValue() {
	value := rand.Float64()
	a.pushMetric(models.Metrics{
		ID:    randomValue,
		MType: models.Gauge,
		Value: &value,
	})
}

func (a *Agent) addStatsMetric(memStats runtime.MemStats) {
	for _, name := range gaugesMemStatNames {
		f, ok := gaugeMemStats[name]
		if !ok {
			continue
		}

		value := f(memStats)
		a.pushMetric(models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &value,
		})
	}
}

type memStatGetter func(runtime.MemStats) float64

var gaugeMemStats = map[string]memStatGetter{
	alloc:         func(ms runtime.MemStats) float64 { return float64(ms.Alloc) },
	buckHashSys:   func(ms runtime.MemStats) float64 { return float64(ms.BuckHashSys) },
	frees:         func(ms runtime.MemStats) float64 { return float64(ms.Frees) },
	gCCPUFraction: func(ms runtime.MemStats) float64 { return ms.GCCPUFraction },
	gCSys:         func(ms runtime.MemStats) float64 { return float64(ms.GCSys) },
	heapAlloc:     func(ms runtime.MemStats) float64 { return float64(ms.HeapAlloc) },
	heapIdle:      func(ms runtime.MemStats) float64 { return float64(ms.HeapIdle) },
	heapInuse:     func(ms runtime.MemStats) float64 { return float64(ms.HeapInuse) },
	heapObjects:   func(ms runtime.MemStats) float64 { return float64(ms.HeapObjects) },
	heapReleased:  func(ms runtime.MemStats) float64 { return float64(ms.HeapReleased) },
	heapSys:       func(ms runtime.MemStats) float64 { return float64(ms.HeapSys) },
	lastGC:        func(ms runtime.MemStats) float64 { return float64(ms.LastGC) },
	lookups:       func(ms runtime.MemStats) float64 { return float64(ms.Lookups) },
	mCacheInuse:   func(ms runtime.MemStats) float64 { return float64(ms.MCacheInuse) },
	mCacheSys:     func(ms runtime.MemStats) float64 { return float64(ms.MCacheSys) },
	mSpanInuse:    func(ms runtime.MemStats) float64 { return float64(ms.MSpanInuse) },
	mSpanSys:      func(ms runtime.MemStats) float64 { return float64(ms.MSpanSys) },
	mallocs:       func(ms runtime.MemStats) float64 { return float64(ms.Mallocs) },
	nextGC:        func(ms runtime.MemStats) float64 { return float64(ms.NextGC) },
	numForcedGC:   func(ms runtime.MemStats) float64 { return float64(ms.NumForcedGC) },
	numGC:         func(ms runtime.MemStats) float64 { return float64(ms.NumGC) },
	otherSys:      func(ms runtime.MemStats) float64 { return float64(ms.OtherSys) },
	pauseTotalNs:  func(ms runtime.MemStats) float64 { return float64(ms.PauseTotalNs) },
	stackInuse:    func(ms runtime.MemStats) float64 { return float64(ms.StackInuse) },
	stackSys:      func(ms runtime.MemStats) float64 { return float64(ms.StackSys) },
	sys:           func(ms runtime.MemStats) float64 { return float64(ms.Sys) },
	totalAlloc:    func(ms runtime.MemStats) float64 { return float64(ms.TotalAlloc) },
}
