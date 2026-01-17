import { createClient } from "@connectrpc/connect";
import { Transport } from "@connectrpc/connect";
import { MailService } from "../gen/mail_pb";
import config from "../config/config";

export default async function performanceTest(transport: Transport) {
  const client = createClient(MailService, transport);
  const concurrencyCount = 10;
  
  console.log("\n--- PERFORMANCE TEST (Concurrent) ---");
  console.log(`Sending ${concurrencyCount} emails concurrently...\n`);

  const sendEmailTask = async (id: number) => {
    const start = performance.now();
    try {
      await client.sendMail({
        to: config.target,
        subject: `Perf Test ${id}`,
        body: "Performance testing body content.",
      });
      const end = performance.now();
      const duration = (end - start).toFixed(2);
      console.log(`Req #${id}: Success (${duration} ms)`);
      return { id, success: true, duration };
    } catch (err) {
      const end = performance.now();
      const duration = (end - start).toFixed(2);
      console.log(`Req #${id}: Failed  (${duration} ms)`);
      return { id, success: false, duration };
    }
  };

  const totalStart = performance.now();

  const promises = [];
  for (let i = 1; i <= concurrencyCount; i++) {
    promises.push(sendEmailTask(i));
  }

  await Promise.all(promises);

  const totalEnd = performance.now();
  const totalTime = (totalEnd - totalStart).toFixed(2);

  console.log("-------------------------------");
  console.log(`Total Time for ${concurrencyCount} requests: ${totalTime} ms`);
  console.log(`Average Time per request: ${(Number(totalTime) / concurrencyCount).toFixed(2)} ms`);
  console.log("-------------------------------");
}