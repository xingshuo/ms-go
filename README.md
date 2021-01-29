# golang工具箱

## 工具列表
```
deadlock: 死锁检测
    功能测试: 
       go test -v
    性能测试: 
       go test -bench=. -run=none
    性能报告:
       测试环境:
         Intel(R) Xeon(R) Platinum 8255C CPU @ 2.50GHz
         cpu逻辑核数: 8核
       测试数据:
         BenchmarkRawMutex-8         	84095577	        13.9 ns/op
         BenchmarkMutexDisable-8     	72938457	        16.3 ns/op
         BenchmarkMutexEnable-8      	 2327222	       514 ns/op
         BenchmarkRawRWMutex-8       	42695430	        28.0 ns/op
         BenchmarkRWMutexDisable-8   	36915428	        32.1 ns/op
         BenchmarkRWMutexEnable-8    	 2330823	       519 ns/op
```