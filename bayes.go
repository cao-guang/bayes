package bayes

import (
	"math"
)

//函数说明:将不同来源和文本的实验样本数据整合成一个[][]string
//Parameters:
//items - 整理的样本数据集
//train_category - 训练类别标签向量
//Returns:
//listPosts - 返回[][]string 类型的样本数据集
//train_category - 原封不动返回 训练类别标签向量（至于为什么要这样 只是为了封装好看而已 手动滑稽）

func LoadDataNB(train_category []int,items...[]string) ([][]string,[]int) {
	listPosts := [][]string{}
	for _,set:=range items{
		listPosts = append(listPosts,set)
	}
	return listPosts,train_category
}


//获取多个[]string求并集

//函数说明:将切分的实验样本词条整理成不重复的词条列表，也就是词汇表
//Parameters:
//items - 整理的样本数据集
//Returns:
//r.List() - 返回不重复的词条列表，也就是词汇表
func CreateUnionList(items ...[]string) []string{
	s := New()
	r:=s.Duplicate() //创建副本
	// 获取并集
	for _, set := range items {
		for _,e := range set {
			r.Add(e)
		}
	}
	return r.SortedList()
}
//创建词向量
//函数说明:根据union_arr词汇表，将input_arr向量化，向量的每个元素为1或0
//Parameters:
//union_arr - CreateUnionList返回的列表
//input_arr - 切分的词条列表
//Returns:
//returnVec - 文档向量,词集模型

func Set2Vec(union_arr []string,input_arr []string) []int{
	union_arr_len := len(union_arr)
	returnVec := []int{}
	for i:=0;i<union_arr_len;i++{
		returnVec = append(returnVec, 0) //向量初始化0
	}
	for _,v:=range input_arr{
		for k,v1:=range union_arr{
			if v == v1{
				returnVec[k] = 1 //训练词集在模型中出现的次数【伯努利模型】将重复的词语都视为其只出现1次 需要用到【多项式模型】的话直接+1
			}
		}
	}
	return returnVec
}


//多项式模型
//函数说明:朴素贝叶斯分类器训练函数
//Parameters:
//train_nb_arr - 训练文档矩阵，即Set2Vec返回的returnVec构成的矩阵
//train_category - 训练类别标签向量，即LoadDataNB返回的train_category
//Returns:
//p_Vect0 - 正常类概率数组
//p_Vect1 - 敏感词类概率数组
//p_pro - 文档属于敏感词的概率
func MultinomialNB(train_nb_arr [][]int,train_category []int)([]float64,[]float64,float64){
    numTrainDocs := len(train_nb_arr) //文档数目
    numDoc := len(train_nb_arr[0]) //词汇表词数目
	p_pro :=float64(sum(train_category))/float64(len(train_category)) //文档属于敏感词类的概率
	p_Num0 :=ones(numDoc)
	p_Num1 :=ones(numDoc) // 初始化为1
	//使用拉普拉斯平滑解决零概率问题；
	//对乘积结果取自然对数避免下溢出问题，采用自然对数进行处理不会有任何损失。
	p_lap0:=2
	p_lap1:=2 //分母初始化为2 ,拉普拉斯平滑
	for i:=0;i<numTrainDocs;i++{
		if train_category[i]==0{ //统计属于正常词汇类的条件概率所需的数据，即P(w0|0),P(w1|0),P(w2|0)···
			p_Num0 = plus_arr(p_Num0,train_nb_arr[i])
			p_lap0 +=sum(train_nb_arr[i])
		}else{ //统计属于敏感词的条件概率所需的数据，即P(w0|1),P(w1|1),P(w2|1)···
			p_Num1 = plus_arr(p_Num1,train_nb_arr[i])
			p_lap1 +=sum(train_nb_arr[i])
		}
	}
	p_Vect0 := division_arr(p_Num0,p_lap0)
	p_Vect1 := division_arr(p_Num1,p_lap1)
	return p_Vect0, p_Vect1, p_pro
}

//函数说明:朴素贝叶斯分类器分类函数
//Parameters:
//doc - 待分类的词条数组
//p_v0 - 正常类的条件概率数组
//p_v1 - 敏感词类的条件概率数组
//p_p - 文档属于敏感词类的概率
//Returns:
//0 - 属于正常类
//1 - 属于敏感词类
func ClassNB(doc []int,p_v0 []float64,p_v1 []float64,p_p float64) int{
	p0 := sum_f(multiplication_arr(doc,p_v0))+math.Log(1.0-p_p)
	p1 := sum_f(multiplication_arr(doc,p_v1))+math.Log(p_p)
	if p1>p0{
		return 1
	}else{
		return 0
	}
	return 0
}