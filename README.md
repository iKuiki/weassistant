# weassistant

针对微信群聊的群主助手，需要搭配wegate使用。weassistant通过mqant定义的mqtt规则连接到wegate，处理wegate传输的事件。对外，weassistant提供http api访问，http由iris框架驱动