package repository

import (
	"bytes"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

// go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1
// go tool pprof repository.test cpu.out
// go tool pprof -http=:8091 repository.test cpu.out
func BenchmarkMultipleMap(b *testing.B) {

	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	dbConnect := NewMockIDatabase(ctrl)
	dbConnect.EXPECT().Exec(gomock.Any(), gomock.Any()).AnyTimes()

	collectorM := NewCollector(dbConnect)

	for i := 0; i < b.N; i++ {
		collectorM.HandleEvent(12, "sample")
	}
}

func BenchmarkSingleMap(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	dbConnect := NewMockIDatabase(ctrl)
	dbConnect.EXPECT().Exec(gomock.Any(), gomock.Any()).AnyTimes()

	collectorSM := NewCollectorSW(dbConnect)
	for i := 0; i < b.N; i++ {
		collectorSM.HandleEvent(12, "sample")
	}
}

func BenchmarkSyncMap(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	dbConnect := NewMockIDatabase(ctrl)
	dbConnect.EXPECT().Exec(gomock.Any(), gomock.Any()).AnyTimes()

	collectorSNG := NewCollectorSNG(dbConnect)
	for i := 0; i < b.N; i++ {
		collectorSNG.HandleEvent(12, "sample")
	}
}

func BenchmarkConcat(b *testing.B) {
	var str string
	for n := 0; n < b.N; n++ {
		str += "x"
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); str != s {
		b.Errorf("unexpected result; got=%s, want=%s", str, s)
	}
}

func BenchmarkBuffer(b *testing.B) {
	var buffer bytes.Buffer
	for n := 0; n < b.N; n++ {
		buffer.WriteString("x")
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); buffer.String() != s {
		b.Errorf("unexpected result; got=%s, want=%s", buffer.String(), s)
	}
}

func BenchmarkCopy(b *testing.B) {
	bs := make([]byte, b.N)
	bl := 0

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bl += copy(bs[bl:], "x")
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); string(bs) != s {
		b.Errorf("unexpected result; got=%s, want=%s", string(bs), s)
	}
}
