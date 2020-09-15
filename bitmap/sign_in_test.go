package bitmap

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestSign(t *testing.T)  {
	sign := NewSign(userKey())

	// 判断今天是否已签到
	assert.Equal(t,false,sign.IsSign())

	// 今天签到
	err := sign.Sign()
	assert.Nil(t,err)

	// 判断今天是否已签到
	assert.Equal(t,true,sign.IsSign())

	// 判断总共的签到日期
	count , err := sign.Count()
	assert.Nil(t,err)
	assert.Equal(t,uint(1),count)

	// 得到最开始签到的日期
	startTime ,err := sign.StartDay()
	assert.Nil(t,err)
	assert.Equal(t,time.Now().Format(`2006-01-02`),startTime.Format(`2006-01-02`))

	// 得到最后一天签到日期
	endTime ,err := sign.EndDay()
	assert.Nil(t,err)
	assert.Equal(t,time.Now().Format(`2006-01-02`),endTime.Format(`2006-01-02`))
}

func userKey() string {
	rand2 := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(rand2.Int())
}

func mustNil( err error)  {
	if err != nil {
		panic(err)
	}
}