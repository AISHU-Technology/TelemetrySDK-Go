package public

import (
	"context"
	"fmt"
	"log"
	"os"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/config"
	"github.com/Shopify/sarama"
)

// StdoutClient 客户端结构体。
type KafkaClient struct {
	producer sarama.AsyncProducer
	stopCh   chan struct{}
	topic    string
	cfg      *config.Config
}

// Path 获取上报地址。
func (k *KafkaClient) Path() string {
	return k.topic
}

// Stop 关闭发送器。
func (k *KafkaClient) Stop(ctx context.Context) error {
	close(k.stopCh)
	if err := k.producer.Close(); err != nil {
		fmt.Printf(err.Error())
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// UploadData 批量发送可观测性数据。
func (k *KafkaClient) UploadData(ctx context.Context, data []byte) error {
	msg := &sarama.ProducerMessage{Topic: k.cfg.KafkaConfig.Topic, Value: sarama.ByteEncoder(data)}
	k.producer.Input() <- msg
	return nil
}

// NewStdoutClient 创建Exporter的Local客户端。
func NewKafkaClient(opts ...config.Option) Client {
	cfg := config.NewConfig(opts...)
	logFile, err := os.Create("sarama.log")
	if err != nil {
		fmt.Printf(err.Error())
	}
	sarama.Logger = log.New(logFile, "[Sarama]", log.Lshortfile|log.Ldate|log.Ltime|log.Lmsgprefix)

	config := sarama.NewConfig()
	config.Producer.Retry.Max = 10
	config.Producer.Idempotent = true
	config.Producer.Return.Errors = false
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Net.MaxOpenRequests = 1
	if cfg.KafkaConfig.Password != "" && cfg.KafkaConfig.User != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cfg.KafkaConfig.User
		config.Net.SASL.Password = cfg.KafkaConfig.Password
	}
	producer, err := sarama.NewAsyncProducer(cfg.KafkaConfig.Address, config)
	if err != nil {
		fmt.Printf("producer init error: " + err.Error())
		return nil
	}
	k := KafkaClient{producer: producer, stopCh: make(chan struct{}), cfg: cfg}
	go k.asyncDetect()
	return &k
}
func (k *KafkaClient) asyncDetect() {
	var flagErrorEmpty bool
	var flagSuccessesEmpty bool
	for {
		if flagErrorEmpty && flagSuccessesEmpty {
			return
		}
		select {
		case err, ok := <-k.producer.Errors():
			if !ok {
				fmt.Printf("kafka_client stop detecting errors...")
				flagErrorEmpty = true
				continue
			}
			fmt.Printf("kafka_client send event err:%v", err.Err.Error())
		case ack, ok := <-k.producer.Successes():
			if !ok {
				fmt.Printf("kafka_client stop detecting successes...")
				flagSuccessesEmpty = true
				continue
			}
			res, _ := ack.Value.Encode()
			fmt.Printf("kafka_client send complete:%v", res)
		}
	}
}
