extern "C" {
    #include "malloc.h"
    #include "lua.h"
    #include "lualib.h"
    #include "lauxlib.h"
}
#include <iostream>
#include <sstream> 
#include <cstring>
#include <map>

lua_State* L;

int initLua() {
    L = luaL_newstate();  // 创建新的 Lua 状态机
    if (!L) {
        std::cerr << "无法创建Lua状态机" << std::endl;
        return -1;
    }
    // 加载标准库
    luaL_openlibs(L);
    // vars 存放变量值,和基础计算库
    const char *code = "local FormulaParser = require('calc')\n"
        "parser = FormulaParser.new()\n" 
        "vars = {}\n"
        "vars.sin = math.sin\n"
        "vars.cos = math.cos\n"
        "vars.tan = math.tan\n"
        "vars.max = math.max\n"
        "vars.sqrt = math.sqrt\n"
        "vars.pi = math.pi\n"
        "vars.min = math.min\n"
        "vars.pow = math.pow\n"
        "vars['if'] = function(cond, t, f) return cond ~= 0 and t or f end\n"
        "vars.sum = function(...)\n"
        "    local s = 0\n"
        "    for _, v in ipairs({...}) do\n"
        "        s = s + v\n"
        "    end\n"
        "    return s\n"
        "end\n"
        "vars.distance = function(x1, y1, x2, y2)\n"
        "    local dx = x2 - x1\n"
        "    local dy = y2 - y1\n"
        "    return math.sqrt(dx*dx + dy*dy)\n"
        "end";

    if (luaL_loadstring(L,code) != LUA_OK) {
        fprintf(stderr, "init load错误: %s\n", lua_tostring(L, -1));
        lua_pop(L, 1);
        lua_close(L);
        return -1;
    }
    // 执行代码
    if (lua_pcall(L, 0, 0, 0) != LUA_OK) {
        fprintf(stderr, "init pcall执行错误: %s\n", lua_tostring(L, -1));
        lua_pop(L, 1);
        lua_close(L);
        return -1;
    }
    return 0;
}

int deinitLua() {
    lua_close(L);
    return 0;
}

// 执行公式计算
double execLua(std::string &expr, std::map<std::string, double> &args) {    

    std::ostringstream  lua_args;
    // // 拼接变量，保存到vars

    for (const auto &pair : args) {
        char buffer[100];
        // 以固定精度输出，避免科学计数法
        std::snprintf(buffer, sizeof(buffer), "%.15g", pair.second);
        lua_args << "vars." << pair.first << "=" << buffer << "\n";
    }

    lua_args << "expr = '" << expr << "'\n";
    lua_args << "return parser:evaluate(expr, vars)\n";

    // 加载代码
    luaL_loadstring(L,lua_args.str().c_str());

    // 执行代码（0个参数，1个返回值，无错误处理函数）
    if (lua_pcall(L, 0, 1, 0) != LUA_OK) {
        fprintf(stderr, "pcall执行错误: %s\n", lua_tostring(L, -1));
        lua_pop(L, 1);
        return -1;
    }
    double ret =  lua_tonumber(L, -1); 
    lua_pop(L, 1);
    return ret;
}
