𝐒𝐞𝐧𝐭𝐢𝐧𝐞𝐥-𝐒𝐑𝐄: 𝐃𝐢𝐬𝐭𝐫𝐢𝐛𝐮𝐭𝐞𝐝 𝐒𝐢𝐭𝐞 𝐑𝐞𝐥𝐢𝐚𝐛𝐢𝐥𝐢𝐭𝐲 𝐏𝐥𝐚𝐭𝐟𝐨𝐫𝐦
Sentinel-SRE is an industrial-grade distributed monitoring system engineered to analyze website health and network performance at scale. Moving beyond traditional monolithic architectures, this platform utilizes an asynchronous, event-driven model to ensure massive throughput, fault tolerance, and millisecond-level precision.

𝐀𝐫𝐜𝐡𝐢𝐭𝐞𝐜𝐭𝐮𝐫𝐚𝐥 𝐁𝐥𝐮𝐞𝐩𝐫𝐢𝐧𝐭
The system is architected as a suite of decoupled microservices, ensuring independent scalability and fault isolation.
𝟏. 𝐒𝐞𝐫𝐯𝐢𝐜𝐞 𝐈𝐧𝐯𝐞𝐧𝐭𝐨𝐫𝐲
𝐀𝐏𝐈 𝐆𝐚𝐭𝐞𝐰𝐚𝐲 (𝐍𝐨𝐝𝐞.𝐣𝐬): Serves as the primary ingress for monitoring configurations, persisting target metadata to the relational store.
𝐀𝐮𝐭𝐨𝐧𝐨𝐦𝐨𝐮𝐬 𝐒𝐜𝐡𝐞𝐝𝐮𝐥𝐞𝐫 (𝐍𝐨𝐝𝐞.𝐣𝐬): The system "heartbeat." It utilizes time-series SQL logic to identify overdue probes and orchestrate tasks via the message broker.
𝐇𝐢𝐠𝐡-𝐂𝐨𝐧𝐜𝐮𝐫𝐫𝐞𝐧𝐜𝐲 𝐖𝐨𝐫𝐤𝐞𝐫 𝐏𝐨𝐨𝐥 (𝐆𝐨𝐥𝐚𝐧𝐠): A performance-tuned execution engine that leverages Go’s M:N scheduler and Goroutines to manage thousands of concurrent network I/O operations.
𝐌𝐞𝐬𝐬𝐚𝐠𝐞 𝐁𝐫𝐨𝐤𝐞𝐫 (𝐑𝐚𝐛𝐛𝐢𝐭𝐌𝐐): Acts as the system’s nervous system, handling task distribution, load leveling, and backpressure management.
𝟐. 𝐃𝐚𝐭𝐚 𝐏𝐞𝐫𝐬𝐢𝐬𝐭𝐞𝐧𝐜𝐞 𝐒𝐭𝐫𝐚𝐭𝐞𝐠𝐲
𝐏𝐨𝐬𝐭𝐠𝐫𝐞𝐒𝐐𝐋: Operational "Source of Truth" storing normalized monitor configurations and historical check telemetry.
𝐑𝐞𝐝𝐢𝐬: In-memory speed layer utilizing a Write-Through Cache pattern to provide  O(1)
real-time status retrieval for the dashboard.


𝐒𝐲𝐬𝐭𝐞𝐦 𝐑𝐞𝐥𝐢𝐚𝐛𝐢𝐥𝐢𝐭𝐲 & 𝐏𝐞𝐫𝐟𝐨𝐫𝐦𝐚𝐧𝐜𝐞 𝐄𝐧𝐠𝐢𝐧𝐞𝐞𝐫𝐢𝐧𝐠
𝐅𝐚𝐮𝐥𝐭 𝐓𝐨𝐥𝐞𝐫𝐚𝐧𝐜𝐞 & 𝐑𝐞𝐬𝐢𝐥𝐢𝐞𝐧𝐜𝐞
𝐆𝐮𝐚𝐫𝐚𝐧𝐭𝐞𝐞𝐝 𝐃𝐞𝐥𝐢𝐯𝐞𝐫𝐲 (𝐌𝐚𝐧𝐮𝐚𝐥 𝐀𝐂𝐊𝐬): Implemented manual message acknowledgments to ensure "At-Least-Once" delivery. If a worker node fails mid-execution, RabbitMQ automatically re-queues the task to preserve data integrity.
𝐈𝐧𝐭𝐞𝐥𝐥𝐢𝐠𝐞𝐧𝐭 𝐅𝐚𝐢𝐥𝐮𝐫𝐞 𝐌𝐢𝐭𝐢𝐠𝐚𝐭𝐢𝐨𝐧: Developed a custom Exponential Backoff algorithm (
2
n
2 
n
 
) to handle transient network noise, significantly reducing false-positive alerts.
𝐏𝐨𝐢𝐬𝐨𝐧 𝐌𝐞𝐬𝐬𝐚𝐠𝐞 𝐇𝐚𝐧𝐝𝐥𝐢𝐧𝐠 (𝐃𝐋𝐐): Architected Dead Letter Queues to isolate "poison pills" and unprocessable tasks for manual auditing without disrupting the primary pipeline.
𝐏𝐞𝐫𝐟𝐨𝐫𝐦𝐚𝐧𝐜𝐞 𝐎𝐩𝐭𝐢𝐦𝐢𝐳𝐚𝐭𝐢𝐨𝐧
𝐈/𝐎 𝐌𝐮𝐥𝐭𝐢𝐩𝐥𝐞𝐱𝐢𝐧𝐠: Go workers utilize non-blocking I/O, allowing a single process to monitor thousands of endpoints with a negligible memory footprint compared to traditional threading models.
𝐁𝐚𝐜𝐤𝐩𝐫𝐞𝐬𝐬𝐮𝐫𝐞 𝐌𝐚𝐧𝐚𝐠𝐞𝐦𝐞𝐧𝐭: Configured QoS Prefetch limits to prevent worker saturation and ensure uniform load distribution across the cluster.
𝐃𝐚𝐭𝐚 𝐍𝐨𝐫𝐦𝐚𝐥𝐢𝐳𝐚𝐭𝐢𝐨𝐧: Separated static configuration from time-series check history to optimize SQL write throughput and minimize storage redundancy.

