# Bug Server

游戏服务器

架构：

- HttpServer
  - Login
    - 登陆、注册、邮箱验证API
  - Hall
    - 玩家各种数据的传输（等级、成就等）
  - Room
    - 对战前准备（选择人物、天赋等）
    - 战斗后结算（统计游戏数据、修改玩家数据）
- GameServer
  - Fighting
    - 使用帧同步技术，用UDP传输战斗数据