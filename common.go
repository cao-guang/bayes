package bayes

import (
	"sync"
	"sort"
	"strconv"
	"fmt"
	"math"
)

type Set struct {
	sync.RWMutex
	m map[string]bool
}

// 新建集合对象
// 可以传入初始元素
func New(items ...string) *Set {
	s := &Set{
		m: make(map[string]bool, len(items)),
	}
	s.Add(items...)
	return s
}

// 创建副本
func (s *Set) Duplicate() *Set {
	s.Lock()
	defer s.Unlock()
	r := &Set{
		m: make(map[string]bool, len(s.m)),
	}
	for e := range s.m {
		r.m[e] = true
	}
	return r
}

// 添加元素
func (s *Set) Add(items ...string) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.m[v] = true
	}
}

// 删除元素
func (s *Set) Remove(items ...string) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		delete(s.m, v)
	}
}

// 判断元素是否存在
func (s *Set) Has(items ...string) bool {
	s.RLock()
	defer s.RUnlock()
	for _, v := range items {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// 统计元素个数
func (s *Set) Count() int {
	s.Lock()
	defer s.Unlock()
	return len(s.m)
}

// 清空集合
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[string]bool{}
}

// 空集合判断
func (s *Set) Empty() bool {
	s.Lock()
	defer s.Unlock()
	return len(s.m) == 0
}

// 获取元素列表（无序）
func (s *Set) List() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// 获取元素列表（有序）
func (s *Set) SortedList() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	sort.Strings(list)
	return list
}

// 并集
// 获取 s 与参数的并集，结果存入 s
func (s *Set) Union(sets ...*Set) {
	// 为了防止多例程死锁，不能同时锁定两个集合
	// 所以这里没有锁定 s，而是创建了一个临时集合
	r := s.Duplicate()
	// 获取并集
	for _, set := range sets {
		set.Lock()
		for e := range set.m {
			r.m[e] = true
		}
		set.Unlock()
	}
	// 将结果转入 s
	s.Lock()
	defer s.Unlock()
	s.m = map[string]bool{}
	for e := range r.m {
		s.m[e] = true
	}
}

// 并集（函数）
// 获取所有参数的并集，并返回
func Union(sets ...*Set) *Set {
	// 处理参数数量
	if len(sets) == 0 {
		return New()
	} else if len(sets) == 1 {
		return sets[0]
	}
	// 获取并集
	r := sets[0].Duplicate()
	for _, set := range sets[1:] {
		set.Lock()
		for e := range set.m {
			r.m[e] = true
		}
		set.Unlock()
	}
	return r
}

// 差集
// 获取 s 与所有参数的差集，结果存入 s
func (s *Set) Minus(sets ...*Set) {
	// 为了防止多例程死锁，不能同时锁定两个集合
	// 所以这里没有锁定 s，而是创建了一个临时集合
	r := s.Duplicate()
	// 获取差集
	for _, set := range sets {
		set.Lock()
		for e := range set.m {
			delete(r.m, e)
		}
		set.Unlock()
	}
	// 将结果转入 s
	s.Lock()
	defer s.Unlock()
	s.m = map[string]bool{}
	for e := range r.m {
		s.m[e] = true
	}
}

// 差集（函数）
// 获取第 1 个参数与其它参数的差集，并返回
func Minus(sets ...*Set) *Set {
	// 处理参数数量
	if len(sets) == 0 {
		return New()
	} else if len(sets) == 1 {
		return sets[0]
	}
	// 获取差集
	r := sets[0].Duplicate()
	for _, set := range sets[1:] {
		for e := range set.m {
			delete(r.m, e)
		}
	}
	return r
}

// 交集
// 获取 s 与其它参数的交集，结果存入 s
func (s *Set) Intersect(sets ...*Set) {
	// 为了防止多例程死锁，不能同时锁定两个集合
	// 所以这里没有锁定 s，而是创建了一个临时集合
	r := s.Duplicate()
	// 获取交集
	for _, set := range sets {
		set.Lock()
		for e := range r.m {
			if _, ok := set.m[e]; !ok {
				delete(r.m, e)
			}
		}
		set.Unlock()
	}
	// 将结果转入 s
	s.Lock()
	defer s.Unlock()
	s.m = map[string]bool{}
	for e := range r.m {
		s.m[e] = true
	}
}

// 交集（函数）
// 获取所有参数的交集，并返回
func Intersect(sets ...*Set) *Set {
	// 处理参数数量
	if len(sets) == 0 {
		return New()
	} else if len(sets) == 1 {
		return sets[0]
	}
	// 获取交集
	r := sets[0].Duplicate()
	for _, set := range sets[1:] {
		for e := range r.m {
			if _, ok := set.m[e]; !ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// 补集
// 获取 s 相对于 full 的补集，结果存入 s
func (s *Set) Complement(full *Set) {
	r := full.Duplicate()
	s.Lock()
	defer s.Unlock()
	// 获取补集
	for e := range s.m {
		delete(r.m, e)
	}
	// 将结果转入 s
	s.m = map[string]bool{}
	for e := range r.m {
		s.m[e] = true
	}
}

// 补集（函数）
// 获取 sub 相对于 full 的补集，并返回
func Complement(sub, full *Set) *Set {
	r := full.Duplicate()
	sub.Lock()
	defer sub.Unlock()
	for e := range sub.m {
		delete(r.m, e)
	}
	return r
}

//数组加运算
func sum(arr []int) int{
	var i int = 0
	for _,v:=range arr{
		i = i+v
	}
	return i
}


//数组加运算 float64
func sum_f(arr []float64) float64  {
	var i float64 = 0
	for _,v:=range arr{
		i = i+v
	}
	return i
}
//取1位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", value), 64)
	return value
}
//初始化一个空数组并用1占位
func ones(inlen int)[]int{
	var a []int
	for i:=0;i<inlen;i++{
		a = append(a,1)
	}
	return a
}


//两个数组值相加运算
func plus_arr(a_arr []int,b_arr []int) []int {
	arr_len :=len(a_arr)
	var c_arr []int
	for i:=0;i<arr_len;i++{
		c_arr = append(c_arr,a_arr[i] + b_arr[i])
	}
	return c_arr
}

//数组内值相除运算求对数值
func division_arr(a_arr []int,b_arr int)[]float64{
	arr_len :=len(a_arr)
	var c_arr []float64
	for i:=0;i<arr_len;i++{
		c_arr = append(c_arr,math.Log(float64(a_arr[i]) / float64(b_arr))) //取对数，防止下溢出
	}
	return c_arr
}
//数组内乘法运算
func multiplication_arr(a_arr []int,b_arr []float64) []float64 {
	arr_len :=len(a_arr)
	var c_arr []float64
	for i:=0;i<arr_len;i++{
		c_arr = append(c_arr,float64(a_arr[i]) * b_arr[i])
	}
	return c_arr
}