package worker

import (
	"log"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type Stats struct {
	MemStats  *mem.VirtualMemoryStat
	DiskStats *disk.UsageStat
	CpuStats  *cpu.TimesStat
	LoadStats *load.AvgStat
	TaskCount int
}

func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.Total / 1024
}

func (s *Stats) MemAvailableKb() uint64 {
	return s.MemStats.Available / 1024
}

func (s *Stats) MemUsedPercent() uint64 {
	return uint64(s.MemStats.UsedPercent)
}

func (s *Stats) DiskTotal() uint64 {
	return s.DiskStats.Total
}

func (s *Stats) DiskFree() uint64 {
	return s.DiskStats.Free
}

func (s *Stats) DiskUsed() uint64 {
	return s.DiskStats.Used
}

func (s *Stats) CpuUsage() float64 {
	idle := s.CpuStats.Idle + s.CpuStats.Iowait
	nonIdle := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System + s.CpuStats.Irq + s.CpuStats.Softirq + s.CpuStats.Steal

	total := idle + nonIdle

	if total == 0 {
		return 0.00
	}

	return (float64(total) - float64(idle)) / float64(total)
}

func GetStats() *Stats {
	return &Stats{
		MemStats:  GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		CpuStats:  GetCpuInfo(),
		LoadStats: GetLoadAvg(),
	}
}

func GetMemoryInfo() *mem.VirtualMemoryStat {
	stats, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Failed to get memory info: %v", err)
		return &mem.VirtualMemoryStat{}
	}
	return stats
}

func GetDiskInfo() *disk.UsageStat {
	stats, err := disk.Usage("/")
	if err != nil {
		log.Printf("Failed to get disk info: %v", err)
		return &disk.UsageStat{}
	}
	return stats
}

func GetCpuInfo() *cpu.TimesStat {
	stats, err := cpu.Times(false)
	if err != nil {
		log.Printf("Failed to get cpu info: %v", err)
		return &cpu.TimesStat{}
	}
	if len(stats) == 0 {
		return &cpu.TimesStat{}
	}
	return &stats[0]
}

func GetLoadAvg() *load.AvgStat {
	stats, err := load.Avg()
	if err != nil {
		log.Printf("Failed to get load avg: %v", err)
		return &load.AvgStat{}
	}
	return stats
}
