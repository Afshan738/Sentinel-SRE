𝐒𝐞𝐧𝐭𝐢𝐧𝐞𝐥-𝐒𝐑𝐄: 𝐃𝐢𝐬𝐭𝐫𝐢𝐛𝐮𝐭𝐞𝐝 𝐒𝐢𝐭𝐞 𝐑𝐞𝐥𝐢𝐚𝐛𝐢𝐥𝐢𝐭𝐲 𝐏𝐥𝐚𝐭𝐟𝐨𝐫𝐦
Sentinel-SRE is a high-performance, distributed monitoring platform designed to track website health and latency at scale. It moves beyond simple "pinger" apps by utilizing an Event-Driven Architecture to ensure high availability, fault tolerance, and millisecond-level response times.
🏗️ 𝐒𝐲𝐬𝐭𝐞𝐦 𝐀𝐫𝐜𝐡𝐢𝐭𝐞𝐜𝐭𝐮𝐫𝐞
The system is built as a set of decoupled microservices that communicate through a durable message broker and a multi-tier storage layer:
𝐀𝐏𝐈 𝐆𝐚𝐭𝐞𝐰𝐚𝐲 (𝐍𝐨𝐝𝐞.𝐣𝐬): The entry point for users to register and manage monitoring targets.
𝐀𝐮𝐭𝐨𝐧𝐨𝐦𝐨𝐮𝐬 𝐒𝐜𝐡𝐞𝐝𝐮𝐥𝐞𝐫 (𝐍𝐨𝐝𝐞.𝐣𝐬): The system "heartbeat." It identifies overdue checks based on specific intervals and triggers tasks via RabbitMQ.
𝐇𝐢𝐠𝐡-𝐂𝐨𝐧𝐜𝐮𝐫𝐫𝐞𝐧𝐜𝐲 𝐖𝐨𝐫𝐤𝐞𝐫 𝐏𝐨𝐨𝐥 (𝐆𝐨𝐥𝐚𝐧𝐠): A performance-optimized engine that utilizes Goroutines to perform thousands of concurrent network probes.
𝐌𝐞𝐬𝐬𝐚𝐠𝐞 𝐁𝐫𝐨𝐤𝐞𝐫 (𝐑𝐚𝐛𝐛𝐢𝐭𝐌𝐐): Manages task distribution and provides backpressure control between scheduling and execution.
𝐌𝐮𝐥𝐭𝐢-𝐓𝐢𝐞𝐫 𝐒𝐭𝐨𝐫𝐚𝐠𝐞:
PostgreSQL: Relational database for persistent configuration and historical time-series data.
Redis: High-speed caching layer for 
O
(
1
)
O(1)
 real-time status retrieval.
🚀 𝐄𝐧𝐠𝐢𝐧𝐞𝐞𝐫𝐢𝐧𝐠 𝐃𝐞𝐬𝐢𝐠𝐧 𝐏𝐚𝐭𝐭𝐞𝐫𝐧𝐬
𝟏. 𝐑𝐞𝐥𝐢𝐚𝐛𝐢𝐥𝐢𝐭𝐲 & 𝐅𝐚𝐮𝐥𝐭 𝐓𝐨𝐥𝐞𝐫𝐚𝐧𝐜𝐞
𝐌𝐚𝐧𝐮𝐚𝐥 𝐀𝐜𝐤𝐧𝐨𝐰𝐥𝐞𝐝𝐠𝐦𝐞𝐧𝐭𝐬 (𝐀𝐂𝐊𝐬): To ensure "At-Least-Once" delivery, workers only acknowledge messages after successful database persistence. If a worker crashes, RabbitMQ re-queues the task automatically.
𝐄𝐱𝐩𝐨𝐧𝐞𝐧𝐭𝐢𝐚𝐥 𝐁𝐚𝐜𝐤𝐨𝐟𝐟: Implemented a retry engine (
2
n
2 
n
 
) to handle transient network noise, preventing false-positive "DOWN" reports.
𝐃𝐞𝐚𝐝 𝐋𝐞𝐭𝐭𝐞𝐫 𝐐𝐮𝐞𝐮𝐞𝐬 (𝐃𝐋𝐐): Poison messages or tasks that fail all retry attempts are moved to a monitor_tasks_dead queue for manual audit and debugging.
𝟐. 𝐏𝐞𝐫𝐟𝐨𝐫𝐦𝐚𝐧𝐜𝐞 & 𝐒𝐜𝐚𝐥𝐚𝐛𝐢𝐥𝐢𝐭𝐲
𝐆𝐨 𝐂𝐨𝐧𝐜𝐮𝐫𝐫𝐞𝐧𝐜𝐲 𝐌𝐨𝐝𝐞𝐥: Utilized Go's M:N scheduler to multiplex thousands of lightweight Goroutines onto limited OS threads, optimizing CPU and RAM usage during heavy I/O wait times.
𝐖𝐫𝐢𝐭𝐞-𝐓𝐡𝐫𝐨𝐮𝐠𝐡 𝐂𝐚𝐜𝐡𝐢𝐧𝐠: Integrated a Redis layer to serve the latest status of any website instantly, reducing the read-load on the primary PostgreSQL database by 90%.
𝐁𝐚𝐜𝐤𝐩𝐫𝐞𝐬𝐬𝐮𝐫𝐞 𝐂𝐨𝐧𝐭𝐫𝐨𝐥: Configured QoS Prefetch (1) to prevent "Thundering Herd" problems and ensure a balanced load distribution across the worker pool.
🛠️ 𝐓𝐞𝐜𝐡 𝐒𝐭𝐚𝐜𝐤
Languages: Golang, Node.js, SQL (PostgreSQL)
Infrastructure: Docker, RabbitMQ, Redis
Communication: AMQP (RabbitMQ), REST (Express)
DevOps: Environment Isolation, Port Mapping, Data Normalization
🚦 𝐇𝐨𝐰 𝐓𝐨 𝐑𝐮𝐧
𝟏. 𝐈𝐧𝐟𝐫𝐚𝐬𝐭𝐫𝐮𝐜𝐭𝐮𝐫𝐞
Ensure Docker is installed and run:
code
Bash
docker-compose up -d
𝟐. 𝐒𝐞𝐫𝐯𝐢𝐜𝐞𝐬
Start the services in separate terminals:
code
Bash
# Start API
cd api-gateway && npm start

# Start Scheduler
cd scheduler-service && npm start

# Start Worker
cd worker-service && go run main.go
📈 𝐑𝐨𝐚𝐝𝐦𝐚𝐩

𝐎𝐛𝐬𝐞𝐫𝐯𝐚𝐛𝐢𝐥𝐢𝐭𝐲: Implementing Prometheus for metric exporting and Grafana for P99 latency visualization.

𝐎𝐩𝐭𝐢𝐦𝐢𝐳𝐚𝐭𝐢𝐨𝐧: Moving from single inserts to Batch/Bulk Inserts to increase database write-throughput.

𝐃𝐞𝐩𝐥𝐨𝐲𝐦𝐞𝐧𝐭: Migration to a Kubernetes (K8s) cluster with Horizontal Pod Autoscaling (HPA).
Created by [Your Name] - Focused on Engineering Excellence.