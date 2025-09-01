package mq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/star-find-cloud/star-mall/conf"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
)

type RocketMQ struct {
	Producer     rocketmq.Producer
	PushConsumer rocketmq.PushConsumer
}

func (mq RocketMQ) initProducer() interface{} {
	var config = conf.GetConfig()
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{config.MQ.RocketMQ.NameServer}),
		producer.WithRetry(3),
		producer.WithGroupName(config.MQ.RocketMQ.ProducerGroup))
	if err != nil {
		log.AppLogger.Errorf("init producer error: %s", err)
	}

	err = p.Start()
	if err != nil {
		log.AppLogger.Errorf("start producer error: %s", err)
	}
	// 生产环境中, 连接常态化, 避免创建连接的开销
	//defer p.Shutdown()

	//msg := &primitive.Message{
	//	Topic: topic,
	//	Body:  []byte("Hello World!"),
	//}
	//
	//ctx := context.Background()
	//res, err := p.SendSync(ctx, msg)
	//if err != nil {
	//	log.AppLogger.Errorf("send message error: %s", err)
	//} else {
	//	fmt.Printf("send message success: result=%s\n", res.String())
	//}

	return p
}

// 消费者：连接集群消费消息
func (mq RocketMQ) initConsumer() interface{} {
	var config = conf.GetConfig()
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{config.MQ.RocketMQ.NameServer}),
		consumer.WithConsumerModel(consumer.Clustering), // 集群消费模式
		consumer.WithGroupName(config.MQ.RocketMQ.ConsumerGroup),
	)
	if err != nil {
		log.AppLogger.Errorf("创建消费者失败: %v\n", err)
		os.Exit(1)
	}

	// 订阅主题并注册回调函数
	err = c.Subscribe(config.MQ.RocketMQ.Topic, consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				log.AppLogger.Infof("收到消息: Topic=%s, Body=%s\n", msg.Topic, string(msg.Body))
			}
			return consumer.ConsumeSuccess, nil
		})
	if err != nil {
		log.AppLogger.Errorf("订阅主题失败: %v\n", err)
		os.Exit(1)
	}

	if err := c.Start(); err != nil {
		log.AppLogger.Errorf("启动消费者失败: %v\n", err)
		os.Exit(1)
	}
	// 生产环境中, 连接常态化, 避免创建连接的开销
	//defer c.Shutdown()

	// 阻塞等待退出信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	return c
}

func (mq RocketMQ) GetProducer() interface{} {
	return mq.initProducer()
}

func (mq RocketMQ) GetConsumer() interface{} {
	return mq.initConsumer()
}

func (mq RocketMQ) ShutDown() error {
	err := mq.Producer.Shutdown()
	if err != nil {
		return err
	}
	err = mq.PushConsumer.Shutdown()
	if err != nil {
		return err
	}
	return nil
}
