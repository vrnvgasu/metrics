package agent

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	models "github.com/vrnvgasu/metrics/internal/model"
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

const ()

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

func (a *Agent) Collect(ctx context.Context, interval time.Duration) error {
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

		time.Sleep(interval * time.Second)
	}
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
		value := a.getMemStatsMetric(name, memStats)
		a.pushMetric(models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &value,
		})
	}
}

func (a *Agent) getMemStatsMetric(metricName string, memStats runtime.MemStats) float64 {
	switch metricName {
	case alloc:
		return float64(memStats.Alloc)
	case buckHashSys:
		return float64(memStats.BuckHashSys)
	case frees:
		return float64(memStats.Frees)
	case gCCPUFraction:
		return memStats.GCCPUFraction
	case gCSys:
		return float64(memStats.GCSys)
	case heapAlloc:
		return float64(memStats.HeapAlloc)
	case heapIdle:
		return float64(memStats.HeapIdle)
	case heapInuse:
		return float64(memStats.HeapInuse)
	case heapObjects:
		return float64(memStats.HeapObjects)
	case heapReleased:
		return float64(memStats.HeapReleased)
	case heapSys:
		return float64(memStats.HeapSys)
	case lastGC:
		return float64(memStats.LastGC)
	case lookups:
		return float64(memStats.Lookups)
	case mCacheInuse:
		return float64(memStats.MCacheInuse)
	case mCacheSys:
		return float64(memStats.MCacheSys)
	case mSpanInuse:
		return float64(memStats.MSpanInuse)
	case mSpanSys:
		return float64(memStats.MSpanSys)
	case mallocs:
		return float64(memStats.Mallocs)
	case nextGC:
		return float64(memStats.NextGC)
	case numForcedGC:
		return float64(memStats.NumForcedGC)
	case numGC:
		return float64(memStats.NumGC)
	case otherSys:
		return float64(memStats.OtherSys)
	case pauseTotalNs:
		return float64(memStats.PauseTotalNs)
	case stackInuse:
		return float64(memStats.StackInuse)
	case stackSys:
		return float64(memStats.StackSys)
	case sys:
		return float64(memStats.Sys)
	case totalAlloc:
		return float64(memStats.TotalAlloc)
	default:
		return 0
	}
}
