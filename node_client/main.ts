import { createGrpcTransport } from "@connectrpc/connect-node";
import config from "./config/config";
import otpTest from "./testcases/otp";
import mailTest from "./testcases/mail";
import performanceTest from "./testcases/performance";

console.log("--- App Configuration ---");
console.table({
  Port: config.port,
  "OTP Cleanup (sec)": config.otpCleanup,
  "OTP Lifespan (sec)": config.otpLifeSpan,
  Workers: config.workers,
  Target: config.target,
});

const transport = createGrpcTransport({
  baseUrl: `http://localhost:${config.port}`,
});

(async () => {
  // otp test
  await otpTest(transport);

  // send email test
  await mailTest(transport);

  // performance test
  await performanceTest(transport);
})();
