没有参数的超组合子，又另有一个名词来描述它：Constant Applicative Forms（简称CAF）

对一个表达式进行求值（当然是惰性的）并对它在图中对应的节点进行原地更新，称为规约
可规约表达式（reducible expression，一般缩写为redex）
无法继续规约的基础表达式称为Normal form
表达式的规约顺序是很重要的—一些程序只在特定的规约顺序下停机
特殊的规约顺序永远选择最外层的redex进行规约，这叫做normal order reduction

用户定义的超组合子，规约只需要做参数替换即可
参数数量不够所以没法作为 redex 处理，是  weak head normal form（一般缩写为WHNF），这种情况下即使它的子表达式中包含redex，也不需要做任何事

WHNF: 整数 &  partial application

G-Machine
堆内存的基本单位不是字节，而是图节点。
栈里只放指向堆的地址，不放实际数据。

coreF里的超组合子会被编译成一系列G-Machine指令，大致可以分为这几种：
//
访问数据的指令，例如PushArg（访问参数用），PushGlobal（访问其他超组合子用）
在堆上构建/更新图节点的指令，如MkApp、PushInt、Update
清理栈内无用地址的Pop指令
表达控制流的Unwind指令
