using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System.Text;

namespace OrderReportService.Services
{
    public class ServiceRabbitMQ : IHostedService
    {
        private IModel channel = null!;
        private IConnection connection = null!;

        private void Run()
        {
            var factory = new ConnectionFactory
            {
                Uri = new Uri("amqp://rabbitmq:1jj395qu@178.128.57.228:5672")
            };

            this.connection = factory.CreateConnection();
            this.channel = this.connection.CreateModel();

            // ตัวอย่างโค้ดการรับส่ง RabbitMQ https://github.com/choudhurynirjhar/rabbitmq-demo
            
            var consumer = new EventingBasicConsumer(channel);
            consumer.Received += (sender, e) =>
            {
                var body = e.Body.ToArray();
                var message = Encoding.UTF8.GetString(body);
                Console.WriteLine(message);
            };

            channel.BasicConsume("q.akenarin.order.report.service", true, consumer);

        }

        public Task StartAsync(CancellationToken cancellationToken)
        {
            this.Run();
            return Task.CompletedTask;
        }

        public Task StopAsync(CancellationToken cancellationToken)
        {
            this.channel.Dispose();
            this.connection.Dispose();
            return Task.CompletedTask;
        }

    }
}