𝐓𝐞𝐜𝐡𝐧𝐨𝐥𝐨𝐠𝐲 𝐒𝐭𝐚𝐜𝐤
𝐋𝐚𝐧𝐠𝐮𝐚𝐠𝐞𝐬: Golang, Node.js, SQL (PostgreSQL)
𝐈𝐧𝐟𝐫𝐚𝐬𝐭𝐫𝐮𝐜𝐭𝐮𝐫𝐞: Docker, RabbitMQ, Redis, PostgreSQL
𝐏𝐫𝐨𝐭𝐨𝐜𝐨𝐥𝐬: AMQP 0-9-1, REST (HTTP/1.1)
𝐒𝐲𝐬𝐭𝐞𝐦 𝐃𝐞𝐬𝐢𝐠𝐧: Microservices, Event-Driven Architecture, Distributed Caching

 𝐄𝐱𝐞𝐜𝐮𝐭𝐢𝐨𝐧 𝐆𝐮𝐢𝐝𝐞
𝟏. 𝐁𝐨𝐨𝐭𝐬𝐭𝐫𝐚𝐩 𝐈𝐧𝐟𝐫𝐚𝐬𝐭𝐫𝐮𝐜𝐭𝐮𝐫𝐞
Deploy the containerized environment using Docker Compose:
code
Bash
docker-compose up -d
𝟐. 𝐈𝐧𝐢𝐭𝐢𝐚𝐥𝐢𝐳𝐞 𝐒𝐞𝐫𝐯𝐢𝐜𝐞𝐬
Execute the services in independent runtimes:
code
Bash
# Gateway
cd api-gateway && npm install && npm start

# Scheduler
cd scheduler-service && npm install && npm start

# Worker Pool
cd worker-service && go run main.go

𝐄𝐧𝐠𝐢𝐧𝐞𝐞𝐫𝐢𝐧𝐠 𝐑𝐨𝐚𝐝𝐦𝐚𝐩

𝐎𝐛𝐬𝐞𝐫𝐯𝐚𝐛𝐢𝐥𝐢𝐭𝐲: Integration of Prometheus exporters and Grafana dashboards for real-time P99 latency visualization.

𝐇𝐢𝐠𝐡 𝐓𝐡𝐫𝐨𝐮𝐠𝐡𝐩𝐮𝐭 𝐖𝐫𝐢𝐭𝐞𝐬: Implementing Batch-Insert Buffering to further optimize database write-cycles.

𝐂𝐥𝐨𝐮𝐝-𝐍𝐚𝐭𝐢𝐯𝐞 𝐎𝐫𝐜𝐡𝐞𝐬𝐭𝐫𝐚𝐭𝐢𝐨𝐧: Transitioning to Kubernetes (K8s) with Horizontal Pod Autoscaling (HPA) for elastic demand management.
𝐀𝐮𝐭𝐡𝐨𝐫: Afshan Qasim
𝐄𝐧𝐠𝐢𝐧𝐞𝐞𝐫𝐢𝐧𝐠 𝐚 𝐫𝐞𝐥𝐢𝐚𝐛𝐥𝐞, 𝐬𝐜𝐚𝐥𝐚𝐛𝐥𝐞, 𝐚𝐧𝐝 𝐨𝐛𝐬𝐞𝐫𝐯𝐚𝐛𝐥𝐞 𝐟𝐮𝐭𝐮𝐫𝐞.
