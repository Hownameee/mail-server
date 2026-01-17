import { createClient } from "@connectrpc/connect";
import { Transport } from "@connectrpc/connect";
import config from "../config/config";
import { OtpService } from "../gen/mail_pb";

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export default async function otpTest(transport: Transport) {
  const client = createClient(OtpService, transport);
  const email = config.target;
  let res, valid;

  console.log("\n--- OTP TEST SUITE ---");

  // --- Test 1: Happy Path ---
  console.log("Test:   1. Send & Validate");
  console.log("Expect: true");
  
  res = await client.sendCode({ email });
  valid = await client.validateCode({ email, otpCode: res.otpCode });
  
  console.log(`Output: ${valid.isValid}`);
  console.log(`Result: ${valid.isValid === true ? "✅ PASS" : "❌ FAIL"}`);
  console.log("-------------------------------");

  // --- Test 2: Double Spend ---
  console.log("Test:   2. Double Spend");
  console.log("Expect: false");
  
  // (Assuming res.otpCode is from previous test)
  valid = await client.validateCode({ email, otpCode: res.otpCode });
  
  console.log(`Output: ${valid.isValid}`);
  console.log(`Result: ${valid.isValid === false ? "✅ PASS" : "❌ FAIL"}`);
  console.log("-------------------------------");

  // --- Test 3: Wrong Code ---
  console.log("Test:   3. Wrong Input Code");
  console.log("Expect: false");

  res = await client.sendCode({ email });
  valid = await client.validateCode({ email, otpCode: "000000" });

  console.log(`Output: ${valid.isValid}`);
  console.log(`Result: ${valid.isValid === false ? "✅ PASS" : "❌ FAIL"}`);
  console.log("-------------------------------");

  // --- Test 4: Retry Correct Code ---
  console.log("Test:   4. Retry Correct Code (Persistence check)");
  console.log("Expect: true");

  // We use the 'res' from Test 3 to check if the OTP still exists
  valid = await client.validateCode({ email, otpCode: res.otpCode });

  console.log(`Output: ${valid.isValid}`);
  console.log(`Result: ${valid.isValid === true ? "✅ PASS" : "❌ FAIL"}`);
  console.log("-------------------------------");

  // --- Test 5: Expiration ---
  console.log("Test:   5. Expiration (LifeSpan)");
  console.log("Expect: false");

  res = await client.sendCode({ email });
  await sleep((Number(config.otpLifeSpan) * 1000));
  valid = await client.validateCode({ email, otpCode: res.otpCode });

  console.log(`Output: ${valid.isValid}`);
  console.log(`Result: ${valid.isValid === false ? "✅ PASS" : "❌ FAIL"}`);
  console.log("-------------------------------");
}