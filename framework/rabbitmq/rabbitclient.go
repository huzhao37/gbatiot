package rabbitmq

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
	"time"
	corex "yuniot/core"
)

var (
	myConfig = new(corex.Config)
	// 定义常量
	RabbitMq_Uri             string
	RabbitMq_User            string
	RabbitMq_Pwd             string
	RabbitMq_QueueName       string
	RabbitMq_QueueBKName     string
	RabbitMq_FailedQueueName string
	RabbitMq_Control       string
	RabbitMq_ControlBK     string
	RabbitMq_ControlFailed string
	RabbitMq_QueueCount      int
	_conn                    *amqp.Connection
	_channel                 *amqp.Channel
)

func init() {
	myConfig.InitConfig("./config/config.txt")//./config/config.txt
	var env = myConfig.Read("global", "env")
	RabbitMq_Uri = myConfig.Read(env, "rabbitmq.host")
	RabbitMq_User = myConfig.Read(env, "rabbitmq.user")
	RabbitMq_Pwd = myConfig.Read(env, "rabbitmq.pwd")
	RabbitMq_QueueName = myConfig.Read(env, "rabbitmq.queuename")
	RabbitMq_QueueBKName = myConfig.Read(env, "rabbitmq.queuebkname")
	RabbitMq_FailedQueueName = myConfig.Read(env, "rabbitmq.failedqueuename")

	RabbitMq_Control = myConfig.Read(env, "rabbitmq.control")
	RabbitMq_ControlBK = myConfig.Read(env, "rabbitmq.controlbk")
	RabbitMq_ControlFailed = myConfig.Read(env, "rabbitmq.controlfailed")
}
func failOnError(err error, msg string) {
	if err != nil {
		corex.Logger.Fatalf("[rabbitmq] %s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
//消费者1
func Read(handler func([]byte)) {
	url := "amqp://" + RabbitMq_User + ":" + RabbitMq_Pwd + "@" + RabbitMq_Uri + "/" //"amqp://guest:guest@localhost:5672/"
	var err error
	_conn, err = amqp.Dial(url)
	failOnError(err, "[rabbitmq]Failed to connect to RabbitMQ")
	defer _conn.Close()

	_channel, err = _conn.Channel()
	failOnError(err, "[rabbitmq]Failed to open a channel")
	defer _channel.Close()

	q, err := _channel.QueueDeclare(
		RabbitMq_QueueName, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	_channel.QueueBind(RabbitMq_QueueName,RabbitMq_QueueName,"amq.topic",false,nil)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	err = _channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "[rabbitmq]Failed to set QoS")

	msgs, err := _channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "[rabbitmq]Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			RabbitMq_QueueCount = q.Messages
			fmt.Printf("[rabbitmq]当前队列%s累积数据量%d个\n", RabbitMq_QueueName, RabbitMq_QueueCount)

			s := fmt.Sprintf("%x", d.Body)                    //corex.BytesConvertHexArr(d.Body)// corex.BinaryStrToHexStr(biu.BytesToBinaryString(d.Body[:]))
			fmt.Printf("[rabbitmq]Received a message: %s", s) //strings.Join(s, "")
			corex.Try(func() {
				handler(d.Body)
			}, func(e interface{}) {
				Write(s, RabbitMq_FailedQueueName) //错误队列
				corex.Logger.Printf("[rabbitmq]解析原始字节数据出错: %s！", e)
			})
			//dot_count := bytes.Count(d.Body, []byte("."))
			//t := time.Duration(dot_count)
			//time.Sleep(1 * time.Second)
			fmt.Printf("[rabbitmq]Done")
			Write(s, RabbitMq_QueueBKName) //备份
			d.Ack(false)
		}
	}()

	fmt.Printf(" [[rabbitmq]] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//消息生产者-持久化消息
func Write(body string, queuename string) {
	url := "amqp://" + RabbitMq_User + ":" + RabbitMq_Pwd + "@" + RabbitMq_Uri + "/" //"amqp://guest:guest@localhost:5672/"
	if _conn==nil{
		_conn, err:= amqp.Dial(url)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer _conn.Close()
	}
	if _channel==nil{
		_channel, err := _conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer _channel.Close()
	}
	q, err := _channel.QueueDeclare(
		queuename, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	_channel.QueueBind(queuename,queuename,"amq.topic",false,nil)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	err = _channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

//消息生产者-持久化消息2
func Write2(body string, user string,pwd string,uri string,queue string) {
	url := "amqp://" + user + ":" + pwd + "@" + uri + "/" //"amqp://guest:guest@localhost:5672/"
	if _conn==nil{
		_conn, err:= amqp.Dial(url)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer _conn.Close()
	}
	if _channel==nil{
		_channel, err := _conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer _channel.Close()
	}
	q, err := _channel.QueueDeclare(
		queue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	_channel.QueueBind(queue,queue,"amq.topic",false,nil)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	err = _channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent %s", body)
}
//消费者2
func Read2(handler func(bytes []byte,mqid int)(bool),user string,pwd string,uri string,queue string,remains int,mqid int) {
	url := "amqp://" + user + ":" + pwd + "@" + uri + "/" //"amqp://guest:guest@localhost:5672/"
	var err error
	_conn, err = amqp.Dial(url)
	failOnError(err, "[rabbitmq]Failed to connect to RabbitMQ")
	defer _conn.Close()

	_channel, err = _conn.Channel()
	failOnError(err, "[rabbitmq]Failed to open a channel")
	defer _channel.Close()

	q, err := _channel.QueueDeclare(
		queue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	_channel.QueueBind(queue,queue,"amq.topic",false,nil)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	err = _channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "[rabbitmq]Failed to set QoS")

	msgs, err := _channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "[rabbitmq]Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			RabbitMq_QueueCount = q.Messages
			fmt.Printf("[rabbitmq]当前队列%s累积数据量%d个\n", queue, remains)

			s := fmt.Sprintf("%x", d.Body)
			fmt.Printf("[rabbitmq]Received a message: %s", s) //strings.Join(s, "")
			var ra =false
			corex.Try(func() {
				ra=handler(d.Body,mqid)
			}, func(e interface{}) {
				Write(s, queue+"Failed") //错误队列
				corex.Logger.Printf("[rabbitmq]解析原始字节数据出错: %s！", e)
			})
			//dot_count := bytes.Count(d.Body, []byte("."))
			//t := time.Duration(dot_count)
			//time.Sleep(1 * time.Second)
			fmt.Printf("[rabbitmq]Done")
			if ra{
				Write(s, queue+"BK") //备份
			}else {
				Write(s, queue+"Failed") //错误队列
			}
			d.Ack(false)
			time.Sleep(100*time.Millisecond)
		}
	}()

	fmt.Printf(" [[rabbitmq]] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//消费者3
func Read3(handler func(bytes []byte,mqid int)(bool),user string,pwd string,uri string,key string,queue string,remains int,mqid int) {
	url := "amqp://" + user + ":" + pwd + "@" + uri + "/" //"amqp://guest:guest@localhost:5672/"
	var err error
	_conn, err = amqp.Dial(url)
	failOnError(err, "[rabbitmq]Failed to connect to RabbitMQ")
	defer _conn.Close()

	_channel, err = _conn.Channel()
	failOnError(err, "[rabbitmq]Failed to open a channel")
	defer _channel.Close()

	q, err := _channel.QueueDeclare(
		queue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	_channel.QueueBind(queue,key,"amq.topic",false,nil)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	err = _channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "[rabbitmq]Failed to set QoS")

	msgs, err := _channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "[rabbitmq]Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			RabbitMq_QueueCount = q.Messages
			fmt.Printf("[rabbitmq]当前队列%s累积数据量%d个\n", queue, remains)

			s := fmt.Sprintf("%x", d.Body)
			fmt.Printf("[rabbitmq]Received a message: %s", s) //strings.Join(s, "")
			var ra =false
			corex.Try(func() {
				ra=handler(d.Body,mqid)
			}, func(e interface{}) {
				Write(s, queue+"Failed") //错误队列
				corex.Logger.Printf("[rabbitmq]解析原始字节数据出错: %s！", e)
			})
			//dot_count := bytes.Count(d.Body, []byte("."))
			//t := time.Duration(dot_count)
			//time.Sleep(1 * time.Second)
			fmt.Printf("[rabbitmq]Done")
			if ra{
				Write(s, queue+"BK") //备份
			}else {
				Write(s, queue+"Failed") //错误队列
			}
			d.Ack(false)
			time.Sleep(100*time.Millisecond)
		}
	}()

	fmt.Printf(" [[rabbitmq]] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//controlresponse
func ReadControl(handler func([]byte,interface{})(bool), user string,pwd string,uri string,queue string,cmd interface{}) {
	url := "amqp://" + user + ":" + pwd + "@" + uri + "/" //"amqp://guest:guest@localhost:5672/"
	var err error
	_conn, err = amqp.Dial(url)
	failOnError(err, "[rabbitmq]Failed to connect to RabbitMQ")
	defer _conn.Close()

	_channel, err = _conn.Channel()
	failOnError(err, "[rabbitmq]Failed to open a channel")
	defer _channel.Close()

	q, err := _channel.QueueDeclare(
		queue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	_channel.QueueBind(queue,queue,"amq.topic",false,nil)
	failOnError(err, "[rabbitmq]Failed to declare a queue")
	err = _channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "[rabbitmq]Failed to set QoS")

	msgs, err := _channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "[rabbitmq]Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			RabbitMq_QueueCount = q.Messages
			fmt.Printf("[rabbitmq]当前队列%s累积数据量%d个\n", queue, RabbitMq_QueueCount)

			s := fmt.Sprintf("%x", d.Body)
			fmt.Printf("[rabbitmq]Received a message: %s", s) //strings.Join(s, "")
			var ra =false
			corex.Try(func() {
				ra=handler(d.Body,cmd)
			}, func(e interface{}) {
				corex.Logger.Fatal("队列解析错误：%s",e)
				Write(s, RabbitMq_ControlFailed) //错误队列
				//corex.Logger.Printf("[rabbitmq]解析控制数据出错: %s！", e)
			})
			fmt.Printf("[rabbitmq]Done")
			if ra{
				Write(s, queue+"BK") //备份
			}else {
				Write(s, RabbitMq_ControlFailed) //错误队列
			}
			d.Ack(false)
			time.Sleep(100*time.Millisecond)
		}
	}()

	fmt.Printf(" [[rabbitmq]] Waiting for messages. To exit press CTRL+C")
	<-forever
}
//demo
func testconsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//demo
func testpub() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent %s", body)
}
