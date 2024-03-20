package gutil

import (
	"strconv"
	"time"
)

// DefaultIdGenerator 是默认的 ID 生成器实现。
type DefaultIdGenerator struct {
	Options    *IdGeneratorOptions // ID 生成器的选项
	SnowWorker ISnowWorker         // Snowflake 算法的工作器
}

// NewDefaultIdGenerator 创建一个新的默认 ID 生成器。
func NewDefaultIdGenerator(options *IdGeneratorOptions) *DefaultIdGenerator {
	// 检查选项是否为空
	if options == nil {
		panic("idgen: options cannot be nil")
	}

	// 检查 BaseTime 是否在有效范围内
	minTime := int64(631123200000) // 1990-01-01 00:00:00
	if options.BaseTime < minTime || options.BaseTime > time.Now().UnixNano()/1e6 {
		panic("idgen: BaseTime out of range")
	}

	// 检查 WorkerIdBitLength 是否在有效范围内
	if options.WorkerIdBitLength <= 0 || options.WorkerIdBitLength > 21 {
		panic("idgen: WorkerIdBitLength out of range (1-21)")
	}

	// 检查 WorkerIdBitLength + SeqBitLength 是否小于等于 22
	if options.WorkerIdBitLength+options.SeqBitLength > 22 {
		panic("idgen: WorkerIdBitLength + SeqBitLength should be less than or equal to 22")
	}

	// 检查 WorkerId 是否在有效范围内
	maxWorkerIdNumber := uint16(1<<options.WorkerIdBitLength) - 1
	if maxWorkerIdNumber == 0 {
		maxWorkerIdNumber = 63
	}
	if options.WorkerId < 0 || options.WorkerId > maxWorkerIdNumber {
		panic("idgen: WorkerId out of range (0-" + strconv.FormatUint(uint64(maxWorkerIdNumber), 10) + ")")
	}

	// 检查 SeqBitLength 是否在有效范围内
	if options.SeqBitLength < 2 || options.SeqBitLength > 21 {
		panic("idgen: SeqBitLength out of range (2-21)")
	}

	// 检查 MaxSeqNumber 是否在有效范围内
	maxSeqNumber := uint32(1<<options.SeqBitLength) - 1
	if maxSeqNumber == 0 {
		maxSeqNumber = 63
	}
	if options.MaxSeqNumber < 0 || options.MaxSeqNumber > maxSeqNumber {
		panic("idgen: MaxSeqNumber out of range (1-" + strconv.FormatUint(uint64(maxSeqNumber), 10) + ")")
	}

	// 检查 MinSeqNumber 是否在有效范围内
	if options.MinSeqNumber < 5 || options.MinSeqNumber > maxSeqNumber {
		panic("idgen: MinSeqNumber out of range (5-" + strconv.FormatUint(uint64(maxSeqNumber), 10) + ")")
	}

	// 检查 TopOverCostCount 是否在有效范围内
	if options.TopOverCostCount < 0 || options.TopOverCostCount > 10000 {
		panic("idgen: TopOverCostCount out of range (0-10000)")
	}

	// 创建 Snowflake 算法的工作器
	var snowWorker = NewSnowWorkerDefault(options)
	time.Sleep(510 * time.Microsecond)

	// 返回默认 ID 生成器实例
	return &DefaultIdGenerator{
		Options:    options,
		SnowWorker: snowWorker,
	}
}

// NewLong 生成一个长整型的 ID。
func (dig DefaultIdGenerator) NewLong() int64 {
	return dig.SnowWorker.NextId()
}

// ExtractTime 从 ID 中提取时间信息。
func (dig DefaultIdGenerator) ExtractTime(id int64) time.Time {
	// 将时间信息从 ID 中提取出来并返回
	return time.UnixMilli(id>>(dig.Options.WorkerIdBitLength+dig.Options.SeqBitLength) + dig.Options.BaseTime)
}
