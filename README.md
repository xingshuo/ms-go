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
         BenchmarkRawMutex-8              	85776766	        13.9 ns/op
         BenchmarkMutexDisable-8          	73861753	        16.3 ns/op
         BenchmarkMutexEnable-8           	 2893243	       412 ns/op
         BenchmarkRawRWMutexLock-8        	42231187	        28.1 ns/op
         BenchmarkRawRWMutexRLock-8       	99902935	        11.7 ns/op
         BenchmarkRWMutexLockDisable-8    	36746524	        32.0 ns/op
         BenchmarkRWMutexRLockDisable-8   	73733905	        16.6 ns/op
         BenchmarkRWMutexLockEnable-8     	 2548990	       478 ns/op
         BenchmarkRWMutexRLockEnable-8    	 3657556	       325 ns/op

```