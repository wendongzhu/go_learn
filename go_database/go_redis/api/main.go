package main

// go get -u github.com/go-redis/redis

//普通连接
//var redisDb *redis.Client
//
//func initClient()(err error) {
//	redisDb = redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//		Password: "",
//		DB: 0,
//	})
//	_, err = redisDb.Ping().Result()
//	if err != nil {
//		return err
//	}
//	return nil
//}
