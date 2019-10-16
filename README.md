# bayes
golang 贝叶斯算法训练过滤敏感词，用了多项式模型和伯努利模型


添加引用
```go
import(
"github.com/cao-guang/bayes"
)
```
初始化调用
```go
func main()  {
        //准备训练样本数据，这里就准备了3组，训练样本数据越多，效果你懂的
        aa := []string{"狗b", "我艹", "愚蠢", "装你大爷", "五毛们"}
    	bb := []string{"法L功", "sb", "杀b", "你吗b", "你妈的", "五毛们", "愚蠢"}
    	cc := []string{"爱国", "奉献", "勤劳", "努力", "奋斗", "汗水"}
    	classVec := []int{1, 1, 0} // 1 代表敏感词 0 代表正常类词汇
    	listPosts, listClasses := bayes.LoadDataNB(classVec, aa, bb, cc) //考虑到数据来源的多样性，需要加载样本数据返回需要的格式
        myVocabList := bayes.CreateUnionList(aa, bb, cc) //求样本数据并集
        //对得到的并集样本数据进行向量化 ps：这里就不进一步封装了，便于大家理解
        trainMat := [][]int{}
        for _, v := range listPosts {
            trainMat = append(trainMat, bayes.Set2Vec(myVocabList, v))
        }
        p0V, p1V, pAb := bayes.MultinomialNB(trainMat, listClasses) //贝叶斯算法计算概率
        test1 := []string{"装你大爷"}
        Doc1 := bayes.Set2Vec(myVocabList, test1)
        result_1 := bayes.ClassNB(Doc1, p0V, p1V, pAb)
        fmt.Println("训练结果：", result_1)
        test2 := []string{"爱国", "奉献"}
        Doc2 := bayes.Set2Vec(myVocabList, test2)      //得到需要训练的词的向量
        result_2 := bayes.ClassNB(Doc2, p0V, p1V, pAb) //分类器分类
        fmt.Println("训练结果：", result_2)
}
```