package rmq

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	exchange          = flag.String("exchange", "ex.supachai.fanout", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType      = flag.String("exchange-type", "fanout", "Exchange type - direct|fanout|topic|x-custom")
	queue             = flag.String("queue", "q.supachai.order.service", "Ephemeral AMQP queue name")
	bindingKey        = flag.String("key", "", "AMQP binding key")
	consumerTag       = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	lifetime          = flag.Duration("lifetime", 0*time.Second, "lifetime of process before shutdown (0s=infinite)")
	verbose           = flag.Bool("verbose", true, "enable verbose output of message data")
	autoAck           = flag.Bool("auto_ack", true, "enable message auto-ack")
	ErrLog            = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)
	Log               = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmsgprefix)
	deliveryCount int = 0
)

func init() {
	flag.Parse()
}

func ConnectToRMQ() {
	c, err := NewConsumer(os.Getenv("RABBITMQ_URL"), *exchange, *exchangeType, *queue, *bindingKey, *consumerTag)
	if err != nil {
		ErrLog.Fatalf("%s", err)
	}

	SetupCloseHandler(c)

	if *lifetime > 0 {
		Log.Printf("running for %s", *lifetime)
		time.Sleep(*lifetime)
	} else {
		Log.Printf("running until Consumer is done")
		<-c.done
	}

	Log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		ErrLog.Fatalf("error during shutdown: %s", err)
	}
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func SetupCloseHandler(consumer *Consumer) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		Log.Printf("Ctrl+C pressed in Terminal")
		if err := consumer.Shutdown(); err != nil {
			ErrLog.Fatalf("error during shutdown: %s", err)
		}
		os.Exit(0)
	}()
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-consumer")
	Log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err.Error())
	}

	go func() {
		Log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	Log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %s", err)
	}

	Log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("exchange declare: %s", err.Error())
	}

	Log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue declare: %s", err.Error())
	}

	Log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("queue bind: %s", err.Error())
	}

	Log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		*autoAck,   // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue consume: %s", err)
	}

	go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer Log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	cleanup := func() {
		Log.Printf("handle: deliveries channel closed")
		done <- nil
	}

	defer cleanup()

	for d := range deliveries {
		deliveryCount++
		if *verbose {
			Log.Printf(
				"got %dB delivery: [%v] %q",
				len(d.Body),
				d.DeliveryTag,
				d.Body,
			)
		} else {
			if deliveryCount%65536 == 0 {
				Log.Printf("delivery count %d", deliveryCount)
			}
		}
		if !*autoAck {
			d.Ack(false)
		}
	}
}
