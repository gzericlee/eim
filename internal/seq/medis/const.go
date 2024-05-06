package medis

const Unit = 1e5                  // 公共数字单位
const ListKey = "seq"             // 存储已生成数据的键名
const MaxKey = "mdx"              // 存储当前最大值的键名
const Capacity = int64(50 * Unit) // 信道容量
const Percent = float64(0.7)      // 信道容量阈值比
const Multiple = int64(5)         // 用于计算补充量的倍数
const RandMax = 250               // 随机值右闭值

var Freedom = 0 // 作为 Goroutine 锁
