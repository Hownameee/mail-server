import { createClient } from "@connectrpc/connect";
import { Transport } from "@connectrpc/connect";
import { MailService } from "../gen/mail_pb";
import config from "../config/config";

export default async function mailTest(transport: Transport) {
  const client = createClient(MailService, transport);

  console.log("\n--- MAIL TEST SUITE ---");

  // --- Test 1: Send Valid Email ---
  console.log("Test:   1. Send Valid Email");
  console.log("Expect: Success (No Error)");

  try {
    const res = await client.sendMail({
      to: config.target,
      subject: "Test from Node Client",
      body: "This is a test email sent via ConnectRPC.",
    });

    console.log(`Output: Success (Status: ${res.status})`);
    console.log("Result: ✅ PASS");
    console.log("Check your email to get test message");
  } catch (err) {
    console.log(`Output: Error - ${err}`);
    console.log("Result: ❌ FAIL");
  }
  console.log("-------------------------------");

  // --- Test 2: Validation Error (Missing To) ---
  console.log("Test:   2. Missing Recipient");
  console.log("Expect: Error (recipient email is required)");

  try {
    await client.sendMail({
      to: "",
      subject: "Should Fail",
      body: "...",
    });

    console.log("Output: Success (Unexpected)");
    console.log("Result: ❌ FAIL");
  } catch (err: any) {
    const msg = err.message || "Unknown Error";
    console.log(`Output: Error Caught ("${msg}")`);
    console.log("Result: ✅ PASS");
  }
  console.log("-------------------------------");
}
