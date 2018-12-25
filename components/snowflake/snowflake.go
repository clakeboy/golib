package snowflake

import (
	"fmt"
	"time"
)

/**
 * Twitter_Snowflake<br>
 * SnowFlake的结构如下(每部分用-分开):<br>
 * 0 - 0000000000 0000000000 0000000000 0000000000 0 - 00000 - 00000 - 000000000000 <br>
 * 1位标识，由于long基本类型在Java中是带符号的，最高位是符号位，正数是0，负数是1，所以id一般是正数，最高位是0<br>
 * 41位时间截(毫秒级)，注意，41位时间截不是存储当前时间的时间截，而是存储时间截的差值（当前时间截 - 开始时间截)
 * 得到的值），这里的的开始时间截，一般是我们的id生成器开始使用的时间，由我们程序来指定的（如下下面程序IdWorker类的startTime属性）。41位的时间截，可以使用69年，年T = (1L << 41) / (1000L * 60 * 60 * 24 * 365) = 69<br>
 * 10位的数据机器位，可以部署在1024个节点，包括5位datacenterId和5位workerId<br>
 * 12位序列，毫秒内的计数，12位的计数顺序号支持每个节点每毫秒(同一机器，同一时间截)产生4096个ID序号<br>
 * 加起来刚好64位，为一个Long型。<br>
 * SnowFlake的优点是，整体上按照时间自增排序，并且整个分布式系统内不会产生ID碰撞(由数据中心ID和机器ID作区分)，并且效率较高，经测试，SnowFlake每秒能够产生26万ID左右。
 */

const (
	//机器id所占的位数
	workerIdBits int64 = 5
	//数据标识id所占的位数
	dataCenterIdBits int64 = 5
	//支持的最大机器ID，结果是31 (这个移位算法可以很快的计算出几位二进制数所能表示的最大十进制数)
	maxWorkerId int64 = ^(-1 << uint(workerIdBits))
	//最大数据标识ID
	maxDataCenterId = maxWorkerId
	//序列在id中占的位数
	sequenceBits int64 = 12
	//机器ID向左移12位
	workerIdShift = sequenceBits
	//数据标识id向左移17位(12+5)
	dataCenterIdShift = sequenceBits + dataCenterIdBits
	//时间截向左移22位(5+5+12)
	timestampLeftShift = sequenceBits + workerIdBits + dataCenterIdBits
	//生成序列的掩码，这里为4095 (0b111111111111=0xfff=4095)
	sequenceMask int64 = ^(-1 << uint(sequenceBits))
)

type SnowFlake struct {
	epoch         int64 //开始的时间戳,毫秒级
	workerId      int64 //工作ID (0-31)
	dataCenterId  int64 //数据中心ID (0-31)
	sequence      int64 //毫秒内序列 (0-4095)
	lastTimestamp int64 //上次生成ID的时间戳,毫秒级
}

type SnowId struct {
	RawId int64
}

func SnowRaw(raw int64) *SnowId {
	return &SnowId{
		RawId: raw,
	}
}

//创建一个 snowid 生成器
//epoch: 开始产生ID的纪元时间戳
//workerId: 工作ID标识
//dataCenterId: 数据ID标识
func NewShowFlake(epoch, workerId, dataCenterId int64) (*SnowFlake, error) {
	if workerId > maxWorkerId || workerId < 0 {
		return nil, fmt.Errorf("worker Id can't be greater than %d or less than 0", maxWorkerId)
	}

	if dataCenterId > maxDataCenterId || dataCenterId < 0 {
		return nil, fmt.Errorf("datacenter Id can't be greater than %d or less than 0", maxWorkerId)
	}

	return &SnowFlake{
		epoch:         epoch,
		workerId:      workerId,
		dataCenterId:  dataCenterId,
		lastTimestamp: -1,
	}, nil
}

//得到下一个ID
func (s *SnowFlake) NextId() (*SnowId, error) {
	currentTimestamp := getTimeMilliSecond()
	//如果当前时间小于上一次ID生成的时间戳，说明系统时钟回退过这个时候应当抛出异常
	if currentTimestamp < s.lastTimestamp {
		return nil, fmt.Errorf("Clock moved backwards.  Refusing to generate id for %d milliseconds", s.lastTimestamp-currentTimestamp)
	}

	//如果是同一时间生成的，则进行毫秒内序列
	if currentTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		//毫秒内序列溢出
		if s.sequence == 0 {
			//阻塞到下一个毫秒,获得新的时间戳
			currentTimestamp = s.tilNextMillis()
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = currentTimestamp

	rawId := ((currentTimestamp - s.epoch) << uint(timestampLeftShift)) |
		(s.dataCenterId << uint(dataCenterIdShift)) |
		(s.workerId << uint(workerIdShift)) |
		s.sequence
	return SnowRaw(rawId), nil
}

func (s *SnowFlake) tilNextMillis() int64 {
	time.Sleep(time.Millisecond)
	return getTimeMilliSecond()
}

//得到毫秒时间戳
func getTimeMilliSecond() int64 {
	return time.Now().UnixNano() / 1e6
}
