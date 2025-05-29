local FormulaParser = {}
FormulaParser.__index = FormulaParser

local unpack = unpack or table.unpack

-- 运算符优先级表（调整!的优先级高于比较运算符）
local precedence = {
    ["u-"] = 8,  -- 一元负号（最高）
    ["!"] = 7,   -- 逻辑非
    ["^"] = 6,   -- 幂运算
    ["*"] = 5, ["/"] = 5, ["%"] = 5,  -- 乘,除,取模
    ["+"] = 4, ["-"] = 4,  -- 加减
    [">"] = 3, ["<"] = 3, [">="] = 3, ["<="] = 3,  -- 比较运算
    ["=="] = 3, ["!="] = 3,
    ["&&"] = 2,  -- 逻辑与
    ["||"] = 1   -- 逻辑或（最低）
}

function FormulaParser.new()
    local self = setmetatable({}, FormulaParser)
    self._cache = {}
    return self
end

-- 清空缓存
function FormulaParser:clearCache()
    self._cache = {}
end

-- 解析函数参数列表（完整支持嵌套表达式）
local function parseArguments(expr, startPos)
    local args = {}
    local currentArg = ""
    local parenLevel = 0
    local i = startPos  -- 参数起始索引
    
    while i <= #expr do
        local c = expr:sub(i, i)
        
        if c == '(' then   -- 参数里面有表达式
            parenLevel = parenLevel + 1
            currentArg = currentArg .. c
        elseif c == ')' then
            if parenLevel == 0 then
                -- 结束参数解析
                if currentArg ~= "" then
                    table.insert(args, currentArg)
                end
                return args, i
            else
                parenLevel = parenLevel - 1
                currentArg = currentArg .. c
            end
        elseif c == ',' and parenLevel == 0 then
            -- 参数分隔符
            table.insert(args, currentArg)   --把完整的参数插入到args
            currentArg = ""
        else
            currentArg = currentArg .. c   -- 拼接普通参数
        end
        
        i = i + 1
    end
    
    error("Unclosed function arguments")
end

