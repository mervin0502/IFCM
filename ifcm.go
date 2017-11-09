package ifcm

import (
	"fmt"
	"math"
	"sort"
)

// func init() {
// 	log.Println("domain")
// 	fmt.Println("domain")
// 	// log.SetOutput(bytes.NewReader(b))
// }

//Interval
type Interval struct {
	Start float64
	End   float64
}

func (d *Interval) String() string {
	return fmt.Sprintf("(%v, %v )\t", d.Start, d.End)
}

//IFCM
type IFCM struct {
	data []float64
	c    int
	m    int
	t    float64
	a    float64
	max  int
	v    []float64
	u    [][]float64
}

//NewIFCM
func NewIFCM(data []float64, c int) *IFCM {
	obj := &IFCM{
		data: data,
		c:    c,
	}
	//init parameters
	obj.m = 2
	obj.t = 0.005
	obj.a = 0.85
	obj.v = obj.initS()
	obj.max = 100
	obj.u = obj.initU()
	return obj
}

//Run
func (ifcm *IFCM) Run() []*Interval {
	preU := [][]float64{}
	for i := 0; i < ifcm.max; i++ {
		ifcm.copy(&preU, ifcm.u)
		ifcm.getU()
		ifcm.getV()
		if ifcm.isStop(preU, ifcm.u) {
			break
		}
	} //for
	return ifcm.getInterval()
}

//initS
func (ifcm IFCM) initS() []float64 {
	_temp := make([]float64, len(ifcm.data))
	copy(_temp, ifcm.data)
	sort.Float64s(_temp)
	_len := len(_temp)
	_skip := _len / (ifcm.c * 2)
	var out []float64
	for i := 1; i < ifcm.c*2; i += 2 {
		out = append(out, _temp[int(i*_skip)])
	}
	return out
}

//initU
func (ifcm *IFCM) initU() [][]float64 {
	u := make([][]float64, ifcm.c)
	for i := 0; i < ifcm.c; i++ {
		u[i] = make([]float64, len(ifcm.data))
	}
	return u
}

//getU
func (ifcm *IFCM) getU() {
	getu := func(i, k int) float64 {
		var _sum float64 = 0
		iv := ifcm.v[i]
		kv := ifcm.data[k]
		var dik, djk, temp float64
		dik = (iv - kv) * (iv - kv)
		for j := 0; j < ifcm.c; j++ {
			jv := ifcm.v[j]
			djk = (jv - kv) * (jv - kv)
			if djk != 0 {
				temp = dik / djk
			} else {
				temp = 0
			}
			_sum += math.Pow(temp, float64(1/(ifcm.m-1)))
		} //for
		if _sum == 0 {
			temp = 0
		} else {
			temp = 1 / _sum
		}
		return temp + ifcm.hesitationDegree(temp)
	}
	for i := 0; i < ifcm.c; i++ {
		for j := 0; j < len(ifcm.data); j++ {
			ifcm.u[i][j] = getu(i, j)
		}
	}
}

//getV
func (ifcm *IFCM) getV() {
	getv := func(i int) float64 {
		var _sum1 float64 = 0
		var _sum2 float64 = 0
		for j := 0; j < len(ifcm.data); j++ {
			_sum1 += math.Pow(ifcm.u[i][j], float64(ifcm.m)) * ifcm.data[j]
			_sum2 += math.Pow(ifcm.u[i][j], float64(ifcm.m))
		}
		if _sum1 == 0 || _sum2 == 0 {
			return 0
		} else {
			return _sum1 / _sum2
		}
	}
	for i := 0; i < ifcm.c; i++ {
		ifcm.v[i] = getv(i)
	}
}

//isStop
func (ifcm *IFCM) isStop(preU, curU [][]float64) bool {
	for i := 0; i < ifcm.c; i++ {
		for j := 0; j < len(ifcm.data); j++ {
			if math.Abs(preU[i][j]-curU[i][j]) > ifcm.t {
				return false
			}
		}
	}
	return true
}

//getInterval
func (ifcm *IFCM) getInterval() []*Interval {
	sort.Float64s(ifcm.v)
	_mid0, _end0 := ifcm.v[0], ifcm.v[1]
	_mid1, _end1 := ifcm.v[len(ifcm.v)-2], ifcm.v[len(ifcm.v)-1]
	// var _data []float64
	_data := make([]float64, len(ifcm.data))
	copy(_data, ifcm.data)
	sort.Float64s(_data)
	var _min float64 = float64(math.MaxFloat32)
	var _max float64 = 0

	var _temp0, _temp1 []float64
	for _, _v := range _data {
		if float64(_v) <= (_mid0+_end0)/2 {
			_temp0 = append(_temp0, float64(_v))
		} else if float64(_v) > (_mid1+_end1)/2 {
			_temp1 = append(_temp1, float64(_v))
		}

		if _min > float64(_v) {
			_min = float64(_v)
		}
		if _max < float64(_v) {
			// fmt.Println(_v)
			_max = float64(_v)
		}

	} //for
	var _std float64
	_std = ifcm.std(_temp0)
	if _std == 0 {
		_std = _min / 100
	}
	_min -= _std

	_std = ifcm.std(_temp1)
	if _std == 0 {
		_std = _max / 100
	}
	_max += _std

	out := make([]*Interval, ifcm.c)
	var i int
	var _pre, _cur float64
	_pre = _min
	for i = 0; i < ifcm.c-1; i++ {
		_cur = (ifcm.v[i] + ifcm.v[i+1]) / 2
		out[i] = &Domain{_pre, _cur}
		_pre = _cur
	}
	out[i] = &Interval{_pre, _max}
	return out
}

//nomemebershipDegree
func (ifcm *IFCM) nomemebershipDegree(x float64) float64 {
	return math.Pow(1-math.Pow(x, ifcm.a), 1/ifcm.a)
}

//hesitationDegree
func (ifcm *IFCM) hesitationDegree(x float64) float64 {
	return 1 - x - ifcm.nomemebershipDegree(x)
}

//std
func (ifcm *IFCM) std(data []float64) float64 {
	var _sum float64 = 0
	_len := len(data)

	if _len == 0 {
		return 0
	}

	for _, _v := range data {
		_sum += float64(_v)
	}

	_avg := _sum / float64(_len)
	_sum = 0
	for _, _v := range data {
		_sum += math.Pow(float64(_v)-_avg, float64(2))
	}

	return math.Sqrt(_sum / float64(_len))

}

//copy
func (ifcm *IFCM) copy(_new *[][]float64, _old [][]float64) {
	var _row []float64
	*_new = make([][]float64, len(_old))
	for i := 0; i < len(_old); i++ {
		_row = _old[i]
		(*_new)[i] = make([]float64, len(_row))
		copy((*_new)[i], _row)
	}
}
