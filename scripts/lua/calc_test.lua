local FormulaParser = require("calc")

local parser = FormulaParser.new()

-- 1. 基础算术运算测试
print("\n=== 基础算术运算 ===")
local tests = {
    ["3 + 4 * 2"] = 11,
    ["(3 + 4) * 2"] = 14,
    ["10 - 3 / 2"] = 8.5,
    ["2 ^ (3 ^ 2)"] = 512, 
    ["-5 + 3"] = -2,      -- 一元负号
    ["1.2e3 + 2.5e-2"] = 1200.025,
    ["5 % 3"] = 2
}

for expr, expected in pairs(tests) do
    local result = parser:evaluate(expr)
    assert(math.abs(result - expected) < 1e-10, 
           string.format("%s = %s (期望: %s)", expr, result, expected))
    print(string.format("%-20s = %-10s (通过)", expr, result))
end

-- 2. 逻辑运算测试
local logicTests = {
    ["5 > 3 && 2 < 4"] = 1,
    ["5 % 2"] = 1,
    ["5 % (-1+4)"] = 2,
    ["!(1 == 0)"] = 1,
    ["!(1 == 1)"] = 0,
    ["(5 >= 5) || (3 < 2)"] = 1,
    ["x != y"] = 0,
    ["x == y"] = 1,
    ["0 && 1"] = 0,
    ["!(5>3)"] = 0,
    ["!(x==y)"] = 0,
    ["!0"] = 1,
    ["!(-2+2)"] = 1,
    ["!1"] = 0
}

local logicVars = {x = 5, y = 5}
print("\n=== 逻辑运算 ===")
print("变量: x=" .. logicVars.x .. " y=" .. logicVars.y)
for expr, expected in pairs(logicTests) do
    local result = parser:evaluate(expr, logicVars)
    assert(result == expected, 
           string.format("%s = %s (期望: %s)", expr, result, expected))
    print(string.format("%-20s = %-10s (通过)", expr, result))
end

-- 3. 函数调用测试
print("\n=== 函数调用 ===")
local funcTests = {
    ["max(3, 5)"] = 5,
    ["sum(1, 2, 3, 4)"] = 10,
    ["pow(2, 3) + sqrt(16)"] = 12,
    ["if(score > 60, pass, nopass)"] = "pass",
    ["nested(sin(pi/2), max(1,2))"] = 2,
    ["if(score > 90, A, if(score > 80, B, C))"] = "C"
}

local funcVars = {
    pi = math.pi,
    sin = math.sin,
    max = math.max,
    sum = function(...) 
        local s = 0
        for _, v in ipairs({...}) do s = s + v end
        return s 
    end,
    pass = "pass",
    nopass = "nopass",
    A = "A",
    B = "B",
    C = "C",
    pow = math.pow,
    sqrt = math.sqrt,
    ["if"] = function(cond, t, f) return cond ~= 0 and t or f end,
    score = 80,
    nested = function(a,b) return math.max(a,b) end
}
print(string.format("变量: score = %d", funcVars.score))

for expr, expected in pairs(funcTests) do
    local result = parser:evaluate(expr, funcVars)
    assert(result == expected or (type(result) == "number" and math.abs(result - expected) < 1e-10),
           string.format("%s = %s (期望: %s)", expr, result, expected))
    print(string.format("%-40s = %-10s (通过)", expr, result))
end

-- 4. 复杂表达式测试
print("\n=== 复杂表达式 ===")
local complexTests = {
    ["3 + sin(pi/2)*2 - sqrt(max(4,9))"] = 2,
    ["factorial(5)"] = 120,
    ["max(sin(pi/2), min(2, 3)) + pow(2, 3)"] = 10,
    ["sum(1+2, 3*4, pow(2,(1+2)))"] = 23,
    ["1+(2)+1"] = 4,
    ["1+(-2+3-3)+3"] = 2,
    ["max(sin(pi), min(1, 2))"] = 1
}

-- 先声明complexVars为空表
local complexVars = {}
complexVars.factorial = function(n)
    return n == 0 and 1 or n * complexVars.factorial(n-1)
end
complexVars.x = 1
complexVars.y = 5
complexVars.sin = math.sin
complexVars.max = math.max
complexVars.sqrt = math.sqrt
complexVars.pi = math.pi
complexVars.min = math.min
complexVars.pow = math.pow
complexVars.sum = function(...) 
    local s = 0
    for _, v in ipairs({...}) do
        s = s + v
    end
    return s
end

for expr, expected in pairs(complexTests) do
    local result = parser:evaluate(expr, complexVars)
    assert(result == expected,
           string.format("%s = %s (期望: %s)", expr, result, expected))
    print(string.format("%-40s = %-10s (通过)", expr, result))
end

-- 5. 错误处理测试
print("\n=== 错误处理 ===")
local errorTests = {
    "1 / 0",                -- 除零错误
    "1 % 0",                -- 除零错误
    "unknown_func(1)",      -- 未定义函数
    "x + y",                -- 未定义变量
    "3 + * 4",              -- 语法错误
    "sin(pi",               -- 括号不匹配
    "factorial(-1)"         -- 自定义错误
}
local errorVars = {}
for _, expr in ipairs(errorTests) do
    local ok, err = pcall(function() return parser:evaluate(expr, errorVars) end)
    print(string.format("%-30s => %s (符合预期)", expr, ok and "通过" or "错误: "..err))
    assert(not ok, "本应出错的表达式却执行成功: "..expr)
end

-- 6. 性能测试
print("\n=== 性能测试 ===")
local perfExpr = "sin(x)*cos(y) + sqrt(z^2)*tan(z)"
local perfVars = {
    x = 0.5, y = 0.3, z = 0.2,
    sin = math.sin, cos = math.cos, sqrt = math.sqrt, tan = math.tan
}

-- 预热
parser:evaluate(perfExpr, perfVars)

local start = os.clock()
local iterations = 100000
for i = 1, iterations do
    parser:evaluate(perfExpr, perfVars)
end
local elapsed = os.clock() - start
print(string.format("表达式: %s", perfExpr))
print(string.format("执行 %d 次, 总耗时: %.3f 秒", iterations, elapsed))
print(string.format("平均每次耗时: %.3f 微秒", elapsed * 1e6 / iterations))

-- 7. 缓存测试
print("\n=== 缓存测试 ===")
parser:clearCache()
local expr = "sin(0.5) + cos(0.3)"
local start1 = os.clock()
for i = 1, 1 do parser:evaluate(expr, perfVars) end
local time1 = os.clock() - start1

local start2 = os.clock()
for i = 1, 100 do parser:evaluate(expr, perfVars) end
local time2 = os.clock() - start2

time1=time1* 1e6
time2=time2* 1e6/100
print(string.format("首次解析后计算: %.3f 微秒", time1))
print(string.format("缓存后100次计算 平均: %.3f 微秒", time2))
print(string.format("性能提升: %.1f%%", (time1-time2)/time1*100))


local mathVars = {
    -- 原有数学函数...
    distance = function(x1, y1, x2, y2)
        local dx = x2 - x1
        local dy = y2 - y1
        return math.sqrt(dx*dx + dy*dy)
    end,
    sqrt = math.sqrt,  -- 确保sqrt函数可用
    x1 = 1,
    y1 = 2,
    x2 = 4,
    y2 = 6,
}


local result = parser:evaluate("distance(x1,y1,x2,y2)", mathVars)
print(string.format("\ndistance (1,2) (4,6)=%d", result))
