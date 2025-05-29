extern "C" {
    #include "stdio.h"
    #include "lua.h"
    #include "lualib.h"
    #include "lauxlib.h"
    #include <sys/time.h>
}
#include <iostream>
#include "luabase.hpp"

int main() {

    struct timeval start, end;
    long long mtime, seconds, useconds;

    if (initLua() < 0 ) {
        return -1;
    }
    std::cerr << "初始化lua环境 ok." << std::endl;

    // std::string expr = "sin(x)*cos(y) + sqrt(z^2)*tan(z)";
    // std::map<std::string, double> args = {
    //     {"x", 0.5}, {"y", 0.5222}, {"z", 0.335}
    // };
    std::string expr = "(x1+x2+x3+x4+x1*x2)*4 + (y1+y2-y3+y4)/4 + (x3/z1+z2+z3*z4)^2";
    std::map<std::string, double> args = {
        {"x1", 0.5}, {"x2", 0.5222}, {"x3", 0.335}, {"x4", 0.545},
        {"y1", 1212}, {"y2", 45454}, {"y3", 3232}, {"y4", 5656546565},
        {"z1", 0.2}, {"z2", 43434},  {"z3", 0.5443}, {"z4", 1233.6577}
    };
    double ret = execLua(expr,args);
    if (ret < 0 ) { 
        deinitLua();
        return -1;
    }
    std::cerr << "公式: " << expr << "  计算结果: "<< ret << std::endl;

    // ---- 性能测试 ----
    // 记录开始时间
    gettimeofday(&start, NULL);
    int i=0;
    for(i=0;i<100000;i++){
        execLua(expr,args); 
    }
    // 记录结束时间
    gettimeofday(&end, NULL);
    mtime = end.tv_sec*1000000-start.tv_sec*1000000 + end.tv_usec-start.tv_usec;
    printf("100000次执行总耗时: %lld 微秒\n", mtime);
    printf("平均耗时: %.3f 微秒\n", double(mtime)/100000);

    deinitLua();
    std::cerr << "释放lua环境 ok."  << std::endl;
    return 0;
}