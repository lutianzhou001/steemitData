package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {

	start := time.Now()
	limit := New(concurrency) // New Limit 控制并发量
	max := 65802725

	ctx := context.TODO()
	cfg, err := OpenConfigFile()
	if err != nil {
		log.Fatalf("open file error:%s", err)
	}
	cl := initializeMongo(cfg, ctx)

	for i := 1; i < max; i++ {
		wg.Add(1)
		value := i
		goFunc := func() {
			fmt.Printf("start func: %d\n", value)
			// 配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
			blockMessage, err := getBlockMessage(ctx, value, cl)
			if err != nil {
				log.Println(err)
			}
			err = processingBlockMessage(ctx, blockMessage, cl)
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}
		limit.Run(goFunc)
	}

	// 阻塞代码防止退出
	wg.Wait()
	fmt.Printf("耗时: %fs", time.Now().Sub(start).Seconds())
}

type Limit struct {
	number  int
	channel chan struct{}
}

// Limit struct 初始化
func New(number int) *Limit {
	return &Limit{
		number:  number,
		channel: make(chan struct{}, number),
	}
}

// Run 方法：创建有限的 go f 函数的 goroutine
func (limit *Limit) Run(f func()) {
	limit.channel <- struct{}{}
	go func() {
		f()
		<-limit.channel
	}()
}

// WaitGroup 对象内部有一个计数器，从0开始
// 有三个方法：Add(), Done(), Wait() 用来控制计数器的数量
var wg = sync.WaitGroup{}

const (
	concurrency = 420 // 控制并发量
)

func initializeMongo(cfg Config, ctx context.Context) *mongo.Client {
	var clientOptions *options.ClientOptions
	clientOptions = options.Client().ApplyURI("mongodb://" + cfg.DataBaseLocal.Host + ":" + cfg.DataBaseLocal.Port + "/" + cfg.DataBaseLocal.Database)
	cl, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("connect to mongo error:%s", err)
	}
	err = cl.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("ping mongo error:%s", err)
	}
	fmt.Println("Connect MongoDb Successfully")
	return cl
}

func OpenConfigFile() (Config, error) {
	absPath, _ := filepath.Abs("./config.yml")
	f, err := os.Open(absPath)
	if err != nil {
		return Config{}, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Closing file error: %v", err)
		}
	}(f)
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}

type Config struct {
	DataBaseLocal struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_local"`
}