function FormulaParser:infixToPostfix(expr)
    local output = {}
    local stack = {}
    local i = 1   -- 全局的解析索引位置
    local n = #expr      -- 获取表达式的长度
    
    local function peek()
        return stack[#stack]    -- 获取栈顶元素，但不弹出
    end
    
    local function isDigit(c)
        return c:match("[0-9%.]") -- 匹配连续的数字
    end
    
    local function isLetter(c)
        return c:match("[a-zA-Z_]")  -- 匹配字母以及下划线
    end
    
    local function isOperator(c)    
        return c:match("[%+%-%*/%%%^><!=&|]")   -- 匹配操作符
    end

    local function getNumber()     -- 解析数值，包括1e-3这种科学计数法
        local start = i
        while i <= n and (isDigit(expr:sub(i, i)) or expr:sub(i, i):lower() == 'e') do
            if expr:sub(i, i):lower() == 'e' then
                i = i + 1
                if expr:sub(i, i) == '+' or expr:sub(i, i) == '-' then
                    i = i + 1
                end
            end
            i = i + 1
        end
        local num = expr:sub(start, i - 1)
        i = i - 1
        return tonumber(num)
    end
    
    local function getIdentifier()  -- 解析变量名称，必须字符开头，可以带有数字
        local start = i
        while i <= n and (isLetter(expr:sub(i, i)) or isDigit(expr:sub(i, i))) do
            i = i + 1
        end
        local id = expr:sub(start, i - 1)
        i = i - 1
        return id
    end
    
    local function getOperator(expr, i)
        local start = i
        local c = expr:sub(i, i)  -- 取出第一个字符
        i = i + 1
        
        -- 多字符运算符  返回第二个字符的位置
        if c == '=' and expr:sub(i, i) == '=' then
            return "==", i
        elseif c == '!' and expr:sub(i, i) == '=' then
            return "!=", i
        elseif c == '>' and expr:sub(i, i) == '=' then
            return ">=", i
        elseif c == '<' and expr:sub(i, i) == '=' then
            return "<=", i
        elseif c == '&' and expr:sub(i, i) == '&' then
            return "&&", i
        elseif c == '|' and expr:sub(i, i) == '|' then
            return "||", i
        elseif c == '!' then
            -- 确保!是一元运算符
            if #output > 0 and 
                (type(output[#output]) == "number" or 
                output[#output] == ")" or
                (type(output[#output]) == "string" and not precedence[output[#output]])) then
                error("Invalid use of ! operator")
            end
            return "!", start 
        elseif c == '-' then
            if #output == 0  or ( #stack > 0 and peek() == '(' ) then  -- 表达式第一个符号是负号， 括号表达式第一个是负号
                return "u-", start
            end
        end
        
        return c, start
    end
    
    -- 执行解析
    while i <= n do
        local c = expr:sub(i, i)  --取出第一个字符
        
        if c:match("%s") then  -- 去掉空格
            i = i + 1
        elseif isDigit(c) then  -- 如果是立即数，append到output
            table.insert(output, getNumber())  
            i = i + 1
        elseif isLetter(c) then  -- 如果是字符，说明是变量
            local id = getIdentifier() -- 解析变量
            if i <= n and expr:sub(i+1, i+1) == '(' then   -- 变量名称后面是(, 说明是函数调用
                -- 处理函数调用
                local args, newPos = parseArguments(expr, i+2) -- 解析参数
                table.insert(output, {
                    func = id,   -- 函数名
                    args = args  -- 参数
                })
                i = newPos
            else
                table.insert(output, id)  -- 普通变量，append到output
            end
            i = i + 1
        elseif c == '(' then   -- 括号表达式
            table.insert(stack, c)  -- 左括号压入临时栈stack
            i = i + 1
            -- 特别处理括号后紧跟负号的情况，如"(-1)"
            if expr:sub(i, i) == '-' then
                local op, newPos = getOperator(expr, i)
                if op == "u-" then
                    table.insert(stack, op)
                    i = newPos + 1
                end
            end
        elseif c == ')' then  -- 括号表达式结束
            while #stack > 0 and peek() ~= '(' do  -- stack保存了非空的括号表达式
                table.insert(output, table.remove(stack))  -- 取出来，append到output
            end
            if #stack == 0 then   -- 没有匹配到左括号
                error("Mismatched parentheses")
            end
            table.remove(stack)  -- 移除括号表达式，包括空的括号表达式：()
            i = i + 1
        elseif isOperator(c) then  -- 解析操作符
            local op, newPos = getOperator(expr, i)
            i = newPos + 1
        
            while #stack > 0 and peek() ~= '(' and   -- 栈不为空 并且 栈顶不是左括号 
                  (precedence[peek()] or 0) >= (precedence[op] or 0) do  -- 并且 栈顶运算符的优先级 ≥ 当前运算符的优先级
                table.insert(output, table.remove(stack))   -- 将栈顶运算符弹出并添加到输出队列
            end
            table.insert(stack, op) 
            --[[
                举例：
                输入表达式：a + b * c - d
                执行到 - 时的状态：
                    Stack: [ + * ]
                    Output: [a b c]
                处理过程：
                    比较 - 和 *：* 优先级更高 → 弹出 *
                    比较 - 和 +：+ 优先级更高 → 弹出 +
                    将 - 压栈
                    最终输出队列：[a b c * + d -]
            ]]
        else
            error("Unknown character: " .. c)
        end
    end
    
    while #stack > 0 do
        local op = table.remove(stack)
        if op == '(' then
            error("Mismatched parentheses")
        end
        table.insert(output, op)
    end
    
    return output
end

function FormulaParser:evaluatePostfix(postfix, variables)
    variables = variables or {}
    local stack = {}
    
    for _, token in ipairs(postfix) do
        if type(token) == "number" then
            table.insert(stack, token)
        elseif type(token) == "string" and not precedence[token] then
            if variables[token] == nil then
                error("Undefined variable: " .. token)
            end
            table.insert(stack, variables[token])
        elseif type(token) == "table" and token.func then
            local func = variables[token.func]
            if not func or type(func) ~= "function" then
                error("Undefined function: " .. token.func)
            end
            
            -- 递归计算所有参数
            local args = {}
            for _, argExpr in ipairs(token.args) do
                table.insert(args, self:evaluate(argExpr, variables))
            end
            
            table.insert(stack, func(unpack(args)))
        else
            -- 运算符处理
            if token == "u-" then
                local a = table.remove(stack)
                if not a then error("Insufficient operands for unary -") end
                table.insert(stack, -a)
            elseif token == "!" then
                local a = table.remove(stack)
                -- 确保正确处理各种类型的"假"值
                local is_true = function()
                    if (type(a) == "number") then 
                        return a ~= 0 
                    elseif (type(a) == "boolean") then 
                        return a 
                    elseif (a ~= nil) then 
                        return true 
                    end
                    return false
                end
                table.insert(stack, is_true() and 0 or 1)
            else
                local b = table.remove(stack)
                local a = table.remove(stack)
                
                if token == "+" then table.insert(stack, a + b)
                elseif token == "-" then table.insert(stack, a - b)
                elseif token == "*" then table.insert(stack, a * b)
                elseif token == "/" then 
                    -- 添加除数是否为0的检查
                    if b == 0 then
                        error("Division by zero")
                    end
                    table.insert(stack, a / b)
                elseif token == "%" then 
                    -- 取模运算也要检查除数
                    if b == 0 then
                        error("Modulo by zero")
                    end
                    table.insert(stack, a % b)
                elseif token == "^" then table.insert(stack, a ^ b)
                elseif token == "&&" then table.insert(stack, (a ~= 0 and b ~= 0) and 1 or 0)
                elseif token == "||" then table.insert(stack, (a ~= 0 or b ~= 0) and 1 or 0)
                elseif token == "==" then table.insert(stack, a == b and 1 or 0)
                elseif token == "!=" then table.insert(stack, a ~= b and 1 or 0)
                elseif token == ">" then table.insert(stack, a > b and 1 or 0)
                elseif token == "<" then table.insert(stack, a < b and 1 or 0)
                elseif token == ">=" then table.insert(stack, a >= b and 1 or 0)
                elseif token == "<=" then table.insert(stack, a <= b and 1 or 0)
                else error("Unknown operator: " .. token)
                end
            end
        end
    end
    
    if #stack ~= 1 then
        error("Invalid expression")
    end
    
    return stack[1]
end

function FormulaParser:evaluate(expr, variables)
    -- 移除所有空格
    local normalizedExpr = expr:gsub("%s+", "")

    if not self._cache[normalizedExpr] then
        self._cache[normalizedExpr] = self:infixToPostfix(normalizedExpr)
    end
    return self:evaluatePostfix(self._cache[normalizedExpr], variables)
end

return  FormulaParser

