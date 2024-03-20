package gutil

import (
	"sync"
	"time"
)

type ISnowWorker interface {
	NextId() int64
}

// SnowWorkerDefault 实现了基于 Snowflake 算法的工作器。
type SnowWorkerDefault struct {
	BaseTime          int64  // 基础时间
	WorkerId          uint16 // 机器码
	WorkerIdBitLength byte   // 机器码位长
	SeqBitLength      byte   // 自增序列数位长
	MaxSeqNumber      uint32 // 最大序列数（含）
	MinSeqNumber      uint32 // 最小序列数（含）
	TopOverCostCount  uint32 // 最大漂移次数
	_TimestampShift   byte
	_CurrentSeqNumber uint32

	_LastTimeTick           int64
	_TurnBackTimeTick       int64
	_TurnBackIndex          byte
	_IsOverCost             bool
	_OverCostCountInOneTerm uint32

	sync.Mutex
}

// NewSnowWorkerDefault 创建一个新的 Snowflake 算法的工作器。
func NewSnowWorkerDefault(options *IdGeneratorOptions) ISnowWorker {
	var workerIdBitLength byte
	var seqBitLength byte
	var maxSeqNumber uint32

	// 1.BaseTime
	var baseTime int64
	if options.BaseTime != 0 {
		baseTime = options.BaseTime
	} else {
		baseTime = 1582136402000
	}

	// 2.WorkerIdBitLength
	if options.WorkerIdBitLength == 0 {
		workerIdBitLength = 6
	} else {
		workerIdBitLength = options.WorkerIdBitLength
	}

	// 3.WorkerId
	var workerId = options.WorkerId

	// 4.SeqBitLength
	if options.SeqBitLength == 0 {
		seqBitLength = 6
	} else {
		seqBitLength = options.SeqBitLength
	}

	// 5.MaxSeqNumber
	if options.MaxSeqNumber <= 0 {
		maxSeqNumber = (1 << seqBitLength) - 1
	} else {
		maxSeqNumber = options.MaxSeqNumber
	}

	// 6.MinSeqNumber
	var minSeqNumber = options.MinSeqNumber

	// 7.TopOverCostCount
	var topOverCostCount = options.TopOverCostCount
	if topOverCostCount == 0 {
		topOverCostCount = 2000
	}

	timestampShift := (byte)(workerIdBitLength + seqBitLength)
	currentSeqNumber := minSeqNumber

	return &SnowWorkerDefault{
		BaseTime:          baseTime,
		WorkerIdBitLength: workerIdBitLength,
		WorkerId:          workerId,
		SeqBitLength:      seqBitLength,
		MaxSeqNumber:      maxSeqNumber,
		MinSeqNumber:      minSeqNumber,
		TopOverCostCount:  topOverCostCount,
		_TimestampShift:   timestampShift,
		_CurrentSeqNumber: currentSeqNumber,

		_LastTimeTick:           0,
		_TurnBackTimeTick:       0,
		_TurnBackIndex:          0,
		_IsOverCost:             false,
		_OverCostCountInOneTerm: 0,
	}
}

// NextOverCostId 生成下一个超出成本的 ID。
func (m1 *SnowWorkerDefault) NextOverCostId() int64 {
	currentTimeTick := m1.GetCurrentTimeTick()
	if currentTimeTick > m1._LastTimeTick {
		m1._LastTimeTick = currentTimeTick
		m1._CurrentSeqNumber = m1.MinSeqNumber
		m1._IsOverCost = false
		m1._OverCostCountInOneTerm = 0
		return m1.CalcId(m1._LastTimeTick)
	}
	if m1._OverCostCountInOneTerm >= m1.TopOverCostCount {
		m1._LastTimeTick = m1.GetNextTimeTick()
		m1._CurrentSeqNumber = m1.MinSeqNumber
		m1._IsOverCost = false
		m1._OverCostCountInOneTerm = 0
		return m1.CalcId(m1._LastTimeTick)
	}
	if m1._CurrentSeqNumber > m1.MaxSeqNumber {
		m1._LastTimeTick++
		m1._CurrentSeqNumber = m1.MinSeqNumber
		m1._IsOverCost = true
		m1._OverCostCountInOneTerm++
		return m1.CalcId(m1._LastTimeTick)
	}
	return m1.CalcId(m1._LastTimeTick)
}

// NextNormalId 生成下一个正常的 ID。
func (m1 *SnowWorkerDefault) NextNormalId() int64 {
	currentTimeTick := m1.GetCurrentTimeTick()
	if currentTimeTick < m1._LastTimeTick {
		if m1._TurnBackTimeTick < 1 {
			m1._TurnBackTimeTick = m1._LastTimeTick - 1
			m1._TurnBackIndex++
			if m1._TurnBackIndex > 4 {
				m1._TurnBackIndex = 1
			}
		}
		return m1.CalcTurnBackId(m1._TurnBackTimeTick)
	}
	if m1._TurnBackTimeTick > 0 {
		m1._TurnBackTimeTick = 0
	}
	if currentTimeTick > m1._LastTimeTick {
		m1._LastTimeTick = currentTimeTick
		m1._CurrentSeqNumber = m1.MinSeqNumber
		return m1.CalcId(m1._LastTimeTick)
	}
	if m1._CurrentSeqNumber > m1.MaxSeqNumber {
		m1._LastTimeTick++
		m1._CurrentSeqNumber = m1.MinSeqNumber
		m1._IsOverCost = true
		m1._OverCostCountInOneTerm = 1
		return m1.CalcId(m1._LastTimeTick)
	}
	return m1.CalcId(m1._LastTimeTick)
}

// CalcId 计算给定时间戳的 ID。
func (m1 *SnowWorkerDefault) CalcId(useTimeTick int64) int64 {
	result := int64(useTimeTick<<m1._TimestampShift) + int64(m1.WorkerId<<m1.SeqBitLength) + int64(m1._CurrentSeqNumber)
	m1._CurrentSeqNumber++
	return result
}

// CalcTurnBackId 计算时间回拨时的 ID。
func (m1 *SnowWorkerDefault) CalcTurnBackId(useTimeTick int64) int64 {
	result := int64(useTimeTick<<m1._TimestampShift) + int64(m1.WorkerId<<m1.SeqBitLength) + int64(m1._TurnBackIndex)
	m1._TurnBackTimeTick--
	return result
}

// GetCurrentTimeTick .
func (m1 *SnowWorkerDefault) GetCurrentTimeTick() int64 {
	var millis = time.Now().UnixNano() / 1e6
	return millis - m1.BaseTime
}

// GetNextTimeTick .
func (m1 *SnowWorkerDefault) GetNextTimeTick() int64 {
	tempTimeTicker := m1.GetCurrentTimeTick()
	for tempTimeTicker <= m1._LastTimeTick {
		time.Sleep(time.Duration(1) * time.Millisecond)
		tempTimeTicker = m1.GetCurrentTimeTick()
	}
	return tempTimeTicker
}

// NextId .
func (m1 *SnowWorkerDefault) NextId() int64 {
	m1.Lock()
	defer m1.Unlock()
	if m1._IsOverCost {
		return m1.NextOverCostId()
	} else {
		return m1.NextNormalId()
	}
}
