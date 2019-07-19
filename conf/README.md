# config

## Conf

config包中定义了weassistant的配置接口，weassistant运行时需要传入一个Conf接口来提供配置，出于最大化解耦的目的，我们约定weassistant读取的所有配置都应当从接口Conf中读取

## ConfExtra

config包还提供了ConfExtra接口来给weassistant提供初始化好的服务。与Conf一样，出于解耦目的，weassistant所调用的所有服务都应当从confExtra中获取

## conf

conf是config包提供的一个Conf的实现，实现了从json创建Conf的功能

## confExtra

confExtra是config包提供的一个ConfExtra的实现，允许根据Conf创建ConfExtra的实例