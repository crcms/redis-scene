package bitmap

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/uniplaces/carbon"
	"redis-scene/server"
	"strconv"
	"time"
)

// 位图的Redis基本用法
// SETBIT key offset value 				// 设置位图的偏移值
// GETBIT key offset 					// 获取当前偏移值
// BITCOUNT key [start end] 			// 统计有值位
// BITPOS key bit [start] [end]			// 获取指定位最开始位置，如 BITPOS key 1 0 365 | BITPOS key 0 0 365

type sign struct {
	userId string
	key string
}

func NewSign(userId string) *sign  {
	return &sign{
		userId:userId,
		key: `sign` + `:` +strconv.Itoa(carbon.Now().Year()) +`:` + userId,
	}
}

// 签到 以offset设置，允许重复签到，不会影响结果lar
func (s sign) Sign() error {
	if s.IsSign() {
		return nil
	}
	return server.Redis.SetBit(context.Background(),s.key,s.currentDay(),1).Err()
}

// 判断今天是否签到
func (s sign) IsSign() bool {
	r , err := server.Redis.GetBit(context.Background(),s.key,s.currentDay()).Uint64()

	if r == 1 && err == nil {
		return true
	}

	return false
}

// 当前在一年中的第N天，也是就offset使用
func (s sign) currentDay() int64 {
	startYear, err := carbon.Create(2020,1,1,0,0,0,0,`PRC`)
	if err != nil {
		panic(err)
	}

	return carbon.Now().DiffInDays(startYear,true)
}

// 一年总的签到次数
func (s sign) Count() (uint ,error) {
	v , err := server.Redis.BitCount(context.Background(),s.key,&redis.BitCount{
		Start: 0,
		End:   s.currentDay(),
	}).Uint64()
	if err != nil {
		return 0,err
	}

	return uint(v),nil
}

// 最开始的签到日期
func (s sign) StartDay() (time.Time,error)  {
	r , err := server.Redis.BitPos(context.Background(),s.key,1,0).Uint64()
	if err != nil {
		return time.Time{},err
	}

	startDay ,err:= carbon.CreateFromDate(carbon.Now().Year(),1,1,`PRC`)
	if err !=nil {
		return time.Time{},err
	}

	return startDay.AddDays(int(r)).Time,nil
}

// 最后一次签到日期
func (s sign) EndDay() (time.Time,error)  {
	r , err := server.Redis.BitPos(context.Background(),s.key,1,-1).Uint64()
	if err != nil {
		return time.Time{},err
	}

	startDay ,err:= carbon.CreateFromDate(carbon.Now().Year(),1,1,`PRC`)
	if err !=nil {
		return time.Time{},err
	}

	return startDay.AddDays(int(r)).Time,nil
}

// TODO: 统计同一时间段共多少个用户的签到？